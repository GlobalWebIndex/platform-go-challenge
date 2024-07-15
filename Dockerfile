# Start with a base image that includes Go
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code into the container
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/app/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
