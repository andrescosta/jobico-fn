########################################################################################

OS := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))
VERSION ?= 1.2.3
DBG ?=
GOFLAGS ?=
GOFLAGS := $(GOFLAGS) -modcacherw
GO_VERSION := 1.21
BUILD_IMAGE := golang:1.21.5-bullseye 

BUILD_DIRS := bin/$(OS)_$(ARCH)                   \
              bin/tools                           \
              .go/bin/$(OS)_$(ARCH)               \
              .go/bin/$(OS)_$(ARCH)/$(OS)_$(ARCH) \
              .go/cache                           \
              .go/pkg
go-build: | $(BUILD_DIRS)
	docker run                                                  \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd)://src                                         \
	    -w //src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env GOCACHE="/.cache/gocache"                         \
	    --env GOMODCACHE="/.cache/gomodcache"                   \
	    --env ARCH="$(ARCH)"                                    \
	    --env OS="$(OS)"                                        \
	    --env VERSION="$(VERSION)"                              \
	    --env DEBUG="$(DBG)"                                    \
	    --env GOFLAGS="$(GOFLAGS)"                              \
	    $(BUILD_IMAGE)                                          \
	    ./build/buildc.sh ./...

$(BUILD_DIRS):
	mkdir -p $@
