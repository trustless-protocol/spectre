# Token Transfer (ICS-20 v2)

ICS-20 is the app most packets carry: moving fungible tokens between chains. Fast-IBC uses the stock ibc-go transfer app, pinned to v10.3.0 ([`relayer/go.mod`](https://github.com/decentrio/fast-ibc/blob/main/relayer/go.mod)), so there is no custom transfer logic to audit on the Cosmos side. It plugs into the packet layer as a v2 `IBCModule` with four callbacks, one per packet stage: `OnSendPacket`, `OnRecvPacket`, `OnTimeoutPacket`, `OnAcknowledgementPacket`, each on a single `Payload` ([`modules/apps/transfer/v2/ibc_module.go`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/v2/ibc_module.go)). Both ports must be `transfer` and both client IDs must be in `{string}-{number}` form. All ibc-go paths are tag v10.3.0.

## The payload

`Value` is decoded per its `Encoding`:

| `Encoding` string | Constant |
|-------------------|----------|
| `application/json` | `EncodingJSON` |
| `application/x-protobuf` | `EncodingProtobuf` |
| `application/x-solidity-abi` | `EncodingABI` |

Whatever the encoding, `UnmarshalPacketData` accepts only version `ics20-1` and yields a `FungibleTokenPacketData{Denom, Amount, Sender, Receiver, Memo}` (all strings), converted into an internal token with a structured `Denom{Base, Trace}`, each `Trace` hop a (port, channel) pair ([`modules/apps/transfer/types/packet.go`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/types/packet.go)). Fast-IBC uses `EncodingABI` both ways: the Ethereum `ICS20Transfer` contract stamps it and rejects any other on receive ([`contracts/ICS20Transfer.sol`](https://github.com/decentrio/fast-ibc/blob/main/contracts/ICS20Transfer.sol), constant in [`contracts/utils/ICS20Lib.sol`](https://github.com/decentrio/fast-ibc/blob/main/contracts/utils/ICS20Lib.sol)), and the relayer copies `Version`, `Encoding`, `Value` across verbatim.

## Sending (Cosmos → Ethereum)

`OnSendPacket` requires the payload `Sender` to equal the tx signer and rejects a base denom containing `/` (a v2 packet has no channel identifier to split trace from base). Then, by direction ([`modules/apps/transfer/keeper/relay.go`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/keeper/relay.go)):

- **Native or outbound token:** **escrowed** to `GetEscrowAddress(sourcePort, sourceChannel)`, the ADR-028 address `sha256("ics20-1" || 0x00 || "{portID}/{channelID}")` truncated to 20 bytes. The version is pinned to `ics20-1`, so a route's escrow address never moves.
- **Voucher returning home:** its denom already carries the outgoing (port, channel) prefix, so the coins are sent to the module account and **burned**.

## Receiving (Ethereum → Cosmos)

`OnRecvPacket` validates ports and ID format, decodes, and fails on any error. Two branches ([`modules/apps/transfer/keeper/relay.go`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/keeper/relay.go)):

- **Token coming home:** prefixed by the source (port, channel), so the first hop is stripped and the coins **unescrowed** from `GetEscrowAddress(destPort, destChannel)`.
- **Foreign token:** a hop `NewHop(destPort, destChannel)` is prepended, the denom recorded (`SetDenom` plus bank metadata), and a **voucher minted**, denominated `ibc/{sha256(trace + "/" + base)}` (a native denom with an empty trace keeps its `Base`).

## Timeouts and acks

`OnAcknowledgementPacket` does nothing on success and refunds on an error ack (the sentinel `ErrorAcknowledgement` routes through the shared refund path, and any other non-success ack is rejected). `OnTimeoutPacket` always refunds. A refund is the inverse of the send. Coins prefixed by the source (port, channel) were burned, so they are **re-minted**. Otherwise they were escrowed, so they are **unescrowed** from `GetEscrowAddress(sourcePort, sourceChannel)` ([`modules/apps/transfer/keeper/relay.go`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/keeper/relay.go)).

## Store

Transfer keeps little state of its own. The escrowed coins live in the **bank** module, held by per-route escrow accounts at `GetEscrowAddress(port, channel)`, not in a transfer key. What transfer itself stores:

| Key | Value | Notes |
| --- | ----- | ----- |
| [`DenomKey`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/types/keys.go#L51) (`0x03`) | `Denom{Base, Trace}` | the recorded denom for each voucher, set when a foreign token is received |
| [`TotalEscrowForDenomKey`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/types/keys.go#L75) (`totalEscrowForDenom/{denom}`) | `sdk.Coin` | running total escrowed per denom |
| [`ParamsKey`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/types/keys.go#L34) (`params`) | `Params` | `SendEnabled`, `ReceiveEnabled` |
| [`PortKey`](https://github.com/cosmos/ibc-go/blob/v10.3.0/modules/apps/transfer/types/keys.go#L47) (`0x01`) | `transfer` | the bound port ID |
