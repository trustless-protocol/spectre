# Quality and Testing

## Solidity (Foundry)

| Scope | Location | Command |
|---|---|---|
| Core | `test/core/` | `forge test --match-path 'test/core/*'` |
| ICS-20 | `test/apps/ics20/` | `forge test --match-path 'test/apps/ics20/*'` |
| Spectre | `test/light-clients/spectre/` | `forge test --match-path 'test/light-clients/spectre/*'` |
| Integration | `test/integration/` | `forge test --match-path 'test/integration/*'` |
| Deployment | `test/deployment/` | `forge test --match-path 'test/deployment/*'` |
| Full non-shadowfork | `test/` | `forge test --no-match-path 'test/shadowfork/*'` |
| Compatibility | `test/compatibility/` | `scripts/solidity-refactor/checker-self-test.sh` |

The handwritten format gate is `scripts/check-solidity-format.sh`; generated verifier implementations are excluded. The complete structural gate is `scripts/check-solidity-compatibility.sh`, which enforces ABI, storage, bytecode, hot-path gas, fixture, binding, verifier, test-inventory, architecture, and protected-source boundaries.

## Foundry Configuration

`foundry.toml` pins Solidity 0.8.28, Cancun, optimizer runs 1,000, `via_ir = true`, no bytecode metadata hash, and 100,000 local fuzz runs. CI pins Foundry v1.7.1.

## Generated Prover Boundary

The supported verifier set is N4 only. Run `scripts/solidity-refactor/build-prover-artifacts.sh` to generate R1CS, proving key, verifying key, and Solidity verifier into ignored staging, smoke-test the pair, publish only N4, and record runtime provenance. Groth16 setup is randomized, so artifacts from different setup runs must never be mixed.

## Cross-Language Consumers

```bash
cargo test --locked -p ibc-eureka-solidity-types --all-features
cd packages/go-abigen && go test ./...
cd relayer && LD_LIBRARY_PATH=../third_party/ecip-gnark go test ./...
```

ABI and binding generation is controlled by `scripts/solidity-refactor/contracts.json` and requires `abigen` 1.17.2-stable. Solidity fixture staging uses `test/fixtures/solidity/` and Spectre fixture staging uses `test/fixtures/spectre/`.

## Automated CI

`.github/workflows/solidity.yml` runs on pushes to `main`, pull requests, and manual dispatch. It checks out the private prover submodules with `SUBMODULE_TOKEN`, builds the Garaga FFI, generates a paired N4 set, checks formatting, force-builds sizes, runs Foundry and Rust, and executes the compatibility self-tests.

The existing Go and E2E workflows retain their documented trigger policies.
