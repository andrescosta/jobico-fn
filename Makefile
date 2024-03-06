TARGETS ?= ctl listener repo recorder queue
FORMAT_FILES = $(shell find . -type f -name '*.go' -not -path "*.pb.go")
OUTBINS = $(foreach bin,$(TARGETS),bin/$(bin))

.PHONY: newbin perf1 perf2 k6 go-build test test_coverage test_html checks hadolint init-coverage dckr_build dckr_up dckr_upobs dckr_down dckr_stop lint vuln build release format local $(FORMAT_FILES) $(TARGETS) dockerbuild

APP?=application
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)

MKDIR_REPO_CMD = mkdir -p reports 
MKDIR_BIN_CMD = mkdir -p bin
BUILD_CMD = ./build/build.sh
ENV_CMD = ./build/env.sh
ifeq ($(OS),Windows_NT)
ifneq ($(MSYSTEM), MSYS)
	MKDIR_REPO_CMD = pwsh -noprofile -command "new-item reports -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	MKDIR_BIN_CMD = pwsh -noprofile -command "new-item bin -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	BUILD_CMD = pwsh -noprofile -command ".\build\build.ps1"
	ENV_CMD = pwsh -noprofile -command ".\build\env.ps1"
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
	@$(BUILD_CMD)

env: init-release
	@$(ENV_CMD)

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
dckr_build:
	docker compose -f .\compose\compose.yml build

dckr_up:
	docker compose -f .\compose\compose.yml up -d

dckr_upobs:
	docker compose -f .\compose\compose.yml --profile obs up -d

dckr_down:
	docker compose -f .\compose\compose.yml down 

dckr_stop:
	docker compose -f .\compose\compose.yml stop

### Kind

dockerbuild: $(TARGETS)

$(TARGETS):
	@docker build -f compose/Dockerfile --target $@ -t jobico/$@ . 
	@kind load docker-image jobico/$@:latest
	@kubectl apply -f .\k8s\config\$@.yaml

kindcluster:
	@kind create cluster --config .\k8s\config\cluster.yaml

kind: kindcluster dockerbuild

kinddel: 
	kind delete cluster


