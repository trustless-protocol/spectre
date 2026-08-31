#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
staging="$repo_root/.artifacts/solidity/prover"
native_lib="$repo_root/third_party/ecip-gnark"
export LD_LIBRARY_PATH="$native_lib${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
export DYLD_LIBRARY_PATH="$native_lib${DYLD_LIBRARY_PATH:+:$DYLD_LIBRARY_PATH}"
force=0
case "${1:-}" in
  "") ;;
  --force) force=1 ;;
  *) echo "usage: $0 [--force]" >&2; exit 2 ;;
esac

buckets=$(
  cd "$repo_root/relayer"
  go run ./prover/cmd -list-buckets
)
if [[ -z "$buckets" ]]; then
  echo "prover generator reported no configured buckets" >&2
  exit 1
fi

complete=1
for bucket in $buckets; do
  for artifact in r1cs.bin pk.bin vk.bin; do
    [[ -s "$repo_root/relayer/bin/n$bucket/$artifact" ]] || complete=0
  done
  [[ -s "$repo_root/contracts/verifiers/Groth16Verifier_N$bucket.sol" ]] || complete=0
done
if ((complete == 1 && force == 0)); then
  echo "Prover artifacts already present for all configured buckets; skipping randomized setup."
  echo "Use $0 --force to regenerate them intentionally."
  exit 0
fi

mkdir -p "$staging"
run_staging=$(mktemp -d "$staging/run.XXXXXX")
generated="$run_staging/generated"
trap 'rm -rf "$run_staging"' EXIT

mkdir -p "$generated/bin" "$generated/verifiers"

(
  cd "$repo_root/relayer"
  go run ./prover/cmd "$generated/bin" "$generated/verifiers"
)

for bucket in $buckets; do
  for artifact in r1cs.bin pk.bin vk.bin; do
    [[ -s "$generated/bin/n$bucket/$artifact" ]] || {
      echo "missing staged n$bucket/$artifact" >&2
      exit 1
    }
  done
  [[ -s "$generated/verifiers/Groth16Verifier_N$bucket.sol" ]] || {
    echo "missing staged Groth16Verifier_N$bucket.sol" >&2
    exit 1
  }
done

mkdir -p "$repo_root/contracts/verifiers"
for bucket in $buckets; do
  mkdir -p "$repo_root/relayer/bin/n$bucket"
  cp "$generated/bin/n$bucket/r1cs.bin" "$repo_root/relayer/bin/n$bucket/r1cs.bin"
  cp "$generated/bin/n$bucket/pk.bin" "$repo_root/relayer/bin/n$bucket/pk.bin"
  cp "$generated/bin/n$bucket/vk.bin" "$repo_root/relayer/bin/n$bucket/vk.bin"
  cp "$generated/verifiers/Groth16Verifier_N$bucket.sol" \
    "$repo_root/contracts/verifiers/Groth16Verifier_N$bucket.sol"
done

python3 - "$repo_root" "$staging/provenance.json" $buckets <<'PY'
import hashlib
import json
import sys
from pathlib import Path

root = Path(sys.argv[1])
output = Path(sys.argv[2])
buckets = [int(value) for value in sys.argv[3:]]
paths = {}
for bucket in buckets:
    prefix = f'n{bucket}'
    paths[f'{prefix}_r1cs'] = f'relayer/bin/{prefix}/r1cs.bin'
    paths[f'{prefix}_proving_key'] = f'relayer/bin/{prefix}/pk.bin'
    paths[f'{prefix}_verifying_key'] = f'relayer/bin/{prefix}/vk.bin'
    paths[f'{prefix}_solidity_verifier'] = f'contracts/verifiers/Groth16Verifier_N{bucket}.sol'
data = {
    'schema_version': 1,
    'scope': 'checked-local-generator',
    'supported_buckets': buckets,
    'generator': 'cd relayer && go run ./prover/cmd <staging-bin> <staging-verifiers>',
    'artifacts': {
        key: {'path': value, 'sha256': hashlib.sha256((root / value).read_bytes()).hexdigest()}
        for key, value in paths.items()
    },
}
output.write_text(json.dumps(data, indent=2) + '\n')
PY

echo "Published paired artifacts for configured buckets: $(printf '%s ' $buckets)"
echo "Provenance: ${staging#"$repo_root/"}"/provenance.json
