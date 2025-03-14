.PHONY: docker-run docker-stop

.DEFAULT_GOAL := help

DOCKER_COMPOSE = docker-compose -f docker/docker-compose.yml
GO = go
APP_NAME = bankapi
VERSION ?= latest

### run app in docker
docker-run: 
	$(DOCKER_COMPOSE) up -d 

### stop all docker containers
docker-stop:
	$(DOCKER_COMPOSE) down -v
