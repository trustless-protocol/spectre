#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/.." && pwd)/setup/04-create-clients.sh" "$@"
