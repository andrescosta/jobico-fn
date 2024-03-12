TARGETS ?= ctl listener repo recorder queue exec 
SUPPORT_TARGETS ?= jaeger prometheus
FORMAT_FILES = $(shell find . -type f -name '*.go' -not -path "*.pb.go")
OUTBINS = $(foreach bin,$(TARGETS),bin/$(bin))

MKDIR_REPO_CMD = mkdir -p reports 
MKDIR_BIN_CMD = mkdir -p bin
BUILD_CMD = ./build/build.sh
ENV_CMD = ./build/env.sh
X509 = ./hacks/c.sh
DO_SLEEP = sleep 10
ifeq ($(OS),Windows_NT)
ifneq ($(MSYSTEM), MSYS)
	MKDIR_REPO_CMD = pwsh -noprofile -command "new-item reports -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	MKDIR_BIN_CMD = pwsh -noprofile -command "new-item bin -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	BUILD_CMD = pwsh -noprofile -command ".\build\build.ps1"
	ENV_CMD = pwsh -noprofile -command ".\build\env.ps1"
	DO_SLEEP = pwsh -noprofile -command "Start-Sleep 10"
	X509 = pwsh -noprofile -command "./hacks/c.ps1"
endif
endif

## Release
.PHONY: init-release
init-release:
	@$(MKDIR_BIN_CMD) 

release: format checks test env build 

build: init-release
	@$(BUILD_CMD)

env: init-release
	@$(ENV_CMD)

## Local environment

local: env build

### Validations
.PHONY: lint vuln

checks: lint vuln 

lint:
	@golangci-lint run ./...

vuln:
	@govulncheck ./...

## Tests
.PHONY: init-coverage test

init-coverage:
	@$(MKDIR_REPO_CMD) 

test:
	go test -count=1 -race -timeout 60s ./internal/test 

test_coverage: init-coverage
	go test ./... -coverprofile=./reports/coverage.out

test_html: test_coverage
	go tool cover -html=./reports/coverage.out

## Performance
.PHONY: k6 perf1 perf2

k6: 
	go install go.k6.io/xk6/cmd/xk6@latest
	xk6 build --with github.com/szkiba/xk6-yaml@latest --output perf/k6.exe

perf1: 
	perf/k6.exe run perf/events.js

perf2: 
	perf/k6.exe run perf/eventsandstream.js

## Format
.PHONY: $(FORMAT_FILES)  

format: $(FORMAT_FILES)  

$(FORMAT_FILES):
	@gofumpt -w $@

## Docker compose targets.
.PONY: hadolint dckr_build dckr_up dckr_upobs dckr_down dckr_stop

hadolint:
	@cat ./compose/Dockerfile | docker run --rm -i hadolint/hadolint

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

## kubernetes targets 

### Kind cluster
.PHONY: kinddel kindcluster waitnginx

kind: kindcluster dockerimages waitnginx deploy

kinddel: 
	kind delete cluster

kindcluster:
	@kind create cluster --config .\k8s\config\cluster.yaml
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

waitnginx:
	@$(DO_SLEEP) 
	@kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

### Container images
dockerimages: $(TARGETS:%=dockerimages/%)
dockerimages/%: SVC=$*
dockerimages/%:
	@docker build -f compose/Dockerfile --target $(SVC) -t jobico/$(SVC) . 
	@kind load docker-image jobico/$(SVC):latest

### K8s manifests
.PHONY: base

deploy: base certs supportcerts supportmanifests manifests

base:
	@kubectl apply -f .\k8s\config\namespace.yaml
	@kubectl apply -f .\k8s\config\configmap.yaml

certs: $(TARGETS:%=certs/%)
certs/%: SVC=$*
certs/%:
	@kubectl delete secret $(SVC)-cert --namespace=jobico --ignore-not-found=true
	@kubectl create secret tls $(SVC)-cert --key .\k8s\certs\$(SVC).key --cert .\k8s\certs\$(SVC).crt --namespace=jobico

supportcerts: $(SUPPORT_TARGETS:%=supportcerts/%)
supportcerts/%: SVC=$*
supportcerts/%:
	@kubectl delete secret $(SVC)-cert --namespace=jobico --ignore-not-found=true
	@kubectl create secret tls $(SVC)-cert --key .\k8s\certs\$(SVC).key --cert .\k8s\certs\$(SVC).crt --namespace=jobico

supportmanifests: $(SUPPORT_TARGETS:%=supportmanifests/%)
supportmanifests/%: SVC=$*
supportmanifests/%:
	@kubectl apply -f .\k8s\config\$(SVC).yaml

manifests: $(TARGETS:%=manifests/%)
manifests/%: SVC=$*
manifests/%:
	@kubectl apply -f .\k8s\config\$(SVC).yaml

rollback: $(TARGETS:%=rollback/%)
rollback/%: SVC=$*
rollback/%:
	@kubectl delete -f .\k8s\config\$(SVC).yaml

## Hacks
.PHONY: x509
x509:
	@$(X509)