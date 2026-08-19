package services

import (
	"math/big"
	"testing"

	updateclientContract "relayer/bindings/UpdateClient"
	"relayer/prover"
)

func TestUpdateCommitToMisbehaviourPreservesValidatorAddress(t *testing.T) {
	wantAddress := [20]byte{0x01, 0x23, 0x45, 0x67, 0x89}
	converted := updateCommitToMisbehaviour(updateclientContract.IICS07TendermintMsgsBlockCommit{
		CommitSigs: []updateclientContract.IICS07TendermintMsgsCommitSig{{
			Flag:             2,
			ValidatorAddress: wantAddress,
		}},
	})

	if len(converted.CommitSigs) != 1 {
		t.Fatalf("commit signatures = %d, want 1", len(converted.CommitSigs))
	}
	if got := converted.CommitSigs[0]; got.Flag != 2 || got.ValidatorAddress != wantAddress {
		t.Fatalf("converted signature = %+v, want flag=2 validatorAddress=%x", got, wantAddress)
	}
}

func TestMisbehaviourEvidenceDetectsPinnedSetEquivocation(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey0, pubkey1, pubkey2},
		votingPowers: []uint64{40, 35, 25},
		totalPower:   100,
	}

	header1 := misbehaviourHeaderCandidate{
		height:    42,
		blockHash: [32]byte{0xAA},
		candidates: []prover.ValidatorSignature{
			{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
			{Index: 1, PublicKey: pubkey1[:], Power: 35, Active: true},
		},
	}
	header2 := misbehaviourHeaderCandidate{
		height:    42,
		blockHash: [32]byte{0xBB},
		candidates: []prover.ValidatorSignature{
			{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
			{Index: 2, PublicKey: pubkey2[:], Power: 25, Active: true},
			{Index: 1, PublicKey: pubkey1[:], Power: 35, Active: true},
		},
	}

	evidence, err := misbehaviourEvidenceFromCandidates(header1, header2, pinned, pinned)
	if err != nil {
		t.Fatalf("build evidence: %v", err)
	}
	if evidence == nil {
		t.Fatal("expected same-height distinct block hashes with pinned quorum to produce evidence")
	}
	if len(evidence.sigs1) != 2 || len(evidence.sigs2) != 2 {
		t.Fatalf("selected signatures = %d/%d, want 2/2", len(evidence.sigs1), len(evidence.sigs2))
	}
	if _, err := misbehaviourEvidenceFromCandidates(header1, header1, pinned, pinned); err == nil {
		t.Fatal("same block hash at the same height must be rejected")
	}
}

func TestMisbehaviourEvidenceIgnoresConflictsWithoutPinnedQuorum(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pubkey2 := [32]byte{0x03}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1, 2},
		pubkeys:      [][32]byte{pubkey0, pubkey1, pubkey2},
		votingPowers: []uint64{40, 35, 25},
		totalPower:   100,
	}

	header1 := misbehaviourHeaderCandidate{
		height:    99,
		blockHash: [32]byte{0xAA},
		candidates: []prover.ValidatorSignature{
			{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
			{Index: 1, PublicKey: pubkey1[:], Power: 35, Active: true},
		},
	}
	header2 := misbehaviourHeaderCandidate{
		height:    99,
		blockHash: [32]byte{0xBB},
		candidates: []prover.ValidatorSignature{
			{Index: 0, PublicKey: pubkey0[:], Power: 40, Active: true},
		},
	}

	if _, err := misbehaviourEvidenceFromCandidates(header1, header2, pinned, pinned); err == nil {
		t.Fatal("conflicting headers without pinned quorum must be rejected")
	}
}

func TestMisbehaviourBatchProofBuildsPinnedMetadata(t *testing.T) {
	pubkey0 := [32]byte{0x01}
	pubkey1 := [32]byte{0x02}
	pinned := pinnedCosmosValidatorSet{
		indices:      []uint32{0, 1},
		pubkeys:      [][32]byte{pubkey0, pubkey1},
		votingPowers: []uint64{50, 50},
		totalPower:   100,
	}
	headerProof := prover.HeaderBatchProof{
		Bucket: 2,
		PaddedSigs: []prover.ValidatorSignature{
			{Index: 7, PublicKey: pubkey1[:], Active: true},
			{Index: 0, PublicKey: make([]byte, 32), Active: false},
		},
		Proof:         [8]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6), big.NewInt(7), big.NewInt(8)},
		Commitments:   [2]*big.Int{big.NewInt(9), big.NewInt(10)},
		CommitmentPok: [2]*big.Int{big.NewInt(11), big.NewInt(12)},
	}

	got, err := misbehaviourBatchProof(headerProof, pinned)
	if err != nil {
		t.Fatalf("misbehaviourBatchProof: %v", err)
	}
	if got.Bucket != 2 {
		t.Fatalf("bucket = %d, want 2", got.Bucket)
	}
	if len(got.SignerIndices) != 2 || got.SignerIndices[0] != 7 {
		t.Fatalf("signer indices = %+v", got.SignerIndices)
	}
	if len(got.PinnedValidatorIndices) != 2 || got.PinnedValidatorIndices[0] != 1 {
		t.Fatalf("pinned indices = %+v", got.PinnedValidatorIndices)
	}
	if len(got.Active) != 2 || !got.Active[0] || got.Active[1] {
		t.Fatalf("active flags = %+v", got.Active)
	}
}
