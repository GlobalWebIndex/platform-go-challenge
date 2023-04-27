# Use the official Golang image as the parent image
FROM golang:1.20.3-alpine3.17 AS build

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . .

# download Go modules and dependencies
RUN go mod download

# Build the Go application
RUN go build .

# Start a new stage from scratch
FROM alpine:3.14 AS runtime

# Copy the executable from the previous stage
COPY --from=build /app/challenge /usr/local/bin/challenge
COPY config.yaml .
# Expose port 8000 for the application
EXPOSE 8000

# Run the application
CMD ["challenge"]
