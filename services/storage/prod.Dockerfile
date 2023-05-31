FROM golang:1.19-alpine

WORKDIR /usr/src/app

ENV CGO_ENABLED=0

RUN apk update && apk add make
RUN go mod download

COPY ./ ./

ENTRYPOINT go build -o ./.bin/app ./cmd/main.go && ./.bin/app