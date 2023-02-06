generate:
	cd api; buf generate

run:
	cd cmd; go run main.go

test:
	go test -v ./...

clean:
	rm -rf pkg

