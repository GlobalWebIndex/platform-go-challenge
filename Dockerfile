# syntax=docker/dockerfile:1

FROM golang:1.25 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o platform-go-challenge ./cmd/server

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=builder /app/platform-go-challenge /platform-go-challenge
EXPOSE 8080
ENTRYPOINT ["/platform-go-challenge"]
