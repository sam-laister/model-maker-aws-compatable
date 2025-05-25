#!/bin/sh
set -e  # Exit if any command fails
echo "Starting application..."

make 

exec "$@"
