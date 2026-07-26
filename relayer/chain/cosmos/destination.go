package cosmos

import (
	"context"
	"errors"
	"fmt"
	"time"

	"relayer/chain"
	"relayer/chain/codec"
	relayerclient "relayer/client"
	"relayer/services"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// Destination is the Cosmos implementation of chain.Destination for the ETH->Cosmos
// direction: it hosts the 08-wasm Ethereum beacon light client and is where the
// beacon update + ETH-origin packets are submitted. It holds the shared worker +
// context because SendCosmosTxBatch and the pre-submit catch-up need them; it is
// constructed by the wiring, not the cfg-only registry.
type Destination struct {
	worker *services.Worker
	svcCtx services.Context
}

// NewDestination wires the Cosmos destination to the shared worker + context.
func NewDestination(worker *services.Worker, svcCtx services.Context) *Destination {
	return &Destination{worker: worker, svcCtx: svcCtx}
}

func (d *Destination) Chain() chain.ChainType { return chain.Cosmos }

// UpdateClient decodes the beacon update payload, waits for the Cosmos node to
// catch up to the signature slot (so the submitted proof height is queryable),
// then submits the MsgUpdateClient batch — exactly as the legacy UpdateEthClient
// path did.
func (d *Destination) UpdateClient(ctx context.Context, _ string, update chain.ClientUpdate) error {
	if err := ctx.Err(); err != nil {
		return err // shutting down — do not start a client-update tx
	}
	msgs, ethClientState, sigSlot, err := codec.DecodeBeaconUpdate(update.Payload)
	if err != nil {
		return fmt.Errorf("cosmos dest: decode beacon update: %w", err)
	}
	if len(msgs) == 0 {
		return nil
	}
	d.worker.WaitForCosmosCatchUp(ctx, d.svcCtx, &ethClientState, sigSlot)
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.svcCtx, msgs); err != nil {
		return fmt.Errorf("cosmos dest: submit beacon update (exec block %d): %w", update.Height, err)
	}
	return nil
}

// RelayPackets builds the Cosmos-bound IBC messages (recv/ack) from the given
// ETH-origin packets + storage proofs and submits them via SendCosmosTxBatch.
// The proof height is the on-chain 08-wasm client's latest slot (read here so it
// matches the execution block the ETH source built the proof against). ETH-origin
// timeouts are handled by the async scanner, not this path.
func (d *Destination) RelayPackets(ctx context.Context, packets []chain.RelayPacket) error {
	if len(packets) == 0 {
		return nil
	}
	if err := ctx.Err(); err != nil {
		return err // shutting down — do not start a packet tx
	}
	signerAddr, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("cosmos dest: signer address: %w", err)
	}
	ethClientState, err := relayerclient.GetEthereumClientState(d.svcCtx.CosmosClient(), d.svcCtx.EthClientID())
	if err != nil {
		return fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	proofHeight := clienttypes.Height{RevisionNumber: 0, RevisionHeight: ethClientState.LatestSlot}

	msgs := make([]any, 0, len(packets))
	for _, rp := range packets {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(rp.Packet); err != nil {
			// A malformed packet is deterministic — it will never decode on retry.
			return chain.Permanent(fmt.Errorf("cosmos dest: decode packet: %w", err))
		}
		switch rp.Type {
		case chain.SendPacket:
			msgs = append(msgs, &channeltypesv2.MsgRecvPacket{
				Packet:          pkt,
				ProofCommitment: rp.Proof,
				ProofHeight:     proofHeight,
				Signer:          signerAddr,
			})
		case chain.AckPacket:
			msgs = append(msgs, &channeltypesv2.MsgAcknowledgement{
				Packet:          pkt,
				Acknowledgement: channeltypesv2.Acknowledgement{AppAcknowledgements: rp.AckBytes},
				ProofAcked:      rp.Proof,
				ProofHeight:     proofHeight,
				Signer:          signerAddr,
			})
		default:
			return chain.Permanent(fmt.Errorf("cosmos dest: unsupported packet type %d (seq=%d)", rp.Type, pkt.Sequence))
		}
	}
	// A deterministic Cosmos DeliverTx revert cannot succeed on retry — surface it
	// as chain.Permanent so the module DROPS the batch rather than re-queueing it
	// forever (each retry re-runs the beacon client update and drains gas). CheckTx/
	// broadcast/RPC errors stay transient and are re-queued.
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.svcCtx, msgs); err != nil {
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			return chain.Permanent(err)
		}
		return err
	}
	return nil
}

// HasPacketReceipt reports whether an ETH-origin packet was already delivered on
// Cosmos. NOTE: not used by the current RelayModule flow (it drives timeout
// scanning, which is not yet wired into the module). A correct implementation
// needs a verified Cosmos-side packet-receipt ABCI query; left explicit rather
// than guessed.
func (d *Destination) HasPacketReceipt(_ context.Context, _ []byte) (bool, error) {
	return false, fmt.Errorf("cosmos dest: HasPacketReceipt not wired (timeout-scan, unused by module)")
}

// ClientExpiresAt returns a conservative expiry for the 08-wasm ETH light client
// on Cosmos. Unlike a Tendermint client there is no single trusting period; the
// client must be advanced within ~2 sync-committee periods of its latest tracked
// slot. We report one period after that slot's time — deliberately early so the
// refresh routine never lets it lapse (a wrong-too-late value would expire the
// client, the failure class we care about most).
func (d *Destination) ClientExpiresAt(_ context.Context, _ string) (time.Time, error) {
	cs, err := relayerclient.GetEthereumClientState(d.svcCtx.CosmosClient(), d.svcCtx.EthClientID())
	if err != nil {
		return time.Time{}, fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	periodSecs := cs.EpochsPerSyncCommitteePeriod * cs.SlotsPerEpoch * cs.SecondsPerSlot
	if periodSecs == 0 {
		return time.Time{}, nil // misconfigured — report no expiry, fall back to periodic refresh
	}
	latestSlotTime := cs.ComputeTimestampAtSlot(cs.LatestSlot)
	return time.Unix(int64(latestSlotTime+periodSecs), 0), nil
}
