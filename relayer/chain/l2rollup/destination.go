// Package l2rollup holds the relayer adapters for the L2->Cosmos path (Arbitrum /
// OP-Stack). Per Dũng's light-client design the Cosmos side is TRUSTLESS — no
// relayer signature: the L2 wasm light client verifies L2 state against the shared
// Ethereum light client + L1 rollup proofs. So this Destination mirrors the
// existing beacon (ETH->Cosmos) Cosmos destination almost exactly:
//
//   - UpdateClient submits MsgUpdateClient wrapping a wasm ClientMessage whose Data
//     is the L2 header (the builder produces it). No ETH-first ordering step is
//     needed: the header builder proves against the ETH client's ALREADY-trusted L1
//     block (it reads EthClientLatestSlotAndBlock), so the update verifies against
//     current ETH state. ETH-client freshness only bounds how recent an L2 update
//     can be, it is not a correctness ordering requirement.
//   - RelayPackets submits the standard channeltypesv2 recv / ack / timeout
//     messages, which the client maps to VerifyMembership / VerifyNonMembership.
//     The L2 client does not decode packets.
package l2rollup

import (
	"context"
	"errors"
	"fmt"
	"time"

	"relayer/chain"
	"relayer/chain/wasmclient"
	relayerclient "relayer/client"
	"relayer/services"
	"relayer/subscriber"

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
	cosmos   services.CosmosEndpoint
	clientID string // the L2 wasm light-client id on Cosmos
}

// NewDestination wires the L2 Cosmos destination. clientID is the L2 wasm client.
func NewDestination(worker *services.Worker, cosmos services.CosmosEndpoint, clientID string) *Destination {
	return &Destination{worker: worker, cosmos: cosmos, clientID: clientID}
}

func (d *Destination) Chain() chain.ChainType { return chain.Cosmos }

// UpdateClient submits MsgUpdateClient for the L2 wasm client, wrapping the L2
// header (the sole update payload, produced by the L2 builder) in a wasm
// ClientMessage.
//
// No ETH-first ordering step is needed: the header builder assembles its L1 proofs
// against the ETH client's already-trusted L1 block (EthClientLatestSlotAndBlock), so
// every L2 update verifies against the ETH state the client already holds. The ETH
// client's freshness bounds how recent an L2 update can be, but advancing it is not a
// prerequisite for this submit to verify.
func (d *Destination) UpdateClient(ctx context.Context, _ string, update chain.ClientUpdate) error {
	if len(update.Payloads) != 1 {
		return fmt.Errorf("l2 dest: expected exactly one client update payload, got %d", len(update.Payloads))
	}
	signer, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("l2 dest: signer address: %w", err)
	}
	msg, err := wasmclient.BuildUpdateClient(signer, d.clientID, update.Payloads[0])
	if err != nil {
		return fmt.Errorf("l2 dest: build update (height %d): %w", update.Height, err)
	}
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.cosmos, []any{msg}); err != nil {
		return fmt.Errorf("l2 dest: submit update (height %d): %w", update.Height, err)
	}
	return nil
}

// RelayPackets builds the standard channeltypesv2 messages (recv / ack / timeout)
// from the L2-origin packets + storage proofs and submits them. Same wrappers as
// the existing Cosmos beacon destination; only the proof height differs (the L2
// client's latest tracked height instead of the beacon slot).
func (d *Destination) RelayPackets(ctx context.Context, packets []chain.RelayPacket) error {
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
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.cosmos, msgs); err != nil {
		if isPermanentCosmosFailure(err) {
			return chain.Permanent(err)
		}
		return err
	}
	return nil
}

// proofHeight is the L2 client's latest tracked height — the height the packet
// proofs must be verified against. The wasm ClientState carries LatestHeight
// directly; for the L2 client its revision height is the L2 block number.
func (d *Destination) proofHeight() (clienttypes.Height, error) {
	h, err := relayerclient.GetWasmClientLatestHeight(d.cosmos.CosmosClient(), d.clientID)
	if err != nil {
		return clienttypes.Height{}, fmt.Errorf("l2 dest: read L2 client latest height: %w", err)
	}
	return h, nil
}

// HasPacketReceipt reports whether an L2-origin packet was already delivered on
// Cosmos, via the same Cosmos receipt ABCI query the ETH recovery path uses.
func (d *Destination) HasPacketReceipt(_ context.Context, packet []byte) (bool, error) {
	var pkt channeltypesv2.Packet
	if err := pkt.Unmarshal(packet); err != nil {
		return false, fmt.Errorf("l2 dest: decode packet: %w", err)
	}
	return subscriber.HasCosmosPacketReceipt(d.cosmos, pkt)
}

// ClientExpiresAt reports when the L2 wasm client would expire on its own timer.
// Unlike a Tendermint client, an optimistic L2 client has NO independent trusting
// period: it advances per-packet via header updates and its trust roots in the
// SHARED L1 (Ethereum) client, whose freshness the ETH path refreshes. So the L2
// client has no self-expiry timer — return a far-future time so the anti-expiry
// refresh routine never force-updates it (its freshness is the L1 client's).
func (d *Destination) ClientExpiresAt(_ context.Context, _ string) (time.Time, error) {
	return time.Now().Add(100 * 365 * 24 * time.Hour), nil
}
