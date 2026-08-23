package prover

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"
)

// ZK-07: the witness commitment is computed in three places — here in Go
// (ComputeWitnessHash), in-circuit (BatchCircuit.Define), and on-chain
// (SignatureVerifier._hashWitness) — and a Groth16 proof is only accepted when
// all three agree. Until now only two of those pairings were tested:
// Solidity-vs-Solidity by SignatureVerifierHashEquivalenceTest, and
// Go-vs-circuit by the solver tests. Nothing compared the Go layout against the
// Solidity one, so a drift between them would first show up in E2E, if at all.
//
// This test closes that by pinning vectors both sides can evaluate: Go writes
// the inputs and the digest it computes, and the Foundry test
// (test/light-clients/spectre/SignatureVerifierHashEquivalenceTest.t.sol) reads the same
// file and asserts _hashWitness reproduces the digest.
//
// Run with -update to regenerate after an intentional layout change. Keeping the
// generator in-repo and wired to a test is deliberate — see SUP-01/SUP-06, where
// golden vectors outlived the generator that produced them and became
// unmaintainable.

var updateWitnessVectors = flag.Bool("update", false, "regenerate the Go↔Solidity witness-hash vectors")

// witnessHashVectorsPath is relative to this package directory.
const witnessHashVectorsPath = "../../test/fixtures/witness_hash_vectors.json"

// witnessHashVector is one cross-check case. The fields are exactly what the
// on-chain verifier receives: SharedBlock (height, round, blockIDHash) plus the
// per-slot pubkeys and active flags. Go derives the same shared fields from the
// signed bytes instead, which is the property under test.
type witnessHashVector struct {
	Name        string   `json:"name"`
	Height      uint64   `json:"height"`
	Round       uint32   `json:"round"`
	BlockIDHash string   `json:"blockIdHash"`
	Pubkeys     []string `json:"pubkeys"`
	Active      []bool   `json:"active"`
	WitnessHash string   `json:"witnessHash"`
}

// witnessHashVectorFile wraps the vectors in an object carrying an explicit
// count. Foundry's JSON cheatcodes address one value at a time and have no
// wildcard, so the Solidity side needs the length before it can walk the array.
type witnessHashVectorFile struct {
	Count   int                 `json:"count"`
	Vectors []witnessHashVector `json:"vectors"`
}

// crossCheckVote builds canonical-vote bytes for an explicit (height, round,
// blockHash) in the #199 body layout, so the Solidity side can rebuild the same
// PrefixHead from SharedBlock alone.
//
//	uvarint(bodyLen) || Type(0x08 0x02) || Height(0x11 || sfixed64) ||
//	  [Round(0x19 || sfixed64)] || BlockID(0x22 0x48 0x0a 0x20 || hash[32]) || trailer
//
// sfixed64 is little-endian, matching Encode.encodeSfixed64. trailLen stands in
// for the timestamp and chain id and drives bodyLen across the 128-byte boundary
// that decides whether the leading varint is one byte or two.
func crossCheckVote(height uint64, round uint32, blockHash [32]byte, trailLen int) []byte {
	body := []byte{0x08, 0x02, 0x11}
	var h [8]byte
	binary.LittleEndian.PutUint64(h[:], height)
	body = append(body, h[:]...)

	if round > 0 {
		var r [8]byte
		binary.LittleEndian.PutUint64(r[:], uint64(round))
		body = append(body, 0x19)
		body = append(body, r[:]...)
	}

	body = append(body, 0x22, 0x48, 0x0a, 0x20)
	body = append(body, blockHash[:]...)
	body = append(body, bytes.Repeat([]byte{0x2a}, trailLen)...)

	prefix := make([]byte, binary.MaxVarintLen64)
	n := binary.PutUvarint(prefix, uint64(len(body)))
	return append(prefix[:n], body...)
}

func buildWitnessHashVectors(t *testing.T) []witnessHashVector {
	t.Helper()

	cases := []struct {
		name     string
		height   uint64
		round    uint32
		bucket   int
		signers  int
		trailLen int
	}{
		// Single-byte height, no round, some padding slots — the common shape.
		{"bucket16_noRound_mixedSlots", 53, 0, 16, 11, 10},
		// Height wide enough to exercise every sfixed64 byte, round present, and a
		// trailer long enough to force the two-byte leading varint.
		{"bucket8_round_wideHeight_twoByteVarint", 0x0102030405060708, 7, 8, 6, 90},
		// Every slot active.
		{"bucket4_noRound_allActive", 1, 0, 4, 4, 10},
		// Every slot padding: no active signer, so there is no shared prefix to
		// derive — see the assertion below.
		{"bucket4_round_singleSigner", 999_999, 1, 4, 1, 10},
	}

	vectors := make([]witnessHashVector, 0, len(cases))
	for _, tc := range cases {
		var blockHash [32]byte
		for i := range blockHash {
			blockHash[i] = byte(0xA0 + i + tc.bucket)
		}

		vote := crossCheckVote(tc.height, tc.round, blockHash, tc.trailLen)

		sigs := make([]ValidatorSignature, tc.bucket)
		pubkeys := make([]string, tc.bucket)
		active := make([]bool, tc.bucket)
		for i := range sigs {
			pub := make([]byte, 32)
			for j := range pub {
				pub[j] = byte(i*7 + j + 1)
			}
			sigs[i] = ValidatorSignature{PublicKey: pub, Index: i, Active: i < tc.signers}
			if sigs[i].Active {
				sigs[i].SignedBytes = vote
			} else {
				sigs[i].SignedBytes = DummyMsgBytes(uint16(tc.bucket), uint16(i))
			}
			pubkeys[i] = "0x" + hex.EncodeToString(pub)
			active[i] = sigs[i].Active
		}

		hash, err := ComputeWitnessHash(sigs)
		if err != nil {
			t.Fatalf("%s: ComputeWitnessHash: %v", tc.name, err)
		}

		vectors = append(vectors, witnessHashVector{
			Name:        tc.name,
			Height:      tc.height,
			Round:       tc.round,
			BlockIDHash: "0x" + hex.EncodeToString(blockHash[:]),
			Pubkeys:     pubkeys,
			Active:      active,
			WitnessHash: "0x" + hex.EncodeToString(hash[:]),
		})
	}

	return vectors
}

// TestWitnessHashVectors regenerates the vectors and compares them to the
// committed file. A failure here means the Go witness layout changed: rerun with
// -update, and expect the Foundry side to fail too until Solidity is changed to
// match (or, if only Go drifted, until Go is fixed back).
func TestWitnessHashVectors(t *testing.T) {
	vectors := buildWitnessHashVectors(t)

	got, err := json.MarshalIndent(witnessHashVectorFile{Count: len(vectors), Vectors: vectors}, "", "  ")
	if err != nil {
		t.Fatalf("marshal vectors: %v", err)
	}
	got = append(got, '\n')

	path := filepath.Clean(witnessHashVectorsPath)
	if *updateWitnessVectors {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("create fixture dir: %v", err)
		}
		if err := os.WriteFile(path, got, 0o644); err != nil {
			t.Fatalf("write vectors: %v", err)
		}
		t.Logf("regenerated %s", path)
		return
	}

	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read vectors (run with -update to generate): %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(want), bytes.TrimSpace(got)) {
		t.Fatalf("witness-hash vectors are stale.\nThe Go witness layout no longer matches %s.\n"+
			"If the change was intended, rerun with -update and update SignatureVerifier._hashWitness "+
			"and BatchCircuit.Define in the same change.", path)
	}
}

// TestComputeWitnessHash_RequiresActiveSigner pins the precondition the vectors
// above rely on: the shared prefix is read off a real signed message, so a batch
// with no active signer has nothing to commit to and must be an error rather
// than a zero-prefix digest.
func TestComputeWitnessHash_RequiresActiveSigner(t *testing.T) {
	sigs := make([]ValidatorSignature, 4)
	for i := range sigs {
		pub := make([]byte, 32)
		pub[0] = byte(i + 1)
		sigs[i] = ValidatorSignature{
			PublicKey:   pub,
			SignedBytes: DummyMsgBytes(4, uint16(i)),
			Active:      false,
		}
	}

	if _, err := ComputeWitnessHash(sigs); err == nil {
		t.Fatal("expected an error for a batch with no active signer")
	}
}
