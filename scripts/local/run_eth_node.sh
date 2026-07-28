#!/bin/sh

set -euxo pipefail

# Local Ethereum PoS devnet (geth + lighthouse) via Kurtosis. This script owns ONLY
# the node: it starts the enclave, discovers the RPC/WS/beacon endpoints, and writes
# them to a handoff file. Contract deployment lives in the companion
# scripts/local/deploy_eth_contracts.sh (which sources that handoff).
#
#   run_eth_node.sh          -> .eth-devnet-run/eth.env  (endpoints)
#   deploy_eth_contracts.sh  -> deploys E2ETestDeploy, patches relayer config
#
# All internal paths (eth-network-params.yaml, relayer/) are repo-root relative — cd up
# so this script works regardless of where it's invoked from.
cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

RUN_DIR=${RUN_DIR:-$REPO_ROOT/.eth-devnet-run}

kurtosis enclave rm -f my-testnet || true
killall gaiad || true
rm -rf "$HOME/.gaia"

# Run eth chain
kurtosis run --enclave my-testnet github.com/ethpandaops/ethereum-package@6.1.0 --args-file eth-network-params.yaml

sleep 30

ETH_RPC=$(kurtosis enclave inspect my-testnet \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "http://$1\n";
    exit;
  }
')

ETH_WS=$(kurtosis enclave inspect my-testnet \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+ws:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "ws://$1\n";
    exit;
  }
')

ETH_BEACON_API=$(kurtosis enclave inspect my-testnet \
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

# Endpoint handoff for deploy_eth_contracts.sh.
mkdir -p "$RUN_DIR"
ENV_FILE=$RUN_DIR/eth.env
{
    printf 'export ETH_RPC=%s\n' "$ETH_RPC"
    printf 'export ETH_WS=%s\n' "$ETH_WS"
    printf 'export ETH_BEACON_API=%s\n' "$ETH_BEACON_API"
} >"$ENV_FILE"

echo "Ethereum node ready. Endpoints written to $ENV_FILE"
echo "Deploy the IBC contracts + patch the relayer config with:"
echo "  scripts/local/deploy_eth_contracts.sh"
