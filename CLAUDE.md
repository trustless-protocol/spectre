# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Production Solidity implementation of **IBC v2** for Ethereum ↔ Cosmos interoperability. Three language layers: Solidity (contracts), Go (relayer), Rust (CosmWasm only).

Key difference from upstream `solidity-ibc-eureka`: uses **gnark Groth16** (Go, Ed25519) instead of SP1 (Rust, RISC-V zkVM) for Tendermint light client verification.

## Commands

```bash
# Build
just build-contracts              # Compile Solidity
just build-go-relayer             # Build Go relayer (go build ./...)

# Test
just test-foundry                 # All Solidity tests
forge test --match-test <name> -vvv  # Single Solidity test
just test-go-relayer              # All Go relayer tests
cd relayer && go test -run <TestName> ./...  # Single Go test
just test-e2e <name>              # E2E test by name

# Lint
just lint                         # All linters (solidity + go + buf)
just lint-solidity                # forge fmt + solhint + natlint

# Security
just slither                      # Slither static analysis

# Generate
just generate-abi                 # Extract ABIs

# Encoding cross-validation
forge test --match-contract EncodeTest -vvv  # Solidity encoding tests

# Node setup (local development) — invoke from any cwd; each script self-anchors to repo root
./scripts/local/run_cosmos_node.sh         # Local Cosmos node with test accounts
./scripts/local/run_cosmos_node_docker.sh  # Docker-based Cosmos node setup
./scripts/local/run_eth_node.sh            # Ethereum testnet via Kurtosis + deploy contracts
./scripts/local/wasm.sh                    # Submit Ethereum light client WASM via governance
./scripts/local/wasm_docker.sh             # Docker-based WASM submission
```

**Prerequisites**: `bun` (not npm/yarn), `just`, Foundry, Go 1.21+. E2E also needs Docker + Kurtosis.

## Relayer CLI

```bash
# One-time setup: deploy light clients on both chains
cd relayer
go run ./cmd/main.go create-clients \
  --config config.json \
  --trust-level 1/3 \
  --wasm-checksum <hex>
# Copy ICS07 address from log into config.json cosmos_to_eth.ics07_client

# Start relay loop (bi-directional: Cosmos ↔ ETH)
go run ./cmd/main.go start --config config.json

# Generate genesis state
go run ./cmd/main.go genesis --trusted-block 0 --trusting-period 0
```

Config: JSON file with `modules` array containing `cosmos_to_eth` and `eth_to_cosmos` entries (see `relayer/config.example.json`).
Secrets: `relayer/.env` ships with default sample keys for the local Kurtosis devnet (`ETH_PRIVATE_KEY`, `COSMOS_PRIVATE_KEY`, `COSMOS_CHAIN_ID`, `PROVER_BIN_DIR`) — replace before any real deployment.

## Documentation Map

**IMPORTANT: Read ALL docs/ files at the start of every conversation for full project context.**

| Document | Contents |
|----------|----------|
| [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | System diagram, contract hierarchy, request flows, encoding pipeline, directory map |
| [docs/DESIGN.md](docs/DESIGN.md) | Coding conventions, protobuf encoding rules, naming, formatting, access control roles |
| [docs/PRODUCT_SENSE.md](docs/PRODUCT_SENSE.md) | IBC domain model, token transfer flow, API surface, business rules |
| [docs/QUALITY.md](docs/QUALITY.md) | Test strategy (Foundry/Go/E2E), CI/CD, linting, fixtures |
| [docs/SECURITY.md](docs/SECURITY.md) | Trust model, proof verification chain, attack surfaces, static analysis |
| [docs/RELIABILITY.md](docs/RELIABILITY.md) | Error handling, observability, known risks, recovery procedures |
| [docs/logging.md](docs/logging.md) | Structured logging conventions |
| [docs/metrics.md](docs/metrics.md) | RED metrics, Prometheus conventions |

## Core Conventions

- **Encoding**: `Encode.sol` must match Go `proto.Marshal()` exactly — cross-validate via `EncodeTest.t.sol`
- **Proxy pattern**: UUPS for core contracts, Beacon for instances
- **Solidity formatting**: line length 120, tab width 4, double quotes (`foundry.toml`)
- **Go relayer**: `go.mod` replace directives for local `ecip-gnark`/`decentrio-gnark` — adjust per dev setup
- **Bindings**: After Solidity changes, regenerate Go bindings with `abigen` (output: `packages/go-abigen/`)
- **Relayer config**: JSON config file (not .env) for `start` command; `.env` only for secrets/prover paths
- **E2E**: interchaintest suites require Docker + Kurtosis + compiled binaries

## Key Architecture

```
ICS26Router (UUPS) ← main IBC entry point
  ├─ ICS20Transfer (UUPS) ← token bridge
  │   ├─ IBCERC20 (Beacon) ← bridged token wrapper
  │   └─ Escrow (Beacon) ← token custody
  └─ Groth16ICS07Tendermint (UUPS) ← ZK light client + 2/3 quorum
      └─ WrapperVerifier ← rebuilds CanonicalVote, hashes witness, dispatches by bucket
          └─ Groth16Verifier_N{N} (one per N ∈ {4,8,16,32,64,128})
```

ZK flow: top-N validator Ed25519 sigs (≥2/3 voting power) → padded to nearest
bucket with deterministic dummy keypairs → BatchCircuit hashes the witness
into a single SHA-256 public input + ECIP batch verify → per-bucket
`Groth16Verifier_N{N}.sol` checks the proof. R/S are bound only by the proof
(not in calldata or witness hash); A is bound to keep the on-chain pubkey
lookup honest.

## Go Relayer Structure

```
relayer/
├── bindings/       # Auto-generated Go bindings for Solidity contracts
├── client/         # Tendermint + Ethereum RPC/Beacon API clients
├── prover/         # Bucketed Ed25519 batch prover
│   ├── buckets.go      # Buckets + smallestBucketGEQ
│   ├── circuit.go      # BatchCircuit (hash-aggregate witness commit)
│   ├── hash_witness.go # Off-chain witness layout (matches WrapperVerifier)
│   ├── dummy.go        # Deterministic dummy keypair padding
│   ├── extractor.go    # Quorum selection from CometBFT commit
│   ├── prover.go       # Bucket registry + GenerateProof
│   ├── cmd/            # One-shot setup tool: compile every bucket + emit Groth16Verifier_N{N}.sol
│   └── bin/            # Per-bucket artifacts: bin/n{N}/{r1cs,pk,vk}.bin
├── runner/         # Service runner utilities
├── services/       # Context, Worker, batch builder, relay loop, PendingPacketTracker
│   └── pending.go      # In-memory tracker for Cosmos-originated packets pending ETH delivery (1h TTL)
├── subscriber/     # Cosmos WebSocket + Ethereum event listeners
├── transaction/    # Ethereum + Cosmos transaction submission
├── utils/          # IBC path helpers, byte utils
├── test/           # Manual test script (reference relay flow)
└── cmd/main.go     # CLI: start, create-clients, genesis, fixtures
```

Setup-circuits entrypoint: from `relayer/`, run
`go run ./prover/cmd ./bin ../contracts/verifiers`. After regeneration, the
per-bucket verifier vk changes — redeploy `Groth16Verifier_N{N}.sol` and
re-register them via `WrapperVerifier.setBucket(...)`.

## Go Bindings Package

```
packages/go-abigen/
├── groth16ics07tendermint/  # ICS07 Tendermint light client
├── ics26router/             # IBC router
├── ics20transfer/           # Token transfer
├── ibcerc20/                # Bridged ERC20 wrapper
└── relayerhelper/           # Helper contract
```
