package services

import (
	"context"
	"errors"
	"testing"
	"time"

	client "relayer/client"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func timeoutRetryPacket(sequence uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{SourceClient: "client-0", Sequence: sequence, TimeoutTimestamp: 1}
}

func TestTimeoutRetryBackoffAndDeadLetter(t *testing.T) {
	tracker := NewPendingPacketTracker()
	packet := timeoutRetryPacket(1)
	tracker.Add(packet, 10)
	now := time.Now()

	if tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, now) {
		t.Fatal("first deterministic failure must not dead-letter")
	}
	if due := tracker.GetDue(now); len(due) != 0 {
		t.Fatalf("packet under retry backoff is immediately due: %d", len(due))
	}
	if due := tracker.GetDue(now.Add(timeoutRetryBackoffMax)); len(due) != 1 {
		t.Fatalf("packet did not become due after its bounded backoff: %d", len(due))
	}

	for attempt := 1; attempt <= maxTimeoutAttempts; attempt++ {
		dead := tracker.recordTimeoutFailureForTest(packet.SourceClient, packet.Sequence, now)
		if attempt < maxTimeoutAttempts && dead {
			t.Fatalf("attempt %d dead-lettered early", attempt+1)
		}
		if attempt == maxTimeoutAttempts && !dead {
			t.Fatalf("attempt %d did not dead-letter", attempt+1)
		}
	}
	if tracker.Len() != 0 || tracker.DeadLetteredTimeouts() != 1 {
		t.Fatalf("dead-letter state = tracked:%d dead:%d, want 0:1", tracker.Len(), tracker.DeadLetteredTimeouts())
	}
}

func TestApplyTimeoutOutcomeDistinguishesDeferralFailureAndCompletion(t *testing.T) {
	tracker := NewPendingPacketTracker()
	now := time.Now()
	packet := timeoutRetryPacket(2)
	tracker.Add(packet, 10)
	info := tracker.GetAll()[0]

	applyTimeoutOutcome(tracker, info, timeoutDeferred, now, "test")
	if due := tracker.GetDue(now); len(due) != 0 {
		t.Fatalf("deferred packet remained due: %d", len(due))
	}
	applyTimeoutOutcome(tracker, info, timeoutNotDue, now, "test")
	if got := tracker.GetAll()[0].Deferrals; got != 2 {
		t.Fatalf("not-due packet deferrals = %d, want 2", got)
	}
	if due := tracker.GetDue(now); len(due) != 0 {
		t.Fatalf("not-due packet remained immediately eligible: %d", len(due))
	}
	applyTimeoutOutcome(tracker, info, timeoutAlreadyReceived, now, "test")
	if tracker.Len() != 0 {
		t.Fatal("already-received packet remained in timeout tracker")
	}

	if got := timeoutOutcomeForError(client.ErrPacketAlreadyReceived); got != timeoutAlreadyReceived {
		t.Fatalf("already received outcome = %d", got)
	}
	if got := timeoutOutcomeForError(ErrPermanentRelayFailure); got != timeoutFailed {
		t.Fatalf("permanent failure outcome = %d", got)
	}
	if got := timeoutOutcomeForError(errors.New("rpc unavailable")); got != timeoutDeferred {
		t.Fatalf("transient failure outcome = %d", got)
	}
}

type timeoutBatchHandler struct {
	TransactionHandler
	singletons [][]any
}

func (h *timeoutBatchHandler) SendCosmosTxBatch(_ context.Context, _ CosmosEndpoint, msgs []any) error {
	h.singletons = append(h.singletons, msgs)
	if len(msgs) != 1 {
		return errors.New("test only accepts singleton probes")
	}
	if msgs[0] == "poison" {
		return ErrPermanentRelayFailure
	}
	return nil
}

func TestPermanentCosmosTimeoutBatchIsolatesThePoisonPacket(t *testing.T) {
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
}

func TestCosmosTimeoutUpdatePrefixFailureDefersAll(t *testing.T) {
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
}

type transientTimeoutBatchHandler struct{ TransactionHandler }

func (h *transientTimeoutBatchHandler) SendCosmosTxBatch(_ context.Context, _ CosmosEndpoint, _ []any) error {
	return errors.New("cosmos rpc unavailable")
}

func TestCosmosTimeoutSingletonProbeDefersTransientFailure(t *testing.T) {
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
}

func TestResolveCosmosTimeoutProofStateDefersSharedFailures(t *testing.T) {
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
}

func TestQueueStateIncludesStuckAndDeadLetterTimeouts(t *testing.T) {
	svc := New(nil, nil, DefaultConfig())
	packet := timeoutRetryPacket(21)
	svc.BatchBuilder.PendingTracker.Add(packet, 10)
	now := time.Now()
	for range deferralAlertThreshold {
		svc.BatchBuilder.PendingTracker.deferTimeoutRetryForTest(packet.SourceClient, packet.Sequence, now)
	}

	dead := timeoutRetryPacket(22)
	svc.BatchBuilder.EthPendingTracker.Add(dead, 10)
	for attempt := 0; attempt <= maxTimeoutAttempts; attempt++ {
		svc.BatchBuilder.EthPendingTracker.recordTimeoutFailureForTest(dead.SourceClient, dead.Sequence, now)
	}

	state := svc.queueStateAt(now)
	if state.cosmosStuck != 1 || state.cosmosWorstDeferrals != deferralAlertThreshold {
		t.Fatalf("stuck state = %d/%d, want 1/%d", state.cosmosStuck, state.cosmosWorstDeferrals, deferralAlertThreshold)
	}
	if state.ethTimeoutDead != 1 {
		t.Fatalf("eth timeout dead letters = %d, want 1", state.ethTimeoutDead)
	}
}

func TestClientUpdateAgesUseTrustedChainTimestamps(t *testing.T) {
	svc := New(nil, nil, DefaultConfig())
	now := time.Now()
	if _, _, ok := svc.clientUpdateAges(now); ok {
		t.Fatal("client update ages reported before either direction was observed")
	}

	svc.ObserveCosmosOnEVMUpdate(now.Add(-2 * time.Minute))
	svc.ObserveEVMOnCosmosUpdate(now.Add(-3 * time.Minute))
	cosmosAge, evmAge, ok := svc.clientUpdateAges(now)
	if !ok || cosmosAge != 2*time.Minute || evmAge != 3*time.Minute {
		t.Fatalf("client update ages = (%s, %s, %v), want (2m, 3m, true)", cosmosAge, evmAge, ok)
	}

	// An older observation must not make the reported client state regress.
	svc.ObserveCosmosOnEVMUpdate(now.Add(-10 * time.Minute))
	cosmosAge, _, _ = svc.clientUpdateAges(now)
	if cosmosAge != 2*time.Minute {
		t.Fatalf("older observation regressed cosmos-on-eth age to %s", cosmosAge)
	}
}
