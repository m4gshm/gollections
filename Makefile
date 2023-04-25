
.PHONY: all
all: build test readme lint bench

.PHONY: test
test:
	$(info #Running tests...)
	go clean -testcache
	go test ./...

.PHONY: build
build:
	$(info #Building...)
	go clean -cache
	go build ./...

.PHONY: builda
builda:
	$(info #Building...)
	go clean -cache
	go build -gcflags -m ./...

.PHONY: bench
bench:
	$(info #Running benchmarks...)
	go test -benchtime 1s -bench . -benchmem ./...

.PHONY: update
update:
	$(info #Undate tools...)
	go install github.com/go-delve/delve/cmd/dlv@latest
	go install golang.org/x/tools/gopls@latest

.PHONY: lint
lint:
	$(info #Lints...)
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .
	# go vet ./...
	# go install github.com/tetafro/godot/cmd/godot@latest
	# godot ./:
	go install github.com/kisielk/errcheck@latest
	errcheck -ignoretests ./...
	go install github.com/alexkohler/nakedret/cmd/nakedret@latest
	nakedret ./...
	# go install golang.org/x/lint/golint@latest
	# golint ./...
	go install github.com/mgechev/revive@latest
	revive ./...

.PHONY: readme
readme:
	$(info #README.md...)
	asciidoctor -b docbook internal/docs/readme.adoc 
	pandoc -f docbook -t gfm internal/docs/readme.xml -o README.md	
