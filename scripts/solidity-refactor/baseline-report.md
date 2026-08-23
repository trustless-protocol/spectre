# Solidity Refactor Baseline Report

- Baseline commit: `b75a613a2d020b1e9717df672a066788bf9ccf30` (`origin/main` at branch creation).
- Toolchain: Solidity 0.8.28, Cancun, optimizer 1,000, via IR; local Forge `1.6.0-v1.7.0`; Rust 1.94.1; Go 1.25.9; abigen 1.17.2-stable.
- Initial blockers reproduced: missing generated N4 verifier and stale Rust Solidity message includes/conversions.
- Rust repair: restored the pinned legacy Groth16 ABI input and reconciled includes/conversions with current Solidity types. `cargo test --locked -p ibc-eureka-solidity-types --all-features` passes.
- Prover repair: N4 R1CS, PK, VK, and Solidity verifier are generated and smoke-tested as one randomized pair. Runtime provenance is written under ignored `.artifacts/solidity-refactor/prover/`.
- Baseline Foundry: 33 suites, 303 tests passed, 0 failed, 0 skipped with 100,000 fuzz runs.
- Forced baseline sizes: ICS26Router 24,465 B runtime (111 B margin), ICS20Transfer 22,224 B (2,352 B margin), SpectreClient 19,390 B (5,186 B margin).
- Frozen evidence: canonical ABI, nine first-party ERC-7201 namespace/struct layouts, initcode, runtime bytecode, hot-path gas ceilings, fixtures, Go bindings, verifier pair, test inventory, and toolchain/lock digests in `baseline.json` and companion manifests.

Commands:

```bash
forge test --no-match-path 'test/shadowfork/*'
forge build --force --skip test --skip script --sizes
cargo test --locked -p ibc-eureka-solidity-types --all-features
scripts/solidity-refactor/build-prover-artifacts.sh
```
