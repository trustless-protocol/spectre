#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$repo_root"

tmp_dir=$(mktemp -d)
trap 'rm -rf "$tmp_dir"' EXIT

cp scripts/solidity-refactor/baseline.json "$tmp_dir/baseline.json"

expect_failure() {
  local label=$1
  if python3 scripts/solidity-refactor/check_compatibility.py --baseline "$tmp_dir/baseline.json" >/dev/null 2>&1; then
    echo "negative case unexpectedly passed: $label" >&2
    exit 1
  fi
  cp scripts/solidity-refactor/baseline.json "$tmp_dir/baseline.json"
  echo "negative case rejected: $label"
}

jq '.artifact_contracts.ICS26Router.abi_sha256 = "drift"' \
  scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure abi

jq '.storage_namespaces.ICS20Transfer.structs.ICS20TransferStorage = "drift"' \
  scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure storage

jq '.artifact_contracts.SpectreClient.runtime_sha256 = "drift"' \
  scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure runtime-bytecode

first_fixture=$(jq -r '.digests | keys[0]' scripts/solidity-refactor/baseline.json)
jq --arg fixture "$first_fixture" '.digests[$fixture] = "drift"' \
  scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure fixture-or-binding

jq '.test_count += 1' scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure test-inventory

jq '.gas_baseline.contracts.UpdateClientGasTest["test_gas_n4()"] = 1' \
  scripts/solidity-refactor/baseline.json > "$tmp_dir/baseline.json"
expect_failure gas-regression

jq '.contracts[0].source = "contracts/missing/ICS26Router.sol"' \
  scripts/solidity-refactor/contracts.json > "$tmp_dir/contracts.json"
if python3 scripts/solidity-refactor/check_compatibility.py --contract-manifest "$tmp_dir/contracts.json" >/dev/null 2>&1; then
  echo "negative case unexpectedly passed: fully-qualified-artifact" >&2
  exit 1
fi
echo "negative case rejected: fully-qualified-artifact"

jq '.provenance_path = ".artifacts/does-not-exist.json" | .artifacts.solidity_verifier.sha256 = "drift"' \
  scripts/solidity-refactor/verifier-manifest.json > "$tmp_dir/verifier.json"
if python3 scripts/solidity-refactor/check_compatibility.py --verifier-manifest "$tmp_dir/verifier.json" >/dev/null 2>&1; then
  echo "negative case unexpectedly passed: verifier-provenance" >&2
  exit 1
fi
echo "negative case rejected: verifier-provenance"

jq '.generated_digests["abi/ICS26Router.json"] = "drift"' \
  scripts/solidity-refactor/tooling-rename-manifest.json > "$tmp_dir/tooling-renames.json"
if python3 scripts/solidity-refactor/check_compatibility.py --tooling-rename-manifest "$tmp_dir/tooling-renames.json" >/dev/null 2>&1; then
  echo "negative case unexpectedly passed: tooling-rename-digest" >&2
  exit 1
fi
echo "negative case rejected: tooling-rename-digest"

jq '.entries[0].target = "contracts/missing/Owned.sol"' \
  scripts/solidity-refactor/ownership-manifest.json > "$tmp_dir/ownership.json"
if python3 scripts/solidity-refactor/check_architecture.py --manifest "$tmp_dir/ownership.json" >/dev/null 2>&1; then
  echo "negative case unexpectedly passed: architecture" >&2
  exit 1
fi
echo "negative case rejected: architecture"

python3 scripts/solidity-refactor/check_compatibility.py >/dev/null
python3 scripts/solidity-refactor/check_architecture.py >/dev/null
echo "all compatibility checker self-tests passed"
