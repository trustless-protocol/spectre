package prover

import (
	"crypto/ed25519"
	"testing"
	"time"

	operatorclient "operator/client"

	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/types"
)

func TestExtractValidatorSignature_NilLightBlock(t *testing.T) {
	_, err := ExtractValidatorSignature(nil, "test-chain")
	if err == nil {
		t.Fatal("expected error for nil light block, got nil")
	}
	if err.Error() != "light block is nil" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestExtractValidatorSignature_NilCommit(t *testing.T) {
	lb := &operatorclient.LightBlock{
		SignedHeader: types.SignedHeader{
			Commit: nil,
		},
	}

	_, err := ExtractValidatorSignature(lb, "test-chain")
	if err == nil {
		t.Fatal("expected error for nil commit, got nil")
	}
	if err.Error() != "commit is nil" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestExtractValidatorSignature_EmptyValidators(t *testing.T) {
	lb := &operatorclient.LightBlock{
		SignedHeader: types.SignedHeader{
			Commit: &types.Commit{
				Signatures: []types.CommitSig{},
			},
		},
		ValSet: types.ValidatorSet{
			Validators: []*types.Validator{},
		},
	}

	_, err := ExtractValidatorSignature(lb, "test-chain")
	if err == nil {
		t.Fatal("expected error for empty validators, got nil")
	}
	if err.Error() != "validator set is empty" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestExtractValidatorSignature_AllAbsent(t *testing.T) {
	privKey := cmted25519.GenPrivKey()
	pubKey := privKey.PubKey()

	val := &types.Validator{
		Address:     pubKey.Address(),
		PubKey:      pubKey,
		VotingPower: 10,
	}

	lb := &operatorclient.LightBlock{
		SignedHeader: types.SignedHeader{
			Commit: &types.Commit{
				Signatures: []types.CommitSig{
					{
						BlockIDFlag: types.BlockIDFlagAbsent,
					},
				},
			},
		},
		ValSet: types.ValidatorSet{
			Validators: []*types.Validator{val},
		},
	}

	_, err := ExtractValidatorSignature(lb, "test-chain")
	if err == nil {
		t.Fatal("expected error when all signatures are absent, got nil")
	}
	expected := "no valid non-absent signatures found in commit"
	if err.Error() != expected {
		t.Fatalf("unexpected error message: got %q, want %q", err.Error(), expected)
	}
}

func TestExtractValidatorSignature_ValidSingleValidator(t *testing.T) {
	chainID := "test-chain-1"

	privKey := cmted25519.GenPrivKey()
	pubKey := privKey.PubKey()

	val := &types.Validator{
		Address:     pubKey.Address(),
		PubKey:      pubKey,
		VotingPower: 10,
	}

	blockID := types.BlockID{
		Hash: make([]byte, 32),
		PartSetHeader: types.PartSetHeader{
			Total: 1,
			Hash:  make([]byte, 32),
		},
	}

	header := &types.Header{
		ChainID: chainID,
		Height:  100,
		Time:    time.Now().UTC(),
	}

	commit := &types.Commit{
		Height:  100,
		Round:   0,
		BlockID: blockID,
		Signatures: []types.CommitSig{
			{
				BlockIDFlag:      types.BlockIDFlagCommit,
				ValidatorAddress: val.Address,
				Timestamp:        header.Time,
			},
		},
	}

	// Compute the vote sign bytes so we can sign them
	voteSignBytes := commit.VoteSignBytes(chainID, 0)

	// Sign the vote bytes with the private key
	sig, err := privKey.Sign(voteSignBytes)
	if err != nil {
		t.Fatalf("failed to sign vote bytes: %v", err)
	}
	commit.Signatures[0].Signature = sig

	lb := &operatorclient.LightBlock{
		SignedHeader: types.SignedHeader{
			Header: header,
			Commit: commit,
		},
		ValSet: types.ValidatorSet{
			Validators: []*types.Validator{val},
		},
		BlockHeight: 100,
	}

	result, err := ExtractValidatorSignature(lb, chainID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify signature length
	if len(result.Signature) != ed25519.SignatureSize {
		t.Fatalf("expected signature length %d, got %d", ed25519.SignatureSize, len(result.Signature))
	}

	// Verify public key length
	if len(result.PublicKey) != ed25519.PublicKeySize {
		t.Fatalf("expected public key length %d, got %d", ed25519.PublicKeySize, len(result.PublicKey))
	}

	// Verify the returned signature matches what we signed
	if !ed25519.Verify(result.PublicKey, result.SignBytes, result.Signature) {
		t.Fatal("returned signature does not verify against the returned public key and sign bytes")
	}

	// Verify sign bytes match what we expect
	expectedSignBytes := commit.VoteSignBytes(chainID, 0)
	if len(result.SignBytes) != len(expectedSignBytes) {
		t.Fatalf("sign bytes length mismatch: got %d, want %d", len(result.SignBytes), len(expectedSignBytes))
	}
	for i := range result.SignBytes {
		if result.SignBytes[i] != expectedSignBytes[i] {
			t.Fatalf("sign bytes differ at index %d", i)
		}
	}
}
