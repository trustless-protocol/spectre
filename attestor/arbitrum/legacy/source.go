// Package legacy ingests pre-BoLD Nitro rollup nodes from finalized L1.
package legacy

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	arbitrum "attestor/arbitrum"
	boldattestor "attestor/arbitrum/bold"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	abiWordSize               = 32
	nodeCreatedDataWords      = 15
	nodeConfirmedDataWords    = 2
	legacyNodeHashWord        = 11
	afterStateBlockHashWord   = 6
	afterStateMachineWord     = 10
	afterInboxAccumulatorWord = 12
)

var (
	nodeCreatedTopic = crypto.Keccak256Hash([]byte(
		"NodeCreated(uint64,bytes32,bytes32,bytes32,(((bytes32[2],uint64[2]),uint8),((bytes32[2],uint64[2]),uint8),uint64),bytes32,bytes32,uint256)",
	))
	nodeConfirmedTopic = crypto.Keccak256Hash([]byte(
		"NodeConfirmed(uint64,bytes32,bytes32)",
	))
	nodeRejectedTopic       = crypto.Keccak256Hash([]byte("NodeRejected(uint64)"))
	chainIDSelector         = crypto.Keccak256([]byte("chainId()"))[:4]
	getNodeSelector         = crypto.Keccak256([]byte("getNode(uint64)"))[:4]
	latestConfirmedSelector = crypto.Keccak256([]byte("latestConfirmed()"))[:4]
)

// L1Client is the finalized Ethereum RPC surface required by RollupCoreSource.
type L1Client interface {
	HeaderByNumber(context.Context, *big.Int) (*types.Header, error)
	ChainID(context.Context) (*big.Int, error)
	FilterLogs(context.Context, ethereum.FilterQuery) ([]types.Log, error)
	CallContract(context.Context, ethereum.CallMsg, *big.Int) ([]byte, error)
}

// RollupCoreSource reads the legacy NodeCreated, NodeConfirmed, and
// NodeRejected lifecycle used by canonical Arbitrum Sepolia.
type RollupCoreSource struct {
	client  L1Client
	address common.Address
}

// NewRollupCoreSource constructs a finalized-L1 legacy node source.
func NewRollupCoreSource(
	client L1Client,
	address common.Address,
) (*RollupCoreSource, error) {
	if client == nil {
		return nil, errors.New("L1 client must not be nil")
	}
	if address == (common.Address{}) {
		return nil, errors.New("RollupCore address must not be zero")
	}
	return &RollupCoreSource{client: client, address: address}, nil
}

// FinalizedBlockNumber returns the L1 execution block covered by Ethereum's
// finalized tag.
func (s *RollupCoreSource) FinalizedBlockNumber(ctx context.Context) (uint64, error) {
	header, err := s.client.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if errors.Is(err, ethereum.NotFound) || (err == nil && header == nil) {
		return 0, fmt.Errorf("finalized L1 block was not found: %w", ethereum.NotFound)
	}
	if err != nil {
		return 0, fmt.Errorf("query finalized L1 block: %w", err)
	}
	if header.Number == nil || !header.Number.IsUint64() {
		return 0, errors.New("finalized L1 block number is not a uint64")
	}
	return header.Number.Uint64(), nil
}

// ChainIDs returns the connected Ethereum L1 chain ID and legacy RollupCore's
// configured Arbitrum L2 chain ID at the supplied finalized L1 block.
func (s *RollupCoreSource) ChainIDs(
	ctx context.Context,
	finalizedL1Block uint64,
) (*big.Int, *big.Int, error) {
	l1ChainID, err := s.client.ChainID(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("query L1 chain ID: %w", err)
	}
	output, err := s.callAt(ctx, chainIDSelector, finalizedL1Block)
	if err != nil {
		return nil, nil, fmt.Errorf("query legacy RollupCore chainId: %w", err)
	}
	if len(output) != abiWordSize {
		return nil, nil, fmt.Errorf(
			"legacy RollupCore chainId returned %d bytes, want %d",
			len(output),
			abiWordSize,
		)
	}
	return new(big.Int).Set(l1ChainID), new(big.Int).SetBytes(output), nil
}

// Assertions reads the complete legacy node lifecycle from one inclusive
// finalized-L1 block range.
func (s *RollupCoreSource) Assertions(
	ctx context.Context,
	fromL1Block uint64,
	toL1Block uint64,
) (
	[]arbitrum.ProposedAssertion,
	[]boldattestor.ConfirmedAssertion,
	[]boldattestor.RejectedAssertion,
	error,
) {
	if fromL1Block > toL1Block {
		return nil, nil, nil, fmt.Errorf(
			"invalid legacy node log range %d..%d",
			fromL1Block,
			toL1Block,
		)
	}
	logs, err := s.client.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromL1Block),
		ToBlock:   new(big.Int).SetUint64(toL1Block),
		Addresses: []common.Address{s.address},
		Topics: [][]common.Hash{{
			nodeCreatedTopic,
			nodeConfirmedTopic,
			nodeRejectedTopic,
		}},
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"filter legacy RollupCore node logs %d..%d: %w",
			fromL1Block,
			toL1Block,
			err,
		)
	}

	var proposals []arbitrum.ProposedAssertion
	var confirmations []boldattestor.ConfirmedAssertion
	var rejections []boldattestor.RejectedAssertion
	for _, event := range logs {
		if event.Removed {
			return nil, nil, nil, fmt.Errorf(
				"finalized legacy RollupCore log at L1 block %d was marked removed",
				event.BlockNumber,
			)
		}
		if len(event.Topics) == 0 {
			return nil, nil, nil, fmt.Errorf(
				"legacy RollupCore log at L1 block %d has no event topic",
				event.BlockNumber,
			)
		}
		switch event.Topics[0] {
		case nodeCreatedTopic:
			proposal, parseErr := parseNodeCreated(event)
			if parseErr != nil {
				return nil, nil, nil, parseErr
			}
			proposals = append(proposals, proposal)
		case nodeConfirmedTopic:
			confirmation, parseErr := parseNodeConfirmed(event)
			if parseErr != nil {
				return nil, nil, nil, parseErr
			}
			confirmations = append(confirmations, confirmation)
		case nodeRejectedTopic:
			rejection, parseErr := parseNodeRejected(event)
			if parseErr != nil {
				return nil, nil, nil, parseErr
			}
			rejections = append(rejections, rejection)
		default:
			return nil, nil, nil, fmt.Errorf(
				"unexpected legacy RollupCore event topic %s",
				event.Topics[0],
			)
		}
	}
	return proposals, confirmations, rejections, nil
}

// AssertionStatus resolves a persisted legacy node against RollupCore at a
// finalized L1 block. Rejected nodes are removed by NodeRejected ingestion;
// the node-hash comparison additionally fails closed on deployment changes.
func (s *RollupCoreSource) AssertionStatus(
	ctx context.Context,
	proposal arbitrum.ProposedAssertion,
	finalizedL1Block uint64,
) (uint8, error) {
	if proposal.LegacyNodeNumber == 0 {
		return 0, errors.New("legacy proposal node number must not be zero")
	}
	nodeHash, err := s.nodeHashAt(ctx, proposal.LegacyNodeNumber, finalizedL1Block)
	if err != nil {
		return 0, err
	}
	if nodeHash == (common.Hash{}) || nodeHash != proposal.AssertionHash {
		return boldattestor.AssertionStatusNone, nil
	}
	latestConfirmed, err := s.latestConfirmedAt(ctx, finalizedL1Block)
	if err != nil {
		return 0, err
	}
	if proposal.LegacyNodeNumber <= latestConfirmed {
		return boldattestor.AssertionStatusConfirmed, nil
	}
	return boldattestor.AssertionStatusPending, nil
}

func (s *RollupCoreSource) nodeHashAt(
	ctx context.Context,
	nodeNumber uint64,
	l1Block uint64,
) (common.Hash, error) {
	output, err := s.callAt(ctx, encodeUint64Call(getNodeSelector, nodeNumber), l1Block)
	if err != nil {
		return common.Hash{}, fmt.Errorf("query legacy node %d: %w", nodeNumber, err)
	}
	const nodeWords = 12
	if len(output) != nodeWords*abiWordSize {
		return common.Hash{}, fmt.Errorf(
			"legacy getNode(%d) returned %d bytes, want %d",
			nodeNumber,
			len(output),
			nodeWords*abiWordSize,
		)
	}
	return wordHash(output, legacyNodeHashWord), nil
}

func (s *RollupCoreSource) latestConfirmedAt(
	ctx context.Context,
	l1Block uint64,
) (uint64, error) {
	output, err := s.callAt(ctx, latestConfirmedSelector, l1Block)
	if err != nil {
		return 0, fmt.Errorf("query legacy latestConfirmed: %w", err)
	}
	value, err := uint64Word(output)
	if err != nil {
		return 0, fmt.Errorf("decode legacy latestConfirmed: %w", err)
	}
	return value, nil
}

func (s *RollupCoreSource) callAt(
	ctx context.Context,
	data []byte,
	l1Block uint64,
) ([]byte, error) {
	return s.client.CallContract(ctx, ethereum.CallMsg{
		To:   &s.address,
		Data: append([]byte(nil), data...),
	}, new(big.Int).SetUint64(l1Block))
}

func parseNodeCreated(event types.Log) (arbitrum.ProposedAssertion, error) {
	if len(event.Topics) != 4 {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated at L1 block %d has %d topics, want 4",
			event.BlockNumber,
			len(event.Topics),
		)
	}
	if len(event.Data) != nodeCreatedDataWords*abiWordSize {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated at L1 block %d has %d data bytes, want %d",
			event.BlockNumber,
			len(event.Data),
			nodeCreatedDataWords*abiWordSize,
		)
	}
	nodeNumber, err := uint64Word(event.Topics[1][:])
	if err != nil {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"decode NodeCreated node number at L1 block %d: %w",
			event.BlockNumber,
			err,
		)
	}
	if nodeNumber == 0 {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated at L1 block %d has a zero node number",
			event.BlockNumber,
		)
	}
	nodeHash := event.Topics[3]
	if nodeHash == (common.Hash{}) {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated %d has a zero node hash",
			nodeNumber,
		)
	}
	machineStatus, err := uint8Word(
		event.Data[afterStateMachineWord*abiWordSize : (afterStateMachineWord+1)*abiWordSize],
	)
	if err != nil {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"decode NodeCreated %d terminal machine status: %w",
			nodeNumber,
			err,
		)
	}
	if machineStatus != 1 && machineStatus != 2 {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated %d has invalid terminal machine status %d",
			nodeNumber,
			machineStatus,
		)
	}
	l2BlockHash := wordHash(event.Data, afterStateBlockHashWord)
	if l2BlockHash == (common.Hash{}) {
		return arbitrum.ProposedAssertion{}, fmt.Errorf(
			"NodeCreated %d has a zero L2 block hash",
			nodeNumber,
		)
	}
	return arbitrum.ProposedAssertion{
		AssertionHash:    nodeHash,
		ParentHash:       event.Topics[2],
		L2BlockHash:      l2BlockHash,
		InboxAccumulator: wordHash(event.Data, afterInboxAccumulatorWord),
		L1BlockNumber:    event.BlockNumber,
		LegacyNodeNumber: nodeNumber,
	}, nil
}

func parseNodeConfirmed(event types.Log) (boldattestor.ConfirmedAssertion, error) {
	if len(event.Topics) != 2 {
		return boldattestor.ConfirmedAssertion{}, fmt.Errorf(
			"NodeConfirmed at L1 block %d has %d topics, want 2",
			event.BlockNumber,
			len(event.Topics),
		)
	}
	if len(event.Data) != nodeConfirmedDataWords*abiWordSize {
		return boldattestor.ConfirmedAssertion{}, fmt.Errorf(
			"NodeConfirmed at L1 block %d has %d data bytes, want %d",
			event.BlockNumber,
			len(event.Data),
			nodeConfirmedDataWords*abiWordSize,
		)
	}
	nodeNumber, err := uint64Word(event.Topics[1][:])
	if err != nil {
		return boldattestor.ConfirmedAssertion{}, fmt.Errorf(
			"decode NodeConfirmed node number at L1 block %d: %w",
			event.BlockNumber,
			err,
		)
	}
	if nodeNumber == 0 {
		return boldattestor.ConfirmedAssertion{}, fmt.Errorf(
			"NodeConfirmed at L1 block %d has a zero node number",
			event.BlockNumber,
		)
	}
	l2BlockHash := wordHash(event.Data, 0)
	if l2BlockHash == (common.Hash{}) {
		return boldattestor.ConfirmedAssertion{}, fmt.Errorf(
			"NodeConfirmed %d has a zero L2 block hash",
			nodeNumber,
		)
	}
	return boldattestor.ConfirmedAssertion{
		L2BlockHash:      l2BlockHash,
		L1BlockNumber:    event.BlockNumber,
		LegacyNodeNumber: nodeNumber,
	}, nil
}

func parseNodeRejected(event types.Log) (boldattestor.RejectedAssertion, error) {
	if len(event.Topics) != 2 {
		return boldattestor.RejectedAssertion{}, fmt.Errorf(
			"NodeRejected at L1 block %d has %d topics, want 2",
			event.BlockNumber,
			len(event.Topics),
		)
	}
	if len(event.Data) != 0 {
		return boldattestor.RejectedAssertion{}, fmt.Errorf(
			"NodeRejected at L1 block %d has %d data bytes, want 0",
			event.BlockNumber,
			len(event.Data),
		)
	}
	nodeNumber, err := uint64Word(event.Topics[1][:])
	if err != nil {
		return boldattestor.RejectedAssertion{}, fmt.Errorf(
			"decode NodeRejected node number at L1 block %d: %w",
			event.BlockNumber,
			err,
		)
	}
	if nodeNumber == 0 {
		return boldattestor.RejectedAssertion{}, fmt.Errorf(
			"NodeRejected at L1 block %d has a zero node number",
			event.BlockNumber,
		)
	}
	return boldattestor.RejectedAssertion{LegacyNodeNumber: nodeNumber}, nil
}

func encodeUint64Call(selector []byte, value uint64) []byte {
	data := make([]byte, len(selector)+abiWordSize)
	copy(data, selector)
	binary.BigEndian.PutUint64(data[len(selector)+abiWordSize-8:], value)
	return data
}

func uint64Word(word []byte) (uint64, error) {
	if len(word) != abiWordSize {
		return 0, fmt.Errorf("uint64 ABI word has %d bytes, want %d", len(word), abiWordSize)
	}
	if !bytes.Equal(word[:abiWordSize-8], make([]byte, abiWordSize-8)) {
		return 0, errors.New("uint64 ABI word has non-zero high bytes")
	}
	return binary.BigEndian.Uint64(word[abiWordSize-8:]), nil
}

func uint8Word(word []byte) (uint8, error) {
	if len(word) != abiWordSize {
		return 0, fmt.Errorf("uint8 ABI word has %d bytes, want %d", len(word), abiWordSize)
	}
	if !bytes.Equal(word[:abiWordSize-1], make([]byte, abiWordSize-1)) {
		return 0, errors.New("uint8 ABI word has non-zero high bytes")
	}
	return word[abiWordSize-1], nil
}

func wordHash(data []byte, word int) common.Hash {
	start := word * abiWordSize
	return common.BytesToHash(data[start : start+abiWordSize])
}

var _ boldattestor.AssertionSource = (*RollupCoreSource)(nil)
