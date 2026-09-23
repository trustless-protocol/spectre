package gnark

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"os"
	"testing"

	"relayer/prover"
	"relayer/prover/witness"
)

// BenchmarkProve is the per-bucket proving baseline the performance
// requirement in 02-Prover is measured against (regressions over 10% need an
// explanation). It needs real artifacts:
//
//	PROVER_BIN_DIR=/abs/path/to/bin go test -bench BenchmarkProve -run '^$' ./prover/gnark/
//
// Skipped when PROVER_BIN_DIR is unset so the normal suite stays fast.
func BenchmarkProve(b *testing.B) {
	binDir := os.Getenv("PROVER_BIN_DIR")
	if binDir == "" {
		b.Skip("PROVER_BIN_DIR not set; benchmark needs real artifacts")
	}
	engine, err := New(binDir, nil)
	if err != nil {
		b.Fatalf("load engine: %v", err)
	}

	for _, n := range engine.Buckets() {
		req := benchRequest(b, engine, n)
		b.Run(fmt.Sprintf("n%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if _, err := engine.Prove(context.Background(), req); err != nil {
					b.Fatalf("prove n=%d: %v", n, err)
				}
			}
		})
	}
}

// benchRequest fills the whole bucket with real signers over a synthetic
// #199-layout vote — the worst case (no dummy slots).
func benchRequest(b *testing.B, engine *Engine, n int) prover.ProofRequest {
	b.Helper()
	msg := benchVote()
	sigs := make([]witness.ValidatorSignature, n)
	for i := range sigs {
		pub, priv, err := ed25519.GenerateKey(nil)
		if err != nil {
			b.Fatalf("generate key: %v", err)
		}
		sigs[i] = witness.ValidatorSignature{
			Signature:   ed25519.Sign(priv, msg),
			PublicKey:   pub,
			SignedBytes: msg,
			Index:       i,
			Active:      true,
		}
	}
	digest, err := witness.ComputeWitnessHash(sigs)
	if err != nil {
		b.Fatalf("witness hash: %v", err)
	}
	return prover.ProofRequest{
		ID:            prover.NewProofID(engine.VKDigest(n), digest),
		Bucket:        n,
		Sigs:          sigs,
		WitnessDigest: digest,
	}
}

// benchVote mirrors the setup smoke vote: 1-byte varint, Type|Height prefix
// head, no round, block hash at the no-round offset.
func benchVote() []byte {
	const voteLen = 120
	msg := make([]byte, voteLen)
	msg[0] = byte(voteLen - 1)
	msg[1], msg[2] = 0x08, 0x02
	msg[3] = 0x11
	for j := 0; j < 8; j++ {
		msg[4+j] = byte(j + 1)
	}
	msg[12], msg[13], msg[14], msg[15] = 0x22, 0x48, 0x0a, 0x20
	for j := 0; j < 32; j++ {
		msg[16+j] = byte(0xA0 + j)
	}
	for j := 48; j < voteLen; j++ {
		msg[j] = byte(j)
	}
	return msg
}
