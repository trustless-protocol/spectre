# Optimistic L2 ICS-08 clients

Fast-IBC builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum` for Arbitrum BoLD v2 assertions and legacy Nitro
  numeric nodes;
- `cw-ics08-wasm-base` for Base dispute games; and
- `cw-ics08-wasm-op` for OP dispute games.

> **Security assumption:** these clients are not safe against a dishonest rollup proposer on their
> own. Operators must run an independent watchdog that detects a challenged or invalid accepted
> proposal and submits two conflicting headers through `UpdateStateOnMisbehaviour` to freeze the
> client before affected packet proofs are used.

These are optimistic clients. Every update still waits until the proposal and all account/storage
proofs are contained in an Ethereum state finalized by the pinned Ethereum light client. It does
not wait for rollup challenge settlement. Consequently, a consensus state accepted by one of
these clients can later be rejected by the rollup challenge system.

Arbitrum accepts any BoLD assertion whose authenticated packed status is nonzero, including a
pending assertion. For legacy Nitro it accepts the latest confirmed node or a pending node between
the authenticated `firstUnresolvedNode` and `latestNodeCreated`, then binds the L2 block hash and
send root to the proven `Node.confirmData`. Resolved losing nodes are excluded. Base and OP accept
any caller-selected factory game-list entry, including an unresolved game. They authenticate the
game runtime, extract its root claim at the runtime-profile offset, bind the output-root preimage
to a canonical L2 block header, and authenticate the L2 `ICS26Router` account. Game type,
resolution, winner, retirement, blacklist, pause state, implementation identity, and settlement
delays are deliberately not checked.

## Creation and runtime profiles

Creation uses three JSON byte fields:

```text
InstantiateMsg { client_state, consensus_state, checksum }
```

The decoded states are stored directly at `client_state.latest_height`. Creation performs no L1
query or client-defined semantic validation. Deployment identities are immutable data inside the
client state's runtime profile; the example Sepolia JSON files under each verifier's `config/`
directory are tooling/test inputs and are not compiled into Wasm. Operators must set the actual
Ethereum client ID and checksum before client creation.

### L2 client creation config

`create-clients-cosmos --l2-config <path>` (repeatable, one per rollup source) creates each L2 wasm
client on Cosmos, anchored to the L1 (08-wasm ETH) client created in the same run. Each
`--l2-config` JSON is:

```json
{
  "wasm_checksum": "<hex checksum of the governance-stored L2 client wasm>",
  "l2_rpc_url": "https://<l2-execution-rpc>",
  "rollup_profile": {
    "common": {
      "l2_router": "0x…",
      "ethereum_client": { "client_id": "", "wasm_checksum": [] }
    }
  },
  "bootstrap_block": 0,
  "counterparty_client_id": "<L2-side client id on the rollup's ICS26Router>"
}
```

`rollup_profile` is the full ICS-08 verifier Profile embedded verbatim as the client state (the L2
router, commitment slot, and chain ids all live in `common`, so they are not repeated elsewhere).
`rollup_profile.common.ethereum_client.client_id` may be left empty — `create-clients-cosmos`
injects the freshly-created L1 client id. `bootstrap_block` `0` (or absent) bootstraps from the L2
latest. **`counterparty_client_id`** registers the L2-side client (the one tracking Cosmos) as this
client's counterparty inline; leave it empty to defer registration until that id is known, then
register it manually. Omitting it silently defers registration, so set it once the L2-side client
id is available.

Updates may advance the latest height or backfill an absent historical height. Replaying identical
state is idempotent; conflicting state at an existing height is rejected but does not freeze the
client. A lone conflict may be a relayer error. Freezing requires the watchdog to submit two
independently valid optimistic headers for the same L2 height through
`UpdateStateOnMisbehaviour`; when they authenticate different consensus states, the client freezes.
Frozen clients reject updates and packet proofs.

OP and Base profiles also pin `l2_header_fork`, which selects the canonical execution-header field
set. Operators must migrate to a reviewed profile at an L2 fork boundary instead of relying on a
fork hardcoded into the verifier.

## Fixture requirements

Reproducible fixtures use fixed numeric L1 and L2 block tags and include full block responses,
`eth_getProof` account/storage responses, the L2 router proof, and `eth_getCode` for OP Stack games.
BoLD Arbitrum fixtures additionally include the complete `AssertionCreated` data required to
recompute the assertion hash. Legacy Nitro fixtures include the numeric node, send root, lifecycle
proof, and `confirmData` proof. The adjacent schema-v2 provenance manifest records both block
hashes, source revisions, the profile SHA-256, and the fixture SHA-256.

## Host requirements and verification

The Cosmos application must retain its BLS custom querier and permit the client-state, status, and
consensus-state IBC query paths. The current adapter uses CosmWasm's deprecated
`QueryRequest::Stargate`, so a stock host that does not whitelist those paths is incompatible; the
deployment must provide the allowlist or a custom Wasm host image. The L2 clients pin the Ethereum
Wasm checksum and inherit its active/inactive status.

```bash
cargo test --locked \
  -p l2-client -p op-stack-verifier -p op-verifier -p base-verifier \
  -p arbitrum-verifier -p cw-ics08-wasm-op -p cw-ics08-wasm-base \
  -p cw-ics08-wasm-arbitrum
cargo build --target wasm32-unknown-unknown --release --locked \
  -p cw-ics08-wasm-op -p cw-ics08-wasm-base -p cw-ics08-wasm-arbitrum
```
