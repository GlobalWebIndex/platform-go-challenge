FROM golang:1.25.4-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

EXPOSE 8080

RUN go build -o server .

EXPOSE 8090

CMD ["./server"]