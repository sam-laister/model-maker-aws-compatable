.DEFAULT_GOAL := build-and-run

GO_TEST_FLAGS ?= -race -count=1 -v -timeout=5m -json

BIN_FILE=main.out

build-and-run: build run

DSN=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable&timezone=$(DB_TIMEZONE)

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${BIN_FILE}"

dsn:
	@echo "DSN: ${DSN}"

clean:
	go clean
	rm --force "cp.out"
	rm --force nohup.out

run:
	./"${BIN_FILE}"

test:
	@godotenv -f .env.test go test ./... -cover fmt

test-cover:
	echo "Running tests with coverage..."
	@godotenv -f .env.test go test $(GO_TEST_FLAGS) -coverprofile=cp.out -coverpkg=$$(go list ./... | grep -v -e /bin -e /cmd -e /examples) ./...
	@go tool cover -html=cp.out -o cover.html
	@go-cover-treemap -coverprofile=cp.out > out.svg

test-packages:
	go test $(GO_TEST_FLAGS) $$(go list ./... | grep -v -e /bin -e /cmd -e /examples) |\
		tparse --follow -sort=elapsed -trimpath=auto -all

doc:
	swag init

lint:
	golangci-lint run --enable-all

run-seed:
	@godotenv -f .env  go run seeds/bin/seeder.go

db-status:
	@GOOSE_DBSTRING="${DSN}" goose -dir=$(MIGRATION_PATH) status

up:
	@GOOSE_DBSTRING="${DSN}" goose -dir=$(MIGRATION_PATH) up

reset:
	@GOOSE_DBSTRING="${DSN}" goose -dir=$(MIGRATION_PATH) reset
