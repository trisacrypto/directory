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

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")


docker build -t trisa/gds:$TAG -f $DIR/gds/Dockerfile $REPO
docker build -t trisa/gds-ui:$TAG -f $DIR/gds-ui/Dockerfile $REPO
docker build -t trisa/gds-testnet-ui:$TAG -f $DIR/gds-testnet-ui/Dockerfile $REPO
docker build -t trisa/grpc-proxy:$TAG -f $DIR/grpc-proxy/Dockerfile $REPO

docker tag trisa/gds:$TAG gcr.io/trisa-gds/gds:$TAG
docker tag trisa/gds-ui:$TAG gcr.io/trisa-gds/gds-ui:$TAG
docker tag trisa/grpc-proxy:$TAG gcr.io/trisa-gds/grpc-proxy:$TAG
docker tag trisa/gds-testnet-ui:$TAG gcr.io/trisa-gds/gds-testnet-ui:$TAG

docker push trisa/gds:$TAG
docker push trisa/gds-ui:$TAG
docker push trisa/gds-testnet-ui:$TAG
docker push trisa/grpc-proxy:$TAG

docker push gcr.io/trisa-gds/gds:$TAG
docker push gcr.io/trisa-gds/gds-ui:$TAG
docker push gcr.io/trisa-gds/grpc-proxy:$TAG
docker push gcr.io/trisa-gds/gds-testnet-ui:$TAG
