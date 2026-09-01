# IBC v2 in Brief

## What IBC does
> TODO: 1 paragraph — trust-minimized messaging; each chain runs a light client of the
> other.

## v2 vs v1
> TODO: short table — v1: connections + channels + 4-step handshakes; v2: none of
> that, packets route by **client ID**, one packet carries app **payloads**. The
> client layer (`02-client`, `08-wasm`) is unchanged from v1 — v2 replaced only the
> packet layer.

## The layer model
> TODO: small diagram — apps → `04-channel/v2` → `02-client` → light clients.

## Fast-IBC is v2-only
> TODO: 2 sentences; evidence: only `04-channel/v2/types` imports, v2 packet msgs
> (`relayer/transaction/handler.go`) — re-verify lines when writing.
