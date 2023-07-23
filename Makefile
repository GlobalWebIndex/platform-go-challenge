generate:
	cd api; buf generate; mv ../pkg/microservice.swagger.json ../static/swagger-ui/swagger.yaml

run:
	cd cmd; go run main.go

test:
	go test -v ./...

clean:
	rm -rf pkg

build:
	go build -o ./build/gwoi_crm ./cmd/main.go
