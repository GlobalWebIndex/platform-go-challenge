FROM golang:1.23.6-alpine3.21 AS build-img

RUN apk add --no-cache git

WORKDIR /tmp/Gwi

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o build/main ./cmd/main.go

FROM alpine:3.21

COPY --from=build-img /tmp/Gwi/build/main /app/main

EXPOSE 8080

CMD ["/app/main"]

