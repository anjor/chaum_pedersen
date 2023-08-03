#!/bin/bash

set -e

host="$1"
port="$2"
shift 2
cmd="$@"

timeout 15 bash -c \
  'until nc -z "$0" "$1"; do echo "$0:$1 is unavailable - sleeping"; sleep 1; done' "$host" "$port"

>&2 echo "$host:$port is up - executing command"
exec $cmd
