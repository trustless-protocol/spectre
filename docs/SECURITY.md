# Security

## Security Model

### Trust Assumptions

1. **Validator set**: At least 2/3 of voting power is honest (standard BFT assumption)
2. **Groth16 circuit**: BatchCircuit correctly verifies a batched Ed25519 commit
   (eddsa.VerifyBatchWithMsgBytes) and binds its full witness via SHA-256
3. **Admin keys**: UUPS upgrades require admin access via AccessManager
4. **Relayer**: Untrusted — relayer cannot forge proofs, only relay valid state transitions
5. **Trusted setup**: Each bucket has its own `(pk, vk)` from `groth16.Setup`.
   Compromise of the per-bucket toxic waste would let an attacker forge
   batch proofs for that bucket size; treat artifact regeneration as a
   security event and redeploy verifiers atomically.

### Proof Verification Chain

```
Cosmos validators sign block (Ed25519 over CanonicalVote bytes)
    → Relayer batches a 2/3-quorum subset, pads to bucket with deterministic
      dummy keypairs, and produces one Groth16 proof
    → SpectreClient sums voting power over UNIQUE active signers
      (skips padding via active flag, rejects duplicate active indices) and
      requires `power * 3 > totalVotingPower * 2`
    → SignatureVerifier recomputes the witness commitment from calldata
      (PrefixHead ‖ roundPresent ‖ blockHash ‖ per-slot active ‖ pubkey —
      one SHA-256 over the whole batch, exposed as two 128-bit public
      inputs) and dispatches to Groth16Verifier_N{bucket}
    → Per-bucket Groth16Verifier checks the proof against that bucket's VK
    → ICS-23 Merkle proofs verify packet commitments against the new app hash
```

### Attack Surfaces

| Surface | Mitigation |
|---------|-----------|
| Forged headers | Groth16 batch proof + 2/3 unique-signer voting-power quorum |
| Calldata pubkey swap | Pubkey A bound into the SHA-256 witness commit; tampering changes the public input and breaks proof |
| Padding inflating quorum | `active=false` slots skipped on-chain; active byte hashed so calldata can't toggle it |
| Duplicate signer in calldata | Sorted-index invariant (`idx > prevIdx`) rejects an active validator index appearing twice |
| Bucket spoofing | `SignatureVerifier.buckets[bucket]` lookup; only owner can `setBucket`; verifier dispatch matches the circuit each `(pk, vk)` was built for |
| Validator-set spoofing (pinned set) | The pinned validator set is written only by the constructor and `updateConsensusState`; re-pinning requires the full header + batch-proof + >2/3 quorum verification AND `hashValSet(newValSet) == header.nextValidatorsHash`. `updateApplicationState` never re-stores validator data, so the frequent per-packet path cannot swap the pinned set |
| Stale pinned set | The app-state path may intentionally use a pinned set whose hash differs from the trusted height's `nextValidatorsHash`, but the pinned snapshot stores the consensus timestamp at which it was pinned and must remain inside the trusting period for every update |
| Stale-revision cache (chain ID vs latestHeight drift) | Constructor invariant: `ChainId.get(clientState.chainId).revisionNumber == clientState.latestHeight.revisionNumber`. `updateApplicationState`/`updateConsensusState`/`misbehaviour` rely on `latestHeight.revisionNumber` instead of re-parsing `chainId` per call |
| BFT-time non-monotonicity evidence | Standalone misbehaviour intentionally freezes only for same-height equivocation: two valid, proof-backed headers at the same height with different block hashes. Cross-height BFT-time non-monotonicity evidence is outside this client's freeze scope. |
| Replay attacks | Packet sequencing in ICS26Router |
| Double-spend | Commitment storage in ICS24Host |
| Validator equivocation | Misbehaviour detection → client freeze |
| Unbounded token minting | Rate limiting (RateLimitUpgradeable) |
| Malicious upgrade | AccessManager RBAC + UUPS auth |

### Permissionless Client Isolation

`ICS02Client.addClient(counterpartyInfo, client)` is intentionally permissionless when the
caller uses the auto-generated `client-N` identifier path. Anyone can register an arbitrary
light-client contract, so the safety boundary is fund isolation rather than client admission.

ICS20 escrows are keyed by `destinationClient`, and voucher denoms are namespaced as
`transfer/<clientId>/<base-denom>`. A rogue client can only affect packets and escrows for
that client identifier, so it can only touch funds users voluntarily route through that
client. It cannot collide with custom client IDs because custom identifiers reject the
reserved `client-` prefix.

### Pause And Replay Semantics

Pausing is asymmetric by design. `sendPacket` and `recvPacket` are paused because they admit
new outbound commitments or inbound application execution. `ackPacket` and `timeoutPacket`
remain callable while paused so in-flight packets can still settle, refund, and release
funds.

Replay handling is also deliberate. Replayed receives, acknowledgements, and timeouts still
verify the relevant proof first, then emit `Noop` instead of reverting if the local packet
state has already been consumed. This lets relayers batch duplicate work without one stale
packet reverting an entire multicall.

### Rate Limit Semantics

Escrow rate limits are net-flow counters over a one-day decay period. Sends out of escrow add
usage; deposits into escrow subtract usage. This allows legitimate inbound flow to restore
outbound capacity, while an attacker must first back any extra outbound capacity with their
own deposited tokens.

Tokens with a rate limit of `0` are untracked. If a limit is enabled later, accounting starts
from zero at that point; previous untracked transfers are not counted retroactively.

### Production Escrow Configuration

`PRODUCTION_ESCROW_CONFIG` supplies `clients`, `tokens`, and `limits` arrays. Each
`tokens[i]`/`limits[i]` pair is installed on every configured client escrow; the format
does not support a distinct launch cap for each client. Client IDs must be nonempty and
unique, and production token limits must be nonzero.

**`RATE_LIMITER_ROLE` is global, not per-escrow.** Escrows are `BeaconProxy` instances
created per client, so a per-target `setTargetFunctionRole` grant cannot exist before the
escrow does — and one covering only the pre-created escrows would leave every later client
uncapped, since an unset limit means *no* limit. `Escrow.setRateLimit` therefore checks the
role on the AccessManager directly, and one role holder governs every escrow, present and
future. It must be granted with an execution delay of **zero**: that path never consumes a
scheduled operation, so a delayed grant is rejected rather than silently ignored.

**Onboarding a client after launch requires escrow activation.** Production deployments call
`ICS20Transfer.enableEscrowLaunchGate()`, after which packet processing will not create an
escrow on demand and will reject a provisioned escrow until it is active. A client therefore
cannot go live uncapped. The required sequence is `createEscrow(clientId)` (ADMIN_ROLE),
`setRateLimit` with non-zero limits for every required token (RATE_LIMITER_ROLE), then
`activateEscrow(clientId, tokens)` (ADMIN_ROLE). The activation step verifies the limits and
the gate cannot be turned off once enabled. Activation is a point-in-time check of only the
listed tokens: setting a limit back to `0` does not deactivate the escrow, and tokens omitted
from the activation list remain uncapped. Existing deployments must activate all live escrows
before enabling the gate because their active state defaults to `false`.

## Static Analysis

```bash
just slither    # Run Slither security scanner
```

Config (`.slither.config.json`):
- Framework: Foundry
- Excludes: informational, low severity, dependencies

## Access Control

Roles defined in `IBCRolesLib.sol`:
- Upgrades require admin (AccessManager)
- Emergency controls via PAUSER/UNPAUSER roles
- Rate limits managed by RATE_LIMITER_ROLE

### Proof Submission (SpectreClient)

`PROOF_SUBMITTER_ROLE` controls who can call `verifyMembership`/`updateApplicationState`/`updateConsensusState`:
- **`roleManager = address(0)`**: Anyone can submit proofs (permissionless, suitable for dev/test). Proofs are still cryptographically verified by Groth16.
- **`roleManager != address(0)`**: Only addresses granted `PROOF_SUBMITTER_ROLE` by the role manager can submit. Use for production to restrict relay to authorized operators.

In permissionless mode (`roleManager = address(0)`), no account receives `DEFAULT_ADMIN_ROLE`. If the client freezes, `unfreeze()` cannot be called; recovery is by migrating the router to a replacement client.

## Known Considerations

- `contracts/verifiers/Groth16Verifier_N{N}.sol` are auto-generated by `relayer/prover/cmd` from each bucket's VK — do not manually edit
- Each `groth16.Setup` run uses fresh randomness; the VK changes every regeneration, so per-bucket verifiers and the SignatureVerifier bucket registry must be redeployed atomically
- Every address registered through `SignatureVerifier.setBucket` must be a gnark-compatible verifier that reverts on invalid proofs. `SignatureVerifier` treats a non-reverting `staticcall` as proof success.
- The executable verifier manifest is intentionally N4-only. Any header whose quorum needs more than four distinct active signers cannot be served. This is a hard launch blocker for normal production validator sets; adding a bucket requires a coordinated circuit setup, paired artifact publication, verifier deployment, `SignatureVerifier` registration, deployment verification, and manifest update. Historical benchmark reports for larger buckets are measurements, not evidence that those buckets are currently published or supported.
- Outbound packet timeouts are capped at 1 day by `ICS26Router.MAX_TIMEOUT_DURATION`; operators should pick default transfer timeouts with expected relayer and counterparty outage windows in mind
- `relayer/go.mod` replace directives point to local paths — verify before building
- `Encode` and `Header` are internal libraries (inlined into contract bytecode), so no link step is needed when stamping the Go-binding bytecode. `contracts/compile.sh` resolves fully qualified manifest entries and fails loudly if an artifact source mismatches or an unlinked library placeholder reappears.
- Pinned-set history grows monotonically — `updateConsensusState` appends one snapshot (`snapshots` + `snapshotHeights` in the ERC-7201 Store) per rotation and each rotation writes a new SSTORE2 blob; nothing is evicted. At the 24 h rotation cadence this is ~365 snapshots/year — modest, but the relayer / operator owns the SSTORE bill, and historical snapshots are load-bearing (misbehaviour proofs resolve against the set pinned at their trusted height). `updateApplicationState` (the frequent per-packet path) reads the pinned set without re-storing it; the relayer must submit `updateConsensusState` before the pinned snapshot exits the trusting period or before signer overlap decays below the >2/3 pinned-quorum floor.
