# Reliability

## Error Handling

### Solidity

- Custom errors defined in `contracts/errors/` (gas-efficient vs string reverts)
- `require()` with descriptive messages for input validation
- Proof verification failures revert the entire transaction (atomic)

### Go Relayer

- `log.Fatal` for unrecoverable startup errors
- `start` command: drives bi-directional relay (Cosmos ↔ ETH) through the generic chain-adapter engine `runAdapterEngine` (`relayer/cmd/run_adapters.go`), which runs one `relay.Module` (`relayer/relay/module.go`) per direction over the `chain.Source`/`chain.Destination` adapters. One engine runs per `cosmos_to_eth` source, so a stall on one source does not block the others. (This replaced the monolithic `services.StartLoop`, removed after the cutover; the reliability properties below are preserved.)
- `create-clients-cosmos` / `create-clients-eth` commands: one-time setup, fail fast on deployment errors; run `create-clients-cosmos` before `create-clients-eth` (each writes its client id back to config), so the flow supports staged retry; `create-clients-eth` skips the deploy if `spectre_client` already has code (idempotent). With multiple `cosmos_to_eth` sources, `--source <ics26_client_id>` selects the module and the write-back targets it; an unknown or ambiguous source fails fast rather than writing the wrong module
- Transaction handler retries with re-queried account sequence on nonce conflicts
- `BatchBuilder` manages packet batching with separate mutexes for Cosmos / ETH queues; it owns separate pending trackers for Cosmos-originated packets awaiting ETH delivery and ETH-originated packets awaiting Cosmos delivery/ack/timeout
- `SubscribeCosmos` performs startup and periodic gap recovery with CometBFT `TxSearch` for `send_packet`, `write_acknowledgement`, and `timeout_packet` event payloads. `COSMOS_STARTUP_LOOKBACK_BLOCKS` defaults to 256; set it to `0` to scan the full indexed history. The adapters wrap this subscriber unchanged (via `Source.Subscribe`), so gap recovery is retained.
- `SubscribeEth` performs startup/reconnect log recovery for `SendPacket` and `WriteAcknowledgement`. `ETH_STARTUP_LOOKBACK_BLOCKS` defaults to 256; set it at least as large as the maximum expected relayer downtime if older ETH logs must be recovered.
- Gap-recovery cursors survive restarts in `.relayer-state/recovery-cursors.json` (override with `RELAYER_RECOVERY_STATE_FILE`). The versioned JSON file stores Cosmos, ETH-send, and ETH-write-ack cursors per `ics26_client_id`; corrupt or unreadable state is fatal at startup because silently falling back to the lookback would recreate the outage gap. Recovery advances and atomically persists after each successful provider-sized chunk (`COSMOS_RECOVERY_CHUNK_HEIGHTS` / `ETH_RECOVERY_CHUNK_BLOCKS`, both default 256), never past the lowest packet still queued or in flight. Deployments running independent relayer processes from one working directory must give each process a distinct state-file path.
- Timeout checking: packets past their timeout are skipped before submission; subscriber normalises IBC v2 timeouts from ns → s at ingest (`normalizeTimeoutSeconds`)
- Provability precondition (non-blocking): before proving, the module partitions a batch by `Source.RelayableHeight` — the highest source height provable right now (Cosmos: `latest-2` for the AppHash H+2 lag, `relayer/chain/cosmos/source.go`; ETH: the finalized execution block, `relayer/chain/evm/source.go`). Packets above it are re-queued with a growing waiting backoff (3→15 s) instead of burning a proof/RPC every flush. This replaced the blocking `waitCosmosAppHash` / `waitBeaconFinality` waits with a value check.
- Permanent-vs-transient classification (`relayer/chain/errors.go`): a transient failure (RPC blip, not-yet-provable) is re-queued; a permanent failure (deterministic on-chain revert, or a send already past its timeout) is DROPPED — funds stay safe because the timeout scanner refunds a dropped packet once it expires. Timeout submissions themselves use a bounded retry budget: a deterministic failure becomes a visible dead letter for 30 days, then returns to the timeout scanner with a fresh budget. The retention period prevents event replays from bypassing the budget while ensuring stale client/sequence-number reuse cannot be hidden permanently.
- Poison-packet isolation: because a multicall is atomic, one bad packet (timed-out / duplicate / already-received) reverts the whole batch. On a permanent batch failure the module retries each packet individually (`relayer/relay/module.go` `relayIsolated`), so the poison packet drops alone while its valid siblings still relay.
- Cosmos→ETH timeout fallback: `scanForCosmosTimeouts` (every 30 s, wired via the module's `WithTimeoutScanner` hook → `ScanCosmosTimeouts`) scans `PendingPacketTracker`, builds `MsgTimeout` for each expired entry, batches all timeouts plus a `MsgUpdateClient` into one Cosmos tx, and `Recover()`s from panic so a single bad packet can't kill the goroutine
- ETH→Cosmos timeout fallback: `scanForEthTimeouts` (every 30 s, via `WithTimeoutScanner` → `ScanEthTimeouts`) scans `EthPendingTracker`, re-checks the ETH source commitment, then reuses the `timeoutEthSend` proof flow for expired ETH-originated packets
- The Cosmos source skips relaying a Cosmos-emitted `timeout_packet` back to ETH for Cosmos-originated packets (`ShouldRelayCosmosTimeoutToEth`, applied in `relayer/chain/cosmos/source.go`): the source-side MsgTimeout on Cosmos already refunded the sender and ETH never held a commitment, so calling `ICS26Router.timeoutPacket` would be a NoOp at best
- Pending-tracker hygiene: the module tracks every Cosmos-origin send for timeout recovery and removes it once the relay succeeds (`untrack`), so a delivered packet does not linger in the tracker; the ETH direction removes on the terminal `EthAck`/`EthTimeout` its subscriber observes
- Graceful shutdown: the signal-derived `Run` context reaches subscribers, reconnect waits, WebSocket watches, gap-recovery chunks, client updates, proofs, storage reads and timeout scanners. Each `relay.Module` registers named `subscriber`, `refresh`, `timeout-scan`, and `periodic-update` workers and drains them for up to 30 seconds; a timeout is an error and names every worker still running. The command closes Cosmos/EVM clients only after module drain returns. CometBFT unsubscribe cleanup deliberately uses `context.WithoutCancel(parent)` with its own 5-second timeout, so cancellation does not skip cleanup and cleanup cannot hang shutdown.
- Client-update idempotence: `BuildCosmosClientUpdateMsg` always re-queries the on-chain SpectreClient state and short-circuits to `HasMsg=false` when on-chain trusted height already covers the chunk's source block — avoids regenerating a Groth16 proof that another chunk (or a frontrunner relayer) already submitted
- Validator-set rotation is governed by two `cosmos_to_eth` config values; `updateConsensusState` (the heavy path: full valset calldata + SSTORE2 write + snapshot push) is submitted when either fires, otherwise updates stay on the cheaper `updateApplicationState`:
  - `rotation_threshold` (default `5/6`, must exceed the `2/3` quorum floor): on any update whose target `nextValidatorsHash` differs from the on-chain pinned set, rotate once the pinned set's voting-power overlap with the target block's signers has decayed to ≤ this fraction. `1/1` rotates on any pinned-set change; smaller values defer rotation so routine voting-power churn doesn't push per-packet flushes onto the heavy path
  - `refresh_interval_seconds` (default 86400): the guaranteed-rotation cadence for quiet chains (on busy chains the threshold is the primary trigger). It runs as the module's periodic-update hook (`relayer/cmd/run_adapters.go` → `svc.RotatePinnedSet`), scheduled from the client's **on-chain freshness** (`PinnedSetRotationDueIn`), not process uptime — so a restart near the rotation deadline rotates promptly instead of waiting a whole fresh interval. An interval that cannot be safely derived from the on-chain trusting period is FATAL at startup rather than silently defaulting (which could exceed the trusting period and let the client expire)
- Pinned Ethereum client freshness (L2 paths): an L2 header builder can only prove the rollup contracts at the L1 block the Cosmos-side Ethereum client trusts, so that client's freshness is a precondition of building, not a background nicety. Two mechanisms, and both are needed — neither substitutes for the other:
  - **On demand** (`chain/l2rollup/header_*.go` → `UpdateEthClientIfStale`): the builder asks for freshness at the moment it needs it. A timer cannot guarantee the client is current *at the point of use*.
  - **Periodic** (`cmd/build_l2_source.go` → `runEthClientFreshnessLoop`, every `ethClientPollInterval` = 8 min): the on-demand path only runs while packets are flowing, so an idle relayer would never advance the client at all. Observed idle at 260 blocks behind. This is the mirror of the Cosmos→ETH `refresh_interval_seconds` routine above.
  - The tick only *checks* (one `eth_blockNumber`, one client-state read) and no-ops unless the client is more than `ethClientUpdateLag` (64 blocks) behind, sharing that threshold, the cooldown and the mutex with the on-demand path. Sizing is bounded by the execution node's **state window** (~128 blocks on a default geth), not by the trusting period: the client sits ~64 blocks back by construction (it only follows finalized headers) and drifts up to another interval between ticks, so worst case is ~104 blocks behind. Shorten the interval before raising the lag threshold.
  - Staleness is not merely latency: once the pinned block leaves the state window, `eth_getProof` there fails, so the *next* packet cannot be relayed either. `historical state ... is not available` is a symptom of a stalled client, never of insufficient L1 retention.
- Sync-committee period crossing: the Ethereum client update path takes the next period's committee from the preceding period's light-client update (`next_sync_committee`, proven by `next_sync_committee_branch`), fetching that single update when it falls outside the requested range. It must NOT come from `light_client/bootstrap` — beacon nodes serve bootstraps only for the checkpoint roots they retain, so the first boundary answers `404 NOT_FOUND: Sync committee for period N not found` and the client stops advancing for good. Periods roll roughly every 27 hours on mainnet, which makes this a client-expiry path.
- Bounded proof assembly: `chain/l2rollup.Builder.Build` puts a `headerBuildTimeout` (90s) around the whole L2 proof assembly. go-ethereum's HTTP client carries no timeout of its own, and the relay module drives one direction on a single goroutine, so one unanswered RPC wedges that entire direction with no error, no retry and nothing in the log. The relay context alone is not enough — it is only cancelled at shutdown. A direction that goes silent while others stay healthy is this failure: `kill -QUIT <pid>` and look for a goroutine blocked in `net/http.(*persistConn).roundTrip`.
- Validator churn / multi-hop recovery (RLY-01): every Cosmos→ETH client update proof still requires >2/3 of the current pinned validator set among the target block's signers. When validator churn breaks that at the chain tip, `BuildCosmosClientUpdateMsg` binary-searches backward from the tip (`findHighestFeasibleHop`) for the highest intermediate height that still clears pinned-set quorum, updates and unconditionally rotates the pinned set to that height's validator set (`RLY01_HOP_FALLBACK`/`RLY01_HOP_FOUND`/`RLY01_HOP_ROTATE` log lines), and lets the caller repeat from there. `RefreshCosmosClient` and `update-client --max-hops` repeat only while the builder reports that it used this hop fallback, up to `maxHopsPerRefresh` / `--max-hops` (16 by default); an ordinary direct update stops after one transaction even if the source chain has produced a newer block by then. `--target-height`/`--force-rotation` remain operator overrides for recovery. If not even one block past the trusted height clears quorum, the client cannot advance (`RLY01_HOP_EXHAUSTED`) and needs manual intervention. `RLY01_QUORUM_WARN` fires as a leading indicator before an update actually fails.


## Observability

### Logging

See `docs/logging.md` for full guide. Key principles:
- Structured JSON to stdout/stderr
- OpenTelemetry semantic conventions
- Correlation fields: `trace_id`, `service_name`, `service_version`, `environment`
- Handle errors once: log OR return, not both

### Metrics

See `docs/metrics.md` for full guide. Focus on:
- RED metrics (Rate, Errors, Duration) per service
- Connection health to RPC nodes and dependencies
- Error categorization (timeout, validation, network)

## Known Risks

| Risk | Impact | Mitigation |
|------|--------|-----------|
| CometBFT WebSocket disconnect | Missed packets | `SubscribeCosmos` retries failed subscriptions and periodically backfills missed Cosmos events with `TxSearch` |
| Ethereum RPC rate limiting | Delayed relaying | Configurable retry backoff |
| Groth16 proof generation timeout | Stalled client update | Relayer restart |
| Circuit artifact mismatch | Proof verification failure | Per-bucket `bin/n{N}/{r1cs,pk,vk}.bin`; redeploy `Groth16Verifier_N{N}.sol` + SignatureVerifier bucket registry whenever artifacts are regenerated |
| Quorum exceeds largest bucket | Cosmos→ETH updates stall | Add a larger entry to `prover.Buckets`, recompile via `prover/cmd`, redeploy verifiers |
| `go.mod` replace directive | Build failure on new dev machine | Document local path setup |
| ETH→Cosmos relay | Relies on Ethereum event subscription + beacon finality availability | Fully implemented via the `evm` source + `cosmos` destination adapters: `RelayableHeight` gates on the finalized exec block, the `beacon` builder + `cosmos` destination advance the 08-wasm client (`WaitForCosmosCatchUp`), the ETH storage proof is fetched by `evm.Source.MembershipProof`, and `MsgRecvPacket` / `MsgAcknowledgement` are broadcast to Cosmos; terminal `EthAck` / `EthTimeout` clear ETH pending state and are otherwise log-only |
| Ethereum sync committee period crossing | Stale Ethereum light client on Cosmos | Multi-period update logic in routine.go handles period boundary transitions |
| Cosmos-originated packet expires on ETH | `MsgRecvPacket` reverts with `IBCInvalidTimeoutTimestamp`; commitment stuck on Cosmos | `scanForCosmosTimeouts` (30 s tick) drains `PendingPacketTracker`, first drops packets already received on ETH, then builds non-membership proof on ETH and bundles `MsgUpdateClient` + `MsgTimeout`s into one Cosmos tx; `PurgeStaleWithoutTimeout(1h)` only removes entries that have no timeout so timeout obligations are not dropped |
| ETH-originated packet expires on Cosmos | ETH source commitment remains locked if the packet left the in-memory relay queue | `EthPendingTracker` is populated from live/recovered ETH `SendPacket` events; `scanForEthTimeouts` (30 s tick) submits `timeoutPacket` on ETH once the Cosmos proof state reaches the timeout |
| Subscriber timeout-timestamp unit drift | `ibc-go` emits TimeoutTimestamp in ns; ETH expects seconds | `normalizeTimeoutSeconds` at subscriber boundary divides by 1e9 when value > 1e12 |

## Health Checks

- **Relayer**: Logs startup success, connection to both chains
- **Contracts**: `SpectreClient` stores `latestHeight` — compare against chain head to detect staleness
- **E2E tests**: 5 interchaintest suites validate full round-trip functionality

## Recovery Procedures

- **Stale client**: Re-run relayer with fresh headers to update client
- **Frozen client**: Requires governance action (admin upgrade or new client deployment)
- **Nonce error**: Relayer re-queries account sequence and retries
- **Counterparty upgrade / hard-fork**: If a counterparty chain undergoes a hard-fork or client upgrade that changes light client rules, the existing Tendermint light client contract will be bricked. Because the `upgradeClient` interface is not supported, recovery requires deploying a new light client contract instance, registering it in the router, and re-configuring the relayer to use the new client ID.
- **Circuit update**: From `relayer/`, run `go run ./prover/cmd ./bin ../contracts/verifiers` to regenerate every bucket's artifacts and emit fresh `Groth16Verifier_N{N}.sol`; for GPU proving use `go run -tags=icicle ./prover/cmd -gpu-prove ./bin ../contracts/verifiers`; redeploy each per-bucket verifier and re-register them via `SignatureVerifier.setBucket(...)` before the next E2E run
