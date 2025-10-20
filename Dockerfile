# Use the official Golang image
FROM golang:latest

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build statically for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -o favourites-app ./cmd/app/main.go

# Expose the server port
EXPOSE 3002

# Run the application
CMD ["./favourites-app"]
