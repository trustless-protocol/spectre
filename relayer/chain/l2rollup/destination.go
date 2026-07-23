// Package l2 holds the relayer adapters for the L2->Cosmos path (Arbitrum /
// OP-Stack). Per Dũng's light-client design the Cosmos side is TRUSTLESS — no
// relayer signature: the L2 wasm light client verifies L2 state against the shared
// Ethereum light client + L1 rollup proofs. So this Destination mirrors the
// existing beacon (ETH->Cosmos) Cosmos destination almost exactly:
//
//   - UpdateClient submits MsgUpdateClient wrapping a wasm ClientMessage whose
//     Data is the L2 header (the builder produces it). The SHARED ETH client must
//     be updated first (ordering dependency, see TODO).
//   - RelayPackets submits the standard channeltypesv2 recv / ack / timeout
//     messages, which the client maps to VerifyMembership / VerifyNonMembership.
//     The L2 client does not decode packets.
//
// Skeleton status: the message shapes are real (same wasm/channeltypesv2 wrappers
// fast-ibc already uses); the proof height needs the L2 client-state read and the
// ETH-first ordering needs wiring — both marked TODO and pending Dũng's exact
// client-state / ClientMessage schema.
package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"relayer/chain"
	"relayer/services"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// Destination satisfies chain.Destination.
var _ chain.Destination = (*Destination)(nil)

// isPermanentCosmosFailure reports whether a Cosmos submit error is a
// deterministic on-chain revert (drop) vs a transient one (retry) — same
// classification the Cosmos beacon destination uses.
func isPermanentCosmosFailure(err error) bool {
	return errors.Is(err, services.ErrPermanentRelayFailure)
}

// Destination is the Cosmos side of the L2->Cosmos path: it hosts the L2 wasm
// light client and is where the L2 client update + L2-origin packets are
// submitted. Chain() is Cosmos (this is the destination chain).
type Destination struct {
	worker   *services.Worker
	svcCtx   services.Context
	clientID string // the L2 wasm light-client id on Cosmos
}

// NewDestination wires the L2 Cosmos destination. clientID is the L2 wasm client.
func NewDestination(worker *services.Worker, svcCtx services.Context, clientID string) *Destination {
	return &Destination{worker: worker, svcCtx: svcCtx, clientID: clientID}
}

func (d *Destination) Chain() chain.ChainType { return chain.Cosmos }

// UpdateClient submits MsgUpdateClient for the L2 wasm client, wrapping the L2
// header (update.Payload, produced by the l2 builder) in a wasm ClientMessage.
//
// TODO(Đức, Dũng P1): the SHARED Ethereum light client must be advanced first (the
// L2 client verifies L2 roots against it + L1 rollup proofs). Wire an
// "update ETH client to cover this L2 update" step before the submit, mirroring
// the beacon path's WaitForCosmosCatchUp ordering.
func (d *Destination) UpdateClient(_ context.Context, _ string, update chain.ClientUpdate) error {
	if len(update.Payload) == 0 {
		return nil // nothing to submit
	}
	signer, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("l2 dest: signer address: %w", err)
	}
	msg, err := buildWasmUpdateClient(signer, d.clientID, update.Payload)
	if err != nil {
		return fmt.Errorf("l2 dest: build update (height %d): %w", update.Height, err)
	}
	if err := d.worker.TxHandler.SendCosmosTxBatch(d.svcCtx, []any{msg}); err != nil {
		return fmt.Errorf("l2 dest: submit update (height %d): %w", update.Height, err)
	}
	return nil
}

// RelayPackets builds the standard channeltypesv2 messages (recv / ack / timeout)
// from the L2-origin packets + storage proofs and submits them. Same wrappers as
// the existing Cosmos beacon destination; only the proof height differs (the L2
// client's latest tracked height instead of the beacon slot).
func (d *Destination) RelayPackets(_ context.Context, packets []chain.RelayPacket) error {
	if len(packets) == 0 {
		return nil
	}
	signer, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("l2 dest: signer address: %w", err)
	}
	proofHeight, err := d.proofHeight()
	if err != nil {
		return err
	}

	msgs := make([]any, 0, len(packets))
	for _, rp := range packets {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(rp.Packet); err != nil {
			return fmt.Errorf("l2 dest: decode packet: %w", err)
		}
		switch rp.Type {
		case chain.SendPacket:
			msgs = append(msgs, &channeltypesv2.MsgRecvPacket{
				Packet:          pkt,
				ProofCommitment: rp.Proof,
				ProofHeight:     proofHeight,
				Signer:          signer,
			})
		case chain.AckPacket:
			msgs = append(msgs, &channeltypesv2.MsgAcknowledgement{
				Packet:          pkt,
				Acknowledgement: channeltypesv2.Acknowledgement{AppAcknowledgements: rp.AckBytes},
				ProofAcked:      rp.Proof,
				ProofHeight:     proofHeight,
				Signer:          signer,
			})
		case chain.TimeoutPacket:
			msgs = append(msgs, &channeltypesv2.MsgTimeout{
				Packet:          pkt,
				ProofUnreceived: rp.Proof,
				ProofHeight:     proofHeight,
				Signer:          signer,
			})
		default:
			return fmt.Errorf("l2 dest: unsupported packet type %d (seq=%d)", rp.Type, pkt.Sequence)
		}
	}
	if err := d.worker.TxHandler.SendCosmosTxBatch(d.svcCtx, msgs); err != nil {
		if isPermanentCosmosFailure(err) {
			return chain.Permanent(err)
		}
		return err
	}
	return nil
}

// proofHeight is the L2 client's latest tracked height — the height the packet
// proofs must be verified against.
//
// TODO(Đức, Dũng P1): read the L2 wasm client state's latest height (Dũng's
// client stores state_root + IBC storage_root + timestamp at the L2 block number).
// Stubbed until the L2 client-state schema is fixed.
func (d *Destination) proofHeight() (clienttypes.Height, error) {
	return clienttypes.Height{}, fmt.Errorf("l2 dest: proofHeight not wired (needs L2 client-state read)")
}

// HasPacketReceipt reports whether an L2-origin packet was already delivered on
// Cosmos. TODO: verified Cosmos receipt ABCI query (mirrors the beacon dest stub).
func (d *Destination) HasPacketReceipt(_ context.Context, _ []byte) (bool, error) {
	return false, fmt.Errorf("l2 dest: HasPacketReceipt not wired")
}

// ClientExpiresAt returns when the L2 wasm client expires. TODO: read the L2
// client state's trusting period / timestamp.
func (d *Destination) ClientExpiresAt(_ context.Context, _ string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("l2 dest: ClientExpiresAt not wired")
}

// buildWasmUpdateClient wraps an L2 header (JSON bytes) in a wasm ClientMessage
// and a MsgUpdateClient — the exact wrapper fast-ibc uses for the beacon client
// (services.buildMsgUpdateClient), reused for the L2 client.
func buildWasmUpdateClient(signer, clientID string, l2Header []byte) (*clienttypes.MsgUpdateClient, error) {
	clientMessage := &ibcwasmtypes.ClientMessage{Data: l2Header}
	anyMsg, err := codectypes.NewAnyWithValue(clientMessage)
	if err != nil {
		return nil, fmt.Errorf("wrap client message: %w", err)
	}
	return &clienttypes.MsgUpdateClient{
		ClientId:      clientID,
		ClientMessage: anyMsg,
		Signer:        signer,
	}, nil
}
