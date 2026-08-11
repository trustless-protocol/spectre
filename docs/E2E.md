# End-to-end runbooks

Every local bring-up: Cosmos↔Ethereum, Cosmos↔OP, Cosmos↔Arbitrum, Cosmos↔Base,
and running against an L2 someone else operates. Split out of the README, which
had grown past the point where any of it was findable.

Prerequisites are in [the README](../README.md#requirements); these runbooks assume
the toolchain is in place and add Docker + Kurtosis on top.

**Read this first, whichever chain you are bringing up.**

Order is load-bearing. The Tendermint light client on the EVM side can only be
created once the L1 beacon has finalized an epoch, and the Ethereum light client on
Cosmos needs its wasm stored by governance first. Out of order, the failures land
several steps later and name the wrong thing.

Two things that are easy to get wrong before you start:

- **The circuit artifacts are not in the repo.** `.gitignore` excludes every
  `contracts/verifiers/Groth16Verifier_N*.sol`, so a fresh clone has none and
  `forge build` fails on `E2ETestDeployL2`'s import. Step 1 of the first runbook
  below is mandatory, not optional — and the `bin/` it writes must stay paired with
  the verifiers it emitted, because each `Setup()` produces a different verifying
  key. Deploying verifiers from one run and proving with another's `pk.bin` makes
  every proof revert on chain.
- **Gaia must be the custom-host build.** The scripts require the checkout to be on
  `test/ibc-host-customs`; see [Gaia binary](#gaia-binary) below.

### One L1, shared between rollups

Every rollup script calls `run_eth_node.sh` to guarantee an L1 exists in its enclave,
and **that script reuses a running enclave rather than rebuilding it**. Three
consequences worth knowing before you start:

- The rollup scripts may be run in **any order**, and re-running one does not disturb
  the L2 services, contracts or chain state of the others sharing that L1.
- `run_eth_node.sh` still re-discovers the endpoints and rewrites `eth.env` on every
  call, so a caller can source it unconditionally.
- `FORCE_RECREATE=1` destroys the enclave and rebuilds the L1 from scratch — and takes
  **every rollup in it** with it. That is the flag for a genuinely clean L1, not for
  routine re-runs.

Sharing is per **enclave**, and the defaults differ: `run_optimism_node.sh` and
`run_arbitrum_node.sh` both use `op-devnet`, so they share out of the box, while the
plain ETH↔Cosmos devnet (`run_eth_node.sh` on its own) uses `my-testnet` and is a
separate L1. Set `ENCLAVE` explicitly to put them together.

## Local Cosmos ↔ Ethereum E2E

End-to-end run on local Cosmos + Ethereum nodes. Requires Docker + Kurtosis on
top of the toolchain in [Requirements](../README.md#requirements).

### Gaia binary

The Cosmos scripts need the custom Gaia build carrying the IBC host changes — most
visibly the 08-wasm Stargate allowlist entry for `ClientStatus`, without which L2
client creation fails with `status Unknown`. They default to `../gaia/build/gaiad`
and require that checkout to be **on** `test/ibc-host-customs`:

```bash
cd ../gaia
git checkout test/ibc-host-customs
GOTOOLCHAIN=go1.25.7 make build
cd ../fast-ibc
```

`GAIAD=/path/to/gaiad` or `GAIA_DIR=/path/to/gaia` overrides that, and `GAIAD=`
takes priority over the branch check entirely — which is what you want when the
checkout is on a detached HEAD, or when the binary is already installed elsewhere.
The scripts never fall back to a PATH `gaiad` on their own, because a stock binary
produces a chain that fails several steps later for unrelated-looking reasons.

If you override, confirm you pointed at the right build:

```bash
$GAIAD version                         # -> test/ibc-host-customs-<sha>
strings $GAIAD | grep -c ClientStatus  # -> non-zero
```

> **Tip for local dev**: `relayer/prover/buckets.go` defaults to
> `Buckets = []int{4, 8, 16, 32, 64, 128}`. Compiling all six takes 30+ minutes
> (bucket 128 alone is ~30 min) and the local single-validator chain only ever
> uses bucket 4. For local iteration, temporarily edit it to:
>
> ```go
> var Buckets = []int{4}
> ```
>
> before step 1 below. Revert before pushing — production needs the full set.
>
> The repo ships **no** generated verifiers — `.gitignore` excludes
> `contracts/verifiers/Groth16Verifier_N*.sol`, so `git ls-tree` lists none and a
> fresh clone has none. Step 1 below writes them, and `forge build` fails without
> it (`scripts/E2ETestDeployL2.s.sol` imports `Groth16Verifier_N4`).
>
> Whatever buckets you compile must cover the chain's quorum: a signer count with
> no bucket makes `SignatureVerifier.verifyBatchProof` revert `UnknownBucket(N)`.
> Chains larger than the local devnet need the bigger buckets built and deployed.

```bash
# 1. REQUIRED on a fresh clone: compile per-bucket circuits + emit
#    Groth16Verifier_N{N}.sol. Neither the circuit artifacts nor the verifiers are
#    committed. Re-run when circuit code changes — and redeploy the verifiers with
#    it, because a new Setup() means a new verifying key and stale verifiers reject
#    every proof.
#    CPU default:
cd relayer
go run ./prover/cmd ./bin ../contracts/verifiers
#    GPU variant:
#    go run -tags=icicle ./prover/cmd -gpu-prove ./bin ../contracts/verifiers

# 2. Build the relayer binary
go build -o relayer ./cmd
#    GPU build:
#    go build -tags=icicle -o relayer ./cmd

# 3. Start Ethereum first and wait until the beacon node finalizes.
#    Beacon RPC is pinned to 32101 via eth-network-params.yaml (public_port_start).
#    run_eth_node.sh owns only the node (endpoints -> .eth-devnet-run/eth.env);
#    deploy_eth_contracts.sh is the deploy half (reads that handoff).
./scripts/local/run_eth_node.sh         # Kurtosis Ethereum testnet (node only)
# Poll until finalized.epoch > 0:
curl -s http://127.0.0.1:32101/eth/v1/beacon/states/head/finality_checkpoints
#    The deploy patches relayer/config.json, so create it first. Each path has a
#    runnable example carrying just that path's module pair; config.example.json
#    is the catalogue of all of them and is NOT runnable on its own.
cp relayer/config.ethereum.example.json relayer/config.json
./scripts/local/deploy_eth_contracts.sh # deploy core contracts + patch relayer config
#    Deploy with the SAME key the relayer runs with: E2ETestDeploy sets
#    relayers[0] = msg.sender, so the deployer receives the ICS26Router relayer role.
#    Defaults to the devnet key relayer/.env ships; on any other network set
#    ETH_DEPLOYER_ADDRESS + ETH_DEPLOYER_PRIVATE_KEY (and E2E_FAUCET_ADDRESS if the
#    test ERC20 should go elsewhere).

# 4. Then start Cosmos and submit the Ethereum LC WASM via governance
./scripts/local/run_cosmos_node.sh   # local Cosmos chain with funded test accounts
./scripts/local/wasm.sh              # submit + vote-pass the Ethereum LC WASM proposal
#
# Multi-validator alternative — spin up N nodes (180 default) with a
# Cosmos-Hub-like staked-power distribution so the relayer exercises a
# real signer-selection path (top ~20 hold ~2/3). Each node's RPC is
# striped from :31000.
#    NUM_NODES=20 ./scripts/local/run_cosmos_node_n.sh
#    NUM_NODES=20 ./scripts/local/wasm_n.sh
# After multi-node bring-up, point `relayer/config.example.json` and any
# `gaiad` --node / --home flags at the val0 home + RPC (defaults
# $HOME/.gaia-multi/val0 + tcp://127.0.0.1:31000).

# 5. Create the light clients, one command per chain. The Cosmos side comes first
#    (creates the 08-wasm ETH client, writing cosmos_wasm_client_id), then the ETH
#    side deploys the Tendermint light client wired to that id (writing
#    spectre_client). Both write back into relayer/config.json automatically.
./relayer create-clients-cosmos --config config.json --wasm-checksum <hex-from-wasm.sh>
./relayer create-clients-eth    --config config.json --trust-level 2/3

# 6. Start the bi-directional relay loop. Use the SAME config the steps above
#    patched — deploy_eth_contracts.sh and create-clients-* write into config.json.
./relayer start --config config.json
#    GPU run:
#    ./relayer start --config config.json --gpu-prove
#
#    Multiple Cosmos sources: add one `cosmos_to_eth` module per source to
#    config.json, each with a distinct `ics26_client_id` — the ETH router's
#    client id for that Cosmos chain (config.example.json uses "cosmoshub-1";
#    a second source might be "osmosis-1"). Create its clients with --source,
#    then start once — `start` runs one independent relay loop per source in
#    the same process (shared prover + ETH endpoint):
#      ./relayer create-clients-cosmos --config config.json --source osmosis-1 --wasm-checksum <hex>
#      ./relayer create-clients-eth    --config config.json --source osmosis-1
#      ./relayer start --config config.json

# 7. send packet

. ./scripts/local/gaiad_binary.sh
ABS_TIMEOUT=$(($(date +%s) + 2000))

"$GAIAD" tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 1000stake \
  --from test1 \
  --home "$HOME/.gaia" \
  --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 \
  --keyring-backend test \
  --gas-prices 1stake \
  --absolute-timeouts \
  --packet-timeout-timestamp "$ABS_TIMEOUT" \
  --generate-only \
| jq '.body.messages[0].encoding = "application/x-solidity-abi"' \
| "$GAIAD" tx sign /dev/stdin \
    --from test1 \
    --home "$HOME/.gaia" \
    --chain-id test-ibc-eth \
    --keyring-backend test \
| "$GAIAD" tx broadcast /dev/stdin \
    --node tcp://127.0.0.1:26657 \
    -y


# 8. Check it landed. Addresses come from step 3's output, not from here —
#    ICS20Transfer is printed by deploy_eth_contracts.sh, and the wrapper is
#    derived from the denom trace (which carries the client the token arrived on).
WRAPPED=$(cast call <ICS20Transfer> 'ibcERC20Contract(string)(address)' \
  'transfer/<cosmos-client-on-eth>/stake' --rpc-url <eth-rpc>)
cast call "$WRAPPED" 'balanceOf(address)(uint256)' <receiver> --rpc-url <eth-rpc>

#    And the ack on the Cosmos side:
"$GAIAD" q txs --query "message.action='/ibc.core.channel.v2.MsgAcknowledgement'" \
  --node tcp://127.0.0.1:26657 -o json
```

Send an ICS-20 transfer from Cosmos to trigger a client update + `recvPacket`
round-trip; the `[UpdateCosmosClient]` log line reports the chosen bucket.

## Local Cosmos ↔ OP E2E

Brings up an L1 + OP Stack L2 + Cosmos and relays both directions. The L1 is shared:
`run_optimism_node.sh` calls `run_eth_node.sh` for it (Fusaka-from-genesis, the same
`eth-network-params.yaml` the ETH↔Cosmos devnet uses), which reuses the enclave's L1
if one is already up — see [One L1, shared between rollups](#one-l1-shared-between-rollups).

```bash
# 1. L1 (Fulu) + OP Stack L2 in one Kurtosis enclave. Ends with games proposed
#    and .op-devnet-run/attestor.env written (L1/L2/op-node/beacon endpoints).
./scripts/local/run_optimism_node.sh

# 2. Attestor — independent verifier over the replica op-node; serves gRPC :3001.
#    Sources attestor.env automatically. Every attestor script defaults to
#    GRPC_PORT=3001, so a host running several needs distinct ports; the
#    config.example.json modules expect op :3001, arbitrum :3002, base :3003.
./scripts/local/run_op_attestor.sh

# 3. Cosmos node, then gov-store the OP light-client wasm. The Ethereum client
#    is NOT needed for a Cosmos<->OP deployment: the attestor-trusted L2 client
#    authenticates nothing through it. Add ./scripts/local/wasm.sh only if this
#    config also relays Cosmos<->Ethereum.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm_op.sh     # -> OP client checksum

# 4. IBC contracts on the L2 (E2ETestDeployL2). Deploy with the SAME key the relayer
#    will run with (relayer/.env ETH_PRIVATE_KEY): E2ETestDeployL2 grants the
#    ICS26Router relayer role to msg.sender, and without that role every
#    updateApplicationState reverts. The values below are the devnet-only key that
#    relayer/.env ships as ETH_PRIVATE_KEY — change BOTH together or they drift apart.
#    The account also needs an L2 balance (the L1 devnet faucet address has none
#    there), but funding alone does NOT fix the role.
DST_CHAIN=opstack \
L2_DEPLOYER_ADDRESS=0x8943545177806ED17B9F23F0a21ee5948eCaa776 \
L2_DEPLOYER_PRIVATE_KEY=bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31 \
  ./scripts/local/deploy_l2_contracts.sh
# It patches relayer/config.json (override with RELAYER_CONFIG), both the forward
# module and the matching l2_to_cosmos module's rollup_profile.common.l2_router.

# 5. Clients. The OP client on Cosmos, then the L2 side (SpectreClient deployed
#    + addClient'd). Copy op-l2-config.example.json and fill in the checksum from
#    step 3, the L2 router address from step 4, and the L2 chain id.
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config op-l2-config.json
./relayer create-clients-eth --config config.json --source <ics26_client_id> --trust-level 2/3

# 6. Relay both directions.
./relayer start --config config.json
```

### Running against an existing L2 (no devnet L2)

The flow above brings up its own L2. To relay against a node someone else operates —
OP Sepolia, or a replica you already run — skip step 1 and point the attestor at it.
There is no local L1 to start either: the attestor-trusted client verifies nothing
against L1, so only the *attestor* needs an L1 RPC, and a public endpoint is enough.

Verified end to end against OP Sepolia; every value below was needed to get there.

```bash
# 1. Attestor against the existing replica. NETWORK must name the real chain, and
#    the dispute-game contracts cannot be auto-resolved for a public network — read
#    them off the chain once (see below) and pass them.
OP_NODE_RPC_URL=http://<host>:9545 \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
NETWORK=op-sepolia SRC_CHAIN=op-sepolia \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 \
DISPUTE_GAME_FACTORY=0x05F9613aDB30026FFd634f38e5C4dFd30a197Fa1 \
RESPECTED_GAME_TYPE=8 \
  ./scripts/local/run_op_attestor.sh

# 2. Cosmos + the OP wasm, unchanged from the devnet flow.
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_op.sh

# 3. config.json: one runnable config per path, already cut to the OP module
#    pair. The deploy below patches this file, so it has to exist first.
cp relayer/config.op.example.json relayer/config.json
# then fill in the fields listed in "Values you fill in by hand" below.

# 4. Deploy the L2 contracts with a key funded ON THAT L2.
#    DO THIS BEFORE STARTING THE ATTESTOR if both share one rate-limited RPC —
#    see "One RPC budget, two processes" in the Arbitrum section.
DST_CHAIN=opstack L2_RPC=http://<host>:8545 RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh

# 5. l2-config, then clients and relay. The relayer must sign with the SAME key
#    that deployed (E2ETestDeployL2 grants the router role to msg.sender). Pass it
#    in the environment rather than editing relayer/.env — godotenv does not
#    override a variable already set, so the export wins.
cp relayer/op-l2-config.example.json relayer/op-l2-config.json
# fill in wasm_checksum (step 2), l2_rpc_url, and l2_router (= ICS26_ADDRESS from step 4)
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config op-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --source <ics26_client_id> --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

#### Values you fill in by hand, and where each comes from

Everything not listed here is written for you — see
[What the tooling fills in for you](#what-the-tooling-fills-in-for-you).

**`relayer/config.json`:**

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-op` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-op` | your L2 WebSocket endpoint; optional, the relayer skips the WS client when empty |
| `l2_rpc_url` | `op-to-cosmos` | the same L2 HTTP endpoint (step 4 overwrites it with `L2_RPC` anyway) |
| `attestor_addr` | `op-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1; the example ships `127.0.0.1:3001`, which matches the script default |
| `attestor_src_chain` | `op-to-cosmos` | the `SRC_CHAIN` you passed in step 1 (`op-sepolia`) — must match, or `AttestedUpTo` answers for another chain |
| `head_kind` | `op-to-cosmos` | the same choice as `ATTESTATION_HEAD` in step 1 |
| `rollup_profile.common.l2_chain_id` | `op-to-cosmos` | `eth_chainId` on your L2 (`11155420` on OP Sepolia) |
| `log_scan_chunk` | `op-to-cosmos` | your provider's `eth_getLogs` span cap — measure it, see the Arbitrum section |

**`relayer/op-l2-config.json`:** `wasm_checksum` (from `wasm_op.sh`), `l2_rpc_url`, and
`rollup_profile.common.l2_router` (= `ICS26_ADDRESS` printed by step 4 — which is why
this file is filled in after the deploy).

**Resolving the dispute-game contracts.** `run_op_attestor.sh` can derive them with
`cast` for a chain it knows; for a public network it exits with
`set DISPUTE_GAME_FACTORY and RESPECTED_GAME_TYPE`. Read them from the chain — the
portal address comes from op-node itself, so nothing is hardcoded:

```bash
PORTAL=$(curl -s -X POST -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"optimism_rollupConfig","params":[]}' \
  http://<host>:9545 | jq -r .result.deposit_contract_address)
cast call $PORTAL "disputeGameFactory()(address)" --rpc-url $L1_RPC_URL
cast call $PORTAL "respectedGameType()(uint32)"  --rpc-url $L1_RPC_URL
```

**The replica must serve `eth_getProof` below its head.** The relayer proves the
router account at the attested height, which trails the head, and op-reth's default
proof window is small — a node serving only `latest` fails every build with
`distance to target block exceeds maximum proof window`. The mode-2 bring-up in
`run_op_attestor.sh` starts op-reth with `--rpc.eth-proof-window=100000`; an external
replica needs the operator to have done the same. Check before starting:

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getProof '["<router>",[],"'$(printf '0x%x' $((H-1000)))'"]' --rpc-url $L2
```

**Gaia must be the custom-host build.** `gaiad_binary.sh` requires the checkout to be
*on* branch `test/ibc-host-customs`, which a detached HEAD fails even when the code is
right. Point at an already-built binary instead — `GAIAD=` takes priority and skips
the branch check. Confirm it is the right one:

```bash
$GAIAD version                       # → test/ibc-host-customs-<sha>
strings $GAIAD | grep -c ClientStatus  # the 08-wasm allowlist entry L2 clients need
```

**Set `DERIVED_GAP_BLOCKS` explicitly.** The local devnet handoffs export `5`; an
attestor attached to an external replica inherits nothing and takes the default 150,
which puts a ~5 minute floor under the return direction. See the latency section above.

## Local Cosmos ↔ Arbitrum E2E

Brings up an L1 + Arbitrum Nitro/BoLD L2 + Cosmos and relays both directions. The L1
is shared with the OP flow: `run_arbitrum_node.sh` attaches Arbitrum to the existing
Kurtosis enclave if OP already brought one up, or creates the same local L1 itself
through `run_eth_node.sh` when the enclave does not exist — see
[One L1, shared between rollups](#one-l1-shared-between-rollups).

```bash
# 1. L1 (Fulu) + Arbitrum Nitro/BoLD L2 in one Kurtosis enclave. Ends with a
#    finalized AssertionCreated check and .arbitrum-devnet-run/attestor.env written
#    (L1/L2/beacon/Nitro feed/RollupCore endpoints).
./scripts/local/run_arbitrum_node.sh

# 2. Attestor — independent non-sequencing Nitro replica; serves gRPC :3001.
#    Sources .arbitrum-devnet-run/attestor.env automatically.
GRPC_PORT=3002 ./scripts/local/run_arbitrum_attestor.sh

# 3. Cosmos node, then gov-store the Arbitrum light-client wasm. The Ethereum
#    client is NOT needed for a Cosmos<->Arbitrum deployment — same as OP, the
#    attestor-trusted L2 client authenticates nothing through it. Add
#    ./scripts/local/wasm.sh only if this config also relays Cosmos<->Ethereum.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm_arb.sh    # -> Arbitrum client checksum

# 4. IBC contracts on the Arbitrum L2 (E2ETestDeployL2). Use the same relayer key
#    rules as OP: the deployer receives the ICS26Router relayer role, so it must be
#    the key in relayer/.env ETH_PRIVATE_KEY unless you change both together.
DST_CHAIN=arbitrum \
L2_DEPLOYER_ADDRESS=0x8943545177806ED17B9F23F0a21ee5948eCaa776 \
L2_DEPLOYER_PRIVATE_KEY=bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31 \
  ./scripts/local/deploy_l2_contracts.sh
# It patches the cosmos_to_l2 module whose dst_chain is arbitrum.

# 5. Clients. Cosmos side first, one invocation per kind, then the L2 side
#    (SpectreClient deployed + addClient'd). The Arbitrum client no longer anchors to
#    an Ethereum client, so create it on its own:
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config <arb-l2-config.json>
#
#    Only run the Ethereum half (--wasm-checksum, no --l2-config) if this config also
#    relays Cosmos<->Ethereum AND that client does not exist yet. Re-running it against
#    a working deployment rewrites cosmos_to_eth.cosmos_wasm_client_id and repoints the
#    path at a client the Sepolia-side SpectreClient was never registered against; every
#    recvPacket then reverts on a counterparty mismatch.
./relayer create-clients-eth --config config.json --source <ics26_client_id> --trust-level 2/3

# 6. Relay both directions.
./relayer start --config config.json
```

Copy `relayer/arb-l2-config.example.json` and fill in three values. It is the same
shape as OP's and Base's `--l2-config` — five profile keys and nothing else. There is
no `protocol` block, no `assertions_mapping_slot` and no `rollup` address in it: since
#345/#347 the client verifies no assertion, so nothing about BoLD reaches it. Those
values configure the **attestor**, not the client.

- `wasm_checksum`: checksum printed by `wasm_arb.sh`.
- `l2_rpc_url`: `.arbitrum-devnet-run/attestor.env` `L2_RPC_URL`.
- `rollup_profile.common.l2_router`: `ICS26_ADDRESS` from `deploy_l2_contracts.sh`.
- `rollup_profile.common.l2_chain_id`: `L2_CHAIN_ID` from the handoff (`412346` on the
  local devnet, `421614` on Arbitrum Sepolia).
- `counterparty_client_id`: must equal the `cosmos_to_l2` module's `ics26_client_id`
  (`arb-client-0` in the example config) — both name the same client on the L2 router.

`l2_header_fork` is `london`, not `prague` as on OP and Base: Nitro emits none of the
post-London optional header fields.

### Running against an existing Arbitrum L2 (no devnet L2)

The Arbitrum counterpart of the OP section above. There is no local L1 and no local
Nitro: the attestor reads a public Arbitrum RPC and a public Sepolia RPC, and the
client verifies nothing against L1, so nothing else needs an L1 endpoint.

Verified end to end against **Arbitrum Sepolia**; every value below was needed.

```bash
# 1. Attestor. CHAIN_PROFILE supplies the RollupCore address, both chain ids, the
#    BoLD slot/offset and a scan start block near the head, all from
#    attestor/arbitrum/config.arbitrum-sepolia.json — pass only the endpoints.
CHAIN_PROFILE=arbitrum-sepolia \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
L2_RPC_URL=<arbitrum-sepolia-rpc> \
L2_WS_URL=<arbitrum-sepolia-ws> \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 GRPC_PORT=3002 DETACH=1 \
  ./scripts/local/run_arbitrum_attestor.sh

# 2. Cosmos + the Arbitrum wasm, unchanged from the devnet flow.
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_arb.sh          # -> checksum

# 3. config.json: one runnable config per path, already cut to the Arbitrum pair.
cp relayer/config.arbitrum.example.json relayer/config.json
# then fill in the seven fields listed in "Values you fill in by hand" below.

# 4. Deploy the L2 contracts with a key funded ON Arbitrum Sepolia.
#    DO THIS BEFORE STARTING THE ATTESTOR if both share one rate-limited RPC —
#    see "One RPC budget, two processes" below.
DST_CHAIN=arbitrum L2_RPC=<arbitrum-sepolia-rpc> RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh
# Patches BOTH modules: ics26_address on the forward one, rollup_profile.common.l2_router
# and l2_rpc_url on the return one.

# 5. l2-config, then clients. The relayer must sign with the SAME key that deployed
#    (E2ETestDeployL2 grants the router role to msg.sender). Pass it in the
#    environment rather than editing relayer/.env — godotenv does not override a
#    variable already set, so the export wins.
cp relayer/arb-l2-config.example.json relayer/arb-l2-config.json
# fill in wasm_checksum (step 2), l2_rpc_url, and l2_router (= ICS26_ADDRESS from step 4)
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config arb-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --source arb-client-0 --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

**`l2_wasm_client_id` is an output, not an input.** `create-clients-cosmos` creates that
client and writes its id back into the module, overwriting whatever is there — the
example's `08-wasm-N` placeholder or an empty string alike. It used to *reject* an empty
one, forcing a made-up id that it then overwrote (#309); it no longer does. Every other
command still requires it, because by then the client must exist.

**Set `log_scan_chunk` to your provider's `eth_getLogs` span cap.** The relayer scans
L2 packet logs over ranges, and a provider that caps the span rejects the call rather
than truncating it. Alchemy's free tier caps at **10 blocks**, drpc at 10 000; `0`
means one call per range and is right for a provider with no cap. Check yours:

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getLogs '[{"fromBlock":"'$(printf '0x%x' $((H-1000)))'","toBlock":"'$(printf '0x%x' $H)'"}]' --rpc-url $L2
```

**The L2 RPC must serve `eth_getProof` well below the head.** The relayer proves the
router account at the attested height, which trails the head. A pruned node fails every
build. Alchemy serves it at archive depth; check before starting:

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getProof '["<router>",[],"'$(printf '0x%x' $((H-5000)))'"]' --rpc-url $L2
```

**A leftover `.arbitrum-devnet-run/` no longer hijacks the run.** The devnet handoff
exports `ROLLUP_CORE_ADDRESS` and both chain ids as real environment variables, which
beat the profile — a stale one used to silently produce `l2_chain_id=412346` and
`src_chain=arbdev` against a public RPC. With `CHAIN_PROFILE` set to anything but
`devnet` the handoff is now ignored, and the script says so.

#### Values you fill in by hand, and where each comes from

Ten values. Everything else in both files is written for you — the table after this one
says by what. Verified by running this section end to end against Arbitrum Sepolia.

**`relayer/config.json`** — seven:

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-arbitrum` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-arbitrum` | your L2 WebSocket endpoint. Optional — the relayer skips the WS client when it is empty, and the L2→Cosmos direction polls regardless |
| `l2_rpc_url` | `arbitrum-to-cosmos` | the same L2 HTTP endpoint. Step 4 overwrites this with `L2_RPC`, so a mismatch here is corrected rather than fatal |
| `attestor_addr` | `arbitrum-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1. The example ships `127.0.0.1:3002`, which matches the `GRPC_PORT=3002` above — **change it if you used a different port**, or the relayer dials a closed socket |
| `attestor_src_chain` | `arbitrum-to-cosmos` | `src_chain` in `attestor/arbitrum/config.<profile>.json` (`arbitrum-sepolia`). Must equal what the attestor reports, or `AttestedUpTo` answers for a chain you did not ask about |
| `head_kind` | `arbitrum-to-cosmos` | `unsafe`, `safe` or `finalized` — the same choice as `ATTESTATION_HEAD` in step 1 |
| `rollup_profile.common.l2_chain_id` | `arbitrum-to-cosmos` | `eth_chainId` on your L2 (`421614` on Arbitrum Sepolia, `412346` on the local devnet) |
| `log_scan_chunk` | `arbitrum-to-cosmos` | your provider's `eth_getLogs` span cap — measure it, see below |

**`relayer/arb-l2-config.json`** — three:

| Field | Value comes from |
|---|---|
| `wasm_checksum` | the hex printed by `wasm_arb.sh` in step 2 |
| `l2_rpc_url` | your L2 HTTP endpoint |
| `rollup_profile.common.l2_router` | `ICS26_ADDRESS` printed by step 4 — so this file is filled in *after* the deploy, which is why step 5 comes after step 4 |

`counterparty_client_id` and `l2_chain_id` in that file already match the example config
(`arb-client-0`, `421614`); change them only if you changed the module's
`ics26_client_id` or are on a different chain.

#### What the tooling fills in for you

Identical for OP, Arbitrum and Base. Do not pre-fill these — the commands overwrite
whatever is there:

| Field | Written by |
|---|---|
| `ics26_address`, `membership`, `update_client`, `signature_verifier`, `misbehaviour` | `deploy_l2_contracts.sh` (forward module) |
| `rollup_profile.common.l2_router`, `l2_rpc_url` | `deploy_l2_contracts.sh` (return module — it prints `Patched the return module "…"`) |
| `cosmos_wasm_client_id`, `l2_wasm_client_id` | `create-clients-cosmos`, which overwrites the example's `08-wasm-N` placeholder |
| `spectre_client` | `create-clients-eth` |

#### One RPC budget, two processes

Applies to every L2 in this document; the numbers below were measured on Arbitrum
Sepolia.

**The attestor and the relayer both hammer the L2 endpoint, and on a free tier they do
not fit together.** The attestor issues `runtime_backfill_concurrency` (8) parallel
header fetches every refresh; the relayer adds `eth_getLogs` scans, `eth_getProof`
builds and transaction sends. Against one Alchemy free-tier key this measurably
exceeds the concurrent-request budget:

- `deploy_l2_contracts.sh` fails outright with `Max retries exceeded HTTP error 429
  with empty body` when the attestor is already running. Pausing the attestor
  (`docker stop fast-ibc-arbitrum-attestor`) makes the same command succeed.
- With both running, the attestor logs `backing off: failures=3 delay=40.5s
  rate_limited=true` and its frontier stops. The relayer then reports `source relayable
  height` frozen for minutes with `DERIVED_GAP_BLOCKS=10` set, which looks like a
  relayer problem and is not one.

Two ways out, in order of preference:

**Give each process its own endpoint.** The relayer does not need a keyed provider.
Measured on Arbitrum Sepolia, both of these serve what it requires with no API key:

| Endpoint | `eth_getProof` at head−500 | `eth_getLogs` span | WebSocket |
|---|---|---|---|
| `https://sepolia-rollup.arbitrum.io/rpc` | yes | ≥10 000 | — |
| `https://arbitrum-sepolia.drpc.org` | yes | ≥10 000 | `wss://arbitrum-sepolia.drpc.org` |
| `https://arbitrum-sepolia-rpc.publicnode.com` | **no** — fails below ~head−128 | ≥10 000 | `wss://…publicnode.com` |

So point `L2_RPC_URL` (attestor) at the keyed endpoint and the relayer's `eth_rpc_url` /
`l2_rpc_url` at a keyless one, and neither starves the other. Do not use publicnode for
the relayer: it cannot prove the router account at the attested height.

**Or serialise the load.** Deploy the contracts first, then start the attestor, and
lower `RUNTIME_BACKFILL_CONCURRENCY` if the attestor still trips the limit.

#### Adding an L2 to a deployment that already relays Cosmos↔Ethereum

`create-clients-cosmos` creates one kind of client per invocation. Without `--l2-config`
it creates the Ethereum light client for the Cosmos↔Ethereum path and rewrites the
`cosmos_to_eth` module's `cosmos_wasm_client_id`. With `--l2-config` it creates only the
named L2 clients and does not touch the Ethereum one:

```bash
# Adding Arbitrum to a running deployment — the Ethereum client is left alone
./relayer create-clients-cosmos --config config.json --l2-config <arb-l2-config.json>
```

`--wasm-checksum` is not needed there (nothing Ethereum-side is being created); the L2
client's own wasm checksum comes from the `--l2-config` file.

Keeping the two apart matters: re-running the Ethereum half against a working deployment
replaces a client the Ethereum-side SpectreClient is already registered against, so every
`recvPacket` starts reverting on a counterparty mismatch while the original client is
orphaned with nothing advancing it until it expires.

Several L2s can still be created in one run:

```bash
./relayer create-clients-cosmos --config config.json \
  --l2-config op.json --l2-config base.json --l2-config arb.json
```

The `--l2-config` file carries the full ICS-08 profile (see
[docs/L2_CLIENTS.md](L2_CLIENTS.md#l2-client-creation-config)). For the
attestor-trusted clients `rollup_profile.common` is five keys — `l2_chain_id`,
`l2_router`, `commitment_slot`, `profile_version`, and `l2_header_fork` — and nothing
else; an unknown key fails `instantiate` on the Rust side. `profile_version` must match
the wasm artifact (`op_attestor_v1`, `base_attestor_v1`, `arbitrum_attestor_v1`), and
`l2_header_fork` must match the chain's execution-header layout: `prague` for OP and
Base, `london` for Arbitrum Nitro, which produces none of the post-London header fields.
`packages/op-verifier/config/op-sepolia.json`,
`packages/base-verifier/config/base-sepolia.json` and
`packages/arbitrum-verifier/config/arbitrum-sepolia.json` are working public-network
templates.

### Sending a test packet

The source client is the Cosmos client that tracks the **L2** (the OP/Arbitrum client, e.g.
`08-wasm-1`), not the Ethereum client. It needs a registered counterparty, and that counterparty
must equal the client id the SpectreClient was added under on the L2 router — the relayer only
picks up packets whose `destination_client` matches its configured `ics26_client_id`.

`create-clients-cosmos` registers it for you when the l2-config carries
`counterparty_client_id`, which is the normal path; the command below is for the case where it
does not, or where you need to repoint one by hand.

```bash
. ./scripts/local/gaiad_binary.sh

"$GAIAD" tx ibc client add-counterparty <l2-client-on-cosmos> <ics26_client_id> "" \
  --from test1 --home "$HOME/.gaia" --chain-id test-ibc-eth \
  --keyring-backend test --gas-prices 1stake --gas 300000 -y

ABS_TIMEOUT=$(($(date +%s) + 2000))
"$GAIAD" tx ibc-transfer transfer transfer <l2-client-on-cosmos> <evm-receiver> 1000stake \
  --from test1 --home "$HOME/.gaia" --chain-id test-ibc-eth \
  --node tcp://127.0.0.1:26657 --keyring-backend test --gas-prices 1stake \
  --absolute-timeouts --packet-timeout-timestamp "$ABS_TIMEOUT" --generate-only \
| jq '.body.messages[0].encoding = "application/x-solidity-abi"' \
| "$GAIAD" tx sign /dev/stdin --from test1 --home "$HOME/.gaia" \
    --chain-id test-ibc-eth --keyring-backend test \
| "$GAIAD" tx broadcast /dev/stdin --node tcp://127.0.0.1:26657 -y
```

`create-clients-cosmos` writes the L2 client id into a `cosmos_to_l2` module's
`cosmos_wasm_client_id` when the l2-config carries `counterparty_client_id`; that value is what
`create-clients-eth` registers as the SpectreClient's counterparty. If it holds the *Ethereum*
client id instead, `recvPacket` reverts with a counterparty mismatch (the revert data carries the
expected and actual client ids), so set `counterparty_client_id` in the l2-config — or fix
`cosmos_wasm_client_id` by hand before `create-clients-eth`.

Two more settings decide whether the **return** direction (acks, and packets sent from the L2)
works at all, and both fail silently — nothing errors, the direction just never moves:

- `l2_ics26_client_id` on the `l2_to_cosmos` module must be the client id the SpectreClient was
  added under on the L2 router. The L2 subscriber filters events by it, so a stale placeholder
  silently drops every `WriteAcknowledgement` the L2 emits.
- Leave the attestor's `disable_derived_roots` at its default (`false`). The header builder proves
  no settlement object any more, so it accepts any attested root; turning derived roots off leaves
  the frontier waiting on games, which are posted long after the L2 block they commit, and the
  return direction simply idles for hours.

#### Return-direction latency is set by the attestor, not the prover

The forward leg costs a Groth16 proof; the return leg costs a wait. Measured on a
Cosmos↔OP Sepolia round trip against an external replica:

| step | time |
|---|---|
| prover loads `bin/n4/pk.bin` (162 MB, once per process) | 17 s |
| Groth16 proof, 604 761 constraints, CPU | 4.3 s |
| forward leg lands on the L2 (`updateApplicationState` + `recvPacket`) | ~3 s |
| **ack waits for the attestor frontier to reach its block** | **4 min 12 s** |

The wait is `DERIVED_GAP_BLOCKS`: the attestor publishes one derived root per gap, and
`RelayableHeight` follows that frontier, so a packet written just after an attestation
waits nearly a full gap. At the default 150 blocks that is ~5 minutes on a 2 s chain.

All three local devnet handoffs (`run_optimism_node.sh`, `run_base_node.sh`,
`run_arbitrum_node.sh`) export `DERIVED_GAP_BLOCKS=5`, so a devnet return leg is
seconds. Arbitrum's did not until the rename made the name uniform, which is why a
local Arbitrum devnet used to pay the 150-block wait that OP and Base did not.
**An attestor attached to an external replica inherits none of this** — it takes the
script default of 150. Pass the knob explicitly for an E2E:

```bash
DERIVED_GAP_BLOCKS=10 OP_NODE_RPC_URL=... ./scripts/local/run_op_attestor.sh
```

`run_arbitrum_attestor.sh` takes the same variable under the same name. It used to call
it `DERIVED_ATTESTATION_GAP_BLOCKS`, so an invocation copied from here was silently
ignored and the run took the 150 default with nothing in the log to explain the wait.

The forward leg costs the same on every stack. Measured on Cosmos↔Arbitrum Sepolia and
Cosmos↔Base Sepolia round trips, same machine and same day as the OP numbers above:

| step | OP Sepolia | Arbitrum Sepolia | Base Sepolia |
|---|---|---|---|
| Groth16 proof (604 761 constraints, CPU, bucket N=4) | 5 s | 6 s | 5 s |
| forward leg lands on the L2 | 3 s | 3 s | 3 s |
| `recvPacket` gas (`updateApplicationState` + `recvPacket`) | 1 509 605 | 1 506 525 | 1 509 792 |
| ack gas (`updateApplicationState` + `ackPacket`) | — | 855 883 | 855 955 |
| return leg, send → `MsgRecvPacket` on Cosmos | 4 min 25 s (gap 150) | — | ~1 min (gap 10) |

The difference in the return leg is `DERIVED_GAP_BLOCKS`, not the chain: OP ran at the
default 150, Base at 10. The Arbitrum figure is absent because its attestor stalled
mid-run (#358), so the one number obtained would measure the stall rather than the gap.

Lower is not free: each derived root is an attestation, and at `unsafe`/`safe` heads
they stay provisional until the finalized head catches up, so the provisional backlog
grows as the gap shrinks.

While a packet waits, the relayer logs one `waiting: N packet(s) not yet relayable` line
per batch period — ~100 identical lines across a default-gap wait. That is expected, not
a stall; `source relayable height` climbing is the thing to watch.

#### Sending back (L2 → Cosmos)

The forward transfer above mints a wrapped token on the L2. Sending it back is a
contract call, not a CLI command — the relayer picks the packet up from the router
event either way. Resolve the wrapper first: the denom is the full trace, so it
carries the client the token arrived on, not just `stake`.

```bash
L2=<l2-rpc>
WRAPPED=$(cast call <ICS20Transfer> "ibcERC20Contract(string)(address)" \
  "transfer/<l2-client-id>/stake" --rpc-url $L2)

cast send "$WRAPPED" "approve(address,uint256)" <ICS20Transfer> 1000 \
  --private-key $ETH_PRIVATE_KEY --rpc-url $L2

ABS_TIMEOUT=$(($(date +%s) + 2000))
cast send <ICS20Transfer> \
  "sendTransfer((address,uint256,string,string,string,uint64,string))" \
  "($WRAPPED,1000,<cosmos1-receiver>,<l2-source-client-id>,transfer,$ABS_TIMEOUT,)" \
  --private-key $ETH_PRIVATE_KEY --rpc-url $L2
```

The tuple is `SendTransferMsg` in field order (`contracts/msgs/IICS20TransferMsgs.sol`):
denom, amount, receiver, sourceClient, destPort, timeoutTimestamp, memo.

Two easy mistakes: `sourceClient` is the client id **on the L2 router** (the one the
SpectreClient was added under), not the Cosmos-side client; and `timeoutTimestamp` is
in **seconds**, the same unit trap `--absolute-timeouts` exists for on the Cosmos side.

The leg is done when the ack is relayed back and the packet leaves the pending
tracker (`[CosmosTimeoutScan] Checking 0 pending packets`) — a forward-only success
proves half the system.

### Failure modes

| Symptom | Cause |
|---|---|
| `beacon api unavailable at ...:59717` | The `eth-to-cosmos` module still has the example's beacon URL. Point `eth_beacon_api_url` at this run's `ETH_BEACON_API` (from `attestor.env`). |
| `MsgStoreCode` fails: `reference-types not enabled` | The wasm was built with a plain `cargo build`. Build through `cosmwasm/optimizer` (`just build-cw-ics08-wasm-*`); a raw release build embeds features the CosmWasm VM rejects. |
| Optimizer fails: `rustc 1.86.0 is not supported ... requires rustc 1.90` | A dependency raised its MSRV above the optimizer image's Rust. Pin the dependency down (e.g. `cargo update -p ruint --precise 1.17.0`) or bump the optimizer image. |
| `MsgCreateClient`: `status Unknown: client state is not active` | 08-wasm Stargate allowlist is missing `ClientStatus` — see [docs/L2_CLIENTS.md](L2_CLIENTS.md#host-requirements-and-verification). |
| L2 client update panics the tx: `returning attributes from a contract is not allowed` | The deployed wasm predates the fix that made the client return data only. Rebuild through `cosmwasm/optimizer` and gov-store it. |
| `updateApplicationState` reverts, ~82k gas, no revert string | The relayer's signer lacks the ICS26Router relayer role on the L2. `cast run <tx>` shows `canCall(...) → false`; funding the address does not help. |
| Send fails: `timeout exceeds the maximum expected value` | Sent without `--absolute-timeouts`, so the CLI writes `timeout_timestamp` in nanoseconds while IBC v2 reads seconds. |
| ETH client stops advancing: `404 NOT_FOUND: Sync committee for period N not found` | The beacon does not serve a `light_client/bootstrap` for that period. The relayer takes the committee from the preceding period's update instead, so this should only appear if that update is also unavailable — check the endpoint serves `/eth/v1/beacon/light_client/updates`. |
| Ack fails: `unknown field account_proof, expected one of key, value, proof` | The L2 membership proof was built in the Ethereum L1 shape. The L2 client wants a bare `EvmStorageProof` — see [docs/L2_CLIENTS.md](L2_CLIENTS.md). |
| `packet at height N not yet covered by the destination client (trusts M)` | Normal wait: the attestor frontier has not reached the packet's block yet, so `RelayableHeight` is still below it. It clears on the next attestation; if it never does, the attestor is stuck — check its log rather than the relayer's. |
| `the attestor does not recognise L2 block N ... as canonical` | The relayer's L2 RPC and the attestor's replica disagree at that height — a reorg past the frontier, or the two pointed at different chains. Not retried away: confirm both use the same L2. |
| A direction goes silent — no error, no retry, other directions healthy | A hung RPC. `kill -QUIT <relayer-pid>` dumps every goroutine; look for one blocked in `net/http.(*persistConn).roundTrip`. Note the dump kills the process. |
| Gov proposal ends `REJECTED` without votes | `wasm.sh` resolves the proposal id after a fixed `sleep`; if indexing is slower the id is empty and the vote step is skipped. Vote manually before the (short devnet) voting period ends. |
| `eth_getProof`: `distance to target block exceeds maximum proof window` | The L2 execution node will not prove below its head. The relayer proves the router account at the *attested* height, which always trails. On local Base, this usually means the relayer was pointed at the sequencer RPC instead of the follower RPC. Otherwise widen the node's proof window (`--rpc.eth-proof-window`); on someone else's node, ask — no relayer setting works around it. |
| Attestor exits: `set DISPUTE_GAME_FACTORY and RESPECTED_GAME_TYPE for network <x>` | The script resolves those from a known chain preset and cannot for a public network. Read them off the chain — see [Running against an existing L2](#running-against-an-existing-l2-no-devnet-l2). |
| `[gaiad_binary] ERROR: ... is on branch unknown; expected test/ibc-host-customs` | The Gaia checkout is on a detached HEAD, which fails the branch check even when the code is right. Set `GAIAD=` to an already-built binary — see [Gaia binary](#gaia-binary). |
| Attestor exits: `failed to listen on attestor grpc address 127.0.0.1:3001: address already in use` | An earlier attestor still holds the port. Find it with `lsof -nP -iTCP:3001 -sTCP:LISTEN` — plain `lsof -ti :3001` also matches the relayer *connected* to that port, and killing that list takes the relayer down with it. |
| `forge script` fails `insufficient funds ... have 0` | The deployer has no balance on the L2 — use an L2-funded account (step 4). Balances do not carry between rollups: funded on L1 Sepolia, OP or Arbitrum still means zero on Base. |
| `eth_getProof` returns nothing at all — no result, no error, the call just hangs | The node prunes state below some depth and does not say so. Measure the boundary (see the external-L2 sections) and keep the attested height inside it; `DERIVED_GAP_BLOCKS=10` trails ~10–20 blocks, `head_kind=finalized` trails ~600. |
| Arbitrum attestor stops attesting: last log line is a normal attestation, container still up | Its sequential per-block back-fill cannot keep pace with the chain — issue #358. Nothing is logged because the head never changes, so no error path is taken. `docker restart fast-ibc-arbitrum-attestor` re-anchors it at the head for roughly another minute. |
| Arbitrum attestor attests the wrong chain (`src_chain=arbdev`, `l2_chain_id=412346`) despite `CHAIN_PROFILE=arbitrum-sepolia` | A leftover `.arbitrum-devnet-run/attestor.env` was sourced; its exports are real env vars and beat the profile. Fixed — a non-`devnet` `CHAIN_PROFILE` now ignores the handoff and logs that it did. Delete the directory if you are on an older checkout. |
| `eth_getLogs` rejected: `you can make eth_getLogs requests with up to a 10 block range` | The provider caps the log span. Set `log_scan_chunk` in the `l2_to_cosmos` module to that cap (Alchemy free tier 10, drpc 10 000). |
| Return direction sees no events at all, forward direction fine | `l2_ics26_client_id` does not equal the forward module's `ics26_client_id`. The L2 subscriber filters events on it, so a mismatch drops every one silently. |

## Local Cosmos ↔ Base E2E

Same shape as the OP flow, with three differences that will bite if you copy the OP
steps and rename the scripts.

**Base is not vanilla OP Stack.** `base/base` ships one unified Rust node (execution
and consensus in a single process) plus its own batcher, so `run_base_node.sh` runs
base's own images from a pinned clone rather than optimism-package. It still speaks
the op-node RPC namespace, which is why the attestor is shared.

**Base gets its own enclave by default** (`base-devnet`), so a Base run cannot
disturb an OP or Arbitrum devnet. Pass `ENCLAVE=op-devnet` to settle it on the same
L1 as those instead.

**Base and Optimism are the same chain type to the relayer.** `src_chain` /
`dst_chain` stay `opstack` — `cmd/main.go` accepts only `cosmos | ethereum | opstack
| arbitrum`. The two are told apart by module **name** and client id, which is why
the deploy step below needs `MODULE_NAME`.

```bash
# 1. L1 (Fulu) + Base L2. Brings the L1 up via run_eth_node.sh if the enclave does
#    not exist. Ends with .base-devnet-run/attestor.env written.
./scripts/local/run_base_node.sh

# 2. Attestor. run_base_attestor.sh is a symlink to run_op_attestor.sh; it picks
#    .base-devnet-run/attestor.env from its own name, so do NOT call
#    run_op_attestor.sh here — that one attaches to the OP devnet.
GRPC_PORT=3003 ./scripts/local/run_base_attestor.sh

# 3. Cosmos node, then gov-store the Base light-client wasm. The Ethereum client
#    is NOT needed for a Cosmos<->Base deployment — like OP and Arbitrum, the
#    attestor-trusted L2 client authenticates nothing through it. Add
#    ./scripts/local/wasm.sh only if this config also relays Cosmos<->Ethereum.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm_base.sh     # -> Base client checksum

# 4. Deploy the L2 IBC contracts onto Base. MODULE_NAME is required: the module's
#    dst_chain is "opstack" (see above), so matching on it alone would also match
#    the Optimism module. DST_CHAIN only picks which devnet handoff to read. The
#    deploy txs go to Base's sequencer RPC, while the relayer config is patched to
#    Base's follower RPC because L2->Cosmos membership proofs need eth_getProof
#    below the head.
L2_ENV_FILE=.base-devnet-run/attestor.env MODULE_NAME=cosmos-to-base \
DST_CHAIN=base ./scripts/local/deploy_l2_contracts.sh

# 5. Create the clients. Copy relayer/base-l2-config.example.json and fill in the
#    checksum from step 3, the L2 router from step 4, and the L2 chain id. There is
#    only ONE checksum: the file carries no ethereum_client since #345/#347, because
#    the client verifies nothing against L1.
cd relayer && go build -o relayer ./cmd
./relayer create-clients-cosmos --config config.json --l2-config base-l2-config.json
./relayer create-clients-eth --config config.json \
    --source base-client-0 --trust-level 2/3

# 6. Relay.
./relayer start --config config.json
```

Everything in the "Local Cosmos ↔ OP E2E" section about the ICS26Router relayer
role and `--absolute-timeouts` on the test transfer applies unchanged — Base uses
the same attestor and the same L2 contracts.

### Running against an existing Base L2 (no devnet L2)

Skip step 1 and point the attestor at the node someone else runs. Verified end to
end against **Base Sepolia**, relaying both directions.

```bash
# 1. Attestor against the existing node. NETWORK must name the real chain, and the
#    dispute-game contracts cannot be auto-resolved for a public network — read
#    them off the chain once (below) and pass them.
OP_NODE_RPC_URL=http://<host>:7545 \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
NETWORK=base-sepolia SRC_CHAIN=base-sepolia \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 \
DISPUTE_GAME_FACTORY=0xd6E6dBf4F7EA0ac412fD8b65ED297e64BB7a06E1 \
RESPECTED_GAME_TYPE=621 \
GRPC_PORT=3003 \
  ./scripts/local/run_base_attestor.sh

# 2. Cosmos + the Base wasm, unchanged from the devnet flow.
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_base.sh

# 3. config.json: one runnable config per path, already cut to the Base pair.
cp relayer/config.base.example.json relayer/config.json
# then fill in the fields listed in "Values you fill in by hand" below.

# 4. Deploy the L2 contracts with a key funded ON Base Sepolia. MODULE_NAME is
#    still required — dst_chain is "opstack" and would also match cosmos-to-op.
#    DO THIS BEFORE STARTING THE ATTESTOR if both share one rate-limited RPC —
#    see "One RPC budget, two processes" in the Arbitrum section.
DST_CHAIN=opstack MODULE_NAME=cosmos-to-base \
L2_RPC=http://<host>:8545 RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh

# 5. l2-config, then clients. The relayer must sign with the SAME key that
#    deployed (E2ETestDeployL2 grants the router role to msg.sender). Pass it in
#    the environment rather than editing relayer/.env — godotenv does not override
#    a variable already set, so the export wins.
cp relayer/base-l2-config.example.json relayer/base-l2-config.json
# fill in wasm_checksum (step 2), l2_rpc_url, and l2_router (= ICS26_ADDRESS from step 4)
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config base-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --source base-client-0 --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

#### Values you fill in by hand, and where each comes from

Everything not listed here is written for you — see
[What the tooling fills in for you](#what-the-tooling-fills-in-for-you).

**`relayer/config.json`:**

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-base` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-base` | your L2 WebSocket endpoint; optional, the relayer skips the WS client when empty |
| `l2_rpc_url` | `base-to-cosmos` | the same L2 HTTP endpoint (step 4 overwrites it with `L2_RPC` anyway) |
| `attestor_addr` | `base-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1. The example ships `127.0.0.1:3003`, matching the `GRPC_PORT=3003` above — Base needs its own port so it can coexist with an OP attestor on 3001 |
| `attestor_src_chain` | `base-to-cosmos` | the `SRC_CHAIN` you passed in step 1 (`base-sepolia`) |
| `head_kind` | `base-to-cosmos` | the same choice as `ATTESTATION_HEAD` in step 1 |
| `rollup_profile.common.l2_chain_id` | `base-to-cosmos` | `eth_chainId` on your L2 (`84532` on Base Sepolia) |
| `log_scan_chunk` | `base-to-cosmos` | your provider's `eth_getLogs` span cap — measure it, see the Arbitrum section |

**`relayer/base-l2-config.json`:** `wasm_checksum` (from `wasm_base.sh`), `l2_rpc_url`,
and `rollup_profile.common.l2_router` (= `ICS26_ADDRESS` printed by step 4).

**Fund the deployer on Base itself.** An account funded on L1 Sepolia, OP Sepolia or
Arbitrum Sepolia has nothing here — balances do not carry across rollups. The deploy
fails at the first transaction with `insufficient funds for gas * price + value: have
0`, after the script has already written `broadcast/`. Check first:
`cast balance <addr> --rpc-url $L2 --ether`.

**Resolving the dispute-game contracts.** The values above are Base Sepolia's, read
from the chain rather than hardcoded anywhere — the portal address comes from the
node itself:

```bash
PORTAL=$(curl -s -X POST -H 'content-type: application/json' \
  --data '{"jsonrpc":"2.0","id":1,"method":"optimism_rollupConfig","params":[]}' \
  http://<host>:7545 | jq -r .result.deposit_contract_address)
cast call $PORTAL "disputeGameFactory()(address)" --rpc-url $L1_RPC_URL
cast call $PORTAL "respectedGameType()(uint32)"  --rpc-url $L1_RPC_URL
```

**Check how far below the head the node will serve `eth_getProof`.** The relayer
proves the router account at the attested height, which trails the head. A node that
prunes state does not error — it can simply stop responding, which looks like a hung
relayer rather than a configuration problem. One measured node served head-500 in
about a second and returned nothing at all at head-800 after 30 s. With
`DERIVED_GAP_BLOCKS=10` the attested height trails by ~10–20 blocks, comfortably
inside that; `head_kind=finalized` trails by ~600 and would not be.

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getProof '["<router>",[],"'$(printf '0x%x' $((H-500)))'"]' --rpc-url $L2
```

Useful:

```bash
./scripts/local/run_base_node.sh --status   # are the Base services up
./scripts/local/run_base_node.sh --reset    # redeploy Base only; L1 preserved
./scripts/local/run_base_node.sh --stop     # stop Base, leave the L1 running
kurtosis enclave rm -f base-devnet          # remove the L1 and everything on it
```
