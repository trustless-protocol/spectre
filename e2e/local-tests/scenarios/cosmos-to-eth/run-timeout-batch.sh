#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Multiple Cosmos→ETH packets timeout, refund on Cosmos.
#
# Corresponds to: Test_10_TimeoutPacketFromCosmos

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"
source "$REPO_ROOT/e2e/local-tests/lib/eth.sh"
source "$REPO_ROOT/e2e/local-tests/lib/relayer.sh"

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
AMOUNT="${AMOUNT:-1000000000}"
COSMOS_NATIVE_DENOM="${COSMOS_NATIVE_DENOM:-stake}"
BATCH_SIZE="${BATCH_SIZE:-5}"

config_value() {
  local key="$1"
  if [ -f "$CONFIG_FILE" ]; then
    jq -er --arg key "$key" '([.. | objects | .[$key]? // empty][0]) // empty' "$CONFIG_FILE" 2>/dev/null || true
  fi
}

load_scenario_contract_addresses() {
  if [ -z "${ICS26_ADDRESS:-}" ]; then
    ICS26_ADDRESS="$(config_value ics26_address)"
  fi
  if [ -z "${ICS26_ADDRESS:-}" ]; then
    ICS26_ADDRESS="$(config_value ics26Router)"
  fi

  if [ -z "${ICS20_ADDRESS:-}" ]; then
    ICS20_ADDRESS="$(config_value ics20_address)"
  fi
  if [ -z "${ICS20_ADDRESS:-}" ]; then
    ICS20_ADDRESS="$(config_value ics20Transfer)"
  fi

  if [ -z "${BROADCAST_JSON:-}" ]; then
    local chain_id chain_broadcast
    chain_id="$(cast chain-id --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
    chain_broadcast="$REPO_ROOT/broadcast/E2ETestDeploy.s.sol/$chain_id/run-latest.json"
    if [ -n "$chain_id" ] && [ -f "$chain_broadcast" ]; then
      BROADCAST_JSON="$chain_broadcast"
    elif [ -d "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" ]; then
      BROADCAST_JSON="$(find "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" -path '*/run-latest.json' -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | awk 'NR == 1 { print $2 }' || true)"
    fi
  fi

  if [ -f "${BROADCAST_JSON:-}" ]; then
    if [ -z "${ERC20_ADDRESS:-}" ]; then
      ERC20_ADDRESS="$(jq -er '(.erc20 // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .erc20))' "$BROADCAST_JSON" 2>/dev/null || true)"
    fi
    if [ -z "${ICS20_ADDRESS:-}" ]; then
      ICS20_ADDRESS="$(jq -er '(.ics20Transfer // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .ics20Transfer))' "$BROADCAST_JSON" 2>/dev/null || true)"
    fi
  fi
}

require_cmd cast
require_cmd jq

load_env_file
ETH_PRIVATE_KEY="${ETH_PRIVATE_KEY:-${ETH_USER_PK:-${PRIVATE_KEY:-}}}"
discover_kurtosis_endpoints

if [ -z "${ETH_RPC_URL:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ETH_RPC_URL="$(config_value eth_rpc_url)"
fi

load_scenario_contract_addresses

if [ -z "${COSMOS_WASM_CLIENT_ID:-}" ] && [ -f "$CONFIG_FILE" ]; then
  COSMOS_WASM_CLIENT_ID="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.cosmos_wasm_client_id // empty' "$CONFIG_FILE" 2>/dev/null || true)"
fi
COSMOS_WASM_CLIENT_ID="${COSMOS_WASM_CLIENT_ID:-08-wasm-0}"

COSMOS_SENDER_ADDRESS="$(key_address "$COSMOS_SIGNER_KEY")"
ETH_SENDER="$(eth_wallet_address)"

if [ -z "${ETH_RPC_URL:-}" ]; then
  log_err "ETH_RPC_URL not set"
  exit 1
fi
if [ -z "${ETH_PRIVATE_KEY:-}" ]; then
  log_err "ETH_PRIVATE_KEY not set"
  exit 1
fi
if [ -z "${COSMOS_SENDER_ADDRESS:-}" ]; then
  log_err "COSMOS_SIGNER_KEY '$COSMOS_SIGNER_KEY' not found in keyring"
  exit 1
fi
if [ ! -x "$SEND_PACKET_BIN" ]; then
  log "Building send-packet helper..."
  (cd "$REPO_ROOT/relayer" && go build -o "$SEND_PACKET_BIN" ./cmd/sendpacket)
fi

relayer_ensure_running

# ─── Record pre-transfer state ───
log_header "Pre-transfer State"
BEFORE_NATIVE_BALANCE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos native balance" "$BEFORE_NATIVE_BALANCE"

TOTAL_AMOUNT="$(echo "$AMOUNT * $BATCH_SIZE" | bc)"
log_kv "Batch size" "$BATCH_SIZE"
log_kv "Total amount" "$TOTAL_AMOUNT"

# ─── Send batch with short timeout ───
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-45}"
SEND_MEMO="cosmos-timeout-batch"

log_header "Sending ${BATCH_SIZE} Cosmos → ETH Transfers (${TIMEOUT_SECONDS}s timeout)"

COSMOS_SEND_TX_HASH="$("$SEND_PACKET_BIN" \
  --receiver "$(echo "$ETH_SENDER" | tr '[:upper:]' '[:lower:]')" \
  --amount "$AMOUNT" \
  --denom "$COSMOS_NATIVE_DENOM" \
  --timeout-seconds "$TIMEOUT_SECONDS" \
  --client-id "$COSMOS_WASM_CLIENT_ID" \
  --memo "$SEND_MEMO" \
  --node "$COSMOS_RPC_URL" \
  --chain-id "$COSMOS_CHAIN_ID" \
  --count "$BATCH_SIZE")"

log_ok "Batch transfer submitted: $COSMOS_SEND_TX_HASH"
log "Waiting for timeout..."

# ─── Poll for refund ───
REFUND_THRESHOLD="$(bc_sub1 "$BEFORE_NATIVE_BALANCE")"
poll=0
MAX_POLLS=120
POLL_INTERVAL=15
REFUNDED=false

log_header "Polling for Refund"
while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_NATIVE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

  if bc_ge "$CURRENT_NATIVE" "$REFUND_THRESHOLD"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] Native refunded │ balance=${CURRENT_NATIVE} (target ≥ ${REFUND_THRESHOLD})"
    REFUNDED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] native=${CURRENT_NATIVE}"
done

# ─── Final state ───
log_header "Final State"
AFTER_NATIVE_BALANCE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos native balance" "$AFTER_NATIVE_BALANCE"

# ─── Summary ───
log_header "Cosmos Timeout Batch Test Summary"
log_kv "Direction" "Cosmos → ETH"
log_kv "Batch size" "$BATCH_SIZE"
log_kv "Amount per packet" "${AMOUNT} ${COSMOS_NATIVE_DENOM}"
log_kv "Total sent" "${TOTAL_AMOUNT} ${COSMOS_NATIVE_DENOM}"
log_kv "Timeout" "${TIMEOUT_SECONDS}s"

if [ "${REFUNDED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • Relayer detected the expired packets"
  log "  • All ${BATCH_SIZE} packets timed out and were refunded"
  log "  • Cosmos native balance returned to original"
else
  log_err "RESULT: FAIL"
  log_err "  • Native balance not refunded after ${MAX_POLLS} polls"
  log_err "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  exit 1
fi
