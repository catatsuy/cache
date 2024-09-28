.PHONY: test
test:
	go test -cover ./...

.PHONY: bench
bench:
	go test -bench=. -benchmem

.PHONY: vet
vet:
	go vet ./...
