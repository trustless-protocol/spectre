#!/usr/bin/env bash
#
# Spin up N gaiad validator nodes on a single host (default N=180).
#
# Usage:
#   ./run_cosmos_node_n.sh                # 180 validators
#   NUM_NODES=64 ./run_cosmos_node_n.sh   # custom count
#
# Notes:
# - All homes live under $GAIA_BASE (default $HOME/.gaia-multi); the directory
#   is wiped on every run so you start from a clean genesis.
# - Ports are allocated with a stride of 1 per node from configurable bases.
#   With N=180 the highest port reached is ~37179. Make sure your OS file
#   descriptor limit is high (e.g. `ulimit -n 65535`).
# - The relayer .env is updated with the FIRST validator's RPC/API so existing
#   tooling keeps working; treat node 0 as the "primary".

set -euo pipefail

NUM_NODES="${NUM_NODES:-180}"
CHAIN_ID="${CHAIN_ID:-test-ibc-eth}"
DENOM="${DENOM:-stake}"
KEYRING="${KEYRING:-test}"
GAIA_BASE="${GAIA_BASE:-$HOME/.gaia-multi}"
RELAYER_ENV_FILE="${RELAYER_ENV_FILE:-relayer/.env}"

# Port bases (each node consumes base+i for its index i in [0, NUM_NODES)).
P2P_BASE="${P2P_BASE:-30000}"
RPC_BASE="${RPC_BASE:-31000}"
PROXY_BASE="${PROXY_BASE:-32000}"
API_BASE="${API_BASE:-33000}"
GRPC_BASE="${GRPC_BASE:-34000}"
GRPC_WEB_BASE="${GRPC_WEB_BASE:-35000}"
PPROF_BASE="${PPROF_BASE:-36000}"
PROM_BASE="${PROM_BASE:-37000}"

# Voting-power distribution mimicking Cosmos Hub mainnet: top ~20 validators
# hold ~2/3 of total voting power. Tier boundaries + per-tier stakes are env-
# configurable so you can reshape the curve (e.g. for stress tests).
#
# Default tiers (N=180):
#   ranks [0, 20)   -> 33.5T stake each  -> 20 * 33.5T  = 670T  (~67%)
#   ranks [20, 60)  -> 6T    stake each  -> 40 * 6T     = 240T  (~24%)
#   ranks [60,180)  -> 0.75T stake each  -> 120 * 0.75T = 90T   (~9%)
#                                                 total = 1000T
TOP_TIER_END="${TOP_TIER_END:-20}"
MID_TIER_END="${MID_TIER_END:-60}"
TOP_TIER_STAKE="${TOP_TIER_STAKE:-33500000000000}"
MID_TIER_STAKE="${MID_TIER_STAKE:-6000000000000}"
TAIL_TIER_STAKE="${TAIL_TIER_STAKE:-750000000000}"

# Buffer added to each validator's genesis balance on top of its self-stake
# so the account can still pay fees after delegating into its own validator.
GENESIS_BAL_BUFFER="${GENESIS_BAL_BUFFER:-1000000000}"

val_stake_amount() {
    local i="$1"
    if [ "$i" -lt "$TOP_TIER_END" ]; then
        printf '%s' "$TOP_TIER_STAKE"
    elif [ "$i" -lt "$MID_TIER_END" ]; then
        printf '%s' "$MID_TIER_STAKE"
    else
        printf '%s' "$TAIL_TIER_STAKE"
    fi
}

val_stake() {
    printf '%s%s' "$(val_stake_amount "$1")" "$DENOM"
}

val_genesis_bal() {
    local stake
    stake="$(val_stake_amount "$1")"
    printf '%s%s' "$((stake + GENESIS_BAL_BUFFER))" "$DENOM"
}

USER_KEYS=(test test1 test2 test3)
USER_BALANCES=(
    "1100000000000${DENOM}"
    "2000000000000${DENOM}"
    "2000000000000${DENOM}"
    "200000000000${DENOM}"
)

if ! command -v gaiad >/dev/null 2>&1; then
    echo "gaiad not found in PATH" >&2
    exit 1
fi

echo "Stopping any existing gaiad processes..."
killall gaiad 2>/dev/null || true

echo "Wiping $GAIA_BASE..."
rm -rf "$GAIA_BASE"
mkdir -p "$GAIA_BASE"

# Build arrays of per-node values.
HOMES=()
VAL_KEYS=()
P2P_PORTS=()
RPC_PORTS=()
PROXY_PORTS=()
API_PORTS=()
GRPC_PORTS=()
GRPC_WEB_PORTS=()
PPROF_PORTS=()
PROM_PORTS=()

for ((i = 0; i < NUM_NODES; i++)); do
    HOMES+=("$GAIA_BASE/val$i")
    VAL_KEYS+=("val$i")
    P2P_PORTS+=("$((P2P_BASE + i))")
    RPC_PORTS+=("$((RPC_BASE + i))")
    PROXY_PORTS+=("$((PROXY_BASE + i))")
    API_PORTS+=("$((API_BASE + i))")
    GRPC_PORTS+=("$((GRPC_BASE + i))")
    GRPC_WEB_PORTS+=("$((GRPC_WEB_BASE + i))")
    PPROF_PORTS+=("$((PPROF_BASE + i))")
    PROM_PORTS+=("$((PROM_BASE + i))")
done

PRIMARY_HOME="${HOMES[0]}"

upsert_env_var() {
    local file="$1" key="$2" value="$3" tmp
    mkdir -p "$(dirname "$file")"
    touch "$file"
    tmp="$(mktemp)"
    awk -v key="$key" -v value="$value" '
        BEGIN { updated = 0 }
        $0 ~ ("^" key "=") { print key "=\"" value "\""; updated = 1; next }
        { print }
        END { if (!updated) print key "=\"" value "\"" }
    ' "$file" > "$tmp"
    mv "$tmp" "$file"
}

echo "Initializing $NUM_NODES validator homes..."
for i in "${!HOMES[@]}"; do
    gaiad init "${VAL_KEYS[$i]}" --chain-id "$CHAIN_ID" --home "${HOMES[$i]}" >/dev/null 2>&1
    gaiad keys add "${VAL_KEYS[$i]}" --keyring-backend "$KEYRING" --home "${HOMES[$i]}" >/dev/null 2>&1
done

echo "Creating user keys on primary home..."
for i in "${!USER_KEYS[@]}"; do
    gaiad keys add "${USER_KEYS[$i]}" --keyring-backend "$KEYRING" --home "$PRIMARY_HOME" >/dev/null 2>&1
done

# Genesis tweaks (governance + feemarket).
tmp_genesis() {
    local f="$PRIMARY_HOME/config/genesis.json"
    jq "$1" "$f" > "$f.tmp" && mv "$f.tmp" "$f"
}
tmp_genesis '.app_state["gov"]["params"]["voting_period"]="30s"'
tmp_genesis '.app_state["gov"]["params"]["expedited_voting_period"]="20s"'
tmp_genesis '.app_state["feemarket"]["params"]["max_block_utilization"]="300000000"'

echo "Adding genesis accounts..."
for i in "${!USER_KEYS[@]}"; do
    gaiad genesis add-genesis-account "${USER_KEYS[$i]}" "${USER_BALANCES[$i]}" \
        --keyring-backend "$KEYRING" --home "$PRIMARY_HOME" >/dev/null
done
for i in "${!VAL_KEYS[@]}"; do
    val_addr=$(gaiad keys show "${VAL_KEYS[$i]}" -a --keyring-backend "$KEYRING" --home "${HOMES[$i]}")
    gaiad genesis add-genesis-account "$val_addr" "$(val_genesis_bal "$i")" --home "$PRIMARY_HOME" >/dev/null
done

# Sync the in-progress genesis to every node so gentx can sign against it.
for ((i = 1; i < NUM_NODES; i++)); do
    cp "$PRIMARY_HOME/config/genesis.json" "${HOMES[$i]}/config/genesis.json"
done

echo "Generating gentx for each validator (tiered stakes)..."
for i in "${!HOMES[@]}"; do
    gaiad genesis gentx "${VAL_KEYS[$i]}" "$(val_stake "$i")" \
        --chain-id "$CHAIN_ID" \
        --keyring-backend "$KEYRING" \
        --home "${HOMES[$i]}" >/dev/null 2>&1
done

# Sanity-check that 2/3 voting power is concentrated as intended.
top_total=$((TOP_TIER_END * TOP_TIER_STAKE))
mid_total=$(((MID_TIER_END - TOP_TIER_END) * MID_TIER_STAKE))
tail_total=$(((NUM_NODES - MID_TIER_END) * TAIL_TIER_STAKE))
grand_total=$((top_total + mid_total + tail_total))
if [ "$grand_total" -gt 0 ]; then
    top_pct=$((top_total * 100 / grand_total))
    mid_pct=$((mid_total * 100 / grand_total))
    tail_pct=$((tail_total * 100 / grand_total))
    echo "Stake distribution: top($TOP_TIER_END)=${top_pct}% mid($((MID_TIER_END - TOP_TIER_END)))=${mid_pct}% tail($((NUM_NODES - MID_TIER_END)))=${tail_pct}%"
fi

mkdir -p "$PRIMARY_HOME/config/gentx"
for ((i = 1; i < NUM_NODES; i++)); do
    cp "${HOMES[$i]}"/config/gentx/*.json "$PRIMARY_HOME/config/gentx/"
done

gaiad genesis collect-gentxs --home "$PRIMARY_HOME" >/dev/null
gaiad genesis validate-genesis --home "$PRIMARY_HOME"

# Distribute final genesis.
for ((i = 1; i < NUM_NODES; i++)); do
    cp "$PRIMARY_HOME/config/genesis.json" "${HOMES[$i]}/config/genesis.json"
done

echo "Collecting node ids..."
NODE_IDS=()
for i in "${!HOMES[@]}"; do
    NODE_IDS+=("$(gaiad tendermint show-node-id --home "${HOMES[$i]}")")
done

# Each node uses up to PEERS_PER_NODE neighbors (rotating window) to avoid an
# O(N^2) persistent_peers string. Set PEERS_PER_NODE=0 to disable the cap and
# include every other node (only feasible for small N).
PEERS_PER_NODE="${PEERS_PER_NODE:-20}"

build_peers() {
    local idx="$1"
    local peers=""
    local count=0
    local limit="$PEERS_PER_NODE"
    if [ "$limit" -le 0 ] || [ "$limit" -ge "$NUM_NODES" ]; then
        limit=$((NUM_NODES - 1))
    fi
    for ((k = 1; k <= limit; k++)); do
        local j=$(((idx + k) % NUM_NODES))
        local peer="${NODE_IDS[$j]}@127.0.0.1:${P2P_PORTS[$j]}"
        if [ -z "$peers" ]; then peers="$peer"; else peers="$peers,$peer"; fi
        count=$((count + 1))
        if [ "$count" -ge "$limit" ]; then break; fi
    done
    printf '%s' "$peers"
}

configure_node() {
    local idx="$1"
    local home="${HOMES[$idx]}"
    local peers
    peers="$(build_peers "$idx")"

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

echo "Configuring node configs..."
for i in "${!HOMES[@]}"; do
    configure_node "$i"
done

COSMOS_PRIVATE_KEY="$(gaiad keys export test1 --unarmored-hex --unsafe --keyring-backend "$KEYRING" --home "$PRIMARY_HOME" -y 2>/dev/null)"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_PRIVATE_KEY" "$COSMOS_PRIVATE_KEY"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_CHAIN_ID" "$CHAIN_ID"
upsert_env_var "$RELAYER_ENV_FILE" "COSMOS_RPC" "tcp://127.0.0.1:${RPC_PORTS[0]}"
echo "Updated $RELAYER_ENV_FILE (primary RPC tcp://127.0.0.1:${RPC_PORTS[0]})"

PIDS=()
echo "Starting $NUM_NODES validators..."
for i in "${!HOMES[@]}"; do
    gaiad start --home "${HOMES[$i]}" > "${HOMES[$i]}/gaiad.log" 2>&1 &
    PIDS+=("$!")
done

cleanup() {
    echo "Shutting down ${#PIDS[@]} validators..."
    for pid in "${PIDS[@]}"; do
        kill "$pid" 2>/dev/null || true
    done
}
trap cleanup EXIT INT TERM

cat <<EOF
Started $NUM_NODES Gaia validators.
  base dir       : $GAIA_BASE
  primary rpc    : tcp://127.0.0.1:${RPC_PORTS[0]}
  primary p2p    : tcp://127.0.0.1:${P2P_PORTS[0]}
  last  rpc      : tcp://127.0.0.1:${RPC_PORTS[$((NUM_NODES - 1))]}
  peers per node : $PEERS_PER_NODE
  stake tiers    : top[0,$TOP_TIER_END)=$TOP_TIER_STAKE mid[$TOP_TIER_END,$MID_TIER_END)=$MID_TIER_STAKE tail[$MID_TIER_END,$NUM_NODES)=$TAIL_TIER_STAKE
Logs: \$GAIA_BASE/val<i>/gaiad.log
EOF

wait "${PIDS[@]}"
