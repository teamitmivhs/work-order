#!/usr/bin/env bash
set -euo pipefail

podman() { return 0; }
curl() {
  [[ " $* " == *"%{http_code}"* ]] && printf '401'
  return 0
}
export -f podman curl

root="$(git rev-parse --show-toplevel)"
WORK_ORDER_COMPOSE_FILE="$root/src/docker-compose.persistent.yml" \
WORK_ORDER_ENV_FILE="$root/README.md" \
WORK_ORDER_BACKUP_COMMAND=/bin/true \
  "$root/scripts/deploy-server.sh" >/dev/null
