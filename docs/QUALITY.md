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
| Groth16 light client | `test/groth16-ics07/` | `forge test --match-path test/groth16-ics07/` |

Run a single test:
```bash
forge test --match-test testSendTransfer -vvv
```

### Go (Relayer)

| Type | Location | Command |
|------|----------|---------|
| Unit tests | `relayer/*/` | `cd relayer && go test ./...` |
| Race detection | — | `cd relayer && go test -race ./...` |
| Single package | — | `cd relayer && go test -v ./prover/...` |
| Single test | — | `cd relayer && go test -run TestName ./pkg/...` |

Test files: `prover/`, `client/`, `subscriber/`, `services/`, `keys/`, `utils/`

### Go (E2E)

| Suite | Test File | Command |
|-------|-----------|---------|
| IBC Eureka | `ibc_eureka_test.go` | `just test-e2e-eureka` |
| Relayer | `relayer_test.go` | `just test-e2e-relayer` |
| Cosmos Relayer | `cosmos_relayer_test.go` | `just test-e2e-cosmos-relayer` |
| Groth16 ICS07 | `groth16_ics07_test.go` | `just test-e2e-groth16-ics07` |
| Multi-chain | `multichain_test.go` | `just test-e2e-multichain` |

Requires: Docker Desktop, Kurtosis, compiled relayer binary, Groth16 network key.

### Rust

```bash
just test-cargo              # All Rust tests
just test-cargo <name>       # Single test
```

## Encoding Cross-Validation

`Encode.sol` and `Header.sol` are cross-validated against Go `proto.Marshal()`
via `forge test --match-contract EncodeTest -vvv`. Test fixtures contain the
expected Go-reference hex; any encoding change MUST keep these tests green.

## Foundry Configuration

From `foundry.toml`:
- Solidity `0.8.28`, EVM `cancun`
- Optimizer: 10,000 runs with `--via-ir`
- Fuzz: 100,000 runs locally, 5,000 in CI
- Fixed block timestamp for reproducibility

## Test Fixtures

Pre-generated Groth16 proofs in `test/solidity-ibc/fixtures/` and `test/groth16-ics07/fixtures/`. Regenerate with:
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
