package transaction

import (
	contractICS26Router "operator/bindings/ICS26Router"
	"operator/services"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

type TxBuilder struct {
}

func (t *TxBuilder) RecvPacketTxBuilder(ctx services.Context, packet channeltypesv2.Packet) (contractICS26Router.IICS26RouterMsgsMsgRecvPacket, error) {
	payloads := make([]contractICS26Router.IICS26RouterMsgsPayload, len(packet.Payloads))
	for _, p := range packet.Payloads {
		payloads = append(payloads, contractICS26Router.IICS26RouterMsgsPayload{
			SourcePort: p.SourcePort,
			DestPort:   p.DestinationPort,
			Version:    p.Version,
			Encoding:   p.Encoding,
			Value:      p.Value,
		})
	}

	msgRecvPacket := contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet: contractICS26Router.IICS26RouterMsgsPacket{
			Sequence:         packet.Sequence,
			SourceClient:     packet.SourceClient,
			DestClient:       packet.DestinationClient,
			TimeoutTimestamp: packet.TimeoutTimestamp,
			Payloads:         payloads,
		},
	}

	return msgRecvPacket, nil
}
