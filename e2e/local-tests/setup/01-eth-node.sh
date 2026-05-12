#!/usr/bin/env bash

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
source "$REPO_ROOT/e2e/local-tests/lib/common.sh"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"

log_header "Cleaning previous state"
kurtosis enclave rm -f my-testnet 2>/dev/null || true
killall gaiad 2>/dev/null || true
rm -rf "$HOME/.gaia"

log_header "Starting Ethereum Kurtosis"
kurtosis run --enclave my-testnet github.com/ethpandaops/ethereum-package@6.1.0 --args-file "$REPO_ROOT/eth-network-params.yaml"

sleep 30

KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
discover_kurtosis_endpoints

log_header "Ethereum Endpoints"
log_kv "ETH_RPC" "$ETH_RPC_URL"
log_kv "ETH_WS" "$ETH_WS_URL"
log_kv "ETH_BEACON_API" "$ETH_BEACON_API_URL"

log "Deploying ETH contracts..."
export E2E_FAUCET_ADDRESS=0x8943545177806ED17B9F23F0a21ee5948eCaa776
RESULT=$(cd "$REPO_ROOT" && forge script scripts/E2ETestDeploy.s.sol:E2ETestDeploy \
    --rpc-url "$ETH_RPC_URL" \
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

ICS20_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.ics20Transfer')

ICS26_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.ics26Router')

VERIFIER_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.wrapperVerifier')

MEMBERSHIP_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.membership')

UPDATE_CLIENT_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.updateClient')

MISBEHAVIOUR_ADDRESS=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g' \
  | jq -r '.misbehaviour')

log_header "Deployed Contracts"
log_kv "ERC20" "$ERC20_ADDRESS"
log_kv "ICS20Transfer" "$ICS20_ADDRESS"
log_kv "ICS26Router" "$ICS26_ADDRESS"
log_kv "Verifier" "$VERIFIER_ADDRESS"
log_kv "Membership" "$MEMBERSHIP_ADDRESS"
log_kv "UpdateClient" "$UPDATE_CLIENT_ADDRESS"
log_kv "Misbehaviour" "$MISBEHAVIOUR_ADDRESS"

ENV_FILE="$REPO_ROOT/relayer/.env"
upsert_env_var "$ENV_FILE" ETH_RPC_URL "$ETH_RPC_URL"
upsert_env_var "$ENV_FILE" ETH_WS_URL "$ETH_WS_URL"
upsert_env_var "$ENV_FILE" ETH_BEACON_API_URL "$ETH_BEACON_API_URL"
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
  "$ETH_RPC_URL" \
  "$ETH_WS_URL" \
  "$ICS26_ADDRESS" \
  "$VERIFIER_ADDRESS" \
  "$MEMBERSHIP_ADDRESS" \
  "$UPDATE_CLIENT_ADDRESS" \
  "$MISBEHAVIOUR_ADDRESS" \
  "$ETH_BEACON_API_URL"

log_header "Waiting for Execution Finality"

FINALIZED_BLOCK=0
for i in $(seq 1 60); do
  FINALIZED_BLOCK_HEX=$(curl -s \
    -H 'content-type: application/json' \
    --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["finalized",false],"id":1}' \
    "$ETH_RPC_URL" \
    | jq -r '.result.number // empty' 2>/dev/null || echo "")

  if [[ "$FINALIZED_BLOCK_HEX" =~ ^0x[0-9a-fA-F]+$ ]]; then
    FINALIZED_BLOCK=$((16#${FINALIZED_BLOCK_HEX#0x}))
  fi

  if [ "$FINALIZED_BLOCK" -gt 1 ]; then
    log_ok "Finalized at block $FINALIZED_BLOCK"
    break
  fi
  printf "  [%2d/60] polling...\r" "$i"
  sleep 5
done
echo ""
if [ "$FINALIZED_BLOCK" -le 1 ]; then
  log_err "finalized block did not become > 1 within 5 minutes"
  exit 1
fi
