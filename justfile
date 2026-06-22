set dotenv-load

# Default task lists all available tasks
default:
  just --list

# Build the contracts using `forge build`
[group('build')]
build-contracts: clean-foundry
	forge build

# Build the relayer using `cargo build`
[group('build')]
build-relayer:
	cargo build --bin relayer --release --locked

# Build the relayer using `go build`
[group('build')]
build-go-relayer:
	cd relayer && go build ./...

# Build and optimize the eth wasm light client using `cosmwasm/optimizer`. Requires `docker` and `gzip`
[group('build')]
build-cw-ics08-wasm-eth:
	docker run --rm -v "$(pwd)":/code --mount type=volume,source="$(basename "$(pwd)")_cache",target=/target --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry cosmwasm/optimizer:0.17.0 ./programs/cw-ics08-wasm-eth
	cp artifacts/cw_ics08_wasm_eth.wasm e2e/interchaintestv8/wasm
	gzip -n e2e/interchaintestv8/wasm/cw_ics08_wasm_eth.wasm -f

# Build the relayer docker image
[group('build')]
build-relayer-image:
    docker build -t eureka-relayer:latest -f programs/relayer/Dockerfile .

# Install the Go relayer for use in the e2e tests.
# Builds with an explicit -o name because the package dir is `cmd/`, so a plain
# `go install ./cmd/...` would land at $GOPATH/bin/cmd, not $GOPATH/bin/relayer.
[group('install')]
install-go-relayer:
	cd relayer && go build -o $(go env GOPATH)/bin/relayer ./cmd

# Generate per-bucket prover artifacts (r1cs/pk/vk + Solidity verifiers) if missing.
# Probes for the bucket-4 vk.bin AND the bucket-4 Solidity verifier — the cheapest
# "present?" signal because bucket 4 is the smallest and is always built first.
[group('build')]
build-prover-artifacts:
	@if [ ! -f relayer/bin/n4/vk.bin ] || [ ! -f contracts/verifiers/Groth16Verifier_N4.sol ]; then \
		echo "Building prover artifacts..."; \
		cd relayer && go run ./prover/cmd ./bin ../contracts/verifiers ; \
	else \
		echo "Prover artifacts already present, skipping setup"; \
	fi

# Run all linters
[group('lint')]
lint:
	@echo "Running all linters..."
	just lint-solidity
	just lint-go
	just lint-buf
	just lint-rust

# Lint the Solidity code using `forge fmt` and `bun:solhint`
[group('lint')]
lint-solidity:
	@echo "Linting the Solidity code..."
	forge fmt --check
	bun solhint -w 0 '{scripts,contracts,test}/**/*.sol'
	natlint run --include 'contracts/**/*.sol'

# Lint the Go code using `golangci-lint`
[group('lint')]
lint-go:
	@echo "Linting the Go code..."
	cd relayer && golangci-lint run
	cd e2e/interchaintestv8 && golangci-lint run
	cd packages/go-abigen && golangci-lint run

# Lint the Protobuf files using `buf lint`
[group('lint')]
lint-buf:
	@echo "Linting the Protobuf files..."
	buf lint

# Lint the Rust code using `cargo fmt` and `cargo clippy`
[group('lint')]
lint-rust:
	@echo "Linting the Rust code..."
	cargo fmt --all -- --check
	cargo clippy --all-targets --all-features -- -D warnings


# Generate the (non-bytecode) ABI files for the contracts
[group('generate')]
generate-abi: build-contracts
	jq '.abi' out/ICS26Router.sol/ICS26Router.json > abi/ICS26Router.json
	jq '.abi' out/ICS20Transfer.sol/ICS20Transfer.json > abi/ICS20Transfer.json
	jq '.abi' out/Groth16ICS07Tendermint.sol/Groth16ICS07Tendermint.json > abi/Groth16ICS07Tendermint.json
	jq '.abi' out/ERC20.sol/ERC20.json > abi/ERC20.json
	jq '.abi' out/IBCERC20.sol/IBCERC20.json > abi/IBCERC20.json
	jq '.abi' out/RelayerHelper.sol/RelayerHelper.json > abi/RelayerHelper.json
	abigen --abi abi/ERC20.json --pkg erc20 --type Contract --out e2e/interchaintestv8/types/erc20/contract.go
	abigen --abi abi/Groth16ICS07Tendermint.json --pkg groth16ics07tendermint --type Contract --out packages/go-abigen/groth16ics07tendermint/contract.go
	abigen --abi abi/ICS20Transfer.json --pkg ics20transfer --type Contract --out packages/go-abigen/ics20transfer/contract.go
	abigen --abi abi/ICS26Router.json --pkg ics26router --type Contract --out packages/go-abigen/ics26router/contract.go
	abigen --abi abi/IBCERC20.json --pkg ibcerc20 --type Contract --out packages/go-abigen/ibcerc20/contract.go
	abigen --abi abi/RelayerHelper.json --pkg relayerhelper --type Contract --out packages/go-abigen/relayerhelper/contract.go

# Generate the ABI files with bytecode for the required contracts (only Groth16ICS07Tendermint)
[group('generate')]
generate-abi-bytecode: build-contracts
	cp out/Groth16ICS07Tendermint.sol/Groth16ICS07Tendermint.json abi/bytecode

# Generate the fixtures for the wasm tests using the e2e tests
[group('generate')]
generate-fixtures-wasm: clean-foundry install-go-relayer build-prover-artifacts
	@echo "Generating fixtures... This may take a while."
	@echo "Generating recvPacket and acknowledgePacket groth16 fixtures..."
	cd e2e/interchaintestv8 && ETH_TESTNET_TYPE=pos GENERATE_WASM_FIXTURES=true E2E_PROOF_TYPE=groth16 go test -v -run '^TestWithIbcEurekaTestSuite/Test_ICS20TransferERC20TokenfromEthereumToCosmosAndBack$' -timeout 60m
	@echo "Generating native SdkCoin recvPacket groth16 fixtures..."
	cd e2e/interchaintestv8 && ETH_TESTNET_TYPE=pos GENERATE_WASM_FIXTURES=true E2E_PROOF_TYPE=groth16 go test -v -run '^TestWithIbcEurekaTestSuite/Test_ICS20TransferNativeCosmosCoinsToEthereumAndBack$' -timeout 60m
	@echo "Generating timeoutPacket groth16 fixtures..."
	cd e2e/interchaintestv8 && ETH_TESTNET_TYPE=pos GENERATE_WASM_FIXTURES=true E2E_PROOF_TYPE=groth16 go test -v -run '^TestWithIbcEurekaTestSuite/Test_TimeoutPacketFromCosmos$' -timeout 60m

# Generate go types for the e2e tests from the ethereum light client code
[group('generate')]
generate-ethereum-types:
	cargo run --bin generate_json_schema --features test-utils
	bun quicktype --src-lang schema --lang go --just-types-and-package --package ethereum --src ethereum_types_schema.json --out e2e/interchaintestv8/types/ethereum/types.gen.go --top-level GeneratedTypes
	rm ethereum_types_schema.json
	sed -i.bak 's/int64/uint64/g' e2e/interchaintestv8/types/ethereum/types.gen.go # quicktype generates int64 instead of uint64 :(
	rm -f e2e/interchaintestv8/types/ethereum/types.gen.go.bak # this is to be linux and mac compatible (coming from the sed command)
	cd e2e/interchaintestv8 && golangci-lint run --fix types/ethereum/types.gen.go

# Generate the code from pritibuf using `buf generate`. (Only used for relayer testing at the moment)
[group('generate')]
generate-buf:
    @echo "Generating Protobuf files for relayer"
    buf generate --template buf.gen.yaml

shadowfork := if env("ETH_RPC_URL", "") == "" { "--no-match-path test/shadowfork/*" } else { "" }

# Run all the foundry tests
[group('test')]
test-foundry testname=".\\*":
	forge test -vvv --show-progress --fuzz-runs 5000 --match-test ^{{testname}}\(.\*\)\$ {{shadowfork}}
	@ {{ if shadowfork == "" { "" } else { 'echo ' + BOLD + YELLOW + 'Ran without shadowfork tests since ETH_RPC_URL was not set' } }}

# Run the benchmark tests
[group('test')]
test-benchmark testname=".\\*":
	forge test -vvv --show-progress --gas-report --match-path test/solidity-ibc/BenchmarkTest.t.sol --match-test {{testname}}

# Run the cargo tests
[group('test')]
test-cargo testname="--all":
	cargo test {{testname}} --locked --no-fail-fast -- --nocapture

# Run the tests in abigen
[group('test')]
test-abigen:
	@echo "Running abigen tests..."
	cd packages/go-abigen && go test -v ./...

# Run Go relayer tests
[group('test')]
test-go-relayer:
	@echo "Running Go relayer tests..."
	cd relayer && go test -v ./...

# Run any e2e test using the test's full name. For example, `just test-e2e TestWithIbcEurekaTestSuite/Test_Deploy`
#
# ETH_TESTNET_TYPE=pos picks the Kurtosis PoS network (with beacon API), which is what
# the wasm light client + the cosmos→eth direction need. Override to "pow" only if the
# test explicitly targets the anvil/PoW path.
[group('test')]
test-e2e testname: clean-foundry install-go-relayer build-prover-artifacts
	@echo "Running {{testname}} test..."
	cd e2e/interchaintestv8 && \
		RELAYER_BINARY="$(go env GOPATH)/bin/relayer" \
		PROVER_BIN_DIR="$(git rev-parse --show-toplevel)/relayer/bin" \
		ETH_TESTNET_TYPE=pos \
		go test -v -run '^{{testname}}$' -timeout 120m

# Run any e2e test in the IbcEurekaTestSuite. For example, `just test-e2e-eureka Test_Deploy`
[group('test')]
test-e2e-eureka testname:
	@echo "Running {{testname}} test..."
	just test-e2e TestWithIbcEurekaTestSuite/{{testname}}

# Run any e2e test in the RelayerTestSuite. For example, `just test-e2e-relayer Test_RelayerInfo`
[group('test')]
test-e2e-relayer testname:
	@echo "Running {{testname}} test..."
	just test-e2e TestWithRelayerTestSuite/{{testname}}

# Run any e2e test in the MultichainTestSuite. For example, `just test-e2e-multichain Test_Deploy`
[group('test')]
test-e2e-multichain testname:
	@echo "Running {{testname}} test..."
	just test-e2e TestWithMultichainTestSuite/{{testname}}

# Clean up the foundry cache and out directories
[group('clean')]
clean-foundry:
	@echo "Cleaning up cache and out directories"
	-rm -rf cache out broadcast # ignore errors

# Clean up the cargo artifacts using `cargo clean`
[group('clean')]
clean-cargo:
	@echo "Cleaning up cargo target directory"
	cargo clean

# Run Slither static analysis on contracts
[group('security')]
slither:
	@echo "Running Slither static analysis..."
	slither . --config-file .slither.config.json
