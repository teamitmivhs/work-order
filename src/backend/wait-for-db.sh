#!/bin/sh

# Wait for database to be ready
host="$1"
port="$2"
shift 2

echo "Waiting for database at $host:$port..."

# Use netcat to check if the port is accessible
while ! timeout 1 nc -z "$host" "$port" 2>/dev/null; do
  echo "Database port is unavailable - sleeping"
  sleep 2
done

echo "Database port is accessible - executing command"
exec "$@"
