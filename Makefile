GOFMT_FILES = $(shell find . -type f -name '*.go' -not -path "./api/types/*")

APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)

lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

build:
	./build/build.sh

env:
	./build/env.sh

gofmt: $(GOFMT_FILES)  

$(GOFMT_FILES):
	@gofmt -s -w $@

release: gofmt lint vuln build env 

local: env build

.PHONY: lint vuln build release gofmt local $(GOFMT_FILES)