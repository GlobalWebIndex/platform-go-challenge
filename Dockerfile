FROM golang:latest
WORKDIR /app
ADD . .
RUN go mod download
CMD cd cmd && go run main.go