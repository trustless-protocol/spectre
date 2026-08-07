#!/usr/bin/env bash

# End-to-end bring-up of the Arbitrum attestor:
#
#   [local L1] --------> [attestor] <-------- [Nitro HTTP + WebSocket]
#                           |
#                           +-> gRPC
#
# The attestor uses external Nitro endpoints and does not launch or own a Nitro
# process. This is the Arbitrum counterpart to run_op_attestor.sh.
#
# Local devnet shortcut: when ROLLUP_CORE_ADDRESS is unset and
# .arbitrum-devnet-run/attestor.env exists (written by
# run_arbitrum_node.sh), that handoff is sourced automatically.
#
# Required env when not using the handoff:
#   L1_RPC_URL, L1_CHAIN_ID, L2_RPC_URL, L2_WS_URL, L2_CHAIN_ID
#   ROLLUP_CORE_ADDRESS, ROLLUP_DEPLOYMENT_BLOCK
#
# Optional env:
#   CHAIN_PROFILE (devnet)             devnet, or a chain whose reference config
#       exists at attestor/arbitrum/config.<profile>.json (e.g. arbitrum-sepolia).
#       Supplies src_chain, the BoLD slot and offset, poll intervals, scan range
#       and start block. Every one is overridable individually below.
#   ASSERTION_START_BLOCK              first L1 block of the assertion scan
#   SRC_CHAIN                          relayer source-chain label
#   ASSERTIONS_MAPPING_SLOT            BoLD _assertions mapping slot
#   ASSERTION_STATUS_OFFSET            BoLD AssertionNode.status byte offset
#   ATTESTOR_L1_RPC_URL                container-visible L1 RPC override
#   ATTESTOR_NITRO_RPC_URL             container-visible Nitro RPC override
#   ATTESTOR_NITRO_WS_URL              container-visible Nitro WebSocket override
#   ATTESTOR_DOCKER_NETWORK            optional Docker network to join
#   GRPC_PORT (3001)
#   RUNTIME_POLL_INTERVAL (2s)
#   ATTESTATION_HEAD (finalized)       Nitro head: unsafe|safe|finalized
#   DISABLE_DERIVED_ROOTS (false)      attest only RollupCore-backed assertions
#       when true. The attestor-trusted L2 client verifies no assertion, so the
#       relayer does not need one, and an assertion lands long after the block it
#       covers — assertions-only leaves the frontier hours behind for no gain.
#       Matches the OP attestor default.
#   DERIVED_ATTESTATION_GAP_BLOCKS (150)
#       Minimum L2-block gap between derived attestations, and therefore the floor
#       on return-direction latency: a packet waits until the frontier reaches its
#       block, and the frontier moves one derived root per gap. Lower it for a
#       devnet or an E2E; see run_op_attestor.sh for the same knob on OP.
#   MAX_DERIVED_ROOTS (1000)
#   ASSERTION_POLL_INTERVAL (2s)
#   ASSERTION_MAX_BLOCK_RANGE (2000)
#   ATTESTOR_RUN_DIR (.arbitrum-attestor-run) — legacy RUN_DIR still honoured,
#       but prefer ATTESTOR_RUN_DIR: the devnet bring-up scripts use RUN_DIR for
#       their own artifacts
#   STARTUP_WAIT_SECS (600)
#   ATTESTOR_DOCKER_IMAGE (fast-ibc-arbitrum-attestor:local)
#   ATTESTOR_CONTAINER_NAME (fast-ibc-arbitrum-attestor)
#   ATTESTOR_STATE_VOLUME (fast-ibc-arbitrum-attestor-data)
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

# CHAIN_PROFILE selects the chain-specific defaults. They are genuinely
# chain-shaped and a wrong one fails silently: the BoLD storage slot differs per
# deployment, and reading the wrong slot yields an empty mapping, so the attestor
# attests nothing and logs no error naming the slot.
#
# A named profile other than "devnet" reads its values from the committed
# reference config for that chain, rather than restating them here. The repo
# already carries those files and they are maintained; duplicating the numbers
# into this script is exactly how the two drifted apart (the script defaulted the
# slot to 0x76 while attestor/arbitrum/config.arbitrum-sepolia.json says 0x75).
#
# An explicit env var always wins, so existing invocations are unaffected.
CHAIN_PROFILE=${CHAIN_PROFILE:-devnet}
if [ "$CHAIN_PROFILE" = devnet ]; then
    # Devnet RollupCreator deploys the legacy layout, where 0x76 is correct, and
    # the devnet posts assertions every 10s so a 2s poll is proportionate.
    PROFILE_SRC_CHAIN=arbdev
    PROFILE_ASSERTIONS_MAPPING_SLOT=0x0000000000000000000000000000000000000000000000000000000000000076
    PROFILE_ASSERTION_STATUS_OFFSET=25
    PROFILE_RUNTIME_POLL_INTERVAL=2s
    PROFILE_ASSERTION_POLL_INTERVAL=2s
    PROFILE_ASSERTION_MAX_BLOCK_RANGE=2000
    PROFILE_ASSERTION_START_BLOCK=   # devnet: fall back to ROLLUP_DEPLOYMENT_BLOCK
else
    PROFILE_FILE=$REPO_ROOT/attestor/arbitrum/config.$CHAIN_PROFILE.json
    if [ ! -f "$PROFILE_FILE" ]; then
        printf '[run_arbitrum_attestor] ERROR: unknown CHAIN_PROFILE %q (no %s)\n' \
            "$CHAIN_PROFILE" "$PROFILE_FILE" >&2
        exit 1
    fi
    command -v jq >/dev/null 2>&1 || {
        printf '[run_arbitrum_attestor] ERROR: jq is required to read %s\n' "$PROFILE_FILE" >&2
        exit 1
    }
    profile_field() {
        value=$(jq -er --arg k "$1" '.[$k] // empty' "$PROFILE_FILE") || {
            printf '[run_arbitrum_attestor] ERROR: %s has no %s\n' "$PROFILE_FILE" "$1" >&2
            exit 1
        }
        printf '%s' "$value"
    }
    PROFILE_SRC_CHAIN=$(profile_field src_chain)
    PROFILE_ASSERTIONS_MAPPING_SLOT=$(profile_field assertions_mapping_slot)
    PROFILE_ASSERTION_STATUS_OFFSET=$(profile_field assertion_status_offset)
    PROFILE_RUNTIME_POLL_INTERVAL=$(profile_field runtime_poll_interval)
    PROFILE_ASSERTION_POLL_INTERVAL=$(profile_field assertion_poll_interval)
    PROFILE_ASSERTION_MAX_BLOCK_RANGE=$(profile_field assertion_max_block_range)
    PROFILE_ASSERTION_START_BLOCK=$(profile_field assertion_start_block)
    printf '[run_arbitrum_attestor] chain profile %s from %s\n' "$CHAIN_PROFILE" "$PROFILE_FILE"
fi

SRC_CHAIN=${SRC_CHAIN:-$PROFILE_SRC_CHAIN}
ASSERTIONS_MAPPING_SLOT=${ASSERTIONS_MAPPING_SLOT:-$PROFILE_ASSERTIONS_MAPPING_SLOT}
ASSERTION_STATUS_OFFSET=${ASSERTION_STATUS_OFFSET:-$PROFILE_ASSERTION_STATUS_OFFSET}
GRPC_PORT=${GRPC_PORT:-3001}
RUNTIME_POLL_INTERVAL=${RUNTIME_POLL_INTERVAL:-$PROFILE_RUNTIME_POLL_INTERVAL}
ASSERTION_POLL_INTERVAL=${ASSERTION_POLL_INTERVAL:-$PROFILE_ASSERTION_POLL_INTERVAL}
ASSERTION_MAX_BLOCK_RANGE=${ASSERTION_MAX_BLOCK_RANGE:-$PROFILE_ASSERTION_MAX_BLOCK_RANGE}
ATTESTATION_HEAD=${ATTESTATION_HEAD:-finalized}
DISABLE_DERIVED_ROOTS=${DISABLE_DERIVED_ROOTS:-false}
DERIVED_ATTESTATION_GAP_BLOCKS=${DERIVED_ATTESTATION_GAP_BLOCKS:-150}
MAX_DERIVED_ROOTS=${MAX_DERIVED_ROOTS:-1000}
# ATTESTOR_RUN_DIR takes precedence over the legacy RUN_DIR: the devnet bring-up
# scripts use RUN_DIR for their OWN artifacts, so a handoff that exports it into
# the caller's shell sends the next bring-up's package clone, downloads and
# handoff into this attestor's directory. Handoffs export ATTESTOR_RUN_DIR
# instead; RUN_DIR still works for anything that already sets it.
RUN_DIR=${ATTESTOR_RUN_DIR:-${RUN_DIR:-$REPO_ROOT/.arbitrum-attestor-run}}
STARTUP_WAIT_SECS=${STARTUP_WAIT_SECS:-600}
ATTESTOR_DOCKER_IMAGE=${ATTESTOR_DOCKER_IMAGE:-fast-ibc-arbitrum-attestor:local}
ATTESTOR_CONTAINER_NAME=${ATTESTOR_CONTAINER_NAME:-fast-ibc-arbitrum-attestor}
ATTESTOR_STATE_VOLUME=${ATTESTOR_STATE_VOLUME:-fast-ibc-arbitrum-attestor-data}
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
: "${L1_CHAIN_ID:?L1_CHAIN_ID is required}"
: "${L2_RPC_URL:?L2_RPC_URL is required (Nitro HTTP RPC)}"
: "${L2_WS_URL:?L2_WS_URL is required (Nitro WebSocket RPC)}"
: "${L2_CHAIN_ID:?L2_CHAIN_ID is required}"
: "${ROLLUP_CORE_ADDRESS:?ROLLUP_CORE_ADDRESS is required}"
: "${ROLLUP_DEPLOYMENT_BLOCK:?ROLLUP_DEPLOYMENT_BLOCK is required}"

if [ "$ROLLUP_DEPLOYMENT_BLOCK" = 0 ]; then
    log "ROLLUP_DEPLOYMENT_BLOCK is 0; using 1 because the attestor scans from a non-genesis L1 block"
    ROLLUP_DEPLOYMENT_BLOCK=1
fi

# Where the assertion scan starts. The rollup's deployment block is right for a
# devnet deployed minutes ago; on a long-lived public rollup it is a back-scan of
# millions of L1 blocks before the first attestation, so a named profile supplies
# a block near the head instead.
ASSERTION_START_BLOCK=${ASSERTION_START_BLOCK:-${PROFILE_ASSERTION_START_BLOCK:-$ROLLUP_DEPLOYMENT_BLOCK}}
[[ "$ASSERTION_START_BLOCK" =~ ^[0-9]+$ ]] ||
    fail "ASSERTION_START_BLOCK must be a decimal L1 block number, got: $ASSERTION_START_BLOCK"

[ -f "$NITRO_SEQUENCER_CONFIG" ] ||
    fail "Nitro sequencer config not found: $NITRO_SEQUENCER_CONFIG"
[[ "$GRPC_PORT" =~ ^[0-9]+$ ]] && [ "$GRPC_PORT" -gt 0 ] && [ "$GRPC_PORT" -le 65535 ] ||
    fail "GRPC_PORT must be between 1 and 65535"
[[ "$ASSERTIONS_MAPPING_SLOT" =~ ^0x[0-9a-fA-F]{64}$ ]] ||
    fail "ASSERTIONS_MAPPING_SLOT must be a 32-byte hexadecimal value"
[[ "$ASSERTION_STATUS_OFFSET" =~ ^[0-9]+$ ]] &&
    [ "$ASSERTION_STATUS_OFFSET" -lt 32 ] ||
    fail "ASSERTION_STATUS_OFFSET must be between 0 and 31"
case "$ATTESTATION_HEAD" in
    unsafe | safe | finalized) ;;
    *) fail "ATTESTATION_HEAD must be unsafe, safe, or finalized" ;;
esac
case "$DISABLE_DERIVED_ROOTS" in
    true | false) ;;
    *) fail "DISABLE_DERIVED_ROOTS must be true or false" ;;
esac
[[ "$DERIVED_ATTESTATION_GAP_BLOCKS" =~ ^[1-9][0-9]*$ ]] ||
    fail "DERIVED_ATTESTATION_GAP_BLOCKS must be greater than zero"
[[ "$MAX_DERIVED_ROOTS" =~ ^[1-9][0-9]*$ ]] ||
    fail "MAX_DERIVED_ROOTS must be greater than zero"

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
CONTAINER_NITRO_RPC_URL=${ATTESTOR_NITRO_RPC_URL:-$(container_url "$L2_RPC_URL")}
CONTAINER_NITRO_WS_URL=${ATTESTOR_NITRO_WS_URL:-$(container_url "$L2_WS_URL")}

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

CONFIG=$RUN_DIR/config.json
jq -n \
    --arg grpc_listen_address "127.0.0.1:50051" \
    --arg runtime_poll_interval "$RUNTIME_POLL_INTERVAL" \
    --arg attestation_head "$ATTESTATION_HEAD" \
    --argjson disable_derived_roots "$DISABLE_DERIVED_ROOTS" \
    --argjson derived_attestation_gap_blocks "$DERIVED_ATTESTATION_GAP_BLOCKS" \
    --argjson max_derived_roots "$MAX_DERIVED_ROOTS" \
    --arg src_chain "$SRC_CHAIN" \
    --arg l1_rpc_url "$CONTAINER_L1_RPC_URL" \
    --argjson l1_chain_id "$L1_CHAIN_ID" \
    --argjson l2_chain_id "$L2_CHAIN_ID" \
    --arg rollup_core_address "$ROLLUP_CORE_ADDRESS" \
    --arg assertions_mapping_slot "$ASSERTIONS_MAPPING_SLOT" \
    --argjson assertion_status_offset "$ASSERTION_STATUS_OFFSET" \
    --argjson assertion_start_block "$ASSERTION_START_BLOCK" \
    --arg assertion_poll_interval "$ASSERTION_POLL_INTERVAL" \
    --argjson assertion_max_block_range "$ASSERTION_MAX_BLOCK_RANGE" \
    --arg nitro_rpc_url "$CONTAINER_NITRO_RPC_URL" \
    --arg nitro_ws_url "$CONTAINER_NITRO_WS_URL" \
    '{
        grpc_listen_address: $grpc_listen_address,
        runtime_poll_interval: $runtime_poll_interval,
        attestation_head: $attestation_head,
        disable_derived_roots: $disable_derived_roots,
        derived_attestation_gap_blocks: $derived_attestation_gap_blocks,
        max_derived_roots: $max_derived_roots,
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
        attestor_state_path: "/var/lib/fast-ibc/attestor/attested-roots.json",
        nitro_rpc_url: $nitro_rpc_url,
        nitro_ws_url: $nitro_ws_url
    }' >"$CONFIG"

log "attestor config written to $CONFIG"

if docker container inspect "$ATTESTOR_CONTAINER_NAME" >/dev/null 2>&1; then
    fail "container $ATTESTOR_CONTAINER_NAME already exists; stop or rename it before starting another attestor"
fi

log "building $ATTESTOR_DOCKER_IMAGE"
docker build \
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
    --volume "$ATTESTOR_STATE_VOLUME:/var/lib/fast-ibc/attestor" \
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

The attestor assertion state is persisted in Docker volume:
  $ATTESTOR_STATE_VOLUME

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
