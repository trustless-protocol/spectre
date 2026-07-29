package main

import (
	"errors"
	"testing"
)

// TestIsSignatureSlotSkew: a submitted update can be rejected purely because its
// signature slot is ahead of the slot the client derives from the Cosmos block
// timestamp, which trails real time by a few seconds. That is transient and must be
// retried with the SAME message; every other rejection must surface immediately.
func TestIsSignatureSlotSkew(t *testing.T) {
	skew := errors.New("cosmos tx failed at CheckTx with code 14 codespace=08-wasm: " +
		"verify client message failed: update signature slot 635 is more recent than " +
		"the calculated current slot 632: wasm contract call failed")
	if !isSignatureSlotSkew(skew) {
		t.Error("host-clock skew should be retryable")
	}

	for name, err := range map[string]error{
		"nil":           nil,
		"frozen client": errors.New("client is frozen at height 42"),
		"sequence":      errors.New("account sequence mismatch, expected 6, got 5"),
		"checksum":      errors.New("wasm checksum has not been stored"),
	} {
		if isSignatureSlotSkew(err) {
			t.Errorf("%s must not be treated as clock skew", name)
		}
	}
}

// TestEthClientUpdateLag documents the freshness window the builders rely on: the
// pinned client only advances in finalized steps, so demanding zero lag would submit
// an update before every header build, while the window must stay well inside a
// default execution node's state retention for the proof to remain servable.
func TestEthClientUpdateLag(t *testing.T) {
	if ethClientUpdateLag == 0 {
		t.Fatal("a zero lag window would trigger an update on every header build")
	}
	if ethClientUpdateLag >= 128 {
		t.Fatalf("lag window %d is at or beyond a default node's ~128-block state window",
			ethClientUpdateLag)
	}
}
