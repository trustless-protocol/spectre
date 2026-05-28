package prover

import (
	"strings"
	"testing"
)

func TestNewProofBackend_DefaultsToNative(t *testing.T) {
	backend, err := NewProofBackend(false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if backend.Name() != proverBackendNative {
		t.Fatalf("got backend %q, want %q", backend.Name(), proverBackendNative)
	}
}

// On a non-icicle build NewProofBackend(true) must return a clear error
// pointing at the missing build tag; on an icicle build it must return the
// ICICLE backend. One test, covers both modes.
func TestNewProofBackend_GPURequiresIcicleBuild(t *testing.T) {
	backend, err := NewProofBackend(true)
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

func TestGPUProveEnvEnabled(t *testing.T) {
	cases := []struct {
		raw  string
		want bool
	}{
		{"", false},
		{"0", false},
		{"false", false},
		{"no", false},
		{"1", true},
		{"true", true},
		{"YES", true},
		{"On", true},
	}
	for _, c := range cases {
		t.Run(c.raw, func(t *testing.T) {
			t.Setenv("GPU_PROVE", c.raw)
			if got := GPUProveEnvEnabled(); got != c.want {
				t.Fatalf("GPU_PROVE=%q → %v, want %v", c.raw, got, c.want)
			}
		})
	}
}
