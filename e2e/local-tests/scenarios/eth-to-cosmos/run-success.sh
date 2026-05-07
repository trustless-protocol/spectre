#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== ETH → Cosmos Transfer: Success Case ==="
echo ""

echo "Step 1: Sending transfer..."
bash "$SCRIPT_DIR/transfer.sh"
echo ""

echo "Step 2: Checking Ethereum state..."
bash "$SCRIPT_DIR/check-eth.sh"
echo ""

echo "Step 3: Checking Cosmos state (voucher balance)..."
bash "$SCRIPT_DIR/check-cosmos.sh"
echo ""

echo "=== Success case complete ==="
echo "If the relayer is running, tokens should have been transferred from"
echo "Ethereum to Cosmos. Check the output above for voucher balances."
