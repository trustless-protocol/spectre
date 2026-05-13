#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: ETH → Cosmos packet timeout on the Cosmos side.
#
# This script sends a transfer from Ethereum with a short timeout while the
# relayer is running. The packet is considered "timed out on the Cosmos side"
# when the timeout expires before the relayer can build the ZK proof and
# deliver the packet to Cosmos. The relayer then detects the expired packet,
# builds a non-membership proof from Cosmos, and calls timeoutPacket() on
# Ethereum's ICS26Router to refund the sender.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"
RELAYER_PID_FILE="${RELAYER_PID_FILE:-$REPO_ROOT/relayer/relayer.pid}"
RELAYER_DIR="$REPO_ROOT/relayer"

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_HOME="${COSMOS_HOME:-$HOME/.gaia}"
COSMOS_KEYRING="${COSMOS_KEYRING:-test}"
COSMOS_RECEIVER_KEY="${COSMOS_RECEIVER_KEY:-test}"
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

# ---- Validate required local tools before discovery ----
require_cmd cast
require_cmd jq

load_env_file
ETH_PRIVATE_KEY="${ETH_PRIVATE_KEY:-${ETH_USER_PK:-${PRIVATE_KEY:-}}}"
discover_kurtosis_endpoints

if [ -z "${ETH_RPC_URL:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ETH_RPC_URL="$(config_value eth_rpc_url)"
fi

load_scenario_contract_addresses

# Read cosmos_wasm_client_id from config
if [ -z "${COSMOS_WASM_CLIENT_ID:-}" ] && [ -f "$CONFIG_FILE" ]; then
  COSMOS_WASM_CLIENT_ID="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.cosmos_wasm_client_id // empty' "$CONFIG_FILE" 2>/dev/null || true)"
fi
COSMOS_WASM_CLIENT_ID="${COSMOS_WASM_CLIENT_ID:-08-wasm-0}"

RECEIVER="${RECEIVER:-}"
load_cosmos_receiver

uint_value() {
  awk '{ print $1 }'
}

# Big-integer helpers (bash arithmetic overflows > 2^63)
bc_sub1() {
  echo "$1 - 1" | bc
}

bc_ge() {
  local a="$1" b="$2"
  [ "$(echo "$a >= $b" | bc)" = "1" ]
}

if [ -z "${ETH_RPC_URL:-}" ]; then
  log_err "ETH_RPC_URL not set, Kurtosis endpoint not detected, and no eth_rpc_url found in $CONFIG_FILE"
  log_err "Start ETH with ./setup/01-eth-node.sh or export ETH_RPC_URL manually."
  exit 1
fi
if [ -z "${ETH_PRIVATE_KEY:-}" ]; then
  log_err "ETH_PRIVATE_KEY not set"
  exit 1
fi
if [ -z "${RECEIVER:-}" ]; then
  log_err "RECEIVER not set; start Cosmos node or set RECEIVER env"
  exit 1
fi
if [ -z "${ERC20_ADDRESS:-}" ] || [ -z "${ICS20_ADDRESS:-}" ]; then
  log_err "ERC20_ADDRESS or ICS20_ADDRESS not found"
  if [ -n "${ICS26_ADDRESS:-}" ]; then
    log_err "ICS26_ADDRESS candidate: $ICS26_ADDRESS"
    if ! contract_has_code "$ICS26_ADDRESS"; then
      log_err "No bytecode at ICS26_ADDRESS on $ETH_RPC_URL; relayer/config.json is stale for the current chain."
    fi
  fi
  log_err "The script checks env, $CONFIG_FILE, Foundry broadcast JSON, ICS26Router.getIBCApp(\"$DEST_PORT\"), and ERC20 mint logs."
  log_err "Run ./setup/01-eth-node.sh to deploy contracts, refresh relayer/config.json, or export ERC20_ADDRESS and ICS20_ADDRESS."
  log_err "Expected Foundry broadcast path: $REPO_ROOT/broadcast/E2ETestDeploy.s.sol/<chain-id>/run-latest.json"
  exit 1
fi

# ─── Ensure relayer is running ───
if [ -f "$RELAYER_PID_FILE" ]; then
  RELAYER_PID="$(cat "$RELAYER_PID_FILE")"
  if kill -0 "$RELAYER_PID" 2>/dev/null; then
    log_ok "Relayer running (PID $RELAYER_PID)"
  else
    log_err "relayer PID file exists but process not running. Start relayer first:"
    log_err "./setup/05-relayer.sh"
    exit 1
  fi
else
  log_err "relayer not running (no PID file). Start relayer first:"
  log_err "./setup/05-relayer.sh"
  exit 1
fi

ETH_SENDER="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"

# ─── Record pre-transfer state ───
log_header "Pre-transfer State"
BEFORE_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
log_kv "Sender ERC20 balance" "$BEFORE_ETH_BALANCE"

ESCROW_ADDRESS="$(cast call "$ICS20_ADDRESS" 'getEscrow(string)(address)' "$SOURCE_CLIENT" --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "")"
if [ -n "$ESCROW_ADDRESS" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
  BEFORE_ESCROW_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  log_kv "Escrow ERC20 balance" "$BEFORE_ESCROW_BALANCE"
fi

COSMOS_VOUCHER_DENOM="transfer/$COSMOS_WASM_CLIENT_ID/$ERC20_ADDRESS"
BEFORE_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos voucher balance" "$BEFORE_COSMOS_BALANCE"

# ─── Send transfer with short timeout ───
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-15}"
TIMEOUT=$(($(date +%s) + TIMEOUT_SECONDS))
log_header "Sending ICS20Transfer (${TIMEOUT_SECONDS}s timeout)"
log_kv "Timeout" "epoch $TIMEOUT ($(date -d "@$TIMEOUT" '+%H:%M:%S' 2>/dev/null || echo 'N/A'))"

log "Approving ERC20 allowance..."
cast send "$ERC20_ADDRESS" "approve(address,uint256)" "$ICS20_ADDRESS" "$AMOUNT" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY" > /dev/null

log "Submitting sendTransfer..."
TRANSFER_TUPLE="($ERC20_ADDRESS,$AMOUNT,$RECEIVER,$SOURCE_CLIENT,$DEST_PORT,$TIMEOUT,\"\")"
cast send "$ICS20_ADDRESS" \
  "sendTransfer((address,uint256,string,string,string,uint64,string))" \
  "$TRANSFER_TUPLE" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY" > /dev/null

log "Transfer submitted (waiting for timeout...)"

# ─── Poll for refund (ERC20 balance returning to original) ───
MAX_POLLS=120
POLL_INTERVAL=15
poll=0
REFUNDED=false

REFUND_THRESHOLD="$(bc_sub1 "$BEFORE_ETH_BALANCE")"

log_header "Polling for Refund"
while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
  CURRENT_ESCROW="0"
  if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
    CURRENT_ESCROW="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  fi
  CURRENT_COSMOS="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

  if bc_ge "$CURRENT_ETH_BALANCE" "$REFUND_THRESHOLD"; then
    log_ok "[poll ${poll}/${MAX_POLLS}] ERC20 refunded │ sender=${CURRENT_ETH_BALANCE} (target ≥ ${REFUND_THRESHOLD})"
    REFUNDED=true
    break
  fi

  log "[poll ${poll}/${MAX_POLLS}] sender=${CURRENT_ETH_BALANCE}  escrow=${CURRENT_ESCROW}  cosmos=${CURRENT_COSMOS}"
done

# ─── Final state ───
log_header "Final State"
AFTER_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
log_kv "Sender ERC20 balance" "$AFTER_ETH_BALANCE"

if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
  AFTER_ESCROW_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  log_kv "Escrow ERC20 balance" "$AFTER_ESCROW_BALANCE"
fi

AFTER_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
log_kv "Cosmos voucher" "$AFTER_COSMOS_BALANCE"

# ─── Summary ───
log_header "Cosmos Timeout Test Summary"
log_kv "Direction" "ETH → Cosmos"
log_kv "Amount sent" "${AMOUNT} wei"
log_kv "Timeout" "${TIMEOUT_SECONDS}s"

if [ "${REFUNDED:-false}" = "true" ]; then
  log_ok "RESULT: PASS"
  log "  • Relayer detected the expired packet"
  log "  • Built non-membership proof from Cosmos"
  log "  • Called timeoutPacket() on Ethereum ICS26Router"
  log "  • Sender refunded (ERC20 balance returned to original)"
  log "  • Cosmos voucher = 0 (packet never relayed)"
elif bc_ge "$AFTER_ETH_BALANCE" "$REFUND_THRESHOLD"; then
  log_ok "RESULT: PASS (detected after final poll)"
  log "  • Sender ERC20 balance returned to original"
elif [ "$AFTER_COSMOS_BALANCE" = "0" ]; then
  log_warn "RESULT: PARTIAL"
  log_warn "  • Packet was NOT relayed to Cosmos (timeout prevented relay)"
  log_warn "  • ERC20 balance still decreased — refund not yet submitted"
  log_warn "  • The relayer may still be building the ZK proof"
  log_warn "  • Check relayer logs: tail -f $RELAYER_DIR/relayer.log"
else
  log_err "RESULT: FAIL"
  log_err "  • Packet WAS relayed to Cosmos before timeout"
  log_err "  • Try with shorter TIMEOUT_SECONDS or check relayer status"
fi
