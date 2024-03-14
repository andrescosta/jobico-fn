TARGETS ?= ctl listener repo recorder queue exec 
SUPPORT_TARGETS ?= jaeger prometheus
FORMAT_FILES = $(shell find . -type f -name '*.go' -not -path "*.pb.go")
OUTBINS = $(foreach bin,$(TARGETS),bin/$(bin))

MKDIR_REPO_CMD = mkdir -p reports 
MKDIR_BIN_CMD = mkdir -p bin
BUILD_CMD = ./build/build.sh
ENV_CMD = ./build/env.sh
LINT_INSTALL_CMD = curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.56.1
X509 = ./hacks/cert/create.sh
X509Install = ./hacks/cert/add.sh
DO_SLEEP = sleep 10
GO_TEST_CMD = CGO_ENABLED=1 go test
CERTS_DIR_CMD = mkdir -p .\k8s\certs
CERTS_CMD = openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\$(SVC).key -out .\k8s\certs\$(SVC).crt -subj "/CN=$(SVC)/O=$(SVC)" -addext "subjectAltName = DNS:$(SVC)"
ifeq ($(OS),Windows_NT)
ifneq ($(MSYSTEM), MSYS)
	MKDIR_REPO_CMD = pwsh -noprofile -command "new-item reports -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	MKDIR_BIN_CMD = pwsh -noprofile -command "new-item bin -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	BUILD_CMD = pwsh -noprofile -command ".\build\build.ps1"
	ENV_CMD = pwsh -noprofile -command ".\build\env.ps1"
	DO_SLEEP = pwsh -noprofile -command "Start-Sleep 10"
	X509 = pwsh -noprofile -command ".\hacks\cert\create.bat"
	X509Install = pwsh -noprofile -command ".\hacks\cert\add.bat"
	LINT_INSTALL_CMD = winget install golangci-lint
	GO_TEST_CMD = go test
	CERTS_DIR_CMD = @pwsh -noprofile -command "new-item .\k8s\certs -ItemType Directory -Force -ErrorAction silentlycontinue | Out-Null"
	CERTS_CMD = pwsh -noprofile -command 'openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout .\k8s\certs\$(SVC).key -out .\k8s\certs\$(SVC).crt -subj "/CN=$(SVC)/O=$(SVC)" -addext "subjectAltName = DNS:$(SVC)"'
endif
endif

## Dependencies
.PHONY: dep

dep:
	go install mvdan.cc/gofumpt@latest
	go install golang.org/x/vuln/cmd/govulncheck@latest
	@$(LINT_INSTALL_CMD)

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
	@$(GO_TEST_CMD) -count=1 -race -timeout 60s ./internal/test 

test-coverage: init-coverage
	@$(GO_TEST_CMD)  ./... -coverprofile=./reports/coverage.out

test-html: test_coverage
	go tool cover -html=./reports/coverage.out

## Performance
.PHONY: k6 perf1/docker perf2/docker perf1/k8s perf2/k8s

k6: 
	go install go.k6.io/xk6/cmd/xk6@latest
	xk6 build --with github.com/szkiba/xk6-yaml@latest --output perf/k6.exe

perf1-local: 
	perf/k6.exe run -e HOST_CTL=ctl:50052 -e HOST_REPO=repo:50053 -e HOST_LISTENER=http://listener:8080 -e TLS=false -e TENANT=tenant_1 perf/events.js

perf2-local: 
	perf/k6.exe run -e HOST_CTL=ctl:50052 -e HOST_REPO=repo:50053 -e HOST_LISTENER=http://listener:8080 -e TLS=false -e TENANT=tenant_1 perf/eventsandstream.js

perf1-k8s: 
	perf/k6.exe run -e HOST_CTL=ctl:443 -e HOST_REPO=repo:443 -e HOST_LISTENER=https://listener -e TLS=true -e TENANT=tenant_1 perf/events.js

perf2-k8s: 
	perf/k6.exe run -e HOST_CTL=ctl:443 -e HOST_REPO=repo:443 -e HOST_LISTENER=https://listener -e TLS=true -e TENANT=tenant_1 perf/eventsandstream.js

## Format
.PHONY: $(FORMAT_FILES)  

format: $(FORMAT_FILES)  

$(FORMAT_FILES):
	@gofumpt -w $@

## Docker compose targets.
.PONY: hadolint docker-build docker-up docker-up-obs docker-down docker-stop

hadolint:
	@cat ./compose/Dockerfile | docker run --rm -i hadolint/hadolint

docker-build:
	docker compose -f .\compose\compose.yml build

docker-up:
	docker compose -f .\compose\compose.yml up -d

docker-up-obs:
	docker compose -f .\compose\compose.yml --profile obs up -d

docker-down:
	docker compose -f .\compose\compose.yml down 

docker-stop:
	docker compose -f .\compose\compose.yml stop

## kubernetes targets 

### Kind cluster
.PHONY: kind-delete kind-cluster wait-ingress ingress

kind: kind-cluster ingress docker-images load-images wait-ingress deploy

kind-delete:
	kind delete cluster

kind-cluster:
	@kind create cluster --config ./k8s/config/cluster.yaml

ingress:
	@kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

wait-ingress:
	@$(DO_SLEEP) 
	@kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s

load-images: $(TARGETS:%=load-images/%)
load-images/%: SVC=$*
load-images/%:
	@kind load docker-image jobico/$(SVC):latest

### Container images

docker-images: $(TARGETS:%=docker-images/%)
docker-images/%: SVC=$*
docker-images/%:
	@docker build -f compose/Dockerfile --target $(SVC) -t jobico/$(SVC) . 

### K8s manifests
.PHONY: base create-certs-dir

deploy: base certs supportcerts supportmanifests manifests

base:
	@kubectl apply -f ./k8s/config/namespace.yaml
	@kubectl apply -f ./k8s/config/configmap.yaml

certs: $(TARGETS:%=certs/%)
certs/%: SVC=$*
certs/%:
	@kubectl delete secret $(SVC)-cert --namespace=jobico --ignore-not-found=true
	@kubectl create secret tls $(SVC)-cert --key ./k8s/certs/$(SVC).key --cert ./k8s/certs/$(SVC).crt --namespace=jobico

supportcerts: $(SUPPORT_TARGETS:%=supportcerts/%)
supportcerts/%: SVC=$*
supportcerts/%:
	@kubectl delete secret $(SVC)-cert --namespace=jobico --ignore-not-found=true
	@kubectl create secret tls $(SVC)-cert --key ./k8s/certs/$(SVC).key --cert ./k8s/certs/$(SVC).crt --namespace=jobico

supportmanifests: $(SUPPORT_TARGETS:%=supportmanifests/%)
supportmanifests/%: SVC=$*
supportmanifests/%:
	@kubectl apply -f ./k8s/config/$(SVC).yaml

manifests: $(TARGETS:%=manifests/%)
manifests/%: SVC=$*
manifests/%:
	@kubectl apply -f ./k8s/config/$(SVC).yaml

rollback: rollback-manifests rollback-support-manifests

rollback-manifests: $(TARGETS:%=rollback-manifests/%)
rollback-manifests/%: SVC=$*
rollback-manifests/%:
	@kubectl delete -f ./k8s/config/$(SVC).yaml

rollback-support-manifests: $(SUPPORT_TARGETS:%=rollback-support-manifests/%)
rollback-support-manifests/%: SVC=$*
rollback-support-manifests/%:
	@kubectl delete -f ./k8s/config/$(SVC).yaml

create-certs: create-certs-dir $(TARGETS:%=create-certs/%)
create-certs/%: SVC=$*
create-certs/%:
	@$(CERTS_CMD)

create-certs-dir:
	@$(CERTS_DIR_CMD)

upload-certs-linux:
	sudo cp ./k8s/certs/*.crt /usr/local/share/ca-certificates
	sudo update-ca-certificates

# This target must be run as "admin" on windows
upload-certs-windows/%: upload-certs-windows $(TARGETS:%=upload-certs-windows/%)
upload-certs-windows/%: SVC=$*
upload-certs-windows/%:
	certutil -enterprise -f -v -AddStore "Root" .\k8s\certs\$(SVC)