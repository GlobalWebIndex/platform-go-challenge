FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Go module download
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
ENV GOTOOLCHAIN=auto
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/favourites-api ./cmd/favourites-api

FROM gcr.io/distroless/base-debian12:nonroot
WORKDIR /app

# Copy binary and public key (verification only)
COPY --from=builder /app/favourites-api /app/favourites-api
COPY public.pem /app/public.pem
COPY configs/prod.yaml /app/configs/prod.yaml

ENV PUBLIC_KEY_PATH=/app/public.pem
ENV APP_ENV=prod
EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["/app/favourites-api"]
