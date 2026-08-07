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
## Running it end to end

Bring-up runbooks live in **[docs/E2E.md](docs/E2E.md)** — Cosmos↔Ethereum,
Cosmos↔OP, Cosmos↔Arbitrum, Cosmos↔Base, running against an L2 you did not
deploy, sending test packets in both directions, and the failure modes each
step produces.

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
