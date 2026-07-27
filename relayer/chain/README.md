# `chain` — multi-chain adapter contract

This package is the **contract** that lets one relayer process serve many
source/destination chains (Cosmos today; EVM L2s Arbitrum / Optimism / Base
next) without the core knowing any chain's specifics. Implement the interfaces
here, register your adapter, and the generic relay loop (`relay.Module`) drives
it.

**Start here, then open `chain.go`** — it is the single source of truth for the
method signatures. This README is the map.

## The three interfaces

| Interface | Who/what it models | Methods |
|---|---|---|
| `Source` | a chain whose state we prove to a destination | `Chain`, `Subscribe`, `LatestHeight`, `RelayableHeight`, `QueryHeader`, `MembershipProof`, `NonMembershipProof` |
| `Destination` | a chain hosting the light-client-of-source, where we submit | `Chain`, `UpdateClient`, `RelayPackets`, `HasPacketReceipt`, `ClientExpiresAt` |
| `ClientUpdateBuilder` | the per-source strategy that builds a `ClientUpdate` | `Name`, `Build` |

Data crossing the boundary is deliberately opaque (`[]byte` / small structs) so
each adapter owns its own encoding. **Do not leak chain-specific types across
the interface** — that is the whole point (`Event.Raw` / `RelayPacket.Packet` are
proto-marshaled `channeltypesv2.Packet`; `ClientUpdate.Payload` is adapter-owned).

## The relay loop (`relay.Module`)

One `Module` drives a single source → destination path, depending only on the
interfaces above (no god-object, unit-testable with mock adapters). Per batch of
source events it:

```
1. track  every SendPacket for timeout recovery                (WithPacketTracker)
2. relayable := Source.RelayableHeight()   — cheap precondition
   • events above relayable  -> re-queue with WAITING backoff, no expensive work
   • none relayable          -> return (no client update, no proof)
3. Destination.UpdateClient(ClientUpdateBuilder.Build(...))   — advance once,
   skipped when the client already covers `relayable`
4. Source.Membership/NonMembershipProof(...) per provable packet
   • chain.IsPermanent(err)  -> DROP (timed out; the scanner refunds it)
   • else                    -> re-queue transient
5. Destination.RelayPackets(...) — one multicall for the batch
   • chain.IsPermanent(err)  -> DROP    • else -> re-queue transient
```

Alongside `Run`, the module runs three background loops: an **expiry-driven
refresh** (`Destination.ClientExpiresAt`), an optional **timeout scanner**
(`WithTimeoutScanner`), and an optional **periodic forced update**
(`WithPeriodicUpdate`, e.g. pinned-set rotation) — see "Reliability" below.

## Error classification (`errors.go`)

The module re-queues by default (never drop a valid packet) and only diverges on
explicitly-tagged errors:

| Tag | Meaning | Module action |
|---|---|---|
| `chain.Retryable` | transient infra (RPC blip) | re-queue immediately |
| *(waiting)* | not yet relayable (finality / AppHash H+2 lag) | re-queue with growing backoff |
| `chain.Permanent` | deterministic (timed-out packet, on-chain revert) | **drop** (funds still safe via the timeout scanner) |

Adapters tag their errors: e.g. a source's `MembershipProof` returns
`chain.Permanent` for a send whose timeout has passed; a destination's
`RelayPackets` wraps `services.ErrPermanentRelayFailure` (on-chain revert) as
`chain.Permanent`.

## Adding a chain = register + config, zero core edits

```go
func init() {
    chain.RegisterSource(chain.OPStack, newOptimismSource)
    chain.RegisterDestination(chain.OPStack, newOptimismDestination)
    chain.RegisterClientUpdateBuilder("l2-opstack", newOPStackBuilder)
}
```

Config selects adapters by `src_chain` / `dst_chain` (a `ChainType`) and
`builder` (a name). An unregistered name fails loudly at startup
(`UnknownAdapterError`). (The Cosmos↔ETH adapters are currently wired directly by
`cmd/run_adapters.go` because they need the shared `services` context/worker; the
cfg-only registry is for the pure-config L2 adapters.)

## Principle: wrap, don't rewrite

The existing Cosmos↔ETH code is battle-tested — **wrap it behind the
interfaces, don't reimplement it.** The `groth16` builder calls the current
signature-extract → pinned-set select → `GenerateProof` → assemble pipeline
**unchanged**; `Build` just returns the assembled message as
`ClientUpdate{Height, Payload}`. The generic reliability machinery (gap-recovery
subscriber, pending tracker, timeout scanner, nonce/sequence serialization, batch
builder) is reused, not rewritten — adapters only supply the per-chain parts
(RPC, event decode, tx build, finality notion). The `services` surface the
adapters reuse is exported in `services/adapter_api.go`.

If a builder needs to read destination trust state (e.g. the `groth16` builder
reads the pinned validator set + trusted consensus from the ETH SpectreClient),
keep that read inside the builder's own package via the shared context — do
**not** add Cosmos-specific reads to the generic `Destination` interface.

## Finality hook: `Source.LatestHeight` / `RelayableHeight`

Two related hooks carry the confirmation policy:

- **`LatestHeight`** — the highest height this adapter is willing to *advance the
  client to* (the finality/trust decision). It **must** be in the same number
  space as `RelayableHeight`, event heights, and `ClientUpdate.Height` (the module
  compares all four against `m.lastHeight`):
  - **Cosmos** (BFT instant finality): the latest committed height.
  - **ETH L1**: the finalized execution block number (the beacon client only
    advances to finalized headers, so latest == relayable here).
  - **EVM L2**: the height at the source's `HeadKind` confirmation policy —
    `Finalized` (L1-finalized L2 height) / `Safe` / `Unsafe` (soft head).
- **`RelayableHeight`** — the highest height whose packets can be *proven right
  now*, ≤ `LatestHeight`, encoding each chain's state-availability lag (Cosmos:
  `latest-2` for the AppHash H+2 lag; ETH: the finalized execution block). The
  module uses it as the cheap precondition so it never burns a proof for a packet
  that isn't yet provable (replaces the legacy blocking `waitCosmosAppHash` /
  `waitBeaconFinality`).

⚠️ **Confirmation policy trades safety for latency.** Proving at the unsafe head
exposes the bridge to L2 reorgs / sequencer equivocation (a released packet on
the destination against a source state that later reverts). The L2 wasm client
verifies the rollup **validity** proof (the state root is canonical per the L1
rollup contract) but that is **not** anti-reorg — anti-reorg comes only from the
`HeadKind` gate on `LatestHeight`. Changing the trust/latency tradeoff later =
changing only the source's `HeadKind`. See the return-path threat model in the
coordination issue.

## Reliability behaviors ported from the legacy `StartLoop`

The generic module reproduces every reliability property of the monolithic
`services.StartLoop`, so the cutover is behavior-preserving:

- **Requeue transient / drop permanent** — see the error table above. There is no
  dead-letter queue or retry budget (removed at cutover): a permanent packet is
  dropped and the timeout scanner refunds it. Because a multicall is atomic, a
  permanent BATCH failure retries each packet individually (`relayIsolated`), so a
  poison packet drops alone while its valid siblings still relay.
- **Timeout recovery** — `WithPacketTracker` records every send and untracks it once
  the relay succeeds (a delivered packet can no longer time out); `WithTimeoutScanner`
  runs the existing `ScanCosmosTimeouts` / `ScanEthTimeouts` on a ticker, refunding
  packets that expired undelivered.
- **Periodic forced update** (`WithPeriodicUpdate`, Cosmos→ETH only) — force-
  rotates the SpectreClient pinned validator set with exponential backoff on
  failure, so a quiet period cannot let the set decay below quorum. It is scheduled
  from the client's **on-chain freshness** (`PinnedSetRotationDueIn` → `initialDelay`),
  not process uptime, so a restart near the deadline rotates promptly; an
  underivable interval is fatal at startup. The beacon client has no pinned set, so
  ETH→Cosmos omits it.
- **Waiting backoff** — a not-yet-relayable packet re-checks with a growing delay
  (3→15s) instead of every batch period, so a long finality wait is quiet and
  does not churn the finality RPC (`services.BatchBuilder` `NotBefore`).
- **Multicall + shared serialization** — packets fold into one `RelayPackets`
  multicall; the shared `TransactionHandler` serializes ETH nonce / Cosmos
  sequence across all sources.
- **Graceful shutdown** — `handleBatch` takes the `Run` ctx and the background
  loops re-check `ctx.Err()` after each timer/tick, so cancellation stops new
  proof/tx work instead of continuing past teardown.

## L2 (rollup) adapters

The L2→Cosmos path is **trustless**: no relayer signature or re-execution. The
per-L2 `HeaderBuilder` (`l2rollup/header_op.go`, `l2rollup/header_arbitrum_bold.go`,
`l2rollup/header_arbitrum_legacy.go`) assembles a JSON
`ClientMessage` — the L1 rollup-contract witnesses (AnchorStateRegistry for
OP-Stack, RollupCore for Arbitrum) + the RLP L2 header + the L2 IBC-handler
account proof — and the L2 wasm client on Cosmos verifies those rollup proofs
against the shared cw-ics08-wasm-eth client (see `l2rollup/clientmessage.go`).
The generic `Builder` (`l2rollup/builder.go`) decodes the target height, delegates
to the `HeaderBuilder`, and packages the JSON as `ClientUpdate.Payload`.

Source/Destination/Builder are shared across all rollups; only the per-L2
`HeaderBuilder` differs. Optimism and Base share the OP-Stack `HeaderBuilder`
(`l2-opstack`), differing only by config; Arbitrum has its own (`l2-arbitrum`).
This is a skeleton — the eth_getProof / RLP / rollup-contract mechanics are TODO,
pending the L2 config and the shared-ETH-client handle.

## Status

- **Contract** (`chain.go`, `errors.go`) + **Cosmos↔ETH adapters**
  (`cosmos/`, `evm/`, `codec/`) + **generic `relay.Module`** + **`cmd` cutover**
  (`runAdapterEngine` replaces `svc.StartLoop` as the `start` engine): **done**.
  Both directions (Cosmos→ETH groth16, ETH→Cosmos beacon) relay via adapters,
  validated end-to-end on the local devnet.
- The legacy `services.StartLoop` + `handleCosmos`/`handleEth` have been removed;
  `runAdapterEngine` is the sole `start` engine.
- **L2 adapters** (`l2rollup/`, `l2-opstack` / `l2-arbitrum`): skeleton in place;
  proof assembly + event listener TODO (see the L2 section above).
- **Next**: wire the L2 proof assembly; dissolve the `services.Context`
  god-object into per-module context; fold the client update into the packet
  multicall (single tx).

Multi-source (many Cosmos chains, one ETH) works unchanged: `cmd` runs one
independent `runAdapterEngine` per `cosmos_to_eth` source, sharing the prover,
`TransactionHandler`, and ETH beacon endpoint; ETH events are partitioned by the
per-source router client id.
