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

upsert_env_var() {
  local file="$1"
  local key="$2"
  local value="$3"
  local tmp

  mkdir -p "$(dirname "$file")"
  touch "$file"
  tmp="$(mktemp)"

  awk -v key="$key" -v value="$value" '
    BEGIN { updated = 0 }
    $0 ~ ("^" key "=") {
      print key "=\"" value "\""
      updated = 1
      next
    }
    { print }
    END {
      if (!updated) {
        print key "=\"" value "\""
      }
    }
  ' "$file" > "$tmp"

  mv "$tmp" "$file"
}

poll_cosmos_balance() {
  local address="$1" denom="$2" target="$3" timeout_sec="${4:-120}" poll_interval="${5:-15}"
  local poll=0 current
  while [ $poll -lt "$timeout_sec" ]; do
    current="$("$COSMOS_BIN" query bank balance "$address" "$denom" \
      --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" \
      --chain-id "${COSMOS_CHAIN_ID:-test-ibc-eth}" \
      --output json 2>/dev/null | jq -r '.balance.amount // "0"' || true)"
    if [ "$current" = "$target" ] 2>/dev/null; then
      echo "$current"
      return 0
    fi
    sleep "$poll_interval"
    poll=$((poll + poll_interval))
  done
  echo "$current"
  return 1
}

wait_for_cosmos_block() {
  local before after
  before="$("$COSMOS_BIN" status --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"
  log "Current block: $before"
  local target=$((before + ${1:-1}))
  while true; do
    after="$("$COSMOS_BIN" status --node "${COSMOS_RPC_URL:-http://127.0.0.1:26657}" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"
    if [ "$after" -ge "$target" ] 2>/dev/null; then
      break
    fi
    sleep 1
  done
  log_ok "Reached block $after"
}
