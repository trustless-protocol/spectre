package transaction

import "testing"

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
