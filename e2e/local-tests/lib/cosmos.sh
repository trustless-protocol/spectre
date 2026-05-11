#!/usr/bin/env bash

set -euo pipefail

cosmos_query() {
  "$COSMOS_BIN" query "$@" \
    --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" \
    --chain-id "${COSMOS_CHAIN_ID:-test-ibc-eth}" \
    --output json
}

key_address() {
  local key="$1"
  "$COSMOS_BIN" keys show "$key" -a \
    --keyring-backend "${COSMOS_KEYRING:-test}" \
    --home "${COSMOS_HOME:-$HOME/.gaia}" 2>/dev/null || true
}



wait_for_cosmos_block() {
  local before after
  before="$("$COSMOS_BIN" status --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"
  echo "Current block: $before" >&2
  local target=$((before + ${1:-1}))
  while true; do
    after="$("$COSMOS_BIN" status --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"
    if [ "$after" -ge "$target" ] 2>/dev/null; then
      break
    fi
    sleep 1
  done
  echo "Reached block $after" >&2
}
