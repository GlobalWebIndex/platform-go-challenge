# ─────────────────────────────────────────────────────
# BUILD STAGE
# ─────────────────────────────────────────────────────
FROM golang:1.18.2-alpine AS builder

WORKDIR /favourites-api

# Install git (required for go modules sometimes)
RUN apk add --no-cache git

# Copy go module files first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Build the binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api .

# ─────────────────────────────────────────────────────
# RUNTIME STAGE
# ─────────────────────────────────────────────────────
FROM alpine:3.16

# Work directory
WORKDIR /favourites-api

# Copy compiled binary from builder stage
COPY --from=builder /favourites-api/server /favourites-api/server

# Expose port
EXPOSE 8080

# Run server
ENTRYPOINT ["/favourites-api/api"]
