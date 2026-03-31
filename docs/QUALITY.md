# Quality & Testing

## Test Strategy

### Solidity (Foundry)

| Type | Location | Command |
|------|----------|---------|
| Unit tests | `test/solidity-ibc/` | `just test-foundry` |
| Fixture tests | `test/solidity-ibc/FixtureTest.t.sol` | `forge test --match-contract FixtureTest` |
| Encoding cross-validation | `test/solidity-ibc/EncodeTest.t.sol` | `forge test --match-contract EncodeTest` |
| Gas benchmarks | `test/solidity-ibc/BenchmarkTest.t.sol` | `just test-benchmark <name>` |
| Shadowfork tests | `test/shadowfork/` | Requires `ETH_RPC_URL` env var |
| SP1 light client | `test/sp1-ics07/` | `forge test --match-path test/sp1-ics07/` |

Run a single test:
```bash
forge test --match-test testSendTransfer -vvv
```

### Go (Operator)

| Type | Location | Command |
|------|----------|---------|
| Unit tests | `operator/*/` | `cd operator && go test ./...` |
| Race detection | — | `cd operator && go test -race ./...` |
| Single package | — | `cd operator && go test -v ./prover/...` |
| Single test | — | `cd operator && go test -run TestName ./pkg/...` |

Test files: `prover/`, `client/`, `subscriber/`, `services/`, `keys/`, `utils/`

### Go (E2E)

| Suite | Test File | Command |
|-------|-----------|---------|
| IBC Eureka | `ibc_eureka_test.go` | `just test-e2e-eureka` |
| Relayer | `relayer_test.go` | `just test-e2e-relayer` |
| Cosmos Relayer | `cosmos_relayer_test.go` | `just test-e2e-cosmos-relayer` |
| SP1 ICS07 | `sp1_ics07_test.go` | `just test-e2e-sp1-ics07` |
| Multi-chain | `multichain_test.go` | `just test-e2e-multichain` |

Requires: Docker Desktop, Kurtosis, compiled operator/relayer binaries, SP1 network key.

### Rust

```bash
just test-cargo              # All Rust tests
just test-cargo <name>       # Single test
```

## Encoding Cross-Validation

`Encode.sol` and `Header.sol` are cross-validated against Go `proto.Marshal()`:

1. `cd operator && go run ./cmd/encode_debug/` — outputs Go reference hex
2. `forge test --match-contract EncodeTest` — 34 tests compare Solidity output vs Go reference
3. Any encoding change MUST pass both Go and Solidity validation

## Foundry Configuration

From `foundry.toml`:
- Solidity `0.8.28`, EVM `cancun`
- Optimizer: 10,000 runs with `--via-ir`
- Fuzz: 100,000 runs locally, 5,000 in CI
- Fixed block timestamp for reproducibility

## Test Fixtures

Pre-generated SP1 proofs in `test/solidity-ibc/fixtures/` and `test/sp1-ics07/fixtures/`. Regenerate with:
```bash
just generate-fixtures-solidity
```

## Linting

```bash
just lint                 # All linters
just lint-solidity        # forge fmt + solhint + natlint
just lint-go              # golangci-lint
just lint-rust            # cargo fmt + cargo clippy
just lint-buf             # Protobuf linting
```

## Security

```bash
just slither              # Slither static analysis
```

Config in `.slither.config.json`: excludes low/informational findings and dependencies.

## CI/CD

8 GitHub workflows in `.github/workflows/`:
- `foundry.yml` — Lint + Solidity unit tests on PR/push
- `rust.yml` — Cargo tests + clippy
- `e2e-minimal.yml` / `e2e-full.yml` / `e2e-mock.yml` — E2E suites
- `docker.yml` — Relayer Docker image build
- `release.yml` — Release automation
- `abigen.yaml` — Go bindings generation
