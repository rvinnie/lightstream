DOCKER_COMPOSE_FILE = docker-compose.yaml

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down

ps:
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

proto:
	protoc \
	--go_out=. \
	--go_opt=paths=source_relative \
	--go-grpc_out=. \
	--go-grpc_opt=paths=source_relative \
	api/proto/v1/image_storage.proto

re: down up

.DEFAULT_GOAL := re
.PHONY: build, up, down, ps, re