package subscriber

import (
	"errors"
	"io"
	"log"
	"path/filepath"
	"testing"

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
	wantErr := errors.New("provider limit")
	_, err := sub.scanEthRangeInChunks(deps, "send", &cursor, 10, 4,
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
