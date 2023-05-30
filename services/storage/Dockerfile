FROM golang:1.19-alpine

WORKDIR /usr/src/app

RUN apk update && apk add make

COPY ./ ./

CMD ["make"]
