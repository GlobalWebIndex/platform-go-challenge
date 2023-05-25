FROM golang:1.20.4-alpine

RUN apk update && apk add git build-base

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/go-delve/delve/cmd/dlv@v1.20.2

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -gcflags='all=-N -l' -o main ./cmd/http/main.go

EXPOSE 8080
EXPOSE 40000

CMD ["dlv", "--listen=:40000", "--headless=true", "--accept-multiclient", "--continue=true", "--api-version=2", "exec", "./main"]
