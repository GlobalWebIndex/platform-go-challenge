# Builder stage
FROM golang:alpine3.14 AS builder

# Install dependencies
RUN apk update && apk add --no-cache \
  git \
  ca-certificates \
  && update-ca-certificates

# Add source files and set the proper workdir |
COPY . $GOPATH/src/github.com/josedelri85/platform-go-challenge/
WORKDIR $GOPATH/src/github.com/josedelri85/platform-go-challenge/

# Enable GO Modules
ENV GO111MODULES=on
# Build the binary
RUN go build -mod=vendor -o /go/bin/platform-go-challenge .

# Final image
FROM alpine

# Copy our static executable
COPY --from=builder /go/bin/platform-go-challenge /go/bin/platform-go-challenge

# Copy the ca-certificates to be able to perform https requests
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs

RUN apk add --no-cache tzdata && apk add curl

# Run the binary
ENTRYPOINT ["/go/bin/platform-go-challenge"]