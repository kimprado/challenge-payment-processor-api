#!/bin/sh

set -e

case "$1" in
    unit)
        go test ./... -tags="unit" -cover -coverprofile=coverage.out
        go tool cover -func=coverage.out | tail -n 1
        ;;
    integration)
        go test -parallel 10 -timeout 1m30s ./... -tags="integration" -cover -coverprofile=coverage.out
        go tool cover -func=coverage.out | tail -n 1
        ;;
    all)
        go test -parallel 10 -timeout 1m30s ./... -tags="test" -cover -coverprofile=coverage.out
        go tool cover -func=coverage.out | tail -n 1
        ;;
    envvars)
        PROCESSOR_ENVIRONMENT_NAME="test_ENV-VARS" \
        PROCESSOR_SERVER_PORT="4033" \
        PROCESSOR_REDISDB_HOST="host-env-test" \
        PROCESSOR_REDISDB_PORT="6523" \
        PROCESSOR_LOGGING_LEVEL="ROOT: WARN-teste" \
        go test ./... -tags="testenvvars" -cover -coverprofile=coverage.out
        go tool cover -func=coverage.out | tail -n 1
        ;;
    *)
        echo "Usage: {unit|integration|all|envvars}" >&2
       ;;
esac