DOCKER_COMPOSE_FILE = docker-compose.yaml

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down

ps:
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

re: down up

.DEFAULT_GOAL := re
.PHONY: build, up, down, ps, re