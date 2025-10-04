# Repository Guidelines

This repository implements Go caching primitives, singleflight variants, and supporting tooling. Follow these practices to keep contributions fast and reliable.

## Project Structure & Module Organization
- Core library lives at the root (`cache.go`, `singleflight.go`, plus matching `*_test.go` suites and `example_test.go`).
- Performance experiments reside in `benchmark/`, a standalone Go module with its own `go.mod`, documentation, and result images.
- CI automation is defined in `.github/workflows/go.yml`; local helper targets are collected in `Makefile`.

## Build, Test, and Development Commands
- `make vet`: runs `go vet ./...` to surface static-analysis issues.
- `make test`: executes `go test -cover ./...` for unit and example coverage.
- `make bench`: invokes `go test -C benchmark -modfile=go.mod -bench=. -benchmem` against the benchmark module.
- `go test -C benchmark -modfile=go.mod -bench=. -benchmem -benchtime=3s -cpu=1,2,4`: mirrors the CI benchmark matrix when you need production-grade numbers.

## Coding Style & Naming Conventions
- Format Go code with `gofmt` or `goimports` before committing; no custom lint exceptions are in place.
- Follow Go’s idiomatic naming: exported types like `WriteHeavyCache` describe behavior, tests use `TestXxx`/`BenchmarkXxx`, and receiver names mirror type initials.
- Group related helpers with short comments only when intent is non-obvious.

## Testing Guidelines
- Keep tests alongside implementations in `*_test.go`; prefer table-driven `t.Run` cases for clarity and reuse.
- Example code belongs in `example_test.go` to keep documentation snippets executable.
- Regenerate coverage with `make test`; avoid committing transient artifacts such as `coverage.out` unless they illustrate a documentation change.

## Benchmark & Performance Notes
- Extend scenarios in `benchmark/benchmark_test.go`; document new cases in `benchmark/README.md` and refresh plots in `benchmark/images/` when results change.
- Include before/after numbers in PR descriptions whenever you touch performance-sensitive paths.

## Commit & Pull Request Guidelines
- Write imperative, single-line commit subjects (e.g., `Add stale flag to ReadHeavyCache`) followed by optional context in the body.
- Reference related issues, summarize functional and performance impacts, and list the validation you ran (`make vet`, `make test`, benchmark command).
- Confirm GitHub Actions is green before requesting review; include screenshots or tables for benchmark-driven changes when helpful.
