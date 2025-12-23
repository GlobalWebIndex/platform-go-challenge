# Multi-stage Dockerfile for Go application
# Stage 1: Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy just the Go file (no go.mod/go.sum yet)
COPY go_impl.go .

# Initialize module and download dependencies
RUN go mod init github.com/gwi-challenge && \
    go get github.com/google/uuid && \
    go get github.com/gorilla/mux && \
    go get github.com/lib/pq


# RUN TESTS (fails build if tests fail)
RUN go test -v

# Build the application with optimizations
# CGO_ENABLED=0: static binary (no libc dependency)
# GOOS=linux: target Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server go_impl.go

# Stage 2: Runtime stage
# Use scratch (empty image) for smallest final size
FROM alpine:latest

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/server .

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./server"]