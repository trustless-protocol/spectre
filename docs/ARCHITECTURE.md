# Architecture

## System Overview

Solidity IBC Eureka is a production IBC v2 implementation for Ethereum-Cosmos interoperability, spanning three language layers:

```
                    Ethereum (EVM)                          Cosmos
              ┌──────────────────────┐             ┌──────────────────┐
              │   ICS26Router (UUPS) │◄── IBC ───► │  CometBFT Node   │
              │         │            │   packets    │   (Tendermint)   │
              │   ICS20Transfer      │             └──────────────────┘
              │    ├─ IBCERC20       │                     ▲
              │    └─ Escrow         │                     │
              │         │            │              ┌──────┴──────┐
              │  Groth16ICS07Tendermint  │◄── proofs ──│  Go Relayer │
              │    └─ WrapperVerifier│              │  (Groth16)  │
              │       └─ Groth16    │              └─────────────┘
              └──────────────────────┘
```

## Contract Hierarchy

```
AccessManager (OpenZeppelin RBAC)
    ↓ governs
ICS26Router (UUPS) ─── Main entry point for all IBC messages
    ├─ ICS20Transfer (UUPS) ─── ICS-20 fungible token transfers
    │   ├─ IBCERC20 (Beacon Proxy) ─── Wrapper ERC20 for bridged tokens
    │   └─ Escrow (Beacon Proxy) ─── Token custody during transfer
    └─ Groth16ICS07Tendermint (UUPS) ─── ZK light client for Cosmos chains
        └─ WrapperVerifier ─── Ed25519 decompression + SHA512
            └─ Groth16Verifier ─── Auto-generated from circuit VK
```

## Request Flow: IBC Transfer (Cosmos → Ethereum)

```
1. Cosmos chain commits IBC packet
2. Go Relayer detects packet via CometBFT WebSocket (subscriber/)
3. Relayer extracts validator Ed25519 signature (prover/extractor.go)
4. Relayer generates Groth16 proof: Ed25519 sig → gnark circuit → proof
5. Relayer submits to Ethereum:
   a. Groth16ICS07Tendermint.updateClient() — verifies header via Groth16
   b. ICS26Router.recvPacket() — routes to ICS20Transfer
   c. ICS20Transfer mints IBCERC20 tokens (or unlocks Escrow)
```

## Request Flow: IBC Transfer (Ethereum → Cosmos)

```
1. User calls ICS20Transfer.sendTransfer() on Ethereum
2. Tokens locked in Escrow (or IBCERC20 burned)
3. ICS26Router records packet commitment
4. Go Relayer detects SendPacket event (subscriber/)
5. Relayer builds MsgRecvPacket with membership proof
6. Cosmos chain verifies via Ethereum light client (CosmWasm)
7. Tokens released on Cosmos side
```

## Request Flow: ACK Relay (Ethereum → Cosmos)

```
1. Ethereum ICS26Router emits WriteAcknowledgement event
2. Go Relayer detects event via Ethereum event subscription (subscriber/)
3. Relayer converts Ethereum packet format to Cosmos IBC v2 format
4. Relayer submits MsgAcknowledgement to Cosmos chain
5. Cosmos chain processes the acknowledgement
```

## Data Flow: Groth16 Proof Generation

```
Ed25519 signature (R, S, pubKey, message)
    ↓
Decompress R, A points (off-chain)
    ↓
SHA512(R || A || message) → hash (off-chain)
    ↓
gnark PreHashCircuit witness (24 public inputs = 6 × 4 limbs)
    ↓
groth16.Prove() → proof[8], commitments[2], commitmentPok[2]
    ↓
On-chain: WrapperVerifier.verifyProof()
    → SHA512(R || A || msg) on-chain
    → Reduce mod L
    → Encode 24 limbs
    → Groth16Verifier.verifyProof()
```

## On-chain Encoding Pipeline

```
Solidity Encode.sol                    Go proto.Marshal() (source of truth)
─────────────────                     ──────────────────────────────────────
encodeVersion(Version)            ↔   cmtversion.Consensus.Marshal()
encodeValidator(SimpleValidator)  ↔   cmttypes.SimpleValidator.Marshal()
encodeBlockId(BlockId)            ↔   cmttypes.BlockID.Marshal()
encodePartSetHeader(PartSetHeader)↔   cmttypes.PartSetHeader.Marshal()
cdcEncodeString(str)              ↔   gogotypes.StringValue{Value: str}.Marshal()
cdcEncodeInt64(n)                 ↔   gogotypes.Int64Value{Value: n}.Marshal()
cdcEncodeBytes32(hash)            ↔   gogotypes.BytesValue{Value: hash}.Marshal()
encodeTimestamp(secs)             ↔   gogotypes.StdTimeMarshal(time)

Header.hashHeader()               ↔   types.Header.Hash() (14 fields → merkle)
Header.hashValSet()               ↔   ValidatorSet.Hash() (proto-encoded → merkle)
Header.merkleHash()                ↔   merkle.HashFromByteSlices() (1-byte prefix)
```

Cross-validated via `relayer/cmd/encode_debug/` + `test/solidity-ibc/EncodeTest.t.sol` (34 tests).

## Directory Map

| Directory | Language | Purpose |
|-----------|----------|---------|
| `contracts/` | Solidity | Core IBC protocol, ICS20, light client, encoding |
| `contracts/utils/` | Solidity | Encoding, hashing, verifiers, helpers |
| `contracts/programs/` | Solidity | UpdateClient, Membership, Misbehaviour verification |
| `contracts/light-clients/` | Solidity | Groth16ICS07Tendermint + message types |
| `relayer/` | Go | Relayer CLI + Groth16 prover |
| `relayer/cmd/` | Go | CLI: start, create-clients, genesis, fixtures |
| `relayer/prover/` | Go | Ed25519 → Groth16 proof generation |
| `relayer/client/` | Go | Tendermint RPC + Ethereum Beacon API + Ethereum light client state |
| `relayer/runner/` | Go | Service runner utilities |
| `relayer/services/` | Go | Context, Worker, batch processing, relay loop |
| `relayer/subscriber/` | Go | CometBFT WebSocket + Ethereum event listeners |
| `relayer/transaction/` | Go | ETH + Cosmos transaction submission (single + batch) |
| `relayer/bindings/` | Go | Auto-generated contract bindings |
| `packages/go-abigen/` | Go | Shared Go bindings for Solidity contracts |
| `packages/ethereum/` | Rust | Ethereum light client for CosmWasm |
| `packages/tendermint-light-client/` | Rust | Tendermint client types and provers |
| `programs/groth16-programs/` | Rust | RISC-V proving programs (legacy) |
| `e2e/interchaintestv8/` | Go | End-to-end tests with real chains |
| `scripts/` | Solidity | Deployment scripts |
| `test/` | Solidity | Foundry unit/integration/benchmark tests |

## Dependency Direction

```
test/ ──────────► contracts/ ◄────── scripts/
                     │
                     ▼ (abigen)
              packages/go-abigen/
              relayer/bindings/
                     │
                     ▼
              relayer/cmd/ ─── CLI entry point
              relayer/services/ ─── Context, Worker, relay loop
              relayer/runner/ ─── Service runner utilities
              relayer/subscriber/ ─── Event listeners (Cosmos + ETH)
              relayer/client/ ─── RPC + Beacon API clients
              relayer/transaction/ ─── Tx submission (single + batch)
              relayer/prover/ ─── Groth16 prover
                     │
                     ▼
              External: CometBFT, ecip-gnark, IBC-Go v10
```

## Upgrade Strategy

| Contract | Authority | Proxy Pattern | Upgrade Function |
|----------|-----------|---------------|------------------|
| ICS26Router | Admins | UUPS | `upgradeToAndCall` |
| ICS20Transfer | Admins | UUPS | `upgradeToAndCall` |
| Escrow | ICS20Transfer | Beacon | `upgradeEscrowTo` |
| IBCERC20 | ICS20Transfer | Beacon | `upgradeIBCERC20To` |

UUPS for core contracts (admin-governed). Beacon for instances (ICS20Transfer upgrades all atomically).
