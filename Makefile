FORMAT_FILES = $(shell find . -type f -name '*.go' -not -path "*.pb.go")

.PHONY: go-build checks hadolint init gosec obs up down stop compose lint vuln build release format local $(FORMAT_FILES)

APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)

MKDIR_REPO_CMD = mkdir -p reports 
ifeq ($(OS),Windows_NT)
ifneq ($(MSYSTEM), MSYS)
	MKDIR_REPO_CMD = pwsh -noprofile -command "new-item reports -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
endif
endif

lint:
	@golangci-lint run ./...

hadolint:
	@cat ./compose/Dockerfile | docker run --rm -i hadolint/hadolint

vuln:
	@govulncheck ./...

gosec: init
	@gosec -quiet -out ./reports/gosec.txt ./... 

build:
	./build/build.sh

format: $(FORMAT_FILES)  

$(FORMAT_FILES):
	@gofumpt -w $@

release: checks build env

checks: format lint vuln gosec 

local: env build

init:
	@$(MKDIR_REPO_CMD) 



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

