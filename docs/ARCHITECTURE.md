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
    └─ Groth16ICS07Tendermint (UUPS) ─── ZK light client + 2/3 quorum check
        ├─ _cachedValidatorSets ─── Storage cache (validatorsHash → ValidatorSet)
        │                          populated only after a full update succeeds; lets
        │                          subsequent updates skip the per-call Header.hashValSet
        │                          for unchanged sets, and lets the relayer omit the
        │                          val-set bytes from calldata on cache hits.
        └─ WrapperVerifier ─── Rebuilds CanonicalVote bytes per slot,
            │                  hashes the full witness (active flag, pubkey,
            │                  msgLen, msg) into a single SHA-256 digest, and
            │                  dispatches to the per-bucket verifier
            └─ Groth16Verifier_N{N} ─── One contract per bucket size
                                        (N ∈ {4, 8, 16, 32, 64, 128}),
                                        auto-generated from each bucket's VK
```

`updateClient` dispatches to one of two `UpdateClient` library entry points
depending on whether the validator-set hashes for this update are already
cached (or implied by the adjacent fast-path):

- `updateClientResolved(msg)` — used when the validator set hashes are known
  good (cache hit, adjacent update, or trustedNext matches current). Skips
  the per-call `Header.hashValSet(...)` re-derivation, since
  `_prepareUpdateClientMessage` has already either loaded the cached set or
  validated the supplied one against `validatorsHash`.
- `updateClient(msg)` (raw calldata pass-through) — used on the cold path
  when neither side is resolvable, performing the full hashValSet check
  inside `validateBasic`.

## Request Flow: IBC Transfer (Cosmos → Ethereum)

```
1. Cosmos chain commits IBC packet
2. Go Relayer detects packet via CometBFT WebSocket (subscriber/)
3. Relayer extracts top-N validator Ed25519 signatures until ≥ 2/3 voting
   power (prover/extractor.go), pads to nearest bucket with deterministic
   dummy keypairs (prover/dummy.go).
4. Relayer generates a single Groth16 batch proof for the bucket
   (prover/circuit.go BatchCircuit + prover/prover.go bucket registry).
5. Relayer submits to Ethereum:
   a. Groth16ICS07Tendermint.updateClient() — checks unique-signer 2/3
      quorum, then dispatches to WrapperVerifier.verifyBatchProof().
      `_prepareUpdateClientMessage` first decides whether each validator set
      can be served from `_cachedValidatorSets` (cache hit) or — for an
      adjacent update — skipped via the `trustedNext == current` shortcut.
      A cache hit serves the on-chain pubkeys + voting powers directly and
      lets the call route through `UPDATE_CLIENT.updateClientResolved`,
      which skips the `Header.hashValSet` re-derivation entirely.
   b. WrapperVerifier rebuilds each slot's CanonicalVote bytes from the
      shared block header + per-slot Timestamp, hashes the witness, and
      forwards the SHA-256 digest as the proof's only public input to the
      bucket-specific Groth16Verifier_N{N}.
   c. ICS26Router.recvPacket() — routes to ICS20Transfer
   d. ICS20Transfer mints IBCERC20 tokens (or unlocks Escrow)
```

The Cosmos→ETH `updateClient` and the packet's `recvPacket` are submitted
as a single `ICS26Router.multicall(...)` so the update + packet apply
atomically. The relayer's `BuildCosmosClientUpdateMsg` queries the
authoritative on-chain trusted height before regenerating a Groth16 proof,
so subsequent chunks of the same source block (or any race where another
relayer or a crashed-replayer has already advanced the client) short-
circuit to `HasMsg=false` and skip the ~90 s proof gen.

## Request Flow: IBC Transfer (Ethereum → Cosmos)

```
1. User calls ICS20Transfer.sendTransfer() on Ethereum
2. Tokens locked in Escrow (or IBCERC20 burned)
3. ICS26Router records packet commitment; emits SendPacket event
4. Go Relayer detects SendPacket via Ethereum log subscription (subscriber/)
5. ethProofHeight: poll Beacon API GetFinalityUpdate (≤ 60 × 10 s) until the
   event's exec-block is finalised
6. UpdateEthClient: relay sync-committee updates into 08-wasm on Cosmos so
   LatestExecutionBlockNumber catches up to the event
7. Fetch ETH storage proof at the finalised block via
   client.GetEthMembershipProof(ICS26_STORAGE_SLOT, clientID||1||seq)
8. Broadcast MsgRecvPacket{packet, proof, height={0, LatestSlot}, signer}
   on Cosmos; 08-wasm verifies membership, then ICS Core dispatches to the
   destination app (mint / unlock)
9. If the packet times out before step 8 lands, timeoutEthSend bounces back
   through Groth16ICS07Tendermint to submit MsgTimeoutPacket on ETH
```

## Request Flow: Cosmos-originated Timeout (background)

```
1. Every cosmos send_packet logged into BatchBuilder.PendingTracker on
   intake
2. scanForCosmosTimeouts (30 s tick) walks the tracker, picks entries whose
   TimeoutTimestamp < current ETH block.time
3. For each expired packet, build a non-membership proof of the receipt
   path on ETH at the finalised height
4. Bundle one MsgUpdateClient + N MsgTimeout into a single Cosmos tx; on
   success, remove the entries from the tracker
5. PurgeStale(1h) drops any entry that never settles, preventing tracker
   growth from leaked packets
```

## Relayer Concurrency Model

`StartLoop` (`relayer/services/main.go`) drains the two batch channels on
**independent goroutines** — one for Cosmos→ETH (`handleCosmos`), one for
ETH→Cosmos (`handleEth`) — so a slow direction (proof gen, beacon-finality
wait) never stalls the other.

The two handlers touch disjoint mutable state:

| State                          | Owner              |
| ------------------------------ | ------------------ |
| `ctx.latestEthTimestamp`       | handleCosmos       |
| `ctx.latestCosmosTimestamp`    | handleEth          |
| `Services.lastCosmosAppHashHeight` | handleCosmos   |
| `Services.lastFinalizedExecBlock`  | handleEth      |

`Timestamp.Snapshot/Set/SetTime` accessors guard the shared
`latestEth/CosmosTimestamp` fields against the routine + timeout-scanner
goroutines. The two `Services.last*` memoization fields are single-writer
per goroutine and need no locking.

Two per-chunk waits are memoized so chunked flushes don't repeat
identical setup work:

- `waitCosmosAppHash(targetHeight)` polls Cosmos height until it covers
  `maxPacketHeight + 2` (so the packet commitment is in the queried
  AppHash). Memoizes the highest confirmed height; later chunks of the
  same source block return without an RPC.
- `waitBeaconFinality(eventBlock)` memoizes the highest finalized exec
  block; subsequent chunks whose event block is already covered return
  immediately without re-pinging the beacon API.

## Request Flow: ACK Relay (Ethereum → Cosmos)

```
1. Ethereum ICS26Router emits WriteAcknowledgement event
2. Go Relayer detects event via Ethereum event subscription (subscriber/)
3. Relayer converts Ethereum packet format to Cosmos IBC v2 format
4. Relayer submits MsgAcknowledgement to Cosmos chain
5. Cosmos chain processes the acknowledgement
```

## Data Flow: Groth16 Batch Proof Generation

```
Pick top-N validator signatures by voting power until ≥ 2/3 (extractor.go)
    ↓
Pad to bucket size with deterministic dummy keypairs (dummy.go); each pad
slot signs `"fast-ibc-dummy" || bucket(2 BE) || slot(2 BE)` so it's a
valid Ed25519 signature with active=false.
    ↓
For each slot i: Sig=(R,S), Pub=A, msg=CanonicalVote bytes, msgLen, active
    ↓
gnark BatchCircuit (circuit.go):
   - In-circuit witness hash: per slot
        active(1) || A(32) || msgLen(2 BE) || msg[MaxMsgLen padded]
     → SHA-256 → 32-byte public input (Hash[32]uints.U8)
   - eddsa.VerifyBatchWithMsgBytes (ECIP aggregate over all slots) using
     each slot's full CanonicalVote bytes truncated to msgLen.
    ↓
groth16.Prove(byBucket[N]) → proof[8], commitments[2], commitmentPok[2]
    ↓
On-chain WrapperVerifier.verifyBatchProof(bucket, ..., pubkeys, ts*, active, shared):
   - Rebuild CanonicalVote bytes from (shared, timestampSeconds[i], timestampNanos[i])
     for active slots, or DummyMsgBytes(bucket, i) for padding slots.
   - _hashWitness recomputes the SHA-256 witness commit byte-for-byte.
   - Pack the 32-byte digest as 32 uint256 public inputs and dispatch via
     selector to Groth16Verifier_N{bucket}.verifyProof().
```

R and S are deliberately NOT in calldata or the witness hash — the Groth16
proof itself binds them via the in-circuit Ed25519 verify, and no on-chain
logic consumes them. A is hashed because Solidity uses `pubkeys[i]` to look
up the validator's voting power for the quorum check; binding A defeats
calldata pubkey swap attacks.

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

Cross-validated via `test/solidity-ibc/EncodeTest.t.sol`.

## Directory Map

| Directory | Language | Purpose |
|-----------|----------|---------|
| `contracts/` | Solidity | Core IBC protocol, ICS20, light client, encoding |
| `contracts/utils/` | Solidity | Encoding, hashing, verifiers, helpers |
| `contracts/programs/` | Solidity | UpdateClient, Membership, Misbehaviour verification |
| `contracts/light-clients/` | Solidity | Groth16ICS07Tendermint + message types |
| `relayer/` | Go | Relayer CLI + Groth16 prover |
| `relayer/cmd/` | Go | CLI: start, create-clients, genesis, fixtures |
| `relayer/prover/` | Go | Bucketed Ed25519 batch prover (BatchCircuit, witness hash, dummy padding) |
| `relayer/prover/cmd/` | Go | One-shot tool: compile + setup every bucket, write artifacts + emit `Groth16Verifier_N{N}.sol` |
| `relayer/prover/bin/` | binary | Per-bucket circuit artifacts (`bin/n{N}/{r1cs,pk,vk}.bin`) |
| `relayer/client/` | Go | Tendermint RPC + Ethereum Beacon API + Ethereum light client state |
| `relayer/runner/` | Go | Service runner utilities |
| `relayer/services/` | Go | Context, Worker, batch processing, relay loop, PendingPacketTracker (Cosmos→ETH timeout fallback) |
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
