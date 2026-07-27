#!/usr/bin/env bash

# Local Arbitrum Nitro devnet backed by a local proof-of-stake Ethereum L1.
#
# This wraps Offchain Labs' official nitro-testnode instead of duplicating its
# RollupCore deployment and Nitro configuration logic. The testnode owns its
# Docker Compose stack and persistent volumes. This wrapper:
#
#   1. checks out one nitro-testnode revision under RUN_DIR;
#   2. initializes its local geth + Prysm L1 and simple Nitro node;
#   3. waits for L1 finality and the L2 JSON-RPC endpoint;
#   4. discovers the deployed RollupCore from Nitro's generated chain config;
#   5. writes the attestor handoff to RUN_DIR/attestor.env.
#
# This mirrors the OP local workflow:
#   run_arbitrum_stack.sh    owns the L1 + L2 devnet
#   run_arbitrum_attestor.sh owns an independent verifier replica + attestor
#
# Re-running without --reset preserves the existing chain. Reinitialization is
# deliberately explicit because nitro-testnode --init-force deletes its Docker
# volumes.

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

RUN_DIR=${RUN_DIR:-$REPO_ROOT/.arbitrum-devnet-run}
TESTNODE_DIR=${TESTNODE_DIR:-$RUN_DIR/nitro-testnode}
NITRO_TESTNODE_REPOSITORY=${NITRO_TESTNODE_REPOSITORY:-https://github.com/OffchainLabs/nitro-testnode.git}
NITRO_TESTNODE_REF=${NITRO_TESTNODE_REF:-release}
STARTUP_WAIT_SECS=${STARTUP_WAIT_SECS:-600}
FINALITY_WAIT_SECS=${FINALITY_WAIT_SECS:-900}
ASSERTION_WAIT_SECS=${ASSERTION_WAIT_SECS:-900}
ASSERTIONS_MAPPING_SLOT=${ASSERTIONS_MAPPING_SLOT:-0x0000000000000000000000000000000000000000000000000000000000000076}
ASSERTION_STATUS_OFFSET=${ASSERTION_STATUS_OFFSET:-25}

# Upstream test-node.bash checks this exact Compose project label when deciding
# whether initialization is needed. Keep it explicit and refuse to take over
# pre-existing volumes unless the caller passes --reset.
COMPOSE_PROJECT_NAME=nitro-testnode
export COMPOSE_PROJECT_NAME

ACTION=up
RESET=false

usage() {
    cat <<'EOF'
Usage: scripts/local/run_arbitrum_stack.sh [--reset|--stop|--status]

Starts an official Nitro testnode with a local proof-of-stake Ethereum L1.

Options:
  --reset   Delete and recreate this testnode's Docker volumes.
  --stop    Stop containers while preserving their volumes.
  --status  Show the testnode's Docker Compose status.
  -h        Show this help.

Environment:
  RUN_DIR                    Runtime checkout/output directory
                             (default: .arbitrum-devnet-run).
  TESTNODE_DIR               nitro-testnode checkout directory.
  NITRO_TESTNODE_REPOSITORY  Source repository.
  NITRO_TESTNODE_REF         Branch, tag, or commit fetched on first setup
                             (default: release).
  STARTUP_WAIT_SECS          RPC startup timeout (default: 600).
  FINALITY_WAIT_SECS         finalized-L1 timeout (default: 900).
  ASSERTION_WAIT_SECS        first finalized assertion timeout (default: 900).
  ASSERTIONS_MAPPING_SLOT    BoLD _assertions storage slot (default: 0x76).
  ASSERTION_STATUS_OFFSET    AssertionNode status byte offset (default: 25).

The first run needs network access to clone the official testnode and pull its
Docker images. Set NITRO_TESTNODE_REF to a full commit SHA in CI for a
reproducible deployment.
EOF
}

fail() {
    printf '[run_arbitrum_stack] ERROR: %s\n' "$*" >&2
    exit 1
}

log() {
    printf '\n[run_arbitrum_stack] %s\n' "$*"
}

while [ "$#" -gt 0 ]; do
    case "$1" in
        --reset)
            [ "$ACTION" = up ] || fail "--reset cannot be combined with --stop or --status"
            RESET=true
            ;;
        --stop)
            [ "$ACTION" = up ] || fail "--stop cannot be combined with another action"
            [ "$RESET" = false ] || fail "--reset cannot be combined with --stop"
            ACTION=stop
            ;;
        --status)
            [ "$ACTION" = up ] || fail "--status cannot be combined with another action"
            [ "$RESET" = false ] || fail "--reset cannot be combined with --status"
            ACTION=status
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            usage >&2
            fail "unknown argument: $1"
            ;;
    esac
    shift
done

require_command() {
    command -v "$1" >/dev/null 2>&1 || fail "$1 is required"
}

require_command docker
require_command git
require_command curl
require_command jq
require_command cast
docker compose version >/dev/null 2>&1 || fail "Docker Compose v2 is required"

COMPOSE_OWNER_MARKER=$RUN_DIR/compose-project
if [ "$ACTION" != up ]; then
    [ -d "$TESTNODE_DIR/.git" ] ||
        fail "nitro-testnode checkout not found at $TESTNODE_DIR"
    [ -f "$COMPOSE_OWNER_MARKER" ] ||
        fail "this run directory has not initialized the $COMPOSE_PROJECT_NAME Compose project"
    case "$ACTION" in
        stop)
            log "stopping the local Arbitrum stack and preserving its volumes"
            (
                cd "$TESTNODE_DIR"
                docker compose down
            )
            ;;
        status)
            (
                cd "$TESTNODE_DIR"
                docker compose ps
            )
            ;;
    esac
    exit 0
fi

mkdir -p "$RUN_DIR"

checkout_testnode() {
    if [ ! -d "$TESTNODE_DIR/.git" ]; then
        [ ! -e "$TESTNODE_DIR" ] ||
            fail "$TESTNODE_DIR exists but is not a git checkout"
        log "cloning official nitro-testnode ref $NITRO_TESTNODE_REF"
        mkdir -p "$TESTNODE_DIR"
        git -C "$TESTNODE_DIR" init --quiet
        git -C "$TESTNODE_DIR" remote add origin "$NITRO_TESTNODE_REPOSITORY"
        git -C "$TESTNODE_DIR" fetch --depth 1 origin "$NITRO_TESTNODE_REF"
        git -C "$TESTNODE_DIR" checkout --detach FETCH_HEAD
        git -C "$TESTNODE_DIR" submodule update --init --recursive --depth 1
        git -C "$TESTNODE_DIR" rev-parse HEAD >"$RUN_DIR/nitro-testnode.commit"
        return
    fi

    [ -z "$(git -C "$TESTNODE_DIR" status --porcelain)" ] ||
        fail "nitro-testnode checkout has local changes: $TESTNODE_DIR"
    if [ ! -f "$RUN_DIR/nitro-testnode.commit" ]; then
        git -C "$TESTNODE_DIR" rev-parse HEAD >"$RUN_DIR/nitro-testnode.commit"
    fi
    log "reusing nitro-testnode commit $(cat "$RUN_DIR/nitro-testnode.commit")"
}

checkout_testnode

has_testnode_volumes() {
    [ -n "$(docker volume ls \
        --filter "label=com.docker.compose.project=$COMPOSE_PROJECT_NAME" \
        --quiet)" ]
}

if [ ! -f "$COMPOSE_OWNER_MARKER" ]; then
    if [ "$RESET" = false ] && has_testnode_volumes; then
        fail "Docker volumes for project $COMPOSE_PROJECT_NAME already exist; use --reset only if they may be replaced"
    fi
    printf '%s\n' "$COMPOSE_PROJECT_NAME" >"$COMPOSE_OWNER_MARKER"
elif [ "$(cat "$COMPOSE_OWNER_MARKER")" != "$COMPOSE_PROJECT_NAME" ]; then
    fail "unexpected Compose project recorded in $COMPOSE_OWNER_MARKER"
fi

INITIALIZED_MARKER=$RUN_DIR/initialized
INITIALIZING=false
if [ "$RESET" = true ] || [ ! -f "$INITIALIZED_MARKER" ]; then
    INITIALIZING=true
    rm -f "$INITIALIZED_MARKER"
    if [ "$RESET" = true ]; then
        log "reset requested; nitro-testnode will recreate its Docker volumes"
    else
        log "initializing the local PoS L1 and Nitro rollup"
    fi
    (
        cd "$TESTNODE_DIR"
        ./test-node.bash \
            --init-force \
            --pos \
            --simple \
            --no-tokenbridge \
            --detach
    )
else
    log "starting the existing local PoS L1 and Nitro rollup"
    (
        cd "$TESTNODE_DIR"
        ./test-node.bash \
            --pos \
            --simple \
            --no-tokenbridge \
            --detach
    )
fi

rpc_result() {
    local url=$1
    local method=$2
    local params=${3:-'[]'}

    curl --fail --silent --show-error \
        --header 'Content-Type: application/json' \
        --data "{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"$method\",\"params\":$params}" \
        "$url" |
        jq -er 'if .error then error(.error.message) else .result end'
}

wait_until() {
    local timeout=$1
    local description=$2
    shift 2
    local deadline=$((SECONDS + timeout))

    until "$@" >/dev/null 2>&1; do
        if [ "$SECONDS" -ge "$deadline" ]; then
            log "TIMEOUT after ${timeout}s waiting for: $description"
            return 1
        fi
        sleep 5
    done
}

L1_RPC_URL=http://127.0.0.1:8545
L1_WS_URL=ws://127.0.0.1:8546
L1_BEACON_API_URL=http://127.0.0.1:3500
L2_RPC_URL=http://127.0.0.1:8547
L2_WS_URL=ws://127.0.0.1:8548
NITRO_FEED_URL=ws://127.0.0.1:9642

l1_ready() {
    rpc_result "$L1_RPC_URL" eth_chainId
}

l2_ready() {
    rpc_result "$L2_RPC_URL" eth_chainId
}

beacon_ready() {
    curl --fail --silent "$L1_BEACON_API_URL/eth/v1/node/version" |
        jq -e '.data.version | length > 0'
}

l1_finalized() {
    local block
    block=$(rpc_result "$L1_RPC_URL" eth_getBlockByNumber '["finalized",false]') ||
        return 1
    [ "$(printf '%s' "$block" | jq -r '.number // "0x0"')" != 0x0 ]
}

log "waiting for local L1 execution RPC"
wait_until "$STARTUP_WAIT_SECS" "L1 execution RPC" l1_ready

log "waiting for local L1 beacon API"
wait_until "$STARTUP_WAIT_SECS" "L1 beacon API" beacon_ready

log "waiting for Nitro L2 RPC"
wait_until "$STARTUP_WAIT_SECS" "Nitro L2 RPC" l2_ready

log "waiting for a non-genesis finalized L1 block"
wait_until "$FINALITY_WAIT_SECS" "L1 finality" l1_finalized

hex_to_decimal() {
    local value=$1
    [[ "$value" =~ ^0x[0-9a-fA-F]+$ ]] ||
        fail "invalid hexadecimal integer returned by JSON-RPC: $value"
    printf '%d\n' "$((value))"
}

L1_CHAIN_ID=$(hex_to_decimal "$(rpc_result "$L1_RPC_URL" eth_chainId)")
L2_CHAIN_ID=$(hex_to_decimal "$(rpc_result "$L2_RPC_URL" eth_chainId)")
[ "$L1_CHAIN_ID" -eq 1337 ] ||
    fail "nitro-testnode PoS L1 chain ID is $L1_CHAIN_ID, expected 1337"

SEQUENCER_CONFIG=$RUN_DIR/sequencer_config.json
(
    cd "$TESTNODE_DIR"
    docker compose exec -T sequencer cat /config/sequencer_config.json
) >"$SEQUENCER_CONFIG"

discover_rollup_from_json() {
    jq -r '
        [
            .. | objects |
            (.rollup? // empty) |
            if type == "string" then .
            elif type == "object" then (.rollup? // empty)
            else empty
            end |
            select(type == "string" and test("^0x[0-9a-fA-F]{40}$"))
        ][0] // empty
    ' "$1"
}

ROLLUP_CORE_ADDRESS=$(discover_rollup_from_json "$SEQUENCER_CONFIG")
if [ -z "$ROLLUP_CORE_ADDRESS" ]; then
    CHAIN_INFO_JSON=$(jq -r '
        first(
            .. | objects |
            .["info-json"]? |
            select(type == "string")
        ) // empty
    ' "$SEQUENCER_CONFIG")
    if [ -n "$CHAIN_INFO_JSON" ]; then
        CHAIN_INFO_FILE=$RUN_DIR/chain-info.json
        printf '%s\n' "$CHAIN_INFO_JSON" >"$CHAIN_INFO_FILE"
        ROLLUP_CORE_ADDRESS=$(discover_rollup_from_json "$CHAIN_INFO_FILE")
    fi
fi

if ! [[ "$ROLLUP_CORE_ADDRESS" =~ ^0x[0-9a-fA-F]{40}$ ]]; then
    fail "could not discover RollupCore from $SEQUENCER_CONFIG"
fi

find_contract_deployment_block() {
    local address=$1
    local latest_hex
    local low=0
    local high
    local middle
    local block_tag
    local code

    latest_hex=$(rpc_result "$L1_RPC_URL" eth_blockNumber)
    high=$(hex_to_decimal "$latest_hex")
    while [ "$low" -lt "$high" ]; do
        middle=$(((low + high) / 2))
        printf -v block_tag '0x%x' "$middle"
        code=$(rpc_result "$L1_RPC_URL" eth_getCode "[\"$address\",\"$block_tag\"]")
        if [ "$code" = 0x ]; then
            low=$((middle + 1))
        else
            high=$middle
        fi
    done
    printf '%d\n' "$low"
}

ROLLUP_DEPLOYMENT_BLOCK=$(find_contract_deployment_block "$ROLLUP_CORE_ADDRESS")

ASSERTION_CREATED_SIGNATURE='AssertionCreated(bytes32,bytes32,((bytes32,bytes32,(bytes32,uint256,address,uint64,uint64)),((bytes32[2],uint64[2]),uint8,bytes32),((bytes32[2],uint64[2]),uint8,bytes32)),bytes32,uint256,bytes32,uint256,address,uint64)'
ASSERTION_CREATED_TOPIC=$(cast keccak "$ASSERTION_CREATED_SIGNATURE")

finalized_assertion_logs() {
    local finalized_block
    local params

    finalized_block=$(rpc_result "$L1_RPC_URL" eth_getBlockByNumber '["finalized",false]') ||
        return 1
    finalized_block=$(printf '%s' "$finalized_block" | jq -er '.number')
    params=$(jq -cn \
        --arg address "$ROLLUP_CORE_ADDRESS" \
        --arg from_block "$(printf '0x%x' "$ROLLUP_DEPLOYMENT_BLOCK")" \
        --arg to_block "$finalized_block" \
        --arg topic "$ASSERTION_CREATED_TOPIC" \
        '[{
            address: $address,
            fromBlock: $from_block,
            toBlock: $to_block,
            topics: [$topic]
        }]')
    rpc_result "$L1_RPC_URL" eth_getLogs "$params"
}

has_finalized_assertion() {
    finalized_assertion_logs | jq -e 'length > 0'
}

validate_assertion_storage_layout() {
    local assertion_hash=$1
    local storage_key
    local storage_word
    local byte_from_left
    local status_hex
    local status

    [[ "$ASSERTIONS_MAPPING_SLOT" =~ ^0x[0-9a-fA-F]{64}$ ]] ||
        fail "ASSERTIONS_MAPPING_SLOT must be a 32-byte hexadecimal value"
    [[ "$ASSERTION_STATUS_OFFSET" =~ ^[0-9]+$ ]] &&
        [ "$ASSERTION_STATUS_OFFSET" -lt 32 ] ||
        fail "ASSERTION_STATUS_OFFSET must be between 0 and 31"

    storage_key=$(cast index bytes32 "$assertion_hash" "$ASSERTIONS_MAPPING_SLOT")
    storage_word=$(rpc_result "$L1_RPC_URL" eth_getStorageAt \
        "[\"$ROLLUP_CORE_ADDRESS\",\"$storage_key\",\"finalized\"]")
    [[ "$storage_word" =~ ^0x[0-9a-fA-F]{64}$ ]] ||
        fail "RollupCore returned a malformed _assertions storage word"

    # Solidity storage offsets count from the least-significant byte, whereas
    # the RPC word is rendered most-significant byte first.
    byte_from_left=$((31 - ASSERTION_STATUS_OFFSET))
    status_hex=${storage_word:$((2 + byte_from_left * 2)):2}
    status=$((16#$status_hex))
    [ "$status" -eq 1 ] || [ "$status" -eq 2 ] ||
        fail "configured BoLD storage layout did not resolve a pending or confirmed assertion (status=$status)"
}

log "waiting for the first AssertionCreated event covered by finalized L1 (max ${ASSERTION_WAIT_SECS}s)"
if wait_until "$ASSERTION_WAIT_SECS" "first finalized BoLD assertion" has_finalized_assertion; then
    ASSERTION_HASH=$(finalized_assertion_logs | jq -er '.[-1].topics[1]')
    validate_assertion_storage_layout "$ASSERTION_HASH"
    log "validated BoLD _assertions layout with finalized assertion $ASSERTION_HASH"
else
    log "WARNING: no finalized assertion yet; keeping the reviewed v3.1.0 storage-layout defaults"
fi

if [ "$INITIALIZING" = true ]; then
    touch "$INITIALIZED_MARKER"
fi

NITRO_IMAGE=${NITRO_NODE_VERSION:-}
if [ -z "$NITRO_IMAGE" ]; then
    NITRO_IMAGE=$(grep -oE 'offchainlabs/nitro-node:[[:alnum:]._-]+' \
        "$TESTNODE_DIR/test-node.bash" | head -1 || true)
fi
[ -n "$NITRO_IMAGE" ] ||
    fail "could not determine the official Nitro image used by nitro-testnode"

ENV_FILE=$RUN_DIR/attestor.env
{
    printf 'export L1_RPC_URL=%q\n' "$L1_RPC_URL"
    printf 'export L1_WS_URL=%q\n' "$L1_WS_URL"
    printf 'export L1_BEACON_URL=%q\n' "$L1_BEACON_API_URL"
    printf 'export L1_BEACON_API_URL=%q\n' "$L1_BEACON_API_URL"
    printf 'export L1_CHAIN_ID=%q\n' "$L1_CHAIN_ID"
    printf 'export L2_RPC_URL=%q\n' "$L2_RPC_URL"
    printf 'export L2_WS_URL=%q\n' "$L2_WS_URL"
    printf 'export L2_CHAIN_ID=%q\n' "$L2_CHAIN_ID"
    printf 'export NITRO_FEED_URL=%q\n' "$NITRO_FEED_URL"
    printf 'export ATTESTOR_L1_RPC_URL=%q\n' http://geth:8545
    printf 'export ATTESTOR_L1_BEACON_URL=%q\n' http://prysm_beacon_chain:3500
    printf 'export ATTESTOR_NITRO_FEED_URL=%q\n' ws://sequencer:9642
    printf 'export ATTESTOR_DOCKER_NETWORK=%q\n' "${COMPOSE_PROJECT_NAME}_default"
    printf 'export ROLLUP_CORE_ADDRESS=%q\n' "$ROLLUP_CORE_ADDRESS"
    printf 'export ROLLUP_DEPLOYMENT_BLOCK=%q\n' "$ROLLUP_DEPLOYMENT_BLOCK"
    printf 'export ROLLUP_PROTOCOL=%q\n' bold-v2
    printf 'export ASSERTIONS_MAPPING_SLOT=%q\n' "$ASSERTIONS_MAPPING_SLOT"
    printf 'export ASSERTION_STATUS_OFFSET=%q\n' "$ASSERTION_STATUS_OFFSET"
    printf 'export SRC_CHAIN=%q\n' arbdev
    printf 'export ATTESTOR_NITRO_IMAGE=%q\n' "$NITRO_IMAGE"
    printf 'export NITRO_TESTNODE_DIR=%q\n' "$TESTNODE_DIR"
    printf 'export NITRO_SEQUENCER_CONFIG=%q\n' "$SEQUENCER_CONFIG"
} >"$ENV_FILE"

log "local Arbitrum stack is ready"
cat <<EOF
L1 execution RPC:      $L1_RPC_URL
L1 WebSocket:          $L1_WS_URL
L1 beacon API:         $L1_BEACON_API_URL
L1 chain ID:           $L1_CHAIN_ID
L2 execution RPC:      $L2_RPC_URL
L2 WebSocket:          $L2_WS_URL
L2 chain ID:           $L2_CHAIN_ID
RollupCore:            $ROLLUP_CORE_ADDRESS
Rollup deployment L1:  $ROLLUP_DEPLOYMENT_BLOCK
BoLD assertions slot:  $ASSERTIONS_MAPPING_SLOT
Attestor environment:  $ENV_FILE

Run the independent Nitro replica and Arbitrum attestor:
  source "$ENV_FILE"
  scripts/local/run_arbitrum_attestor.sh

Stop the stack without deleting data:
  scripts/local/run_arbitrum_stack.sh --stop

Recreate the chain and contracts:
  scripts/local/run_arbitrum_stack.sh --reset
EOF
