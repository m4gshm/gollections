
.PHONY: all
all: clean build test readme lint bench

.PHONY: test
test:
	$(info #Running tests...)
	go test ./...

.PHONY: cover
cover:
	$(info #Running cover tests...)
	go test -coverprofile=coverage.out -coverpkg=github.com/m4gshm/gollections/...  ./...
	go tool cover -html=coverage.out

.PHONY: cover-console-out
cover-console-out:
	$(info #Running cover tests...)
	go test -coverprofile=coverage.out -coverpkg=github.com/m4gshm/gollections/...  ./...
	go tool cover -func=coverage.out

.PHONY: clean
clean:
	$(info #Building...)
	go clean -cache
	go clean -testcache

.PHONY: build
build:
	$(info #Building...)
	go env -w GOEXPERIMENT=newinliner
	go fmt ./...
	go vet ./...
	go build ./...

.PHONY: builda
builda:
	$(info #Building...)
	go clean -cache
	go build -gcflags -m ./...

.PHONY: bench
bench:
	$(info #Running benchmarks...)
	# go test -gcflags=-d=loopvar=3 -benchtime 1s -bench . -benchmem ./...
	go env -w GOEXPERIMENT=newinliner
	go test -benchtime 1s -bench . -benchmem ./...

.PHONY: update
update:
	$(info #Undate tools...)
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/tools/gopls@latest

.PHONY: lint
lint:
	$(info #Lints...)
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	golangci-lint run

.PHONY: readme
readme:
	$(info #README.md...)
	asciidoctor -b docbook internal/docs/readme.adoc 
	pandoc -f docbook -t gfm internal/docs/readme.xml -o README.md	

.PHONY: pkgsite
pkgsite:
	go install golang.org/x/pkgsite/cmd/pkgsite@latest && pkgsite 