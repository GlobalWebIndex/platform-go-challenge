# Multi-stage build for smaller image size
FROM golang:1.24-alpine AS builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0 GOOS=linux
RUN go mod download \
 && go build -ldflags="-s -w" -o server main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /build/server .
EXPOSE 8080
CMD ["./server"]
