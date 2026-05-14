#!/usr/bin/env bash
set -euo pipefail

# Inspect Cosmos state for Cosmos → ETH path.
# Checks native and voucher balances.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

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
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"
COSMOS_NATIVE_DENOM="${COSMOS_NATIVE_DENOM:-stake}"

require_cmd "$COSMOS_BIN"
require_cmd jq
require_cmd sha256sum
load_contract_addresses

if [ -z "${ADDRESS:-}" ]; then
  ADDRESS="$(key_address "$COSMOS_RECEIVER_KEY")"
fi

SIGNER_ADDRESS="${COSMOS_SIGNER_ADDRESS:-$(key_address "$COSMOS_SIGNER_KEY")}"

log_header "Cosmos Endpoint"
log_kv "COSMOS_BIN" "$COSMOS_BIN"
log_kv "COSMOS_RPC_URL" "$COSMOS_RPC_URL"
log_kv "COSMOS_CHAIN_ID" "$COSMOS_CHAIN_ID"
log_kv "latest_block" "$("$COSMOS_BIN" status --node "$COSMOS_RPC_URL" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')"

log_header "Keys"
log_kv "receiver_key" "$COSMOS_RECEIVER_KEY"
log_kv "receiver_address" "${ADDRESS:-}"
log_kv "signer_key" "$COSMOS_SIGNER_KEY"
log_kv "signer_address" "${SIGNER_ADDRESS:-}"

if [ -n "${ADDRESS:-}" ]; then
  log_header "Account Balances"
  cosmos_query bank balances "$ADDRESS" | jq .
fi

if [ -n "${SIGNER_ADDRESS:-}" ]; then
  log_header "Signer Native Balance"
  cosmos_query bank balance "$SIGNER_ADDRESS" "$COSMOS_NATIVE_DENOM" | jq .
fi

if [ -n "${ADDRESS:-}" ] && [ -n "${ERC20_ADDRESS:-}" ]; then
  # Read cosmos_wasm_client_id from config
  if [ -z "${COSMOS_WASM_CLIENT_ID:-}" ] && [ -f "$REPO_ROOT/relayer/config.json" ]; then
    COSMOS_WASM_CLIENT_ID="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.cosmos_wasm_client_id // empty' "$REPO_ROOT/relayer/config.json" 2>/dev/null || true)"
  fi
  COSMOS_WASM_CLIENT_ID="${COSMOS_WASM_CLIENT_ID:-08-wasm-0}"

  VOUCHER_TRACE="transfer/$COSMOS_WASM_CLIENT_ID/$ERC20_ADDRESS"
  VOUCHER_DENOM="$(derive_ibc_denom "$VOUCHER_TRACE")"
  log_header "ETH → Cosmos Voucher Balance"
  log_kv "trace" "$VOUCHER_TRACE"
  log_kv "denom" "$VOUCHER_DENOM"
  cosmos_query bank balance "$ADDRESS" "$VOUCHER_DENOM" | jq .
fi

log_header "IBC Clients"
cosmos_query ibc client states | jq .
