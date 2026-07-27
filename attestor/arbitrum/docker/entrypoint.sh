#!/bin/sh

set -eu

fail() {
	printf 'attestor-entrypoint: %s\n' "$*" >&2
	exit 1
}

if [ "${1:-}" != "start" ]; then
	exec /usr/local/bin/attestor "$@"
fi

config_path="${ATTESTOR_CONFIG:-/config/config.json}"
grpc_listen_address="${ATTESTOR_GRPC_LISTEN_ADDRESS:-0.0.0.0:50051}"
nitro_binary="${NITRO_BINARY_PATH:-/usr/local/bin/nitro}"
<<<<<<< HEAD
nitro_binary_hash_file="${NITRO_BINARY_SHA256_FILE:-/usr/local/share/fast-ibc/nitro.sha256}"
expected_nitro_binary_hash="${NITRO_BINARY_SHA256:-}"
nitro_allowed_wasm_module_roots="${NITRO_ALLOWED_WASM_MODULE_ROOTS:-/home/user/nitro-legacy/machines,/home/user/target/machines}"
=======
expected_nitro_binary_hash="${NITRO_BINARY_SHA256:-}"
>>>>>>> main
nitro_data_dir="${NITRO_DATA_DIR:-/var/lib/fast-ibc/nitro}"
nitro_ipc_path="${NITRO_IPC_PATH:-/run/fast-ibc/nitro.ipc}"
attestor_state_path="${ATTESTOR_STATE_PATH:-/var/lib/fast-ibc/nitro/attested-roots.json}"
rendered_config="/run/fast-ibc/attestor.json"

[ -r "$config_path" ] || fail "config file is not readable: $config_path"
[ -x "$nitro_binary" ] || fail "Nitro binary is not executable: $nitro_binary"
<<<<<<< HEAD
if [ -z "$expected_nitro_binary_hash" ]; then
	[ -r "$nitro_binary_hash_file" ] ||
		fail "Nitro binary SHA-256 file is not readable: $nitro_binary_hash_file"
	expected_nitro_binary_hash="$(tr -d '[:space:]' <"$nitro_binary_hash_file")"
fi
[ "${#expected_nitro_binary_hash}" -eq 64 ] ||
	fail "expected Nitro binary SHA-256 must contain 64 hexadecimal characters"
[ -n "$nitro_allowed_wasm_module_roots" ] ||
	fail "NITRO_ALLOWED_WASM_MODULE_ROOTS must not be empty"
=======
[ -n "$expected_nitro_binary_hash" ] || fail "NITRO_BINARY_SHA256 is required"
>>>>>>> main

config_dir="$(dirname "$config_path")"
mkdir -p "$nitro_data_dir" "$(dirname "$nitro_ipc_path")" "$(dirname "$rendered_config")"

nitro_binary_hash="$(sha256sum "$nitro_binary" | awk '{print $1}')"
[ "${#nitro_binary_hash}" -eq 64 ] || fail "could not calculate Nitro SHA-256"
[ "$nitro_binary_hash" = "$expected_nitro_binary_hash" ] ||
	fail "Nitro binary SHA-256 does not match the image build pin"

jq \
	--arg grpc_listen_address "$grpc_listen_address" \
	--arg nitro_binary_path "$nitro_binary" \
	--arg nitro_binary_sha256 "$expected_nitro_binary_hash" \
<<<<<<< HEAD
	--arg nitro_allowed_wasm_module_roots "$nitro_allowed_wasm_module_roots" \
=======
>>>>>>> main
	--arg nitro_work_dir "$config_dir" \
	--arg nitro_data_dir "$nitro_data_dir" \
	--arg nitro_ipc_path "$nitro_ipc_path" \
	--arg attestor_state_path "$attestor_state_path" \
	'
	.grpc_listen_address = $grpc_listen_address
	| .nitro_binary_path = $nitro_binary_path
	| .nitro_binary_sha256 = $nitro_binary_sha256
	| .nitro_work_dir = $nitro_work_dir
	| .nitro_data_dir = $nitro_data_dir
	| .nitro_ipc_path = $nitro_ipc_path
	| .attestor_state_path = $attestor_state_path
<<<<<<< HEAD
	| .nitro_arguments = (
		(.nitro_arguments // []) as $arguments
		| if any(
			$arguments[];
			startswith("--validation.wasm.allowed-wasm-module-roots")
		)
		then $arguments
		else $arguments + [
			"--validation.wasm.allowed-wasm-module-roots=" + $nitro_allowed_wasm_module_roots
		]
		end
	)
=======
>>>>>>> main
	' \
	"$config_path" >"$rendered_config"

chmod 600 "$rendered_config"

exec /usr/local/bin/attestor "$@" --config "$rendered_config"
