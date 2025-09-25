# Dockerfile
FROM golang:1.23-alpine

WORKDIR /app

# Cache deps
COPY go.mod go.sum ./
RUN go mod download

# Build
COPY . .
RUN CGO_ENABLED=0 go build -o server ./src/scripts/server

ENV ADDR=:8080
EXPOSE 8080
CMD ["./server"]