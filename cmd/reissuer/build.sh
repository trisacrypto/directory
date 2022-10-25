#!/bin/bash
# A helper for building the reissuer command and loading it into the GDS pod since it
# is not part of our default deployment. This script will also cleanup after itself.

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-n NAMESPACE] [clean|build|check]
A helper for common docker compose operations to run GDS
services locally. Flags are as follows (getopt required):

    -h            display this help and exit
    -n NAMESPACE  specify kubernetes namespace (default: "testnet")

There are two ways to use this script. Build and deploy the reissuer:

    ${0##*/} build

Clean up the reissuer and remove it from the pod:

    ${0##*/} clean

You can also specify the namespace to deploy it to the trisa or testnet GDS.
EOF
}

# Parse command line options with getopt
NAMESPACE=testnet
OPTIND=1

while getopts hn: opt; do
    case $opt in
        h)
            show_help
            exit 0
            ;;
        n)  NAMESPACE=$OPTARG
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
KUBECTX="$(kubectl config current-context)"

# Check if the argument is specified
if [[ $# -eq 1 ]]; then
    case $1 in
        clean)
            echo "your k8s cluster context is $KUBECTX"
            echo "removing reissuer from GDS pod in $NAMESPACE"
            kubectl -n $NAMESPACE exec -it gds-0 -- rm /usr/local/bin/reissuer
            exit 0
            ;;
        check)
            echo "your k8s cluster context is $KUBECTX"
            echo "checking if reissuer is on the GDS pod in $NAMESPACE"
            kubectl -n $NAMESPACE exec -it gds-0 -- ls -la /usr/local/bin/reissuer
            exit 0
            ;;
        build)
            echo "building reissuer for linux and sending to k8s $NAMESPACE"
            ;;
        *)
            show_help >&2
            exit 2
            ;;
    esac
fi

echo "your k8s cluster context is $KUBECTX"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$(git rev-parse --short HEAD)'" -o $PWD/reissuer $DIR
kubectl -n $NAMESPACE cp $PWD/reissuer gds-0:/usr/local/bin/reissuer
rm $PWD/reissuer