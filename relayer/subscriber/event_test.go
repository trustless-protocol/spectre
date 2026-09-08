package subscriber

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gogo/protobuf/proto"
)

func TestEthPacketToCosmosPacket_SinglePayload(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "eth-client-0",
		DestClient:       "cosmos-client-0",
		TimeoutTimestamp: 1700000000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
			{
				SourcePort: "transfer",
				DestPort:   "transfer",
				Version:    "ics20-1",
				Encoding:   "proto3",
				Value:      []byte("test-value"),
			},
		},
	}

	seq := big.NewInt(42)
	result := EthPacketToCosmosPacket(ethPacket, seq)

	if result.Sequence != 42 {
		t.Errorf("Sequence: got %d, want 42", result.Sequence)
	}
	if result.SourceClient != "eth-client-0" {
		t.Errorf("SourceClient: got %q, want %q", result.SourceClient, "eth-client-0")
	}
	if result.DestinationClient != "cosmos-client-0" {
		t.Errorf("DestinationClient: got %q, want %q", result.DestinationClient, "cosmos-client-0")
	}
	if result.TimeoutTimestamp != 1700000000 {
		t.Errorf("TimeoutTimestamp: got %d, want 1700000000", result.TimeoutTimestamp)
	}
	if len(result.Payloads) != 1 {
		t.Fatalf("Payloads length: got %d, want 1", len(result.Payloads))
	}

	p := result.Payloads[0]
	if p.SourcePort != "transfer" {
		t.Errorf("SourcePort: got %q", p.SourcePort)
	}
	if p.DestinationPort != "transfer" {
		t.Errorf("DestinationPort: got %q", p.DestinationPort)
	}
	if p.Version != "ics20-1" {
		t.Errorf("Version: got %q", p.Version)
	}
	if p.Encoding != "proto3" {
		t.Errorf("Encoding: got %q", p.Encoding)
	}
	if string(p.Value) != "test-value" {
		t.Errorf("Value: got %q", p.Value)
	}
}

func TestEthPacketToCosmosPacket_MultiplePayloads(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "src",
		DestClient:       "dst",
		TimeoutTimestamp: 999,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
			{SourcePort: "port1", DestPort: "port1d", Version: "v1", Encoding: "e1", Value: []byte("val1")},
			{SourcePort: "port2", DestPort: "port2d", Version: "v2", Encoding: "e2", Value: []byte("val2")},
			{SourcePort: "port3", DestPort: "port3d", Version: "v3", Encoding: "e3", Value: []byte("val3")},
		},
	}

	result := EthPacketToCosmosPacket(ethPacket, big.NewInt(1))

	if len(result.Payloads) != 3 {
		t.Fatalf("Payloads length: got %d, want 3", len(result.Payloads))
	}
	for i, p := range result.Payloads {
		if p.SourcePort == "" {
			t.Errorf("payload[%d] SourcePort is empty", i)
		}
		if p.DestinationPort == "" {
			t.Errorf("payload[%d] DestinationPort is empty", i)
		}
	}
}

func TestEthPacketToCosmosPacket_EmptyPayloads(t *testing.T) {
	ethPacket := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient:     "src",
		DestClient:       "dst",
		TimeoutTimestamp: 0,
		Payloads:         []contractICS26Router.IICS26RouterMsgsPayload{},
	}

	result := EthPacketToCosmosPacket(ethPacket, big.NewInt(0))

	if len(result.Payloads) != 0 {
		t.Errorf("expected 0 payloads, got %d", len(result.Payloads))
	}
}

func TestNormalizeTimeoutSeconds(t *testing.T) {
	cases := []struct {
		input    uint64
		expected uint64
	}{
		{1_700_000_000_000_000_000, 1_700_000_000}, // 2023 ns → 2023 seconds
		{1_000_000_000_000, 1_000},                 // nanoseconds → seconds
		{1_000, 1_000},                             // already seconds (year 2003), unchanged
		{0, 0},                                     // no timeout, unchanged
	}

	for _, c := range cases {
		got := normalizeTimeoutSeconds(c.input)
		if got != c.expected {
			t.Errorf("normalizeTimeoutSeconds(%d) = %d, want %d", c.input, got, c.expected)
		}
	}
}

func TestTxHeightFromEvent(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		data   commettypes.TMEventData
		events map[string][]string
		want   uint64
	}{
		{
			name:   "uses tx result height when available",
			data:   commettypes.EventDataTx{TxResult: abcitypes.TxResult{Height: 77}},
			events: map[string][]string{EVENT_TX_HEIGHT_FIELD: {"12"}},
			want:   77,
		},
		{
			name:   "falls back to tx height event field",
			events: map[string][]string{EVENT_TX_HEIGHT_FIELD: {"77"}},
			want:   77,
		},
		{
			name:   "missing tx height returns zero",
			events: map[string][]string{},
			want:   0,
		},
		{
			name:   "invalid tx height returns zero",
			events: map[string][]string{EVENT_TX_HEIGHT_FIELD: {"not-a-height"}},
			want:   0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := txHeightFromEvent(tc.data, tc.events); got != tc.want {
				t.Fatalf("txHeightFromEvent(%v, %v) = %d, want %d", tc.data, tc.events, got, tc.want)
			}
		})
	}
}

func TestEthStartupRecoveryLookbackBlocksFromEnv(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want uint64
	}{
		{name: "empty uses default", raw: "", want: defaultEthStartupRecoveryLookbackBlocks},
		{name: "invalid uses default", raw: "not-a-number", want: defaultEthStartupRecoveryLookbackBlocks},
		{name: "zero allowed", raw: "0", want: 0},
		{name: "custom value", raw: "1024", want: 1024},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ethStartupRecoveryLookbackBlocksFromEnv(tc.raw)
			if got != tc.want {
				t.Fatalf("ethStartupRecoveryLookbackBlocksFromEnv(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

func TestEthStartupRecoveryStartBlock(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		latest    uint64
		lookback  uint64
		wantStart uint64
	}{
		{name: "lookback smaller than latest", latest: 1000, lookback: 256, wantStart: 744},
		{name: "lookback equals latest", latest: 256, lookback: 256, wantStart: 0},
		{name: "lookback greater than latest", latest: 42, lookback: 100, wantStart: 0},
		{name: "zero lookback", latest: 42, lookback: 0, wantStart: 42},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ethStartupRecoveryStartBlock(tc.latest, tc.lookback)
			if got != tc.wantStart {
				t.Fatalf("ethStartupRecoveryStartBlock(%d, %d) = %d, want %d",
					tc.latest, tc.lookback, got, tc.wantStart)
			}
		})
	}
}

func TestCosmosStartupRecoveryLookbackBlocksFromEnv(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want uint64
	}{
		{name: "empty uses bounded default", raw: "", want: 256},
		{name: "invalid uses bounded default", raw: "not-a-number", want: 256},
		{name: "zero means full history", raw: "0", want: 0},
		{name: "custom bounded value", raw: "5000", want: 5000},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := cosmosStartupRecoveryLookbackBlocksFromEnv(tc.raw)
			if got != tc.want {
				t.Fatalf("cosmosStartupRecoveryLookbackBlocksFromEnv(%q) = %d, want %d", tc.raw, got, tc.want)
			}
		})
	}
}

func TestCosmosStartupRecoveryStartHeight(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		latestHeight uint64
		lookback     uint64
		wantStart    uint64
	}{
		{name: "no latest height", latestHeight: 0, lookback: 0, wantStart: 0},
		{name: "zero lookback scans full history", latestHeight: 100, lookback: 0, wantStart: 1},
		{name: "lookback greater than latest scans from first block", latestHeight: 42, lookback: 100, wantStart: 1},
		{name: "bounded lookback is inclusive", latestHeight: 100, lookback: 10, wantStart: 91},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := cosmosStartupRecoveryStartHeight(tc.latestHeight, tc.lookback)
			if got != tc.wantStart {
				t.Fatalf("cosmosStartupRecoveryStartHeight(%d, %d) = %d, want %d",
					tc.latestHeight, tc.lookback, got, tc.wantStart)
			}
		})
	}
}

func TestCosmosTxSearchQuery(t *testing.T) {
	t.Parallel()

	got := cosmosTxSearchQuery(cometBFTSendPacketTxSearch, 0, 25)
	want := "send_packet.encoded_packet_hex EXISTS AND tx.height >= 1 AND tx.height <= 25"
	if got != want {
		t.Fatalf("cosmosTxSearchQuery = %q, want %q", got, want)
	}
}

func TestCosmosEventsFromTxResult(t *testing.T) {
	t.Parallel()

	tx := &coretypes.ResultTx{
		Height: 66,
		Index:  3,
		Tx:     []byte("tx-bytes"),
		TxResult: abcitypes.ExecTxResult{
			Events: []abcitypes.Event{
				{
					Type: "send_packet",
					Attributes: []abcitypes.EventAttribute{
						{Key: "encoded_packet_hex", Value: "abc"},
					},
				},
			},
		},
	}

	data, events := cosmosEventsFromTxResult(tx)
	if got := events[EVENT_TX_HEIGHT_FIELD]; len(got) != 1 || got[0] != "66" {
		t.Fatalf("tx.height event = %v, want [66]", got)
	}
	if got := events[EVENT_SEND_PACKET_FIELD]; len(got) != 1 || got[0] != "abc" {
		t.Fatalf("send packet event = %v, want [abc]", got)
	}
	if got := txHeightFromEvent(data, events); got != 66 {
		t.Fatalf("txHeightFromEvent = %d, want 66", got)
	}
}

func TestCosmosPacketMatchesConfiguredClient(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		ethClientID    string
		routerClientID string
		packet         *channeltypesv2.Packet
		want           bool
	}{
		{
			name: "nil packet does not match",
			want: false,
		},
		{
			name: "no configured IDs preserves legacy behavior",
			packet: &channeltypesv2.Packet{
				SourceClient:      "unrelated-src",
				DestinationClient: "unrelated-dst",
			},
			want: true,
		},
		{
			name:        "source matches eth client",
			ethClientID: "08-wasm-0",
			packet: &channeltypesv2.Packet{
				SourceClient:      "08-wasm-0",
				DestinationClient: "cosmoshub-1",
			},
			want: true,
		},
		{
			name:        "destination matches eth client",
			ethClientID: "08-wasm-0",
			packet: &channeltypesv2.Packet{
				SourceClient:      "cosmoshub-1",
				DestinationClient: "08-wasm-0",
			},
			want: true,
		},
		{
			name:           "source matches router client",
			routerClientID: "cosmoshub-1",
			packet: &channeltypesv2.Packet{
				SourceClient:      "cosmoshub-1",
				DestinationClient: "08-wasm-0",
			},
			want: true,
		},
		{
			name:           "destination matches router client",
			routerClientID: "cosmoshub-1",
			packet: &channeltypesv2.Packet{
				SourceClient:      "08-wasm-0",
				DestinationClient: "cosmoshub-1",
			},
			want: true,
		},
		{
			name:           "configured IDs skip unrelated packet",
			ethClientID:    "08-wasm-0",
			routerClientID: "cosmoshub-1",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-src",
				DestinationClient: "other-dst",
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := cosmosDeps{IDs: services.ClientIDs{EVMOnCosmos: tc.ethClientID, CosmosOnEVM: tc.routerClientID}}

			if got := cosmosPacketMatchesConfiguredClient(ctx, tc.packet); got != tc.want {
				t.Fatalf("cosmosPacketMatchesConfiguredClient = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestEthEventClientIDFilter(t *testing.T) {
	t.Parallel()

	ctx := ethDeps{}
	if got := ethEventClientIDFilter(ctx); got != nil {
		t.Fatalf("ethEventClientIDFilter without router client ID = %v, want nil", got)
	}

	ctx.IDs.CosmosOnEVM = "simd-2-client-0"
	got := ethEventClientIDFilter(ctx)
	if len(got) != 1 || got[0] != "simd-2-client-0" {
		t.Fatalf("ethEventClientIDFilter = %v, want [simd-2-client-0]", got)
	}
}

func TestDecodeCosmosPacketsFromEvents(t *testing.T) {
	t.Parallel()

	logger := log.New(io.Discard, "", 0)
	sendPacketHex := mustPacketHex(t, channeltypesv2.Packet{
		Sequence:          11,
		SourceClient:      "cosmos-client",
		DestinationClient: "eth-client",
		TimeoutTimestamp:  1_700_000_000_000_000_000,
	})
	ackPacketHex := mustPacketHex(t, channeltypesv2.Packet{
		Sequence:          12,
		SourceClient:      "eth-client",
		DestinationClient: "cosmos-client",
		TimeoutTimestamp:  1234,
	})
	timeoutPacketHex := mustPacketHex(t, channeltypesv2.Packet{
		Sequence:          13,
		SourceClient:      "eth-client",
		DestinationClient: "other-client",
	})
	settledPacketHex := mustPacketHex(t, channeltypesv2.Packet{
		Sequence:          14,
		SourceClient:      "cosmos-client",
		DestinationClient: "eth-client",
	})
	ackHex := mustAckHex(t, channeltypesv2.Acknowledgement{
		AppAcknowledgements: [][]byte{[]byte("ack")},
	})

	packets := decodeCosmosPacketsFromEvents(
		logger,
		commettypes.EventDataTx{TxResult: abcitypes.TxResult{Height: 55}},
		map[string][]string{
			EVENT_SEND_PACKET_FIELD:      {sendPacketHex},
			EVENT_WRITE_ACK_PACKET_FIELD: {ackPacketHex},
			EVENT_ACKNOWLEDGEMENT_FIELD:  {ackHex},
			EVENT_TIMEOUT_PACKET_FIELD:   {timeoutPacketHex},
			EVENT_ACK_PACKET_FIELD:       {settledPacketHex},
			EVENT_TX_HEIGHT_FIELD:        {"44"},
		},
		"test",
	)

	if len(packets) != 4 {
		t.Fatalf("decoded packet count = %d, want 4", len(packets))
	}
	if packets[0].Type != services.CosmosSend || packets[0].Packet.Sequence != 11 {
		t.Fatalf("packet[0] = type %v seq %d, want CosmosSend seq 11", packets[0].Type, packets[0].Packet.Sequence)
	}
	if packets[0].Packet.TimeoutTimestamp != 1_700_000_000 {
		t.Fatalf("normalized timeout = %d, want 1700000000", packets[0].Packet.TimeoutTimestamp)
	}
	if packets[1].Type != services.CosmosAck || packets[1].Packet.Sequence != 12 {
		t.Fatalf("packet[1] = type %v seq %d, want CosmosAck seq 12", packets[1].Type, packets[1].Packet.Sequence)
	}
	if len(packets[1].AckBytes) != 1 || string(packets[1].AckBytes[0]) != "ack" {
		t.Fatalf("ack bytes = %q, want [ack]", packets[1].AckBytes)
	}
	if packets[2].Type != services.CosmosTimeout || packets[2].Packet.Sequence != 13 {
		t.Fatalf("packet[2] = type %v seq %d, want CosmosTimeout seq 13", packets[2].Type, packets[2].Packet.Sequence)
	}
	// acknowledge_packet is the fourth event and the only terminal one on this
	// side: Cosmos consumed the ack for a packet it sent, so that packet is done.
	// Without it a packet another relayer acknowledged stays in the pending
	// tracker until a timeout scan happens to query for it.
	if packets[3].Type != services.CosmosAcknowledged || packets[3].Packet.Sequence != 14 {
		t.Fatalf("packet[3] = type %v seq %d, want CosmosAcknowledged seq 14", packets[3].Type, packets[3].Packet.Sequence)
	}
	for i, p := range packets {
		if p.BlockNumber != 55 {
			t.Fatalf("packet[%d] block = %d, want 55", i, p.BlockNumber)
		}
	}
}

func TestAdvanceRecoveryStart(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		initial   uint64
		candidate uint64
		want      uint64
	}{
		{name: "advance forward", initial: 10, candidate: 15, want: 15},
		{name: "equal keeps current", initial: 10, candidate: 10, want: 10},
		{name: "smaller keeps current", initial: 10, candidate: 5, want: 10},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.initial
			advanceRecoveryStart(&got, tc.candidate)
			if got != tc.want {
				t.Fatalf("advanceRecoveryStart(%d, %d) = %d, want %d", tc.initial, tc.candidate, got, tc.want)
			}
		})
	}
}

func TestEnqueueEthSendPacket(t *testing.T) {
	t.Parallel()

	bb := services.NewBatchBuilder()
	ev := &contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(7),
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient:     "eth-client-0",
			DestClient:       "cosmos-client-0",
			TimeoutTimestamp: 1234,
			Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
				{
					SourcePort: "transfer",
					DestPort:   "transfer",
					Version:    "ics20-1",
					Encoding:   "proto3",
					Value:      []byte("payload"),
				},
			},
		},
		Raw: gethtypes.Log{BlockNumber: 88},
	}

	if _, err := enqueueEthSendPacket(bb, ev, nil); err != nil {
		t.Fatal(err)
	}

	got := flushSingleEthPacket(t, bb)
	if got.Type != services.EthSend {
		t.Fatalf("packet type = %v, want %v", got.Type, services.EthSend)
	}
	if got.Packet.Sequence != 7 {
		t.Fatalf("packet sequence = %d, want 7", got.Packet.Sequence)
	}
	if got.BlockNumber != 88 {
		t.Fatalf("block number = %d, want 88", got.BlockNumber)
	}
	if bb.EthPendingTracker.Len() != 1 {
		t.Fatalf("eth pending tracker len = %d, want 1", bb.EthPendingTracker.Len())
	}
	pending := bb.EthPendingTracker.GetAll()
	if pending[0].Packet.Sequence != 7 || pending[0].BlockNumber != 88 {
		t.Fatalf("eth pending entry = seq %d block %d, want seq 7 block 88",
			pending[0].Packet.Sequence, pending[0].BlockNumber)
	}
}

func TestEnqueueEthSendPacketDedupesSeenEvent(t *testing.T) {
	t.Parallel()

	bb := services.NewBatchBuilder()
	ev := &contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(7),
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient: "eth-client-0",
			DestClient:   "cosmos-client-0",
		},
		Raw: gethtypes.Log{BlockNumber: 88, Index: 3},
	}
	seen := make(map[ethEventKey]struct{})

	if enqueued, err := enqueueEthSendPacket(bb, ev, seen); err != nil || !enqueued {
		t.Fatal("first enqueue should be accepted")
	}
	if enqueued, err := enqueueEthSendPacket(bb, ev, seen); err != nil || enqueued {
		t.Fatal("duplicate enqueue should be skipped")
	}

	got := flushSingleEthPacket(t, bb)
	if got.Packet.Sequence != 7 {
		t.Fatalf("packet sequence = %d, want 7", got.Packet.Sequence)
	}
	if bb.EthPendingTracker.Len() != 1 {
		t.Fatalf("eth pending tracker len = %d, want 1", bb.EthPendingTracker.Len())
	}
}

func TestEnqueueEthSendPacketRetriesWhenPendingStateCannotPersist(t *testing.T) {
	dir := t.TempDir()
	bb, err := services.NewPersistentBatchBuilder(dir)
	if err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(dir, "eth.json")
	if err := os.Remove(statePath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(statePath, 0o700); err != nil {
		t.Fatal(err)
	}

	ev := &contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(7),
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient: "eth-client-0",
			DestClient:   "cosmos-client-0",
		},
		Raw: gethtypes.Log{BlockNumber: 88, Index: 3},
	}
	seen := make(map[ethEventKey]struct{})

	if enqueued, err := enqueueEthSendPacket(bb, ev, seen); err == nil || enqueued {
		t.Fatal("enqueue succeeded after pending-state persistence failed")
	}
	if len(seen) != 0 {
		t.Fatal("failed pending-state write marked the source event seen")
	}
	if bb.EthPendingTracker.Len() != 0 {
		t.Fatal("failed pending-state write left the packet tracked in memory")
	}

	ch := make(chan services.EthBatch, 1)
	bb.CheckEth(context.Background(), services.BatchConfig{BatchSize: 1, BatchPeriods: time.Hour}, ch)
	select {
	case batch := <-ch:
		t.Fatalf("failed pending-state write still enqueued relay batch: %+v", batch)
	default:
	}
}

func TestFailedEthPendingAddKeepsRecoveryCursorRetryable(t *testing.T) {
	dir := t.TempDir()
	bb, err := services.NewPersistentBatchBuilder(dir)
	if err != nil {
		t.Fatal(err)
	}
	statePath := filepath.Join(dir, "eth.json")
	if err := os.Remove(statePath); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(statePath, 0o700); err != nil {
		t.Fatal(err)
	}

	router := common.HexToAddress("0x00000000000000000000000000000000000000AA")
	server := ethRecoveryRPCServer(t, ethRecoverySendPacketLog(t, router))
	ethClient, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(ethClient.Close)
	cosmosClient, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatal(err)
	}
	filterer, err := contractICS26Router.NewContractICS26RouterFilterer(router, ethClient)
	if err != nil {
		t.Fatal(err)
	}

	cursor := uint64(88)
	sub := NewSubscriber(nil)
	deps := ethDeps{
		Cosmos: services.CosmosEndpoint{Client: cosmosClient},
		EVM: services.EVMEndpoint{
			Client:    ethClient,
			Contracts: services.EVMContracts{Router: router},
		},
		Logger: log.New(io.Discard, "", 0),
	}
	_, err = sub.scanEthRangeInChunks(context.Background(), deps, "SendPacket", &cursor, 90, 3,
		func(from, to uint64) (ethRecoveryStats, error) {
			return recoverEthSendPackets(context.Background(), deps, bb, filterer, from, to, nil)
		}, func() {})
	if err == nil {
		t.Fatal("failed durable pending add did not fail the recovery chunk")
	}
	if cursor != 88 {
		t.Fatalf("recovery cursor advanced to %d after failed durable add, want 88", cursor)
	}
	if bb.EthPendingTracker.Len() != 0 {
		t.Fatal("failed durable pending add left packet tracked")
	}
}

func ethRecoveryRPCServer(t *testing.T, sendPacketLog gethtypes.Log) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var result any
		switch request.Method {
		case "eth_getLogs":
			result = []gethtypes.Log{sendPacketLog}
		case "eth_getStorageAt":
			// The pending commitment exists, so recovery reaches the durable Add.
			result = "0x01"
		case "abci_query":
			// No Cosmos receipt exists for the SendPacket yet.
			result = map[string]any{"response": map[string]any{"code": 0, "value": ""}}
		default:
			http.Error(w, "unexpected RPC method: "+request.Method, http.StatusNotImplemented)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"jsonrpc": "2.0",
			"id":      request.ID,
			"result":  result,
		}); err != nil {
			t.Error(err)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func ethRecoverySendPacketLog(t *testing.T, router common.Address) gethtypes.Log {
	t.Helper()
	parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	event := parsed.Events["SendPacket"]
	packet := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient: "eth-client-0",
		DestClient:   "cosmos-client-0",
	}
	data, err := event.Inputs.NonIndexed().Pack(packet)
	if err != nil {
		t.Fatal(err)
	}
	return gethtypes.Log{
		Address: router,
		Topics: []common.Hash{
			event.ID,
			crypto.Keccak256Hash([]byte(packet.SourceClient)),
			common.BigToHash(big.NewInt(8)),
		},
		Data:        data,
		BlockNumber: 88,
		Index:       4,
	}
}

func TestEnqueueEthWriteAcknowledgement(t *testing.T) {
	t.Parallel()

	bb := services.NewBatchBuilder()
	ev := &contractICS26Router.ContractICS26RouterWriteAcknowledgement{
		Sequence:         big.NewInt(9),
		Acknowledgements: [][]byte{[]byte("ack")},
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient:     "eth-client-0",
			DestClient:       "cosmos-client-0",
			TimeoutTimestamp: 5678,
			Payloads: []contractICS26Router.IICS26RouterMsgsPayload{
				{
					SourcePort: "transfer",
					DestPort:   "transfer",
					Version:    "ics20-1",
					Encoding:   "proto3",
					Value:      []byte("payload"),
				},
			},
		},
		Raw: gethtypes.Log{BlockNumber: 99},
	}

	enqueueEthWriteAcknowledgement(bb, ev, nil)

	got := flushSingleEthPacket(t, bb)
	if got.Type != services.EthWriteAck {
		t.Fatalf("packet type = %v, want %v", got.Type, services.EthWriteAck)
	}
	if got.Packet.Sequence != 9 {
		t.Fatalf("packet sequence = %d, want 9", got.Packet.Sequence)
	}
	if got.BlockNumber != 99 {
		t.Fatalf("block number = %d, want 99", got.BlockNumber)
	}
	if len(got.AckBytes) != 1 || string(got.AckBytes[0]) != "ack" {
		t.Fatalf("ack bytes = %q, want [ack]", got.AckBytes)
	}
}

func flushSingleEthPacket(t *testing.T, bb *services.BatchBuilder) services.EthPacket {
	t.Helper()

	ch := make(chan services.EthBatch, 1)
	bb.CheckEth(context.Background(), services.BatchConfig{
		BatchSize:    1,
		BatchPeriods: time.Hour,
	}, ch)

	select {
	case batch := <-ch:
		if len(batch.Packets) != 1 {
			t.Fatalf("batch packet count = %d, want 1", len(batch.Packets))
		}
		return batch.Packets[0]
	default:
		t.Fatal("expected one flushed eth batch")
		return services.EthPacket{}
	}
}

func mustPacketHex(t *testing.T, packet channeltypesv2.Packet) string {
	t.Helper()
	bz, err := proto.Marshal(&packet)
	if err != nil {
		t.Fatalf("failed to marshal packet: %v", err)
	}
	return hex.EncodeToString(bz)
}

func mustAckHex(t *testing.T, ack channeltypesv2.Acknowledgement) string {
	t.Helper()
	bz, err := proto.Marshal(&ack)
	if err != nil {
		t.Fatalf("failed to marshal acknowledgement: %v", err)
	}
	return hex.EncodeToString(bz)
}

// The bug this pins, found in review: the ETH driver folded only the SendPacket
// branch's stats into `found` and discarded the WriteAcknowledgement branch's, so
// a pass that recovered acks still counted as quiet — and after 20 of them printed
// "found nothing" directly beside the line reporting what it found. Both branches
// now go through foundSomething(), and any future branch must too.
func TestEthRecoveryStatsFoundSomething(t *testing.T) {
	for _, tc := range []struct {
		name  string
		stats ethRecoveryStats
		want  bool
	}{
		{"nothing", ethRecoveryStats{}, false},
		{"recovered", ethRecoveryStats{recovered: 1}, true},
		{"skipped only", ethRecoveryStats{skipped: 1}, true},
		{"both", ethRecoveryStats{recovered: 2, skipped: 3}, true},
	} {
		if got := tc.stats.foundSomething(); got != tc.want {
			t.Errorf("%s: foundSomething() = %v, want %v", tc.name, got, tc.want)
		}
	}

	// The aggregation itself: a quiet SendPacket pass must not mask a productive
	// WriteAcknowledgement one.
	found := false
	found = found || ethRecoveryStats{}.foundSomething()
	found = found || ethRecoveryStats{recovered: 1}.foundSomething()
	if !found {
		t.Fatal("a pass that recovered WriteAcknowledgements must not count as quiet")
	}
}

// The heartbeat is shared by both directions so they cannot drift: one line every
// quietScanHeartbeat passes, and the counter keeps climbing so the printed number
// is the real streak.
func TestAdvanceQuietScans(t *testing.T) {
	var quiet uint64
	beats := 0
	for i := 1; i <= quietScanHeartbeat*2; i++ {
		if advanceQuietScans(&quiet) {
			beats++
			if quiet%quietScanHeartbeat != 0 {
				t.Fatalf("beat at %d, which is not a multiple of %d", quiet, quietScanHeartbeat)
			}
		}
	}
	if beats != 2 {
		t.Fatalf("beats = %d over %d passes, want 2", beats, quietScanHeartbeat*2)
	}
	if quiet != uint64(quietScanHeartbeat*2) {
		t.Fatalf("counter = %d, want it to keep counting to %d", quiet, quietScanHeartbeat*2)
	}
}

// The scan must never read up to the chain head. Status reports a height as soon
// as the block commits, but TxSearch reads the tx indexer, which CometBFT writes
// from a separate goroutine after commit — so the newest block is routinely
// committed and not yet indexed. Since the scan advances its cursor past
// whatever it read, a tx there is skipped permanently and silently.
func TestCosmosIndexedHeightStopsShortOfTheHead(t *testing.T) {
	t.Parallel()

	cases := []struct {
		latest uint64
		want   uint64
	}{
		{0, 0},
		{1, 0},
		{cosmosIndexerLagBlocks, 0},
		{cosmosIndexerLagBlocks + 1, 1},
		{202_197, 202_197 - cosmosIndexerLagBlocks},
	}
	for _, tc := range cases {
		got := cosmosIndexedHeight(tc.latest)
		if got != tc.want {
			t.Fatalf("cosmosIndexedHeight(%d) = %d, want %d", tc.latest, got, tc.want)
		}
	}

	// The property that actually protects packets: the range one pass scans must
	// never reach the head, at any cursor position. Asserted through
	// cosmosScanRange -- the function the scanner calls -- so removing the lag
	// from the scan fails here, which a test of cosmosIndexedHeight alone would
	// not.
	const latest = 202_197
	for cursor := uint64(1); cursor <= latest; cursor += 977 {
		from, to, ok := cosmosScanRange(cursor, latest)
		if !ok {
			continue
		}
		if to >= latest {
			t.Fatalf("cosmosScanRange(%d, %d) scans to %d, which reaches the head; "+
				"a tx in the newest block may be committed but not yet indexed, and the cursor "+
				"advances past it permanently", cursor, latest, to)
		}
		if from != cursor {
			t.Fatalf("cosmosScanRange(%d, %d) starts at %d; the cursor must not be skipped forward",
				cursor, latest, from)
		}
	}

	// Nothing scannable yet must be reported, not silently scanned as [x, 0].
	if _, _, ok := cosmosScanRange(0, latest); ok {
		t.Fatal("an unseeded cursor must not produce a scan range")
	}
	if _, _, ok := cosmosScanRange(latest, latest); ok {
		t.Fatal("a cursor at the head has nothing safely indexed to scan yet")
	}
}

// A live subscription that stops delivering is invisible by construction: it
// looks exactly like a quiet chain. These pin the two signals that tell them
// apart -- the chain still producing blocks, and how long the silence has run.
func TestLivePathStale(t *testing.T) {
	t.Parallel()

	// fresh returns a health tracker whose delivery clock starts at now, as
	// subscribeCosmosOnce sets it immediately after subscribing.
	fresh := func(now time.Time) *cosmosLiveHealth {
		h := &cosmosLiveHealth{}
		h.recordSubscribed(now)
		return h
	}
	start := time.Unix(1_700_000_000, 0)

	// An idle chain delivers nothing because there is nothing to deliver.
	// Reconnecting here would tear down a healthy subscription every few minutes.
	t.Run("a chain that is not advancing is never stale", func(t *testing.T) {
		h := fresh(start)
		stale, _, report := h.livePathStale(false, start.Add(24*time.Hour))
		if stale || report {
			t.Fatalf("stale=%v report=%v on an idle chain; silence proves nothing without block production", stale, report)
		}
	})

	t.Run("silence inside the window is not yet evidence", func(t *testing.T) {
		h := fresh(start)
		if stale, _, _ := h.livePathStale(true, start.Add(cosmosLivePathStaleAfter-time.Second)); stale {
			t.Fatal("reported a dead subscription while it was only briefly quiet")
		}
	})

	// The case the watchdog exists for. Note the clock runs from the SUBSCRIBE,
	// so a subscription that never delivered a single event is caught -- the old
	// version measured from the first live event and so had nothing to measure.
	t.Run("a subscription that never delivers is caught", func(t *testing.T) {
		h := fresh(start)
		stale, silence, report := h.livePathStale(true, start.Add(cosmosLivePathStaleAfter+time.Second))
		if !stale || !report {
			t.Fatalf("stale=%v report=%v; a subscription silent past the window while the chain advances must reconnect and say so", stale, report)
		}
		if silence < cosmosLivePathStaleAfter {
			t.Fatalf("silence = %s, want at least %s: the line has to carry the real number", silence, cosmosLivePathStaleAfter)
		}
	})

	// The half that was missing. A latched bool cleared only by a live event
	// means an outage that never ends produces exactly one line for the whole
	// life of the process -- and the log then reads as though it recovered.
	t.Run("a permanent outage keeps reporting on the ladder", func(t *testing.T) {
		h := fresh(start)
		now := start.Add(cosmosLivePathStaleAfter)
		if _, _, report := h.livePathStale(true, now); !report {
			t.Fatal("the first stale pass must report")
		}

		// Every 30s tick up to (not including) the five-minute step: same
		// threshold, so every one of them must be silent.
		staleSince := now
		for tick := 30 * time.Second; tick < 5*time.Minute; tick += 30 * time.Second {
			if _, _, report := h.livePathStale(true, staleSince.Add(tick)); report {
				t.Fatalf("repeated the warning %s into one threshold; that is 120 lines an hour", tick)
			}
		}

		// Crossing 5m, then 15m, then an hour each release exactly one more line.
		for _, step := range []time.Duration{5 * time.Minute, 15 * time.Minute, time.Hour, 2 * time.Hour} {
			if _, _, report := h.livePathStale(true, start.Add(cosmosLivePathStaleAfter+step)); !report {
				t.Fatalf("an outage still dead at %s said nothing; silence reads as recovery", step)
			}
		}
	})

	t.Run("a live event ends the outage and resets the ladder", func(t *testing.T) {
		h := fresh(start)
		h.livePathStale(true, start.Add(time.Hour)) // outage under way, thresholds burnt
		h.recordEvent()
		if !h.staleSince.IsZero() || h.reported != 0 {
			t.Fatalf("a delivered event must end the outage: staleSince=%v reported=%d", h.staleSince, h.reported)
		}
		if stale, _, _ := h.livePathStale(true, h.lastSeen.Add(time.Second)); stale {
			t.Fatal("a subscription that just delivered must not be torn down")
		}
		// And a NEW outage after it reports from the first step again.
		if _, _, report := h.livePathStale(true, h.lastSeen.Add(cosmosLivePathStaleAfter+time.Second)); !report {
			t.Fatal("a fresh outage must report again rather than inherit the old one's ladder")
		}
	})

	// Reconnecting is an attempt to fix the outage, not evidence it is over.
	// Resetting the ladder here would make a reconnect loop silent after one line.
	t.Run("resubscribing does not end the outage", func(t *testing.T) {
		h := fresh(start)
		h.livePathStale(true, start.Add(cosmosLivePathStaleAfter))
		h.recordSubscribed(start.Add(cosmosLivePathStaleAfter))
		if h.staleSince.IsZero() {
			t.Fatal("resubscribing cleared the outage; only a delivered event may do that")
		}
	})
}

// Block production is what makes silence meaningful, so the baseline pass must
// not claim an advance it cannot know about.
func TestObserveChainHeight(t *testing.T) {
	t.Parallel()

	var h cosmosLiveHealth
	if h.observeChainHeight(100) {
		t.Fatal("the first observation has nothing to compare against and must not report an advance")
	}
	if h.observeChainHeight(100) {
		t.Fatal("an unchanged head is not an advance")
	}
	if !h.observeChainHeight(101) {
		t.Fatal("a rising head is the signal that the chain is producing")
	}
}

func TestAgeReportThreshold(t *testing.T) {
	t.Parallel()

	cases := []struct {
		age  time.Duration
		want int
		why  string
	}{
		{0, 1, "the first stale pass always reports"},
		{4 * time.Minute, 1, "no new information yet"},
		{5 * time.Minute, 2, ""},
		{15 * time.Minute, 3, ""},
		{time.Hour, 4, ""},
		{3 * time.Hour, 6, "past the ladder it keeps counting hourly, so an overnight outage leaves a trail"},
	}
	for _, tc := range cases {
		if got := ageReportThreshold(tc.age); got != tc.want {
			t.Fatalf("ageReportThreshold(%s) = %d, want %d: %s", tc.age, got, tc.want, tc.why)
		}
	}
}

// TestEnqueueEthTerminalRecordsBlockNumber pins that a settling event carries
// the block it was observed at. Both enqueues used to omit it, which is
// invisible today only because the source bridge drops these types before
// reading Height — a trap for the next change, not a live bug.
func TestEnqueueEthTerminalRecordsBlockNumber(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		typ      services.EthPacketType
		ackBytes [][]byte
	}{
		{"ack", services.EthAck, [][]byte{[]byte("ack")}},
		{"timeout", services.EthTimeout, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			bb := services.NewBatchBuilder()
			pkt := contractICS26Router.IICS26RouterMsgsPacket{
				SourceClient: "eth-client-0",
				DestClient:   "cosmos-client-0",
			}
			enqueueEthTerminal(bb, log.New(io.Discard, "", 0), tc.typ, pkt, big.NewInt(4), tc.ackBytes, 4242)

			ch := make(chan services.EthBatch, 1)
			bb.CheckEth(context.Background(), services.BatchConfig{BatchSize: 1}, ch)

			select {
			case batch := <-ch:
				if len(batch.Packets) != 1 {
					t.Fatalf("got %d packets, want 1", len(batch.Packets))
				}
				if got := batch.Packets[0].BlockNumber; got != 4242 {
					t.Fatalf("BlockNumber = %d, want 4242", got)
				}
				if got := batch.Packets[0].Type; got != tc.typ {
					t.Fatalf("Type = %v, want %v", got, tc.typ)
				}
			default:
				t.Fatal("no batch flushed")
			}
		})
	}
}

// TestEnqueueEthTerminalSettlesThePendingTracker: a settling event must clear
// the tracker entry, or the timeout scanner keeps querying a packet that can
// never time out.
func TestEnqueueEthTerminalSettlesThePendingTracker(t *testing.T) {
	t.Parallel()

	bb := services.NewBatchBuilder()
	pkt := contractICS26Router.IICS26RouterMsgsPacket{
		SourceClient: "eth-client-0",
		DestClient:   "cosmos-client-0",
	}
	cosmosPacket := EthPacketToCosmosPacket(pkt, big.NewInt(4))
	bb.EthPendingTracker.Add(cosmosPacket, 100)
	if bb.EthPendingTracker.Len() != 1 {
		t.Fatalf("tracker not seeded")
	}

	enqueueEthTerminal(bb, log.New(io.Discard, "", 0), services.EthAck, pkt, big.NewInt(4), [][]byte{[]byte("ack")}, 4242)

	if got := bb.EthPendingTracker.Len(); got != 0 {
		t.Fatalf("tracker length = %d after settlement, want 0", got)
	}
}

// --- dedupe keys and their retention window ---
//
// These three functions are the whole of the subscriber's "have I already handled
// this event" memory. They had no test. What makes them worth pinning is not the
// bookkeeping but the two opposite failures they sit between: a key that is too
// coarse suppresses an event that genuinely needs relaying, and a window that
// prunes too eagerly lets an already-relayed event come back around.

// Pruning bounds the memory of a process that runs for weeks. The boundary is what
// matters: prune one height too far and an event still inside the reorg window is
// forgotten, so a re-delivery is relayed a second time.
func TestPruneCosmosSeenEvents_PrunesAtTheWindowBoundary(t *testing.T) {
	const current = cosmosSeenEventRetentionHeights + 1_000
	minKept := current - cosmosSeenEventRetentionHeights

	seen := map[cosmosEventKey]struct{}{
		{Sequence: 1, Height: minKept - 1}: {}, // just outside the window
		{Sequence: 2, Height: minKept}:     {}, // exactly on the boundary
		{Sequence: 3, Height: minKept + 1}: {}, // inside
		{Sequence: 4, Height: current}:     {}, // current height
		{Sequence: 5, Height: 0}:           {}, // unknown height
	}
	pruneCosmosSeenEvents(seen, current)

	for _, want := range []struct {
		seq  uint64
		h    uint64
		kept bool
		why  string
	}{
		{1, minKept - 1, false, "outside the retention window"},
		{2, minKept, true, "exactly on the boundary — the window is inclusive"},
		{3, minKept + 1, true, "inside the window"},
		{4, current, true, "the current height"},
		{5, 0, true, "unknown height: no way to tell if it aged out, so it is kept"},
	} {
		_, ok := seen[cosmosEventKey{Sequence: want.seq, Height: want.h}]
		if ok != want.kept {
			t.Errorf("height %d kept=%v, want %v (%s)", want.h, ok, want.kept, want.why)
		}
	}
}

// Below the retention window there is nothing old enough to prune, and the
// subtraction must not be attempted: this is unsigned arithmetic, so an unguarded
// currentHeight - retention on a young chain wraps to an enormous minimum and
// deletes every entry — on a chain that has produced nothing old enough to forget.
//
// The height chosen here is well below the window, not on its boundary. At exactly
// currentHeight == retention the subtraction yields 0 and prunes nothing anyway,
// so a test there passes with or without the guard and proves nothing.
func TestPruneCosmosSeenEvents_YoungChainKeepsEverything(t *testing.T) {
	const young = 10 // vs a retention window of thousands of heights
	if young >= cosmosSeenEventRetentionHeights {
		t.Fatalf("this test needs a height below the %d-height window", cosmosSeenEventRetentionHeights)
	}
	seen := map[cosmosEventKey]struct{}{
		{Sequence: 1, Height: 1}:     {},
		{Sequence: 2, Height: young}: {},
	}
	pruneCosmosSeenEvents(seen, young)
	if len(seen) != 2 {
		t.Fatalf("pruned %d of 2 entries on a chain younger than the retention window", 2-len(seen))
	}
}

// Mirror of the Cosmos side. The two windows are different sizes (blocks vs
// heights, different chains) but the boundary rule must be the same, or one
// direction re-relays where the other does not.
func TestPruneEthSeenEvents_PrunesAtTheWindowBoundary(t *testing.T) {
	const current = ethSeenEventRetentionBlocks + 1_000
	minKept := current - ethSeenEventRetentionBlocks

	seen := map[ethEventKey]struct{}{
		{TxHash: "0x1", BlockNumber: minKept - 1}: {},
		{TxHash: "0x2", BlockNumber: minKept}:     {},
		{TxHash: "0x3", BlockNumber: current}:     {},
		{TxHash: "0x4", BlockNumber: 0}:           {},
	}
	pruneEthSeenEvents(seen, current)

	for _, want := range []struct {
		hash string
		b    uint64
		kept bool
	}{
		{"0x1", minKept - 1, false},
		{"0x2", minKept, true},
		{"0x3", current, true},
		{"0x4", 0, true},
	} {
		_, ok := seen[ethEventKey{TxHash: want.hash, BlockNumber: want.b}]
		if ok != want.kept {
			t.Errorf("block %d kept=%v, want %v", want.b, ok, want.kept)
		}
	}
}

// Mirror of the Cosmos young-chain case, and for the same reason: below the window
// the subtraction wraps and takes every entry with it.
func TestPruneEthSeenEvents_YoungChainKeepsEverything(t *testing.T) {
	const young = 10
	if young >= ethSeenEventRetentionBlocks {
		t.Fatalf("this test needs a block below the %d-block window", ethSeenEventRetentionBlocks)
	}
	seen := map[ethEventKey]struct{}{
		{TxHash: "0x1", BlockNumber: 1}:     {},
		{TxHash: "0x2", BlockNumber: young}: {},
	}
	pruneEthSeenEvents(seen, young)
	if len(seen) != 2 {
		t.Fatalf("pruned %d of 2 entries on a chain younger than the retention window", 2-len(seen))
	}
}

// --- gap-recovery range guards ---
//
// Both recovery scans refuse an empty or backwards range before dialing. The
// guard is what keeps a cursor that has run ahead of the head from turning into a
// backwards scan; without it the range is passed to the RPC as-is, and what comes
// back is either nothing (silent event loss, the failure this whole path exists to
// prevent) or an error on every recovery tick.

func TestRecoverCosmosEvents_RefusesEmptyRangeBeforeDialing(t *testing.T) {
	// Zero-value deps have no Cosmos client: reaching the RPC would panic rather
	// than return the empty stats asserted here.
	for _, tt := range []struct {
		name       string
		start, end uint64
	}{
		{"end below start", 200, 100},
		{"end is zero", 0, 0},
		{"end is zero with a real start", 100, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			stats, err := recoverCosmosEvents(context.Background(), cosmosDeps{}, nil, tt.start, tt.end, nil)
			if err != nil {
				t.Fatalf("an empty range is not an error, got %v", err)
			}
			if stats.recovered != 0 || stats.skipped != 0 {
				t.Fatalf("reported work on an empty range: %+v", stats)
			}
		})
	}
}

// The ETH side scans SendPacket and WriteAcknowledgement in two independent
// functions, and each carries its own copy of the range guard. They get one test
// each, because "recoverEth" is not a unit -- there is no such function, and E9
// wants the prefix to name something real.
//
// The pair is what matters, though: a guard present on one scan and missing on
// the other is the one-sided asymmetry this repo keeps finding. Keeping them
// adjacent, over the same table, is what makes that visible. Change one, look at
// the other.
var ethBackwardsRangeScans = map[string]func(context.Context) (ethRecoveryStats, error){
	"send packets": func(ctx context.Context) (ethRecoveryStats, error) {
		return recoverEthSendPackets(ctx, ethDeps{}, nil, nil, 200, 100, nil)
	},
	"write acknowledgements": func(ctx context.Context) (ethRecoveryStats, error) {
		return recoverEthWriteAcknowledgements(ctx, ethDeps{}, nil, nil, 200, 100, nil)
	},
}

// assertRefusesBackwardsRange checks that the scan returns without dialing. A nil
// filterer means any attempt to scan would panic rather than return.
func assertRefusesBackwardsRange(t *testing.T, scan func(context.Context) (ethRecoveryStats, error)) {
	t.Helper()
	stats, err := scan(context.Background())
	if err != nil {
		t.Fatalf("a backwards range is not an error, got %v", err)
	}
	if stats.recovered != 0 {
		t.Fatalf("reported work on a backwards range: %+v", stats)
	}
}

func TestRecoverEthSendPackets_RefusesBackwardsRangeBeforeDialing(t *testing.T) {
	assertRefusesBackwardsRange(t, ethBackwardsRangeScans["send packets"])
}

func TestRecoverEthWriteAcknowledgements_RefusesBackwardsRangeBeforeDialing(t *testing.T) {
	assertRefusesBackwardsRange(t, ethBackwardsRangeScans["write acknowledgements"])
}

// The dedupe key decides whether an event is treated as one already handled.
// Too coarse and a real event is suppressed; the reorg case below is where
// that costs a packet.
func TestCosmosEventKeyForPacket(t *testing.T) {
	// A reorg re-emits the same logical packet at a different height. The key must
	// treat that as a NEW event, or the second emission is silently dropped as a
	// duplicate and the packet is never relayed from the height that actually stuck.
	//
	// Everything else about the packet being equal, height alone must change the key.
	t.Run("height separates reorged events", func(t *testing.T) {
		packet := func(height uint64) services.CosmosPacket {
			return services.CosmosPacket{
				Type:        services.CosmosSend,
				BlockNumber: height,
				Packet: &channeltypesv2.Packet{
					Sequence:          9,
					SourceClient:      "cosmos-client-0",
					DestinationClient: "eth-router-0",
				},
			}
		}

		if cosmosEventKeyForPacket(packet(100)) != cosmosEventKeyForPacket(packet(100)) {
			t.Fatal("the same event at the same height must produce the same key")
		}
		if cosmosEventKeyForPacket(packet(100)) == cosmosEventKeyForPacket(packet(101)) {
			t.Fatal("a reorged event re-emitted at a new height was treated as a duplicate")
		}
	})

	// Each identifying field must participate in the key on its own. A key that
	// ignores, say, the packet type would let a WriteAcknowledgement be swallowed by
	// the SendPacket already seen for the same sequence.
	t.Run("every identifying field counts", func(t *testing.T) {
		base := services.CosmosPacket{
			Type:        services.CosmosSend,
			BlockNumber: 100,
			Packet: &channeltypesv2.Packet{
				Sequence:          9,
				SourceClient:      "cosmos-client-0",
				DestinationClient: "eth-router-0",
			},
		}

		tests := []struct {
			field  string
			mutate func(p *services.CosmosPacket)
		}{
			{"packet type", func(p *services.CosmosPacket) { p.Type = services.CosmosAck }},
			{"sequence", func(p *services.CosmosPacket) { p.Packet.Sequence = 10 }},
			{"source client", func(p *services.CosmosPacket) { p.Packet.SourceClient = "cosmos-client-1" }},
			{"destination client", func(p *services.CosmosPacket) { p.Packet.DestinationClient = "eth-router-1" }},
			{"height", func(p *services.CosmosPacket) { p.BlockNumber = 101 }},
		}
		for _, tt := range tests {
			t.Run(tt.field, func(t *testing.T) {
				other := base
				inner := *base.Packet
				other.Packet = &inner
				tt.mutate(&other)
				if cosmosEventKeyForPacket(base) == cosmosEventKeyForPacket(other) {
					t.Fatalf("%s does not affect the dedupe key; two distinct events collide", tt.field)
				}
			})
		}
	})

	// A packet that failed to decode carries no identity. The key must still be
	// well-formed rather than panicking — the caller records it like any other.
	t.Run("handles a packet that did not decode", func(t *testing.T) {
		key := cosmosEventKeyForPacket(services.CosmosPacket{Type: services.CosmosSend, BlockNumber: 5})
		if key.Height != 5 || key.Sequence != 0 || key.SourceClient != "" {
			t.Fatalf("nil packet produced %+v, want only the height filled in", key)
		}
	})
}

// --- the live enqueue decision ---
//
// enqueueCosmosPackets decides which of the packets decoded out of a Cosmos event
// actually reach the relay queue. Everything it drops is dropped SILENTLY as far
// as the packet is concerned -- there is no retry behind it, so a packet wrongly
// filtered here is simply never relayed.

// cosmosTestDeps returns deps with the two client ids set and a quiet logger, so
// a test asserts on the queue rather than on log output.
func cosmosTestDeps(ethOnCosmos, cosmosOnEVM string) cosmosDeps {
	return cosmosDeps{
		IDs:    services.ClientIDs{EVMOnCosmos: ethOnCosmos, CosmosOnEVM: cosmosOnEVM},
		Logger: log.New(io.Discard, "", 0),
	}
}

func cosmosTestPacket(seq uint64, src, dst string) services.CosmosPacket {
	return services.CosmosPacket{
		Type:        services.CosmosSend,
		BlockNumber: 100,
		Packet:      &channeltypesv2.Packet{Sequence: seq, SourceClient: src, DestinationClient: dst},
	}
}

func TestEnqueueCosmosPackets(t *testing.T) {
	const (
		ethOnCosmos = "08-wasm-0"
		cosmosOnEVM = "eth-router-0"
	)

	t.Run("enqueues a packet on a configured client", func(t *testing.T) {
		bb := services.NewBatchBuilder()
		seen := map[cosmosEventKey]struct{}{}

		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{cosmosTestPacket(1, "cosmos-client", cosmosOnEVM)}, seen, false)
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if stats.recovered != 1 || stats.skipped != 0 {
			t.Fatalf("stats = %+v, want recovered=1 skipped=0", stats)
		}
		if cosmos, _ := bb.QueueDepths(); cosmos != 1 {
			t.Fatalf("queued %d packets, want 1", cosmos)
		}
	})

	// One pass can carry BOTH a send and the acknowledge_packet that closes it --
	// ordinary during recovery, where the window is sized to cover the downtime.
	// Settling only the tracker leaves the send QUEUED, so it is relayed for
	// nothing and, worse, tracked again by that relay: after settling removed it,
	// and after the cursor moved past the terminal event that would clear it a
	// second time. The entry then stays pending for the life of the process.
	t.Run("a send settled later in the same pass is removed from the queue too", func(t *testing.T) {
		bb := services.NewBatchBuilder()
		send := cosmosTestPacket(7, "cosmos-client", cosmosOnEVM)
		settled := cosmosTestPacket(7, "cosmos-client", cosmosOnEVM)
		settled.Type = services.CosmosAcknowledged

		if _, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{send, settled}, map[cosmosEventKey]struct{}{}, false); err != nil {
			t.Fatalf("enqueue: %v", err)
		}

		if cosmos, _ := bb.QueueDepths(); cosmos != 0 {
			t.Fatalf("queued %d packet(s); the send settled in this very pass and must not be relayed", cosmos)
		}
		if got := bb.PendingTracker.Len(); got != 0 {
			t.Fatalf("pending tracker holds %d packet(s) after the settlement", got)
		}
	})

	// The rule that makes terminal events safe to read at all: they close a
	// packet's lifecycle, so they settle the tracker and must NEVER reach the
	// batch. A terminal event in the relay queue is an empty "relay" of a packet
	// with nothing left to do -- it would be proven, submitted, and rejected.
	t.Run("a terminal acknowledge_packet settles the tracker and never reaches the batch", func(t *testing.T) {
		bb := services.NewBatchBuilder()
		packet := cosmosTestPacket(1, "cosmos-client", cosmosOnEVM)
		packet.Type = services.CosmosAcknowledged
		bb.PendingTracker.Add(*packet.Packet, 40)
		if bb.PendingTracker.Len() != 1 {
			t.Fatal("the packet must be tracked before the terminal event arrives")
		}

		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{packet}, map[cosmosEventKey]struct{}{}, false)
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if cosmos, _ := bb.QueueDepths(); cosmos != 0 {
			t.Fatalf("queued %d terminal event(s); a settled packet has no relay work left", cosmos)
		}
		if got := bb.PendingTracker.Len(); got != 0 {
			t.Fatalf("pending tracker holds %d packet(s) after the settlement; the whole point is dropping it without a query", got)
		}
		if stats.recovered != 1 {
			t.Fatalf("stats = %+v: a settlement is progress, not a skip", stats)
		}
	})

	t.Run("drops a packet belonging to another client pair", func(t *testing.T) {
		// A Cosmos chain relays for more than one counterparty. Enqueuing another
		// pair's packet does not merely waste a proof: it is submitted to a router
		// that has no commitment for it, and fails permanently.
		bb := services.NewBatchBuilder()
		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{cosmosTestPacket(1, "other-src", "other-dst")},
			map[cosmosEventKey]struct{}{}, false)
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if stats.skipped != 1 || stats.recovered != 0 {
			t.Fatalf("stats = %+v, want skipped=1 recovered=0", stats)
		}
		if cosmos, _ := bb.QueueDepths(); cosmos != 0 {
			t.Fatalf("queued %d packets from an unrelated client pair", cosmos)
		}
	})

	t.Run("matches on either side of either configured client", func(t *testing.T) {
		// The four ways a packet can belong to this relayer's path. Missing one
		// silently drops a whole direction -- an ack travels with the clients
		// swapped relative to the send it answers.
		for name, packet := range map[string]services.CosmosPacket{
			"source is the eth client":         cosmosTestPacket(1, ethOnCosmos, "x"),
			"destination is the eth client":    cosmosTestPacket(2, "x", ethOnCosmos),
			"source is the router client":      cosmosTestPacket(3, cosmosOnEVM, "x"),
			"destination is the router client": cosmosTestPacket(4, "x", cosmosOnEVM),
		} {
			t.Run(name, func(t *testing.T) {
				bb := services.NewBatchBuilder()
				stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
					[]services.CosmosPacket{packet}, map[cosmosEventKey]struct{}{}, false)
				if err != nil || stats.recovered != 1 {
					t.Fatalf("stats = %+v, err = %v; want recovered=1", stats, err)
				}
			})
		}
	})

	t.Run("drops a packet that did not decode", func(t *testing.T) {
		bb := services.NewBatchBuilder()
		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{{Type: services.CosmosSend}}, map[cosmosEventKey]struct{}{}, false)
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if stats.skipped != 1 {
			t.Fatalf("stats = %+v, want skipped=1", stats)
		}
	})

	t.Run("enqueues a packet once across repeated deliveries", func(t *testing.T) {
		// The live subscription and the gap recovery overlap by design, so the
		// same packet arrives twice. Relaying it twice costs a proof and a tx for
		// a submission the destination rejects.
		bb := services.NewBatchBuilder()
		seen := map[cosmosEventKey]struct{}{}
		packet := cosmosTestPacket(1, "cosmos-client", cosmosOnEVM)

		for i := 0; i < 3; i++ {
			if _, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
				[]services.CosmosPacket{packet}, seen, false); err != nil {
				t.Fatalf("delivery %d: %v", i, err)
			}
		}
		if cosmos, _ := bb.QueueDepths(); cosmos != 1 {
			t.Fatalf("queued %d copies of one packet", cosmos)
		}
	})

	t.Run("an enqueued packet is recorded, which is what the dedupe above rests on", func(t *testing.T) {
		seen := map[cosmosEventKey]struct{}{}
		packet := cosmosTestPacket(1, "cosmos-client", cosmosOnEVM)
		bb := services.NewBatchBuilder()

		if _, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{packet}, seen, false); err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if _, ok := seen[cosmosEventKeyForPacket(packet)]; !ok {
			t.Fatal("an enqueued packet was not recorded in the seen set")
		}
	})

	// Which drops are remembered is not uniform, and the difference is the cost of
	// the decision that produced them. Both branches below were untested; the
	// subtest above used to carry the name of the first one while exercising
	// neither, which review caught.

	t.Run("a packet declined during recovery is marked seen", func(t *testing.T) {
		// This is the drop worth remembering: deciding it cost a round trip to the
		// counterparty (HasEthPacketReceipt / HasPendingEthPacketCommitment), and a
		// recovery pass re-reads the same range every time. A Cosmos-originated
		// timeout reaches the same branch without an RPC -- it is declined by
		// ShouldRelayCosmosTimeoutToEth, which is pure -- so the branch is driven
		// here with no server standing in for the chain.
		seen := map[cosmosEventKey]struct{}{}
		bb := services.NewBatchBuilder()
		packet := cosmosTestPacket(1, "cosmos-client", cosmosOnEVM)
		packet.Type = services.CosmosTimeout

		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{packet}, seen, true)
		if err != nil {
			t.Fatalf("recovery enqueue: %v", err)
		}
		if stats.recovered != 0 || stats.skipped != 1 {
			t.Fatalf("stats = %+v; want the packet declined, not enqueued", stats)
		}
		if cosmos, _ := bb.QueueDepths(); cosmos != 0 {
			t.Fatalf("a declined packet was queued anyway (depth %d)", cosmos)
		}
		if _, ok := seen[cosmosEventKeyForPacket(packet)]; !ok {
			t.Fatal("a packet declined during recovery was not marked seen, so every later " +
				"recovery pass pays for the same decision again")
		}
	})

	t.Run("an unrelated packet is NOT marked seen", func(t *testing.T) {
		// The opposite of the branch above, and deliberately so. Rejecting an
		// unrelated client is a field comparison: stateless, free, and identical
		// on every pass. Caching it would buy nothing and would grow the map with
		// packets this relayer will never relay -- unbounded in the number of
		// other clients on the chain, not in its own traffic.
		seen := map[cosmosEventKey]struct{}{}
		bb := services.NewBatchBuilder()
		packet := cosmosTestPacket(1, "somebody-else", "not-our-router")

		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps(ethOnCosmos, cosmosOnEVM), bb,
			[]services.CosmosPacket{packet}, seen, false)
		if err != nil {
			t.Fatalf("enqueue: %v", err)
		}
		if stats.recovered != 0 || stats.skipped != 1 {
			t.Fatalf("stats = %+v; want the packet skipped as unrelated", stats)
		}
		if _, ok := seen[cosmosEventKeyForPacket(packet)]; ok {
			t.Fatal("an unrelated packet was cached in the seen set; that map then grows with " +
				"other clients' traffic rather than this relayer's")
		}
	})

	t.Run("with no client ids configured every packet matches", func(t *testing.T) {
		// The permissive default: an unconfigured relayer relays whatever it sees
		// rather than silently nothing, which is the failure that looks like a
		// dead relayer with a clean log.
		bb := services.NewBatchBuilder()
		stats, err := enqueueCosmosPackets(context.Background(), cosmosTestDeps("", ""), bb,
			[]services.CosmosPacket{cosmosTestPacket(1, "anything", "at-all")},
			map[cosmosEventKey]struct{}{}, false)
		if err != nil || stats.recovered != 1 {
			t.Fatalf("stats = %+v, err = %v; want recovered=1", stats, err)
		}
	})
}

// --- the live event handler ---

// failingCosmosDeps returns deps whose Cosmos RPC answers every method with an
// error, so a gap recovery triggered from the live path fails rather than
// panicking on a nil client.
func failingCosmosDeps(t *testing.T, ethOnCosmos, cosmosOnEVM string) cosmosDeps {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID json.RawMessage `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(req.ID) +
			`,"error":{"code":-32603,"message":"node is catching up"}}`))
	}))
	t.Cleanup(server.Close)

	client, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatalf("create CometBFT client: %v", err)
	}
	deps := cosmosTestDeps(ethOnCosmos, cosmosOnEVM)
	deps.Cosmos = services.CosmosEndpoint{Client: client}
	return deps
}

// liveEventAtHeight builds the ResultEvent shape the CometBFT subscription
// delivers, carrying only the tx height -- which is all the cursor logic reads.
func liveEventAtHeight(height int64) coretypes.ResultEvent {
	return coretypes.ResultEvent{
		Data:   commettypes.EventDataTx{TxResult: abcitypes.TxResult{Height: height}},
		Events: map[string][]string{},
	}
}

// The recovery cursor is what bounds the window a restart re-scans. Advancing it
// past a gap the relayer did NOT successfully re-scan means those heights are
// never looked at again: the events in them are lost silently, which is the exact
// failure the gap recovery exists to prevent.
//
// This is a named failure mode for this repo -- never advance a cursor on a failed
// operation -- and the live path is where it is easiest to get wrong, because the
// success and failure branches sit either side of one `else`.
func TestProcessLiveCosmosEvent_CursorSurvivesAFailedGapRecovery(t *testing.T) {
	const (
		cursorBefore = uint64(100)
		liveHeight   = int64(150)
	)
	deps := failingCosmosDeps(t, "08-wasm-0", "eth-router-0")
	s := NewSubscriber()
	cursor := cursorBefore
	bb := services.NewBatchBuilder()

	s.processLiveCosmosEvent(context.Background(), deps, bb, &cursor,
		map[cosmosEventKey]struct{}{}, &cosmosLiveHealth{}, liveEventAtHeight(liveHeight))

	if cursor != cursorBefore {
		t.Fatalf("cursor advanced to %d after a failed gap recovery; heights %d-%d would never be re-scanned",
			cursor, cursorBefore, liveHeight-1)
	}
}

// recordedQueries collects the tx_search queries a stub node was asked, so a test
// can assert on the height range the relayer computed rather than only on what
// the stub chose to answer.
type recordedQueries struct {
	mu      sync.Mutex
	queries []string
}

func (r *recordedQueries) add(q string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.queries = append(r.queries, q)
}

func (r *recordedQueries) all() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]string(nil), r.queries...)
}

// emptyTxSearchCosmosDeps answers tx_search with a valid, empty result -- what a
// node returns for a range containing no IBC transactions, which is the ordinary
// case for most heights. It records every query it was asked.
func emptyTxSearchCosmosDeps(t *testing.T, ethOnCosmos, cosmosOnEVM string) (cosmosDeps, *recordedQueries) {
	t.Helper()
	seen := &recordedQueries{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage `json:"id"`
			Method string          `json:"method"`
			Params struct {
				Query string `json:"query"`
			} `json:"params"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		w.Header().Set("Content-Type", "application/json")
		if req.Method != "tx_search" {
			_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(req.ID) +
				`,"error":{"code":-32601,"message":"unexpected method ` + req.Method + `"}}`))
			return
		}
		seen.add(req.Params.Query)
		_, _ = w.Write([]byte(`{"jsonrpc":"2.0","id":` + string(req.ID) +
			`,"result":{"txs":[],"total_count":"0"}}`))
	}))
	t.Cleanup(server.Close)

	client, err := rpchttp.New(server.URL, "/websocket")
	if err != nil {
		t.Fatalf("create CometBFT client: %v", err)
	}
	deps := cosmosTestDeps(ethOnCosmos, cosmosOnEVM)
	deps.Cosmos = services.CosmosEndpoint{Client: client}
	return deps, seen
}

// The mirror of the above: when the gap scan succeeds the cursor MUST move, or
// every later event re-scans the same range and the relayer spends its time
// re-reading history instead of relaying.
//
// The gap here is real -- heights 100 to 149 are scanned against a node that
// answers -- so this exercises the success branch rather than the guard that
// skips an empty range.
func TestProcessLiveCosmosEvent_CursorAdvancesOnASuccessfulGapRecovery(t *testing.T) {
	const (
		cursorBefore = uint64(100)
		liveHeight   = int64(150)
	)
	deps, queries := emptyTxSearchCosmosDeps(t, "08-wasm-0", "eth-router-0")
	s := NewSubscriber()
	cursor := cursorBefore

	s.processLiveCosmosEvent(context.Background(), deps, services.NewBatchBuilder(), &cursor,
		map[cosmosEventKey]struct{}{}, &cosmosLiveHealth{}, liveEventAtHeight(liveHeight))

	if cursor != uint64(liveHeight) {
		t.Fatalf("cursor = %d after a successful scan of %d-%d, want %d; a cursor that "+
			"does not move makes every later event re-scan the same range",
			cursor, cursorBefore, liveHeight-1, liveHeight)
	}

	// The gap ENDS one below the live height. Including the live height would
	// re-scan the very event being processed here, enqueuing it twice -- once from
	// the live path and once from the recovery beside it.
	asked := queries.all()
	if len(asked) == 0 {
		t.Fatal("no tx_search was issued, so no gap was scanned")
	}
	wantRange := fmt.Sprintf("tx.height >= %d AND tx.height <= %d", cursorBefore, liveHeight-1)
	for _, q := range asked {
		if !strings.Contains(q, wantRange) {
			t.Fatalf("scanned %q, want the range %q", q, wantRange)
		}
	}
}

// A live event must count towards liveness even when the gap recovery beside it
// fails. The health signal answers "is the subscription still delivering", and
// conflating it with "is recovery healthy" makes a working subscription look dead
// and triggers a reconnect that loses in-flight events.
func TestProcessLiveCosmosEvent_RecordsLivenessDespiteRecoveryFailure(t *testing.T) {
	deps := failingCosmosDeps(t, "08-wasm-0", "eth-router-0")
	s := NewSubscriber()
	cursor := uint64(100)
	health := &cosmosLiveHealth{}

	s.processLiveCosmosEvent(context.Background(), deps, services.NewBatchBuilder(), &cursor,
		map[cosmosEventKey]struct{}{}, health, liveEventAtHeight(150))

	if health.events != 1 {
		t.Fatalf("live health recorded %d events, want 1", health.events)
	}
	if health.lastSeen.IsZero() {
		t.Fatal("live health did not record when the event arrived")
	}
}

// An event with no height still has to be processed: the packets in it are real.
// What must not happen is a prune keyed on height zero, which would compute a
// retention floor from nothing.
func TestProcessLiveCosmosEvent_HeightlessEventDoesNotPrune(t *testing.T) {
	deps := cosmosTestDeps("08-wasm-0", "eth-router-0")
	s := NewSubscriber()
	cursor := uint64(100)
	seen := map[cosmosEventKey]struct{}{
		{Sequence: 1, Height: 1}: {}, // far below any plausible retention floor
	}

	s.processLiveCosmosEvent(context.Background(), deps, services.NewBatchBuilder(), &cursor,
		seen, &cosmosLiveHealth{}, coretypes.ResultEvent{Events: map[string][]string{}})

	if len(seen) != 1 {
		t.Fatalf("an event with no height pruned the seen set (%d entries left)", len(seen))
	}
	if cursor != 100 {
		t.Fatalf("cursor moved to %d on an event with no height", cursor)
	}
}
