#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RELAYER_DIR="$REPO_ROOT/relayer"
KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"

CONFIG_FILE="${CONFIG_FILE:-$RELAYER_DIR/config.json}"
RELAYER_PID_FILE="${RELAYER_PID_FILE:-$RELAYER_DIR/relayer.pid}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

build_relayer_if_needed() {
  if [ ! -x "$RELAYER_DIR/relayer" ] || find "$RELAYER_DIR" -name '*.go' -newer "$RELAYER_DIR/relayer" | grep -q .; then
    echo "  Building relayer binary..."
    (cd "$RELAYER_DIR" && go build -o relayer ./cmd)
  fi
}

require_cmd cast
require_cmd jq
ensure_relayer_config "$CONFIG_FILE" "$RELAYER_DIR/config.example.json"
refresh_relayer_eth_endpoints "$CONFIG_FILE"

cd "$RELAYER_DIR"
build_relayer_if_needed

PROVER_BIN_DIR="${PROVER_BIN_DIR:-$RELAYER_DIR/bin}" \
LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}" \
  ./relayer start --config "$CONFIG_FILE" &
RELAYER_PID="$!"
echo "$RELAYER_PID" > "$RELAYER_PID_FILE"
cleanup() {
  if kill -0 "$RELAYER_PID" 2>/dev/null; then
    kill "$RELAYER_PID" 2>/dev/null || true
  fi
  rm -f "$RELAYER_PID_FILE"
}
trap cleanup EXIT INT TERM
wait "$RELAYER_PID"
