package server

import (
	"testing"

	"attestor/arbitrum"
	attestorpb "attestor/types/attestor"
)

func TestRunModeFromProto(t *testing.T) {
	tests := []struct {
		proto attestorpb.RunMode
		want  arbitrum.RunMode
	}{
		{proto: attestorpb.RunMode_RUN_MODE_UNSAFE, want: arbitrum.RunModeUnsafe},
		{proto: attestorpb.RunMode_RUN_MODE_SAFE, want: arbitrum.RunModeSafe},
		{proto: attestorpb.RunMode_RUN_MODE_FINALIZED, want: arbitrum.RunModeFinalized},
	}
	for _, test := range tests {
		t.Run(string(test.want), func(t *testing.T) {
			got, err := runModeFromProto(test.proto)
			if err != nil {
				t.Fatalf("convert protobuf run mode: %v", err)
			}
			if got != test.want {
				t.Fatalf("converted mode: got %q want %q", got, test.want)
			}
		})
	}

	for _, mode := range []attestorpb.RunMode{
		attestorpb.RunMode_RUN_MODE_UNSPECIFIED,
		attestorpb.RunMode(99),
	} {
		if _, err := runModeFromProto(mode); err == nil {
			t.Fatalf("unsupported protobuf mode %d was accepted", mode)
		}
	}
}
