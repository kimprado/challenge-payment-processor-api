#!/bin/sh

set -e

case "$1" in
    build)
        go build -ldflags \
        '-X github.com/prometheus/common/version.Version='$GIT_VERSION'
        -X github.com/prometheus/common/version.BuildDate='$DATE' 
        -X github.com/prometheus/common/version.Branch='$BRANCH' 
        -X github.com/prometheus/common/version.Revision='$GIT_REVISION'
        -X github.com/prometheus/common/version.BuildUser='$USER'' \
        -o ./payment-processor-api.bin github.com/challenge/payment-processor/cmd/processorAPI 
        ;;
    build-static)
        WORKBUILD=${2:-/usr/dist}
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo \
        -ldflags \
        '-X github.com/prometheus/common/version.Version='$GIT_VERSION'
        -X github.com/prometheus/common/version.BuildDate='$DATE' 
        -X github.com/prometheus/common/version.Branch='$BRANCH' 
        -X github.com/prometheus/common/version.Revision='$GIT_REVISION'
        -X github.com/prometheus/common/version.BuildUser='$USER'' \
        -o $WORKBUILD/payment-processor-api.bin github.com/challenge/payment-processor/cmd/processorAPI 
        ;;
    wire)
        wire ./cmd/processorAPI/
        ;;
    wire-testes)
        wire \
        ./internal/pkg/commom/config \
        ./internal/pkg/infra/redis \
        ;;
    generate)
        go generate ./cmd/processorAPI/
        ;;
    generate-testes)
        go generate \
        ./internal/pkg/commom/config \
        ./internal/pkg/infra/redis \
        ;;
    *)
        echo "Usage: {build|wire}" >&2
       ;;
esac
