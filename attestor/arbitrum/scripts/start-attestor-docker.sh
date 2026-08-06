#!/usr/bin/env bash

set -euo pipefail

usage() {
	cat >&2 <<'EOF'
Usage: start-attestor-docker.sh [CONFIG_PATH]

Builds and runs the attestor against the Nitro HTTP and WebSocket endpoints in
CONFIG_PATH.

Environment:
  ATTESTOR_DOCKER_IMAGE       Image name (default: fast-ibc-attestor:local).
  ATTESTOR_CONTAINER_NAME     Container name (default: fast-ibc-attestor).
  ATTESTOR_GRPC_PUBLISH       Host publish address (default: 127.0.0.1:50051).
  ATTESTOR_STATE_VOLUME       Persistent attestor-state volume
                              (default: fast-ibc-attestor-data).
EOF
}

fail() {
	printf 'start-attestor-docker: %s\n' "$*" >&2
	exit 1
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

command -v docker >/dev/null 2>&1 || fail "docker is required"
[[ -f "$requested_config" ]] || fail "config file not found: $requested_config"

config_dir="$(cd "$(dirname "$requested_config")" && pwd)"
config_name="$(basename "$requested_config")"
image="${ATTESTOR_DOCKER_IMAGE:-fast-ibc-attestor:local}"
container_name="${ATTESTOR_CONTAINER_NAME:-fast-ibc-attestor}"
grpc_publish="${ATTESTOR_GRPC_PUBLISH:-127.0.0.1:50051}"
state_volume="${ATTESTOR_STATE_VOLUME:-fast-ibc-attestor-data}"

printf 'Building attestor image %s\n' "$image"
docker build \
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
	--volume "$state_volume:/var/lib/fast-ibc/attestor" \
	--env "ATTESTOR_CONFIG=/config/$config_name" \
	"$image"
