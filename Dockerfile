FROM golang:1.18-rc-alpine3.14

RUN mkdir /platform-go-challenge
COPY . /platform-go-challenge
WORKDIR /platform-go-challenge
RUN go build cmd/platform-go-challenge/main.go
CMD ["/platform-go-challenge/main"]
