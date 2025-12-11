# --- Build stage ---
FROM golang:1.25.5-alpine AS builder

WORKDIR /src

# Enable Go modules and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the API binary to a unique path
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/favourites-api ./cmd/api

# --- Runtime stage ---
FROM alpine:3.20

WORKDIR /app

# Add a non-root user for safety
RUN adduser -D -g '' appuser
USER appuser

# Copy binary from builder
COPY --from=builder /bin/favourites-api /app/favourites-api

# Config env for your app
ENV HTTP_ADDR=:8080

EXPOSE 8080

ENTRYPOINT ["/app/favourites-api"]
    