.PHONY: build build-debug up up-debug down logs dbshell start-db stop-db test

build:
	docker-compose build

build-debug:
	docker-compose build --build-arg DEBUG=true --build-arg GIN_MODE=debug

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f app

test:
	go test ./... -v
