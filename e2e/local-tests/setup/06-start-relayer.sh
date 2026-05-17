#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RELAYER_DIR="$REPO_ROOT/relayer"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

cd "$RELAYER_DIR"
: > relayer.log
export LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}"
./relayer start --config config.json >> relayer.log 2>&1 &
echo $! > relayer.pid
log_ok "Started relayer PID: $(cat relayer.pid)"
