# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Install SQLite and any necessary tools
RUN apk add --no-cache gcc musl-dev libc-dev sqlite sqlite-dev

# Copy application files
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/challenge ./cmd/challenge

# Final stage
FROM alpine:3.18

# Copy the built application from the builder stage
COPY --from=builder /app/challenge /challenge

# Install SQLite3 in the final image
RUN apk add --no-cache sqlite

# Specify the entrypoint
ENTRYPOINT ["/challenge"]
