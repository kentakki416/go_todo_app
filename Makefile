.PHONY: help build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build:
	docker build -t buidougumi0617/gotodo:${DOCKER_TAG} \
		--target deploy ./

build-local:
	docker compose build --no-cache

up:
	docker compose up -d

down: 
	docker compose down

logs:	
	docker compose logs -f

ps:
	docker compose ps

test:
	go test -race -shuffle=on ./...

dry-migrate: ## Try migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo --dry-run < ./_tools/mysql/schema.sql

migrate:  ## Execute migration
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < ./_tools/mysql/schema.sql
