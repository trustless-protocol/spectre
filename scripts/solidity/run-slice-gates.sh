#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
cd "$repo_root"

slice=${1:-full}
case "$slice" in
  core) focused='test/core/*' ;;
  ics20) focused='test/apps/ics20/*' ;;
  spectre) focused='test/light-clients/spectre/*' ;;
  integration) focused='test/integration/*' ;;
  deployment) focused='test/deployment/*' ;;
  full) focused='test/**/*' ;;
  *) echo "unknown slice: $slice" >&2; exit 2 ;;
esac

scripts/check-solidity-format.sh
forge build --force --skip test --skip script --sizes
forge test --match-path "$focused"
scripts/check-solidity-compatibility.sh
git diff --check
