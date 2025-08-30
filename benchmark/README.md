# Benchmarking Singleflight Implementations

This repository benchmarks different implementations of the singleflight pattern in Go. The singleflight pattern reduces redundant operations by collapsing multiple concurrent requests for the same key into a single request.

## Implementations

- **StandardSingleflight**
  Baseline using `golang.org/x/sync/singleflight`. `interface{}` API, supports panic/Goexit propagation, a shared-result flag, and **synchronous cleanup** after `fn` finishes.

- **StandardSingleflightCast**
  Same as StandardSingleflight; the benchmark additionally does a type assertion (e.g., `v.(int)`) to measure that overhead. Not a different library—just a benchmark variant.

- **GenericsSingleflight**
  Minimal patch of the standard implementation to add generics (`Group[T]`), hosted at `github.com/catatsuy/sync/singleflight`. Semantics match the standard version (panic/Goexit, shared flag, **synchronous delete**), with slightly fewer allocations.

- **CustomSingleflight**
  Fully custom, generics-based implementation in `github.com/catatsuy/cache` focused on latency and zero allocations. Key differences: **return-first with asynchronous map delete**, per-call mutex to guarantee single execution, **no shared flag**, and **no panic/Goexit handling** (non-goal). Intended for idempotent, finite operations (e.g., cache fills).

> **Contract for CustomSingleflight:** `fn` does not panic, is idempotent, and completes in finite time. If you need panic propagation or the shared flag, use the standard implementation.

## Benchmark Results

EC2 c7g.xlarge (Graviton3, 4 vCPU) / Debian 13 / Go 1.25.0

```
goos: linux
goarch: arm64
pkg: github.com/catatsuy/cache/benchmark
BenchmarkSingleflight/std/keys=1                18591908               196.6 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=1-2              16908397               213.3 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=1-4              10718676               331.7 ns/op            82 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1           18222015               197.3 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1-2         16597257               212.4 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=1-4         10827585               329.7 ns/op            82 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=1           17918259               202.7 ns/op            80 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=1-2         16597784               217.6 ns/op            79 B/op          0 allocs/op
BenchmarkSingleflight/generics/keys=1-4         11186880               308.9 ns/op            75 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1             79731640                42.67 ns/op            0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1-2           24880732               127.3 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=1-4           25826616               152.3 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/std/keys=10               18039288               199.9 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=10-2             17164132               206.7 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std/keys=10-4             12818521               282.4 ns/op            86 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10          18266598               199.8 ns/op            88 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10-2        17554419               204.8 ns/op            87 B/op          1 allocs/op
BenchmarkSingleflight/std-cast/keys=10-4        12830840               281.9 ns/op            86 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=10          17554911               200.8 ns/op            80 B/op          1 allocs/op
BenchmarkSingleflight/generics/keys=10-2        16990652               211.6 ns/op            79 B/op          0 allocs/op
BenchmarkSingleflight/generics/keys=10-4        12554142               286.3 ns/op            78 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10            69217628                49.34 ns/op            0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10-2          32092108               113.9 ns/op             0 B/op          0 allocs/op
BenchmarkSingleflight/custom/keys=10-4          21601365               164.6 ns/op             7 B/op          0 allocs/op
PASS
```

## Observations (EC2 c7g.xlarge, linux/arm64)

- Setup: `go test -bench=. -benchmem -benchtime=3s -cpu=1,2,4` (RunParallel), `keys=1,10`, trivial `fn` (`return i, nil`).
- **CustomSingleflight is consistently fastest.**
  - `keys=1` (worst contention): **42.67 ns/op** vs std **196.6** (@P=1 → ~**4.6×**), **152.3** vs **331.7** (@P=4 → ~**2.2×**).
  - `keys=10` (moderate contention): **49.34** vs **199.9** (@P=1 → ~**4.0×**), **164.6** vs **282.4** (@P=4 → ~**1.7×**).
- **Allocations / memory**
  - CustomSingleflight: **0 allocs/op (≈0 B/op)**.
  - GenericsSingleflight: **0–1 allocs/op (~75–80 B/op)**.
  - Standard / StandardSingleflightCast: **1 alloc/op (~86–88 B/op)**.
- Standard vs StandardSingleflightCast are essentially identical → type assertion cost is negligible.

> Absolute ns/op varies by machine, but the ordering and relative gaps are consistent in our tests.

## Setup and Run

To build the Docker image and run the benchmark:

```bash
docker build -t benchmark-runner .
docker run --rm benchmark-runner
```
