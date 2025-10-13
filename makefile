.PHONY: help build run test clean docker-build docker-up docker-down

help:
	@echo "Available commands:"
	@echo "  make build         - Build the Go binary"
	@echo "  make run           - Run the server locally"
	@echo "  make test          - Run tests"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-up     - Start Docker Compose"
	@echo "  make docker-down   - Stop Docker Compose"
	@echo "  make clean         - Remove binaries and containers"

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test -cover ./...

docker-build:
	docker compose build

docker-up:
	docker compose up

docker-down:
	docker compose down

clean:
	rm -rf bin/
	docker compose down
	docker rmi platform-go-challenge-platform-go-challenge:latest || true