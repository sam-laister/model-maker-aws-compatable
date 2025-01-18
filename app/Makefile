include .env

run: 
	@go run main.go

build:
	@go build -o bin/app .

run-seed:
	@go run seeder/seeder.go

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) status

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) up

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DATABASE_URL) goose -dir=$(MIGRATION_PATH) reset

test:
	@go test ./... -coverprofile cover.out

open-test:
	@go tool cover -html=cover.out

run-test:
	- make test
	- make open-test