#!/bin/bash

GOOS=linux GOARCH=amd64 go build -o hackintosh ./cmd/hack
kubectl -n testnet cp $PWD/hackintosh gds-0:/usr/local/bin/hackintosh
#kubectl -n trisa cp $PWD/hackintosh gds-0:/usr/local/bin/hackintosh
rm hackintosh