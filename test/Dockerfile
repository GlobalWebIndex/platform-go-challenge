FROM golang:1.20.5-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY go.mod ./

# Download all dependencies.
RUN go mod download

# Copy the source code into the container
COPY main.go ./

RUN go build -o main .

# Command to run the Go application
CMD ["./main"]