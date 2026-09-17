package evmflush_test

import (
	"context"
	"io"
	"log"
	"math/big"
	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/chain/evmflush"
	relayerclient "relayer/client"
	"relayer/rpcmock"
	"relayer/services"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var routerAddr = ethcommon.HexToAddress("0x000000000000000000000000000000000000c26e")

func routerABI(t *testing.T) *abi.ABI {
	t.Helper()
	parsed, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("router ABI: %v", err)
	}
	return parsed
}

func packet(seq uint64, src, dst string) contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence: seq, SourceClient: src, DestClient: dst, TimeoutTimestamp: 1_900_000_000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort: "transfer", DestPort: "transfer", Version: "ics20-1", Encoding: "application/json", Value: []byte("{}"),
		}},
	}
}

// sendLog is a SendPacket log as the router emits it: the event id, the
// keccak of the indexed client id, the sequence, and the ABI-encoded packet.
func sendLog(t *testing.T, a *abi.ABI, clientID string, pkt contractICS26Router.IICS26RouterMsgsPacket, block uint64) rpcmock.RawLog {
	t.Helper()
	ev := a.Events["SendPacket"]
	data, err := ev.Inputs.NonIndexed().Pack(pkt)
	if err != nil {
		t.Fatalf("pack SendPacket: %v", err)
	}
	return rpcmock.RawLog{
		Address:     routerAddr,
		Topics:      []ethcommon.Hash{ev.ID, crypto.Keccak256Hash([]byte(clientID)), ethcommon.BigToHash(new(big.Int).SetUint64(pkt.Sequence))},
		Data:        data,
		BlockNumber: block,
	}
}

func writeAckLog(t *testing.T, a *abi.ABI, clientID string, pkt contractICS26Router.IICS26RouterMsgsPacket, acks [][]byte, block uint64) rpcmock.RawLog {
	t.Helper()
	ev := a.Events["WriteAcknowledgement"]
	data, err := ev.Inputs.NonIndexed().Pack(pkt, acks)
	if err != nil {
		t.Fatalf("pack WriteAcknowledgement: %v", err)
	}
	return rpcmock.RawLog{
		Address:     routerAddr,
		Topics:      []ethcommon.Hash{ev.ID, crypto.Keccak256Hash([]byte(clientID)), ethcommon.BigToHash(new(big.Int).SetUint64(pkt.Sequence))},
		Data:        data,
		BlockNumber: block,
	}
}

func newRouter(t *testing.T, node *rpcmock.EVMNode) evmflush.Router {
	t.Helper()
	return evmflush.Router{
		Client: node.Client(t), Address: routerAddr, Span: &relayerclient.LogSpan{},
		Label: "[FlushTest]", Logger: log.New(io.Discard, "", 0),
	}
}

func commitmentOf(pkt contractICS26Router.IICS26RouterMsgsPacket) ethcommon.Hash {
	return evmflush.PacketCommitment(pkt, new(big.Int).SetUint64(pkt.Sequence))
}

func commitmentPath(seq uint64) []byte {
	return services.EthPath("client-0", seq, services.EthPathCommitment)
}

// The probe is bounded above by prevSequenceSend and below by the cap, walks
// downward across passes, and wraps to the newest -- never 1..N.
func TestWindow(t *testing.T) {
	for _, tc := range []struct {
		name         string
		prev, cursor uint64
		cap          int
		lo, hi       uint64
		ok           bool
	}{
		{name: "a router that never sent has nothing to probe", prev: 0, cursor: 0, cap: 200},
		{name: "fewer sends than the cap: the whole space, once", prev: 5, cursor: 0, cap: 200, lo: 1, hi: 5, ok: true},
		{name: "first pass takes the newest cap", prev: 1000, cursor: 0, cap: 200, lo: 801, hi: 1000, ok: true},
		{name: "next pass continues below the last", prev: 1000, cursor: 801, cap: 200, lo: 601, hi: 800, ok: true},
		{name: "the last slice is short", prev: 1000, cursor: 150, cap: 200, lo: 1, hi: 149, ok: true},
		{name: "reaching sequence 1 wraps to the newest", prev: 1000, cursor: 1, cap: 200, lo: 801, hi: 1000, ok: true},
		{name: "a cursor above prev (client re-created) restarts from the newest", prev: 300, cursor: 801, cap: 200, lo: 101, hi: 300, ok: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			lo, hi, ok := evmflush.Window(tc.prev, tc.cursor, tc.cap)
			if ok != tc.ok || lo != tc.lo || hi != tc.hi {
				t.Fatalf("Window(prev=%d, cursor=%d, cap=%d) = [%d,%d] ok=%v, want [%d,%d] ok=%v",
					tc.prev, tc.cursor, tc.cap, lo, hi, ok, tc.lo, tc.hi, tc.ok)
			}
			if ok && hi-lo+1 > uint64(tc.cap) {
				t.Fatalf("window [%d,%d] is wider than the cap %d", lo, hi, tc.cap)
			}
		})
	}
}

func TestSequences_AreNewestFirst(t *testing.T) {
	got := evmflush.Sequences(3, 5)
	if len(got) != 3 || got[0] != 5 || got[2] != 3 {
		t.Fatalf("Sequences(3,5) = %v, want [5 4 3]", got)
	}
	if got := evmflush.Sequences(5, 3); got != nil {
		t.Fatalf("an inverted range must be empty, got %v", got)
	}
}

func TestPrevSequenceSend_ReadsTheRouterSlot(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	node.SetStorage(routerAddr, services.EthPrevSequenceSendKey("client-0"), ethcommon.BigToHash(big.NewInt(1234)))

	prev, err := newRouter(t, node).PrevSequenceSend(context.Background(), "client-0")
	if err != nil || prev != 1234 {
		t.Fatalf("PrevSequenceSend = (%d, %v), want 1234", prev, err)
	}
	if prev, err := newRouter(t, node).PrevSequenceSend(context.Background(), "client-9"); err != nil || prev != 0 {
		t.Fatalf("a client that never sent must read 0, got (%d, %v)", prev, err)
	}
}

// One pass over 150 sequences is TWO round trips (batches of 100), not 150, and
// only the slots that hold a value come back -- with the value, which is what
// says which packet the sequence names.
func TestOutstandingSlots_IsBatchedAndReturnsOnlyHeldSlots(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	held7 := ethcommon.HexToHash("0x07")
	held120 := ethcommon.HexToHash("0x0120")
	node.SetStorage(routerAddr, services.EthIBCStorageKey(commitmentPath(7)), held7)
	node.SetStorage(routerAddr, services.EthIBCStorageKey(commitmentPath(120)), held120)

	slots, err := newRouter(t, node).OutstandingSlots(context.Background(), evmflush.Sequences(1, 150), commitmentPath)
	if err != nil {
		t.Fatalf("OutstandingSlots: %v", err)
	}
	if len(slots) != 2 || slots[0].Sequence != 120 || slots[0].Value != held120 || slots[1].Sequence != 7 || slots[1].Value != held7 {
		t.Fatalf("slots = %+v, want [120, 7] with their values, newest first", slots)
	}
	if got := node.MethodCalls("eth_getStorageAt"); got != 150 {
		t.Fatalf("eth_getStorageAt calls = %d, want one per probed sequence (150)", got)
	}
	if got := node.Batches(); got != 2 {
		t.Fatalf("JSON-RPC batches = %d, want 2: the probe must be one round trip per 100 slots, not one per slot", got)
	}
}

func TestSendLog(t *testing.T) {
	t.Run("finds the event by topic over the whole chain", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := routerABI(t)
		node.SetBlockNumber(5000)
		node.AddLog(sendLog(t, a, "client-0", packet(41, "client-0", "08-wasm-0"), 10))
		node.AddLog(sendLog(t, a, "client-0", packet(42, "client-0", "08-wasm-0"), 11))
		node.AddLog(sendLog(t, a, "client-1", packet(42, "client-1", "08-wasm-0"), 12)) // same sequence, other client

		ev, found, err := newRouter(t, node).SendLog(context.Background(), "client-0", 42, commitmentOf(packet(42, "client-0", "08-wasm-0")))
		if err != nil || !found {
			t.Fatalf("SendLog = (found=%v, %v), want the log at block 11", found, err)
		}
		if ev.Packet.SourceClient != "client-0" || ev.Sequence.Uint64() != 42 || ev.Raw.BlockNumber != 11 {
			t.Fatalf("wrong log: client=%s seq=%d block=%d", ev.Packet.SourceClient, ev.Sequence.Uint64(), ev.Raw.BlockNumber)
		}
		if got := node.MethodCalls("eth_getLogs"); got != 1 {
			t.Fatalf("eth_getLogs calls = %d, want 1: an uncapped provider is asked once for the whole chain", got)
		}
	})

	// A provider that caps the range is asked for less and the walk starts at the
	// head, where an outstanding packet is most likely to be.
	t.Run("narrows when the provider refuses and walks back from the head", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := routerABI(t)
		node.SetBlockNumber(1000)
		node.CapLogRange(100, "query returned more than 10000 results")
		node.AddLog(sendLog(t, a, "client-0", packet(3, "client-0", "08-wasm-0"), 850))

		router := newRouter(t, node)
		ev, found, err := router.SendLog(context.Background(), "client-0", 3, commitmentOf(packet(3, "client-0", "08-wasm-0")))
		if err != nil || !found || ev.Raw.BlockNumber != 850 {
			t.Fatalf("SendLog = (found=%v, %v), want the log at block 850", found, err)
		}
		if router.Span.Chunk == 0 || router.Span.Chunk > 100 {
			t.Fatalf("span chunk = %d after a refusal, want narrowed to at most the provider's 100", router.Span.Chunk)
		}
		if got := node.MethodCalls("eth_getLogs"); got > 8 {
			t.Fatalf("eth_getLogs calls = %d; a log 150 blocks below the head must not cost a walk of the whole chain", got)
		}
	})

	// The walk is bounded. A packet whose event is older than MaxLogChunks provider
	// chunks is reported and left, not chased to genesis one chunk at a time.
	t.Run("gives up after max log chunks", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := routerABI(t)
		node.SetBlockNumber(100_000)
		node.CapLogRange(10, "block range too large")
		node.AddLog(sendLog(t, a, "client-0", packet(1, "client-0", "08-wasm-0"), 5))

		_, found, err := newRouter(t, node).SendLog(context.Background(), "client-0", 1, commitmentOf(packet(1, "client-0", "08-wasm-0")))
		if err != nil {
			t.Fatalf("giving up is not an error: %v", err)
		}
		if found {
			t.Fatal("a log 100k blocks back was found through a 10-block cap; the walk is not bounded")
		}
		// MaxLogChunks walks plus the handful of refused requests it took to find the cap.
		if got := node.MethodCalls("eth_getLogs"); got > evmflush.MaxLogChunks+16 {
			t.Fatalf("eth_getLogs calls = %d, want at most MaxLogChunks (%d) plus the ladder", got, evmflush.MaxLogChunks)
		}
	})

	// Reported by @neitdung on #496. A client re-created on the same id restarts
	// its sequences, so the topic filter matches one SendPacket log per generation
	// and a full-range query returns the OLDEST first. Stopping at the first log
	// rebuilt the closed packet, its commitment mismatched, and the live send behind
	// the same sequence was never examined -- on every pass. The held commitment
	// picks the generation; a commitment no log matches is "not found", not an
	// error, and costs one query.
	t.Run("picks the generation whose commitment the router holds", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		a := routerABI(t)
		node.SetBlockNumber(5000)
		old := packet(1, "client-0", "08-wasm-0")
		old.TimeoutTimestamp = 1_700_000_000
		current := packet(1, "client-0", "08-wasm-0")
		node.AddLog(sendLog(t, a, "client-0", old, 10))
		node.AddLog(sendLog(t, a, "client-0", current, 4000))

		router := newRouter(t, node)
		ev, found, err := router.SendLog(context.Background(), "client-0", 1, commitmentOf(current))
		if err != nil || !found || ev.Raw.BlockNumber != 4000 {
			t.Fatalf("SendLog(current) = (found=%v, block=%v, %v), want the live generation at block 4000", found, ev, err)
		}
		ev, found, err = router.SendLog(context.Background(), "client-0", 1, commitmentOf(old))
		if err != nil || !found || ev.Raw.BlockNumber != 10 {
			t.Fatalf("SendLog(old) = (found=%v, %v), want the old generation at block 10", found, err)
		}
		calls := node.MethodCalls("eth_getLogs")
		_, found, err = router.SendLog(context.Background(), "client-0", 1, ethcommon.HexToHash("0xdead"))
		if err != nil || found {
			t.Fatalf("SendLog(unknown commitment) = (found=%v, %v), want not found and no error", found, err)
		}
		if got := node.MethodCalls("eth_getLogs") - calls; got != 1 {
			t.Fatalf("an unmatched commitment cost %d eth_getLogs calls, want 1: the whole chain was already asked", got)
		}
	})
}

func TestWriteAckLog_CarriesTheAcknowledgements(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	a := routerABI(t)
	node.SetBlockNumber(50)
	acks := [][]byte{[]byte(`{"result":"AQ=="}`)}
	node.AddLog(writeAckLog(t, a, "client-0", packet(9, "08-wasm-0", "client-0"), acks, 40))

	ev, found, err := newRouter(t, node).WriteAckLog(context.Background(), "client-0", 9, commitmentOf(packet(9, "08-wasm-0", "client-0")).Bytes())
	if err != nil || !found {
		t.Fatalf("WriteAckLog = (found=%v, %v)", found, err)
	}
	if len(ev.Acknowledgements) != 1 || string(ev.Acknowledgements[0]) != string(acks[0]) {
		t.Fatalf("acknowledgements = %q, want the log's", ev.Acknowledgements)
	}
	if ev.Packet.DestClient != "client-0" || ev.Sequence.Uint64() != 9 {
		t.Fatalf("wrong log: dest=%s seq=%d", ev.Packet.DestClient, ev.Sequence.Uint64())
	}
}

// The same rule on the acknowledgement side, where the commitment is the one the
// SENDING chain holds -- and the walk-back path, where the live generation sits
// in a newer chunk than the stale one.
func TestWriteAckLog_PicksTheGenerationTheSenderStillHolds(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	a := routerABI(t)
	node.SetBlockNumber(1000)
	node.CapLogRange(100, "block range too large")
	acks := [][]byte{[]byte(`{"result":"AQ=="}`)}
	old := packet(2, "08-wasm-0", "client-0")
	old.TimeoutTimestamp = 1_700_000_000
	current := packet(2, "08-wasm-0", "client-0")
	node.AddLog(writeAckLog(t, a, "client-0", old, acks, 300))
	node.AddLog(writeAckLog(t, a, "client-0", current, acks, 950))

	ev, found, err := newRouter(t, node).WriteAckLog(context.Background(), "client-0", 2, commitmentOf(current).Bytes())
	if err != nil || !found || ev.Raw.BlockNumber != 950 {
		t.Fatalf("WriteAckLog(current) = (found=%v, %v), want the live generation at block 950", found, err)
	}
	ev, found, err = newRouter(t, node).WriteAckLog(context.Background(), "client-0", 2, commitmentOf(old).Bytes())
	if err != nil || !found || ev.Raw.BlockNumber != 300 {
		t.Fatalf("WriteAckLog(old) = (found=%v, %v), want the old generation at block 300", found, err)
	}
}
