# OP Stack Attestor

The attestor answers one question for the relayer: **is this OP Stack chain's
state legit, so packets can be relayed against it?** It is the OP Stack
counterpart of our Groth16 prover for Cosmos sources — an **independent
verifier the relay path consumes**, running as its own binary with its own
lifecycle. It drives an operator-owned **verify-mode replica** (op-node +
op-reth) that re-derives the L2 chain from data posted on Ethereum L1, and
**byte-compares** output roots against that self-computed state — replay,
not trust.

- To deploy, configure, and operate it: the **[Operator Guide](GUIDE.md)**.
- Design rationale, trust model, and the OP Stack IBC roadmap: the Ethereum
  L2 architecture doc (upcoming; local draft `OPTIMISM_IBC_DESIGN.md`).

## Architecture

```
[L1 node] ← derives ← [op-node + op-reth replica] ← RPC ← [attestor] ← gRPC ← [relayer]
   batches/blobs        verify mode, operator-owned        feed + metrics      (op_to_cosmos, future)
```

One attestation pass: ingest new factory games (exact index cursor) → read
replica heads → resolve provisional rechecks the finalized head now covers →
decide pending games → self-attest a derived root at the gating head.

```
attestor/
├── attestor.go       Attestor interface + failure backoff
├── opstack/          the OP Stack implementation
│   ├── opstack.go    attestation loop (ingest → heads → rechecks → verdicts → derived)
│   ├── games.go      DisputeGameFactory ingest (exact index cursor; l2SequenceNumber
│   │                 with l2BlockNumber fallback across the interface rename)
│   ├── replica.go    op-node client (optimism_syncStatus / optimism_outputAtBlock)
│   ├── store.go      persistent state: cursor, pending, rechecks, attested feed
│   ├── challenge.go  mismatch hook (log-only placeholder for the on-chain challenger)
│   └── config.go     op_source config + validation
├── server/           read-only gRPC service over the running attestors
├── client/           typed Go client mirroring the store's read interface
├── cmd/              the `attestor` binary (own config loader, metrics, gRPC)
├── bindings/         GENERATED abigen bindings — never hand-edit
└── types/attestor/   GENERATED buf output of proto/attestor/attestor.proto
```

Its own Go module, **no cgo** — builds and tests without `libgaraga_rs.so`.
The gRPC proto (`../proto/attestor/attestor.proto`, generated via
`../buf.gen.attestor.yaml`) is the **only contract** with the relayer.

## gRPC API (`attestor.AttestorService`)

The attested-root feed it serves holds two entry types: **derived roots**
(`source: "derived"` — the replica's own output at the gating head; the feed
the relay path consumes) and **game-verified roots** (`source: "game"` —
`DisputeGameFactory` proposals byte-compared against honest derivation).
Provisional semantics are defined in the Operator Guide's
[trust levels](GUIDE.md#trust-levels).

| Method | Purpose |
|---|---|
| `Info` | Per-chain status: head config, replica heads, pending count, confirmed + provisional frontiers |
| `AttestedUpTo` | The attested frontier (`include_provisional` opt-in) |
| `AttestedRootAtOrBelow` | The relay-path query: highest attested root covering a packet's height |
| `WatchAttested` | Server-streamed frontier advances (streams the frontier, not every entry) |

Go consumers use `attestor/client` (`client.Dial(addr)`), which mirrors the
store's read interface — in-process and sidecar consumption are
interchangeable.

## Development

```bash
just build-attestor           # binary at attestor/attestor
just test-attestor            # go test -race -count=1 ./... — no network needed
```

Unit tests cover the verdict state machine, cursor rules,
provisional→confirmed plus both divergence directions, derived
gap/reverify/prune, store persistence/corruption, bootstrap binary search,
and op-node JSON-RPC marshalling — all with fakes, no network.
