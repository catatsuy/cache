.PHONY: test
test:
	go test -cover ./...

.PHONY: bench
bench:
	go test -C benchmark -modfile=go.mod -bench=. -benchmem

.PHONY: vet
vet:
	go vet ./...
