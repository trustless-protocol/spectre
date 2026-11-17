# Testnet end-to-end runbooks

Relaying between a local Cosmos chain and an L2 that someone else operates: OP Sepolia,
Arbitrum Sepolia or Base Sepolia. The Cosmos chain is still the local Gaia node, and the
L2 contracts are still deployed with the test script `E2ETestDeployL2`. For a production
deployment see [PRODUCTION_DEPLOYMENT.md](PRODUCTION_DEPLOYMENT.md).

Read [E2E_LOCAL.md](E2E_LOCAL.md) first, and complete steps 1 and 2 of its Ethereum flow (prover
artifacts and relayer binary). Each flow below is complete in itself: it runs the attestor against
an existing L2 where the local flow starts its own. These sections of the local runbook apply
unchanged:

- [Gaia binary](E2E_LOCAL.md#gaia-binary)
- [Prover key and verifier must be a pair](E2E_LOCAL.md#prover-key-and-verifier-must-be-a-pair)
- [Signing keys when several processes share a chain](E2E_LOCAL.md#signing-keys-when-several-processes-share-a-chain)
- [Spending caps](E2E_LOCAL.md#spending-caps): `start` needs all four; set Sepolia-sized values
- [Sending a test packet](E2E_LOCAL.md#sending-a-test-packet)
- [Failure modes](E2E_LOCAL.md#failure-modes)

## Table of contents

- [OP Sepolia](#op-sepolia)
  - [Values to fill in by hand (OP)](#values-to-fill-in-by-hand-op)
- [Arbitrum Sepolia](#arbitrum-sepolia)
  - [Values to fill in by hand (Arbitrum)](#values-to-fill-in-by-hand-arbitrum)
- [Base Sepolia](#base-sepolia)
  - [Values to fill in by hand (Base)](#values-to-fill-in-by-hand-base)
- [Values the tooling fills in](#values-the-tooling-fills-in)
- [One RPC budget, two processes](#one-rpc-budget-two-processes)

## OP Sepolia

The [local OP flow](E2E_LOCAL.md#local-cosmos--op-e2e) brings up its own L2. To relay against a
node someone else operates (OP Sepolia, or a replica you already run), skip its step 1 and
point the attestor at that node. There is no local L1 to start either: the attestor-trusted client verifies nothing
against L1, so only the *attestor* needs an L1 RPC, and a public endpoint is enough.

An earlier version of this flow ran against OP Sepolia. The current flow, which adds signed
attestations, has not yet been run end to end.

**Step 1.** Attestor against the existing replica. NETWORK must name the real chain, and the
dispute-game contracts cannot be auto-resolved for a public network: read them off the chain
once (see below) and pass them.

```bash
OP_NODE_RPC_URL=http://<host>:9545 \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
NETWORK=op-sepolia SRC_CHAIN=op-sepolia \
L2_CHAIN_ID=11155420 \
ATTESTOR_SIGNING_KEY=<32-byte-ed25519-seed-hex> \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 \
DISPUTE_GAME_FACTORY=0x05F9613aDB30026FFd634f38e5C4dFd30a197Fa1 \
RESPECTED_GAME_TYPE=8 \
  ./scripts/local/run_op_attestor.sh
```

**Step 2.** Cosmos + the OP wasm, as in the local flow.

```bash
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_op.sh
```

**Step 3.** `config.json`: one runnable config per path, already cut to the OP module pair. The
deploy below patches this file, so it has to exist first.

```bash
cp relayer/config.op.example.json relayer/config.json
```

Then fill in the fields listed in "Values to fill in by hand (OP)" below.

**Step 4.** Deploy the L2 contracts with a key funded ON THAT L2. If the attestor and this deploy share one rate-limited RPC, run this step before step 1: see "One RPC budget, two processes" below.

```bash
DST_CHAIN=opstack L2_RPC=http://<host>:8545 RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh
```

**Step 5.** The `--l2-config` file, then clients and relay. The relayer must sign with the SAME key that
deployed (E2ETestDeployL2 grants the router role to msg.sender). Pass it in the environment
rather than editing `relayer/.env`: godotenv does not override a variable already set, so the
export wins.

```bash
cp relayer/op-l2-config.example.json relayer/op-l2-config.json
```

Fill in `wasm_checksum` (step 2), `l2_rpc_url`, `l2_router` (= `ICS26_ADDRESS` from step 4), and
the attestor public key as described below.

```bash
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config op-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

### Values to fill in by hand (OP)

Everything not listed here is written for you. See
[Values the tooling fills in](#values-the-tooling-fills-in).

The attestor logs its public key twice at startup: hex for human comparison and standard base64
for the relayer. Copy the base64 value from
`.op-attestor-run/attestor.log` into both locations below before creating the client. They must
contain the public half of the exact `ATTESTOR_SIGNING_KEY` supplied in step 1; leaving the
checked-in disposable devnet key in either file lets the daemon start but makes every certificate
fail against the client or relayer configuration.

**`relayer/config.json`:**

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-op` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-op` | your L2 WebSocket endpoint; optional, the relayer skips the WS client when empty |
| `l2_rpc_url` | `op-to-cosmos` | the same L2 HTTP endpoint (step 4 overwrites it with `L2_RPC` anyway) |
| `attestor_endpoints[0].address` | `op-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1; the example ships `127.0.0.1:3001`, which matches the script default. The list takes one `{address, attestor_index}` entry per attestor daemon, and the config is rejected when it has fewer entries than `attestors.threshold` |
| `attestor_endpoints[0].attestor_index` | `op-to-cosmos` | the position of that attestor's key in `attestors.public_keys` (`0` for the single local attestor). Each index may appear once and must lie inside the key list |
| `attestor_src_chain` | `op-to-cosmos` | the `SRC_CHAIN` you passed in step 1 (`op-sepolia`): must match, or `AttestedUpTo` answers for another chain |
| `head_kind` | `op-to-cosmos` | the same choice as `ATTESTATION_HEAD` in step 1. `start` asks each attestor's `Info` for the head it gates on and refuses to boot when the two differ: the attestor signs only its configured head, so a mismatch relays nothing |
| `rollup_profile.common.l2_chain_id` | `op-to-cosmos` | `eth_chainId` on your L2 (`11155420` on OP Sepolia). `start` verifies this against the RPC and refuses to boot on a mismatch |
| `attestors.public_keys` | `op-to-cosmos` | the base64 public key logged by the attestor started with `ATTESTOR_SIGNING_KEY` |
| `log_scan_chunk` | `op-to-cosmos` | your provider's `eth_getLogs` span cap: measure it, see [Arbitrum Sepolia](#arbitrum-sepolia) |
| `router_deploy_block` | `op-to-cosmos` | optional: block the `l2_router` proxy was deployed at; see [`router_deploy_block`](E2E_LOCAL.md#router_deploy_block-where-packet-flush-log-lookups-start) |

**`relayer/op-l2-config.json`:** `wasm_checksum` (from `wasm_op.sh`), `l2_rpc_url`, and
`rollup_profile.common.l2_router` (= `ICS26_ADDRESS` printed by step 4, which is why
this file is filled in after the deploy), plus the same base64 key in
`attestors.public_keys`.

**Resolving the dispute-game contracts.** `run_op_attestor.sh` can derive them with
`cast` for a chain it knows; for a public network it exits with
`set DISPUTE_GAME_FACTORY and RESPECTED_GAME_TYPE`. Read them from the chain: the
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
proof window is small: a node serving only `latest` fails every build with
`distance to target block exceeds maximum proof window`. When `run_op_attestor.sh`
starts its own replica (given `OP_RETH_BIN` and `OP_NODE_BIN`), it runs op-reth with
`--rpc.eth-proof-window=100000`; an external
replica needs the operator to have done the same. Check before starting:

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getProof '["<router>",[],"'$(printf '0x%x' $((H-1000)))'"]' --rpc-url $L2
```

**Gaia must be the custom-host build.** `gaiad_binary.sh` requires the checkout to be
*on* branch `test/ibc-host-customs`, which a detached HEAD fails even when the code is
right. Point at an already-built binary instead: `GAIAD=` takes priority and skips
the branch check. Confirm it is the right one:

```bash
$GAIAD version                       # → test/ibc-host-customs-<sha>
strings $GAIAD | grep -c ClientStatus  # the 08-wasm allowlist entry L2 clients need
```

**Set `DERIVED_GAP_BLOCKS` explicitly.** The local devnet handoffs export `5`; an
attestor attached to an external replica inherits nothing and takes the default 150,
which puts a ~5 minute floor under the return direction. See [return-direction
latency](E2E_LOCAL.md#return-direction-latency-is-set-by-the-attestor-not-the-prover).

## Arbitrum Sepolia

The Arbitrum counterpart of the OP section above, replacing step 1 of the [local Arbitrum
flow](E2E_LOCAL.md#local-cosmos--arbitrum-e2e). There is no local L1 and no local
Nitro: the attestor reads a public Arbitrum RPC and a public Sepolia RPC, and the
client verifies nothing against L1, so nothing else needs an L1 endpoint.

An earlier version of this flow ran against Arbitrum Sepolia. The current flow, which adds signed
attestations, has not yet been run end to end.

**Step 1.** Attestor. `CHAIN_PROFILE` supplies the RollupCore address, both chain ids, the BoLD
slot/offset and a scan start block near the head, all from
`attestor/arbitrum/config.arbitrum-sepolia.json`: pass only the endpoints.

```bash
CHAIN_PROFILE=arbitrum-sepolia \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
L2_RPC_URL=<arbitrum-sepolia-rpc> \
L2_WS_URL=<arbitrum-sepolia-ws> \
ATTESTOR_SIGNING_KEY=<32-byte-ed25519-seed-hex> \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 GRPC_PORT=3002 DETACH=1 \
  ./scripts/local/run_arbitrum_attestor.sh
```

**Step 2.** Cosmos + the Arbitrum wasm, as in the local flow.

```bash
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_arb.sh          # -> checksum
```

**Step 3.** `config.json`: one runnable config per path, already cut to the Arbitrum pair.

```bash
cp relayer/config.arbitrum.example.json relayer/config.json
```

Then fill in the fields listed in "Values to fill in by hand (Arbitrum)" below.

**Step 4.** Deploy the L2 contracts with a key funded ON Arbitrum Sepolia. If the attestor and this deploy share one rate-limited RPC, run this step before step 1: see "One RPC budget, two processes"
below.

```bash
DST_CHAIN=arbitrum L2_RPC=<arbitrum-sepolia-rpc> RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh
```

Patches BOTH modules: `ics26_address` on the forward one, `rollup_profile.common.l2_router` and
`l2_rpc_url` on the return one.

**Step 5.** The `--l2-config` file, then clients. The relayer must sign with the SAME key that deployed
(E2ETestDeployL2 grants the router role to msg.sender). Pass it in the environment rather than
editing `relayer/.env`: godotenv does not override a variable already set, so the export wins.

```bash
cp relayer/arb-l2-config.example.json relayer/arb-l2-config.json
```

Fill in `wasm_checksum` (step 2), `l2_rpc_url`, `l2_router` (= `ICS26_ADDRESS` from step 4), and
the attestor public key as described below.

```bash
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config arb-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

**`l2_wasm_client_id` is an output, not an input.** `create-clients-cosmos` creates that
client and writes its id back into the module, overwriting whatever is there: the
example's `08-wasm-N` placeholder or an empty string alike. Every other command still
requires it, because by then the client must exist.

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

**With `CHAIN_PROFILE` set to anything but `devnet`, a leftover `.arbitrum-devnet-run/`
is ignored.** The devnet handoff exports `ROLLUP_CORE_ADDRESS` and both chain ids as real
environment variables, which would otherwise beat the profile and silently produce
`l2_chain_id=412346` and `src_chain=arbdev` against a public RPC. The script says so when
it skips the handoff.

### Values to fill in by hand (Arbitrum)

Everything else in both files is written for you. See
[Values the tooling fills in](#values-the-tooling-fills-in).

The attestor logs its public key twice at startup: hex for human comparison and standard base64
for the relayer. Copy the base64 value from
`.arbitrum-attestor-run/attestor.log` into both locations below before creating the client. They
must contain the public half of the exact `ATTESTOR_SIGNING_KEY` supplied in step 1; leaving the
checked-in disposable devnet key in either file lets the daemon start but makes every certificate
fail against the client or relayer configuration.

**`relayer/config.json`:**

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-arbitrum` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-arbitrum` | your L2 WebSocket endpoint. Optional: the relayer skips the WS client when it is empty, and the L2→Cosmos direction polls regardless |
| `l2_rpc_url` | `arbitrum-to-cosmos` | the same L2 HTTP endpoint. Step 4 overwrites this with `L2_RPC`, so a mismatch here is corrected rather than fatal |
| `attestor_endpoints[0].address` | `arbitrum-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1. The example ships `127.0.0.1:3002`, which matches the `GRPC_PORT=3002` above: **change it if you used a different port**, or the relayer dials a closed socket. The list takes one `{address, attestor_index}` entry per attestor daemon, and the config is rejected when it has fewer entries than `attestors.threshold` |
| `attestor_endpoints[0].attestor_index` | `arbitrum-to-cosmos` | the position of that attestor's key in `attestors.public_keys` (`0` for the single local attestor). Each index may appear once and must lie inside the key list |
| `attestor_src_chain` | `arbitrum-to-cosmos` | `src_chain` in `attestor/arbitrum/config.<profile>.json` (`arbitrum-sepolia`). Must equal what the attestor reports, or `AttestedUpTo` answers for a chain you did not ask about |
| `head_kind` | `arbitrum-to-cosmos` | `unsafe`, `safe` or `finalized`: the same choice as `ATTESTATION_HEAD` in step 1. `start` asks each attestor's `Info` for the head it gates on and refuses to boot when the two differ: the attestor signs only its configured head, so a mismatch relays nothing |
| `rollup_profile.common.l2_chain_id` | `arbitrum-to-cosmos` | `eth_chainId` on your L2 (`421614` on Arbitrum Sepolia, `412346` on the local devnet). `start` verifies this against the RPC and refuses to boot on a mismatch |
| `attestors.public_keys` | `arbitrum-to-cosmos` | the base64 public key logged by the attestor started with `ATTESTOR_SIGNING_KEY` |
| `log_scan_chunk` | `arbitrum-to-cosmos` | your provider's `eth_getLogs` span cap: measure it, see the `log_scan_chunk` note above |
| `router_deploy_block` | `arbitrum-to-cosmos` | optional: block the `l2_router` proxy was deployed at; see [`router_deploy_block`](E2E_LOCAL.md#router_deploy_block-where-packet-flush-log-lookups-start) |

**`relayer/arb-l2-config.json`:**

| Field | Value comes from |
|---|---|
| `wasm_checksum` | the hex printed by `wasm_arb.sh` in step 2 |
| `l2_rpc_url` | your L2 HTTP endpoint |
| `rollup_profile.common.l2_router` | `ICS26_ADDRESS` printed by step 4, so this file is filled in *after* the deploy, which is why step 5 comes after step 4 |
| `attestors.public_keys` | the same base64 public key logged by the attestor |

`counterparty_client_id` and `l2_chain_id` in that file already match the example config
(`arb-client-0`, `421614`); change them only if you changed the module's
`ics26_client_id` or are on a different chain.

## Base Sepolia

Skip step 1 of the [local Base flow](E2E_LOCAL.md#local-cosmos--base-e2e) and point the attestor
at the node someone else runs. An earlier version of this flow ran against Base Sepolia. The
current flow, which adds signed
attestations, has not yet been run end to end.

**Step 1.** Attestor against the existing node. NETWORK must name the real chain, and the
dispute-game contracts cannot be auto-resolved for a public network: read them off the chain
once (below) and pass them.

```bash
OP_NODE_RPC_URL=http://<host>:7545 \
L1_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
NETWORK=base-sepolia SRC_CHAIN=base-sepolia \
L2_CHAIN_ID=84532 \
ATTESTOR_SIGNING_KEY=<32-byte-ed25519-seed-hex> \
ATTESTOR_RUN_DIR=.base-attestor-run \
ATTESTATION_HEAD=unsafe DERIVED_GAP_BLOCKS=10 \
DISPUTE_GAME_FACTORY=0xd6E6dBf4F7EA0ac412fD8b65ED297e64BB7a06E1 \
RESPECTED_GAME_TYPE=621 \
GRPC_PORT=3003 \
  ./scripts/local/run_base_attestor.sh
```

**Step 2.** Cosmos + the Base wasm, as in the local flow.

```bash
GAIAD=~/go/bin/gaiad ./scripts/local/run_cosmos_node.sh
GAIAD=~/go/bin/gaiad ./scripts/local/wasm_base.sh
```

**Step 3.** `config.json`: one runnable config per path, already cut to the Base pair.

```bash
cp relayer/config.base.example.json relayer/config.json
```

Then fill in the fields listed in "Values to fill in by hand (Base)" below.

**Step 4.** Deploy the L2 contracts with a key funded ON Base Sepolia. `MODULE_NAME` is still
required: `dst_chain` is "opstack" and would also match cosmos-to-op. With `L2_RPC` and
`MODULE_NAME` set, `DST_CHAIN` selects nothing else, so `opstack` and `base` behave the same here. If the attestor and this deploy share one rate-limited RPC, run this step before step 1: see "One RPC budget, two processes" below.

```bash
DST_CHAIN=opstack MODULE_NAME=cosmos-to-base \
L2_RPC=http://<host>:8545 RELAYER_CONFIG=relayer/config.json \
L2_DEPLOYER_ADDRESS=0x... L2_DEPLOYER_PRIVATE_KEY=... \
  ./scripts/local/deploy_l2_contracts.sh
```

**Step 5.** The `--l2-config` file, then clients. The relayer must sign with the SAME key that deployed
(E2ETestDeployL2 grants the router role to msg.sender). Pass it in the environment rather than
editing `relayer/.env`: godotenv does not override a variable already set, so the export wins.

```bash
cp relayer/base-l2-config.example.json relayer/base-l2-config.json
```

Fill in `wasm_checksum` (step 2), `l2_rpc_url`, `l2_router` (= `ICS26_ADDRESS` from step 4), and
the attestor public key as described below.

```bash
cd relayer
./relayer create-clients-cosmos --config config.json --l2-config base-l2-config.json
ETH_PRIVATE_KEY=<deployer-key> ./relayer create-clients-eth --config config.json \
  --trust-level 2/3
ETH_PRIVATE_KEY=<deployer-key> ./relayer start --config config.json
```

### Values to fill in by hand (Base)

Everything not listed here is written for you. See
[Values the tooling fills in](#values-the-tooling-fills-in).

Base uses the OP Stack attestor, which logs its public key twice at startup: hex for human
comparison and standard base64 for the relayer. Copy the base64 value from
`.base-attestor-run/attestor.log` into both locations below before creating the client. They must
contain the public half of the exact `ATTESTOR_SIGNING_KEY` supplied in step 1; leaving the
checked-in disposable devnet key in either file lets the daemon start but makes every certificate
fail against the client or relayer configuration.

**`relayer/config.json`:**

| Field | Module | Value comes from |
|---|---|---|
| `eth_rpc_url` | `cosmos-to-base` | your L2 HTTP endpoint |
| `eth_ws_url` | `cosmos-to-base` | your L2 WebSocket endpoint; optional, the relayer skips the WS client when empty |
| `l2_rpc_url` | `base-to-cosmos` | the same L2 HTTP endpoint (step 4 overwrites it with `L2_RPC` anyway) |
| `attestor_endpoints[0].address` | `base-to-cosmos` | `127.0.0.1:$GRPC_PORT` from step 1. The example ships `127.0.0.1:3003`, matching the `GRPC_PORT=3003` above: Base needs its own port so it can coexist with an OP attestor on 3001. The list takes one `{address, attestor_index}` entry per attestor daemon, and the config is rejected when it has fewer entries than `attestors.threshold` |
| `attestor_endpoints[0].attestor_index` | `base-to-cosmos` | the position of that attestor's key in `attestors.public_keys` (`0` for the single local attestor). Each index may appear once and must lie inside the key list |
| `attestor_src_chain` | `base-to-cosmos` | the `SRC_CHAIN` you passed in step 1 (`base-sepolia`) |
| `head_kind` | `base-to-cosmos` | the same choice as `ATTESTATION_HEAD` in step 1. `start` asks each attestor's `Info` for the head it gates on and refuses to boot when the two differ: the attestor signs only its configured head, so a mismatch relays nothing |
| `rollup_profile.common.l2_chain_id` | `base-to-cosmos` | `eth_chainId` on your L2 (`84532` on Base Sepolia). `start` verifies this against the RPC and refuses to boot on a mismatch |
| `attestors.public_keys` | `base-to-cosmos` | the base64 public key logged by the attestor started with `ATTESTOR_SIGNING_KEY` |
| `log_scan_chunk` | `base-to-cosmos` | your provider's `eth_getLogs` span cap: measure it, see [Arbitrum Sepolia](#arbitrum-sepolia) |
| `router_deploy_block` | `base-to-cosmos` | optional: block the `l2_router` proxy was deployed at; see [`router_deploy_block`](E2E_LOCAL.md#router_deploy_block-where-packet-flush-log-lookups-start) |

**`relayer/base-l2-config.json`:** `wasm_checksum` (from `wasm_base.sh`), `l2_rpc_url`,
and `rollup_profile.common.l2_router` (= `ICS26_ADDRESS` printed by step 4), plus the same
base64 key in `attestors.public_keys`.

**Fund the deployer on Base itself.** An account funded on L1 Sepolia, OP Sepolia or
Arbitrum Sepolia has nothing here: balances do not carry across rollups. The deploy
fails at the first transaction with `insufficient funds for gas * price + value: have
0`, after the script has already written `broadcast/`. Check first:
`cast balance <addr> --rpc-url $L2 --ether`.

**Resolving the dispute-game contracts.** The values above are Base Sepolia's, read
from the chain rather than hardcoded anywhere: the portal address comes from the
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
prunes state does not error: it can simply stop responding, which looks like a hung
relayer rather than a configuration problem. One measured node served head-500 in
about a second and returned nothing at all at head-800 after 30 s. With
`DERIVED_GAP_BLOCKS=10` the attested height trails by ~10–20 blocks, comfortably
inside that; `head_kind=finalized` trails by ~600 and would not be.

```bash
H=$(cast block-number --rpc-url $L2)
cast rpc eth_getProof '["<router>",[],"'$(printf '0x%x' $((H-500)))'"]' --rpc-url $L2
```

## Values the tooling fills in

Identical for OP, Arbitrum and Base. Do not pre-fill these: the commands overwrite
whatever is there:

| Field | Written by |
|---|---|
| `ics26_address`, `membership`, `update_client`, `signature_verifier`, `misbehaviour` | `deploy_l2_contracts.sh` (forward module) |
| `rollup_profile.common.l2_router`, `l2_rpc_url` | `deploy_l2_contracts.sh` (return module: it prints `Patched the return module "…"`) |
| `cosmos_wasm_client_id`, `l2_wasm_client_id` | `create-clients-cosmos`, which overwrites the example's `08-wasm-N` placeholder |
| `spectre_client` | `create-clients-eth` |

## One RPC budget, two processes

Applies to every L2 in this document. The numbers below were measured on Arbitrum
Sepolia.

**The attestor and the relayer both load the L2 endpoint, and on a free tier they do
not fit together.** The attestor issues `runtime_backfill_concurrency` (8) parallel
header fetches every refresh; the relayer adds `eth_getLogs` scans, `eth_getProof`
builds and transaction sends. Against one Alchemy free-tier key this measurably
exceeds the concurrent-request budget:

- `deploy_l2_contracts.sh` fails outright with `Max retries exceeded HTTP error 429
  with empty body` when the attestor is already running. Pausing the attestor
  (`docker stop spectre-arbitrum-attestor`) makes the same command succeed.
- With both running, the attestor logs `backing off: failures=3 delay=40.5s
  rate_limited=true` and its frontier stops. The relayer then reports `source relayable
  height` frozen for minutes with `DERIVED_GAP_BLOCKS=10` set, which looks like a
  relayer problem and is not one.

Two ways out, in order of preference:

**Give each process its own endpoint.** The relayer does not need a keyed provider.
Measured on Arbitrum Sepolia, both of these serve what it requires with no API key:

| Endpoint | `eth_getProof` at head−500 | `eth_getLogs` span | WebSocket |
|---|---|---|---|
| `https://sepolia-rollup.arbitrum.io/rpc` | yes | ≥10 000 | none |
| `https://arbitrum-sepolia.drpc.org` | yes | ≥10 000 | `wss://arbitrum-sepolia.drpc.org` |
| `https://arbitrum-sepolia-rpc.publicnode.com` | **no**: fails below ~head−128 | ≥10 000 | `wss://…publicnode.com` |

So point `L2_RPC_URL` (attestor) at the keyed endpoint and the relayer's `eth_rpc_url` /
`l2_rpc_url` at a keyless one, and neither starves the other. Do not use publicnode for
the relayer: it cannot prove the router account at the attested height.

**Or serialise the load.** Deploy the contracts first, then start the attestor, and
lower `RUNTIME_BACKFILL_CONCURRENCY` if the attestor still trips the limit.
