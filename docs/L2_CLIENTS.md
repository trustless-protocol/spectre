# Attestor-trusted L2 ICS-08 clients

Fast-IBC builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum`;
- `cw-ics08-wasm-base`; and
- `cw-ics08-wasm-op`.

The artifacts share one client lifecycle and differ only in their deployment-profile version.

> **Security assumption:** these bring-up clients trust the relayer to forward data produced by the
> selected attestor. They do not verify an attestor signature, L1 consensus, an OP dispute game, or
> an Arbitrum assertion. The router account proof establishes only that the supplied router storage
> root belongs to the supplied L2 execution state.

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

The current Go `relayer/chain/l2rollup` implementation is still the settlement-proof client for
existing deployments. Its OP builder emits `beacon_slot`, an authenticated L1 state root, factory
and game proofs, an output-root preimage, and the L2 header/router proof. Its Arbitrum builder emits
the corresponding RollupCore/BoLD assertion shape. Those Go builders and their command wiring are
live legacy entry points and have deliberately not been deleted by the L2 Wasm refactor.

Those messages are **not wire-compatible** with the new clients: the Rust message uses strict
unknown-field rejection and accepts only the canonical L2 execution header and router account
proof, while the old Go messages contain settlement fields. A new relayer builder must fetch the
exact attested L2 execution header and router account proof, package the common envelope above, and
commit the header's L2 height. Until that builder exists, the new Wasm artifacts cannot be used by
the current `l2_to_cosmos` module end to end.

The current attestor `AttestedRoot` protobuf supplies `l2_block_number`, `root`, `source`, optional
game/assertion provenance, `provisional`, and `attested_at`. These fields gate which execution
header the relayer selects, but they are not copied into `AttestedL2Header`. OP exposes its
configured `attestation_head` separately through `Info`; Arbitrum currently leaves that field unset
and remains assertion-gated.

The required integration work is therefore:

1. build `AttestedL2Header` in the relayer instead of selecting and proving a settlement object;
2. retain the old Go builders under an explicit legacy path while existing old-client deployments
   still use them; and
3. add a cross-language fixture that the Go encoder and Rust decoder both accept.

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
