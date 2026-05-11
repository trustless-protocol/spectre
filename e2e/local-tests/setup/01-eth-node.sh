#!/usr/bin/env bash

set -euxo pipefail

kurtosis enclave rm -f my-testnet || true
killall gaiad || true
rm -rf $HOME/.gaia

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

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

ETH_WS=$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" \
| perl -ne '
  if (/el-1-geth-lighthouse/) { $in=1 }
  elsif ($in && /^\S/) { $in=0 }
  elsif ($in && /^\s+ws:.*->\s+(127\.0\.0\.1:\d+)/) {
    print "ws://$1\n";
    exit;
  }
')

ETH_BEACON_API=$(kurtosis enclave inspect "$KURTOSIS_ENCLAVE" \
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
    --optimizer-runs 200 \
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

ENV_FILE="$REPO_ROOT/relayer/.env"
upsert_env_var "$ENV_FILE" ETH_RPC_URL "$ETH_RPC"
upsert_env_var "$ENV_FILE" ETH_WS_URL "$ETH_WS"
upsert_env_var "$ENV_FILE" ETH_BEACON_API_URL "$ETH_BEACON_API"
upsert_env_var "$ENV_FILE" ERC20_ADDRESS "$ERC20_ADDRESS"
upsert_env_var "$ENV_FILE" ICS20_ADDRESS "$ICS20_ADDRESS"
upsert_env_var "$ENV_FILE" ICS26_ADDRESS "$ICS26_ADDRESS"
upsert_env_var "$ENV_FILE" ETH_PRIVATE_KEY "bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31"
upsert_env_var "$ENV_FILE" PROVER_BIN_DIR "./bin"


# Configure runtime relayer config
cd "$REPO_ROOT/relayer"
ensure_relayer_config config.json config.example.json
update_relayer_deploy_config \
  config.json \
  "$ETH_RPC" \
  "$ETH_WS" \
  "$ICS26_ADDRESS" \
  "$VERIFIER_ADDRESS" \
  "$MEMBERSHIP_ADDRESS" \
  "$UPDATE_CLIENT_ADDRESS" \
  "$MISBEHAVIOUR_ADDRESS" \
  "$ETH_BEACON_API"

echo "Waiting for execution finality (polling until finalized block > 1)..."
FINALIZED_BLOCK=0
for i in $(seq 1 60); do
  FINALIZED_BLOCK_HEX=$(curl -s \
    -H 'content-type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["finalized",false],"id":1}' \
    "$ETH_RPC" \
    | jq -r '.result.number // empty' 2>/dev/null || echo "")

  if [[ "$FINALIZED_BLOCK_HEX" =~ ^0x[0-9a-fA-F]+$ ]]; then
    FINALIZED_BLOCK=$((16#${FINALIZED_BLOCK_HEX#0x}))
  fi

  if [ "$FINALIZED_BLOCK" -gt 1 ]; then
    echo "Execution finalized at block $FINALIZED_BLOCK"
    break
  fi
  sleep 5
done
if [ "$FINALIZED_BLOCK" -le 1 ]; then
  echo "Error: finalized block did not become > 1 within 5 minutes" >&2
  exit 1
fi
