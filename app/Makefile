include .env

.DEFAULT_GOAL := build-and-run

BIN_FILE=main.out

build-and-run: build run

build:
	@go build -o "${BIN_FILE}"

clean:
	go clean
	rm --force "cp.out"
	rm --force nohup.out

run:
	./"${BIN_FILE}"

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

run:
	./"${BIN_FILE}"

doc:
	swag init

lint:
	golangci-lint run --enable-all

run-seed:
	@go run seeder/seeder.go

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) status

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) up

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) reset


run-test:
	- make test
	- make open-test