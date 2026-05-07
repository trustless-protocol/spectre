#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RELAYER_DIR="$REPO_ROOT/relayer"

CONFIG_FILE="${CONFIG_FILE:-$RELAYER_DIR/config.json}"

# Auto-copy config.example.json if config.json doesn't exist
if [ ! -f "$CONFIG_FILE" ] && [ -f "$RELAYER_DIR/config.example.json" ]; then
  cp "$RELAYER_DIR/config.example.json" "$CONFIG_FILE"
  echo "Created $CONFIG_FILE from config.example.json"
fi

cd "$RELAYER_DIR"

LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}" \
  ./relayer start --config "$(basename "$CONFIG_FILE")"
