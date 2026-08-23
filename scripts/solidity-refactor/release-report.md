# Solidity Refactor Release Report

## Outcome

The handwritten Solidity source tree now has explicit core, ICS-20, light-client compatibility, Spectre, shared, generated-verifier, and periphery owners. Tests mirror those packages. Global `contracts/utils`, `contracts/interfaces`, `contracts/msgs`, and `contracts/errors` no longer contain Solidity sources.

The final semantic diff review found no protocol behavior change. Non-mechanical production edits are limited to the reviewed `BytesPrefix` shared extraction and the selector-manifest split; exact runtime bytecode, ABI, storage, protocol fixtures, and gas measurements demonstrate equivalence. Generated verifier implementations were not hand-edited, and no unrelated tracked change was identified.

## Compatibility

An isolated build of origin/main and the refactored tree matched exactly for canonical ABI, creation bytecode, and deployed bytecode across 28 entrypoints, modules, stores, and libraries; nine first-party ERC-7201 namespace constants and ordered struct layouts also match exactly. The compatibility checker also locks hot-path gas ceilings, protocol fixtures, tracked ABI, Go bindings, verifier provenance, EIP-170, and test inventory. Its negative suite rejects ABI, storage, runtime, gas, fixture/binding, test-inventory, fully qualified artifact, verifier-provenance, and architecture drift.

The only accepted generated artifact delta is fully qualified source/build metadata in `abi/bytecode/SpectreClient.json`; ABI, initcode, and runtime objects remain exact.

## Deployment and Prover

Production deploy/verify scripts require and register only `VERIFIER_N4`, matching `prover.Buckets` and the verifier manifest. Unsupported local N8/N16 verifier sources and N8/N16/N32/N64 key sets were moved into ignored, recoverable `.artifacts/solidity-refactor/prover/unsupported-archive/`.

Groth16 setup is randomized. `build-prover-artifacts.sh` stages, smoke-tests, publishes, and records one paired set. Gnark still emits its existing hash-to-field exporter warning; this release does not claim a live on-chain proof E2E.

## Size and Optimization Decision

Fresh sizes are unchanged: ICS26Router 24,465 B (111 B margin), ICS20Transfer 22,224 B (2,352 B), SpectreClient 19,390 B (5,186 B). Reaching 1 KiB router margin requires 913 safe bytes. No reviewed candidate met that threshold, so no optional bytecode optimization was attempted.

## Repeated Release Matrix

| Gate | Pass 1 | Pass 2 |
|---|---|---|
| Binding regeneration | Exit 0; tracked Go bindings unchanged | Exit 0; tracked Go bindings unchanged |
| Full non-shadowfork Foundry | Exit 0; 34 suites, 305 passed, 0 failed, 0 skipped | Exit 0 after `--force`; 34 suites, 305 passed, 0 failed, 0 skipped |
| Rust Solidity types | Exit 0 with all features | Exit 0 with all features |
| Shared Go bindings | Exit 0 | Exit 0 |
| Relayer Go | Exit 0 | Exit 0 |
| E2E Go packages | Exit 0, compile-only | Exit 0, compile-only |
| Compatibility self-test | Exit 0; all intentional drift rejected | Exit 0; ABI, storage, bytecode, gas, fixtures, bindings, tests, verifier pair, architecture, and protected docs pass |

Both Foundry passes used the configured 100,000 fuzz runs. The second pass rebuilt 201 files from forced artifacts. The final production-only size build rebuilt 162 files and reproduced the baseline sizes. Seventeen locked `RecvPacketGasTest` and `UpdateClientGasTest` hot-path measurements matched origin/main exactly, so the maximum recorded regression is 0%.

## Validation Boundaries

Format, shell syntax, JSON, stale-path, and whitespace gates pass. Docker/Kurtosis chains and shadowfork RPC tests were not run. Groth16 setup is randomized, so the selected paired N4 artifacts were provenance-checked twice rather than independently regenerated and incorrectly expected to have deterministic hashes.

`git diff --name-only -- docs/refactor` is empty, all four protected-source SHA-256 digests pass, and their modification times remain 2026-08-20: source documents were read only.
