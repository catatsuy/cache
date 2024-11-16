# Benchmarking Singleflight Implementations

This repository benchmarks different implementations of the singleflight pattern in Go. The singleflight pattern reduces redundant operations by collapsing multiple concurrent requests for the same key into a single request.

## Implementations

- **StandardSingleflight**: Uses the `golang.org/x/sync/singleflight` package.
- **StandardSingleflightCast**: Similar to `StandardSingleflight` but includes type assertions.
- **GenericsSingleflight**: A patched version of `golang.org/x/sync/singleflight` to support Go generics, implemented in the https://github.com/catatsuy/sync repository.
- **CustomSingleflight**: A fully custom implementation optimized for speed and memory usage.

## Benchmark Results

```
goos: linux
goarch: arm64
pkg: github.com/catatsuy/cache/benchmark
BenchmarkStandardSingleflight-8          1000000              1081 ns/op             168 B/op          2 allocs/op
BenchmarkStandardSingleflightCast-8      1000000              1115 ns/op             126 B/op          2 allocs/op
BenchmarkGenericsSingleflight-8          1000000              1145 ns/op             102 B/op          2 allocs/op
BenchmarkCustomSingleflight-8            1974880               573.0 ns/op           110 B/op          2 allocs/op
PASS
ok      github.com/catatsuy/cache/benchmark     5.526s
```

## Observations

- `CustomSingleflight` is:
  - **1.89x faster** than `StandardSingleflight`.
  - **1.95x faster** than `StandardSingleflightCast`.
  - **2.00x faster** than `GenericsSingleflight`.
- `CustomSingleflight` also uses the least memory per operation, with only 110 B/op.

## Setup and Run

To build the Docker image and run the benchmark:

```bash
docker build -t benchmark-runner .
docker run --rm benchmark-runner
```
