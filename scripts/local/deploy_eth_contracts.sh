#!/bin/sh

set -euxo pipefail

# Deploys the IBC v2 contracts (E2ETestDeploy) to a running local Ethereum devnet
# and patches the relayer config with the resulting addresses + endpoints. The node is
# owned by scripts/local/run_eth_node.sh; this script is the deploy half that was split
# out of it.
#
# Endpoint source: when ETH_RPC is unset, the .eth-devnet-run/eth.env handoff written by
# run_eth_node.sh is sourced automatically. Override any of ETH_RPC / ETH_WS /
# ETH_BEACON_API in the environment to target a different node.
#
# Env:
#   ETH_RPC / ETH_WS / ETH_BEACON_API   endpoints; sourced from the run_eth_node.sh
#                                       handoff when ETH_RPC is unset
#   RELAYER_CONFIG                      relayer config to patch
#                                       (default: relayer/config.json)
#   ETH_DEPLOYER_ADDRESS / _PRIVATE_KEY deployer (defaults to the devnet-only key
#                                       relayer/.env ships). MUST match the relayer's
#                                       ETH_PRIVATE_KEY — see the note below.
#   E2E_FAUCET_ADDRESS                  test ERC20 recipient (defaults to the deployer)
#
# All internal paths (scripts/E2ETestDeploy.s.sol, relayer/) are repo-root relative.
cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

RUN_DIR=${RUN_DIR:-$REPO_ROOT/.eth-devnet-run}
ENV_FILE=${ENV_FILE:-$RUN_DIR/eth.env}
# Patch the operator's own config, not the committed example. The default used
# to be config.example.json, which wrote deployment addresses into a tracked
# file — and left the runbook's later steps reading a config nothing had
# patched. deploy_l2_contracts.sh has always defaulted to config.json.
RELAYER_CONFIG=${RELAYER_CONFIG:-$REPO_ROOT/relayer/config.json}
case "$RELAYER_CONFIG" in
    /*) ;;
    *) RELAYER_CONFIG=$REPO_ROOT/$RELAYER_CONFIG ;;
esac

if [ -z "${ETH_RPC:-}" ]; then
    [ -f "$ENV_FILE" ] || {
        echo "ERROR: ETH_RPC is unset and no handoff at $ENV_FILE; run scripts/local/run_eth_node.sh first" >&2
        exit 1
    }
    # shellcheck disable=SC1090
    . "$ENV_FILE"
fi

: "${ETH_RPC:?ETH_RPC is required}"

# Deployer. E2ETestDeploy sets relayers[0] = msg.sender, so whoever deploys receives
# the ICS26Router relayer role — it MUST be the key the relayer runs with
# (relayer/.env ETH_PRIVATE_KEY). A different deployer leaves the relayer
# unauthorized and every updateApplicationState reverts with no reason string, which
# is indistinguishable from an unfunded signer until you run `cast run` on the tx.
#
# The defaults are the devnet-only key relayer/.env ships; override both on any real
# network. E2E_FAUCET_ADDRESS receives the test ERC20 and follows the deployer unless
# set explicitly.
ETH_DEPLOYER_ADDRESS=${ETH_DEPLOYER_ADDRESS:-0x8943545177806ED17B9F23F0a21ee5948eCaa776}
ETH_DEPLOYER_PRIVATE_KEY=${ETH_DEPLOYER_PRIVATE_KEY:-bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31}
E2E_FAUCET_ADDRESS=${E2E_FAUCET_ADDRESS:-$ETH_DEPLOYER_ADDRESS}
export E2E_FAUCET_ADDRESS

echo "ETH_RPC: $ETH_RPC"
echo "ETH_WS: ${ETH_WS:-}"
echo "ETH_BEACON_API: ${ETH_BEACON_API:-}"
echo "deployer: $ETH_DEPLOYER_ADDRESS"
echo "faucet:   $E2E_FAUCET_ADDRESS"

# Deploy ETH contracts
RESULT=$(forge script scripts/E2ETestDeploy.s.sol:E2ETestDeploy \
    --rpc-url "$ETH_RPC" \
    --broadcast \
    --ffi \
    --sender "$ETH_DEPLOYER_ADDRESS" --private-key "$ETH_DEPLOYER_PRIVATE_KEY" \
    2>/dev/null
)

# The script returns the deployment JSON as `0: string "<escaped-json>"`; unwrap it once.
DEPLOYMENT_JSON=$(echo "$RESULT" \
  | sed -n 's/^0: string "\(.*\)".*/\1/p' \
  | sed 's/\\"/"/g')

ERC20_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.erc20')
ICS20_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.ics20Transfer')
ICS26_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.ics26Router')
VERIFIER_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.signatureVerifier')
MEMBERSHIP_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.membership')
UPDATE_CLIENT_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.updateClient')
MISBEHAVIOUR_ADDRESS=$(echo "$DEPLOYMENT_JSON" | jq -r '.misbehaviour')

echo "ERC20_ADDRESS: $ERC20_ADDRESS"
echo "ICS20_ADDRESS: $ICS20_ADDRESS"
echo "ICS26_ADDRESS: $ICS26_ADDRESS"
echo "VERIFIER_ADDRESS: $VERIFIER_ADDRESS"
echo "MEMBERSHIP_ADDRESS: $MEMBERSHIP_ADDRESS"
echo "UPDATE_CLIENT_ADDRESS: $UPDATE_CLIENT_ADDRESS"
echo "MISBEHAVIOUR_ADDRESS: $MISBEHAVIOUR_ADDRESS"

# Patch the relayer config with the deployed addresses + endpoints.
jq \
  --arg ETH_RPC "$ETH_RPC" \
  --arg ETH_WS "${ETH_WS:-}" \
  --arg ICS26 "$ICS26_ADDRESS" \
  --arg WRAP "$VERIFIER_ADDRESS" \
  --arg MEMB "$MEMBERSHIP_ADDRESS" \
  --arg UPCL "$UPDATE_CLIENT_ADDRESS" \
  --arg MIS "$MISBEHAVIOUR_ADDRESS" \
  --arg ETH_BEACON "${ETH_BEACON_API:-}" '
    (.. | objects | select(has("eth_rpc_url")) | .eth_rpc_url) = $ETH_RPC
  | (.. | objects | select(has("eth_ws_url")) | .eth_ws_url) = $ETH_WS
  | (.. | objects | select(has("ics26_address")) | .ics26_address) = $ICS26
  | (.. | objects | select(has("signature_verifier")) | .signature_verifier) = $WRAP
  | (.. | objects | select(has("membership")) | .membership) = $MEMB
  | (.. | objects | select(has("update_client")) | .update_client) = $UPCL
  | (.. | objects | select(has("misbehaviour")) | .misbehaviour) = $MIS
  | (.. | objects | select(has("eth_beacon_api_url")) | .eth_beacon_api_url) = $ETH_BEACON
  ' "$RELAYER_CONFIG" > "$RELAYER_CONFIG.tmp" && mv "$RELAYER_CONFIG.tmp" "$RELAYER_CONFIG"

echo "Patched $RELAYER_CONFIG with the deployed addresses."
