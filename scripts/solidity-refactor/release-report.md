# Solidity Refactor Release Report

## Outcome

The handwritten Solidity source tree now matches `refactor/design-contracts@115397d`: the ICS-02 registry is under `core/client-registry`, light-client implementations remain under `light-clients`, and the only shared access libraries are `IBCRolesLib` and `IBCIdentifiers`. Deployment-only selector grouping lives under `scripts/deployments`; `contracts/shared/bytes`, `contracts/shared/encoding`, and `contracts/periphery/access` are absent. Tests mirror the final production packages.

Message and error files/containers use the reviewed `<Protocol>Msgs` and `<Protocol>Errors` naming. The final semantic diff review found no protocol behavior change. Runtime bytecode, canonical ABI shape/selectors/topics, storage, protocol fixtures, and locked gas measurements demonstrate equivalence. Generated verifier implementations and historical Groth16 artifacts were not hand-edited, and `docs/refactor` remained read-only.

## Compatibility

An isolated build of origin/main and the refactored tree matched exactly for canonical runtime ABI, creation bytecode, and deployed bytecode across 28 entrypoints, modules, stores, and libraries; nine first-party ERC-7201 namespace constants and ordered struct layouts also match exactly. Runtime comparison normalizes only the 16 reviewed container names listed in `tooling-rename-manifest.json`; tuple order/types, selectors, event indexing/topics, errors, bytecode, and storage remain locked. The manifest also pins the 14 regenerated ABI/Go-binding outputs by exact digest, and all Rust, Go, E2E compile-only, and relayer consumers pass with the new type names.

The compatibility negative suite rejects ABI, storage, runtime, gas, fixture/binding, test-inventory, fully qualified artifact, verifier-provenance, tooling-rename digest, and architecture drift. Full metadata changes in `abi/bytecode/SpectreClient.json` are accepted only where the semantic ABI, initcode, and runtime objects remain exact after the reviewed name normalization.

## Deployment and Prover

Production deploy/verify scripts require and register only `VERIFIER_N4`, matching `prover.Buckets` and the verifier manifest. Unsupported local N8/N16 verifier sources and N8/N16/N32/N64 key sets were moved into ignored, recoverable `.artifacts/solidity-refactor/prover/unsupported-archive/`.

Groth16 setup is randomized. `build-prover-artifacts.sh` stages, smoke-tests, publishes, and records one paired set. Gnark still emits its existing hash-to-field exporter warning; this release does not claim a live on-chain proof E2E.

## Size and Optimization Decision

Fresh sizes are unchanged: ICS26Router 24,465 B (111 B margin), ICS20Transfer 22,224 B (2,352 B), SpectreClient 19,390 B (5,186 B). Reaching 1 KiB router margin requires 913 safe bytes. No reviewed candidate met that threshold, so no optional bytecode optimization was attempted.

## Repeated Release Matrix

| Gate | Pass 1 | Pass 2 |
|---|---|---|
| Binding regeneration | Exit 0; approved ABI/binding digests reproduced | Exit 0; approved ABI/binding digests reproduced |
| Full non-shadowfork Foundry | Exit 0 after `--force`; 34 suites, 305 passed, 0 failed, 0 skipped | Exit 0 after `--force`; 34 suites, 305 passed, 0 failed, 0 skipped |
| Rust Solidity types | Exit 0 for locked package | Exit 0 for locked package |
| Shared Go bindings | Exit 0 | Exit 0 |
| Relayer Go | Exit 0 | Exit 0 |
| E2E Go packages | Exit 0, compile-only | Exit 0, compile-only |
| Compatibility self-test | Exit 0; all intentional drift including tooling rename rejected | Exit 0; all intentional drift rejected and positive compatibility/architecture/protected-doc gates pass |

Both Foundry passes used the configured 100,000 fuzz runs and independently rebuilt 200 files from forced artifacts. Both production-only size builds rebuilt 161 files and reproduced the baseline sizes. Seventeen locked `RecvPacketGasTest` and `UpdateClientGasTest` hot-path measurements matched origin/main exactly, so the maximum recorded regression is 0%.

## Validation Boundaries

Format, shell syntax, JSON, stale-path, and whitespace gates pass. Docker/Kurtosis chains and shadowfork RPC tests were not run. Groth16 setup is randomized, so the selected paired N4 artifacts were provenance-checked twice rather than independently regenerated and incorrectly expected to have deterministic hashes.

`git diff --name-only -- docs/refactor` is empty and the frozen compatibility baseline is unchanged: the reviewed source documents were read only.
