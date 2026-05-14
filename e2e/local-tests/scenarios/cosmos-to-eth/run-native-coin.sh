#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Native Cosmos coin → ETH (minting ibcERC20), then return to Cosmos.
#
# Corresponds to: Test_ICS20TransferNativeCosmosCoinsToEthereumAndBack

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

# ─── Record pre-transfer state ───
log_header "Pre-transfer State"
BEFORE_NATIVE_BALANCE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos native balance" "$BEFORE_NATIVE_BALANCE"

IBCERC20_DENOM_PATH="${COSMOS_NATIVE_DENOM}/transfer/${SOURCE_CLIENT}"
IBCERC20_ADDRESS="$(get_ibc_erc20_address "$ICS20_ADDRESS" "$IBCERC20_DENOM_PATH")"
if [ -n "$IBCERC20_ADDRESS" ] && [ "$IBCERC20_ADDRESS" != "0x0000000000000000000000000000000000000000" ]; then
  BEFORE_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
  log_kv "ibcERC20 balance" "${BEFORE_IBCERC20_BALANCE:-0}"
else
  log "ibcERC20 not yet created (will be minted on first transfer)"
  BEFORE_IBCERC20_BALANCE="0"
fi

# ─── Phase 1: Cosmos → ETH ───
log_header "Phase 1: Cosmos → ETH Native Coin Transfer"
log "Sending native coin to Ethereum..."

TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-1800}"
SEND_MEMO="native-coin-send"

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

# Poll for ibcERC20 creation and balance
poll=0
MAX_POLLS=120
POLL_INTERVAL=15
IBCERC20_FOUND=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  IBCERC20_ADDRESS="$(get_ibc_erc20_address "$ICS20_ADDRESS" "$IBCERC20_DENOM_PATH")"
  if [ -n "$IBCERC20_ADDRESS" ] && [ "$IBCERC20_ADDRESS" != "0x0000000000000000000000000000000000000000" ]; then
    CURRENT_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
    if [ "$CURRENT_IBCERC20_BALANCE" = "$AMOUNT" ]; then
      log_ok "[poll ${poll}/${MAX_POLLS}] ibcERC20 minted │ balance=${CURRENT_IBCERC20_BALANCE}"
      IBCERC20_FOUND=true
      break
    fi
    log "[poll ${poll}/${MAX_POLLS}] ibcERC20 balance=${CURRENT_IBCERC20_BALANCE} (target=$AMOUNT)"
  else
    log "[poll ${poll}/${MAX_POLLS}] ibcERC20 contract not yet created..."
  fi
done

if [ "${IBCERC20_FOUND:-false}" != "true" ]; then
  log_err "ibcERC20 not minted after ${MAX_POLLS} polls"
  log_err "Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  exit 1
fi

# Verify ibcERC20 metadata
IBCERC20_NAME="$(safe_cast_call "$IBCERC20_ADDRESS" 'name()(string)')"
IBCERC20_SYMBOL="$(safe_cast_call "$IBCERC20_ADDRESS" 'symbol()(string)')"
IBCERC20_FULL_DENOM="$(safe_cast_call "$IBCERC20_ADDRESS" 'fullDenomPath()(string)')"
log_kv "ibcERC20 address" "$IBCERC20_ADDRESS"
log_kv "ibcERC20 name" "$IBCERC20_NAME"
log_kv "ibcERC20 symbol" "$IBCERC20_SYMBOL"
log_kv "ibcERC20 fullDenomPath" "$IBCERC20_FULL_DENOM"

if [ "$IBCERC20_FULL_DENOM" != "$IBCERC20_DENOM_PATH" ]; then
  log_warn "ibcERC20 fullDenomPath ($IBCERC20_FULL_DENOM) does not match expected ($IBCERC20_DENOM_PATH)"
fi

# ─── Phase 2: ETH → Cosmos return transfer ───
log_header "Phase 2: ETH → Cosmos Return Transfer"

RETURN_TIMEOUT_SECONDS="${RETURN_TIMEOUT_SECONDS:-1800}"
RETURN_TIMEOUT=$(($(date +%s) + RETURN_TIMEOUT_SECONDS))

log "Approving ibcERC20 allowance..."
eth_approve "$IBCERC20_ADDRESS" "$ICS20_ADDRESS" "$AMOUNT"

log "Submitting sendTransfer of ibcERC20 back to Cosmos..."
eth_send_transfer "$ICS20_ADDRESS" "$IBCERC20_ADDRESS" "$AMOUNT" "$COSMOS_SENDER_ADDRESS" "$SOURCE_CLIENT" "$DEST_PORT" "$RETURN_TIMEOUT"
log "Return transfer submitted (waiting for relay...)"

# Poll for native balance restoration on Cosmos
TARGET_NATIVE_BALANCE="$BEFORE_NATIVE_BALANCE"
REFUND_THRESHOLD="$(bc_sub1 "$TARGET_NATIVE_BALANCE")"
poll=0
REFUNDED=false

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_NATIVE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
  CURRENT_IBCERC20="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"

  if bc_ge "$CURRENT_NATIVE" "$REFUND_THRESHOLD"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] Native balance restored │ balance=${CURRENT_NATIVE} (target ≥ ${REFUND_THRESHOLD})"
    REFUNDED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] native=${CURRENT_NATIVE}  ibcERC20=${CURRENT_IBCERC20}"
done

# ─── Final state ───
log_header "Final State"
AFTER_NATIVE_BALANCE="$("$COSMOS_BIN" query bank balance "$COSMOS_SENDER_ADDRESS" "$COSMOS_NATIVE_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos native balance" "$AFTER_NATIVE_BALANCE"

AFTER_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ETH_SENDER")"
log_kv "ibcERC20 balance" "$AFTER_IBCERC20_BALANCE"

ICS20_IBCERC20_BALANCE="$(eth_balance_of "$IBCERC20_ADDRESS" "$ICS20_ADDRESS")"
log_kv "ibcERC20 balance in ICS20" "$ICS20_IBCERC20_BALANCE"

# ─── Summary ───
log_header "Native Coin Roundtrip Test Summary"
log_kv "Direction" "Cosmos → ETH → Cosmos"
log_kv "Amount" "${AMOUNT} ${COSMOS_NATIVE_DENOM}"

if [ "${REFUNDED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • Cosmos→ETH transfer relayed successfully (ibcERC20 minted)"
  log "  • ETH→Cosmos return transfer relayed successfully"
  log "  • Cosmos native balance restored"
  log "  • ibcERC20 balance on ETH: $AFTER_IBCERC20_BALANCE"
else
  log_err "RESULT: FAIL"
  log_err "  • Native balance not restored after ${MAX_POLLS} polls"
  log_err "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  exit 1
fi
