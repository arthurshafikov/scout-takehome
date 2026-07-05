BIN := "./.bin/app"

include main.env

DOCKER_COMPOSE=docker compose

ifeq ($(APP_ENV), production)
	DOCKER_COMPOSE_FILE := "deployments/docker-compose.prod.yml"
else
	DOCKER_COMPOSE_FILE := "deployments/docker-compose.yml"
endif
DOCKER_COMPOSE_TEST_FILE := "./deployments/docker-compose.tests.yml"

APP_NAME := "boilerplate"
APP_TEST_NAME := "boilerplate_test"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

POSTGRESQL_DSN="host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=${DB_NAME} port=${DB_PORT} sslmode=${DB_SSL_MODE} TimeZone=${DB_TIMEZONE}"

build:
	go build -a -o $(BIN) -ldflags "$(LDFLAGS)" cmd/app/main.go

run: build 
	 $(BIN) -cfgFolder ./configs -env ./

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

integration-tests:
	$(DOCKER_COMPOSE) -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} up pgsql_test --wait
	goose -dir migrations postgres "host=localhost user=test password=test dbname=test port=54400 sslmode=prefer TimeZone=UTC" up
	$(DOCKER_COMPOSE) -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} up app --attach app --abort-on-container-exit --exit-code-from app
	$(DOCKER_COMPOSE) -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes

reset-integration-tests:
	$(DOCKER_COMPOSE) -f ${DOCKER_COMPOSE_TEST_FILE} -p ${APP_TEST_NAME} down --volumes

refresh: reset 
	make migrate

# DB commands for docker
migrate-in-docker:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} exec --workdir /go/src app goose -dir migrations postgres ${POSTGRESQL_DSN} up

rollback-in-docker:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} exec --workdir /go/src app goose -dir migrations postgres ${POSTGRESQL_DSN} down

reset-in-docker:
	$(DOCKER_COMPOSE) --env-file ./main.env -f ${DOCKER_COMPOSE_FILE} -p ${APP_NAME} exec --workdir /go/src app goose -dir migrations postgres ${POSTGRESQL_DSN} reset

# DB commands for working without docker
migrate:
	goose -dir migrations postgres ${POSTGRESQL_DSN} up
	
migrate-tests:
	goose -dir migrations postgres "host=${DB_HOST} user=${DB_USER} password=${DB_PASSWORD} dbname=postgres port=${DB_PORT} sslmode=${DB_SSL_MODE} TimeZone=${DB_TIMEZONE}" up

migrate-one:
	goose -dir migrations postgres ${POSTGRESQL_DSN} up-by-one

rollback:
	goose -dir migrations postgres ${POSTGRESQL_DSN} down

reset:
	goose -dir migrations postgres ${POSTGRESQL_DSN} reset

newmig:
	goose -dir migrations create ${NAME} sql
