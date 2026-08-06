package l2rollup

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// semanticJSONEqual compares two JSON documents ignoring field order and
// insignificant whitespace — the equality Dũng specified for the cross-language
// contract (do NOT assert raw text equality).
func semanticJSONEqual(t *testing.T, got, want []byte) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal(got, &g); err != nil {
		t.Fatalf("unmarshal got: %v\n%s", err, got)
	}
	if err := json.Unmarshal(want, &w); err != nil {
		t.Fatalf("unmarshal want: %v", err)
	}
	if !reflect.DeepEqual(g, w) {
		t.Fatalf("JSON not semantically equal:\n got: %s\nwant: %s", got, want)
	}
}

// attestedExampleJSON is the wire-shape fixture for the only header the
// attestor-trusted clients accept. Dummy proofs — it locks field names and value
// representations, not any real chain state. The 256-byte logs_bloom comes from a
// 512-zero-hex-digit string.
func attestedExampleJSON() string {
	return `{
  "type": "header",
  "value": {
    "l2_header": {
      "parent_hash": "0x1111111111111111111111111111111111111111111111111111111111111111",
      "ommers_hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
      "beneficiary": "0x3333333333333333333333333333333333333333",
      "state_root": "0x4444444444444444444444444444444444444444444444444444444444444444",
      "transactions_root": "0x5555555555555555555555555555555555555555555555555555555555555555",
      "receipts_root": "0x6666666666666666666666666666666666666666666666666666666666666666",
      "logs_bloom": "0x` + strings.Repeat("0", 512) + `",
      "difficulty": "0x0",
      "number": 42,
      "gas_limit": 30000000,
      "gas_used": 21000,
      "timestamp": 1700000000,
      "extra_data": "0x",
      "mix_hash": "0x7777777777777777777777777777777777777777777777777777777777777777",
      "nonce": "0x0000000000000000",
      "base_fee_per_gas": "0x9",
      "withdrawals_root": "0x8888888888888888888888888888888888888888888888888888888888888888",
      "blob_gas_used": 11,
      "excess_blob_gas": 12,
      "parent_beacon_block_root": "0x9999999999999999999999999999999999999999999999999999999999999999"
    },
    "router_proof": { "proof": [[1, 2], [3]] }
  }
}`
}

// TestAttestedL2Header_RoundTrip decodes the canonical example into the Go types and
// re-encodes it, asserting the wire form is semantically identical. This is the
// freeze test for the encoding contract now shared by every chain — the per-chain
// fixtures it replaces died with the settlement headers.
func TestAttestedL2Header_RoundTrip(t *testing.T) {
	want := attestedExampleJSON()

	var env clientMessageEnvelope
	if err := json.Unmarshal([]byte(want), &env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if env.Type != "header" {
		t.Fatalf("type = %q, want header", env.Type)
	}
	var h AttestedL2Header
	if err := json.Unmarshal(env.Value, &h); err != nil {
		t.Fatalf("decode header: %v", err)
	}

	// Spot-check the tricky representations: hex quantities, optional fields, and the
	// number-array proof matrix that must never serialize as base64.
	if h.L2Header.Number != 42 || h.L2Header.Difficulty.Uint64() != 0 || h.L2Header.BaseFeePerGas.Uint64() != 9 {
		t.Fatalf("scalar decode mismatch")
	}
	if h.L2Header.BlobGasUsed == nil || *h.L2Header.BlobGasUsed != 11 {
		t.Fatalf("blob_gas_used decode mismatch")
	}
	if len(h.RouterProof.Proof) != 2 || h.RouterProof.Proof[0][1] != 2 {
		t.Fatalf("router_proof decode mismatch: %v", h.RouterProof.Proof)
	}

	got, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	semanticJSONEqual(t, got, []byte(want))
}

// The []byte→base64 gotcha at the struct level, not just the scalar unit test: an
// account proof must go out as arrays of numbers.
func TestAttestedL2Header_ProofIsNumberArrays(t *testing.T) {
	h := &AttestedL2Header{RouterProof: EvmAccountProof{Proof: byteMatrix{{222, 173}}}}
	raw, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !strings.Contains(string(raw), `"proof":[[222,173]]`) {
		t.Fatalf("router_proof not number arrays: %s", raw)
	}
}
