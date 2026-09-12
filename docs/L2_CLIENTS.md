# Authenticated L2 ICS-08 clients

Spectre builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum`;
- `cw-ics08-wasm-base`; and
- `cw-ics08-wasm-op`.

The artifacts share one authenticated client kernel. Their data-only profiles select the chain ID,
router, commitment slot and canonical execution-header fork.

> **Development boundary:** the contracts, relayer and attestor use the authenticated signed wire
> contract described below. The stack remains development-only until the manual release gate,
> operational key-custody review and authenticated OP/Base/Arbitrum E2E all pass against the exact
> candidate diff.

## Live 08-wasm surface

Each L2 artifact exports only the entry points used by the current client lifecycle:

- `instantiate` stores explicitly supplied client and consensus state;
- `sudo` handles update, misbehaviour, membership and non-membership operations; and
- `query` handles client-message verification, misbehaviour checks, timestamps and status; and
- `migrate` validates or atomically replaces the active attestor set during a governed code
  migration.

There is no `execute` entry point. Client upgrade and substitute-client recovery sudo messages
return typed unsupported-operation errors.

Light-client responses contain data only. ibc-go's 08-wasm keeper rejects attributes, events and
messages returned by a light client.

## Creation and profile

Creation uses the standard direct byte fields:

```text
InstantiateMsg { client_state, consensus_state, checksum }
```

The decoded client state contains `latest_height`, an optional `frozen_height`, one immutable
artifact profile and the active attestor configuration:

```json
{
  "latest_height": 1,
  "frozen_height": null,
  "profile": {
    "common": {
      "l2_chain_id": 11155420,
      "l2_router": "0x645280885749dc97ea461de280eb3273c91d36df",
      "commitment_slot": "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600",
      "profile_version": "op_attestor_v1",
      "l2_header_fork": "prague"
    }
  },
  "attestors": {
    "public_keys": ["<base64 Ed25519 public key>"],
    "threshold": 1
  }
}
```

The only accepted profile identifiers are `op_attestor_v1`, `base_attestor_v1` and
`arbitrum_attestor_v1`. Public keys must be unique, byte-sorted, exactly 32 bytes each and contain
at most 32 members. The threshold must be between one and the member count.

These development artifacts are fresh-authenticated-state-only. A client state missing `attestors`
is rejected with a typed error before any migration write; create a fresh authenticated client
instead. The strict signed-header and required-attestors schemas prevent unsigned state or messages
from crossing into this client; the `_v1` profile suffix is not treated as a security boundary.

Creation also validates revision-zero nonzero height, nonzero roots and block hash, and equality
between the client latest height and bootstrap consensus height. It performs no host or L1 query.

## Governed attestor rotation

Attestors can change only while governance migrates the client to a different checksum through
IBC-go's authority-only `ibc.lightclients.wasm.v1.MsgMigrateContract`. The `msg` bytes use one of
these strict JSON payloads:

```json
{"keep_attestors":{}}
```

```json
{
  "replace_attestors": {
    "public_keys": ["<base64 Ed25519 public key>"],
    "threshold": 1
  }
}
```

`keep_attestors` validates the loaded client and leaves storage byte-identical.
`replace_attestors` validates the replacement with the same rules as instantiation, then rewrites
only the client-state envelope. Height, frozen status, profile, checksum bytes visible to the
contract, and every consensus state remain unchanged. IBC-go commits the new checksum after the
entry point succeeds. Invalid replacements fail before any write, and the old set stops being
accepted immediately after a successful replacement.

Rotation applies only to an already-valid authenticated client. It cannot add attestors to legacy
unsigned state; both `keep_attestors` and `replace_attestors` fail without writes for that input.

This operation is a proactive rotation mechanism: use it before a retiring key can be abused, or
to restore liveness when an unavailable key has not already produced untrusted updates. It is not
post-compromise recovery. A frozen client stays frozen, and consensus states accepted under the old
key remain in the client store. If a key may already have been abused, the current contract must be
replaced by a newly created client, including the packet sequence and receipt coordination described
in the E2E runbook. A follow-up recovery PR will implement IBC-go's authority-only
`MsgRecoverClient`/`migrate_client_store` substitute-client flow; that lifecycle operation remains
unsupported in these artifacts.

## Signed update contract

Update and client-message verification consume this shape inside the tagged `ClientMessage`
envelope:

```text
ClientMessage::Header(SignedAttestedL2Header {
    l2_header: CanonicalEvmHeader,
    router_proof: EvmAccountProof,
    attestor_signature: Vec<IndexedAttestorSignature>
})
```

Every update carries exactly `threshold` signatures in strictly increasing attestor-index order.
The contract signs the fixed-width statement:

```text
SHA256("SPECTRE_L2_ATTESTATION_V1")
|| u64be(l2_chain_id)
|| l2_router[20]
|| SHA256(u16be(threshold) || u16be(member_count) || sorted_public_keys)
|| u64be(block_number)
|| block_hash[32]
|| state_root[32]
```

The statement remains exactly 164 bytes. Finality (`unsafe`, `safe`, or `finalized`) is an
off-chain selection policy and is intentionally not another signed field. Therefore every
finality tier must use a separate attestor/KMS key set. Config loading rejects the same attestor-set
hash when it appears under different `head_kind` values in one relayer config.

Signature structure and Ed25519 verification complete before router-proof traversal. The contract
then validates the configured header fork, derives the L2 block hash and verifies the router
account proof against the authenticated state root. Unsigned, malformed, duplicate-index,
under-threshold and wrong-context messages fail without writes.

Updates may advance the latest height or backfill an absent historical height. Re-submitting the
same block is idempotent; a different identity at an existing height is rejected. Misbehaviour
requires two independently authenticated, conflicting headers before the client can freeze.

## Packet proofs

Membership and non-membership load the exact stored L2 consensus height and verify a bounded
`EvmStorageProof` against its authenticated router storage root and configured commitment mapping
slot. Delay values other than zero fail closed.

## Fixtures and validation

The deterministic Go generator owns the cross-language Ed25519 vectors:

```bash
go run scripts/generate-l2-attestation-fixtures.go
```

The combined validator regenerates those vectors to a temporary file, checks the committed digest,
runs Rust and Go tests, builds all four Wasm artifacts and verifies their interfaces:

```bash
scripts/validate-l2-clients.sh
```

Set `VALIDATE_OPTIMIZED=1` to build each artifact twice with the pinned optimizer and run the pinned
wasmvm conformance and gas suites. Generated validation and gas reports are written under
`target/`; release manifests are not tracked while the application remains under development.

Run the complete release gate locally with:

```bash
just validate-wasm-release
```

The local gate requires a complete gas baseline and a fresh four-artifact conformance report. It
does not enable or depend on GitHub Actions while repository workflows remain disabled.
