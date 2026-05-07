# ETH → Cosmos Scenarios

| Script | Description |
|--------|-------------|
| `transfer.sh` | Send an ICS20 transfer from Ethereum to Cosmos |
| `check-eth.sh` | Inspect Ethereum state (ERC20, ICS20, account) |
| `check-cosmos.sh` | Inspect Cosmos state (voucher balance, IBC clients) |
| `run-success.sh` | Orchestrate full success flow (transfer + verify) |
| `run-timeout.sh` | Test timeout scenario (short timeout, verify packet skipped) |

## Prerequisites

Run `setup/` steps 01–05 first to start infrastructure and the relayer.
