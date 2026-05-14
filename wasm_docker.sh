#!/bin/sh

set -euxo pipefail

CONTAINER_NAME="ibc-wasm-simd"
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
WASM_FILE="${WASM_FILE:-$SCRIPT_DIR/e2e/interchaintestv8/wasm/cw_ics08_wasm_eth.wasm.gz}"

if [ ! -f "$WASM_FILE" ]; then
  echo "error: wasm file not found: $WASM_FILE" >&2
  exit 1
fi

# Copy wasm into container
docker cp "$WASM_FILE" $CONTAINER_NAME:/root/cw_ics08_wasm_eth.wasm.gz

# Submit store-code proposal through the current simd ibc-wasm CLI
docker exec $CONTAINER_NAME simd tx ibc-wasm store-code /root/cw_ics08_wasm_eth.wasm.gz \
  --from test1 \
  --keyring-backend test \
  --gas 200000000 \
  --fees 200000000stake \
  --chain-id test-ibc-eth \
  --title ibc-eureka \
  --summary ibc-eureka \
  --deposit 10000000stake \
  -y

sleep 5

# Vote yes (use test account which has staking power as validator)
docker exec $CONTAINER_NAME simd tx gov vote 1 yes \
  --from test --keyring-backend test --fees 1000stake \
  --chain-id test-ibc-eth -y

echo "Waiting 35s for proposal to pass..."
sleep 35

# Get the checksum
CHECKSUM=$(docker exec $CONTAINER_NAME simd q ibc-wasm checksums -o json | grep -o '"[a-f0-9]\{64\}"' | head -1 | tr -d '"')

echo ""
echo "=== Wasm code stored successfully ==="
echo "Checksum: 0x$CHECKSUM"
