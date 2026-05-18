# fast-ibc

A Solidity implementation of IBC Eureka (IBC v2) with a Go relayer using gnark
Groth16 for Tendermint light client verification. Each `updateClient` proves a
2/3+ voting-power quorum of validator Ed25519 signatures in a single Groth16
proof, with in-circuit CanonicalVote reconstruction so only ~32 bytes of public
input land on-chain.

## Architecture

```
┌─────────────────┐         ┌──────────────────────────────────────┐
│   Cosmos Chain  │         │           Ethereum Chain             │
│                 │         │                                      │
│  IBC v2 Module  │────────▶│  ICS26Router.sol                     │
│  (send_packet)  │         │      │                               │
└─────────────────┘         │      ▼                               │
                            │  Groth16ICS07Tendermint.sol          │
         ┌──────────────────│      │  (light client + 2/3 quorum)  │
         │                  │      ▼                               │
         │  Go Relayer      │  WrapperVerifier.sol                 │
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
matching `Groth16Verifier_N{N}.sol`; `WrapperVerifier` dispatches by bucket.

The circuit reconstructs each validator's `CanonicalVote` bytes from a shared
block header + per-slot `Timestamp`, hashes the full witness (active flag,
pubkey, msg) into a single SHA-256 digest, and exposes that digest as the only
public input. This keeps the on-chain verifier well under EIP-170. Padding
slots carry `active=false`; both the in-circuit hash and the on-chain quorum
check skip them.


## Requirements

- [Go](https://golang.org/) >= 1.21
- [Foundry](https://getfoundry.sh/)
- [Bun](https://bun.sh/)
- [Just](https://github.com/casey/just)

## Sibling repos (required to build the relayer)

The Go relayer's `relayer/go.mod` has `replace` directives pointing at two sibling
repos via relative paths:

```
replace (
    0x5ea000000/ecip-gnark      => ../../ecip-gnark
    github.com/consensys/gnark  => ../../decentrio-gnark
)
```

Clone both **next to** `fast-ibc` (so they sit two directories up from
`relayer/`) before running `just install-go-relayer` / `just build-prover-artifacts`:

```bash
# from the parent directory that contains fast-ibc
git clone https://github.com/decentrio/ecip-gnark
git clone https://github.com/decentrio/gnark decentrio-gnark
```

Resulting layout:

```
parent/
├── fast-ibc/
├── ecip-gnark/        # provides 0x5ea000000/ecip-gnark (Ed25519 in-circuit ops, garaga_rs FFI)
└── decentrio-gnark/   # fork of consensys/gnark v0.13.0 with hash-aggregate verifier tweaks
```

Without these, `go build ./...` under `relayer/` fails with
`replacement directory ../../ecip-gnark does not exist`.

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

```bash
# 1. (One-time) compile per-bucket circuits + emit Groth16Verifier_N{N}.sol.
#    Re-run only when circuit code changes. After this, redeploy contracts.
cd relayer
go run ./prover/cmd ./bin ../contracts/verifiers

# 2. Build the relayer binary
go build -o relayer ./cmd

# 3. Start Ethereum first and wait until the beacon node finalizes.
#    Replace 56246 with your Kurtosis-mapped beacon RPC port.
./scripts/local/run_eth_node.sh   # Kurtosis Ethereum testnet + deploys core contracts
# Poll until finalized.epoch > 0:
curl -s http://127.0.0.1:56246/eth/v1/beacon/states/head/finality_checkpoints

# 4. Then start Cosmos and submit the Ethereum LC WASM via governance
./scripts/local/run_cosmos_node.sh   # local Cosmos chain with funded test accounts
./scripts/local/wasm.sh              # submit + vote-pass the Ethereum LC WASM proposal

# 5. Deploy Tendermint light client on Ethereum.
#    Copies the ICS07 address back into relayer/config.json automatically.
./relayer create-clients \
  --config config.json \
  --wasm-checksum <hex-from-wasm.sh>

# 6. Start the bi-directional relay loop
./relayer start --config config.json
```

Send an ICS-20 transfer from Cosmos to trigger an `updateClient` + `recvPacket`
round-trip; the `[UpdateCosmosClient]` log line reports the chosen bucket.

## Contracts

Core IBC protocol contracts:

- `ICS26Router.sol` — IBC packet routing
- `ICS20Transfer.sol` — Fungible token transfer (ICS-20)
- `Groth16ICS07Tendermint.sol` — Tendermint light client (2/3 quorum + batch verify)
- `WrapperVerifier.sol` — Rebuilds CanonicalVote bytes, hashes witness, dispatches per bucket
- `Groth16Verifier_N{N}.sol` — Per-bucket Groth16 verifiers (N ∈ {4,8,16,32,64,128})
- `Membership.sol` — On-chain ICS23 Merkle proof verification

## License

MIT
