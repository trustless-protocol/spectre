package services

import (
	"fmt"
	"slices"
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// These tests guard the #106 regression: a per-packet proof-build failure must be
// routed to the re-queue set (and, for CosmosSend, dropped from the pending tracker)
// instead of being silently dropped. They exercise the exact control flow that owns
// the five proof-build `continue` sites, with fake proof builders — no live RPC.

func cosmosPkt(seq uint64, typ CosmosPacketType, src, dst string, timeout uint64, ack [][]byte) CosmosPacket {
	return CosmosPacket{
		Type: typ,
		Packet: &channeltypesv2.Packet{
			Sequence:          seq,
			SourceClient:      src,
			DestinationClient: dst,
			TimeoutTimestamp:  timeout,
		},
		AckBytes: ack,
	}
}

func ethPkt(seq uint64, typ EthPacketType, src, dst string, timeout uint64, ack [][]byte) EthPacket {
	return EthPacket{
		Type: typ,
		Packet: &channeltypesv2.Packet{
			Sequence:          seq,
			SourceClient:      src,
			DestinationClient: dst,
			TimeoutTimestamp:  timeout,
		},
		AckBytes: ack,
	}
}

func cosmosSeqs(ps []CosmosPacket) []uint64 {
	out := make([]uint64, len(ps))
	for i, p := range ps {
		out[i] = p.Packet.Sequence
	}
	return out
}

func ethSeqs(ps []EthPacket) []uint64 {
	out := make([]uint64, len(ps))
	for i, p := range ps {
		out[i] = p.Packet.Sequence
	}
	return out
}

func cosmosMsgSeqs(ms []cosmosBatchMsg) []uint64 {
	out := make([]uint64, len(ms))
	for i, m := range ms {
		out[i] = m.sequence
	}
	return out
}

func ethMsgSeqs(ms []ethBatchMsg) []uint64 {
	out := make([]uint64, len(ms))
	for i, m := range ms {
		out[i] = m.sequence
	}
	return out
}

// assertSeqs compares as a sorted slice so duplicates (a packet routed twice) and
// missing/extra entries both fail.
func assertSeqs(t *testing.T, label string, got, want []uint64) {
	t.Helper()
	g := append([]uint64(nil), got...)
	w := append([]uint64(nil), want...)
	slices.Sort(g)
	slices.Sort(w)
	if fmt.Sprint(g) != fmt.Sprint(w) {
		t.Fatalf("%s: got %v, want %v", label, g, w)
	}
}

func TestPlanCosmosPacketMsgs_RoutesFailuresAndSkips(t *testing.T) {
	const router = "router-client"
	failSeqs := map[uint64]bool{2: true, 5: true, 8: true}
	build := func(p channeltypesv2.Packet, _ string, _ []byte) ([]byte, error) {
		if failSeqs[p.Sequence] {
			return nil, fmt.Errorf("boom seq=%d", p.Sequence)
		}
		return []byte{0xAA}, nil
	}

	packets := []CosmosPacket{
		cosmosPkt(1, CosmosSend, "src", "dst", 0, nil),             // ok
		cosmosPkt(2, CosmosSend, "src", "dst", 0, nil),             // proof fail
		cosmosPkt(3, CosmosSend, "src", "dst", 500, nil),           // timed out (ethBlockTime>=500)
		cosmosPkt(4, CosmosAck, "src", "dst", 0, [][]byte{{0x01}}), // ok
		cosmosPkt(5, CosmosAck, "src", "dst", 0, [][]byte{{0x01}}), // proof fail
		cosmosPkt(6, CosmosAck, "src", "dst", 0, nil),              // ack bytes missing → skip
		cosmosPkt(7, CosmosTimeout, "src", "other", 0, nil),        // shouldRelay true, ok
		cosmosPkt(8, CosmosTimeout, "src", "other", 0, nil),        // shouldRelay true, proof fail
		cosmosPkt(9, CosmosTimeout, "src", router, 0, nil),         // shouldRelay false → skip
	}

	msgs, trackerAdds, transientFailures, trackerRemoves := planCosmosPacketMsgs(
		packets, 1000, router, build, build)

	assertSeqs(t, "msgs", cosmosMsgSeqs(msgs), []uint64{1, 4, 7})
	assertSeqs(t, "trackerAdds", cosmosSeqs(trackerAdds), []uint64{1, 2, 3}) // every CosmosSend
	assertSeqs(t, "transientFailures", cosmosSeqs(transientFailures), []uint64{2, 5, 8})
	assertSeqs(t, "trackerRemoves", cosmosSeqs(trackerRemoves), []uint64{2}) // only the failed CosmosSend
}

func TestPlanCosmosPacketMsgs_AllSucceed_NoRequeueNoRemove(t *testing.T) {
	build := func(channeltypesv2.Packet, string, []byte) ([]byte, error) { return []byte{0xAA}, nil }
	packets := []CosmosPacket{
		cosmosPkt(1, CosmosSend, "src", "dst", 0, nil),
		cosmosPkt(2, CosmosAck, "src", "dst", 0, [][]byte{{0x01}}),
	}
	msgs, trackerAdds, transientFailures, trackerRemoves := planCosmosPacketMsgs(packets, 0, "router", build, build)

	assertSeqs(t, "msgs", cosmosMsgSeqs(msgs), []uint64{1, 2})
	assertSeqs(t, "trackerAdds", cosmosSeqs(trackerAdds), []uint64{1})
	if len(transientFailures) != 0 || len(trackerRemoves) != 0 {
		t.Fatalf("all-success batch must not re-queue or remove: failures=%v removes=%v",
			cosmosSeqs(transientFailures), cosmosSeqs(trackerRemoves))
	}
}

func TestPlanEthPacketMsgs_RoutesFailuresAndExpired(t *testing.T) {
	failPaths := map[string]bool{
		string(EthPath("clientA", 2, 1)): true, // EthSend seq=2
		string(EthPath("clientA", 5, 3)): true, // EthWriteAck seq=5
	}
	build := func(path []byte) ([]byte, error) {
		if failPaths[string(path)] {
			return nil, fmt.Errorf("boom")
		}
		return []byte{0xBB}, nil
	}

	relayable := []relayablePacket{
		{packet: ethPkt(1, EthSend, "clientA", "", 0, nil)},                  // ok
		{packet: ethPkt(2, EthSend, "clientA", "", 0, nil)},                  // proof fail
		{packet: ethPkt(3, EthSend, "clientA", "", 1, nil)},                  // expired (timeout in 1970)
		{packet: ethPkt(4, EthWriteAck, "", "clientA", 0, [][]byte{{0x01}})}, // ok
		{packet: ethPkt(5, EthWriteAck, "", "clientA", 0, [][]byte{{0x01}})}, // proof fail
	}

	msgs, expired, transientFailures := planEthPacketMsgs(relayable, 99, "signer", build)

	assertSeqs(t, "msgs", ethMsgSeqs(msgs), []uint64{1, 4})
	assertSeqs(t, "expired", ethSeqs(expired), []uint64{3})
	assertSeqs(t, "transientFailures", ethSeqs(transientFailures), []uint64{2, 5})
}
