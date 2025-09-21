# syntax=docker/dockerfile:1

# Podman uses the same Dockerfile syntax. This Containerfile mirrors the Dockerfile.

# ----- Build stage -----
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /bin/favourites .

# ----- Runtime stage -----
FROM scratch
EXPOSE 8080
WORKDIR /app
COPY --from=builder /bin/favourites /favourites
ENTRYPOINT ["/favourites"]
