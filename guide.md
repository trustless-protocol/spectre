```
# prerequisites (kurtosis, gaiad)
cd relayer
LD_LIBRARY_PATH=$HOME/works/ecip-gnark go run ./prover/cmd ./bin ../contracts/verifiers

# Build the relayer binary
go build -o relayer ./cmd

./run_eth_node.sh
# Poll until finalized.epoch > 0:
curl -s http://127.0.0.1:32774/eth/v1/beacon/states/head/finality_checkpoints

./run_cosmos_node.sh 
# wait to cosmos up
./wasm.sh

# Copies the ICS07 address back into relayer/config.json automatically.
LD_LIBRARY_PATH=$HOME/works/ecip-gnark ./relayer create-clients \
  --config config.example.json \
  --wasm-checksum d24688886ed8cec00c667fa69c173fbab9a08c75900ce18afe10517c82e55592

# copy config.example.json to config.json
LD_LIBRARY_PATH=$HOME/works/ecip-gnark ./relayer start --config config.example.json

# wait to success

./scripts/eth_to_cosmos_transfer.sh

# catch done on relayer log
./scripts/check_eth_data.sh
./scripts/check_cosmos_data.sh
```