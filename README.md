# fast-ibc

A Solidity implementation of IBC Eureka (IBC v2) with a Go relayer using gnark Groth16 for Tendermint light client verification.

## Architecture

```
┌─────────────────┐         ┌──────────────────────────────────────┐
│   Cosmos Chain  │         │           Ethereum Chain             │
│                 │         │                                      │
│  IBC v2 Module  │────────▶│  ICS26Router.sol                     │
│  (send_packet)  │         │      │                               │
└─────────────────┘         │      ▼                               │
                            │  Groth16ICS07Tendermint.sol              │
         ┌──────────────────│      │  (light client)               │
         │                  │      ▼                               │
         │  Go Operator     │  Groth16Verifier.sol                 │
         │  ┌────────────┐  │      │  (Ed25519 gnark circuit)      │
         │  │ UpdateClient│  │      ▼                               │
         │  │  (gnark ZK) │  │  Membership.sol                     │
         │  │             │  │      │  (ICS23 on-chain proof)       │
         │  │ RecvPacket  │  │      ▼                               │
         │  └────────────┘  │  ICS20Transfer.sol                   │
         └──────────────────└──────────────────────────────────────┘
```


## Requirements

- [Go](https://golang.org/) >= 1.21
- [Foundry](https://getfoundry.sh/)
- [Bun](https://bun.sh/)
- [Just](https://github.com/casey/just)

## Quick Start

```bash
# Install dependencies
bun install

# Build Solidity contracts
just build-contracts

# Build the Go relayer
just build-relayer
```

## Contracts

Core IBC protocol contracts:

- `ICS26Router.sol` — IBC packet routing
- `ICS20Transfer.sol` — Fungible token transfer (ICS-20)
- `Groth16ICS07Tendermint.sol` — Tendermint light client (uses gnark Groth16)
- `Groth16Verifier.sol` — Custom Groth16 verifier (Ed25519 gnark circuit)
- `Membership.sol` — On-chain ICS23 Merkle proof verification

## Go Relayer

The relayer handles the relay loop between Cosmos and Ethereum:

1. **Subscribe** to Cosmos `send_packet` events
2. **UpdateClient** — generate gnark Groth16 ZK proof of Tendermint consensus
3. **RecvPacket** — submit packet with ICS23 membership proof to Ethereum

```bash
cd relayer
go build ./...
```

## Development

```bash
# Run Solidity tests
just test-foundry

# Run Go relayer tests
cd relayer && go test ./...

# Lint
just lint
```

## License

MIT
