#!/bin/sh
set -eu

if [ "$#" -ne 3 ]; then
  echo "usage: $0 <optimized.wasm> <proposal-title> <summary>" >&2
  exit 2
fi

artifact_path="$1"
proposal_title="$2"
proposal_summary="$3"
if [ ! -f "$artifact_path" ]; then
  echo "missing optimized artifact: $artifact_path" >&2
  echo "run VALIDATE_OPTIMIZED=1 scripts/validate-l2-clients.sh first" >&2
  exit 1
fi

. ./scripts/local/gaiad_binary.sh
. ./scripts/local/wasm_checksum.sh

chain_id="test-ibc-eth"
keyring="test"
signer="cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn"
wasm_byte_code="$(gzip -n -c "$artifact_path" | base64 -w 0)"

jq -n \
  --arg signer "$signer" \
  --arg wasm_byte_code "$wasm_byte_code" \
  --arg title "$proposal_title" \
  --arg summary "$proposal_summary" \
  '{
    messages: [{
      "@type": "/ibc.lightclients.wasm.v1.MsgStoreCode",
      signer: $signer,
      wasm_byte_code: $wasm_byte_code
    }],
    title: $title,
    deposit: "10000000stake",
    summary: $summary,
    expedited: false
  }' > proposal.json

"$GAIAD" tx gov submit-proposal proposal.json \
  --from val1 \
  --home "$HOME/.gaia" \
  --chain-id "$chain_id" \
  --keyring-backend "$keyring" \
  --gas 200000000 \
  --gas-prices 1stake \
  -y

sleep 5
proposal_id="$(
  "$GAIAD" q gov proposals -o json \
    | jq -r '.proposals | sort_by(.id | tonumber) | last | .id'
)"
sleep 5
vote_default_validators "$proposal_id"
sleep 30

checksum="$(wasm_checksum_from_proposal proposal.json)"
assert_wasm_checksum_stored "$checksum"
echo "Checksum: 0x$checksum"
