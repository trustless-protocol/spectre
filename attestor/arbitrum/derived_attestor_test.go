package arbitrum

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

func TestDerivedRootAttestorPublishesConfiguredUnsafeHead(t *testing.T) {
	store, err := NewAttestedRootStore("arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	runtime := &derivedTestRuntime{
		snapshot: RuntimeSnapshot{
			Unsafe:         testStoreCommitment(120),
			UnsafeObserved: true,
			Safe:           testStoreCommitment(110),
			SafeObserved:   true,
			Finalized:      testStoreCommitment(100),
			FinalizedSeen:  true,
		},
		commitments: map[uint64]BlockCommitment{},
	}
	attestor, err := NewDerivedRootAttestor(runtime, store, DerivedAttestorConfig{
		AttestationHead: RunModeUnsafe,
		GapBlocks:       10,
		MaxRoots:        10,
	})
	if err != nil {
		t.Fatalf("create derived attestor: %v", err)
	}
	attestor.now = func() time.Time { return time.Unix(1_700_000_000, 0) }

	if err := attestor.SyncOnce(context.Background()); err != nil {
		t.Fatalf("attest unsafe head: %v", err)
	}
	root, found := store.HighestAttested(true)
	if !found || root.Source != SourceDerived || root.L2BlockNumber != 120 || !root.Provisional {
		t.Fatalf("unsafe derived frontier: found=%t root=%+v", found, root)
	}
	if _, found := store.HighestAttested(false); found {
		t.Fatal("unsafe derived root entered confirmed frontier before finalization")
	}
}

func TestDerivedRootAttestorFinalizedRecheckCorrectsDivergence(t *testing.T) {
	store, err := NewAttestedRootStore("arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	provisional := testStoreCommitment(100)
	if err := store.AppendDerived(provisional, time.Unix(1, 0), true); err != nil {
		t.Fatalf("append provisional root: %v", err)
	}
	finalized := provisional
	finalized.BlockHash = common.HexToHash("0xf1")
	finalized.StateRoot = common.HexToHash("0xf2")
	runtime := &derivedTestRuntime{
		snapshot: RuntimeSnapshot{
			Unsafe:         testStoreCommitment(120),
			UnsafeObserved: true,
			Safe:           testStoreCommitment(110),
			SafeObserved:   true,
			Finalized:      finalized,
			FinalizedSeen:  true,
		},
		commitments: map[uint64]BlockCommitment{100: finalized},
	}
	attestor, err := NewDerivedRootAttestor(runtime, store, DerivedAttestorConfig{
		AttestationHead: RunModeUnsafe,
		Disabled:        true,
		GapBlocks:       10,
		MaxRoots:        10,
	})
	if err != nil {
		t.Fatalf("create derived attestor: %v", err)
	}
	attestor.now = func() time.Time { return time.Unix(2, 0) }

	if err := attestor.SyncOnce(context.Background()); err != nil {
		t.Fatalf("recheck derived root: %v", err)
	}
	root, found := store.HighestAttested(false)
	if !found || root.Provisional || root.Root != finalized.StateRoot || root.L2BlockHash != finalized.BlockHash {
		t.Fatalf("finalized corrected root: found=%t root=%+v", found, root)
	}
}

func TestDerivedRootAttestorTracksUnsafeReorgBeforeFinalization(t *testing.T) {
	store, err := NewAttestedRootStore("arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	oldHead := testStoreCommitment(120)
	if err := store.AppendDerived(oldHead, time.Unix(1, 0), true); err != nil {
		t.Fatalf("append old head: %v", err)
	}
	newHead := oldHead
	newHead.BlockHash = common.HexToHash("0xa1")
	newHead.StateRoot = common.HexToHash("0xa2")
	runtime := &derivedTestRuntime{snapshot: RuntimeSnapshot{
		Unsafe:         newHead,
		UnsafeObserved: true,
		Finalized:      testStoreCommitment(100),
		FinalizedSeen:  true,
	}}
	attestor, err := NewDerivedRootAttestor(runtime, store, DerivedAttestorConfig{
		AttestationHead: RunModeUnsafe,
		GapBlocks:       10,
		MaxRoots:        10,
	})
	if err != nil {
		t.Fatalf("create derived attestor: %v", err)
	}
	if err := attestor.SyncOnce(context.Background()); err != nil {
		t.Fatalf("sync reorged unsafe head: %v", err)
	}
	root, found := store.HighestAttested(true)
	if !found || root.Root != newHead.StateRoot || root.L2BlockHash != newHead.BlockHash || !root.Provisional {
		t.Fatalf("reorged unsafe root: found=%t root=%+v", found, root)
	}
}

func TestDerivedRootAttestorDropsUnavailableEntryAndPersistsRemainingWork(t *testing.T) {
	path := filepath.Join(t.TempDir(), "attested-roots.json")
	store, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	for _, height := range []uint64{90, 100} {
		if err := store.AppendDerived(testStoreCommitment(height), time.Unix(1, 0), true); err != nil {
			t.Fatalf("append provisional root at %d: %v", height, err)
		}
	}
	if err := store.Save(); err != nil {
		t.Fatalf("save initial store: %v", err)
	}

	runtime := &derivedTestRuntime{
		snapshot: RuntimeSnapshot{
			Unsafe:         testStoreCommitment(120),
			UnsafeObserved: true,
			Finalized:      testStoreCommitment(100),
			FinalizedSeen:  true,
		},
		commitments: map[uint64]BlockCommitment{100: testStoreCommitment(100)},
		failures:    map[uint64]error{90: ethereum.NotFound},
	}
	attestor, err := NewDerivedRootAttestor(runtime, store, DerivedAttestorConfig{
		AttestationHead: RunModeUnsafe,
		GapBlocks:       10,
		MaxRoots:        10,
	})
	if err != nil {
		t.Fatalf("create derived attestor: %v", err)
	}
	if err := attestor.SyncOnce(context.Background()); err != nil {
		t.Fatalf("sync after unavailable historical block: %v", err)
	}

	reloaded, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("reload store: %v", err)
	}
	if _, found := reloaded.DerivedAt(90); found {
		t.Fatal("unavailable provisional root was retained")
	}
	if root, found := reloaded.DerivedAt(100); !found || root.Provisional {
		t.Fatalf("later finalized root was not confirmed: found=%t root=%+v", found, root)
	}
	if root, found := reloaded.DerivedAt(120); !found || !root.Provisional {
		t.Fatalf("new unsafe root was not persisted: found=%t root=%+v", found, root)
	}
}

func TestDerivedRootAttestorPersistsProgressAcrossTransientRecheckFailure(t *testing.T) {
	path := filepath.Join(t.TempDir(), "attested-roots.json")
	store, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	for _, height := range []uint64{90, 100} {
		if err := store.AppendDerived(testStoreCommitment(height), time.Unix(1, 0), true); err != nil {
			t.Fatalf("append provisional root at %d: %v", height, err)
		}
	}
	if err := store.Save(); err != nil {
		t.Fatalf("save initial store: %v", err)
	}

	runtime := &derivedTestRuntime{
		snapshot: RuntimeSnapshot{
			Unsafe:         testStoreCommitment(120),
			UnsafeObserved: true,
			Finalized:      testStoreCommitment(100),
			FinalizedSeen:  true,
		},
		commitments: map[uint64]BlockCommitment{100: testStoreCommitment(100)},
		failures:    map[uint64]error{90: errors.New("temporary Nitro RPC failure")},
	}
	attestor, err := NewDerivedRootAttestor(runtime, store, DerivedAttestorConfig{
		AttestationHead: RunModeUnsafe,
		GapBlocks:       10,
		MaxRoots:        10,
	})
	if err != nil {
		t.Fatalf("create derived attestor: %v", err)
	}
	if err := attestor.SyncOnce(context.Background()); err == nil {
		t.Fatal("transient recheck failure was not reported")
	}

	reloaded, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("reload store: %v", err)
	}
	if root, found := reloaded.DerivedAt(90); !found || !root.Provisional {
		t.Fatalf("transiently unavailable root was not retained: found=%t root=%+v", found, root)
	}
	if root, found := reloaded.DerivedAt(100); !found || root.Provisional {
		t.Fatalf("later finalized root was not persisted: found=%t root=%+v", found, root)
	}
	if root, found := reloaded.DerivedAt(120); !found || !root.Provisional {
		t.Fatalf("new unsafe root was not persisted: found=%t root=%+v", found, root)
	}
}

type derivedTestRuntime struct {
	snapshot    RuntimeSnapshot
	commitments map[uint64]BlockCommitment
	failures    map[uint64]error
}

func (r *derivedTestRuntime) Snapshot() RuntimeSnapshot { return r.snapshot }

func (r *derivedTestRuntime) CommitmentAt(
	_ context.Context,
	height uint64,
) (BlockCommitment, error) {
	if err := r.failures[height]; err != nil {
		return BlockCommitment{}, err
	}
	return r.commitments[height], nil
}
