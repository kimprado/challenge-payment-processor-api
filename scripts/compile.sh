#!/bin/sh

set -e

DATE=$(date +%FT%T%z)
USER=$(whoami)
GIT_VERSION=$(git --no-pager describe --tags --always)
GIT_REVISION=$(git --no-pager describe --long --always)
BRANCH=$(git branch | grep \* | cut -d ' ' -f2)

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
        wire gen --output_file_prefix=tmp_ \
        ./internal/pkg/commom/config \
        ./internal/pkg/infra/redis \
        ./internal/pkg/processor \
        ./internal/pkg/processor/api

        mv ./internal/pkg/commom/config/tmp_wire_gen.go ./internal/pkg/commom/config/wire_gen_test.go
        mv ./internal/pkg/infra/redis/tmp_wire_gen.go ./internal/pkg/infra/redis/wire_gen_test.go
        mv ./internal/pkg/processor/tmp_wire_gen.go ./internal/pkg/processor/wire_gen_test.go
        mv ./internal/pkg/processor/api/tmp_wire_gen.go ./internal/pkg/processor/api/wire_gen_test.go
        ;;
    *)
        echo "Usage: {build|wire}" >&2
       ;;
esac
