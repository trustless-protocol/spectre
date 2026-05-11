#!/usr/bin/env bash

# Reusable polling/waiting functions for E2E test scripts.
# Source this file: source "$(dirname "$0")/lib/sleep.sh"

set -euo pipefail

# wait_for_condition <desc> <max_polls> <poll_interval> <condition_cmd>
#   Generic polling: runs condition_cmd repeatedly (as a bash -c eval) until
#   it exits 0.  desc is a human-readable label for log messages.
wait_for_condition() {
  local desc="$1"
  local max_polls="$2"
  local poll_interval="$3"
  local condition_cmd="$4"
  local poll=0

  while [ $poll -lt "$max_polls" ]; do
    if eval "$condition_cmd" >/dev/null 2>&1; then
      echo "wait_for_condition[$desc] succeeded after ${poll} polls"
      return 0
    fi
    echo "[poll ${poll}/${max_polls}] waiting for: $desc"
    sleep "$poll_interval"
    poll=$((poll + 1))
  done

  echo "ERROR: wait_for_condition[$desc] timed out after ${max_polls} polls" >&2
  return 1
}

# wait_for_beacon_api <enclave> <max_polls> <poll_interval>
#   Prints the beacon API URL to stdout once the endpoint is available.
wait_for_beacon_api() {
  local enclave="${1:-my-testnet}"
  local max_polls="${2:-60}"
  local poll_interval="${3:-10}"

  local api_url=""
  local poll=0

  while [ $poll -lt "$max_polls" ]; do
    api_url=$(kurtosis enclave inspect "$enclave" 2>/dev/null \
    | awk '
      $0 ~ /cl-1-lighthouse-geth/ {in_service=1}
      in_service && /http:/ {
          match($0, /127\.0\.0\.1:[0-9]+/)
          print "http://" substr($0, RSTART, RLENGTH)
          exit
      }
      in_service && /^[^[:space:]]/ {in_service=0}
    ') || true

    if [ -n "$api_url" ]; then
      echo "$api_url"
      return 0
    fi

    echo "[poll ${poll}/${max_polls}] beacon API not yet available" >&2
    sleep "$poll_interval"
    poll=$((poll + 1))
  done

  echo "ERROR: beacon API did not become available" >&2
  return 1
}

# wait_for_beacon_finality <beacon_api_url> <max_polls> <poll_interval>
#   Polls /eth/v1/beacon/states/head/finality_checkpoints until
#   finalized.epoch > 0.
wait_for_beacon_finality() {
  local beacon_api="${1:?beacon_api_url required}"
  local max_polls="${2:-120}"
  local poll_interval="${3:-10}"
  local poll=0

  echo "Waiting for beacon finality on $beacon_api ..." >&2

  while [ $poll -lt "$max_polls" ]; do
    local finality finalized_epoch
    finality=$(curl -sf "$beacon_api/eth/v1/beacon/states/head/finality_checkpoints" 2>/dev/null || echo "")
    finalized_epoch=$(printf '%s' "$finality" | jq -r '.data.finalized.epoch // empty' 2>/dev/null || echo "")

    if [ -n "$finalized_epoch" ] && [ "$finalized_epoch" -gt 0 ] 2>/dev/null; then
      echo "Beacon finalized at epoch $finalized_epoch" >&2
      return 0
    fi

    echo "[poll ${poll}/${max_polls}] beacon not yet finalized (epoch=${finalized_epoch:-unavailable})" >&2
    sleep "$poll_interval"
    poll=$((poll + 1))
  done

  echo "ERROR: beacon did not finalize within ${max_polls} polls" >&2
  return 1
}

# wait_for_proposal_exists <chain_id> <rpc_url> <max_polls> <poll_interval>
#   Polls until a governance proposal exists. Prints the latest proposal ID.
wait_for_proposal_exists() {
  local chain_id="${1:?chain_id required}"
  local rpc_url="${2:-http://127.0.0.1:26657}"
  local max_polls="${3:-30}"
  local poll_interval="${4:-2}"
  local poll=0

  while [ $poll -lt "$max_polls" ]; do
    local proposals proposal_id
    proposals=$(gaiad q gov proposals --chain-id "$chain_id" --node "$rpc_url" -o json 2>/dev/null || echo "")
    proposal_id=$(printf '%s' "$proposals" | jq -er '.proposals | sort_by(.id | tonumber) | last | .id' 2>/dev/null || echo "")

    if [ -n "$proposal_id" ]; then
      echo "$proposal_id"
      return 0
    fi

    sleep "$poll_interval"
    poll=$((poll + 1))
  done

  echo "ERROR: no governance proposal appeared" >&2
  return 1
}

# wait_for_proposal_status <proposal_id> <desired_status> <chain_id> <rpc_url> <max_polls> <poll_interval>
#   Polls until a governance proposal reaches the desired status (e.g. PASSED).
wait_for_proposal_status() {
  local proposal_id="${1:?proposal_id required}"
  local desired_status="${2:-PASSED}"
  local chain_id="${3:?chain_id required}"
  local rpc_url="${4:-http://127.0.0.1:26657}"
  local max_polls="${5:-60}"
  local poll_interval="${6:-5}"
  local poll=0

  while [ $poll -lt "$max_polls" ]; do
    local status
    status=$(gaiad q gov proposal "$proposal_id" --chain-id "$chain_id" --node "$rpc_url" -o json 2>/dev/null \
      | jq -r '.status // .proposal.status // "PROPOSAL_STATUS_UNSPECIFIED"' 2>/dev/null || echo "")

    if [ "$status" = "$desired_status" ]; then
      echo "Proposal $proposal_id reached status: $desired_status"
      return 0
    fi

    echo "[poll ${poll}/${max_polls}] proposal $proposal_id status=$status (waiting for $desired_status)"
    sleep "$poll_interval"
    poll=$((poll + 1))
  done

  echo "ERROR: proposal $proposal_id did not reach $desired_status within ${max_polls} polls" >&2
  return 1
}
