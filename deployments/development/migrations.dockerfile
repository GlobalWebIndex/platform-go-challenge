FROM golang:1.20.4-alpine

ENV ACTION=""

WORKDIR /app

COPY deployments/development/.migrations migrations

RUN apk update

RUN go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

ENTRYPOINT [ "sh", "-c", "migrate -source file://./migrations -database mongodb://mongodb:27017/gwi $ACTION" ]
