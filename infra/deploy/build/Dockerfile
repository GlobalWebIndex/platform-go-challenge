FROM golang:1.18 as builder
WORKDIR /home/platform-go-challenge
COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-X main.version=$version" -o platform-go-challenge ./cmd/platform-go-challenge/main.go

FROM alpine
COPY --from=builder /home/platform-go-challenge/platform-go-challenge .
COPY --from=builder /home/platform-go-challenge/script/run.sh .

ENTRYPOINT ["./run.sh"]
