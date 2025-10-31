FROM golang:latest AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api/main

RUN chmod +x main

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"]