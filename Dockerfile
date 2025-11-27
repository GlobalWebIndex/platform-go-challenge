FROM golang:1.24 AS builder

WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o favorites ./cmd/main.go

FROM alpine:3
WORKDIR /root/
COPY --from=builder /app/favorites .
EXPOSE 8080
CMD ["./favorites"]