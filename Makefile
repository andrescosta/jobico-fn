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

## {+} Docker targets
docker-build: build
	docker build -t ${APP} .
	docker tag ${APP} ${APP}:${COMMIT_SHA}

docker-push: check-environment docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

check-environment:
ifndef APP_ENV
    $(error ENV not set, allowed values - `staging` or `production`)
endif

## {-} Docker Targets

env:
	./build/env.sh

gofmt: $(GOFMT_FILES)  

$(GOFMT_FILES):
	@gofmt -s -w $@

release: gofmt lint vuln build env 

local: env build

## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

.PHONY: lint vuln build release gofmt local help $(GOFMT_FILES)