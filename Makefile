GOFMT_FILES = $(shell find . -type f -name '*.go' -not -path "./api/types/*")

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

build:
	./build/build.sh

gofmt: $(GOFMT_FILES)  

$(GOFMT_FILES):
	@gofmt -s -w $@

release: gofmt lint vuln build 

.PHONY: lint vuln build release gofmt $(GOFMT_FILES)