package l2rollup

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

// opExampleJSON is Dũng's canonical serializer-shape example from PR #245 (the OP
// `{"type":"header","value":{...}}` envelope). The 256-byte logs_bloom is built
// from a 512-zero-hex-digit string. It is a wire-shape fixture (dummy proofs), used
// only to lock the Go encoder's field names + value representations.
func opExampleJSON() string {
	bloom := "0x" + strings.Repeat("0", 512)
	return `{
  "type": "header",
  "value": {
    "beacon_slot": 123,
    "l1_state_root": "0x0000000000000000000000000000000000000000000000000000000000000010",
    "factory_proof": {"proof": [[248, 81]]},
    "game_index": 17,
    "game_proof": {
      "key": "0x0000000000000000000000000000000000000000000000000000000000000012",
      "value": [1],
      "proof": [[249, 0, 1]]
    },
    "game_account_proof": {"proof": [[170]]},
    "game_runtime": [222, 173],
    "output_root_proof": {
      "version": "0x0000000000000000000000000000000000000000000000000000000000000000",
      "state_root": "0x0000000000000000000000000000000000000000000000000000000000000004",
      "message_passer_storage_root": "0x0000000000000000000000000000000000000000000000000000000000000014",
      "latest_blockhash": "0x0000000000000000000000000000000000000000000000000000000000000015"
    },
    "l2_header": {
      "parent_hash": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "ommers_hash": "0x0000000000000000000000000000000000000000000000000000000000000002",
      "beneficiary": "0x0000000000000000000000000000000000000003",
      "state_root": "0x0000000000000000000000000000000000000000000000000000000000000004",
      "transactions_root": "0x0000000000000000000000000000000000000000000000000000000000000005",
      "receipts_root": "0x0000000000000000000000000000000000000000000000000000000000000006",
      "logs_bloom": "` + bloom + `",
      "difficulty": "0x0",
      "number": 7,
      "gas_limit": 30000000,
      "gas_used": 21000,
      "timestamp": 1700000000,
      "extra_data": "0xabcd",
      "mix_hash": "0x0000000000000000000000000000000000000000000000000000000000000008",
      "nonce": "0x0000000000000000",
      "base_fee_per_gas": "0x9",
      "withdrawals_root": "0x000000000000000000000000000000000000000000000000000000000000000a",
      "blob_gas_used": 11,
      "excess_blob_gas": 12,
      "parent_beacon_block_root": "0x000000000000000000000000000000000000000000000000000000000000000d",
      "requests_hash": "0x000000000000000000000000000000000000000000000000000000000000000e"
    },
    "router_proof": {"proof": [[187, 204]]}
  }
}`
}

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

// TestOpStackHeader_RoundTrip decodes Dũng's canonical example into the Go types
// and re-encodes it, asserting the wire form is byte-semantically identical. This
// is the freeze test for the OP/Base encoding contract.
func TestOpStackHeader_RoundTrip(t *testing.T) {
	want := opExampleJSON()

	var env clientMessageEnvelope
	if err := json.Unmarshal([]byte(want), &env); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	if env.Type != "header" {
		t.Fatalf("type = %q, want header", env.Type)
	}
	var h OpStackHeader
	if err := json.Unmarshal(env.Value, &h); err != nil {
		t.Fatalf("decode OP header: %v", err)
	}

	// Spot-check a few decoded values across the tricky representations.
	if h.GameRuntime[0] != 222 || h.GameRuntime[1] != 173 {
		t.Fatalf("game_runtime = %v, want [222 173]", h.GameRuntime)
	}
	if len(h.GameProof.Value) != 1 || h.GameProof.Value[0] != 1 {
		t.Fatalf("game_proof.value = %v, want [1]", h.GameProof.Value)
	}
	if h.L2Header.Difficulty.Uint64() != 0 || h.L2Header.BaseFeePerGas.Uint64() != 9 {
		t.Fatalf("difficulty/base_fee decode mismatch")
	}
	if h.L2Header.BlobGasUsed == nil || *h.L2Header.BlobGasUsed != 11 {
		t.Fatalf("blob_gas_used decode mismatch")
	}

	got, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	semanticJSONEqual(t, got, []byte(want))
}

// TestOpStackHeader_GameRuntimeNotBase64 guards the whole header path against the
// []byte→base64 gotcha at the struct level (not just the scalar unit test).
func TestOpStackHeader_GameRuntimeNotBase64(t *testing.T) {
	h := &OpStackHeader{GameRuntime: byteList{222, 173}}
	raw, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if !strings.Contains(string(raw), `"game_runtime":[222,173]`) {
		t.Fatalf("game_runtime not a number array: %s", raw)
	}
}

// TestArbitrumBoldHeader_RoundTrip locks the Arbitrum encoding contract, including the
// nested assertion (global_state arrays + snake_case machine_status) and the
// l1_state_root / l2_header fields the code carries.
func TestArbitrumBoldHeader_RoundTrip(t *testing.T) {
	// Arbitrum's verifier Header is a tagged enum, so the ClientMessage value is itself
	// tagged: {"type":"header","value":{"type":"bold_v2","value":<header>}}.
	want := `{
  "type": "header",
  "value": {
    "type": "bold_v2",
    "value": {
      "beacon_slot": 123,
      "l1_state_root": "0x0000000000000000000000000000000000000000000000000000000000000010",
      "rollup_proof": {"proof": [[248, 81]]},
      "assertion_hash": "0x0000000000000000000000000000000000000000000000000000000000000020",
      "assertion_proof": {
        "key": "0x0000000000000000000000000000000000000000000000000000000000000021",
        "value": [1],
        "proof": [[249, 0, 1]]
      },
      "assertion": {
        "parent_assertion_hash": "0x0000000000000000000000000000000000000000000000000000000000000030",
        "after_state": {
          "global_state": {
            "bytes32_vals": [
              "0x0000000000000000000000000000000000000000000000000000000000000031",
              "0x0000000000000000000000000000000000000000000000000000000000000032"
            ],
            "u64_vals": [123, 0]
          },
          "machine_status": "finished",
          "end_history_root": "0x0000000000000000000000000000000000000000000000000000000000000033"
        },
        "inbox_acc": "0x0000000000000000000000000000000000000000000000000000000000000034"
      },
      "l2_header": ` + minimalCancunL2Header() + `,
      "router_proof": {"proof": [[187, 204]]}
    }
  }
}`

	var outer clientMessageEnvelope
	if err := json.Unmarshal([]byte(want), &outer); err != nil {
		t.Fatalf("decode outer envelope: %v", err)
	}
	var inner clientMessageEnvelope
	if err := json.Unmarshal(outer.Value, &inner); err != nil {
		t.Fatalf("decode inner (tagged Arbitrum header) envelope: %v", err)
	}
	if inner.Type != "bold_v2" {
		t.Fatalf("inner tag = %q, want bold_v2", inner.Type)
	}
	var h ArbitrumBoldHeader
	if err := json.Unmarshal(inner.Value, &h); err != nil {
		t.Fatalf("decode Arbitrum header: %v", err)
	}
	if h.Assertion.AfterState.MachineStatus != MachineStatusFinished {
		t.Fatalf("machine_status = %q, want finished", h.Assertion.AfterState.MachineStatus)
	}
	if h.Assertion.AfterState.GlobalState.U64Vals != [2]uint64{123, 0} {
		t.Fatalf("u64_vals = %v", h.Assertion.AfterState.GlobalState.U64Vals)
	}

	got, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("re-encode: %v", err)
	}
	semanticJSONEqual(t, got, []byte(want))
}

// minimalCancunL2Header is a Cancun-fork canonical header JSON reused by the
// Arbitrum fixture (same CanonicalEvmHeader type as OP).
func minimalCancunL2Header() string {
	bloom := "0x" + strings.Repeat("0", 512)
	return `{
      "parent_hash": "0x0000000000000000000000000000000000000000000000000000000000000001",
      "ommers_hash": "0x0000000000000000000000000000000000000000000000000000000000000002",
      "beneficiary": "0x0000000000000000000000000000000000000003",
      "state_root": "0x0000000000000000000000000000000000000000000000000000000000000004",
      "transactions_root": "0x0000000000000000000000000000000000000000000000000000000000000005",
      "receipts_root": "0x0000000000000000000000000000000000000000000000000000000000000006",
      "logs_bloom": "` + bloom + `",
      "difficulty": "0x0",
      "number": 7,
      "gas_limit": 30000000,
      "gas_used": 21000,
      "timestamp": 1700000000,
      "extra_data": "0xabcd",
      "mix_hash": "0x0000000000000000000000000000000000000000000000000000000000000008",
      "nonce": "0x0000000000000000",
      "base_fee_per_gas": "0x9",
      "withdrawals_root": "0x000000000000000000000000000000000000000000000000000000000000000a",
      "blob_gas_used": 11,
      "excess_blob_gas": 12,
      "parent_beacon_block_root": "0x000000000000000000000000000000000000000000000000000000000000000d",
      "requests_hash": "0x000000000000000000000000000000000000000000000000000000000000000e"
    }`
}
