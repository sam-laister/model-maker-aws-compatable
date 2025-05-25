# Go Backend Rest Api for the project

[![Go Tests](https://github.com/2024-dissertation/go-api/actions/workflows/test.yml/badge.svg)](https://github.com/2024-dissertation/go-api/actions/workflows/test.yml)

## New Method

As of 17/01/25, migrations are now handled with [Goose](https://github.com/pressly/goose) and [Upper.io](https://upper.io/v4/adapter/postgresql/) Postgresql Adapter. 

To migrate up, reset or check status of the database, use the following commands:

 - `make db-status` to check the status of the database
 - `make up` to migrate up
 - `make reset` to reset the database

Other commands in the Makefile:

 - `make run` to run the server
 - `make build` to build the server  
 - `make seed` to seed the database

The .env in app/ is necessary for Makefile to work. Structure is as follows:

```env
DB_HOST=db
DB_NAME=appdb
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=changeme
DB_TIMEZONE=UTC

PORT=3333
MIGRATION_PATH=db/migrations
GEMINI_API_KEY=
```

This repo also contains a postres docker, but migrations are still handled through Goose and Go.

## Old Method

Database is handled with migrations using [Go Migrations CLI Tools](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md) and [Upper.io](https://upper.io/v4/adapter/postgresql/) Postgresql Adapter.

`brew install golang-migrate`

mirgate with `source .env && migrate -database ${DATABASE_URL} -path migrations up`

coverter https://developer.apple.com/augmented-reality/tools/

datasets: https://github.com/natowi/photogrammetry_datasets

## Running the server

The current routes are as follows:

 - POST   /verify
 - GET    /tasks
 - POST   /tasks
 - GET    /tasks/:taskID
 - POST   /tasks/:taskID/upload
 - POST   /tasks/:taskID/start
 - POST   /uploads
 - GET    /uploads/:taskId/:filename
 - GET    /objects/:taskID/:filename

 A postman collection exists to demo these routes.