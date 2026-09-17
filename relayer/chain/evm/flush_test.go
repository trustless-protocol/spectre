package evm

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"math/big"
	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/chain"
	"relayer/rpcmock"
	"relayer/services"
	"relayer/subscriber"
	"testing"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtbytes "github.com/cometbft/cometbft/libs/bytes"
	cmtjson "github.com/cometbft/cometbft/libs/json"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/cosmos/gogoproto/proto"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	flushCosmosOnEVM = "cosmoshub-1" // the router's client for Cosmos
	flushEVMOnCosmos = "08-wasm-0"   // Cosmos's client for ETH
)

var flushRouter = ethcommon.HexToAddress("0x000000000000000000000000000000000000c26e")

func flushABI(t *testing.T) *abi.ABI {
	t.Helper()
	parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("router ABI: %v", err)
	}
	return parsed
}

func flushEthPacket(seq uint64, src, dst string, timeout uint64) contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence: seq, SourceClient: src, DestClient: dst, TimeoutTimestamp: timeout,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort: "transfer", DestPort: "transfer", Version: "ics20-1", Encoding: "application/json", Value: []byte("{}"),
		}},
	}
}

func topics(a *abi.ABI, event, clientID string, seq uint64) []ethcommon.Hash {
	return []ethcommon.Hash{a.Events[event].ID, crypto.Keccak256Hash([]byte(clientID)), ethcommon.BigToHash(new(big.Int).SetUint64(seq))}
}

func commitmentOf(pkt contractICS26Router.IICS26RouterMsgsPacket) ethcommon.Hash {
	return ethcommon.BytesToHash(channeltypesv2.CommitPacket(subscriber.EthPacketToCosmosPacket(pkt, new(big.Int).SetUint64(pkt.Sequence))))
}

// cosmosCommitments serves the PacketCommitments query the ack leg pages, for
// one client, from a fixed list.
func cosmosCommitments(t *testing.T, node *rpcmock.CosmosNode, states []channeltypesv2.PacketState) {
	t.Helper()
	node.Handle("abci_query", func(params json.RawMessage) (any, error) {
		var q struct {
			Path string            `json:"path"`
			Data cmtbytes.HexBytes `json:"data"`
		}
		if err := cmtjson.Unmarshal(params, &q); err != nil {
			return nil, err
		}
		var request channeltypesv2.QueryPacketCommitmentsRequest
		if err := proto.Unmarshal(q.Data, &request); err != nil {
			return nil, err
		}
		response := channeltypesv2.QueryPacketCommitmentsResponse{}
		for i := range states {
			if states[i].ClientId == request.ClientId {
				response.Commitments = append(response.Commitments, &states[i])
			}
		}
		raw, err := proto.Marshal(&response)
		if err != nil {
			return nil, err
		}
		return &coretypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: raw}}, nil
	})
}

func newFlushSource(t *testing.T, evmNode *rpcmock.EVMNode, cosmosNode *rpcmock.CosmosNode) *Source {
	t.Helper()
	s := &Source{
		evm:    services.EVMEndpoint{Client: evmNode.Client(t), Contracts: services.EVMContracts{Router: flushRouter}},
		ids:    services.ClientIDs{CosmosOnEVM: flushCosmosOnEVM, EVMOnCosmos: flushEVMOnCosmos},
		logger: log.New(io.Discard, "", 0),
	}
	if cosmosNode != nil {
		s.cosmos = services.CosmosEndpoint{Client: cosmosNode.Client(t)}
	}
	return s
}

func TestUnrelayedPackets(t *testing.T) {
	// The send leg: the router says which sequences are still committed, the log
	// says what packet each one is, and only a packet that matches its commitment,
	// and can still be received, is handed to the relay loop.
	t.Run("lists the sends the router still holds", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := flushABI(t)
		node.SetBlockNumber(500)
		node.SetStorage(flushRouter, services.EthPrevSequenceSendKey(flushCosmosOnEVM), ethcommon.BigToHash(big.NewInt(4)))

		// seq 1: acknowledged long ago -- no commitment, no work.
		// seq 2: outstanding, matches its log, live -> relayed.
		// seq 3: outstanding but past its timeout -> the scanner's job, not this one.
		// seq 4: the slot holds a hash that is NOT the logged packet's (re-created client) -> skipped.
		live := flushEthPacket(2, flushCosmosOnEVM, flushEVMOnCosmos, 4_000_000_000)
		expired := flushEthPacket(3, flushCosmosOnEVM, flushEVMOnCosmos, 1)
		replayed := flushEthPacket(4, flushCosmosOnEVM, flushEVMOnCosmos, 4_000_000_000)
		for _, pkt := range []contractICS26Router.IICS26RouterMsgsPacket{live, expired} {
			node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, pkt.Sequence, services.EthPathCommitment)), commitmentOf(pkt))
		}
		node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, 4, services.EthPathCommitment)), ethcommon.HexToHash("0xdead"))
		for _, pkt := range []contractICS26Router.IICS26RouterMsgsPacket{live, expired, replayed} {
			data, err := a.Events["SendPacket"].Inputs.NonIndexed().Pack(pkt)
			if err != nil {
				t.Fatal(err)
			}
			node.AddLog(rpcmock.RawLog{Address: flushRouter, Topics: topics(a, "SendPacket", flushCosmosOnEVM, pkt.Sequence), Data: data, BlockNumber: 100 + pkt.Sequence})
		}

		s := newFlushSource(t, node, nil)
		events, err := s.UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Type != chain.SendPacket || events[0].Sequence != 2 || events[0].Height != 102 {
			t.Fatalf("events = %+v, want exactly the live, matching send (seq 2 at block 102)", events)
		}
		if events[0].ClientID != flushEVMOnCosmos {
			t.Fatalf("event client = %q, want the packet's destination client %q", events[0].ClientID, flushEVMOnCosmos)
		}
		if s.flush.sendCursor != 1 {
			t.Fatalf("send cursor = %d, want 1: the pass covered the whole space and the next one wraps", s.flush.sendCursor)
		}
		if got := node.Batches(); got != 1 {
			t.Fatalf("storage probe took %d batches, want 1 round trip for 4 sequences", got)
		}
	})

	// The acknowledgement leg: Cosmos lists the sends it still holds a commitment
	// for, the router is probed for an ack on exactly those, and the ack's packet
	// must match the commitment Cosmos holds.
	t.Run("lists the acks the router wrote for open cosmos sends", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cosmosNode := rpcmock.NewCosmosNode(t)
		a := flushABI(t)
		node.SetBlockNumber(500)

		// Cosmos->ETH packets: seq 7 acknowledged on the router, seq 8 not yet
		// received there, seq 9 acknowledged but the log's packet does not match
		// what Cosmos committed.
		acked := flushEthPacket(7, flushEVMOnCosmos, flushCosmosOnEVM, 4_000_000_000)
		pending := flushEthPacket(8, flushEVMOnCosmos, flushCosmosOnEVM, 4_000_000_000)
		mismatched := flushEthPacket(9, flushEVMOnCosmos, flushCosmosOnEVM, 4_000_000_000)
		cosmosCommitments(t, cosmosNode, []channeltypesv2.PacketState{
			{ClientId: flushEVMOnCosmos, Sequence: 7, Data: commitmentOf(acked).Bytes()},
			{ClientId: flushEVMOnCosmos, Sequence: 8, Data: commitmentOf(pending).Bytes()},
			{ClientId: flushEVMOnCosmos, Sequence: 9, Data: []byte("not what the log says")},
		})
		acks := [][]byte{[]byte(`{"result":"AQ=="}`)}
		for _, pkt := range []contractICS26Router.IICS26RouterMsgsPacket{acked, mismatched} {
			node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, pkt.Sequence, services.EthPathAck)), ethcommon.HexToHash("0xac"))
			data, err := a.Events["WriteAcknowledgement"].Inputs.NonIndexed().Pack(pkt, acks)
			if err != nil {
				t.Fatal(err)
			}
			node.AddLog(rpcmock.RawLog{Address: flushRouter, Topics: topics(a, "WriteAcknowledgement", flushCosmosOnEVM, pkt.Sequence), Data: data, BlockNumber: 200 + pkt.Sequence})
		}

		s := newFlushSource(t, node, cosmosNode)
		events, err := s.UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Type != chain.AckPacket || events[0].Sequence != 7 || events[0].Height != 207 {
			t.Fatalf("events = %+v, want exactly the acknowledged, matching send (seq 7 at block 207)", events)
		}
		if len(events[0].AckBytes) != 1 || string(events[0].AckBytes[0]) != string(acks[0]) {
			t.Fatalf("ack bytes = %q, want the log's", events[0].AckBytes)
		}
		if s.flush.ackCursor != 9 {
			t.Fatalf("ack cursor = %d, want the last Cosmos commitment taken (9)", s.flush.ackCursor)
		}
		// Only the three listed sequences were probed for an ack, not a range.
		if got := node.MethodCalls("eth_getStorageAt"); got != 1+3 {
			t.Fatalf("eth_getStorageAt calls = %d, want prevSequenceSend + one per Cosmos commitment (4)", got)
		}
	})

	// A Cosmos node that cannot be reached must not hide the router's sends.
	t.Run("a failing ack leg keeps the sends", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cosmosNode := rpcmock.NewCosmosNode(t) // no abci_query handler: the query errors
		a := flushABI(t)
		node.SetBlockNumber(50)
		node.SetStorage(flushRouter, services.EthPrevSequenceSendKey(flushCosmosOnEVM), ethcommon.BigToHash(big.NewInt(1)))
		live := flushEthPacket(1, flushCosmosOnEVM, flushEVMOnCosmos, 4_000_000_000)
		node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, 1, services.EthPathCommitment)), commitmentOf(live))
		data, _ := a.Events["SendPacket"].Inputs.NonIndexed().Pack(live)
		node.AddLog(rpcmock.RawLog{Address: flushRouter, Topics: topics(a, "SendPacket", flushCosmosOnEVM, 1), Data: data, BlockNumber: 10})

		events, err := newFlushSource(t, node, cosmosNode).UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("a failing ack leg must be reported, not fail the pass: %v", err)
		}
		if len(events) != 1 || events[0].Sequence != 1 {
			t.Fatalf("events = %+v, want the send the ack leg's failure must not hide", events)
		}
	})

	// @neitdung's repro on #496: after a client re-creation reuses a sequence, two
	// SendPacket logs share every topic and the router's slot holds the NEWER
	// packet's commitment. The enumeration must return that packet, not skip the
	// sequence because the older log came back first.
	t.Run("a send behind a reused sequence is found", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := flushABI(t)
		node.SetBlockNumber(900)
		node.SetStorage(flushRouter, services.EthPrevSequenceSendKey(flushCosmosOnEVM), ethcommon.BigToHash(big.NewInt(1)))
		stale := flushEthPacket(1, flushCosmosOnEVM, flushEVMOnCosmos, 3_000_000_000)
		current := flushEthPacket(1, flushCosmosOnEVM, flushEVMOnCosmos, 4_000_000_000)
		node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, 1, services.EthPathCommitment)), commitmentOf(current))
		// Stale first: that is the order a full-range eth_getLogs returns them in.
		for _, l := range []struct {
			block uint64
			pkt   contractICS26Router.IICS26RouterMsgsPacket
		}{{10, stale}, {800, current}} {
			data, err := a.Events["SendPacket"].Inputs.NonIndexed().Pack(l.pkt)
			if err != nil {
				t.Fatal(err)
			}
			node.AddLog(rpcmock.RawLog{Address: flushRouter, Topics: topics(a, "SendPacket", flushCosmosOnEVM, 1), Data: data, BlockNumber: l.block})
		}

		events, err := newFlushSource(t, node, nil).UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Sequence != 1 || events[0].Height != 800 {
			t.Fatalf("events = %+v, want the current generation of seq 1 (block 800)", events)
		}
	})

	// The acknowledgement leg under the same reuse: two WriteAcknowledgement logs
	// for one sequence, and Cosmos's commitment names the newer packet.
	t.Run("an ack behind a reused sequence is found", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		cosmosNode := rpcmock.NewCosmosNode(t)
		a := flushABI(t)
		node.SetBlockNumber(900)
		stale := flushEthPacket(5, flushEVMOnCosmos, flushCosmosOnEVM, 3_000_000_000)
		current := flushEthPacket(5, flushEVMOnCosmos, flushCosmosOnEVM, 4_000_000_000)
		cosmosCommitments(t, cosmosNode, []channeltypesv2.PacketState{{ClientId: flushEVMOnCosmos, Sequence: 5, Data: commitmentOf(current).Bytes()}})
		node.SetStorage(flushRouter, services.EthIBCStorageKey(services.EthPath(flushCosmosOnEVM, 5, services.EthPathAck)), ethcommon.HexToHash("0xac"))
		acks := [][]byte{[]byte(`{"result":"AQ=="}`)}
		for _, l := range []struct {
			block uint64
			pkt   contractICS26Router.IICS26RouterMsgsPacket
		}{{20, stale}, {700, current}} {
			data, err := a.Events["WriteAcknowledgement"].Inputs.NonIndexed().Pack(l.pkt, acks)
			if err != nil {
				t.Fatal(err)
			}
			node.AddLog(rpcmock.RawLog{Address: flushRouter, Topics: topics(a, "WriteAcknowledgement", flushCosmosOnEVM, 5), Data: data, BlockNumber: l.block})
		}

		events, err := newFlushSource(t, node, cosmosNode).UnrelayedPackets(context.Background())
		if err != nil {
			t.Fatalf("UnrelayedPackets: %v", err)
		}
		if len(events) != 1 || events[0].Type != chain.AckPacket || events[0].Sequence != 5 || events[0].Height != 700 {
			t.Fatalf("events = %+v, want the acknowledgement of the current generation of seq 5 (block 700)", events)
		}
	})
}
