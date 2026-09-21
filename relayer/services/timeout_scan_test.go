package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	client "relayer/client"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func okUpdate() (*EthClientUpdateResult, error) {
	return &EthClientUpdateResult{EthClientState: &client.EthereumClientState{LatestSlot: 1}}, nil
}

// prereqServices builds a Services holding n tracked packets, and returns the
// pendingPacketInfo set the scan would be carrying at the point the shared
// prerequisites are resolved.
func prereqServices(t *testing.T, seqs ...uint64) (*Services, []pendingPacketInfo) {
	t.Helper()
	s := &Services{BatchBuilder: NewBatchBuilder()}
	expired := make([]pendingPacketInfo, 0, len(seqs))
	for _, seq := range seqs {
		pkt := channeltypesv2.Packet{
			Sequence:          seq,
			SourceClient:      "08-wasm-0",
			DestinationClient: "client-0",
			TimeoutTimestamp:  1,
		}
		s.BatchBuilder.PendingTracker.Add(pkt, 100)
		expired = append(expired, pendingPacketInfo{Packet: pkt, BlockNumber: 100})
	}
	return s, expired
}

type transientTimeoutBatchHandler struct{ TransactionHandler }

type timeoutBatchHandler struct {
	TransactionHandler
	singletons [][]any
}

func timeoutRetryPacket(sequence uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{SourceClient: "client-0", Sequence: sequence, TimeoutTimestamp: 1}
}

func TestResolveCosmosTimeoutProofState(t *testing.T) {
	t.Run("defers shared failures", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			buildUpdate      func() (*EthClientUpdateResult, error)
			fetchClientState func() (*client.EthereumClientState, error)
		}{
			{
				name:        "update build",
				buildUpdate: func() (*EthClientUpdateResult, error) { return nil, errors.New("beacon unavailable") },
				fetchClientState: func() (*client.EthereumClientState, error) {
					t.Fatal("client state fetch must not run")
					return nil, nil
				},
			},
			{
				name:             "client state fetch",
				buildUpdate:      func() (*EthClientUpdateResult, error) { return &EthClientUpdateResult{}, nil },
				fetchClientState: func() (*client.EthereumClientState, error) { return nil, errors.New("cosmos unavailable") },
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				tracker := NewPendingPacketTracker()
				expired := make([]pendingPacketInfo, 0, 2)
				for _, sequence := range []uint64{1, 2} {
					packet := timeoutRetryPacket(sequence)
					tracker.Add(packet, sequence)
					expired = append(expired, pendingPacketInfo{Packet: packet, BlockNumber: sequence})
				}

				svc := New(nil, nil, DefaultConfig())
				_, _, ok := svc.resolveCosmosTimeoutProofState(tracker, expired, tc.buildUpdate, tc.fetchClientState)
				if ok {
					t.Fatal("failed shared prerequisite reported success")
				}
				for _, info := range tracker.GetAll() {
					if info.Deferrals != 1 || info.TimeoutAttempts != 0 {
						t.Fatalf("seq=%d retry state = deferrals:%d attempts:%d, want 1:0",
							info.Packet.Sequence, info.Deferrals, info.TimeoutAttempts)
					}
				}
			})
		}
	})

	// TestResolveCosmosTimeoutProofState/prerequisite failures defer without charging is the P1 from
	// review, raised twice: both shared-prerequisite exits returned bare.
	//
	// Neither failure is attributable to any packet -- they fail identically for the
	// whole scan -- so nothing may be charged. But NOT deferring is what made the
	// bug expensive and invisible at once: the next 30s scan rebuilds the client
	// state, finality and update construction from scratch forever, and because
	// Deferrals stays at zero the STUCK line and the stuck gauge cannot report it.
	t.Run("prerequisite failures defer without charging", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name             string
			buildUpdate      func() (*EthClientUpdateResult, error)
			fetchClientState func() (*client.EthereumClientState, error)
		}{
			{
				name:        "BuildEthClientUpdateMsgs fails",
				buildUpdate: func() (*EthClientUpdateResult, error) { return nil, fmt.Errorf("beacon api unavailable") },
				fetchClientState: func() (*client.EthereumClientState, error) {
					t.Fatal("client state must not be fetched after the update build failed")
					return nil, nil
				},
			},
			{
				// updateResult with a nil EthClientState sends the scan to the
				// fallback fetch -- the second exit, which survived the first fix.
				name:             "GetEthereumClientState fallback fails",
				buildUpdate:      func() (*EthClientUpdateResult, error) { return &EthClientUpdateResult{}, nil },
				fetchClientState: func() (*client.EthereumClientState, error) { return nil, fmt.Errorf("cosmos rpc down") },
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				s, expired := prereqServices(t, 11, 12)
				tracker := s.BatchBuilder.PendingTracker

				_, _, ok := s.resolveCosmosTimeoutProofState(tracker, expired, tc.buildUpdate, tc.fetchClientState)
				if ok {
					t.Fatal("a failed shared prerequisite must abort the scan")
				}

				// Deferred: nothing may be due again immediately, or the next 30s
				// scan rebuilds the client state, finality and update construction.
				if due := tracker.GetDue(time.Now()); len(due) != 0 {
					t.Fatalf("%d packet(s) still due immediately; the next scan would rebuild everything", len(due))
				}

				// Not charged: a shared failure is not evidence about any packet.
				if got := tracker.DeadLetteredTimeouts(); got != 0 {
					t.Fatalf("DeadLetteredTimeouts() = %d; a shared prerequisite failure must not charge packets", got)
				}
				if tracker.Len() != len(expired) {
					t.Fatalf("tracker holds %d packet(s), want %d; nothing was delivered", tracker.Len(), len(expired))
				}
			})
		}
	})

	// TestResolveCosmosTimeoutProofState/repeated failures become visible: deferrals never
	// dead-letter, so the stuck gauge is the only thing that can surface a
	// prerequisite that never recovers. Without the deferral it reads zero forever.
	t.Run("repeated failures become visible", func(t *testing.T) {
		t.Parallel()

		s, expired := prereqServices(t, 31)
		failing := func() (*EthClientUpdateResult, error) { return nil, fmt.Errorf("beacon api unavailable") }
		unreachable := func() (*client.EthereumClientState, error) {
			t.Fatal("unreachable")
			return nil, nil
		}

		for range deferralAlertThreshold {
			s.resolveCosmosTimeoutProofState(s.BatchBuilder.PendingTracker, expired, failing, unreachable)
		}

		count, worst := s.BatchBuilder.PendingTracker.StuckTimeouts()
		if count != 1 || worst != deferralAlertThreshold {
			t.Fatalf("StuckTimeouts() = (%d, %d), want (1, %d): a prerequisite that never recovers must become visible",
				count, worst, deferralAlertThreshold)
		}
	})

	// A healthy resolve must not touch the tracker at all.
	t.Run("success leaves packets due", func(t *testing.T) {
		t.Parallel()

		s, expired := prereqServices(t, 41)
		_, state, ok := s.resolveCosmosTimeoutProofState(s.BatchBuilder.PendingTracker, expired, okUpdate, func() (*client.EthereumClientState, error) {
			t.Fatal("fallback must not run when the update carries a client state")
			return nil, nil
		})
		if !ok || state == nil {
			t.Fatal("a successful resolve must report ok with a client state")
		}
		if due := s.BatchBuilder.PendingTracker.GetDue(time.Now()); len(due) != len(expired) {
			t.Fatalf("GetDue() = %d, want %d: a successful resolve must not defer anything", len(due), len(expired))
		}
	})
}

func TestHandleCosmosTimeoutBatchFailure(t *testing.T) {
	t.Run("permanent cosmos timeout batch isolates the poison packet", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		poison := timeoutRetryPacket(11)
		healthy := timeoutRetryPacket(12)
		tracker.Add(poison, 10)
		tracker.Add(healthy, 10)

		handler := &timeoutBatchHandler{}
		svc := New(handler, nil, DefaultConfig())
		processed := []pendingPacketInfo{{Packet: poison}, {Packet: healthy}}
		svc.handleCosmosTimeoutBatchFailure(context.Background(), CosmosEndpoint{}, tracker, nil, []any{"poison", "healthy"}, processed, ErrPermanentRelayFailure)

		if len(handler.singletons) != 2 {
			t.Fatalf("singleton probes = %d, want 2", len(handler.singletons))
		}
		if tracker.Len() != 1 {
			t.Fatalf("tracked packets = %d, want only poison packet", tracker.Len())
		}
		remaining := tracker.GetAll()
		if len(remaining) != 1 || remaining[0].Packet.Sequence != poison.Sequence || remaining[0].TimeoutAttempts != 1 {
			t.Fatalf("remaining tracker state = %#v, want one charged poison packet", remaining)
		}
	})

	t.Run("update prefix failure defers all", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		first := timeoutRetryPacket(31)
		second := timeoutRetryPacket(32)
		tracker.Add(first, 10)
		tracker.Add(second, 10)

		svc := New(nil, nil, DefaultConfig())
		svc.handleCosmosTimeoutBatchFailure(
			context.Background(), CosmosEndpoint{}, tracker,
			[]any{"update"}, []any{"first", "second"},
			[]pendingPacketInfo{{Packet: first}, {Packet: second}},
			&BatchPartialError{SucceededCount: 0, Err: ErrPermanentRelayFailure},
		)

		for _, info := range tracker.GetAll() {
			if info.TimeoutAttempts != 0 || info.Deferrals != 1 {
				t.Fatalf("seq=%d retry state = attempts:%d deferrals:%d, want 0:1", info.Packet.Sequence, info.TimeoutAttempts, info.Deferrals)
			}
		}
	})

	t.Run("singleton probe defers transient failure", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		packet := timeoutRetryPacket(41)
		tracker.Add(packet, 10)

		svc := New(&transientTimeoutBatchHandler{}, nil, DefaultConfig())
		svc.handleCosmosTimeoutBatchFailure(
			context.Background(), CosmosEndpoint{}, tracker,
			nil, []any{"timeout"}, []pendingPacketInfo{{Packet: packet}}, ErrPermanentRelayFailure,
		)

		info := tracker.GetAll()
		if len(info) != 1 || info[0].TimeoutAttempts != 0 || info[0].Deferrals != 1 {
			t.Fatalf("singleton transient retry state = %#v, want one deferred uncharged packet", info)
		}
	})
}

// TestUpdateMsgsFailed is the second P1. SendCosmosTxBatch submits
// [updateMsgs..., timeoutMsgs...] and SucceededCount is a prefix length, so
// fewer successes than update messages proves no timeout was ever executed.
func TestUpdateMsgsFailed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		updateMsgCount int
		want           bool
	}{
		{
			name:           "update prefix incomplete: no timeout could have run",
			err:            &BatchPartialError{SucceededCount: 1, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           true,
		},
		{
			name:           "update prefix complete: the failure is in the timeouts",
			err:            &BatchPartialError{SucceededCount: 2, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           false,
		},
		{
			name:           "whole tx failed with update msgs present: nothing executed",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 1,
			want:           true,
		},
		{
			name:           "no update msgs: the timeouts are the only suspects",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 0,
			want:           false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := updateMsgsFailed(tc.err, tc.updateMsgCount); got != tc.want {
				t.Fatalf("updateMsgsFailed = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestScanForCosmosTimeouts_HonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	svc := New(nil, nil, DefaultConfig())
	svc.TrackCosmosPending(channeltypesv2.Packet{
		SourceClient:      "07-tendermint-0",
		DestinationClient: "client-0",
		Sequence:          1,
		TimeoutTimestamp:  1,
	}, 100)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		svc.scanForCosmosTimeouts(ctx, CosmosEndpoint{}, EVMEndpoint{
			Client:    client,
			Contracts: EVMContracts{Router: common.HexToAddress("0x1")},
		}, "08-wasm-0")
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Cosmos timeout scan ignored cancellation")
	}
}

func TestScanForEVMTimeouts(t *testing.T) {
	t.Run("tracker removal", func(t *testing.T) {
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
	})

	t.Run("shared update failure runs once and defers all", func(t *testing.T) {
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
	})

	t.Run("clears shared deferrals after recovery", func(t *testing.T) {
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
	})
}

// TestShouldRelayCosmosTimeoutToEth pins the direction-detection logic that
// decides whether a CosmosTimeout event needs a follow-up timeoutPacket on
// ETH. The old guard compared SourceClient against the router client ID — a
// mismatch in IBC v2 because:
//
//   - For a Cosmos→ETH packet, the Cosmos-emitted event has
//     SourceClient="08-wasm-0" (Cosmos's client identifying the ETH chain) and
//     DestinationClient="cosmoshub-1" (ETH's client identifying Cosmos).
//   - For an ETH→Cosmos packet, it's mirrored: source=cosmoshub-1, dest=08-wasm-0.
//
// `CosmosRouterClientID` in the relayer config is the ETH-side router's name
// for this Cosmos chain ("cosmoshub-1"), so the correct field to compare is
// `DestinationClient`. The old comparison silently inverted both branches.
func TestShouldRelayCosmosTimeoutToEth(t *testing.T) {
	tests := []struct {
		name                 string
		packet               *channeltypesv2.Packet
		cosmosRouterClientID string
		want                 bool
	}{
		{
			name: "Cosmos-originated packet skips ETH relay",
			packet: &channeltypesv2.Packet{
				SourceClient:      "08-wasm-0",   // Cosmos's client for ETH
				DestinationClient: "cosmoshub-1", // ETH's client for Cosmos
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false, // Cosmos's MsgTimeout already refunded the sender
		},
		{
			name: "ETH-originated packet is relayed",
			packet: &channeltypesv2.Packet{
				SourceClient:      "cosmoshub-1", // ETH's client for Cosmos
				DestinationClient: "08-wasm-0",   // Cosmos's client for ETH
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true, // ETH still holds the commitment, needs ICS26Router.timeoutPacket
		},
		{
			name:                 "nil packet is not relayed",
			packet:               nil,
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false,
		},
		{
			name: "unrelated client IDs are relayed (treated as ETH-originated)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true,
		},
		{
			name: "empty cosmosRouterClientID is not relayed (misconfiguration)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "",
			want:                 false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldRelayCosmosTimeoutToEth(tc.packet, tc.cosmosRouterClientID)
			if got != tc.want {
				t.Fatalf("ShouldRelayCosmosTimeoutToEth=%v want=%v", got, tc.want)
			}
		})
	}
}
