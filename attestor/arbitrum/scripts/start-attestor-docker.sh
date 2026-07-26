#!/usr/bin/env bash

set -euo pipefail

usage() {
	cat >&2 <<'EOF'
Usage: start-attestor-docker.sh [CONFIG_PATH]

Builds one image containing the attestor and a prebuilt Linux Nitro binary,
then runs the attestor with Nitro as its managed child process.

Environment:
  ATTESTOR_DOCKER_IMAGE       Image name (default: fast-ibc-attestor:local).
  ATTESTOR_CONTAINER_NAME     Container name (default: fast-ibc-attestor).
  ATTESTOR_GRPC_PUBLISH       Host publish address (default: 127.0.0.1:50051).
  ATTESTOR_NITRO_VOLUME       Persistent Nitro volume (default: fast-ibc-nitro-data).
EOF
}

fail() {
	printf 'start-attestor-docker: %s\n' "$*" >&2
	exit 1
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
attestor_root="$(cd "$attestor_dir/.." && pwd)"
requested_config="${1:-$attestor_dir/config.json}"
nitro_binary="$attestor_dir/bin/nitro"

command -v docker >/dev/null 2>&1 || fail "docker is required"
[[ -f "$requested_config" ]] || fail "config file not found: $requested_config"
[[ -f "$nitro_binary" ]] || fail "prebuilt Linux Nitro binary not found: $nitro_binary"
[[ -x "$nitro_binary" ]] || fail "Nitro binary is not executable: $nitro_binary"

config_dir="$(cd "$(dirname "$requested_config")" && pwd)"
config_name="$(basename "$requested_config")"
image="${ATTESTOR_DOCKER_IMAGE:-fast-ibc-attestor:local}"
container_name="${ATTESTOR_CONTAINER_NAME:-fast-ibc-attestor}"
grpc_publish="${ATTESTOR_GRPC_PUBLISH:-127.0.0.1:50051}"
nitro_volume="${ATTESTOR_NITRO_VOLUME:-fast-ibc-nitro-data}"
nitro_binary_hash="$(sha256_file "$nitro_binary")"
[[ "$nitro_binary_hash" =~ ^[[:xdigit:]]{64}$ ]] || fail "could not calculate Nitro SHA-256"
nitro_binary_hash="$(printf '%s' "$nitro_binary_hash" | tr '[:upper:]' '[:lower:]')"

printf 'Building attestor image %s\n' "$image"
docker build \
	--build-arg "NITRO_BINARY_SHA256=$nitro_binary_hash" \
	--file "$attestor_dir/Dockerfile" \
	--tag "$image" \
	"$attestor_root"

printf 'Starting %s with config %s\n' "$container_name" "$requested_config"
exec docker run \
	--rm \
	--init \
	--name "$container_name" \
	--stop-timeout 40 \
	--publish "$grpc_publish:50051" \
	--volume "$config_dir:/config:ro" \
	--volume "$nitro_volume:/var/lib/fast-ibc/nitro" \
	--env "ATTESTOR_CONFIG=/config/$config_name" \
	"$image"
