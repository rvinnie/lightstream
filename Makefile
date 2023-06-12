ifeq ($(version), production)
        DOCKER_COMPOSE_FILE = docker-compose.production.yaml
else
        DOCKER_COMPOSE_FILE = docker-compose.yaml
endif

DOCKER_COMPOSE_FILE_STRESS = docker-compose.stress.yaml

up:
	docker-compose -f ${DOCKER_COMPOSE_FILE} up --build

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} down

ps:
	docker-compose -f ${DOCKER_COMPOSE_FILE} ps

re: down up

stress:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -f ${DOCKER_COMPOSE_FILE_STRESS} up --build

stress_down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -f ${DOCKER_COMPOSE_FILE_STRESS} down

stress_ps:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -f ${DOCKER_COMPOSE_FILE_STRESS} ps

stress_re: stress_down stress

.DEFAULT_GOAL := re
.PHONY: up down ps re stress stress_down stress_ps stress_re