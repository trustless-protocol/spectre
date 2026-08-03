# Helpers for local 08-wasm proposal scripts.

wasm_checksum_from_proposal() {
    proposal_file=$1
    python3 - "$proposal_file" <<'PY'
import base64
import gzip
import hashlib
import json
import sys

with open(sys.argv[1], "r", encoding="utf-8") as fh:
    proposal = json.load(fh)

for msg in proposal.get("messages", []):
    blob = msg.get("wasm_byte_code")
    if not blob:
        continue
    wasm = base64.b64decode(blob)
    if wasm[:2] == b"\x1f\x8b":
        wasm = gzip.decompress(wasm)
    print(hashlib.sha256(wasm).hexdigest())
    break
else:
    raise SystemExit("wasm_byte_code not found in proposal")
PY
}

wasm_checksum_is_stored() {
    checksum=$1
    shift
    "$GAIAD" q ibc-wasm checksums "$@" -o json \
        | jq -e --arg checksum "$checksum" '(.checksums // []) | index($checksum) != null' >/dev/null
}

assert_wasm_checksum_stored() {
    checksum=$1
    shift
    if ! wasm_checksum_is_stored "$checksum" "$@"; then
        echo "ERROR: computed wasm checksum $checksum was not stored on-chain" >&2
        echo "       Check the proposal status and execution logs before using this checksum." >&2
        return 1
    fi
}
