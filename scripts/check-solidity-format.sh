#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
cd "$repo_root"

mapfile -t solidity_files < <(find contracts scripts test -type f -name '*.sol' ! -path 'contracts/verifiers/*' -print | sort)
if ((${#solidity_files[@]} == 0)); then
  echo "No tracked Solidity files found"
  exit 0
fi
forge fmt --check "${solidity_files[@]}"
