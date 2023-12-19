GOFMT_FILES:= $(shell find . -type f -name '*.go' -not -path "./api/types/*")

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

.PHONY: gofmt
gofmt:  
	@gofmt -s -l -w $(SRC)