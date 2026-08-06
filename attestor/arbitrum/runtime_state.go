package arbitrum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// NitroHeaderReader is the canonical Nitro RPC surface required by the
// attestor runtime and verification service.
type NitroHeaderReader interface {
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
}

type nitroHeaderByHashReader interface {
	HeaderByHash(context.Context, common.Hash) (*types.Header, error)
}

var (
	// ErrCommitmentNotReady means Nitro has not yet derived the assertion's
	// block at the requested confirmation level.
	ErrCommitmentNotReady = errors.New("Nitro commitment is not ready")

	// ErrCommitmentMismatch means the assertion's block is not canonical in
	// the configured Nitro view.
	ErrCommitmentMismatch = errors.New("Nitro commitment does not match")
)

// RunMode selects the Nitro confirmation level used for a state-root check.
type RunMode string

const (
	RunModeUnsafe    RunMode = "unsafe"
	RunModeSafe      RunMode = "safe"
	RunModeFinalized RunMode = "finalized"
)

// Normalize validates an internal run mode.
func (m RunMode) Normalize() (RunMode, error) {
	switch m {
	case RunModeUnsafe:
		return RunModeUnsafe, nil
	case RunModeSafe:
		return RunModeSafe, nil
	case RunModeFinalized:
		return RunModeFinalized, nil
	default:
		return "", fmt.Errorf(
			"attestor run mode must be one of unsafe, safe, or finalized, got %q",
			m,
		)
	}
}

func (m RunMode) headSelector() *big.Int {
	switch m {
	case RunModeSafe:
		return big.NewInt(int64(rpc.SafeBlockNumber))
	case RunModeFinalized:
		return big.NewInt(int64(rpc.FinalizedBlockNumber))
	default:
		return nil
	}
}

// BlockCommitment is the canonical Nitro commitment observed at one L2 height.
type BlockCommitment struct {
	BlockNumber uint64
	BlockHash   common.Hash
	StateRoot   common.Hash
}

// FinalizedConsistency records whether the root that became finalized matches
// the roots this attestor previously observed while the same height was unsafe
// and safe. Missing observations are represented explicitly so a future alert
// monitor can distinguish incomplete history from a commitment mismatch.
type FinalizedConsistency struct {
	Finalized BlockCommitment

	Unsafe         BlockCommitment
	UnsafeObserved bool
	UnsafeMatches  bool

	Safe         BlockCommitment
	SafeObserved bool
	SafeMatches  bool

	PreviousFinalized    BlockCommitment
	HasPreviousFinalized bool
	FinalizedReorg       bool
}

// Consistent reports whether both pre-finality observations exist and match the
// finalized state root, and the finalized head itself did not regress or change.
func (c FinalizedConsistency) Consistent() bool {
	return c.UnsafeObserved && c.SafeObserved && c.UnsafeMatches && c.SafeMatches && !c.FinalizedReorg
}

// RuntimeSnapshot is a point-in-time copy of the tracked Nitro heads and the
// latest finalized consistency result.
type RuntimeSnapshot struct {
	Unsafe         BlockCommitment
	UnsafeObserved bool
	Safe           BlockCommitment
	SafeObserved   bool
	Finalized      BlockCommitment
	FinalizedSeen  bool

	LastFinalizedCheck    FinalizedConsistency
	HasLastFinalizedCheck bool
}

// RuntimeState tracks Nitro's unsafe, safe, and finalized L2 commitments in
// memory. The tagged safe/finalized heads reflect the corresponding L1 data
// availability/finality view exposed by Nitro.
//
// Unsafe and safe observations are retained only until their heights finalize.
// This bounds memory by the current unsafe-to-finalized window while preserving
// the first root observed at each height for later consistency monitoring.
type RuntimeState struct {
	reader NitroHeaderReader

	refreshMu sync.Mutex
	mu        sync.RWMutex

	initialized bool
	heads       map[RunMode]BlockCommitment
	observed    map[RunMode]map[uint64]BlockCommitment

	lastFinalizedCheck    FinalizedConsistency
	hasLastFinalizedCheck bool
}

// NewRuntimeState constructs an empty tracker backed by the configured Nitro
// RPC endpoint.
func NewRuntimeState(reader NitroHeaderReader) (*RuntimeState, error) {
	if reader == nil {
		return nil, errors.New("Nitro header reader must not be nil")
	}
	return &RuntimeState{
		reader: reader,
		heads:  make(map[RunMode]BlockCommitment, 3),
		observed: map[RunMode]map[uint64]BlockCommitment{
			RunModeUnsafe: make(map[uint64]BlockCommitment),
			RunModeSafe:   make(map[uint64]BlockCommitment),
		},
	}, nil
}

// Head reads the current Nitro head for the requested confirmation mode.
func (s *RuntimeState) Head(ctx context.Context, mode RunMode) (BlockCommitment, error) {
	if s == nil || s.reader == nil {
		return BlockCommitment{}, errors.New("runtime state is not initialized")
	}
	normalized, err := mode.Normalize()
	if err != nil {
		return BlockCommitment{}, err
	}
	header, err := s.reader.HeaderByNumber(ctx, normalized.headSelector())
	if isNitroBlockNotFound(err) || (err == nil && header == nil) {
		return BlockCommitment{}, fmt.Errorf("Nitro %s head was not found: %w", normalized, ethereum.NotFound)
	}
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("query Nitro %s head: %w", normalized, err)
	}
	commitment, err := commitmentFromHeader(header)
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("invalid Nitro %s head: %w", normalized, err)
	}
	return commitment, nil
}

// CommitmentAt reads Nitro's current canonical commitment at one L2 height.
func (s *RuntimeState) CommitmentAt(ctx context.Context, height uint64) (BlockCommitment, error) {
	if s == nil || s.reader == nil {
		return BlockCommitment{}, errors.New("runtime state is not initialized")
	}
	header, err := s.reader.HeaderByNumber(ctx, new(big.Int).SetUint64(height))
	if isNitroBlockNotFound(err) || (err == nil && header == nil) {
		return BlockCommitment{}, fmt.Errorf("Nitro block %d was not found: %w", height, ethereum.NotFound)
	}
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("query Nitro block %d: %w", height, err)
	}
	commitment, err := commitmentFromHeader(header)
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("invalid Nitro block %d: %w", height, err)
	}
	if commitment.BlockNumber != height {
		return BlockCommitment{}, fmt.Errorf(
			"invalid Nitro header number: got %d want %d",
			commitment.BlockNumber,
			height,
		)
	}
	return commitment, nil
}

// ResolveCanonicalBlockHash resolves an assertion's L2 block hash and proves
// that the same block is canonical at the requested Nitro confirmation level.
//
// Looking the block up by hash recovers its L2 height from the assertion. The
// subsequent height lookup prevents a retained, reorged header from being
// accepted merely because Nitro still has it in its database.
func (s *RuntimeState) ResolveCanonicalBlockHash(
	ctx context.Context,
	blockHash common.Hash,
	mode RunMode,
) (BlockCommitment, error) {
	if s == nil || s.reader == nil {
		return BlockCommitment{}, errors.New("runtime state is not initialized")
	}
	if blockHash == (common.Hash{}) {
		return BlockCommitment{}, fmt.Errorf("%w: block hash is zero", ErrCommitmentMismatch)
	}
	normalized, err := mode.Normalize()
	if err != nil {
		return BlockCommitment{}, err
	}
	hashReader, ok := s.reader.(nitroHeaderByHashReader)
	if !ok {
		return BlockCommitment{}, errors.New("Nitro reader does not support block lookup by hash")
	}

	header, err := hashReader.HeaderByHash(ctx, blockHash)
	if isNitroBlockNotFound(err) || (err == nil && header == nil) {
		return BlockCommitment{}, fmt.Errorf(
			"%w: Nitro block %s was not found",
			ErrCommitmentNotReady,
			blockHash,
		)
	}
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("query Nitro block %s: %w", blockHash, err)
	}
	resolved, err := commitmentFromHeader(header)
	if err != nil {
		return BlockCommitment{}, fmt.Errorf("invalid Nitro block %s: %w", blockHash, err)
	}
	if resolved.BlockHash != blockHash {
		return BlockCommitment{}, fmt.Errorf(
			"%w: block-hash lookup returned %s for %s",
			ErrCommitmentMismatch,
			resolved.BlockHash,
			blockHash,
		)
	}

	head, err := s.Head(ctx, normalized)
	if errors.Is(err, ethereum.NotFound) {
		return BlockCommitment{}, fmt.Errorf(
			"%w: Nitro %s head is unavailable",
			ErrCommitmentNotReady,
			normalized,
		)
	}
	if err != nil {
		return BlockCommitment{}, err
	}
	if resolved.BlockNumber > head.BlockNumber {
		return BlockCommitment{}, fmt.Errorf(
			"%w: block %d is above Nitro %s head %d",
			ErrCommitmentNotReady,
			resolved.BlockNumber,
			normalized,
			head.BlockNumber,
		)
	}

	canonical, err := s.CommitmentAt(ctx, resolved.BlockNumber)
	if errors.Is(err, ethereum.NotFound) {
		return BlockCommitment{}, fmt.Errorf(
			"%w: canonical Nitro block %d is unavailable",
			ErrCommitmentNotReady,
			resolved.BlockNumber,
		)
	}
	if err != nil {
		return BlockCommitment{}, err
	}
	if canonical.BlockHash != blockHash {
		return BlockCommitment{}, fmt.Errorf(
			"%w: assertion block %s is not canonical at height %d (canonical %s)",
			ErrCommitmentMismatch,
			blockHash,
			resolved.BlockNumber,
			canonical.BlockHash,
		)
	}
	return canonical, nil
}

func isNitroBlockNotFound(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, ethereum.NotFound) {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "block not found")
}

// Refresh records all newly crossed unsafe and safe heights and compares every
// newly finalized commitment with the first roots observed at those levels.
// The initial refresh records only the current heads; it does not perform a
// historical gap backfill against the public RPC. The returned checks are
// suitable for logging now and alerting later.
func (s *RuntimeState) Refresh(ctx context.Context) ([]FinalizedConsistency, error) {
	if s == nil || s.reader == nil {
		return nil, errors.New("runtime state is not initialized")
	}

	s.refreshMu.Lock()
	defer s.refreshMu.Unlock()

	heads := make(map[RunMode]BlockCommitment, 3)
	for _, mode := range []RunMode{RunModeUnsafe, RunModeSafe, RunModeFinalized} {
		head, err := s.Head(ctx, mode)
		if err != nil {
			return nil, err
		}
		heads[mode] = head
	}
	if heads[RunModeFinalized].BlockNumber > heads[RunModeSafe].BlockNumber {
		return nil, fmt.Errorf(
			"Nitro finalized head %d is above safe head %d",
			heads[RunModeFinalized].BlockNumber,
			heads[RunModeSafe].BlockNumber,
		)
	}
	if heads[RunModeSafe].BlockNumber > heads[RunModeUnsafe].BlockNumber {
		return nil, fmt.Errorf(
			"Nitro safe head %d is above unsafe head %d",
			heads[RunModeSafe].BlockNumber,
			heads[RunModeUnsafe].BlockNumber,
		)
	}

	s.mu.RLock()
	initialized := s.initialized
	previousHeads := make(map[RunMode]BlockCommitment, len(s.heads))
	for mode, head := range s.heads {
		previousHeads[mode] = head
	}
	s.mu.RUnlock()

	unsafeObservations := []BlockCommitment{heads[RunModeUnsafe]}
	safeObservations := []BlockCommitment{heads[RunModeSafe]}
	var err error
	if initialized {
		unsafeObservations, err = s.commitmentRange(
			ctx,
			nextHeight(previousHeads[RunModeUnsafe].BlockNumber),
			heads[RunModeUnsafe],
		)
		if err != nil {
			return nil, fmt.Errorf("record unsafe commitments: %w", err)
		}
		safeObservations, err = s.commitmentRange(
			ctx,
			nextHeight(previousHeads[RunModeSafe].BlockNumber),
			heads[RunModeSafe],
		)
		if err != nil {
			return nil, fmt.Errorf("record safe commitments: %w", err)
		}
	}

	var newlyFinalized []BlockCommitment
	if initialized {
		previousFinalized := previousHeads[RunModeFinalized]
		switch {
		case heads[RunModeFinalized].BlockNumber > previousFinalized.BlockNumber:
			newlyFinalized, err = s.commitmentRange(
				ctx,
				nextHeight(previousFinalized.BlockNumber),
				heads[RunModeFinalized],
			)
			if err != nil {
				return nil, fmt.Errorf("record finalized commitments: %w", err)
			}
		case heads[RunModeFinalized] != previousFinalized:
			newlyFinalized = []BlockCommitment{heads[RunModeFinalized]}
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	recordFirstObservations(s.observed[RunModeUnsafe], unsafeObservations)
	recordFirstObservations(s.observed[RunModeSafe], safeObservations)

	checks := make([]FinalizedConsistency, 0, len(newlyFinalized))
	for _, finalized := range newlyFinalized {
		check := s.compareFinalized(finalized, previousHeads[RunModeFinalized], initialized)
		checks = append(checks, check)
		s.lastFinalizedCheck = check
		s.hasLastFinalizedCheck = true
	}

	for mode, head := range heads {
		s.heads[mode] = head
	}
	s.initialized = true
	pruneFinalized(s.observed[RunModeUnsafe], heads[RunModeFinalized].BlockNumber)
	pruneFinalized(s.observed[RunModeSafe], heads[RunModeFinalized].BlockNumber)

	return checks, nil
}

// Snapshot returns the latest tracked heads and finalization comparison.
func (s *RuntimeState) Snapshot() RuntimeSnapshot {
	if s == nil {
		return RuntimeSnapshot{}
	}
	s.mu.RLock()
	defer s.mu.RUnlock()

	snapshot := RuntimeSnapshot{
		LastFinalizedCheck:    s.lastFinalizedCheck,
		HasLastFinalizedCheck: s.hasLastFinalizedCheck,
	}
	snapshot.Unsafe, snapshot.UnsafeObserved = s.heads[RunModeUnsafe]
	snapshot.Safe, snapshot.SafeObserved = s.heads[RunModeSafe]
	snapshot.Finalized, snapshot.FinalizedSeen = s.heads[RunModeFinalized]
	return snapshot
}

func (s *RuntimeState) commitmentRange(
	ctx context.Context,
	from uint64,
	head BlockCommitment,
) ([]BlockCommitment, error) {
	if from > head.BlockNumber {
		return nil, nil
	}
	commitments := make([]BlockCommitment, 0, head.BlockNumber-from+1)
	for height := from; ; height++ {
		if height == head.BlockNumber {
			commitments = append(commitments, head)
			break
		}
		commitment, err := s.CommitmentAt(ctx, height)
		if err != nil {
			return nil, err
		}
		commitments = append(commitments, commitment)
	}
	return commitments, nil
}

func (s *RuntimeState) compareFinalized(
	finalized BlockCommitment,
	previousFinalized BlockCommitment,
	hadPrevious bool,
) FinalizedConsistency {
	check := FinalizedConsistency{Finalized: finalized}
	if hadPrevious {
		check.PreviousFinalized = previousFinalized
		check.HasPreviousFinalized = true
		check.FinalizedReorg = finalized.BlockNumber < previousFinalized.BlockNumber ||
			(finalized.BlockNumber == previousFinalized.BlockNumber && finalized != previousFinalized)
	}

	check.Unsafe, check.UnsafeObserved = s.observed[RunModeUnsafe][finalized.BlockNumber]
	check.Safe, check.SafeObserved = s.observed[RunModeSafe][finalized.BlockNumber]
	check.UnsafeMatches = check.UnsafeObserved && check.Unsafe.StateRoot == finalized.StateRoot
	check.SafeMatches = check.SafeObserved && check.Safe.StateRoot == finalized.StateRoot
	return check
}

func commitmentFromHeader(header *types.Header) (BlockCommitment, error) {
	if header == nil {
		return BlockCommitment{}, errors.New("header must not be nil")
	}
	if header.Number == nil || !header.Number.IsUint64() {
		return BlockCommitment{}, errors.New("header number must be a uint64")
	}
	return BlockCommitment{
		BlockNumber: header.Number.Uint64(),
		BlockHash:   header.Hash(),
		StateRoot:   header.Root,
	}, nil
}

func recordFirstObservations(destination map[uint64]BlockCommitment, commitments []BlockCommitment) {
	for _, commitment := range commitments {
		if _, exists := destination[commitment.BlockNumber]; !exists {
			destination[commitment.BlockNumber] = commitment
		}
	}
}

func pruneFinalized(commitments map[uint64]BlockCommitment, finalizedHeight uint64) {
	for height := range commitments {
		if height <= finalizedHeight {
			delete(commitments, height)
		}
	}
}

func nextHeight(height uint64) uint64 {
	if height == ^uint64(0) {
		return height
	}
	return height + 1
}
