#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Voucher (originally Cosmos native) sent from ETH→Cosmos times out
# → voucher (ibcERC20) reminted on ETH.
#
# Corresponds to: Test_TimeoutPacketEthRemintsVouchers

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
SOURCE_CLIENT="${SOURCE_CLIENT:-cosmoshub-1}"
DEST_PORT="${DEST_PORT:-transfer}"
AMOUNT="${AMOUNT:-1000000000}"
COSMOS_NATIVE_DENOM="${COSMOS_NATIVE_DENOM:-stake}"

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
  if [ -z "${ERC20_ADDRESS:-}" ]; then
    ERC20_ADDRESS="$(config_value erc20_address)"
  fi
  if [ -z "${ERC20_ADDRESS:-}" ]; then
    ERC20_ADDRESS="$(config_value erc20)"
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
    if [ -z "${ICS26_ADDRESS:-}" ]; then
      ICS26_ADDRESS="$(jq -er '(.ics26Router // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .ics26Router))' "$BROADCAST_JSON" 2>/dev/null || true)"
    fi
  fi

  if [ -z "${ICS20_ADDRESS:-}" ] && [ -n "${ICS26_ADDRESS:-}" ] && contract_has_code "$ICS26_ADDRESS"; then
    ICS20_ADDRESS="$(cast call "$ICS26_ADDRESS" 'getIBCApp(string)(address)' "$DEST_PORT" --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
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
if [ -z "${ERC20_ADDRESS:-}" ] || [ -z "${ICS20_ADDRESS:-}" ]; then
  log_err "ERC20_ADDRESS or ICS20_ADDRESS not found"
  exit 1
fi
if [ ! -x "$SEND_PACKET_BIN" ]; then
  log "Building send-packet helper..."
  (cd "$REPO_ROOT/relayer" && go build -o "$SEND_PACKET_BIN" ./cmd/sendpacket)
fi

relayer_ensure_running

IBCERC20_DENOM_PATH="transfer/${SOURCE_CLIENT}/${COSMOS_NATIVE_DENOM}"

# ─── Phase 1: Establish ibcERC20 voucher (Cosmos → ETH) ───
log_header "Phase 1: Establish ibcERC20 (Cosmos → ETH)"
IBCERC20_ADDRESS="$(get_ibc_erc20_address "$ICS20_ADDRESS" "$IBCERC20_DENOM_PATH")"
if [ -n "$IBCERC20_ADDRESS" ] && [ "$IBCERC20_ADDRESS" != "0x0000000000000000000000000000000000000000" ]; then
  BEFORE_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
  log_kv "ibcERC20 balance (before)" "${BEFORE_IBCERC20_BALANCE:-0}"
else
  BEFORE_IBCERC20_BALANCE="0"
  log "ibcERC20 not yet created (will be minted on first transfer)"
fi

TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-1800}"
SEND_MEMO="create-voucher"

log "Sending native Cosmos coin to Ethereum..."
COSMOS_SEND_TX_HASH="$("$SEND_PACKET_BIN" \
  --receiver "$(echo "$ETH_SENDER" | tr '[:upper:]' '[:lower:]')" \
  --amount "$AMOUNT" \
  --denom "$COSMOS_NATIVE_DENOM" \
  --timeout-seconds "$TIMEOUT_SECONDS" \
  --client-id "$COSMOS_WASM_CLIENT_ID" \
  --memo "$SEND_MEMO" \
  --node "$COSMOS_RPC_URL" \
  --chain-id "$COSMOS_CHAIN_ID" \
  --count 1)"

log_ok "Native coin transfer submitted: $COSMOS_SEND_TX_HASH"
log "Waiting for relay to Ethereum..."

poll=0
MAX_POLLS=120
POLL_INTERVAL=15
VOUCHER_ESTABLISHED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  IBCERC20_ADDRESS="$(get_ibc_erc20_address "$ICS20_ADDRESS" "$IBCERC20_DENOM_PATH")"
  if [ -n "$IBCERC20_ADDRESS" ] && [ "$IBCERC20_ADDRESS" != "0x0000000000000000000000000000000000000000" ]; then
    CURRENT_IBCERC20="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
    if bc_ge "$CURRENT_IBCERC20" "$AMOUNT"; then
      log_ok "[poll ${poll}/${MAX_POLLS}] ibcERC20 minted │ balance=${CURRENT_IBCERC20}"
      VOUCHER_ESTABLISHED=true
      break
    fi
    log "[poll ${poll}/${MAX_POLLS}] ibcERC20 balance=${CURRENT_IBCERC20} (target=$AMOUNT)"
  else
    log "[poll ${poll}/${MAX_POLLS}] ibcERC20 contract not yet created..."
  fi
done

if [ "${VOUCHER_ESTABLISHED:-false}" != "true" ]; then
  log_err "ibcERC20 not minted after ${MAX_POLLS} polls"
  exit 1
fi

# ─── Phase 2: Send ibcERC20 back with timeout (ETH → Cosmos) ───
log_header "Phase 2: ibcERC20 ETH → Cosmos (timeout)"

RETURN_TIMEOUT_SECONDS="${RETURN_TIMEOUT_SECONDS:-30}"
RETURN_TIMEOUT=$(($(date +%s) + RETURN_TIMEOUT_SECONDS))

log "Approving ibcERC20 allowance..."
eth_approve "$IBCERC20_ADDRESS" "$ICS20_ADDRESS" "$AMOUNT"

log "Submitting sendTransfer of ibcERC20 back to Cosmos..."
eth_send_transfer "$ICS20_ADDRESS" "$IBCERC20_ADDRESS" "$AMOUNT" "$COSMOS_SENDER_ADDRESS" "$SOURCE_CLIENT" "$DEST_PORT" "$RETURN_TIMEOUT"
log "Return transfer submitted (waiting for timeout...)"

# ─── Poll for ibcERC20 remint on ETH ───
REFUND_THRESHOLD="$(bc_sub1 "$AMOUNT")"
poll=0
REFUNDED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_IBCERC20="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"

  if bc_ge "$CURRENT_IBCERC20" "$REFUND_THRESHOLD"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] ibcERC20 reminted │ balance=${CURRENT_IBCERC20} (target ≥ ${REFUND_THRESHOLD})"
    REFUNDED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] ibcERC20=${CURRENT_IBCERC20} (target ≥ ${REFUND_THRESHOLD})"
done

# ─── Final state ───
log_header "Final State"
AFTER_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
log_kv "ibcERC20 balance" "$AFTER_IBCERC20_BALANCE"

ICS20_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ICS20_ADDRESS")"
log_kv "ibcERC20 in ICS20" "$ICS20_IBCERC20_BALANCE"

# ─── Summary ───
log_header "Voucher Timeout (ETH Remint) Test Summary"
log_kv "Direction" "ETH → Cosmos (timeout)"
log_kv "Amount" "${AMOUNT} ${COSMOS_NATIVE_DENOM}"
log_kv "Timeout" "${RETURN_TIMEOUT_SECONDS}s"

if [ "${REFUNDED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • Native Cosmos→ETH transfer established ibcERC20"
  log "  • ibcERC20 sent back ETH→Cosmos with short timeout"
  log "  • Relayer detected timeout and submitted timeoutPacket on ETH"
  log "  • ibcERC20 balance restored on ETH (voucher reminted)"
else
  log_err "RESULT: FAIL"
  log_err "  • ibcERC20 balance not restored after ${MAX_POLLS} polls"
  log_err "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  exit 1
fi
