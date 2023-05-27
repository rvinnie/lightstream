build:
	go mod download && go build -o ./.bin/app ./cmd/main.go

run: build
	./.bin/app

proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	pb/image_storage.proto

.DEFAULT_GOAL := run
.PHONY: build, run