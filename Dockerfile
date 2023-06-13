# [1] build
FROM docker.io/library/golang:1.20 AS build

WORKDIR /src/app
COPY . .

RUN go mod tidy && go build -o /app

# [2] run
FROM scratch
COPY --from=build /app /

ENTRYPOINT ["/app"]
LABEL app.version="0.0.1"

ENV SERVER_GRPC_NETWORK: ${SERVER_GRPC_NETWORK:-tcp} \
    SERVER_GRPC_ADDRESS: ${SERVER_GRPC_ADDRESS:-":9090"} \
    SERVER_RESTGW_ADDRESS: ${SERVER_RESTGW_ADDRESS:-":9080"}

EXPOSE 8080 8090
