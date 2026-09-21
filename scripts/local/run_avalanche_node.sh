#!/usr/bin/env bash

# Local Avalanche devnet for the Cosmos<->C-Chain path: a single avalanchego
# node on --network-id=local with sybil protection off. That yields an
# operational C-Chain (chain id 43112) where every block is accepted (= final)
# within a second — all this path needs, since the attested-header client
# verifies no consensus artifact. No subnets, no multi-node tmpnet.
#
# Writes .avalanche-devnet-run/attestor.env with the endpoints and the
# well-known prefunded local key, mirroring the OP/Base/Arbitrum handoff files
# consumed by deploy_l2_contracts.sh (DST_CHAIN=avalanche) and
# run_avalanche_attestor.sh.
#
# Env:
#   AVALANCHEGO_BIN   path to an avalanchego binary (REQUIRED; pin >= v1.14 so
#                     the C-Chain runs the Granite header layout the client
#                     understands — the coreth reader also accepts Helicon)
#   DATA_DIR          node database/config dir (default: .avalanche-devnet-run/db)
#   HTTP_PORT         node HTTP port (default: 9650)

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$(pwd)

: "${AVALANCHEGO_BIN:?set AVALANCHEGO_BIN to an avalanchego binary (>= v1.14)}"
RUN_DIR=$REPO_ROOT/.avalanche-devnet-run
DATA_DIR=${DATA_DIR:-$RUN_DIR/db}
HTTP_PORT=${HTTP_PORT:-9650}
C_RPC="http://127.0.0.1:${HTTP_PORT}/ext/bc/C/rpc"
C_WS="ws://127.0.0.1:${HTTP_PORT}/ext/bc/C/ws"

mkdir -p "$RUN_DIR" "$DATA_DIR"

"$AVALANCHEGO_BIN" \
  --network-id=local \
  --sybil-protection-enabled=false \
  --http-port="$HTTP_PORT" \
  --staking-port=0 \
  --public-ip=127.0.0.1 \
  --data-dir="$DATA_DIR" \
  --log-level=info \
  >"$RUN_DIR/avalanchego.log" 2>&1 &
NODE_PID=$!
echo "$NODE_PID" > "$RUN_DIR/avalanchego.pid"
echo "avalanchego started (pid $NODE_PID); log: $RUN_DIR/avalanchego.log"

# Poll the C-Chain until it answers with the local chain id (0xa868 = 43112).
# Bootstrap of a fresh single local node takes a few seconds.
for i in $(seq 1 60); do
    CHAIN_ID=$(curl -sf -m 2 -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
        "$C_RPC" | sed -n 's/.*"result":"\(0x[0-9a-f]*\)".*/\1/p' || true)
    if [ "$CHAIN_ID" = "0xa868" ]; then
        echo "C-Chain is up at $C_RPC (chain id 43112)"
        break
    fi
    if ! kill -0 "$NODE_PID" 2>/dev/null; then
        echo "ERROR: avalanchego exited; see $RUN_DIR/avalanchego.log" >&2
        exit 1
    fi
    sleep 2
    if [ "$i" = 60 ]; then
        echo "ERROR: C-Chain did not become healthy in 120s; see $RUN_DIR/avalanchego.log" >&2
        exit 1
    fi
done

# The standard prefunded key of every avalanchego local-network C-Chain genesis
# ("ewoq"). DEVNET ONLY — publicly known, never fund it anywhere real.
cat > "$RUN_DIR/attestor.env" <<EOF
L2_RPC_URL=$C_RPC
L2_WS_URL=$C_WS
L2_CHAIN_ID=43112
L2_DEPLOYER_PRIVATE_KEY=0x56289e99c94b6912bfc12adc093c9b51124f0dc54ac7a766b2bc5ccf558d8027
L2_DEPLOYER_ADDRESS=0x8db97C7cEcE249c2b98bDC0226Cc4C2A57BF52FC
EOF
echo "Wrote $RUN_DIR/attestor.env"
echo "Next: ./scripts/local/run_avalanche_attestor.sh, then DST_CHAIN=avalanche ./scripts/local/deploy_l2_contracts.sh"
