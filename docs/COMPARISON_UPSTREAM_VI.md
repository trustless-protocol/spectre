# So Sánh Kiến Trúc: Fast-IBC vs Upstream solidity-ibc-eureka

## 1. Tổng Quan

| Tiêu chí | Fast-IBC (Decentrio) | Upstream (cosmos/solidity-ibc-eureka) |
|-----------|---------------------|--------------------------------------|
| Ngôn ngữ chính | Solidity + Go + Rust (CosmWasm) | Solidity + Rust + TypeScript |
| ZK proving | gnark Groth16 (Go, Ed25519) | SP1 zkVM (Rust, RISC-V) |
| Relayer | Go monolithic, event-driven | Rust stateless, gRPC operator |
| Light client | 2 loại (Groth16, CW ETH) | 4+ loại (SP1, Attestation, Solana TM, CW ETH) |
| Mục tiêu | Ethereum ↔ Cosmos (IBC v2) | Đa chuỗi (ETH, Solana, Cosmos) |

## 2. Khác Biệt Cốt Lõi: Hệ Thống Chứng Minh ZK

### 2.1. Fast-IBC: gnark Groth16

```
Chữ ký Ed25519 (R ∥ S, 64 bytes)
        │
        ▼
Giải nén điểm Ed25519 → toạ độ Weierstrass (aX, aY, rX, rY)
Tính H = SHA512(R ∥ A ∥ msg) off-chain
        │
        ▼
PreHashCircuit (gnark) → groth16.Prove() → proof [8]*big.Int
        │
        ▼
On-chain: WrapperVerifier → Groth16Verifier (BN254 pairing)
```

**Đặc điểm:**
- Chỉ chứng minh **một phép tính**: xác minh chữ ký Ed25519 trong mạch ZK
- Phần còn lại (giải nén điểm, SHA512, reduce mod L) thực hiện **on-chain** trong `WrapperVerifier.sol`
- Prover chạy **in-process** trong Go relayer, không cần dịch vụ ngoài
- Mạch cố định (PreHashCircuit), trusted setup một lần
- Chi phí gas cao hơn do tính toán on-chain nhiều (~100k gas cho SHA512)

### 2.2. Upstream: SP1 zkVM

```
Toàn bộ logic Tendermint Light Client
  ├── Xác minh chữ ký Ed25519
  ├── Kiểm tra voting power ≥ 2/3
  ├── Xác minh header chain
  └── Tính Merkle root
        │
        ▼
SP1 Program (Rust, biên dịch sang RISC-V) → SP1 Prover → proof
        │
        ▼
On-chain: SP1ICS07Tendermint → SP1Verifier (precompile/contract)
```

**Đặc điểm:**
- Chứng minh **toàn bộ logic** light client trong zkVM
- On-chain chỉ cần xác minh proof duy nhất, chi phí gas thấp hơn
- Prover là dịch vụ **riêng biệt** (SP1 network hoặc local prover)
- Linh hoạt: thay đổi logic light client chỉ cần sửa Rust program, không cần thay mạch
- Thời gian proving lâu hơn (zkVM overhead)

### 2.3. Bảng So Sánh Chi Tiết

| Tiêu chí | Fast-IBC (Groth16) | Upstream (SP1) |
|-----------|-------------------|----------------|
| Phạm vi ZK | Chỉ Ed25519 signature | Toàn bộ LC logic |
| Tính toán on-chain | Nhiều (SHA512, decompress, reduce) | Ít (chỉ verify proof) |
| Gas cost ước tính | Cao hơn (~300-500k gas) | Thấp hơn (~200-300k gas) |
| Proving time | Nhanh hơn (mạch nhỏ) | Chậm hơn (zkVM overhead) |
| Trusted setup | Cần (Groth16 ceremony) | Không cần (STARK-based) |
| Thay đổi logic | Phải sửa mạch + contract | Chỉ sửa Rust program |
| Dependency | gnark (Go native) | SP1 SDK + network |
| Toolchain | Go thuần | Rust + RISC-V toolchain |

## 3. Kiến Trúc Smart Contract

### 3.1. Hệ Thống Contract Chung

Cả hai dự án đều triển khai IBC v2 với cấu trúc tương tự:

```
                Fast-IBC                          Upstream
                ────────                          ────────
ICS26Router (UUPS)                    ICS26Router (UUPS)
  ├── ICS20Transfer (UUPS)              ├── ICS20Transfer (UUPS)
  │   ├── IBCERC20 (Beacon)             │   ├── IBCERC20 (Beacon)
  │   └── Escrow (Beacon)               │   └── Escrow (Beacon)
  └── Groth16ICS07Tendermint            ├── SP1ICS07Tendermint
                                        ├── AttestationLightClient (mới)
                                        ├── ICS27GMP (mới)
                                        └── IFT (Interchain Fungible Token)
```

### 3.2. Khác Biệt Contract

| Thành phần | Fast-IBC | Upstream |
|------------|----------|----------|
| Light client TM | `Groth16ICS07Tendermint` (4 contract con) | `SP1ICS07Tendermint` (1 verifier) |
| Xác minh chữ ký | `WrapperVerifier` + `Groth16Verifier` | SP1Verifier (đơn giản) |
| Membership proof | `Membership.sol` riêng | Tích hợp trong SP1 program |
| Misbehaviour | `Misbehaviour.sol` riêng | Tích hợp trong SP1 program |
| UpdateClient | `UpdateClient.sol` riêng | Tích hợp trong SP1 program |
| GMP (General Message Passing) | Chưa có | `ICS27GMP.sol` |
| IFT (Interchain Fungible Token) | Chưa có | Có |
| Attestation LC | Chưa có | `AttestationLightClient.sol` |

### 3.3. Phân Tách Verification Logic

**Fast-IBC** tách logic xác minh thành 4 contract riêng biệt:
- `Membership.sol`: Xác minh Merkle proof (ICS-23)
- `Misbehaviour.sol`: Phát hiện validator gian lận
- `UpdateClient.sol`: Xác nhận header Tendermint mới
- `WrapperVerifier.sol`: Chuyển đổi Ed25519 → BN254

**Upstream** gộp tất cả logic vào SP1 program (off-chain), on-chain chỉ cần:
- `SP1Verifier`: Xác minh proof duy nhất
- Merkle, misbehaviour, update đều nằm trong Rust program

→ **Nhận xét**: Fast-IBC có bề mặt tấn công on-chain lớn hơn nhưng dễ audit từng phần. Upstream có bề mặt on-chain nhỏ nhưng logic ZK phức tạp hơn.

## 4. Kiến Trúc Relayer

### 4.1. Fast-IBC: Go Monolithic

```
relayer/
├── cmd/main.go           ← CLI đơn giản
├── services/
│   ├── main.go           ← Vòng lặp chính, 4 goroutine
│   ├── routine.go        ← UpdateCosmosClient, UpdateEthClient
│   ├── batch.go          ← BatchBuilder (flush by size/time)
│   └── context.go        ← Shared state với mutex
├── subscriber/           ← WebSocket + Event filter
├── transaction/          ← Ký và gửi tx cả hai chuỗi
├── prover/               ← gnark in-process prover
└── client/               ← RPC clients
```

**Đặc điểm:**
- **Monolithic**: Tất cả trong 1 process Go
- **Stateful**: `Context` struct giữ trạng thái, mutex bảo vệ
- **Event-driven**: WebSocket (Cosmos) + Log filter (ETH)
- **In-process prover**: gnark chạy cùng process, không cần gRPC
- **Batch processing**: Gom packet theo kích thước (10) hoặc thời gian (3s)

### 4.2. Upstream: Rust Stateless Operator

```
programs/operator/
├── src/
│   ├── main.rs           ← gRPC server
│   ├── cli/              ← Rich CLI (clap)
│   ├── operator/
│   │   ├── cosmos_to_eth.rs
│   │   ├── eth_to_cosmos.rs
│   │   └── relay.rs      ← Relay logic
│   ├── clients/          ← Chain-specific clients
│   └── prover/           ← SP1 prover client (remote)
├── proto/                ← Protobuf definitions
└── Cargo.toml
```

**Đặc điểm:**
- **Stateless**: Mỗi relay call là request độc lập
- **gRPC service**: Có thể scale horizontally
- **Remote prover**: SP1 prover chạy riêng (SP1 network hoặc local)
- **Multi-chain**: Hỗ trợ Cosmos, Ethereum, Solana
- **Rich CLI**: clap với nhiều subcommand

### 4.3. So Sánh Relayer

| Tiêu chí | Fast-IBC (Go) | Upstream (Rust) |
|-----------|--------------|----------------|
| Ngôn ngữ | Go | Rust |
| Kiến trúc | Monolithic, stateful | Stateless, gRPC |
| Prover | In-process (gnark) | Remote (SP1 network) |
| Scaling | Vertical (1 process) | Horizontal (nhiều instance) |
| Event handling | Push (WebSocket/filter) | Pull (polling) |
| Batch support | Có (BatchBuilder) | Không rõ |
| Chuỗi hỗ trợ | Cosmos ↔ ETH | Cosmos ↔ ETH ↔ Solana |
| Complexity | Đơn giản hơn | Phức tạp hơn |
| Deployment | Binary đơn | Binary + SP1 prover service |

## 5. Light Client

### 5.1. Tendermint Light Client trên Ethereum

| Tiêu chí | Fast-IBC | Upstream |
|-----------|----------|----------|
| Tên contract | `Groth16ICS07Tendermint` | `SP1ICS07Tendermint` |
| Loại proof | Groth16 (BN254) | SP1 proof (STARK-based) |
| Số contract con | 4 (Wrapper, Groth16, Membership, Misbehaviour) | 1 (SP1Verifier) |
| Merkle verification | On-chain (Membership.sol) | Off-chain (trong SP1 program) |
| Update client | On-chain (UpdateClient.sol) | Off-chain (trong SP1 program) |
| Trusted setup | Cần | Không cần |

### 5.2. Ethereum Light Client trên Cosmos

Cả hai đều dùng **CosmWasm contract** (08-wasm) với logic tương tự:
- Xác minh Beacon Chain sync committee signature (BLS)
- Lưu trạng thái: latest slot, execution block number, sync committee
- Xác minh ETH storage proof (Merkle Patricia Trie)

### 5.3. Light Client Bổ Sung (Chỉ Upstream)

| Light Client | Mô tả |
|-------------|--------|
| `AttestationLightClient` | Multi-sig attestation, không cần ZK |
| Solana Tendermint LC | Tendermint LC trên Solana (SVM program) |

## 6. Hệ Thống Quyền và Bảo Mật

### 6.1. Access Control

| Vai trò | Fast-IBC | Upstream |
|---------|----------|----------|
| ADMIN_ROLE (0) | ✅ | ✅ |
| RELAYER_ROLE (1) | ✅ | ✅ |
| PAUSER_ROLE (2) | ✅ | ✅ |
| UNPAUSER_ROLE (3) | ✅ | ✅ |
| RATE_LIMITER_ROLE (5) | ✅ | ✅ |
| IBC_STORE_ROLE | ❌ | ✅ |
| CLIENT_UPDATER_ROLE | ❌ | ✅ |
| CUSTOM_MIGRATOR_ROLE | ❌ | ✅ |
| TIMELOCK_ROLE | ❌ | ✅ |

→ Upstream có hệ thống phân quyền chi tiết hơn (9 vai trò vs 5 vai trò).

### 6.2. Cơ Chế Bảo Vệ

| Cơ chế | Fast-IBC | Upstream |
|--------|----------|----------|
| Rate limiting | ✅ (Escrow) | ✅ (Escrow + IBCERC20) |
| Pause/Unpause | ✅ | ✅ |
| Reentrancy guard | ✅ (Transient) | ✅ (Transient) |
| Misbehaviour detection | ✅ | ✅ |
| Chống replay | ✅ | ✅ |
| Timelock | ❌ | ✅ |

## 7. Tính Năng và Module

| Tính năng | Fast-IBC | Upstream |
|-----------|----------|----------|
| ICS-20 Token Transfer | ✅ | ✅ |
| ICS-26 Router | ✅ | ✅ |
| ICS-07 Tendermint LC | ✅ (Groth16) | ✅ (SP1) |
| ICS-27 GMP | ❌ | ✅ |
| IFT (Interchain Fungible Token) | ❌ | ✅ |
| Attestation Light Client | ❌ | ✅ |
| Solana Support | ❌ | ✅ |
| E2E Testing | ✅ (interchaintest) | ✅ (interchaintest) |
| OpenTelemetry | ❌ | ✅ |
| Multicall batching | ❌ | ✅ |

## 8. Ưu Điểm và Nhược Điểm

### 8.1. Fast-IBC

**Ưu điểm:**
- **Go-native**: Toàn bộ relayer + prover viết bằng Go, không cần Rust toolchain
- **In-process prover**: Không cần dịch vụ prover riêng, triển khai đơn giản
- **Proving time nhanh**: Mạch Groth16 nhỏ (chỉ Ed25519), tạo proof nhanh
- **Dễ debug**: Monolithic process, dễ theo dõi luồng dữ liệu
- **Ít dependency**: Không phụ thuộc SP1 network hay Rust toolchain

**Nhược điểm:**
- **Gas cost cao hơn**: SHA512, decompress, reduce mod L thực hiện on-chain
- **Trusted setup**: Groth16 yêu cầu ceremony, keys phải bảo mật
- **Bề mặt tấn công on-chain lớn**: 4 contract xác minh vs 1 verifier
- **Khó mở rộng logic**: Thay đổi mạch ZK cần trusted setup mới
- **Thiếu tính năng**: Chưa có GMP, IFT, Solana, Attestation LC
- **Scalability hạn chế**: Monolithic, không scale horizontal

### 8.2. Upstream (solidity-ibc-eureka)

**Ưu điểm:**
- **Gas cost thấp**: Hầu hết logic trong ZK, on-chain chỉ verify proof
- **Không cần trusted setup**: SP1 dùng STARK (transparent)
- **Linh hoạt**: Thay đổi logic LC chỉ cần sửa Rust program
- **Đa chuỗi**: Hỗ trợ Cosmos, Ethereum, Solana
- **Tính năng phong phú**: GMP, IFT, Attestation, OpenTelemetry
- **Scale horizontal**: Stateless operator + remote prover

**Nhược điểm:**
- **Toolchain phức tạp**: Cần Rust + SP1 SDK + prover service
- **Proving time lâu**: zkVM overhead cho toàn bộ LC logic
- **Dependency nặng**: SP1 network, Rust toolchain, nhiều crate
- **Triển khai phức tạp**: Nhiều service cần chạy đồng thời
- **Debug khó hơn**: Logic ZK trong RISC-V VM, khó inspect

## 9. Tóm Tắt Quyết Định Kiến Trúc

```
                    Fast-IBC                    Upstream
                    ────────                    ────────
Triết lý:          "Đơn giản, Go-native"       "Toàn diện, đa chuỗi"

ZK scope:          Nhỏ (chỉ Ed25519)           Lớn (toàn bộ LC logic)
On-chain logic:    Nhiều                        Ít
Off-chain logic:   Ít                           Nhiều

Deployment:        1 binary                     Nhiều service
Complexity:        Thấp                         Cao
Feature set:       Core IBC                     Core IBC + GMP + IFT + Solana

Gas efficiency:    Thấp hơn                     Cao hơn
Proving speed:     Nhanh hơn                    Chậm hơn
Trust assumption:  Trusted setup                Transparent (STARK)
```

**Kết luận**: Fast-IBC ưu tiên sự đơn giản trong triển khai và vận hành, phù hợp cho cầu nối Cosmos ↔ Ethereum chuyên dụng. Upstream ưu tiên tính tổng quát và hiệu quả gas, phù hợp cho hạ tầng IBC đa chuỗi quy mô lớn. Sự khác biệt cốt lõi nằm ở phạm vi ZK proof: Fast-IBC chứng minh ít hơn trong ZK nhưng đổi lại tính toán on-chain nhiều hơn, trong khi upstream đẩy hầu hết logic vào zkVM để giảm thiểu chi phí on-chain.
