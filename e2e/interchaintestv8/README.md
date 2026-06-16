# End-to-end testing suite (interchaintest)

The e2e tests are built on the [interchaintest](https://github.com/cosmos/interchaintest)
library. Each run boots a real Ethereum PoS devnet (Kurtosis + ethereum-package),
a real IBC-Go simd chain (Docker), starts the fast-ibc relayer binary against
them, then drives ICS-20 packet flows end-to-end.

Each test takes **8–15 minutes** wall-time. Designed for both CI and local dev.

## Prerequisites

On top of the toolchain in the [repo root README](../../README.md#requirements):

- **Docker** daemon running
- **[Kurtosis CLI](https://docs.kurtosis.com/install/)** installed
- The two sibling repos (`ecip-gnark`, `decentrio-gnark`) cloned per
  [Sibling repos](../../README.md#sibling-repos-required-to-build-the-relayer)
- `relayer/prover/buckets.go` trimmed to `var Buckets = []int{4}` for local
  runs — the suite only ever exercises bucket 4 (single-validator simd chain),
  and building the full set takes 30+ minutes (n=128 alone ≈ 30 min). Revert
  before pushing — production needs the full set.

The first invocation of `just test-e2e` auto-runs:

- `install-go-relayer` — `go build -o $GOPATH/bin/relayer ./cmd`
- `build-prover-artifacts` — `go run ./prover/cmd ./bin ../contracts/verifiers`
  (compiles every bucket in `Buckets`, emits per-bucket `r1cs/pk/vk.bin` and
  `contracts/verifiers/Groth16Verifier_N{N}.sol`)
- `clean-foundry` — wipes `cache/`, `out/`, `broadcast/`

Subsequent runs reuse the cached prover artifacts (a `relayer/bin/n4/vk.bin`
probe keeps `build-prover-artifacts` a no-op).

## Run a test

From the repo root:

```bash
# Single test by full name (SUITE/CASE)
just test-e2e TestWithIbcEurekaTestSuite/Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack

# Shortcut by suite (drops the TestWith...Suite/ prefix)
just test-e2e-eureka Test_Deploy
just test-e2e-relayer Test_RelayerInfo
just test-e2e-cosmos-relayer Test_ICS20RecvAndAckPacket
just test-e2e-multichain Test_Deploy
```

The recipe sets `ETH_TESTNET_TYPE=pos`, `RELAYER_BINARY=$GOPATH/bin/relayer`,
and `PROVER_BIN_DIR=<repo>/relayer/bin` for you.

## Currently passing `IbcEurekaTestSuite` cases

| Test | What it verifies |
|------|------------------|
| `Test_Deploy` | contracts deployed, wasm-eth client + ICS07 + AddClient + counterparty register all wired |
| `Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` | full ETH→Cosmos→ETH round-trip (5 phases: recv, ack, return-send, return-recv, return-ack) |
| `Test_ICS20TransferNativeCosmosCoinsToEthereumAndBack` | full Cosmos→ETH→Cosmos round-trip for native coin |
| `Test_ICS20TransferLargeAmountFromEthereumToCosmosAndBack` | bigint round-trip (`transferAmount = 5×10²²`, half of `StartingERC20Balance`) |
| `Test_TimeoutPacketFromEth` | short-deadline ETH→Cosmos packets settle by either auto-delivery or timeout refund with balanced cross-chain accounting |
| `Test_TimeoutPacketFromCosmos` | short-deadline Cosmos→ETH packets settle by either auto-delivery or timeout refund with balanced cross-chain accounting |
| `Test_ErrorAckToEthereum` | Cosmos transfer module writes error ack → `MsgAck` on ETH → escrow refund |

## Tests that don't run against fast-ibc's daemon

These cases assume the upstream Rust relayer's gRPC `RelayByTx` model and have
no auto-relay equivalent. They skip at runtime and are excluded from the
reusable E2E matrix pending a refactor.

- `Test_5_FilteredTimeoutPacketFromEth`, `Test_10_FilteredTimeoutPacketFromCosmos`,
  and other partial-timeout-filter variants — fast-ibc's scanner can't
  selectively skip expired packets; it times all of them out together.
- `Test_25_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` (and the 50-packet
  variant) — relayer signs and submits each Cosmos→ETH `recvPacket` sequentially,
  so 25 packets × ~6 s each + initial beacon-finality wait blows past the
  5-min `Eventually` timeout in the phase-4 wait. Pending batching improvements.
- `Test_ICS20TransferERC20TokenFromEthereumToCosmosAndBackFails` and
  `Test_5_FinalizedTimeoutPacketFromEth` — still use gRPC `RelayByTx` to retrieve
  raw relay transactions for manual broadcast.
- `Test_TimeoutPacketEthRemintsVouchers` and
  `Test_TimeoutPacketCosmosRemintsVouchers` — still mix auto-relay setup with
  manual `RelayByTx` timeout/ack broadcasts.
- `groth16_ics07_test.go` (fixture-generation flow) — references the deleted
  `fixtures membership` CLI subcommand.
- `cosmos_relayer_test.go` (Cosmos↔Cosmos via the Rust gRPC relayer) — out of
  scope for fast-ibc (ETH↔Cosmos only).
- `multichain_test.go` `MultichainTestSuite` cases — still rely on the upstream
  gRPC `CreateClient`, `RelayByTx`, and `Info` service plus per-chain Cosmos
  signing semantics. fast-ibc's Go daemon currently runs a single configured
  Cosmos↔ETH path without exposing that gRPC API.
- `relayer_test.go` `RelayerTestSuite` cases — still use gRPC `RelayByTx` +
  manual broadcast pattern.

## Architecture note

fast-ibc's relayer is an **event-driven daemon**, not the upstream Rust gRPC
service. Tests therefore *wait* for the daemon to drive each step (via
`require.Eventually` polling the post-relay end state) instead of pulling a
relay tx via `s.RelayerClient.RelayByTx(...)` and broadcasting it themselves.

The `s.RelayerClient` / `s.EthRelayerSubmitter` plumbing kept from upstream
is still wired in `e2esuite/suite.go` for the subset of tests that use raw-tx
broadcasts, but the IbcEurekaTestSuite cases above avoid them entirely.
