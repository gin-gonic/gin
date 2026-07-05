# CPU Optimization

CPU-bound bottlenecks show up as functions dominating the CPU profile. The patterns below target the most common causes: missed inlining opportunities, poor cache utilization, and unnecessary computation.

## Function Inlining

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for hot functions with high cumulative CPU time; if a small helper dominates the profile, it's likely not being inlined 2- `go build -gcflags="-m"` — grep for `"cannot inline"` on your hot-path functions; the reason (e.g., `"function too complex"`, `"unhandled op"`) tells you what to simplify

The Go compiler inlines small functions, eliminating call overhead. Functions that are too complex (loops, many statements, or calls to non-inlineable functions) won't be inlined — this matters in tight loops called millions of times.

```go
// Bad — log call prevents inlining
func abs(x int) int {
    if x < 0 {
        log.Printf("negative: %d", x) // blocks inlining
        return -x
    }
    return x
}

// Good — simple enough to inline
func abs(x int) int {
    if x < 0 { return -x }
    return x
}
```

**Check inlining decisions:**

```bash
go build -gcflags="-m" ./... 2>&1 | grep "can inline"
go build -gcflags="-m" ./... 2>&1 | grep "inlining call"
```

Move side effects (logging, metrics) outside hot-path functions or guard them with conditional checks.

### Value receivers enable inlining

Value receivers allow the compiler to fully inline fluent method chains. Pointer receivers add indirection that blocks inlining:

```go
// Pointer receiver — indirection prevents inlining, constant overhead per call
func (c *config) WithTimeout(d time.Duration) *config { c.timeout = d; return c }

// Value receiver — fully inlined, -80% time in fluent chains
func (c config) WithTimeout(d time.Duration) config { c.timeout = d; return c }
```

## Cache Locality

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for loops over slices/matrices consuming disproportionate CPU; cache-miss-heavy code shows high `runtime.memmove` or flat time in simple index operations 2- `go test -bench` — benchmark row-first vs column-first traversal; expect 10-50x difference on large matrices purely from cache effects

Modern CPUs fetch data in 64-byte cache lines. Sequential memory access is dramatically faster than random access because the prefetcher can load the next cache line before you need it.

### Row-major traversal

Go stores 2D arrays in row-major order. Column-first traversal jumps across memory, causing cache misses:

```go
// Bad — column-first, jumps across memory (~10M cache misses)
for col := 0; col < 1024; col++ {
    for row := 0; row < 1024; row++ {
        sum += matrix[row][col]
    }
}

// Good — row-first, sequential access (~125K cache misses)
for row := 0; row < 1024; row++ {
    for col := 0; col < 1024; col++ {
        sum += matrix[row][col]
    }
}
```

Performance difference: 10-50x purely from cache effects.

### Contiguous 2D allocation

Allocating each row separately scatters data across the heap:

```go
// Bad — N separate allocations, poor cache locality
matrix := make([][]float64, rows)
for i := range matrix { matrix[i] = make([]float64, cols) }

// Good — single contiguous allocation, cache-friendly
data := make([]float64, rows*cols)
matrix := make([][]float64, rows)
for i := range matrix { matrix[i] = data[i*cols : (i+1)*cols] }
```

### Struct of Arrays (SoA) vs Array of Structs (AoS)

When iterating over a single field of a struct, AoS wastes cache space loading unused fields:

```go
// AoS — loading each Point (24 bytes) to read only x (8 bytes) = 66% cache waste
type Point struct { x, y, z float64 }
points := make([]Point, n)
for i := range points { sum += points[i].x }

// SoA — all x values contiguous, 100% cache utilization
type Points struct { xs, ys, zs []float64 }
for i := range ps.xs { sum += ps.xs[i] }
```

Use SoA when iterating over a subset of fields (physics, graphics, analytics). AoS is fine when accessing all fields together or for small structs.

### Pointer-heavy vs value-heavy data

Index-based data structures (nodes stored in a contiguous array, referenced by index) beat pointer-based structures for cache locality:

```go
// Pointer-based tree — each node scattered in heap, random cache misses
type Node struct { value int; left, right *Node }

// Index-based tree — nodes in contiguous array, cache-friendly
type Tree struct { nodes []Node }
type Node struct { value int; left, right int } // indices into nodes
```

## False Sharing

**Diagnose:** 1- `go tool pprof` (CPU profile + mutex profile) — look for atomic operations or counter updates consuming unexpectedly high CPU; in the mutex profile, look for contention on variables that shouldn't need locking 2- `go test -bench` — benchmark concurrent counter increments; if adding goroutines makes it _slower_ instead of faster, false sharing is likely

When goroutines update variables that share the same 64-byte CPU cache line, each write invalidates the other core's cache, causing severe degradation:

```go
// Bad — a and b on same cache line, cores fight for it
type Counters struct { a, b int64 }

// Good — separate cache lines, no interference
type Counters struct {
    a int64    // 8 bytes
    _ [56]byte // 64 - 8 = 56 bytes padding
    b int64    // 8 bytes
}
```

Only apply cache-line padding when profiling confirms contention on concurrent counters/flags.

## Instruction-Level Parallelism

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for tight arithmetic loops (sum, dot product) where the loop body itself dominates CPU; these are candidates for multi-accumulator optimization 2- `go test -bench` — benchmark single vs multi-accumulator versions; expect 2-4x improvement when the loop is truly CPU-bound with a dependency chain

Modern CPUs execute multiple independent instructions simultaneously. A single accumulator creates a dependency chain — each addition waits for the previous one:

```go
// Bad — sequential dependency, CPU pipeline stalls
var total int64
for _, v := range data { total += v }

// Good — 4 independent accumulators, CPU pipelines all 4 in parallel
var s0, s1, s2, s3 int64
limit := len(data) - len(data)%4
for i := 0; i < limit; i += 4 {
    s0 += data[i]; s1 += data[i+1]; s2 += data[i+2]; s3 += data[i+3]
}
for i := limit; i < len(data); i++ { s0 += data[i] }
total := s0 + s1 + s2 + s3
```

Expect 2-4x improvement for tight arithmetic loops. Only use when profiling shows the loop is a bottleneck.

## SIMD (Single Instruction, Multiple Data)

**Diagnose:** 1- `go tool pprof` (CPU profile) — confirm a numeric inner loop consumes >20% of CPU; SIMD only helps CPU-bound numeric work, not allocation or I/O bottlenecks 2- `go test -bench` — measure the loop's baseline ns/op; provides the reference point to validate SIMD gains 3- `go build -gcflags="-d=ssa/prove/debug=2"` — check if the compiler already auto-vectorized the loop; look for `"Proved"` bounds-check eliminations that enable vectorization 4- `GOSSAFUNC=MyFunc go build` — generate SSA dump (`ssa.html`) to inspect whether the compiler produces vector instructions for the hot loop 5- `go tool objdump -s MyFunc ./binary` — verify the final assembly contains SIMD instructions (e.g., `VMOVAPD`, `VADDPD` on amd64) rather than scalar equivalents

Go 1.26+ includes an experimental `simd/archsimd` package (requires `GOEXPERIMENT=simd` flag) providing low-level SIMD intrinsics for amd64 with 128/256/512-bit vectors. For broader portability, the compiler auto-vectorizes simple loops, and several strategies exist.

**Options for explicit SIMD in Go:**

- **Experimental `simd/archsimd` (Go 1.26+, speculative)** — Direct SIMD intrinsics via vector types with CPU feature detection. Limited to AMD64. Use with caution: this is an experimental, in-progress API (`GOEXPERIMENT=simd`) whose package path and type names are subject to change before stabilization. Not covered by Go 1 compatibility guarantees, and should never be exposed in public APIs. Verify the actual import path and API against the Go toolchain you are using.

  ```go
  // Requires: GOEXPERIMENT=simd go build
  // WARNING: experimental API — package path and types may change
  import "simd/archsimd"

  v := archsimd.Int32x4{1, 2, 3, 4}
  ```

- **Let the compiler do it** — write simple, idiomatic loops on `[]float64`/`[]int32` slices. Check auto-vectorization: `go build -gcflags="-d=ssa/prove/debug=2" ./...`
- **`math/bits`** — operations like `OnesCount`, `LeadingZeros`, `RotateLeft` map directly to hardware instructions (POPCNT, CLZ, ROL)
- **Hand-written assembly** — `.s` files with AVX2/NEON instructions for critical inner loops. Libraries like `klauspost/compress` and `minio/sha256-simd` use this approach
- **Third-party vectorized libraries** — for common operations (hashing, compression, encoding), use libraries that already have optimized SIMD implementations rather than writing your own

### Handling CPU-specific instruction sets

Hand-written assembly unlocks higher performance but couples code to specific CPU features (AVX2, NEON, etc.). Three strategies exist:

**1. Compile on a production-similar machine**

Build binaries on hardware matching your deployment target, so the compiler generates code for the exact CPU instruction set available at runtime:

```bash
# Compiling on production hardware ensures optimal code generation
# for that specific CPU architecture and generation
ssh prod-server "cd /path && go build -o app ."
```

**Tradeoff:** Simplest approach, but requires access to production hardware and different binaries per CPU type (Intel vs AMD vs Apple Silicon). Breaks CI/CD portability.

**2. Runtime CPU feature detection + multiple implementations**

Implement the function multiple times — one for each CPU capability — and dispatch at runtime:

```go
// dispatch.go
var sumImpl func([]int64) int64

func init() {
    if cpu.X86.HasAVX2 {
        sumImpl = sumAVX2
    } else {
        sumImpl = sumGeneric
    }
}

func Sum(data []int64) int64 {
    return sumImpl(data)
}

// sum_generic.go
func sumGeneric(data []int64) int64 {
    var total int64
    for _, v := range data { total += v }
    return total
}

// sum_amd64.s
TEXT ·sumAVX2(SB), NOSPLIT, $0-32
    // AVX2 implementation
    VMOVAPD (SI), Y0
    // ...
```

**Tradeoff:** Single binary works everywhere; trades one function-call dispatch overhead for full CPU feature utilization. Libraries like `encoding/base64` and `sha256` use this pattern.

**3. Compile-time selection with `//go:build` tags**

Use conditional compilation to generate different code at build time for each target:

```go
// sum_fast.go
//go:build amd64 && !nosimd

package mylib

// AVX2 assembly via cgo or inline
func Sum(data []int64) int64 {
    return sumAVX2(data) // or calls to .s file
}

// sum_generic.go
//go:build !amd64 || nosimd

package mylib

func Sum(data []int64) int64 {
    var total int64
    for _, v := range data { total += v }
    return total
}
```

Build different binaries per target:

```bash
GOOS=linux GOARCH=amd64 go build -o app-avx2 .     # Uses sum_fast.go
GOOS=darwin GOARCH=arm64 go build -o app-neon .    # Uses sum_generic.go
go build -tags=nosimd -o app-safe .                # Fallback everywhere
```

**Tradeoff:** Zero runtime overhead; each binary is fully optimized for its target. Requires shipping multiple binaries and coordinating which binary runs where.

**When SIMD is NOT worth pursuing:**

- Go's lack of intrinsics means SIMD requires assembly — high maintenance burden, platform-specific, and harder to debug
- Auto-vectorization covers the most common cases (simple numeric loops)
- If your bottleneck is allocations or I/O, SIMD won't help

**Recommendation:** Start with auto-vectorization. For Go 1.26+, evaluate `simd/archsimd` for AMD64-only workloads (remembering it's experimental). Move to runtime detection (option 2 above) if profiling shows a bottleneck and the code needs to run on heterogeneous hardware. Only use compile-time selection (option 3) if you control the deployment environment and can test each per-binary variant.

Only invest in hand-written SIMD when profiling shows a numeric inner loop consuming >20% of CPU and the compiler isn't auto-vectorizing it.

## Tight Loops and the Scheduler

**Diagnose:** 1- `go tool pprof` (goroutine profile) — look for many goroutines stuck in `"runnable"` state (waiting for CPU) while one goroutine monopolizes execution 2- `go tool trace` — visualize goroutine scheduling over time; look for long uninterrupted execution spans on one goroutine while others show scheduling gaps 3- `GODEBUG=schedtrace=1000` — print scheduler state every second; look for unbalanced `runqueue` counts across P's indicating one P is starved 4- `runtime/metrics` (`/sched/latencies:seconds`) — measure how long goroutines wait before getting CPU; high p99 latencies confirm starvation 5- Prometheus `rate(process_cpu_seconds_total[2m])` — monitor if CPU usage hits GOMAXPROCS ceiling; if saturated while other goroutines are starved, a tight loop is monopolizing P's

A goroutine running a CPU-intensive tight loop without function calls may not yield to the scheduler, starving other goroutines. Go 1.14+ added asynchronous preemption, but very tight loops with fully inlined operations can still cause issues:

```go
// Potential starvation — pure computation, no function calls
for { x = x*a + b }

// Safe — non-inlined call triggers preemption check
for item := range work {
    processBatch(item) // function call = preemption point
}
```

**When to use non-inlined calls for scheduling:** Use non-inlined function calls when:

- The loop runs for a long time (hundreds of milliseconds or more of uninterrupted computation)
- Other goroutines are waiting to run (e.g., handling requests, I/O completion, channel operations)
- The loop contains only arithmetic or memory operations with no function calls

For short bursts of computation (< 10ms), preemption isn't critical and inlining for CPU efficiency takes priority.

**Detecting scheduler starvation:** Use these tools to confirm goroutines are being starved:

- **`go tool pprof` goroutine profile** — shows goroutines stuck in "runnable" state (waiting for CPU). If many goroutines are runnable while one dominates CPU, starvation is happening
- **`go tool trace`** — visualizes goroutine scheduling over time. Look for gaps where goroutines aren't running because one goroutine monopolized the scheduler
- **`runtime/metrics` (Go 1.19+)** — measure `/sched/latencies:seconds` to quantify how long goroutines wait for CPU
- **Observable symptoms** — high response latency, requests timing out, uneven request distribution, goroutine counts climbing

**Preventing inlining with `//go:noinline`:** If you have a function that's normally inlinable (small, hot) but you specifically want it to not inline to force scheduler preemption checks, use the `//go:noinline` compiler directive:

```go
//go:noinline
func processBatch(item WorkItem) {
    // CPU-intensive work here
    // This call site will NOT be inlined, even if the function is small
    // The function call itself becomes a preemption point for the scheduler
}

// In tight loop
for item := range work {
    processBatch(item) // Guaranteed preemption point
}
```

**Trade-off:** Using `//go:noinline` prevents inlining, which:

- **Pros:** Guarantees scheduler preemption checks; prevents goroutine starvation
- **Cons:** Adds function call overhead (~10-30 CPU cycles); reduces instruction-level parallelism (ILP) in the caller

Only use `//go:noinline` if profiling shows that scheduler preemption starvation is actually blocking other goroutines. Unnecessary `//go:noinline` directives penalize throughput and latency.

## Reflection and Type Assertions

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for `reflect.Value.*`, `reflect.DeepEqual`, or `fmt.Sprintf` (which uses reflect internally) appearing in hot paths 2- `go test -bench` — compare reflection-based vs typed versions; expect 10-200x difference depending on the reflection operation

- **`reflect` in hot paths** — 10-100x slower due to type introspection and boxing. Replace with generics or typed code
- **`reflect.DeepEqual`** — 50-200x slower than typed comparisons. Use `slices.Equal`, `maps.Equal`, `bytes.Equal` (Go 1.21+)
- **Type switch vs repeated assertions** — type switch dispatches in one evaluation:

```go
// Bad — evaluates interface multiple times
if s, ok := v.(string); ok { return s }
if i, ok := v.(int); ok { return strconv.Itoa(i) }

// Good — single dispatch
switch v := v.(type) {
case string: return v
case int:    return strconv.Itoa(v)
}
```

## Monotonic Time

**Diagnose:** 1- `go test -bench` — benchmark `time.Since(start)` vs `time.Now().Sub(start)`; expect a small but consistent improvement from monotonic clock avoiding wall-clock syscall

`time.Since(start)` uses the monotonic clock, which is immune to wall-clock adjustments (NTP, DST) and slightly faster:

```go
var appStart = time.Now() // captures monotonic time + wall-clock on program start

func myFunc() {
    // Compare durations, not wall-clock times
    elapsed := time.Since(appStart)
    if elapsed > threshold { ... }
}
```
