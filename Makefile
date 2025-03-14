.PHONY: setup docker-run docker-stop test 

.DEFAULT_GOAL := help

DOCKER_COMPOSE = docker-compose -f docker/docker-compose.yml
GO = go
APP_NAME = bankapi
VERSION ?= latest

### setup for run
setup:
	cp .env.example .env

### run app in docker
docker-run: 
	$(DOCKER_COMPOSE) up -d 

### stop all docker containers
docker-stop:
	$(DOCKER_COMPOSE) down -v

### run tests
test:
	$(GO) test ./...
