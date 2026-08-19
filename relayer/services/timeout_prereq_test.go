package services

import (
	"fmt"
	"testing"
	"time"

	client "relayer/client"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

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

func okUpdate() (*EthClientUpdateResult, error) {
	return &EthClientUpdateResult{EthClientState: &client.EthereumClientState{LatestSlot: 1}}, nil
}

// TestTimeoutProofStatePrerequisiteFailuresDeferWithoutCharging is the P1 from
// review, raised twice: both shared-prerequisite exits returned bare.
//
// Neither failure is attributable to any packet -- they fail identically for the
// whole scan -- so nothing may be charged. But NOT deferring is what made the
// bug expensive and invisible at once: the next 30s scan rebuilds the client
// state, finality and update construction from scratch forever, and because
// Deferrals stays at zero the STUCK line and the stuck gauge cannot report it.
func TestTimeoutProofStatePrerequisiteFailuresDeferWithoutCharging(t *testing.T) {
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
}

// TestTimeoutProofStateRepeatedFailuresBecomeVisible: deferrals never
// dead-letter, so the stuck gauge is the only thing that can surface a
// prerequisite that never recovers. Without the deferral it reads zero forever.
func TestTimeoutProofStateRepeatedFailuresBecomeVisible(t *testing.T) {
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
}

// A healthy resolve must not touch the tracker at all.
func TestTimeoutProofStateSuccessLeavesPacketsDue(t *testing.T) {
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
}
