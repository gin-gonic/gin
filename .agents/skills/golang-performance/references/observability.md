# Production Observability for Performance

Third-party monitoring tools complement local profiling (pprof, benchmarks) by providing continuous monitoring, historical trends, and regression detection in production.

## Prometheus Metrics for Go

**Setup:** `github.com/prometheus/client_golang` — expose `/metrics` endpoint with `promhttp.Handler()`. Default collectors automatically export Go runtime metrics (`go_goroutines`, `go_memstats_*`, `go_gc_duration_seconds`, `process_cpu_seconds_total`, etc.).

→ See `samber/cc-skills-golang@golang-benchmark` skill (investigation-session.md) for the full runtime metrics table, investigation session setup (scrape interval tuning, env-var toggling), and cost warnings for profiling tools.

### PromQL Queries for Performance Diagnosis

#### GC pressure

| PromQL | What to look for |
| --- | --- |
| `rate(go_gc_duration_seconds_count[5m])` | GC cycles/s — >2/s sustained suggests excessive allocation rate |
| `rate(go_gc_duration_seconds_sum[5m]) / rate(go_gc_duration_seconds_count[5m])` | Average GC pause — increasing trend means heap is growing or has too many pointers |
| `go_gc_duration_seconds{quantile="1"}` | Worst-case GC pause — spikes here cause tail latency |

#### Memory leaks

| PromQL | What to look for |
| --- | --- |
| `go_memstats_alloc_bytes` | Should be roughly stable under constant load; continuous increase = memory leak |
| `rate(go_memstats_alloc_bytes_total[5m])` | Allocation rate (bytes/s) — drives GC frequency; compare before/after deploy for regressions |
| `process_resident_memory_bytes - go_memstats_sys_bytes` | Gap = non-Go memory (cgo, mmap); growing gap = non-Go leak |

#### Goroutine leaks

| PromQL | What to look for |
| --- | --- |
| `go_goroutines` | Should correlate with load; growing independently of traffic = leak |
| `delta(go_goroutines[1h])` | Net goroutine change over 1h; positive without load increase = leak |

#### CPU saturation

| PromQL | What to look for |
| --- | --- |
| `rate(process_cpu_seconds_total[5m])` | CPU cores consumed; compare to GOMAXPROCS to detect saturation |
| `rate(process_cpu_seconds_total[5m]) / <GOMAXPROCS>` | CPU utilization ratio; >0.8 sustained = CPU-saturated |

#### Regression detection (after deploy)

| PromQL | What to look for |
| --- | --- |
| `rate(go_memstats_alloc_bytes_total[5m])` | Compare before/after deploy; significant increase = new allocation pattern introduced |
| `histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))` | p99 latency increase after deploy = regression (requires app-level histogram) |

### Alerting rules (examples)

[Example alerting rules](../assets/prometheus-alerts.yml) — adjust thresholds to your application; a high-throughput data pipeline will have different baselines than a lightweight API server.

→ See `samber/cc-skills@promql-cli` skill for interactively testing these PromQL expressions against your Prometheus instance from the CLI.

### Grafana Dashboards

→ See `samber/cc-skills-golang@golang-observability` skill for recommended community Grafana dashboards that visualize Go runtime metrics out of the box.

## Continuous Profiling

Continuous profiling collects low-overhead samples in production and stores them for historical comparison. Use it to detect regressions across deploys, compare flamegraphs over time, and feed PGO (see [Runtime Tuning](./runtime.md#profile-guided-optimization-pgo)).

| Tool | Model | Overhead | Best for |
| --- | --- | --- | --- |
| **Grafana Pyroscope** | push SDK or pull (via Alloy) | ~2-5% | Grafana ecosystem, historical flamegraph comparison |
| **Parca** (Polar Signals) | eBPF-based pull | <1% | Infrastructure-wide profiling, no code changes |
| **Datadog Continuous Profiler** | push (agent) | ~1-2% | Existing Datadog users |
| **Google Cloud Profiler** | push (agent) | ~1-2% | GCP-hosted Go services |

### Pyroscope push mode

```go
import "github.com/grafana/pyroscope-go"

pyroscope.Start(pyroscope.Config{
    ApplicationName: "myapp",
    ServerAddress:   "http://pyroscope:4040",
    ProfileTypes: []pyroscope.ProfileType{
        pyroscope.ProfileCPU,
        pyroscope.ProfileAllocObjects,
        pyroscope.ProfileAllocSpace,
        pyroscope.ProfileInuseObjects,
        pyroscope.ProfileInuseSpace,
        pyroscope.ProfileGoroutines,
    },
})
```

### Pyroscope pull mode (via Grafana Alloy)

No code changes required — Alloy scrapes `/debug/pprof/*` endpoints periodically. Configure Alloy to target your service's pprof endpoint.

When using third-party profiling libraries, refer to the library's official documentation for current API signatures.

## Real-Time Visualization (Development)

| Tool | What it does |
| --- | --- |
| **statsviz** (`github.com/arl/statsviz`) | Real-time browser dashboard at `/debug/statsviz` — heap, GC pauses, goroutines, scheduler. Register with `statsviz.Register(mux)`. Great for local development |
| **expvar** (stdlib `expvar`) | JSON metrics at `/debug/vars` — lightweight, no dependencies. Integrates with Netdata, Telegraf, or custom dashboards |
