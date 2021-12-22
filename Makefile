
.PHONY: all
all: build test bench

.PHONY: test
test: gofmt
	$(info #Running tests...)
	go test ./...


.PHONY: build
build: gofmt
	$(info #Building...)
	go build -gcflags -m ./...

.PHONY: bench
bench: gofmt
	$(info #Running benchmarks...)
	go test -gcflags -m -bench . ./...

.PHONY: gofmt
gofmt:
	go fmt ./...