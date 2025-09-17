FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

RUN go build -o platform-go-challenge main.go

EXPOSE 8080

CMD ["./platform-go-challenge"]