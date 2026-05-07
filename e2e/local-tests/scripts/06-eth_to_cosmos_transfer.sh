#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/.." && pwd)/scenarios/eth-to-cosmos/transfer.sh" "$@"
