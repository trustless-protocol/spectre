#!/usr/bin/env bash

set -euxo pipefail

kurtosis enclave rm -f my-testnet || true
killall gaiad || true
rm -rf $HOME/.gaia

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"

# Run eth chain
kurtosis run --enclave my-testnet github.com/ethpandaops/ethereum-package@6.1.0 --args-file "$REPO_ROOT/eth-network-params.yaml"

sleep 30

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"

ETH_RPC=$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+rpc:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "http://$1\n";
    exit;
  }
')

# Auto-detect beacon port for finality check
ETH_BEACON_PORT=$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" \
| awk '
  $0 ~ /cl-1-lighthouse-geth/ {in_service=1}
  in_service && /http:/ {
      match($0, /127\.0\.0\.1:[0-9]+/)
      print substr($0, RSTART+10, RLENGTH-10)
      exit
  }
  in_service && /^[^[:space:]]/ {in_service=0}
')

ETH_WS=$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" \
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

# Deploy ETH contracts
export E2E_FAUCET_ADDRESS=0x8943545177806ED17B9F23F0a21ee5948eCaa776
RESULT=$(cd "$REPO_ROOT" && forge script scripts/E2ETestDeploy.s.sol:E2ETestDeploy \
    --rpc-url $ETH_RPC \
    --broadcast \
    --ffi \
    --sender 0x8943545177806ED17B9F23F0a21ee5948eCaa776 --private-key bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31 \
    2>/dev/null
)

ERC20_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.erc20')

echo "ERC20_ADDRESS: $ERC20_ADDRESS"

ICS20_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.ics20Transfer')

echo "ICS20_ADDRESS: $ICS20_ADDRESS"

ICS26_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.ics26Router')

echo "ICS26_ADDRESS: $ICS26_ADDRESS"

VERIFIER_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.wrapperVerifier')

echo "VERIFIER_ADDRESS: $VERIFIER_ADDRESS"

MEMBERSHIP_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.membership')

echo "MEMBERSHIP_ADDRESS: $MEMBERSHIP_ADDRESS"

UPDATE_CLIENT_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.updateClient')
echo "UPDATE_CLIENT_ADDRESS: $UPDATE_CLIENT_ADDRESS"

MISBEHAVIOUR_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.misbehaviour')
echo "MISBEHAVIOUR_ADDRESS: $MISBEHAVIOUR_ADDRESS"


# Configure relayer config
cd "$REPO_ROOT/relayer"
jq \
  --arg ETH_RPC "$ETH_RPC" \
  --arg ETH_WS "$ETH_WS" \
  --arg ICS26 "$ICS26_ADDRESS" \
  --arg WRAP "$VERIFIER_ADDRESS" \
  --arg MEMB "$MEMBERSHIP_ADDRESS" \
  --arg UPCL "$UPDATE_CLIENT_ADDRESS" \
  --arg MIS "$MISBEHAVIOUR_ADDRESS" \
  --arg ETH_BEACON "$ETH_BEACON_API" '
    (.. | objects | select(has("eth_rpc_url")) | .eth_rpc_url) = $ETH_RPC
  | (.. | objects | select(has("eth_ws_url")) | .eth_ws_url) = $ETH_WS
  | (.. | objects | select(has("ics26_address")) | .ics26_address) = $ICS26
  | (.. | objects | select(has("wrapper_verifier")) | .wrapper_verifier) = $WRAP
  | (.. | objects | select(has("membership")) | .membership) = $MEMB
  | (.. | objects | select(has("update_client")) | .update_client) = $UPCL
  | (.. | objects | select(has("misbehaviour")) | .misbehaviour) = $MIS
  | (.. | objects | select(has("eth_beacon_api_url")) | .eth_beacon_api_url) = $ETH_BEACON
  ' config.example.json > config.tmp && mv config.tmp config.example.json
