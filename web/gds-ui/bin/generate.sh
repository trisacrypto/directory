#!/bin/bash

# Locate the directory of the script in order to compute relative paths correctly
BINDIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
PROJECT="$BINDIR/../../.."

# Check to ensure the trisacrypto/trisa repository has been cloned
if [ ! -d "$PROJECT/../trisa/proto" ]; then
    echo "Please ensure that the github.com/trisacrypto/trisa repo has been cloned into the directory containing this repository"
    exit 1
fi

protoc -I "$PROJECT/../trisa/proto" -I $PROJECT/proto \
    --js_out=import_style=commonjs:"$PROJECT/web/gds-ui/src/api" \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:"$PROJECT/web/gds-ui/src/api" \
    gds/admin/v1/admin.proto \
    trisa/gds/api/v1beta1/api.proto \
    trisa/gds/models/v1beta1/models.proto \
    trisa/gds/models/v1beta1/ca.proto
