# ---------- Build stage ----------
FROM golang:1.22-alpine AS build
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# ---------- Runtime stage ----------
FROM alpine:3.20
WORKDIR /app

COPY --from=build /app/server /app/server

EXPOSE 8080
ENTRYPOINT ["/app/server"]
