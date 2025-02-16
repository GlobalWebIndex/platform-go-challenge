# Stage 1: Build the Go binary
FROM golang:1.24 AS builder

WORKDIR /app

# Copy Go module files first
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application with static linking for Alpine compatibility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gwi-app main.go

# Debugging: Check if binary was created
RUN ls -l /app

# Stage 2: Minimal production image
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS support
RUN apk --no-cache add ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gwi-app .

# Copy the generated SSL certificates
COPY cert.pem key.pem ./

# Set execute permissions for the Go binary
RUN chmod +x gwi-app

# Expose the HTTPS port
EXPOSE 8080

# Run the application with HTTPS
ENTRYPOINT ["./gwi-app"]
