# Local E2E Test Guide

End-to-end test for fast-ibc on local Cosmos + Ethereum nodes.
Requires Docker, Kurtosis, Go >= 1.21, Foundry, Bun, and Just.

## Prerequisites

```bash
# Install dependencies (one-time setup)
# Kurtosis: https://docs.kurtosis.com/install
# Foundry: https://getfoundry.sh/
# Bun: https://bun.sh/
# Just: https://github.com/casey/just
```

## Step 0: Compile Prover Circuits (One-time)

```bash
cd relayer
LD_LIBRARY_PATH=$HOME/works/ecip-gnark go run ./prover/cmd ./bin ../contracts/verifiers
```

Re-run only when circuit code changes. After this, redeploy contracts.

## Step 1: Build Relayer Binary

```bash
cd relayer
go build -o relayer ./cmd
```

## Step 2: Start Ethereum Node

```bash
./scripts/01-run_eth_node.sh
```

This script:
- Starts a Kurtosis Ethereum testnet
- Deploys core Solidity contracts (ICS26, ICS20, ICS07, Verifiers)
- Configures `relayer/config.example.json` with deployed addresses

Wait until beacon finalizes (poll until `finalized.epoch > 0`):
```bash
# Auto-detect beacon port from Kurtosis
BEACON_PORT=$(kurtosis enclave inspect my-testnet | awk '/cl-1-lighthouse-geth/ {in_s=1} in_s && /http:/ {match($0, /127\.0\.0\.1:[0-9]+/); print substr($0, RSTART+10, RLENGTH-10); exit}')
curl -s http://127.0.0.1:${BEACON_PORT:-32774}/eth/v1/beacon/states/head/finality_checkpoints
```

## Step 3: Start Cosmos Node

```bash
./scripts/02-run_cosmos_node.sh
```

This script:
- Initializes 4 Gaia validators with test accounts
- Configures genesis with fast governance voting periods
- Exports `COSMOS_PRIVATE_KEY` to `relayer/.env`
- Starts all validators

## Step 4: Submit Ethereum Light Client WASM

```bash
./scripts/03-wasm.sh
```

Submits the Ethereum ICS08 WASM light client via Cosmos governance.
Waits for proposal to pass and outputs the WASM checksum.

## Step 5: Create IBC Clients

```bash
./scripts/04-create-clients.sh
```

Auto-detects WASM checksum from `03-wasm.sh` (saved to `.state/wasm_checksum`).
Copies `config.example.json` to `config.json` if needed.
Deploys the ICS07 Tendermint light client on Ethereum and copies the address back into `config.json`.

## Step 6: Start Relayer

```bash
./scripts/05-start-relayer.sh
```

Auto-copies `config.example.json` to `config.json` if needed.
The relayer runs a bi-directional relay loop (Cosmos ↔ ETH).

## Step 7: Send ETH → Cosmos Transfer

```bash
./scripts/06-eth_to_cosmos_transfer.sh
```

Auto-detects `RECEIVER` from Cosmos key (`test` by default).
Auto-loads `ERC20` and `ICS20Transfer` addresses from `broadcast/E2ETestDeploy.s.sol/*/run-latest.json`.
Default amount: `1000000000` wei. Override with `AMOUNT=<value> ./scripts/06-eth_to_cosmos_transfer.sh`.

## Step 8: Verify Results

Check Ethereum state:
```bash
./scripts/07-check_eth_data.sh
```

Check Cosmos state (voucher balance):
```bash
./scripts/08-check_cosmos_data.sh
```

## Script Ordering

| Order | Script | Purpose |
|-------|--------|---------|
| 0 | (prover/cmd) | Compile ZK circuits |
| 1 | `01-run_eth_node.sh` | Start Ethereum + deploy contracts |
| 2 | `02-run_cosmos_node.sh` | Start Cosmos validators |
| 3 | `03-wasm.sh` | Submit WASM via governance (saves checksum) |
| 4 | `04-create-clients.sh` | Deploy ICS07 light client (auto-reads checksum) |
| 5 | `05-start-relayer.sh` | Run relayer loop |
| 6 | `06-eth_to_cosmos_transfer.sh` | Send ICS20 transfer (auto-detects receiver) |
| 7 | `07-check_eth_data.sh` | Inspect Ethereum state |
| 8 | `08-check_cosmos_data.sh` | Inspect Cosmos state + voucher balance |

## State Files

The `.state/` directory stores outputs from previous steps:
- `.state/wasm_checksum` - Written by `03-wasm.sh`, read by `04-create-clients.sh`
