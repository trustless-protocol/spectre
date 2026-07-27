package l2rollup

import (
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// These vectors are the canonical Arbitrum Sepolia legacy node 10764, cross-checked
// against arbitrum-verifier's canonical_sepolia_legacy_layout_matches_live_node_10764
// test AND verified live on-chain (RollupCore 0xd808…, slots 0x75/0x76).
const (
	sepoliaNodesSlot       = "0x0000000000000000000000000000000000000000000000000000000000000076"
	sepoliaNode10764Base   = "0xe4251833ee6ae54b9b2f5e13ed8882e17a2845215254f11f2b85a94b2f9b59b7"
	sepoliaNode10764CDSlot = "0xe4251833ee6ae54b9b2f5e13ed8882e17a2845215254f11f2b85a94b2f9b59b9"
	// packed {_latestConfirmed@0, _firstUnresolvedNode@8, _latestNodeCreated@16, staked@24}
	sepoliaLifecycleWord = "0x00000000003f5a360000000000002a0c0000000000002a0d0000000000002a0c"
)

func TestMappingSlotU64MatchesLiveSepolia(t *testing.T) {
	got := mappingSlotU64(10_764, ethcommon.HexToHash(sepoliaNodesSlot))
	if got != ethcommon.HexToHash(sepoliaNode10764Base) {
		t.Fatalf("mappingSlotU64(10764, 0x76) = %s, want %s", got, sepoliaNode10764Base)
	}
}

func TestAddStorageOffset(t *testing.T) {
	base := ethcommon.HexToHash(sepoliaNode10764Base)
	got, err := addStorageOffset(base, 2) // confirm_data_offset
	if err != nil {
		t.Fatalf("addStorageOffset: %v", err)
	}
	if got != ethcommon.HexToHash(sepoliaNode10764CDSlot) {
		t.Errorf("addStorageOffset(base, 2) = %s, want %s", got, sepoliaNode10764CDSlot)
	}

	// Overflow past 32 bytes must error.
	maxSlot := ethcommon.HexToHash("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	if _, err := addStorageOffset(maxSlot, 1); err == nil {
		t.Error("expected overflow error for max slot + 1")
	}
}

func TestPackedU64Lifecycle(t *testing.T) {
	word := ethcommon.FromHex(sepoliaLifecycleWord)
	cases := []struct {
		offset uint8
		want   uint64
	}{{0, 10_764}, {8, 10_765}, {16, 10_764}}
	for _, c := range cases {
		got, err := packedU64(word, c.offset)
		if err != nil {
			t.Fatalf("packedU64(offset=%d): %v", c.offset, err)
		}
		if got != c.want {
			t.Errorf("packedU64(offset=%d) = %d, want %d", c.offset, got, c.want)
		}
	}
	// Offset that would read past the low end of the word must error.
	if _, err := packedU64(word, 25); err == nil {
		t.Error("expected error for offset 25 (field exceeds word)")
	}
}

func TestPackedU64AcceptsMinimalValue(t *testing.T) {
	// eth_getProof returns minimal big-endian storage values; storageWord left-pads.
	// The lifecycle word above trims to its high byte 0x3f...; packedU64 must still
	// read the low-end fields correctly after padding.
	minimal := ethcommon.FromHex("0x3f5a360000000000002a0c0000000000002a0d0000000000002a0c")
	got, err := packedU64(minimal, 0)
	if err != nil {
		t.Fatalf("packedU64(minimal, 0): %v", err)
	}
	if got != 10_764 {
		t.Errorf("packedU64(minimal, 0) = %d, want 10764", got)
	}
}

func TestLegacyNodeActive(t *testing.T) {
	// Mirror arbitrum-verifier legacy_lifecycle_rejects_invalid_state_and_resolved_losers.
	active := []uint64{10, 12, 15}  // latest_confirmed / first_unresolved-edge / latest_created
	inactive := []uint64{9, 11, 16} // below confirmed, resolved loser, above created
	for _, n := range active {
		if !legacyNodeActive(n, 10, 12, 15) {
			t.Errorf("node %d should be active (confirmed=10, first=12, created=15)", n)
		}
	}
	for _, n := range inactive {
		if legacyNodeActive(n, 10, 12, 15) {
			t.Errorf("node %d should be inactive (confirmed=10, first=12, created=15)", n)
		}
	}
}

func TestStorageWord(t *testing.T) {
	w, err := storageWord(ethcommon.FromHex("0x2a0c"))
	if err != nil {
		t.Fatalf("storageWord: %v", err)
	}
	if w[30] != 0x2a || w[31] != 0x0c {
		t.Errorf("storageWord left-pad wrong: %x", w)
	}
	if _, err := storageWord(make([]byte, 33)); err == nil {
		t.Error("expected error for value wider than one word")
	}
}

func TestLegacyConfirmDataBinding(t *testing.T) {
	// confirmData = keccak256(blockHash || sendRoot) for node 10764 — the exact binding
	// buildLegacyHeaderFor self-checks, cross-referenced to the live confirmData value.
	blockHash := ethcommon.FromHex("0xfec7821ffcbba26c10f3e60964e2b7288262203c7dfa390285a61dc92d654ee5")
	sendRoot := ethcommon.FromHex("0xe51b1f0c5c8fcd2155d335c309a49c82807ecc19e159d3f60eeba1d3e63ecb77")
	got := crypto.Keccak256Hash(blockHash, sendRoot)
	want := ethcommon.HexToHash("0x9e3f9c983ceb4e5683856972d3302811aff89b0a9b5802a173aaa3e3b4700765")
	if got != want {
		t.Fatalf("confirmData = %s, want %s", got, want)
	}
}
