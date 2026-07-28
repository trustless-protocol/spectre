#!/usr/bin/env bash

set -euxo pipefail

cd "$(dirname "$0")/../.."

CHAIN_ID="test-ibc-eth"
KEYRING="test"

# Build L2 wasm clients if artifacts don't already exist
WASM_TARGET="target/wasm32-unknown-unknown/release"
if ! ls "$WASM_TARGET"/cw-ics08-wasm-{arbitrum,base,op}.wasm >/dev/null 2>&1; then
  cargo build --target wasm32-unknown-unknown --release \
    -p cw-ics08-wasm-arbitrum -p cw-ics08-wasm-base -p cw-ics08-wasm-op --locked
fi

# Store each L2 client wasm bytecode via a separate governance proposal,
# then vote and wait for passage. One proposal per L2 because each
# MsgStoreCode stores a distinct wasm blob.
for L2 in arbitrum base op; do
  WASM_FILE="$WASM_TARGET/cw-ics08-wasm-${L2}.wasm"
  WASM_B64=$(gzip -c "$WASM_FILE" | base64 -w0)
  EXPECTED_CHECKSUM=$(sha256sum "$WASM_FILE" | cut -d' ' -f1)

  PROPOSAL_FILE="/tmp/proposal-l2-${L2}.json"
  cat > "$PROPOSAL_FILE" <<EOF
{
  "messages": [
    {
      "@type": "/ibc.lightclients.wasm.v1.MsgStoreCode",
      "signer": "cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
      "wasm_byte_code": "${WASM_B64}"
    }
  ],
  "title": "L2 wasm client: ${L2}",
  "deposit": "10000000stake",
  "summary": "Store cw-ics08-wasm-${L2} bytecode",
  "expedited": false
}
EOF

  # Record the current frontier before submitting. The async transaction may not
  # be indexed immediately, so only a strictly newer proposal can belong to this
  # loop iteration.
  PREVIOUS_PROPOSAL_ID=$(gaiad q gov proposals -o json \
    | jq -r '[.proposals[]? | (.id // .proposal_id) | tonumber] | max // 0')

  gaiad tx gov submit-proposal "$PROPOSAL_FILE" \
    --from val1 \
    --home "$HOME/.gaia" \
    --chain-id "$CHAIN_ID" \
    --keyring-backend "$KEYRING" \
    --gas 200000000 \
    --gas-prices 1stake \
    -y

  # Poll until the newly submitted proposal is indexed.
  PROPOSAL_ID=""
  for attempt in $(seq 1 30); do
    PROPOSAL_ID=$(gaiad q gov proposals -o json 2>/dev/null \
      | jq --argjson previous "$PREVIOUS_PROPOSAL_ID" -r \
        '[.proposals[]? | (.id // .proposal_id) | tonumber | select(. > $previous)] | min // empty')
    if [ -n "$PROPOSAL_ID" ]; then
      break
    fi
    sleep 2
  done

  if [ -z "$PROPOSAL_ID" ]; then
    echo "ERROR: L2-${L2}: no proposal indexed after 60s" >&2
    exit 1
  fi
  echo "L2-${L2} proposal ID: $PROPOSAL_ID"

  # Vote yes from all 3 validators
  for VAL_HOME in "$HOME/.gaia" "$HOME/.gaia-val2" "$HOME/.gaia-val3"; do
    VAL=$(basename "$VAL_HOME" | sed 's/^\.gaia-//; s/^\.gaia$/val1/')
    gaiad tx gov vote "$PROPOSAL_ID" yes \
      --from "$VAL" \
      --home "$VAL_HOME" \
      --chain-id "$CHAIN_ID" \
      --keyring-backend "$KEYRING" \
      --gas-prices 1stake \
      -y
  done

  # Poll until the checksum appears (proposal executed)
  CHECKSUM=""
  for attempt in $(seq 1 60); do
    CHECKSUM=$(gaiad q ibc-wasm checksums -o json 2>/dev/null \
      | jq -r "[.checksums[] | select(. == \"$EXPECTED_CHECKSUM\")][0] // empty")
    if [ -n "$CHECKSUM" ] && [ "$CHECKSUM" != "null" ]; then
      break
    fi
    echo "L2-${L2}: checksum not yet stored (attempt $attempt/60)..."
    sleep 5
  done

  if [ -z "$CHECKSUM" ] || [ "$CHECKSUM" = "null" ]; then
    echo "ERROR: L2-${L2}: wasm checksum never stored" >&2
    exit 1
  fi

  echo "cw-ics08-wasm-${L2}: 0x$CHECKSUM"
done

echo ""
echo "=== All L2 wasm clients stored successfully ==="
