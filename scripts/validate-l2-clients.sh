#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

cargo test --locked \
  -p l2-client -p op-stack-verifier -p arbitrum-verifier -p base-verifier -p op-verifier \
  -p cw-ics08-wasm-arbitrum -p cw-ics08-wasm-base -p cw-ics08-wasm-op
cargo build --target wasm32-unknown-unknown --release \
  -p cw-ics08-wasm-arbitrum -p cw-ics08-wasm-base -p cw-ics08-wasm-op --locked
