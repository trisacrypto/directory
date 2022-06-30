#!/bin/bash
# A helper for common docker compose operations to run GDS locally.

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-p PROFILE] [clean|build|up]
A helper for common docker compose operations to run GDS
services locally. Flags are as follows (getopt required):

    -h          display this help and exit
    -p PROFILE  specify the docker compose profile to use

There are two ways to use this script. Run docker compose:

    ${0##*/} up

Build the images cleaning the docker cache first:

    ${0##*/} clean
    ${0##*/} build

You can also specify the profile to run only some services.
EOF
}

# Parse command line options with getopt
PROFILE="all"
OPTIND=1

while getopts hp: opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        p)  PROFILE=$OPTARG
            ;;
        *)
            show_help >&2
            exit 2
            ;;
    esac
done
shift "$((OPTIND-1))"

# Ensure only zero or one arguments are passed to the script
if [ $# -gt 1 ]; then
    show_help >&2
    exit 2
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
        docker compose -p gds -f $DIR/docker-compose.yaml --profile=$PROFILE build
        exit 0
    elif [[ $1 == "up" ]]; then
        echo "starting docker compose services"
    else
        show_help >&2
        exit 2
    fi
fi

# By default just bring docker compose up
docker compose -p gds -f $DIR/docker-compose.yaml --profile=$PROFILE up