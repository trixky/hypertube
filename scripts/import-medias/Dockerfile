FROM golang:1.18.3-alpine3.16

WORKDIR /app
COPY ./ /app

RUN go mod download

ENTRYPOINT go run main.go
