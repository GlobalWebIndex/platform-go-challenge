# --- Stage 1: Builder ---
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy dependency files first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build static binary
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server ./cmd/server

# --- Stage 2: Runner ---
FROM scratch
# If your app makes HTTPS requests, you need certs:
# FROM alpine:latest
# RUN apk add --no-cache ca-certificates

COPY --from=builder /server /server

EXPOSE 8080
ENTRYPOINT ["/server"]
