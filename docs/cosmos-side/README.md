# Fast-IBC — Cosmos-Side Specs

What runs **on the Cosmos chain** for Fast-IBC. For readers who don't know IBC v2.
Scope: on-chain components only. Status: scaffold — `TODO` sections unwritten.

| # | Document | Covers |
|---|----------|--------|
| 1 | [01-ibc-v2.md](01-ibc-v2.md) | IBC v2 vs v1 in brief; the layer model |
| 2 | [02-packet-layer.md](02-packet-layer.md) | `04-channel/v2`: packets, client-ID routing |
| 3 | [03-client-layer.md](03-client-layer.md) | `02-client` + `08-wasm` client host |
| 4 | [04-eth-light-client.md](04-eth-light-client.md) | The Ethereum light client contract |
| 5 | [05-token-transfer.md](05-token-transfer.md) | ICS-20 v2: denoms, escrow, mint |

Rules: English; terminology matches the code; every claim cites `file:line`; code wins
over docs.
