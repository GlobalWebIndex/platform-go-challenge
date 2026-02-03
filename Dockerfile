# Stage 1: Builder
FROM golang:1.25 AS builder

WORKDIR /app

# Download dependencies first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build statically linked binary
# CGO_ENABLED=0 for static binary, needed for distroless
RUN CGO_ENABLED=0 GOOS=linux go build -o favorites-app ./cmd/server

# Stage 2: Runtime
# using distroless static image for security and size
FROM gcr.io/distroless/static-debian12

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/favorites-app /favorites-app

# Expose port (adjust if your app uses a different one)
EXPOSE 8080

# Run as non-root (distroless default user is non-root usually, but making sure)
USER nonroot:nonroot

ENTRYPOINT ["/favorites-app"]
