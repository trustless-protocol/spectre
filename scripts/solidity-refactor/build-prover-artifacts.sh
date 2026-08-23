#!/usr/bin/env bash
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)
staging="$repo_root/.artifacts/solidity-refactor/prover"
mkdir -p "$staging"
run_staging=$(mktemp -d "$staging/run.XXXXXX")
generated="$run_staging/generated"
native_lib="$repo_root/third_party/ecip-gnark"
export LD_LIBRARY_PATH="$native_lib${LD_LIBRARY_PATH:+:$LD_LIBRARY_PATH}"
trap 'rm -rf "$run_staging"' EXIT

mkdir -p "$generated/bin" "$generated/verifiers"

(
  cd "$repo_root/relayer"
  go run ./prover/cmd "$generated/bin" "$generated/verifiers"
)

mapfile -t generated_verifiers < <(find "$generated/verifiers" -maxdepth 1 -type f -name "Groth16Verifier_N*.sol" -printf "%f\n" | sort)
if [[ "${generated_verifiers[*]}" != "Groth16Verifier_N4.sol" ]]; then
  echo "generator published unsupported verifier set: ${generated_verifiers[*]}" >&2
  exit 1
fi

for artifact in r1cs.bin pk.bin vk.bin; do
  [[ -s "$generated/bin/n4/$artifact" ]] || { echo "missing staged n4/$artifact" >&2; exit 1; }
done

mkdir -p "$repo_root/relayer/bin/n4" "$repo_root/contracts/verifiers"
cp "$generated/bin/n4/r1cs.bin" "$repo_root/relayer/bin/n4/r1cs.bin"
cp "$generated/bin/n4/pk.bin" "$repo_root/relayer/bin/n4/pk.bin"
cp "$generated/bin/n4/vk.bin" "$repo_root/relayer/bin/n4/vk.bin"
cp "$generated/verifiers/Groth16Verifier_N4.sol" "$repo_root/contracts/verifiers/Groth16Verifier_N4.sol"

python3 - "$repo_root" "$staging/provenance.json" <<'PY'
import hashlib
import json
import sys
from pathlib import Path

root = Path(sys.argv[1])
output = Path(sys.argv[2])
paths = {
    'r1cs': 'relayer/bin/n4/r1cs.bin',
    'proving_key': 'relayer/bin/n4/pk.bin',
    'verifying_key': 'relayer/bin/n4/vk.bin',
    'solidity_verifier': 'contracts/verifiers/Groth16Verifier_N4.sol',
}
data = {
    'schema_version': 1,
    'scope': 'checked-local-generator',
    'supported_buckets': [4],
    'generator': 'cd relayer && go run ./prover/cmd <staging-bin> <staging-verifiers>',
    'artifacts': {
        key: {'path': value, 'sha256': hashlib.sha256((root / value).read_bytes()).hexdigest()}
        for key, value in paths.items()
    },
}
output.write_text(json.dumps(data, indent=2) + '\n')
PY

echo "Published paired N4 artifacts; provenance: ${staging#"$repo_root/"}"/provenance.json
