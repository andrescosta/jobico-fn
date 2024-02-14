FORMAT_FILES = $(shell find . -type f -name '*.go' -not -path "*.pb.go")

.PHONY: newbin perf1 perf2 k6 go-build test test_coverage test_html checks hadolint init-coverage obs up down stop compose lint vuln build release format local $(FORMAT_FILES)

APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)

MKDIR_REPO_CMD = mkdir -p reports 
MKDIR_BIN_CMD = mkdir bin
ifeq ($(OS),Windows_NT)
ifneq ($(MSYSTEM), MSYS)
	MKDIR_REPO_CMD = pwsh -noprofile -command "new-item reports -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	MKDIR_BIN_CMD = pwsh -noprofile -command "new-item bin -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
endif
endif



lint:
	@golangci-lint run ./...

hadolint:
	@cat ./compose/Dockerfile | docker run --rm -i hadolint/hadolint
test:
	go test -count=1 -race -timeout 60s ./internal/test 

test_coverage: init-coverage
	go test ./... -coverprofile=./reports/coverage.out

test_html: test_coverage
	go tool cover -html=./reports/coverage.out

vuln:
	@govulncheck ./...

build: init-release
	./build/build.sh

k6: 
	go install go.k6.io/xk6/cmd/xk6@latest
	xk6 build --with github.com/szkiba/xk6-yaml@latest --output perf/k6.exe

perf1:
	perf/k6.exe run perf/events.js

perf2:
	perf/k6.exe run perf/eventsandstream.js

format: $(FORMAT_FILES)  

$(FORMAT_FILES):
	@gofumpt -w $@

release: checks test env build 

checks: format lint vuln 

local: env build

init-coverage:
	@$(MKDIR_REPO_CMD) 

init-release:
	@$(MKDIR_BIN_CMD) 

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

env: init-release
	./build/env.sh

