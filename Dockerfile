# Use the latest official Golang image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the workspace
COPY . .

# Expose port 3000 for the Fiber app
EXPOSE 3000

# Command to run the Go app
CMD ["go", "run", "main.go"]