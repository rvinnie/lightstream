PATH_TO_PROTO = "../../api/proto/v1"

build:
	go mod download && go build -o ./.bin/app ./cmd/main.go

run: build
	./.bin/app

proto:
	protoc \
	--go_out=pb \
	--go_opt=paths=source_relative \
	--go-grpc_out=pb \
	--go-grpc_opt=paths=source_relative \
	--proto_path=$(PATH_TO_PROTO) $(PATH_TO_PROTO)/*.proto

.DEFAULT_GOAL := run
.PHONY: build, run