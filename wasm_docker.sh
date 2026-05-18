#!/bin/sh

set -euxo pipefail

CONTAINER_NAME="ibc-wasm-simd"

# Generate proposal.json locally from wasm.sh (lines 6-17 are the JSON body, line 5 is echo '{', line 18 is }' > ...)
echo '{' > /tmp/proposal.json
sed -n '6,17p' wasm.sh >> /tmp/proposal.json
echo '}' >> /tmp/proposal.json

# Copy proposal.json into container
docker cp /tmp/proposal.json $CONTAINER_NAME:/root/proposal.json

# Submit governance proposal
docker exec $CONTAINER_NAME simd tx gov submit-proposal /root/proposal.json \
  --from test1 --keyring-backend test --gas 200000000 --fees 200000000stake \
  --chain-id test-ibc-eth -y

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
