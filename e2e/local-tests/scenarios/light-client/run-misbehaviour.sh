#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Misbehaviour detection and client freezing.
# This is an advanced P2 scenario that requires conflicting header generation.
#
# Corresponds to: Test_DoubleSignMisbehaviour, Test_BreakingTimeMonotonicityMisbehaviour

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

ENV_FILE="${ENV_FILE:-$REPO_ROOT/relayer/.env}"
CONFIG_FILE="${CONFIG_FILE:-$REPO_ROOT/relayer/config.json}"

config_value() {
  local key="$1"
  if [ -f "$CONFIG_FILE" ]; then
    jq -er --arg key "$key" '([.. | objects | .[$key]? // empty][0]) // empty' "$CONFIG_FILE" 2>/dev/null || true
  fi
}

load_env_file
discover_kurtosis_endpoints

if [ -z "${ETH_RPC_URL:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ETH_RPC_URL="$(config_value eth_rpc_url)"
fi

if [ -z "${ICS07_CLIENT:-}" ] && [ -f "$CONFIG_FILE" ]; then
  ICS07_CLIENT="$(jq -er '.. | objects | select(.name == "cosmos_to_eth") | .config.ics07_client // empty' "$CONFIG_FILE" 2>/dev/null || true)"
fi

if [ -z "${ETH_RPC_URL:-}" ]; then
  log_err "ETH_RPC_URL not set"
  exit 1
fi
if [ -z "${ICS07_CLIENT:-}" ]; then
  log_err "ICS07_CLIENT address not found in $CONFIG_FILE"
  exit 1
fi

log_header "Misbehaviour Test"
log_warn "This scenario requires generating conflicting CometBFT headers."
log_warn "It is not feasible to implement in pure bash and requires:"
log_warn "  1. A forked Cosmos chain producing double-signed blocks, OR"
log_warn "  2. A custom Go test harness that crafts invalid headers."
log_warn ""
log_warn "For automated testing, run the Go test instead:"
log_warn "  cd e2e/interchaintestv8 && go test -run Test_DoubleSignMisbehaviour"

# Query current frozen state
FROZEN_RAW="$(cast call "$ICS07_CLIENT" 'clientState()(string,(uint8,uint8),(uint64,uint64),uint32,uint32,bool,uint8)' --rpc-url "$ETH_RPC_URL" 2>/dev/null || true)"
if [ -n "$FROZEN_RAW" ]; then
  IS_FROZEN="$(echo "$FROZEN_RAW" | awk 'NR==6 { print $1 }')"
  log_kv "Client frozen state" "${IS_FROZEN:-unknown}"
fi

log_ok "RESULT: INFO (manual/advanced scenario)"
