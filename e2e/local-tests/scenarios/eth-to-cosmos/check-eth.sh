#!/usr/bin/env bash
set -euo pipefail

# Inspect local Ethereum data for the ETH -> Cosmos path.
# Optional inputs:
#   ADDRESS=0x...              account to inspect
#   ETH_TX_HASH=0x...          Ethereum tx hash to inspect
#   ERC20_ADDRESS=0x...        token contract override
#   ICS20_ADDRESS=0x...        transfer contract override
#   KURTOSIS_ENCLAVE=name      defaults to my-testnet

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
ETH_RPC_URL="${ETH_RPC_URL:-${ETH_RPC:-}}"
ETH_PRIVATE_KEY="${ETH_PRIVATE_KEY:-${ETH_USER_PK:-${PRIVATE_KEY:-}}}"
SOURCE_CLIENT="${SOURCE_CLIENT:-${ETH_SOURCE_CLIENT_ID:-cosmoshub-1}}"
ETH_TX_HASH="${ETH_TX_HASH:-${TX_HASH:-}}"

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd" >&2
    exit 1
  fi
}

discover_kurtosis_eth_rpc() {
  if [ -n "$ETH_RPC_URL" ] || ! command -v kurtosis >/dev/null 2>&1; then
    return 0
  fi

  local inspect
  if ! inspect="$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" 2>/dev/null)"; then
    return 0
  fi

  ETH_RPC_URL="$(printf '%s\n' "$inspect" | perl -ne '
    if (/el-1-geth-lighthouse/) { $in=1 }
    elsif ($in && /^\S/) { $in=0 }
    elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
      print "http://$1\n";
      exit;
    }
  ')"
}

load_contract_addresses() {
  if [ -n "${ERC20_ADDRESS:-}" ] && [ -n "${ICS20_ADDRESS:-}" ]; then
    return 0
  fi

  if [ -z "$BROADCAST_JSON" ]; then
    local chain_id chain_broadcast
    chain_id="$(cast chain-id --rpc-url "$ETH_RPC_URL")"
    chain_broadcast="$REPO_ROOT/broadcast/E2ETestDeploy.s.sol/$chain_id/run-latest.json"
    if [ -f "$chain_broadcast" ]; then
      BROADCAST_JSON="$chain_broadcast"
    elif [ -d "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" ]; then
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

safe_cast_call() {
  cast call "$@" --rpc-url "$ETH_RPC_URL" 2>/dev/null || true
}

contract_has_code() {
  local address="$1"
  local code
  code="$(cast code "$address" --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
  [ -n "$code" ] && [ "$code" != "0x" ]
}

require_cmd cast
require_cmd jq
discover_kurtosis_eth_rpc
load_contract_addresses

if [ -z "$ETH_RPC_URL" ]; then
  echo "missing ETH_RPC_URL; set it in $ENV_FILE or start Kurtosis with ./run_eth_node.sh" >&2
  exit 1
fi

if [ -z "${ADDRESS:-}" ] && [ -n "$ETH_PRIVATE_KEY" ]; then
  ADDRESS="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"
fi

echo "━━━ Ethereum Endpoint ━━━"
printf "  %-20s %s\n" "KURTOSIS_ENCLAVE:" "$KURTOSIS_ENCLAVE"
printf "  %-20s %s\n" "ETH_RPC_URL:"      "$ETH_RPC_URL"
printf "  %-20s %s\n" "chain_id:"         "$(cast chain-id --rpc-url "$ETH_RPC_URL")"
printf "  %-20s %s\n" "latest_block:"     "$(cast block-number --rpc-url "$ETH_RPC_URL")"
printf "  %-20s %s\n" "BROADCAST_JSON:"   "${BROADCAST_JSON:-}"

if [ -n "${ADDRESS:-}" ]; then
  echo ""
  echo "━━━ Account ━━━"
  printf "  %-20s %s\n" "address:"             "$ADDRESS"
  printf "  %-20s %s\n" "native_balance_wei:"  "$(cast balance "$ADDRESS" --rpc-url "$ETH_RPC_URL")"
fi

if [ -n "${ERC20_ADDRESS:-}" ]; then
  echo ""
  echo "━━━ ERC20 ━━━"
  printf "  %-20s %s\n" "address:" "$ERC20_ADDRESS"
  if contract_has_code "$ERC20_ADDRESS"; then
    printf "  %-20s %s\n" "symbol:"        "$(safe_cast_call "$ERC20_ADDRESS" 'symbol()(string)')"
    printf "  %-20s %s\n" "decimals:"      "$(safe_cast_call "$ERC20_ADDRESS" 'decimals()(uint8)')"
    printf "  %-20s %s\n" "total_supply:"  "$(safe_cast_call "$ERC20_ADDRESS" 'totalSupply()(uint256)')"
    if [ -n "${ADDRESS:-}" ]; then
      printf "  %-20s %s\n" "account_balance:" "$(safe_cast_call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ADDRESS")"
    fi
    if [ -n "${ADDRESS:-}" ] && [ -n "${ICS20_ADDRESS:-}" ]; then
      printf "  %-20s %s\n" "allowance_to_ICS20:" "$(safe_cast_call "$ERC20_ADDRESS" 'allowance(address,address)(uint256)' "$ADDRESS" "$ICS20_ADDRESS")"
    fi
  else
    echo "  ! no bytecode at ERC20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "${ICS20_ADDRESS:-}" ]; then
  echo ""
  echo "━━━ ICS20Transfer ━━━"
  printf "  %-20s %s\n" "address:" "$ICS20_ADDRESS"
  if contract_has_code "$ICS20_ADDRESS"; then
    printf "  %-20s %s\n" "ics26:" "$(safe_cast_call "$ICS20_ADDRESS" 'ics26()(address)')"
    printf "  %-20s %s\n" "escrow (${SOURCE_CLIENT}):" "$(safe_cast_call "$ICS20_ADDRESS" 'getEscrow(string)(address)' "$SOURCE_CLIENT")"
  else
    echo "  ! no bytecode at ICS20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "$ETH_TX_HASH" ]; then
  echo ""
  echo "━━━ Ethereum Transaction ━━━"
  printf "  %-20s %s\n" "tx_hash:" "$ETH_TX_HASH"
  cast tx "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
  echo ""
  echo "━━━ Ethereum Receipt ━━━"
  cast receipt "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
fi
