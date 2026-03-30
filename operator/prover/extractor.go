package prover

import (
	"crypto/ed25519"
	"fmt"

	operatorclient "operator/client"

	"github.com/cometbft/cometbft/types"
)

// ValidatorSignature contains the extracted Ed25519 signature data from a LightBlock.
type ValidatorSignature struct {
	Signature []byte // 64 bytes: R || S
	PublicKey []byte // 32 bytes: compressed Ed25519 public key
	SignBytes []byte // canonical vote sign bytes
}

// ExtractValidatorSignature extracts the first non-absent commit signature from a LightBlock.
// For single validator mode, this returns the first valid signature found.
func ExtractValidatorSignature(lightBlock *operatorclient.LightBlock, chainID string) (*ValidatorSignature, error) {
	if lightBlock == nil {
		return nil, fmt.Errorf("light block is nil")
	}

	commit := lightBlock.SignedHeader.Commit
	if commit == nil {
		return nil, fmt.Errorf("commit is nil")
	}

	validators := lightBlock.ValSet
	if len(validators.Validators) == 0 {
		return nil, fmt.Errorf("validator set is empty")
	}

	for i, sig := range commit.Signatures {
		if sig.BlockIDFlag == types.BlockIDFlagAbsent {
			continue
		}

		if i >= len(validators.Validators) {
			return nil, fmt.Errorf("validator index %d out of range (have %d validators)", i, len(validators.Validators))
		}

		validator := validators.Validators[i]
		pubKeyBytes := validator.PubKey.Bytes()
		if len(pubKeyBytes) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("invalid public key size: %d, expected %d", len(pubKeyBytes), ed25519.PublicKeySize)
		}

		sigData := sig.Signature
		if len(sigData) != ed25519.SignatureSize {
			continue // skip invalid signatures
		}

		// Get canonical vote sign bytes
		voteData := commit.VoteSignBytes(chainID, int32(i))

		// Verify the signature before using it
		if !ed25519.Verify(pubKeyBytes, voteData, sigData) {
			return nil, fmt.Errorf("signature verification failed for validator %d", i)
		}

		return &ValidatorSignature{
			Signature: sigData,
			PublicKey: pubKeyBytes,
			SignBytes: voteData,
		}, nil
	}

	return nil, fmt.Errorf("no valid non-absent signatures found in commit")
}
