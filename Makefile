
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
	go clean -cache
	go build ./...

.PHONY: bench
bench:
	$(info #Running benchmarks...)
	go test -benchtime 1s -bench . -benchmem ./...

.PHONY: gofmt
gofmt:
	go version
	go fmt ./...

.PHONY: govet
govet:
	go vet ./...

.PHONY: govet
godot:
#	go install github.com/tetafro/godot/cmd/godot@latest
	godot ./

.PHONY: errcheck
errcheck:
#	go install github.com/kisielk/errcheck@latest
	errcheck ./...