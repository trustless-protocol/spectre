# Avalanche ↔ IBC — Architecture Overview

**Cosmos ↔ Avalanche C-Chain interoperability without a bridge committee.**
Descriptive document: what the system is, how each side verifies the other, and what it assumes.
Last reviewed 2026-09-22.

---

## 1. What this is

Two chains that cannot read each other, connected through IBC light-client interfaces instead of a
bridge multisig.

| Direction | Verifier | What it verifies |
|---|---|---|
| Cosmos → C-Chain | **Spectre**, a ZK light client deployed on the C-Chain | Tendermint validator signatures reaching ⅔ of voting power, compressed into one Groth16 proof |
| C-Chain → Cosmos | **Avalanche Warp light client**, an ICS-08 CosmWasm module on the Cosmos chain | A ≥67% stake-weighted BLS aggregate from Avalanche primary-network validators over a C-Chain block hash, plus a Merkle proof of IBC state under the authenticated state root |

Both directions use the same IBC interfaces, but their underlying security proofs differ: Spectre
proves a Tendermint consensus statement; the Avalanche client verifies a stake-weighted Warp
attestation and binds it to C-Chain state (§8 draws that line exactly).

Both carry standard IBC v2 traffic, with ICS-20 as the initial application. No independent bridge
attestor committee sits between the two chains.

### 1.1 Current trust boundary

Security-critical verification happens on-chain, in the two light clients. The **relayer and the
signature aggregator are untrusted transport**: compromising them delays or omits updates but cannot
produce an authenticated state transition (I7).

One qualification, stated up front rather than buried: the **host chain's bootstrap and governance
path is inside the trust boundary today.** It selects the initial Avalanche validator-set snapshot
and authorizes every rotation of it (§9.1). Protocol-attested rotation is designed, not implemented.

---

## 2. Why the Avalanche direction is the hard one

Cosmos consensus produces a portable artifact: a block commit carrying validator signatures. Anyone
can verify it outside the source chain — that is what Spectre compresses into a ZK proof.

Avalanche's Snowman finality exposes no equivalent portable certificate. Finality emerges from
repeated randomized sampling among validators, and the votes are transient network messages, never
aggregated into a commit a third party can verify later. A Tendermint-style light client for
Avalanche therefore cannot be built.

What Avalanche does expose is **Warp Messaging**: validators register BLS keys on the P-Chain and, on
request (ACP-118), sign messages about chain data. Those signatures aggregate into one 96-byte
BLS12-381 signature plus a signer bitset over a known validator set — checkable by anyone holding
that set and its stake weights.

This architecture uses that primitive as the external verification mechanism:

> **The Cosmos client verifies a stake-weighted attestation by Avalanche validators. It does not
> claim to verify Snowman consensus itself.**

---

## 3. Architecture at a glance

```text
              Cosmos chain                                      Avalanche C-Chain
   ┌──────────────────────────────────┐                ┌──────────────────────────────────┐
   │ IBC core (ibc-go)                │                │ ICS26Router ─ ICS20Transfer      │
   │  └ 08-wasm client:               │                │      │            └ Escrow       │
   │    Avalanche Warp light client   │                │      └ Spectre ZK light client   │
   │      • pinned validator set      │                │        • pinned Cosmos valset    │
   │      • quorum + BLS verification │                │        • Groth16 verifier        │
   │      • Coreth header identity    │                │                                  │
   │      • MPT state proof           │                │                                  │
   └──────────────┬───────────────────┘                └───────────────┬──────────────────┘
                  │                                                    │
        header + aggregate + proof                          Groth16 proof + header
                  │                                                    │
   ┌──────────────┴────────────────────────────────────────────────────┴──────────────────┐
   │ Relayer (untrusted transport)                                                        │
   │   • C-Chain event subscription + gap recovery   • Cosmos event subscription           │
   │   • Warp header builder ──► signature-aggregator sidecar ──► validator p2p (ACP-118)  │
   │   • Groth16 prover for the Cosmos → C-Chain leg  • timeout scanners on both legs      │
   └──────────────────────────────────────────────────────────────────────────────────────┘
```

### 3.1 Protocol permissionlessness vs reference deployment

The connection is permissionless at the protocol level: anyone who can construct a valid update or
packet proof may submit it. The reference deployment happens to run one relayer process per
source→destination pairing — an operational choice, not a protocol requirement.

---

## 4. Security invariants

These define the intended security properties. Everything in §6 is an elaboration of them.

**I1 — Authenticated state only.** No unverified header, state root, or router storage root enters
client state:

```text
stored_consensus_state ⇒ full_update_verification_passed
```

The authenticated-header type has a private constructor, so a transition cannot be invoked with a
deserialized or hand-built header.

**I2 — Proofs are bound to authenticated state.** Every accepted storage proof verifies against a
root authenticated by a verified header:

```text
accepted_proof ⇒ proof.root == authenticated_header.state_root
```

**I3 — Quorum uses the full snapshot weight as denominator.** `total_weight` is the whole primary
network's stake at the snapshot height, including validators with no registered BLS key — those are
not listed in the set (they cannot sign and must not occupy a bitset index), but their stake stays in
the denominator:

```text
signed_stake · quorum_den  ≥  total_weight · quorum_num        (67/100)
```

Keyless stake therefore reduces available signing capacity instead of silently vanishing from the
denominator. Loading a snapshot whose `total_weight` is below the sum of listed weights fails.

**I4 — Signatures are domain-separated.** The client rebuilds the signed message from its own
configured network ID and source chain ID plus the locally derived block hash. A signature for
another Avalanche network or another chain cannot be replayed here.

**I5 — Misbehaviour requires authenticated conflicts.** Two malformed headers are not evidence:

```text
verify(A) ∧ verify(B) ∧ conflict(A, B)  ⇒  freeze
```

**I6 — Client height never moves backward.** `latest_height` is monotonic (an older header may still
be submitted; it simply cannot lower the frontier), and stored adjacent heights must chain by parent
hash.

**I7 — Relayer and aggregator hold no verification authority.** They can affect liveness only; they
cannot create a state transition without data that passes the full pipeline.

**I8 — Validator-set snapshots are monotonic.** A replacement snapshot must carry a strictly greater
P-Chain height than the one it replaces, so no old set can be replayed as a new one.

---

## 5. Cosmos → C-Chain

The C-Chain is a geth-derived EVM chain with BN254 precompiles and Cancun semantics, so this
direction reuses the stack already running against Ethereum and three L2 families: IBC router, ICS-20
transfer application, escrow, Spectre client, EVM relayer path.

Per update Spectre proves Tendermint validator signatures reaching the required voting power against
the pinned Cosmos validator set; the C-Chain verifies the Groth16 proof, then the membership proof of
the packet commitment. The pinned Cosmos set rotates on its own schedule, independent of packet
traffic.

Avalanche-specific work here is deployment and configuration.

---

## 6. C-Chain → Cosmos: the Warp light client

### 6.1 Trust anchor — the pinned validator-set snapshot

The client stores a canonical snapshot of Avalanche's primary-network validators at a fixed P-Chain
height:

- the validators **with** registered BLS keys, in canonical P-Chain order — each a 48-byte
  compressed G1 key plus a stake weight;
- `total_weight`, the snapshot's whole stake **including keyless validators** (I3);
- the P-Chain height the snapshot was taken at (I8).

Mainnet is ~600 validators, ≈33 KB of client state; the stored set is capped at 4096 entries as a
bound against a hostile rotation payload.

Three details are security-critical:

- **Keyless stake dilutes the quorum** rather than disappearing from it (I3).
- **Canonical order is trusted as given; an ordering mismatch fails closed.** The signer bitset
  indexes directly into the stored order, so a stored order that disagrees with the order the
  signatures were collected against selects the wrong keys and the aggregate fails to verify — a
  liveness failure, not a false accept.
- **The membership and weights themselves are a trust-root assumption.** A snapshot that lists the
  wrong validators or the wrong stake is treated as canonical by the client, so its correctness rests
  on the bootstrap/governance path below, not on any check the client can perform.

The snapshot is canonical *because the bootstrap/governance path selected it* — the client does not
yet derive or rotate it trustlessly from Avalanche protocol state:

```text
Avalanche P-Chain snapshot → bootstrap / host-chain governance → pinned set → Warp verification
```

### 6.2 Update message

| Field | Purpose |
|---|---|
| `header` | Canonical Coreth execution header (Coreth adds `ExtDataHash`, `BlockGasCost`, `settledHeight`, … so its RLP differs from an Ethereum header) |
| `signer_bit_set` | Which pinned validators signed, big-endian, one bit per canonical index |
| `signature` | 96-byte aggregate BLS12-381 signature |
| `router_proof` | Merkle-Patricia account proof of the ICS26Router under the header's state root |

The signed Warp bytes are never transmitted. The client rebuilds them from its own configuration and
the locally derived block hash, in AvalancheGo's linear codec:

```text
unsigned message = codec(2)=0 ‖ network_id(4) ‖ source_chain_id(32) ‖ payload_len(4) ‖ payload
payload          = codec(2)=0 ‖ type_id(4)=0 (Hash) ‖ block_hash(32)
```

Because `network_id` and `source_chain_id` come from trusted client state, cross-network and
cross-chain replay are rejected by construction (I4).

### 6.3 Update verification algorithm

Given `U = (header, signer_bit_set, signature, router_proof)`:

```text
Shape
  1. Validate header form (Coreth fork), height bounds, signature length,
     bitset encoding, and proof-size limits.

Identity
  2. block_hash = CorethHash(header)                  — derived, never trusted as input.

Authority
  3. warp_message = WarpMessage(client.network_id, client.source_chain_id, block_hash)
  4. Parse signer_bit_set against the canonical set length; reject non-minimal
     encodings and any bit beyond the set.
  5. signed_stake = Σ weight[i] for each selected i   — overflow-checked.
  6. Require signed_stake · quorum_den ≥ total_weight · quorum_num.
  7. Aggregate exactly the selected public keys and verify the BLS signature
     over warp_message.

State
  8. Verify router_proof against header.state_root; derive the authenticated
     ICS26Router storage root that packet proofs will use.

Transition
  9. Apply transition rules: monotonic frontier, parent linkage with stored
     adjacent heights, idempotent re-submission.
 10. Same height already stored?  identical → no-op;  conflicting → misbehaviour.
 11. Commit atomically (all fallible encoding done before the first write).
```

The ordering is the point:

> **Authentication precedes state commitment.** A header does not become trusted because its hash and
> its Merkle proof are internally consistent — only because a stake-weighted quorum signed it.

Steps 3–7 mirror AvalancheGo's own algorithm (`vms/platformvm/warp/signature.go`), reimplemented in
Rust.

### 6.4 Stored consensus state

Per authenticated block: state root, router storage root, timestamp, height, block hash, parent hash,
first-accepted time. The shape matches the other EVM light clients in this codebase, so packet
verification and host semantics are identical across them.

### 6.5 State transitions

Every update is fully validated before the first persistent write:

```text
untrusted update → verify ──failure──→ no state change
                      │
                      └──success──→ authenticated transition → commit
```

Re-submitting already-stored state is idempotent (first write wins), and adjacent stored heights must
chain by parent hash.

### 6.6 Misbehaviour

On Avalanche, acceptance is finality: no reorgs, no challenge window. Two **independently verified**
headers that disagree at one height therefore mean the validator set signed two histories — the
client freezes, and every later update and proof query fails (I5). Malformed or unverifiable headers
are rejected, never treated as evidence.

The same property sets the client's timing: a verified header is usable immediately, with no
optimistic delay of the kind an OP-Stack-style client needs.

### 6.7 BLS verification on the host chain

The client verifies the aggregate through a single custom query whose wire shape is identical to this
repo's Ethereum light client (`aggregate_verify { public_keys, message, signature }`). Avalanche Warp
and Ethereum consensus share the BLS ciphersuite and domain-separation tag, so a Cosmos chain that
already hosts the Ethereum client verifies Warp aggregates **with no chain-side change**; any other
host routes one handler.

---

## 7. C-Chain state roots under asynchronous execution

Avalanche's asynchronous execution (ACP-194, Helicon) separates acceptance from execution:

> **`header(N).stateRoot` is not the post-state of executing block N** — it commits the state executed
> through `settledHeight(N)`, some blocks behind.

```text
accepted block N ─┬─ block hash ──► Warp attestation
                  └─ header.stateRoot ──► state executed through settledHeight(N) ──► ICS26Router
```

Consequence: a proof fetched at height N does not verify against header N's root. The **relayer**
resolves the pairing — it reads `settledHeight` from the authenticated header itself, not from an
upgrade schedule or an assumed network version, and fetches every paired proof at that height. A
commitment that lands beyond the settled frontier is treated as temporarily unavailable and retried
once the frontier advances.

The split of responsibility matters and does not weaken I7: **choosing the right proof height is a
relayer responsibility; accepting a proof remains an on-chain safety check.** The client never
interprets `settledHeight` — it verifies the router proof against `header.state_root` (I2), so a
proof fetched at the wrong height simply fails to verify. A relayer that gets the pairing wrong
stalls progress; it cannot make a mismatched proof accepted.

This mapping was established empirically against a live network — probing candidate heights N-6…N,
only `settledHeight(N)` matches the header's root — so the path is already correct across the mainnet
upgrade instead of needing a migration when it lands.

---

## 8. What a Warp attestation proves — and what it does not

The verification chain:

```text
C-Chain header → local Coreth hash reconstruction → domain-separated Warp message
   → signer bitset over the pinned set → stake-weighted quorum → BLS aggregate verification
   → authenticated header → router/state proof → authenticated IBC state
```

| Property | Cryptographically verified by the Cosmos client? | Basis |
|---|---:|---|
| Header hashes to `block_hash` | Yes | Local Coreth hash reconstruction |
| Signers belong to the pinned set | Yes | Canonical set + bitset |
| Signed stake reaches the threshold | Yes | Integer quorum arithmetic |
| Aggregate signature is valid | Yes | On-chain BLS verification |
| Router proof belongs to the authenticated state root | Yes | MPT verification |
| Signatures are bound to this Avalanche network and chain | Yes | Reconstructed Warp message |
| Validators sign only blocks meeting the required acceptance semantics | **No — protocol assumption** | Avalanche Warp signing behaviour |
| The Warp quorum equals a portable Snowman consensus certificate | **No** | Explicitly not claimed |

The open question behind row 7 — *under exactly what conditions does Avalanche validator software
sign a C-Chain block-hash Warp message?* — must be traced against the production Avalanche
implementation before this assumption is treated as a settled security property. It is carried as
open item 1 in §15.

Where this sits: unlike an application-specific bridge committee, this design derives its signer
authority from Avalanche's primary-network validator set, weighted by stake, against the configured
stake-weighted Warp threshold. It does not provide the security property of a slashing-backed
consensus proof: no mechanism in this protocol punishes a validator for a false attestation. If
Avalanche later exposes an acceptance or consensus proof (ACP-75, Simplex), it replaces steps 3–7 of
§6.3 inside the same client and nothing else changes.

---

## 9. Validator-set rotation

### 9.1 Today — governance-gated

The snapshot is rotated by a host-chain governance migration supplying a newer snapshot, with the
client enforcing `new.p_chain_height > current.p_chain_height` (I8), so an old set cannot be replayed
as a newer one.

Governance-controlled rotation is therefore a **top-level security assumption**. Between rotations,
Avalanche → Cosmos security depends on the pinned snapshot being correct, on governance selecting the
correct next snapshot, and on the migration preserving canonical order and stake accounting.

### 9.2 Target — protocol-attested rotation

The intended design removes governance from routine set selection:

```text
current set S_t ──stake-weighted Warp quorum──► authenticated description of S_(t+1)
     ──epoch / activation rule──► pending set ──delayed activation──► active set S_(t+1)
```

To be specified before implementation: the P-Chain epoch/snapshot identifier (ACP-181), the
authenticated representation of the next set, canonical ordering, which set is authorized to attest
the transition, activation timing, overlap and continuity rules, rollback and replay protection, and
the failure behaviour when the next set cannot reach quorum.

Until that ships, rotation must not be described as trustless.

---

## 10. Trust model

**Verified, every update:** §4 (I1–I8) and the algorithm in §6.3.

**Assumed:**

1. Avalanche validators issue the relevant Warp signature only under the acceptance semantics this
   client relies on (§8, row 7).
2. The configured 67/100 stake quorum represents the intended Warp security threshold.
3. The pinned snapshot corresponds to the canonical primary-network set at its snapshot height.
4. Host-chain governance rotates the snapshot honestly, until §9.2 ships.
5. The CosmWasm client, the host's BLS verification, and the proof-verification code are correct.

No mechanism in this light-client protocol independently punishes a validator for a false Warp
attestation, so the system must not be presented as equivalent to a slashing-backed consensus proof.

**Relayer and aggregator (I7).** They may omit signatures, return incomplete aggregates, delay
updates, supply malformed messages, or submit stale headers. Each of those is rejected by the
pipeline:

```text
relayer / aggregator compromise ─┬─ can break liveness
                                 └─ cannot forge authenticated state
```

---

## 11. Signature collection

Warp signatures are served over Avalanche's validator networking stack (ACP-118), not over public
EVM JSON-RPC — public API nodes do not expose `warp_*` methods. The reference relayer therefore runs
Ava Labs' **signature-aggregator as a sidecar** and speaks HTTP to it:

```text
relayer ──unsigned message──► aggregator ──ACP-118──► validators
   ◄──aggregate + bitset──────────┘
relayer: re-derive the unsigned bytes, reject any answer over different bytes, build the update
```

Three deliberate consequences:

- Avalanche's node stack (AvalancheGo, TLS p2p identity) stays out of the relayer's dependency tree.
- Aggregation is stateless and free to retry: signatures are never cached and the builder
  re-aggregates on each attempt.
- The aggregator is transport, not a trusted party — its output is re-verified on-chain, and the
  relayer itself refuses an answer produced over unsigned bytes it did not construct.

A failed or partial aggregation advances nothing; the update stays retryable.

---

## 12. Operating the connection

| Component | Where it runs | Notes |
|---|---|---|
| IBC contracts + Spectre client | Avalanche C-Chain | Same EVM-side architecture as the existing stack |
| Avalanche Warp light client | Cosmos chain | ICS-08 / 08-wasm, stored and instantiated through governance |
| Relayer | Operator infrastructure | Reference deployment: one process per pairing; protocol stays permissionless |
| Signature aggregator | Beside the relayer | Needs outbound reach to validator p2p |
| Groth16 prover | Operator infrastructure | Cosmos → C-Chain leg only |

Bootstrap state — canonical validator snapshot, bootstrap header, router storage root at the settled
height, and the Avalanche network identity — is generated from live chain data and validated against
the chain's EVM chain id before instantiation, rather than maintained by hand.

**Observed on Fuji** (2026-09-21; 6 attempts, one host — indicative, not a benchmark): every attempt
reached a verified 67% aggregate in ~2 s, with 12 signers covering 67.13% of stake out of a 73-entry
canonical set. A same-day crawl found 75 validators on Fuji and 597 on mainnet, all with registered
BLS keys.

---

## 13. Implementation status

| Area | Status |
|---|---|
| Cosmos → C-Chain contracts / Spectre / relayer | Implemented |
| Warp message construction and signer-bitset parsing | Implemented |
| Stake-weighted quorum + BLS verification | Implemented |
| Coreth header hash reconstruction | Implemented |
| Router / state proof verification | Implemented |
| Misbehaviour freeze | Implemented |
| Relayer Warp path + signature-aggregator integration | Implemented |
| Bootstrap-state generation | Implemented |
| `settledHeight` proof pairing | Implemented |
| Validator-set rotation | Governance-controlled, monotonic (§9.1) |
| Protocol-attested rotation | Designed, not implemented (§9.2) |
| Full live-network end-to-end run | Not yet performed |
| External audit of the Avalanche client | Not yet performed |
| ZK-wrapped Warp verification (hosts without cheap BLS) | Future |

The implemented rows are covered by unit and cross-language tests running against live-network
captures (a real Fuji aggregate, header, and 73-validator bootstrap).

---

## 14. Code map

| Path | Responsibility |
|---|---|
| `packages/avalanche-warp/` | Warp message construction, signer bitset, quorum arithmetic, BLS verification abstraction |
| `packages/avalanche-light-client/` | Client state, update verification, transitions, validator-set handling |
| `programs/cw-ics08-wasm-avalanche/` | The ICS-08 / 08-wasm contract artifact |
| `relayer/chain/l2rollup/warp_message.go` | Warp message encoding, mirrored byte-for-byte with the Rust side |
| `relayer/chain/l2rollup/header_warp.go` | Warp-backed header construction and aggregator client |
| `relayer/chain/l2rollup/coreth_header.go` | Coreth header encoding/hash and the settled-height mapping |
| `tools/warp-spike/` | Aggregation measurement and bootstrap-state generation |
| `docs/E2E.md` | Cosmos ↔ Avalanche bring-up and integration runbook |

---

## 15. Open items

1. The exact protocol conditions under which an Avalanche validator signs a C-Chain block-hash Warp
   message (§8, row 7).
2. The justification for treating the configured Warp quorum as sufficient evidence for this client's
   acceptance model.
3. Protocol-attested validator-set rotation (§9.2).
4. Full end-to-end live-network validation.
5. External security audit.
6. Gas and resource benchmarks for worst-case validator-set and proof sizes.
7. Liveness behaviour under sustained partial signature availability.
8. Upgrade compatibility across future Coreth header and execution-model changes.

---

## References

- Avalanche Warp Messaging / ICM documentation
- ACP-118 (signature requests), ACP-181 (epoched validator views), ACP-194 (asynchronous execution)
- AvalancheGo `vms/platformvm/warp/signature.go` — the reference verification algorithm
- Coreth header implementation
- IBC v2 specification; ICS-08 / 08-wasm light-client interface
