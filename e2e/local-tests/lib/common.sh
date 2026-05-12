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

ensure_relayer_config() {
  local config_file="$1"
  local template_file="${2:-$REPO_ROOT/relayer/config.example.json}"

  if [ -f "$config_file" ]; then
    return 0
  fi

  if [ ! -f "$template_file" ]; then
    echo "Error: $config_file not found and template $template_file is missing" >&2
    exit 1
  fi

  cp "$template_file" "$config_file"
}

# Backfill new config fields (ics26_client_id, cosmos_wasm_client_id) that older
# config.json files may be missing. Defaults to src_chain for ics26_client_id
# and "08-wasm-0" for cosmos_wasm_client_id.
backfill_relayer_config_fields() {
  local config_file="$1"

  if [ ! -f "$config_file" ]; then
    return 0
  fi

  jq '
    (.. | objects | select(.name == "cosmos_to_eth") | select(.config.cosmos_wasm_client_id == null) | .config.cosmos_wasm_client_id) = "08-wasm-0"
  | (.. | objects | select(.name == "cosmos_to_eth") | select(.config.ics26_client_id == null) | .config.ics26_client_id) = (.src_chain // "test-ibc-eth")
  ' "$config_file" > "$config_file.tmp" && mv "$config_file.tmp" "$config_file"
}

refresh_relayer_eth_endpoints() {
  local config_file="$1"

  backfill_relayer_config_fields "$config_file"

  ETH_RPC_URL=""
  ETH_WS_URL=""
  ETH_BEACON_API_URL=""
  discover_kurtosis_endpoints

  if [ -z "${ETH_RPC_URL:-}" ]; then
    ETH_RPC_URL="$(jq -er '[.. | objects | .eth_rpc_url? // empty][0]' "$config_file" 2>/dev/null || true)"
  fi
  if [ -z "${ETH_WS_URL:-}" ]; then
    ETH_WS_URL="$(jq -er '[.. | objects | .eth_ws_url? // empty][0]' "$config_file" 2>/dev/null || true)"
  fi
  if [ -z "${ETH_BEACON_API_URL:-}" ]; then
    ETH_BEACON_API_URL="$(jq -er '[.. | objects | .eth_beacon_api_url? // empty][0]' "$config_file" 2>/dev/null || true)"
  fi

  if [ -z "${ETH_RPC_URL:-}" ]; then
    echo "  ✗ ETH RPC URL not found in Kurtosis or $config_file" >&2
    exit 1
  fi

  if ! cast chain-id --rpc-url "$ETH_RPC_URL" >/dev/null 2>&1; then
    echo "  ✗ Ethereum RPC not responding at $ETH_RPC_URL" >&2
    echo "    Run ./setup/01-eth-node.sh or refresh $config_file with the current Kurtosis RPC URL." >&2
    exit 1
  fi

  echo "  ETH_RPC_URL:        $ETH_RPC_URL"
  [ -n "${ETH_WS_URL:-}" ] && echo "  ETH_WS_URL:         $ETH_WS_URL"
  [ -n "${ETH_BEACON_API_URL:-}" ] && echo "  ETH_BEACON_API_URL: $ETH_BEACON_API_URL"

  jq \
    --arg ETH_RPC "$ETH_RPC_URL" \
    --arg ETH_WS "${ETH_WS_URL:-}" \
    --arg ETH_BEACON "${ETH_BEACON_API_URL:-}" '
      (.. | objects | select(has("eth_rpc_url")) | .eth_rpc_url) = $ETH_RPC
    | (.. | objects | select(has("eth_ws_url")) | .eth_ws_url) = $ETH_WS
    | (.. | objects | select(has("eth_beacon_api_url")) | .eth_beacon_api_url) = $ETH_BEACON
    ' "$config_file" > "$config_file.tmp"
  mv "$config_file.tmp" "$config_file"
}

update_relayer_deploy_config() {
  local config_file="$1"
  local eth_rpc="$2"
  local eth_ws="$3"
  local ics26="$4"
  local wrapper="$5"
  local membership="$6"
  local update_client="$7"
  local misbehaviour="$8"
  local eth_beacon="$9"

  backfill_relayer_config_fields "$config_file"

  jq \
    --arg ETH_RPC "$eth_rpc" \
    --arg ETH_WS "$eth_ws" \
    --arg ICS26 "$ics26" \
    --arg WRAP "$wrapper" \
    --arg MEMB "$membership" \
    --arg UPCL "$update_client" \
    --arg MIS "$misbehaviour" \
    --arg ETH_BEACON "$eth_beacon" '
      (.. | objects | select(has("eth_rpc_url")) | .eth_rpc_url) = $ETH_RPC
    | (.. | objects | select(has("eth_ws_url")) | .eth_ws_url) = $ETH_WS
    | (.. | objects | select(has("ics26_address")) | .ics26_address) = $ICS26
    | (.. | objects | select(has("wrapper_verifier")) | .wrapper_verifier) = $WRAP
    | (.. | objects | select(has("membership")) | .membership) = $MEMB
    | (.. | objects | select(has("update_client")) | .update_client) = $UPCL
    | (.. | objects | select(has("misbehaviour")) | .misbehaviour) = $MIS
    | (.. | objects | select(has("eth_beacon_api_url")) | .eth_beacon_api_url) = $ETH_BEACON
    ' "$config_file" > "$config_file.tmp"
  mv "$config_file.tmp" "$config_file"
}
