#!/usr/bin/env bash
set -euo pipefail

export TIMEOUT_SECONDS="${TIMEOUT_SECONDS:-15}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
"$SCRIPT_DIR/run-timeout.sh" "$@"
