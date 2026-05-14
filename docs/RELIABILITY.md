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
- `BatchBuilder` manages packet batching with separate mutexes for Cosmos / ETH queues; also owns a `PendingPacketTracker` for Cosmos-originated packets awaiting ETH delivery
- Timeout checking: packets past their timeout are skipped before submission; subscriber normalises IBC v2 timeouts from ns → s at ingest (`normalizeTimeoutSeconds`)
- Cosmos→ETH timeout fallback: `scanForCosmosTimeouts` goroutine (every 30 s) scans `PendingPacketTracker`, builds `MsgTimeout` for each expired entry, batches all timeouts plus a `MsgUpdateClient` into one Cosmos tx, and `Recover()`s from panic so a single bad packet can't kill the goroutine


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
| CometBFT WebSocket disconnect | Missed packets | Subscriber reconnection logic |
| Ethereum RPC rate limiting | Delayed relaying | Configurable retry backoff |
| Groth16 proof generation timeout | Stalled client update | Relayer restart |
| Circuit artifact mismatch | Proof verification failure | Per-bucket `bin/n{N}/{r1cs,pk,vk}.bin`; redeploy `Groth16Verifier_N{N}.sol` + WrapperVerifier bucket registry whenever artifacts are regenerated |
| Quorum exceeds largest bucket | Cosmos→ETH updates stall | Add a larger entry to `prover.Buckets`, recompile via `prover/cmd`, redeploy verifiers |
| `go.mod` replace directive | Build failure on new dev machine | Document local path setup |
| ETH→Cosmos relay | Relies on Ethereum event subscription + beacon finality availability | Fully implemented (`handleEth`): poll `GetFinalityUpdate` until exec-block finalised, `UpdateEthClient` on Cosmos, fetch ETH storage proof, broadcast `MsgRecvPacket` / `MsgAcknowledgement`; terminal `EthAck` / `EthTimeout` only logged |
| Ethereum sync committee period crossing | Stale Ethereum light client on Cosmos | Multi-period update logic in routine.go handles period boundary transitions |
| Cosmos-originated packet expires on ETH | `MsgRecvPacket` reverts with `IBCInvalidTimeoutTimestamp`; commitment stuck on Cosmos | `scanForCosmosTimeouts` (30 s tick) drains `PendingPacketTracker`, builds non-membership proof on ETH, bundles `MsgUpdateClient` + `MsgTimeout`s into one Cosmos tx; `PurgeStale(1h)` GC prevents tracker bloat if anything leaks |
| Subscriber timeout-timestamp unit drift | `ibc-go` emits TimeoutTimestamp in ns; ETH expects seconds | `normalizeTimeoutSeconds` at subscriber boundary divides by 1e9 when value > 1e12 |

## Health Checks

- **Relayer**: Logs startup success, connection to both chains
- **Contracts**: `Groth16ICS07Tendermint` stores `latestHeight` — compare against chain head to detect staleness
- **E2E tests**: 5 interchaintest suites validate full round-trip functionality

## Recovery Procedures

- **Stale client**: Re-run relayer with fresh headers to update client
- **Frozen client**: Requires governance action (admin upgrade or new client deployment)
- **Nonce error**: Relayer re-queries account sequence and retries
- **Circuit update**: From `relayer/`, run `go run ./prover/cmd ./bin ../contracts/verifiers` to regenerate every bucket's artifacts and emit fresh `Groth16Verifier_N{N}.sol`; redeploy each per-bucket verifier and re-register them via `WrapperVerifier.setBucket(...)` before the next E2E run
