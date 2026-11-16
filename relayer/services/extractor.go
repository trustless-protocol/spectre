package services

import (
	"crypto/ed25519"
	"fmt"
	"log"
	"sort"
	"sync"

	relayerclient "relayer/client"
	"relayer/prover/witness"

	"github.com/cometbft/cometbft/types"
)

// ExtractValidatorSignatures lives in services rather than the prover package
// because signature selection is relayer policy (greedy quorum, pinned-set
// allowedIndices) and it needs cometbft + the RPC client; the prover only
// consumes the resulting witness.ValidatorSignature slice.
//
// ExtractValidatorSignatures collects locally-verified COMMIT precommit
// signatures the proof circuit can prove (witness.CheckCircuitEncoding) until
// their cumulative voting power exceeds 2/3 of TotalVotingPower. Candidates are
// sorted by power for selection, while the selected proof slots are returned in
// validator-index order.
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
	maxBucket int,
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

	candidates := make([]witness.ValidatorSignature, 0, len(commit.Signatures))
	for i, sig := range commit.Signatures {
		// Only a COMMIT precommit is a proof signer. A NIL precommit verifies over
		// an empty BlockID, so it commits to no block: the witness hash, the
		// circuit and SpectreClient all reject it, and a top-power validator
		// voting nil would be selected first and fail every proof. Any other flag
		// makes VoteSignBytes panic, so this also runs before sign bytes are built.
		if sig.BlockIDFlag != types.BlockIDFlagCommit {
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
		// crypto/ed25519 accepts a few key and signature encodings the proof
		// circuit rejects (#579). Selecting one would fail every proof over this
		// quorum, and the relayer would select it again on every retry, so skip
		// the signer and build the quorum from the others.
		if err := witness.CheckCircuitEncoding(pubKeyBytes, sig.Signature); err != nil {
			reportUnprovableSigner(unprovableSigners, validator, i, commit.Height, err)
			continue
		}
		// A validator keyed with a public padding key still proves, so it stays
		// a candidate, but anyone can sign as it (#579 item 4).
		if bucket, slot, isDummy := witness.DummyKeySlot(pubKeyBytes); isDummy {
			reportDummyKeySigner(dummyKeySigners, validator, i, commit.Height, bucket, slot)
		}
		candidates = append(candidates, witness.ValidatorSignature{
			Signature:   sig.Signature,
			PublicKey:   pubKeyBytes,
			Index:       i,
			Power:       validator.VotingPower,
			SignedBytes: voteData,
			Active:      true,
		})
	}
	if len(candidates) == 0 {
		return nil, fmt.Errorf("no valid COMMIT precommit signatures found in commit")
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

	allCandidates := append([]witness.ValidatorSignature(nil), candidates...)
	selected := append([]witness.ValidatorSignature(nil), candidates[:cutoff]...)
	sortSelectedSignaturesByIndex(selected)
	if len(selected) > maxBucket {
		return nil, fmt.Errorf("quorum requires %d signers but largest available bucket is %d", len(selected), maxBucket)
	}

	return &ExtractorResult{
		Candidates: allCandidates,
		Signatures: selected,
	}, nil
}

func sortSelectedSignaturesByIndex(sigs []witness.ValidatorSignature) {
	sort.Slice(sigs, func(i, j int) bool {
		return sigs[i].Index < sigs[j].Index
	})
}

// ExtractorResult bundles the valid signer candidates and the selected signer
// set. Selection policy lives in the relayer; this type is the shape it hands
// back.
type ExtractorResult struct {
	Candidates []witness.ValidatorSignature
	Signatures []witness.ValidatorSignature
}

// unprovableSigners and dummyKeySigners hold the keys already reported by
// reportUnprovableSigner and reportDummyKeySigner. Every caller of
// ExtractValidatorSignatures shares them (the relay modules and the misbehaviour
// path run on their own goroutines), so reportedKeys guards each with its own
// mutex.
var (
	unprovableSigners = &reportedKeys{keys: make(map[string]struct{})}
	dummyKeySigners   = &reportedKeys{keys: make(map[string]struct{})}
)

// maxReportedKeys bounds reportedKeys. Reaching it clears the set, which costs
// one repeated WARN per key and nothing else.
const maxReportedKeys = 1024

// reportedKeys remembers which keys were reported, so that a validator that
// signs every block in a shape the circuit rejects is reported at WARN once and
// at DEBUG after that (logbridge: DEBUG is for repeats with nothing new to say).
type reportedKeys struct {
	mu   sync.Mutex
	keys map[string]struct{}
}

// firstReport records key and reports whether it was not recorded before.
func (r *reportedKeys) firstReport(key []byte) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, seen := r.keys[string(key)]; seen {
		return false
	}
	if len(r.keys) >= maxReportedKeys {
		clear(r.keys)
	}
	r.keys[string(key)] = struct{}{}
	return true
}

// reportUnprovableSigner logs a validator whose signature verifies but cannot be
// proven. An honestly generated key and nonce never take these shapes, so the
// line names the validator for the operator to follow up.
func reportUnprovableSigner(reported *reportedKeys, validator *types.Validator, index int, height int64, err error) {
	level := "DEBUG"
	if reported.firstReport(validator.PubKey.Bytes()) {
		level = "WARN"
	}
	log.Printf("%s [ExtractValidatorSignatures] skipping validator %s (index %d, power %d, pubkey %X) at height %d: "+
		"the proof circuit cannot prove its signature: %v. An honestly generated key and nonce never do this, so this "+
		"validator is most likely trying to stall the light client; the quorum is built from the other signers",
		level, validator.Address, index, validator.VotingPower, validator.PubKey.Bytes(), height, err)
}

// reportDummyKeySigner logs a validator whose consensus key is one of the public
// padding keys witness.PadWithDummies signs with. Its private key follows from a
// public seed, so this is never an accident; the line names the validator for
// the operator to follow up.
func reportDummyKeySigner(reported *reportedKeys, validator *types.Validator, index int, height int64, bucket, slot int) {
	level := "DEBUG"
	if reported.firstReport(validator.PubKey.Bytes()) {
		level = "WARN"
	}
	log.Printf("%s [ExtractValidatorSignatures] validator %s (index %d, power %d, pubkey %X) at height %d uses the "+
		"public padding key of bucket %d slot %d. Its private key is derived from a public seed, so anyone can sign as "+
		"this validator, and no honest key generation produces it. Proofs over it still succeed, so it stays selectable",
		level, validator.Address, index, validator.VotingPower, validator.PubKey.Bytes(), height, bucket, slot)
}
