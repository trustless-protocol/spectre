# Local E2E Test Guide

End-to-end test for fast-ibc on local Cosmos + Ethereum nodes.
Requires Docker, Kurtosis, Go >= 1.21, Foundry, Bun, and Just.

## Directory Structure

```
e2e/local-tests/
├── lib/                    # Shared helper libraries
│   ├── common.sh           # ETH/Kurtosis/contract helpers
│   └── cosmos.sh           # Cosmos chain helpers
├── setup/                  # Infrastructure (shared across scenarios)
│   ├── 01-eth-node.sh      # Start Ethereum + deploy contracts
│   ├── 02-cosmos-node.sh   # Start Cosmos validators
│   ├── 03-wasm.sh          # Submit WASM via governance
│   ├── 04-create-clients.sh# Deploy light clients
│   └── 05-relayer.sh       # Run relayer loop
├── scenarios/              # Scenario-specific scripts
│   ├── eth-to-cosmos/
│   │   ├── transfer.sh     # Send ICS20 transfer (ETH→Cosmos)
│   │   ├── check-eth.sh    # Inspect Ethereum state
│   │   ├── check-cosmos.sh # Inspect Cosmos state
│   │   ├── run-success.sh  # Full success flow (transfer + verify)
│   │   └── run-timeout.sh  # Timeout test (short timeout, verify skip)
│   └── cosmos-to-eth/
│       └── README.md       # Placeholder for future scenarios
└── .state/                 # Inter-script state (wasm checksum, etc.)
```

**Backward compatibility:** Old paths in `scripts/` are wrappers that delegate
to the new locations. You can use either `scripts/<name>.sh` or the new paths.

## Configuration

The relayer uses two configuration files:

| File | Purpose |
|------|---------|
| `relayer/config.json` | Module addresses, endpoints, client IDs (auto-generated) |
| `relayer/.env` | Private keys, chain IDs, prover paths (written by scripts) |

### `config.json` fields (`cosmos_to_eth` module)

| Field | Description |
|-------|-------------|
| `tm_rpc_url` | Cosmos Tendermint RPC endpoint |
| `ics26_address` | ICS26 Router contract address on Ethereum |
| `ics26_client_id` | Client ID of the Cosmos light client on Ethereum's ICS26 Router |
| `cosmos_wasm_client_id` | Client ID of the Ethereum wasm light client on Cosmos |
| `eth_rpc_url` | Ethereum execution RPC URL |
| `eth_ws_url` | Ethereum WebSocket URL (for event subscriptions) |
| `ics07_client` | ICS07 Tendermint light client address on Ethereum |
| `wrapper_verifier` | WrapperVerifier contract address |
| `membership` | Membership contract address |
| `misbehaviour` | Misbehaviour contract address |
| `update_client` | UpdateClient contract address |

### `.env` file

| Variable | Set by | Used by |
|----------|--------|---------|
| `ETH_RPC_URL` | `01-eth-node.sh` | Relayer start, genesis, fixtures |
| `ETH_WS_URL` | `01-eth-node.sh` | Relayer start |
| `ETH_BEACON_API_URL` | `01-eth-node.sh` | Relayer create-clients, start |
| `ETH_PRIVATE_KEY` | `01-eth-node.sh` | Relayer start, create-clients (Ethereum tx signing) |
| `ERC20_ADDRESS` | `01-eth-node.sh` | Transfer/check scripts |
| `ICS20_ADDRESS` | `01-eth-node.sh` | Transfer/check scripts |
| `ICS26_ADDRESS` | `01-eth-node.sh` | Transfer/check scripts |
| `COSMOS_PRIVATE_KEY` | `02-cosmos-node.sh` | Relayer create-clients, start (Cosmos tx signing) |
| `COSMOS_CHAIN_ID` | `02-cosmos-node.sh` | Relayer create-clients, start |
| `PROVER_BIN_DIR` | `01-eth-node.sh` | Relayer create-clients, start (ZK circuit artifacts) |
| `WASM_CHECKSUM` | `04-create-clients.sh` | Relayer create-clients (Ethereum light client checksum) |
| `PROVER_BIN_DIR` | default: `./bin` | Path to per-bucket ZK artifacts (r1cs, pk, vk) |

### `LD_LIBRARY_PATH`

The gnark prover requires the ECIP shared library at runtime. Set via:
```bash
export LD_LIBRARY_PATH=$HOME/works/ecip-gnark${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}
```
Both `04-create-clients.sh` and `05-relayer.sh` set this automatically, defaulting to `$HOME/works/ecip-gnark`.

## Prerequisites

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
./setup/01-eth-node.sh
```

This script:
- Starts a Kurtosis Ethereum testnet
- Deploys core Solidity contracts (ICS26, ICS20, ICS07, Verifiers)
- Configures `relayer/config.json` with deployed addresses

Wait until beacon finalizes (poll until `finalized.epoch > 0`):
```bash
BEACON_PORT=$(kurtosis enclave inspect my-testnet | awk '/cl-1-lighthouse-geth/ {in_s=1} in_s && /http:/ {match($0, /127\.0\.0\.1:[0-9]+/); print substr($0, RSTART+10, RLENGTH-10); exit}')
curl -s http://127.0.0.1:${BEACON_PORT:-32774}/eth/v1/beacon/states/head/finality_checkpoints
```

## Step 3: Start Cosmos Node

```bash
./setup/02-cosmos-node.sh
```

This script:
- Initializes 4 Gaia validators with test accounts
- Configures genesis with fast governance voting periods
- Exports `COSMOS_PRIVATE_KEY` to `relayer/.env`
- Starts all validators

## Step 4: Submit Ethereum Light Client WASM

```bash
./setup/03-wasm.sh
```

Submits the Ethereum ICS08 WASM light client via Cosmos governance.
Waits for proposal to pass and outputs the WASM checksum.

## Step 5: Create IBC Clients

```bash
./setup/04-create-clients.sh
```

Auto-detects WASM checksum from `03-wasm.sh` (saved to `.state/wasm_checksum`).
Creates `config.json` from `config.example.json` if needed, refreshes Kurtosis endpoints, deploys the ICS07 Tendermint light client on Ethereum, and copies the address back into `config.json`.

## Step 6: Start Relayer

```bash
./setup/05-relayer.sh
```

Creates `config.json` from `config.example.json` if needed, refreshes Kurtosis endpoints, and starts
the relayer with `LD_LIBRARY_PATH` and `PROVER_BIN_DIR` configured. The relayer runs a bi-directional
relay loop (Cosmos ↔ ETH). Requires the `.env` file to have `ETH_PRIVATE_KEY` and `COSMOS_PRIVATE_KEY`.

## Scenarios

After infrastructure is up and the relayer is running, choose a scenario:

### ETH → Cosmos Success

```bash
./scenarios/eth-to-cosmos/run-success.sh
```

Or run steps individually:
```bash
./scenarios/eth-to-cosmos/transfer.sh      # Send ICS20 transfer
./scenarios/eth-to-cosmos/check-eth.sh     # Verify Ethereum state
./scenarios/eth-to-cosmos/check-cosmos.sh  # Verify Cosmos voucher
```

### ETH → Cosmos Timeout

```bash
./scenarios/eth-to-cosmos/run-timeout.sh
```

Sends a transfer with a 10-second timeout, waits for expiry, then verifies the
packet was NOT relayed to Cosmos. The relayer is stopped during this test so
the packet cannot be relayed before the timeout.

The script also documents the steps needed for the full on-chain timeout refund
(non-membership proof from Cosmos + `timeoutPacket()` on Ethereum).

## Script Ordering

| Order | Script | Purpose |
|-------|--------|---------|
| 0 | (prover/cmd) | Compile ZK circuits |
| 1 | `setup/01-eth-node.sh` | Start Ethereum + deploy contracts |
| 2 | `setup/02-cosmos-node.sh` | Start Cosmos validators |
| 3 | `setup/03-wasm.sh` | Submit WASM via governance (saves checksum) |
| 4 | `setup/04-create-clients.sh` | Deploy ICS07 light client (auto-reads checksum) |
| 5 | `setup/05-relayer.sh` | Run relayer loop |
| 6 | `scenarios/eth-to-cosmos/transfer.sh` | Send ICS20 transfer (auto-detects receiver) |
| 7 | `scenarios/eth-to-cosmos/check-eth.sh` | Inspect Ethereum state |
| 8 | `scenarios/eth-to-cosmos/check-cosmos.sh` | Inspect Cosmos state + voucher balance |
| — | `scenarios/eth-to-cosmos/run-timeout.sh` | Timeout test (short timeout + verify skip) |

## State Files

The `.state/` directory stores outputs from previous steps:
- `.state/wasm_checksum` - Written by `03-wasm.sh`, read by `04-create-clients.sh`
