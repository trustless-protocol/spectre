# Spectre

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
`Groth16Verifier_N{N}.sol`. The prover engine discovers which buckets are
servable by scanning `bin/` for complete artifact triples at startup
(`relayer/prover/gnark/engine.go`) — no code change is needed once a bucket's
artifacts exist. No artifact triples are checked in (`relayer/bin/` is
gitignored), so a fresh clone must generate or provision them before the prover
can serve any bucket. By default the generator (`relayer/prover/cmd`,
`-buckets` flag) builds N=4; serving a larger quorum requires generating that
bucket's artifacts, then publishing, deploying, and registering its
`Groth16Verifier_N{N}.sol`.

The store keeps one 32-byte `keccak256(abi.encode(consensusState))` per height,
not the root itself. A packet carries the consensus state and its commitment
root, and `SpectreClient._validateMembershipInput` rehashes what it was given,
compares that against the stored hash, checks the root matches the state, and
checks the trusting period — all before the ICS23 proof is looked at.

The update above is drawn router-managed, which is the access-controlled path
(`ICS02ClientUpgradeable.updateApplicationState` is `restricted` and forwards to
the client). A deployment that is not router-managed has the relayer call
`SpectreClient` directly instead — `relayer/transaction/ethereum.go:670-694`
picks between the two. The packet path always goes through the router.

## Building

### Requirements

- [Go](https://golang.org/) >= 1.26 (`relayer/go.mod` pins 1.26.0, toolchain 1.26.9)
- [Foundry](https://getfoundry.sh/)
- [Bun](https://bun.sh/)
- [Just](https://github.com/casey/just)
- Optional for GPU proving: ICICLE runtime/libs installed on the host, plus an
  `icicle` build of the relayer/prover tool

### gnark forks

The relayer/prover depend on two private forks of gnark, pinned as ordinary Go
modules in `relayer/go.mod`:

```
require github.com/decentrio/ecip-gnark v0.1.0-rc.3

replace github.com/consensys/gnark => github.com/decentrio/gnark v0.16.4-decentrio-fork.2
```

- [`decentrio/ecip-gnark`](https://github.com/decentrio/ecip-gnark) — Ed25519 in-circuit ops + the `garaga_rs` Rust FFI.
- [`decentrio/gnark`](https://github.com/decentrio/gnark) — fork of consensys/gnark with hash-aggregate verifier tweaks.

Both repositories are private, so the Go proxy and checksum DB cannot serve
them. One-time setup:

```bash
go env -w GOPRIVATE='github.com/decentrio/*'
gh auth setup-git        # or any other git credential for github.com/decentrio
```

To hack on a fork, point a local `replace` at your checkout while you work and
never commit it (machine-local paths do not belong in `go.mod`).

#### Garaga FFI

`ecip-gnark`'s prover links a Rust FFI (`libgaraga_rs`). Its cgo directive is
`-L..`, which lands in the read-only Go module cache, so the library is built
into `.artifacts/native/` instead. Build it once (and after bumping the
`ecip-gnark` version):

```bash
scripts/build-garaga-ffi.sh
```

Then point the linker and the loader at it (the `just` recipes and the scripts
under `scripts/` already set `CGO_LDFLAGS`; the loader path is still yours to set):

```bash
export CGO_LDFLAGS="-L$PWD/.artifacts/native"
export LD_LIBRARY_PATH="$PWD/.artifacts/native"       # DYLD_LIBRARY_PATH on macOS
```

Without the built FFI, `go build ./...` under `relayer/` fails at link with
`library 'garaga_rs' not found`.

### Build and test

```bash
bun install                  # Solidity dependencies (never npm/yarn)
just build-go-relayer        # go build ./... in relayer/
just build-prover-artifacts  # compile circuits into relayer/bin/n{N}/ and emit contracts/verifiers/
just build-contracts         # forge build — needs the generated verifier (see below)

just test-foundry            # all Solidity tests
just test-go-relayer         # all Go relayer tests
just lint                    # solidity + go + buf (`just lint-rust` is separate)
```

The order matters on a fresh clone. `contracts/verifiers/Groth16Verifier_N4.sol`
is produced by `just build-prover-artifacts` and is not checked in, while the
deploy scripts (`scripts/E2ETestDeploy.s.sol`, `scripts/E2ETestDeployL2.s.sol`)
and two integration tests import it. Until it exists, `just build-contracts`
stops with `Source "contracts/verifiers/Groth16Verifier_N4.sol" not found`.

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

## One relay path per process

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

## Optional GPU proving

CPU proving remains the default. GPU proving is opt-in: build with
`-tags=icicle`, then enable the ICICLE backend via env or flag when needed.

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

**Regenerating prover artifacts changes the verifying key.** The `go run` example
above writes straight into `relayer/bin` and `contracts/verifiers`. The supported
path is `just build-prover-artifacts` (`scripts/solidity/build-prover-artifacts.sh`),
which covers the configured buckets (N4 by default): it generates into ignored
staging (`.artifacts/solidity/prover`), checks that the staged set is complete,
publishes it, records provenance, and does nothing when every artifact already
exists unless you pass `--force`. It does not pass `-tags=icicle`, so it cannot use
the GPU backend. However you regenerate, every deployed `Groth16Verifier_N{N}` must
be redeployed and re-registered with `SignatureVerifier.setBucket(...)`, or every
proof fails on-chain. Do not regenerate for a live deployment without agreeing that
redeploy first.

### ICICLE Environment

Optional ICICLE tuning env vars read by `relayer/prover/gnark/backend_icicle.go`:

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
`[bench] enabled: detailed gas/timing logs are active` at startup so it's
obvious which mode you're in.

When enabled, the relayer emits the following lines so a single E2E run can be
diffed for performance regressions without rebuilding:

| Prefix | Where | Fields |
|---|---|---|
| `[prover]` | per proof job in the async prover (`relayer/prover/prover.go`), always on | `proof_id`, `bucket`, `witness` (queued) / `done in <duration>` |
| `[bench] eth` | per `SendEthTx` / `SendEthTxBatch` (`relayer/transaction/ethereum.go`) | label or `multicall labels=...`, `gasUsed`, `submit`, `wait`, `total`, `tx` |
| `[bench] eth inner[i]` | per multicall inner | `label`, `gas` from `debug_traceTransaction` callTracer |
| `[bench] cosmos` | per `SendCosmosTx` / `SendCosmosTxBatch` (`relayer/transaction/cosmos.go`) | msg type or `batch msgs=N`, `gasWanted`, `gasUsed`, `broadcast`, `total`, `height`, `hash` |

`SendEthTxBatch` packs N inner calls but the receipt only reports the total.
For each multicall, after the tx confirms the relayer calls
`debug_traceTransaction` with the `callTracer`, descends past the UUPS
proxy → impl wrapper frame, and logs one `inner[i]` line per direct child of
the multicall body. The lines look like this (values elided):

```
[bench] eth multicall labels=updateApplicationState,recvPacket:1 gasUsed=<total> ...
[bench] eth inner[0] updateApplicationState gas=<n> (from trace)
[bench] eth inner[1] recvPacket:1 gas=<n> (from trace)
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
