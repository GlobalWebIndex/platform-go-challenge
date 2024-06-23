# Builder
FROM golang:1.22.4-alpine3.20 as builder

RUN apk update && apk upgrade && \
    apk --update add git make bash

WORKDIR /home/go/server

RUN go install github.com/air-verse/air@latest

COPY . .

RUN go mod tidy
