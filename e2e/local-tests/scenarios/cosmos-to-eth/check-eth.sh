#!/usr/bin/env bash
set -euo pipefail

# Inspect Ethereum state for Cosmos → ETH path.
# Checks ibcERC20 contracts and balances.

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/eth.sh"

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
COSMOS_NATIVE_DENOM="${COSMOS_NATIVE_DENOM:-stake}"

require_cmd cast
require_cmd jq
discover_kurtosis_endpoints

if [ -z "$ETH_RPC_URL" ]; then
  log_err "missing ETH_RPC_URL"
  exit 1
fi

load_contract_addresses

if [ -z "${ADDRESS:-}" ] && [ -n "$ETH_PRIVATE_KEY" ]; then
  ADDRESS="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"
fi

log_header "Ethereum Endpoint"
log_kv "ETH_RPC_URL" "$ETH_RPC_URL"
log_kv "chain_id" "$(cast chain-id --rpc-url "$ETH_RPC_URL")"
log_kv "latest_block" "$(cast block-number --rpc-url "$ETH_RPC_URL")"

if [ -n "${ADDRESS:-}" ]; then
  log_header "Account"
  log_kv "address" "$ADDRESS"
  log_kv "native_balance_wei" "$(cast balance "$ADDRESS" --rpc-url "$ETH_RPC_URL")"
fi

if [ -n "${ICS20_ADDRESS:-}" ]; then
  log_header "ICS20Transfer"
  log_kv "address" "$ICS20_ADDRESS"
  if contract_has_code "$ICS20_ADDRESS"; then
    log_kv "ics26" "$(safe_cast_call "$ICS20_ADDRESS" 'ics26()(address)')"

    # Check native denom ibcERC20
    NATIVE_DENOM_PATH="${COSMOS_NATIVE_DENOM}/transfer/${SOURCE_CLIENT}"
    NATIVE_IBCERC20="$(get_ibc_erc20_address "$ICS20_ADDRESS" "$NATIVE_DENOM_PATH")"
    if [ -n "$NATIVE_IBCERC20" ] && [ "$NATIVE_IBCERC20" != "0x0000000000000000000000000000000000000000" ]; then
      log_header "ibcERC20 (native Cosmos coin)"
      log_kv "denom_path" "$NATIVE_DENOM_PATH"
      log_kv "address" "$NATIVE_IBCERC20"
      log_kv "name" "$(safe_cast_call "$NATIVE_IBCERC20" 'name()(string)')"
      log_kv "symbol" "$(safe_cast_call "$NATIVE_IBCERC20" 'symbol()(string)')"
      log_kv "totalSupply" "$(safe_cast_call "$NATIVE_IBCERC20" 'totalSupply()(uint256)')"
      if [ -n "${ADDRESS:-}" ]; then
        log_kv "account_balance" "$(safe_cast_call "$NATIVE_IBCERC20" 'balanceOf(address)(uint256)' "$ADDRESS")"
      fi
    fi
  fi
fi
