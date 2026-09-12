# Authenticated L2 legacy-test mapping

This inventory corrects the review count: the removed list contains **9 tests** and **2 helper
functions**, not 11 tests. The authenticated replacements below retain or strengthen each behavior while
removing the obsolete on-chain finality classifier.

| Legacy symbol | Kind | Authenticated replacement |
|---|---|---|
| `signing_bytes_are_domain_separated_and_big_endian` | test | `state_tests::v1_domain_and_statement_offsets_are_wire_locked` pins the V1 development domain, 164-byte length, every offset and both big-endian integers. |
| `accepts_the_go_attestor_signature_vector` | test | `state_tests::decodes_every_go_vector_byte_for_byte` plus Go `TestWasmAttestationVectors` verify the same regenerated statements and Ed25519 signatures. |
| `rejects_an_unsigned_header_before_any_router_proof_is_accepted` | test | `verification_tests::rejects_certificate_structure_before_any_crypto_or_trie_work` and `entrypoints_tests::host_paths_reject_conflicts_without_authenticated_attestations` reject unsigned input before malformed proofs and preserve storage. |
| `host_paths_freeze_on_conflicting_finalized_attestations` | test | wasmvm `TestOptimizedArtifactsConformToPinnedWasmVM` submits two independently signed conflicting certificates through the real sudo path, asserts Frozen, and asserts the trusted consensus bytes are unchanged. |
| `host_paths_reject_reorgable_conflicts_without_freezing` | test | Finality classification is retired on-chain. `entrypoints_tests::host_paths_reject_conflicts_without_authenticated_attestations` proves unauthenticated evidence cannot freeze; config `TestLoadConfig_RejectsAttestorSetReuseAcrossFinalityTiers` enforces the replacement trust boundary. |
| `same_height_finalized_conflict_freezes_without_overwriting` | test | `runtime_tests::only_explicit_authenticated_misbehaviour_freezes` plus the wasmvm conformance assertion cover freeze and no consensus overwrite. |
| `reorgable_same_height_conflicts_are_rejected_without_freezing` | test | Finality classification is retired on-chain. `runtime_tests::same_height_conflict_is_rejected_without_freezing_or_overwriting` keeps direct-update behavior; config rejects cross-tier key reuse. |
| `finalized_conflict_with_a_bad_parent_is_rejected_before_the_parent_check` | test | `runtime_tests::same_height_conflict_with_a_bad_parent_is_rejected_before_the_parent_check` retains precedence and no overwrite. |
| `only_finalized_conflicts_prove_misbehaviour` | test | `verification_tests::misbehaviour_authenticates_both_certificates_before_either_proof`, real wasmvm signed-conflict coverage, and unsigned/invalid-certificate zero-write cases establish that only two active-set certificates are actionable. |
| `profile_with_attestation_head` | helper | Removed: profiles contain no finality field; finality is an off-chain relayer/attestor policy. |
| `bootstrap_with_attestation_head` | helper | Removed: bootstrap stores an attestor set, while finality-tier isolation is enforced at relayer config load. |

The helpers are deliberately not counted as tests. No unsigned profile/state migration is implied
by this mapping: missing-attestors state fails with a typed error before any write, and a fresh
authenticated client must be created. Finality-specific deployments must use disjoint attestor
keys; there is no statement byte that silently substitutes for the removed helper.
