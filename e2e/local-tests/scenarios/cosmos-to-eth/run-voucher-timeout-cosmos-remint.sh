#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Voucher (originally ERC20) sent from Cosmos→ETH times out
# → voucher reminted on Cosmos.
#
# Corresponds to: Test_TimeoutPacketCosmosRemintsVouchers

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
COSMOS_RECEIVER_KEY="${COSMOS_RECEIVER_KEY:-test}"
COSMOS_SIGNER_KEY="${COSMOS_SIGNER_KEY:-test1}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"
SOURCE_CLIENT="${SOURCE_CLIENT:-cosmoshub-1}"
DEST_PORT="${DEST_PORT:-transfer}"
AMOUNT="${AMOUNT:-1000000000}"

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

  if [ -z "${ERC20_ADDRESS:-}" ]; then
    local mint_receiver
    mint_receiver="${E2E_FAUCET_ADDRESS:-}"
    if [ -z "$mint_receiver" ] && [ -n "${ETH_PRIVATE_KEY:-}" ]; then
      mint_receiver="$(cast wallet address --private-key "$ETH_PRIVATE_KEY" 2>/dev/null || true)"
    fi
    if [ -n "$mint_receiver" ]; then
      ERC20_ADDRESS="$(cast logs \
        --from-block 0 \
        --to-block latest \
        --rpc-url "$ETH_RPC_URL" \
        'Transfer(address indexed from,address indexed to,uint256 value)' \
        0x0000000000000000000000000000000000000000 \
        "$mint_receiver" 2>/dev/null \
        | awk '/^- address:/ { print $3; exit }')"
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

RECEIVER="${RECEIVER:-}"
load_cosmos_receiver

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
if [ -z "${RECEIVER:-}" ]; then
  log_err "RECEIVER not set"
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

VOUCHER_TRACE="transfer/$COSMOS_WASM_CLIENT_ID/$ERC20_ADDRESS"
VOUCHER_DENOM="$(derive_ibc_denom "$VOUCHER_TRACE")"

# ─── Phase 1: Establish voucher (ETH → Cosmos) ───
log_header "Phase 1: Establish Voucher (ETH → Cosmos)"
BEFORE_VOUCHER_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos voucher balance (before)" "$BEFORE_VOUCHER_BALANCE"

TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-1800}"
TIMEOUT=$(($(date +%s) + TIMEOUT_SECONDS))

log "Approving ERC20 allowance..."
eth_approve "$ERC20_ADDRESS" "$ICS20_ADDRESS" "$AMOUNT"

log "Submitting sendTransfer..."
eth_send_transfer "$ICS20_ADDRESS" "$ERC20_ADDRESS" "$AMOUNT" "$RECEIVER" "$SOURCE_CLIENT" "$DEST_PORT" "$TIMEOUT"
log "Transfer submitted (waiting for relay to Cosmos...)"

TARGET_VOUCHER_BALANCE="$(echo "$BEFORE_VOUCHER_BALANCE + $AMOUNT" | bc)"
log "Polling Cosmos for voucher balance (target ≥ $TARGET_VOUCHER_BALANCE)..."

poll=0
MAX_POLLS=120
POLL_INTERVAL=15
VOUCHER_ESTABLISHED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_VOUCHER="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

  if bc_ge "$CURRENT_VOUCHER" "$TARGET_VOUCHER_BALANCE"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] Voucher established │ balance=${CURRENT_VOUCHER}"
    VOUCHER_ESTABLISHED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] voucher=${CURRENT_VOUCHER} (target ≥ ${TARGET_VOUCHER_BALANCE})"
done

if [ "${VOUCHER_ESTABLISHED:-false}" != "true" ]; then
  log_err "Voucher not established after ${MAX_POLLS} polls"
  exit 1
fi

# ─── Phase 2: Send voucher back with timeout (Cosmos → ETH) ───
log_header "Phase 2: Voucher Cosmos → ETH (timeout)"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-45}"
SEND_MEMO="voucher-timeout-cosmos"

COSMOS_SEND_TX_HASH="$("$SEND_PACKET_BIN" \
  --receiver "$(echo "$ETH_SENDER" | tr '[:upper:]' '[:lower:]')" \
  --amount "$AMOUNT" \
  --denom "$VOUCHER_TRACE" \
  --timeout-seconds "$TIMEOUT_SECONDS" \
  --client-id "$COSMOS_WASM_CLIENT_ID" \
  --memo "$SEND_MEMO" \
  --node "$COSMOS_RPC_URL" \
  --chain-id "$COSMOS_CHAIN_ID" \
  --count 1)"

log_ok "Voucher return transfer submitted: $COSMOS_SEND_TX_HASH"
log "Waiting for timeout..."

# ─── Poll for voucher remint on Cosmos ───
REFUND_THRESHOLD="$(bc_sub1 "$TARGET_VOUCHER_BALANCE")"
poll=0
REFUNDED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_VOUCHER="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

  if bc_ge "$CURRENT_VOUCHER" "$REFUND_THRESHOLD"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] Voucher reminted │ balance=${CURRENT_VOUCHER} (target ≥ ${REFUND_THRESHOLD})"
    REFUNDED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] voucher=${CURRENT_VOUCHER} (target ≥ ${REFUND_THRESHOLD})"
done

# ─── Final state ───
log_header "Final State"
AFTER_VOUCHER_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos voucher" "$AFTER_VOUCHER_BALANCE"

# ─── Summary ───
log_header "Voucher Timeout (Cosmos Remint) Test Summary"
log_kv "Direction" "Cosmos → ETH (timeout)"
log_kv "Amount" "${AMOUNT} wei"
log_kv "Timeout" "${TIMEOUT_SECONDS}s"

if [ "${REFUNDED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • ERC20→Cosmos voucher established"
  log "  • Voucher sent back Cosmos→ETH with short timeout"
  log "  • Relayer detected timeout and submitted timeoutPacket on Cosmos"
  log "  • Voucher balance restored on Cosmos"
else
  log_err "RESULT: FAIL"
  log_err "  • Voucher balance not restored after ${MAX_POLLS} polls"
  log_err "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  exit 1
fi
