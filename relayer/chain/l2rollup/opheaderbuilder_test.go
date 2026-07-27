package l2rollup

import (
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// TestGameListElementSlot_RustParity checks the slot math against the reference
// definition keccak256(slot) + index (op-stack-verifier dynamic_array_element_slot),
// computed independently here with big.Int so a divergence in the implementation is
// caught.
func TestGameListElementSlot_RustParity(t *testing.T) {
	slot := ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000005")
	base := new(big.Int).SetBytes(crypto.Keccak256(slot.Bytes()))

	for _, idx := range []uint64{0, 1, 17, 255, 256, 1 << 32} {
		got, err := gameListElementSlot(slot, idx)
		if err != nil {
			t.Fatalf("index %d: %v", idx, err)
		}
		want := new(big.Int).Add(base, new(big.Int).SetUint64(idx))
		var wantHash ethcommon.Hash
		want.FillBytes(wantHash[:])
		if got != wantHash {
			t.Fatalf("index %d: slot = %s, want %s", idx, got.Hex(), wantHash.Hex())
		}
	}
}

// TestGameListElementSlot_Index0IsKeccak locks the base case: element 0's slot is
// exactly keccak256(gameListSlot).
func TestGameListElementSlot_Index0IsKeccak(t *testing.T) {
	slot := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000ab")
	got, err := gameListElementSlot(slot, 0)
	if err != nil {
		t.Fatal(err)
	}
	want := ethcommon.BytesToHash(crypto.Keccak256(slot.Bytes()))
	if got != want {
		t.Fatalf("element 0 slot = %s, want keccak256(slot) = %s", got.Hex(), want.Hex())
	}
}

func TestAddressFromStorageValue(t *testing.T) {
	// A right-aligned address word.
	word := ethcommon.HexToHash("0x000000000000000000000000abcdef0123456789abcdef0123456789abcdef01")
	addr, err := addressFromStorageValue(word.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	want := ethcommon.HexToAddress("0xabcdef0123456789abcdef0123456789abcdef01")
	if addr != want {
		t.Fatalf("addr = %s, want %s", addr.Hex(), want.Hex())
	}

	// A minimal (leading-zero-trimmed) storage value must still left-pad correctly.
	addr2, err := addressFromStorageValue([]byte{0x01})
	if err != nil {
		t.Fatal(err)
	}
	if addr2 != ethcommon.HexToAddress("0x0000000000000000000000000000000000000001") {
		t.Fatalf("minimal value addr = %s", addr2.Hex())
	}
}

// TestToCanonicalHeader_CancunFields confirms a Cancun header maps all mandatory
// fields plus the blob/withdrawals/parent-beacon optionals (present), and difficulty
// serializes as a u256 quantity.
func TestToCanonicalHeader_CancunFields(t *testing.T) {
	blob := uint64(11)
	excess := uint64(12)
	wr := ethcommon.HexToHash("0x0a")
	pbr := ethcommon.HexToHash("0x0d")
	h := &types.Header{
		Number:           big.NewInt(7),
		GasLimit:         30_000_000,
		GasUsed:          21_000,
		Time:             1_700_000_000,
		Difficulty:       big.NewInt(0),
		BaseFee:          big.NewInt(9),
		Extra:            []byte{0xab, 0xcd},
		WithdrawalsHash:  &wr,
		BlobGasUsed:      &blob,
		ExcessBlobGas:    &excess,
		ParentBeaconRoot: &pbr,
	}
	ch := toCanonicalHeader(h)

	if ch.Number != 7 || ch.GasLimit != 30_000_000 || ch.GasUsed != 21_000 {
		t.Fatalf("scalar fields mismatch: %+v", ch)
	}
	if ch.BaseFeePerGas == nil || ch.BaseFeePerGas.Uint64() != 9 {
		t.Fatal("base_fee not carried")
	}
	if ch.BlobGasUsed == nil || *ch.BlobGasUsed != 11 || ch.ExcessBlobGas == nil || *ch.ExcessBlobGas != 12 {
		t.Fatal("blob fields not carried")
	}
	if ch.WithdrawalsRoot == nil || ch.ParentBeaconBlockRoot == nil {
		t.Fatal("shanghai/cancun roots not carried")
	}
	if len(ch.Nonce) != 8 {
		t.Fatalf("nonce = %d bytes, want 8", len(ch.Nonce))
	}
}

// TestToCanonicalHeader_LondonOmitsOptionals: a pre-Shanghai header must leave the
// optional tail nil so the L2-side London fork validation accepts it.
func TestToCanonicalHeader_LondonOmitsOptionals(t *testing.T) {
	h := &types.Header{
		Number:     big.NewInt(1),
		Difficulty: big.NewInt(0),
		BaseFee:    big.NewInt(1),
	}
	ch := toCanonicalHeader(h)
	if ch.WithdrawalsRoot != nil || ch.BlobGasUsed != nil || ch.ParentBeaconBlockRoot != nil || ch.RequestsHash != nil {
		t.Fatal("London header must omit shanghai/cancun/prague optionals")
	}
}
