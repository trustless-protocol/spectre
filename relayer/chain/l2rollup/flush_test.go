package l2rollup

import (
	"context"
	"errors"
	"math/big"
	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/chain"
	"relayer/rpcmock"
	"relayer/services"
	"relayer/subscriber"
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	l2FlushClient   = "client-0"  // the L2 router's client for Cosmos
	l2FlushWasmName = "08-wasm-3" // Cosmos's wasm client for this L2
)

var l2FlushRouter = ethcommon.HexToAddress("0x00000000000000000000000000000000000000f2")

func l2FlushABI(t *testing.T) *abi.ABI {
	t.Helper()
	parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("router ABI: %v", err)
	}
	return parsed
}

func flushL2Packet(seq uint64, src, dst string) contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence: seq, SourceClient: src, DestClient: dst, TimeoutTimestamp: 4_000_000_000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort: "transfer", DestPort: "transfer", Version: "ics20-1", Encoding: "application/json", Value: []byte("{}"),
		}},
	}
}

func l2Commitment(pkt contractICS26Router.IICS26RouterMsgsPacket) ethcommon.Hash {
	return ethcommon.BytesToHash(channeltypesv2.CommitPacket(subscriber.EthPacketToCosmosPacket(pkt, new(big.Int).SetUint64(pkt.Sequence))))
}

func l2Log(t *testing.T, a *abi.ABI, event, clientID string, block uint64, args ...any) rpcmock.RawLog {
	t.Helper()
	data, err := a.Events[event].Inputs.NonIndexed().Pack(args...)
	if err != nil {
		t.Fatalf("pack %s: %v", event, err)
	}
	pkt := args[0].(contractICS26Router.IICS26RouterMsgsPacket)
	return rpcmock.RawLog{
		Address: l2FlushRouter,
		Topics:  []ethcommon.Hash{a.Events[event].ID, crypto.Keccak256Hash([]byte(clientID)), ethcommon.BigToHash(new(big.Int).SetUint64(pkt.Sequence))},
		Data:    data, BlockNumber: block,
	}
}

// listingState is a CosmosPacketState that can also page a commitment window,
// which is what the acknowledgement leg needs from the Cosmos side.
type listingState struct {
	fakeCosmosState
	window []services.PacketCommitment
	asked  []string
}

func (l *listingState) CommitmentWindow(_ context.Context, clientID string, _ uint64, _ int) ([]services.PacketCommitment, error) {
	l.asked = append(l.asked, clientID)
	return l.window, nil
}

func TestUnrelayedPackets(t *testing.T) {
	t.Run("L2 sends filtered to the paired cosmos client", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := l2FlushABI(t)
		node.SetBlockNumber(300)
		node.SetStorage(l2FlushRouter, services.EthPrevSequenceSendKey(l2FlushClient), ethcommon.BigToHash(big.NewInt(2)))
		toUs := flushL2Packet(1, l2FlushClient, l2FlushWasmName)
		toOther := flushL2Packet(2, l2FlushClient, "08-wasm-9")
		for _, pkt := range []contractICS26Router.IICS26RouterMsgsPacket{toUs, toOther} {
			node.SetStorage(l2FlushRouter, services.EthIBCStorageKey(services.EthPath(l2FlushClient, pkt.Sequence, services.EthPathCommitment)), l2Commitment(pkt))
			node.AddLog(l2Log(t, a, "SendPacket", l2FlushClient, 10+pkt.Sequence, pkt))
		}

		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName, logScanChunk: 0}
		events, err := src.UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Type != chain.SendPacket || events[0].Sequence != 1 || events[0].Height != 11 {
			t.Fatalf("events = %+v, want only the send addressed to %s (seq 1)", events, l2FlushWasmName)
		}
	})

	t.Run("L2 acks for open cosmos sends", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := l2FlushABI(t)
		node.SetBlockNumber(300)

		fromCosmos := flushL2Packet(9, l2FlushWasmName, l2FlushClient)
		state := &listingState{window: []services.PacketCommitment{{Sequence: 9, Data: l2Commitment(fromCosmos).Bytes()}}}
		acks := [][]byte{[]byte("ok")}
		node.SetStorage(l2FlushRouter, services.EthIBCStorageKey(services.EthPath(l2FlushClient, 9, services.EthPathAck)), ethcommon.HexToHash("0xac"))
		node.AddLog(l2Log(t, a, "WriteAcknowledgement", l2FlushClient, 42, fromCosmos, acks))

		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName, cosmosState: state}
		events, err := src.UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Type != chain.AckPacket || events[0].Sequence != 9 || events[0].Height != 42 {
			t.Fatalf("events = %+v, want the acknowledgement for seq 9 at block 42", events)
		}
		if len(state.asked) != 1 || state.asked[0] != l2FlushWasmName {
			t.Fatalf("Cosmos commitments were paged for %v, want the paired wasm client %q: that is the client Cosmos->L2 sends are committed under", state.asked, l2FlushWasmName)
		}
		if src.flush.ackCursor != 9 {
			t.Fatalf("ack cursor = %d, want 9", src.flush.ackCursor)
		}
	})

	// A Cosmos seam that cannot page commitments (the reconciliation-only doubles,
	// or no seam at all) simply has no acknowledgements to list.
	t.Run("L2 without a commitment lister lists sends only", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.SetBlockNumber(10)
		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName, cosmosState: &fakeCosmosState{}}
		events, err := src.UnrelayedPackets(context.Background())
		if err != nil || len(events) != 0 {
			t.Fatalf("UnrelayedPackets = (%+v, %v), want nothing and no error", events, err)
		}
	})

	// The L2 half of @neitdung's #496 repro: a reused sequence leaves two logs per
	// leg, and the held commitment must pick the current one on both.
	t.Run("L2 reused sequence resolves to the current generation", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := l2FlushABI(t)
		node.SetBlockNumber(900)
		node.SetStorage(l2FlushRouter, services.EthPrevSequenceSendKey(l2FlushClient), ethcommon.BigToHash(big.NewInt(1)))
		staleSend := flushL2Packet(1, l2FlushClient, l2FlushWasmName)
		staleSend.TimeoutTimestamp = 3_000_000_000
		currentSend := flushL2Packet(1, l2FlushClient, l2FlushWasmName)
		node.SetStorage(l2FlushRouter, services.EthIBCStorageKey(services.EthPath(l2FlushClient, 1, services.EthPathCommitment)), l2Commitment(currentSend))
		node.AddLog(l2Log(t, a, "SendPacket", l2FlushClient, 10, staleSend))
		node.AddLog(l2Log(t, a, "SendPacket", l2FlushClient, 800, currentSend))

		staleAck := flushL2Packet(3, l2FlushWasmName, l2FlushClient)
		staleAck.TimeoutTimestamp = 3_000_000_000
		currentAck := flushL2Packet(3, l2FlushWasmName, l2FlushClient)
		state := &listingState{window: []services.PacketCommitment{{Sequence: 3, Data: l2Commitment(currentAck).Bytes()}}}
		acks := [][]byte{[]byte("ok")}
		node.SetStorage(l2FlushRouter, services.EthIBCStorageKey(services.EthPath(l2FlushClient, 3, services.EthPathAck)), ethcommon.HexToHash("0xac"))
		node.AddLog(l2Log(t, a, "WriteAcknowledgement", l2FlushClient, 20, staleAck, acks))
		node.AddLog(l2Log(t, a, "WriteAcknowledgement", l2FlushClient, 700, currentAck, acks))

		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName, cosmosState: state}
		events, err := src.UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 2 || events[0].Type != chain.SendPacket || events[0].Height != 800 || events[1].Type != chain.AckPacket || events[1].Height != 700 {
			t.Fatalf("events = %+v, want the current send (block 800) and the current ack (block 700)", events)
		}
	})

	// An L2 send past its timeout is the timeout scanner's, not a proof to build
	// and fail every pass -- the same pre-filter the ETH source applies.
	t.Run("L2 skips an expired send", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := l2FlushABI(t)
		node.SetBlockNumber(300)
		node.SetStorage(l2FlushRouter, services.EthPrevSequenceSendKey(l2FlushClient), ethcommon.BigToHash(big.NewInt(1)))
		expired := flushL2Packet(1, l2FlushClient, l2FlushWasmName)
		expired.TimeoutTimestamp = 1
		node.SetStorage(l2FlushRouter, services.EthIBCStorageKey(services.EthPath(l2FlushClient, 1, services.EthPathCommitment)), l2Commitment(expired))
		node.AddLog(l2Log(t, a, "SendPacket", l2FlushClient, 11, expired))

		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName}
		events, err := src.UnrelayedPackets(context.Background())
		if err != nil || len(events) != 0 {
			t.Fatalf("UnrelayedPackets = (%+v, %v), want nothing: the send is expired", events, err)
		}
	})

	// A probe that fails must not move the cursors (the ETH source's rule, mirrored).
	t.Run("a failed probe leaves both cursors where they were", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.SetBlockNumber(300)
		node.SetStorage(l2FlushRouter, services.EthPrevSequenceSendKey(l2FlushClient), ethcommon.BigToHash(big.NewInt(2)))
		fromCosmos := flushL2Packet(9, l2FlushWasmName, l2FlushClient)
		state := &listingState{window: []services.PacketCommitment{{Sequence: 9, Data: l2Commitment(fromCosmos).Bytes()}}}
		src := &Source{eth: node.Client(t), router: l2FlushRouter, l2ClientID: l2FlushClient, cosmosWasmClientID: l2FlushWasmName, cosmosState: state}

		node.FailBatches(errors.New("probe down"))
		if _, err := src.UnrelayedPackets(context.Background()); err == nil {
			t.Fatal("the send-leg probe failed and the pass reported success")
		}
		if src.flush.sendCursor != 0 || src.flush.ackCursor != 0 {
			t.Fatalf("cursors advanced on a failed probe: send=%d ack=%d", src.flush.sendCursor, src.flush.ackCursor)
		}

		node.SetStorage(l2FlushRouter, services.EthPrevSequenceSendKey(l2FlushClient), ethcommon.Hash{})
		if _, err := src.UnrelayedPackets(context.Background()); err != nil {
			t.Fatalf("a failing ack leg is reported, not returned: %v", err)
		}
		if src.flush.ackCursor != 0 {
			t.Fatalf("ack cursor advanced to %d on a failed probe", src.flush.ackCursor)
		}

		node.FailBatches(nil)
		if _, err := src.UnrelayedPackets(context.Background()); err != nil {
			t.Fatal(err)
		}
		if src.flush.ackCursor != 9 {
			t.Fatalf("ack cursor = %d after a successful probe, want 9", src.flush.ackCursor)
		}
	})
}
