#!/bin/sh
set -eu

cd "$(dirname "$0")/../.."
exec ./scripts/local/store_wasm_artifact.sh \
  artifacts/cw_ics08_wasm_avalanche.wasm \
  avalanche-light-client \
  "cw-ics08-wasm Avalanche C-Chain signed light client"
