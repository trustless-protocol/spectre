# Optimistic L2 ICS-08 clients

Fast-IBC builds three checksum-distinct 08-wasm artifacts with unchanged filenames:

- `cw-ics08-wasm-arbitrum` for Arbitrum BoLD v2 assertions;
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
pending assertion. Base and OP accept any caller-selected factory game-list entry, including an
unresolved game. They authenticate the game runtime, extract its root claim at the runtime-profile
offset, bind the output-root preimage to a canonical L2 block header, and authenticate the L2
`ICS26Router` account. Game type, resolution, winner, retirement, blacklist, pause state,
implementation identity, and settlement delays are deliberately not checked.

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
state is idempotent. A conflicting update at an existing height is handled by the stored
`FinalityPolicy`, not solely by `UpdateStateOnMisbehaviour`.

> **Unsafe-head reorg warning:** the default `FinalityPolicy` has
> `minimum_membership_level: unsafe` and `freeze_on_trusted_conflict: true`. Therefore every
> non-`resolved_invalid` Unsafe consensus state is trusted for conflict handling. With
> `head_kind: unsafe`, a same-height sequencer reorg after the first block was relayed is a
> trusted conflict and freezes the client. This is intentional current behavior, not an automatic
> reorg correction; frozen clients reject updates and packet proofs until governance recovery.
>
> The replacement path is available only when the existing state is below the configured
> membership threshold and the incoming state reaches it. In the current implementation Safe
> evidence fails closed (`SafeVerificationUnavailable`), so a practical non-freezing correction
> policy requires `minimum_membership_level: finalized` (and suitable finalized evidence). Do not
> use `head_kind: unsafe` with the default policy if normal sequencer reorgs must not freeze the
> client.

To deliberately replace Unsafe conflicts, set
`freeze_on_trusted_conflict: false`; a Finalized conflict still freezes unconditionally. This
trades the freeze/liveness failure for the normal safety risk of accepting and using an Unsafe
state. Alternatively, setting `minimum_membership_level: finalized` keeps the freeze-on-trusted-
conflict safeguard while making Unsafe updates ineligible for membership proofs and conflict
freezing.

`UpdateStateOnMisbehaviour` remains the watchdog mechanism for freezing on two independently
valid conflicting headers. Frozen clients reject updates and packet proofs.

OP and Base profiles also pin `l2_header_fork`, which selects the canonical execution-header field
set. Operators must migrate to a reviewed profile at an L2 fork boundary instead of relying on a
fork hardcoded into the verifier.

## Fixture requirements

Reproducible fixtures use fixed numeric L1 and L2 block tags and include full block responses,
`eth_getProof` account/storage responses, the L2 router proof, and `eth_getCode` for OP Stack games.
BoLD Arbitrum fixtures additionally include the complete `AssertionCreated` data required to
recompute the assertion hash. The adjacent schema-v2 provenance manifest records both block hashes,
source revisions, the profile SHA-256, and the fixture SHA-256.

## Host requirements and verification

The Cosmos application must retain its BLS custom querier and permit the client-state, status, and
consensus-state IBC query paths. The current adapter uses CosmWasm's deprecated
`QueryRequest::Stargate`, so a stock host that does not whitelist those paths is incompatible; the
deployment must provide the allowlist or a custom Wasm host image. The L2 clients pin the Ethereum
Wasm checksum and inherit its active/inactive status.

The allowlist that matters is the **08-wasm** one (light clients run there), not the `x/wasm` one
used by ordinary contracts — a host can have the paths in the latter and still reject the former.
All three paths are required:

```go
// gaia app/keepers/keepers.go — ibcwasmkeeper.QueryPlugins
Stargate: ibcwasmkeeper.AcceptListStargateQuerier([]string{
    "/ibc.core.client.v1.Query/ClientState",
    "/ibc.core.client.v1.Query/ClientStatus",     // required by L2 clients
    "/ibc.core.client.v1.Query/ConsensusState",
}, bApp.GRPCQueryRouter()),
```

`ClientStatus` is easy to miss because the Ethereum client never queries another client — the L2
clients are the first to do so, in `validate_l1_client` (they must confirm the pinned Ethereum
client is still `Active`). Its absence does **not** surface as a permission error:

| Symptom | Cause |
|---|---|
| `MsgCreateClient` fails: `cannot create client (08-wasm-N) with status Unknown: client state is not active` | The L2 client's `Status{}` query errored inside the contract; ibc-go maps the error to `Unknown`. Check the 08-wasm Stargate allowlist before suspecting the client state. |

If the pinned Ethereum client is genuinely `Active` and its checksum matches
`profile.common.ethereum_client.wasm_checksum` (verify both with `gaiad q ibc client status` /
`state`), the allowlist is the remaining suspect.

A second, more confusing symptom comes from the same query path:

| Symptom | Cause |
|---|---|
| Every update fails: `IBC host query failed: codespace: undefined, code: 1: wasm contract call failed` | The contract asked the host for the Ethereum consensus state at the header's `beacon_slot` and got `NotFound`. ibc-go answers with a gRPC status error, which wasmd redacts to a bare codespace/code — so neither the client nor the slot appears anywhere. |

The usual cause is a client-id mismatch: the relayer builds headers pinned to the Ethereum client
its config names, while the contract queries the one baked into its own client state at creation.
The client state is authoritative. `create-clients-cosmos` writes the created id into both the
client state and the `l2_to_cosmos` module's `rollup_profile.common.ethereum_client.client_id`, and
`start` refuses to boot when they disagree, naming both ids. A config assembled by hand — or carried
over from an earlier devnet run — is what re-opens this.

The same mismatch has a slower failure mode with no error at all: the client the config names keeps
being advanced while the client the contract actually reads goes stale, and eventually expires.

### Packet proofs are a bare storage proof, not the L1 shape

The Ethereum L1 membership proof (`client.GetEthMembershipProof`) bundles an account proof
with the storage proof and hex-encodes both, because the ETH light client re-derives the
account from the L1 state root. An L2 client must **not** be given that: it already
authenticated the router's storage root through the header's `router_proof`, so it expects a
bare `EvmStorageProof` and rejects anything else outright (`deny_unknown_fields`):

```
unknown field `account_proof`, expected one of `key`, `value`, `proof`
```

Two encoding details that are easy to get wrong, and only fail on-chain:

| Field | Wire form | Why |
|---|---|---|
| `key` | `0x`-hex string | serde `B256` |
| `value` | JSON **number array**, full **32 bytes** | Go `[]byte` marshals to base64 by default — use the package's `byteList`. The client compares the value byte-for-byte against the commitment the IBC host expects, which is a full 32-byte word, not `eth_getProof`'s minimal big-endian form. The trie check is unaffected: `encode_storage_value` strips leading zeros itself |
| `proof` | array of number arrays | same `[]uint8` trap — use `byteMatrix` |

### Game selection must follow finality, not the frontier

Dispute games are posted at the L1 **head**, while the pinned Ethereum client only advances to
**finalized** L1. On a rollup that posts games about as fast as finality advances, the
attestor's frontier stays permanently ahead of what is provable, so demanding the frontier
game deadlocks — the target rises exactly as fast as the pinned block does and the client never
moves. Demanding the packet's own height instead selects a game at or *below* it, which by
construction cannot cover that packet.

The builder therefore walks **down** from the requested height to the highest attested game
that is actually present in the factory's list at the pinned block (`gameCount()` read there).
Every candidate is still a root the attestor approved, so this only ever picks an older
verdict — never an unverified one — and the client advances monotonically until it crosses any
given packet height.

`attested game N is not visible at the pinned L1 block M yet` is the normal wait for the newest
game, not an error; it clears once that game's creation block finalizes.

### Light clients must return nothing but data

ibc-go's 08-wasm keeper rejects a light client whose response carries attributes, events, or
messages, and it does so by panicking inside the VM call — so the entire transaction fails, not just
the update:

```
recovered: checksum (...): returning attributes from a contract is not allowed
  [08-wasm/keeper/contract_keeper.go:155]
```

The shared `l2_client_entrypoints!` macro therefore returns `Response::default().set_data(data)` and
nothing else. Observability that would naturally be an event (state accepted vs promoted, freeze on
a trusted conflict, finality level) has to come from the caller or from querying the client state
after the update. Do not re-add `add_attributes`: it compiles, it passes every unit test, and it
fails only on-chain.

```bash
cargo test --locked \
  -p l2-client -p op-stack-verifier -p op-verifier -p base-verifier \
  -p arbitrum-verifier -p cw-ics08-wasm-op -p cw-ics08-wasm-base \
  -p cw-ics08-wasm-arbitrum
cargo build --target wasm32-unknown-unknown --release --locked \
  -p cw-ics08-wasm-op -p cw-ics08-wasm-base -p cw-ics08-wasm-arbitrum
```
