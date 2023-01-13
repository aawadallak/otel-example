FROM golang:alpine

WORKDIR /usr/app

ADD go.mod go.sum

RUN go mod tidy

RUN go build -o app .

