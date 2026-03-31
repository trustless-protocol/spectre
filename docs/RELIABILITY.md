# Reliability

## Error Handling

### Solidity

- Custom errors defined in `contracts/errors/` (gas-efficient vs string reverts)
- `require()` with descriptive messages for input validation
- Proof verification failures revert the entire transaction (atomic)

### Go Operator

- `log.Fatal` for unrecoverable startup errors
- Service loops return errors up the chain
- Transaction handler retries with re-queried account sequence on nonce conflicts
- `BatchBuilder` uses mutex-guarded packet queue for concurrent safety

### Rust Relayer

- Structured logging via `tracing` crate (see `docs/logging.md`)
- Error propagation with `anyhow`/`thiserror`

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
| Groth16 proof generation timeout | Stalled client update | Operator restart |
| Circuit artifact mismatch | Proof verification failure | Versioned `r1cs.bin`, `pk.bin`, `vk.bin` |
| `go.mod` replace directive | Build failure on new dev machine | Document local path setup |

## Health Checks

- **Operator**: Logs startup success, connection to both chains
- **Contracts**: `SP1ICS07Tendermint` stores `latestHeight` — compare against chain head to detect staleness
- **E2E tests**: 5 interchaintest suites validate full round-trip functionality

## Recovery Procedures

- **Stale client**: Re-run operator with fresh headers to update client
- **Frozen client**: Requires governance action (admin upgrade or new client deployment)
- **Nonce error**: Operator re-queries account sequence and retries
- **Circuit update**: Run `go run ./prover/cmd/ <output_dir>` to regenerate artifacts, redeploy `Groth16Verifier.sol`
