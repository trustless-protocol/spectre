#!/usr/bin/env bash
set -euo pipefail

# Cosmos → ETH transfer helper.
# Wraps the send-packet Go helper with local-test defaults.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"
SEND_PACKET_BIN="${SEND_PACKET_BIN:-$REPO_ROOT/e2e/local-tests/bin/send-packet}"

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_HOME="${COSMOS_HOME:-$HOME/.gaia}"
COSMOS_KEYRING="${COSMOS_KEYRING:-test}"
COSMOS_SIGNER_KEY="${COSMOS_SIGNER_KEY:-test1}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"

ETH_RPC_URL="${ETH_RPC_URL:-${ETH_RPC:-}}"
ETH_PRIVATE_KEY="${ETH_PRIVATE_KEY:-${ETH_USER_PK:-${PRIVATE_KEY:-}}}"

AMOUNT="${AMOUNT:-${TRANSFER_AMOUNT:-1000000000}}"
DENOM="${DENOM:-stake}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-1800}"
MEMO="${MEMO:-}"
COUNT="${COUNT:-1}"

if [ -f "$ENV_FILE" ]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

discover_kurtosis_endpoints

if [ -z "${ETH_RPC_URL:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ETH_RPC_URL="$(jq -er '[.. | objects | .eth_rpc_url? // empty][0]' "$CONFIG_FILE" 2>/dev/null || true)"
fi

if [ -z "$ETH_RPC_URL" ]; then
  log_err "missing ETH_RPC_URL"
  exit 1
fi
if [ -z "$ETH_PRIVATE_KEY" ]; then
  log_err "missing ETH_PRIVATE_KEY"
  exit 1
fi

ETH_SENDER="$(cast wallet address --private-key "$ETH_PRIVATE_KEY" 2>/dev/null)"
if [ -z "$ETH_SENDER" ]; then
  log_err "failed to derive ETH sender address"
  exit 1
fi

if [ -z "${COSMOS_WASM_CLIENT_ID:-}" ] && [ -f "$CONFIG_FILE" ]; then
  COSMOS_WASM_CLIENT_ID="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.cosmos_wasm_client_id // empty' "$CONFIG_FILE" 2>/dev/null || true)"
fi
COSMOS_WASM_CLIENT_ID="${COSMOS_WASM_CLIENT_ID:-08-wasm-0}"

if [ ! -x "$SEND_PACKET_BIN" ]; then
  log "Building send-packet helper..."
  (cd "$REPO_ROOT/relayer" && go build -o "$SEND_PACKET_BIN" ./cmd/sendpacket)
fi

log_header "Cosmos → ETH Transfer"
log_kv "COSMOS_RPC_URL" "$COSMOS_RPC_URL"
log_kv "COSMOS_CHAIN_ID" "$COSMOS_CHAIN_ID"
log_kv "ETH_SENDER" "$ETH_SENDER"
log_kv "AMOUNT" "$AMOUNT"
log_kv "DENOM" "$DENOM"
log_kv "TIMEOUT" "${TIMEOUT_SECONDS}s"
log_kv "COUNT" "$COUNT"

RECEIVER_LOWER="$(echo "$ETH_SENDER" | tr '[:upper:]' '[:lower:]')"

TX_HASH="$("$SEND_PACKET_BIN" \
  --receiver "$RECEIVER_LOWER" \
  --amount "$AMOUNT" \
  --denom "$DENOM" \
  --timeout-seconds "$TIMEOUT_SECONDS" \
  --client-id "$COSMOS_WASM_CLIENT_ID" \
  --memo "$MEMO" \
  --node "$COSMOS_RPC_URL" \
  --chain-id "$COSMOS_CHAIN_ID" \
  --count "$COUNT")"

log_ok "Transfer submitted: $TX_HASH"
