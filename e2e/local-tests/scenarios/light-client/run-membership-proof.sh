#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Verify on-chain membership proof.
# Uses the relayer `fixtures membership` command to build and verify a proof.
#
# Corresponds to: Test_Membership

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"
COSMOS_NATIVE_DENOM="${COSMOS_NATIVE_DENOM:-stake}"

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

# Build relayer binary if needed
RELAYER_DIR="${RELAYER_DIR:-$REPO_ROOT/relayer}"
RELAYER_BIN="${RELAYER_BIN:-$RELAYER_DIR/relayer}"
if [ ! -x "$RELAYER_BIN" ]; then
  log "Building relayer..."
  (cd "$RELAYER_DIR" && go build -o relayer ./cmd)
fi

# Query a known bank balance key to verify membership for
QUERY_ADDRESS="${QUERY_ADDRESS:-$(key_address test)}"
if [ -z "$QUERY_ADDRESS" ]; then
  log_err "Could not determine query address"
  exit 1
fi

log_header "Membership Proof Test"
log_kv "ICS07_CLIENT" "$ICS07_CLIENT"
log_kv "Query address" "$QUERY_ADDRESS"
log_kv "Query denom" "$COSMOS_NATIVE_DENOM"

# The relayer fixtures membership command needs:
# - TENDERMINT_RPC_URL (from .env)
# - ETH_RPC_URL (from .env)
# - CONTRACT_ADDRESS (ICS07 client address)
# - PRIVATE_KEY (ETH private key for tx submission)
# It takes args: <key_path> <is_base64> <membership_type>

# We need to construct the key path for bank balance. This is a merkle path in the
# IBC store. The fixtures command expects the key path as a string.
# For a simple test, we'll use the bank balance prefix approach.
# Note: constructing exact merkle paths requires knowledge of the store key layout.
# This script is a skeleton that documents the required flow.

log_warn "Membership proof construction requires exact store key paths."
log_warn "The relayer fixtures membership command can be used as:"
log_warn "  TENDERMINT_RPC_URL=... ETH_RPC_URL=... CONTRACT_ADDRESS=$ICS07_CLIENT PRIVATE_KEY=..."
log_warn "  $RELAYER_BIN fixtures membership <key_path> <is_base64> <membership_type>"

log_ok "RESULT: INFO"
log "  • To fully automate this, the key_path must be constructed from the target state key."
log "  • See relayer/cmd/main.go MembershipCmd for the exact command format."
