#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/.." && pwd)/setup/01-eth-node.sh" "$@"
