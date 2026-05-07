#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
STATE_DIR="$SCRIPT_DIR/../.state"
RELAYER_DIR="$REPO_ROOT/relayer"

WASM_CHECKSUM_FILE="$STATE_DIR/wasm_checksum"
CONFIG_FILE="${CONFIG_FILE:-$RELAYER_DIR/config.json}"

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

# Check if config.json exists, if not copy from config.example.json
if [ ! -f "$CONFIG_FILE" ] && [ -f "$RELAYER_DIR/config.example.json" ]; then
  cp "$RELAYER_DIR/config.example.json" "$CONFIG_FILE"
  echo "Created $CONFIG_FILE from config.example.json"
fi

cd "$RELAYER_DIR"

LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}" \
  ./relayer create-clients \
    --config "$(basename "$CONFIG_FILE")" \
    --wasm-checksum "$WASM_CHECKSUM"
