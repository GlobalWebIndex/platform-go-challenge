#!/bin/sh
set -e

# lint
echo "Checking lint"
golint -set_exit_status=1 `go list -mod=vendor ./...`
echo "Lint success!"

# test
echo "Running tests"
go test ./... -race -coverprofile=coverage.txt -covermode=atomic -p=1
echo "Testing success!"