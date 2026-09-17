package services

import (
	"context"
	"errors"
	updateClientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
	"testing"
	"time"
)

func TestBuildCosmosClientUpdateMsg_RejectsAlreadyCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	worker := &Worker{}
	_, err := worker.BuildCosmosClientUpdateMsg(ctx, CosmosEndpoint{}, EVMEndpoint{}, time.Second, "2/3", "", 0, "2/3", false, 0)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

// TestWaitForCosmosCatchUp_AbortsOnCancelledContext verifies the catch-up wait
// returns promptly when the context is already cancelled, instead of polling the
// Cosmos client (which is never reached — a zero endpoint would panic if it were).
func TestWaitForCosmosCatchUp_AbortsOnCancelledContext(t *testing.T) {
	w := &Worker{}
	stdCtx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	done := make(chan struct{})
	go func() {
		// A zero endpoint set is fine: the cancelled stdCtx short-circuits
		// before ctx.CosmosClient() is ever called.
		w.WaitForCosmosCatchUp(stdCtx, CosmosEndpoint{}, &relayerclient.EthereumClientState{}, 0)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("WaitForCosmosCatchUp did not abort on a cancelled context")
	}
}

// The on-chain height is uint64 and the Tendermint RPC takes int64. The guard is
// the conversion boundary; without it a height above 2^63-1 wraps to a negative
// number, and a negative height reads to the RPC as "latest" rather than as an
// error.
func TestClientStateRevisionHeightInt64(t *testing.T) {
	const maxInt64 = uint64(1)<<63 - 1

	tests := []struct {
		name    string
		height  uint64
		want    int64
		wantErr bool
	}{
		{"zero", 0, 0, false},
		{"ordinary height", 12345, 12345, false},
		{"largest representable", maxInt64, int64(maxInt64), false},
		{"one past int64", maxInt64 + 1, 0, true},
		{"max uint64", ^uint64(0), 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := clientStateRevisionHeightInt64(relayerclient.ClientState{
				LatestHeight: updateClientContract.IICS02ClientMsgsHeight{RevisionHeight: tt.height},
			})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("height %d was accepted; it does not fit in an int64", tt.height)
				}
				return
			}
			if err != nil {
				t.Fatalf("height %d rejected: %v", tt.height, err)
			}
			if got != tt.want {
				t.Fatalf("height = %d, want %d", got, tt.want)
			}
		})
	}
}
