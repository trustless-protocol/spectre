package cosmos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"relayer/chain"
	"relayer/chain/wasmclient"
	relayerclient "relayer/client"
	"relayer/services"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// Destination is the Cosmos implementation of chain.Destination for the ETH->Cosmos
// direction: it hosts the 08-wasm Ethereum beacon light client and is where the
// beacon update + ETH-origin packets are submitted. It holds the worker, Cosmos
// endpoint, and Ethereum client ID needed by that path; it is constructed by the
// wiring, not the cfg-only registry.
type Destination struct {
	worker   *services.Worker
	cosmos   services.CosmosEndpoint
	clientID string
}

type atomicCosmosTxHandler interface {
	SendCosmosTxBatchAtomic(context.Context, services.CosmosEndpoint, []any) error
}

// NewDestination wires the Cosmos destination to the worker and Cosmos endpoint.
func NewDestination(worker *services.Worker, cosmos services.CosmosEndpoint, clientID string) *Destination {
	return &Destination{worker: worker, cosmos: cosmos, clientID: clientID}
}

func (d *Destination) Chain() chain.ChainType { return chain.Cosmos }

func (d *Destination) SupportsUpdatePacketFolding() bool {
	if d.worker == nil || d.worker.TxHandler == nil {
		return false
	}
	_, ok := d.worker.TxHandler.(atomicCosmosTxHandler)
	return ok
}

// UpdateClient wraps each beacon JSON header in a wasm ClientMessage, waits for
// the Cosmos node to reach the latest signature slot, then submits the ordered
// MsgUpdateClient batch. This preserves the legacy pre-submit timing guard
// without carrying relayer-only metadata inside the payload.
func (d *Destination) UpdateClient(ctx context.Context, clientID string, update chain.ClientUpdate) error {
	if err := ctx.Err(); err != nil {
		return err // shutting down — do not start a client-update tx
	}
	if len(update.Payloads) == 0 {
		return nil // client already current
	}
	if clientID == "" {
		clientID = d.clientID
	}
	signer, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("cosmos dest: signer address: %w", err)
	}
	sigSlot, err := highestBeaconSignatureSlot(update.Payloads)
	if err != nil {
		return fmt.Errorf("cosmos dest: decode beacon update: %w", err)
	}
	msgs, err := buildBeaconUpdateMessages(signer, clientID, update.Payloads)
	if err != nil {
		return err
	}
	readCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	ethClientState, err := relayerclient.GetEthereumClientStateWithContext(readCtx, d.cosmos.CosmosClient(), d.clientID)
	if err != nil {
		return fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	if err := d.worker.WaitForCosmosCatchUp(ctx, d.cosmos, ethClientState, sigSlot); err != nil {
		return fmt.Errorf("cosmos dest: wait for chain catch-up: %w", err)
	}
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.cosmos, msgs); err != nil {
		return fmt.Errorf("cosmos dest: submit beacon update (exec block %d): %w", update.Height, err)
	}
	return nil
}

func buildBeaconUpdateMessages(signer, clientID string, payloads [][]byte) ([]any, error) {
	msgs := make([]any, 0, len(payloads))
	for i, payload := range payloads {
		msg, err := wasmclient.BuildUpdateClient(signer, clientID, payload)
		if err != nil {
			return nil, fmt.Errorf("cosmos dest: wrap beacon update %d: %w", i, err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// highestBeaconSignatureSlot reads the timing prerequisite from the same JSON
// payload that the wasm client verifies. The last header normally has the latest
// slot, but taking the maximum makes ordered period-crossing batches explicit.
func highestBeaconSignatureSlot(payloads [][]byte) (uint64, error) {
	var latest uint64
	for i, payload := range payloads {
		var header relayerclient.EthereumHeader
		if err := json.Unmarshal(payload, &header); err != nil {
			return 0, fmt.Errorf("header %d: unmarshal JSON: %w", i, err)
		}
		slot, err := strconv.ParseUint(header.ConsensusUpdate.SignatureSlot, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("header %d: parse signature slot: %w", i, err)
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest, nil
}

// highestBeaconFinalizedSlot returns the destination consensus height installed
// by the ordered beacon-update payloads. Packet messages folded behind those
// updates must name this slot even though it is not on-chain yet.
func highestBeaconFinalizedSlot(payloads [][]byte) (uint64, error) {
	var latest uint64
	for i, payload := range payloads {
		var header relayerclient.EthereumHeader
		if err := json.Unmarshal(payload, &header); err != nil {
			return 0, fmt.Errorf("header %d: unmarshal JSON: %w", i, err)
		}
		slot, err := strconv.ParseUint(header.ConsensusUpdate.FinalizedHeader.Beacon.Slot, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("header %d: parse finalized slot: %w", i, err)
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest, nil
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
	readCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	ethClientState, err := relayerclient.GetEthereumClientStateWithContext(readCtx, d.cosmos.CosmosClient(), d.clientID)
	if err != nil {
		return fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	proofHeight := clienttypes.Height{RevisionNumber: 0, RevisionHeight: ethClientState.LatestSlot}
	msgs, err := buildPacketMessages(signerAddr, proofHeight, packets)
	if err != nil {
		return err
	}
	return d.sendPacketBatch(ctx, msgs)
}

// RelayWithUpdate submits the ordered beacon updates and all packet messages in
// one Cosmos transaction. Packet proof heights use the finalized slot installed
// by the update rather than querying the still-stale on-chain client state.
func (d *Destination) RelayWithUpdate(ctx context.Context, clientID string, update chain.ClientUpdate, packets []chain.RelayPacket) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(update.Payloads) == 0 {
		return fmt.Errorf("cosmos dest: folded relay requires a client update payload")
	}
	if d.worker == nil || d.worker.TxHandler == nil {
		return fmt.Errorf("cosmos dest: transaction handler is not configured")
	}
	atomicSender, ok := d.worker.TxHandler.(atomicCosmosTxHandler)
	if !ok {
		return fmt.Errorf("cosmos dest: transaction handler does not support atomic update/packet batches")
	}
	if clientID == "" {
		clientID = d.clientID
	}
	signer, err := d.worker.TxHandler.CosmosSignerAddress()
	if err != nil {
		return fmt.Errorf("cosmos dest: signer address: %w", err)
	}
	updateMsgs, err := buildBeaconUpdateMessages(signer, clientID, update.Payloads)
	if err != nil {
		return err
	}
	sigSlot, err := highestBeaconSignatureSlot(update.Payloads)
	if err != nil {
		return fmt.Errorf("cosmos dest: decode beacon update: %w", err)
	}
	finalizedSlot, err := highestBeaconFinalizedSlot(update.Payloads)
	if err != nil {
		return fmt.Errorf("cosmos dest: decode beacon update: %w", err)
	}
	packetMsgs, err := buildPacketMessages(signer, clienttypes.Height{RevisionNumber: 0, RevisionHeight: finalizedSlot}, packets)
	if err != nil {
		return err
	}
	readCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	ethClientState, err := relayerclient.GetEthereumClientStateWithContext(readCtx, d.cosmos.CosmosClient(), d.clientID)
	if err != nil {
		return fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	if err := d.worker.WaitForCosmosCatchUp(ctx, d.cosmos, ethClientState, sigSlot); err != nil {
		return fmt.Errorf("cosmos dest: wait for chain catch-up: %w", err)
	}

	msgs := make([]any, 0, len(updateMsgs)+len(packetMsgs))
	msgs = append(msgs, updateMsgs...)
	msgs = append(msgs, packetMsgs...)
	if err := atomicSender.SendCosmosTxBatchAtomic(ctx, d.cosmos, msgs); err != nil {
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			return chain.Permanent(err)
		}
		return err
	}
	return nil
}

func buildPacketMessages(signer string, proofHeight clienttypes.Height, packets []chain.RelayPacket) ([]any, error) {
	msgs := make([]any, 0, len(packets))
	for _, rp := range packets {
		var pkt channeltypesv2.Packet
		if err := pkt.Unmarshal(rp.Packet); err != nil {
			// A malformed packet is deterministic — it will never decode on retry.
			return nil, chain.Permanent(fmt.Errorf("cosmos dest: decode packet: %w", err))
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
		default:
			return nil, chain.Permanent(fmt.Errorf("cosmos dest: unsupported packet type %d (seq=%d)", rp.Type, pkt.Sequence))
		}
	}
	return msgs, nil
}

func (d *Destination) sendPacketBatch(ctx context.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}
	// A deterministic Cosmos DeliverTx revert cannot succeed on retry — surface it
	// as chain.Permanent so the module DROPS the batch rather than re-queueing it
	// forever (each retry re-runs the beacon client update and drains gas). CheckTx/
	// broadcast/RPC errors stay transient and are re-queued.
	if err := d.worker.TxHandler.SendCosmosTxBatch(ctx, d.cosmos, msgs); err != nil {
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
func (d *Destination) ClientExpiresAt(ctx context.Context, _ string) (time.Time, error) {
	readCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	cs, err := relayerclient.GetEthereumClientStateWithContext(readCtx, d.cosmos.CosmosClient(), d.clientID)
	if err != nil {
		return time.Time{}, fmt.Errorf("cosmos dest: eth client state: %w", err)
	}
	return ethClientExpiry(cs), nil
}

// ethClientExpiry derives the conservative expiry of the 08-wasm Ethereum light
// client from its on-chain state.
//
// A beacon client has no single trusting period the way a Tendermint client does:
// it must be advanced within about two sync-committee periods of its latest
// tracked slot. This reports ONE period after that slot's time, which is
// deliberately early -- the refresh routine then acts sooner than strictly
// necessary, and being wrong in the other direction lets the client lapse, which
// is the failure that costs the most.
//
// A zero period means the client state is misconfigured (any of the three factors
// unset). Reporting the zero time rather than an error is what the module reads as
// "no expiry", so it falls back to the periodic refresh instead of treating a
// malformed client state as an expiry emergency.
func ethClientExpiry(cs *relayerclient.EthereumClientState) time.Time {
	periodSecs := cs.EpochsPerSyncCommitteePeriod * cs.SlotsPerEpoch * cs.SecondsPerSlot
	if periodSecs == 0 {
		return time.Time{}
	}
	return time.Unix(int64(cs.ComputeTimestampAtSlot(cs.LatestSlot)+periodSecs), 0)
}
