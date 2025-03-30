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

# Get the tag as the first argument or from git if none is supplied
if [ -z "$1" ]; then
    TAG=$(git rev-parse --short HEAD)
else
    TAG=$1
fi

# Get the platform as the second argument or use linux/amd64 by default
if [ -z "$2" ]; then
    PLATFORM="linux/amd64"
else
    PLATFORM=$2
fi

# Confirm that we're continuing with the tag
if ! ask "Continue with tag $TAG?" N; then
    exit 1
fi

# Set some helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
REPO=$(realpath "$DIR/..")
DOTENV="$REPO/.env"

# Load .env file from project root if it exists
if [ -f $DOTENV ]; then
    set -o allexport
    source $DOTENV
    set +o allexport
fi

# Check the environment variables
REQUIRED_ENV=(
    "REACT_APP_VASPDIRECTORY_CLIENT_ID"
    "REACT_APP_VASPDIRECTORY_ANALYTICS_ID"
    "REACT_APP_TRISATEST_CLIENT_ID"
    "REACT_APP_TRISATEST_ANALYTICS_ID"
    "REACT_APP_STAGING_VASPDIRECTORY_ANALYTICS_ID"
    "REACT_APP_STAGING_VASPDIRECTORY_CLIENT_ID"
    "REACT_APP_STAGING_TESTNET_ANALYTICS_ID"
    "REACT_APP_STAGING_TRISATEST_CLIENT_ID"
    "REACT_APP_AUTH0_CLIENT_ID"
    "REACT_APP_STAGING_AUTH0_CLIENT_ID"
)

for reqvar in ${REQUIRED_ENV[@]}; do
    if [ -z "${!reqvar+x}" ]; then
        echo "$reqvar environment variable is required"
        exit 1
    fi
done

if [ -z "$GIT_REVISION" ]; then
    export GIT_REVISION=$(git rev-parse --short HEAD)
fi

if [ -z "$REACT_APP_GIT_REVISION" ]; then
    export REACT_APP_GIT_REVISION=$(git rev-parse --short HEAD)
fi

if [ -z "$REACT_APP_VERSION_NUMBER" ]; then
    export REACT_APP_VERSION_NUMBER="$(git describe --exact-match --abbrev=0).dev"
fi

# Build the primary backend images
docker buildx build --platform $PLATFORM -t trisa/gds:$TAG -f $DIR/gds/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
docker buildx build --platform $PLATFORM -t trisa/gds-bff:$TAG -f $DIR/bff/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
docker buildx build --platform $PLATFORM -t trisa/trtl:$TAG -f $DIR/trtl/Dockerfile --build-arg GIT_REVISION=${GIT_REVISION} $REPO
docker buildx build --platform $PLATFORM -t trisa/trtl-init:$TAG -f $DIR/trtl-init/Dockerfile $DIR/trtl-init
docker buildx build --platform $PLATFORM -t trisa/trtlsim:$TAG -f $DIR/trtlsim/Dockerfile .
docker buildx build --platform $PLATFORM -t trisa/cathy:$TAG -f $DIR/cathy/Dockerfile .
docker buildx build --platform $PLATFORM -t trisa/maintenance:$TAG -f $DIR/maintenance/Dockerfile .

# Build the UI image for trisa.directory
docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-user-ui:$TAG -f $DIR/gds-user-ui/Dockerfile \
    --build-arg REACT_APP_TRISA_BASE_URL=https://bff.trisa.directory/v1/ \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_VASPDIRECTORY_ANALYTICS_ID} \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_SENTRY_DSN} \
    --build-arg REACT_APP_AUTH0_DOMAIN=${REACT_APP_AUTH0_DOMAIN} \
    --build-arg REACT_APP_AUTH0_CLIENT_ID=${REACT_APP_AUTH0_CLIENT_ID} \
    --build-arg REACT_APP_AUTH0_REDIRECT_URI=https://trisa.directory/auth/callback \
    --build-arg REACT_APP_AUTH0_SCOPE="openid profile email" \
    --build-arg REACT_APP_AUTH0_AUDIENCE=https://bff.trisa.directory \
    $REPO

# Build the UI image for vaspdirectory.dev
docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-staging-user-ui:$TAG -f $DIR/gds-user-ui/Dockerfile \
    --build-arg REACT_APP_TRISA_BASE_URL=https://bff.vaspdirectory.dev/v1/ \
    --build-arg REACT_APP_ANALYTICS_ID=${REACT_APP_STAGING_VASPDIRECTORY_ANALYTICS_ID} \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_SENTRY_DSN} \
    --build-arg REACT_APP_SENTRY_ENVIRONMENT="staging" \
    --build-arg REACT_APP_AUTH0_DOMAIN=${REACT_APP_AUTH0_DOMAIN} \
    --build-arg REACT_APP_AUTH0_CLIENT_ID=${REACT_APP_STAGING_AUTH0_CLIENT_ID} \
    --build-arg REACT_APP_AUTH0_REDIRECT_URI=https://vaspdirectory.dev/auth/callback \
    --build-arg REACT_APP_AUTH0_SCOPE="openid profile email" \
    --build-arg REACT_APP_AUTH0_AUDIENCE=https://bff.vaspdirectory.dev \
    --build-arg REACT_APP_USE_DASH_LOCALE=true \
    $REPO

# Build the Admin UI images for admin.testnet.directory and admin.trisa.directory
docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.trisa.directory/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=false \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_VASPDIRECTORY_CLIENT_ID} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_ADMIN_SENTRY_DSN} \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    $REPO

docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-testnet-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.testnet.directory/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=true \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_TRISATEST_CLIENT_ID} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_ADMIN_SENTRY_DSN} \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    $REPO

# Build the Admin UI images for admin.trisatest.dev and admin.vaspdirectory.dev
docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-staging-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.vaspdirectory.dev/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=false \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_STAGING_VASPDIRECTORY_CLIENT_ID} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_ADMIN_SENTRY_DSN} \
    --build-arg REACT_APP_SENTRY_ENVIRONMENT="staging" \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    $REPO

docker buildx build \
    --platform $PLATFORM \
    -t trisa/gds-staging-testnet-admin-ui:$TAG -f $DIR/gds-admin-ui/Dockerfile \
    --build-arg REACT_APP_GDS_API_ENDPOINT=https://api.admin.trisatest.dev/v2 \
    --build-arg REACT_APP_GDS_IS_TESTNET=true \
    --build-arg REACT_APP_GOOGLE_CLIENT_ID=${REACT_APP_STAGING_TRISATEST_CLIENT_ID} \
    --build-arg REACT_APP_SENTRY_DSN=${REACT_APP_ADMIN_SENTRY_DSN} \
    --build-arg REACT_APP_SENTRY_ENVIRONMENT="staging" \
    --build-arg REACT_APP_VERSION_NUMBER=${REACT_APP_VERSION_NUMBER} \
    --build-arg REACT_APP_GIT_REVISION=${REACT_APP_GIT_REVISION} \
    $REPO

# Retag the images to push to gcr.io
docker tag trisa/gds:$TAG gcr.io/trisa-gds/gds:$TAG
docker tag trisa/gds-bff:$TAG gcr.io/trisa-gds/gds-bff:$TAG
docker tag trisa/trtl:$TAG gcr.io/trisa-gds/trtl:$TAG
docker tag trisa/trtl-init:$TAG gcr.io/trisa-gds/trtl-init:$TAG
docker tag trisa/trtlsim:$TAG gcr.io/trisa-gds/trtlsim:$TAG
docker tag trisa/gds-user-ui:$TAG gcr.io/trisa-gds/gds-user-ui:$TAG
docker tag trisa/gds-admin-ui:$TAG gcr.io/trisa-gds/gds-admin-ui:$TAG
docker tag trisa/gds-testnet-admin-ui:$TAG gcr.io/trisa-gds/gds-testnet-admin-ui:$TAG
docker tag trisa/gds-staging-user-ui:$TAG gcr.io/trisa-gds/gds-staging-user-ui:$TAG
docker tag trisa/gds-staging-admin-ui:$TAG gcr.io/trisa-gds/gds-staging-admin-ui:$TAG
docker tag trisa/gds-staging-testnet-admin-ui:$TAG gcr.io/trisa-gds/gds-staging-testnet-admin-ui:$TAG
docker tag trisa/cathy:$TAG gcr.io/trisa-gds/cathy:$TAG
docker tag trisa/maintenance:$TAG gcr.io/trisa-gds/maintenance:$TAG

# Push to DockerHub
docker push trisa/gds:$TAG
docker push trisa/gds-bff:$TAG
docker push trisa/trtl:$TAG
docker push trisa/trtl-init:$TAG
docker push trisa/trtlsim:$TAG
docker push trisa/gds-user-ui:$TAG
docker push trisa/gds-admin-ui:$TAG
docker push trisa/gds-testnet-admin-ui:$TAG
docker push trisa/gds-staging-user-ui:$TAG
docker push trisa/gds-staging-admin-ui:$TAG
docker push trisa/gds-staging-testnet-admin-ui:$TAG
docker push trisa/cathy:$TAG
docker push trisa/maintenance:$TAG

# Push to GCR
docker push gcr.io/trisa-gds/gds:$TAG
docker push gcr.io/trisa-gds/gds-bff:$TAG
docker push gcr.io/trisa-gds/trtl:$TAG
docker push gcr.io/trisa-gds/trtl-init:$TAG
docker push gcr.io/trisa-gds/trtlsim:$TAG
docker push gcr.io/trisa-gds/gds-user-ui:$TAG
docker push gcr.io/trisa-gds/gds-admin-ui:$TAG
docker push gcr.io/trisa-gds/gds-testnet-admin-ui:$TAG
docker push gcr.io/trisa-gds/gds-staging-user-ui:$TAG
docker push gcr.io/trisa-gds/gds-staging-admin-ui:$TAG
docker push gcr.io/trisa-gds/gds-staging-testnet-admin-ui:$TAG
docker push gcr.io/trisa-gds/cathy:$TAG
docker push gcr.io/trisa-gds/maintenance:$TAG