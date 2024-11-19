#!/bin/bash
# A helper for building the bffutil command and loading it into the BFF pod.

# Print usage and exit
show_help() {
cat << EOF
Usage: ${0##*/} [-h] [-n NAMESPACE] [-a APPNAME] [clean|build|check|connect]
A helper for common docker compose operations to run GDS
services locally. Flags are as follows (getopt required):

    -h            display this help and exit
    -n NAMESPACE  specify kubernetes namespace (default: "trisa")
    -a APPNAME    specify the deployment app label (default: "bff")

There are two ways to use this script. Build and deploy the bffutil:

    ${0##*/} build

Clean up the bffutil and remove it from the pod:

    ${0##*/} clean

You can also specify the namespace to deploy it to the trisa or staging BFF.
EOF
}

# Parse command line options with getopt
NAMESPACE=trisa
APPNAME=bff
OPTIND=1

while getopts ":hn:a:" opt; do
    case ${opt} in
        h)
            show_help
            exit 0
            ;;
        n)  NAMESPACE=$OPTARG
            ;;
        a)  APPNAME=$OPTARG
            ;;
        \?)
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

if POD="$(kubectl -n $NAMESPACE get pod -l app=$APPNAME -o jsonpath="{.items[0].metadata.name}" 2>/dev/null)"; then
    read -p "continue with pod $POD in $NAMESPACE ($KUBECTX)? " -n 1 -r
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo
        exit 1
    fi
    echo
else
    echo "could not find pod in namespace $NAMESPACE with deployment $APPNAME"
    exit 1
fi

# Check if the argument is specified
if [[ $# -eq 1 ]]; then
    case $1 in
        clean)
            echo "your k8s cluster context is $KUBECTX"
            echo "removing bffutil from $POD in $NAMESPACE"
            kubectl -n $NAMESPACE exec -it $POD -- rm /usr/local/bin/bffutil
            exit 0
            ;;
        check)
            echo "your k8s cluster context is $KUBECTX"
            echo "checking if bffutil is on $POD in $NAMESPACE"
            kubectl -n $NAMESPACE exec -it $POD -- ls -la /usr/local/bin/bffutil
            exit 0
            ;;
        build)
            echo "building bffutil for linux and sending to $POD in $NAMESPACE"
            ;;
        connect)
            echo "your k8s cluster context is $KUBECTX"
            echo "connecting terminal to $POD in $NAMESPACE"
            kubectl -n $NAMESPACE exec -it $POD -- bash
            exit 0
            ;;
        *)
            show_help >&2
            exit 2
            ;;
    esac
fi

echo "your k8s cluster context is $KUBECTX"
GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/trisacrypto/directory/pkg.GitVersion=$(git rev-parse --short HEAD)'" -o $PWD/bffutil $DIR
kubectl -n $NAMESPACE cp $PWD/bffutil $POD:/usr/local/bin/bffutil
rm $PWD/bffutil