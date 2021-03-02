build:
	go build -o main ./cmd/main.go

install:
	go mod tidy
	go mod vendor

run: build
	./main

run-dev:
	env SERVER_PORT=5000 \
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=postgres \
	DB_PASSWORD=dJ42YeQeneP3y8y3 \
	DB_NAME=gocrypto \
	DB_DRIVER=postgres \
	APP_KEYSOURCE_POOL_SIZE=10 \
	APP_KEYSOURCE_RSAKEY_SIZE=2048 \
	air -c air.toml

start-local-db:
	docker run --detach --publish 5432:5432 \
		--env POSTGRES_USER=postgres \
		--env POSTGRES_PASSWORD=dJ42YeQeneP3y8y3 \
		--env POSTGRES_DB=gocrypto \
		--name gocryptopg \
		postgres:alpine

stop-local-db:
	docker stop gocryptopg
	docker rm gocryptopg

test-unit:
	go test ./internal/...

test-full:
	docker-compose -f docker-compose.test.yml up -d db
	docker-compose -f docker-compose.test.yml up --build test
	docker-compose -f docker-compose.test.yml down

test-unit-watch:
	watcher -cmd="make unit-test" -keepalive=true

test-full-watch:
	watcher -cmd="make full-test" -keepalive=true