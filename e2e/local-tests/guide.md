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
./setup/01-eth-node.sh
```

This script:
- Starts a Kurtosis Ethereum testnet
- Deploys core Solidity contracts (ICS26, ICS20, ICS07, Verifiers)
- Configures `relayer/config.example.json` with deployed addresses

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
Copies `config.example.json` to `config.json` if needed.
Deploys the ICS07 Tendermint light client on Ethereum and copies the address back into `config.json`.

## Step 6: Start Relayer

```bash
./setup/05-relayer.sh
```

Auto-copies `config.example.json` to `config.json` if needed.
The relayer runs a bi-directional relay loop (Cosmos ↔ ETH).

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

Sends a transfer with a 5-second timeout, waits for expiry, then verifies the
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
