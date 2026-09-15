package main

import (
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"

	"relayer/services"
)

func ackTestPacket(sequence uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{
		Sequence: sequence, SourceClient: "08-wasm-0", DestinationClient: "eth-router-0",
	}
}

// The two halves of an owed acknowledgement land on OPPOSITE modules: a
// Cosmos-origin send is delivered by cosmos->eth (the debt) and the
// acknowledgement written for it comes back as an ETH WriteAcknowledgement
// relayed by eth->cosmos (the settlement), and the mirror image for an
// ETH-origin send. With the hooks on one module only, every module recorded
// debts the other settled, so nothing ever cleared.
//
// This asserts the shared tracker makes cross-module settlement work: record
// through one module's hook, clear through the other's.
func TestOwedAckSettlesAcrossModules(t *testing.T) {
	svc := services.New(nil, nil, services.DefaultConfig())

	packet := ackTestPacket(7)
	if !svc.RecordAckDue(packet) {
		t.Fatal("recording the owed acknowledgement failed")
	}
	if got := svc.OverdueAcks(0); len(got) != 1 {
		t.Fatalf("overdue = %v, want the freshly recorded debt", got)
	}

	// The settlement arrives from the other direction's module, through the same
	// Services and therefore the same tracker.
	svc.ClearAckDue(packet)
	if got := svc.OverdueAcks(0); len(got) != 0 {
		t.Fatalf("overdue = %v after the acknowledgement was relayed, want none", got)
	}
}

// A nil AckDueTracker panics on the first RecordAckDue, and Services.New uses
// the default constructor -- so this is the ordinary adapter-engine path.
func TestDefaultServicesCanRecordAnOwedAck(t *testing.T) {
	svc := services.New(nil, nil, services.DefaultConfig())
	if svc.BatchBuilder.AckDueTracker == nil {
		t.Fatal("the default batch builder left AckDueTracker nil; RecordAckDue would panic")
	}
	if !svc.RecordAckDue(ackTestPacket(1)) {
		t.Fatal("recording an owed acknowledgement on a default Services failed")
	}
}
