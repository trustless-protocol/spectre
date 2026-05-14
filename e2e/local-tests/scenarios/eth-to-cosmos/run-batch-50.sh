#!/usr/bin/env bash
set -euo pipefail

# Local E2E scenario: Batch of 50 concurrent ERC20 transfers ETH→Cosmos via multicall.
#
# Corresponds to: Test_50_ICS20TransferERC20TokenfromEthereumToCosmosAndBack

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BATCH_SIZE=50 "$SCRIPT_DIR/run-batch-25.sh"
