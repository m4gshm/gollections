
.PHONY: all
all: build test bench

.PHONY: test
test:
	$(info #Running tests...)
	go clean -testcache
	go test ./...


.PHONY: build
build: gofmt govet errcheck
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

.PHONY: govet
govet:
	go vet ./...


.PHONY: errcheck
errcheck:
	go install github.com/kisielk/errcheck@latest
	errcheck ./...