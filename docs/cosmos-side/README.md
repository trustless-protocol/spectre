# The Cosmos Side of Fast-IBC

> [Repository overview](../../README.md)

Fast-IBC is a trust-minimized bridge between a Cosmos chain and Ethereum. Each side runs a light client of the other and verifies every message against it, so no committee or multisig sits in the middle. This set covers only what runs **on the Cosmos chain**, for a reader new to IBC v2.

| # | Document | What it covers |
|---|----------|----------------|
| 1 | [IBC v2](01-ibc-v2.md) | IBC v2 and the layer model |
| 2 | [The Channel Layer (04-channel/v2)](02-channel-layer.md) | `04-channel/v2`: the packet, client-ID routing, commitments |
| 3 | [The Client Registry (02-client)](03-client-registry.md) | `02-client`: the light-client registry, IDs, and counterparty pairing |
| 4 | [The Wasm Host (08-wasm)](04-wasm-host.md) | `08-wasm`: light clients shipped as CosmWasm contracts |
| 5 | [The Ethereum Light Client (cw-ics08-wasm-eth)](05-eth-light-client.md) | The Ethereum light client Fast-IBC ships |
| 6 | [Token Transfer (ICS-20 v2)](06-token-transfer.md) | ICS-20 v2: denoms, escrow, mint, refunds |

On-chain Cosmos components only. The relayer and Ethereum-side Solidity are out of scope. Every claim links its source, and code is the source of truth. ibc-go links point at tag **v10.3.0**, the version the relayer pins.
