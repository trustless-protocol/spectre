#!/usr/bin/env bash

# Local Avalanche devnet for the Cosmos<->C-Chain warp path: the canonical
# FIVE-validator local network (--network-id=local), each node running one of
# avalanchego's published local staking + BLS signer identities.
#
# Five real validators are load-bearing, not ceremony: the warp light client
# verifies a >=67% stake-weighted aggregate BLS signature per update, and the
# signature-aggregator collects those signatures from the registered validators
# over p2p. A single node with sybil protection off is NOT one of the genesis
# validators — the aggregator cannot even bootstrap against it (verified), so
# warp signing needs the actual local validator set live.
#
# Every block is accepted (= final) within ~2s. The v1.15 local genesis runs
# the Granite header shape (no Helicon settled fields); the relayer's proof
# pairing reads that from the headers themselves either way.
#
# Writes .avalanche-devnet-run/attestor.env with the endpoints and the
# well-known prefunded local key, mirroring the other stacks' handoff files
# (consumed by deploy_l2_contracts.sh with DST_CHAIN=avalanche).
#
# Env:
#   AVALANCHEGO_BIN   path to an avalanchego binary (REQUIRED; pin >= v1.14)
#   AVALANCHEGO_TAG   tag whose canonical local keys are fetched when absent
#                     (default v1.15.0; the keys are public test data)
#   DATA_DIR          node database root (default: .avalanche-devnet-run/db)
#   HTTP_PORT         node1 HTTP port (default: 9650); node i uses +2(i-1),
#                     staking ports are HTTP+1 per node

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$(pwd)

: "${AVALANCHEGO_BIN:?set AVALANCHEGO_BIN to an avalanchego binary (>= v1.14)}"
AVALANCHEGO_TAG=${AVALANCHEGO_TAG:-v1.15.0}
RUN_DIR=$REPO_ROOT/.avalanche-devnet-run
DATA_DIR=${DATA_DIR:-$RUN_DIR/db}
HTTP_PORT=${HTTP_PORT:-9650}
KEY_DIR=$RUN_DIR/keys
C_RPC="http://127.0.0.1:${HTTP_PORT}/ext/bc/C/rpc"
C_WS="ws://127.0.0.1:${HTTP_PORT}/ext/bc/C/ws"

mkdir -p "$RUN_DIR" "$DATA_DIR" "$KEY_DIR"

# The canonical local-network validator identities (public test keys shipped in
# the avalanchego repo). Fetched once, pinned to the tag.
for i in 1 2 3 4 5; do
    for file in "staker$i.crt" "staker$i.key" "signer$i.key"; do
        if [ ! -s "$KEY_DIR/$file" ]; then
            curl -sfL -m 30 -o "$KEY_DIR/$file" \
                "https://raw.githubusercontent.com/ava-labs/avalanchego/$AVALANCHEGO_TAG/staking/local/$file" ||
                { echo "ERROR: failed to fetch canonical local key $file" >&2; exit 1; }
        fi
    done
done

rm -f "$RUN_DIR"/avalanchego*.pid

start_node() {
    local index=$1 http=$2 staking=$3 bootstrap_ips=$4 bootstrap_ids=$5
    "$AVALANCHEGO_BIN" \
        --network-id=local \
        --http-port="$http" \
        --staking-port="$staking" \
        --public-ip=127.0.0.1 \
        --data-dir="$DATA_DIR/node$index" \
        --staking-tls-cert-file="$KEY_DIR/staker$index.crt" \
        --staking-tls-key-file="$KEY_DIR/staker$index.key" \
        --staking-signer-key-file="$KEY_DIR/signer$index.key" \
        --bootstrap-ips="$bootstrap_ips" \
        --bootstrap-ids="$bootstrap_ids" \
        --log-level=info \
        >"$RUN_DIR/avalanchego$index.log" 2>&1 &
    echo $! > "$RUN_DIR/avalanchego$index.pid"
}

node_alive() {
    kill -0 "$(cat "$RUN_DIR/avalanchego$1.pid")" 2>/dev/null
}

start_node 1 "$HTTP_PORT" "$((HTTP_PORT + 1))" "" ""
echo "avalanchego node1 started; log: $RUN_DIR/avalanchego1.log"

# Node1's id anchors the other four's bootstrap.
NODE1_ID=""
for _ in $(seq 1 30); do
    NODE1_ID=$(curl -sf -m 2 -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","method":"info.getNodeID","params":{},"id":1}' \
        "http://127.0.0.1:${HTTP_PORT}/ext/info" | sed -n 's/.*"nodeID":"\([^"]*\)".*/\1/p' || true)
    [ -n "$NODE1_ID" ] && break
    node_alive 1 || { echo "ERROR: node1 exited; see $RUN_DIR/avalanchego1.log" >&2; exit 1; }
    sleep 2
done
[ -n "$NODE1_ID" ] || { echo "ERROR: node1 did not answer info.getNodeID" >&2; exit 1; }
echo "node1 id: $NODE1_ID"

for i in 2 3 4 5; do
    start_node "$i" "$((HTTP_PORT + 2 * (i - 1)))" "$((HTTP_PORT + 2 * (i - 1) + 1))" \
        "127.0.0.1:$((HTTP_PORT + 1))" "$NODE1_ID"
done
echo "nodes 2-5 started"

# Poll until the C-Chain answers with the local chain id (0xa868 = 43112) and
# node1 sees the other four validators.
for attempt in $(seq 1 90); do
    CHAIN_ID=$(curl -sf -m 2 -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}' \
        "$C_RPC" | sed -n 's/.*"result":"\(0x[0-9a-f]*\)".*/\1/p' || true)
    PEERS=$(curl -sf -m 2 -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","method":"info.peers","params":{},"id":1}' \
        "http://127.0.0.1:${HTTP_PORT}/ext/info" | grep -o '"nodeID"' | wc -l || true)
    if [ "$CHAIN_ID" = "0xa868" ] && [ "${PEERS:-0}" -ge 4 ]; then
        echo "C-Chain is up at $C_RPC (chain id 43112); node1 peers: $PEERS"
        break
    fi
    for i in 1 2 3 4 5; do
        node_alive "$i" || { echo "ERROR: node$i exited; see $RUN_DIR/avalanchego$i.log" >&2; exit 1; }
    done
    sleep 2
    if [ "$attempt" = 90 ]; then
        echo "ERROR: local network did not become healthy in 180s" >&2
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
echo "Next: run the signature-aggregator sidecar (docs/E2E.md), then DST_CHAIN=avalanche ./scripts/local/deploy_l2_contracts.sh"
