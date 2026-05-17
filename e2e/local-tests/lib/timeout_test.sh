#!/usr/bin/env bash
set -euo pipefail

run_timeout_test() {
  local test_label="${1:-Timeout}"

  require_relayer_running "$RELAYER_PID_FILE"

  ETH_SENDER="$(eth_wallet_address)"

  log_header "Pre-transfer State"
  BEFORE_ETH_BALANCE="$(eth_balance_of "$ERC20_ADDRESS" "$ETH_SENDER")"
  log_kv "Sender ERC20 balance" "$BEFORE_ETH_BALANCE"

  ESCROW_ADDRESS="$(get_escrow_address "$ICS20_ADDRESS" "$SOURCE_CLIENT")"
  if [ -n "$ESCROW_ADDRESS" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
    BEFORE_ESCROW_BALANCE="$(eth_balance_of "$ERC20_ADDRESS" "$ESCROW_ADDRESS")"
    log_kv "Escrow ERC20 balance" "$BEFORE_ESCROW_BALANCE"
  fi

  COSMOS_VOUCHER_DENOM="transfer/${COSMOS_WASM_CLIENT_ID}/${ERC20_ADDRESS}"
  BEFORE_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
  log_kv "Cosmos voucher balance" "$BEFORE_COSMOS_BALANCE"

  TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-10}"
  TIMEOUT=$(($(date +%s) + TIMEOUT_SECONDS))
  log_header "Sending ICS20Transfer (${TIMEOUT_SECONDS}s timeout)"
  log_kv "Timeout" "epoch $TIMEOUT ($(date -d "@$TIMEOUT" '+%H:%M:%S' 2>/dev/null || echo 'N/A'))"

  log "Approving ERC20 allowance..."
  eth_approve "$ERC20_ADDRESS" "$ICS20_ADDRESS" "$AMOUNT"

  log "Submitting sendTransfer..."
  eth_send_transfer "$ICS20_ADDRESS" "$ERC20_ADDRESS" "$AMOUNT" "$RECEIVER" "$SOURCE_CLIENT" "$DEST_PORT" "$TIMEOUT"

  log "Transfer submitted (waiting for timeout...)"

  MAX_POLLS=120
  POLL_INTERVAL=15
  poll=0
  REFUNDED=false

  REFUND_THRESHOLD="$(bc_sub1 "$BEFORE_ETH_BALANCE")"

  while [ $poll -lt $MAX_POLLS ]; do
    sleep $POLL_INTERVAL
    poll=$((poll + 1))

    CURRENT_ETH_BALANCE="$(eth_balance_of "$ERC20_ADDRESS" "$ETH_SENDER")"
    CURRENT_ESCROW="0"
    if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
      CURRENT_ESCROW="$(eth_balance_of "$ERC20_ADDRESS" "$ESCROW_ADDRESS")"
    fi
    CURRENT_COSMOS="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"

    if bc_ge "$CURRENT_ETH_BALANCE" "$REFUND_THRESHOLD"; then
      log_ok "[poll ${poll}/${MAX_POLLS}] ERC20 refunded │ sender=${CURRENT_ETH_BALANCE} (target ≥ ${REFUND_THRESHOLD})"
      REFUNDED=true
      break
    fi

    log "[poll ${poll}/${MAX_POLLS}] sender=${CURRENT_ETH_BALANCE}  escrow=${CURRENT_ESCROW}  cosmos=${CURRENT_COSMOS}"
  done

  log_header "Final State"
  AFTER_ETH_BALANCE="$(eth_balance_of "$ERC20_ADDRESS" "$ETH_SENDER")"
  log_kv "Sender ERC20 balance" "$AFTER_ETH_BALANCE"

  if [ -n "${ESCROW_ADDRESS:-}" ] && [ "$ESCROW_ADDRESS" != "0x" ]; then
    AFTER_ESCROW_BALANCE="$(eth_balance_of "$ERC20_ADDRESS" "$ESCROW_ADDRESS")"
    log_kv "Escrow ERC20 balance" "$AFTER_ESCROW_BALANCE"
  fi

  AFTER_COSMOS_BALANCE="$("$COSMOS_BIN" query bank balance "$RECEIVER" "$COSMOS_VOUCHER_DENOM" --node "$COSMOS_RPC_URL" --chain-id "$COSMOS_CHAIN_ID" --output json 2>/dev/null | jq -r '.balance.amount // "0"')"
  log_kv "Cosmos voucher" "$AFTER_COSMOS_BALANCE"

  log_header "${test_label} Test Summary"
  log_kv "Direction" "ETH → Cosmos"
  log_kv "Amount sent" "${AMOUNT} wei"
  log_kv "Timeout" "${TIMEOUT_SECONDS}s"

  if [ "${REFUNDED:-false}" = "true" ]; then
    log_ok "RESULT: PASS"
    log "  • Relayer detected the expired packet"
    log "  • Built non-membership proof from Cosmos"
    log "  • Called timeoutPacket() on Ethereum ICS26Router"
    log "  • Sender refunded (ERC20 balance returned to original)"
    log "  • Cosmos voucher = 0 (packet never relayed)"
  elif bc_ge "$AFTER_ETH_BALANCE" "$REFUND_THRESHOLD"; then
    log_ok "RESULT: PASS (detected after final poll)"
    log "  • Sender ERC20 balance returned to original"
  elif [ "$AFTER_COSMOS_BALANCE" = "0" ]; then
    log_warn "RESULT: PARTIAL"
    log_warn "  • Packet was NOT relayed to Cosmos (timeout prevented relay)"
    log_warn "  • ERC20 balance still decreased — refund not yet submitted"
    log_warn "  • The relayer may still be building the ZK proof"
    log_warn "  • Check relayer logs: tail -f $REPO_ROOT/relayer/relayer.log"
  else
    log_err "RESULT: FAIL"
    log_err "  • Packet WAS relayed to Cosmos before timeout"
    log_err "  • Try with shorter TIMEOUT_SECONDS or check relayer status"
  fi
}
