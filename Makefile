GO = go
GOFLAGS = -ldflags="-s -w"
CTL_PACKAGE_PATH := ./ctl/cmd
QUEUE_PACKAGE_PATH := ./srv/cmd/queue
RECORDER_PACKAGE_PATH := ./recorder/cmd
LISTENER_PACKAGE_PATH := ./srv/cmd/listener
EXEC_PACKAGE_PATH := ./srv/cmd/executor
REPO_PACKAGE_PATH := ./repo/cmd
TOOLS_CLI_PACKAGE_PATH := ./tools/cmd/cli
TOOLS_DASHB_PACKAGE_PATH := ./tools/cmd/dashboard

.PHONY: build
build:
    go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

