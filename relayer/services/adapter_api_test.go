package services

import (
	"bytes"
	"errors"
	"log"
	"path/filepath"
	"strings"
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func TestUntrackPendingLogsPersistenceFailure(t *testing.T) {
	tests := []struct {
		name    string
		track   func(*Services, channeltypesv2.Packet) bool
		untrack func(*Services, channeltypesv2.Packet)
		tracker func(*Services) *PendingPacketTracker
		kind    string
	}{
		{
			name:    "cosmos",
			track:   func(s *Services, packet channeltypesv2.Packet) bool { return s.TrackCosmosPending(packet, 1) },
			untrack: (*Services).UntrackCosmosPending,
			tracker: func(s *Services) *PendingPacketTracker { return s.BatchBuilder.PendingTracker },
			kind:    "Cosmos",
		},
		{
			name:    "l2",
			track:   func(s *Services, packet channeltypesv2.Packet) bool { return s.TrackL2Pending(packet, 1) },
			untrack: (*Services).UntrackL2Pending,
			tracker: func(s *Services) *PendingPacketTracker { return s.BatchBuilder.L2PendingTracker },
			kind:    "L2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc, err := NewWithPendingState(nil, nil, DefaultConfig(), filepath.Join(t.TempDir(), "state"))
			if err != nil {
				t.Fatal(err)
			}
			packet := channeltypesv2.Packet{SourceClient: "source", DestinationClient: "destination", Sequence: 42}
			if !tc.track(svc, packet) {
				t.Fatal("could not persist pending packet")
			}
			tc.tracker(svc).writeState = func(string, []byte) error { return errors.New("disk full") }

			var logs bytes.Buffer
			originalWriter := log.Writer()
			log.SetOutput(&logs)
			t.Cleanup(func() { log.SetOutput(originalWriter) })

			tc.untrack(svc, packet)

			want := "[Services][ATTENTION] failed to persist removal of settled " + tc.kind + " packet seq=42"
			if !strings.Contains(logs.String(), want) {
				t.Fatalf("persistence failure was not logged with context %q: %q", want, logs.String())
			}
			if got := tc.tracker(svc).Len(); got != 1 {
				t.Fatalf("failed persistence removed %s pending packet", tc.name)
			}
		})
	}
}
