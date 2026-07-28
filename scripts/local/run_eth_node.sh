#!/bin/sh

set -eux

# Local Ethereum PoS L1 (geth + lighthouse, Fulu-from-genesis) via Kurtosis. Owns
# ONLY the L1 node: clones+patches the ethereum-package, starts the enclave,
# discovers the RPC/WS/beacon endpoints, and writes them to a handoff file. Contract
# deployment lives in scripts/local/deploy_eth_contracts.sh (which sources that file).
#
#   run_eth_node.sh          -> $RUN_DIR/eth.env  (ETH_RPC / ETH_WS / ETH_BEACON_API)
#   deploy_eth_contracts.sh  -> deploys E2ETestDeploy, patches relayer config
#
# scripts/local/run_optimism_node.sh reuses this to bring up the OP Stack's L1 in the
# op-devnet enclave (same package, same Fulu recipe) before layering its L2 on top.
#
# The ethereum-package is cloned+patched (not run straight from github) for two
# reasons: (1) the pinned version still passes --allow-insecure-unlock, which geth
# v1.17 removed; (2) a local clone is its own kurtosis package identity, so when the
# OP flow runs optimism-package in the same enclave it does NOT clobber the copy
# optimism-package imports Starlark helpers from (kurtosis caches by package NAME).
#
# Env (defaults in parens):
#   ENCLAVE (my-testnet)                 kurtosis enclave name
#   ETH_PIN (Fulu-validated ref below)   ethereum-package git ref; run_optimism_node
#                                        passes the exact ref optimism-package pins
#   L1_PARAMS (eth-network-params.yaml)  --args-file for the L1
#   RUN_DIR ($REPO_ROOT/.eth-devnet-run) clone + handoff dir
#   SKIP_COSMOS_RESET (unset)            when set, do NOT kill gaiad / wipe ~/.gaia
cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

ENCLAVE=${ENCLAVE:-my-testnet}
# ethereum-package the optimism-package pins (Fulu-era, lighthouse-v8 genesis) —
# used everywhere so the ETH<->Cosmos and OP devnets share one validated L1.
ETH_PIN=${ETH_PIN:-e964e305a19d56b798800e84d264f97b87952c55}
L1_PARAMS=${L1_PARAMS:-eth-network-params.yaml}
RUN_DIR=${RUN_DIR:-$REPO_ROOT/.eth-devnet-run}

log() { printf '\n[run_eth_node] %s\n' "$*"; }

kurtosis enclave rm -f "$ENCLAVE" || true
if [ -z "${SKIP_COSMOS_RESET:-}" ]; then
    killall gaiad || true
    rm -rf "$HOME/.gaia"
fi

# ------------------------------------------------- package clone + patch ---
mkdir -p "$RUN_DIR"
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
# Drop --allow-insecure-unlock (removed in geth v1.17). Portable edit (no sed -i,
# which differs on BSD/macOS vs GNU).
GETH_STAR=$ETH_DIR/src/el/geth/geth_launcher.star
if grep -q -- '--allow-insecure-unlock' "$GETH_STAR"; then
    log "patching geth_launcher.star (drop --allow-insecure-unlock; removed in geth v1.17)"
    grep -v -- '--allow-insecure-unlock' "$GETH_STAR" > "$GETH_STAR.tmp"
    mv "$GETH_STAR.tmp" "$GETH_STAR"
fi

# ------------------------------------------------------------------- run ---
log "starting L1 (enclave $ENCLAVE, package $ETH_DIR) — first run pulls images, be patient"
kurtosis run --enclave "$ENCLAVE" "$ETH_DIR" --args-file "$L1_PARAMS"

sleep 30

ETH_RPC=$(kurtosis enclave inspect "$ENCLAVE" \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "http://$1\n";
    exit;
  }
')

ETH_WS=$(kurtosis enclave inspect "$ENCLAVE" \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+ws:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "ws://$1\n";
    exit;
  }
')

ETH_BEACON_API=$(kurtosis enclave inspect "$ENCLAVE" \
| awk '
  $0 ~ /cl-1-lighthouse-geth/ {in_service=1}
  in_service && /http:/ {
      match($0, /127\.0\.0\.1:[0-9]+/)
      print "http://" substr($0, RSTART, RLENGTH)
      exit
  }
  in_service && /^[^[:space:]]/ {in_service=0}
')

echo "ETH_RPC: $ETH_RPC"
echo "ETH_WS: $ETH_WS"
echo "ETH_BEACON_API: $ETH_BEACON_API"

# Endpoint handoff for deploy_eth_contracts.sh / run_optimism_node.sh.
ENV_FILE=$RUN_DIR/eth.env
{
    printf 'export ETH_RPC=%s\n' "$ETH_RPC"
    printf 'export ETH_WS=%s\n' "$ETH_WS"
    printf 'export ETH_BEACON_API=%s\n' "$ETH_BEACON_API"
} >"$ENV_FILE"

echo "Ethereum L1 (enclave $ENCLAVE) ready. Endpoints written to $ENV_FILE"
echo "Deploy the IBC contracts + patch the relayer config with:"
echo "  scripts/local/deploy_eth_contracts.sh"
