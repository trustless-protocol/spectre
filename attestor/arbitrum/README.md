# Arbitrum Attestor

The attestor is a standalone gRPC backend for the relayer. Like the OP
attestor, it publishes a proposal-independent root from a configured Nitro
`unsafe`, `safe`, or `finalized` head. It also ingests BoLD v2 Arbitrum
assertions from finalized Ethereum L1 blocks and checks each assertion's L2
block against configured Nitro HTTP and WebSocket endpoints. It never starts,
supervises, or stops a Nitro process.

The configured Nitro endpoint performs execution and owns all chain state. The
attestor stores only its RollupCore log cursor, unresolved assertions, verified
frontier entries, and mismatch records; it does not maintain an execution
database.

## Trust boundary

For `source="assertion"` entries, the attested root comes from the assertion
read out of `RollupCore` on Ethereum L1. Nitro supplies the runtime head and
canonical-block view used to decide when that assertion can enter or advance in
the feed. The attestor re-reads assertion status at finalized L1 during
reconciliation and removes an entry that is no longer supported there.
Consequently, an untrusted Nitro endpoint cannot fabricate an assertion-backed
root; it can delay, suppress, or reject an otherwise valid assertion by
misreporting its head or canonical view, which is a liveness failure rather
than a new L1 commitment.

`source="derived"` is deliberately different and mirrors OP's configured-head
feed: its root is read directly from Nitro. `attestation_head="unsafe"` trusts
the configured endpoint immediately; `safe` and `finalized` select Nitro's
corresponding L1-derived heads. Set `disable_derived_roots=true` when the feed
must contain only RollupCore-backed assertions.

The L2 client does not see this distinction. It accepts one header shape — the
canonical L2 execution header plus the router account proof — and verifies no
settlement object, so `source` and `attestation_head` govern only what this
daemon puts in the feed, not what the client will admit. This daemon does not
authenticate the eventual `MsgUpdateClient` submitter either; see
`docs/L2_CLIENTS.md` for what that leaves open.

## Prerequisites

Provide an Ethereum L1 RPC and Nitro endpoints for the target Arbitrum chain.
The Nitro HTTP endpoint must support `eth_getBlockByNumber` with `latest`,
`safe`, and `finalized` tags. The WebSocket endpoint must support
`eth_subscribe("newHeads")`. Both endpoints are checked against `l2_chain_id`
at startup.

Copy `config.example.json` to `config.json` and configure:

- `grpc_listen_address`: attestor gRPC bind address. The example uses loopback
  because the server does not provide TLS or authentication.
- `runtime_poll_interval`: fallback reconciliation interval for Nitro's unsafe,
  safe, and finalized L2 commitments. New unsafe heads normally trigger
  immediate reconciliation through the WebSocket subscription.
- `attestation_head`: Nitro head used for proposal-independent attestations.
  `unsafe` trusts the configured Nitro node and does not wait for an L1
  assertion; `safe` uses Nitro's L1-posted view; `finalized` uses Nitro's
  L1-finalized view. The default is `finalized` when omitted.
- `disable_derived_roots`: disables proposal-independent roots while retaining
  BoLD assertion verification and finalized rechecks of existing derived
  entries.
- `derived_attestation_gap_blocks`: minimum L2 block gap between derived feed
  entries. Omitting it uses the conservative default `150`.
- `max_derived_roots`: maximum confirmed derived entries retained. Provisional
  derived entries and assertion-backed entries are never pruned.
- `src_chain`: source-chain key used by relayer requests.
- `l1_rpc_url`: Ethereum execution RPC used for finalized RollupCore logs and
  assertion status reads. The endpoint must support the `finalized` block tag.
- `l1_chain_id` and `l2_chain_id`: expected Ethereum and Arbitrum IDs. Startup
  fails if the connected L1 or RollupCore reports a different chain.
- `rollup_core_address`: deployed Arbitrum RollupCore proxy address.
- `assertions_mapping_slot` and `assertion_status_offset`: reviewed BoLD
  `_assertions` storage layout values. They must match the Cosmos Arbitrum
  verifier profile.
- `assertion_start_block`: first L1 block scanned on a new state database.
  Configure the RollupCore deployment block for full history, or the creation
  block of a known confirmed object when intentionally bootstrapping a frontier.
- `assertion_poll_interval` and `assertion_max_block_range`: finalized-L1 scan
  cadence and maximum inclusive block count per log query.
- `attestor_state_path`: durable JSON state for the assertion cursor and
  verified feed.
- `nitro_rpc_url`: Nitro HTTP/HTTPS endpoint used for chain state, block
  lookups, and unsafe/safe/finalized head reads.
- `nitro_ws_url`: Nitro WS/WSS endpoint used for `newHeads` notifications.

## Shared local L1 and Nitro devnet

For local bridge testing, Arbitrum and OP settle to the same Ethereum L1 in one
Kurtosis enclave. Start the OP stack first because it owns the shared
ethereum-package deployment, then layer Arbitrum onto that enclave:

```sh
# run_optimism_node.sh is optional: include it to share one L1 between OP and
# Arbitrum; omit it and run_arbitrum_node.sh brings the L1 up itself.
./scripts/local/run_optimism_node.sh
./scripts/local/run_arbitrum_node.sh
./scripts/local/run_arbitrum_attestor.sh
```

Run all commands from the repository root. `run_optimism_node.sh` starts
ethereum-package plus OP in the `op-devnet` enclave; when it is skipped,
`run_arbitrum_node.sh` brings the same L1 up in that enclave itself. The local package under
`scripts/local/kurtosis/arbitrum/` accepts that L1's internal execution and
beacon endpoints; it deploys RollupCore and starts a simple Nitro
sequencer/batch-poster/staker without creating another L1.

`run_arbitrum_node.sh` drives that package, waits for shared-L1 finality,
validates the configured BoLD storage layout against a finalized assertion
when one is available, and writes the published endpoints, chain IDs, and
RollupCore metadata to
`.arbitrum-devnet-run/attestor.env`.

`run_arbitrum_attestor.sh` automatically sources that handoff, builds the
attestor image, and connects it to the Nitro HTTP and WebSocket endpoints
published by Kurtosis. Its runtime files live in `.arbitrum-attestor-run/`;
assertion state lives in the persistent
`fast-ibc-arbitrum-attestor-data` Docker volume.

The resulting local flow matches OP:

```text
run_*_stack.sh    -> .*-devnet-run/attestor.env
run_*_attestor.sh -> external rollup-node endpoint + attestor gRPC on :3001
```

Stopping Arbitrum leaves Ethereum and OP running:

```sh
./scripts/local/run_arbitrum_node.sh --stop
./scripts/local/run_arbitrum_node.sh
```

Redeploying Arbitrum removes and recreates only the two `arb-*` Kurtosis
services. The shared L1 and OP services are preserved, although the previous
RollupCore contracts remain as unreachable history on the development L1:

```sh
./scripts/local/run_arbitrum_node.sh --reset
```

Removing the enclave removes Ethereum, OP, Arbitrum, and their local state:

```sh
kurtosis enclave rm -f op-devnet
```

The first Arbitrum run builds a RollupCreator image from the pinned
`NITRO_CONTRACTS_REF` and pulls the pinned `NITRO_IMAGE`. Override both pins
together only after checking their compatibility. The package uses public
development keys and funds them from the ethereum-package account already
used by the OP external-L1 configuration; none of those keys are suitable for
public networks.

To stop only the attestor, press Ctrl-C in its terminal. This removes its
container but preserves the assertion-state volume. To discard that state,
remove the volume explicitly after confirming it is no longer needed. The
attestor reaches Kurtosis' published L1 and Nitro endpoints
through `host.docker.internal`; it does not need to join Kurtosis' internal
Docker network.

## Run

From this directory:

```sh
cp config.example.json config.json
./scripts/start-attestor.sh ./config.json
```

For canonical Arbitrum Sepolia, start from the checked testnet profile:

```sh
./scripts/start-attestor-docker.sh ./config.arbitrum-sepolia.json
```

The tracked `config.arbitrum-sepolia.json` profile pins Sepolia L1
(`11155111`), Arbitrum Sepolia L2 (`421614`), and
the BoLD Rollup proxy (`0x042B...0Cf4`). It bootstraps at L1 block
`11379731`, which includes a complete confirmed assertion lifecycle. The
public RPC endpoints are convenient defaults but should be replaced by
operator-owned endpoints for sustained log scanning and Nitro synchronization.

The underlying Cobra command is:

```sh
./bin/attestor start --config ./config.json
```

The native launcher builds the attestor and starts it with the selected
configuration. It does not download, launch, or stop Nitro.

## Run in Docker

Keep `config.json` in a directory mounted into the container, then run:

```sh
./scripts/start-attestor-docker.sh ./config.json
```

The image contains only the attestor. `attested-roots.json` is persisted in the
`fast-ibc-attestor-data` Docker volume. The attestor gRPC service is published
at `127.0.0.1:50051` by default. The configured L1 and Nitro endpoints must be
reachable from inside the container.

The launcher can be customized with:

- `ATTESTOR_DOCKER_IMAGE`: image name.
- `ATTESTOR_CONTAINER_NAME`: container name.
- `ATTESTOR_GRPC_PUBLISH`: host gRPC publish address.
- `ATTESTOR_STATE_VOLUME`: persistent attestor-state volume name.

## gRPC API

The protobuf contract is in
`../../proto/attestor/attestor.proto`, with generated Go client code in
`../types/attestor`. Regenerate it from the repository root with
`buf generate --template buf.gen.attestor.yaml`. The reusable Nitro runtime state is in `runtime_state.go`;
the gRPC transport implementation is in `server`. The service exposes:

```text
attestor.AttestorService/AttestedUpTo
attestor.AttestorService/AttestedRootAtOrBelow
attestor.AttestorService/VerifyStateRoot
```

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

`VerifyStateRoot` remains available as a lower-level compatibility and
diagnostic RPC. The caller supplies:

- `block_number`
- `expected_state_root` as exactly 32 bytes
- optional `expected_block_hash` as exactly 32 bytes
- `run_mode` as `RUN_MODE_UNSAFE`, `RUN_MODE_SAFE`, or `RUN_MODE_FINALIZED`

The response returns `valid` plus Nitro's canonical block hash and state root.
A mismatch is a successful gRPC response with `valid=false`; malformed requests,
missing blocks, and Nitro availability failures use gRPC status errors.

Every verification request is bounded by its requested run-mode head: unsafe uses Nitro's
latest head, safe uses its safe head, and finalized uses its finalized head.
The service supports all three modes concurrently and rejects an unspecified
mode or a requested block above the selected head.

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

The standard gRPC health service is also registered for readiness checks.

The relayer builds the same header at every head kind: the canonical Nitro L2
header at the target height plus the router account proof. It never proves a
RollupCore assertion — the client verifies none — so this feed bounds the relay
on two axes instead. `AttestedUpTo` bounds how far it may advance, and
`VerifyStateRoot` answers whether the block it is about to package is the one
this daemon's Nitro connection holds at that height, at the run mode matching
the configured head kind.

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
