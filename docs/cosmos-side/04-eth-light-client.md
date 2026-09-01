# Ethereum Light Client (`cw-ics08-wasm-eth`)

The CosmWasm contract inside 08-wasm that verifies Ethereum. Logic in
`packages/ethereum/light-client`; types in `packages/ethereum/types`.

## Consensus verification
> TODO: sync-committee light-client protocol (Altair spec) + the historical-updates
> modification (why: multiple independent relayers). Trust assumptions stated plainly.

## State verification
> TODO: each update carries an account proof for `IBCStore.sol` → client tracks its
> storage root; membership/non-membership via MPT proofs against that root.

## Misbehaviour
> TODO: two conflicting valid updates at one height ⇒ freeze; recovery via governance.

## Client & consensus state
> TODO: what each stores (fork schedule, sync committee, storage root, timestamps) —
> cite the type definitions.
