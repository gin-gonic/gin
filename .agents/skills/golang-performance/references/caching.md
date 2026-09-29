# Caching Patterns

The fastest code is code that doesn't run. Caching pre-computed results, deduplicating concurrent requests, and avoiding unnecessary work are often the highest-leverage performance improvements.

## Compiled Pattern Caching

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for `regexp.Compile`, `regexp.MustCompile`, or `template.Parse` appearing in hot paths; their presence means patterns are being recompiled per call instead of once 2- `go test -bench -benchmem` — benchmark per-call compilation vs cached version; expect 10-12x improvement and allocs/op dropping to zero for the compilation step

### Regexp at package level

`regexp.Compile` parses a pattern into a state machine — ~5,700ns per compilation. Match operations on a compiled regexp cost ~450ns. Compiling per-call wastes 10-12x:

```go
// Bad — compiled on every call
func isValid(email string) bool {
    re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
    return re.MatchString(email)
}

// Good — compiled once, safe for concurrent use
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)

func isValid(email string) bool { return emailRegex.MatchString(email) }
```

Note: `regexp.MustCompile` panics on invalid patterns — fine for package-level constants (caught at startup). Use `regexp.Compile` for user-provided patterns. Go's regexp uses linear-time matching (no backtracking).

### Template caching

`template.Parse` is equally expensive. Parse once at startup:

```go
var reportTmpl = template.Must(template.ParseFiles("templates/report.html"))
```

### Precomputed lookup tables

When a computation is pure (same input → same output) and the input space is small, replace calculation with array lookup:

```go
var hexDigit = [16]byte{'0','1','2','3','4','5','6','7','8','9','a','b','c','d','e','f'}

func byteToHex(b byte) (byte, byte) {
    return hexDigit[b>>4], hexDigit[b&0x0f] // two array lookups vs branching logic
}
```

If the table fits in L1/L2 cache, lookup is faster than even simple computation.

## Request-Level Caching

**Diagnose:** 1- `go tool pprof` (goroutine profile) — look for many goroutines blocked on the same external call (HTTP fetch, DB query); this signals a cache stampede where N goroutines all miss the cache simultaneously 2- `fgprof` — shows off-CPU wait time; look for the same fetch function dominating wall-clock time across many goroutines, confirming duplicated concurrent work 3- `go tool pprof -alloc_objects` — check if cache miss handling allocates heavily; high alloc counts on fetch functions confirm the stampede is also generating GC pressure

### singleflight for cache stampede prevention

When a cache entry expires, many goroutines may simultaneously discover the miss and all request the same expensive computation. `singleflight` ensures only one goroutine fetches while others wait:

```go
import "golang.org/x/sync/singleflight"

var (
    cache sync.Map
    sf    singleflight.Group
)

func GetWeather(city string) (string, error) {
    if val, ok := cache.Load(city); ok {
        return val.(string), nil
    }

    // Only one goroutine fetches; others block on the same key
    result, err, _ := sf.Do(city, func() (any, error) {
        data, err := fetchFromAPI(city)
        if err == nil { cache.Store(city, data) }
        return data, err
    })
    return result.(string), err
}
```

→ See `samber/cc-skills-golang@golang-concurrency` skill for `singleflight` API details and `sync.Map` vs `RWMutex` decision guidance. → **Generics alternative:** Use `github.com/samber/go-singleflightx` to avoid interface{} boxing overhead; expect 2-4x faster result retrieval compared to the standard library's `singleflight.Group`.

### LRU caches

For bounded caches with eviction, the standard library's `container/list` works but has poor cache locality (each node is a separate heap allocation). For high-performance LRU:

- **`github.com/hashicorp/golang-lru`** — thread-safe, simple API
- **`github.com/elastic/go-freelru`** — merges hashmap and ringbuffer into contiguous memory, ~37x faster than sharded implementations

When using third-party cache libraries, refer to the library's official documentation for current API signatures.

## Algorithmic Complexity

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for functions with high cumulative time that contain nested loops or repeated linear scans; these are algorithmic complexity bottlenecks 2- `go test -bench` — benchmark with different input sizes (100, 1K, 10K, 100K); if time grows quadratically (10x input → 100x time), the algorithm is O(n²) and needs replacement

Before micro-optimizing, check that the algorithm itself isn't the bottleneck. A constant-factor improvement on an O(n²) algorithm loses to a naive O(n log n) implementation at scale.

**Common complexity traps in Go:**

| Pattern | Complexity | Fix | Fixed complexity |
| --- | --- | --- | --- |
| `slices.Contains` in a loop | O(n·m) | Build `map[T]struct{}` first, then lookup | O(n+m) |
| Nested loops for matching | O(n²) | Index with a map, sort+binary search, or `slices.BinarySearch` | O(n log n) or O(n) |
| Repeated `append` without prealloc | O(n²) amortized copies | `make([]T, 0, n)` | O(n) |
| String concatenation with `+=` | O(n²) total copies | `strings.Builder` | O(n) |
| Linear scan for min/max/dedup | O(n) per query | Sort once, query many times | O(n log n) + O(log n) per query |

**Think in Big-O first, then optimize constants.** A 10x constant-factor improvement matters; switching from O(n²) to O(n) matters more.

## Work Avoidance

**Diagnose:** 1- `go tool pprof` (CPU profile) — look for linear scan functions (`slices.Contains`, `slices.Index`) or iterator chains (`Filter`, `Map`) consuming CPU in hot paths 2- `go test -bench` — benchmark the current approach vs a map-based or early-return version; expect O(n) → O(1) for membership tests, significant improvement for short-circuit loops

### Map lookups over slice scanning

`Contains(slice, element)` is O(n). Map lookups are O(1). When doing multiple membership tests against the same collection, build a map once:

```go
// Bad — O(n*m), checking Contains per element
for _, item := range subset {
    if !Contains(collection, item) { return false } // O(n) per check
}

// Good — O(n+m), build map once, O(1) lookups
seen := make(map[T]struct{}, len(collection))
for _, item := range collection { seen[item] = struct{}{} }
for _, item := range subset {
    if _, ok := seen[item]; !ok { return false }
}
```

Use `struct{}` (0 bytes) instead of `bool` (1 byte) for set maps.

### Early returns and short-circuit loops

Return immediately when the answer is known. Finding the target on iteration 3 of 1000 saves 997 iterations:

```go
// Bad — always iterates full collection
found := false
for _, item := range collection {
    if item == target { found = true }
}
return found

// Good — returns on first match
for i := range collection {
    if collection[i] == target { return true }
}
return false
```

### Avoid iterator chains

Chaining iterator operations (`Filter → Map → First`) creates closures and intermediate machinery. A direct loop is simpler and faster:

```go
// Bad — creates 2 iterators with closures
result, ok := First(Filter(collection, predicate))

// Good — single pass, early return, no closures
for i := range collection {
    if predicate(collection[i]) { return collection[i], true }
}
```

### Replace indirect function calls with direct loops

When a function wraps another function (e.g., `FromSlicePtr` calling `Map` with a closure), the closure indirection prevents inlining. Replace with a direct loop:

```go
// Bad — Map() with closure, per-element function call overhead
func FromSlicePtr(items []*T) []T {
    return Map(items, func(p *T) T { return *p })
}

// Good — direct loop, inlineable, -13% to -17% time
func FromSlicePtr(items []*T) []T {
    result := make([]T, len(items))
    for i := range items { result[i] = *items[i] }
    return result
}
```
