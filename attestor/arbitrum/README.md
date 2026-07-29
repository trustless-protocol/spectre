# Arbitrum Attestor

The attestor is a standalone gRPC backend for the relayer. It launches a pinned
Nitro binary, keeps Nitro's chain database on persistent storage, ingests
Arbitrum assertions from finalized Ethereum L1 blocks, and independently checks
each assertion's L2 block against Nitro over private IPC. It supports both BoLD
v2 `AssertionCreated`/`AssertionConfirmed` deployments and the legacy Nitro
`NodeCreated`/`NodeConfirmed` lifecycle still used by canonical Arbitrum
Sepolia.

Nitro performs execution and owns all chain state. The attestor stores only its
RollupCore log cursor, unresolved assertions, verified frontier entries, and
mismatch records; it does not maintain a second execution database.

## Prerequisites

For a native run, download a trusted Nitro release or build Nitro from source
before starting the attestor and set `nitro_binary_path` to that executable.
The Docker launcher instead builds on the pinned official
`offchainlabs/nitro-node:v3.11.2-3599aca` image, including its Nitro binary and
validation-machine files.

Nitro also needs its normal chain configuration, including the chain ID, parent
chain RPC and beacon endpoints, sequencer feed, and first-start snapshot or
genesis settings. The example passes `--conf.file=./nitro.json`; provide that
file or replace `nitro_arguments` with the required Nitro flags.

Copy `config.example.json` to `config.json` and configure:

- `grpc_listen_address`: attestor gRPC bind address. The example uses loopback
  because the server does not provide TLS or authentication.
- `runtime_poll_interval`: fallback reconciliation interval for Nitro's unsafe,
  safe, and finalized L2 commitments. New unsafe heads normally trigger
  immediate reconciliation through the private IPC subscription.
- `src_chain`: source-chain key used by relayer requests.
- `l1_rpc_url`: Ethereum execution RPC used for finalized RollupCore logs and
  assertion status reads. The endpoint must support the `finalized` block tag.
- `l1_chain_id` and `l2_chain_id`: expected Ethereum and Arbitrum IDs. Startup
  fails if the connected L1 or RollupCore reports a different chain.
- `rollup_core_address`: deployed Arbitrum RollupCore proxy address.
- `rollup_protocol`: `bold-v2` for current BoLD deployments or `legacy-nitro`
  for canonical Arbitrum Sepolia. An omitted value retains the historical
  `bold-v2` default.
- `assertions_mapping_slot` and `assertion_status_offset`: reviewed BoLD
  `_assertions` storage layout values. They are required only for
  `rollup_protocol="bold-v2"` and must match the Cosmos Arbitrum verifier
  profile.
- `assertion_start_block`: first L1 block scanned on a new state database.
  Configure the RollupCore deployment block for full history, or the creation
  block of a known confirmed object when intentionally bootstrapping a frontier.
- `assertion_poll_interval` and `assertion_max_block_range`: finalized-L1 scan
  cadence and maximum inclusive block count per log query.
- `attestor_state_path`: durable JSON state for the assertion cursor and
  verified feed.
- `nitro_binary_path`: previously downloaded or built Nitro executable.
- `nitro_data_dir`: durable Nitro chain database directory.
- `nitro_ipc_path`: temporary private IPC socket used by the attestor.
- `nitro_arguments`: Nitro network, feed, snapshot, and RPC configuration.

The attestor creates `nitro_data_dir` when needed and always appends
`--persistent.chain=<nitro_data_dir>` when launching Nitro. It never removes
that directory during shutdown. Do not add `--persistent.chain` or `--ipc.path`
to `nitro_arguments`; those paths are owned by the attestor.

For production sizing, snapshot initialization, pruning, and archive retention,
follow the [official Nitro node documentation](https://docs.arbitrum.io/run-arbitrum-node/run-full-node).

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
when one is available, and writes the published endpoints, chain IDs,
RollupCore metadata, Nitro image pin, and self-contained sequencer config to
`.arbitrum-devnet-run/attestor.env`.

`run_arbitrum_attestor.sh` automatically sources that handoff, builds the
attestor image from the same official Nitro release, and starts an independent
non-sequencing Nitro replica plus the gRPC attestor. Its runtime files live in
`.arbitrum-attestor-run/`; its chain database lives in the persistent
`fast-ibc-arbitrum-attestor-nitro` Docker volume. The devnet sequencer and
verifier replica therefore never share an L2 database.

The resulting local flow matches OP:

```text
run_*_stack.sh    -> .*-devnet-run/attestor.env
run_*_attestor.sh -> independent verifier replica + attestor gRPC on :3001
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
container but preserves the verifier database volume. To discard that
database, remove the volume explicitly after confirming it is no longer
needed. The attestor reaches Kurtosis' published L1 and sequencer endpoints
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
the canonical Rollup proxy (`0xd808...81c8`). It bootstraps at L1 block
`7258441`, where currently confirmed legacy node `10764` was created; the node
was confirmed at block `7258462`. Use deployment block `4139226` instead when
full legacy history is required. The public RPC endpoints are convenient
defaults but should be replaced by operator-owned endpoints for sustained log
scanning and Nitro synchronization.

The underlying Cobra command is:

```sh
./bin/attestor start --config ./config.json
```

The launcher calculates the configured Nitro binary's SHA-256, atomically
writes it to `nitro_binary_sha256` for native runs, builds the attestor, and
starts it. The Docker entrypoint injects the hash recorded in its image without
modifying the host configuration. The
attestor verifies the binary again before launching Nitro. `jq` and either
`sha256sum` or `shasum` must be installed.

Nitro is stopped gracefully when the attestor exits. Only the IPC socket is
removed; the chain database remains in `nitro_data_dir` for the next startup.

## Run in Docker

Keep `config.json`, `nitro.json`, and any other files referenced by relative
Nitro arguments in the same configuration directory, then run:

```sh
./scripts/start-attestor-docker.sh ./config.json
```

The script pulls the official multi-platform Nitro base image while building
the attestor image. The image retains Nitro's validation-machine directories
and injects the official `--validation.wasm.allowed-wasm-module-roots` setting
when it is absent from `nitro_arguments`. The build records the bundled Nitro
binary's SHA-256; the container entrypoint verifies it again and renders an
internal configuration without modifying the host file.

Nitro data and `attested-roots.json` are persisted in the
`fast-ibc-nitro-data` Docker volume. The attestor gRPC service is published at
`127.0.0.1:50051` by default; Nitro RPC remains private inside the container.
The L1 endpoint and all endpoints in `nitro_arguments`, including parent-chain
and beacon URLs, must be reachable from inside the container.

The launcher can be customized with:

- `ATTESTOR_DOCKER_IMAGE`: image name.
- `ATTESTOR_NITRO_IMAGE`: official Nitro base image; defaults to the pinned
  `offchainlabs/nitro-node:v3.11.2-3599aca` release.
- `ATTESTOR_CONTAINER_NAME`: container name.
- `ATTESTOR_GRPC_PUBLISH`: host gRPC publish address.
- `ATTESTOR_NITRO_VOLUME`: persistent Nitro volume name.

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

For Arbitrum, an `AttestedRoot` has `source="assertion"` and `root` equal to the
Nitro L2 state root. BoLD provenance uses `assertion_hash`; legacy Nitro
provenance uses `legacy_node`, containing both the numeric `_nodes` mapping key
and the actual `nodeHash` indexed by `NodeCreated`. A pending assertion enters
the feed as `provisional=true` only after its block is canonical under Nitro's
safe head. It becomes non-provisional only after RollupCore reports it confirmed
in finalized L1 state and the same block is canonical under Nitro's finalized
head. Rejected, challenged, or locally mismatched assertions are never returned.

The chain-specific provenance is a protobuf `oneof`: OP entries carry
`game_index`, Arbitrum BoLD entries carry `assertion_hash`, legacy Arbitrum
entries carry `legacy_node`, and derived entries may carry neither.

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
the existing private IPC connection. Each unsafe-head event triggers a refresh
of the unsafe, safe, and finalized tags. The attestor reconnects and
reconciles after a subscription failure, and retains `runtime_poll_interval` as
a fallback because safe or finalized can advance without a new unsafe block.

The attestor records the first unsafe and safe root observed at each height
until it finalizes. When Nitro advances the finalized head, the runtime compares
the finalized root with both earlier observations and logs a consistency
result. This in-memory result is intended as the source for future alerting;
Nitro remains the persistent owner of chain state.

The standard gRPC health service is also registered for readiness checks.

## Assertion feed flow

On every assertion poll, the attestor:

1. Reads either the BoLD assertion lifecycle or the legacy
   `NodeCreated`/`NodeConfirmed`/`NodeRejected` lifecycle only through
   Ethereum's finalized L1 head and advances a persisted cursor in bounded
   ranges.
2. Recomputes every BoLD assertion hash. In legacy mode it retains both the
   indexed node hash and numeric node ID, allowing status checks to survive
   daemon restarts and rejected nodes to be removed exactly.
3. Resolves the asserted L2 block by hash in Nitro, then checks the canonical
   block at the recovered height. Pending assertions are checked against the
   safe head.
4. Re-reads either `AssertionNode.status` or the legacy `getNode` and
   `latestConfirmed` state at finalized L1. Confirmed assertions are promoted
   only after the Nitro finalized head covers and matches the block.
5. Serves the resulting frontier through the same oracle-shaped RPCs used by
   the OP attestor.
