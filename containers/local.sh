#!/bin/bash

# Print usage and exit
help() {
    echo "Usage $DIR/local.sh [clean|build|up]"
    exit 2
}

# Ensure only zero or one arguments are passed to the script
if [ $# -gt 1 ]; then
    help
fi

# Set some helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Set environment for docker compose
export GIT_REVISION=$(git rev-parse --short HEAD)
export REACT_APP_GIT_REVISION=$GIT_REVISION

# Check if the build or clean arguments are specified
if [[ $# -eq 1 ]]; then
    if [[ $1 == "clean" ]]; then
        docker system prune --all
        exit 0
    elif [[ $1 == "build" ]]; then
        docker compose -p gds -f $DIR/docker-compose.yaml --profile=all build
        exit 0
    elif [[ $1 == "up" ]]; then
        echo "starting docker server"
    else
        help
    fi
fi

# By default just bring docker compose up
docker compose -p gds -f $DIR/docker-compose.yaml --profile=all up