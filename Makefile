ifeq ($(version), prod)
	DOCKER_COMPOSE_FILE = -f docker-compose.production.yaml
else ifeq ($(version), stress)
	DOCKER_COMPOSE_FILE = -f docker-compose.yaml -f docker-compose.stress.yaml
else
	DOCKER_COMPOSE_FILE = -f docker-compose.yaml
endif

up:
	docker-compose ${DOCKER_COMPOSE_FILE} up --build

down:
	docker-compose ${DOCKER_COMPOSE_FILE} down

ps:
	docker-compose ${DOCKER_COMPOSE_FILE} ps

re: down up

.DEFAULT_GOAL := re
.PHONY: up down ps re