# Model Maker Backend API

<img align="right" width="125" src="assets/app-icon.png">

[![Model Maker CI](https://github.com/2024-dissertation/model-maker-docker/actions/workflows/ci.yml/badge.svg)](https://github.com/2024-dissertation/model-maker-docker/actions/workflows/ci.yml)

Model Maker is an accessible and Open Source solution to Photogrammetry. This is the backend for the also open source [Flutter App](https://github.com/2024-dissertation/model-maker-app).

Included is a Go REST API with a Dockerfile environment ready for deployment.

#### Features

- Works against multiple databases:
  - Postgres, MySQL, and more.
- Cross platform.
- Development and Production ready Docker environments.
  - `docker-compose-test.yml` will setup a Postgres container and deploy an Ubuntu based enironment for the server executable.
  - `docker-compose-dev.yml` will create a continous environment for internal compiling with `make`
- Uses [OpenMVG](https://github.com/openMVG/openMVG) and [OpenMVS](https://github.com/cdcseacave/openMVS).
- Makefile for handling deploying, migrations, running tests.
- Firebase for authentication.
- Environment variables integrations.
- ... and more!

#### Setup

Before setting up, a [Google Firebase](https://firebase.google.com/) project must be setup. This project is designed to use the free plan for authentication, so all that'd needed is the `service-account-key.json`.

Environment variables must be configured, check [/app/.env.example](/app/.env.example) for required keys. It's important to provide the raw json as a string in `GOOGLE_CREDENTIALS`. If running the development environment, the json file should be placed in [/app](/app) and the path updated accordingly. If using the production environment, only the string value needs to be configured.

To run the development environment:

```bash
docker compose -f docker-compose-dev.yml up -d
```

and the production environment:

```bash
docker compose -f docker-compose-test.yml up -d
```

**(Note, build times take up to 40 mins on M1 Pro Macbook due to the size of the Photogrammetry tools)**

To rebuild the application in development, the Docker container will stay open without executing the binary. Instead exec in with `docker compose -f docker-compose-dev.yml exec modelmaker bash` and run `make`.

In both scenarios, the API will run on port `3333` unless changed in the `.env`.

Example scripts are provided in [/app](/app) called `run.sh` and `build.sh` for deploying and building the docker images.

#### Migrations

Migrations are handled automatically. Incase manual control is needed, install the tool [Goose](https://github.com/pressly/goose).

To migrate up, reset or check status of the database, use the following commands:

- `make db-status` to check the status of the database
- `make up` to migrate up
- `make reset` to reset the database

Other commands in the Makefile:

- `make run` to run the server
- `make build` to build the server
- `make seed` to seed the database

#### Documentation

Documentation is currently in development for the API, as well as a Postman collection. In the meantime, refer to the frontend App.

#### Special Thanks

- **Datasets index**: https://github.com/natowi/photogrammetry_datasets
- **Testing image dataset**: https://www.youtube.com/watch?v=4LexaqdxdiU
