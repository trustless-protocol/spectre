# Arbitrum optimistic L2 client

This artifact accepts pending or confirmed BoLD v2 assertions once the assertion exists in a
finalized state of the pinned Ethereum light client. It does not wait for BoLD challenge
settlement, so an authenticated proposal may later be rejected by the rollup.

Safe operation therefore requires an independent watchdog to submit conflicting authenticated
headers through the misbehaviour path before affected packet proofs are used.

Rollup, assertion-layout, router, commitment-slot, chain, version, and Ethereum-client identities
come from the immutable runtime profile stored in client state; no Sepolia profile is compiled into
the Wasm.
