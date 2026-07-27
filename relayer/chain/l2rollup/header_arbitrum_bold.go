package l2rollup

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/big"

	relayerclient "relayer/client"

	"github.com/ethereum/go-ethereum"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// defaultAssertionScanRange bounds how far back (in L1 blocks) BuildHeader scans for
// AssertionCreated events when the profile does not pin a range. A production build
// should instead consult the attestor's assertion_hash provenance for an O(1) lookup
// (mirrors gameIndexForHeight's note on the OP side).
const defaultAssertionScanRange = 50_000

// assertionCreatedDataWords is the fixed non-indexed word count of the BoLD v2
// AssertionCreated event. Its layout is frozen by the RollupCore ABI; the decode
// below indexes into it, so a length mismatch is ABI drift and fails loud.
const assertionCreatedDataWords = 25

const evmWord = 32

// assertionCreatedTopic is keccak256 of the BoLD v2 AssertionCreated event signature.
// It must stay byte-identical to the attestor's (attestor/arbitrum/bold/source.go) so
// both decode the same log.
var assertionCreatedTopic = crypto.Keccak256Hash([]byte(
	"AssertionCreated(bytes32,bytes32,((bytes32,bytes32,(bytes32,uint256,address,uint64,uint64)),((bytes32[2],uint64[2]),uint8,bytes32),((bytes32[2],uint64[2]),uint8,bytes32)),bytes32,uint256,bytes32,uint256,address,uint64)",
))

// ArbBoldProfile is the subset of the Arbitrum rollup profile the header builder needs to
// assemble proofs (parsed from rollup_profile by the caller). The verifier-only fields
// (assertion_status_offset, bold_version, commitment_slot) are not needed here — the
// relayer sends witnesses and the wasm client re-derives + checks them.
type ArbBoldProfile struct {
	RollupCore            ethcommon.Address // L1 RollupCore proxy (profile.rollup)
	AssertionsMappingSlot ethcommon.Hash    // Solidity slot of the BoLD _assertions mapping
	L2Router              ethcommon.Address // L2 ICS26Router (router_proof target)
	L1ClientID            string            // shared ETH (08-wasm) client id on Cosmos
	AssertionScanRange    uint64            // L1 blocks scanned back for AssertionCreated (0 → default)
}

// arbitrumBoldHeaderBuilder assembles an ArbitrumBoldHeader: it proves a nonzero-status BoLD
// assertion in the L1 state the shared ETH client trusts, reconstructs the assertion
// claim from its AssertionCreated event, binds it to the canonical L2 header the
// assertion commits, and proves the L2 router state against that header.
type arbitrumBoldHeaderBuilder struct {
	l1      *ethclient.Client // L1 exec (RollupCore eth_getProof + AssertionCreated logs)
	l2      *ethclient.Client // L2 exec (l2_header + router eth_getProof)
	cosmos  cosmosClientStateReader
	profile ArbBoldProfile
}

// NewArbitrumBoldHeaderBuilder wires the Arbitrum builder to the L1/L2 exec RPCs, the
// Cosmos client-state reader (for the trusted L1 block), and the parsed profile.
func NewArbitrumBoldHeaderBuilder(l1, l2 *ethclient.Client, cosmos cosmosClientStateReader, profile ArbBoldProfile) *arbitrumBoldHeaderBuilder {
	return &arbitrumBoldHeaderBuilder{l1: l1, l2: l2, cosmos: cosmos, profile: profile}
}

func (a *arbitrumBoldHeaderBuilder) Name() string { return "l2-arbitrum" }

// BuildHeader assembles the ArbitrumBoldHeader for the highest attested BoLD assertion
// at or below l2Height, and returns the L2 block that assertion commits (which can be
// lower than l2Height — e.g. the attested assertion postdates the ETH client's trusted
// L1 block — so the client advances by the committed block, not the request).
// RPC/availability failures are returned plain (the generic Builder wraps them
// chain.Retryable).
func (a *arbitrumBoldHeaderBuilder) BuildHeader(ctx context.Context, l2Height uint64) (ClientMessage, uint64, error) {
	// 1. The L1 block the shared ETH client trusts — prove the RollupCore there.
	beaconSlot, l1Block, err := a.cosmos.EthClientLatestSlotAndBlock(a.profile.L1ClientID)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-arbitrum: read shared ETH client height: %w", err)
	}
	l1BlockBig := new(big.Int).SetUint64(l1Block)
	l1Header, err := a.l1.HeaderByNumber(ctx, l1BlockBig)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-arbitrum: L1 header at %d: %w", l1Block, err)
	}

	// 2. Pick the assertion committing the largest L2 block <= l2Height, present in the
	//    trusted L1 state.
	claim, assertionHash, assertionL2, err := a.assertionForHeight(ctx, l1Block, l2Height)
	if err != nil {
		return nil, 0, err
	}

	// 3. Prove the RollupCore account + the packed assertion node at its mapping slot.
	assertionSlot := assertionStorageSlot(assertionHash, a.profile.AssertionsMappingSlot)
	rollupProof, err := relayerclient.EthGetProof(a.l1, a.profile.RollupCore, []ethcommon.Hash{assertionSlot}, l1BlockBig)
	if err != nil {
		return nil, 0, err
	}
	if len(rollupProof.Storage) == 0 {
		return nil, 0, fmt.Errorf("l2-arbitrum: rollup assertion-slot proof missing storage entry")
	}

	// 4. Canonical L2 header (the block the assertion commits) + 5. router proof.
	l2HeightBig := new(big.Int).SetUint64(assertionL2)
	l2Header, err := a.l2.HeaderByNumber(ctx, l2HeightBig)
	if err != nil {
		return nil, 0, fmt.Errorf("l2-arbitrum: L2 header at %d: %w", assertionL2, err)
	}
	routerProof, err := relayerclient.EthGetProof(a.l2, a.profile.L2Router, nil, l2HeightBig)
	if err != nil {
		return nil, 0, err
	}

	return &ArbitrumBoldHeader{
		BeaconSlot:     beaconSlot,
		L1StateRoot:    l1Header.Root.Bytes(),
		RollupProof:    EvmAccountProof{Proof: rollupProof.AccountProof},
		AssertionHash:  assertionHash.Bytes(),
		AssertionProof: EvmStorageProof{Key: assertionSlot.Bytes(), Value: rollupProof.Storage[0].Value, Proof: rollupProof.Storage[0].Proof},
		Assertion:      claim,
		L2Header:       toCanonicalHeader(l2Header),
		RouterProof:    EvmAccountProof{Proof: routerProof.AccountProof},
	}, assertionL2, nil
}

// assertionForHeight scans AssertionCreated events on the RollupCore and returns the
// newest one whose committed L2 block is <= l2Height (never exceeding the attested
// frontier the caller passed). BoLD assertions commit monotonically increasing L2
// blocks, so a newest-first scan returns the correct assertion after the fewest L2
// lookups. Malformed events (ABI drift) fail loud; assertions whose block is > l2Height,
// still running, or not resolvable on the L2 node are skipped.
func (a *arbitrumBoldHeaderBuilder) assertionForHeight(ctx context.Context, l1Block, l2Height uint64) (AssertionClaim, ethcommon.Hash, uint64, error) {
	scan := a.profile.AssertionScanRange
	if scan == 0 {
		scan = defaultAssertionScanRange
	}
	from := uint64(0)
	if l1Block > scan {
		from = l1Block - scan
	}

	logs, err := a.l1.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(from),
		ToBlock:   new(big.Int).SetUint64(l1Block),
		Addresses: []ethcommon.Address{a.profile.RollupCore},
		Topics:    [][]ethcommon.Hash{{assertionCreatedTopic}},
	})
	if err != nil {
		return AssertionClaim{}, ethcommon.Hash{}, 0, fmt.Errorf("l2-arbitrum: filter AssertionCreated: %w", err)
	}
	if len(logs) == 0 {
		return AssertionClaim{}, ethcommon.Hash{}, 0, fmt.Errorf("l2-arbitrum: no AssertionCreated in L1 blocks [%d,%d]", from, l1Block)
	}

	for i := len(logs) - 1; i >= 0; i-- {
		claim, blockHash, assertionHash, decErr := decodeAssertionCreated(logs[i])
		if decErr != nil {
			return AssertionClaim{}, ethcommon.Hash{}, 0, decErr
		}
		if claim.AfterState.MachineStatus == MachineStatusRunning {
			continue // unfinished assertion — not attestable
		}
		l2h, err := a.l2.HeaderByHash(ctx, blockHash)
		if err != nil {
			continue // committed L2 block not on this node (reorged/unavailable)
		}
		if n := l2h.Number.Uint64(); n <= l2Height {
			return claim, assertionHash, n, nil
		}
	}
	return AssertionClaim{}, ethcommon.Hash{}, 0, fmt.Errorf("l2-arbitrum: no attested assertion at or below L2 height %d", l2Height)
}

// This file holds the DETERMINISTIC pieces of the Arbitrum header builder — the ones
// that must match the arbitrum-verifier byte-for-byte and so are unit-tested here. The
// RPC assembly lives in BuildHeader and is validated end-to-end on a Nitro devnet.

// assertionStorageSlot derives the BoLD _assertions mapping slot for assertionHash,
// matching arbitrum-verifier's mapping_slot_bytes32: keccak256(key || slot).
func assertionStorageSlot(assertionHash, mappingSlot ethcommon.Hash) ethcommon.Hash {
	return crypto.Keccak256Hash(assertionHash.Bytes(), mappingSlot.Bytes())
}

// decodeAssertionCreated reconstructs the AssertionClaim from one AssertionCreated log
// and returns it with the committed L2 block hash and the assertion hash (topic 1). It
// self-checks the reconstructed hash against the event topic — the same identity
// arbitrum-verifier's AssertionClaim::hash enforces — so any offset drift fails loud.
func decodeAssertionCreated(log types.Log) (AssertionClaim, ethcommon.Hash, ethcommon.Hash, error) {
	if len(log.Topics) != 3 {
		return AssertionClaim{}, ethcommon.Hash{}, ethcommon.Hash{}, fmt.Errorf("l2-arbitrum: AssertionCreated has %d topics, want 3", len(log.Topics))
	}
	if len(log.Data) != assertionCreatedDataWords*evmWord {
		return AssertionClaim{}, ethcommon.Hash{}, ethcommon.Hash{}, fmt.Errorf("l2-arbitrum: AssertionCreated has %d data bytes, want %d", len(log.Data), assertionCreatedDataWords*evmWord)
	}

	assertionHash := log.Topics[1]
	parentHash := log.Topics[2]

	// afterState is the 6-word AssertionState = words 13..18 (abi.encode of the static
	// nested struct). Its keccak feeds the BoLD assertion hash.
	afterState := log.Data[13*evmWord : 19*evmWord]
	blockHash := ethcommon.BytesToHash(afterState[0:evmWord])                     // GlobalState.bytes32Vals[0]
	sendRoot := ethcommon.BytesToHash(afterState[evmWord : 2*evmWord])            // GlobalState.bytes32Vals[1]
	inboxPosition := binary.BigEndian.Uint64(afterState[3*evmWord-8 : 3*evmWord]) // u64Vals[0] (word 15, low 8)
	positionInMessage := binary.BigEndian.Uint64(afterState[4*evmWord-8 : 4*evmWord])
	statusByte := afterState[5*evmWord-1]                                      // machineStatus (word 17, last byte)
	endHistoryRoot := ethcommon.BytesToHash(afterState[5*evmWord : 6*evmWord]) // word 18
	inboxAcc := ethcommon.BytesToHash(log.Data[19*evmWord : 20*evmWord])       // afterInboxBatchAcc

	machineStatus, err := machineStatusFromByte(statusByte)
	if err != nil {
		return AssertionClaim{}, ethcommon.Hash{}, ethcommon.Hash{}, err
	}

	claim := AssertionClaim{
		ParentAssertionHash: parentHash.Bytes(),
		AfterState: AssertionState{
			GlobalState: GlobalState{
				Bytes32Vals: [2]hexBytes{blockHash.Bytes(), sendRoot.Bytes()},
				U64Vals:     [2]uint64{inboxPosition, positionInMessage},
			},
			MachineStatus:  machineStatus,
			EndHistoryRoot: endHistoryRoot.Bytes(),
		},
		InboxAcc: inboxAcc.Bytes(),
	}

	if got := boldAssertionHash(parentHash, afterState, inboxAcc); got != assertionHash {
		return AssertionClaim{}, ethcommon.Hash{}, ethcommon.Hash{}, fmt.Errorf("l2-arbitrum: AssertionCreated hash mismatch: event=%s computed=%s", assertionHash, got)
	}
	return claim, blockHash, assertionHash, nil
}

// boldAssertionHash reproduces BoLD v2 RollupLib.assertionHash:
// keccak256(parent || keccak256(afterState) || inboxAcc).
func boldAssertionHash(parent ethcommon.Hash, afterState []byte, inboxAcc ethcommon.Hash) ethcommon.Hash {
	afterStateHash := crypto.Keccak256(afterState)
	return crypto.Keccak256Hash(parent.Bytes(), afterStateHash, inboxAcc.Bytes())
}

// machineStatusFromByte maps the Solidity MachineStatus enum byte to the wire value
// (0 → running, 1 → finished, 2 → errored). An unknown value is ABI drift.
func machineStatusFromByte(b byte) (MachineStatus, error) {
	switch b {
	case 0:
		return MachineStatusRunning, nil
	case 1:
		return MachineStatusFinished, nil
	case 2:
		return MachineStatusErrored, nil
	default:
		return "", fmt.Errorf("l2-arbitrum: unknown machine status %d", b)
	}
}

var _ HeaderBuilder = (*arbitrumBoldHeaderBuilder)(nil)
