#!/usr/bin/env bash
# Builds libgaraga_rs, the Rust FFI behind decentrio/ecip-gnark, from the version pinned in
# relayer/go.mod and installs it into .artifacts/native/.
#
# ecip-gnark's cgo directive is `-L..`, which resolves inside the read-only Go module cache, so the
# library cannot live next to the sources. Point the linker and the loader at .artifacts/native
# instead:
#   export CGO_LDFLAGS="-L$PWD/.artifacts/native"
#   export LD_LIBRARY_PATH="$PWD/.artifacts/native"       # DYLD_LIBRARY_PATH on macOS
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
out="$repo_root/.artifacts/native"

case "$(uname -s)" in
  Darwin) lib=libgaraga_rs.dylib ;;
  *) lib=libgaraga_rs.so ;;
esac

module_dir=$(
  cd "$repo_root/relayer"
  go mod download github.com/decentrio/ecip-gnark
  go list -m -f '{{.Dir}}' github.com/decentrio/ecip-gnark
)
if [[ ! -f "$module_dir/ffi/garaga_rs/Cargo.toml" ]]; then
  echo "garaga_rs crate not found under $module_dir" >&2
  exit 1
fi

mkdir -p "$out"
CARGO_TARGET_DIR="$out/target" cargo build --release --locked \
  --manifest-path "$module_dir/ffi/garaga_rs/Cargo.toml"

cp "$out/target/release/$lib" "$out/$lib"
echo "Installed $out/$lib from ${module_dir##*/}"
