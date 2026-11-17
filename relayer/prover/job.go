package prover

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"relayer/prover/witness"
)

// ProofID identifies one proof job. It is derived, not assigned:
//
//	ProofID = SHA-256(vkDigest ‖ witnessDigest)
//
// The witness digest says what is being proven; the vk digest says with which
// circuit. Namespacing by vk is load-bearing: after a circuit re-setup, a
// proof built under the old vk must never be returned for a new submission —
// it would revert on-chain instead of failing here. Determinism gives
// idempotent submits (same signatures + same circuit → same id, no duplicate
// prove) and crash-safety (a re-submit after restart reproduces the id).
type ProofID [32]byte

// ClaimToken fences writes from an earlier owner after a job is reclaimed.
// A store assigns a new token on every successful Claim.
type ClaimToken uint64

// NewProofID derives the job id for a witness digest under a circuit version.
func NewProofID(vkDigest, witnessDigest [32]byte) ProofID {
	h := sha256.New()
	h.Write(vkDigest[:])
	h.Write(witnessDigest[:])
	var id ProofID
	copy(id[:], h.Sum(nil))
	return id
}

func (id ProofID) String() string { return hex.EncodeToString(id[:]) }

// Status is the lifecycle of a proof job:
//
//	queued → claimed → proving → done | failed
//
// A worker takes a job with an atomic Claim carrying an expiring lease;
// a worker dying mid-job leaves the lease to expire and the job reclaimable,
// so nothing is ever stranded at claimed/proving.
type Status string

const (
	StatusQueued  Status = "queued"
	StatusClaimed Status = "claimed"
	StatusProving Status = "proving"
	StatusDone    Status = "done"
	StatusFailed  Status = "failed"
)

// FailureClass tells the relayer what a failed job means for its retry
// policy — the prover itself never retries.
type FailureClass string

const (
	// FailureTransient covers timeouts and flaky backends: submitting the same
	// request again is reasonable.
	FailureTransient FailureClass = "transient"
	// FailureInternal covers circuit/artifact/assignment bugs: retrying is
	// useless, a human needs to look.
	FailureInternal FailureClass = "internal"
)

// ProofRequest is the unit of work that travels through the JobQueue. All
// fields are plain values — no pointers to live state, no handles — so any
// queue backend can serialize it.
type ProofRequest struct {
	ID            ProofID
	Bucket        int
	Sigs          []witness.ValidatorSignature // padded to Bucket slots
	WitnessDigest [32]byte
	EnqueuedAt    time.Time // set at submit; workers derive the queue-wait metric from it
}

// ProofJob is what GetProof returns: the job's current status plus its result
// or failure, whichever applies.
type ProofJob struct {
	ID     ProofID
	Status Status
	Result *ProofResult // set when Status == done
	Class  FailureClass // set when Status == failed
	Cause  string       // set when Status == failed
}

// ProofResult is a finished Groth16 batch proof, shaped for the on-chain
// verifier dispatch: the bucket selects Groth16Verifier_N{bucket}, PaddedSigs
// carries the bucket-sized signer metadata that travels with the proof.
type ProofResult struct {
	Bucket        int
	PaddedSigs    []witness.ValidatorSignature
	Proof         [8]*big.Int
	Commitments   [4]*big.Int // C0 then C1, each an (x, y) pair
	CommitmentPok [2]*big.Int
}

// ErrProofFailed is returned by WaitProof when a worker stored a terminal
// failure. The class remains machine-readable so one-shot callers can apply
// the same retry table as callers that poll GetProof directly.
type ErrProofFailed struct {
	ID    ProofID
	Class FailureClass
	Cause string
}

func (e *ErrProofFailed) Error() string {
	return fmt.Sprintf("proof %s failed (%s): %s", e.ID, e.Class, e.Cause)
}

// Port errors. Adapters return these sentinel values (wrapped is fine) so the
// service and the relayer can classify without knowing the backend.
var (
	// ErrQueueFull: the bounded queue rejected the job. The caller backs off;
	// nothing was silently dropped.
	ErrQueueFull = errors.New("proof queue full")
	// ErrNotFound: the store has no entry for this id — never submitted,
	// TTL-collected, or lost to a restart. Re-submitting reproduces the id.
	ErrNotFound = errors.New("proof not found")
	// ErrStoreUnavailable: the store backend cannot be reached right now.
	ErrStoreUnavailable = errors.New("proof store unavailable")
	// ErrStaleClaim means the caller no longer owns the job or its transition is invalid.
	ErrStaleClaim = errors.New("proof claim superseded or transition invalid")
)
