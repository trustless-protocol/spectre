package l2rollup

import (
	"testing"

	"relayer/chain"

	contractICS26Router "relayer/bindings/ICS26Router"
)

func TestKeepIndexed(t *testing.T) {
	events := []chain.Event{
		{Height: 1}, {Height: 2}, {Height: 3},
	}
	cases := []struct {
		name    string
		indices []int
		want    []uint64
	}{
		{"none re-queued", nil, nil},
		{"middle re-queued", []int{1}, []uint64{2}},
		{"all re-queued preserves order", []int{0, 1, 2}, []uint64{1, 2, 3}},
		{"out-of-range ignored", []int{2, 9, -1}, []uint64{3}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := keepIndexed(events, tc.indices)
			if len(got) != len(tc.want) {
				t.Fatalf("kept %d events, want %d", len(got), len(tc.want))
			}
			for i, h := range tc.want {
				if got[i].Height != h {
					t.Fatalf("kept[%d].Height = %d, want %d", i, got[i].Height, h)
				}
			}
		})
	}
}

// TestL2AckToEvent_EmptyAckSkipped guards the no-ack skip: a WriteAcknowledgement
// with no acknowledgement bytes cannot become a MsgAcknowledgement, so it must be
// dropped rather than relayed.
func TestL2AckToEvent_EmptyAckSkipped(t *testing.T) {
	ev := &contractICS26Router.ContractICS26RouterWriteAcknowledgement{
		Acknowledgements: nil,
	}
	if _, ok := l2AckToEvent(ev); ok {
		t.Fatal("empty acknowledgement must be skipped")
	}
}
