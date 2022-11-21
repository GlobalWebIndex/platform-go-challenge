FROM golang:1.18
ENV GO111MODULE=on
WORKDIR /go/src/github.com/ghal/platform-go-challenge
COPY . ./
RUN chmod +x ./script/ci.sh
RUN go install golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616
