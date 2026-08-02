package bold

import (
	"context"
	"math/big"
	"testing"

	arbitrum "attestor/arbitrum"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestParseAssertionCreatedRecomputesAssertionHash(t *testing.T) {
	data := make([]byte, assertionCreatedDataWords*abiWordSize)
	l2BlockHash := common.HexToHash("0x1234")
	copy(data[13*abiWordSize:14*abiWordSize], l2BlockHash[:])
	data[17*abiWordSize+abiWordSize-1] = 1
	afterInboxBatchAcc := common.HexToHash("0x5678")
	copy(data[19*abiWordSize:20*abiWordSize], afterInboxBatchAcc[:])
	parentHash := common.HexToHash("0xabcd")
	afterStateHash := crypto.Keccak256Hash(data[13*abiWordSize : 19*abiWordSize])
	assertionHash := crypto.Keccak256Hash(
		parentHash[:],
		afterStateHash[:],
		afterInboxBatchAcc[:],
	)

	proposal, err := parseAssertionCreated(types.Log{
		Topics:      []common.Hash{assertionCreatedTopic, assertionHash, parentHash},
		Data:        data,
		BlockNumber: 99,
	})
	if err != nil {
		t.Fatalf("parse AssertionCreated: %v", err)
	}
	if proposal.AssertionHash != assertionHash ||
		proposal.ParentHash != parentHash ||
		proposal.L2BlockHash != l2BlockHash ||
		proposal.InboxAccumulator != afterInboxBatchAcc ||
		proposal.L1BlockNumber != 99 {
		t.Fatalf("parsed proposal: %+v", proposal)
	}

	badTopics := []common.Hash{assertionCreatedTopic, common.HexToHash("0xdead"), parentHash}
	if _, err := parseAssertionCreated(types.Log{
		Topics: badTopics,
		Data:   data,
	}); err == nil {
		t.Fatal("expected assertion hash mismatch")
	}
}

func TestRollupCoreSourceReadsStatusAndFinalizedLogs(t *testing.T) {
	rollupAddress := common.HexToAddress("0x0000000000000000000000000000000000001234")
	data := make([]byte, assertionConfirmedDataWords*abiWordSize)
	l2BlockHash := common.HexToHash("0xbeef")
	copy(data[:abiWordSize], l2BlockHash[:])
	assertionHash := common.HexToHash("0xface")
	client := &sourceTestL1Client{
		logs: []types.Log{{
			Address:     rollupAddress,
			Topics:      []common.Hash{assertionConfirmedTopic, assertionHash},
			Data:        data,
			BlockNumber: 100,
		}},
		storageOutput: make([]byte, abiWordSize),
	}
	const statusOffset = uint8(25)
	client.storageOutput[abiWordSize-1-int(statusOffset)] = assertionStatusConfirmed
	mappingSlot := common.HexToHash("0x76")

	source, err := NewRollupCoreSource(client, rollupAddress, mappingSlot, statusOffset)
	if err != nil {
		t.Fatalf("create source: %v", err)
	}
	_, confirmations, err := source.Assertions(context.Background(), 90, 100)
	if err != nil {
		t.Fatalf("read assertion logs: %v", err)
	}
	if len(confirmations) != 1 ||
		confirmations[0].AssertionHash != assertionHash ||
		confirmations[0].L2BlockHash != l2BlockHash {
		t.Fatalf("confirmations: %+v", confirmations)
	}
	status, err := source.AssertionStatus(context.Background(), arbitrum.ProposedAssertion{
		AssertionHash: assertionHash,
	}, 100)
	if err != nil {
		t.Fatalf("read assertion status: %v", err)
	}
	if status != assertionStatusConfirmed {
		t.Fatalf("status: got %d want %d", status, assertionStatusConfirmed)
	}
	statusPreimage := append(assertionHash.Bytes(), mappingSlot.Bytes()...)
	if client.storageKey != crypto.Keccak256Hash(statusPreimage) ||
		client.storageBlock.Uint64() != 100 {
		t.Fatalf(
			"status storage lookup: key=%s block=%v",
			client.storageKey,
			client.storageBlock,
		)
	}
	if client.filterQuery.FromBlock.Uint64() != 90 ||
		client.filterQuery.ToBlock.Uint64() != 100 ||
		len(client.filterQuery.Addresses) != 1 ||
		client.filterQuery.Addresses[0] != rollupAddress {
		t.Fatalf("filter query: %+v", client.filterQuery)
	}
}

func TestRollupCoreSourceSkipsZeroHashAssertionCreated(t *testing.T) {
	rollupAddress := common.HexToAddress("0x0000000000000000000000000000000000001234")
	data := make([]byte, assertionCreatedDataWords*abiWordSize)
	data[17*abiWordSize+abiWordSize-1] = 1
	client := &sourceTestL1Client{
		logs: []types.Log{{
			Address: rollupAddress,
			Topics: []common.Hash{
				assertionCreatedTopic,
				common.HexToHash("0x75"),
				common.HexToHash("0x50"),
			},
			Data:        data,
			BlockNumber: 100,
		}},
	}
	source, err := NewRollupCoreSource(
		client,
		rollupAddress,
		common.HexToHash("0x76"),
		25,
	)
	if err != nil {
		t.Fatalf("create source: %v", err)
	}
	proposals, confirmations, err := source.Assertions(
		context.Background(),
		90,
		100,
	)
	if err != nil {
		t.Fatalf("read assertion logs: %v", err)
	}
	if len(proposals) != 0 || len(confirmations) != 0 {
		t.Fatalf(
			"unexpected assertions: proposals=%+v confirmations=%+v",
			proposals,
			confirmations,
		)
	}
}

type sourceTestL1Client struct {
	logs          []types.Log
	callOutput    []byte
	storageOutput []byte
	storageKey    common.Hash
	storageBlock  *big.Int
	filterQuery   ethereum.FilterQuery
}

func (c *sourceTestL1Client) HeaderByNumber(
	context.Context,
	*big.Int,
) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(100)}, nil
}

func (c *sourceTestL1Client) ChainID(context.Context) (*big.Int, error) {
	return big.NewInt(1), nil
}

func (c *sourceTestL1Client) FilterLogs(
	_ context.Context,
	query ethereum.FilterQuery,
) ([]types.Log, error) {
	c.filterQuery = query
	return append([]types.Log(nil), c.logs...), nil
}

func (c *sourceTestL1Client) CallContract(
	context.Context,
	ethereum.CallMsg,
	*big.Int,
) ([]byte, error) {
	return append([]byte(nil), c.callOutput...), nil
}

func (c *sourceTestL1Client) StorageAt(
	_ context.Context,
	_ common.Address,
	key common.Hash,
	block *big.Int,
) ([]byte, error) {
	c.storageKey = key
	c.storageBlock = new(big.Int).Set(block)
	return append([]byte(nil), c.storageOutput...), nil
}
