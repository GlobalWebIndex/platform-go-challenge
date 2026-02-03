# Go specific variables
BINARY_NAME=favorites-api
MAIN_PATH=./cmd/server/main.go
LOCAL_BIN:=$(CURDIR)/bin

# Linter
GOLANGCI_LINT_VERSION=v1.64.5
GOLANGCI_LINT:=$(LOCAL_BIN)/golangci-lint

.PHONY: all build run test-unit test-integration test-load lint lint-install docker-up docker-down clean

all: build

## Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

## Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	go run $(MAIN_PATH)

## Run the application with tracing enabled
run-trace:
	@echo "Running $(BINARY_NAME) with tracing..."
	OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4317 go run $(MAIN_PATH)

## Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -short ./...

## Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -run Integration ./...

## Run load tests
test-load:
	@echo "Running load tests..."
	k6 run tests/k6/load.js

## Install linter
lint-install:
	@echo "Installing golangci-lint via go install..."
	@mkdir -p $(LOCAL_BIN)
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

## Lint the code
lint: lint-install
	@echo "Linting..."
	$(GOLANGCI_LINT) run --timeout=5m

## Start dependencies (Postgres & Redis)
docker-up:
	@echo "Starting services..."
	docker compose up -d

## Stop dependencies
docker-down:
	@echo "Stopping services..."
	docker compose down

## Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin
	go clean
