FROM golang:1.19-alpine

WORKDIR /usr/src/app

ENV CGO_ENABLED=0

COPY ./ ./

RUN apk update && apk add make
RUN go mod download
RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon

ENTRYPOINT CompileDaemon -build="go build -o ./.bin/app ./cmd/main.go" -command="./.bin/app"