#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/.." && pwd)/setup/02-cosmos-node.sh" "$@"
