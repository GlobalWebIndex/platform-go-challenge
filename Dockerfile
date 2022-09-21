FROM golang:1.19-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY . ./
RUN go build -o /gwiserver


FROM alpine
WORKDIR /
COPY --from=build /gwiserver /gwiserver
CMD ["./gwiserver"]