// Package arbitrum manages a pinned Nitro node and exposes independently
// verified Arbitrum assertion commitments through gRPC. Nitro owns persistent
// execution state; the attestor persists only its finalized-L1 assertion
// cursor, verified frontier, and mismatch metadata.
package arbitrum
