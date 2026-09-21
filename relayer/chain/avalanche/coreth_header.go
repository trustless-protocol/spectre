package avalanche

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"relayer/chain/l2rollup"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	gethrpc "github.com/ethereum/go-ethereum/rpc"
)

// This file is the coreth (Avalanche C-Chain) counterpart of l2rollup/evm_header.go.
// Coreth headers extend the geth header with a fixed ExtDataHash plus a
// cascading optional tail (ExtDataGasUsed .. SettledExcess), so they cannot be
// read through go-ethereum's types.Header — ethclient would silently drop the
// extra fields and Hash() would compute the wrong block hash. Instead the raw
// RPC JSON is decoded here and the block hash recomputed by mirroring coreth's
// generated RLP encoder (avalanchego graft/coreth
// plugin/evm/customtypes/gen_header_serializable_rlp.go): a tail field is
// emitted iff it or any later tail field is present, an absent field emitted
// that way encodes as the RLP empty string (0x80), and trailing absents are
// omitted. The Rust mirror is l2rollup.CanonicalEvmHeader::coreth_rlp_bytes in
// packages/l2-client/src/canonical_header.rs; the shared fixture test keeps the
// two byte-identical.
//
// The recomputed hash is checked against the RPC-reported one, so a future
// coreth header change fails loudly here instead of producing an update the
// attestor (which compares against its own replica's hash) would reject with a
// less specific divergence error.

// corethRPCHeader is the JSON shape coreth's eth_getBlockByNumber returns,
// restricted to header fields. JSON names come from coreth's HeaderSerializable
// gencodec tags.
type corethRPCHeader struct {
	ParentHash       ethcommon.Hash    `json:"parentHash"`
	Sha3Uncles       ethcommon.Hash    `json:"sha3Uncles"`
	Miner            ethcommon.Address `json:"miner"`
	StateRoot        ethcommon.Hash    `json:"stateRoot"`
	TransactionsRoot ethcommon.Hash    `json:"transactionsRoot"`
	ReceiptsRoot     ethcommon.Hash    `json:"receiptsRoot"`
	LogsBloom        hexutil.Bytes     `json:"logsBloom"`
	Difficulty       *hexutil.Big      `json:"difficulty"`
	Number           hexutil.Uint64    `json:"number"`
	GasLimit         hexutil.Uint64    `json:"gasLimit"`
	GasUsed          hexutil.Uint64    `json:"gasUsed"`
	Timestamp        hexutil.Uint64    `json:"timestamp"`
	ExtraData        hexutil.Bytes     `json:"extraData"`
	MixHash          ethcommon.Hash    `json:"mixHash"`
	Nonce            hexutil.Bytes     `json:"nonce"`
	ExtDataHash      ethcommon.Hash    `json:"extDataHash"`

	BaseFeePerGas         *hexutil.Big    `json:"baseFeePerGas"`
	ExtDataGasUsed        *hexutil.Big    `json:"extDataGasUsed"`
	BlockGasCost          *hexutil.Big    `json:"blockGasCost"`
	BlobGasUsed           *hexutil.Uint64 `json:"blobGasUsed"`
	ExcessBlobGas         *hexutil.Uint64 `json:"excessBlobGas"`
	ParentBeaconBlockRoot *ethcommon.Hash `json:"parentBeaconBlockRoot"`
	TimeMilliseconds      *hexutil.Uint64 `json:"timestampMilliseconds"`
	MinDelayExcess        *hexutil.Uint64 `json:"minDelayExcess"`
	TargetExponent        *hexutil.Uint64 `json:"targetExponent"`
	MinPriceExponent      *hexutil.Uint64 `json:"minPriceExponent"`
	SettledHeight         *hexutil.Uint64 `json:"settledHeight"`
	SettledGasUnix        *hexutil.Uint64 `json:"settledGasUnix"`
	SettledGasNumerator   *hexutil.Uint64 `json:"settledGasNumerator"`
	SettledExcess         *hexutil.Uint64 `json:"settledExcess"`

	// Hash is the node-reported block hash the recomputation is checked against.
	Hash ethcommon.Hash `json:"hash"`
}

// corethProofQueryHeight maps a header to the height its state root commits.
//
// Measured on Fuji (post-Helicon, ACP-194 streaming asynchronous execution):
// header(N).stateRoot equals the post-execution state of settledHeight(N)
// (N-5 observed), and eth_getProof(Q) anchors to the post-execution state of
// Q — probing Q ∈ {N-6..N} matches only at Q = settledHeight(N). So every
// proof paired with a header must be fetched at the header's settled height.
// Pre-Helicon headers carry no settled fields and keep the synchronous rule
// (proofs at the header's own height), which makes this future-proof in both
// directions: it reads the pairing from the header instead of the network's
// upgrade calendar.
func corethProofQueryHeight(h corethRPCHeader) uint64 {
	if h.SettledHeight != nil {
		return uint64(*h.SettledHeight)
	}
	return uint64(h.Number)
}

// corethSettledQueryHeight resolves the proof-query height for the header at
// `height` with one light RPC read (only the fields the mapping needs).
func corethSettledQueryHeight(ctx context.Context, client *gethrpc.Client, height uint64) (uint64, error) {
	var header *struct {
		Number        hexutil.Uint64  `json:"number"`
		SettledHeight *hexutil.Uint64 `json:"settledHeight"`
	}
	if err := client.CallContext(ctx, &header, "eth_getBlockByNumber", hexutil.EncodeUint64(height), false); err != nil {
		return 0, fmt.Errorf("l2rollup: coreth header at %d: %w", height, err)
	}
	if header == nil {
		return 0, fmt.Errorf("l2rollup: coreth header at %d: block not found", height)
	}
	if header.SettledHeight != nil {
		return uint64(*header.SettledHeight), nil
	}
	return uint64(header.Number), nil
}

// NewSettledProofHeightResolver returns the proof-height mapping installed on
// the shared Source (WithProofHeightResolver): a client-update height resolves
// to its header's settled height — identity for pre-Helicon headers. This
// keeps the asynchronous-execution rule in the Avalanche provider; the shared
// rollup machinery holds no chain-specific knowledge.
func NewSettledProofHeightResolver(client *gethrpc.Client) func(context.Context, uint64) (uint64, error) {
	return func(ctx context.Context, height uint64) (uint64, error) {
		return corethSettledQueryHeight(ctx, client, height)
	}
}

// readCorethHeader fetches the coreth header at height over raw RPC and returns
// the wire header, its verified block hash, its state root, and the height any
// paired proof must be queried at (the settled height under asynchronous
// execution; see corethProofQueryHeight).
func readCorethHeader(ctx context.Context, client *gethrpc.Client, height *big.Int) (l2rollup.CanonicalEvmHeader, ethcommon.Hash, ethcommon.Hash, uint64, error) {
	var raw json.RawMessage
	if err := client.CallContext(ctx, &raw, "eth_getBlockByNumber", hexutil.EncodeBig(height), false); err != nil {
		return l2rollup.CanonicalEvmHeader{}, ethcommon.Hash{}, ethcommon.Hash{}, 0, fmt.Errorf("l2rollup: coreth header at %s: %w", height, err)
	}
	if len(raw) == 0 || string(raw) == "null" {
		return l2rollup.CanonicalEvmHeader{}, ethcommon.Hash{}, ethcommon.Hash{}, 0, fmt.Errorf("l2rollup: coreth header at %s: block not found", height)
	}
	var header corethRPCHeader
	if err := json.Unmarshal(raw, &header); err != nil {
		return l2rollup.CanonicalEvmHeader{}, ethcommon.Hash{}, ethcommon.Hash{}, 0, fmt.Errorf("l2rollup: decode coreth header at %s: %w", height, err)
	}
	computed, err := corethHeaderHash(header)
	if err != nil {
		return l2rollup.CanonicalEvmHeader{}, ethcommon.Hash{}, ethcommon.Hash{}, 0, err
	}
	if computed != header.Hash {
		return l2rollup.CanonicalEvmHeader{}, ethcommon.Hash{}, ethcommon.Hash{}, 0, fmt.Errorf(
			"l2rollup: coreth header %d: recomputed hash %s does not match reported %s (unknown header field set?)",
			uint64(header.Number), computed, header.Hash)
	}
	return corethWireHeader(header), computed, header.StateRoot, corethProofQueryHeight(header), nil
}

// corethWireHeader maps the RPC header to the wire l2rollup.CanonicalEvmHeader. Optional
// fields are carried through exactly as present or absent, because the Rust
// side reproduces the RLP cascade from that presence.
func corethWireHeader(h corethRPCHeader) l2rollup.CanonicalEvmHeader {
	ch := l2rollup.CanonicalEvmHeader{
		ParentHash:       h.ParentHash.Bytes(),
		OmmersHash:       h.Sha3Uncles.Bytes(),
		Beneficiary:      h.Miner.Bytes(),
		StateRoot:        h.StateRoot.Bytes(),
		TransactionsRoot: h.TransactionsRoot.Bytes(),
		ReceiptsRoot:     h.ReceiptsRoot.Bytes(),
		LogsBloom:        []byte(h.LogsBloom),
		Difficulty:       corethU256((*big.Int)(h.Difficulty)),
		Number:           uint64(h.Number),
		GasLimit:         uint64(h.GasLimit),
		GasUsed:          uint64(h.GasUsed),
		Timestamp:        uint64(h.Timestamp),
		ExtraData:        []byte(h.ExtraData),
		MixHash:          h.MixHash.Bytes(),
		Nonce:            []byte(h.Nonce),
	}
	edh := l2rollup.HexBytes(h.ExtDataHash.Bytes())
	ch.ExtDataHash = &edh
	if h.BaseFeePerGas != nil {
		v := corethU256((*big.Int)(h.BaseFeePerGas))
		ch.BaseFeePerGas = &v
	}
	if h.ExtDataGasUsed != nil {
		v := corethU256((*big.Int)(h.ExtDataGasUsed))
		ch.ExtDataGasUsed = &v
	}
	if h.BlockGasCost != nil {
		v := corethU256((*big.Int)(h.BlockGasCost))
		ch.BlockGasCost = &v
	}
	ch.BlobGasUsed = uint64Ptr(h.BlobGasUsed)
	ch.ExcessBlobGas = uint64Ptr(h.ExcessBlobGas)
	if h.ParentBeaconBlockRoot != nil {
		pb := l2rollup.HexBytes(h.ParentBeaconBlockRoot.Bytes())
		ch.ParentBeaconBlockRoot = &pb
	}
	ch.TimeMilliseconds = uint64Ptr(h.TimeMilliseconds)
	ch.MinDelayExcess = uint64Ptr(h.MinDelayExcess)
	ch.TargetExponent = uint64Ptr(h.TargetExponent)
	ch.MinPriceExponent = uint64Ptr(h.MinPriceExponent)
	ch.SettledHeight = uint64Ptr(h.SettledHeight)
	ch.SettledGasUnix = uint64Ptr(h.SettledGasUnix)
	ch.SettledGasNumerator = uint64Ptr(h.SettledGasNumerator)
	ch.SettledExcess = uint64Ptr(h.SettledExcess)
	return ch
}

// corethU256 builds the shared wire u256 from a big.Int (nil -> 0).
func corethU256(v *big.Int) l2rollup.U256 {
	var u l2rollup.U256
	if v != nil {
		u.Int.Set(v)
	}
	return u
}

func uint64Ptr(v *hexutil.Uint64) *uint64 {
	if v == nil {
		return nil
	}
	u := uint64(*v)
	return &u
}

// corethHeaderHash recomputes keccak256(rlp(header)) with coreth's field order
// and optional-tail cascade.
func corethHeaderHash(h corethRPCHeader) (ethcommon.Hash, error) {
	fixed := []any{
		h.ParentHash, h.Sha3Uncles, h.Miner, h.StateRoot, h.TransactionsRoot,
		h.ReceiptsRoot, []byte(h.LogsBloom), (*big.Int)(h.Difficulty),
		new(big.Int).SetUint64(uint64(h.Number)), uint64(h.GasLimit),
		uint64(h.GasUsed), uint64(h.Timestamp), []byte(h.ExtraData), h.MixHash,
		[]byte(h.Nonce), h.ExtDataHash,
	}
	fields := make([][]byte, 0, len(fixed)+14)
	for _, f := range fixed {
		enc, err := rlp.EncodeToBytes(f)
		if err != nil {
			return ethcommon.Hash{}, fmt.Errorf("l2rollup: rlp-encode coreth header field: %w", err)
		}
		fields = append(fields, enc)
	}

	// Coreth optional tail in declaration order; nil is only encodable when a
	// later field is present, and then encodes as the RLP empty string.
	tail := []any{
		(*big.Int)(h.BaseFeePerGas), (*big.Int)(h.ExtDataGasUsed), (*big.Int)(h.BlockGasCost),
		h.BlobGasUsed, h.ExcessBlobGas, h.ParentBeaconBlockRoot,
		h.TimeMilliseconds, h.MinDelayExcess, h.TargetExponent, h.MinPriceExponent,
		h.SettledHeight, h.SettledGasUnix, h.SettledGasNumerator, h.SettledExcess,
	}
	last := -1
	for i, f := range tail {
		if !isNilTailField(f) {
			last = i
		}
	}
	for i := 0; i <= last; i++ {
		if isNilTailField(tail[i]) {
			fields = append(fields, []byte{0x80})
			continue
		}
		enc, err := rlp.EncodeToBytes(derefTailField(tail[i]))
		if err != nil {
			return ethcommon.Hash{}, fmt.Errorf("l2rollup: rlp-encode coreth tail field %d: %w", i, err)
		}
		fields = append(fields, enc)
	}

	payload := 0
	for _, f := range fields {
		payload += len(f)
	}
	out := make([]byte, 0, payload+9)
	out = appendRlpListHeader(out, payload)
	for _, f := range fields {
		out = append(out, f...)
	}
	return crypto.Keccak256Hash(out), nil
}

// isNilTailField reports whether a typed-nil tail pointer is absent.
func isNilTailField(f any) bool {
	switch v := f.(type) {
	case *big.Int:
		return v == nil
	case *hexutil.Uint64:
		return v == nil
	case *ethcommon.Hash:
		return v == nil
	default:
		return f == nil
	}
}

// derefTailField converts a present tail pointer to the value coreth encodes.
func derefTailField(f any) any {
	switch v := f.(type) {
	case *big.Int:
		return v
	case *hexutil.Uint64:
		return uint64(*v)
	case *ethcommon.Hash:
		return *v
	default:
		return f
	}
}

// appendRlpListHeader appends the RLP list prefix for a payload of the given size.
func appendRlpListHeader(dst []byte, payload int) []byte {
	if payload < 56 {
		return append(dst, byte(0xc0+payload))
	}
	size := payload
	var lenBytes []byte
	for size > 0 {
		lenBytes = append([]byte{byte(size)}, lenBytes...)
		size >>= 8
	}
	dst = append(dst, byte(0xf7+len(lenBytes)))
	return append(dst, lenBytes...)
}
