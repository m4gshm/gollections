
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

.PHONY: lint
lint:
	$(info #Lints...)
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	go vet ./...
	# go install github.com/tetafro/godot/cmd/godot@latest
	# godot ./:
	go install github.com/kisielk/errcheck@latest
	errcheck ./...
	go install github.com/alexkohler/nakedret@latest
	nakedret ./...