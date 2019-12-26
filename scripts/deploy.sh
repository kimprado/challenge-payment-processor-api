#!/bin/sh

set -e

case "$1" in
    start)
        docker-compose up -d --build \
        prometheus \
        grafana \
        acquirers \
        nginx \
        api
        ;;
    start-safe)
        docker-compose up -d --build \
        prometheus \
        grafana \
        acquirers \
        nginx \
        api-safe
        ;;
    stop)
        docker-compose rm -fsv \
        acquirers \
        swagger \
        nginx \
        api \
        api-safe \
        redisdb \
        prometheus \
        grafana
        ;;
    *)
        echo "Usage: {start|stop}" >&2
       ;;
esac
