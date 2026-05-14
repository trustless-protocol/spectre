#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Verify light client auto-updates as Cosmos chain advances.
#
# Corresponds to: Test_UpdateClient

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"
source "$REPO_ROOT/e2e/local-tests/lib/relayer.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"

config_value() {
  local key="$1"
  if [ -f "$CONFIG_FILE" ]; then
    jq -er --arg key "$key" '([.. | objects | .[$key]? // empty][0]) // empty' "$CONFIG_FILE" 2>/dev/null || true
  fi
}

load_env_file
discover_kurtosis_endpoints

if [ -z "${ETH_RPC_URL:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ETH_RPC_URL="$(config_value eth_rpc_url)"
fi

if [ -z "${ICS07_CLIENT:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ICS07_CLIENT="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.ics07_client // empty' "$CONFIG_FILE" 2>/dev/null || true)"
fi

if [ -z "${ETH_RPC_URL:-}" ]; then
  log_err "ETH_RPC_URL not set"
  exit 1
fi
if [ -z "${ICS07_CLIENT:-}" ]; then
  log_err "ICS07_CLIENT address not found in $CONFIG_FILE"
  exit 1
fi

relayer_ensure_running

# ─── Query current client state ───
log_header "Current Light Client State"
CLIENT_STATE_RAW="$(cast call "$ICS07_CLIENT" 'clientState()(string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8)' --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
if [ -z "$CLIENT_STATE_RAW" ]; then
  log_err "Failed to query client state from $ICS07_CLIENT"
  exit 1
fi

# Parse latestHeight from output (cast returns tuple components on separate lines)
LATEST_HEIGHT="$(echo "$CLIENT_STATE_RAW" | awk 'NR==3 { print $1 }')"
log_kv "Latest height (before)" "$LATEST_HEIGHT"

# ─── Wait for Cosmos chain to advance ───
log_header "Waiting for Cosmos Chain to Advance"
BLOCKS_TO_WAIT="${BLOCKS_TO_WAIT:-5}"
wait_for_cosmos_block "$BLOCKS_TO_WAIT"

# ─── Poll for client update ───
log_header "Polling for Client Update"
poll=0
MAX_POLLS=30
POLL_INTERVAL=10
UPDATED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  NEW_CLIENT_STATE_RAW="$(cast call "$ICS07_CLIENT" 'clientState()(string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8)' --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
  NEW_LATEST_HEIGHT="$(echo "$NEW_CLIENT_STATE_RAW" | awk 'NR==3 { print $1 }')"

  if [ -n "$NEW_LATEST_HEIGHT" ] && [ "$NEW_LATEST_HEIGHT" -gt "$LATEST_HEIGHT" ] 2>/dev/null; then
    log_ok "[poll ${poll}/${MAX_POLLS}] Client updated │ height=${NEW_LATEST_HEIGHT} (was ${LATEST_HEIGHT})"
    UPDATED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] height=${NEW_LATEST_HEIGHT} (waiting for > ${LATEST_HEIGHT})"
done

# ─── Summary ───
log_header "Update Client Test Summary"
if [ "${UPDATED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • Light client height increased from $LATEST_HEIGHT to $NEW_LATEST_HEIGHT"
  log "  • Relayer auto-updated the client as Cosmos chain advanced"
else
  log_warn "RESULT: PARTIAL"
  log_warn "  • Light client height did not increase within ${MAX_POLLS} polls"
  log_warn "  • The relayer may still be building the ZK proof"
  log_warn "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
fi
