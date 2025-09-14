FROM golang:1.23.8

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module definition files from your project into the container
COPY go.mod go.sum ./
# Download all the dependencies listed in go.mod and go.sum
RUN go mod download

# Add rest of the code to the /app directory inside the container
COPY . ./

# Run all unittests!
# The image can not be built if at least one test is failing
RUN go test ./...

# Compile the Go app.
RUN go build -o platform-go-challenge ./cmd/main.go

# Tell Docker to execute your compiled Go binary.
CMD ["./platform-go-challenge"]