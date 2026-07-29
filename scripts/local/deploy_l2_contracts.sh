#!/bin/sh

set -euxo pipefail

# Deploys the L2-side receive infrastructure (E2ETestDeployL2) to a running local
# rollup's exec layer and patches the relayer's cosmos_to_l2 module with the
# resulting addresses + L2 endpoint. This is the L2 twin of deploy_eth_contracts.sh:
# same deploy-then-patch shape, but it targets the L2 exec RPC and the cosmos_to_l2
# module (a SpectreClient on the L2), and it runs E2ETestDeployL2 (no ERC20 faucet).
#
# Like the L1 flow it is DEPLOY-ONLY: the SpectreClient itself is created afterwards
# by `relayer create-clients-eth --source <ics26_client_id>`, which reads a fresh
# Cosmos genesis and writes spectre_client back into the module.
#
# Endpoint source (in priority order):
#   1. L2_RPC in the environment (explicit override)
#   2. L2_ENV_FILE, if set, is sourced and its L2_RPC_URL used
#   3. the stack handoff for DST_CHAIN is sourced automatically:
#        opstack   -> .op-devnet-run/attestor.env       (L2_RPC_URL)
#        arbitrum  -> .arbitrum-devnet-run/attestor.env  (L2_RPC_URL)
#
# Env:
#   DST_CHAIN (opstack)                 which cosmos_to_l2 module to patch
#                                       (matched on src_chain=cosmos + dst_chain)
#   RELAYER_CONFIG                      relayer config to patch
#                                       (default: relayer/config.example.json)
#   L2_RPC / L2_ENV_FILE                see endpoint source above
#   L2_WS                               optional L2 exec WS to patch (eth_ws_url)
#   L2_DEPLOYER_PRIVATE_KEY / _ADDRESS  deployer funded on the L2 (defaults to the
#                                       well-known devnet key, prefunded on Kurtosis
#                                       OP + Arbitrum local stacks)
#   PERMIT2                             optional Permit2 address (default: unset)
#
# All internal paths (scripts/E2ETestDeployL2.s.sol, relayer/) are repo-root relative.
cd "$(dirname "$0")/../.."
REPO_ROOT=$PWD

DST_CHAIN=${DST_CHAIN:-opstack}
RELAYER_CONFIG=${RELAYER_CONFIG:-$REPO_ROOT/relayer/config.example.json}
case "$RELAYER_CONFIG" in
    /*) ;;
    *) RELAYER_CONFIG=$REPO_ROOT/$RELAYER_CONFIG ;;
esac

# ---------------------------------------------------------------- endpoint ---
if [ -z "${L2_RPC:-}" ]; then
    if [ -z "${L2_ENV_FILE:-}" ]; then
        case "$DST_CHAIN" in
            opstack) L2_ENV_FILE=$REPO_ROOT/.op-devnet-run/attestor.env ;;
            arbitrum) L2_ENV_FILE=$REPO_ROOT/.arbitrum-devnet-run/attestor.env ;;
            *) echo "ERROR: unknown DST_CHAIN=$DST_CHAIN (want opstack|arbitrum)" >&2; exit 1 ;;
        esac
    fi
    [ -f "$L2_ENV_FILE" ] || {
        echo "ERROR: L2_RPC is unset and no handoff at $L2_ENV_FILE; run the $DST_CHAIN stack first (scripts/local/run_optimism_node.sh or run_arbitrum_node.sh) or set L2_RPC" >&2
        exit 1
    }
    # shellcheck disable=SC1090
    . "$L2_ENV_FILE"
    L2_RPC=${L2_RPC:-${L2_RPC_URL:-}}
    L2_WS=${L2_WS:-${L2_WS_URL:-}}
fi

: "${L2_RPC:?L2_RPC is required (no L2 exec RPC discovered)}"

# Deployer. E2ETestDeployL2 sets relayers[0] = msg.sender, so whoever deploys receives
# the ICS26Router relayer role — it MUST be the key the relayer runs with
# (relayer/.env ETH_PRIVATE_KEY). A different deployer leaves the relayer unauthorized
# and every updateApplicationState reverts with no reason string, which looks exactly
# like an unfunded signer until you run `cast run` on the tx and see canCall -> false.
# Funding the wrong address does not fix it.
#
# Default is the well-known devnet key, prefunded on both stacks; override for a real L2.
L2_DEPLOYER_ADDRESS=${L2_DEPLOYER_ADDRESS:-0x8943545177806ED17B9F23F0a21ee5948eCaa776}
L2_DEPLOYER_PRIVATE_KEY=${L2_DEPLOYER_PRIVATE_KEY:-bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31}

echo "DST_CHAIN: $DST_CHAIN"
echo "L2_RPC:    $L2_RPC"
echo "L2_WS:     ${L2_WS:-}"
echo "deployer:  $L2_DEPLOYER_ADDRESS"

# Fail loud if there is no cosmos_to_l2 module to patch, BEFORE spending a deploy.
MATCHES=$(jq --arg DST "$DST_CHAIN" \
  '[.modules[] | select(.src_chain=="cosmos" and .dst_chain==$DST)] | length' \
  "$RELAYER_CONFIG")
if [ "$MATCHES" -eq 0 ]; then
    echo "ERROR: no cosmos_to_l2 module with src_chain=cosmos dst_chain=$DST_CHAIN in $RELAYER_CONFIG; add one (see the cosmos-to-op example) before deploying" >&2
    exit 1
fi

# ------------------------------------------------------------------ deploy ---
DEPLOY_OUT=$REPO_ROOT/.l2-deployment.json
export E2E_L2_DEPLOYMENT_OUT=$DEPLOY_OUT
[ -n "${PERMIT2:-}" ] && export PERMIT2

forge script scripts/E2ETestDeployL2.s.sol:E2ETestDeployL2 \
    --rpc-url "$L2_RPC" \
    --broadcast \
    --ffi \
    --sender "$L2_DEPLOYER_ADDRESS" --private-key "$L2_DEPLOYER_PRIVATE_KEY"

# E2ETestDeployL2 writes the address handoff to E2E_L2_DEPLOYMENT_OUT.
[ -f "$DEPLOY_OUT" ] || { echo "ERROR: deploy output $DEPLOY_OUT not written" >&2; exit 1; }

ICS26_ADDRESS=$(jq -r '.ics26Router' "$DEPLOY_OUT")
ICS20_ADDRESS=$(jq -r '.ics20Transfer' "$DEPLOY_OUT")
VERIFIER_ADDRESS=$(jq -r '.signatureVerifier' "$DEPLOY_OUT")
MEMBERSHIP_ADDRESS=$(jq -r '.membership' "$DEPLOY_OUT")
UPDATE_CLIENT_ADDRESS=$(jq -r '.updateClient' "$DEPLOY_OUT")
MISBEHAVIOUR_ADDRESS=$(jq -r '.misbehaviour' "$DEPLOY_OUT")

echo "ICS26_ADDRESS:      $ICS26_ADDRESS"
echo "ICS20_ADDRESS:      $ICS20_ADDRESS"
echo "VERIFIER_ADDRESS:   $VERIFIER_ADDRESS"
echo "MEMBERSHIP_ADDRESS: $MEMBERSHIP_ADDRESS"
echo "UPDATE_CLIENT:      $UPDATE_CLIENT_ADDRESS"
echo "MISBEHAVIOUR:       $MISBEHAVIOUR_ADDRESS"

# ------------------------------------------------------------------- patch ---
# Patch ONLY the matching cosmos_to_l2 module (leave the L1 cosmos_to_eth module and
# any other L2 family untouched). spectre_client stays empty — create-clients-eth
# fills it once it deploys the SpectreClient against a fresh Cosmos genesis.
jq \
  --arg DST "$DST_CHAIN" \
  --arg RPC "$L2_RPC" \
  --arg WS "${L2_WS:-}" \
  --arg ICS26 "$ICS26_ADDRESS" \
  --arg VERIF "$VERIFIER_ADDRESS" \
  --arg MEMB "$MEMBERSHIP_ADDRESS" \
  --arg UPCL "$UPDATE_CLIENT_ADDRESS" \
  --arg MIS "$MISBEHAVIOUR_ADDRESS" '
    .modules |= map(
      if .src_chain == "cosmos" and .dst_chain == $DST then
          .config.eth_rpc_url = $RPC
        | (if $WS != "" then .config.eth_ws_url = $WS else . end)
        | .config.ics26_address = $ICS26
        | .config.signature_verifier = $VERIF
        | .config.membership = $MEMB
        | .config.update_client = $UPCL
        | .config.misbehaviour = $MIS
      else . end
    )
  ' "$RELAYER_CONFIG" > "$RELAYER_CONFIG.tmp" && mv "$RELAYER_CONFIG.tmp" "$RELAYER_CONFIG"

echo "Patched $RELAYER_CONFIG cosmos_to_l2 (dst_chain=$DST_CHAIN) with the deployed L2 addresses."
echo "Next: relayer create-clients-eth --source <ics26_client_id>  (deploys the SpectreClient on the L2)."
