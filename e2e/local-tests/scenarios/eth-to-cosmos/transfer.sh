#!/usr/bin/env bash
set -euo pipefail

# Local E2E helper: approve the deployed ICS20Transfer contract and send an
# ERC20 packet from the local Kurtosis Ethereum chain to the local Cosmos chain.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"
BROADCAST_JSON="${BROADCAST_JSON:-}"

if [ -f "$ENV_FILE" ]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
COSMOS_BIN="${COSMOS_BIN:-gaiad}"
COSMOS_HOME="${COSMOS_HOME:-$HOME/.gaia}"
COSMOS_KEYRING="${COSMOS_KEYRING:-test}"
COSMOS_RECEIVER_KEY="${COSMOS_RECEIVER_KEY:-test}"
COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
COSMOS_GRPC_ADDR="${COSMOS_GRPC_ADDR:-127.0.0.1:9090}"
COSMOS_CHAIN_ID="${COSMOS_CHAIN_ID:-test-ibc-eth}"

ETH_RPC_URL="${ETH_RPC_URL:-${ETH_RPC:-}}"
ETH_BEACON_API_URL="${ETH_BEACON_API_URL:-}"
ETH_PRIVATE_KEY="${ETH_PRIVATE_KEY:-${ETH_USER_PK:-${PRIVATE_KEY:-}}}"
AMOUNT="${AMOUNT:-${TRANSFER_AMOUNT:-1000000000}}"
SOURCE_CLIENT="${SOURCE_CLIENT:-${ETH_SOURCE_CLIENT_ID:-cosmoshub-1}}"
DEST_PORT="${DEST_PORT:-${IBC_TRANSFER_PORT:-transfer}}"
MEMO="${MEMO:-}"
TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-1800}"
TIMEOUT="${TIMEOUT:-$(($(date +%s) + TIMEOUT_SECONDS))}"
RECEIVER="${RECEIVER:-${COSMOS_RECEIVER_ADDRESS:-}}"

# Auto-detect RECEIVER from Cosmos key if not set
if [ -z "$RECEIVER" ] && command -v "$COSMOS_BIN" >/dev/null 2>&1; then
  RECEIVER="$("$COSMOS_BIN" keys show "$COSMOS_RECEIVER_KEY" -a \
    --keyring-backend "$COSMOS_KEYRING" \
    --home "$COSMOS_HOME" 2>/dev/null || true)"
fi

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd" >&2
    exit 1
  fi
}

discover_kurtosis_endpoints() {
  if ! command -v kurtosis >/dev/null 2>&1; then
    return 0
  fi

  local inspect
  if ! inspect="$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" 2>/dev/null)"; then
    return 0
  fi

  if [ -z "$ETH_RPC_URL" ]; then
    ETH_RPC_URL="$(printf '%s\n' "$inspect" | perl -ne '
      if (/el-1-geth-lighthouse/) { $in=1 }
      elsif ($in && /^\S/) { $in=0 }
      elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
        print "http://$1\n";
        exit;
      }
    ')"
  fi

  if [ -z "$ETH_BEACON_API_URL" ]; then
    ETH_BEACON_API_URL="$(printf '%s\n' "$inspect" | awk '
      $0 ~ /cl-1-lighthouse-geth/ { in_service = 1 }
      in_service && /http:/ {
        match($0, /127\.0\.0\.1:[0-9]+/)
        print "http://" substr($0, RSTART, RLENGTH)
        exit
      }
      in_service && /^[^[:space:]]/ { in_service = 0 }
    ')"
  fi
}

load_contract_addresses() {
  if [ -z "$BROADCAST_JSON" ]; then
    local chain_id chain_broadcast
    chain_id="$(cast chain-id --rpc-url "$ETH_RPC_URL")"
    chain_broadcast="$REPO_ROOT/broadcast/E2ETestDeploy.s.sol/$chain_id/run-latest.json"
    if [ -f "$chain_broadcast" ]; then
      BROADCAST_JSON="$chain_broadcast"
    else
      BROADCAST_JSON="$(find "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" -path '*/run-latest.json' -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | awk 'NR == 1 { print $2 }')"
    fi
  fi

  if [ ! -f "$BROADCAST_JSON" ]; then
    return 0
  fi

  if [ -z "${ERC20_ADDRESS:-}" ]; then
    ERC20_ADDRESS="$(jq -er '(.erc20 // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .erc20))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi

  if [ -z "${ICS20_ADDRESS:-}" ]; then
    ICS20_ADDRESS="$(jq -er '(.ics20Transfer // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .ics20Transfer))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi
}

load_cosmos_receiver() {
  if [ -n "$RECEIVER" ]; then
    return 0
  fi

  if command -v "$COSMOS_BIN" >/dev/null 2>&1; then
    RECEIVER="$("$COSMOS_BIN" keys show "$COSMOS_RECEIVER_KEY" -a \
      --keyring-backend "$COSMOS_KEYRING" \
      --home "$COSMOS_HOME" 2>/dev/null || true)"
  fi
}

require_cmd jq
require_cmd cast

discover_kurtosis_endpoints
load_contract_addresses
load_cosmos_receiver

if [ -z "$ETH_RPC_URL" ]; then
  cat >&2 <<EOF
missing ETH_RPC_URL.
Set ETH_RPC_URL in $ENV_FILE or start Kurtosis with ./run_eth_node.sh.
EOF
  exit 1
fi

if [ -z "$ETH_PRIVATE_KEY" ]; then
  cat >&2 <<EOF
missing ETH_PRIVATE_KEY.
Set ETH_PRIVATE_KEY in $ENV_FILE.
EOF
  exit 1
fi

if [ -z "${ERC20_ADDRESS:-}" ] || [ -z "${ICS20_ADDRESS:-}" ]; then
  cat >&2 <<EOF
missing ERC20_ADDRESS or ICS20_ADDRESS.
Set them in the environment, or make sure this deployment file exists:
  $BROADCAST_JSON
EOF
  exit 1
fi

if [ -z "$RECEIVER" ]; then
  cat >&2 <<EOF
missing RECEIVER.
Set RECEIVER=<cosmos-bech32-address>, or start Cosmos with ./run_cosmos_node.sh
so the script can read key '$COSMOS_RECEIVER_KEY' from $COSMOS_HOME.
EOF
  exit 1
fi

ETH_CHAIN_ID="$(cast chain-id --rpc-url "$ETH_RPC_URL")"
ETH_SENDER="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"

printf 'ENV_FILE=%s\n' "$ENV_FILE"
printf 'KURTOSIS_ENCLAVE=%s\n' "$KURTOSIS_ENCLAVE"
printf 'ETH_RPC_URL=%s\n' "$ETH_RPC_URL"
printf 'ETH_BEACON_API_URL=%s\n' "${ETH_BEACON_API_URL:-}"
printf 'ETH_CHAIN_ID=%s\n' "$ETH_CHAIN_ID"
printf 'ETH_SENDER=%s\n' "$ETH_SENDER"
printf 'BROADCAST_JSON=%s\n' "${BROADCAST_JSON:-}"
printf 'COSMOS_RPC_URL=%s\n' "$COSMOS_RPC_URL"
printf 'COSMOS_GRPC_ADDR=%s\n' "$COSMOS_GRPC_ADDR"
printf 'COSMOS_CHAIN_ID=%s\n' "$COSMOS_CHAIN_ID"
printf 'ERC20_ADDRESS=%s\n' "$ERC20_ADDRESS"
printf 'ICS20_ADDRESS=%s\n' "$ICS20_ADDRESS"
printf 'RECEIVER=%s\n' "$RECEIVER"
printf 'AMOUNT=%s\n' "$AMOUNT"
printf 'SOURCE_CLIENT=%s\n' "$SOURCE_CLIENT"
printf 'DEST_PORT=%s\n' "$DEST_PORT"
printf 'TIMEOUT=%s\n' "$TIMEOUT"

echo
echo "Approving ERC20 to ICS20Transfer contract..."
cast send "$ERC20_ADDRESS" "approve(address,uint256)" "$ICS20_ADDRESS" "$AMOUNT" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY"

echo
echo "Submitting ICS20Transfer.sendTransfer..."
TRANSFER_TUPLE="($ERC20_ADDRESS,$AMOUNT,$RECEIVER,$SOURCE_CLIENT,$DEST_PORT,$TIMEOUT,\"$MEMO\")"
cast send "$ICS20_ADDRESS" \
  "sendTransfer((address,uint256,string,string,string,uint64,string))" \
  "$TRANSFER_TUPLE" \
  --rpc-url "$ETH_RPC_URL" \
  --private-key "$ETH_PRIVATE_KEY"

echo
echo "Done. If the relayer is running, the packet should be relayed to Cosmos."
