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
attestor_state_path="${ATTESTOR_STATE_PATH:-/var/lib/fast-ibc/attestor/attested-roots.json}"
rendered_config="/run/fast-ibc/attestor.json"

[ -r "$config_path" ] || fail "config file is not readable: $config_path"
mkdir -p "$(dirname "$attestor_state_path")" "$(dirname "$rendered_config")"

jq \
	--arg grpc_listen_address "$grpc_listen_address" \
	--arg attestor_state_path "$attestor_state_path" \
	'
	.grpc_listen_address = $grpc_listen_address
	| .attestor_state_path = $attestor_state_path
	' \
	"$config_path" >"$rendered_config"

chmod 600 "$rendered_config"

exec /usr/local/bin/attestor "$@" --config "$rendered_config"
