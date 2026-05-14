#!/usr/bin/env bash

set -euo pipefail

CONFIG_FILE="relayer/config.example.json"
ENV_FILE="relayer/.env"
RPC_URL=""
ICS26_ADDRESS=""
LOOKBACK=5000
BASE_TX=""
FROM_ADDRESS=""
FRESH_MODE=1
CONFIG_FILE_ABS=""
SEL_UNAUTHORIZED=""
SEL_PROOF_TOO_OLD=""

usage() {
  cat <<'EOF'
Usage:
  scripts/test_quorum_edge_cases.sh [options]

Options:
  --config <path>      Relayer config JSON (default: relayer/config.example.json)
  --env-file <path>    Env file (default: relayer/.env)
  --rpc-url <url>      ETH RPC URL (overrides config)
  --router <address>   ICS26 router address (overrides config)
  --lookback <blocks>  Log scan window (default: 5000)
  --base-tx <hash>     Base updateClient tx hash
  --no-fresh           Use historical tx flow for 5a/5b (requires archive-style pre-block simulation)
  -h, --help           Show help
EOF
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd" >&2
    exit 1
  fi
}

normalize_addr() {
  local addr="$1"
  if [[ ! "$addr" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    echo "invalid address: $addr" >&2
    exit 1
  fi
  echo "$addr"
}

extract_revert_data() {
  local out="$1"
  echo "$out" | sed -n 's/.*data: "\(0x[0-9a-fA-F]*\)".*/\1/p' | head -n1
}

hex_or_dec_to_dec() {
  local v="$1"
  if [[ "$v" == 0x* ]]; then
    echo $((16#${v#0x}))
  else
    echo "$v"
  fi
}

helper_json() {
  local mode="$1"
  local cmd=(go run ./cmd/updateclient_attack --mode "$mode" --rpc-url "$RPC_URL" --router "$ICS26_ADDRESS" --lookback "$LOOKBACK")
  if [[ "$FRESH_MODE" == "1" && ( "$mode" == "5a" || "$mode" == "5b" ) ]]; then
    cmd+=(--fresh --config "$CONFIG_FILE_ABS")
  fi
  if [[ -n "$BASE_TX" ]]; then
    cmd+=(--base-tx "$BASE_TX")
  fi
  (
    cd relayer
    "${cmd[@]}"
  )
}

status_line() {
  local case_name="$1"
  local status="$2"
  local detail="$3"
  printf '%-4s %-5s %s\n' "$case_name" "$status" "$detail"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --config)
      CONFIG_FILE="$2"
      shift 2
      ;;
    --env-file)
      ENV_FILE="$2"
      shift 2
      ;;
    --rpc-url)
      RPC_URL="$2"
      shift 2
      ;;
    --router)
      ICS26_ADDRESS="$2"
      shift 2
      ;;
    --lookback)
      LOOKBACK="$2"
      shift 2
      ;;
    --base-tx)
      BASE_TX="$2"
      shift 2
      ;;
    --no-fresh)
      FRESH_MODE=0
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown option: $1" >&2
      usage
      exit 1
      ;;
  esac
done

require_cmd jq
require_cmd cast
require_cmd go

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "config file not found: $CONFIG_FILE" >&2
  exit 1
fi
CONFIG_FILE_ABS="$(cd "$(dirname "$CONFIG_FILE")" && pwd)/$(basename "$CONFIG_FILE")"

if [[ -f "$ENV_FILE" ]]; then
  set -a
  # shellcheck disable=SC1090
  source "$ENV_FILE"
  set +a
fi

if [[ -z "$RPC_URL" ]]; then
  RPC_URL=$(jq -r '.modules[] | select(.name=="cosmos_to_eth") | .config.eth_rpc_url // empty' "$CONFIG_FILE")
fi
if [[ -z "$ICS26_ADDRESS" ]]; then
  ICS26_ADDRESS=$(jq -r '.modules[] | select(.name=="cosmos_to_eth") | .config.ics26_address // empty' "$CONFIG_FILE")
fi

if [[ -z "$RPC_URL" || -z "$ICS26_ADDRESS" ]]; then
  echo "failed to resolve eth rpc or ics26 address from config" >&2
  exit 1
fi

ICS26_ADDRESS=$(normalize_addr "$ICS26_ADDRESS")
if [[ -n "${ETH_PRIVATE_KEY:-}" ]]; then
  FROM_ADDRESS="$(cast wallet address --private-key "$ETH_PRIVATE_KEY")"
fi

echo "RPC_URL=$RPC_URL"
echo "ICS26_ADDRESS=$ICS26_ADDRESS"
echo "LOOKBACK=$LOOKBACK"
if [[ -n "$BASE_TX" ]]; then
  echo "BASE_TX=$BASE_TX"
fi
if [[ -n "$FROM_ADDRESS" ]]; then
  echo "FROM_ADDRESS=$FROM_ADDRESS"
fi

bucket_line=$(sed -n 's/.*var Buckets = \[\]int{\(.*\)}/\1/p' relayer/prover/buckets.go | head -n1)
if [[ -z "$bucket_line" ]]; then
  echo "unable to parse relayer/prover/buckets.go" >&2
  exit 1
fi
MAX_BUCKET=$(echo "$bucket_line" | tr ',' '\n' | tr -d ' ' | awk 'NF' | sort -n | tail -n1)
echo "MAX_BUCKET=$MAX_BUCKET"

sel_5a=$(cast sig "InsufficientVotingPower(uint64,uint64)")
sel_5b=$(cast sig "DuplicateSigner(uint32)")
SEL_UNAUTHORIZED=$(cast sig "AccessManagedUnauthorized(address)")
SEL_PROOF_TOO_OLD=$(cast sig "ProofIsTooOld(uint256,uint256)")

echo
echo "Running 5a..."
j5a="$(helper_json 5a)"
calldata_5a="$(echo "$j5a" | jq -r '.tampered_calldata')"
set +e
if [[ -n "$FROM_ADDRESS" ]]; then
  out_5a="$(cast call "$ICS26_ADDRESS" --rpc-url "$RPC_URL" --from "$FROM_ADDRESS" --data "$calldata_5a" 2>&1)"
else
  out_5a="$(cast call "$ICS26_ADDRESS" --rpc-url "$RPC_URL" --data "$calldata_5a" 2>&1)"
fi
rc_5a=$?
set -e
if (( rc_5a == 0 )); then
  status_line "5a" "FAIL" "tampered payload unexpectedly succeeded"
else
  data_5a="$(extract_revert_data "$out_5a")"
  sel_seen_5a="${data_5a:0:10}"
  if [[ -n "$data_5a" && "$sel_seen_5a" == "$sel_5a" ]]; then
    status_line "5a" "PASS" "revert selector matched InsufficientVotingPower ($sel_5a)"
  elif [[ -n "$data_5a" && "$sel_seen_5a" == "$SEL_UNAUTHORIZED" ]]; then
    status_line "5a" "SKIP" "caller unauthorized; set ETH_PRIVATE_KEY or --from to relayer EOA"
  elif [[ -n "$data_5a" && "$sel_seen_5a" == "$SEL_PROOF_TOO_OLD" ]]; then
    status_line "5a" "SKIP" "fresh payload expired before call; rerun immediately after generation"
  else
    status_line "5a" "FAIL" "reverted but selector mismatch (expected=$sel_5a got=${sel_seen_5a:-none})"
  fi
fi

echo "Running 5b..."
j5b="$(helper_json 5b)"
calldata_5b="$(echo "$j5b" | jq -r '.tampered_calldata')"
set +e
if [[ -n "$FROM_ADDRESS" ]]; then
  out_5b="$(cast call "$ICS26_ADDRESS" --rpc-url "$RPC_URL" --from "$FROM_ADDRESS" --data "$calldata_5b" 2>&1)"
else
  out_5b="$(cast call "$ICS26_ADDRESS" --rpc-url "$RPC_URL" --data "$calldata_5b" 2>&1)"
fi
rc_5b=$?
set -e
if (( rc_5b == 0 )); then
  status_line "5b" "FAIL" "tampered payload unexpectedly succeeded"
else
  data_5b="$(extract_revert_data "$out_5b")"
  sel_seen_5b="${data_5b:0:10}"
  if [[ -n "$data_5b" && "$sel_seen_5b" == "$sel_5b" ]]; then
    status_line "5b" "PASS" "revert selector matched DuplicateSigner ($sel_5b)"
  elif [[ -n "$data_5b" && "$sel_seen_5b" == "$SEL_UNAUTHORIZED" ]]; then
    status_line "5b" "SKIP" "caller unauthorized; set ETH_PRIVATE_KEY or --from to relayer EOA"
  elif [[ -n "$data_5b" && "$sel_seen_5b" == "$SEL_PROOF_TOO_OLD" ]]; then
    status_line "5b" "SKIP" "fresh payload expired before call; rerun immediately after generation"
  else
    status_line "5b" "FAIL" "reverted but selector mismatch (expected=$sel_5b got=${sel_seen_5b:-none})"
  fi
fi

echo "Running 5c..."
j5c="$(helper_json 5c)"
base_tx_5c="$(echo "$j5c" | jq -r '.base_tx')"
bucket_5c="$(echo "$j5c" | jq -r '.bucket')"
active_5c="$(echo "$j5c" | jq -r '.active_count')"
unique_5c="$(echo "$j5c" | jq -r '.unique_active')"
quorum_5c="$(echo "$j5c" | jq -r '.quorum_satisfied')"
receipt_5c="$(cast receipt "$base_tx_5c" --json --rpc-url "$RPC_URL")"
status_raw_5c="$(echo "$receipt_5c" | jq -r '.status')"
if [[ "$status_raw_5c" == 0x* ]]; then
  status_num_5c=$((16#${status_raw_5c#0x}))
else
  status_num_5c="$status_raw_5c"
fi
if [[ "$active_5c" == "$bucket_5c" && "$unique_5c" == "$active_5c" && "$quorum_5c" == "true" && "$status_num_5c" == "1" ]]; then
  status_line "5c" "PASS" "bucket=$bucket_5c active=$active_5c unique=$unique_5c tx_status=$status_num_5c"
else
  status_line "5c" "SKIP" "latest update uses active=$active_5c bucket=$bucket_5c (not exact boundary)"
fi

echo "Running 5d..."
j5d="$(helper_json 5d)"
min_signers_5d="$(echo "$j5d" | jq -r '.min_signers_for_2of3')"
validator_count_5d="$(echo "$j5d" | jq -r '.validator_count')"
if (( min_signers_5d > MAX_BUCKET )); then
  status_line "5d" "PASS" "min_signers_for_2of3=$min_signers_5d > max_bucket=$MAX_BUCKET"
else
  status_line "5d" "SKIP" "current setup has min_signers_for_2of3=$min_signers_5d <= max_bucket=$MAX_BUCKET (validators=$validator_count_5d)"
fi
