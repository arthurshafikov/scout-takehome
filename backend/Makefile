BIN := "./.bin/app"
SEED := "./.bin/seed"

include main.env

DOCKER_COMPOSE=docker compose

ifeq ($(APP_ENV), production)
	DOCKER_COMPOSE_FILE := "deployments/docker-compose.prod.yml"
else
	DOCKER_COMPOSE_FILE := "deployments/docker-compose.yml"
endif

APP_NAME := "scout"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/app/main.go

build-seed:
	go build -o $(SEED) cmd/seed/main.go

run: build 
	 $(BIN) -cfgFolder ./configs -env ./

seed: build-seed
	$(SEED) -endpoint $(MINIO_ENDPOINT) -access-key $(MINIO_ACCESS_KEY) -secret-key $(MINIO_SECRET_KEY) -bucket $(MINIO_BUCKET)

r:
	go run ./... -cfgFolder ./configs -env ./

test: 
	go test --short ./internal/...

test-race: 
	go test --short -race ./internal/...

lint: 
	golangci-lint run ./...

.PHONY: build test

ps:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} ps

go:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up --build -d

rebuild:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up app --force-recreate --no-deps --build -d

up:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} up -d

stop:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} stop

restart:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} restart

down:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} down --volumes

mocks:
	mockgen -source=./internal/repository/repository.go -destination ./internal/repository/mocks/mock.go
	mockgen -source=./internal/services/services.go -destination ./internal/services/mocks/mock.go
	mockgen -source=./internal/transport/http/server.go -destination ./internal/transport/http/mocks/mock.go

.PHONY: build test run lint mocks ps go up stop down restart
