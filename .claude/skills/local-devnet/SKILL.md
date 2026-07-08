---
name: local-devnet
description: Bring up the full local ETH↔Cosmos devnet and relayer in the proven order (eth → beacon finality → cosmos → wasm gov → build relayer → circuits → create-clients → start), with wait conditions, failure modes, and recovery. Use when running E2E manually, testing the relayer against live local chains, or debugging devnet bring-up.
---

# fast-ibc local devnet runbook

The order below is load-bearing: the Tendermint light-client on ETH can only be created after Ethereum has **finalized** at least one epoch, and the ETH light-client on Cosmos needs the WASM code stored via governance. Doing these out of order produces confusing downstream failures, not clear errors at the wrong step.

Build setup (`bun install`, `just build-contracts`) is assumed done — it is separate from this runbook.

## 1. Ethereum node (Kurtosis) + contracts

```bash
./scripts/local/run_eth_node.sh
```

Deploys the contract suite. Note the ports Kurtosis mapped this run (they vary): the EL RPC port and the **beacon (CL) REST port** — you need the beacon port for the next step.

## 2. Wait for beacon finality — poll, never sleep blindly

```bash
watch -n 5 'curl -s http://127.0.0.1:<beacon-port>/eth/v1/beacon/states/head/finality_checkpoints | jq .data.finalized'
```

Proceed only when `finalized.epoch > 0`. On a fresh Kurtosis net this takes a few minutes. Starting the Cosmos side before this wastes the whole bring-up: `create-clients` will fail against a chain with no finalized checkpoint.

## 3. Cosmos node, then WASM light client via governance

```bash
./scripts/local/run_cosmos_node.sh     # local node with test accounts
./scripts/local/wasm.sh                # stores the ETH light-client WASM through a gov proposal
```

Docker variants exist (`run_cosmos_node_docker.sh`, `wasm_docker.sh`) if running the node natively is not possible. `wasm.sh` prints the stored **checksum** — record it; `create-clients` needs it.

Checksum facts (matter when targeting a real testnet, not just local):
- The 08-wasm checksum is the sha256 of the **uncompressed** wasm.
- `MsgStoreCode` is gov-gated; `MsgCreateClient` referencing an existing checksum is permissionless.
- Do NOT reuse checksums already stored on public testnets without verifying the ConsensusState schema: current schema is **5-field** (slot, state_root, timestamp, current_sync_committee, next_sync_committee — no `storage_root`). Stale 6-field builds fail at DeliverTx with `missing field storage_root`. Inspect a stored build: fetch `GET <rest>/ibc/lightclients/wasm/v1/checksums/{c}/code`, base64-decode, `strings | grep 'struct ConsensusState'`.

## 4. Build the relayer binary and circuit artifacts

```bash
cd relayer
go build -o relayer ./cmd                              # use the binary, not `go run`, for all ops
go run ./prover/cmd ./bin ../contracts/verifiers       # only if bin/n{N}/ artifacts are missing or R1CS changed
```

Circuit regeneration changes the vk → the generated `Groth16Verifier_N{N}.sol` must be **redeployed** and re-registered via `WrapperVerifier.setBucket(...)`. If contracts were just deployed by `run_eth_node.sh` from freshly generated verifiers, they already match — regeneration mid-session is what breaks the pairing.

Bucket sizing: the prover picks the smallest N ∈ {4,8,16,32,64,128} that covers a greedy-by-power ≥⅔ quorum. Local devnets and even the Cosmos Hub `provider` testnet (top-4 validators > ⅔) need only N=4 — don't compile big buckets you won't use.

## 5. Create clients

```bash
./relayer create-clients --config config.json --trust-level 2/3 --wasm-checksum <hex-from-step-3>
# or split: create-clients-cosmos first (needs --wasm-checksum), then create-clients-eth
```

Runs the Cosmos half first (08-wasm client → learns its client id), then the ETH half (ICS07 wired to that id), and writes `cosmos_wasm_client_id` + `ics07_client` back into `config.json` automatically — no manual copying. Config is JSON (`config.example.json` shape: `modules` array with `cosmos_to_eth` + `eth_to_cosmos`); `.env` holds only secrets/prover paths (`ETH_PRIVATE_KEY`, `COSMOS_PRIVATE_KEY`, `COSMOS_CHAIN_ID`, `PROVER_BIN_DIR` — the shipped values are devnet-only).

## 6. Start the relay loop

```bash
./relayer start --config config.json              # add --benchmark for per-inner-call gas + timing logs
```

Sanity signals in the log: both subscribers connected (`Successfully subscribed to ICS26Router events`), startup recovery scans complete, no repeating error line every second.

## 7. Verify with a transfer

Send a token transfer in one direction, watch it relay, then send the reverse direction. **Always test both directions** — the relayer is a mirror and one-sided success proves half the system. A packet from ETH relays only after its block is beacon-finalized (minutes on devnet) — "not relayed yet" right after send is normal, not a bug.

## Failure modes → cause

| Symptom | Cause / fix |
|---|---|
| `create-clients` fails creating Tendermint LC on ETH | ETH not finalized yet — redo step 2 |
| DeliverTx `missing field storage_root` | Stale 6-field 08-wasm build — store a current 5-field build (step 3) |
| Every `updateClient` proof reverts on-chain | Circuit artifacts regenerated after verifiers were deployed — redeploy verifiers + `setBucket` |
| `account sequence mismatch` on Cosmos | Something bypassed the serialized send path; check for a second process using the same key |
| Relayer idle though events fire | Wrong client-ID filter in config, or subscriber connected to the wrong port — check the mapped Kurtosis ports again |
| Beacon finality wait loops forever in relayer | Kurtosis net stalled — check CL container logs; restart the network rather than the relayer |

## Teardown

Kurtosis enclaves and the Cosmos node keep running until stopped: `kurtosis enclave ls` / `kurtosis enclave rm -f <name>`, and kill the cosmos node process. Stale enclaves are the usual cause of "port already in use" on the next bring-up.
