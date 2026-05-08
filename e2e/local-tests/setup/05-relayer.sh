#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
RELAYER_DIR="$REPO_ROOT/relayer"

CONFIG_FILE="${CONFIG_FILE:-$RELAYER_DIR/config.json}"
RELAYER_PID_FILE="${RELAYER_PID_FILE:-$RELAYER_DIR/relayer.pid}"

build_relayer_if_needed() {
  if [ ! -x "$RELAYER_DIR/relayer" ] || find "$RELAYER_DIR" -name '*.go' -newer "$RELAYER_DIR/relayer" | grep -q .; then
    echo "Building current relayer binary..."
    (cd "$RELAYER_DIR" && go build -o relayer ./cmd)
  fi
}

# config.json is the default config file; verify it exists
if [ ! -f "$CONFIG_FILE" ]; then
  echo "Error: $CONFIG_FILE not found. Create it from the example template." >&2
  exit 1
fi

cd "$RELAYER_DIR"
build_relayer_if_needed

LD_LIBRARY_PATH="${LD_LIBRARY_PATH:-$HOME/works/ecip-gnark}" \
  ./relayer start --config "$(basename "$CONFIG_FILE")" &
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
