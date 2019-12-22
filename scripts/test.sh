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
    *)
        echo "Usage: {unit|integration|all}" >&2
       ;;
esac