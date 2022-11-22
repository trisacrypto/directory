#!/bin/bash

# Set some helpful variables
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
KUBECTX="$(kubectl config current-context)"

echo "your k8s cluster context is $KUBECTX"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$(git rev-parse --short HEAD)'" -o $PWD/trtl $DIR

echo "testnet"
kubectl -n testnet cp $PWD/trtl gds-0:/usr/local/bin/trtl
kubectl -n testnet exec -it gds-0 -- trtl metrics > $PWD/fixtures/trtl-usage/testnet.json

echo "mainnet"
kubectl -n trisa cp $PWD/trtl gds-0:/usr/local/bin/trtl
kubectl -n trisa exec -it gds-0 -- trtl metrics > $PWD/fixtures/trtl-usage/mainnet.json

echo "staging mainnet"
kubectl -n staging cp $PWD/trtl gds-mainnet-0:/usr/local/bin/trtl
kubectl -n staging exec -it gds-mainnet-0 -- trtl metrics > $PWD/fixtures/trtl-usage/staging-mainnet.json

echo "staging testnet"
kubectl -n staging cp $PWD/trtl gds-testnet-0:/usr/local/bin/trtl
kubectl -n staging exec -it gds-testnet-0 -- trtl metrics > $PWD/fixtures/trtl-usage/staging-testnet.json

rm $PWD/trtl