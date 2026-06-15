package prover

import (
	"crypto/ed25519"
	"fmt"
	"sort"

	relayerclient "relayer/client"

	"github.com/cometbft/cometbft/types"
)

// ValidatorSignature is the extracted Ed25519 signature data for one validator
// in a block commit, alongside the per-validator timestamp + the canonical
// vote bytes that validator actually signed (ready to feed into the circuit).
type ValidatorSignature struct {
	Signature        []byte // 64 bytes: R || S
	PublicKey        []byte // 32 bytes: compressed Ed25519 public key
	Index            int    // index in the block's validator set
	Power            int64  // validator voting power
	TimestampSeconds int64  // google.protobuf.Timestamp seconds
	TimestampNanos   int32  // google.protobuf.Timestamp nanos
	SignedBytes      []byte // cometbft.Commit.VoteSignBytes(chainID, idx)
	// Active distinguishes real signers (true) from deterministic padding
	// (false). Padding slots carry distinct dummy data so the in-circuit ECIP
	// divisor stays well-formed; their contribution is gated to zero.
	Active bool
}

// SharedBlockData is the subset of CanonicalVote fields that are identical
// across every validator signing the same block. The on-chain WrapperVerifier
// uses these fields to recompute each slot's canonical vote bytes via Encode.sol;
// the circuit itself doesn't see them — it consumes only the per-slot signed
// bytes via ValidatorSignature.SignedBytes.
type SharedBlockData struct {
	Height       int64
	Round        int64
	BlockIDHash  []byte // 32 bytes
	PartSetTotal uint32
	PartSetHash  []byte // 32 bytes
	ChainID      string
}

// ExtractorResult bundles the valid signer candidates, the selected signer
// prefix, and the shared block data needed by the on-chain quorum +
// canonical-vote rebuild path.
type ExtractorResult struct {
	Shared     SharedBlockData
	Candidates []ValidatorSignature
	Signatures []ValidatorSignature
}

// ExtractValidatorSignatures collects non-absent, locally-verified commit
// signatures until their cumulative voting power exceeds 2/3 of
// TotalVotingPower. Candidates are sorted by power for selection, while the
// selected proof slots are returned in validator-index order.
//
// allowedIndices restricts which validator slots may be selected. When nil,
// all signing validators are candidates. Pass the set of cached signer indices
// to guarantee the proof only references indices already cached on-chain.
//
// Returns an error if no quorum can be reached or the required signer count
// would exceed the largest configured bucket.
func ExtractValidatorSignatures(
	lightBlock *relayerclient.LightBlock,
	chainID string,
	allowedIndicesOpt ...map[uint32]bool,
) (*ExtractorResult, error) {
	var allowedIndices map[uint32]bool
	if len(allowedIndicesOpt) > 0 {
		allowedIndices = allowedIndicesOpt[0]
	}

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

	candidates := make([]ValidatorSignature, 0, len(commit.Signatures))
	for i, sig := range commit.Signatures {
		if sig.BlockIDFlag == types.BlockIDFlagAbsent {
			continue
		}
		if allowedIndices != nil && !allowedIndices[uint32(i)] {
			continue
		}
		if i >= len(validators.Validators) {
			return nil, fmt.Errorf("validator index %d out of range (have %d validators)", i, len(validators.Validators))
		}
		validator := validators.Validators[i]
		pubKeyBytes := validator.PubKey.Bytes()
		if len(pubKeyBytes) != ed25519.PublicKeySize {
			continue
		}
		if len(sig.Signature) != ed25519.SignatureSize {
			continue
		}
		voteData := commit.VoteSignBytes(chainID, int32(i))
		if !ed25519.Verify(pubKeyBytes, voteData, sig.Signature) {
			continue
		}
		candidates = append(candidates, ValidatorSignature{
			Signature:        sig.Signature,
			PublicKey:        pubKeyBytes,
			Index:            i,
			Power:            validator.VotingPower,
			TimestampSeconds: sig.Timestamp.Unix(),
			TimestampNanos:   int32(sig.Timestamp.Nanosecond()),
			SignedBytes:      voteData,
			Active:           true,
		})
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no valid non-absent signatures found in commit")
	}

	// Greedy by voting power so the smallest prefix hits quorum and the
	// smallest bucket can be used.
	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Power > candidates[j].Power
	})

	totalPower := validators.TotalVotingPower()
	quorum := totalPower*2/3 + 1

	var accumulated int64
	cutoff := len(candidates)
	for i, c := range candidates {
		accumulated += c.Power
		if accumulated >= quorum {
			cutoff = i + 1
			break
		}
	}
	if accumulated < quorum {
		return nil, fmt.Errorf("insufficient voting power: have %d, need %d of %d", accumulated, quorum, totalPower)
	}

	allCandidates := append([]ValidatorSignature(nil), candidates...)
	selected := append([]ValidatorSignature(nil), candidates[:cutoff]...)
	sortSelectedSignaturesByIndex(selected)
	if len(selected) > MaxBucket() {
		return nil, fmt.Errorf("quorum requires %d signers but largest bucket is %d", len(selected), MaxBucket())
	}

	return &ExtractorResult{
		Shared: SharedBlockData{
			Height:       commit.Height,
			Round:        int64(commit.Round),
			BlockIDHash:  commit.BlockID.Hash,
			PartSetTotal: commit.BlockID.PartSetHeader.Total,
			PartSetHash:  commit.BlockID.PartSetHeader.Hash,
			ChainID:      chainID,
		},
		Candidates: allCandidates,
		Signatures: selected,
	}, nil
}

func sortSelectedSignaturesByIndex(sigs []ValidatorSignature) {
	sort.Slice(sigs, func(i, j int) bool {
		return sigs[i].Index < sigs[j].Index
	})
}
