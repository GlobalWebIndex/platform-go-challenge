FROM golang:1.25.1-alpine AS build-stage

WORKDIR /favAssetApp

COPY go.mod ./
RUN go mod download

COPY data ./data
COPY endpoints ./endpoints
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

# Test stage
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Release stage (small runtime image)
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /docker-gs-ping /docker-gs-ping

EXPOSE 3333

USER nonroot:nonroot

ENTRYPOINT ["/docker-gs-ping"]