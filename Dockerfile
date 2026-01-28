# ---------- build stage ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for go modules sometimes)
RUN apk add --no-cache git

# Copy go mod files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api

# ---------- runtime stage ----------
FROM alpine:3.19

WORKDIR /app

# Certificates for HTTPS / DB connections
RUN apk add --no-cache ca-certificates

# Copy binary from builder
COPY --from=builder /app/api /app/api

# Expose app port
EXPOSE 8080

# Run the binary
CMD ["/app/api"]
