package avalanche

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"relayer/chain/l2rollup"
)

// The fixture is the live-captured, avalanchego-verified Fuji aggregate shared
// with packages/avalanche-warp — pinning the Go builder/splitter to the same
// bytes the wasm client verifies.
type warpFixture struct {
	NetworkID       uint32 `json:"network_id"`
	SourceChainID   string `json:"source_chain_id"`
	BlockHash       string `json:"block_hash"`
	UnsignedMessage string `json:"unsigned_message"`
	SignedMessage   string `json:"signed_message"`
	BitSet          string `json:"bit_set"`
	Signature       string `json:"signature"`
}

func loadWarpFixture(t *testing.T) warpFixture {
	t.Helper()
	raw, err := os.ReadFile("testdata/warp_fuji_fixture.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	var fix warpFixture
	if err := json.Unmarshal(raw, &fix); err != nil {
		t.Fatalf("decode fixture: %v", err)
	}
	return fix
}

func unhex(t *testing.T, s string) []byte {
	t.Helper()
	raw, err := hex.DecodeString(strings.TrimPrefix(s, "0x"))
	if err != nil {
		t.Fatalf("decode %q: %v", s, err)
	}
	return raw
}

func TestBuildBlockHashWarpMessageMatchesTheFixture(t *testing.T) {
	fix := loadWarpFixture(t)
	var chainID, blockHash [32]byte
	copy(chainID[:], unhex(t, fix.SourceChainID))
	copy(blockHash[:], unhex(t, fix.BlockHash))
	built := buildBlockHashWarpMessage(fix.NetworkID, chainID, blockHash)
	if !bytes.Equal(built, unhex(t, fix.UnsignedMessage)) {
		t.Fatalf("built message %x, want %x", built, unhex(t, fix.UnsignedMessage))
	}
}

func TestSplitSignedWarpMessageMatchesTheFixture(t *testing.T) {
	fix := loadWarpFixture(t)
	unsigned, bitSet, signature, err := splitSignedWarpMessage(unhex(t, fix.SignedMessage))
	if err != nil {
		t.Fatalf("split: %v", err)
	}
	if !bytes.Equal(unsigned, unhex(t, fix.UnsignedMessage)) {
		t.Fatal("unsigned prefix mismatch")
	}
	if !bytes.Equal(bitSet, unhex(t, fix.BitSet)) {
		t.Fatal("bitset mismatch")
	}
	if !bytes.Equal(signature[:], unhex(t, fix.Signature)) {
		t.Fatal("signature mismatch")
	}

	// Truncation and wrong type ids fail loud.
	if _, _, _, err := splitSignedWarpMessage(unhex(t, fix.SignedMessage)[:50]); err == nil {
		t.Fatal("expected truncation error")
	}
	tampered := unhex(t, fix.SignedMessage)
	tampered[len(unsigned)] = 1 // signature type id
	if _, _, _, err := splitSignedWarpMessage(tampered); err == nil {
		t.Fatal("expected signature-type error")
	}
}

// The wire shape the wasm client parses: bitset/signature as base64 (CosmWasm
// Binary), the header under "header", inside the {type, value} envelope.
func TestWarpSignedHeaderEncodesTheClientEnvelope(t *testing.T) {
	fix := loadWarpFixture(t)
	msg := &WarpSignedHeader{
		Header:       l2rollup.CanonicalEvmHeader{},
		SignerBitSet: unhex(t, fix.BitSet),
		Signature:    unhex(t, fix.Signature),
		RouterProof:  l2rollup.EvmAccountProof{Proof: [][]byte{{0x01}}},
	}
	encoded, err := msg.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var envelope struct {
		Type  string `json:"type"`
		Value struct {
			SignerBitSet []byte          `json:"signer_bit_set"`
			Signature    []byte          `json:"signature"`
			Header       json.RawMessage `json:"header"`
			RouterProof  json.RawMessage `json:"router_proof"`
		} `json:"value"`
	}
	if err := json.Unmarshal(encoded, &envelope); err != nil {
		t.Fatalf("reparse: %v", err)
	}
	if envelope.Type != "header" {
		t.Fatalf("envelope type %q", envelope.Type)
	}
	if !bytes.Equal(envelope.Value.SignerBitSet, unhex(t, fix.BitSet)) || !bytes.Equal(envelope.Value.Signature, unhex(t, fix.Signature)) {
		t.Fatal("base64 binary round-trip mismatch")
	}
}
