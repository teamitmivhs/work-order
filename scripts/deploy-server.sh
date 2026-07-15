#!/usr/bin/env bash
set -euo pipefail

root="$(git rev-parse --show-toplevel)"
compose="${WORK_ORDER_COMPOSE_FILE:-$root/src/docker-compose.server.yml}"
env_file="${WORK_ORDER_ENV_FILE:-$root/src/.env}"
backup="${WORK_ORDER_BACKUP_COMMAND:-$HOME/.local/bin/work-order-backup}"

[[ -f "$compose" ]] || { echo "Missing server compose: $compose" >&2; exit 1; }
[[ -f "$env_file" ]] || { echo "Missing production env: $env_file" >&2; exit 1; }
[[ -x "$backup" ]] || { echo "Missing backup command: $backup" >&2; exit 1; }

cd "$root/src"
podman compose -f "$compose" config >/dev/null

echo "Backing up production data..."
"$backup"

echo "Building and deploying $(git -C "$root" rev-parse --short HEAD)..."
podman compose -f "$compose" up -d --build --remove-orphans
podman restart work-order-nginx >/dev/null

curl -fsS --retry 30 --retry-delay 2 --retry-all-errors \
  http://127.0.0.1:4323/login >/dev/null
api_status=000
for _ in {1..30}; do
  api_status="$(curl -sS -o /dev/null -w '%{http_code}' \
    http://127.0.0.1:4323/api/profile || true)"
  [[ "$api_status" == 401 ]] && break
  sleep 2
done
[[ "$api_status" == 401 ]] || { echo "API smoke test returned HTTP $api_status" >&2; exit 1; }

podman ps --filter name=work-order --format 'table {{.Names}}\t{{.Status}}'
echo "Deployment complete."
