# fast-ibc

A Solidity implementation of IBC Eureka (IBC v2) with a Go relayer using gnark
Groth16 for Tendermint light client verification. Each client update proves a
2/3+ voting-power quorum of validator Ed25519 signatures in a single Groth16
proof, with in-circuit CanonicalVote reconstruction so only ~32 bytes of public
input land on-chain.

## L2 ICS-08 clients

The repository also contains independent CosmWasm clients for OP Sepolia, Base
Sepolia, and Arbitrum Sepolia. Their deployment evidence requirements,
host-query wiring, and verification commands are documented in
[docs/L2_CLIENTS.md](docs/L2_CLIENTS.md).

## Architecture

```
┌─────────────────┐         ┌──────────────────────────────────────┐
│   Cosmos Chain  │         │           Ethereum Chain             │
│                 │         │                                      │
│  IBC v2 Module  │────────▶│  ICS26Router.sol                     │
│  (send_packet)  │         │      │                               │
└─────────────────┘         │      ▼                               │
                            │  SpectreClient.sol                   │
         ┌──────────────────│      │  (light client + 2/3 quorum)  │
         │                  │      ▼                               │
         │  Go Relayer      │  SignatureVerifier.sol               │
         │  ┌────────────┐  │      │  (rebuilds CanonicalVote +    │
         │  │ Extractor  │  │      │   SHA-256 witness commit)     │
         │  │ + Prover   │  │      ▼                               │
         │  │ (gnark ZK) │  │  Groth16Verifier_N{N}.sol            │
         │  │            │  │      │  (one per bucket size)        │
         │  │ RecvPacket │  │      ▼                               │
         │  └────────────┘  │  Membership.sol → ICS20Transfer.sol  │
         └──────────────────└──────────────────────────────────────┘
```

### Multi-validator batching

The Tendermint commit is verified by batch-proving N Ed25519 signatures whose
voting power sums to ≥ 2/3 of the validator set. To keep Groth16 circuits
fixed-size, the prover picks the smallest **bucket** (N ∈ {4, 8, 16, 32, 64,
128}) that fits the required signers and pads the rest with deterministic
dummy keypairs. Each bucket has its own `(r1cs, pk, vk)` artifacts and a
matching `Groth16Verifier_N{N}.sol`; `SignatureVerifier` dispatches by bucket.

The circuit reconstructs each validator's `CanonicalVote` bytes from a shared
block header + per-slot `Timestamp`, hashes the full witness (active flag,
pubkey, msg) into a single SHA-256 digest, and exposes that digest as the only
public input. This keeps the on-chain verifier well under EIP-170. Padding
slots carry `active=false`; both the in-circuit hash and the on-chain quorum
check skip them.

The current Tendermint light-client cache supports at most **180 active
validators**. Any update whose current validator set has more than 180 active
validators reverts with `ValidatorCountExceedsLimit(count, 180)`, so chains
above that bound need a larger cache layout before they can use this client.

## Requirements

- [Go](https://golang.org/) >= 1.21
- [Foundry](https://getfoundry.sh/)
- [Bun](https://bun.sh/)
- [Just](https://github.com/casey/just)
- Optional for GPU proving: ICICLE runtime/libs installed on the host, plus an
  `icicle` build of the relayer/prover tool

## gnark submodules (required to build the relayer)

The relayer/prover depend on two forks of gnark, vendored as **git submodules**
under `third_party/` and wired through `relayer/go.mod`:

```
replace (
    0x5ea000000/ecip-gnark      => ../third_party/ecip-gnark
    github.com/consensys/gnark  => ../third_party/decentrio-gnark
)
```

- `third_party/ecip-gnark` → [`decentrio/ecip-gnark`](https://github.com/decentrio/ecip-gnark) — Ed25519 in-circuit ops + the `garaga_rs` Rust FFI.
- `third_party/decentrio-gnark` → [`decentrio/gnark`](https://github.com/decentrio/gnark) — fork of consensys/gnark v0.13.0 with hash-aggregate verifier tweaks.

fast-ibc stores only a **pinned commit** of each fork (a submodule pointer), not
their files.

### First checkout

Clone with submodules, or initialise them in an existing clone:

```bash
git clone --recurse-submodules https://github.com/decentrio/fast-ibc
# or, after a plain clone:
git submodule update --init --recursive
```

`ecip-gnark`'s prover links a Rust FFI (`libgaraga_rs`) that cgo expects at the
ecip-gnark root. Build it once (and after any garaga change):

```bash
cd third_party/ecip-gnark/ffi/garaga_rs
cargo build --release
cp target/release/libgaraga_rs.so ../..      # Linux; on macOS: libgaraga_rs.dylib
```

Without `submodule update --init`, `go build ./...` under `relayer/` fails with
`replacement directory ../third_party/ecip-gnark does not exist`; without the
built FFI it fails at link with `library 'garaga_rs' not found`.

> **CI:** the fork repos are **private**, so the `Go` and `E2E (manual)`
> workflows check the submodules out with a `SUBMODULE_TOKEN` secret (a
> fine-grained PAT with read access to `decentrio/gnark` + `decentrio/ecip-gnark`)
> and build `libgaraga_rs` before `go build`. See `.github/workflows/go.yml`.

### Updating the forks (bump workflow)

Develop in the fork repos themselves and push there as usual — fast-ibc only
pins a commit, so do **not** edit inside `third_party/` and forget to push the
fork (CI fetches by SHA and an unpushed commit fails checkout). When a fork has
changes fast-ibc should pick up, bump the pinned commit:

```bash
# pull the tracked branch's latest into the submodule…
git submodule update --remote third_party/ecip-gnark
#   …or pin an exact commit/tag:
#   cd third_party/ecip-gnark && git fetch && git checkout <sha-or-tag> && cd -

git add third_party/ecip-gnark      # stage the new pointer
git commit -m "chore: bump ecip-gnark to <sha>"
git push                            # the PR diff is just a one-line submodule pointer
```

- fast-ibc builds against the **pinned** commit, not your in-progress fork work, until you push + bump. For a tight fork↔fast-ibc co-development loop, temporarily point the `go.mod` replace at a local clone (e.g. `=> ../../ecip-gnark`) and do **not** commit that change; revert + bump when stable.
- After someone else bumps, run `git submodule update --init --recursive` to sync your local tree.

## Optional GPU Proving

CPU proving remains the default. GPU proving is opt-in and follows the
`test/gnark-gpu` approach: build with `-tags=icicle`, then enable the ICICLE
backend via env or flag when needed.

Build requirements for the GPU path:

- NVIDIA GPU with a working CUDA driver/toolkit
- ICICLE runtime libraries must be installed and visible to the linker/runtime
- the relayer and prover tool must be built or run with `-tags=icicle`

The linker/runtime must be able to find libraries such as:

- `libicicle_device`
- `libicicle_field_bn254`
- `libicicle_curve_bn254`

If they are not in a default loader path, export `LD_LIBRARY_PATH` before
building or running:

```bash
export LD_LIBRARY_PATH=/usr/local/lib:${LD_LIBRARY_PATH}
```

Backend selection:

- default: native CPU backend
- env: `GPU_PROVE=1`
- flag: `--gpu-prove`

Examples:

```bash
# Regenerate prover artifacts with GPU proving enabled
cd relayer
go run -tags=icicle ./prover/cmd -gpu-prove ./bin ../contracts/verifiers

# Build a relayer binary with ICICLE support
go build -tags=icicle -o relayer ./cmd

# Run the relayer on GPU
./relayer start --config config.example.json --gpu-prove

# Equivalent env-based run
GPU_PROVE=1 ./relayer start --config config.example.json
```

### ICICLE Environment

Optional ICICLE tuning env vars read by `relayer/prover/backend_icicle.go`:

```bash
export GNARK_ICICLE_DEVICE_ID=0
export GNARK_ICICLE_BACKEND_LIBS=/usr/local/lib
export GNARK_ICICLE_PIN_KEYS=true
```

- `GNARK_ICICLE_DEVICE_ID`: GPU device index.
- `GNARK_ICICLE_BACKEND_LIBS`: backend library location passed into ICICLE.
- `GNARK_ICICLE_PIN_KEYS`: whether proving keys should be pinned to GPU memory.

If the ICICLE runtime is missing, the `icicle` build typically fails at link or
startup with errors such as `library 'icicle_device' not found`. In that case,
ensure ICICLE shared libraries are installed and `LD_LIBRARY_PATH` covers the
directory containing the `libicicle_*` files.

## Local E2E Test

End-to-end run on local Cosmos + Ethereum nodes. Requires Docker + Kurtosis on
top of the toolchain in [Requirements](#requirements).

The local Cosmos scripts must run against the custom Gaia branch that carries the
IBC host changes. By default they use `../gaia/build/gaiad` and require that
checkout to be on `test/ibc-host-customs`:

```bash
cd ../gaia
git checkout test/ibc-host-customs
GOTOOLCHAIN=go1.25.7 make build
cd ../fast-ibc
```

Set `GAIAD=/path/to/gaiad` or `GAIA_DIR=/path/to/gaia` to override this. The
scripts do not fall back to a `gaiad` binary from PATH unless you explicitly set
`GAIAD=gaiad`.

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
> The repo currently ships pre-generated `Groth16Verifier_N{4,8,16}.sol` only
> (`contracts/verifiers/`). If you keep the larger buckets in
> `relayer/prover/buckets.go` or expect chains with larger quorum sets,
> re-run `go run ./prover/cmd ./bin ../contracts/verifiers` from the
> `relayer/` directory to emit the missing `Groth16Verifier_N{32,64,128}.sol`
> before redeploying contracts — otherwise `SignatureVerifier.verifyBatchProof`
> will revert with `UnknownBucket(N)` for any signer count > 16.

```bash
# 1. (One-time) compile per-bucket circuits + emit Groth16Verifier_N{N}.sol.
#    Re-run only when circuit code changes. After this, redeploy contracts.
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
./relayer create-clients-eth    --config config.json

# 6. Start the bi-directional relay loop
./relayer start --config config.example.json
#    GPU run:
#    ./relayer start --config config.example.json --gpu-prove
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


# 8. check
cast call 0xee0fcb8e5ccad0b4197baabd633333886f5c364d \
  'ibcERC20Contract(string)(address)' \
  'transfer/cosmoshub-1/stake' \
  --rpc-url http://127.0.0.1:32003

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'fullDenomPath()(string)' \
  --rpc-url http://127.0.0.1:32003

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'balanceOf(address)(uint256)' \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 \
  --rpc-url http://127.0.0.1:32003

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'escrow()(address)' \
  --rpc-url http://127.0.0.1:32003


"$GAIAD" q txs \
  --query "message.action='/ibc.core.channel.v2.MsgAcknowledgement'" \
  --node tcp://127.0.0.1:26657 \
  -o json


cast receipt 0x4d611d65a802bea81865e7f0e1f0413064518a79883b29481968804d0efa1692--rpc-url http://127.0.0.1:56310 | grep -A3 "IBCAppRecvPacket\|topics\|data"


```

Send an ICS-20 transfer from Cosmos to trigger a client update + `recvPacket`
round-trip; the `[UpdateCosmosClient]` log line reports the chosen bucket.

## Benchmark mode

Detailed per-step gas + timing logs are off by default (production noise) and
opt-in via flag or env:

```bash
# CLI flag
./relayer start --config config.example.json --benchmark

# Env (equivalent)
RELAYER_BENCHMARK=1 ./relayer start --config config.example.json
```

Either source turns it on; flag is the override. The relayer logs
`[benchmark] enabled: detailed gas/timing logs are active` at startup so it's
obvious which mode you're in.

When enabled, the relayer emits the following lines so a single E2E run can be
diffed for performance regressions without rebuilding:

| Prefix | Where | Fields |
|---|---|---|
| `[bench][prover]` | per `GenerateProof` | `sigs`, `bucket`, `witness`, `prove`, `verify`, `total` |
| `[bench][eth]` | per `SendEthTx` / `SendEthTxBatch` | label or `multicall labels=...`, `gasUsed`, `submit`, `wait`, `total`, `tx` |
| `[bench][eth] inner[i]` | per multicall inner | `label`, `gas` from `debug_traceTransaction` callTracer |
| `[bench][cosmos]` | per `SendCosmosTx` / `SendCosmosTxBatch` | msg type or `batch msgs=N`, `gasWanted`, `gasUsed`, `broadcast`, `total`, `height`, `hash` |

`SendEthTxBatch` packs N inner calls but the receipt only reports the total.
For each multicall, after the tx confirms the relayer calls
`debug_traceTransaction` with the `callTracer`, descends past the UUPS
proxy → impl wrapper frame, and logs one `inner[i]` line per direct child of
the multicall body. Sample output:

```
[bench][eth] multicall labels=updateApplicationState,recvPacket:1 gasUsed=3322538 ...
[bench][eth] inner[0] updateApplicationState gas=2491085 (from trace)
[bench][eth] inner[1] recvPacket:1 gas=833989 (from trace)
```

`gas=... (from trace)` is the **real on-chain gasUsed** of that inner CALL,
not an estimate.

The inner-gas trace requires the RPC endpoint to expose the `debug_` namespace
(`--http.api=...,debug` on geth). Verify with:

```bash
curl -X POST <rpc> -d '{"jsonrpc":"2.0","method":"rpc_modules","params":[],"id":1}'
# look for "debug": "1.0" in the response
```

Kurtosis ethereum-package presets typically include it; production endpoints
typically do not. If unavailable, the relayer logs one error line and
continues — the rest of the benchmark output is unaffected.

`utils.SetBenchEnabled(true|false)` lets tests force the flag without touching
env. Definitions live in `relayer/utils/bench.go`.

## Local Cosmos ↔ OP E2E

Brings up an L1 + OP Stack L2 + Cosmos and relays both directions. The L1 is shared:
`run_optimism_node.sh` reuses `run_eth_node.sh` for it (Fusaka-from-genesis, the same
`eth-network-params.yaml` the ETH↔Cosmos devnet uses), then layers the L2 on top.

```bash
# 1. L1 (Fulu) + OP Stack L2 in one Kurtosis enclave. Ends with games proposed
#    and .op-devnet-run/attestor.env written (L1/L2/op-node/beacon endpoints).
./scripts/local/run_optimism_node.sh

# 2. Attestor — independent verifier over the replica op-node; serves gRPC :3001.
#    Sources attestor.env automatically.
./scripts/local/run_op_attestor.sh

# 3. Cosmos node, then gov-store BOTH light-client wasms:
#    the Ethereum client (L2 clients authenticate L1 through it) and the OP client.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm.sh        # -> ETH client checksum
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
# It patches relayer/config.example.json — copy the addresses into config.json.

# 5. Clients. Cosmos side first (ETH client, then the OP client anchored to it via
#    --l2-config), then the L2 side (SpectreClient deployed + addClient'd).
cd relayer
./relayer create-clients-cosmos --config config.json \
  --wasm-checksum <eth-checksum> --l2-config <op-l2-config.json>
./relayer create-clients-eth --config config.json --source <ics26_client_id> --trust-level 2/3

# 6. Relay both directions.
./relayer start --config config.json
```

## Local Cosmos ↔ Arbitrum E2E

Brings up an L1 + Arbitrum Nitro/BoLD L2 + Cosmos and relays both directions. The L1
is shared with the OP flow: `run_arbitrum_node.sh` attaches Arbitrum to the existing
Kurtosis enclave if OP already brought one up, or creates the same local L1 itself
through `run_eth_node.sh` when the enclave does not exist.

```bash
# 1. L1 (Fulu) + Arbitrum Nitro/BoLD L2 in one Kurtosis enclave. Ends with a
#    finalized AssertionCreated check and .arbitrum-devnet-run/attestor.env written
#    (L1/L2/beacon/Nitro feed/RollupCore endpoints).
./scripts/local/run_arbitrum_node.sh

# 2. Attestor — independent non-sequencing Nitro replica; serves gRPC :3001.
#    Sources .arbitrum-devnet-run/attestor.env automatically.
./scripts/local/run_arbitrum_attestor.sh

# 3. Cosmos node, then gov-store BOTH light-client wasms:
#    the Ethereum client (L2 clients authenticate L1 through it) and the Arbitrum client.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm.sh        # -> ETH client checksum
./scripts/local/wasm_arb.sh    # -> Arbitrum client checksum

# 4. IBC contracts on the Arbitrum L2 (E2ETestDeployL2). Use the same relayer key
#    rules as OP: the deployer receives the ICS26Router relayer role, so it must be
#    the key in relayer/.env ETH_PRIVATE_KEY unless you change both together.
DST_CHAIN=arbitrum \
L2_DEPLOYER_ADDRESS=0x8943545177806ED17B9F23F0a21ee5948eCaa776 \
L2_DEPLOYER_PRIVATE_KEY=bcdf20249abf0ed6d944c0288fad489e33f66b3960d9e6229c1cd214ed3bbe31 \
  ./scripts/local/deploy_l2_contracts.sh
# It patches the cosmos_to_l2 module whose dst_chain is arbitrum.

# 5. Clients. Cosmos side first (ETH client, then the Arbitrum client anchored to it
#    via --l2-config), then the L2 side (SpectreClient deployed + addClient'd).
cd relayer
./relayer create-clients-cosmos --config config.json \
  --wasm-checksum <eth-checksum> --l2-config <arb-l2-config.json>
#
#    If the OP flow already ran, an Ethereum client exists — anchor to it instead:
#      ./relayer create-clients-cosmos --config config.json \
#        --l1-client-id <08-wasm-N> --l2-config <arb-l2-config.json>
#    Without --l1-client-id this creates a SECOND Ethereum client and rewrites
#    cosmos_to_eth.cosmos_wasm_client_id, which repoints the Cosmos<->Ethereum path at
#    a client the Sepolia-side SpectreClient was never registered against; every
#    recvPacket then reverts on a counterparty mismatch. Since this section shares the
#    L1 enclave with the OP flow, running both in sequence hits exactly that.
./relayer create-clients-eth --config config.json --source <ics26_client_id> --trust-level 2/3

# 6. Relay both directions.
./relayer start --config config.json
```

The local Arbitrum l2-config is the same shape as OP's `--l2-config`, but its
`rollup_profile.protocol.type` is `bold_v2`. Fill it from the local handoffs:

- `l2_rpc_url`: `.arbitrum-devnet-run/attestor.env` `L2_RPC_URL`.
- `wasm_checksum`: checksum printed by `wasm_arb.sh`.
- `rollup_profile.common.l1_chain_id`, `l2_chain_id`, `rollup`: `L1_CHAIN_ID`,
  `L2_CHAIN_ID`, and `ROLLUP_CORE_ADDRESS` from `.arbitrum-devnet-run/attestor.env`.
- `rollup_profile.common.l2_router`: `ics26Router` deployed by `deploy_l2_contracts.sh`.
- `rollup_profile.protocol.value.assertions_mapping_slot` and
  `assertion_status_offset`: `ASSERTIONS_MAPPING_SLOT` and `ASSERTION_STATUS_OFFSET`.
- `counterparty_client_id`: the Arbitrum L2 router client id configured in the
  `cosmos_to_l2` module, for example `arb-client-0`.

#### Adding an L2 to a deployment that already relays Cosmos↔Ethereum

Every L2 client is anchored to an Ethereum light client on Cosmos, and that anchor is
baked into the L2 client's state at creation — it cannot be changed afterwards.

`create-clients-cosmos` normally **creates** that Ethereum client, and rewrites the
`cosmos_to_eth` module's `cosmos_wasm_client_id` with the new id. Re-running it just to
add a second L2 therefore does real damage to a working L1 path: the Ethereum-side
SpectreClient is still registered against the *old* client, so every `recvPacket` starts
reverting on a counterparty mismatch, and the original client is left orphaned with
nothing advancing it until it expires.

Pass `--l1-client-id` instead. It anchors the new L2 client(s) to an Ethereum client that
already exists, creates nothing else, and leaves the `cosmos_to_eth` module alone:

```bash
# Adding Arbitrum to a deployment whose Ethereum client is already 08-wasm-0
./relayer create-clients-cosmos --config config.json \
  --l1-client-id 08-wasm-0 --l2-config <arb-l2-config.json>
```

`--wasm-checksum` is not needed with it (nothing Ethereum-side is being created), and the
id is read back off-chain before anything is created, so a typo fails immediately rather
than producing an L2 client permanently anchored to a client that does not exist.

Creating several L2s in one go does not need the flag — a single run injects the same
newly-created Ethereum client into every `--l2-config`:

```bash
./relayer create-clients-cosmos --config config.json --wasm-checksum <eth-checksum> \
  --l2-config op.json --l2-config base.json --l2-config arb.json
```

The `--l2-config` file carries the full ICS-08 profile (see
[docs/L2_CLIENTS.md](docs/L2_CLIENTS.md#l2-client-creation-config)); `rollup_profile.common`
needs `l1_chain_id`, `l2_chain_id`, `ethereum_client` (`client_id` is injected by the command,
`wasm_checksum` is not — supply the Ethereum client's checksum as a byte array), `l2_router`,
`commitment_slot`, and `rollup_version`, plus the rollup-specific fields. For OP/Base use the
dispute-game fields (`dispute_game_factory`, `game_list_slot`, `root_claim_bytecode_offset`,
`output_root_format`, `l2_header_fork`). For Arbitrum use the tagged `protocol` object from
`packages/arbitrum-verifier/config/README.md`; the local BoLD devnet uses `bold_v2`.
`packages/op-verifier/config/op-sepolia.json` and
`packages/arbitrum-verifier/config/arbitrum-sepolia.json` are working public-network templates.

### Sending a test packet

The source client is the Cosmos client that tracks the **L2** (the OP/Arbitrum client, e.g.
`08-wasm-1`), not the Ethereum client. It needs a registered counterparty, and that
counterparty must equal the client id the SpectreClient was added under on the L2 router — the
relayer only picks up packets whose `destination_client` matches its configured `ics26_client_id`:

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
works at all — both are easy to leave at their placeholder and then see nothing happen:

- `l2_ics26_client_id` on the `l2_to_cosmos` module must be the client id the SpectreClient was
  added under on the L2 router. The L2 subscriber filters events by it, so a stale placeholder
  silently drops every `WriteAcknowledgement` the L2 emits.
- The attestor must run with `disable_derived_roots: true`. The OP header builder can only prove
  **game-backed** roots, while `RelayableHeight` follows the attestor's newest frontier — under
  the low-latency profile that frontier is a self-derived root, and the two never agree
  (`attested root at height N is source "derived", want game`, retried forever).

### Failure modes

| Symptom | Cause |
|---|---|
| `beacon api unavailable at ...:59717` | The `eth-to-cosmos` module still has the example's beacon URL. Point `eth_beacon_api_url` at this run's `ETH_BEACON_API` (from `attestor.env`). |
| `attested root at height N is source "derived", want game` | The OP header builder can only prove game-backed roots, but the attestor's frontier is a self-derived root. Set `disable_derived_roots: true` in the attestor config so its frontier is game-driven. |
| `eth_getProof(<DisputeGameFactory>): historical state ... is not available` | The pinned Ethereum client on Cosmos has stopped advancing, so the builder keeps proving at one ageing L1 block (the error names the block and how far behind it is). `start` advances it on demand (`[UpdateEthClient]` / `[EthClientUpdate]` in the log) right before each header build; if the message persists, check those lines and confirm the client id there is the one the L2 client was created against. Widening the L1's state retention does not help — a game covering a recent L2 block does not exist at an old L1 block at all. |
| `MsgStoreCode` fails: `reference-types not enabled` | The wasm was built with a plain `cargo build`. Build through `cosmwasm/optimizer` (`just build-cw-ics08-wasm-*`); a raw release build embeds features the CosmWasm VM rejects. |
| Optimizer fails: `rustc 1.86.0 is not supported ... requires rustc 1.90` | A dependency raised its MSRV above the optimizer image's Rust. Pin the dependency down (e.g. `cargo update -p ruint --precise 1.17.0`) or bump the optimizer image. |
| `MsgCreateClient`: `status Unknown: client state is not active` | 08-wasm Stargate allowlist is missing `ClientStatus` — see [docs/L2_CLIENTS.md](docs/L2_CLIENTS.md#host-requirements-and-verification). |
| Every L2 client update fails: `IBC host query failed: codespace: undefined, code: 1` | The `l2_to_cosmos` module's `rollup_profile.common.ethereum_client.client_id` is not the client the L2 client was created against, so the contract looks for a consensus state at a slot only the *other* client has. `create-clients-cosmos` writes the id back and `start` now refuses to boot on a mismatch — see [docs/L2_CLIENTS.md](docs/L2_CLIENTS.md#host-requirements-and-verification). |
| L2 client update panics the tx: `returning attributes from a contract is not allowed` | The deployed wasm predates the fix that made the client return data only. Rebuild through `cosmwasm/optimizer` and gov-store it. |
| `updateApplicationState` reverts, ~82k gas, no revert string | The relayer's signer lacks the ICS26Router relayer role on the L2. `cast run <tx>` shows `canCall(...) → false`; funding the address does not help. |
| Send fails: `timeout exceeds the maximum expected value` | Sent without `--absolute-timeouts`, so the CLI writes `timeout_timestamp` in nanoseconds while IBC v2 reads seconds. |
| ETH client stops advancing: `404 NOT_FOUND: Sync committee for period N not found` | The beacon does not serve a `light_client/bootstrap` for that period. The relayer takes the committee from the preceding period's update instead, so this should only appear if that update is also unavailable — check the endpoint serves `/eth/v1/beacon/light_client/updates`. |
| Ack fails: `unknown field account_proof, expected one of key, value, proof` | The L2 membership proof was built in the Ethereum L1 shape. The L2 client wants a bare `EvmStorageProof` — see [docs/L2_CLIENTS.md](docs/L2_CLIENTS.md). |
| `attested game N is not visible at the pinned L1 block M yet` | Normal wait, not an error: the game is posted at the L1 head and becomes provable once its creation block finalizes (~2 epochs). |
| `packet at height N not yet covered by the destination client (trusts M)` | Normal wait: the client landed on a game committing below the packet; the next game covers it. |
| A direction goes silent — no error, no retry, other directions healthy | A hung RPC. `kill -QUIT <relayer-pid>` dumps every goroutine; look for one blocked in `net/http.(*persistConn).roundTrip`. Note the dump kills the process. |
| Gov proposal ends `REJECTED` without votes | `wasm.sh` resolves the proposal id after a fixed `sleep`; if indexing is slower the id is empty and the vote step is skipped. Vote manually before the (short devnet) voting period ends. |
| `forge script` fails `insufficient funds ... have 0` | The deployer has no balance on the L2 — use an L2-funded account (step 4). |

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
./scripts/local/run_base_attestor.sh

# 3. Cosmos node, then gov-store BOTH light-client wasms: the Ethereum client
#    (L2 clients authenticate L1 through it) and the Base client.
./scripts/local/run_cosmos_node.sh
./scripts/local/wasm.sh          # -> ETH client checksum
./scripts/local/wasm_base.sh     # -> Base client checksum

# 4. Deploy the L2 IBC contracts onto Base. MODULE_NAME is required: matching on
#    dst_chain alone would also match the Optimism module.
L2_ENV_FILE=.base-devnet-run/attestor.env MODULE_NAME=cosmos-to-base \
DST_CHAIN=base ./scripts/local/deploy_l2_contracts.sh

# 5. Create the clients. Copy relayer/base-l2-config.example.json first and fill in
#    the two checksums from step 3 (top-level is the BASE client as hex;
#    ethereum_client.wasm_checksum is the ETH client as a byte array).
cd relayer && go build -o relayer ./cmd
./relayer create-clients-cosmos --config config.json \
    --wasm-checksum <eth-hex> --l2-config base-l2-config.json
./relayer create-clients-eth --config config.json \
    --source base-client-0 --trust-level 2/3

# 6. Relay.
./relayer start --config config.json
```

Everything in the "Local Cosmos ↔ OP E2E" section about the attestor's
`disable_derived_roots`, the ICS26Router relayer role, and `--absolute-timeouts` on
the test transfer applies unchanged — Base uses the same attestor and the same L2
contracts.

Useful:

```bash
./scripts/local/run_base_node.sh --status   # are the Base services up
./scripts/local/run_base_node.sh --reset    # redeploy Base only; L1 preserved
./scripts/local/run_base_node.sh --stop     # stop Base, leave the L1 running
kurtosis enclave rm -f base-devnet          # remove the L1 and everything on it
```
## Contracts

Core IBC protocol contracts:

- `ICS26Router.sol` — IBC packet routing
- `ICS20Transfer.sol` — Fungible token transfer (ICS-20)
- `SpectreClient.sol` — Tendermint light client (2/3 quorum + batch verify; owns client state in an ERC-7201 Store)
- `SignatureVerifier.sol` — Rebuilds CanonicalVote bytes, hashes witness, dispatches per bucket
- `Groth16Verifier_N{N}.sol` — Per-bucket Groth16 verifiers (N ∈ {4,8,16,32,64,128})
- `Membership.sol` — On-chain ICS23 Merkle proof verification

## License

MIT
