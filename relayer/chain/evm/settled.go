package evm

import (
	"context"
	"fmt"

	"relayer/chain"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
)

// PacketSettled implements chain.SettlementChecker for the ETH router.
//
//   - a receive is settled once the router holds the packet's receipt;
//   - an acknowledgement (or timeout) is settled once the router no longer
//     holds the original send's commitment -- ackPacket and timeoutPacket both
//     delete it, and a second submission is answered with Noop.
func (d *Destination) PacketSettled(ctx context.Context, e chain.Event) (bool, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(e.Raw); err != nil {
		return false, fmt.Errorf("evm: decode packet: %w", err)
	}
	switch e.Type {
	case chain.SendPacket:
		return services.HasEthPacketReceipt(ctx, d.evm, pkt)
	case chain.AckPacket, chain.TimeoutPacket:
		pending, err := services.HasPendingEthPacketCommitment(ctx, d.evm, pkt)
		if err != nil {
			return false, err
		}
		return !pending, nil
	}
	return false, fmt.Errorf("evm: settlement of %s is not defined", e.Type)
}

var _ chain.SettlementChecker = (*Destination)(nil)

// HasPacketCommitment implements chain.CommitmentChecker: whether the router
// still holds the commitment of a packet it sent (ETH->Cosmos).
func (s *Source) HasPacketCommitment(ctx context.Context, packet []byte) (bool, error) {
	// An empty slot reads as "gone", so a missing router address must not be
	// asked: every answer would be "gone".
	if s.evm.EthClient() == nil || s.evm.Contracts.Router == (ethcommon.Address{}) {
		return false, fmt.Errorf("eth source: no ETH client or router address to query the commitment on")
	}
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return false, fmt.Errorf("eth source: decode packet: %w", err)
	}
	return services.HasPendingEthPacketCommitment(ctx, s.evm, pkt)
}

var _ chain.CommitmentChecker = (*Source)(nil)
