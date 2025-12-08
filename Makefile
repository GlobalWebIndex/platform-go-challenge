MODULE_NAME := github.com/GlobalWebIndex/platform-go-challenge

.PHONY: all tidy fmt vet build test run

all: tidy fmt vet build test

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet ./...

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

run:
	go run ./cmd/api