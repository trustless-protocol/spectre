#!/bin/sh

set -euxo pipefail

killall gaiad || true
rm -rf $HOME/.gaia

# Set up Gaia node
CHAIN_ID="test-ibc-eth"
gaiad init test-ibc --chain-id $CHAIN_ID
gaiad keys add test --keyring-backend test
gaiad keys add test1 --keyring-backend test
gaiad keys add test2 --keyring-backend test
gaiad keys add test3 --keyring-backend test

cat $HOME/.gaia/config/genesis.json | jq '.app_state["gov"]["params"]["voting_period"]="30s"' > $HOME/.gaia/config/tmp_genesis.json && mv $HOME/.gaia/config/tmp_genesis.json $HOME/.gaia/config/genesis.json
cat $HOME/.gaia/config/genesis.json | jq '.app_state["gov"]["params"]["expedited_voting_period"]="20s"' > $HOME/.gaia/config/tmp_genesis.json && mv $HOME/.gaia/config/tmp_genesis.json $HOME/.gaia/config/genesis.json
cat $HOME/.gaia/config/genesis.json | jq '.app_state["feemarket"]["params"]["max_block_utilization"]="300000000"' > $HOME/.gaia/config/tmp_genesis.json && mv $HOME/.gaia/config/tmp_genesis.json $HOME/.gaia/config/genesis.json
sed -i'' -e "s/^minimum-gas-prices *= .*/minimum-gas-prices = \"1stake\"/" "$HOME/.gaia/config/app.toml"

gaiad genesis add-genesis-account test 1100000000000stake --keyring-backend test
gaiad genesis add-genesis-account test1 2000000000000stake --keyring-backend test
gaiad genesis add-genesis-account test2 2000000000000stake --keyring-backend test
gaiad genesis add-genesis-account test3 200000000000stake --keyring-backend test

GAIA_TEST_ADDRESS=$(gaiad keys show test -a --keyring-backend test)
GAIA_TEST1_ADDRESS=$(gaiad keys show test1 -a --keyring-backend test)
gaiad genesis gentx test 1000000000000stake --chain-id $CHAIN_ID --keyring-backend test
gaiad genesis collect-gentxs
gaiad genesis validate-genesis

gaiad keys export test1 --unarmored-hex --unsafe --keyring-backend test -y

gaiad start 

