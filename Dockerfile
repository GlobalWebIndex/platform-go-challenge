# Simple, single-stage build & run
FROM golang:1.25.0-alpine

# Install what `go mod` needs
RUN apk add --no-cache git ca-certificates

# Workdir inside the container
WORKDIR /app

# 1) Copy only mod files first to leverage Docker layer cache
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy

# 2) Copy the rest of source
COPY . .

# 3) Build API
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o app ./cmd/api

# 4) Expose the ports app listens on
EXPOSE 8080

# 5) Start the server
CMD ["./app"]
