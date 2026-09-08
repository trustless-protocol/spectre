package services

import (
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func queuedPacket(sequence uint64, src, dst string) CosmosPacket {
	return CosmosPacket{
		Type: CosmosSend,
		Packet: &channeltypesv2.Packet{
			Sequence: sequence, SourceClient: src, DestinationClient: dst,
		},
	}
}

// queuedPayload is the same packet with a distinguishing payload and timeout --
// what makes a REPLAYED sequence a different packet rather than the same one.
func queuedPayload(sequence uint64, src, dst, value string, timeout uint64) CosmosPacket {
	p := queuedPacket(sequence, src, dst)
	p.Packet.TimeoutTimestamp = timeout
	p.Packet.Payloads = []channeltypesv2.Payload{{
		SourcePort: "transfer", DestinationPort: "transfer",
		Version: "ics20-2", Encoding: "application/x-solidity-abi",
		Value: []byte(value),
	}}
	return p
}

// The mirror of the L2 fix. A recovery pass can queue a send and then read the
// acknowledge_packet that closes it, further down the same list. Settling only
// the tracker leaves the queue to relay a finished packet -- and the relay
// re-adds the tracker entry that settling just removed, after the cursor moved
// past the terminal event that would clear it again.
func TestDropCosmosQueued(t *testing.T) {
	newBuilder := func(t *testing.T) *BatchBuilder {
		t.Helper()
		return NewBatchBuilder()
	}

	t.Run("removes the settled packet and reports it", func(t *testing.T) {
		b := newBuilder(t)
		b.AddCosmos(queuedPacket(5, "08-wasm-7", "client-0"))
		b.AddCosmos(queuedPacket(6, "08-wasm-7", "client-0"))

		if !b.DropCosmosQueued(*queuedPacket(5, "08-wasm-7", "client-0").Packet) {
			t.Fatal("a queued packet that settled must be reported as removed")
		}
		cosmos, _ := b.QueueDepths()
		if cosmos != 1 {
			t.Fatalf("queue depth = %d, want 1", cosmos)
		}
	})

	t.Run("a packet that is not queued reports false", func(t *testing.T) {
		b := newBuilder(t)
		b.AddCosmos(queuedPacket(6, "08-wasm-7", "client-0"))
		if b.DropCosmosQueued(*queuedPacket(5, "08-wasm-7", "client-0").Packet) {
			t.Fatal("nothing was queued for that identity; claiming a removal makes the log line a lie")
		}
		if cosmos, _ := b.QueueDepths(); cosmos != 1 {
			t.Fatalf("queue depth = %d, want the untouched 1", cosmos)
		}
	})

	// Sequence alone is not an identity: two clients number their packets
	// independently, so matching on it would drop an unrelated live packet.
	t.Run("the same sequence on another client survives", func(t *testing.T) {
		b := newBuilder(t)
		b.AddCosmos(queuedPacket(5, "08-wasm-7", "client-0"))
		b.AddCosmos(queuedPacket(5, "08-wasm-9", "client-0"))

		b.DropCosmosQueued(*queuedPacket(5, "08-wasm-7", "client-0").Packet)
		if cosmos, _ := b.QueueDepths(); cosmos != 1 {
			t.Fatalf("queue depth = %d, want 1 -- the other client's packet must remain", cosmos)
		}
	})
}

// A client migration keeps the client id and restarts sequences from 1 --
// create-clients-eth repoints the existing client rather than adding one, which
// is a documented operational step here, not a hypothetical. So seq=N names two
// different packets over the life of one client, and a stale terminal event for
// the OLD seq=N must not drop the NEW one: that packet would never be relayed,
// and nothing downstream is left to notice.
func TestDropCosmosQueuedKeepsAReplayedSequence(t *testing.T) {
	b := NewBatchBuilder()
	fresh := queuedPayload(5, "08-wasm-7", "client-0", "new-transfer", 2000)
	b.AddCosmos(fresh)

	stale := queuedPayload(5, "08-wasm-7", "client-0", "old-transfer", 1000)
	if b.DropCosmosQueued(*stale.Packet) {
		t.Fatal("a terminal event for the pre-migration seq=5 claimed to remove the post-migration one")
	}
	if cosmos, _ := b.QueueDepths(); cosmos != 1 {
		t.Fatalf("queue depth = %d, want the replayed packet still queued: it has never been relayed", cosmos)
	}

	// The packet it really names is still removable.
	if !b.DropCosmosQueued(*fresh.Packet) {
		t.Fatal("the queued packet's own identity must still match")
	}
	if cosmos, _ := b.QueueDepths(); cosmos != 0 {
		t.Fatalf("queue depth = %d, want 0", cosmos)
	}
}
