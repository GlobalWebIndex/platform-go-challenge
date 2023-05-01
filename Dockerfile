# syntax=docker/dockerfile:1

FROM golang:1.18

# Set destination for COPY
WORKDIR /app

# Copy source code
COPY . .

# Download Go modules
RUN go mod download

# Build
#RUN CGO_ENABLED=0 GOOS=linux go build -o platform-go-challenge
RUN go build .

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080

# Run
CMD ["./platform-go-challenge"]