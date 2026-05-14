# Local E2E Test Guide

This directory contains bash-driven end-to-end tests for the fast-ibc IBC v2 implementation. The tests run against local Ethereum (Kurtosis) and Cosmos (multi-validator `gaiad`) nodes, orchestrated by the Go relayer.

> **Architecture context**: See `docs/ARCHITECTURE.md` for the system diagram, contract hierarchy, and request flows. This guide focuses on the local test scripts only.

---

## Prerequisites

| Tool | Purpose | Version hint |
|------|---------|-------------|
| `bash` + coreutils | Script runtime | GNU bash 4+ |
| `foundry` (`forge`, `cast`) | Solidity deployment & ETH interaction | Latest |
| `jq` | JSON parsing everywhere | 1.6+ |
| `bc` | Big-integer arithmetic (bash overflows > 2^63) | Any |
| `kurtosis` | Local Ethereum testnet | 0.88+ |
| `gaiad` | Cosmos node binary | Cosmos SDK v0.50+ |
| `go` | Build relayer & send-packet helper | 1.21+ |
| `sha256sum` | IBC denom derivation | Any |
| Docker Desktop | Kurtosis containers | Latest |

> **Note**: E2E tests require the relayer binary and Groth16 circuit artifacts (`relayer/bin/n{N}/{r1cs,pk,vk}.bin`). See `docs/ARCHITECTURE.md` § "Go Relayer Structure" and `relayer/prover/cmd/` for artifact generation.

---

## Directory Structure

```
e2e/local-tests/
├── bin/
│   └── send-packet              # Go helper: submits Cosmos MsgSendPacket
├── lib/
│   ├── common.sh                # Logging, Kurtosis discovery, contract loading, env helpers
│   ├── cosmos.sh                # Cosmos queries, key helpers, balance polling, block waiting
│   ├── eth.sh                   # ETH tx helpers: approve, sendTransfer, multicall, balanceOf
│   └── relayer.sh               # Relayer process health, log inspection
├── setup/
│   ├── 00-cleanup.sh            # Tear down all local infrastructure
│   ├── 01-eth-node.sh           # Start Kurtosis ETH + deploy contracts
│   ├── 02-cosmos-node.sh        # Start 4-validator gaiad localnet
│   ├── 03-wasm.sh               # Submit Ethereum light client WASM via governance
│   ├── 04-create-clients.sh   # One-time light client setup (relayer create-clients)
│   └── 05-relayer.sh            # Start bi-directional relay loop
└── scenarios/
    ├── eth-to-cosmos/           # ETH -> Cosmos IBC transfer scenarios
    ├── cosmos-to-eth/           # Cosmos -> ETH IBC transfer scenarios
    └── light-client/            # Light client verification scenarios
```

---

## Quick Start (Full Flow)

Run the entire stack from scratch in sequence:

```bash
cd e2e/local-tests

# 1. Tear down any previous state
./setup/00-cleanup.sh

# 2. Start Ethereum (Kurtosis) and deploy contracts
./setup/01-eth-node.sh

# 3. Start Cosmos 4-validator localnet
./setup/02-cosmos-node.sh

# 4. Store Ethereum light client WASM code on Cosmos (governance proposal)
./setup/03-wasm.sh

# 5. Create IBC light clients on both chains
./setup/04-create-clients.sh

# 6. Start the bi-directional relayer
./setup/05-relayer.sh
```

After setup, run scenarios in a **new terminal** (the relayer stays foreground in step 6 unless backgrounded).

---

## Setup Scripts Reference

### `setup/00-cleanup.sh`
Kills relayer, gaiad processes, removes Kurtosis enclave, and wipes `e2e/local-tests/.state`.

### `setup/01-eth-node.sh`
- Removes stale Kurtosis enclave `my-testnet`
- Spins up Ethereum via `kurtosis run ... ethereum-package@6.1.0`
- Deploys contracts via `forge script scripts/E2ETestDeploy.s.sol`
- Extracts deployed addresses (ERC20, ICS20Transfer, ICS26Router, WrapperVerifier, Membership, UpdateClient, Misbehaviour)
- Writes endpoints + addresses to `relayer/.env`
- Updates `relayer/config.json` with ETH RPC / WS / Beacon / contract addresses
- Polls until `finalized` block > 1 (REQUIRED BEFORE RUNNING OTHER SCRIPTS)

### `setup/02-cosmos-node.sh`
- Initializes 4 `gaiad` validators with unique home dirs and ports
- Configures genesis: short gov voting periods, fee market params, funded accounts
- Sets up persistent P2P peer connections
- Exports `test1` private key to `relayer/.env` as `COSMOS_PRIVATE_KEY`
- Starts all 4 validators in background (logs to `<home>/gaiad.log`)
- **Runs foreground** — open a new terminal for subsequent steps

### `setup/03-wasm.sh`
- Submits an `MsgStoreCode` governance proposal for the Ethereum light client WASM
- Votes `yes` from validators
- Waits for proposal passage (`PROPOSAL_STATUS_PASSED`)
- Saves the WASM checksum hex to `e2e/local-tests/.state/wasm_checksum`

### `setup/04-create-clients.sh`
- Reads WASM checksum from `.state/wasm_checksum`
- Builds relayer binary if stale
- Runs `./relayer create-clients --config config.json --wasm-checksum <hex>`
- This establishes the Tendermint light client on Ethereum and the Ethereum light client on Cosmos

### `setup/05-relayer.sh`
- Refreshes ETH endpoints in `config.json`
- Builds relayer binary if needed
- Starts `./relayer start --config config.json` in background
- Writes PID to `relayer/relayer.pid`
- **Runs foreground** with a cleanup trap — background it with `&` or open a new terminal

---

## Scenario Scripts

All scenarios assume the full stack from [Quick Start](#quick-start-full-flow) is running. They share a common pattern:

1. Source `lib/{common,cosmos,eth,relayer}.sh`
2. Load contract addresses from `relayer/config.json`, `.env`, or Foundry broadcast JSON
3. Auto-discover Kurtosis endpoints if available
4. Record pre-transfer balances
5. Submit transfer(s)
6. Poll for expected state change
7. Print pass/fail summary

> **Environment overrides**: Every scenario respects env vars for `ETH_RPC_URL`, `ETH_PRIVATE_KEY`, `AMOUNT`, `TIMEOUT_SECONDS`, `COSMOS_RPC_URL`, etc. See individual scripts for the full list.

---

### ETH -> Cosmos Scenarios (`scenarios/eth-to-cosmos/`)

| Script | Description | What it verifies |
|--------|-------------|----------------|
| `run-success.sh` | Single ERC20 transfer ETH->Cosmos | Tokens escrowed on ETH, voucher minted on Cosmos |
| `run-success-roundtrip.sh` | ETH->Cosmos then Cosmos->ETH return | Full roundtrip: ERC20 locked, voucher created, voucher burned, escrow released |
| `run-timeout.sh` | Single ETH->Cosmos with short timeout | Packet expires, relayer submits `timeoutPacket`, sender refunded |
| `run-timeout-batch.sh` | Multiple ETH->Cosmos packets timeout | N packets expire, all refunded via multicall timeout |
| `run-error-ack.sh` | Invalid receiver causes error ack | Cosmos rejects packet, error ack relayed back, ERC20 refunded |
| `run-batch-25.sh` | 25 concurrent transfers via `multicall` | All 25 packets relayed, total voucher balance correct |
| `run-batch-50.sh` | 50 concurrent transfers (wraps `run-batch-25.sh`) | Same as 25 but with `BATCH_SIZE=50` |
| `run-uint256-max.sh` | Large amount edge case | Transfers `StartingERC20Balance / 2`, verifies no overflow |
| `run-voucher-timeout-eth-remint.sh` | ibcERC20 (Cosmos native) timeout remint | Native coin -> ETH creates ibcERC20; ibcERC20 -> Cosmos times out; ibcERC20 reminted on ETH |
| `transfer.sh` | Raw helper: approve + `sendTransfer` | Low-level helper, checks ETH/Cosmos state |
| `check-eth.sh` | Inspect ETH-side state | ERC20, ICS20, Escrow, allowances, tx details |
| `check-cosmos.sh` | Inspect Cosmos-side state | Voucher balances, IBC clients, tx details |

**Usage examples:**

```bash
# Basic success case
./scenarios/eth-to-cosmos/run-success.sh

# Full roundtrip (most comprehensive)
./scenarios/eth-to-cosmos/run-success-roundtrip.sh

# Timeout with custom amount and 10-second expiry
AMOUNT=500000000 TIMEOUT_SECONDS=10 ./scenarios/eth-to-cosmos/run-timeout.sh

# 25-packet batch
./scenarios/eth-to-cosmos/run-batch-25.sh

# Check current state without sending anything
./scenarios/eth-to-cosmos/check-eth.sh
./scenarios/eth-to-cosmos/check-cosmos.sh
```

---

### Cosmos -> ETH Scenarios (`scenarios/cosmos-to-eth/`)

| Script | Description | What it verifies |
|--------|-------------|----------------|
| `run-native-coin.sh` | Native Cosmos coin -> ETH (minting ibcERC20), then return | ibcERC20 created with correct metadata; return transfer burns ibcERC20 and restores native balance |
| `run-timeout.sh` | Single Cosmos->ETH with short timeout | Packet expires, relayer submits timeout on Cosmos, native balance refunded |
| `run-timeout-batch.sh` | Multiple Cosmos->ETH packets timeout | N packets expire, all native balance refunded |
| `run-voucher-timeout-cosmos-remint.sh` | Voucher (originally ERC20) timeout remint | ERC20 -> Cosmos creates voucher; voucher -> ETH times out; voucher reminted on Cosmos |
| `transfer.sh` | Raw helper: submit `MsgSendPacket` from Cosmos | Wraps `bin/send-packet` Go helper |
| `check-eth.sh` | Inspect ETH-side ibcERC20 state | ibcERC20 addresses, names, symbols, balances |
| `check-cosmos.sh` | Inspect Cosmos-side native + voucher state | Native balances, voucher balances, IBC clients |

**Usage examples:**

```bash
# Native coin roundtrip (most comprehensive Cosmos->ETH test)
./scenarios/cosmos-to-eth/run-native-coin.sh

# Timeout with custom parameters
AMOUNT=1000000000 TIMEOUT_SECONDS=45 ./scenarios/cosmos-to-eth/run-timeout.sh

# Batch timeout
BATCH_SIZE=5 ./scenarios/cosmos-to-eth/run-timeout-batch.sh

# Send raw Cosmos->ETH transfer
DENOM=stake AMOUNT=1000000000 ./scenarios/cosmos-to-eth/transfer.sh
```

---

### Light Client Scenarios (`scenarios/light-client/`)

| Script | Description | Priority |
|--------|-------------|----------|
| `run-update-client.sh` | Verify auto-update as Cosmos advances | Polls `Groth16ICS07Tendermint.clientState()` until height increases |
| `run-membership-proof.sh` | Skeleton for on-chain membership verification | Documents the `relayer fixtures membership` command flow |
| `run-misbehaviour.sh` | Misbehaviour detection / client freeze | Documents approach; recommends Go tests for actual automation |

**Usage:**

```bash
# Verify light client auto-updates
./scenarios/light-client/run-update-client.sh

# Check membership proof capability (informational)
./scenarios/light-client/run-membership-proof.sh
```

> **Note**: `run-misbehaviour.sh` is informational only. Generating conflicting CometBFT headers requires a custom Go test harness. For automated misbehaviour testing, run `cd e2e/interchaintestv8 && go test -run Test_DoubleSignMisbehaviour`.

---

## Shared Libraries

### `lib/common.sh`
Core utilities used by every script:

| Function | Purpose |
|----------|---------|
| `log`, `log_ok`, `log_warn`, `log_err`, `log_kv`, `log_header` | Structured colored stdout/stderr output |
| `require_cmd <cmd>` | Exit if command missing |
| `discover_kurtosis_endpoints` | Auto-populate `ETH_RPC_URL`, `ETH_WS_URL`, `ETH_BEACON_API_URL` from Kurtosis inspect |
| `load_contract_addresses` | Read ERC20/ICS20/ICS26 from Foundry broadcast JSON |
| `load_cosmos_receiver` | Derive `RECEIVER` bech32 from `gaiad keys show` |
| `safe_cast_call ...` | `cast call` that returns empty on failure |
| `contract_has_code <addr>` | Check if address has bytecode |
| `load_env_file` | Source `relayer/.env` |
| `ensure_relayer_config <file> [template]` | Copy example config if missing |
| `backfill_relayer_config_fields <file>` | Add newer config keys to old configs |
| `refresh_relayer_eth_endpoints <file>` | Update config with live Kurtosis endpoints |
| `build_relayer_if_needed [dir]` | Rebuild relayer binary when `.go` files are newer |
| `uint_value` | Strip uint256 hex formatting |
| `bc_sub1 <n>`, `bc_ge <a> <b>` | Big-int subtraction and comparison |
| `poll_eth_balance <token> <addr> <target> [timeout] [interval]` | Poll ETH ERC20 balance until >= target |
| `get_escrow_address <ics20> <client_id>` | Query ICS20Transfer escrow address |
| `derive_ibc_denom <trace>` | Compute `ibc/<sha256(trace)>` denom |
| `require_relayer_running [pid_file]` | Verify relayer process is alive |
| `update_relayer_deploy_config <file> ...` | Bulk-update config with deployed addresses |

### `lib/cosmos.sh`
Cosmos-specific helpers:

| Function | Purpose |
|----------|---------|
| `cosmos_query <subcommand> ...` | Wrap `gaiad query` with node/chain-id/output json |
| `key_address <key>` | Get bech32 address from keyring |
| `upsert_env_var <file> <key> <value>` | Atomic append-or-replace in `.env` file |
| `poll_cosmos_balance <addr> <denom> <target> [timeout] [interval]` | Poll bank balance until matches target |
| `wait_for_cosmos_block <n>` | Wait until chain advances by N blocks |

### `lib/eth.sh`
Ethereum transaction builders:

| Function | Purpose |
|----------|---------|
| `eth_approve <token> <spender> <amount>` | ERC20 `approve` |
| `eth_send_transfer <ics20> <erc20> <amount> <receiver> <src_client> <dest_port> <timeout> [memo]` | Submit `ICS20Transfer.sendTransfer` |
| `eth_send_transfer_with_receipt ...` | Same but returns JSON receipt |
| `eth_call_multicall <ics20> <calldatas...>` | Submit `ICS20Transfer.multicall(bytes[])` |
| `get_escrow_address <ics20> <client_id>` | Read escrow mapping |
| `get_ibc_erc20_address <ics20> <denom_path>` | Read ibcERC20 contract mapping |
| `get_packet_commitment <ics26> <client_id> <sequence>` | Read packet commitment |
| `eth_balance_of <token> <addr>` | ERC20 `balanceOf` |
| `eth_wallet_address` | Derive address from `ETH_PRIVATE_KEY` |

### `lib/relayer.sh`
Relayer process helpers:

| Function | Purpose |
|----------|---------|
| `relayer_is_running` | Check PID file and process existence |
| `relayer_ensure_running` | Exit with helpful message if relayer down |
| `relayer_wait_for_log <pattern> [timeout]` | Poll log file for pattern |
| `relayer_tail_log [lines]` | Tail relayer log |

---

## Configuration & Environment

### Primary config files
- **`relayer/.env`** — Secrets and endpoints: `ETH_PRIVATE_KEY`, `COSMOS_PRIVATE_KEY`, `ETH_RPC_URL`, `ETH_WS_URL`, `ETH_BEACON_API_URL`, `WASM_CHECKSUM`, `PROVER_BIN_DIR`
- **`relayer/config.json`** — Relayer JSON config with `modules` array (`cosmos_to_eth`, `eth_to_cosmos`). Auto-generated from `config.example.json` by setup scripts.

### State files
- **`e2e/local-tests/.state/wasm_checksum`** — Hex checksum from `03-wasm.sh`, consumed by `04-create-clients.sh`
- **`relayer/relayer.pid`** — Background relayer process ID
- **`relayer/relayer.log`** — Relayer stdout/stderr (when started by `05-relayer.sh`)

### Key environment variables
All scenarios honor these overrides:

| Variable | Default | Used by |
|----------|---------|---------|
| `ETH_RPC_URL` | Kurtosis-discovered or config | All ETH scenarios |
| `ETH_PRIVATE_KEY` | `relayer/.env` | All ETH tx scenarios |
| `ETH_BEACON_API_URL` | Kurtosis-discovered | Setup, relayer |
| `COSMOS_RPC_URL` | `http://127.0.0.1:26657` | All Cosmos scenarios |
| `COSMOS_CHAIN_ID` | `test-ibc-eth` | All Cosmos scenarios |
| `COSMOS_SIGNER_KEY` | `test1` | Cosmos->ETH sender |
| `COSMOS_RECEIVER_KEY` | `test` | ETH->Cosmos receiver |
| `AMOUNT` | `1000000000` | Transfer amount |
| `TIMEOUT_SECONDS` | `1800` (success) / `10-45` (timeout) | Packet timeout |
| `BATCH_SIZE` | `25` or `5` | Batch scenarios |
| `KURTOSIS_ENCLAVE` | `my-testnet` | Setup |
| `SEND_PACKET_BIN` | `bin/send-packet` | Cosmos->ETH scenarios |

---

## Scenario Correspondence with Go E2E Tests

The local bash scenarios mirror the `e2e/interchaintestv8` Go test suite. This provides faster iteration without Docker/Kurtosis complexity for the Go test harness.

| Bash Scenario | Go Test Counterpart | File |
|---------------|--------------------|------|
| `run-success.sh` | `Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` (phase 1) | `ibc_eureka_test.go` |
| `run-success-roundtrip.sh` | `Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` | `ibc_eureka_test.go` |
| `run-batch-25.sh` | `Test_25_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` | `ibc_eureka_test.go` |
| `run-batch-50.sh` | `Test_50_ICS20TransferERC20TokenfromEthereumToCosmosAndBack` | `ibc_eureka_test.go` |
| `run-uint256-max.sh` | `Test_ICS20TransferUint256TokenfromEthereumToCosmosAndBack` | `ibc_eureka_test.go` |
| `run-error-ack.sh` | `Test_ErrorAckToEthereum` | `ibc_eureka_test.go` |
| `run-timeout.sh` | `Test_TimeoutPacketFromEth` | `ibc_eureka_test.go` |
| `run-timeout-batch.sh` | `Test_5_TimeoutPacketFromEth`, `Test_10_TimeoutPacketFromEth` | `ibc_eureka_test.go` |
| `run-native-coin.sh` | `Test_ICS20TransferNativeCosmosCoinsToEthereumAndBack` | `ibc_eureka_test.go` |
| `run-timeout.sh` (cosmos-to-eth) | `Test_TimeoutPacketFromCosmos` | `ibc_eureka_test.go` |
| `run-timeout-batch.sh` (cosmos-to-eth) | `Test_10_TimeoutPacketFromCosmos` | `ibc_eureka_test.go` |
| `run-voucher-timeout-eth-remint.sh` | `Test_TimeoutPacketEthRemintsVouchers` | `ibc_eureka_test.go` |
| `run-voucher-timeout-cosmos-remint.sh` | `Test_TimeoutPacketCosmosRemintsVouchers` | `ibc_eureka_test.go` |
| `run-update-client.sh` | `Test_UpdateClient` | `groth16_ics07_test.go` |
| `run-membership-proof.sh` | `Test_Membership` | `groth16_ics07_test.go` |
| `run-misbehaviour.sh` | `Test_DoubleSignMisbehaviour` | `groth16_ics07_test.go` |

---

## Troubleshooting

### "Relayer not running"
Start it first: `./setup/05-relayer.sh`. The relayer must be active for all relay-dependent scenarios (success, roundtrip, timeout, error-ack, native-coin).

### "ETH_RPC_URL not set"
Either:
- Kurtosis is not running → run `./setup/01-eth-node.sh`
- Or manually export `ETH_RPC_URL=http://127.0.0.1:<rpc_port>`

### "ERC20_ADDRESS or ICS20_ADDRESS not found"
Contracts have not been deployed, or `relayer/config.json` is stale. Re-run `./setup/01-eth-node.sh`.

### Timeout scenario fails (packet relayed before timeout)
The relayer may be too fast. Increase polling or reduce `TIMEOUT_SECONDS`:
```bash
TIMEOUT_SECONDS=5 ./scenarios/eth-to-cosmos/run-timeout.sh
```

### "ibcERC20 not minted after N polls"
Check relayer health and ZK proof generation:
```bash
tail -f relayer/relayer.log
```
Groth16 proof generation can take 30-120s depending on bucket size and CPU.

### "Cosmos native balance not restored"
For Cosmos->ETH return transfers, ensure the relayer has updated the Ethereum light client on Cosmos. The first Cosmos->ETH transfer after setup may experience extra latency.

### Cleanup and restart
```bash
./setup/00-cleanup.sh
# Then re-run setup steps 1-6
```

---

## Implementation Notes

- All scripts use `set -euo pipefail` for strict error handling.
- Big-integer math uses `bc` because bash arithmetic overflows at 2^63.
- Contract address discovery follows a cascade: env var -> `config.json` -> Foundry broadcast JSON -> `ICS26Router.getIBCApp()` -> ERC20 mint logs.
- The `send-packet` Go binary (`bin/send-packet`) is auto-built from `relayer/cmd/sendpacket` when missing.
- Multicall batching encodes identical `sendTransfer` calldata N times and submits via `ICS20Transfer.multicall(bytes[])` for gas efficiency.

---

## Further Reading

- `ETH_TO_COSMOS_SCENARIOS_PLAN.md` — Original analysis and implementation plan (phased approach, P0/P1/P2 priorities)
- `docs/ARCHITECTURE.md` — System diagram, contract hierarchy, encoding pipeline
- `docs/PRODUCT_SENSE.md` — IBC domain model, token transfer flow, API surface
- `docs/QUALITY.md` — Test strategy, linting, CI/CD
- `docs/SECURITY.md` — Trust model, proof verification chain, attack surfaces
- `docs/RELIABILITY.md` — Error handling, observability, recovery procedures
