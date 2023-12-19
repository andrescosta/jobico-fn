GOFMT_FILES = $(shell go list -f '{{.Dir}}' ./... | grep -v '\types')

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: vuln
vuln:
	govulncheck ./...

.PHONY: build
build:
	./build/build.sh

.PHONY: release
release: lint vuln build

.PHONE: gofmt
gofmt: $(GOFMT_FILES) 

$(GOFMT_FILES):
    gofmt -s $@