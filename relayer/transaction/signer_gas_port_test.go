package transaction

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"testing"

	"relayer/keys"
)

type testSigner struct {
	eth    *ecdsa.PrivateKey
	cosmos []byte
	err    error
	calls  int
}

func (s *testSigner) EthKey() (*ecdsa.PrivateKey, error) {
	s.calls++
	return s.eth, s.err
}

func (s *testSigner) CosmosKeyBytes() ([]byte, error) {
	s.calls++
	return append([]byte(nil), s.cosmos...), s.err
}

func TestHandlerValidateKeysUsesInstalledSigner(t *testing.T) {
	t.Setenv("ETH_PRIVATE_KEY", "")
	t.Setenv("COSMOS_PRIVATE_KEY", "")
	eth, err := keys.RestoreKey(testCosmosPrivKeyHex)
	if err != nil {
		t.Fatalf("RestoreKey: %v", err)
	}
	cosmos, err := hex.DecodeString(testCosmosPrivKeyHex)
	if err != nil {
		t.Fatalf("decode Cosmos key: %v", err)
	}
	signer := &testSigner{eth: eth, cosmos: cosmos}
	h := &Handler{}
	h.SetSigner(signer)
	if err := h.ValidateKeys(); err != nil {
		t.Fatalf("ValidateKeys: %v", err)
	}
	if signer.calls != 2 {
		t.Fatalf("signer calls = %d, want 2", signer.calls)
	}
}

func TestSetSignerNilLeavesDefaultAvailable(t *testing.T) {
	h := &Handler{}
	h.SetSigner(nil)
	if h.signer != nil {
		t.Fatal("nil SetSigner installed a signer")
	}
	if h.keySigner() == nil {
		t.Fatal("default signer is nil after SetSigner(nil)")
	}
}

func TestEnvSignerCosmosKeyBytesAreCopied(t *testing.T) {
	t.Setenv("COSMOS_PRIVATE_KEY", testCosmosPrivKeyHex)
	signer := NewEnvSigner()
	first, err := signer.CosmosKeyBytes()
	if err != nil {
		t.Fatalf("first CosmosKeyBytes: %v", err)
	}
	first[0] ^= 0xff
	second, err := signer.CosmosKeyBytes()
	if err != nil {
		t.Fatalf("second CosmosKeyBytes: %v", err)
	}
	if first[0] == second[0] {
		t.Fatal("CosmosKeyBytes returned shared backing storage")
	}
}

func TestProbeErrorClassesAreDisjoint(t *testing.T) {
	for _, text := range []string{"out of gas", "gas required exceeds allowance", "exceeds block gas limit", "intrinsic gas too low", "dial tcp: timeout"} {
		err := errors.New(text)
		count := 0
		for _, class := range []bool{isExecutionOutOfGas(err), isGasAllowanceRejection(err), isProbeTransportError(err)} {
			if class {
				count++
			}
		}
		if count != 1 {
			t.Fatalf("%q matched %d classes, want exactly one", text, count)
		}
	}
}

// An atomic batch keeps its all-or-nothing shape: it must never be split, only
// given more gas, however many messages it carries.
func TestPlanCosmosOutOfGasTerminatesAndPreservesAtomicity(t *testing.T) {
	if next, _, ok := oogPlan(t, 2, 0, 1, 10, true); !ok || next.What != knobBatchSize {
		t.Fatalf("splittable batch planned %v (ok=%v), want a split", next, ok)
	}
	if next, _, ok := oogPlan(t, 2, 0, 1, 10, false); !ok || next.What != knobCosmosGas {
		t.Fatalf("atomic batch planned %v (ok=%v), want a gas escalation and never a split", next, ok)
	}
	if next, _, ok := oogPlan(t, 1, len(cosmosGasHeadroom)-1, 10, 10, false); ok {
		t.Fatalf("at the ceiling planned %v, want permanent", next)
	}
}

func TestFallbackCosmosGasLimitIncreasesAcrossOOGRetries(t *testing.T) {
	const fallback = uint64(200_000)
	previous := uint64(0)
	for step := range cosmosGasHeadroom {
		got := applyCosmosGasHeadroom(fallback, step)
		if got <= previous {
			t.Fatalf("fallback gas at step %d = %d, want greater than previous %d", step, got, previous)
		}
		previous = got
	}
}

func TestEthGasLimitRejectsZeroAndMalformedOverrides(t *testing.T) {
	for _, gasText := range []string{"0", "not-a-number"} {
		if _, err := EthGasLimit(gasText); err == nil {
			t.Fatalf("EthGasLimit(%q) accepted an invalid override", gasText)
		}
	}
	if got, err := EthGasLimit(""); err != nil || got != defaultEthGasLimit {
		t.Fatalf("EthGasLimit(\"\") = (%d, %v), want (%d, nil)", got, err, defaultEthGasLimit)
	}
	if got, err := EthGasLimit("4000000"); err != nil || got != 4_000_000 {
		t.Fatalf("EthGasLimit(valid) = (%d, %v), want (4000000, nil)", got, err)
	}
}
