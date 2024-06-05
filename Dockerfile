FROM golang:1.22-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV JWT_SECRET=your_jwt_secret

RUN go test ./...

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
