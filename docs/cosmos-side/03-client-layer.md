# Client Layer — `02-client` + `08-wasm`

Stock ibc-go v10 modules, unmodified. Fast-IBC supplies the contracts *inside* 08-wasm.

## 02-client
> TODO: client registry, client IDs, client/consensus state storage — 3–4 sentences.

## 08-wasm host
> TODO: ICS-08 — a CosmWasm contract behind the standard light-client interface;
> checksum-addressed code, uploaded by governance; entry points (`instantiate`,
> `sudo`: update/membership, `query`: verify/status).

## Clients hosted in this project
> TODO: table — Ethereum client (`programs/cw-ics08-wasm-eth`, production path, doc 04)
> and L2 clients (`programs/cw-ics08-wasm-{op,base,arbitrum}`, **devnet only**, see
> `docs/L2_CLIENTS.md`).
