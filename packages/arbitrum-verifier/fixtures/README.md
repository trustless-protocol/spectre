# Arbitrum optimistic fixtures

A complete fixture captures a fixed finalized Ethereum execution block, the Rollup assertion
`eth_getProof`, the canonical committed Arbitrum block, and the L2 router account proof.

- A BoLD fixture includes the full `AssertionCreated` data used to recompute its hash.

Every fixture has a schema-v2 provenance manifest as described in `docs/L2_CLIENTS.md`.
