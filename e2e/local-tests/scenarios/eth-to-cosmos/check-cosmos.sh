#!/usr/bin/env bash
set -euo pipefail

# Inspect local Cosmos data for the ETH -> Cosmos path.
# Optional inputs:
#   ADDRESS=cosmos1...         account to inspect
#   COSMOS_TX_HASH=...         Cosmos tx hash to inspect
#   ERC20_ADDRESS=0x...        used to derive transfer/08-wasm-0/<erc20> denom
#   COSMOS_RECEIVER_KEY=test   key to inspect when ADDRESS is unset

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"
BROADCAST_JSON="${BROADCAST_JSON:-}"

if [ -f "$ENV_FILE" ]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_HOME="${COSMOS_HOME:-$HOME/.gaia}"
COSMOS_KEYRING="${COSMOS_KEYRING:-test}"
COSMOS_RECEIVER_KEY="${COSMOS_RECEIVER_KEY:-test}"
COSMOS_SIGNER_KEY="${COSMOS_SIGNER_KEY:-test1}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_GRPC_ADDR="${COSMOS_GRPC_ADDR:-127.0.0.1:9090}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"
COSMOS_TX_HASH="${COSMOS_TX_HASH:-${TX_HASH:-}}"

# Read cosmos_wasm_client_id from config
if [ -z "${COSMOS_WASM_CLIENT_ID:-}" ] && [ -f "$REPO_ROOT/relayer/config.json" ]; then
  COSMOS_WASM_CLIENT_ID="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.cosmos_wasm_client_id // empty' "$REPO_ROOT/relayer/config.json" 2>/dev/null || true)"
fi
COSMOS_WASM_CLIENT_ID="${COSMOS_WASM_CLIENT_ID:-08-wasm-0}"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd" >&2
    exit 1
  fi
}

load_contract_addresses() {
  if [ -n "${ERC20_ADDRESS:-}" ]; then
    return 0
  fi

  if [ -z "$BROADCAST_JSON" ]; then
    if [ -d "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" ]; then
      BROADCAST_JSON="$(find "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" -path '*/run-latest.json' -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | awk 'NR == 1 { print $2 }')"
    fi
  fi

  if [ ! -f "$BROADCAST_JSON" ]; then
    return 0
  fi

  if [ -z "${ERC20_ADDRESS:-}" ]; then
    ERC20_ADDRESS="$(jq -er '(.erc20 // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .erc20))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi
}

cosmos_query() {
  "$COSMOS_BIN" query "$@" \
    --node "$COSMOS_RPC_URL" \
    --chain-id "$COSMOS_CHAIN_ID" \
    --output json
}

key_address() {
  local key="$1"
  "$COSMOS_BIN" keys show "$key" -a \
    --keyring-backend "$COSMOS_KEYRING" \
    --home "$COSMOS_HOME" 2>/dev/null || true
}

ibc_denom() {
  local trace="$1"
  printf 'ibc/%s\n' "$(printf '%s' "$trace" | sha256sum | awk '{ print toupper($1) }')"
}

require_cmd "$COSMOS_BIN"
require_cmd jq
require_cmd sha256sum
load_contract_addresses

if [ -z "${ADDRESS:-}" ]; then
  ADDRESS="${COSMOS_RECEIVER_ADDRESS:-}"
fi

if [ -z "${ADDRESS:-}" ]; then
  ADDRESS="$(key_address "$COSMOS_RECEIVER_KEY")"
fi

SIGNER_ADDRESS="${COSMOS_SIGNER_ADDRESS:-$(key_address "$COSMOS_SIGNER_KEY")}"

echo "━━━ Cosmos Endpoint ━━━"
printf "  %-20s %s\n" "COSMOS_BIN:"      "$COSMOS_BIN"
printf "  %-20s %s\n" "COSMOS_RPC_URL:"  "$COSMOS_RPC_URL"
printf "  %-20s %s\n" "COSMOS_GRPC:"     "$COSMOS_GRPC_ADDR"
printf "  %-20s %s\n" "COSMOS_CHAIN_ID:" "$COSMOS_CHAIN_ID"
printf "  %-20s %s\n" "latest_block:"    "$("$COSMOS_BIN" status --node "$COSMOS_RPC_URL" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"
printf "  %-20s %s\n" "BROADCAST_JSON:"  "${BROADCAST_JSON:-}"

echo ""
echo "━━━ Keys ━━━"
printf "  %-20s %s\n" "receiver_key:"    "$COSMOS_RECEIVER_KEY"
printf "  %-20s %s\n" "receiver_address:" "${ADDRESS:-}"
printf "  %-20s %s\n" "signer_key:"      "$COSMOS_SIGNER_KEY"
printf "  %-20s %s\n" "signer_address:"  "${SIGNER_ADDRESS:-}"

if [ -n "${ADDRESS:-}" ]; then
  echo ""
  echo "━━━ Account Balances ━━━"
  cosmos_query bank balances "$ADDRESS" | jq .

  if [ -n "${ERC20_ADDRESS:-}" ]; then
    VOUCHER_TRACE="transfer/$COSMOS_WASM_CLIENT_ID/$ERC20_ADDRESS"
    VOUCHER_DENOM="$(ibc_denom "$VOUCHER_TRACE")"
    echo ""
    echo "━━━ ETH → Cosmos Voucher Balance ━━━"
    printf "  %-20s %s\n" "trace:" "$VOUCHER_TRACE"
    printf "  %-20s %s\n" "denom:" "$VOUCHER_DENOM"
    cosmos_query bank balance "$ADDRESS" "$VOUCHER_DENOM" | jq .
  fi
fi

echo ""
echo "━━━ IBC Clients ━━━"
cosmos_query ibc client states | jq .

if [ -n "$COSMOS_TX_HASH" ]; then
  echo ""
  echo "━━━ Cosmos Transaction ━━━"
  printf "  %-20s %s\n" "tx_hash:" "$COSMOS_TX_HASH"
  cosmos_query tx "$COSMOS_TX_HASH" | jq .
fi
