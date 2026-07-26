#!/usr/bin/env bash

set -euo pipefail

usage() {
	cat >&2 <<'EOF'
Usage: start-attestor.sh [CONFIG_PATH]

Synchronizes the configured Nitro binary SHA-256, builds the attestor, and
starts it. The attestor launches and owns the Nitro child process.

Environment:
  ATTESTOR_BINARY       Output/path of the attestor binary.
  ATTESTOR_SKIP_BUILD   Set to 1 to run an existing binary without rebuilding.
EOF
}

fail() {
	printf 'start-attestor: %s\n' "$*" >&2
	exit 1
}

resolve_binary_path() {
	local configured_binary="$1"

	case "$configured_binary" in
		/*)
			printf '%s\n' "$configured_binary"
			;;
		*/*)
			printf '%s/%s\n' "$config_dir" "$configured_binary"
			;;
		*)
			command -v "$configured_binary" || true
			;;
	esac
}

sha256_file() {
	local path="$1"

	if command -v sha256sum >/dev/null 2>&1; then
		sha256sum "$path" | awk '{print $1}'
	elif command -v shasum >/dev/null 2>&1; then
		shasum -a 256 "$path" | awk '{print $1}'
	else
		fail "sha256sum or shasum is required"
	fi
}

if [[ "${1:-}" == "-h" || "${1:-}" == "--help" ]]; then
	usage
	exit 0
fi
if [[ "$#" -gt 1 ]]; then
	usage
	exit 2
fi

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
attestor_dir="$(cd "$script_dir/.." && pwd)"
requested_config="${1:-$attestor_dir/config.json}"
[[ -f "$requested_config" ]] || fail "config file not found: $requested_config"

config_dir="$(cd "$(dirname "$requested_config")" && pwd)"
config_path="$config_dir/$(basename "$requested_config")"

command -v jq >/dev/null 2>&1 || fail "jq is required"
configured_binary="$(jq -er '.nitro_binary_path | select(type == "string" and length > 0)' "$config_path")" ||
	fail "nitro_binary_path must be a non-empty string in $config_path"
nitro_binary="$(resolve_binary_path "$configured_binary")"
[[ -n "$nitro_binary" && -f "$nitro_binary" ]] || fail "Nitro binary not found: $configured_binary"
[[ -x "$nitro_binary" ]] || fail "Nitro binary is not executable: $nitro_binary"

nitro_binary_hash="$(sha256_file "$nitro_binary")"
[[ "$nitro_binary_hash" =~ ^[[:xdigit:]]{64}$ ]] || fail "could not calculate Nitro SHA-256"
nitro_binary_hash="$(printf '%s' "$nitro_binary_hash" | tr '[:upper:]' '[:lower:]')"

temporary_config="$(mktemp "$config_dir/.attestor-config.XXXXXX")"
cleanup() {
	if [[ -n "${temporary_config:-}" && -f "$temporary_config" ]]; then
		rm -f "$temporary_config"
	fi
}
trap cleanup EXIT

jq --arg hash "$nitro_binary_hash" '.nitro_binary_sha256 = $hash' "$config_path" >"$temporary_config"
chmod 600 "$temporary_config"
mv "$temporary_config" "$config_path"
temporary_config=""
trap - EXIT
printf 'Nitro SHA-256 synchronized in %s: %s\n' "$config_path" "$nitro_binary_hash"

attestor_binary="${ATTESTOR_BINARY:-$attestor_dir/bin/attestor}"
if [[ "$attestor_binary" != /* ]]; then
	attestor_binary="$attestor_dir/$attestor_binary"
fi

if [[ "${ATTESTOR_SKIP_BUILD:-0}" != "1" ]]; then
	command -v go >/dev/null 2>&1 || fail "go is required to build the attestor"
	mkdir -p "$(dirname "$attestor_binary")"
	printf 'Building attestor at %s\n' "$attestor_binary"
	(
		cd "$attestor_dir"
		go build -o "$attestor_binary" ./cmd
	)
fi

[[ -x "$attestor_binary" ]] || fail "attestor binary is not executable: $attestor_binary"
printf 'Starting attestor with %s\n' "$config_path"
exec "$attestor_binary" start --config "$config_path"
