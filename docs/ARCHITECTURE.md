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
              │     SpectreClient    │◄── proofs ──│  Go Relayer │
              │    └─ SignatureVerifier│            │  (Groth16)  │
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
    └─ SpectreClient (immutable, no proxy) ─── ZK light client + 2/3 quorum check
        │                       Owns ALL client state in an ERC-7201 namespaced
        │                       Store (consensus states, the pinned validator
        │                       set, client state). The ⅔ quorum accounting and
        │                       every Store write live here, not in the modules.
        ├─ modules/ ─── Stateless verification logic, invoked by SpectreClient:
        │   ├─ UpdateClient (delegatecall) ─── reads the shared Store, never writes
        │   ├─ Misbehaviour (delegatecall) ─── reads the shared Store, never writes
        │   └─ Membership (staticcall) ─── pure ICS-23, no Store access
        ├─ pinned validator set ─── current set stored as an SSTORE2 blob
        │                          (≤ 180 validators) plus per-height snapshots
        │                          binary-searched by trusted height; rotated only
        │                          by updateConsensusState after verifying the new
        │                          set against the header's nextValidatorsHash
        └─ SignatureVerifier ─── Recomputes the witness commitment from calldata
            │                  (PrefixHead ‖ roundPresent ‖ blockHash ‖ per-slot
            │                  active ‖ pubkey), exposes the SHA-256 digest as two
            │                  128-bit public inputs, and dispatches to the
            │                  per-bucket verifier
            └─ Groth16Verifier_N{N} ─── One contract per bucket size
                                        (N ∈ {4, 8, 16, 32, 64, 128}),
                                        auto-generated from each bucket's VK
```

SpectreClient exposes two update entry points, both of which verify the
per-block Groth16 signature proof but differ in how they treat the validator
set:

- `updateApplicationState(bytes)` — the frequent, per-packet path. Advances
  the header's appHash against the already-pinned validator set without
  re-storing any validator data. This is the hot path folded into every
  packet multicall.
- `updateConsensusState(bytes)` — the rare, ~24h path. Rotates and re-pins the
  validator set and advances consensus state, requiring
  `hashValSet(newValSet) == header.nextValidatorsHash` before the new set is
  pinned into the Store (new SSTORE2 blob + a per-height snapshot for
  historical anchoring).

The `UpdateClient` module (invoked via delegatecall) performs the header /
signature checks against the shared Store; SpectreClient itself does the ⅔
quorum accounting and all Store writes.

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
   a. SpectreClient.updateApplicationState() — the frequent per-packet path:
      advances the header's appHash against the already-pinned validator set.
      SpectreClient resolves each proof signer against the pinned validator
      set (SSTORE2 blob; duplicate-signer bitmask), requires the summed voting
      power to exceed 2/3 of the pinned total, and writes the new consensus
      state into its Store, invoking the `UpdateClient` module via
      delegatecall (the module reads the shared Store but never writes) and
      dispatching to SignatureVerifier.verifyBatchProof(). The rare
      validator-set rotation goes through `updateConsensusState()` instead
      (see below).
   b. SignatureVerifier recomputes the witness commitment from calldata
      (PrefixHead ‖ roundPresent ‖ blockHash ‖ per-slot active ‖ pubkey)
      and forwards the SHA-256 digest, packed into two 128-bit public
      inputs, to the bucket-specific Groth16Verifier_N{N}.
   c. ICS26Router.recvPacket() — routes to ICS20Transfer
   d. ICS20Transfer mints IBCERC20 tokens (or unlocks Escrow)
```

The Cosmos→ETH client update and the packet's `recvPacket` are submitted
as a single `ICS26Router.multicall(...)` so the update + packet apply
atomically. Because both `updateApplicationState` and `updateConsensusState`
are exposed as router passthroughs, the relayer folds whichever one the block
requires into the packet multicall. The relayer's
`BuildCosmosClientUpdateMsg` queries the authoritative on-chain trusted height
before regenerating a Groth16 proof, so subsequent chunks of the same source
block (or any race where another relayer or a crashed-replayer has already
advanced the client) short-circuit to `HasMsg=false` and skip the ~90 s proof
gen.

## Request Flow: IBC Transfer (Ethereum → Cosmos)

```
1. User calls ICS20Transfer.sendTransfer() on Ethereum
2. Tokens locked in Escrow (or IBCERC20 burned)
3. ICS26Router records packet commitment; emits SendPacket event
4. Go Relayer detects SendPacket via Ethereum log subscription (subscriber/)
5. ethProofHeight: poll Beacon API GetFinalityUpdate (≤ 60 × 10 s) until the
   event's exec-block is finalised
6. Advance the 08-wasm client: the `beacon` builder gathers the sync-committee
   update and the `cosmos` destination submits it (WaitForCosmosCatchUp) so
   LatestExecutionBlockNumber catches up to the event
7. Fetch ETH storage proof at the finalised block via
   client.GetEthMembershipProof(ICS26_STORAGE_SLOT, clientID||1||seq)
8. Broadcast MsgRecvPacket{packet, proof, height={0, LatestSlot}, signer}
   on Cosmos; 08-wasm verifies membership, then ICS Core dispatches to the
   destination app (mint / unlock)
9. If the packet times out before step 8 lands, the ETH-timeout scanner
   (`scanForEthTimeouts`, 30 s tick) submits MsgTimeoutPacket on ETH; a send
   already past its timeout is classified permanent at the source and dropped
   from the relay path so the scanner owns the refund
```

## Request Flow: Cosmos-originated Timeout (background)

```
1. Every cosmos send_packet logged into BatchBuilder.PendingTracker on
   intake
2. scanForCosmosTimeouts (30 s tick) walks the tracker, picks entries whose
   TimeoutTimestamp < current ETH block.time
3. Drop entries whose ETH receipt already exists; they were delivered but
   lingered in local pending state
4. For each remaining expired packet, build a non-membership proof of the receipt
   path on ETH at the finalised height
5. Bundle one MsgUpdateClient + N MsgTimeout into a single Cosmos tx; on
   success, remove the entries from the tracker
6. PurgeStaleWithoutTimeout(1h) drops only no-timeout entries; timeout-bearing
   packets remain tracked until a receipt pre-check or timeout tx clears them
```

## Relayer Concurrency Model

`runAdapterEngine` (`relayer/cmd/run_adapters.go`) runs the two directions as two
`relay.Module` instances (`relayer/relay/module.go`) on **independent goroutines**
— one for Cosmos→ETH, one for ETH→Cosmos — sharing a child context so the first
fatal error stops both. A slow direction (proof gen, finality wait) never stalls
the other. One engine runs per `cosmos_to_eth` source; the shared prover,
`TransactionHandler`, and ETH endpoint are reused, with ETH events partitioned by
the per-source router client id. (This replaced the monolithic `StartLoop` +
`handleCosmos`/`handleEth`, removed after the cutover.)

Each module owns a single append-only cursor:

| State | Owner |
| ----- | ----- |
| `Module.lastHeight` — highest source height whose `ClientUpdate` was submitted | the module for that direction |

`lastHeight` is guarded by the module's mutex (touched by the event and refresh
goroutines) and advances only after a successful client update — never on failure.
The legacy `latestEth/CosmosTimestamp` and `Services.last*` memoization fields are
gone; refresh scheduling now reads on-chain state directly (`Destination.ClientExpiresAt`
for expiry, `PinnedSetRotationDueIn` for the guaranteed pinned-set rotation).

Instead of blocking per-chunk waits, the module gates on `Source.RelayableHeight`
— a non-blocking value (Cosmos `latest-2` for the AppHash H+2 lag; ETH the finalized
execution block). A packet above it is re-queued with a growing waiting backoff
(3→15 s) rather than blocking a goroutine on an RPC poll (this replaced
`waitCosmosAppHash` / `waitBeaconFinality`).

Serialization of chain writes is unchanged: all ETH sends go through the nonce
block under `h.mu`; all Cosmos sends through `SendCosmosTxBatch` under `cosmosMu`.

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
   - In-circuit witness hash (#199 prefix-binding layout):
        PrefixHead(11) || roundPresent(1) || BlockHash(32)
        || per slot: active(1) || A(32)
     → SHA-256 → two 128-bit public field elements. Msgs, MsgLens and Sigs
     stay private witness — the in-circuit Ed25519 verify binds the full
     signed bytes, so hashing them too would be redundant.
   - eddsa.VerifyBatchWithMsgBytes (ECIP aggregate over all slots) using
     each slot's full CanonicalVote bytes truncated to msgLen.
    ↓
groth16.Prove(byBucket[N]) → proof[8], commitments[2], commitmentPok[2]
    ↓
On-chain SignatureVerifier.verifyBatchProof(bucket, proof, commitments,
commitmentPok, pubkeys, active, shared):
   - _hashWitness recomputes the same PrefixHead ‖ roundPresent ‖ blockHash
     ‖ per-slot active ‖ pubkey layout from calldata byte-for-byte — no
     per-slot timestamps or vote bytes travel in calldata.
   - Split the 32-byte digest into two 128-bit uint256 public inputs and
     dispatch via the registered selector to Groth16Verifier_N{bucket}.verifyProof().
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
Header internal merkle hashing    ↔   merkle.HashFromByteSlices() (1-byte prefix)
```

Cross-validated via `test/light-clients/spectre/EncodeTest.t.sol`.

## Directory Map

| Directory | Language | Purpose |
|-----------|----------|---------|
| `contracts/` | Solidity | Core IBC protocol, ICS20, light client, encoding |
| `contracts/core/` | Solidity | ICS-02 registry, ICS-24 commitments, ICS-26 routing, migration |
| `contracts/apps/ics20/` | Solidity | ICS-20 transfer, escrow, voucher, callbacks, and rate limits |
| `contracts/light-clients/` | Solidity | Compatibility interface and message leaf |
| `contracts/light-clients/spectre/` | Solidity | Spectre client, verifier boundary, modules, stores, messages, and libraries |
| `contracts/shared/` | Solidity | Cross-package access libraries and compatibility interfaces |
| `contracts/periphery/` | Solidity | Read-only relayer helper contracts |
| `scripts/deployments/` | Solidity | Deployment-only AccessManager selector and configuration helpers |
| `relayer/` | Go | Relayer CLI + Groth16 prover |
| `relayer/cmd/` | Go | CLI: start, create-clients-{cosmos,eth}, update-client, genesis, fixtures |
| `relayer/prover/` | Go | Bucketed Ed25519 batch prover (BatchCircuit, witness hash, dummy padding) |
| `relayer/prover/cmd/` | Go | One-shot tool: compile + setup each bucket enabled in `prover.Buckets` (currently N4 locally), write paired artifacts, and emit `Groth16Verifier_N{N}.sol` |
| `relayer/prover/bin/` | binary | Per-bucket circuit artifacts (`bin/n{N}/{r1cs,pk,vk}.bin`) |
| `relayer/client/` | Go | Tendermint RPC + Ethereum Beacon API + Ethereum light client state |
| `relayer/runner/` | Go | Service runner utilities |
| `relayer/services/` | Go | Context, Worker, batch processing, relay loop, PendingPacketTracker (Cosmos→ETH timeout fallback) |
| `relayer/subscriber/` | Go | CometBFT WebSocket + Ethereum event listeners |
| `relayer/transaction/` | Go | ETH + Cosmos transaction submission (single + batch) |
| `relayer/bindings/` | Go | Auto-generated contract bindings |
| `packages/go-abigen/` | Go | Shared Go bindings for Solidity contracts |
| `packages/ethereum/` | Rust | Ethereum light client for CosmWasm |
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
| SpectreClient | per-clientId migrator role | none (immutable) | replaced via `ICS26Router.migrateClient` |

UUPS for core contracts (admin-governed). Beacon for instances (ICS20Transfer upgrades all atomically).
The light client is deliberately not upgradeable: it is initialized in its constructor and, if it must
be replaced (bug, frozen client, circuit change), a new instance is deployed and swapped in with
`migrateClient`, which is gated by a client-scoped role on the AccessManager.
