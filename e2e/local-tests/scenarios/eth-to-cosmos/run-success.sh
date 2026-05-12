#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../../.." && pwd)"
source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

log_header "ETH → Cosmos Transfer: Success Case"

log "Step 1: Sending transfer..."
bash "$SCRIPT_DIR/transfer.sh"

log "Step 2: Checking Ethereum state..."
bash "$SCRIPT_DIR/check-eth.sh"

log "Step 3: Checking Cosmos state (voucher balance)..."
bash "$SCRIPT_DIR/check-cosmos.sh"

log_header "Success case complete"
log "If the relayer is running, tokens should have been transferred"
log "from Ethereum to Cosmos. Check the output above for voucher balances."
