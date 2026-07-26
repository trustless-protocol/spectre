package server

import (
	"fmt"

	"attestor/arbitrum"
	attestorpb "attestor/types/attestor"
)

func runModeFromProto(mode attestorpb.RunMode) (arbitrum.RunMode, error) {
	switch mode {
	case attestorpb.RunMode_RUN_MODE_UNSAFE:
		return arbitrum.RunModeUnsafe, nil
	case attestorpb.RunMode_RUN_MODE_SAFE:
		return arbitrum.RunModeSafe, nil
	case attestorpb.RunMode_RUN_MODE_FINALIZED:
		return arbitrum.RunModeFinalized, nil
	default:
		return "", fmt.Errorf(
			"run_mode must be explicitly set to unsafe, safe, or finalized, got %s",
			mode,
		)
	}
}
