# Cache Library for Go

This repository provides Go caching primitives tuned for contrasting read and write workloads, plus concurrency helpers and a faster singleflight implementation.

## Overview

- Distinct cache implementations optimized for write-heavy (`WriteHeavyCache`) and read-heavy (`ReadHeavyCache`) access patterns.
- Expiration-aware variants with stale-while-revalidate helpers (`GetWithExpireStatus`) for serving stale data while refreshing asynchronously.
- Integer-specific caches with atomic-like increment operations.
- `RollingCache` for append-and-rotate workloads.
- A generics-based singleflight that trades optional features for lower latency and zero allocations, plus a faster lock manager for keyed locking.
- Benchmarks and Docker automation under `benchmark/` demonstrating performance gains over standard singleflight.

## Installation

Install the library using `go get`:

```sh
go get github.com/catatsuy/cache
```

Ensure Go 1.25+ to match CI (`.github/workflows/go.yml`).

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/catatsuy/cache"
)

func main() {
	c := cache.NewWriteHeavyCache[int, string]()

	c.Set(1, "apple")
	value, found := c.Get(1)

	if found {
		fmt.Println("Found:", value) // Output: Found: apple
	}
}
```

## Cache Implementations

### Write-Heavy Workloads

`WriteHeavyCache` uses `sync.Mutex` for both reads and writes, prioritizing write throughput.

```go
c := cache.NewWriteHeavyCache[int, string]()
c.Set(1, "apple")
value, found := c.Get(1)
```

### Read-Heavy Workloads

`ReadHeavyCache` relies on `sync.RWMutex` to allow concurrent readers while protecting writes.

```go
c := cache.NewReadHeavyCache[int, string]()
c.Set(1, "orange")
value, found := c.Get(1)
```

### Expiration & Stale-While-Revalidate

The expiration variants accept TTLs per entry and expose `GetWithExpireStatus` to support stale-while-revalidate flows.

#### WriteHeavyCacheExpired Example

```go
c := cache.NewWriteHeavyCacheExpired[int, string]() // assumes import "time"
c.Set(1, "apple", 1*time.Second)
fmt.Println(c.Get(1)) // Found: apple

time.Sleep(2 * time.Second)
_, found := c.Get(1)
fmt.Println(found) // false
```

#### ReadHeavyCacheExpired Example

```go
c := cache.NewReadHeavyCacheExpired[int, string]() // assumes import "time"
c.Set(1, "orange", 1*time.Second)
fmt.Println(c.Get(1)) // Found: orange

time.Sleep(2 * time.Second)
_, found := c.Get(1)
fmt.Println(found) // false
```

#### Stale-While-Revalidate Pattern

```go
if v, found, expired := c.GetWithExpireStatus(key); found {
	if expired {
		go func() {
			fresh := fetch(ctx, key)
			c.Set(key, fresh, 5*time.Minute)
		}()
	}
	return v
}
```

### Integer-Specific Caches

`WriteHeavyCacheInteger` and `ReadHeavyCacheInteger` embed increment helpers for counters.

```go
c := cache.NewWriteHeavyCacheInteger[int, int]()
c.Set(1, 100)
c.Incr(1, 10)
value, _ := c.Get(1)
fmt.Println(value) // 110
```

### RollingCache

`RollingCache` maintains ordered slices with efficient append and rotate operations.

```go
c := cache.NewRollingCache[int](10)
c.Append(1)
c.Append(2)
fmt.Println(c.GetItems()) // [1 2]
rotated := c.Rotate()
fmt.Println(rotated)      // [1 2]
fmt.Println(c.GetItems()) // []
```

## Concurrency Utilities

### LockManager

`LockManager` provides keyed locks for coordinating access across goroutines.

```go
lm := cache.NewLockManager[int]()
lm.Lock(1)
// work
lm.Unlock(1)
```

### SingleflightGroup

`SingleflightGroup` prevents duplicate in-flight work for the same key.

```go
sf := cache.NewSingleflightGroup[string]()
value, err, _ := sf.Do("key", func() (string, error) {
	return "Data for key key", nil
})
if err != nil {
	fmt.Println("Error:", err)
} else {
	fmt.Println("Result:", value)
}
```

Combine it with a cache to coalesce heavy loads:

```go
func Get(key int) int {
	if value, found := c.Get(key); found {
		return value
	}
	v, err, _ := sf.Do(fmt.Sprintf("cacheGet_%d", key), func() (int, error) {
		value := HeavyGet(key)
		c.Set(key, value)
		return value, nil
	})
	if err != nil {
		panic(err)
	}
	return v
}
```

## Practical Examples

`github.com/catatsuy/cache` also ships a lightweight cache API that pairs well with Singleflight. The snippets below show how to compose them. Import helper packages such as `fmt` and `time` as needed.

```go
var (
  c  = cache.NewWriteHeavyCache[int, int]()
  sf = cache.NewSingleflightGroup[int]()
)

// Get returns the cached value when present; otherwise it loads it by calling HeavyGet.
// Singleflight makes sure HeavyGet only runs once per key when multiple callers race.
func Get(key int) (int, error) {
  if value, found := c.Get(key); found {
    return value, nil
  }

  v, err := sf.Do(fmt.Sprintf("cacheGet_%d", key), func() (int, error) {
    value := HeavyGet(key)
    c.Set(key, value)
    return value, nil
  })
  if err != nil {
    return 0, err
  }

  return v, nil
}
```

The pattern below serves stale data immediately using `GetWithExpireStatus` and refreshes it once per key via Singleflight.

```go
var (
  c  = cache.NewWriteHeavyCacheExpired[int, int]()
  sf = cache.NewSingleflightGroup[int]()
)

func Get(key int) (int, error) {
  if v, found, expired := c.GetWithExpireStatus(key); found {
    if !expired {
      return v, nil
    }

    go func(k int) {
      sf.Do(fmt.Sprintf("cacheGet_%d", k), func() (int, error) {
        value := HeavyGet(k)
        c.Set(k, value, 1*time.Minute)
        return value, nil
      })
    }(key)
    return v, nil
  }

  v, err := sf.Do(fmt.Sprintf("cacheGet_%d", key), func() (int, error) {
    value := HeavyGet(key)
    c.Set(key, value, 1*time.Minute)
    return value, nil
  })
  if err != nil {
    return 0, err
  }
  return v, nil
}
```

## Benchmarking Singleflight Implementations

The `benchmark/` module compares several singleflight variants in Go. Singleflight collapses concurrent requests sharing a key into a single execution.

### Implementations

- **StandardSingleflight**: Baseline `golang.org/x/sync/singleflight` using `interface{}`, with panic/Goexit propagation, a shared-result flag, and synchronous cleanup after `fn` completes.
- **StandardSingleflightCast**: Same as the baseline, but the benchmark performs a type assertion (for example `v.(int)`) to measure that overhead. This is just a benchmark variant.
- **GenericsSingleflight**: Lightly patched generic port (`Group[T]`) hosted at `github.com/catatsuy/sync/singleflight`. Matches the standard semantics (panic/Goexit, shared flag, synchronous delete) with slightly fewer allocations.
- **CustomSingleflight**: The generics-based implementation shipped in this repository (`github.com/catatsuy/cache`). It focuses on latency and zero allocations via return-first with asynchronous map delete, per-call mutexes, no shared flag, and no panic/Goexit handling. Intended for idempotent, finite work (e.g., cache fills).

> **Contract for CustomSingleflight:** `fn` must not panic, must be idempotent, and must finish in finite time. If you need panic propagation or the shared flag, prefer the standard implementation.

### Benchmark Results

Environment: EC2 c7g.xlarge (Graviton3, 4 vCPU) / Debian 13 / Go 1.25.1

```
goos: linux
goarch: arm64
BenchmarkSingleflight/std/keys=1                18832320               195.1 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=1-2              15887760               225.8 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=1-4              10460737               337.7 ns/op            82 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1           18096949               198.7 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1-2         16042627               221.6 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1-4         10191168               331.6 ns/op            82 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=1           18848503               191.7 ns/op            80 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=1-2         16614574               217.5 ns/op            79 B/op          0 allocs/op
BenchmarkSingleflight/generics/keys=1-4         11035903               323.9 ns/op            75 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1             91318575                42.49 ns/op            0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1-2           26094780               149.9 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1-4           23411012               151.2 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/std/keys=10               18525980               197.5 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=10-2             16850523               215.0 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=10-4             12107134               302.3 ns/op            86 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10          18550858               197.3 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10-2        16768419               214.9 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10-4        12467149               296.0 ns/op            86 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=10          18988800               193.7 ns/op            80 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=10-2        16899808               211.1 ns/op            79 B/op          0 allocs/op
BenchmarkSingleflight/generics/keys=10-4        12377605               286.6 ns/op            78 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10            75470974                49.51 ns/op            0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10-2          28253089               135.4 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10-4          17369714               199.8 ns/op             8 B/op          0 allocs/op
PASS
```

### Observations (EC2 c7g.xlarge, linux/arm64)

![keys=1](benchmark/images/ns_op%20-%20keys=1.png)
![keys=10](benchmark/images/ns_op%20-%20keys=10.png)

- Setup: `go test -bench=. -benchmem -benchtime=3s -cpu=1,2,4` (RunParallel), `keys=1,10`, trivial `fn` (`return i, nil`).
- **CustomSingleflight is consistently fastest.**
  - `keys=1` (worst contention): **42.49 ns/op** vs std **195.1** (@P=1 → ~**4.6×**), **151.2** vs **337.7** (@P=4 → ~**2.2×**).
  - `keys=10` (moderate contention): **49.51** vs **197.5** (@P=1 → ~**4.0×**), **199.8** vs **302.3** (@P=4 → ~**1.5×**).
- **Allocations / memory**
  - CustomSingleflight: **0 allocs/op (≈0 B/op)**.
  - GenericsSingleflight: **0–1 allocs/op (~75–80 B/op)**.
  - Standard / StandardSingleflightCast: **1 alloc/op (~86–88 B/op)**.
- Standard vs StandardSingleflightCast are essentially identical; type assertion cost is negligible.

> Absolute ns/op varies by machine, but the ordering and relative gaps remain similar in our tests.

### Setup and Run

From `benchmark/`, build and run the Dockerized benchmark harness:

```bash
cd benchmark
docker build -t benchmark-runner .
docker run --rm benchmark-runner
```

Or, run the Go benchmarks directly:

```bash
cd benchmark
go test -bench=. -benchmem -benchtime=3s -cpu=1,2,4
```

## Documentation

For full API documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/catatsuy/cache).
