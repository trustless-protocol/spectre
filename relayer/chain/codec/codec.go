// Package codec holds the gob codecs for the opaque chain.ClientUpdate.Payload of
// each adapter pair. It is a neutral leaf package so the source-side builder and
// the destination-side adapter can share an encoding without importing each
// other (which would cycle: evm decodes cosmos updates, cosmos decodes beacon
// updates).
package codec

import (
	"bytes"
	"encoding/gob"
	"fmt"

	spectreContract "relayer/bindings/SpectreClient"
	updateclientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

// --- Cosmos -> ETH (groth16) update ---

type cosmosUpdate struct {
	Kind      int // services.ClientUpdateKind, kept as int so codec stays leaf
	AppMsg    updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState
	NewValSet spectreContract.IICS07TendermintMsgsValidatorSet
}

// EncodeCosmosUpdate encodes a Cosmos->ETH SpectreClient update. kind carries the
// services.ClientUpdateKind value.
func EncodeCosmosUpdate(kind int, appMsg updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState, newValSet spectreContract.IICS07TendermintMsgsValidatorSet) ([]byte, error) {
	return encode(cosmosUpdate{Kind: kind, AppMsg: appMsg, NewValSet: newValSet})
}

// DecodeCosmosUpdate reverses EncodeCosmosUpdate.
func DecodeCosmosUpdate(payload []byte) (kind int, appMsg updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState, newValSet spectreContract.IICS07TendermintMsgsValidatorSet, err error) {
	var u cosmosUpdate
	if derr := decode(payload, &u); derr != nil {
		return 0, appMsg, newValSet, fmt.Errorf("codec: decode cosmos update: %w", derr)
	}
	return u.Kind, u.AppMsg, u.NewValSet, nil
}

// --- ETH -> Cosmos (beacon) update ---

type beaconUpdate struct {
	MsgBytes       [][]byte // each = proto.Marshal(*clienttypes.MsgUpdateClient)
	EthClientState relayerclient.EthereumClientState
	ProofTimestamp uint64
	SigSlot        uint64
}

// EncodeBeaconUpdate encodes an ETH->Cosmos beacon update. msgs must be
// *clienttypes.MsgUpdateClient (proto-marshaled — the nested Any does not gob
// cleanly).
func EncodeBeaconUpdate(msgs []any, clientState relayerclient.EthereumClientState, proofTimestamp, sigSlot uint64) ([]byte, error) {
	msgBytes := make([][]byte, 0, len(msgs))
	for i, m := range msgs {
		msg, ok := m.(*clienttypes.MsgUpdateClient)
		if !ok {
			return nil, fmt.Errorf("codec: unexpected msg type %T at %d", m, i)
		}
		raw, err := msg.Marshal()
		if err != nil {
			return nil, fmt.Errorf("codec: marshal MsgUpdateClient %d: %w", i, err)
		}
		msgBytes = append(msgBytes, raw)
	}
	return encode(beaconUpdate{MsgBytes: msgBytes, EthClientState: clientState, ProofTimestamp: proofTimestamp, SigSlot: sigSlot})
}

// DecodeBeaconUpdate reverses EncodeBeaconUpdate, returning the messages as []any
// ready for SendCosmosTxBatch plus the state the pre-submit catch-up needs.
func DecodeBeaconUpdate(payload []byte) (msgs []any, clientState relayerclient.EthereumClientState, sigSlot uint64, err error) {
	var u beaconUpdate
	if derr := decode(payload, &u); derr != nil {
		return nil, clientState, 0, fmt.Errorf("codec: decode beacon update: %w", derr)
	}
	msgs = make([]any, 0, len(u.MsgBytes))
	for i, raw := range u.MsgBytes {
		var msg clienttypes.MsgUpdateClient
		if uerr := msg.Unmarshal(raw); uerr != nil {
			return nil, clientState, 0, fmt.Errorf("codec: unmarshal MsgUpdateClient %d: %w", i, uerr)
		}
		msgs = append(msgs, &msg)
	}
	return msgs, u.EthClientState, u.SigSlot, nil
}

func encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(payload []byte, v any) error {
	return gob.NewDecoder(bytes.NewReader(payload)).Decode(v)
}
