# Go Backend Rest Api for the project

Database is handled with migrations using [Go Migrations CLI Tools](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md) and [Upper.io](https://upper.io/v4/adapter/postgresql/) Postgresql Adapter.

`brew install golang-migrate`

mirgate with `source .env && migrate -database ${DATABASE_URL} -path migrations up`

coverter https://developer.apple.com/augmented-reality/tools/

datasets: https://github.com/natowi/photogrammetry_datasets