package opstack

import (
	"context"
	"errors"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

type fakeGames struct {
	games    []ProposedRoot
	errAt    map[uint64]error
	countErr error
}

func (f *fakeGames) GameCount(_ context.Context) (uint64, error) {
	if f.countErr != nil {
		return 0, f.countErr
	}
	return uint64(len(f.games)), nil
}

func (f *fakeGames) GameAtIndex(_ context.Context, idx uint64) (ProposedRoot, error) {
	if err, ok := f.errAt[idx]; ok {
		return ProposedRoot{}, err
	}
	if idx >= uint64(len(f.games)) {
		return ProposedRoot{}, errors.New("index out of range")
	}
	return f.games[idx], nil
}

type fakeReplica struct {
	status      SyncStatus
	statusErr   error
	outputs     map[uint64][32]byte
	outputErr   map[uint64]error
	commitments map[uint64]L2Commitment
}

func (f *fakeReplica) SyncStatus(_ context.Context) (SyncStatus, error) {
	if f.statusErr != nil {
		return SyncStatus{}, f.statusErr
	}
	return f.status, nil
}

func (f *fakeReplica) OutputAtBlock(_ context.Context, l2Block uint64) ([32]byte, error) {
	if err, ok := f.outputErr[l2Block]; ok {
		return [32]byte{}, err
	}
	root, ok := f.outputs[l2Block]
	if !ok {
		return [32]byte{}, errors.New("no output for block")
	}
	return root, nil
}

// commitments overrides what CommitmentAt returns; when a block is absent it
// falls back to the output-root map so existing tests need no new fixtures.
func (f *fakeReplica) CommitmentAt(ctx context.Context, l2Block uint64) (L2Commitment, error) {
	if c, ok := f.commitments[l2Block]; ok {
		return c, nil
	}
	root, err := f.OutputAtBlock(ctx, l2Block)
	if err != nil {
		return L2Commitment{}, err
	}
	return L2Commitment{BlockNumber: l2Block, OutputRoot: root}, nil
}

type fakeHook struct {
	calls []uint64 // game indices
	err   error
}

func (f *fakeHook) OnMismatch(_ context.Context, game ProposedRoot, _ [32]byte) error {
	f.calls = append(f.calls, game.GameIndex)
	return f.err
}

func root(b byte) [32]byte {
	var r [32]byte
	r[0] = b
	return r
}

func game(idx uint64, gameType uint32, l2Block uint64, claim [32]byte) ProposedRoot {
	return ProposedRoot{
		GameIndex:     idx,
		GameAddress:   common.BytesToAddress([]byte{byte(idx + 1)}),
		GameType:      gameType,
		RootClaim:     claim,
		L2BlockNumber: l2Block,
		L1Timestamp:   1000 + idx,
	}
}

func newTestAttestor(t *testing.T, head Head, games *fakeGames, rep *fakeReplica, hook ChallengeHook) *OpStackAttestor {
	t.Helper()
	store, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	store.Bootstrap(0)
	if hook == nil {
		hook = &fakeHook{}
	}
	cfg := Config{
		SrcChain:           "op-test",
		L1RpcUrl:           "http://localhost:1",
		OpNodeRpcUrl:       "http://localhost:2",
		DisputeGameFactory: common.BytesToAddress([]byte{0xfa}),
		RespectedGameType:  0,
		AttestationHead:    head,
		StatePath:          store.path,
		// Verdict-focused tests disable self-derived attestation; derived
		// tests flip this off explicitly.
		DisableDerivedRoots: true,
	}
	cfg.applyDefaults()
	a := New(cfg, games, rep, store, hook, nil, zap.NewNop().Sugar())
	a.now = func() time.Time { return time.Unix(5000, 0) }
	return a
}

func TestRunOnceMatchAttests(t *testing.T) {
	claim := root(0xaa)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 150, UnsafeL2: 160, FinalizedL2: 150}, outputs: map[uint64][32]byte{100: claim}}
	hook := &fakeHook{}
	a := newTestAttestor(t, HeadSafe, games, rep, hook)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if got := a.store.Pending(); len(got) != 0 {
		t.Fatalf("pending = %d entries, want 0", len(got))
	}
	cur, ok := a.AttestedUpTo()
	if !ok || cur.Height != 100 || cur.ID != claim {
		t.Fatalf("AttestedUpTo = (%+v, %t), want height 100 root %x", cur, ok, claim)
	}
	if len(hook.calls) != 0 {
		t.Fatalf("hook called %d times on a match, want 0", len(hook.calls))
	}
}

func TestRunOnceMismatchRecordsAndHooks(t *testing.T) {
	claim := root(0xaa)
	derived := root(0xab) // single-byte divergence
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 150, FinalizedL2: 150}, outputs: map[uint64][32]byte{100: derived}}
	hook := &fakeHook{}
	a := newTestAttestor(t, HeadSafe, games, rep, hook)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("mismatch must not enter the attested feed")
	}
	if len(a.store.Pending()) != 0 {
		t.Fatal("mismatch is a terminal verdict; game must leave pending")
	}
	if len(hook.calls) != 1 || hook.calls[0] != 0 {
		t.Fatalf("hook calls = %v, want exactly [0]", hook.calls)
	}
	mismatches := a.store.state.Mismatches
	if len(mismatches) != 1 || mismatches[0].LocalRoot != derived {
		t.Fatalf("mismatch record = %+v, want one entry with local root %x", mismatches, derived)
	}
}

func TestMismatchIsRetainedWhenChallengeHookFails(t *testing.T) {
	claim := root(0xaa)
	local := root(0xab)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	replica := &fakeReplica{
		status:  SyncStatus{SafeL2: 100, FinalizedL2: 100},
		outputs: map[uint64][32]byte{100: local},
	}
	hook := &fakeHook{err: errors.New("challenge transaction rejected")}
	a := newTestAttestor(t, HeadSafe, games, replica, hook)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: mismatch hook errors must not erase verdict: %v", err)
	}
	if got := a.store.state.Mismatches; len(got) != 1 || got[0].LocalRoot != local {
		t.Fatalf("mismatches after hook failure = %+v", got)
	}
	if pending := a.store.Pending(); len(pending) != 0 {
		t.Fatalf("pending after hook failure = %+v", pending)
	}
}

func TestIngestErrorStopsCursor(t *testing.T) {
	claim := root(0x01)
	games := &fakeGames{
		games: []ProposedRoot{game(0, 0, 100, claim), game(1, 0, 200, claim), game(2, 0, 300, claim)},
		errAt: map[uint64]error{1: errors.New("rpc down")},
	}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 1000, FinalizedL2: 1000}, outputs: map[uint64][32]byte{100: claim}}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err == nil {
		t.Fatal("runOnce must surface the ingest error")
	}
	if got := a.store.NextGameIndex(); got != 1 {
		t.Fatalf("ingest cursor = %d, want 1 (stopped at the failed index)", got)
	}
	if got := a.store.Pending(); len(got) != 1 || got[0].GameIndex != 0 {
		t.Fatalf("pending = %+v, want only game 0", got)
	}

	// Recovery: the error clears, the scan resumes exactly at index 1.
	delete(games.errAt, 1)
	rep.outputs[200] = claim
	rep.outputs[300] = claim
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce after recovery: %v", err)
	}
	if got := a.store.NextGameIndex(); got != 3 {
		t.Fatalf("ingest cursor = %d after recovery, want 3", got)
	}
	if cur, ok := a.AttestedUpTo(); !ok || cur.Height != 300 {
		t.Fatalf("AttestedUpTo = (%+v, %t), want height 300", cur, ok)
	}
}

func TestGameCountErrorStopsBeforeAnyStateChange(t *testing.T) {
	games := &fakeGames{countErr: errors.New("factory unavailable")}
	replica := &fakeReplica{}
	a := newTestAttestor(t, HeadFinalized, games, replica, nil)

	if err := a.runOnce(context.Background()); err == nil || !strings.Contains(err.Error(), "factory unavailable") {
		t.Fatalf("runOnce error = %v, want GameCount failure", err)
	}
	if got := a.store.NextGameIndex(); got != 0 {
		t.Fatalf("ingest cursor = %d, want 0", got)
	}
	if pending := a.store.Pending(); len(pending) != 0 {
		t.Fatalf("pending after GameCount failure = %+v", pending)
	}
}

func TestReadOnlyReplicaFailureMatrix(t *testing.T) {
	store, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	withoutReplica := New(Config{SrcChain: "op-test"}, nil, nil, store, nil, nil, zap.NewNop().Sugar())
	if _, err := withoutReplica.HeadAt(context.Background(), HeadSafe); !errors.Is(err, ErrNoReplica) {
		t.Fatalf("HeadAt without replica error = %v, want ErrNoReplica", err)
	}
	if _, err := withoutReplica.CommitmentAt(context.Background(), 100); !errors.Is(err, ErrNoReplica) {
		t.Fatalf("CommitmentAt without replica error = %v, want ErrNoReplica", err)
	}

	replica := &fakeReplica{statusErr: errors.New("sync status unavailable")}
	withReplica := New(Config{SrcChain: "op-test"}, nil, replica, store, nil, nil, zap.NewNop().Sugar())
	if _, err := withReplica.HeadAt(context.Background(), HeadSafe); err == nil || !strings.Contains(err.Error(), "sync status unavailable") {
		t.Fatalf("HeadAt sync error = %v", err)
	}
	replica.statusErr = nil
	replica.status = SyncStatus{SafeL2: 100}
	if _, err := withReplica.HeadAt(context.Background(), Head("unknown")); err == nil || !strings.Contains(err.Error(), "unknown attestation head") {
		t.Fatalf("HeadAt invalid head error = %v", err)
	}
	replica.outputErr = map[uint64]error{100: errors.New("commitment unavailable")}
	if _, err := withReplica.CommitmentAt(context.Background(), 100); err == nil || !strings.Contains(err.Error(), "commitment unavailable") {
		t.Fatalf("CommitmentAt error = %v", err)
	}
}

func TestReplicaLagKeepsPending(t *testing.T) {
	claim := root(0x02)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 500, claim)}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 400, UnsafeL2: 450, FinalizedL2: 400}} // safe head behind the proposal
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("a lagging replica is not an error: %v", err)
	}
	if got := a.store.Pending(); len(got) != 1 {
		t.Fatalf("pending = %d entries, want 1 (undecided while safe head lags)", len(got))
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("nothing may be affirmed while the safe head lags")
	}
}

func TestReplicaErrorKeepsPendingAndFails(t *testing.T) {
	claim := root(0x03)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	rep := &fakeReplica{
		status:    SyncStatus{SafeL2: 150, FinalizedL2: 150},
		outputErr: map[uint64]error{100: errors.New("op-node down")},
	}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err == nil {
		t.Fatal("replica output error must fail the pass")
	}
	if got := a.store.Pending(); len(got) != 1 {
		t.Fatalf("pending = %d entries, want 1 (no verdict on error)", len(got))
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("nothing may be affirmed on a replica error")
	}
}

func TestSyncStatusErrorDecidesNothing(t *testing.T) {
	claim := root(0x04)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	rep := &fakeReplica{statusErr: errors.New("sync status down")}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err == nil {
		t.Fatal("sync status error must fail the pass")
	}
	if got := a.store.Pending(); len(got) != 1 {
		t.Fatalf("pending = %d entries, want 1", len(got))
	}
}

func TestUnrespectedGameTypeSkipped(t *testing.T) {
	claim := root(0x05)
	games := &fakeGames{games: []ProposedRoot{
		game(0, 42, 100, claim), // wrong type
		game(1, 0, 200, claim),
	}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 1000, FinalizedL2: 1000}, outputs: map[uint64][32]byte{200: claim}}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if got := a.store.NextGameIndex(); got != 2 {
		t.Fatalf("cursor = %d, want 2 (skipped game still advances ingest)", got)
	}
	cur, ok := a.AttestedUpTo()
	if !ok || cur.Height != 200 {
		t.Fatalf("AttestedUpTo = (%+v, %t), want height 200 only", cur, ok)
	}
}

func TestUnsafeProvisionalThenConfirmed(t *testing.T) {
	claim := root(0x06)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 500, claim)}}
	// Unsafe head covers the block, finalized head does not: provisional verdict.
	rep := &fakeReplica{status: SyncStatus{SafeL2: 400, UnsafeL2: 600, FinalizedL2: 400}, outputs: map[uint64][32]byte{500: claim}}
	a := newTestAttestor(t, HeadUnsafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("provisional roots must not count as attested by default")
	}
	if cur, ok := a.store.AttestedUpTo(true); !ok || cur.Height != 500 {
		t.Fatalf("includeProvisional AttestedUpTo = (%+v, %t), want height 500", cur, ok)
	}
	if got := a.store.RecheckList(); len(got) != 1 || !got[0].ProvisionalMatch {
		t.Fatalf("recheck list = %+v, want one provisional match", got)
	}

	// Finalized head catches up and agrees: provisional flag clears.
	rep.status.SafeL2 = 600
	rep.status.FinalizedL2 = 600
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (recheck): %v", err)
	}
	if cur, ok := a.AttestedUpTo(); !ok || cur.Height != 500 {
		t.Fatalf("AttestedUpTo after confirmation = (%+v, %t), want height 500", cur, ok)
	}
	if got := a.store.RecheckList(); len(got) != 0 {
		t.Fatalf("recheck list not drained: %+v", got)
	}
}

func TestUnsafeProvisionalMatchDivergesAtSafe(t *testing.T) {
	claim := root(0x07)
	divergent := root(0x08)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 500, claim)}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 400, UnsafeL2: 600, FinalizedL2: 400}, outputs: map[uint64][32]byte{500: claim}}
	hook := &fakeHook{}
	a := newTestAttestor(t, HeadUnsafe, games, rep, hook)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	// Finalized derivation disagrees with what the sequencer gossiped.
	rep.outputs[500] = divergent
	rep.status.SafeL2 = 600
	rep.status.FinalizedL2 = 600
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (recheck): %v", err)
	}
	if _, ok := a.store.AttestedUpTo(true); ok {
		t.Fatal("diverged provisional match must be removed from the feed")
	}
	if len(a.store.state.Mismatches) != 1 {
		t.Fatalf("mismatches = %+v, want the corrected safe-head verdict", a.store.state.Mismatches)
	}
	if len(hook.calls) != 1 {
		t.Fatalf("hook calls = %v, want exactly one (from the authoritative safe verdict)", hook.calls)
	}
	if got := a.store.RecheckList(); len(got) != 0 {
		t.Fatalf("recheck list not drained: %+v", got)
	}
}

func TestUnsafeProvisionalMismatchConfirmedLegitAtSafe(t *testing.T) {
	claim := root(0x09)
	gossiped := root(0x0a)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 500, claim)}}
	// Gossip disagrees with the claim → provisional mismatch...
	rep := &fakeReplica{status: SyncStatus{SafeL2: 400, UnsafeL2: 600, FinalizedL2: 400}, outputs: map[uint64][32]byte{500: gossiped}}
	hook := &fakeHook{}
	a := newTestAttestor(t, HeadUnsafe, games, rep, hook)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if len(hook.calls) != 1 {
		t.Fatalf("provisional mismatch should raise the hook once, got %v", hook.calls)
	}
	// ...but honest finalized derivation agrees with the claim: proposer was
	// honest, the sequencer gossip diverged.
	rep.outputs[500] = claim
	rep.status.SafeL2 = 600
	rep.status.FinalizedL2 = 600
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (recheck): %v", err)
	}
	cur, ok := a.AttestedUpTo()
	if !ok || cur.Height != 500 || cur.ID != claim {
		t.Fatalf("AttestedUpTo = (%+v, %t), want confirmed height 500", cur, ok)
	}
	if len(a.store.state.Mismatches) != 0 {
		t.Fatalf("provisional mismatch must be superseded, got %+v", a.store.state.Mismatches)
	}
	if len(hook.calls) != 1 {
		t.Fatalf("no additional hook call on a safe-head match, got %v", hook.calls)
	}
}

func TestSafeModeNeverPopulatesRecheck(t *testing.T) {
	claim := root(0x0b)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 100, claim)}}
	rep := &fakeReplica{status: SyncStatus{SafeL2: 150, UnsafeL2: 200, FinalizedL2: 150}, outputs: map[uint64][32]byte{100: claim}}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if got := a.store.RecheckList(); len(got) != 0 {
		t.Fatalf("safe mode populated the recheck list: %+v", got)
	}
	feed := a.store.state.Attested
	if len(feed) != 1 || feed[0].Provisional {
		t.Fatalf("feed = %+v, want one non-provisional entry", feed)
	}
}

func TestBootstrapStartIndex(t *testing.T) {
	games := &fakeGames{games: []ProposedRoot{
		{GameIndex: 0, L1Timestamp: 100},
		{GameIndex: 1, L1Timestamp: 200},
		{GameIndex: 2, L1Timestamp: 300},
		{GameIndex: 3, L1Timestamp: 400},
	}}
	ctx := context.Background()
	cases := []struct {
		cutoff uint64
		want   uint64
	}{
		{cutoff: 0, want: 0},
		{cutoff: 150, want: 1},
		{cutoff: 300, want: 2},
		{cutoff: 999, want: 4}, // everything older → start at the live edge
	}
	for _, tc := range cases {
		got, err := bootstrapStartIndex(ctx, games, 4, tc.cutoff)
		if err != nil {
			t.Fatalf("bootstrapStartIndex(cutoff=%d): %v", tc.cutoff, err)
		}
		if got != tc.want {
			t.Errorf("bootstrapStartIndex(cutoff=%d) = %d, want %d", tc.cutoff, got, tc.want)
		}
	}
	if got, err := bootstrapStartIndex(ctx, games, 0, 100); err != nil || got != 0 {
		t.Errorf("empty factory: got (%d, %v), want (0, nil)", got, err)
	}
}

func TestSafeVerdictAboveFinalizedIsProvisional(t *testing.T) {
	claim := root(0x0c)
	games := &fakeGames{games: []ProposedRoot{game(0, 0, 500, claim)}}
	// Safe covers the block, finalized does not: verdict must be provisional
	// (the safe head can still reorg with L1).
	rep := &fakeReplica{status: SyncStatus{SafeL2: 600, UnsafeL2: 700, FinalizedL2: 400}, outputs: map[uint64][32]byte{500: claim}}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("safe verdict above finalized must not count as confirmed")
	}
	if cur, ok := a.store.AttestedUpTo(true); !ok || cur.Height != 500 {
		t.Fatalf("provisional AttestedUpTo = (%+v, %t), want height 500", cur, ok)
	}
	if got := a.store.RecheckList(); len(got) != 1 {
		t.Fatalf("recheck list = %+v, want one entry", got)
	}

	// Finalized catches up and agrees.
	rep.status.FinalizedL2 = 600
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (recheck): %v", err)
	}
	if cur, ok := a.AttestedUpTo(); !ok || cur.Height != 500 {
		t.Fatalf("AttestedUpTo after confirmation = (%+v, %t), want height 500", cur, ok)
	}
}

func TestDerivedAttestationConfirmedAtFinalized(t *testing.T) {
	rootA, rootB := root(0x20), root(0x21)
	games := &fakeGames{}
	rep := &fakeReplica{
		status:  SyncStatus{SafeL2: 1000, UnsafeL2: 1000, FinalizedL2: 1000},
		outputs: map[uint64][32]byte{1000: rootA, 1200: rootB},
	}
	a := newTestAttestor(t, HeadFinalized, games, rep, nil)
	a.cfg.DisableDerivedRoots = false

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	cur, ok := a.AttestedUpTo()
	if !ok || cur.Height != 1000 || cur.ID != rootA {
		t.Fatalf("AttestedUpTo = (%+v, %t), want confirmed derived root at 1000", cur, ok)
	}

	// Head advances less than the gap: no new derived root.
	rep.status.FinalizedL2 = 1000 + a.cfg.DerivedGapBlocks - 1
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (below gap): %v", err)
	}
	if cur, _ := a.AttestedUpTo(); cur.Height != 1000 {
		t.Fatalf("derived attested below the gap: %+v", cur)
	}

	// Head advances past the gap: new derived root.
	rep.status.FinalizedL2 = 1200
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (past gap): %v", err)
	}
	if cur, _ := a.AttestedUpTo(); cur.Height != 1200 || cur.ID != rootB {
		t.Fatalf("AttestedUpTo = %+v, want derived root at 1200", cur)
	}
}

func TestDerivedProvisionalReverify(t *testing.T) {
	rootA := root(0x22)
	games := &fakeGames{}
	// Safe head ahead of finalized: derived attestation is provisional.
	rep := &fakeReplica{
		status:  SyncStatus{SafeL2: 1000, UnsafeL2: 1100, FinalizedL2: 800},
		outputs: map[uint64][32]byte{1000: rootA},
	}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)
	a.cfg.DisableDerivedRoots = false

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	if _, ok := a.AttestedUpTo(); ok {
		t.Fatal("provisional derived root must not count as confirmed")
	}
	if cur, ok := a.store.AttestedUpTo(true); !ok || cur.Height != 1000 {
		t.Fatalf("provisional AttestedUpTo = (%+v, %t), want 1000", cur, ok)
	}

	// Finalized covers the block, derivation agrees: confirmed in place.
	rep.status.FinalizedL2 = 1000
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (reverify): %v", err)
	}
	if cur, ok := a.AttestedUpTo(); !ok || cur.Height != 1000 || cur.ID != rootA {
		t.Fatalf("AttestedUpTo after reverify = (%+v, %t), want confirmed 1000", cur, ok)
	}
}

func TestDerivedDivergenceCorrected(t *testing.T) {
	provisionalRoot, finalRoot := root(0x23), root(0x24)
	games := &fakeGames{}
	rep := &fakeReplica{
		status:  SyncStatus{SafeL2: 1000, UnsafeL2: 1100, FinalizedL2: 800},
		outputs: map[uint64][32]byte{1000: provisionalRoot},
	}
	a := newTestAttestor(t, HeadSafe, games, rep, nil)
	a.cfg.DisableDerivedRoots = false

	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce: %v", err)
	}
	// L1 reorg scenario: the finalized chain derives a different root.
	rep.outputs[1000] = finalRoot
	rep.status.FinalizedL2 = 1000
	if err := a.runOnce(context.Background()); err != nil {
		t.Fatalf("runOnce (reverify): %v", err)
	}
	cur, ok := a.AttestedUpTo()
	if !ok || cur.Height != 1000 || cur.ID != finalRoot {
		t.Fatalf("AttestedUpTo = (%+v, %t), want corrected finalized root %x", cur, ok, finalRoot)
	}
	if got := a.store.DerivedProvisionalAtOrBelow(1000); len(got) != 0 {
		t.Fatalf("provisional derived entries not drained: %+v", got)
	}
}

func TestPruneDerivedKeepsGamesAndProvisional(t *testing.T) {
	s, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(4000, 0)
	s.RecordMatch(testGame(0, 50, root(0x30)), now, false) // game entry, never pruned
	s.AppendDerived(100, root(0x31), now, false)
	s.AppendDerived(200, root(0x32), now, false)
	s.AppendDerived(300, root(0x33), now, false)
	s.AppendDerived(400, root(0x34), now, true) // provisional, never pruned

	if removed := s.PruneDerived(2); removed != 1 {
		t.Fatalf("PruneDerived removed %d, want 1", removed)
	}
	// The derived entry at 100 is gone; the query falls back to the game
	// entry at 50, which pruning never touches.
	if got, ok := s.HighestAttestedAtOrBelow(100, false); !ok || got.L2BlockNumber != 50 || got.Source != SourceGame {
		t.Fatalf("at 100 after prune: (%+v, %t), want the game entry at 50", got, ok)
	}
	if got, ok := s.HighestAttestedAtOrBelow(250, false); !ok || got.L2BlockNumber != 200 {
		t.Fatalf("at 250 after prune: (%+v, %t), want the surviving derived entry at 200", got, ok)
	}
	if cur, ok := s.AttestedUpTo(true); !ok || cur.Height != 400 {
		t.Fatalf("provisional derived entry must survive pruning, got (%+v, %t)", cur, ok)
	}
	if removed := s.PruneDerived(2); removed != 0 {
		t.Fatalf("second prune removed %d, want 0", removed)
	}
}
