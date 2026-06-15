# Reliability

## Error Handling

### Solidity

- Custom errors defined in `contracts/errors/` (gas-efficient vs string reverts)
- `require()` with descriptive messages for input validation
- Proof verification failures revert the entire transaction (atomic)

### Go Relayer

- `log.Fatal` for unrecoverable startup errors
- `start` command: uses `services.StartLoop()` with batch processing (size/time thresholds), bi-directional relay (Cosmos ↔ ETH)
- `create-clients` command: one-time setup, fails fast on deployment errors
- Transaction handler retries with re-queried account sequence on nonce conflicts
- `BatchBuilder` manages packet batching with separate mutexes for Cosmos / ETH queues; it owns separate pending trackers for Cosmos-originated packets awaiting ETH delivery and ETH-originated packets awaiting Cosmos delivery/ack/timeout
- `SubscribeCosmos` performs startup and periodic gap recovery with CometBFT `TxSearch` for `send_packet`, `write_acknowledgement`, and `timeout_packet` event payloads. `COSMOS_STARTUP_LOOKBACK_BLOCKS=0` (default) scans the full indexed history; set it to a positive block count to bound startup work.
- `SubscribeEth` performs startup/reconnect log recovery for `SendPacket` and `WriteAcknowledgement`. `ETH_STARTUP_LOOKBACK_BLOCKS` defaults to 256; set it at least as large as the maximum expected relayer downtime if older ETH logs must be recovered.
- Timeout checking: packets past their timeout are skipped before submission; subscriber normalises IBC v2 timeouts from ns → s at ingest (`normalizeTimeoutSeconds`)
- Cosmos→ETH timeout fallback: `scanForCosmosTimeouts` goroutine (every 30 s) scans `PendingPacketTracker`, builds `MsgTimeout` for each expired entry, batches all timeouts plus a `MsgUpdateClient` into one Cosmos tx, and `Recover()`s from panic so a single bad packet can't kill the goroutine
- ETH→Cosmos timeout fallback: `scanForEthTimeouts` goroutine (every 30 s) scans `EthPendingTracker`, re-checks the ETH source commitment, then reuses the `timeoutEthSend` proof flow for expired ETH-originated packets
- `handleCosmos` skips the auto-relay of a Cosmos-emitted `timeout_packet` event back to ETH for Cosmos-originated packets (`shouldRelayCosmosTimeoutToEth`): the source-side MsgTimeout on Cosmos already refunded the sender and ETH never held a commitment, so calling `ICS26Router.timeoutPacket` would be a NoOp at best
- Cross-chunk dedup: `waitCosmosAppHash` and `waitBeaconFinality` each memoize their last-confirmed height on the `Services` struct; chunked flushes (`BatchSize` chunking) reuse the wait for the same source block instead of repeating it per chunk
- updateClient idempotence: `BuildCosmosClientUpdateMsg` always re-queries the on-chain ICS07 client state and short-circuits to `HasMsg=false` when on-chain trusted height already covers the chunk's source block — avoids regenerating a Groth16 proof that another chunk (or a frontrunner relayer) already submitted


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
| Circuit artifact mismatch | Proof verification failure | Per-bucket `bin/n{N}/{r1cs,pk,vk}.bin`; redeploy `Groth16Verifier_N{N}.sol` + WrapperVerifier bucket registry whenever artifacts are regenerated |
| Quorum exceeds largest bucket | Cosmos→ETH updates stall | Add a larger entry to `prover.Buckets`, recompile via `prover/cmd`, redeploy verifiers |
| `go.mod` replace directive | Build failure on new dev machine | Document local path setup |
| ETH→Cosmos relay | Relies on Ethereum event subscription + beacon finality availability | Fully implemented (`handleEth`): poll `GetFinalityUpdate` until exec-block finalised, `UpdateEthClient` on Cosmos, fetch ETH storage proof, broadcast `MsgRecvPacket` / `MsgAcknowledgement`; terminal `EthAck` / `EthTimeout` clear ETH pending state and are otherwise log-only |
| Ethereum sync committee period crossing | Stale Ethereum light client on Cosmos | Multi-period update logic in routine.go handles period boundary transitions |
| Cosmos-originated packet expires on ETH | `MsgRecvPacket` reverts with `IBCInvalidTimeoutTimestamp`; commitment stuck on Cosmos | `scanForCosmosTimeouts` (30 s tick) drains `PendingPacketTracker`, builds non-membership proof on ETH, bundles `MsgUpdateClient` + `MsgTimeout`s into one Cosmos tx; `PurgeStale(1h)` GC prevents tracker bloat if anything leaks |
| ETH-originated packet expires on Cosmos | ETH source commitment remains locked if the packet left the in-memory relay queue | `EthPendingTracker` is populated from live/recovered ETH `SendPacket` events; `scanForEthTimeouts` (30 s tick) submits `timeoutPacket` on ETH once the Cosmos proof state reaches the timeout |
| Subscriber timeout-timestamp unit drift | `ibc-go` emits TimeoutTimestamp in ns; ETH expects seconds | `normalizeTimeoutSeconds` at subscriber boundary divides by 1e9 when value > 1e12 |

## Health Checks

- **Relayer**: Logs startup success, connection to both chains
- **Contracts**: `Groth16ICS07Tendermint` stores `latestHeight` — compare against chain head to detect staleness
- **E2E tests**: 5 interchaintest suites validate full round-trip functionality

## Recovery Procedures

- **Stale client**: Re-run relayer with fresh headers to update client
- **Frozen client**: Requires governance action (admin upgrade or new client deployment)
- **Nonce error**: Relayer re-queries account sequence and retries
- **Counterparty upgrade / hard-fork**: If a counterparty chain undergoes a hard-fork or client upgrade that changes light client rules, the existing Tendermint light client contract will be bricked. Because the `upgradeClient` interface is not supported, recovery requires deploying a new light client contract instance, registering it in the router, and re-configuring the relayer to use the new client ID.
- **Circuit update**: From `relayer/`, run `go run ./prover/cmd ./bin ../contracts/verifiers` to regenerate every bucket's artifacts and emit fresh `Groth16Verifier_N{N}.sol`; for GPU proving use `go run -tags=icicle ./prover/cmd -gpu-prove ./bin ../contracts/verifiers`; redeploy each per-bucket verifier and re-register them via `WrapperVerifier.setBucket(...)` before the next E2E run
