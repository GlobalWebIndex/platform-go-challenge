APP_NAME := favourites-api
CMD_DIR  := ./cmd/api
ADDR     ?= :8080
API_URL  ?= http://localhost:8080

.PHONY: help
help:
	@echo "Targets:"
	@echo "  make run              Run API locally (ADDR=$(ADDR))"
	@echo "  make build            Build local binary into ./bin"
	@echo "  make test             Run all tests with coverage"
	@echo "  make lint             Run basic checks (fmt + vet)"
	@echo "  make tidy             Go module tidy"
	@echo "  make swagger          Generate swagger docs into ./docs"
	@echo "  make docker-build     Build docker image"
	@echo "  make docker-run       Run docker container"
	@echo "  make compose-up       Run docker compose (build + up)"
	@echo "  make compose-down     Stop docker compose"
	@echo "  make seed             Seed realistic data (TOTAL=300 USER_ID=user-1)"
	@echo "  make paginate         Walk pagination (USER_ID=user-1 LIMIT=50 EXPECTED=0)"
	@echo "  make cache-check      Verify cache hits/misses (USER_ID=user-1)"
	@echo "  make golden           Run golden path script (USER_ID=user-1)"
	@echo "  make clean            Remove local build artifacts"

.PHONY: run
run:
	ADDR=$(ADDR) go run $(CMD_DIR)

.PHONY: build
build:
	@mkdir -p bin
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/$(APP_NAME) $(CMD_DIR)

.PHONY: test
test:
	go test ./... -cover -count=1

.PHONY: lint
lint:
	go fmt ./...
	go vet ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: swagger
swagger:
	swag init -g cmd/api/main.go -d cmd/api,internal -o docs

.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME):latest .

.PHONY: docker-run
docker-run:
	docker run --rm -p 8080:8080 -e ADDR=:8080 -e ENV=development $(APP_NAME):latest

.PHONY: compose-up
compose-up:
	docker compose up --build

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: seed
seed:
	API_URL=$(API_URL) USER_ID=$${USER_ID:-user-1} TOTAL=$${TOTAL:-300} ./scripts/01_seed_many_assets.sh

.PHONY: paginate
paginate:
	API_URL=$(API_URL) USER_ID=$${USER_ID:-user-1} LIMIT=$${LIMIT:-50} EXPECTED=$${EXPECTED:-0} ./scripts/02_paginate_walk.sh

.PHONY: cache-check
cache-check:
	API_URL=$(API_URL) USER_ID=$${USER_ID:-user-1} ./scripts/03_cache_check.sh

.PHONY: golden
golden:
	API_URL=$(API_URL) USER_ID=$${USER_ID:-user-1} ./scripts/04_golden_path.sh

.PHONY: clean
clean:
	rm -rf bin
