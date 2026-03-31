# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Production Solidity implementation of **IBC v2** for Ethereum ↔ Cosmos interoperability. Three language layers: Solidity (contracts), Go (operator), Rust (CosmWasm only).

Key difference from upstream `solidity-ibc-eureka`: uses **gnark Groth16** (Go, Ed25519) instead of SP1 (Rust, RISC-V zkVM) for Tendermint light client verification.

## Commands

```bash
# Build
just build-contracts              # Compile Solidity
just build-operator               # Build Go operator (go build)

# Test
just test-foundry                 # All Solidity tests
forge test --match-test <name> -vvv  # Single Solidity test
just test-operator                # All Go operator tests
cd operator && go test -run <TestName> ./...  # Single Go test
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
```

**Prerequisites**: `bun` (not npm/yarn), `just`, Foundry, Go 1.21+. E2E also needs Docker + Kurtosis.

## Documentation Map

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
- **Go operator**: `go.mod` replace directives for local `ecip-gnark`/`decentrio-gnark` — adjust per dev setup
- **Bindings**: After Solidity changes, regenerate Go bindings with `abigen`
- **E2E**: interchaintest suites require Docker + Kurtosis + compiled binaries

## Key Architecture

```
ICS26Router (UUPS) ← main IBC entry point
  ├─ ICS20Transfer (UUPS) ← token bridge
  │   ├─ IBCERC20 (Beacon) ← bridged token wrapper
  │   └─ Escrow (Beacon) ← token custody
  └─ Groth16ICS07Tendermint (UUPS) ← ZK light client
      └─ WrapperVerifier → Groth16Verifier
```

ZK flow: Ed25519 sig → Groth16 proof (Go operator/gnark) → on-chain verification (Groth16Verifier.sol)

## Go Operator Structure

```
operator/
├── bindings/       # Auto-generated Go bindings for Solidity contracts
├── client/         # Tendermint + Ethereum RPC clients
├── keys/           # Key management
├── prover/         # gnark Groth16 circuit + prover (Ed25519)
├── services/       # Relay loop: subscriber → batch builder → transaction handler
├── subscriber/     # Cosmos event subscriber
├── transaction/    # Ethereum transaction submission
└── cmd/main.go     # CLI: start, genesis, membership
```
