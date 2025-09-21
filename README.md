# GlobalWebIndex Engineering Challenge

## Introduction

This service follows domain-driven design architecture. The main application can be found in `cmd/main.go`.  In presentation, app and persistence folders there are the handlers, the services and the repositories respectively. Inside the first two there are fakes folders in which you can found fake implementations of methods as obtained by using counterfeiter package. Furthermore, unit tests have been implemented and can be found in there as well. The domain folder contains the definitions of structures and interfaces. Also, a Dockerfile is created to build the image of this service and a docker compose file to run it. 

An ArangoDB database with elements is created automatically when the service is up and running.  

## Packages Installation

- `go get -u github.com/maxbrunsfeld/counterfeiter/v6` 
- `go get github.com/gorilla/mux` 
- `github.com/arangodb/go-driver`


## Environment Variables

- `ARANGO_URL`: ArangoDB url 
- `ARANGO_USERNAME`: ArangoDB username, default value is 'root'.
- `ARANGO_PASSWORD`: ArangoDB password
- `ARANGO_DB`: ArangoDB database
- `SERVE_ON_PORT`: The port that this service runs on, default value is 8000.


## Unit tests

- `Counterfeiter` package is used in order to generate fake implementations of the methods for both service and repository.
- Unit tests for both handlers and services have been implemented in presentation and app/services folders respectively.
- Run `go test ./...` in command line to run all the tests at once.


## Docker

- Run `docker-compose up --build` to build the image and run the service.


## Endpoints

    curl --request POST
         --url http://localhost:8000/assets/{asset_id}/starred 
         --header 'user_id: {user_id}'

    curl --request GET 
         --url http://localhost:8000/assets/starred
         --header 'user_id: {user_id}'

    curl --request PUT
         --url http://localhost:8000/assets/{asset_id}/edit 
         --header 'Content-Type: application/json' 
         --header 'user_id: {user_id}' 
         --data '{"description": "The updated description"}'

    curl --request DELETE 
         --url http://localhost:8000/assets/{asset_id}/delete 
         --header 'Content-Type: application/json' 
         --header 'user_id: {user_id}'
