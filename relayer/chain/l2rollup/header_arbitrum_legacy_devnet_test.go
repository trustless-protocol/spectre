package l2rollup

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestArbLegacyBuildHeader_Devnet assembles a real ArbitrumLegacyHeader against
// canonical Arbitrum Sepolia (legacy Nitro). It is skipped unless ARB_DEVNET_L1_RPC is
// set. It calls buildLegacyHeaderFor directly (no live attestor) with an env-provided
// node number + its committed L2 block — the builder's internal self-checks (node
// active + confirmData == keccak(blockHash||sendRoot)) are the assertion.
//
// Defaults target the canonical Sepolia RollupCore 0xd808… / slots 0x75,0x76. Provide
// ARB_DEVNET_NODE + ARB_DEVNET_L2_BLOCK from the attestor (node → L2 block mapping).
//
//	Run: ARB_DEVNET_L1_RPC=https://sepolia.drpc.org \
//	     ARB_DEVNET_L2_RPC=https://sepolia-rollup.arbitrum.io/rpc \
//	     ARB_DEVNET_NODE=10764 ARB_DEVNET_L2_BLOCK=<committed-l2-block> \
//	     go test -run TestArbLegacyBuildHeader_Devnet -v ./chain/l2rollup/
func TestArbLegacyBuildHeader_Devnet(t *testing.T) {
	l1URL := os.Getenv("ARB_DEVNET_L1_RPC")
	if l1URL == "" {
		t.Skip("ARB_DEVNET_L1_RPC not set; skipping live Arbitrum Sepolia legacy test")
	}
	node := envUint(t, "ARB_DEVNET_NODE")
	l2Block := envUint(t, "ARB_DEVNET_L2_BLOCK")
	if node == 0 || l2Block == 0 {
		t.Skip("ARB_DEVNET_NODE and ARB_DEVNET_L2_BLOCK required (from the attestor's node→L2-block mapping)")
	}
	ctx := context.Background()

	l1, err := ethclient.Dial(l1URL)
	if err != nil {
		t.Fatalf("dial L1: %v", err)
	}
	l2, err := ethclient.Dial(os.Getenv("ARB_DEVNET_L2_RPC"))
	if err != nil {
		t.Fatalf("dial L2: %v", err)
	}
	l1Head, err := l1.BlockNumber(ctx)
	if err != nil {
		t.Fatalf("L1 head: %v", err)
	}

	profile := ArbLegacyProfile{
		RollupCore:            envAddr("ARB_DEVNET_ROLLUP", "0xd80810638dbdf9081b72c1b33c65375e807281c8"),
		L2Router:              envAddr("ARB_DEVNET_ROUTER", "0x645280885749dC97Ea461DE280Eb3273C91D36Df"),
		L1ClientID:            "08-wasm-0",
		NodeLifecycleSlot:     ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000075"),
		NodesMappingSlot:      ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000076"),
		LatestConfirmedOffset: 0,
		FirstUnresolvedOffset: 8,
		LatestCreatedOffset:   16,
		ConfirmDataOffset:     2,
	}
	builder := NewArbitrumLegacyHeaderBuilder(l1, l2, fakeCosmosReader{slot: 1, block: l1Head}, nil, "arbitrum-sepolia", false, profile)

	h, err := builder.buildLegacyHeaderFor(ctx, node, l2Block)
	if err != nil {
		// A confirmData / node-active mismatch surfaces here — the exact bindings verify_legacy checks.
		t.Fatalf("buildLegacyHeaderFor(node=%d, l2Block=%d, l1Block=%d): %v", node, l2Block, l1Head, err)
	}

	if len(h.RollupProof.Proof) == 0 {
		t.Error("rollup account proof is empty")
	}
	if len(h.NodeLifecycleProof.Proof) == 0 || len(h.ConfirmDataProof.Proof) == 0 {
		t.Error("node lifecycle / confirmData storage proof is empty")
	}
	if h.NodeNumber != node {
		t.Errorf("node number = %d, want %d", h.NodeNumber, node)
	}
	if len(h.SendRoot) != 32 {
		t.Errorf("send root width = %d, want 32", len(h.SendRoot))
	}

	raw, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	// Double envelope: {"type":"header","value":{"type":"legacy_nitro","value":{...}}}.
	var outer struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}
	if err := json.Unmarshal(raw, &outer); err != nil || outer.Type != "header" {
		t.Fatalf("outer envelope wrong: type=%q err=%v", outer.Type, err)
	}
	var inner struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(outer.Value, &inner); err != nil || inner.Type != "legacy_nitro" {
		t.Fatalf("inner envelope wrong: type=%q err=%v", inner.Type, err)
	}
	t.Logf("assembled ArbitrumLegacyHeader: node=%d l2Block=%d l1Block=%d bytes=%d", node, l2Block, l1Head, len(raw))
}

func envUint(t *testing.T, key string) uint64 {
	t.Helper()
	s := os.Getenv(key)
	if s == "" {
		return 0
	}
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		t.Fatalf("%s=%q is not a uint64: %v", key, s, err)
	}
	return n
}

func envAddr(key, fallback string) ethcommon.Address {
	if s := os.Getenv(key); s != "" {
		return ethcommon.HexToAddress(s)
	}
	return ethcommon.HexToAddress(fallback)
}
