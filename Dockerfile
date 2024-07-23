FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM gcr.io/distroless/static

WORKDIR /root/

COPY --from=build /app/main .

EXPOSE 8080

CMD ["./main"]