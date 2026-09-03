# Attestor-trusted L2 ICS-08 clients

Spectre builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum`;
- `cw-ics08-wasm-base`; and
- `cw-ics08-wasm-op`.

The artifacts share one client lifecycle and differ only in their deployment-profile version.

> **Permissionless submission, authenticated state.** `MsgUpdateClient` remains permissionless:
> anyone may pay to relay an update. The client accepts it only when the `AttestedL2Header` carries
> a valid Ed25519 signature from the immutable `attestor_public_key` in its profile. A fabricated
> L2 header or a header signed for another L2 chain fails before it can create a consensus state.
>
> The attestor signs the domain-separated tuple `l2_chain_id`, immutable `attestation_head`, L2 block
> height, state root and canonical block hash. The client derives that hash from the complete execution header and verifies
> the router account proof against its signed state root, so timestamp, parent hash and every other
> header field are bound too. The relayer is transport only; it has no authority to manufacture a
> valid update.

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
    "profile_version": "op_attestor_v2",
    "l2_header_fork": "prague",
    "attestor_public_key": "0x<32-byte-ed25519-public-key>",
    "attestation_head": "safe"
  }
}
```

The expected profile versions are `op_attestor_v2`, `base_attestor_v2`, and
`arbitrum_attestor_v2`. The public key and `attestation_head` are immutable for the client lifetime; key rotation therefore
uses a new client or a reviewed client-recovery path. No pinned Ethereum client, beacon slot, game
factory, RollupCore contract, or settlement proof belongs in new client state.

Creation validates the profile version, revision-zero nonzero height, nonzero roots and block hash,
and equality between the client latest height and bootstrap consensus height. It performs no host
or L1 query.

## Client update wire contract

`UpdateState` and `VerifyClientMessage` consume this JSON-compatible shape inside the tagged
`ClientMessage` envelope:

```text
ClientMessage::Header(AttestedL2Header {
    l2_header: CanonicalEvmHeader,
    router_proof: EvmAccountProof,
    attestor_signature: [u8; 64]
})
```

The client validates the configured canonical-header fork, derives the L2 block hash, verifies the
signature over that hash, the configured L2 chain ID, and immutable `attestation_head`, then verifies
the configured `ICS26Router` account proof against the header state root before storing the resulting
consensus state. This prevents an unsafe signature made by the same key from satisfying a client
pinned to `safe` or `finalized`.

Updates may advance the latest height or backfill an absent historical height. An identical update
is idempotent. A different block/state/router identity at the same height is authenticated
misbehaviour and freezes the client. Recovery/unfreeze policy remains a separate governance task.

## Relayer and attestor integration status

The Go `relayer/chain/l2rollup` builder emits exactly the wire format above. It asks the attestor to
compare the candidate block identity with its own replica at the profile's `attestation_head`, verifies
the returned signature against the same public key and head pinned in `rollup_profile`, then packages
the canonical L2 execution header, router account proof and signature. CosmWasm independently repeats
that verification.

`Source.RelayableHeight` bounds how far the relayer may advance from `AttestedUpTo`. The builder
then calls `VerifyStateRoot` at the configured head kind and fails closed if the attestor refuses,
is unavailable, is too old to sign, or returns a malformed signature. Thus a hostile relayer cannot
bypass the attestor by skipping its local check: the client verifies the returned signature itself.

The attestor `AttestedRoot` protobuf supplies `l2_block_number`, `root`, `source`, optional
game/assertion provenance, `provisional`, and `attested_at`. Those fields gate which execution
header the relayer selects; none of them is copied into `AttestedL2Header`. `root` in particular is
chain-specific — an OP output root but an Arbitrum L2 state root — which is why the binding goes
through `VerifyStateRoot` rather than comparing it directly. OP exposes its configured
`attestation_head` separately through `Info`; Arbitrum leaves that field unset and remains
assertion-gated.

An attestor key is configured per L2 source. OP and Base use the OP Stack signer configuration;
Arbitrum uses its Nitro/BoLD signer configuration. The key's public half must be copied exactly into
the profile used to create the Cosmos client.

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
