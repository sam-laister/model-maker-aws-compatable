#!/bin/bash
# This script builds and runs the prod image for the modelmaker service.
docker build -f Dockerfile.prod -t modelmaker:prod --platform=linux/amd64 .