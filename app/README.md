# AWS Compatable fork of Model Maker Golang's API

#### Setup

Copy `.env.example` in `/app` to `.env` and fill with real values. The project was originally setup for Katapult S3 compatable buckets so the env names reflect this. This will be altered soon

To spin up the docker, run the dev environment with:

`docker compose -f docker-compose-dev.yml up -d && docker compose -f docker-compose-dev.yml exec api sh`

Once in the container run the api with `make`.

For production, just spin up the docker with `docker compose up -d`.