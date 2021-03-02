build:
	go build -o main ./cmd/main.go

install:
	go mod tidy
	go mod vendor

run: build
	./main

run-dev:
	air -c air.toml

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