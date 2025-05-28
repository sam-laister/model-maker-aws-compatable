#!/bin/bash
# This script runs a Docker container for the ModelMaker API.
docker run -d \
    --platform=linux/amd64 \
    --name api-modelmaker \
    --env-file .env \
    -p 3333:3333 \
    -e PORT=3333 \
    -e LOG_LEVEL=info \
    modelmaker:prod