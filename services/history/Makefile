include .env

HOST=localhost
PORT=5433

build:
	go mod download && go build -o ./.bin/app ./cmd/main.go

run: build
	./.bin/app

select:
	psql postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${HOST}:${PORT}/${POSTGRES_DB} -c 'SELECT * FROM notifications'

.DEFAULT_GOAL := run
.PHONY: build, run, select