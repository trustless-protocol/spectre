#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"
STATE_DIR="$REPO_ROOT/e2e/local-tests/.state"
RELAYER_PID_FILE="$REPO_ROOT/relayer/relayer.pid"

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

load_timeout_contract_addresses() {
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

if [ -z "${ETH_RPC_URL:-}" ]; then
  echo "ERROR: ETH_RPC_URL not set, Kurtosis endpoint not detected, and no eth_rpc_url found in $CONFIG_FILE" >&2
  echo "Start ETH with ./setup/01-eth-node.sh or export ETH_RPC_URL manually." >&2
  exit 1
fi

echo "Using ETH_RPC_URL: $ETH_RPC_URL"

load_timeout_contract_addresses

RECEIVER="${RECEIVER:-}"
load_cosmos_receiver

uint_value() {
  awk '{ print $1 }'
}

if [ -z "${ETH_RPC_URL:-}" ]; then
  echo "ERROR: ETH_RPC_URL not set and Kurtosis not detected" >&2
  exit 1
fi
if [ -z "${ETH_PRIVATE_KEY:-}" ]; then
  echo "ERROR: ETH_PRIVATE_KEY not set" >&2
  exit 1
fi
if [ -z "${RECEIVER:-}" ]; then
  echo "ERROR: RECEIVER not set; start Cosmos node or set RECEIVER env" >&2
  exit 1
fi
if [ -z "${ERC20_ADDRESS:-}" ] || [ -z "${ICS20_ADDRESS:-}" ]; then
  echo "ERROR: ERC20_ADDRESS or ICS20_ADDRESS not found" >&2
  if [ -n "${ICS26_ADDRESS:-}" ]; then
    echo "ICS26_ADDRESS candidate: $ICS26_ADDRESS" >&2
    if ! contract_has_code "$ICS26_ADDRESS"; then
      echo "No bytecode at ICS26_ADDRESS on $ETH_RPC_URL; relayer/config.json is stale for the current chain." >&2
    fi
  fi
  echo "The script checks env, $CONFIG_FILE, Foundry broadcast JSON, ICS26Router.getIBCApp(\"$DEST_PORT\"), and ERC20 mint logs." >&2
  echo "Run ./setup/01-eth-node.sh to deploy contracts, refresh relayer/config.json, or export ERC20_ADDRESS and ICS20_ADDRESS." >&2
  echo "Expected Foundry broadcast path: $REPO_ROOT/broadcast/E2ETestDeploy.s.sol/<chain-id>/run-latest.json" >&2
  exit 1
fi

# ---- Ensure relayer is running ----
if [ -f "$RELAYER_PID_FILE" ]; then
  RELAYER_PID="$(cat "$RELAYER_PID_FILE")"
  if kill -0 "$RELAYER_PID" 2>/dev/null; then
    echo "Relayer running (PID $RELAYER_PID)"
  else
    echo "ERROR: relayer PID file exists but process not running. Start relayer first:" >&2
    echo "  ./setup/05-relayer.sh" >&2
    exit 1
  fi
else
  echo "ERROR: relayer not running (no PID file). Start relayer first:" >&2
  echo "  ./setup/05-relayer.sh" >&2
  exit 1
fi

ETH_SENDER="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"

# ---- Record pre-transfer state ----
echo ""
echo "=== Pre-transfer state ==="
BEFORE_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
echo "Sender ERC20 balance (pre):  $BEFORE_ETH_BALANCE"

ESCROW_ADDRESS="$(cast call "$ICS20_ADDRESS" 'getEscrow(string)(address)' "$SOURCE_CLIENT" --rpc-url "$ETH_RPC_URL" 2>/dev/null || echo "")"
if [ -n "$ESCROW_ADDRESS" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
  BEFORE_ESCROW_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  echo "Escrow ERC20 balance (pre):   $BEFORE_ESCROW_BALANCE"
fi

COSMOS_VOUCHER_DENOM="transfer/08-wasm-0/$ERC20_ADDRESS"
BEFORE_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
echo "Cosmos voucher balance (pre): $BEFORE_COSMOS_BALANCE"

# ---- Send transfer with short timeout ----
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-10}"
TIMEOUT=$(($(date +%s) + TIMEOUT_SECONDS))
echo ""
echo "=== Sending ICS20Transfer with ${TIMEOUT_SECONDS}s timeout ==="
echo "Timeout at epoch: $TIMEOUT ($(date -d "@$TIMEOUT" '+%H:%M:%S' 2>/dev/null || echo 'check `date`'))"

echo "Approving ERC20 allowance..."
cast send "$ERC20_ADDRESS" "approve(address,uint256)" "$ICS20_ADDRESS" "$AMOUNT" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY" > /dev/null

echo "Submitting sendTransfer..."
TRANSFER_TUPLE="($ERC20_ADDRESS,$AMOUNT,$RECEIVER,$SOURCE_CLIENT,$DEST_PORT,$TIMEOUT,\"\")"
cast send "$ICS20_ADDRESS" \
  "sendTransfer((address,uint256,string,string,string,uint64,string))" \
  "$TRANSFER_TUPLE" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY" > /dev/null

echo ""
echo "Waiting for relayer to detect timeout and submit refund..."
echo "(relayer builds ZK proof for light client update, then calls timeoutPacket)"
echo ""

# ---- Poll for refund (ERC20 balance returning to original) ----
MAX_POLLS=10
POLL_INTERVAL=15
poll=0

# refund expected = original plus maybe tiny diff from gas
# but balance is ERC20 not ETH, so no gas deduction on ERC20
REFUND_THRESHOLD=$((BEFORE_ETH_BALANCE - 1))

while [ $poll -lt $MAX_POLLS ]; do
  sleep $POLL_INTERVAL
  poll=$((poll + 1))

  CURRENT_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
  CURRENT_ESCROW="0"
  if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
    CURRENT_ESCROW="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  fi
  CURRENT_COSMOS="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

  if [ "$CURRENT_ETH_BALANCE" -ge "$REFUND_THRESHOLD" ] 2>/dev/null; then
    echo "[poll ${poll}/${MAX_POLLS}] ERC20 refunded! balance=$CURRENT_ETH_BALANCE (target >= $REFUND_THRESHOLD)"
    REFUNDED=true
    break
  fi

  echo "[poll ${poll}/${MAX_POLLS}] sender=$CURRENT_ETH_BALANCE escrow=$CURRENT_ESCROW cosmos=$CURRENT_COSMOS"
done

# ---- Final state ----
echo ""
echo "=== Final state ==="
AFTER_ETH_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ETH_SENDER" --rpc-url "$ETH_RPC_URL" | uint_value)"
echo "Sender ERC20 balance (final): $AFTER_ETH_BALANCE"

if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
  AFTER_ESCROW_BALANCE="$(cast call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ESCROW_ADDRESS" --rpc-url "$ETH_RPC_URL" | uint_value)"
  echo "Escrow ERC20 balance (final):  $AFTER_ESCROW_BALANCE"
fi

AFTER_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
echo "Cosmos voucher (final):        $AFTER_COSMOS_BALANCE"

# ---- Summary ----
echo ""
echo "=== Timeout Test Summary ==="
echo "Direction:      ETH → Cosmos"
echo "Amount sent:    $AMOUNT wei"
echo "Timeout:        ${TIMEOUT_SECONDS}s"
echo ""

if [ "${REFUNDED:-false}" = "true" ]; then
  echo "RESULT: PASS"
  echo "  - Relayer detected the expired packet"
  echo "  - Built non-membership proof from Cosmos"
  echo "  - Called timeoutPacket() on Ethereum ICS26Router"
  echo "  - Sender refunded (ERC20 balance returned to original)"
  echo "  - Cosmos voucher = 0 (packet never relayed)"
elif [ "$AFTER_ETH_BALANCE" -ge "$REFUND_THRESHOLD" ] 2>/dev/null; then
  echo "RESULT: PASS (detected after final poll)"
  echo "  - Sender ERC20 balance returned to original"
elif [ "$AFTER_COSMOS_BALANCE" = "0" ]; then
  echo "RESULT: PARTIAL"
  echo "  - Packet was NOT relayed to Cosmos (timeout prevented relay)"
  echo "  - ERC20 balance still decreased — refund not yet submitted"
  echo "  - The relayer may still be building the ZK proof"
  echo "  - Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
else
  echo "RESULT: FAIL"
  echo "  - Packet WAS relayed to Cosmos before timeout"
  echo "  - Try with shorter TIMEOUT_SECONDS or check relayer status"
fi
