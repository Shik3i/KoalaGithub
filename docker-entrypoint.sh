#!/bin/sh
set -eu

if [ "$(id -u)" = "0" ]; then
	chown -R koala:koala /data
	exec su-exec koala:koala "$@"
fi

exec "$@"
