# Dockerfile

FROM registry.semaphoreci.com/golang:1.19 as builder

COPY . .

RUN mkdir -p /go/go
ENV GOPATH /go/go

RUN go mod download
RUN go mod verify
RUN go build -o gocaptcha-service

RUN mkdir -p /go/captchas/
RUN mkdir -p /go/logs/
RUN touch /go/logs/common.log

CMD ["./gocaptcha-service"]
