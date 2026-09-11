# OP attestor-trusted L2 client

**DEVNET ONLY.** This client accepts an L2 execution header from *any* submitter, not just the
relayer: `MsgUpdateClient` is permissionless and the header carries no attestor signature, so
nothing distinguishes one sender from another. Anyone can fabricate a self-consistent header at an
unseen height and then prove arbitrary membership or non-membership against it.

It verifies the header hash, the configured fork layout, and the router account proof — which
establishes only that the supplied storage root belongs to the supplied state, not that the state
is the L2's. It verifies no L1 consensus, no dispute game, and no attestor signature.

See [docs/L2_CLIENTS.md](../../docs/L2_CLIENTS.md) for the full statement and the
`ATTESTATIONS_ARE_AUTHENTICATED` gate that closes it.
