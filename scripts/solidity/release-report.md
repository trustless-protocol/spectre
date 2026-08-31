# Solidity Refactor Release Report

## Outcome

The handwritten Solidity source tree now matches `refactor/design-contracts@115397d`: the ICS-02 registry is under `core/client-registry`, light-client implementations remain under `light-clients`, and the only shared access libraries are `IBCRolesLib` and `IBCIdentifiers`. The core-owned `IPausable` compatibility interface is under `core/interfaces`, not a non-target `shared/interfaces` leaf. Deployment-only selector grouping lives under `scripts/deployments`; `contracts/shared/bytes`, `contracts/shared/encoding`, and `contracts/periphery/access` are absent. Tests mirror the final production packages.

Message and error interfaces move to their protocol owners without changing their existing `I...Msgs` and `I...Errors` public names. This keeps tooling ABI names and generated Go/Rust binding types stable and avoids consumer-only rename churn. The Rust Solidity-types package also drops four stale `sol!` inputs and seven conversion implementations whose Solidity interfaces were removed in #223; no in-repository consumer references those removed types, and the locked package now compiles. The stale `test-benchmark` target is repaired and covers both current `*GasTest.t.sol` suites. The final semantic diff review found no protocol behavior change. Runtime bytecode, ABI shape/selectors/topics, storage, protocol fixtures, and locked gas measurements demonstrate equivalence. Generated verifier implementations were not hand-edited.

## Compatibility

An isolated build of origin/main and the refactored tree matched exactly for ABI, creation bytecode, and deployed bytecode across 28 entrypoints, modules, stores, and libraries; nine first-party ERC-7201 namespace constants and ordered struct layouts also match exactly. No container-name normalization is required. Existing ABI and Go-binding files remain byte-for-byte compatible, while the manifest locks only the required Rust source-path and E2E fixture-path updates by exact digest.

The compatibility negative suite rejects ABI, storage, runtime, gas, fixture/binding, test-inventory, fully qualified artifact, deterministic R1CS, tooling-rename digest, and architecture drift. Full metadata changes in `abi/bytecode/SpectreClient.json` are accepted only where the semantic ABI, initcode, and runtime objects remain exact after the reviewed name normalization.

## Deployment and Prover

Production deploy/verify compatibility preserves the six-bucket topology N ∈ {4, 8, 16, 32, 64, 128}: all six verifier addresses must contain code, be distinct, and be registered with the expected selector. The checked local generator/prover manifest remains N4-only, so larger production buckets require coordinated prover enablement and paired artifact publication before launch; the structural refactor does not silently narrow the production configuration.

Groth16 setup is randomized. `build-prover-artifacts.sh` skips a complete configured artifact set, while `--force` stages, smoke-tests, publishes, and records every configured bucket without moving or deleting disabled larger local verifier/key artifacts. The compatibility gate locks only deterministic R1CS output and requires the randomized artifacts to exist; it does not treat self-written runtime provenance as baseline evidence. Gnark still emits its existing hash-to-field exporter warning; this release does not claim a live on-chain proof E2E.

## Size and Optimization Decision

Fresh sizes are unchanged: ICS26Router 24,465 B (111 B margin), ICS20Transfer 22,224 B (2,352 B), SpectreClient 19,390 B (5,186 B). Reaching 1 KiB router margin requires 913 safe bytes. No reviewed candidate met that threshold, so no optional bytecode optimization was attempted.

## Validation Matrix

The branch recorded two release-matrix passes before this review. After restoring the production bucket boundary and strengthening the gates, the final non-shadowfork Foundry suite was also run twice from clean builds; affected and cross-language checks were rerun as follows:

| Gate | Post-review result |
|---|---|
| Binding regeneration | Exit 0 with abigen 1.17.2-stable; ABI and Go consumers reproduced without generated type-name changes, and Spectre source metadata was refreshed while ABI/initcode/runtime remained exact |
| Production size build | Exit 0; ICS26Router 24,465 B, ICS20Transfer 22,224 B, SpectreClient 19,390 B |
| Full non-shadowfork Foundry | Two clean passes; each exited 0 with 34 suites, 305 passed, 0 failed, 0 skipped under the configured 100,000 fuzz runs |
| Production deployment focus | Exit 0; 8 tests, including six-bucket registration/verification, missing-code rejection, and duplicate-verifier rejection |
| Rust Solidity types | Exit 0 for the locked package |
| Shared Go bindings | Exit 0 |
| Relayer Go | Exit 0 |
| E2E Go packages | Exit 0, compile-only |
| Compatibility self-test | Exit 0; all intentional drift, including deterministic-R1CS drift, rejected; positive compatibility/architecture/protected-doc gates pass |

Seventeen locked `RecvPacketGasTest` and `UpdateClientGasTest` hot-path measurements still match origin/main, so the maximum recorded regression is 0%. An independent checkout/build of baseline commit `b75a613` also reproduced the recorded initcode/runtime hashes for the behavior-critical deployables.

## Validation Boundaries

Format, shell syntax, JSON, stale-path, toolchain-lock, and whitespace gates pass. The handwritten format script intentionally excludes generated verifier implementations. Docker/Kurtosis chains and shadowfork RPC tests were not run. Groth16 setup is randomized, so PK/VK/verifier identity is not compared with historical hashes; the staged generator smoke-tests each new pair before publishing it.

The protected-doc gate checks working-tree, index, and committed `baseline..HEAD` changes. The authoritative `01-Solidity-Contracts.md` intentionally remains on `origin/refactor/design-contracts@115397d`; the temporary in-branch copy and its revert do not remove the design authority and are omitted from the squash result to avoid a stale duplicate. Four local-only stale derivative design files, which were never present on the design branch or in committed history, were removed from the worktree and from the invalid evidence-digest inventory; runtime, storage, fixture, gas, and generated-binding baseline values remain unchanged.

Two requirements are historical/external evidence boundaries rather than source-tree mismatches. This repository records the prover generator command, deterministic N4 R1CS digest, and staged smoke-test boundary, but it does not contain the T-2 timing or tracking-issue evidence for the Hải readiness checkpoint. Published commit history also does not preserve the document's production-move-before-test-cleanup sequencing; correcting that record would require rewriting the already-published branch. Neither boundary is claimed as resolved by the current-tree validation.
