# Arbitrum Attestor

The attestor is a standalone gRPC backend for the relayer. Like the OP
attestor, it publishes a proposal-independent root from a configured Nitro
`unsafe`, `safe`, or `finalized` head. It also ingests BoLD v2 Arbitrum
assertions from finalized Ethereum L1 blocks by default and checks each assertion's L2
block against configured Nitro HTTP and WebSocket endpoints. It never starts,
supervises, or stops a Nitro process.

The configured Nitro endpoint performs execution and owns all chain state. The
attestor stores only its RollupCore log cursor, unresolved assertions, verified
frontier entries, and mismatch records; it does not maintain an execution
database.

- To configure, run and operate it: the **[Operator Guide](GUIDE.md)**.
- The signed update contract and finality tiers of the L2 clients that consume this
  attestor: [docs/runbooks/L2_CLIENTS.md](../../docs/runbooks/L2_CLIENTS.md).

## Trust boundary

For `source="assertion"` entries, the attested root comes from the assertion
read out of `RollupCore` on Ethereum L1. Nitro supplies the runtime head and
canonical-block view used to decide when that assertion can enter or advance in
the feed. The attestor re-reads assertion status at finalized L1 once that
snapshot reaches the assertion's creation block. Opt-in early ingestion also
rebuilds the complete unfinalized L1 window and withdraws orphaned assertions.
Consequently, an untrusted Nitro endpoint cannot fabricate an assertion-backed
root; it can delay, suppress, or reject an otherwise valid assertion by
misreporting its head or canonical view, which is a liveness failure rather
than a new L1 commitment.

`source="derived"` is deliberately different and mirrors OP's configured-head
feed: its root is read directly from Nitro. `attestation_head="unsafe"` trusts
the configured endpoint immediately; `safe` and `finalized` select Nitro's
corresponding L1-derived heads. Set `disable_derived_roots=true` when the feed
must contain only RollupCore-backed assertions.

The L2 client does not verify the chain-specific assertion itself. It accepts a
canonical L2 execution header only with this daemon's Ed25519 signature and
verifies the router account proof against the signed header state root. Thus
`source` and `attestation_head` govern what this daemon will sign; the
permissionless Cosmos transaction sender cannot fabricate the signed block
identity. See `docs/runbooks/L2_CLIENTS.md` for the exact signature format.

## gRPC API

The protobuf contract is in
`../../proto/attestor/attestor.proto`, with generated Go client code in
`../types/attestor`. Regenerate it from the repository root with
`buf generate --template buf.gen.attestor.yaml`. The reusable Nitro runtime state is in `runtime_state.go`;
the gRPC transport implementation lives in the shared `attestor/grpc` and `attestor/host` packages. The service exposes:

```text
attestor.AttestorService/Info
attestor.AttestorService/AttestedUpTo
attestor.AttestorService/AttestedRootAtOrBelow
attestor.AttestorService/WatchAttested
attestor.AttestorService/VerifyStateRoot
```

A failed call returns its gRPC code and a classified message only (for example
`query Nitro finalized head`). The upstream error behind it is not returned,
because it quotes the Nitro or L1 URL and an RPC provider's API key can be part
of that URL; it is in the attestor's log with URLs reduced to scheme and host
(see the [runbook](../../docs/runbooks/ATTESTOR.md#error-redaction)).

`AttestedUpTo(src_chain, include_provisional)` returns the highest qualifying
verified commitment. `AttestedRootAtOrBelow` applies an inclusive L2-height
bound so the relayer's header builder can resolve the assertion provenance for
a target height. A response with `found=false` means the attestor has not
accepted a qualifying assertion yet.

For Arbitrum, `root` is the Nitro L2 state root. A proposal-independent entry
has `source="derived"` and no provenance object. When `attestation_head` is
`unsafe` or `safe`, it is provisional until Nitro's finalized head reaches the
same height. The attestor then recomputes the commitment, confirms it on a
match, or replaces it with the finalized block hash/state root and emits a
`HEAD DIVERGENCE` log on a mismatch.

An assertion-backed entry has `source="assertion"` and carries its BoLD
identifier in `assertion_hash`. A pending assertion enters the feed as
`provisional=true` only after its block is canonical under Nitro's safe head.
It becomes non-provisional only after RollupCore reports it confirmed in
finalized L1 state and the same block is canonical under Nitro's finalized
head. Rejected, challenged, or locally mismatched assertions are never
returned.

The chain-specific provenance is a protobuf `oneof`: OP entries carry
`game_index`, Arbitrum entries carry `assertion_hash`, and derived entries may
carry neither.

The relayer calls `VerifyStateRoot` for every header it builds (see
[Operations](GUIDE.md#operations)). The caller supplies:

- `block_number`
- `expected_state_root` as exactly 32 bytes
- optional `expected_block_hash` as exactly 32 bytes
- `run_mode` as `RUN_MODE_UNSAFE`, `RUN_MODE_SAFE`, or `RUN_MODE_FINALIZED`

The response returns `valid` plus Nitro's canonical block hash and state root.
A block-identity mismatch is a successful gRPC response with `valid=false`; malformed requests,
missing blocks, and Nitro availability failures use gRPC status errors.

Every verification request must name the daemon's configured attestation head exactly.
The service rejects an unspecified or different `run_mode`, and rejects a requested block above
that configured head. Deploy separate daemon/key configurations for different finality tiers.

In the background, the attestor subscribes to Nitro's `newHeads` stream over
the configured WebSocket endpoint. Each unsafe-head event triggers a refresh
of the unsafe, safe, and finalized tags followed by a derived-root attestation
pass. The attestor reconnects and reconciles after a subscription failure, and
retains `runtime_poll_interval` as a fallback because safe or finalized can
advance without a new unsafe block.

The attestor records the first unsafe and safe root observed at each height
until it finalizes. When Nitro advances the finalized head, the runtime compares
the finalized root with both earlier observations and logs a consistency
result. This in-memory result is intended as the source for future alerting;
Nitro remains the persistent owner of chain state.

The standard gRPC health service is also registered for readiness checks. The
route is ready once Nitro has reported the configured `attestation_head` and the
assertion loop is making progress: it has come back around, or one of its calls
to L1 or Nitro has returned, within `assertion_call_timeout` plus two
`assertion_poll_interval`s. A loop whose passes fail still counts as running:
those failures are in the log and in `assertion_sync_failures_total`. So does a
pass that runs for a long time over a long list, because it makes progress with
every call. A loop parked on something that never returns does not, and neither
does a feed that is frozen behind it.

## Assertion feed flow

On every assertion poll, the attestor:

1. Reads the BoLD `AssertionCreated`/`AssertionConfirmed` lifecycle only
   through Ethereum's finalized L1 head and advances a persisted cursor in
   bounded ranges.
2. Recomputes every BoLD assertion hash from its event data.
3. Resolves the asserted L2 block by hash in Nitro, then checks the canonical
   block at the recovered height. Pending assertions are checked against the
   safe head.
4. Re-reads `AssertionNode.status` at finalized L1. Confirmed assertions are
   promoted only after the Nitro finalized head covers and matches the block.
5. Serves the resulting frontier through the same oracle-shaped RPCs used by
   the OP attestor.

## Development

```sh
just build-arbitrum-attestor   # binary at attestor/arbitrum/bin/attestor
```
