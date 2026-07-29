package l2rollup

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/ethclient"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// hexUint64 renders n as a 0x-prefixed hex quantity for JSON-RPC params.
func hexUint64(n uint64) string { return fmt.Sprintf("0x%x", n) }

// This file holds the DETERMINISTIC pieces of the OP-Stack header builder — the ones
// that must match the op-stack-verifier byte-for-byte and so are unit-tested here.
// The RPC assembly (eth_getProof / optimism_outputAtBlock / beacon-slot resolution)
// lives in the builder's BuildHeader and is validated end-to-end on the OP devnet.

// gameListElementSlot computes the storage slot of element `index` of the
// DisputeGameFactory game list, matching op-stack-verifier's
// dynamic_array_element_slot exactly: keccak256(gameListSlot) + index, as a 32-byte
// big-endian value, erroring on overflow past 32 bytes. `gameListSlot` is the
// Solidity slot of the dynamic array header.
func gameListElementSlot(gameListSlot ethcommon.Hash, index uint64) (ethcommon.Hash, error) {
	base := crypto.Keccak256(gameListSlot.Bytes()) // keccak256(slot word)
	sum := new(big.Int).Add(new(big.Int).SetBytes(base), new(big.Int).SetUint64(index))
	if sum.BitLen() > 256 {
		return ethcommon.Hash{}, fmt.Errorf("game-list slot arithmetic overflow")
	}
	var out ethcommon.Hash
	sum.FillBytes(out[:]) // left-padded big-endian 32 bytes
	return out, nil
}

// addressFromStorageValue extracts a 20-byte address from a 32-byte storage word
// (right-aligned, as Solidity stores an address), matching the verifier's
// address_from_storage_value.
func addressFromStorageValue(word []byte) (ethcommon.Address, error) {
	if len(word) > 32 {
		return ethcommon.Address{}, fmt.Errorf("storage value exceeds 32 bytes")
	}
	// Left-pad to 32, take the low 20 bytes.
	var buf [32]byte
	copy(buf[32-len(word):], word)
	var addr ethcommon.Address
	copy(addr[:], buf[12:])
	return addr, nil
}

// toCanonicalHeader maps a go-ethereum execution header to the CanonicalEvmHeader the
// L2 client hashes. Optional Cancun/Prague fields are carried through as present so
// the L2-side fork validation (validate_for_fork) and RLP hashing reproduce the same
// block hash the output root commits to.
func toCanonicalHeader(h *types.Header) CanonicalEvmHeader {
	ch := CanonicalEvmHeader{
		ParentHash:       h.ParentHash.Bytes(),
		OmmersHash:       h.UncleHash.Bytes(),
		Beneficiary:      h.Coinbase.Bytes(),
		StateRoot:        h.Root.Bytes(),
		TransactionsRoot: h.TxHash.Bytes(),
		ReceiptsRoot:     h.ReceiptHash.Bytes(),
		LogsBloom:        h.Bloom.Bytes(),
		Difficulty:       u256FromBig(h.Difficulty),
		Number:           h.Number.Uint64(),
		GasLimit:         h.GasLimit,
		GasUsed:          h.GasUsed,
		Timestamp:        h.Time,
		ExtraData:        h.Extra,
		MixHash:          h.MixDigest.Bytes(),
		Nonce:            h.Nonce[:],
	}
	if h.BaseFee != nil {
		bf := u256FromBig(h.BaseFee)
		ch.BaseFeePerGas = &bf
	}
	if h.WithdrawalsHash != nil {
		wr := hexBytes(h.WithdrawalsHash.Bytes())
		ch.WithdrawalsRoot = &wr
	}
	if h.BlobGasUsed != nil {
		v := *h.BlobGasUsed
		ch.BlobGasUsed = &v
	}
	if h.ExcessBlobGas != nil {
		v := *h.ExcessBlobGas
		ch.ExcessBlobGas = &v
	}
	if h.ParentBeaconRoot != nil {
		pb := hexBytes(h.ParentBeaconRoot.Bytes())
		ch.ParentBeaconBlockRoot = &pb
	}
	if h.RequestsHash != nil {
		rh := hexBytes(h.RequestsHash.Bytes())
		ch.RequestsHash = &rh
	}
	return ch
}

// u256FromBig builds the wire u256 from a big.Int (nil → 0).
func u256FromBig(v *big.Int) u256 {
	var u u256
	if v != nil {
		u.Int.Set(v)
	}
	return u
}

// factoryGameCount reads DisputeGameFactory.gameCount() at a specific L1 block. The
// builder uses it to tell "this game does not exist yet at the block we can prove
// against" apart from "this game-list slot is empty", which look identical in a raw
// storage proof but mean very different things.
func factoryGameCount(ctx context.Context, l1 *ethclient.Client, factory ethcommon.Address, block *big.Int) (uint64, error) {
	selector := crypto.Keccak256([]byte("gameCount()"))[:4]
	out, err := l1.CallContract(ctx, ethereum.CallMsg{To: &factory, Data: selector}, block)
	if err != nil {
		return 0, fmt.Errorf("eth_call gameCount() at block %s: %w", block, err)
	}
	return new(big.Int).SetBytes(out).Uint64(), nil
}
