.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: build
build:
	./build/build.sh