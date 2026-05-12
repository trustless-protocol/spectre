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
  log_err "missing ETH_RPC_URL; set it in $ENV_FILE or start Kurtosis with ./run_eth_node.sh"
  exit 1
fi

if [ -z "${ADDRESS:-}" ] && [ -n "$ETH_PRIVATE_KEY" ]; then
  ADDRESS="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"
fi

log_header "Ethereum Endpoint"
log_kv "KURTOSIS_ENCLAVE" "$KURTOSIS_ENCLAVE"
log_kv "ETH_RPC_URL" "$ETH_RPC_URL"
log_kv "chain_id" "$(cast chain-id --rpc-url "$ETH_RPC_URL")"
log_kv "latest_block" "$(cast block-number --rpc-url "$ETH_RPC_URL")"
log_kv "BROADCAST_JSON" "${BROADCAST_JSON:-}"

if [ -n "${ADDRESS:-}" ]; then
  log_header "Account"
  log_kv "address" "$ADDRESS"
  log_kv "native_balance_wei" "$(cast balance "$ADDRESS" --rpc-url "$ETH_RPC_URL")"
fi

if [ -n "${ERC20_ADDRESS:-}" ]; then
  log_header "ERC20"
  log_kv "address" "$ERC20_ADDRESS"
  if contract_has_code "$ERC20_ADDRESS"; then
    log_kv "symbol" "$(safe_cast_call "$ERC20_ADDRESS" 'symbol()(string)')"
    log_kv "decimals" "$(safe_cast_call "$ERC20_ADDRESS" 'decimals()(uint8)')"
    log_kv "total_supply" "$(safe_cast_call "$ERC20_ADDRESS" 'totalSupply()(uint256)')"
    if [ -n "${ADDRESS:-}" ]; then
      log_kv "account_balance" "$(safe_cast_call "$ERC20_ADDRESS" 'balanceOf(address)(uint256)' "$ADDRESS")"
    fi
    if [ -n "${ADDRESS:-}" ] && [ -n "${ICS20_ADDRESS:-}" ]; then
      log_kv "allowance_to_ICS20" "$(safe_cast_call "$ERC20_ADDRESS" 'allowance(address,address)(uint256)' "$ADDRESS" "$ICS20_ADDRESS")"
    fi
  else
    log_warn "no bytecode at ERC20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "${ICS20_ADDRESS:-}" ]; then
  log_header "ICS20Transfer"
  log_kv "address" "$ICS20_ADDRESS"
  if contract_has_code "$ICS20_ADDRESS"; then
    log_kv "ics26" "$(safe_cast_call "$ICS20_ADDRESS" 'ics26()(address)')"
    log_kv "escrow (${SOURCE_CLIENT})" "$(safe_cast_call "$ICS20_ADDRESS" 'getEscrow(string)(address)' "$SOURCE_CLIENT")"
  else
    log_warn "no bytecode at ICS20_ADDRESS on ETH_RPC_URL"
  fi
fi

if [ -n "$ETH_TX_HASH" ]; then
  log_header "Ethereum Transaction"
  log_kv "tx_hash" "$ETH_TX_HASH"
  cast tx "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
  log_header "Ethereum Receipt"
  cast receipt "$ETH_TX_HASH" --rpc-url "$ETH_RPC_URL"
fi
