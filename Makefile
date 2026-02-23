.PHONY: build run test lint migrate-up migrate-down migrate-create

build:
	go build -o bin/api ./cmd/api

run:
	go run ./cmd/api

test:
	go test ./... -v -count=1

lint:
	golangci-lint run ./...

migrate-up:
	goose -dir db/migrations postgres "$(DATABASE_URL)" up

migrate-down:
	goose -dir db/migrations postgres "$(DATABASE_URL)" down

migrate-create:
	goose -dir db/migrations create $(name) sql
