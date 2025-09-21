# syntax=docker/dockerfile:1

# ----- Build stage -----
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Copy the entire project (small repo, no external deps)
COPY . .

# Build a static binary for Linux
ENV CGO_ENABLED=0
RUN go build -ldflags="-s -w" -o /bin/favourites .

# ----- Runtime stage -----
FROM scratch AS runtime

# App runs on 8080
EXPOSE 8080

# Working directory for relative paths like ./data/favourites.json
WORKDIR /app

# Copy binary
COPY --from=builder /bin/favourites /favourites

# Default to in-memory store. To enable file store, set env at runtime:
#   -e FAV_STORE=file -e FAV_PATH=/data/favourites.json -v $(pwd)/data:/data
ENTRYPOINT ["/favourites"]
