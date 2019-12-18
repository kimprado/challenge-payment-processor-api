#!/bin/sh

set -e

case "$1" in
    start)
        docker-compose up -d --build \
        nginx \
        api
        ;;
    start-safe)
        docker-compose up -d --build \
        nginx \
        api-safe
        ;;
    stop)
        docker-compose rm -fsv \
        swagger \
        nginx \
        api \
        api-safe \
        redisdb
        ;;
    *)
        echo "Usage: {start|stop}" >&2
       ;;
esac
