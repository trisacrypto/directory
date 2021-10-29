#!/bin/bash

ask() {
    local prompt default reply

    if [[ ${2:-} = 'Y' ]]; then
        prompt='Y/n'
        default='Y'
    elif [[ ${2:-} = 'N' ]]; then
        prompt='y/N'
        default='N'
    else
        prompt='y/n'
        default=''
    fi

    while true; do

        # Ask the question (not using "read -p" as it uses stderr not stdout)
        echo -n "$1 [$prompt] "

        # Read the answer (use /dev/tty in case stdin is redirected from somewhere else)
        read -r reply </dev/tty

        # Default?
        if [[ -z $reply ]]; then
            reply=$default
        fi

        # Check if the reply is valid
        case "$reply" in
            Y*|y*) return 0 ;;
            N*|n*) return 1 ;;
        esac

    done
}


if [ -z "$1" ]; then
    TAG=$(git rev-parse --short HEAD)
else
    TAG=$1
fi


if ! ask "Continue with tag $TAG?" N; then
    exit 1
fi

if [ -z "$REACT_APP_VASPDIRECTORY_CLIENT_ID" ]; then
    echo "REACT_APP_VASPDIRECTORY_CLIENT_ID environment variable required"
    exit 1
fi

if [ -z "$REACT_APP_VASPDIRECTORY_ANALYTICS_ID" ]; then
    echo "REACT_APP_VASPDIRECTORY_ANALYTICS_ID environment variable required"
    exit 1
fi

if [ -z "$REACT_APP_TRISATEST_CLIENT_ID" ]; then
    echo "REACT_APP_TRISATEST_CLIENT_ID environment variable required"
    exit 1
fi

if [ -z "$REACT_APP_TRISATEST_ANALYTICS_ID" ]; then
    echo "REACT_APP_TRISATEST_ANALYTICS_ID environment variable required"
    exit 1
fi



DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")


# Build the primary backend images
docker build -t trisa/gds:$TAG -f $DIR/gds/Dockerfile $REPO
docker build -t trisa/grpc-proxy:$TAG -f $DIR/grpc-proxy/Dockerfile $REPO

# Build the UI images for trisatest.net and vaspdirectory.net
docker build \
    -t trisa/gds-ui:$TAG -f $DIR/gds-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://proxy.vaspdirectory.net \
    --build-arg REACT_APP_GDS_IS_TESTNET=false \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_VASPDIRECTORY_ANALYTICS_ID} \
    $REPO

docker build \
    -t trisa/gds-testnet-ui:$TAG -f $DIR/gds-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://proxy.trisatest.net \
    --build-arg REACT_APP_GDS_IS_TESTNET=true \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_TRISATEST_ANALYTICS_ID} \
    $REPO

# Build the Admin UI images for admin.trisatest.net and admin.vaspdirectory.net
docker build \
    -t trisa/gds-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.vaspdirectory.net/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=false \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_VASPDIRECTORY_CLIENT_ID} \
    $REPO

docker build \
    -t trisa/gds-testnet-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.trisatest.net/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=true \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_TRISATEST_CLIENT_ID} \
    $REPO

# Retag the images to push to gcr.io
docker tag trisa/gds:$TAG gcr.io/trisa-gds/gds:$TAG
docker tag trisa/gds-ui:$TAG gcr.io/trisa-gds/gds-ui:$TAG
docker tag trisa/grpc-proxy:$TAG gcr.io/trisa-gds/grpc-proxy:$TAG
docker tag trisa/gds-testnet-ui:$TAG gcr.io/trisa-gds/gds-testnet-ui:$TAG
docker tag trisa/gds-admin-ui:$TAG gcr.io/trisa-gds/gds-admin-ui:$TAG
docker tag trisa/gds-testnet-admin-ui:$TAG gcr.io/trisa-gds/gds-testnet-admin-ui:$TAG

# Push to DockerHub
docker push trisa/gds:$TAG
docker push trisa/grpc-proxy:$TAG
docker push trisa/gds-ui:$TAG
docker push trisa/gds-testnet-ui:$TAG
docker push trisa/gds-admin-ui:$TAG
docker push trisa/gds-testnet-admin-ui:$TAG

# Push to GCR
docker push gcr.io/trisa-gds/gds:$TAG
docker push gcr.io/trisa-gds/grpc-proxy:$TAG
docker push gcr.io/trisa-gds/gds-ui:$TAG
docker push gcr.io/trisa-gds/gds-testnet-ui:$TAG
docker push gcr.io/trisa-gds/gds-admin-ui:$TAG
docker push gcr.io/trisa-gds/gds-testnet-admin-ui:$TAG
