package legacy

import (
	"bytes"
	"context"
	"math/big"
	"testing"

	boldattestor "attestor/arbitrum/bold"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func TestLegacyEventTopicsMatchCanonicalNitroABI(t *testing.T) {
	if got, want := nodeCreatedTopic, common.HexToHash(
		"0x4f4caa9e67fb994e349dd35d1ad0ce23053d4323f83ce11dc817b5435031d096",
	); got != want {
		t.Fatalf("NodeCreated topic: got %s want %s", got, want)
	}
	if got, want := nodeConfirmedTopic, common.HexToHash(
		"0x22ef0479a7ff660660d1c2fe35f1b632cf31675c2d9378db8cec95b00d8ffa3c",
	); got != want {
		t.Fatalf("NodeConfirmed topic: got %s want %s", got, want)
	}
	if got, want := nodeRejectedTopic, common.HexToHash(
		"0xeaffa3d968707ec919a2fc9f31d5ab2b86c905881ff561725d5a82fc95ad4640",
	); got != want {
		t.Fatalf("NodeRejected topic: got %s want %s", got, want)
	}
}

func TestParseLegacyNodeLifecycle(t *testing.T) {
	const nodeNumber = uint64(10_764)
	nodeHash := common.HexToHash("0x1111")
	parentHash := common.HexToHash("0x2222")
	l2BlockHash := common.HexToHash("0x3333")
	inboxAccumulator := common.HexToHash("0x4444")
	createdData := make([]byte, nodeCreatedDataWords*abiWordSize)
	copyWord(createdData, afterStateBlockHashWord, l2BlockHash)
	createdData[(afterStateMachineWord+1)*abiWordSize-1] = 1
	copyWord(createdData, afterInboxAccumulatorWord, inboxAccumulator)

	proposal, err := parseNodeCreated(types.Log{
		Topics: []common.Hash{
			nodeCreatedTopic,
			uint64Topic(nodeNumber),
			parentHash,
			nodeHash,
		},
		Data:        createdData,
		BlockNumber: 100,
	})
	if err != nil {
		t.Fatalf("parse NodeCreated: %v", err)
	}
	if proposal.AssertionHash != nodeHash ||
		proposal.ParentHash != parentHash ||
		proposal.L2BlockHash != l2BlockHash ||
		proposal.InboxAccumulator != inboxAccumulator ||
		proposal.LegacyNodeNumber != nodeNumber ||
		proposal.L1BlockNumber != 100 {
		t.Fatalf("parsed proposal: %+v", proposal)
	}

	confirmedData := make([]byte, nodeConfirmedDataWords*abiWordSize)
	copyWord(confirmedData, 0, l2BlockHash)
	confirmation, err := parseNodeConfirmed(types.Log{
		Topics:      []common.Hash{nodeConfirmedTopic, uint64Topic(nodeNumber)},
		Data:        confirmedData,
		BlockNumber: 101,
	})
	if err != nil {
		t.Fatalf("parse NodeConfirmed: %v", err)
	}
	if confirmation.LegacyNodeNumber != nodeNumber ||
		confirmation.L2BlockHash != l2BlockHash ||
		confirmation.L1BlockNumber != 101 {
		t.Fatalf("parsed confirmation: %+v", confirmation)
	}

	rejection, err := parseNodeRejected(types.Log{
		Topics:      []common.Hash{nodeRejectedTopic, uint64Topic(nodeNumber + 1)},
		BlockNumber: 102,
	})
	if err != nil {
		t.Fatalf("parse NodeRejected: %v", err)
	}
	if rejection.LegacyNodeNumber != nodeNumber+1 {
		t.Fatalf("parsed rejection: %+v", rejection)
	}
}

func TestParseCanonicalArbitrumSepoliaNode10764(t *testing.T) {
	data := common.FromHex(
		"0xb549ad556e95f848e63a2647890959bd7d8c0f205da4156b5f4db885889e3153" +
			"5ee95f66e142b3ab1db095d1205b88e230f2e8cd96b28e9aab07c537895189e9" +
			"63da69bd66a7e62f6d9c15ce6cc768c29dc5f792af78b7e08c0cd62242299f46" +
			"000000000000000000000000000000000000000000000000000000000005700f" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000001" +
			"fec7821ffcbba26c10f3e60964e2b7288262203c7dfa390285a61dc92d654ee5" +
			"e51b1f0c5c8fcd2155d335c309a49c82807ecc19e159d3f60eeba1d3e63ecb77" +
			"0000000000000000000000000000000000000000000000000000000000057022" +
			"0000000000000000000000000000000000000000000000000000000000000000" +
			"0000000000000000000000000000000000000000000000000000000000000001" +
			"0000000000000000000000000000000000000000000000000000000000003995" +
			"880be29479af37e93713f4b567a9fea3b4bc4ff74ca176fe93150bf8e77b7624" +
			"184884e1eb9fefdc158f6c8ac912bb183bf3cf83f0090317e0bc4ac5860baa39" +
			"0000000000000000000000000000000000000000000000000000000000057022",
	)
	proposal, err := parseNodeCreated(types.Log{
		Topics: []common.Hash{
			nodeCreatedTopic,
			uint64Topic(10_764),
			common.HexToHash(
				"0xa832396c917f4c88b023e6e945c25dd347b6985773983760064181dbc6ffe812",
			),
			common.HexToHash(
				"0x62d57802e8e11f418ddc8a54071ad73f1e1ec63c8ff40085fcd9fb34ce4d396f",
			),
		},
		Data:        data,
		BlockNumber: 7_258_441,
	})
	if err != nil {
		t.Fatalf("parse canonical Sepolia NodeCreated: %v", err)
	}
	if proposal.LegacyNodeNumber != 10_764 ||
		proposal.L2BlockHash != common.HexToHash(
			"0xfec7821ffcbba26c10f3e60964e2b7288262203c7dfa390285a61dc92d654ee5",
		) ||
		proposal.InboxAccumulator != common.HexToHash(
			"0x880be29479af37e93713f4b567a9fea3b4bc4ff74ca176fe93150bf8e77b7624",
		) {
		t.Fatalf("canonical Sepolia proposal: %+v", proposal)
	}
}

func TestLegacySourceReadsLifecycleAndStatus(t *testing.T) {
	const nodeNumber = uint64(77)
	rollupAddress := common.HexToAddress("0x0000000000000000000000000000000000001234")
	nodeHash := common.HexToHash("0x1234")
	l2BlockHash := common.HexToHash("0x5678")
	createdData := make([]byte, nodeCreatedDataWords*abiWordSize)
	copyWord(createdData, afterStateBlockHashWord, l2BlockHash)
	createdData[(afterStateMachineWord+1)*abiWordSize-1] = 2
	confirmedData := make([]byte, nodeConfirmedDataWords*abiWordSize)
	copyWord(confirmedData, 0, l2BlockHash)

	client := &sourceTestL1Client{
		logs: []types.Log{
			{
				Topics: []common.Hash{
					nodeCreatedTopic,
					uint64Topic(nodeNumber),
					common.HexToHash("0xabcd"),
					nodeHash,
				},
				Data: createdData,
			},
			{
				Topics: []common.Hash{nodeConfirmedTopic, uint64Topic(nodeNumber)},
				Data:   confirmedData,
			},
			{
				Topics: []common.Hash{nodeRejectedTopic, uint64Topic(nodeNumber + 1)},
			},
		},
		nodeHash:        nodeHash,
		latestConfirmed: nodeNumber,
	}
	source, err := NewRollupCoreSource(client, rollupAddress)
	if err != nil {
		t.Fatalf("create legacy source: %v", err)
	}
	proposals, confirmations, rejections, err := source.Assertions(
		context.Background(),
		90,
		100,
	)
	if err != nil {
		t.Fatalf("read legacy lifecycle: %v", err)
	}
	if len(proposals) != 1 ||
		len(confirmations) != 1 ||
		len(rejections) != 1 {
		t.Fatalf(
			"lifecycle lengths: proposals=%d confirmations=%d rejections=%d",
			len(proposals),
			len(confirmations),
			len(rejections),
		)
	}
	if client.filterQuery.FromBlock.Uint64() != 90 ||
		client.filterQuery.ToBlock.Uint64() != 100 ||
		len(client.filterQuery.Addresses) != 1 ||
		client.filterQuery.Addresses[0] != rollupAddress {
		t.Fatalf("filter query: %+v", client.filterQuery)
	}

	status, err := source.AssertionStatus(context.Background(), proposals[0], 100)
	if err != nil {
		t.Fatalf("read confirmed status: %v", err)
	}
	if status != boldattestor.AssertionStatusConfirmed {
		t.Fatalf("confirmed status: got %d", status)
	}
	client.latestConfirmed = nodeNumber - 1
	status, err = source.AssertionStatus(context.Background(), proposals[0], 100)
	if err != nil {
		t.Fatalf("read pending status: %v", err)
	}
	if status != boldattestor.AssertionStatusPending {
		t.Fatalf("pending status: got %d", status)
	}
	client.nodeHash = common.HexToHash("0xdead")
	status, err = source.AssertionStatus(context.Background(), proposals[0], 100)
	if err != nil {
		t.Fatalf("read missing status: %v", err)
	}
	if status != boldattestor.AssertionStatusNone {
		t.Fatalf("missing status: got %d", status)
	}
}

type sourceTestL1Client struct {
	logs            []types.Log
	filterQuery     ethereum.FilterQuery
	nodeHash        common.Hash
	latestConfirmed uint64
}

func (c *sourceTestL1Client) HeaderByNumber(
	context.Context,
	*big.Int,
) (*types.Header, error) {
	return &types.Header{Number: big.NewInt(100)}, nil
}

func (c *sourceTestL1Client) ChainID(context.Context) (*big.Int, error) {
	return big.NewInt(11_155_111), nil
}

func (c *sourceTestL1Client) FilterLogs(
	_ context.Context,
	query ethereum.FilterQuery,
) ([]types.Log, error) {
	c.filterQuery = query
	return append([]types.Log(nil), c.logs...), nil
}

func (c *sourceTestL1Client) CallContract(
	_ context.Context,
	message ethereum.CallMsg,
	_ *big.Int,
) ([]byte, error) {
	switch {
	case bytes.Equal(message.Data, chainIDSelector):
		return abiUint64(421_614), nil
	case bytes.Equal(message.Data, latestConfirmedSelector):
		return abiUint64(c.latestConfirmed), nil
	case len(message.Data) == len(getNodeSelector)+abiWordSize &&
		bytes.Equal(message.Data[:len(getNodeSelector)], getNodeSelector):
		output := make([]byte, 12*abiWordSize)
		copyWord(output, legacyNodeHashWord, c.nodeHash)
		return output, nil
	default:
		return nil, ethereum.NotFound
	}
}

func copyWord(data []byte, word int, value common.Hash) {
	copy(data[word*abiWordSize:(word+1)*abiWordSize], value[:])
}

func uint64Topic(value uint64) common.Hash {
	return common.BytesToHash(abiUint64(value))
}

func abiUint64(value uint64) []byte {
	return encodeUint64Call(nil, value)
}

var _ L1Client = (*sourceTestL1Client)(nil)
