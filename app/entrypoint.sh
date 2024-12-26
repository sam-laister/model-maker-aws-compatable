#!/bin/sh

# install dependencies
apk add --no-cache curl

curl -L https://github.com/golang-migrate/migrate/releases/download/$MIGRATE_VERSION/migrate.linux-arm64.tar.gz | tar xvz -C /usr/local/bin

# run migrations as required
# migrate -path=./migrations -database="$VET_DB_DSN" up

# run the command passed as arguments
exec "$@"