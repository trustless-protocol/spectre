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



require_cmd cast
require_cmd jq
discover_kurtosis_endpoints
load_contract_addresses

if [ -z "$ETH_RPC_URL" ]; then
  echo "missing ETH_RPC_URL; set it in $ENV_FILE or start Kurtosis with ./run_eth_node.sh" >&2
  exit 1
fi

if [ -z "${ADDRESS:-}" ] && [ -n "$ETH_PRIVATE_KEY" ]; then
  ADDRESS="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"
fi

echo "== Ethereum endpoint =="
printf 'KURTOSIS_ENCLAVE=%s\n' "$KURTOSIS_ENCLAVE"
printf 'ETH_RPC_URL=%s\n' "$ETH_RPC_URL"
printf 'chain_id=%s\n' "$(cast chain-id --rpc-url "$ETH_RPC_URL")"
printf 'latest_block=%s\n' "$(cast block-number --rpc-url "$ETH_RPC_URL")"
printf 'BROADCAST_JSON=%s\n' "${BROADCAST_JSON:-}"

if [ -n "${ADDRESS:-}" ]; then
  echo
  echo "== Account =="
  printf 'address=%s\n' "$ADDRESS"
  printf 'native_balance_wei=%s\n' "$(cast balance "$ADDRESS" --rpc-url "$ETH_RPC_URL")"
fi

if [ -n "${ERC20_ADDRESS:-}" ]; then
  echo
  echo "== ERC20 =="
  printf 'ERC20_ADDRESS=%s\n' "$ERC20_ADDRESS"
  if contract_has_code "$ERC20_ADDRESS"; then
    printf 'symbol=%s\n' "$(safe_cast_call "$ERC20_ADDRESS" 'symbol()(string)')"
    printf 'decimals=%s\n' "$(safe_cast_call "$ERC20_ADDRESS" 'decimals()(uint8)')"
    printf 'total_supply=%s\n' "$(safe_cast_call "$ERC20_ADDRESS" 'totalSupply()(uint256)')"
    if [ -n "${ADDRESS:-}" ]; then
      printf 'account_balance=%s\n' "$(safe_cast_call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ADDRESS")"
    fi
    if [ -n "${ADDRESS:-}" ] && [ -n "${ICS20_ADDRESS:-}" ]; then
      printf 'allowance_to_ICS20=%s\n' "$(safe_cast_call "$ERC20_ADDRESS" 'allowance(address,address)(uint256)' "$ADDRESS" "$ICS20_ADDRESS")"
    fi
  else
    echo "warning=no bytecode at ERC20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "${ICS20_ADDRESS:-}" ]; then
  echo
  echo "== ICS20Transfer =="
  printf 'ICS20_ADDRESS=%s\n' "$ICS20_ADDRESS"
  if contract_has_code "$ICS20_ADDRESS"; then
    printf 'ics26=%s\n' "$(safe_cast_call "$ICS20_ADDRESS" 'ics26()(address)')"
    printf 'escrow_for_%s=%s\n' "$SOURCE_CLIENT" "$(safe_cast_call "$ICS20_ADDRESS" 'getEscrow(string)(address)' "$SOURCE_CLIENT")"
  else
    echo "warning=no bytecode at ICS20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "$ETH_TX_HASH" ]; then
  echo
  echo "== Ethereum transaction =="
  printf 'tx_hash=%s\n' "$ETH_TX_HASH"
  cast tx "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
  echo
  echo "== Ethereum receipt =="
  cast receipt "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
fi
