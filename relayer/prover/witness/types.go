// Package proof holds the proof data model and the witness byte contract the
// prover core, the gnark engine, and the queue/store adapters all share. It is
// a leaf package (stdlib only) so every prover subpackage can import it
// without cycles; the root prover package aliases everything here, so callers
// keep using prover.X names.
package witness

// ValidatorSignature is one slot of a proof batch: the Ed25519 signature a
// validator produced over its canonical vote bytes, plus the metadata the
// on-chain quorum check needs. It is the unit of data crossing the boundary
// between the relayer (which selects signatures) and the prover.
type ValidatorSignature struct {
	Signature   []byte // 64 bytes: R || S
	PublicKey   []byte // 32 bytes: compressed Ed25519 A
	SignedBytes []byte // Commit.VoteSignBytes(chainID, idx) — what was actually signed
	Index       int    // position in the validator set
	Power       int64  // voting power, consumed by the on-chain quorum check
	// Active distinguishes real signers (true) from deterministic dummy
	// padding (false). Every slot's signature verifies in-circuit; Active is
	// committed into the witness hash and is what the on-chain quorum check
	// reads to decide which slots carry voting power.
	Active bool
}
