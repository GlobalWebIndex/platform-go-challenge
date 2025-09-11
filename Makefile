APP_NAME=server

build:
	go build -o $(APP_NAME) .

run:
	go run .

docker-build:
	docker build -t platform-go-challenge .

docker-run:
	docker run -p 8080:8080 platform-go-challenge

test:
	go test ./...
