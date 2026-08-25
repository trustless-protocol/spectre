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
| Compatibility | `test/compatibility/` | `scripts/solidity/checker-self-test.sh` |

The handwritten format gate is `scripts/check-solidity-format.sh`; generated verifier implementations are excluded. The complete structural gate is `scripts/check-solidity-compatibility.sh`, which enforces ABI, storage, bytecode, hot-path gas, fixture, binding, verifier, test-inventory, architecture, and protected-source boundaries.

## Foundry Configuration

`foundry.toml` pins Solidity 0.8.28, Cancun, optimizer runs 1,000, `via_ir = true`, no bytecode metadata hash, and 100,000 local fuzz runs. The observed local tool versions are recorded in `scripts/solidity/toolchain-manifest.json`.

## Generated Prover Boundary

The checked local generator/prover manifest enables N4 only. Run `scripts/solidity/build-prover-artifacts.sh` to generate R1CS, proving key, verifying key, and Solidity verifier into ignored staging, smoke-test the pair, publish N4, and record runtime provenance. Production deployment preserves the six-bucket N ∈ {4, 8, 16, 32, 64, 128} compatibility surface; every enabled bucket needs its own coordinated artifact set. Groth16 setup is randomized, so artifacts from different setup runs must never be mixed.

## Cross-Language Consumers

```bash
cargo test --locked -p ibc-eureka-solidity-types --all-features
cd packages/go-abigen && go test ./...
cd relayer && LD_LIBRARY_PATH=../third_party/ecip-gnark go test ./...
```

ABI and binding generation is controlled by `scripts/solidity/contracts.json` and requires `abigen` 1.17.2-stable. Solidity fixture staging uses `test/fixtures/solidity/` and Spectre fixture staging uses `test/fixtures/spectre/`.

## Manual Release Gates

This Solidity refactor does not add required CI. Run the commands above and `scripts/check-solidity-format.sh`, `forge build --force --skip test --skip script --sizes`, and `scripts/solidity/checker-self-test.sh` locally. Record exact commands, exit codes, sizes, and compatibility digests in the PR body and release report.
