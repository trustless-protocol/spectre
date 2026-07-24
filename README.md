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
#    Replace 56246 with your Kurtosis-mapped beacon RPC port.
./scripts/local/run_eth_node.sh   # Kurtosis Ethereum testnet + deploys core contracts
# Poll until finalized.epoch > 0:
curl -s http://127.0.0.1:59717/eth/v1/beacon/states/head/finality_checkpoints

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

# 5. Create the light clients on both chains. Runs the Cosmos side first
#    (creates the 08-wasm ETH client), then deploys the Tendermint light
#    client on Ethereum. Writes cosmos_wasm_client_id + spectre_client back
#    into relayer/config.json automatically.
#    (Split alternative: create-clients-cosmos then create-clients-eth.)
./relayer create-clients \
  --config config.json \
  --wasm-checksum <hex-from-wasm.sh>

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
#      ./relayer create-clients --config config.json --source osmosis-1 --wasm-checksum <hex>
#      ./relayer start --config config.json

# 7. send packet

ABS_TIMEOUT=$(($(date +%s) + 2000))

gaiad tx ibc-transfer transfer transfer 08-wasm-0 0x8943545177806ed17b9f23f0a21ee5948ecaa776 1000stake \
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
| gaiad tx sign /dev/stdin \
    --from test1 \
    --home "$HOME/.gaia" \
    --chain-id test-ibc-eth \
    --keyring-backend test \
| gaiad tx broadcast /dev/stdin \
    --node tcp://127.0.0.1:26657 \
    -y


# 8. check
cast call 0xee0fcb8e5ccad0b4197baabd633333886f5c364d \
  'ibcERC20Contract(string)(address)' \
  'transfer/cosmoshub-1/stake' \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'fullDenomPath()(string)' \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'balanceOf(address)(uint256)' \
  0x8943545177806ed17b9f23f0a21ee5948ecaa776 \
  --rpc-url http://127.0.0.1:59619

cast call 0x016f5f33DbCb653e6393698Beba9DC19d828D75e \
  'escrow()(address)' \
  --rpc-url http://127.0.0.1:59619


gaiad q txs \
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
