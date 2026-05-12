#!/usr/bin/env bash

set -euo pipefail
set -x

REPO_ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"

killall gaiad || true
rm -rf "$HOME/.gaia" "$HOME/.gaia-val2" "$HOME/.gaia-val3" "$HOME/.gaia-val4"

CHAIN_ID="test-ibc-eth"
DENOM="stake"
KEYRING="test"

HOMES=(
    "$HOME/.gaia"
    "$HOME/.gaia-val2"
    "$HOME/.gaia-val3"
    "$HOME/.gaia-val4"
)
VAL_KEYS=(val1 val2 val3 val4)
VAL_STAKES=(
    "1000000000000${DENOM}"
    "1000000000000${DENOM}"
    "1000000000000${DENOM}"
    "1000000000000${DENOM}"
)
USER_KEYS=(test test1 test2 test3)
USER_BALANCES=(
    "1100000000000${DENOM}"
    "2000000000000${DENOM}"
    "2000000000000${DENOM}"
    "200000000000${DENOM}"
)

P2P_PORTS=(26656 26756 26856 26956)
RPC_PORTS=(26657 26757 26857 26957)
PROXY_PORTS=(26658 26758 26858 26958)
API_PORTS=(1317 1327 1337 1347)
GRPC_PORTS=(9090 9100 9110 9120)
GRPC_WEB_PORTS=(9091 9101 9111 9121)
PPROF_PORTS=(6060 6070 6080 6090)
PROM_PORTS=(26660 26760 26860 26960)

PRIMARY_HOME="${HOMES[0]}"
RELAYER_ENV_FILE="${RELAYER_ENV_FILE:-$REPO_ROOT/relayer/.env}"

upsert_env_var() {
    local file="$1"
    local key="$2"
    local value="$3"
    local tmp

    mkdir -p "$(dirname "$file")"
    touch "$file"
    tmp="$(mktemp)"

    awk -v key="$key" -v value="$value" '
        BEGIN { updated = 0 }
        $0 ~ ("^" key "=") {
            print key "=\"" value "\""
            updated = 1
            next
        }
        { print }
        END {
            if (!updated) {
                print key "=\"" value "\""
            }
        }
    ' "$file" > "$tmp"

    mv "$tmp" "$file"
}

for i in "${!HOMES[@]}"; do
    gaiad init "${VAL_KEYS[$i]}" --chain-id "$CHAIN_ID" --home "${HOMES[$i]}"
    gaiad keys add "${VAL_KEYS[$i]}" --keyring-backend "$KEYRING" --home "${HOMES[$i]}"
done

# Keep the old local testing accounts available in the default Gaia home.
for i in "${!USER_KEYS[@]}"; do
    gaiad keys add "${USER_KEYS[$i]}" --keyring-backend "$KEYRING" --home "$PRIMARY_HOME"
done

jq '.app_state["gov"]["params"]["voting_period"]="30s"' "$PRIMARY_HOME/config/genesis.json" > "$PRIMARY_HOME/config/tmp_genesis.json"
mv "$PRIMARY_HOME/config/tmp_genesis.json" "$PRIMARY_HOME/config/genesis.json"
jq '.app_state["gov"]["params"]["expedited_voting_period"]="20s"' "$PRIMARY_HOME/config/genesis.json" > "$PRIMARY_HOME/config/tmp_genesis.json"
mv "$PRIMARY_HOME/config/tmp_genesis.json" "$PRIMARY_HOME/config/genesis.json"
jq '.app_state["feemarket"]["params"]["max_block_utilization"]="300000000"' "$PRIMARY_HOME/config/genesis.json" > "$PRIMARY_HOME/config/tmp_genesis.json"
mv "$PRIMARY_HOME/config/tmp_genesis.json" "$PRIMARY_HOME/config/genesis.json"

for i in "${!USER_KEYS[@]}"; do
    gaiad genesis add-genesis-account "${USER_KEYS[$i]}" "${USER_BALANCES[$i]}" --keyring-backend "$KEYRING" --home "$PRIMARY_HOME"
done

for i in "${!VAL_KEYS[@]}"; do
    val_addr=$(gaiad keys show "${VAL_KEYS[$i]}" -a --keyring-backend "$KEYRING" --home "${HOMES[$i]}")
    gaiad genesis add-genesis-account "$val_addr" "2000000000000${DENOM}" --home "$PRIMARY_HOME"
done

for i in 1 2 3; do
    cp "$PRIMARY_HOME/config/genesis.json" "${HOMES[$i]}/config/genesis.json"
done

for i in "${!HOMES[@]}"; do
    gaiad genesis gentx "${VAL_KEYS[$i]}" "${VAL_STAKES[$i]}" \
        --chain-id "$CHAIN_ID" \
        --keyring-backend "$KEYRING" \
        --home "${HOMES[$i]}"
done

mkdir -p "$PRIMARY_HOME/config/gentx"
for i in 1 2 3; do
    cp "${HOMES[$i]}"/config/gentx/*.json "$PRIMARY_HOME/config/gentx/"
done

gaiad genesis collect-gentxs --home "$PRIMARY_HOME"
gaiad genesis validate-genesis --home "$PRIMARY_HOME"

for i in 1 2 3; do
    cp "$PRIMARY_HOME/config/genesis.json" "${HOMES[$i]}/config/genesis.json"
done

NODE_IDS=()
for i in "${!HOMES[@]}"; do
    NODE_IDS+=("$(gaiad tendermint show-node-id --home "${HOMES[$i]}")")
done

configure_node() {
    local idx="$1"
    local home="${HOMES[$idx]}"
    local peers=""

    for j in "${!HOMES[@]}"; do
        if [ "$j" != "$idx" ]; then
            peer="${NODE_IDS[$j]}@127.0.0.1:${P2P_PORTS[$j]}"
            if [ -z "$peers" ]; then
                peers="$peer"
            else
                peers="$peers,$peer"
            fi
        fi
    done

    sed -i'' -e "s#^proxy_app = .*#proxy_app = \"tcp://127.0.0.1:${PROXY_PORTS[$idx]}\"#" "$home/config/config.toml"
    sed -i'' -e "/^\[rpc\]/,/^\[/{s#^laddr = .*#laddr = \"tcp://127.0.0.1:${RPC_PORTS[$idx]}\"#;}" "$home/config/config.toml"
    sed -i'' -e "/^\[p2p\]/,/^\[/{s#^laddr = .*#laddr = \"tcp://0.0.0.0:${P2P_PORTS[$idx]}\"#; s#^persistent_peers = .*#persistent_peers = \"$peers\"#; s#^addr_book_strict = .*#addr_book_strict = false#; s#^allow_duplicate_ip = .*#allow_duplicate_ip = true#;}" "$home/config/config.toml"
    sed -i'' -e "s#^pprof_laddr = .*#pprof_laddr = \"localhost:${PPROF_PORTS[$idx]}\"#" "$home/config/config.toml"
    sed -i'' -e "s#^prometheus_listen_addr = .*#prometheus_listen_addr = \":${PROM_PORTS[$idx]}\"#" "$home/config/config.toml"
    sed -i'' -e "s/^minimum-gas-prices *= .*/minimum-gas-prices = \"1${DENOM}\"/" "$home/config/app.toml"
    sed -i'' -e "/^\[api\]/,/^\[/{s#^address = .*#address = \"tcp://0.0.0.0:${API_PORTS[$idx]}\"#;}" "$home/config/app.toml"
    sed -i'' -e "/^\[grpc\]/,/^\[/{s#^address = .*#address = \"0.0.0.0:${GRPC_PORTS[$idx]}\"#;}" "$home/config/app.toml"
    sed -i'' -e "/^\[grpc-web\]/,/^\[/{s#^address = .*#address = \"0.0.0.0:${GRPC_WEB_PORTS[$idx]}\"#;}" "$home/config/app.toml"
}

for i in "${!HOMES[@]}"; do
    configure_node "$i"
done

COSMOS_PRIVATE_KEY="$(gaiad keys export test1 --unarmored-hex --unsafe --keyring-backend "$KEYRING" --home "$PRIMARY_HOME" -y)"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_PRIVATE_KEY" "$COSMOS_PRIVATE_KEY"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_CHAIN_ID" "$CHAIN_ID"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_GAS_LIMIT" "500000"
echo "  Written to $RELAYER_ENV_FILE: COSMOS_PRIVATE_KEY, COSMOS_CHAIN_ID, COSMOS_GAS_LIMIT"

PIDS=()
for i in "${!HOMES[@]}"; do
    gaiad start --home "${HOMES[$i]}" > "${HOMES[$i]}/gaiad.log" 2>&1 &
    PIDS+=("$!")
done

cleanup() {
    for pid in "${PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
    done
}
trap cleanup EXIT INT TERM

echo ""
echo "━━━ Cosmos Validators Started ━━━"
for i in "${!HOMES[@]}"; do
    printf "  %-6s  rpc=tcp://127.0.0.1:%-5s  p2p=tcp://127.0.0.1:%-5s  home=%s\n" \
        "${VAL_KEYS[$i]}" "${RPC_PORTS[$i]}" "${P2P_PORTS[$i]}" "${HOMES[$i]}"
done
echo "  (logs: <home>/gaiad.log)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

wait "${PIDS[@]}"
