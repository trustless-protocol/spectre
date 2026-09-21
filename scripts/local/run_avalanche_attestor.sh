#!/usr/bin/env bash

# Builds and runs the Avalanche C-Chain attestor sidecar against a C-Chain RPC.
#
# With no env set, attaches to the local devnet: sources
# .avalanche-devnet-run/attestor.env (written by run_avalanche_node.sh).
#
# Env:
#   C_CHAIN_RPC_URL       C-Chain execution RPC (default: local devnet handoff)
#   L2_CHAIN_ID           C-Chain EVM chain id (default: from handoff / 43112)
#   ATTESTOR_SIGNING_KEY  32-byte Ed25519 seed in hex; its public half must be
#                         pinned in the Cosmos client profile. Defaults to the
#                         repo's DEVNET-ONLY sample identity when the local
#                         handoff file is in use (devnet_attestor_identity.sh).
#   SRC_CHAIN             attestor src_chain key (default: avaxdev; must match
#                         the relayer module's attestor_src_chain)
#   GRPC_PORT             sidecar gRPC port (default: 3001)

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$(pwd)
RUN_DIR=$REPO_ROOT/.avalanche-devnet-run

if [ -z "${C_CHAIN_RPC_URL:-}" ] && [ -f "$RUN_DIR/attestor.env" ]; then
    # shellcheck disable=SC1091
    . "$RUN_DIR/attestor.env"
    C_CHAIN_RPC_URL=$L2_RPC_URL
    if [ -z "${ATTESTOR_SIGNING_KEY:-}" ]; then
        # DEVNET-ONLY identity for the avalanche (finalized) tier. It is
        # deliberately distinct from devnet_attestor_identity.sh: the relayer
        # refuses one attestor key reused across different head_kind values,
        # and the OP/Base/Arbitrum devnets run "safe". Public test data —
        # seed = sha256("fast-ibc avalanche devnet attestor identity v1");
        # public half (in the example configs' attestors.public_keys):
        # base64 o+IGHINpSfBe/aq/3e++IpdnKD4pW41m387fjxUy5+4=
        ATTESTOR_SIGNING_KEY=5779646464bd3d379b53b38a776ad5a3c4567e3a5f83199d0660d4807c378ee8
    fi
fi
: "${C_CHAIN_RPC_URL:?set C_CHAIN_RPC_URL or run run_avalanche_node.sh first}"
: "${L2_CHAIN_ID:=43112}"
: "${ATTESTOR_SIGNING_KEY:?set ATTESTOR_SIGNING_KEY (32-byte Ed25519 seed, hex)}"
SRC_CHAIN=${SRC_CHAIN:-avaxdev}
GRPC_PORT=${GRPC_PORT:-3001}

mkdir -p "$RUN_DIR" attestor/avalanche/bin
( cd attestor/avalanche && go build -o bin/attestor ./cmd )

cat > "$RUN_DIR/avalanche-attestor.json" <<EOF
{
    "modules": [
        {
            "name": "avalanche_source",
            "src_chain": "$SRC_CHAIN",
            "config": {
                "c_chain_rpc_url": "$C_CHAIN_RPC_URL",
                "chain_id": $L2_CHAIN_ID,
                "poll_interval_seconds": 2,
                "attestation_signing_key": "env:ATTESTOR_SIGNING_KEY"
            }
        }
    ],
    "server": {
        "address": "127.0.0.1",
        "log_level": "info",
        "grpc_port": $GRPC_PORT
    }
}
EOF

export ATTESTOR_SIGNING_KEY
exec ./attestor/avalanche/bin/attestor --config "$RUN_DIR/avalanche-attestor.json"
