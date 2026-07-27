# Arbitrum optimistic L2 client

This artifact supports two profile-selected RollupCore layouts:

- BoLD v2 assertions authenticated through the `_assertions` mapping; and
- legacy Nitro numeric nodes authenticated through the packed lifecycle slot and
  `_nodes[nodeNumber].confirmData`.

It accepts pending or confirmed commitments once their RollupCore state exists in a finalized
state of the pinned Ethereum light client. It does not wait for rollup challenge settlement, so an
authenticated pending proposal may later be rejected by the rollup.

Safe operation therefore requires an independent watchdog to submit conflicting authenticated
headers through the misbehaviour path before affected packet proofs are used.

Protocol, RollupCore storage layout, router, commitment slot, chain, version, and Ethereum-client
identities come from the immutable runtime profile stored in client state. A header's tagged
protocol must match that profile; no Sepolia profile is compiled into the Wasm.
