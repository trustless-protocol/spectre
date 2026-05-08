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

## Local E2E Test

End-to-end run on local Cosmos + Ethereum nodes. Requires Docker + Kurtosis on
top of the toolchain in [Requirements](#requirements).

```bash
# 1. (One-time) compile per-bucket circuits + emit Groth16Verifier_N{N}.sol.
#    Re-run only when circuit code changes. After this, redeploy contracts.
cd relayer
go run ./prover/cmd ./bin ../contracts/verifiers

# 2. Build the relayer binary
go build -o relayer ./cmd

# 3. Start Ethereum first and wait until the beacon node finalizes.
#    Replace 56246 with your Kurtosis-mapped beacon RPC port.
./run_eth_node.sh        # Kurtosis Ethereum testnet + deploys core contracts
# Poll until finalized.epoch > 0:
curl -s <eth_beacon_api_ur>/eth/v1/beacon/states/head/finality_checkpoints
curl -s http://127.0.0.1:51180/eth/v1/beacon/states/head/finality_checkpoints

# 4. Then start Cosmos and submit the Ethereum LC WASM via governance
#    This requires a wasm-enabled Cosmos binary (08-wasm), e.g. simd:
#    COSMOS_BIN=simd ./run_cosmos_node.sh
#    COSMOS_BIN=simd ./wasm.sh
#    If you only have stock gaiad, use the container flow instead:
#    ./run_cosmos_node_docker.sh
#    ./wasm_docker.sh
./run_cosmos_node.sh
./wasm.sh

# 5. Deploy Tendermint light client on Ethereum.
#    Copies the ICS07 address back into relayer/config.json automatically.
./relayer create-clients \
  --config config.json \
  --wasm-checksum <hex-from-wasm.sh>

# 6. Start the bi-directional relay loop
./relayer start --config config.example.json

# 7. send tx

gaiad tx ibc-transfer transfer cosmoshub-1 08-wasm-0 0x8943545177806ED17B9F23F0a21ee5948eCaa776 1000stake --from test1 --keyring-backend test --gas-prices 1stake --packet-timeout-timestamp 600 --generate-only | jq '.body.messages[0].encoding = "application/x-solidity-abi"' | gaiad tx sign /dev/stdin --from test1 --keyring-backend test | gaiad tx broadcast /dev/stdin -y
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
