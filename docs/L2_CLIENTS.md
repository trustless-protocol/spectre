# Attestor-trusted L2 ICS-08 clients

Spectre builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum`;
- `cw-ics08-wasm-base`; and
- `cw-ics08-wasm-op`.

The artifacts share one client lifecycle and differ only in their deployment-profile version.

> **DEVNET ONLY — these clients trust ANY SUBMITTER, not just the relayer.**
>
> `MsgUpdateClient` is permissionless and the header carries no attestor signature, so the client
> cannot tell the relayer's headers from anyone else's. Anyone can build a self-consistent
> `AttestedL2Header` at a height the client has not seen — a fabricated state root with a matching
> router account proof is cheap, because nothing ties either to the real L2 — get it accepted, and
> then prove arbitrary membership or non-membership against it. That is enough to mint tokens on the
> Cosmos side or to time out packets that were in fact delivered.
>
> The clients verify no attestor signature, no L1 consensus, no OP dispute game and no Arbitrum
> assertion. The router account proof establishes only that the supplied router storage root belongs
> to the supplied L2 execution state — it says nothing about whether that state is the L2's.
>
> The relayer does gate its own submissions (`Source.RelayableHeight` bounds the height from the
> attestor frontier, and the header builder asks the attestor's `VerifyStateRoot` whether the block
> it packaged is canonical), but that is a relayer-side check against accidental divergence. The
> client cannot re-check it, so it constrains an honest relayer, not an attacker.
>
> The attestor-only redesign (PR #343, not yet merged) states the target invariant — nothing
> unauthenticated may ever enter the client — and this interim format does not meet it. **Do not
> deploy these artifacts on a network holding real value.** The gate that closes it is
> `ATTESTATIONS_ARE_AUTHENTICATED` in `packages/l2-client/src/runtime.rs`, flipped in the same
> change that adds signature verification and the signed-attestation wire version.

## Live 08-wasm surface

Each artifact exports only the entry points used by the current client lifecycle:

- `instantiate` stores explicitly supplied client and consensus state;
- `sudo` handles `update_state`, `update_state_on_misbehaviour`, `verify_membership`, and
  `verify_non_membership`; and
- `query` handles `verify_client_message`, `check_for_misbehaviour`, `timestamp_at_height`, and
  `status`.

There is no `execute` entry point because an ICS-08 client is driven by the host, not by ordinary
contract messages. There is currently no contract `migrate` entry point and no implementation of
the optional IBC client-upgrade or substitute-client recovery sudo messages. Deployments that need
those governance paths must implement and review them separately; ordinary update and packet relay
does not call them.

Light-client responses must contain data only. ibc-go's 08-wasm keeper rejects responses carrying
attributes, events, or messages, so the shared entry points intentionally return
`Response::default().set_data(data)`.

## Creation and runtime profile

Creation uses the standard direct byte fields:

```text
InstantiateMsg { client_state, consensus_state, checksum }
```

The decoded client state contains `latest_height`, an optional `frozen_height`, and one immutable
artifact profile. A profile contains only:

```json
{
  "common": {
    "l2_chain_id": 11155420,
    "l2_router": "0x645280885749dc97ea461de280eb3273c91d36df",
    "commitment_slot": "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600",
    "profile_version": "op_attestor_v1",
    "l2_header_fork": "prague"
  }
}
```

The expected profile versions are `op_attestor_v1`, `base_attestor_v1`, and
`arbitrum_attestor_v1`. No pinned Ethereum client, beacon slot, game factory, RollupCore contract,
finality policy, or settlement proof belongs in new client state.

Creation validates the profile version, revision-zero nonzero height, nonzero roots and block hash,
and equality between the client latest height and bootstrap consensus height. It performs no host
or L1 query.

## Client update wire contract

`UpdateState` and `VerifyClientMessage` consume this JSON-compatible shape inside the tagged
`ClientMessage` envelope:

```text
ClientMessage::Header(AttestedL2Header {
    l2_header: CanonicalEvmHeader,
    router_proof: EvmAccountProof
})
```

The client validates the configured canonical-header fork, derives the L2 block hash, verifies the
configured `ICS26Router` account proof against the header state root, and stores the resulting
consensus roots and block identity. Attestor provenance and authentication are not carried in this
wire version.

Updates may advance the latest height or backfill an absent historical height. An identical update
is idempotent. A different block/state/router identity at the same height is rejected as a
conflict. Because updates currently carry no authenticated attestation, conflicting headers do not
constitute actionable misbehaviour; every host validation path refuses them instead of freezing the
client. The `frozen_height` field remains reserved for the signed-attestation protocol and recovery
path.

## Relayer and attestor integration status

The Go `relayer/chain/l2rollup` builder emits exactly the wire format above: the canonical L2
execution header at the requested height plus the router account proof, and nothing else. The
per-chain settlement builders it replaced — the OP one proving a DisputeGameFactory game against an
authenticated L1 state root, the Arbitrum one proving a RollupCore/BoLD assertion — were deleted
rather than kept behind a legacy path, because the clients they fed no longer exist either.

Because the client verifies no L1 object, the attestor is the entire trust boundary and it has to
bound two separate things. `Source.RelayableHeight` bounds how far the relayer may advance, from
`AttestedUpTo`. The builder then calls the attestor's `VerifyStateRoot` with the state root and
block hash it is about to package, at the run mode matching the configured head kind, and refuses a
block the attestor's replica does not hold as canonical at that height — otherwise an L2 RPC that
reorged past the frontier would supply a replacement block at an approved height and the client
would accept it.

That binding is best-effort by construction: nothing in the wire format lets the client re-check the
answer, so it defends against divergence between the relayer's L2 RPC and the attestor, not against
a relayer that simply skips the call. Both in-tree attestors implement `VerifyStateRoot`; an
attestor binary older than that answers `Unimplemented`, and the relayer degrades to the unbound
path with a one-per-process warning rather than refusing to relay on a version skew.

The attestor `AttestedRoot` protobuf supplies `l2_block_number`, `root`, `source`, optional
game/assertion provenance, `provisional`, and `attested_at`. Those fields gate which execution
header the relayer selects; none of them is copied into `AttestedL2Header`. `root` in particular is
chain-specific — an OP output root but an Arbitrum L2 state root — which is why the binding goes
through `VerifyStateRoot` rather than comparing it directly. OP exposes its configured
`attestation_head` separately through `Info`; Arbitrum leaves that field unset and remains
assertion-gated.

Remaining integration work:

1. implement `VerifyStateRoot` in the OP-Stack attestor, after which the degraded path above stops
   being reachable;
2. add a cross-language fixture that the Go encoder and Rust decoder both accept; and
3. rebuild `relayer/l2fixtures`, which still captures settlement evidence for the Rust verifier
   fixtures.

The signed attestation fields and verification specified by the redesign return in the next wire
version. Until that protocol lands, this unsigned bring-up format cannot supply actionable
misbehaviour evidence.

## Packet proofs

Membership and non-membership load the exact stored L2 consensus height and verify a bare
`EvmStorageProof` against its authenticated router storage root and configured commitment mapping
slot. No Ethereum-client query or finality threshold is involved.

The proof wire shape is:

```text
EvmStorageProof { key: bytes32, value: byte array, proof: array<byte array> }
```

The relayer must not send the L1 client's combined account-and-storage proof shape. The router
account proof was already checked during the client update.

## Validation

```bash
cargo test --locked \
  -p l2-client -p op-verifier -p base-verifier \
  -p arbitrum-verifier -p cw-ics08-wasm-op -p cw-ics08-wasm-base \
  -p cw-ics08-wasm-arbitrum

cargo build --target wasm32-unknown-unknown --release --locked \
  -p cw-ics08-wasm-op -p cw-ics08-wasm-base -p cw-ics08-wasm-arbitrum
```
