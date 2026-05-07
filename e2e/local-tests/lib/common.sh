#!/usr/bin/env bash

set -euo pipefail

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

  if [ -z "${ETH_RPC_URL:-}" ]; then
    ETH_RPC_URL="$(printf '%s\n' "$inspect" | perl -ne '
      if (/el-1-geth-lighthouse/) { $in=1 }
      elsif ($in && /^\S/) { $in=0 }
      elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
        print "http://$1\n";
        exit;
      }
    ')"
  fi

  if [ -z "${ETH_WS_URL:-}" ]; then
    ETH_WS_URL="$(printf '%s\n' "$inspect" | perl -ne '
      if (/el-1-geth-lighthouse/) { $in=1 }
      elsif ($in && /^\S/) { $in=0 }
      elsif ($in && /^\s+ws:.*->\s+(127\.0\.0\.1:\d+)/) {
        print "ws://$1\n";
        exit;
      }
    ')"
  fi

  if [ -z "${ETH_BEACON_API_URL:-}" ]; then
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
  if [ -z "${BROADCAST_JSON:-}" ]; then
    local chain_id chain_broadcast
    chain_id="$(cast chain-id --rpc-url "${ETH_RPC_URL:?}")"
    chain_broadcast="$REPO_ROOT/broadcast/E2ETestDeploy.s.sol/$chain_id/run-latest.json"
    if [ -f "$chain_broadcast" ]; then
      BROADCAST_JSON="$chain_broadcast"
    else
      BROADCAST_JSON="$(find "$REPO_ROOT/broadcast/E2ETestDeploy.s.sol" -path '*/run-latest.json' -type f -printf '%T@ %p\n' 2>/dev/null | sort -nr | awk 'NR == 1 { print $2 }')"
    fi
  fi

  if [ ! -f "${BROADCAST_JSON:-}" ]; then
    return 0
  fi

  if [ -z "${ERC20_ADDRESS:-}" ]; then
    ERC20_ADDRESS="$(jq -er '(.erc20 // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .erc20))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi

  if [ -z "${ICS20_ADDRESS:-}" ]; then
    ICS20_ADDRESS="$(jq -er '(.ics20Transfer // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .ics20Transfer))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi

  if [ -z "${ICS26_ADDRESS:-}" ]; then
    ICS26_ADDRESS="$(jq -er '(.ics26Router // (.returns["0"].value | gsub("\\\\\""; "\"") | fromjson | .ics26Router))' "$BROADCAST_JSON" 2>/dev/null || true)"
  fi
}

load_cosmos_receiver() {
  if [ -n "${RECEIVER:-}" ]; then
    return 0
  fi

  if command -v "${COSMOS_BIN:-gaiad}" >/dev/null 2>&1; then
    RECEIVER="$("$COSMOS_BIN" keys show "$COSMOS_RECEIVER_KEY" -a \
      --keyring-backend "${COSMOS_KEYRING:-test}" \
      --home "${COSMOS_HOME:-$HOME/.gaia}" 2>/dev/null || true)"
  fi
}

safe_cast_call() {
  cast call "$@" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null || true
}

contract_has_code() {
  local address="$1"
  local code
  code="$(cast code "$address" --rpc-url "${ETH_RPC_URL:?}" 2>/dev/null || true)"
  [ -n "$code" ] && [ "$code" != "0x" ]
}

load_env_file() {
  local env_file="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
  if [ -f "$env_file" ]; then
    set -a
    # shellcheck disable=SC1090
    source "$env_file"
    set +a
  fi
}
