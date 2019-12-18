#!/bin/sh

set -e

case "$1" in
    package)
        docker image build -t challenge/payment-processor-api -f Dockerfile.package . 
        ;;
    package-safe)
        docker image build -t challenge/payment-processor-api -f Dockerfile.package-safe . 
        ;;
    *)
        echo "Usage: {package}" >&2
       ;;
esac
