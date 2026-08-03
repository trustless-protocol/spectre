package services

import (
	"context"
	"errors"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func TestScanForEVMTimeoutsTrackerRemoval(t *testing.T) {
	expiredPacket := channeltypesv2.Packet{
		SourceClient:      "l2-client-0",
		DestinationClient: "08-wasm-1",
		Sequence:          7,
		TimeoutTimestamp:  uint64(time.Now().Add(-time.Minute).Unix()),
	}

	tests := []struct {
		name             string
		commitment       bool
		commitmentErr    error
		timeoutSucceeded bool
		wantRemaining    bool
	}{
		{
			name:          "commitment check error keeps pending packet",
			commitmentErr: errors.New("rpc unavailable"),
			wantRemaining: true,
		},
		{
			name:             "timeout submit failure keeps pending packet",
			commitment:       true,
			timeoutSucceeded: false,
			wantRemaining:    true,
		},
		{
			name:          "cleared commitment removes pending packet",
			commitment:    false,
			wantRemaining: false,
		},
		{
			name:             "successful timeout removes pending packet",
			commitment:       true,
			timeoutSucceeded: true,
			wantRemaining:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tracker := NewPendingPacketTracker()
			tracker.Add(expiredPacket, 123)

			timeoutCalled := false
			svc := New(nil, nil, DefaultConfig())
			svc.scanForEVMTimeouts(context.Background(), Context{}, evmTimeoutScanOptions{
				tag:     "Test",
				tracker: tracker,
				hasPendingCommitment: func(Context, channeltypesv2.Packet) (bool, error) {
					return tc.commitment, tc.commitmentErr
				},
				timeoutSend: func(_ context.Context, _ Context, packet EthPacket) bool {
					timeoutCalled = true
					if packet.BlockNumber != 123 {
						t.Fatalf("timeout packet block number = %d, want 123", packet.BlockNumber)
					}
					return tc.timeoutSucceeded
				},
			})

			gotRemaining := tracker.Len() == 1
			if gotRemaining != tc.wantRemaining {
				t.Fatalf("tracker remaining = %v, want %v", gotRemaining, tc.wantRemaining)
			}
			wantTimeoutCalled := tc.commitmentErr == nil && tc.commitment
			if timeoutCalled != wantTimeoutCalled {
				t.Fatalf("timeoutCalled = %v, want %v", timeoutCalled, wantTimeoutCalled)
			}
		})
	}
}

func TestTrackL2Pending(t *testing.T) {
	svc := New(nil, nil, DefaultConfig())
	pkt := channeltypesv2.Packet{SourceClient: "l2-client-0", Sequence: 9}

	svc.TrackL2Pending(pkt, 456)
	pending := svc.BatchBuilder.L2PendingTracker.GetAll()
	if len(pending) != 1 {
		t.Fatalf("pending len = %d, want 1", len(pending))
	}
	if pending[0].BlockNumber != 456 {
		t.Fatalf("pending block = %d, want 456", pending[0].BlockNumber)
	}

	svc.UntrackL2Pending(pkt)
	if got := svc.BatchBuilder.L2PendingTracker.Len(); got != 0 {
		t.Fatalf("pending len after untrack = %d, want 0", got)
	}
}
