ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: --setpath, migrateup, migratedown, update, run, docker-start, run-all

update:
	go get ./...
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

run-all: docker-start update migrateup run

run:
	@go run ./cmd/api

test:
	go test ./...

migrateup: --setpath
	migrate -path db/migrations -database $(dbpath) -verbose up

migratedown: --setpath
	migrate -path db/migrations -database $(dbpath) -verbose down

docker-start:
	docker compose build --no-cache
	docker compose up -d

--setpath:
	$(eval dbpath = $(DB_TYPE)://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable)