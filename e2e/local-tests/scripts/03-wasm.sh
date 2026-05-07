#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/.." && pwd)/setup/03-wasm.sh" "$@"
