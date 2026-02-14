# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -o favourites-api \
    ./cmd/favourites-api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary, config, and OpenAPI spec from builder
COPY --from=builder /build/favourites-api .
COPY --from=builder /build/config ./config
COPY --from=builder /build/swagger.yaml .

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the server (ensure swagger.yaml is in cwd for Swagger UI)
CMD ["./favourites-api"]
