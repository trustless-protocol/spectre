# fast-ibc

A ZK light client that lets Ethereum follow Tendermint consensus, and the Go
relayer and prover that feed it. A purpose-built gnark Groth16 circuit
batch-verifies the validator Ed25519 signatures behind one CometBFT commit — a
2/3+ voting-power quorum, with each `CanonicalVote` reconstructed in-circuit —
so a client update settles on-chain as a single proof carrying ~32 bytes of
public input. Signers are padded up to the nearest fixed circuit size, one
`Groth16Verifier_N{N}` per bucket, for validator sets up to 180 active members.

The client plugs into an [IBC v2](https://github.com/cosmos/ibc/tree/main/spec/IBC_V2)
Solidity stack for packet routing and ICS-20 transfers, and the same relayer
carries the reverse direction against a CosmWasm Ethereum light client on the
Cosmos side (`programs/cw-ics08-wasm-eth`), which follows the beacon sync
committee rather than a Groth16 proof. CosmWasm ICS-08 clients extend the path
to OP, Base and Arbitrum, backed by attestor sidecars.

## Architecture

```
Cosmos                  Go relayer                  Ethereum
──────                  ──────────                  ────────

commit @ H         ──▶  extract top-N Ed25519  ──▶  ICS26Router.updateApplicationState
(validator sigs)        sigs holding ≥ 2/3          └──▶ SpectreClient.updateApplicationState
                        voting power                     └──▶ UpdateClient.sol   (delegatecall)
                        gnark: one Groth16                    └──▶ SignatureVerifier
                        proof, ~32 bytes                           └──▶ Groth16Verifier_N{N}
                        of public input             ⇒ keccak(consensusState) stored per height

send_packet        ──▶  packet + ICS23 proof   ──▶  ICS26Router.recvPacket
                                                    ├──▶ SpectreClient.verifyMembership
                                                    │    └──▶ Membership.sol      (staticcall)
                                                    │         ICS23 against the commitment root
                                                    └──▶ ICS20Transfer.onRecvPacket
                                                         └──▶ Escrow / IBCERC20
```

The two flows are independent, and that is the point of the design: a Groth16
proof is paid once per client update, and every packet relayed against that
update costs only an ICS23 membership check. Nothing on the packet path touches
the verifier contracts.

The Tendermint commit is verified by batch-proving N Ed25519 signatures whose
voting power sums to ≥ 2/3 of the validator set. To keep Groth16 circuits
fixed-size, production retains the bucket topology N ∈ {4, 8, 16, 32, 64,
128}. Each bucket has one inseparable `(r1cs, pk, vk)` set and matching
`Groth16Verifier_N{N}.sol`. The checked-in repository generator/prover manifest
currently enables N=4 only; serving a larger quorum requires enabling that
bucket in the prover and publishing, deploying, and registering its matching
artifact set.

The store keeps one 32-byte `keccak256(abi.encode(consensusState))` per height,
not the root itself. A packet carries the consensus state and its commitment
root, and `SpectreClient._validateMembershipInput` rehashes what it was given,
compares that against the stored hash, checks the root matches the state, and
checks the trusting period — all before the ICS23 proof is looked at.

The update above is drawn router-managed, which is the access-controlled path
(`ICS02ClientUpgradeable.updateApplicationState` is `restricted` and forwards to
the client). A deployment that is not router-managed has the relayer call
`SpectreClient` directly instead — `relayer/transaction/handler.go:638-642`
picks between the two. The packet path always goes through the router.

## Building

### Requirements

- [Go](https://golang.org/) >= 1.25 (`relayer/go.mod` pins 1.25.7)
- [Foundry](https://getfoundry.sh/)
- [Bun](https://bun.sh/)
- [Just](https://github.com/casey/just)
- Optional for GPU proving: ICICLE runtime/libs installed on the host, plus an
  `icicle` build of the relayer/prover tool

### gnark submodules

The relayer/prover depend on two forks of gnark, vendored as **git submodules**
under `third_party/` and wired through `relayer/go.mod`:

```
replace (
    0x5ea000000/ecip-gnark      => ../third_party/ecip-gnark
    attestor/types              => ../attestor/types
    github.com/consensys/gnark  => ../third_party/decentrio-gnark
)
```

- `third_party/ecip-gnark` → [`decentrio/ecip-gnark`](https://github.com/decentrio/ecip-gnark) — Ed25519 in-circuit ops + the `garaga_rs` Rust FFI.
- `third_party/decentrio-gnark` → [`decentrio/gnark`](https://github.com/decentrio/gnark) — fork of consensys/gnark v0.13.0 with hash-aggregate verifier tweaks.

fast-ibc stores only a **pinned commit** of each fork (a submodule pointer), not
their files.

#### First checkout

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

### Build and test

```bash
bun install                  # Solidity dependencies (never npm/yarn)
just build-contracts         # forge build
just build-go-relayer        # go build ./... in relayer/
just build-prover-artifacts  # compile circuits into relayer/bin/n{N}/ and emit the verifiers

just test-foundry            # all Solidity tests
just test-go-relayer         # all Go relayer tests
just lint                    # solidity + go + buf (`just lint-rust` is separate)
```

`just --list` shows every recipe, including the per-L2 wasm builds, the
attestor sidecars, and Slither.

## Running it end to end

Bring-up runbooks live in **[docs/E2E.md](docs/E2E.md)** — Cosmos↔Ethereum,
Cosmos↔OP, Cosmos↔Arbitrum, Cosmos↔Base, running against an L2 you did not
deploy, sending test packets in both directions, and the failure modes each
step produces.

## L2 ICS-08 clients

The repository also contains independent CosmWasm clients for OP Sepolia, Base
Sepolia, and Arbitrum Sepolia. Their deployment evidence requirements,
host-query wiring, and verification commands are documented in
[docs/L2_CLIENTS.md](docs/L2_CLIENTS.md).

## Production deployment

See [docs/PRODUCTION_DEPLOYMENT.md](docs/PRODUCTION_DEPLOYMENT.md) before using
the production AccessManager deployment scripts. It documents the governance
timelock requirement, the escrow launch gate, and the scheduled light-client
provisioning flow.

## Updating the gnark forks

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

## Optional GPU proving

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
./relayer start --config config.json --gpu-prove

# Equivalent env-based run
GPU_PROVE=1 ./relayer start --config config.json
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

### One relay path per process

`start` serves exactly one source→destination path and refuses a config declaring
more than one, so any real deployment runs several processes side by side.

**Give each `start` process its own `COSMOS_PRIVATE_KEY`.** Every path writes to
Cosmos — even a Cosmos→L2 path that appears to write only to the L2, because its
timeout scanner refunds expired Cosmos-origin packets on Cosmos. Two processes
sharing a key share an account sequence with nothing coordinating them, and the
loser of that race can have a single-packet transaction dropped, stranding the
transfer in escrow until it times out.

`ETH_PRIVATE_KEY` is a weaker rule: nonces are tracked per `{chain id, address}`,
so a Cosmos↔Ethereum process (L1) and a Cosmos↔L2 process (L2) can share one key.
Only processes writing the **same** EVM chain need separate ones.

Set them per process rather than editing `relayer/.env` — `godotenv` does not
override an already-set variable:

```bash
COSMOS_PRIVATE_KEY=$KEY_ETH_PATH ./relayer start --config config.json &
COSMOS_PRIVATE_KEY=$KEY_OP_PATH  ./relayer start --config config.op.json &
```

Full reasoning, including the one other command that shares the Cosmos key:
[docs/E2E.md](docs/E2E.md#signing-keys-when-several-processes-share-a-chain).

## Benchmark mode

Detailed per-step gas + timing logs are off by default (production noise) and
opt-in via flag or env:

```bash
# CLI flag
./relayer start --config config.json --benchmark

# Env (equivalent)
RELAYER_BENCHMARK=1 ./relayer start --config config.json
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

## License

MIT
