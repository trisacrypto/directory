#!/bin/sh

HOSTNAME=$(hostname)
echo "configuring $HOSTNAME certificates"

CERT="/data/secret/$HOSTNAME.cert.pem"
CHAIN="/data/secret/$HOSTNAME.chain.pem"

if [[ -f $CERT && -f $CHAIN ]]; then
    cp $CERT /data/certs/mtls_cert.pem
    cp $CHAIN /data/certs/mtls_chain.pem
    echo "mtls certificates configured"
    exit 0
else
    echo "could not find certificate and/or chain for configuration"
    exit 1
fi