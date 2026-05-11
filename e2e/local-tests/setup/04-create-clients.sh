#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
STATE_DIR="$SCRIPT_DIR/../.state"
RELAYER_DIR="$REPO_ROOT/relayer"
KURTOSIS_ENCLAVE="${KURTOSIS_ENCLAVE:-my-testnet}"

WASM_CHECKSUM_FILE="$STATE_DIR/wasm_checksum"
CONFIG_FILE="${CONFIG_FILE:-$RELAYER_DIR/config.json}"

source "$REPO_ROOT/e2e/local-tests/lib/common.sh"

build_relayer_if_needed() {
  if [ ! -x "$RELAYER_DIR/relayer" ] || find "$RELAYER_DIR" -name '*.go' -newer "$RELAYER_DIR/relayer" | grep -q .; then
    echo "Building current relayer binary..."
    (cd "$RELAYER_DIR" && go build -o relayer ./cmd)
  fi
}

# Auto-detect wasm checksum if not provided
WASM_CHECKSUM="${WASM_CHECKSUM:-}"
if [ -z "$WASM_CHECKSUM" ] && [ -f "$WASM_CHECKSUM_FILE" ]; then
  WASM_CHECKSUM="$(cat "$WASM_CHECKSUM_FILE")"
  echo "Auto-detected WASM checksum: $WASM_CHECKSUM"
fi

if [ -z "$WASM_CHECKSUM" ]; then
  echo "Error: WASM checksum not found." >&2
  echo "Please run ./scripts/03-wasm.sh first, or set WASM_CHECKSUM env var." >&2
  exit 1
fi

require_cmd cast
require_cmd jq
ensure_relayer_config "$CONFIG_FILE" "$RELAYER_DIR/config.example.json"
refresh_relayer_eth_endpoints "$CONFIG_FILE"

cd "$RELAYER_DIR"
build_relayer_if_needed

# Persist wasm checksum in .env so the relayer can find it
ENV_FILE="$RELAYER_DIR/.env"
source "$REPO_ROOT/e2e/local-tests/lib/cosmos.sh"
upsert_env_var "$ENV_FILE" WASM_CHECKSUM "$WASM_CHECKSUM"

PROVER_BIN_DIR="${PROVER_BIN_DIR:-$RELAYER_DIR/bin}" \
LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}" \
  ./relayer create-clients \
    --config "$CONFIG_FILE" \
    --wasm-checksum "$WASM_CHECKSUM"
