# Start from golang base image
FROM golang:1.19.3-alpine3.16 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git build-base bash

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY ../go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and the go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the working Directory inside the container
COPY .. .

# Build the Go app
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/http/main.go

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage. Observe we also copied the .env file
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

# Expose port 8080 to the outside world
EXPOSE 8080

#Command to run the executable
RUN chmod +x ./main
CMD ["./main"]