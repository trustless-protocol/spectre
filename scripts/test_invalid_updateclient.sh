#!/usr/bin/env bash

set -euo pipefail

CONFIG_FILE="relayer/config.example.json"
ENV_FILE="relayer/.env"
RPC_URL=""
ICS26_ADDRESS=""
UPDATE_TX=""
DONOR_TX=""
LOOKBACK=5000
SEND_MODE=0
FROM_ADDRESS=""
MUTATE_OFFSET=""
CASE_MODE="4a"

usage() {
  cat <<'EOF'
Usage:
  scripts/test_invalid_updateclient.sh [options]

Options:
  --case <4a|4b|4c>      Test case mode (default: 4a)
  --config <path>         Relayer config JSON (default: relayer/config.example.json)
  --env-file <path>       Env file for ETH_PRIVATE_KEY (default: relayer/.env)
  --rpc-url <url>         ETH RPC URL (overrides config)
  --router <address>      ICS26 router address (overrides config)
  --update-tx <hash>      Base successful updateClient tx hash
  --donor-tx <hash>       Donor tx hash for 4b (optional if auto-detect can find previous tx)
  --lookback <blocks>     Log scan window for auto-detect (default: 5000)
  --from <address>        from address for eth_call (derived from ETH_PRIVATE_KEY if omitted)
  --offset <nibble-pos>   4a only: nibble offset in updateMsg body to mutate
  --send                  Broadcast tampered calldata as a real tx (default: eth_call only)
  -h, --help              Show help

Cases:
  4a: Tampered proof bytes (random nibble flip in updateMsg)
  4b: Donor proof + base appHash (witness hash mismatch)
  4c: Swap two signer pubkeys in calldata (proof/pubkey mismatch)

Notes:
  - Default mode is safe: uses eth_call only.
  - --send requires ETH_PRIVATE_KEY in env file or current shell env.
EOF
}

require_cmd() {
  local cmd="$1"
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "missing required command: $cmd" >&2
    exit 1
  fi
}

dec_to_hex() {
  local dec="$1"
  printf '0x%x' "$dec"
}

to_lower() {
  echo "$1" | tr '[:upper:]' '[:lower:]'
}

normalize_addr() {
  local addr="$1"
  if [[ ! "$addr" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    echo "invalid address: $addr" >&2
    exit 1
  fi
  echo "$addr"
}

run_eth_call_and_expect_revert() {
  local calldata="$1"
  local call_cmd=(cast call "$ICS26_ADDRESS" --rpc-url "$RPC_URL" --data "$calldata")
  if [[ -n "$FROM_ADDRESS" ]]; then
    call_cmd+=(--from "$FROM_ADDRESS")
  fi

  echo "Running eth_call with tampered updateClient calldata..."
  set +e
  local call_output
  call_output="$("${call_cmd[@]}" 2>&1)"
  local call_exit=$?
  set -e

  if (( call_exit == 0 )); then
    echo "unexpected: tampered eth_call succeeded"
    echo "$call_output"
    exit 2
  fi

  echo "PASS: tampered eth_call reverted as expected"
  echo "$call_output"
}

auto_find_latest_update_tx() {
  local latest_block from_block from_block_hex topic filter logs_json tx
  latest_block=$(cast block-number --rpc-url "$RPC_URL")
  from_block=$((latest_block - LOOKBACK))
  if (( from_block < 0 )); then
    from_block=0
  fi
  from_block_hex=$(dec_to_hex "$from_block")

  topic=$(cast sig-event "ICS02ClientUpdated(string,uint8)")
  filter=$(jq -cn \
    --arg address "$ICS26_ADDRESS" \
    --arg topic "$topic" \
    --arg from "$from_block_hex" \
    '{"address":$address,"fromBlock":$from,"toBlock":"latest","topics":[ $topic ]}')
  logs_json=$(cast rpc eth_getLogs "$filter" --rpc-url "$RPC_URL")
  tx=$(echo "$logs_json" | jq -r 'if type=="array" and length>0 then .[-1].transactionHash else empty end')
  echo "$tx"
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --case)
      CASE_MODE="$2"
      shift 2
      ;;
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
    --update-tx)
      UPDATE_TX="$2"
      shift 2
      ;;
    --donor-tx)
      DONOR_TX="$2"
      shift 2
      ;;
    --lookback)
      LOOKBACK="$2"
      shift 2
      ;;
    --from)
      FROM_ADDRESS="$2"
      shift 2
      ;;
    --offset)
      MUTATE_OFFSET="$2"
      shift 2
      ;;
    --send)
      SEND_MODE=1
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

if [[ "$CASE_MODE" != "4a" && "$CASE_MODE" != "4b" && "$CASE_MODE" != "4c" ]]; then
  echo "--case must be one of: 4a, 4b, 4c" >&2
  exit 1
fi

require_cmd jq
require_cmd cast
if [[ "$CASE_MODE" == "4b" || "$CASE_MODE" == "4c" ]]; then
  require_cmd go
fi

if [[ ! -f "$CONFIG_FILE" ]]; then
  echo "config file not found: $CONFIG_FILE" >&2
  exit 1
fi

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

if [[ -z "$FROM_ADDRESS" && -n "${ETH_PRIVATE_KEY:-}" ]]; then
  FROM_ADDRESS=$(cast wallet address --private-key "$ETH_PRIVATE_KEY")
fi
if [[ -n "$FROM_ADDRESS" ]]; then
  FROM_ADDRESS=$(normalize_addr "$FROM_ADDRESS")
fi

echo "CASE_MODE=$CASE_MODE"
echo "RPC_URL=$RPC_URL"
echo "ICS26_ADDRESS=$ICS26_ADDRESS"
if [[ -n "$FROM_ADDRESS" ]]; then
  echo "FROM_ADDRESS=$FROM_ADDRESS"
fi

TAMPERED_CALLDATA=""

if [[ "$CASE_MODE" == "4a" ]]; then
  if [[ -z "$UPDATE_TX" ]]; then
    UPDATE_TX=$(auto_find_latest_update_tx)
  fi
  if [[ -z "$UPDATE_TX" ]]; then
    echo "no ICS02ClientUpdated tx found in lookback=$LOOKBACK blocks" >&2
    exit 1
  fi
  if [[ ! "$UPDATE_TX" =~ ^0x[0-9a-fA-F]{64}$ ]]; then
    echo "invalid update tx hash: $UPDATE_TX" >&2
    exit 1
  fi

  echo "UPDATE_TX=$UPDATE_TX"

  TX_JSON=$(cast rpc eth_getTransactionByHash "$UPDATE_TX" --rpc-url "$RPC_URL")
  TO_ADDR=$(echo "$TX_JSON" | jq -r '.to // empty')
  INPUT=$(echo "$TX_JSON" | jq -r '.input // empty')

  if [[ -z "$INPUT" || "$INPUT" == "0x" ]]; then
    echo "empty calldata for tx: $UPDATE_TX" >&2
    exit 1
  fi
  if [[ -n "$TO_ADDR" && "$(to_lower "$TO_ADDR")" != "$(to_lower "$ICS26_ADDRESS")" ]]; then
    echo "warning: tx.to ($TO_ADDR) != ICS26_ADDRESS ($ICS26_ADDRESS)" >&2
  fi

  DECODED=$(cast decode-calldata --json "updateClient(string,bytes)" "$INPUT")
  CLIENT_ID=$(echo "$DECODED" | jq -r '.[0]')
  UPDATE_MSG=$(echo "$DECODED" | jq -r '.[1]')

  if [[ ! "$UPDATE_MSG" =~ ^0x[0-9a-fA-F]+$ ]]; then
    echo "decoded updateMsg is not hex bytes" >&2
    exit 1
  fi

  MSG_BODY="${UPDATE_MSG#0x}"
  MSG_LEN=${#MSG_BODY}
  if (( MSG_LEN < 8 )); then
    echo "updateMsg too short to tamper" >&2
    exit 1
  fi

  if [[ -n "$MUTATE_OFFSET" ]]; then
    if ! [[ "$MUTATE_OFFSET" =~ ^[0-9]+$ ]]; then
      echo "invalid --offset, must be integer nibble index in updateMsg body" >&2
      exit 1
    fi
    if (( MUTATE_OFFSET < 0 || MUTATE_OFFSET >= MSG_LEN )); then
      echo "offset out of range: $MUTATE_OFFSET (msg nibbles: $MSG_LEN)" >&2
      exit 1
    fi
    POS="$MUTATE_OFFSET"
  else
    POS=$((MSG_LEN / 2))
  fi

  OLD_NIBBLE="${MSG_BODY:$POS:1}"
  if [[ "$OLD_NIBBLE" == "0" ]]; then
    NEW_NIBBLE="1"
  else
    NEW_NIBBLE="0"
  fi

  TAMPERED_BODY="${MSG_BODY:0:$POS}${NEW_NIBBLE}${MSG_BODY:$((POS+1))}"
  TAMPERED_UPDATE_MSG="0x${TAMPERED_BODY}"
  TAMPERED_CALLDATA=$(cast calldata "updateClient(string,bytes)" "$CLIENT_ID" "$TAMPERED_UPDATE_MSG")

  echo "CLIENT_ID=$CLIENT_ID"
  echo "MUTATE_POS=$POS"
  echo "MUTATE_${OLD_NIBBLE}->${NEW_NIBBLE}"
  echo "TAMPERED_CALLDATA_LEN=${#TAMPERED_CALLDATA}"
else
  HELPER_CMD=(go run ./cmd/updateclient_attack --mode "$CASE_MODE" --rpc-url "$RPC_URL" --router "$ICS26_ADDRESS" --lookback "$LOOKBACK")
  if [[ -n "$UPDATE_TX" ]]; then
    HELPER_CMD+=(--base-tx "$UPDATE_TX")
  fi
  if [[ -n "$DONOR_TX" ]]; then
    HELPER_CMD+=(--donor-tx "$DONOR_TX")
  fi

  ATTACK_JSON=$(
    cd relayer
    "${HELPER_CMD[@]}"
  )

  echo "$ATTACK_JSON" | jq '{mode, base_tx, donor_tx, client_id, mutation, base_app_hash, donor_app_hash}'

  UPDATE_TX=$(echo "$ATTACK_JSON" | jq -r '.base_tx')
  DONOR_TX=$(echo "$ATTACK_JSON" | jq -r '.donor_tx // empty')
  CLIENT_ID=$(echo "$ATTACK_JSON" | jq -r '.client_id')
  MUTATION=$(echo "$ATTACK_JSON" | jq -r '.mutation')
  TAMPERED_CALLDATA=$(echo "$ATTACK_JSON" | jq -r '.tampered_calldata')

  echo "UPDATE_TX=$UPDATE_TX"
  if [[ -n "$DONOR_TX" ]]; then
    echo "DONOR_TX=$DONOR_TX"
  fi
  echo "CLIENT_ID=$CLIENT_ID"
  echo "MUTATION=$MUTATION"
  echo "TAMPERED_CALLDATA_LEN=${#TAMPERED_CALLDATA}"
fi

if [[ -z "$TAMPERED_CALLDATA" || "$TAMPERED_CALLDATA" == "null" ]]; then
  echo "failed to produce tampered calldata" >&2
  exit 1
fi

run_eth_call_and_expect_revert "$TAMPERED_CALLDATA"

if (( SEND_MODE == 1 )); then
  if [[ -z "${ETH_PRIVATE_KEY:-}" ]]; then
    echo "--send requested but ETH_PRIVATE_KEY is empty" >&2
    exit 1
  fi

  echo "Broadcasting tampered tx (--send)..."
  cast send "$ICS26_ADDRESS" \
    --rpc-url "$RPC_URL" \
    --private-key "$ETH_PRIVATE_KEY" \
    --data "$TAMPERED_CALLDATA"
fi
