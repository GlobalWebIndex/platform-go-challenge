FROM golang:1.25

WORKDIR /app

COPY . .

RUN go mod tidy || true

CMD ["/bin/bash"]
