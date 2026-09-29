# Runtime Tuning

Runtime settings control garbage collection frequency, memory limits, CPU scheduling, and compiler optimizations. Tune them after profiling — the defaults are well-chosen for most workloads.

## Garbage Collector Tuning

**Diagnose:** 1- `GODEBUG=gctrace=1` — print one line per GC cycle; look for high GC frequency (cycles/s), high CPU% (>5% means GC is competing for CPU), or heap growing faster than expected 2- `runtime.ReadMemStats` — inspect `Alloc`, `TotalAlloc`, `NumGC`, `PauseNs`; compare `Alloc` vs `Sys` to see how much memory the GC is reclaiming vs how much the OS allocated 3- `go tool trace` — visualize GC stop-the-world pauses and GC assist stealing CPU from application goroutines; look for long STW bars or frequent assist marks 4- `debug.ReadGCStats` — get pause time percentiles (p50, p95, p99); high p99 pauses indicate large heap scans or too many pointers 5- `runtime/metrics` — programmatic access to GC stats for dashboards; monitor `/gc/cycles/total`, `/gc/heap/allocs`, `/gc/pauses` 6- `GODEBUG=gcpacertrace=1` — trace the GC pacer's decisions; useful to understand why GC triggers earlier or later than expected 7- Prometheus `rate(go_gc_duration_seconds_count[5m])` — monitor GC frequency in production; >2 cycles/s sustained suggests excessive allocation rate

### GOGC (default: 100)

Controls the heap growth ratio that triggers the next GC cycle. `GOGC=100` means GC runs when the heap doubles since the last collection. Higher values reduce GC frequency but use more memory:

```bash
GOGC=50  ./myapp  # latency-sensitive: more frequent, shorter GC pauses
GOGC=200 ./myapp  # throughput-oriented: less frequent GC, more memory used
GOGC=off ./myapp  # disable GC entirely (testing only!)
```

### GOMEMLIMIT (Go 1.19+)

Soft memory limit — the runtime increases GC frequency to stay under this limit. Essential for containerized applications where exceeding the container limit triggers an OOM kill:

```bash
# Container with 512MB limit: leave headroom for non-heap memory (goroutine stacks, OS buffers)
GOMEMLIMIT=450MiB ./myapp

# Container with 1GB limit
GOMEMLIMIT=900MiB ./myapp
```

The GC pacer adjusts collection timing based on both GOGC and GOMEMLIMIT. When the heap approaches the limit, the GC runs more aggressively regardless of GOGC.

### Programmatic control

```go
import "runtime/debug"

debug.SetGCPercent(200)                    // equivalent to GOGC=200
debug.SetMemoryLimit(450 * 1024 * 1024)   // 450 MiB soft limit
```

Use programmatic control for dynamic tuning based on observed workload, or when environment variables cannot be set.

### Ballast pattern (pre-Go 1.19)

Before GOMEMLIMIT, teams allocated a large byte array at startup to inflate the live heap size, reducing GC frequency:

```go
var ballast [1 << 30]byte // 1 GB — obsolete pattern
```

**GOMEMLIMIT is strictly better** — it provides the same benefit (fewer GC cycles) without wasting physical memory. Use GOMEMLIMIT instead.

## GC Profiling and Diagnostics

### GODEBUG=gctrace=1

Prints a line per GC cycle to stderr:

```bash
GODEBUG=gctrace=1 ./myapp 2>&1 | head -20
```

Sample output:

```
gc 5 @1.234s 2%: 0.012+12+0.9 ms clock, 0.25+8.9/20+18 ms cpu, 45->92->50 MB, 200 MB goal, 8 P
```

Key fields:

- `gc 5` — 5th GC cycle
- `@1.234s` — time since program start
- `2%` — total CPU time spent in GC
- `45->92->50 MB` — heap before → peak during collection → after
- `200 MB goal` — target heap size (based on GOGC and GOMEMLIMIT)
- `8 P` — number of processors

Watch for: GC frequency (too often = too many allocations), pause times (high = large heap or many pointers), CPU% (high = tune GOGC or reduce allocations).

### runtime.ReadMemStats

Programmatic monitoring for dashboards and alerting:

```go
var m runtime.MemStats
runtime.ReadMemStats(&m)

fmt.Printf("Alloc: %d MB\n", m.Alloc/1024/1024)       // currently allocated
fmt.Printf("TotalAlloc: %d MB\n", m.TotalAlloc/1024/1024) // cumulative
fmt.Printf("Sys: %d MB\n", m.Sys/1024/1024)            // requested from OS
fmt.Printf("NumGC: %d\n", m.NumGC)                      // completed collections
fmt.Printf("LastPause: %d ms\n", m.PauseNs[(m.NumGC+255)%256]/1_000_000)
```

### GC pacing

The GC pacer predicts when to start the next collection based on:

1. **Live heap size** after the last collection
2. **GOGC percentage** — how much growth to allow
3. **GOMEMLIMIT** — soft ceiling (if set)
4. **Current allocation rate** — how fast the heap is growing

The pacer starts collection early enough to finish before hitting the target. Fast allocation rates cause earlier starts.

## Allocation Rate Reduction

**Diagnose:** 1- `go tool pprof -alloc_objects` — rank functions by allocation count; the top allocators are where allocation reduction will have the biggest GC impact 2- `GODEBUG=gctrace=1` — monitor GC frequency before and after reducing allocations; expect fewer GC cycles per second as allocation rate drops 3- Prometheus `rate(go_memstats_alloc_bytes_total[5m])` — track allocation rate trend in production; compare before/after deploy to detect regressions

Reducing allocations helps more than tuning GOGC — it addresses the root cause instead of managing the symptom:

- **Value types over pointer types** where possible — values stay on the stack (no GC), pointers escape to the heap
- **Pool frequently allocated objects** with `sync.Pool` (see [memory.md](./memory.md))
- **Preallocate slices and maps** — → See `samber/cc-skills-golang@golang-data-structures` skill
- **Avoid interface boxing** in hot paths — use typed parameters or generics

## GOMAXPROCS in Containers

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for high `runtime.schedule` or `runtime.findRunnable` overhead; this indicates too many P's competing for work or too few P's starving goroutines 2- `go tool trace` — check if goroutines are evenly distributed across P's; uneven distribution suggests GOMAXPROCS is misconfigured for the container 3- `GODEBUG=schedtrace=1000` — print scheduler state every second; look for `runqueue` imbalances or idle P's when work is available 4- `runtime.GOMAXPROCS(0)` — query the current value; if it returns the host CPU count (e.g., 64) instead of the container limit (e.g., 2), the runtime is over-scheduling 5- Prometheus `rate(process_cpu_seconds_total[5m])` — monitor CPU cores consumed in production; if consistently near GOMAXPROCS value, the app is CPU-saturated

**Go 1.25+** improves container CPU detection, particularly for cgroup v2. The runtime sets `GOMAXPROCS` based on:

- Logical CPUs on the machine
- Process CPU affinity mask
- cgroup CPU quota limits (on Linux)

In a container with 2 CPU cores on a 64-core host running Go 1.25+ with **cgroup v2**, `GOMAXPROCS` is correctly set to 2 by default. For **cgroup v1** environments, validate the detected value at startup and consider using `go.uber.org/automaxprocs` to ensure correctness.

**For Go 1.24 and earlier**, use the `go.uber.org/automaxprocs` library to handle container CPU detection:

```go
// Pre-Go 1.25: explicit container-aware detection
import _ "go.uber.org/automaxprocs"

func main() {
    // GOMAXPROCS is now correctly set to container CPU limit
    startServer()
}
```

**Manual override** (if needed):

```bash
GOMAXPROCS=2 ./myapp
GODEBUG=updatemaxprocs=0 ./myapp  # disable dynamic updates (Go 1.25+)
```

**Known limitations (Go 1.25)**: cgroup v1 on certain systems (Oracle OCPUs) may not properly detect Kubernetes CPU limits. Manually set `GOMAXPROCS` as a workaround in these cases.

## Profile-Guided Optimization (PGO)

**Diagnose:** 1- `go tool pprof` (CPU profile) — collect a representative production profile (30+ seconds); look for hot interface method calls and deep call chains that PGO can optimize via devirtualization and inlining 2- `go test -bench` — benchmark before and after placing `default.pgo`; expect 2-7% improvement on interface-heavy code, less on already-optimized paths

Go 1.21+ supports PGO — the compiler uses a production CPU profile to make better inlining and devirtualization decisions. Expected improvement: 2-7% for minimal effort.

**Workflow:**

1. Collect a production CPU profile (30+ seconds of representative load):

   ```bash
   curl http://localhost:6060/debug/pprof/profile?seconds=60 > cpu.pprof
   ```

2. Place as `default.pgo` in the main package directory:

   ```bash
   cp cpu.pprof ./cmd/myapp/default.pgo
   ```

3. Build — `go build` auto-detects `default.pgo`:

   ```bash
   go build ./cmd/myapp
   ```

**What the compiler optimizes:**

- **Inlining** — hot function calls are inlined more aggressively
- **Devirtualization** — interface method calls with high probability of targeting specific types become direct calls

**When it helps most:** code with many interface calls, hot inlining opportunities, deep call stacks. **When it helps least:** already-optimized code, memory-bound workloads.

Rebuild profiles after significant code changes — stale profiles can mislead the compiler.

## Logging Overhead in Hot Paths

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for `fmt.Sprintf`, `log.Printf`, or `slog.(*Logger).log` appearing in hot paths; these indicate log formatting consuming CPU even when the log level filters the message 2- `go build -gcflags="-m"` — check if log arguments escape to the heap; expect `"moved to heap"` for arguments boxed into `any` interface by logging functions 3- `go test -bench -benchmem` — benchmark with logging enabled vs disabled; if allocs/op doesn't change, the logger is allocating even when the level is off

Log formatting allocates memory and consumes CPU even when the message is discarded because it's below the configured level:

```go
// Bad — fmt.Sprintf runs BEFORE the logger checks the level
logger.Debug(fmt.Sprintf("processing item %d with data %v", item.ID, item.Data))

// Good — slog defers formatting until level check passes (Go 1.21+)
slog.Debug("processing item", slog.Int("id", item.ID), slog.Any("data", item.Data))

// Best — LogAttrs: zero allocations when level is disabled
slog.LogAttrs(ctx, slog.LevelDebug, "processing item",
    slog.Int("id", item.ID))
```

In hot paths, even `slog.Any` can allocate. Prefer typed attributes: `slog.Int`, `slog.String`, `slog.Bool`.

## Panic/Recover Cost

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for `runtime.gopanic` or `runtime.gorecover` in the profile; their presence in hot paths means panic/recover is being used for control flow 2- `go test -bench` — benchmark panic/recover vs error-return versions; expect 10-100x overhead from stack unwinding and defer execution

`panic` triggers stack unwinding, running all deferred functions up the call stack. `recover` catches the panic but the unwinding itself is expensive. Never use panic/recover for control flow:

```go
// Bad — panic overhead for a normal condition
defer func() { recover() }()
v, _ := strconv.Atoi(s) // relies on panic for invalid input

// Good — explicit error check, no panic overhead
v, err := strconv.Atoi(s)
if err != nil { continue }
```

Panic is appropriate only for truly unrecoverable situations (programmer errors, corrupted state). Always convert panics to errors at package boundaries.
