# Solidity Refactor Consumer Report

| Consumer | Validation | Result |
|---|---|---|
| Rust Solidity types | `cargo test --locked -p ibc-eureka-solidity-types --all-features` | Pass |
| Shared Go bindings | `cd packages/go-abigen && go test ./...` | Pass |
| Relayer Go | `cd relayer && LD_LIBRARY_PATH=../third_party/ecip-gnark go test ./...` | Pass |
| E2E Go packages | `cd e2e/interchaintestv8 && go test -run '^$' ./...` | Compile pass; no chains started |
| Shell deployment callers | `bash -n scripts/local/deploy_eth_contracts.sh scripts/local/deploy_l2_contracts.sh ...` | Pass |
| Go bindings | pinned manifest regeneration with abigen 1.17.2-stable | No tracked Go binding diff |
| Spectre bytecode artifact | manifest regeneration | ABI/initcode/runtime identical; fully qualified source metadata changed as expected |

The generator resolves `source:contract` pairs from `contracts.json`, verifies each Foundry artifact's metadata source, stages temporary ABI/bin files outside `contracts/data`, and rejects unlinked placeholders.

The binding regeneration, Rust, shared Go binding, relayer Go, and E2E compile-only checks were repeated after the final forced build. Both passes exited 0 and produced the same tracked binding digests.
