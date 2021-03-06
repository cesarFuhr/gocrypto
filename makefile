# local development environments
SERVER_PORT=5000
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=dJ42YeQeneP3y8y3
DB_NAME=gocrypto
DB_DRIVER=postgres
DB_MAX_OPEN_CONNS=5
APP_KEYSOURCE_POOL_SIZE=10
APP_KEYSOURCE_RSAKEY_SIZE=2048

APP_ENV_STRING = SERVER_PORT=$(SERVER_PORT) \
	DB_HOST=$(DB_HOST) \
	DB_PORT=$(DB_PORT) \
	DB_USER=$(DB_USER) \
	DB_PASSWORD=$(DB_PASSWORD) \
	DB_NAME=$(DB_NAME) \
	DB_DRIVER=$(DB_DRIVER) \
	APP_KEYSOURCE_POOL_SIZE=$(APP_KEYSOURCE_POOL_SIZE) \
	APP_KEYSOURCE_RSAKEY_SIZE=$(APP_KEYSOURCE_RSAKEY_SIZE)

build:
	go build -o main ./cmd/main.go

install:
	go mod tidy
	go mod vendor

run: build
	./main

run-dev: build
	env $(APP_ENV_STRING) ./main

watch-dev: build
	env $(APP_ENV_STRING) air -c air.toml

start-local-db:
	docker run --detach --publish 127.0.0.1:$(DB_PORT):$(DB_PORT) \
		--env POSTGRES_USER=$(DB_USER) \
		--env POSTGRES_PASSWORD=$(DB_PASSWORD) \
		--env POSTGRES_DB=$(DB_NAME) \
		--name gocryptodb \
		postgres:alpine

stop-local-db:
	docker stop gocryptodb
	docker rm gocryptodb

test-unit:
	go test ./internal/...

test-full:
	docker-compose -f docker-compose.test.yml up -d db
	docker-compose -f docker-compose.test.yml up --build test
	docker-compose -f docker-compose.test.yml down

watch-test:
	watcher -cmd="make test-unit" -keepalive=true

watch-test-full:
	watcher -cmd="make test-full" -keepalive=true