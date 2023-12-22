GOFMT_FILES = $(shell find . -type f -name '*.go' -not -path "./api/types/*")

.PHONY: obs up down stop compose lint vuln build release gofmt local $(GOFMT_FILES)

APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)


lint:
	golangci-lint run ./...

vuln:
	govulncheck ./...

build:
	./build/build.sh

gofmt: $(GOFMT_FILES)  

$(GOFMT_FILES):
	@gofmt -s -w $@

release: gofmt lint vuln build env 

local: env build

### Docker compose targets.
compose:
	docker compose -f .\compose\compose.yml up

up:
	docker compose -f .\compose\compose.yml up -d

obs:
	docker compose -f .\compose\compose.yml --profile obs up -d

down:
	docker compose -f .\compose\compose.yml down 

stop:
	docker compose -f .\compose\compose.yml stop

env:
	./build/env.sh

