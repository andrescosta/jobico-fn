#!/usr/bin/env bash

set -o errexit
set -o nounset

go install github.com/golangci/golangci-lint/cmd/golangci-lint

printf "Running golangci-lint: "
ERRS=$(golangci-lint run "$@" 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"
echo