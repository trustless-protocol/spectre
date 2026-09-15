package transaction

import (
	"context"
	"errors"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"relayer/services"
)

func TestDuplicateDropIsOnlySafeForOneMessage(t *testing.T) {
	cases := []struct {
		name     string
		msgCount int
		want     bool
	}{
		{"single message is the whole transaction", 1, true},
		{"two-message batch", 2, false},
		{"large batch", 25, false},
		{"empty batch", 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := duplicateDropIsSafe(tc.msgCount); got != tc.want {
				t.Fatalf("duplicateDropIsSafe(%d) = %v, want %v", tc.msgCount, got, tc.want)
			}
		})
	}
}

// TestIsCosmosDuplicatePacketError covers the deterministic duplicate/redundant
// signals that must be dropped (not re-queued): channelv2 11/12 and, critically,
// channel/22 ErrRedundantTx — the ante-handler rejection whose omission caused an
// infinite re-submit gas drain.
func TestIsCosmosDuplicatePacketError(t *testing.T) {
	cases := []struct {
		name      string
		codespace string
		code      uint32
		want      bool
	}{
		{"channelv2 ack-exists", "channelv2", 11, true},
		{"channelv2 noop", "channelv2", 12, true},
		{"channel redundant-tx", "channel", 22, true},
		{"channel noop (23) not matched", "channel", 23, false},
		{"channelv2 other code", "channelv2", 5, false},
		{"unrelated codespace", "ibc", 22, false},
		{"zero", "", 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isCosmosDuplicatePacketError(tc.codespace, tc.code); got != tc.want {
				t.Fatalf("isCosmosDuplicatePacketError(%q,%d) = %v, want %v", tc.codespace, tc.code, got, tc.want)
			}
		})
	}
}

func TestDuplicateFailuresBeforeDeliverTxArePermanent(t *testing.T) {
	for _, stage := range []string{"Simulation", "CheckTx"} {
		t.Run(stage, func(t *testing.T) {
			err := newCosmosPreDeliverFailure(stage, 22, "channel", "redundant tx", nil)
			if !errors.Is(err, services.ErrPermanentRelayFailure) {
				t.Fatalf("duplicate %s failure is not permanent: %v", stage, err)
			}
		})
	}

	err := newCosmosPreDeliverFailure("CheckTx", 32, "sdk", "account sequence mismatch", nil)
	if errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("transient CheckTx failure was marked permanent: %v", err)
	}
}

func TestSplitAroundDuplicateStopsAtTheFirstFailingHalf(t *testing.T) {
	msgs := make([]sdk.Msg, 4)
	firstHalfErr := errors.New("first half reverted")

	var calls int
	send := func(
		_ context.Context,
		_ services.CosmosEndpoint,
		sub []sdk.Msg,
		_, sequence uint64,
		allowSplit bool,
	) (uint64, int, error) {
		calls++
		if !allowSplit {
			t.Fatal("duplicate halves must retain recursive splitting")
		}
		if calls == 1 {
			return sequence, 0, firstHalfErr
		}
		return sequence + uint64(len(sub)), len(sub), nil
	}

	_, succeeded, err := splitCosmosBatchAfterDuplicateWith(
		send, context.Background(), services.CosmosEndpoint{}, msgs, 7, 100,
	)
	if !errors.Is(err, firstHalfErr) {
		t.Fatalf("want the first half error propagated, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("the second half must not be attempted after the first fails, got %d calls", calls)
	}
	if succeeded != 0 {
		t.Fatalf("nothing was delivered, so the prefix count must be 0, got %d", succeeded)
	}
}

func TestSplitAroundDuplicateThreadsSequenceAndSumsOnSuccess(t *testing.T) {
	msgs := make([]sdk.Msg, 4)
	var seenSequences []uint64

	send := func(
		_ context.Context,
		_ services.CosmosEndpoint,
		sub []sdk.Msg,
		_, sequence uint64,
		allowSplit bool,
	) (uint64, int, error) {
		if !allowSplit {
			t.Fatal("duplicate halves must retain recursive splitting")
		}
		seenSequences = append(seenSequences, sequence)
		return sequence + uint64(len(sub)), len(sub), nil
	}

	finalSequence, succeeded, err := splitCosmosBatchAfterDuplicateWith(
		send, context.Background(), services.CosmosEndpoint{}, msgs, 7, 100,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if succeeded != 4 {
		t.Fatalf("want all 4 messages reported delivered, got %d", succeeded)
	}
	if finalSequence != 104 {
		t.Fatalf("want the final sequence threaded out (104), got %d", finalSequence)
	}
	if len(seenSequences) != 2 || seenSequences[0] != 100 || seenSequences[1] != 102 {
		t.Fatalf("halves must be signed with consecutive sequences [100 102], got %v", seenSequences)
	}
}

// The batch rung and the split that carries it out are computed in two
// different places: CosmosBatchLadder returns a size, and
// splitCosmosBatchAfterDuplicateWith halves the slice itself. Both were pinned
// -- the ladder by TestLaddersChangeSomethingEachRungAndEnd, the split by the
// hard-coded sequences above -- but nothing tied them to each other, so either
// could be changed alone and stay green. The rung is what the operator reads in
// "splitting ... of %d to %d"; when the two disagree that line reports a size
// the retry never submitted, and C1's claim that the call site only describes
// the ladder stops being true.
func TestSplitAfterDuplicateSubmitsTheSizeTheLadderNames(t *testing.T) {
	// Odd counts are the interesting ones: integer division is where two
	// independent halvings drift first.
	for _, count := range []int{2, 3, 4, 5, 9, 32} {
		t.Run(fmt.Sprintf("%d messages", count), func(t *testing.T) {
			rung, ok := CosmosBatchLadder(count)(0)
			if !ok {
				t.Fatalf("the ladder reports no rung for %d messages, so the split must not be reached at all", count)
			}

			var firstSubmitted int
			send := func(
				_ context.Context,
				_ services.CosmosEndpoint,
				sub []sdk.Msg,
				_, sequence uint64,
				_ bool,
			) (uint64, int, error) {
				if firstSubmitted == 0 {
					firstSubmitted = len(sub)
				}
				return sequence + uint64(len(sub)), len(sub), nil
			}

			if _, _, err := splitCosmosBatchAfterDuplicateWith(
				send, context.Background(), services.CosmosEndpoint{}, make([]sdk.Msg, count), 7, 100,
			); err != nil {
				t.Fatalf("split: %v", err)
			}

			if firstSubmitted != int(rung.To) {
				t.Fatalf("the ladder named %d messages but the split submitted %d; the log line and the retry disagree",
					rung.To, firstSubmitted)
			}
		})
	}
}
