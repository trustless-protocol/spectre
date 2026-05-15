#!/usr/bin/env bash

set -euo pipefail

# ETH-specific helpers for local E2E tests.
# Requires: ETH_RPC_URL, ETH_PRIVATE_KEY to be set.
# Source after common.sh.

eth_approve() {
  local token="$1" spender="$2" amount="$3"
  cast send "$token" "approve(address,uint256)" "$spender" "$amount" \
    --rpc-url "${ETH_RPC_URL:?}" \
    --private-key "${ETH_PRIVATE_KEY:?}" > /dev/null
}

eth_send_transfer() {
  local ics20="$1" erc20="$2" amount="$3" receiver="$4" source_client="$5" dest_port="$6" timeout="$7" memo="${8:-}"
  local tuple
  tuple="($erc20,$amount,$receiver,$source_client,$dest_port,$timeout,\"$memo\")"
  cast send "$ics20" \
    "sendTransfer((address,uint256,string,string,string,uint64,string))" \
    "$tuple" \
    --rpc-url "${ETH_RPC_URL:?}" \
    --private-key "${ETH_PRIVATE_KEY:?}" > /dev/null
}

eth_send_transfer_with_receipt() {
  local ics20="$1" erc20="$2" amount="$3" receiver="$4" source_client="$5" dest_port="$6" timeout="$7" memo="${8:-}"
  local tuple
  tuple="($erc20,$amount,$receiver,$source_client,$dest_port,$timeout,\"$memo\")"
  cast send "$ics20" \
    "sendTransfer((address,uint256,string,string,string,uint64,string))" \
    "$tuple" \
    --rpc-url "${ETH_RPC_URL:?}" \
    --private-key "${ETH_PRIVATE_KEY:?}" \
    --json 2>/dev/null
}

eth_call_multicall() {
  local ics20="$1"
  shift
  local arr=""
  local sep=""
  for arg in "$@"; do
    arr="${arr}${sep}${arg}"
    sep=", "
  done
  cast send "$ics20" \
    "multicall(bytes[])" \
    "[$arr]" \
    --rpc-url "${ETH_RPC_URL:?}" \
    --private-key "${ETH_PRIVATE_KEY:?}" > /dev/null
}

get_escrow_address() {
  local ics20="$1" client_id="$2"
  cast call "$ics20" 'getEscrow(string)(address)' "$client_id" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null || true
}

get_ibc_erc20_address() {
  local ics20="$1" denom_path="$2"
  cast call "$ics20" 'ibcERC20Contract(string)(address)' "$denom_path" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null || true
}

get_packet_commitment() {
  local ics26="$1" client_id="$2" sequence="$3"
  cast call "$ics26" 'getPacketCommitment(string,uint64)(bytes32)' "$client_id" "$sequence" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null || true
}

eth_balance_of() {
  local token="$1" address="$2"
  cast call "$token" 'balanceOf(address)(uint256)' "$address" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null | uint_value || true
}

eth_wallet_address() {
  cast wallet address --private-key "${ETH_PRIVATE_KEY:?}" 2>/dev/null
}
