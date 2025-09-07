# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /server ./cmd/server

# Stage 2: Create the final, minimal image
FROM alpine:latest

WORKDIR /

COPY --from=builder /server .

EXPOSE 8080

CMD ["/server"]