#!/usr/bin/env bash

# Local OP Stack devnet for attestor testing, via Kurtosis, in two layers
# inside one enclave:
#   1. L1 — ethereum-package (same recipe as run_eth_node.sh: geth +
#      lighthouse, 2s slots, finality ~2min), scripts/local/op-l1-params.yaml
#   2. L2 — ethpandaops/optimism-package attached to that L1 via
#      external_l1_network_params: sequencer + an independent replica
#      op-node + batcher + proposer posting real DisputeGameFactory games
#      every 60s.
#
# The OP Stack components (op-geth, op-node, op-batcher, op-proposer,
# op-deployer, contracts) are pinned to the LATEST ethereum-optimism releases
# via the `registry:` map in scripts/local/op-network-params.yaml — the
# kurtosis package's own defaults lag the monorepo by ~a year. Bump versions
# there, not here.
#
# Ends by printing (and writing to $RUN_DIR/attestor.env) everything
# scripts/local/run_op_attestor.sh needs to run against this devnet.
#
# Optional env:
#   ENCLAVE (op-devnet)         kurtosis enclave name
#   OP_PACKAGE_REF              ethpandaops/optimism-package git ref to clone
#                               (default: a pinned main commit — 1.4.0, the
#                               latest release, imports helpers from a 2025
#                               ethereum-package that conflicts with the L1
#                               below; see ETH_PACKAGE note)
#   ETH_PACKAGE                 L1 kurtosis package. Default: a local clone of
#                               the exact version the optimism-package clone's
#                               kurtosis.yml replace directive pins, patched to
#                               drop --allow-insecure-unlock so the latest geth
#                               (>= v1.17) boots. If overriding by github URL,
#                               it MUST stay the same version optimism-package
#                               imports from: kurtosis caches dependencies by
#                               package NAME, so running a different
#                               ethereum-package version for the L1 silently
#                               replaces the copy optimism-package imports
#                               Starlark helpers from and it crashes on renamed
#                               symbols (e.g. get_client_tolerations). (A local
#                               clone is a separate package identity — safe.)
#   RUN_DIR (.op-devnet-run)    package clone + downloaded artifacts + attestor.env
#   GAME_WAIT_SECS (900)        max wait for the first proposed game

set -euo pipefail

cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

ENCLAVE=${ENCLAVE:-op-devnet}
# main @ 2025-09-19 ("chore: Update ethereum-package (#398)") — newest
# optimism-package; its ethereum-package pin is Fulu-era, which lighthouse v8
# genesis handling requires.
OP_PACKAGE_REF=${OP_PACKAGE_REF:-7bef190d7c0b9f619438ed08b17bd5e5f51e72ff}
RUN_DIR=${RUN_DIR:-$REPO_ROOT/.op-devnet-run}
GAME_WAIT_SECS=${GAME_WAIT_SECS:-900}

mkdir -p "$RUN_DIR"

log() { printf '\n[run_op_stack] %s\n' "$*"; }

# ------------------------------------------------- package clone + patch ---
# The kurtosis package is cloned locally (instead of run straight from
# github) so we can patch it: op-deployer >= v0.5 requires a non-zero
# operatorFeeVaultRecipient (Isthmus operator fee) in the intent, a field
# optimism-package 1.4.0 predates. The patch fills it with the same dev
# wallet the package generates for the sequencer fee vault.
PKG_DIR=$RUN_DIR/optimism-package
if [ ! -e "$PKG_DIR/.pinned-ref" ] || [ "$(cat "$PKG_DIR/.pinned-ref")" != "$OP_PACKAGE_REF" ]; then
    log "cloning ethpandaops/optimism-package@$OP_PACKAGE_REF into $PKG_DIR"
    rm -rf "$PKG_DIR"
    mkdir -p "$PKG_DIR"
    git -C "$PKG_DIR" init --quiet
    git -C "$PKG_DIR" fetch --quiet --depth 1 \
        https://github.com/ethpandaops/optimism-package "$OP_PACKAGE_REF"
    git -C "$PKG_DIR" checkout --quiet FETCH_HEAD
    echo "$OP_PACKAGE_REF" > "$PKG_DIR/.pinned-ref"
fi

# L1 package = exactly the ethereum-package version optimism-package pins in
# its kurtosis.yml replace directive (see ETH_PACKAGE note above) — but
# cloned locally so the geth launcher can be patched: the pinned package
# still passes --allow-insecure-unlock, which geth v1.17 removed, and
# op-l1-params.yaml pins the latest geth. The local clone is its own
# kurtosis package identity, so it does NOT clobber the name-cached copy
# optimism-package imports Starlark helpers from.
if [ -z "${ETH_PACKAGE:-}" ]; then
    ETH_PIN=$(grep -oE 'github.com/ethpandaops/ethereum-package@[0-9a-f]+' \
        "$PKG_DIR/kurtosis.yml" | head -1 | cut -d@ -f2)
    [ -n "$ETH_PIN" ] || { log "ERROR: no ethereum-package pin in $PKG_DIR/kurtosis.yml — set ETH_PACKAGE"; exit 1; }
    ETH_DIR=$RUN_DIR/ethereum-package
    if [ ! -e "$ETH_DIR/.pinned-ref" ] || [ "$(cat "$ETH_DIR/.pinned-ref")" != "$ETH_PIN" ]; then
        log "cloning ethpandaops/ethereum-package@$ETH_PIN into $ETH_DIR"
        rm -rf "$ETH_DIR"
        mkdir -p "$ETH_DIR"
        git -C "$ETH_DIR" init --quiet
        git -C "$ETH_DIR" fetch --quiet --depth 1 \
            https://github.com/ethpandaops/ethereum-package "$ETH_PIN"
        git -C "$ETH_DIR" checkout --quiet FETCH_HEAD
        echo "$ETH_PIN" > "$ETH_DIR/.pinned-ref"
    fi
    GETH_STAR=$ETH_DIR/src/el/geth/geth_launcher.star
    if grep -q -- '--allow-insecure-unlock' "$GETH_STAR"; then
        log "patching geth_launcher.star (drop --allow-insecure-unlock; removed in geth v1.17)"
        sed -i '/--allow-insecure-unlock/d' "$GETH_STAR"
    fi
    ETH_PACKAGE=$ETH_DIR
fi

DEPLOYER_STAR=$PKG_DIR/src/contracts/contract_deployer.star
if ! grep -q operatorFeeVaultRecipient "$DEPLOYER_STAR"; then
    log "patching contract_deployer.star (operatorFeeVaultRecipient for op-deployer >= v0.5)"
    python3 - "$DEPLOYER_STAR" <<'PYEOF'
import re
import sys
path = sys.argv[1]
src = open(path).read()
pat = re.compile(
    r'([ \t]+)"sequencerFeeVaultRecipient": read_chain_cmd\(\s*\n'
    r'[ \t]+"sequencerFeeVaultRecipient", chain_id\s*\n'
    r'[ \t]+\),'
)
m = pat.search(src)
if not m:
    sys.exit("anchor for operatorFeeVaultRecipient patch not found")
indent = m.group(1)
insert = (
    f'\n{indent}"operatorFeeVaultRecipient": read_chain_cmd(\n'
    f'{indent}    "sequencerFeeVaultRecipient", chain_id\n'
    f'{indent}),'
)
src = src[: m.end()] + insert + src[m.end() :]
open(path, "w").write(src)
PYEOF
fi
# op-node >= v1.18 must be told the L1 chain config when the L1 chain id is
# not a known network (it needs the L1 blob schedule); without it startup
# dies with "failed to read chain spec: open : no such file or directory".
# op-node parses that file strictly, and the ethereum-package L1 genesis
# carries the legacy terminalTotalDifficultyPassed key geth dropped in
# v1.15 — so the fund step (which has jq and writes into the
# op-deployer-configs artifact the op-nodes already mount at
# /network-configs) emits a cleaned copy as l1-chain-config.json, and the
# op-node launcher passes --rollup.l1-chain-config pointing at it.
if ! grep -q 'l1-chain-config.json' "$DEPLOYER_STAR"; then
    log "patching contract_deployer.star (emit cleaned l1-chain-config.json)"
    python3 - "$DEPLOYER_STAR" <<'PYEOF'
import sys
path = sys.argv[1]
src = open(path).read()
files_anchor = '''"/network-data": op_deployer_init.files_artifacts[0],
            "/fund-script": fund_script_artifact,
        },'''
run_anchor = """run='bash /fund-script/fund.sh "{0}"'.format(l2_chain_ids),"""
if files_anchor not in src or run_anchor not in src:
    sys.exit("anchor for l1-chain-config patch not found in contract_deployer.star")
src = src.replace(files_anchor, files_anchor.replace(
    '"/fund-script": fund_script_artifact,',
    '"/fund-script": fund_script_artifact,\n            "/l1-genesis": "el_cl_genesis_data",'), 1)
src = src.replace(run_anchor, """run='bash /fund-script/fund.sh "{0}" && jq "del(.config.terminalTotalDifficultyPassed)" /l1-genesis/genesis.json > /network-data/l1-chain-config.json'.format(l2_chain_ids),""", 1)
open(path, "w").write(src)
PYEOF
fi
# The fund script sends with --priority-gas-price 1gwei but no explicit max
# fee; on a quiet devnet L1 the basefee decays to a few wei within minutes,
# cast then computes maxFeePerGas < maxPriorityFeePerGas and every funding tx
# is rejected ("max priority fee per gas higher than max fee per gas") — the
# batcher/proposer wallets stay empty and the safe head never advances. Pin
# an explicit max fee.
FUND_SH=$(find "$PKG_DIR/static_files" -name fund.sh | head -1)
if [ -n "$FUND_SH" ] && ! grep -q -- '--gas-price 2gwei' "$FUND_SH"; then
    log "patching fund.sh (explicit --gas-price for a near-zero-basefee L1)"
    sed -i 's/--priority-gas-price 1gwei/--gas-price 2gwei --priority-gas-price 1gwei/' "$FUND_SH"
fi

OPNODE_STAR=$PKG_DIR/src/cl/op-node/launcher.star
if ! grep -q 'rollup.l1-chain-config' "$OPNODE_STAR"; then
    log "patching op-node launcher.star (--rollup.l1-chain-config for op-node >= v1.18)"
    python3 - "$OPNODE_STAR" <<'PYEOF'
import sys
path = sys.argv[1]
src = open(path).read()
cmd_anchor = '"--rpc.enable-admin",'
if cmd_anchor not in src:
    sys.exit("anchor for l1-chain-config patch not found in launcher.star")
src = src.replace(cmd_anchor, cmd_anchor + '''
        "--rollup.l1-chain-config=/network-configs/l1-chain-config.json",''', 1)
open(path, "w").write(src)
PYEOF
fi

# wait_until <timeout-secs> <description> <command...> — poll every 5s.
wait_until() {
    local timeout=$1 desc=$2; shift 2
    local deadline=$((SECONDS + timeout))
    until "$@" >/dev/null 2>&1; do
        if [ "$SECONDS" -ge "$deadline" ]; then
            log "TIMEOUT after ${timeout}s waiting for: $desc"
            return 1
        fi
        sleep 5
    done
}

log "removing any previous '$ENCLAVE' enclave"
kurtosis enclave rm -f "$ENCLAVE" 2>/dev/null || true

# ------------------------------------------------------------------- L1 ---
# Same recipe as run_eth_node.sh (ethereum-package, geth + lighthouse, fast
# slots) but in this enclave, so the OP package can reach it by service DNS
# (external_l1_network_params in op-network-params.yaml).
log "starting L1 (package: $ETH_PACKAGE, 2s slots) — first run pulls images, be patient"
kurtosis run --enclave "$ENCLAVE" "$ETH_PACKAGE" \
    --args-file scripts/local/op-l1-params.yaml

CL_HTTP="http://$(kurtosis port print "$ENCLAVE" cl-1-lighthouse-geth http | sed 's|^http://||')"
l1_finalized() {
    local epoch
    epoch=$(curl -sf "$CL_HTTP/eth/v1/beacon/states/head/finality_checkpoints" \
        | jq -r '.data.finalized.epoch') || return 1
    [ "${epoch:-0}" -gt 0 ]
}
log "waiting for the first finalized L1 epoch (~2min at 2s slots) before layering the L2 on top"
wait_until 360 "first finalized L1 epoch" l1_finalized

# ------------------------------------------------------------------- L2 ---
log "starting OP Stack L2 (package: $PKG_DIR @ $OP_PACKAGE_REF, patched)"
kurtosis run --enclave "$ENCLAVE" "$PKG_DIR" \
    --args-file scripts/local/op-network-params.yaml

# ------------------------------------------------------------ discovery ---
# L1 execution RPC (ethereum-package service inside the enclave).
L1_SVC=$(kurtosis enclave inspect "$ENCLAVE" | grep -oE 'el-1-[a-z-]+' | head -1)
L1_RPC_URL="http://$(kurtosis port print "$ENCLAVE" "$L1_SVC" rpc | sed 's|^http://||')"

# op-node services: one per participant. The attestor uses the *replica*
# participant's op-node (independent verifier); the sequencer's is first.
mapfile -t OP_NODE_SVCS < <(kurtosis enclave inspect "$ENCLAVE" | grep -oE 'op-cl-[a-zA-Z0-9-]+' | sort -u)
if [ "${#OP_NODE_SVCS[@]}" -eq 0 ]; then
    log "ERROR: no op-node (op-cl-*) services found in enclave $ENCLAVE"
    kurtosis enclave inspect "$ENCLAVE"
    exit 1
fi
REPLICA_SVC=""
for svc in "${OP_NODE_SVCS[@]}"; do
    case $svc in *replica*) REPLICA_SVC=$svc ;; esac
done
# Fall back to the last op-node service if none is named "replica".
[ -n "$REPLICA_SVC" ] || REPLICA_SVC=${OP_NODE_SVCS[${#OP_NODE_SVCS[@]}-1]}
OP_NODE_RPC_URL="http://$(kurtosis port print "$ENCLAVE" "$REPLICA_SVC" rpc | sed 's|^http://||')"

log "L1 RPC:            $L1_RPC_URL   ($L1_SVC)"
log "replica op-node:   $OP_NODE_RPC_URL   ($REPLICA_SVC)"

# Contract addresses from op-deployer's state artifact.
rm -rf "$RUN_DIR/op-deployer"
kurtosis files download "$ENCLAVE" op-deployer-configs "$RUN_DIR/op-deployer"
STATE_JSON=$(find "$RUN_DIR/op-deployer" -name 'state.json' | head -1)
if [ -z "$STATE_JSON" ]; then
    log "ERROR: state.json not found in the op-deployer-configs artifact"
    find "$RUN_DIR/op-deployer" -maxdepth 2 -type f
    exit 1
fi
DISPUTE_GAME_FACTORY=$(jq -r '.opChainDeployments[0].DisputeGameFactoryProxy' "$STATE_JSON")
OPTIMISM_PORTAL=$(jq -r '.opChainDeployments[0].OptimismPortalProxy' "$STATE_JSON")
if [ -z "$DISPUTE_GAME_FACTORY" ] || [ "$DISPUTE_GAME_FACTORY" = "null" ]; then
    log "ERROR: DisputeGameFactoryProxy not in $STATE_JSON — dumping keys:"
    jq '.opChainDeployments[0] | keys' "$STATE_JSON"
    exit 1
fi
# The proposer posts the game type from op-network-params.yaml.
RESPECTED_GAME_TYPE=$(grep -A2 'proposer_params:' scripts/local/op-network-params.yaml | grep 'game_type:' | awk '{print $2}')

log "DisputeGameFactory: $DISPUTE_GAME_FACTORY"
log "OptimismPortal:     $OPTIMISM_PORTAL"
log "game type:          $RESPECTED_GAME_TYPE (proposer_params in op-network-params.yaml)"

# ----------------------------------------------------------- readiness ---
log "waiting for the replica op-node to answer optimism_syncStatus"
sync_status() {
    curl -sf -X POST -H 'Content-Type: application/json' \
        --data '{"jsonrpc":"2.0","id":1,"method":"optimism_syncStatus","params":[]}' \
        "$OP_NODE_RPC_URL" | jq -e '.result.finalized_l2.number'
}
wait_until 180 "optimism_syncStatus on the replica" sync_status

game_count() {
    local n
    n=$(cast call "$DISPUTE_GAME_FACTORY" 'gameCount()(uint256)' --rpc-url "$L1_RPC_URL" | awk '{print $1}')
    [ "${n:-0}" -ge 1 ]
}
log "waiting for the proposer's first DisputeGameFactory game (max ${GAME_WAIT_SECS}s; interval 60s but the proposer waits for a safe/finalized L2 block first)"
if wait_until "$GAME_WAIT_SECS" "first proposed game" game_count; then
    log "games proposed: $(cast call "$DISPUTE_GAME_FACTORY" 'gameCount()(uint256)' --rpc-url "$L1_RPC_URL")"
else
    log "WARNING: no game proposed yet — check: kurtosis service logs $ENCLAVE \$(kurtosis enclave inspect $ENCLAVE | grep -oE 'op-proposer[a-zA-Z0-9-]*' | head -1)"
fi

# ------------------------------------------------------------- handoff ---
ENV_FILE=$RUN_DIR/attestor.env
cat > "$ENV_FILE" <<EOF
export L1_RPC_URL=$L1_RPC_URL
export OP_NODE_RPC_URL=$OP_NODE_RPC_URL
export DISPUTE_GAME_FACTORY=$DISPUTE_GAME_FACTORY
export OPTIMISM_PORTAL=$OPTIMISM_PORTAL
export RESPECTED_GAME_TYPE=$RESPECTED_GAME_TYPE
export SRC_CHAIN=opdev
# Low-latency profile (attestor/optimism/GUIDE.md "fastest safe profile"): attest
# unsafe-head roots ~seconds after block production; roots stay provisional
# until the finalized recheck confirms them (consumers pass
# include_provisional=true for the fresh frontier).
export ATTESTATION_HEAD=unsafe
export POLL_INTERVAL_SECONDS=2
export DERIVED_GAP_BLOCKS=5
export LOOKBACK_BLOCKS=1000
EOF

log "devnet is up. Attestor env written to $ENV_FILE"
cat <<EOF

Run the attestor against it:

  source $ENV_FILE
  ./scripts/local/run_op_attestor.sh

Inspect the devnet:

  kurtosis enclave inspect $ENCLAVE
  cast call $DISPUTE_GAME_FACTORY 'gameCount()(uint256)' --rpc-url $L1_RPC_URL

Tear down:

  kurtosis enclave rm -f $ENCLAVE

EOF
