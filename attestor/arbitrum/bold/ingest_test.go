package bold

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	arbitrum "attestor/arbitrum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type headTestSource struct {
	assertionTestSource
	head            uint64
	events          []arbitrum.ProposedAssertion
	confirmed       bool
	headErr         error
	canonicalCheck  func()
	hashErr         error
	changedAt       uint64
	failRangeFrom   uint64
	statusHashErr   error
	statusNumberErr error
	heads           []arbitrum.L1IngestHead
	statusBlocks    []uint64
	statusHashes    []common.Hash
	ranges          [][2]uint64
}

func blockHash(n uint64) common.Hash { return common.BigToHash(new(big.Int).SetUint64(n)) }

func (s *headTestSource) Head(_ context.Context, h arbitrum.L1IngestHead) (L1Block, error) {
	s.heads = append(s.heads, h)
	if h == arbitrum.L1IngestFinalized {
		return L1Block{Number: s.finalizedBlock, Hash: blockHash(s.finalizedBlock)}, s.finalizedErr
	}
	return L1Block{Number: s.head, Hash: blockHash(s.head)}, s.headErr
}

func (s *headTestSource) CanonicalHash(_ context.Context, n uint64) (common.Hash, error) {
	if s.canonicalCheck != nil {
		s.canonicalCheck()
	}
	if s.changedAt == n {
		return common.HexToHash("0xbad"), nil
	}
	return blockHash(n), s.hashErr
}

func (s *headTestSource) Assertions(_ context.Context, from, to uint64) ([]arbitrum.ProposedAssertion, []ConfirmedAssertion, error) {
	s.ranges = append(s.ranges, [2]uint64{from, to})
	if s.failRangeFrom != 0 && from >= s.failRangeFrom {
		return nil, nil, errors.New("log range failed")
	}
	var proposals []arbitrum.ProposedAssertion
	var confirmations []ConfirmedAssertion
	for _, p := range s.events {
		if p.L1BlockNumber >= from && p.L1BlockNumber <= to {
			proposals = append(proposals, p)
		}
		if s.confirmed && s.head >= from && s.head <= to {
			confirmations = append(confirmations, ConfirmedAssertion{AssertionHash: p.AssertionHash, L2BlockHash: p.L2BlockHash, L1BlockNumber: s.head})
		}
	}
	return proposals, confirmations, nil
}

func (s *headTestSource) status(p arbitrum.ProposedAssertion, n uint64) uint8 {
	for _, event := range s.events {
		if event.AssertionHash == p.AssertionHash && event.L1BlockNumber <= n {
			if s.confirmed && s.head <= n {
				return assertionStatusConfirmed
			}
			return assertionStatusPending
		}
	}
	return assertionStatusNone
}

func (s *headTestSource) AssertionStatus(_ context.Context, p arbitrum.ProposedAssertion, n uint64) (uint8, error) {
	s.statusBlocks = append(s.statusBlocks, n)
	return s.status(p, n), s.statusNumberErr
}

func (s *headTestSource) AssertionStatusAtHash(_ context.Context, p arbitrum.ProposedAssertion, h common.Hash) (uint8, error) {
	s.statusHashes = append(s.statusHashes, h)
	return s.status(p, h.Big().Uint64()), s.statusHashErr
}

type headTestResolver struct {
	finalized bool
	err       error
	modes     []arbitrum.RunMode
}

func (r *headTestResolver) ResolveCanonicalBlockHash(_ context.Context, hash common.Hash, mode arbitrum.RunMode) (arbitrum.BlockCommitment, error) {
	r.modes = append(r.modes, mode)
	if r.err != nil {
		return arbitrum.BlockCommitment{}, r.err
	}
	if mode == arbitrum.RunModeFinalized && !r.finalized {
		return arbitrum.BlockCommitment{}, arbitrum.ErrCommitmentNotReady
	}
	return arbitrum.BlockCommitment{BlockNumber: 500, BlockHash: hash, StateRoot: common.HexToHash("0x777")}, nil
}

func earlyProposal(hash int64, block uint64) arbitrum.ProposedAssertion {
	return arbitrum.ProposedAssertion{AssertionHash: common.BigToHash(big.NewInt(hash)), L2BlockHash: common.BigToHash(big.NewInt(hash + 1000)), L1BlockNumber: block}
}

func newHeadTestLoop(t *testing.T, mode arbitrum.L1IngestHead) (*AssertionAttestor, *headTestSource, *headTestResolver) {
	t.Helper()
	store, err := arbitrum.NewAttestedRootStore("arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	s := &headTestSource{assertionTestSource: assertionTestSource{finalizedBlock: 100}, head: 110, events: []arbitrum.ProposedAssertion{earlyProposal(1, 110)}}
	r := &headTestResolver{}
	a, err := NewAssertionAttestor(s, r, store, AssertionAttestorConfig{L1IngestHead: mode, MaxL1BlockRange: 5})
	if err != nil {
		t.Fatal(err)
	}
	a.now = func() time.Time { return time.Unix(1000, 0) }
	return a, s, r
}

func TestEarlyIngestSurvivesUntilFinalityAndPromotesOnlyAtFinalized(t *testing.T) {
	for _, mode := range []arbitrum.L1IngestHead{arbitrum.L1IngestSafe, arbitrum.L1IngestLatest} {
		t.Run(string(mode), func(t *testing.T) {
			a, s, r := newHeadTestLoop(t, mode)
			s.confirmed = true // Even a confirmation in the speculative window is insufficient.
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			root, found := a.store.HighestAttested(true)
			if !found || !root.Provisional || root.AssertionConfirmed {
				t.Fatalf("early root = %+v, found=%v", root, found)
			}
			if _, found := a.store.HighestAttested(false); found {
				t.Fatal("early assertion entered finalized frontier")
			}
			if a.store.NextL1Block() != 101 {
				t.Fatal("speculation advanced durable cursor")
			}
			if !reflect.DeepEqual(s.heads, []arbitrum.L1IngestHead{arbitrum.L1IngestFinalized, mode}) {
				t.Fatalf("heads = %v", s.heads)
			}
			if !reflect.DeepEqual(s.statusHashes, []common.Hash{blockHash(110)}) {
				t.Fatalf("hash reads = %v", s.statusHashes)
			}
			a.now = func() time.Time { return time.Unix(2000, 0) }
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			repeated, found := a.store.HighestAttested(true)
			if !found || repeated.L2BlockNumber < root.L2BlockNumber || repeated.AttestedAt != root.AttestedAt {
				t.Fatal("unchanged entry lost timestamp")
			}
			s.finalizedBlock = 110
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			if _, found := a.store.HighestAttested(false); found {
				t.Fatal("promoted before Nitro finalized")
			}
			beforePromotion, found := a.store.HighestAttested(true)
			if !found || beforePromotion.L2BlockNumber < repeated.L2BlockNumber {
				t.Fatal("frontier regressed while waiting for Nitro finalization")
			}
			r.finalized = true
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			promoted, found := a.store.HighestAttested(false)
			if !found || promoted.L2BlockNumber != root.L2BlockNumber || promoted.AttestedAt != root.AttestedAt {
				t.Fatalf("promotion = %+v", promoted)
			}
			frontier, found := a.store.HighestAttested(true)
			if !found || frontier.L2BlockNumber < beforePromotion.L2BlockNumber {
				t.Fatal("frontier regressed on promotion")
			}
			for _, n := range s.statusBlocks {
				if n != 110 {
					t.Fatalf("finalized reconciliation used %d", n)
				}
			}
		})
	}
}

func TestEarlyIngestReorgRepostAndHeadRegression(t *testing.T) {
	a, s, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	ctx := context.Background()
	metrics := arbitrum.NewMetrics(prometheus.NewRegistry())
	a.store.SetMetrics(metrics)
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	original, _ := a.store.HighestAttested(true)
	s.events[0].L1BlockNumber = 115
	s.head = 115
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	reposted, _ := a.store.HighestAttested(true)
	if reposted.L1BlockNumber != 115 || reposted.AttestedAt != original.AttestedAt {
		t.Fatalf("reposted = %+v", reposted)
	}
	s.events = []arbitrum.ProposedAssertion{earlyProposal(2, 113)}
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	replaced, _ := a.store.HighestAttested(true)
	if replaced.AssertionHash != s.events[0].AssertionHash || len(a.store.ProvisionalAssertions()) != 1 {
		t.Fatalf("replacement = %+v", replaced)
	}
	if a.store.NextL1Block() != 101 {
		t.Fatal("reorg skipped durable window")
	}
	if metricValue(metrics.AssertionRevocations.WithLabelValues("arb")) != 1 {
		t.Fatal("replacement revocation not counted")
	}
	s.head = 105
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if _, found := a.store.HighestAttested(true); found {
		t.Fatal("root above regressed head retained")
	}
	if metricValue(metrics.AssertionIngestHead.WithLabelValues("arb")) != 105 {
		t.Fatal("successful head not observed")
	}
	if metricValue(metrics.AssertionRevocations.WithLabelValues("arb")) != 2 {
		t.Fatal("regression revocation not counted")
	}
}

func TestEarlyIngestFailuresDoNotPublishPartialSnapshots(t *testing.T) {
	for _, tc := range []struct {
		name   string
		mutate func(*headTestSource, *headTestResolver)
	}{
		{"unsupported head", func(s *headTestSource, _ *headTestResolver) { s.headErr = errors.New("unsupported safe tag") }},
		{"second log chunk", func(s *headTestSource, _ *headTestResolver) { s.failRangeFrom = 106 }},
		{"status at hash", func(s *headTestSource, _ *headTestResolver) { s.statusHashErr = errors.New("hash unavailable") }},
		{"canonical read", func(s *headTestSource, _ *headTestResolver) { s.hashErr = errors.New("RPC unavailable") }},
		{"ingest snapshot changed", func(s *headTestSource, _ *headTestResolver) { s.changedAt = 110 }},
		{"finalized snapshot changed", func(s *headTestSource, _ *headTestResolver) { s.changedAt = 100 }},
		{"Nitro RPC failure", func(_ *headTestSource, r *headTestResolver) { r.err = errors.New("Nitro unavailable") }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a, s, r := newHeadTestLoop(t, arbitrum.L1IngestSafe)
			metrics := arbitrum.NewMetrics(prometheus.NewRegistry())
			a.store.SetMetrics(metrics)
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			before := a.store.ProvisionalAssertions()
			s.events = []arbitrum.ProposedAssertion{earlyProposal(2, 103), earlyProposal(3, 108)}
			tc.mutate(s, r)
			if err := a.SyncOnce(context.Background()); err == nil {
				t.Fatal("expected failure")
			}
			if !reflect.DeepEqual(before, a.store.ProvisionalAssertions()) || len(a.store.Proposals()) != 0 || a.store.NextL1Block() != 101 {
				t.Fatal("failed speculative pass published partial state")
			}
			if metricValue(metrics.AssertionRevocations.WithLabelValues("arb")) != 0 {
				t.Fatal("failed pass counted withdrawals")
			}
		})
	}
}

func TestEarlyIngestPersistenceAndLocalStartupCleanup(t *testing.T) {
	a, s, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	path := filepath.Join(t.TempDir(), "state.json")
	store, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	a.store = store
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	before := a.store.ProvisionalAssertions()
	reloaded, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	a.store = reloaded
	if !reloaded.HasUnfinalizedAssertions() {
		t.Fatal("persisted speculative root not detected")
	}
	s.headErr = errors.New("head unsupported")
	callsBefore := len(s.heads) + len(s.ranges) + len(s.statusBlocks) + len(s.statusHashes)
	if err := a.ReconcileStartup(ctx); err != nil {
		t.Fatalf("startup depended on unavailable RPC: %v", err)
	}
	if got := len(s.heads) + len(s.ranges) + len(s.statusBlocks) + len(s.statusHashes); got != callsBefore {
		t.Fatal("startup made RPC calls")
	}
	cleaned, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil || cleaned.HasUnfinalizedAssertions() || cleaned.NextL1Block() != 101 {
		t.Fatalf("startup cleanup was not durable: %v", err)
	}
	s.headErr = nil
	s.events = nil
	if err := a.ReconcileStartup(ctx); err != nil {
		t.Fatal(err)
	}
	if reloaded.HasUnfinalizedAssertions() {
		t.Fatal("restart retained orphaned root")
	}
	s.events = []arbitrum.ProposedAssertion{earlyProposal(1, 110)}
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	} // Force durable rename failure.
	s.events = []arbitrum.ProposedAssertion{earlyProposal(2, 110)}
	if err := a.SyncOnce(ctx); err == nil {
		t.Fatal("expected persistence failure")
	}
	if !reflect.DeepEqual(before, a.store.ProvisionalAssertions()) {
		t.Fatal("published speculative state despite persistence failure")
	}
}

func TestFinalizedModeDefaultsAndPurgesSpeculativeState(t *testing.T) {
	for _, mode := range []arbitrum.L1IngestHead{"", arbitrum.L1IngestFinalized} {
		t.Run(string(mode), func(t *testing.T) {
			a, s, _ := newHeadTestLoop(t, mode)
			if err := a.SyncOnce(context.Background()); err != nil {
				t.Fatal(err)
			}
			if len(s.heads) != 0 || len(s.statusHashes) != 0 {
				t.Fatal("default queried speculative views")
			}
			if _, found := a.store.HighestAttested(true); found {
				t.Fatal("default ingested unfinalized event")
			}
		})
	}
	a, _, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	if err := a.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	a.config.L1IngestHead = arbitrum.L1IngestFinalized
	if err := a.ReconcileStartup(context.Background()); err != nil {
		t.Fatal(err)
	}
	if a.store.HasUnfinalizedAssertions() {
		t.Fatal("switching to finalized kept speculative roots")
	}
}

func TestSpeculativeCommitPreservesConcurrentDerivedWrites(t *testing.T) {
	a, _, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	var wg sync.WaitGroup
	errs := make(chan error, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		for n := uint64(1); n <= 20; n++ {
			err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
				return true, store.AppendDerived(arbitrum.BlockCommitment{BlockNumber: n, BlockHash: blockHash(n), StateRoot: blockHash(n + 1000)}, time.Unix(1000, 0), true)
			})
			if err != nil {
				errs <- err
				return
			}
		}
	}()
	for n := 0; n < 5; n++ {
		if err := a.SyncOnce(context.Background()); err != nil {
			t.Fatal(err)
		}
	}
	wg.Wait()
	select {
	case err := <-errs:
		t.Fatal(err)
	default:
	}
	for n := uint64(1); n <= 20; n++ {
		root, found := a.store.HighestAttestedAtOrBelow(n, true)
		if !found || root.Source != arbitrum.SourceDerived || root.L2BlockNumber != n {
			t.Fatalf("lost derived root %d", n)
		}
	}
}

func TestSyncFailureCannotHideUncheckedRootBehindAdvancedCursor(t *testing.T) {
	a, s, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	path := filepath.Join(t.TempDir(), "state.json")
	store, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	a.store = store
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	s.finalizedBlock = 110
	s.events = nil // The assertion disappeared before this block finalized.
	s.statusNumberErr = errors.New("finalized status unavailable")
	if err := a.SyncOnce(ctx); err == nil {
		t.Fatal("sync accepted failed finalized recheck")
	}
	reloaded, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.NextL1Block() != 101 || !reloaded.HasUnfinalizedAssertions() {
		t.Fatal("failed startup hid unchecked root behind cursor")
	}
	a.store = reloaded
	s.statusNumberErr = nil
	if err := a.ReconcileStartup(ctx); err != nil {
		t.Fatal(err)
	}
	if _, found := a.store.HighestAttested(true); found {
		t.Fatal("orphan survived authoritative finalized check")
	}
	if a.store.NextL1Block() != 101 {
		t.Fatal("startup cleanup advanced the cursor")
	}
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if a.store.NextL1Block() != 111 {
		t.Fatal("background finalized cursor did not resume")
	}
}

func TestEarlyProposalNotReadyIsPersistedAndDroppedAtStartup(t *testing.T) {
	a, s, r := newHeadTestLoop(t, arbitrum.L1IngestSafe)
	r.err = arbitrum.ErrCommitmentNotReady
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if len(a.store.Proposals()) != 1 || !a.store.HasUnfinalizedAssertions() {
		t.Fatal("unresolved speculative proposal not retained")
	}
	s.events = nil
	if err := a.ReconcileStartup(ctx); err != nil {
		t.Fatal(err)
	}
	if len(a.store.Proposals()) != 0 {
		t.Fatal("orphaned unresolved proposal survived startup")
	}
}

func TestRejectedSpeculativeAssertionIsWithdrawnAndNeverPromoted(t *testing.T) {
	a, _, r := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	r.err = arbitrum.ErrCommitmentMismatch
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if _, found := a.store.HighestAttested(true); found {
		t.Fatal("Nitro-mismatched speculative assertion retained")
	}
	if len(a.store.Mismatches()) != 1 || len(a.store.Proposals()) != 0 {
		t.Fatal("missing terminal mismatch record")
	}
}

func TestSpeculativeSnapshotCannotDemoteExistingFinalizedAssertion(t *testing.T) {
	a, s, r := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	proposal := s.events[0]
	proposal.L1BlockNumber = 90
	proposal.Confirmed = true
	commitment, err := r.ResolveCanonicalBlockHash(context.Background(), proposal.L2BlockHash, arbitrum.RunModeSafe)
	if err != nil {
		t.Fatal(err)
	}
	if err := a.store.RecordAssertionAttestation(proposal, commitment, time.Unix(1000, 0), false); err != nil {
		t.Fatal(err)
	}
	s.events = append(s.events, earlyProposal(2, 109))
	for pass := 0; pass < 2; pass++ {
		if err := a.SyncOnce(context.Background()); err != nil {
			t.Fatalf("speculative replay stalled ingestion: %v", err)
		}
	}
	if len(a.store.ProvisionalAssertions()) != 1 || a.store.ProvisionalAssertions()[0].AssertionHash != s.events[1].AssertionHash {
		t.Fatal("replay prevented unrelated assertion ingestion")
	}
	root, found := a.store.HighestAttested(false)
	if !found || root.Provisional || root.L1BlockNumber != 90 {
		t.Fatalf("finalized assertion changed: %+v", root)
	}
}

func metricValue(metric prometheus.Metric) float64 {
	var out dto.Metric
	if err := metric.Write(&out); err != nil {
		panic(err)
	}
	if out.Gauge != nil {
		return out.GetGauge().GetValue()
	}
	return out.GetCounter().GetValue()
}

func TestCanonicalRPCDoesNotBlockDerivedCommit(t *testing.T) {
	a, s, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	entered := make(chan struct{})
	release := make(chan struct{})
	var once sync.Once
	s.canonicalCheck = func() {
		once.Do(func() { close(entered) })
		<-release
	}
	syncDone := make(chan error, 1)
	go func() { syncDone <- a.SyncOnce(context.Background()) }()
	select {
	case <-entered:
	case <-time.After(5 * time.Second):
		close(release)
		t.Fatal("canonical RPC was not reached")
	}
	commitDone := make(chan error, 1)
	go func() {
		commitDone <- a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
			return true, store.AppendDerived(arbitrum.BlockCommitment{BlockNumber: 1, BlockHash: blockHash(1), StateRoot: blockHash(2)}, time.Unix(1000, 0), true)
		})
	}()
	select {
	case err := <-commitDone:
		close(release)
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		close(release)
		<-syncDone
		<-commitDone
		t.Fatal("canonical RPC held the save lock")
	}
	if err := <-syncDone; err != nil {
		t.Fatal(err)
	}
	if _, found := a.store.DerivedAt(1); !found {
		t.Fatal("speculative commit lost concurrent derived root")
	}
}

func TestFinalizedReconciliationRunsOncePerAssertionPerPass(t *testing.T) {
	for _, initiallyReady := range []bool{false, true} {
		t.Run(fmt.Sprint(initiallyReady), func(t *testing.T) {
			a, s, r := newHeadTestLoop(t, arbitrum.L1IngestLatest)
			if !initiallyReady {
				r.err = arbitrum.ErrCommitmentNotReady
			}
			ctx := context.Background()
			if err := a.SyncOnce(ctx); err != nil {
				t.Fatal(err)
			}
			r.err = nil
			s.finalizedBlock = 110
			s.statusBlocks = nil
			r.modes = nil
			if err := a.SyncOnce(ctx); err != nil {
				t.Fatal(err)
			}
			if len(s.statusBlocks) != 1 || len(r.modes) != 1 {
				t.Fatalf("duplicate reconciliation: status=%v Nitro=%v", s.statusBlocks, r.modes)
			}
			s.confirmed = true
			r.finalized = true
			if err := a.SyncOnce(ctx); err != nil {
				t.Fatal(err)
			}
			if _, found := a.store.HighestAttested(false); !found {
				t.Fatal("deduplication prevented later promotion")
			}
		})
	}
}

func TestStartupPersistenceFailureRetainsUncheckedState(t *testing.T) {
	a, _, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	path := filepath.Join(t.TempDir(), "state.json")
	store, err := arbitrum.LoadAttestedRootStore(path, "arb", 101)
	if err != nil {
		t.Fatal(err)
	}
	a.store = store
	if err := a.SyncOnce(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatal(err)
	}
	if err := a.ReconcileStartup(context.Background()); err == nil {
		t.Fatal("startup accepted non-durable cleanup")
	}
	if !a.store.HasUnfinalizedAssertions() || a.store.NextL1Block() != 101 {
		t.Fatal("failed cleanup published state or advanced cursor")
	}
}

func TestStartupCleanupPreservesFinalizedEvidenceAndDerivedRoots(t *testing.T) {
	a, s, r := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	for _, block := range []uint64{90, 100} {
		proposal := earlyProposal(int64(block), block)
		commitment, err := r.ResolveCanonicalBlockHash(ctx, proposal.L2BlockHash, arbitrum.RunModeSafe)
		if err != nil {
			t.Fatal(err)
		}
		if err := a.store.RecordAssertionAttestation(proposal, commitment, time.Unix(1000, 0), block == 100); err != nil {
			t.Fatal(err)
		}
	}
	if err := a.store.AppendDerived(arbitrum.BlockCommitment{BlockNumber: 600, BlockHash: blockHash(600), StateRoot: blockHash(601)}, time.Unix(1000, 0), true); err != nil {
		t.Fatal(err)
	}
	s.headErr = errors.New("RPC unavailable")
	r.err = errors.New("Nitro unavailable")
	if err := a.ReconcileStartup(ctx); err != nil {
		t.Fatal(err)
	}
	if !a.store.IsFinalizedAssertion(earlyProposal(90, 90).AssertionHash) {
		t.Fatal("startup removed finalized root")
	}
	if roots := a.store.ProvisionalAssertions(); len(roots) != 1 || roots[0].L1BlockNumber != 100 {
		t.Fatalf("startup removed provisional root with finalized L1 evidence: %v", roots)
	}
	if _, found := a.store.DerivedAt(600); !found {
		t.Fatal("startup removed derived root")
	}
	if a.store.NextL1Block() != 101 {
		t.Fatal("startup changed cursor")
	}
}

func TestFinalizedReconciliationKeepsRepostedCreationBlock(t *testing.T) {
	a, s, _ := newHeadTestLoop(t, arbitrum.L1IngestLatest)
	ctx := context.Background()
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	// The same assertion is reposted before the old creation block finalizes.
	s.events[0].L1BlockNumber = 113
	s.finalizedBlock, s.head = 115, 115
	if err := a.SyncOnce(ctx); err != nil {
		t.Fatal(err)
	}
	root, found := a.store.HighestAttested(true)
	if !found || root.L1BlockNumber != 113 || len(s.statusBlocks) != 1 {
		t.Fatalf("reposted root=%+v status reads=%v", root, s.statusBlocks)
	}
}
