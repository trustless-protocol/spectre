package services

import (
	"log"

	contractICS26Router "relayer/bindings/ICS26Router"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// cosmosBatchMsg is one ETH-bound inner call planned from a Cosmos batch.
type cosmosBatchMsg struct {
	msg          any
	label        string        // log label, e.g. "RecvPacket"
	sequence     uint64        // for log lines
	sourceClient string        // for PendingTracker.Remove
	isRecv       bool          // true → call PendingTracker.Remove after success
	origin       *CosmosPacket // source packet to re-queue on failure; nil for folded updateClient
}

// ethBatchMsg is one Cosmos-bound inner msg planned from an ETH batch.
type ethBatchMsg struct {
	msg      any
	label    string // log label, e.g. "EthSend", "EthWriteAck", "UpdateClient"
	sequence uint64
	origin   *EthPacket // source packet to re-queue on failure; nil for folded updateClient
}

// relayablePacket is an ETH packet that survived the pre-filter and is ready for
// proof construction.
type relayablePacket struct {
	packet EthPacket
	signer string
}

// planCosmosPacketMsgs decides, per packet, the ETH-bound message to build and how to
// route failures. The proof builders are injected so every per-packet failure branch
// is unit-testable without a live Context/RPC (issue #106).
//
// Returns:
//   - msgs:              successfully-built ETH-bound inner calls.
//   - trackerAdds:       CosmosSend packets to record in the pending tracker (timeout
//     fallback) — every CosmosSend, regardless of proof outcome.
//   - transientFailures: packets whose proof build failed; the caller re-queues them
//     (transient) instead of dropping them.
//   - trackerRemoves:    CosmosSend packets whose proof build failed and must be dropped
//     from the pending tracker (they are now back on the retry queue; they are re-added
//     on the next attempt, Add being idempotent).
func planCosmosPacketMsgs(
	packets []CosmosPacket,
	ethBlockTime uint64,
	routerClientID string,
	buildMembership func(packet channeltypesv2.Packet, clientID string, pathType []byte) ([]byte, error),
	buildNonMembership func(packet channeltypesv2.Packet, clientID string, pathType []byte) ([]byte, error),
) (msgs []cosmosBatchMsg, trackerAdds, transientFailures, trackerRemoves []CosmosPacket) {
	for _, packet := range packets {
		switch packet.Type {
		case CosmosSend:
			// Track every CosmosSend so the async timeout scanner can refund it on
			// the source chain if ETH delivery never completes.
			trackerAdds = append(trackerAdds, packet)

			if ethBlockTime > 0 && packet.Packet.TimeoutTimestamp > 0 && ethBlockTime >= packet.Packet.TimeoutTimestamp {
				log.Printf("[RecvPacket] Packet seq=%d timed out (timeout=%d <= eth_block_time=%d), deferring to async timeout scanner",
					packet.Packet.Sequence, packet.Packet.TimeoutTimestamp, ethBlockTime)
				continue
			}

			calldata, err := buildMembership(*packet.Packet, packet.Packet.SourceClient, []byte{1})
			if err != nil {
				log.Printf("[RecvPacket] seq=%d: %v", packet.Packet.Sequence, err)
				// Transient proof failure: re-queue for retry instead of dropping, and
				// drop the tracker entry added above so the timeout scanner does not act
				// on a packet that is back on the relay queue (issue #106).
				transientFailures = append(transientFailures, packet)
				trackerRemoves = append(trackerRemoves, packet)
				continue
			}

			msgs = append(msgs, cosmosBatchMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
					Packet:        toEthPacket(*packet.Packet),
					MembershipMsg: calldata,
				},
				label:        "RecvPacket",
				sequence:     packet.Packet.Sequence,
				sourceClient: packet.Packet.SourceClient,
				isRecv:       true,
				origin:       &packet,
			})
		case CosmosAck:
			if len(packet.AckBytes) == 0 {
				log.Printf("[AckPacket] seq=%d: acknowledgement bytes missing, skipping", packet.Packet.Sequence)
				continue
			}

			calldata, err := buildMembership(*packet.Packet, packet.Packet.DestinationClient, []byte{3})
			if err != nil {
				log.Printf("[AckPacket] seq=%d: %v", packet.Packet.Sequence, err)
				transientFailures = append(transientFailures, packet)
				continue
			}

			msgs = append(msgs, cosmosBatchMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgAckPacket{
					Packet:          toEthPacket(*packet.Packet),
					Acknowledgement: packet.AckBytes[0],
					MembershipMsg:   calldata,
				},
				label:    "AckPacket",
				sequence: packet.Packet.Sequence,
				origin:   &packet,
			})
		case CosmosTimeout:
			if !shouldRelayCosmosTimeoutToEth(packet.Packet, routerClientID) {
				log.Printf("[Timeout] seq=%d: Cosmos-originated packet timeout already handled locally, skipping ETH relay", packet.Packet.Sequence)
				continue
			}

			calldata, err := buildNonMembership(*packet.Packet, packet.Packet.DestinationClient, []byte{2})
			if err != nil {
				log.Printf("[Timeout] seq=%d: %v", packet.Packet.Sequence, err)
				transientFailures = append(transientFailures, packet)
				continue
			}

			msgs = append(msgs, cosmosBatchMsg{
				msg: contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
					Packet:           toEthPacket(*packet.Packet),
					NonMembershipMsg: calldata,
				},
				label:    "Timeout",
				sequence: packet.Packet.Sequence,
				origin:   &packet,
			})
		default:
			log.Printf("[StartLoop] Unknown cosmos packet type: %d (seq=%d)", packet.Type, packet.Packet.Sequence)
		}
	}
	return msgs, trackerAdds, transientFailures, trackerRemoves
}

// planEthPacketMsgs decides, per relayable packet, the Cosmos-bound msg to build and how
// to route failures. buildProof is injected for testability (issue #106).
//
// Returns:
//   - msgs:              successfully-built Cosmos-bound inner msgs.
//   - expired:           EthSend packets past their timeout; the caller routes them to
//     the local timeout flow (timeoutEthSend) — kept out of this pure fn as a side effect.
//   - transientFailures: packets whose proof build failed; the caller re-queues them.
func planEthPacketMsgs(
	relayable []relayablePacket,
	proofSlot uint64,
	signerAddr string,
	buildProof func(path []byte) ([]byte, error),
) (msgs []ethBatchMsg, expired, transientFailures []EthPacket) {
	for _, r := range relayable {
		packet := r.packet
		origin := packet
		switch packet.Type {
		case EthSend:
			if ethPacketExpired(packet) {
				expired = append(expired, packet)
				continue
			}
			proofBytes, err := buildProof(ethPath(packet.Packet.SourceClient, packet.Packet.Sequence, 1))
			if err != nil {
				log.Printf("[EthSend] seq=%d: failed to get ETH membership proof: %v", packet.Packet.Sequence, err)
				transientFailures = append(transientFailures, packet)
				continue
			}
			msgs = append(msgs, ethBatchMsg{
				msg: &channeltypesv2.MsgRecvPacket{
					Packet:          *packet.Packet,
					ProofCommitment: proofBytes,
					ProofHeight:     clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
					Signer:          signerAddr,
				},
				label:    "EthSend",
				sequence: packet.Packet.Sequence,
				origin:   &origin,
			})
		case EthWriteAck:
			proofBytes, err := buildProof(ethPath(packet.Packet.DestinationClient, packet.Packet.Sequence, 3))
			if err != nil {
				log.Printf("[EthWriteAck] seq=%d: failed to get ETH membership proof: %v", packet.Packet.Sequence, err)
				transientFailures = append(transientFailures, packet)
				continue
			}
			msgs = append(msgs, ethBatchMsg{
				msg: &channeltypesv2.MsgAcknowledgement{
					Packet: *packet.Packet,
					Acknowledgement: channeltypesv2.Acknowledgement{
						AppAcknowledgements: packet.AckBytes,
					},
					ProofAcked:  proofBytes,
					ProofHeight: clienttypes.Height{RevisionNumber: 0, RevisionHeight: proofSlot},
					Signer:      signerAddr,
				},
				label:    "EthWriteAck",
				sequence: packet.Packet.Sequence,
				origin:   &origin,
			})
		}
	}
	return msgs, expired, transientFailures
}
