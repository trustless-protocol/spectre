package prover

import (
	"strings"
	"testing"
)

func TestNewProofBackendFromSelection_DefaultsToNative(t *testing.T) {
	backend, err := NewProofBackendFromSelection("", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend.Name() != proverBackendNative {
		t.Fatalf("got backend %q, want %q", backend.Name(), proverBackendNative)
	}
}

func TestNewProofBackendFromEnv_RejectsUnsupportedBackend(t *testing.T) {
	t.Setenv("GPU_PROVE", "")
	t.Setenv("GNARK_PROVER_BACKEND", "bogus")

	_, err := NewProofBackendFromEnv()
	if err == nil || !strings.Contains(err.Error(), "unsupported GNARK_PROVER_BACKEND") {
		t.Fatalf("want unsupported backend error, got %v", err)
	}
}

func TestNewProofBackendFromEnv_GPUProveAliasWins(t *testing.T) {
	t.Setenv("GPU_PROVE", "1")
	t.Setenv("GNARK_PROVER_BACKEND", "native")

	backend, err := NewProofBackendFromEnv()
	if err != nil {
		if !strings.Contains(err.Error(), "requires building with -tags=icicle") {
			t.Fatalf("unexpected error: %v", err)
		}
		return
	}
	if backend.Name() != proverBackendICICLE {
		t.Fatalf("got backend %q, want %q", backend.Name(), proverBackendICICLE)
	}
}
