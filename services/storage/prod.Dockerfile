FROM golang:1.19-alpine

WORKDIR /usr/src/app

ENV CGO_ENABLED=0

COPY ./ ./

RUN apk add --no-cache make && go mod download

ENTRYPOINT go build -o ./.bin/app ./cmd/main.go && ./.bin/app