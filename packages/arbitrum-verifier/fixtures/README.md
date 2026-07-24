# Arbitrum optimistic fixtures

A complete fixture captures a fixed finalized Ethereum execution block, the Rollup assertion
`eth_getProof`, the full `AssertionCreated` event data used to recompute its hash, the canonical
committed Arbitrum block, and the L2 router account proof. The assertion must still be pending at
the captured Ethereum block to exercise the optimistic acceptance path.

Every fixture has a schema-v2 provenance manifest as described in `docs/L2_CLIENTS.md`.
