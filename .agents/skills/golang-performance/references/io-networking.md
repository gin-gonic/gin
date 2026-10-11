# I/O & Networking Optimization

Network and I/O bottlenecks show up as goroutines blocked on syscalls or waiting for responses. The key levers are connection reuse, proper timeouts, and streaming instead of buffering.

## HTTP Transport Configuration

**Diagnose:** 1- `go tool pprof` (goroutine + block profile) — look for goroutines blocked on `net/http.(*Transport).dialConn` or `net/http.(*persistConn).readLoop`; many goroutines waiting here means connection pool exhaustion 2- `fgprof` — captures both on-CPU and off-CPU wait time; look for HTTP calls dominating wall-clock time even when CPU profile shows them as cheap 3- `go tool trace` — visualize goroutine lifecycles; look for long gaps where goroutines wait for network I/O instead of processing 4- Prometheus `go_goroutines` — monitor goroutine count in production; steadily rising under stable load suggests connection or goroutine leaks from misconfigured HTTP clients

### Connection pooling

The default `http.Transport` has conservative pool settings — `MaxIdleConnsPerHost` defaults to 2. Under high concurrency, requests queue waiting for connections instead of running in parallel:

```go
// Bad — default transport, only 2 idle connections per host
client := &http.Client{}

// Good — tuned for high-concurrency service-to-service calls
var apiClient = &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:          100,             // total idle connections across all hosts
        MaxIdleConnsPerHost:   20,              // per-host idle connections (default is 2!)
        MaxConnsPerHost:       50,              // cap total connections per host (0 = unlimited)
        IdleConnTimeout:       90 * time.Second,
        TLSHandshakeTimeout:  5 * time.Second,
        ResponseHeaderTimeout: 10 * time.Second,
    },
}
```

For web crawlers hitting many different hosts, disable keep-alive to avoid accumulating idle connections:

```go
crawlerClient := &http.Client{
    Transport: &http.Transport{DisableKeepAlives: true},
}
```

### Timeouts

The zero-value `http.Client` and `http.Server` have NO timeouts. A slow or malicious peer holds connections open indefinitely, exhausting file descriptors and memory:

```go
// Server — always set timeouts to prevent Slowloris attacks
server := &http.Server{
    Addr:         ":8080",
    Handler:      handler,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
```

### Drain response body for connection reuse

Connections are only returned to the pool when the body is fully read. Even if you don't need the body, drain it:

```go
resp, err := client.Get(url)
if err != nil { return err }
defer resp.Body.Close()
_, _ = io.Copy(io.Discard, resp.Body) // drain to enable connection reuse
```

## Streaming vs Buffering

**Diagnose:** 1- `go tool pprof -inuse_space` — look for large single allocations (MB-sized) from `io.ReadAll`, `bytes.Buffer.Grow`, or `json.Unmarshal`; these indicate buffering entire payloads instead of streaming

### Avoid io.ReadAll for large payloads

`io.ReadAll` loads the entire stream into memory. For large files or HTTP responses, this causes massive memory spikes:

```go
// Bad — 2GB file = 2GB allocation
data, _ := io.ReadAll(f)

// Good — process line by line, O(1) memory
scanner := bufio.NewScanner(f)
for scanner.Scan() { processLine(scanner.Bytes()) }

// Good — stream between reader and writer (32KB internal buffer)
io.Copy(w, resp.Body)
```

`io.ReadAll` is fine for small, bounded payloads (< 1MB) where the size is known.

### Streaming JSON

Use `json.NewDecoder` for large JSON payloads instead of `json.Unmarshal` (which buffers the entire body):

```go
dec := json.NewDecoder(r)
for dec.More() {
    var item Item
    if err := dec.Decode(&item); err != nil { return err }
    process(item) // one item at a time
}
```

## JSON Performance

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for `encoding/json.(*Decoder).Decode`, `reflect.Value.*`, or `encoding/json.Marshal` consuming significant CPU; these indicate reflection-based JSON is the bottleneck 2- `go test -bench -benchmem` — measure ns/op and allocs/op for marshal/unmarshal; expect high alloc counts from reflection; code-gen alternatives should show 2-5x fewer allocs

The standard `encoding/json` package uses reflection to inspect struct fields at runtime. For high-throughput services, this creates significant CPU and allocation overhead.

**Options for faster JSON:**

- **Custom `MarshalJSON`/`UnmarshalJSON`** — hand-written methods for hot-path types eliminate reflection
- **Code-generation libraries** — `easyjson`, `ffjson` generate marshal/unmarshal methods at build time, no reflection at runtime
- **Drop-in replacements** — `github.com/goccy/go-json`, `github.com/json-iterator/go`, `github.com/bytedance/sonic` offer 2-5x better performance
- **`encoding/json/v2`** (experimental, behind `GOEXPERIMENT=jsonv2`) — evaluate deliberately; most production code should keep `encoding/json` unless the project explicitly opts into the experiment

When using third-party JSON libraries, refer to the library's official documentation for up-to-date API signatures.

## Cgo Overhead

**Diagnose:** 1- `go tool pprof` (CPU profile + threadcreate profile) — look for `runtime.cgocall` or `runtime.asmcgocall` consuming CPU; high threadcreate count means cgo calls are pinning goroutines to OS threads 2- `go test -bench` — benchmark the cgo call loop vs a pure Go equivalent; expect ~50-100ns overhead per cgo crossing

Each Go-to-C call via cgo costs ~50-100ns due to stack switching, signal mask manipulation, and scheduler coordination:

```go
// Bad — cgo overhead per element dominates for tight loops
for i, v := range values {
    values[i] = float64(C.sqrt(C.double(v))) // ~100ns overhead PER CALL
}

// Good — use pure Go stdlib (math.Sqrt is as fast as C and inlineable)
for i, v := range values { values[i] = math.Sqrt(v) }

// Good — batch when C code is unavoidable
C.batch_sqrt((*C.double)(&values[0]), C.int(len(values))) // amortize overhead
```

Additional cgo costs: goroutine is pinned to an OS thread, C code cannot be preempted (may delay GC), and function inlining is blocked at the boundary.

## Buffered I/O

**Diagnose:** 1- `go test -bench` — benchmark buffered vs unbuffered I/O; expect 3-10x improvement from reducing syscall count 2- `go tool trace` — look for frequent short syscalls (`pread`, `pwrite`) in rapid succession; many tiny I/O operations indicate unbuffered access

Unbuffered file reads/writes issue a syscall per operation. `bufio.Reader` and `bufio.Writer` batch small operations, reducing syscalls by 10x or more:

```go
// Bad — syscall per line
for _, line := range lines { f.WriteString(line + "\n") }

// Good — buffered, batches writes into larger chunks
w := bufio.NewWriter(f)
for _, line := range lines { w.WriteString(line + "\n") }
w.Flush()
```

## Concurrent Multi-Stage Pipelines

**Diagnose:** 1- `go tool trace` — visualize resource utilization across stages; look for sequential idle gaps where CPU, disk, or network sit unused while another resource is busy 2- `go tool pprof` (CPU + goroutine profile) — confirm each stage saturates a _different_ resource; if multiple stages compete for the same resource (e.g., both CPU-bound), concurrency won't help

In rare scenarios where each pipeline stage saturates a _different_ resource (CPU, disk I/O, network), running stages concurrently instead of sequentially can improve throughput — even with batching between stages.

### The unusual scenario

Imagine processing records: Stage A compresses (CPU-bound), Stage B writes to disk (I/O-bound), Stage C uploads to network (network-bound). Sequential execution wastes resources:

```
Time:    0       10      20      30      40      50
CPU:     AAAAAAAAAA|..........|..........|..........|
Disk:    ..........|BBBBBBBBBB|..........|..........|
Network: ..........|..........|CCCCCCCCCC|..........|
```

Concurrent stages let resources work in parallel:

```
Time:    0       10      20      30      40      50
CPU:     AAAAAAAAAA|AA........|
Disk:    ..........|BBBBBBBBBB|BB........|
Network: ..........|..........|CCCCCCCCCC|CC........|
```

**Code pattern:**

```go
// Each stage runs in its own goroutine, bounded by channel buffers
compressedCh := make(chan []byte, 100)    // A → B buffer
uploadedCh := make(chan bool, 100)        // B → C buffer

// Stage A: CPU-bound compression
go func() {
    for record := range inputCh {
        compressed := compress(record)    // saturates CPU
        compressedCh <- compressed
    }
    close(compressedCh)
}()

// Stage B: I/O-bound disk writes
go func() {
    for compressed := range compressedCh {
        diskFile.Write(compressed)        // saturates disk I/O
        uploadedCh <- true
    }
    close(uploadedCh)
}()

// Stage C: network-bound uploads
go func() {
    for <-uploadedCh {
        client.Post(uploadURL, ...)       // saturates network
    }
}()
```

With batching per stage, total throughput = min(A_throughput, B_throughput, C_throughput). Without concurrency, throughput = sequential sum of stages. **Concurrent stages only help when bottlenecks don't overlap.**

### When to use this (and when NOT to)

**Use concurrent pipelines only when ALL of these are true:**

1. **Resource saturation is predictable and non-overlapping** — You measured that A saturates one resource (e.g., CPU = 95%), B saturates another (disk I/O = 90%), C saturates a third (network = 85%). Overlapping saturation means concurrency adds no benefit.
2. **Bottleneck shifts don't hurt latency** — Processing order doesn't matter, or records can flow out-of-order through stages.
3. **Buffering overhead is acceptable** — Inter-stage channels consume memory. For large records, channel buffers can overflow system limits.
4. **You've benchmarked the alternative** — Profile both sequential and concurrent versions. Sequential + batching often wins because it is simpler and avoids context-switching overhead.

**Avoid concurrent pipelines if:**

- **Records must be ordered** — Concurrent processing may reorder records; if downstream expects order, you need synchronization that kills the speedup.
- **Resources overlap** — If A and B both compete for CPU (e.g., both compress), concurrency causes context-switching overhead with no resource utilization gain.
- **Latency matters more than throughput** — A single record now travels through 3 stages in parallel, increasing per-record latency.
- **Memory is tight** — Each stage's channel buffer is a memory budget; deeply buffered channels can exhaust available RAM.

→ See `samber/cc-skills-golang@golang-concurrency` skill for detailed channel patterns and when to use worker pools instead.

## Batch Operations

**Diagnose:** 1- `go test -bench` — benchmark single-item vs batched operations; expect N-fold improvement in throughput when amortizing per-operation overhead (syscalls, round-trips) 2- `go tool trace` — look for repeated short network/disk operations with idle gaps between them; these gaps represent wasted round-trip time that batching eliminates

Batching amortizes per-operation overhead (syscalls, network round-trips, transaction costs) across many items. The pattern applies everywhere: I/O, database, network, and even in-memory processing.

### Database: batch inserts over row-by-row

Inserting 1,000 rows one at a time means 1,000 round-trips, 1,000 query parses, and 1,000 transaction commits. A single batch insert does it in one round-trip:

```go
// Bad — 1,000 round-trips, ~500ms
for _, user := range users {
    db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
}

// Good — 1 round-trip with multi-row VALUES, ~5ms
const batchSize = 1000
for i := 0; i < len(users); i += batchSize {
    end := min(i+batchSize, len(users))
    batch := users[i:end]
    // Build multi-row INSERT or use COPY protocol
    tx, _ := db.Begin()
    stmt, _ := tx.Prepare(pq.CopyIn("users", "name", "email"))
    for _, u := range batch { stmt.Exec(u.Name, u.Email) }
    stmt.Exec()
    tx.Commit()
}
```

→ See `samber/cc-skills-golang@golang-database` skill for detailed batch patterns and connection pool configuration.

### HTTP: batch API calls

Instead of N individual HTTP requests, send one request with N items when the API supports it:

```go
// Bad — 100 HTTP round-trips
for _, id := range ids {
    resp, _ := client.Get(fmt.Sprintf("/api/users/%s", id))
    // ...
}

// Good — 1 HTTP request with all IDs
resp, _ := client.Post("/api/users/batch", "application/json",
    bytes.NewReader(marshalIDs(ids)))
```

### Channel: batch processing from a stream

Accumulate items from a channel and process in bulk to reduce per-item overhead:

```go
func batchProcessor(in <-chan Item, batchSize int) {
    batch := make([]Item, 0, batchSize)
    ticker := time.NewTicker(100 * time.Millisecond) // flush on timeout too
    defer ticker.Stop()
    for {
        select {
        case item, ok := <-in:
            if !ok { flush(batch); return }
            batch = append(batch, item)
            if len(batch) >= batchSize { flush(batch); batch = batch[:0] }
        case <-ticker.C:
            if len(batch) > 0 { flush(batch); batch = batch[:0] }
        }
    }
}
```
