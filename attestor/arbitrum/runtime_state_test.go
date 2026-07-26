package arbitrum

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func TestRunModeUsesFinalizedNaming(t *testing.T) {
	mode, err := RunMode("finalized").Normalize()
	if err != nil {
		t.Fatalf("normalize finalized mode: %v", err)
	}
	if mode != RunModeFinalized {
		t.Fatalf("normalized mode: got %q want %q", mode, RunModeFinalized)
	}

	if _, err := RunMode("finalize").Normalize(); err == nil {
		t.Fatal("legacy finalize mode must be rejected")
	}
	if _, err := RunMode("").Normalize(); err == nil {
		t.Fatal("empty internal run mode must be rejected")
	}
}

func TestRuntimeStateTracksFinalizedConsistency(t *testing.T) {
	reader := newRuntimeTestReader()
	reader.setHeader(testHeader(1, 0x01))
	reader.setHeader(testHeader(2, 0x02))
	reader.setHeader(testHeader(3, 0x03))
	reader.setHeader(testHeader(4, 0x04))
	reader.setHeads(3, 2, 1)

	state, err := NewRuntimeState(reader)
	if err != nil {
		t.Fatalf("create runtime state: %v", err)
	}
	checks, err := state.Refresh(context.Background())
	if err != nil {
		t.Fatalf("initial refresh: %v", err)
	}
	if len(checks) != 0 {
		t.Fatalf("initial refresh produced checks: %+v", checks)
	}

	reader.setHeads(3, 2, 2)
	checks, err = state.Refresh(context.Background())
	if err != nil {
		t.Fatalf("finalize height 2: %v", err)
	}
	if len(checks) != 1 || !checks[0].Consistent() {
		t.Fatalf("height 2 consistency: %+v", checks)
	}

	// Height 3 was first observed with root 0x03 while unsafe. Before it
	// becomes safe, replace the canonical header with root 0x33 to model an
	// unsafe reorg. The first unsafe observation must be retained so finalization
	// reports the mismatch while the safe observation matches.
	reader.setHeader(testHeader(3, 0x33))
	reader.setHeads(4, 3, 2)
	if _, err := state.Refresh(context.Background()); err != nil {
		t.Fatalf("record reorged safe height 3: %v", err)
	}

	reader.setHeads(4, 3, 3)
	checks, err = state.Refresh(context.Background())
	if err != nil {
		t.Fatalf("finalize height 3: %v", err)
	}
	if len(checks) != 1 {
		t.Fatalf("height 3 checks: got %d want 1", len(checks))
	}
	check := checks[0]
	if !check.UnsafeObserved || !check.SafeObserved {
		t.Fatalf("height 3 observations missing: %+v", check)
	}
	if check.UnsafeMatches {
		t.Fatalf("reorged unsafe root unexpectedly matched: %+v", check)
	}
	if !check.SafeMatches {
		t.Fatalf("safe root did not match final root: %+v", check)
	}
	if check.Consistent() {
		t.Fatalf("mismatched finalization reported consistent: %+v", check)
	}

	snapshot := state.Snapshot()
	if !snapshot.UnsafeObserved || snapshot.Unsafe.BlockNumber != 4 {
		t.Fatalf("unsafe snapshot: %+v", snapshot)
	}
	if !snapshot.SafeObserved || snapshot.Safe.BlockNumber != 3 {
		t.Fatalf("safe snapshot: %+v", snapshot)
	}
	if !snapshot.FinalizedSeen || snapshot.Finalized.BlockNumber != 3 {
		t.Fatalf("finalized snapshot: %+v", snapshot)
	}
	if !snapshot.HasLastFinalizedCheck || snapshot.LastFinalizedCheck.Finalized.BlockNumber != 3 {
		t.Fatalf("finalized check snapshot: %+v", snapshot)
	}
}

func TestRuntimeStateRejectsInconsistentHeadOrdering(t *testing.T) {
	reader := newRuntimeTestReader()
	reader.setHeader(testHeader(1, 0x01))
	reader.setHeader(testHeader(2, 0x02))
	reader.setHeads(1, 2, 1)

	state, err := NewRuntimeState(reader)
	if err != nil {
		t.Fatalf("create runtime state: %v", err)
	}
	if _, err := state.Refresh(context.Background()); err == nil {
		t.Fatal("expected safe-above-unsafe ordering error")
	}
}

func TestResolveCanonicalBlockHashHonorsConfirmationMode(t *testing.T) {
	reader := newRuntimeTestReader()
	for height := uint64(1); height <= 5; height++ {
		reader.setHeader(testHeader(height, byte(height)))
	}
	reader.setHeads(5, 5, 4)
	state, err := NewRuntimeState(reader)
	if err != nil {
		t.Fatalf("create runtime state: %v", err)
	}

	blockFive := reader.header(5)
	commitment, err := state.ResolveCanonicalBlockHash(
		context.Background(),
		blockFive.Hash(),
		RunModeSafe,
	)
	if err != nil {
		t.Fatalf("resolve safe block: %v", err)
	}
	if commitment.BlockNumber != 5 || commitment.StateRoot != blockFive.Root {
		t.Fatalf("safe commitment: %+v", commitment)
	}

	_, err = state.ResolveCanonicalBlockHash(
		context.Background(),
		blockFive.Hash(),
		RunModeFinalized,
	)
	if !errors.Is(err, ErrCommitmentNotReady) {
		t.Fatalf("block above finalized head: got %v want not ready", err)
	}
}

func TestResolveCanonicalBlockHashRejectsReorgedHeader(t *testing.T) {
	reader := newRuntimeTestReader()
	oldHeader := testHeader(3, 0x03)
	reader.setHeader(oldHeader)
	reader.setHeader(testHeader(4, 0x04))
	reader.setHeads(4, 3, 3)
	reader.setHeader(testHeader(3, 0x33))

	state, err := NewRuntimeState(reader)
	if err != nil {
		t.Fatalf("create runtime state: %v", err)
	}
	_, err = state.ResolveCanonicalBlockHash(
		context.Background(),
		oldHeader.Hash(),
		RunModeSafe,
	)
	if !errors.Is(err, ErrCommitmentMismatch) {
		t.Fatalf("reorged block: got %v want mismatch", err)
	}
}

func testHeader(height uint64, root byte) *types.Header {
	return &types.Header{
		Number: new(big.Int).SetUint64(height),
		Root:   [32]byte{31: root},
	}
}

type runtimeTestReader struct {
	mu        sync.RWMutex
	unsafe    uint64
	safe      uint64
	finalized uint64
	headers   map[uint64]*types.Header
	byHash    map[common.Hash]*types.Header
}

func newRuntimeTestReader() *runtimeTestReader {
	return &runtimeTestReader{
		headers: make(map[uint64]*types.Header),
		byHash:  make(map[common.Hash]*types.Header),
	}
}

func (r *runtimeTestReader) setHeads(unsafe, safe, finalized uint64) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.unsafe = unsafe
	r.safe = safe
	r.finalized = finalized
}

func (r *runtimeTestReader) setHeader(header *types.Header) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.headers[header.Number.Uint64()] = types.CopyHeader(header)
	r.byHash[header.Hash()] = types.CopyHeader(header)
}

func (r *runtimeTestReader) header(height uint64) *types.Header {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return types.CopyHeader(r.headers[height])
}

func (r *runtimeTestReader) HeaderByNumber(_ context.Context, number *big.Int) (*types.Header, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var height uint64
	switch {
	case number == nil:
		height = r.unsafe
	case number.Int64() == int64(rpc.SafeBlockNumber):
		height = r.safe
	case number.Int64() == int64(rpc.FinalizedBlockNumber):
		height = r.finalized
	case number.Sign() >= 0 && number.IsUint64():
		height = number.Uint64()
	default:
		return nil, fmt.Errorf("unsupported block selector %s", number)
	}
	header, ok := r.headers[height]
	if !ok {
		return nil, fmt.Errorf("header %d not configured", height)
	}
	return types.CopyHeader(header), nil
}

func (r *runtimeTestReader) HeaderByHash(
	_ context.Context,
	hash common.Hash,
) (*types.Header, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	header, ok := r.byHash[hash]
	if !ok {
		return nil, ethereum.NotFound
	}
	return types.CopyHeader(header), nil
}
