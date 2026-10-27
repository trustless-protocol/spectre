package prover

import (
	"relayer/prover/witness"

	"crypto/ed25519"
	"testing"
)

// MakeTestVote builds a synthetic canonical-vote byte slice satisfying the
// #199 layout: a 1-byte leading length varint, the Type|Height prefix head at
// body offset 0, no round field, and the 32-byte block hash at the no-round
// offset. The exact field values are irrelevant — witness.CommonPrefixFields and the
// witness hash only read the structure.
func MakeTestVote() []byte {
	const voteLen = 120
	msg := make([]byte, voteLen)
	msg[0] = byte(voteLen - 1) // 1-byte varint, continuation bit clear
	msg[1] = 0x08              // Type field
	msg[2] = 0x02              // precommit
	msg[3] = 0x11              // Height sfixed64
	for j := 0; j < 8; j++ {
		msg[4+j] = byte(j + 1)
	}
	msg[12] = 0x22 // BlockID tag (not 0x19 → no round, Round = 0)
	msg[13] = 0x48
	msg[14] = 0x0a
	msg[15] = 0x20
	for j := 0; j < 32; j++ {
		msg[16+j] = byte(0xA0 + j)
	}
	for j := 48; j < voteLen; j++ {
		msg[j] = byte(j)
	}
	return msg
}

// MakeTestSigs returns n real, verifying Ed25519 signatures over the same
// synthetic vote — the shape SubmitProof validation expects.
func MakeTestSigs(t *testing.T, n int) []witness.ValidatorSignature {
	t.Helper()
	msg := MakeTestVote()
	out := make([]witness.ValidatorSignature, n)
	for i := range out {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			t.Fatalf("generate key: %v", err)
		}
		out[i] = witness.ValidatorSignature{
			Signature:   ed25519.Sign(priv, msg),
			PublicKey:   pub,
			SignedBytes: msg,
			Index:       i,
			Power:       10,
			Active:      true,
		}
	}
	return out
}
