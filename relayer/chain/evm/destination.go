// Package evm adapts an EVM chain (Ethereum L1, and the L2s Arbitrum/Optimism/
// Base) as a chain.Destination — where a Cosmos-sourced SpectreClient update is
// submitted. It wraps the existing services/transaction code.
package evm

import (
	"context"
	"errors"
	"fmt"
	"time"

	contractICS26Router "relayer/bindings/ICS26Router"
	"relayer/chain"
	"relayer/chain/codec"
	"relayer/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// Destination is the EVM implementation of chain.Destination for a chain that
// hosts a SpectreClient (ZK Tendermint) light client of Cosmos. It holds the
// shared worker + context because the tx handler and on-chain reads need them;
// it is constructed by the wiring, not the cfg-only registry.
type Destination struct {
	worker *services.Worker
	svcCtx services.Context
}

// NewDestination wires the EVM destination to the shared worker + context.
func NewDestination(worker *services.Worker, svcCtx services.Context) *Destination {
	return &Destination{worker: worker, svcCtx: svcCtx}
}

func (d *Destination) Chain() chain.ChainType { return chain.Ethereum }

// UpdateClient decodes the Cosmos->ETH update payload back into the typed message
// and submits it. SendEthTx routes updateApplicationState vs updateConsensusState
// off the result's Kind, exactly as the legacy RefreshCosmosClient path did.
func (d *Destination) UpdateClient(ctx context.Context, _ string, update chain.ClientUpdate) error {
	if err := ctx.Err(); err != nil {
		return err // shutting down — do not start a client-update tx
	}
	kind, appMsg, newValSet, err := codec.DecodeCosmosUpdate(update.Payload)
	if err != nil {
		return fmt.Errorf("evm: decode client update: %w", err)
	}
	result := services.CosmosClientUpdateBuildResult{
		Kind:      services.ClientUpdateKind(kind),
		AppMsg:    appMsg,
		NewValSet: newValSet,
		HasMsg:    true,
	}
	if err := d.worker.TxHandler.SendEthTx(ctx, d.svcCtx, result); err != nil {
		return fmt.Errorf("evm: submit client update (height %d): %w", update.Height, err)
	}
	return nil
}

// RelayPackets builds the concrete ICS26 messages (recv/ack/timeout) from the
// given source packets + injected proofs and submits them via SendEthTxBatch.
// The proofs are already built by the source (Membership/NonMembershipProof), so
// unlike the legacy planCosmosPacketMsgs this does no proof generation.
func (d *Destination) RelayPackets(ctx context.Context, packets []chain.RelayPacket) error {
	if err := ctx.Err(); err != nil {
		return err // shutting down — do not start a packet tx
	}
	msgs := make([]any, 0, len(packets))
	for _, rp := range packets {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(rp.Packet); err != nil {
			// A malformed packet is deterministic — it will never decode on retry.
			return chain.Permanent(fmt.Errorf("evm: decode packet: %w", err))
		}
		ethPkt := services.ToEthPacket(pkt)
		switch rp.Type {
		case chain.SendPacket:
			msgs = append(msgs, contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
				Packet:        ethPkt,
				MembershipMsg: rp.Proof,
			})
		case chain.AckPacket:
			if len(rp.AckBytes) == 0 {
				return chain.Permanent(fmt.Errorf("evm: ack packet seq=%d missing acknowledgement bytes", pkt.Sequence))
			}
			msgs = append(msgs, contractICS26Router.IICS26RouterMsgsMsgAckPacket{
				Packet:          ethPkt,
				Acknowledgement: rp.AckBytes[0], // ETH MsgAckPacket takes a single ack
				MembershipMsg:   rp.Proof,
			})
		case chain.TimeoutPacket:
			msgs = append(msgs, contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
				Packet:           ethPkt,
				NonMembershipMsg: rp.Proof,
			})
		default:
			return chain.Permanent(fmt.Errorf("evm: unknown packet type %d (seq=%d)", rp.Type, pkt.Sequence))
		}
	}
	if len(msgs) == 0 {
		return nil
	}
	// A deterministic on-chain revert (status=0) cannot succeed on retry — surface
	// it as chain.Permanent so the module DROPS the batch instead of re-queueing it
	// forever (each retry re-runs the client update and drains gas). Everything else
	// (nonce, RPC, broadcast) stays transient and is re-queued.
	if err := d.worker.TxHandler.SendEthTxBatch(ctx, d.svcCtx, msgs); err != nil {
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			return chain.Permanent(err)
		}
		return err
	}
	return nil
}

// HasPacketReceipt reports whether a packet was already delivered on ETH. The
// packet is the proto-marshaled channeltypesv2.Packet (the convention shared
// with Event.Raw and RelayPackets).
func (d *Destination) HasPacketReceipt(_ context.Context, packet []byte) (bool, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return false, fmt.Errorf("evm: decode packet: %w", err)
	}
	return services.HasEthPacketReceipt(d.svcCtx, pkt)
}

// ClientExpiresAt returns when the SpectreClient-of-Cosmos on this chain expires
// (trusted consensus timestamp + trusting period), so the RelayModule refresh
// routine can advance the client before it lapses.
func (d *Destination) ClientExpiresAt(_ context.Context, _ string) (time.Time, error) {
	return services.CosmosClientExpiry(d.svcCtx)
}
