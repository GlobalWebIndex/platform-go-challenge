FROM golang:1.25-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o server ./cmd/service

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server .
COPY --from=build /app/config.yaml .
CMD ["./server"]
