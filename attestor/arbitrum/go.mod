module attestor/arbitrum

go 1.25.7

require (
	attestor/core v0.0.0
	attestor/grpc v0.0.0
	attestor/host v0.0.0
	attestor/types v0.0.0
	github.com/ethereum/go-ethereum v1.15.5
	github.com/spf13/cobra v1.10.2
	golang.org/x/sync v0.20.0
	google.golang.org/grpc v1.82.1
)

replace attestor/types => ../types

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/ProjectZKM/Ziren/crates/go-runtime/zkvm_runtime v0.0.0-20251001021608-1fe7b43fc4d6 // indirect
	github.com/StackExchange/wmi v1.2.1 // indirect
	github.com/bits-and-blooms/bitset v1.20.0 // indirect
	github.com/consensys/gnark-crypto v0.18.1 // indirect
	github.com/crate-crypto/go-eth-kzg v1.4.0 // indirect
	github.com/crate-crypto/go-ipa v0.0.0-20240724233137-53bbb0ceb27a // indirect
	github.com/deckarep/golang-set/v2 v2.6.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/ethereum/c-kzg-4844/v2 v2.1.5 // indirect
	github.com/ethereum/go-verkle v0.2.2 // indirect
	github.com/gammazero/deque v1.1.0 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/holiman/uint256 v1.3.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/shirou/gopsutil v3.21.4-0.20210419000835-c7a38de76ee5+incompatible // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	github.com/supranational/blst v0.3.16-0.20250831170142-f48500c1fdbe // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace attestor/core => ../core

replace attestor/grpc => ../grpc

replace attestor/host => ../host

// Nitro's geth fork retains github.com/ethereum/go-ethereum as its declared
// module path, so the fork is selected with a module replacement.
replace github.com/ethereum/go-ethereum => github.com/offchainlabs/go-ethereum v0.0.0-20260616211651-0f618f330b8d
