package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	relayerclient "relayer/client"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ArbLegacyProfile is the subset of the Arbitrum `legacy_nitro` rollup profile the
// header builder needs (parsed from rollup_profile.protocol.value). It mirrors
// arbitrum-verifier LegacyProfile: the packed node-lifecycle slot, the `_nodes` mapping
// slot, the lifecycle field offsets, and the confirmData word offset.
type ArbLegacyProfile struct {
	RollupCore            ethcommon.Address // L1 RollupCore proxy (profile.rollup)
	L2Router              ethcommon.Address // L2 ICS26Router (router_proof target)
	L1ClientID            string            // shared ETH (08-wasm) client id on Cosmos
	NodeLifecycleSlot     ethcommon.Hash    // packed {_latestConfirmed,_firstUnresolvedNode,_latestNodeCreated}
	NodesMappingSlot      ethcommon.Hash    // Solidity slot of mapping(uint64 => Node) _nodes
	LatestConfirmedOffset uint8             // byte offset of _latestConfirmed in the lifecycle slot
	FirstUnresolvedOffset uint8             // byte offset of _firstUnresolvedNode
	LatestCreatedOffset   uint8             // byte offset of _latestNodeCreated
	ConfirmDataOffset     uint8             // storage-word offset of Node.confirmData from the node base
}

// arbitrumLegacyHeaderBuilder assembles an ArbitrumLegacyHeader: it resolves the target
// legacy Nitro node from the attestor, proves the RollupCore account + the packed
// node-lifecycle slot + `_nodes[node].confirmData` in the trusted L1 state, and binds
// the committed L2 block via confirmData = keccak(blockHash || sendRoot).
//
// Unlike the OP and BoLD builders (self-contained), the legacy builder consumes the
// attestor's LegacyNode provenance: the committed L2 block number is not recoverable
// from RollupCore storage alone (confirmData only commits keccak(blockHash||sendRoot)),
// and the attestor already maps node → L2 block by running Nitro.
type arbitrumLegacyHeaderBuilder struct {
	l1       *ethclient.Client // L1 exec (RollupCore eth_getProof)
	l2       *ethclient.Client // L2 exec (l2_header + sendRoot + router eth_getProof)
	cosmos   cosmosClientStateReader
	attestor AttestorClient
	srcChain string
	profile  ArbLegacyProfile
	// includeProvisional must match the Source's: the source decides WHICH heights
	// are relayable, this decides which node is proven for them. Two different
	// answers would gate on one rule and prove with another.
	includeProvisional bool
}

// NewArbitrumLegacyHeaderBuilder wires the legacy Arbitrum builder. includeProvisional
// must be the Source's value so the builder resolves the same node the source gated on.
func NewArbitrumLegacyHeaderBuilder(l1, l2 *ethclient.Client, cosmos cosmosClientStateReader, attestor AttestorClient, srcChain string, profile ArbLegacyProfile, includeProvisional bool) *arbitrumLegacyHeaderBuilder {
	return &arbitrumLegacyHeaderBuilder{
		l1: l1, l2: l2, cosmos: cosmos, attestor: attestor,
		srcChain: srcChain, profile: profile,
		includeProvisional: includeProvisional,
	}
}

func (b *arbitrumLegacyHeaderBuilder) Name() string { return "l2-arbitrum-legacy" }

// BuildHeader resolves the attested legacy node at or below request.Height, assembles its
// header, and returns the L2 block that node commits (at or below the request, so the
// client advances by it). RPC/availability failures are returned plain (the generic
// Builder wraps them chain.Retryable).
func (b *arbitrumLegacyHeaderBuilder) BuildHeader(ctx context.Context, request HeaderRequest) (ClientMessage, uint64, error) {
	root, found, err := b.attestor.AttestedRootAtOrBelow(ctx, b.srcChain, request.Height, b.includeProvisional)
	if err != nil {
		return nil, 0, err
	}
	if !found {
		return nil, 0, fmt.Errorf("l2-arbitrum-legacy: attestor has no %s node at or below L2 height %d", request.Finality, request.Height)
	}
	// The legacy provenance is the LegacyNode oneof; GetLegacyNode() is nil for a
	// non-legacy (game/assertion) root, so guard it before reading the node number.
	node := root.GetLegacyNode()
	if node == nil || node.GetNodeNumber() == 0 {
		return nil, 0, fmt.Errorf("l2-arbitrum-legacy: attested root at height %d carries no legacy node provenance", request.Height)
	}
	committedHeight := root.GetL2BlockNumber()
	header, err := b.buildLegacyHeaderFor(ctx, node.GetNodeNumber(), committedHeight)
	if err != nil {
		return nil, 0, err
	}
	return header, committedHeight, nil
}

// buildLegacyHeaderFor assembles the header for an explicit (nodeNumber, l2BlockNumber).
// It is the pure core BuildHeader delegates to, and is exercised directly by the devnet
// test without a live attestor.
func (b *arbitrumLegacyHeaderBuilder) buildLegacyHeaderFor(ctx context.Context, nodeNumber, l2BlockNumber uint64) (*ArbitrumLegacyHeader, error) {
	if nodeNumber == 0 {
		return nil, fmt.Errorf("l2-arbitrum-legacy: node number must not be zero")
	}

	// 1. The L1 block the shared ETH client trusts — prove the RollupCore there.
	// Ask for the pinned ETH client to be current first: this builder can only
	// prove at the block that client trusts, so its freshness is a precondition,
	// not a background nicety (#276). A failure here is logged by the
	// implementation and does not abort — the client may still be fresh enough,
	// or another direction may be advancing it.
	_ = b.cosmos.UpdateEthClientIfStale(ctx, b.profile.L1ClientID)
	beaconSlot, l1Block, err := b.cosmos.EthClientLatestSlotAndBlock(b.profile.L1ClientID)
	if err != nil {
		return nil, fmt.Errorf("l2-arbitrum-legacy: read shared ETH client height: %w", err)
	}
	l1BlockBig := new(big.Int).SetUint64(l1Block)
	l1Header, err := b.l1.HeaderByNumber(ctx, l1BlockBig)
	if err != nil {
		return nil, fmt.Errorf("l2-arbitrum-legacy: L1 header at %d: %w", l1Block, err)
	}

	// 2. Prove the RollupCore account + two storage slots: the packed node-lifecycle
	//    slot and _nodes[node].confirmData.
	nodeBase := mappingSlotU64(nodeNumber, b.profile.NodesMappingSlot)
	confirmDataSlot, err := addStorageOffset(nodeBase, b.profile.ConfirmDataOffset)
	if err != nil {
		return nil, err
	}
	rollupProof, err := relayerclient.EthGetProof(ctx, b.l1, b.profile.RollupCore,
		[]ethcommon.Hash{b.profile.NodeLifecycleSlot, confirmDataSlot}, l1BlockBig)
	if err != nil {
		return nil, annotatePinnedL1(ctx, b.l1, "l2-arbitrum-legacy: rollup node proof", b.profile.L1ClientID, l1Block, err)
	}
	lifecycle, err := findStorageProof(rollupProof, b.profile.NodeLifecycleSlot)
	if err != nil {
		return nil, err
	}
	confirmData, err := findStorageProof(rollupProof, confirmDataSlot)
	if err != nil {
		return nil, err
	}

	// Self-check: the node must be active (latest confirmed or a pending node) in the
	// proven lifecycle slot — never build a header for a resolved-loser node.
	if err := b.checkNodeActive(nodeNumber, lifecycle.Value); err != nil {
		return nil, err
	}

	// 3. Canonical L2 header + its send root (Arbitrum exposes sendRoot on the block).
	l2HeightBig := new(big.Int).SetUint64(l2BlockNumber)
	l2Header, err := b.l2.HeaderByNumber(ctx, l2HeightBig)
	if err != nil {
		return nil, fmt.Errorf("l2-arbitrum-legacy: L2 header at %d: %w", l2BlockNumber, err)
	}
	sendRoot, err := l2SendRoot(ctx, b.l2, l2BlockNumber)
	if err != nil {
		return nil, err
	}

	// Self-check: confirmData == keccak256(blockHash || sendRoot) — the same binding
	// verify_legacy enforces; catches a wrong node/block or a header-hash mismatch.
	cdWord, err := storageWord(confirmData.Value)
	if err != nil {
		return nil, fmt.Errorf("l2-arbitrum-legacy: confirmData value: %w", err)
	}
	want := crypto.Keccak256Hash(l2Header.Hash().Bytes(), sendRoot.Bytes())
	if ethcommon.BytesToHash(cdWord[:]) != want {
		return nil, fmt.Errorf("l2-arbitrum-legacy: confirmData mismatch for node %d block %d: onchain=%x computed=%s",
			nodeNumber, l2BlockNumber, cdWord, want)
	}

	// 4. Router proof against the committed L2 header's state root.
	routerProof, err := relayerclient.EthGetProof(ctx, b.l2, b.profile.L2Router, nil, l2HeightBig)
	if err != nil {
		return nil, err
	}

	return &ArbitrumLegacyHeader{
		BeaconSlot:         beaconSlot,
		L1StateRoot:        l1Header.Root.Bytes(),
		RollupProof:        EvmAccountProof{Proof: rollupProof.AccountProof},
		NodeNumber:         nodeNumber,
		NodeLifecycleProof: EvmStorageProof{Key: b.profile.NodeLifecycleSlot.Bytes(), Value: lifecycle.Value, Proof: lifecycle.Proof},
		ConfirmDataProof:   EvmStorageProof{Key: confirmDataSlot.Bytes(), Value: confirmData.Value, Proof: confirmData.Proof},
		SendRoot:           sendRoot.Bytes(),
		L2Header:           toCanonicalHeader(l2Header),
		RouterProof:        EvmAccountProof{Proof: routerProof.AccountProof},
	}, nil
}

// checkNodeActive decodes the packed lifecycle slot and rejects a node that is neither
// the latest confirmed nor within the pending [firstUnresolved, latestCreated] range.
func (b *arbitrumLegacyHeaderBuilder) checkNodeActive(nodeNumber uint64, lifecycleValue []byte) error {
	latestConfirmed, err := packedU64(lifecycleValue, b.profile.LatestConfirmedOffset)
	if err != nil {
		return fmt.Errorf("l2-arbitrum-legacy: decode latest confirmed: %w", err)
	}
	firstUnresolved, err := packedU64(lifecycleValue, b.profile.FirstUnresolvedOffset)
	if err != nil {
		return fmt.Errorf("l2-arbitrum-legacy: decode first unresolved: %w", err)
	}
	latestCreated, err := packedU64(lifecycleValue, b.profile.LatestCreatedOffset)
	if err != nil {
		return fmt.Errorf("l2-arbitrum-legacy: decode latest created: %w", err)
	}
	if !legacyNodeActive(nodeNumber, latestConfirmed, firstUnresolved, latestCreated) {
		return fmt.Errorf("l2-arbitrum-legacy: node %d is neither latest-confirmed (%d) nor pending [%d,%d]",
			nodeNumber, latestConfirmed, firstUnresolved, latestCreated)
	}
	return nil
}

// This file holds the DETERMINISTIC pieces of the legacy Arbitrum header builder — the
// ones that must match arbitrum-verifier byte-for-byte and so are unit-tested here.

// mappingSlotU64 derives the base slot of `_nodes[key]`, matching arbitrum-verifier's
// mapping_slot_u64: keccak256(pad32(key) || slot) with key right-aligned in the word.
func mappingSlotU64(key uint64, slot ethcommon.Hash) ethcommon.Hash {
	var keyWord [32]byte
	binary.BigEndian.PutUint64(keyWord[24:], key)
	return crypto.Keccak256Hash(keyWord[:], slot.Bytes())
}

// addStorageOffset adds a byte offset to a 32-byte storage slot, erroring on overflow
// past 32 bytes (matching arbitrum-verifier's add_storage_offset).
func addStorageOffset(base ethcommon.Hash, offset uint8) (ethcommon.Hash, error) {
	sum := new(big.Int).Add(new(big.Int).SetBytes(base.Bytes()), new(big.Int).SetUint64(uint64(offset)))
	if sum.BitLen() > 256 {
		return ethcommon.Hash{}, fmt.Errorf("l2-arbitrum-legacy: storage-slot offset overflows")
	}
	var out ethcommon.Hash
	sum.FillBytes(out[:])
	return out, nil
}

// storageWord left-pads a canonical (minimal big-endian) EVM storage value to one word.
func storageWord(value []byte) ([32]byte, error) {
	var w [32]byte
	if len(value) > 32 {
		return w, fmt.Errorf("storage value is %d bytes, wider than one word", len(value))
	}
	copy(w[32-len(value):], value)
	return w, nil
}

// packedU64 reads a big-endian uint64 packed at byte `offset` from the least-significant
// end of a storage word, matching arbitrum-verifier's packed_u64.
func packedU64(value []byte, offset uint8) (uint64, error) {
	w, err := storageWord(value)
	if err != nil {
		return 0, err
	}
	end := 32 - int(offset)
	start := end - 8
	if start < 0 || end > 32 {
		return 0, fmt.Errorf("l2-arbitrum-legacy: uint64 field at offset %d exceeds the storage word", offset)
	}
	return binary.BigEndian.Uint64(w[start:end]), nil
}

// legacyNodeActive mirrors arbitrum-verifier's legacy_node_is_active: the node is the
// latest confirmed one, or a pending node in [firstUnresolved, latestCreated].
func legacyNodeActive(node, latestConfirmed, firstUnresolved, latestCreated uint64) bool {
	return node == latestConfirmed || (node >= firstUnresolved && node <= latestCreated)
}

// findStorageProof returns the storage proof for slot from an eth_getProof result,
// matching by key so it does not depend on the RPC preserving request order.
func findStorageProof(proof *relayerclient.RawEvmProof, slot ethcommon.Hash) (relayerclient.RawStorageProof, error) {
	for _, sp := range proof.Storage {
		if sp.Key == slot {
			return sp, nil
		}
	}
	return relayerclient.RawStorageProof{}, fmt.Errorf("l2-arbitrum-legacy: eth_getProof missing storage entry for slot %s", slot)
}

// l2SendRoot reads the Arbitrum block's send root (a non-standard header field the RPC
// exposes on eth_getBlockByNumber). It is the send root confirmData commits to.
func l2SendRoot(ctx context.Context, l2 *ethclient.Client, blockNumber uint64) (ethcommon.Hash, error) {
	var blk struct {
		SendRoot *ethcommon.Hash `json:"sendRoot"`
	}
	if err := l2.Client().CallContext(ctx, &blk, "eth_getBlockByNumber", hexUint64(blockNumber), false); err != nil {
		return ethcommon.Hash{}, fmt.Errorf("l2-arbitrum-legacy: eth_getBlockByNumber(%d) send root: %w", blockNumber, err)
	}
	if blk.SendRoot == nil {
		return ethcommon.Hash{}, fmt.Errorf("l2-arbitrum-legacy: L2 block %d has no sendRoot (not an Arbitrum chain?)", blockNumber)
	}
	return *blk.SendRoot, nil
}

var _ HeaderBuilder = (*arbitrumLegacyHeaderBuilder)(nil)
