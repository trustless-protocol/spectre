#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

forge build --skip test --skip script
python3 scripts/solidity/check_architecture.py
python3 scripts/solidity/check_compatibility.py
