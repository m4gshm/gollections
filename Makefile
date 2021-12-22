
.PHONY: all
all: build test

.PHONY: test
test:
	$(info #Running tests...)
	go test ./...


.PHONY: build
build:
	$(info #Building...)
	go build -gcflags -m ./...

.PHONY: bench
bench:
	$(info #Running benchmarks...)
	go test -gcflags -m -bench . ./...

