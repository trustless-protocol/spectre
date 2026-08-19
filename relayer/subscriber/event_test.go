package subscriber

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"math/big"
	"os"
	"testing"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/services"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
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
			EVENT_TX_HEIGHT_FIELD:        {"44"},
		},
		"test",
	)

	if len(packets) != 3 {
		t.Fatalf("decoded packet count = %d, want 3", len(packets))
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

	enqueueEthSendPacket(bb, ev, nil)

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

	if !enqueueEthSendPacket(bb, ev, seen) {
		t.Fatal("first enqueue should be accepted")
	}
	if enqueueEthSendPacket(bb, ev, seen) {
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
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	ev := &contractICS26Router.ContractICS26RouterSendPacket{
		Sequence: big.NewInt(7),
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			SourceClient: "eth-client-0",
			DestClient:   "cosmos-client-0",
		},
		Raw: gethtypes.Log{BlockNumber: 88, Index: 3},
	}
	seen := make(map[ethEventKey]struct{})

	if enqueueEthSendPacket(bb, ev, seen) {
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

// Recovery quietly doing the live path's job is the exact state that hid a lost
// packet for 14 minutes. It must be reported — once per outage, not per scan.
func TestLiveHealthReportsRecoveryDoingLiveWork(t *testing.T) {
	t.Parallel()

	var liveHealth cosmosLiveHealth
	start := time.Now()

	// First recovery hit with no live history: take the current time as the
	// reference rather than crying wolf at startup.
	if liveHealth.recoveryIsCoveringForLive(1, start) {
		t.Fatal("warned before any live baseline existed")
	}
	// Still inside the window: not yet evidence of a dead subscription.
	if liveHealth.recoveryIsCoveringForLive(1, start.Add(cosmosLivePathStaleAfter-time.Second)) {
		t.Fatal("warned while the live path was only briefly quiet")
	}
	// Past the window with recovery still finding packets: report.
	if !liveHealth.recoveryIsCoveringForLive(1, start.Add(cosmosLivePathStaleAfter+time.Second)) {
		t.Fatal("did not report recovery doing the live path's job")
	}
	// ...but only once, or a long outage floods the log.
	if liveHealth.recoveryIsCoveringForLive(1, start.Add(cosmosLivePathStaleAfter+time.Minute)) {
		t.Fatal("repeated the warning for the same outage")
	}

	// A live event clears it, and a later outage reports again.
	liveHealth.recordEvent()
	if liveHealth.events != 1 {
		t.Fatalf("events = %d, want 1", liveHealth.events)
	}
	if !liveHealth.recoveryIsCoveringForLive(1, liveHealth.lastSeen.Add(cosmosLivePathStaleAfter+time.Second)) {
		t.Fatal("a fresh outage after recovery must report again")
	}
}

// A scan that finds nothing says nothing about the live path — recovery finding
// no work is the healthy state, not a symptom.
func TestLiveHealthIgnoresEmptyRecoveryScans(t *testing.T) {
	t.Parallel()

	var liveHealth cosmosLiveHealth
	liveHealth.recordEvent()
	if liveHealth.recoveryIsCoveringForLive(0, liveHealth.lastSeen.Add(24*time.Hour)) {
		t.Fatal("an empty scan must not be read as the live path failing")
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
			enqueueEthTerminal(bb, tc.typ, pkt, big.NewInt(4), tc.ackBytes, 4242)

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

	enqueueEthTerminal(bb, services.EthAck, pkt, big.NewInt(4), [][]byte{[]byte("ack")}, 4242)

	if got := bb.EthPendingTracker.Len(); got != 0 {
		t.Fatalf("tracker length = %d after settlement, want 0", got)
	}
}
