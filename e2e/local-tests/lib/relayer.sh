#!/usr/bin/env bash

set -euo pipefail

# Relayer helpers for local E2E tests.
# The relayer runs as a background daemon; bash scripts rely on it auto-relaying.
# These helpers focus on process health and log inspection.

RELAYER_DIR="${RELAYER_DIR:-$REPO_ROOT/relayer}"
RELAYER_PID_FILE="${RELAYER_PID_FILE:-$RELAYER_DIR/relayer.pid}"
RELAYER_LOG_FILE="${RELAYER_LOG_FILE:-$RELAYER_DIR/relayer.log}"

relayer_is_running() {
  if [ -f "$RELAYER_PID_FILE" ]; then
    local pid
    pid="$(cat "$RELAYER_PID_FILE")"
    kill -0 "$pid" 2>/dev/null
  else
    false
  fi
}

relayer_ensure_running() {
  if ! relayer_is_running; then
    log_err "Relayer is not running. Start it first:"
    log_err "  ./setup/05-relayer.sh"
    exit 1
  fi
  local pid
  pid="$(cat "$RELAYER_PID_FILE")"
  log_ok "Relayer running (PID $pid)"
}

relayer_wait_for_log() {
  local pattern="$1" timeout_sec="${2:-120}"
  if [ ! -f "$RELAYER_LOG_FILE" ]; then
    log_warn "Relayer log file not found: $RELAYER_LOG_FILE"
    return 1
  fi
  local poll=0
  while [ $poll -lt "$timeout_sec" ]; do
    if grep -q "$pattern" "$RELAYER_LOG_FILE" 2>/dev/null; then
      return 0
    fi
    sleep 2
    poll=$((poll + 2))
  done
  return 1
}

relayer_tail_log() {
  local lines="${1:-20}"
  if [ -f "$RELAYER_LOG_FILE" ]; then
    tail -n "$lines" "$RELAYER_LOG_FILE"
  else
    log_warn "Relayer log file not found: $RELAYER_LOG_FILE"
  fi
}
