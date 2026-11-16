package cosmos

import (
	"context"
	"fmt"

	"relayer/chain"
	"relayer/subscriber"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// PacketSettled implements chain.SettlementChecker for the Cosmos side of the
// ETH->Cosmos path, against the same ABCI store the L2 destination and the ETH
// recovery scan already query.
//
//   - a receive is settled once Cosmos holds the packet's receipt;
//   - an acknowledgement (or timeout) is settled once Cosmos no longer holds
//     the original send's commitment.
func (d *Destination) PacketSettled(ctx context.Context, e chain.Event) (bool, error) {
	if d.cosmos.CosmosClient() == nil {
		return false, fmt.Errorf("cosmos dest: no Cosmos client to query settlement on")
	}
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(e.Raw); err != nil {
		return false, fmt.Errorf("cosmos dest: decode packet: %w", err)
	}
	switch e.Type {
	case chain.SendPacket:
		return subscriber.HasCosmosPacketReceipt(ctx, d.cosmos, pkt)
	case chain.AckPacket, chain.TimeoutPacket:
		pending, err := subscriber.HasPendingCosmosPacketCommitment(ctx, d.cosmos, pkt)
		if err != nil {
			return false, err
		}
		return !pending, nil
	}
	return false, fmt.Errorf("cosmos dest: settlement of %s is not defined", e.Type)
}

var _ chain.SettlementChecker = (*Destination)(nil)

// HasPacketCommitment implements chain.CommitmentChecker: whether Cosmos still
// holds the commitment of a packet it sent (Cosmos->ETH, Cosmos->L2).
func (s *Source) HasPacketCommitment(ctx context.Context, packet []byte) (bool, error) {
	if s.cosmos.CosmosClient() == nil {
		return false, fmt.Errorf("cosmos source: no Cosmos client to query the commitment on")
	}
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return false, fmt.Errorf("cosmos source: decode packet: %w", err)
	}
	return subscriber.HasPendingCosmosPacketCommitment(ctx, s.cosmos, pkt)
}

var _ chain.CommitmentChecker = (*Source)(nil)
