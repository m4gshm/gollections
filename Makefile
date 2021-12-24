
.PHONY: all
all: build test bench

.PHONY: test
test:
	$(info #Running tests...)
	go clean -testcache
	go test ./...


.PHONY: build
build:
	$(info #Building...)
	# go build -gcflags -m ./...
	go build ./...

.PHONY: bench
bench:
	$(info #Running benchmarks...)
	go test -bench . -benchmem ./...

.PHONY: gofmt
gofmt:
	go fmt ./...