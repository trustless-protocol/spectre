package subscriber

import (
	"context"
	"errors"
	"io"
	"log"
	"path/filepath"
	"testing"

	"relayer/chain"
	relayerclient "relayer/client"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func newRecoveryTestSubscriber(t *testing.T) (*Subscriber, *services.RecoveryStateStore, cosmosDeps) {
	t.Helper()
	store, err := services.LoadRecoveryState(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	deps := cosmosDeps{IDs: services.ClientIDs{CosmosOnEVM: "source-a"}, Logger: log.New(io.Discard, "", 0)}
	return NewSubscriber(store), store, deps
}

func TestRecoveryChunkSize(t *testing.T) {
	t.Setenv(cosmosRecoveryChunkEnv, "17")
	if got := recoveryChunkSize(cosmosRecoveryChunkEnv, defaultCosmosRecoveryChunkHeights); got != 17 {
		t.Fatalf("chunk size = %d, want 17", got)
	}
	t.Setenv(cosmosRecoveryChunkEnv, "0")
	if got := recoveryChunkSize(cosmosRecoveryChunkEnv, defaultCosmosRecoveryChunkHeights); got != defaultCosmosRecoveryChunkHeights {
		t.Fatalf("invalid chunk fallback = %d", got)
	}
}

func TestSubscriberResumesFromPersistedCursors(t *testing.T) {
	sub, store, deps := newRecoveryTestSubscriber(t)
	if err := store.SaveCheckpoint("source-a", services.RecoveryCursors{CosmosHeight: 100, EthSendBlock: 120, EthWriteAckBlk: 900}); err != nil {
		t.Fatal(err)
	}
	if got := sub.resumeCosmosCursor(deps, 1000, 256); got != 100 {
		t.Fatalf("cosmos resume = %d, want 100", got)
	}
	send, ack := sub.resumeEthCursors(deps, 1000, 256)
	if send != 120 || ack != 744 {
		t.Fatalf("eth resume = (%d,%d), want (120,744)", send, ack)
	}
}

func TestPersistCursorsClampToOutstandingPackets(t *testing.T) {
	sub, store, deps := newRecoveryTestSubscriber(t)
	// Seed progress above the outstanding packets. The safety clamp must be able
	// to rewind an already-persisted cursor, not only constrain its first write.
	if err := store.SaveCheckpoint("source-a", services.RecoveryCursors{CosmosHeight: 90, EthSendBlock: 100, EthWriteAckBlk: 100}); err != nil {
		t.Fatal(err)
	}
	bb := services.NewBatchBuilder()
	bb.AddCosmos(services.CosmosPacket{Type: services.CosmosSend, Packet: &channeltypesv2.Packet{Sequence: 1}, BlockNumber: 50})
	bb.AddEth(services.EthPacket{Type: services.EthSend, Packet: &channeltypesv2.Packet{Sequence: 2}, BlockNumber: 60})

	sub.persistCosmosCursor(deps, bb, 100)
	sub.persistEthCursors(deps, bb, 110, 120)
	got, ok := store.Get("source-a")
	if !ok {
		t.Fatal("source cursors not saved")
	}
	if got.CosmosHeight != 50 || got.EthSendBlock != 60 || got.EthWriteAckBlk != 60 {
		t.Fatalf("clamped cursors = %+v", got)
	}
}

func TestScanEthRangePersistsOnlySuccessfulChunks(t *testing.T) {
	sub, store, deps := newRecoveryTestSubscriber(t)
	cursor := uint64(1)
	calls := 0
	persist := func() {
		if err := store.SaveCheckpoint("source-a", services.RecoveryCursors{EthSendBlock: cursor}); err != nil {
			t.Fatal(err)
		}
	}
	// "provider limit" is deliberately NOT a range rejection: it names no range, no
	// block window and no result-set size, so it stays an ordinary failure and the
	// scan stops at the chunk that failed rather than narrowing.
	wantErr := errors.New("provider limit")
	span := &relayerclient.LogSpan{Chunk: 4}
	_, err := sub.scanEthRangeInChunks(context.Background(), deps, "send", &cursor, 10, span,
		func(from, to uint64) (ethRecoveryStats, error) {
			calls++
			if calls == 2 {
				return ethRecoveryStats{}, wantErr
			}
			return ethRecoveryStats{recovered: 1}, nil
		}, persist)
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want %v", err, wantErr)
	}
	if cursor != 5 {
		t.Fatalf("cursor = %d, want first failed chunk start 5", cursor)
	}
	got, _ := store.Get("source-a")
	if got.EthSendBlock != 5 {
		t.Fatalf("persisted cursor = %d, want 5", got.EthSendBlock)
	}
}

// The ETH mirror of the L2 narrowing branch, and the one place the two policies
// differ on purpose: this side persists a cursor after every piece that
// succeeded, so a narrowing must resume at that cursor rather than re-scan the
// range. Re-scanning would re-offer events already handed to the batch builder.
func TestScanEthRangeNarrowsAndResumesAtTheCursor(t *testing.T) {
	sub, store, deps := newRecoveryTestSubscriber(t)
	cursor := uint64(1)
	span := &relayerclient.LogSpan{Chunk: 8}
	persist := func() {
		if err := store.SaveCheckpoint("source-a", services.RecoveryCursors{EthSendBlock: cursor}); err != nil {
			t.Fatal(err)
		}
	}

	const cap = uint64(2)
	var asked [][2]uint64
	_, err := sub.scanEthRangeInChunks(context.Background(), deps, "send", &cursor, 16, span,
		func(from, to uint64) (ethRecoveryStats, error) {
			asked = append(asked, [2]uint64{from, to})
			// Bounded on purpose: a scan that does not narrow asks for the same
			// refused width forever, and an unbounded test reports that as a
			// wall-clock timeout rather than as the defect it is.
			if len(asked) > 64 {
				t.Fatalf("the scan asked 64 times without settling: %v", asked)
			}
			if to-from+1 > cap {
				return ethRecoveryStats{}, errors.New("query returned more than 10000 results")
			}
			return ethRecoveryStats{recovered: 1}, nil
		}, persist)
	if err != nil {
		t.Fatalf("scan never found a servable span: %v", err)
	}
	if span.Chunk > cap {
		t.Fatalf("span settled at %d, want <= %d", span.Chunk, cap)
	}
	if cursor != 17 {
		t.Fatalf("cursor = %d, want 17 (the range was not finished)", cursor)
	}
	// Every refused attempt must have started at the SAME block: narrowing changes
	// the width, never the position, and the cursor never moves on a failure.
	if asked[0][0] != 1 {
		t.Fatalf("first attempt started at %d, want 1", asked[0][0])
	}
	for i := 1; i < len(asked); i++ {
		prev, cur := asked[i-1], asked[i]
		if prev[1]-prev[0]+1 > cap && cur[0] != prev[0] {
			t.Fatalf("after refusing [%d,%d] the scan moved to [%d,%d]; a failed piece must be retried in place, not skipped",
				prev[0], prev[1], cur[0], cur[1])
		}
	}
	// The pieces that succeeded must tile the range with no gap — a gap here is a
	// silently missed packet, which is what this recovery scan exists to prevent.
	var next uint64 = 1
	for _, a := range asked {
		if a[1]-a[0]+1 > cap {
			continue // refused, contributed nothing
		}
		if a[0] != next {
			t.Fatalf("gap or overlap: expected a piece starting at %d, got [%d,%d]", next, a[0], a[1])
		}
		next = a[1] + 1
	}
	if next != 17 {
		t.Fatalf("scanned pieces ended at %d, want 17", next)
	}
}

// The bottom of the ladder on the ETH side. A one-block span still refused is not
// a sizing problem; halving further is how the recovery scan would spin forever.
func TestScanEthRangeStopsWhenTheSpanCannotShrinkFurther(t *testing.T) {
	sub, _, deps := newRecoveryTestSubscriber(t)
	cursor := uint64(1)
	span := &relayerclient.LogSpan{Chunk: 4}
	calls := 0
	_, err := sub.scanEthRangeInChunks(context.Background(), deps, "send", &cursor, 8, span,
		func(from, to uint64) (ethRecoveryStats, error) {
			calls++
			if calls > 32 {
				t.Fatal("the scan is halving forever instead of bottoming out")
			}
			return ethRecoveryStats{}, errors.New("block range too large")
		}, func() {})
	if err == nil {
		t.Fatal("a provider refusing every span was reported as a successful scan")
	}
	if !chain.IsPermanent(err) {
		t.Fatalf("error = %v, want permanent at the bottom of the ladder", err)
	}
	if cursor != 1 {
		t.Fatalf("cursor = %d, want 1 — it must never advance past a piece that failed", cursor)
	}
}
