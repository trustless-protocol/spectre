#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
source "$SCRIPT_DIR/../lib/sleep.sh"

usage() {
  cat >&2 <<'HELP'
Usage: sleep.sh <command> [args...]

Commands:
  beacon-finality [enclave]         Wait until beacon finalizes (finalized.epoch > 0)
  proposal-exists [chain-id]        Wait until a governance proposal exists (prints ID)
  proposal-status <id> <status>     Wait until a proposal reaches given status
  cosmos-block [target]             Wait for N blocks (default: 1)
  condition <desc> <cmd...>         Generic wait: run cmd repeatedly until exit 0

Examples:
  sleep.sh beacon-finality
  sleep.sh beacon-finality my-testnet
  sleep.sh proposal-exists test-ibc-eth
  sleep.sh proposal-status 1 PROPOSAL_STATUS_PASSED
  sleep.sh cosmos-block 2
  sleep.sh condition "beacon API available" "curl -sf http://127.0.0.1:32774/eth/v1/node/syncing"
HELP
  exit 1
}

cmd="${1:-}"
shift 2>/dev/null || usage

case "$cmd" in
  beacon-finality)
    enclave="${1:-my-testnet}"
    beacon_api=$(wait_for_beacon_api "$enclave" 60 10)
    wait_for_beacon_finality "$beacon_api" 120 10
    ;;
  proposal-exists)
    chain_id="${1:-test-ibc-eth}"
    rpc_url="${2:-http://127.0.0.1:26657}"
    proposal_id=$(wait_for_proposal_exists "$chain_id" "$rpc_url" 30 2)
    echo "$proposal_id"
    ;;
  proposal-status)
    proposal_id="${1:?proposal_id required}"
    status="${2:-PROPOSAL_STATUS_PASSED}"
    chain_id="${3:-test-ibc-eth}"
    rpc_url="${4:-http://127.0.0.1:26657}"
    wait_for_proposal_status "$proposal_id" "$status" "$chain_id" "$rpc_url" 60 5
    ;;
  cosmos-block)
    target="${1:-1}"
    COSMOS_BIN="${COSMOS_BIN:-gaiad}"
    COSMOS_RPC_URL="${COSMOS_RPC_URL:-http://127.0.0.1:26657}"
    before=$("$COSMOS_BIN" status --node "$COSMOS_RPC_URL" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')
    echo "Current block: $before" >&2
    target_block=$((before + target))
    while true; do
      after=$("$COSMOS_BIN" status --node "$COSMOS_RPC_URL" 2>/dev/null | jq -r '.sync_info.latest_block_height // .SyncInfo.latest_block_height // empty')
      if [ "$after" -ge "$target_block" ] 2>/dev/null; then
        break
      fi
      sleep 1
    done
    echo "Reached block $after" >&2
    ;;
  condition)
    desc="${1:?description required}"
    shift
    wait_for_condition "$desc" 60 10 "$*"
    ;;
  *)
    usage
    ;;
esac
