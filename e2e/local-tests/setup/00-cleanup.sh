#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

log_header "Cleanup: Killing relayer"
RELAYER_PID_FILE="${RELAYER_PID_FILE:-$REPO_ROOT/relayer/relayer.pid}"
if [ -f "$RELAYER_PID_FILE" ]; then
  RELAYER_PID="$(cat "$RELAYER_PID_FILE" 2>/dev/null || true)"
  if [ -n "$RELAYER_PID" ] && kill -0 "$RELAYER_PID" 2>/dev/null; then
    kill "$RELAYER_PID" 2>/dev/null
    log_ok "Killed relayer (PID $RELAYER_PID)"
  fi
  rm -f "$RELAYER_PID_FILE"
else
  pkill -f "relayer start" 2>/dev/null && log_ok "Killed relayer process" || log "No relayer process found"
fi

log_header "Cleanup: Killing Cosmos validators"
if killall gaiad 2>/dev/null; then
  log_ok "Killed gaiad processes"
else
  log "No gaiad processes running"
fi

log_header "Cleanup: Removing Cosmos state"
rm -rf "$HOME/.gaia" "$HOME/.gaia-val2" "$HOME/.gaia-val3" "$HOME/.gaia-val4"
log_ok "Removed Gaia home directories"

log_header "Cleanup: Removing Kurtosis enclave"
KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"
if command -v kurtosis &>/dev/null; then
  kurtosis enclave rm -f "$KURTOSIS_ENCLAVE" 2>/dev/null && log_ok "Removed Kurtosis enclave: $KURTOSIS_ENCLAVE" || log "No Kurtosis enclave: $KURTOSIS_ENCLAVE"
else
  log "kurtosis not installed, skipping"
fi

log_header "Cleanup: Removing state files"
STATE_DIR="$REPO_ROOT/e2e/local-tests/.state"
rm -rf "$STATE_DIR"
log_ok "Removed $STATE_DIR"

log "Cleanup complete"
