# Build the Go binary
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files first for better caching
COPY go.mod go.sum ./

# Copy vendor directory and source code
COPY vendor/ ./vendor/

# Copy the rest of the source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/

# Create the final image
FROM alpine:latest

# Add ca-certificates for HTTPS requests and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user to run the application
RUN adduser -D -H -h /app appuser

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Make the binary executable
RUN chmod +x ./main

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose the port the application runs on
EXPOSE 8080

# Command to run the application
CMD ["./main"]