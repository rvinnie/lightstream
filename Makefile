build:
	go mod download && go build -o ./.bin/app ./services/streaming/cmd/main.go

run: build
	./.bin/app

.DEFAULT_GOAL := run
.PHONY: build, run