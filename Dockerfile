# Start with a base image containing Go
FROM golang:1.23

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN go build -o main .

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
