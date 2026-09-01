# Packet Layer — `04-channel/v2`

Stock ibc-go v10 module, unmodified.

## Packet structure
> TODO: packet fields + payload fields (port, version, encoding, value); timeout is
> timestamp-based.

## Client-ID routing
> TODO: how a packet names its counterparty client; counterparty registration
> (`MsgRegisterCounterparty`); no handshake.

## Message flow
> TODO: `MsgSendPacket` / `MsgRecvPacket` / `MsgAcknowledgement` / `MsgTimeout` — one
> line each; which proof each carries and against what commitment path.

## Commitments
> TODO: what is stored under which key on the Cosmos side (packet commitment, receipt,
> ack).
