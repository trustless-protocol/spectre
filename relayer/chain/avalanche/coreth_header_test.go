package avalanche

import (
	"encoding/json"
	"os"
	"testing"

	"relayer/chain/l2rollup"

	ethcommon "github.com/ethereum/go-ethereum/common"
)

// The fixture is a real Fuji C-Chain header (#58513648, post-Helicon: full
// optional tail), fetched 2026-09-21 from https://api.avax-test.network. The
// same header backs the Rust test
// canonical_header_tests::hashes_a_real_fuji_coreth_header_to_its_reported_block_hash,
// so the Go and Rust coreth RLP encoders are pinned to one another through it.
func fujiHeader(t *testing.T) corethRPCHeader {
	t.Helper()
	raw, err := os.ReadFile("testdata/coreth_fuji_header.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var header corethRPCHeader
	if err := json.Unmarshal(raw, &header); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return header
}

func TestCorethHeaderHashMatchesReportedFujiBlockHash(t *testing.T) {
	header := fujiHeader(t)
	computed, err := corethHeaderHash(header)
	if err != nil {
		t.Fatalf("corethHeaderHash: %v", err)
	}
	want := ethcommon.HexToHash("0x68c074dadf040a70a81825dfdce0581bc604beca6cf8778a6f99c4a5acfb74e5")
	if computed != want {
		t.Fatalf("computed %s, want %s", computed, want)
	}
	if header.Hash != want {
		t.Fatalf("fixture reported hash drifted: %s", header.Hash)
	}
}

// A Granite-era header (Helicon fields absent) must still hash: the cascade
// stops at min_delay_excess and never pads trailing absents.
func TestCorethHeaderHashHandlesShorterOptionalTail(t *testing.T) {
	header := fujiHeader(t)
	header.TargetExponent = nil
	header.MinPriceExponent = nil
	header.SettledHeight = nil
	header.SettledGasUnix = nil
	header.SettledGasNumerator = nil
	header.SettledExcess = nil
	full, err := corethHeaderHash(fujiHeader(t))
	if err != nil {
		t.Fatalf("full tail: %v", err)
	}
	short, err := corethHeaderHash(header)
	if err != nil {
		t.Fatalf("short tail: %v", err)
	}
	if full == short {
		t.Fatal("shortening the optional tail must change the hash")
	}
}

// The wire mapping must carry every coreth field so the Rust side reproduces
// the same RLP; spot-check presence and JSON names against canonical_header.rs.
func TestCorethWireHeaderCarriesCorethFields(t *testing.T) {
	wire := corethWireHeader(fujiHeader(t))
	if wire.ExtDataHash == nil || wire.TimeMilliseconds == nil || wire.SettledExcess == nil {
		t.Fatal("coreth fields missing from wire header")
	}
	if wire.WithdrawalsRoot != nil || wire.RequestsHash != nil {
		t.Fatal("coreth headers must not carry geth withdrawals/requests fields")
	}
	raw, err := json.Marshal(wire)
	if err != nil {
		t.Fatalf("marshal wire header: %v", err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatalf("reparse wire header: %v", err)
	}
	for _, key := range []string{"ext_data_hash", "ext_data_gas_used", "block_gas_cost", "time_milliseconds", "min_delay_excess", "settled_excess"} {
		if _, ok := keys[key]; !ok {
			t.Fatalf("wire JSON missing %q", key)
		}
	}
	if _, ok := keys["withdrawals_root"]; ok {
		t.Fatal("wire JSON must omit absent geth optionals")
	}
}

// A geth-family wire header must omit every coreth key, or the deployed
// OP/Base/Arbitrum wasm clients (deny_unknown_fields) would reject it.
func TestGethWireHeaderOmitsCorethFields(t *testing.T) {
	raw, err := json.Marshal(l2rollup.CanonicalEvmHeader{})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var keys map[string]json.RawMessage
	if err := json.Unmarshal(raw, &keys); err != nil {
		t.Fatalf("reparse: %v", err)
	}
	for _, key := range []string{"ext_data_hash", "ext_data_gas_used", "block_gas_cost", "time_milliseconds", "min_delay_excess", "target_exponent", "min_price_exponent", "settled_height", "settled_gas_unix", "settled_gas_numerator", "settled_excess"} {
		if _, ok := keys[key]; ok {
			t.Fatalf("geth wire JSON must omit %q", key)
		}
	}
}
