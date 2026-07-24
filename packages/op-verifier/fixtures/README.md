# OP optimistic fixtures

A complete fixture captures a fixed finalized Ethereum execution block, the factory and selected
game-list `eth_getProof` responses, game `eth_getCode`, the canonical committed OP block, its
output-root preimage, and the L2 router account proof. The selected game must be unresolved at the
captured Ethereum block to exercise the optimistic acceptance path.

Every fixture has a schema-v2 provenance manifest as described in `docs/L2_CLIENTS.md`.
