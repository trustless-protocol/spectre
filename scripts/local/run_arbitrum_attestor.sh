#!/usr/bin/env bash

# End-to-end bring-up of the Arbitrum attestor:
#
#   [local L1 + sequencer feed] <- [independent Nitro replica] <- [attestor]
#                                                               |
#                                                               +-> gRPC
#
# The attestor and the Nitro replica it owns run in one Docker container. The
# replica has its own persistent chain database and cannot sequence, post
# batches, or stake. This is the Arbitrum counterpart to run_op_attestor.sh.
#
# Local devnet shortcut: when ROLLUP_CORE_ADDRESS is unset and
# .arbitrum-devnet-run/attestor.env exists (written by
# run_arbitrum_node.sh), that handoff is sourced automatically.
#
# Required env when not using the handoff:
#   L1_RPC_URL, L1_BEACON_URL, L1_CHAIN_ID, L2_CHAIN_ID
#   ROLLUP_CORE_ADDRESS, ROLLUP_DEPLOYMENT_BLOCK
#   NITRO_SEQUENCER_CONFIG
#
# Optional env:
#   SRC_CHAIN (arbdev)                 relayer source-chain label
#   ASSERTIONS_MAPPING_SLOT (0x76)     BoLD _assertions mapping slot
#   ASSERTION_STATUS_OFFSET (25)       BoLD AssertionNode.status byte offset
#   NITRO_FEED_URL (ws://127.0.0.1:9642)
#   ATTESTOR_L1_RPC_URL               container-visible L1 RPC override
#   ATTESTOR_L1_BEACON_URL            container-visible beacon API override
#   ATTESTOR_NITRO_FEED_URL           container-visible sequencer feed override
#   ATTESTOR_DOCKER_NETWORK           optional Docker network to join
#   GRPC_PORT (3001)
#   RUNTIME_POLL_INTERVAL (2s)
#   ASSERTION_POLL_INTERVAL (2s)
#   ASSERTION_MAX_BLOCK_RANGE (2000)
#   ATTESTOR_RUN_DIR (.arbitrum-attestor-run) — legacy RUN_DIR still honoured,
#       but prefer ATTESTOR_RUN_DIR: the devnet bring-up scripts use RUN_DIR for
#       their own artifacts
#   STARTUP_WAIT_SECS (600)
#   ATTESTOR_DOCKER_IMAGE (fast-ibc-arbitrum-attestor:local)
#   ATTESTOR_NITRO_IMAGE (offchainlabs/nitro-node:v3.11.2-3599aca)
#   ATTESTOR_CONTAINER_NAME (fast-ibc-arbitrum-attestor)
#   ATTESTOR_NITRO_VOLUME (fast-ibc-arbitrum-attestor-nitro)
#   DETACH (0)                         exit after readiness and leave container running

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

DEVNET_ENV=$REPO_ROOT/.arbitrum-devnet-run/attestor.env
if [ -z "${ROLLUP_CORE_ADDRESS:-}" ] && [ -f "$DEVNET_ENV" ]; then
    printf '[run_arbitrum_attestor] sourcing %s\n' "$DEVNET_ENV"
    # shellcheck disable=SC1090
    . "$DEVNET_ENV"
fi

SRC_CHAIN=${SRC_CHAIN:-arbdev}
ASSERTIONS_MAPPING_SLOT=${ASSERTIONS_MAPPING_SLOT:-0x0000000000000000000000000000000000000000000000000000000000000076}
ASSERTION_STATUS_OFFSET=${ASSERTION_STATUS_OFFSET:-25}
NITRO_FEED_URL=${NITRO_FEED_URL:-ws://127.0.0.1:9642}
GRPC_PORT=${GRPC_PORT:-3001}
RUNTIME_POLL_INTERVAL=${RUNTIME_POLL_INTERVAL:-2s}
ASSERTION_POLL_INTERVAL=${ASSERTION_POLL_INTERVAL:-2s}
ASSERTION_MAX_BLOCK_RANGE=${ASSERTION_MAX_BLOCK_RANGE:-2000}
# ATTESTOR_RUN_DIR takes precedence over the legacy RUN_DIR: the devnet bring-up
# scripts use RUN_DIR for their OWN artifacts, so a handoff that exports it into
# the caller's shell sends the next bring-up's package clone, downloads and
# handoff into this attestor's directory. Handoffs export ATTESTOR_RUN_DIR
# instead; RUN_DIR still works for anything that already sets it.
RUN_DIR=${ATTESTOR_RUN_DIR:-${RUN_DIR:-$REPO_ROOT/.arbitrum-attestor-run}}
STARTUP_WAIT_SECS=${STARTUP_WAIT_SECS:-600}
ATTESTOR_DOCKER_IMAGE=${ATTESTOR_DOCKER_IMAGE:-fast-ibc-arbitrum-attestor:local}
ATTESTOR_NITRO_IMAGE=${ATTESTOR_NITRO_IMAGE:-offchainlabs/nitro-node:v3.11.2-3599aca}
ATTESTOR_CONTAINER_NAME=${ATTESTOR_CONTAINER_NAME:-fast-ibc-arbitrum-attestor}
ATTESTOR_NITRO_VOLUME=${ATTESTOR_NITRO_VOLUME:-fast-ibc-arbitrum-attestor-nitro}
DETACH=${DETACH:-0}

fail() {
    printf '[run_arbitrum_attestor] ERROR: %s\n' "$*" >&2
    exit 1
}

log() {
    printf '\n[run_arbitrum_attestor] %s\n' "$*"
}

require_command() {
    command -v "$1" >/dev/null 2>&1 || fail "$1 is required"
}

require_command docker
require_command jq

: "${L1_RPC_URL:?L1_RPC_URL is required (Ethereum L1 execution RPC)}"
: "${L1_BEACON_URL:?L1_BEACON_URL is required (Ethereum beacon API)}"
: "${L1_CHAIN_ID:?L1_CHAIN_ID is required}"
: "${L2_CHAIN_ID:?L2_CHAIN_ID is required}"
: "${ROLLUP_CORE_ADDRESS:?ROLLUP_CORE_ADDRESS is required}"
: "${ROLLUP_DEPLOYMENT_BLOCK:?ROLLUP_DEPLOYMENT_BLOCK is required}"
: "${NITRO_SEQUENCER_CONFIG:?NITRO_SEQUENCER_CONFIG is required}"

if [ "$ROLLUP_DEPLOYMENT_BLOCK" = 0 ]; then
    log "ROLLUP_DEPLOYMENT_BLOCK is 0; using 1 because the attestor scans from a non-genesis L1 block"
    ROLLUP_DEPLOYMENT_BLOCK=1
fi

[ -f "$NITRO_SEQUENCER_CONFIG" ] ||
    fail "Nitro sequencer config not found: $NITRO_SEQUENCER_CONFIG"
[[ "$GRPC_PORT" =~ ^[0-9]+$ ]] && [ "$GRPC_PORT" -gt 0 ] && [ "$GRPC_PORT" -le 65535 ] ||
    fail "GRPC_PORT must be between 1 and 65535"
[[ "$ASSERTIONS_MAPPING_SLOT" =~ ^0x[0-9a-fA-F]{64}$ ]] ||
    fail "ASSERTIONS_MAPPING_SLOT must be a 32-byte hexadecimal value"
[[ "$ASSERTION_STATUS_OFFSET" =~ ^[0-9]+$ ]] &&
    [ "$ASSERTION_STATUS_OFFSET" -lt 32 ] ||
    fail "ASSERTION_STATUS_OFFSET must be between 0 and 31"

mkdir -p "$RUN_DIR"
RUN_DIR=$(cd "$RUN_DIR" && pwd)

container_url() {
    local url=$1
    case "$url" in
        http://* | https://* | ws://* | wss://*) ;;
        *) url="http://$url" ;;
    esac
    printf '%s\n' "$url" |
        sed -E 's#(://)(127\.0\.0\.1|localhost)([:/])#\1host.docker.internal\3#'
}

CONTAINER_L1_RPC_URL=${ATTESTOR_L1_RPC_URL:-$(container_url "$L1_RPC_URL")}
CONTAINER_L1_BEACON_URL=${ATTESTOR_L1_BEACON_URL:-$(container_url "$L1_BEACON_URL")}
CONTAINER_NITRO_FEED_URL=${ATTESTOR_NITRO_FEED_URL:-$(container_url "$NITRO_FEED_URL")}

DOCKER_NETWORK_ARGS=()
if [ -n "${ATTESTOR_DOCKER_NETWORK:-}" ]; then
    docker network inspect "$ATTESTOR_DOCKER_NETWORK" >/dev/null 2>&1 ||
        fail "Docker network not found: $ATTESTOR_DOCKER_NETWORK"
    DOCKER_NETWORK_ARGS=(--network "$ATTESTOR_DOCKER_NETWORK")
fi

DOCKER_RM_ARGS=(--rm)
if [ "$DETACH" = 1 ]; then
    DOCKER_RM_ARGS=()
fi

NITRO_CONFIG=$RUN_DIR/nitro-replica.json
jq '
    .node.sequencer = false
    | .execution.sequencer.enable = false
    | .execution["forwarding-target"] = "null"
    | .node["delayed-sequencer"].enable = false
    | .node["batch-poster"].enable = false
    | .node.staker.enable = false
    | .node.feed.output.enable = false
' "$NITRO_SEQUENCER_CONFIG" >"$NITRO_CONFIG"

CONFIG=$RUN_DIR/config.json
jq -n \
    --arg grpc_listen_address "127.0.0.1:50051" \
    --arg runtime_poll_interval "$RUNTIME_POLL_INTERVAL" \
    --arg src_chain "$SRC_CHAIN" \
    --arg l1_rpc_url "$CONTAINER_L1_RPC_URL" \
    --argjson l1_chain_id "$L1_CHAIN_ID" \
    --argjson l2_chain_id "$L2_CHAIN_ID" \
    --arg rollup_core_address "$ROLLUP_CORE_ADDRESS" \
    --arg assertions_mapping_slot "$ASSERTIONS_MAPPING_SLOT" \
    --argjson assertion_status_offset "$ASSERTION_STATUS_OFFSET" \
    --argjson assertion_start_block "$ROLLUP_DEPLOYMENT_BLOCK" \
    --arg assertion_poll_interval "$ASSERTION_POLL_INTERVAL" \
    --argjson assertion_max_block_range "$ASSERTION_MAX_BLOCK_RANGE" \
    --arg parent_chain_url "$CONTAINER_L1_RPC_URL" \
    --arg beacon_url "$CONTAINER_L1_BEACON_URL" \
    --arg feed_url "$CONTAINER_NITRO_FEED_URL" \
    '{
        grpc_listen_address: $grpc_listen_address,
        runtime_poll_interval: $runtime_poll_interval,
        src_chain: $src_chain,
        l1_rpc_url: $l1_rpc_url,
        l1_chain_id: $l1_chain_id,
        l2_chain_id: $l2_chain_id,
        rollup_core_address: $rollup_core_address,
        assertions_mapping_slot: $assertions_mapping_slot,
        assertion_status_offset: $assertion_status_offset,
        assertion_start_block: $assertion_start_block,
        assertion_poll_interval: $assertion_poll_interval,
        assertion_max_block_range: $assertion_max_block_range,
        attestor_state_path: "/var/lib/fast-ibc/nitro/attested-roots.json",
        nitro_binary_path: "/usr/local/bin/nitro",
        nitro_binary_sha256: "0000000000000000000000000000000000000000000000000000000000000000",
        nitro_arguments: [
            "--conf.file=/config/nitro-replica.json",
            "--parent-chain.connection.url=" + $parent_chain_url,
            "--parent-chain.blob-client.beacon-url=" + $beacon_url,
            "--node.feed.input.url=" + $feed_url,
            "--node.feed.output.enable=false",
            "--node.sequencer=false",
            "--execution.sequencer.enable=false",
            "--execution.forwarding-target=null",
            "--node.delayed-sequencer.enable=false",
            "--node.batch-poster.enable=false",
            "--node.staker.enable=false"
        ],
        nitro_work_dir: "/config",
        nitro_data_dir: "/var/lib/fast-ibc/nitro",
        nitro_ipc_path: "/run/fast-ibc/nitro.ipc",
        nitro_startup_timeout: "5m",
        nitro_shutdown_timeout: "30s"
    }' >"$CONFIG"

log "replica config written to $NITRO_CONFIG"
log "attestor config written to $CONFIG"

if docker container inspect "$ATTESTOR_CONTAINER_NAME" >/dev/null 2>&1; then
    fail "container $ATTESTOR_CONTAINER_NAME already exists; stop or rename it before starting another attestor"
fi

log "building $ATTESTOR_DOCKER_IMAGE with Nitro runtime $ATTESTOR_NITRO_IMAGE"
docker build \
    --build-arg "NITRO_IMAGE=$ATTESTOR_NITRO_IMAGE" \
    --file attestor/arbitrum/Dockerfile \
    --tag "$ATTESTOR_DOCKER_IMAGE" \
    attestor

LOG_FOLLOWER_PID=
cleanup() {
    trap - EXIT INT TERM
    log "shutting down the Arbitrum attestor"
    if [ -n "$LOG_FOLLOWER_PID" ]; then
        kill "$LOG_FOLLOWER_PID" >/dev/null 2>&1 || true
    fi
    docker stop --time 40 "$ATTESTOR_CONTAINER_NAME" >/dev/null 2>&1 || true
}
trap cleanup EXIT INT TERM

log "starting $ATTESTOR_CONTAINER_NAME (log: $RUN_DIR/attestor.log)"
docker run \
    --detach \
    ${DOCKER_RM_ARGS[@]+"${DOCKER_RM_ARGS[@]}"} \
    --init \
    --name "$ATTESTOR_CONTAINER_NAME" \
    --stop-timeout 40 \
    --add-host host.docker.internal:host-gateway \
    ${DOCKER_NETWORK_ARGS[@]+"${DOCKER_NETWORK_ARGS[@]}"} \
    --publish "127.0.0.1:$GRPC_PORT:50051" \
    --volume "$RUN_DIR:/config:ro" \
    --volume "$ATTESTOR_NITRO_VOLUME:/var/lib/fast-ibc/nitro" \
    --env ATTESTOR_CONFIG=/config/config.json \
    "$ATTESTOR_DOCKER_IMAGE" >/dev/null

printf '%s\n' "$ATTESTOR_CONTAINER_NAME" >"$RUN_DIR/attestor.container"
docker logs --follow "$ATTESTOR_CONTAINER_NAME" >"$RUN_DIR/attestor.log" 2>&1 &
LOG_FOLLOWER_PID=$!

wait_until() {
    local timeout=$1
    local description=$2
    shift 2
    local deadline=$((SECONDS + timeout))

    until "$@" >/dev/null 2>&1; do
        if ! docker container inspect \
            --format '{{.State.Running}}' "$ATTESTOR_CONTAINER_NAME" 2>/dev/null |
            grep -q true; then
            docker logs "$ATTESTOR_CONTAINER_NAME" >&2 || true
            fail "$ATTESTOR_CONTAINER_NAME exited while waiting for $description"
        fi
        if [ "$SECONDS" -ge "$deadline" ]; then
            docker logs "$ATTESTOR_CONTAINER_NAME" >&2 || true
            fail "timed out after ${timeout}s waiting for $description"
        fi
        sleep 2
    done
}

grpc_ready() {
    (exec 3<>"/dev/tcp/127.0.0.1/$GRPC_PORT")
}

log "waiting for attestor gRPC on 127.0.0.1:$GRPC_PORT"
wait_until "$STARTUP_WAIT_SECS" "attestor gRPC" grpc_ready

log "attestor is ready"
cat <<EOF

  tail -f $RUN_DIR/attestor.log
  grpcurl -plaintext -d '{}' 127.0.0.1:$GRPC_PORT attestor.AttestorService/Info
  grpcurl -plaintext -d '{"src_chain":"$SRC_CHAIN","include_provisional":true}' \
    127.0.0.1:$GRPC_PORT attestor.AttestorService/AttestedUpTo

The verifier Nitro database is persisted in Docker volume:
  $ATTESTOR_NITRO_VOLUME

Ctrl-C stops the attestor container without deleting that volume.

EOF

if [ "$DETACH" = 1 ]; then
    trap - EXIT INT TERM
    log "detached; container $ATTESTOR_CONTAINER_NAME remains running"
    exit 0
fi

CONTAINER_EXIT_STATUS=$(docker wait "$ATTESTOR_CONTAINER_NAME")
[ "$CONTAINER_EXIT_STATUS" -eq 0 ] ||
    fail "$ATTESTOR_CONTAINER_NAME exited with status $CONTAINER_EXIT_STATUS"
