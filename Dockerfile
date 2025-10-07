# ---------- Stage 1: Build ----------
FROM golang:1.25-alpine AS builder

# # Set environment variables
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# Set work directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all source code
COPY . .

# Build the binary
RUN go build -o main .

# ---------- Stage 2: Run ----------
FROM alpine:latest

# Install certificates (needed for HTTPS requests)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the compiled binary from builder
COPY --from=builder /app/main .

# Expose port (same as your app)
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
