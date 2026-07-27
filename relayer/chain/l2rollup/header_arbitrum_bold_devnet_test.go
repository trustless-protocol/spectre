package l2rollup

import (
	"context"
	"encoding/json"
	"math/big"
	"os"
	"strconv"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TestArbBuildHeader_Devnet exercises the full RPC assembly against a live local
// Arbitrum (Nitro) devnet. It is skipped unless ARB_DEVNET_L1_RPC is set (the devnet is
// ephemeral, not part of CI). It confirms BuildHeader assembles a well-formed
// ArbitrumBoldHeader from a real RollupCore assertion + L2 header, and that the assertion
// commits exactly the L2 block hash of the packaged header — the binding verify()
// checks in arbitrum-verifier.
//
//	Run: ARB_DEVNET_L1_RPC=... ARB_DEVNET_L2_RPC=... ARB_DEVNET_ROLLUP=... \
//	     ARB_DEVNET_ROUTER=... ARB_DEVNET_ASSERTIONS_SLOT=0x... \
//	     go test -run TestArbBuildHeader_Devnet -v ./chain/l2rollup/
func TestArbBuildHeader_Devnet(t *testing.T) {
	l1URL := os.Getenv("ARB_DEVNET_L1_RPC")
	if l1URL == "" {
		t.Skip("ARB_DEVNET_L1_RPC not set; skipping live Arbitrum devnet integration test")
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
	l2Head, err := l2.BlockNumber(ctx)
	if err != nil {
		t.Fatalf("L2 head: %v", err)
	}

	// Public RPCs cap eth_getLogs block ranges; ARB_DEVNET_SCAN_RANGE tunes the
	// AssertionCreated scan window (0/unset → builder default).
	var scanRange uint64
	if s := os.Getenv("ARB_DEVNET_SCAN_RANGE"); s != "" {
		scanRange, _ = strconv.ParseUint(s, 10, 64)
	}
	profile := ArbBoldProfile{
		RollupCore:            ethcommon.HexToAddress(os.Getenv("ARB_DEVNET_ROLLUP")),
		AssertionsMappingSlot: ethcommon.HexToHash(os.Getenv("ARB_DEVNET_ASSERTIONS_SLOT")),
		L2Router:              ethcommon.HexToAddress(os.Getenv("ARB_DEVNET_ROUTER")),
		L1ClientID:            "08-wasm-0",
		AssertionScanRange:    scanRange,
	}
	builder := NewArbitrumBoldHeaderBuilder(l1, l2, fakeCosmosReader{slot: 1, block: l1Head}, profile)

	msg, committedHeight, err := builder.BuildHeader(ctx, l2Head)
	if err != nil {
		t.Fatalf("BuildHeader(l2Head=%d, l1Block=%d): %v", l2Head, l1Head, err)
	}
	if committedHeight == 0 || committedHeight > l2Head {
		t.Errorf("committed height = %d, want in (0, %d]", committedHeight, l2Head)
	}
	h := msg.(*ArbitrumBoldHeader)

	if len(h.RollupProof.Proof) == 0 {
		t.Error("rollup account proof is empty")
	}
	if len(h.AssertionProof.Proof) == 0 {
		t.Error("assertion storage proof is empty")
	}
	if len(h.AssertionHash) != 32 || ethcommon.BytesToHash(h.AssertionHash) == (ethcommon.Hash{}) {
		t.Errorf("assertion hash is zero/short: %x", h.AssertionHash)
	}

	// The assertion must commit exactly the L2 block hash of the header we packaged —
	// the binding verify() enforces (global_state.bytes32_vals[0] == l2_header.hash()).
	l2Header, err := l2.HeaderByNumber(ctx, new(big.Int).SetUint64(h.L2Header.Number))
	if err != nil {
		t.Fatalf("L2 header at %d: %v", h.L2Header.Number, err)
	}
	committed := ethcommon.BytesToHash(h.Assertion.AfterState.GlobalState.Bytes32Vals[0])
	if committed != l2Header.Hash() {
		t.Errorf("assertion committed block hash != l2 header hash:\n committed=%s\n l2hash   =%s", committed, l2Header.Hash())
	}

	raw, err := h.EncodeClientMessage()
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatalf("envelope not JSON: %v", err)
	}
	t.Logf("assembled ArbitrumBoldHeader: l2Height=%d assertion=%s l1Block=%d bytes=%d",
		h.L2Header.Number, ethcommon.BytesToHash(h.AssertionHash), l1Head, len(raw))
}
