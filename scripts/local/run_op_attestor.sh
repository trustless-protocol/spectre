#!/usr/bin/env bash

# End-to-end bring-up of the OP Stack attestor (attestor/GUIDE.md):
#   [op-reth + op-node replica] ← RPC ← [attestor] → metrics + gRPC feed
#
# Two modes:
#   1. Existing replica:  OP_NODE_RPC_URL=http://host:9545 ./scripts/local/run_op_attestor.sh
#      (replica bring-up skipped; only L1_RPC_URL is needed besides it)
#      Local devnet shortcut: when OP_NODE_RPC_URL is unset and
#      .op-devnet-run/attestor.env exists (written by run_op_stack.sh), it is
#      sourced automatically — a bare ./scripts/local/run_op_attestor.sh
#      attaches to the local devnet.
#   2. Full bring-up:     provide OP_RETH_BIN + OP_NODE_BIN (built per attestor/GUIDE.md, Replica setup)
#      plus L1_RPC_URL + L1_BEACON_URL; the script starts both and waits for
#      optimism_syncStatus before pointing the attestor at it.
#
# Required env:
#   L1_RPC_URL        Ethereum L1 execution RPC (both modes)
#   OP_NODE_RPC_URL   existing replica op-node RPC (mode 1), OR
#   OP_RETH_BIN, OP_NODE_BIN, L1_BEACON_URL (mode 2)
#
# Optional env (defaults in parentheses):
#   NETWORK (op-mainnet)            op-node --network / chain preset
#   SRC_CHAIN ($NETWORK)            src_chain label in the config
#   ATTESTATION_HEAD (finalized)    finalized | safe | unsafe
#   METRICS_PORT (3000)  GRPC_PORT (3001)  OP_NODE_RPC_PORT (9545)  AUTHRPC_PORT (8551)
#   POLL_INTERVAL_SECONDS (30)  LOOKBACK_BLOCKS (600)
#   DERIVED_GAP_BLOCKS (150)        min L2-block gap between derived attestations;
#       lower = fresher feed, larger provisional backlog at safe/unsafe heads
#   OPTIMISM_PORTAL / DISPUTE_GAME_FACTORY / RESPECTED_GAME_TYPE
#       (read fresh from the chain via `cast` when unset; see resolve step)
#   DATADIR ($RUN_DIR/op-reth-data)  L1_RPC_KIND (standard)
#   RUN_DIR (.op-attestor-run)      logs, pids, jwt, config.json, state file
#   REPLICA_WAIT_SECS (300)         max wait for optimism_syncStatus

set -euo pipefail

# All internal paths are repo-root relative — cd up so this script works
# regardless of where it is invoked from.
cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

# Auto-source the local devnet handoff (written by run_op_stack.sh) when the
# caller hasn't wired the env themselves — forgetting to `source` it was a
# recurring foot-gun. An explicit OP_NODE_RPC_URL means the caller brought
# their own environment; leave it untouched.
DEVNET_ENV=$REPO_ROOT/.op-devnet-run/attestor.env
if [ -z "${OP_NODE_RPC_URL:-}" ] && [ -f "$DEVNET_ENV" ]; then
    printf '[run_op_attestor] sourcing %s\n' "$DEVNET_ENV"
    # shellcheck disable=SC1090
    . "$DEVNET_ENV"
fi

NETWORK=${NETWORK:-op-mainnet}
SRC_CHAIN=${SRC_CHAIN:-$NETWORK}
ATTESTATION_HEAD=${ATTESTATION_HEAD:-finalized}
METRICS_PORT=${METRICS_PORT:-3000}
GRPC_PORT=${GRPC_PORT:-3001}
OP_NODE_RPC_PORT=${OP_NODE_RPC_PORT:-9545}
AUTHRPC_PORT=${AUTHRPC_PORT:-8551}
POLL_INTERVAL_SECONDS=${POLL_INTERVAL_SECONDS:-30}
LOOKBACK_BLOCKS=${LOOKBACK_BLOCKS:-600}
DERIVED_GAP_BLOCKS=${DERIVED_GAP_BLOCKS:-150}
L1_RPC_KIND=${L1_RPC_KIND:-standard}
RUN_DIR=${RUN_DIR:-$REPO_ROOT/.op-attestor-run}
REPLICA_WAIT_SECS=${REPLICA_WAIT_SECS:-300}

: "${L1_RPC_URL:?L1_RPC_URL is required (Ethereum L1 execution RPC)}"

mkdir -p "$RUN_DIR"
PIDS=()

log() { printf '\n[run_op_attestor] %s\n' "$*"; }

cleanup() {
    trap - EXIT INT TERM
    log "shutting down..."
    for pid in "${PIDS[@]:-}"; do
        [ -n "$pid" ] && kill "$pid" 2>/dev/null || true
    done
    # Give processes a moment to flush state, then insist.
    for pid in "${PIDS[@]:-}"; do
        [ -n "$pid" ] || continue
        for _ in $(seq 1 20); do
            kill -0 "$pid" 2>/dev/null || break
            sleep 0.5
        done
        kill -9 "$pid" 2>/dev/null || true
    done
}
trap cleanup EXIT INT TERM

# wait_until <timeout-secs> <description> <command...>
# Polls the command every 2s until it succeeds — never a blind sleep.
wait_until() {
    local timeout=$1 desc=$2; shift 2
    local deadline=$((SECONDS + timeout))
    until "$@" >/dev/null 2>&1; do
        if [ "$SECONDS" -ge "$deadline" ]; then
            log "TIMEOUT after ${timeout}s waiting for: $desc"
            return 1
        fi
        sleep 2
    done
}

sync_status() {
    curl -sf -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","id":1,"method":"optimism_syncStatus","params":[]}' \
        "$1" | jq -e '.result.finalized_l2.number'
}

# ---------------------------------------------------------------- replica ---
if [ -n "${OP_NODE_RPC_URL:-}" ]; then
    log "using existing replica op-node at $OP_NODE_RPC_URL (bring-up skipped)"
else
    : "${OP_RETH_BIN:?OP_RETH_BIN is required when OP_NODE_RPC_URL is not set (attestor/GUIDE.md, Replica setup)}"
    : "${OP_NODE_BIN:?OP_NODE_BIN is required when OP_NODE_RPC_URL is not set}"
    : "${L1_BEACON_URL:?L1_BEACON_URL is required when OP_NODE_RPC_URL is not set (must serve blob sidecars)}"
    DATADIR=${DATADIR:-$RUN_DIR/op-reth-data}
    JWT_FILE=$RUN_DIR/jwt.txt
    [ -s "$JWT_FILE" ] || openssl rand -hex 32 > "$JWT_FILE"

    log "starting op-reth (log: $RUN_DIR/op-reth.log)"
    # --rpc.eth-proof-window: optimism_outputAtBlock needs eth_getProof at
    # ~hours-old blocks; 100000 blocks ≈ 2.3 days at 2s (attestor/GUIDE.md, Replica setup).
    "$OP_RETH_BIN" node \
        --chain="$NETWORK" \
        --datadir="$DATADIR" \
        --http --ws \
        --authrpc.port="$AUTHRPC_PORT" \
        --authrpc.jwtsecret="$JWT_FILE" \
        --rpc.eth-proof-window=100000 \
        > "$RUN_DIR/op-reth.log" 2>&1 &
    PIDS+=($!)
    echo $! > "$RUN_DIR/op-reth.pid"

    wait_until 60 "op-reth engine API listening on :$AUTHRPC_PORT" \
        bash -c "exec 3<>/dev/tcp/127.0.0.1/$AUTHRPC_PORT"

    log "starting op-node (log: $RUN_DIR/op-node.log)"
    # Sequencer mode stays off (default): the replica must remain a pure
    # verifier. RPC bound to localhost — keep it private.
    "$OP_NODE_BIN" \
        --l1="$L1_RPC_URL" --l1.rpckind="$L1_RPC_KIND" \
        --l1.beacon="$L1_BEACON_URL" \
        --l2="ws://localhost:$AUTHRPC_PORT" --l2.jwt-secret="$JWT_FILE" --l2.enginekind=reth \
        --network="$NETWORK" \
        --syncmode=execution-layer \
        --rpc.addr=127.0.0.1 --rpc.port="$OP_NODE_RPC_PORT" \
        > "$RUN_DIR/op-node.log" 2>&1 &
    PIDS+=($!)
    echo $! > "$RUN_DIR/op-node.pid"
    OP_NODE_RPC_URL="http://127.0.0.1:$OP_NODE_RPC_PORT"
fi

log "waiting for optimism_syncStatus at $OP_NODE_RPC_URL (max ${REPLICA_WAIT_SECS}s)"
if wait_until "$REPLICA_WAIT_SECS" "optimism_syncStatus" sync_status "$OP_NODE_RPC_URL"; then
    curl -s -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","id":1,"method":"optimism_syncStatus","params":[]}' \
        "$OP_NODE_RPC_URL" \
        | jq '.result | {unsafe: .unsafe_l2.number, safe: .safe_l2.number, finalized: .finalized_l2.number}'
    log "replica heads above; a syncing replica is fine — the attestor holds games in pending until heads cover them"
else
    log "WARNING: replica not answering optimism_syncStatus — continuing anyway;"
    log "the attestor will ingest proposals and keep everything pending (designed fail-safe posture)"
fi

# --------------------------------------------- factory + respected game type ---
# Read both fresh from the chain (governance can change them; attestor/GUIDE.md).
if [ -z "${DISPUTE_GAME_FACTORY:-}" ] || [ -z "${RESPECTED_GAME_TYPE:-}" ]; then
    if [ -z "${OPTIMISM_PORTAL:-}" ] && [ "$NETWORK" = "op-mainnet" ]; then
        OPTIMISM_PORTAL=0xbEb5Fc579115071764c7423A4f12eDde41f106Ed
    fi
    if command -v cast >/dev/null 2>&1 && [ -n "${OPTIMISM_PORTAL:-}" ]; then
        log "reading DisputeGameFactory + respectedGameType from OptimismPortal $OPTIMISM_PORTAL"
        DISPUTE_GAME_FACTORY=${DISPUTE_GAME_FACTORY:-$(cast call "$OPTIMISM_PORTAL" 'disputeGameFactory()(address)' --rpc-url "$L1_RPC_URL")}
        RESPECTED_GAME_TYPE=${RESPECTED_GAME_TYPE:-$(cast call "$OPTIMISM_PORTAL" 'respectedGameType()(uint32)' --rpc-url "$L1_RPC_URL")}
    elif [ "$NETWORK" = "op-mainnet" ]; then
        # Documented values (attestor/GUIDE.md, measured 2026-07-22) — chain read preferred.
        log "WARNING: cast unavailable — falling back to documented OP Mainnet values"
        DISPUTE_GAME_FACTORY=${DISPUTE_GAME_FACTORY:-0xe5965Ab5962eDc7477C8520243A95517CD252fA9}
        RESPECTED_GAME_TYPE=${RESPECTED_GAME_TYPE:-8}
    else
        log "ERROR: set DISPUTE_GAME_FACTORY and RESPECTED_GAME_TYPE (or OPTIMISM_PORTAL + cast) for network $NETWORK"
        exit 1
    fi
fi
log "factory=$DISPUTE_GAME_FACTORY respected_game_type=$RESPECTED_GAME_TYPE"

# ------------------------------------------------------------- attestor ---
log "building the attestor binary"
if command -v just >/dev/null 2>&1; then
    just build-attestor
else
    (cd attestor && go build -o attestor ./cmd)
fi

CONFIG=$RUN_DIR/config.json
STATE_PATH=$RUN_DIR/$SRC_CHAIN.attested-roots.json
cat > "$CONFIG" <<EOF
{
  "modules": [
    {
      "name": "op_source",
      "src_chain": "$SRC_CHAIN",
      "config": {
        "l1_rpc_url": "$L1_RPC_URL",
        "op_node_rpc_url": "$OP_NODE_RPC_URL",
        "dispute_game_factory": "$DISPUTE_GAME_FACTORY",
        "respected_game_type": $RESPECTED_GAME_TYPE,
        "attestation_head": "$ATTESTATION_HEAD",
        "poll_interval_seconds": $POLL_INTERVAL_SECONDS,
        "derived_attestation_gap_blocks": $DERIVED_GAP_BLOCKS,
        "state_path": "$STATE_PATH",
        "l1_bootstrap_lookback_blocks": $LOOKBACK_BLOCKS
      }
    }
  ],
  "server": { "address": "127.0.0.1", "port": $METRICS_PORT, "grpc_port": $GRPC_PORT }
}
EOF
log "config written to $CONFIG"

log "starting the attestor (log: $RUN_DIR/attestor.log)"
./attestor/attestor --config "$CONFIG" > "$RUN_DIR/attestor.log" 2>&1 &
ATTESTOR_PID=$!
PIDS+=("$ATTESTOR_PID")
echo "$ATTESTOR_PID" > "$RUN_DIR/attestor.pid"

wait_until 30 "attestor metrics on :$METRICS_PORT" \
    curl -sf "http://127.0.0.1:$METRICS_PORT/metrics"
wait_until 30 "attestor gRPC on :$GRPC_PORT" \
    bash -c "exec 3<>/dev/tcp/127.0.0.1/$GRPC_PORT"

log "attestor is up. Useful commands:"
cat <<EOF

  tail -f $RUN_DIR/attestor.log
  curl -s localhost:$METRICS_PORT/metrics | grep fast_ibc_op_attestor
  grpcurl -plaintext -d '{}' 127.0.0.1:$GRPC_PORT attestor.AttestorService/Info
  grpcurl -plaintext -d '{"src_chain":"$SRC_CHAIN"}' 127.0.0.1:$GRPC_PORT attestor.AttestorService/AttestedUpTo

Healthy operation (attestor/GUIDE.md): 'ingested proposal: game ...' during the
bootstrap window, then replica heads climbing on /metrics, a derived root
every ~5 min once the gating head is covered, and 'attested game N' ~hourly.
Ctrl-C stops the attestor$( [ -f "$RUN_DIR/op-node.pid" ] && echo " and the replica" ).

EOF

tail -n +1 -f "$RUN_DIR/attestor.log" &
PIDS+=($!)
wait "$ATTESTOR_PID"
