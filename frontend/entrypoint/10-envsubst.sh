#!/bin/sh
set -eu

# Default upstream if not provided (docker compose dev network)
: "${API_UPSTREAM:=http://backend:18888}"

TEMPLATE="/etc/nginx/conf.d/default.conf.template"
TARGET="/etc/nginx/conf.d/default.conf"

if [ -f "$TEMPLATE" ]; then
  echo "[entrypoint] Rendering nginx config with API_UPSTREAM=$API_UPSTREAM"
  envsubst '\$API_UPSTREAM' < "$TEMPLATE" > "$TARGET"
fi

exit 0



