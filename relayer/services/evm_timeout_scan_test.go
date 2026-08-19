package services

import (
	"context"
	"errors"
	client "relayer/client"
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
			svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
				tag:     "Test",
				tracker: tracker,
				prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
					return &client.LightBlock{}, true
				},
				hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
					return tc.commitment, tc.commitmentErr
				},
				timeoutSend: func(_ context.Context, _ evmTimeoutDeps, packet EthPacket, _ *client.LightBlock) timeoutSendOutcome {
					timeoutCalled = true
					if packet.BlockNumber != 123 {
						t.Fatalf("timeout packet block number = %d, want 123", packet.BlockNumber)
					}
					if tc.timeoutSucceeded {
						return timeoutSent
					}
					return timeoutDeferred
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

func TestScanForEVMTimeoutsSharedUpdateFailureRunsOnceAndDefersAll(t *testing.T) {
	tracker := NewPendingPacketTracker()
	for _, sequence := range []uint64{1, 2, 3} {
		tracker.Add(channeltypesv2.Packet{
			SourceClient: "l2-client-0", Sequence: sequence,
			TimeoutTimestamp: uint64(time.Now().Add(-time.Minute).Unix()),
		}, sequence)
	}

	prepareCalls := 0
	timeoutCalls := 0
	svc := New(nil, nil, DefaultConfig())
	svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
		tag:     "Test",
		tracker: tracker,
		prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
			prepareCalls++
			return nil, false
		},
		hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
			return true, nil
		},
		timeoutSend: func(context.Context, evmTimeoutDeps, EthPacket, *client.LightBlock) timeoutSendOutcome {
			timeoutCalls++
			return timeoutSent
		},
	})

	if prepareCalls != 1 {
		t.Fatalf("shared client update attempts = %d, want 1", prepareCalls)
	}
	if timeoutCalls != 0 {
		t.Fatalf("timeout submissions after failed shared update = %d, want 0", timeoutCalls)
	}
	for _, info := range tracker.GetAll() {
		if info.Deferrals != 1 || info.TimeoutAttempts != 0 {
			t.Fatalf("seq=%d retry state = deferrals:%d attempts:%d, want 1:0",
				info.Packet.Sequence, info.Deferrals, info.TimeoutAttempts)
		}
	}
}

func TestScanForEVMTimeoutsClearsSharedDeferralsAfterRecovery(t *testing.T) {
	tracker := NewPendingPacketTracker()
	packet := channeltypesv2.Packet{
		SourceClient: "l2-client-0", Sequence: 1,
		TimeoutTimestamp: uint64(time.Now().Add(-time.Minute).Unix()),
	}
	tracker.Add(packet, 123)
	// Make the packet due again with a saturated shared-outage history.
	past := time.Now().Add(-time.Hour)
	for i := 0; i < 6; i++ {
		tracker.deferTimeoutRetryForTest(packet.SourceClient, packet.Sequence, past)
	}

	svc := New(nil, nil, DefaultConfig())
	svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
		tag:     "Test",
		tracker: tracker,
		prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
			return &client.LightBlock{}, true
		},
		hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
			return true, nil
		},
		timeoutSend: func(context.Context, evmTimeoutDeps, EthPacket, *client.LightBlock) timeoutSendOutcome {
			return timeoutDeferred
		},
	})

	got := tracker.GetAll()
	if len(got) != 1 || got[0].Deferrals != 1 {
		t.Fatalf("deferrals after recovered shared prerequisite = %#v, want a fresh single deferral", got)
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
