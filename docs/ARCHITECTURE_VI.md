# Kiến Trúc Tổng Quan Hệ Thống (High-Level Architecture)

## 1. Tổng Quan

Fast-IBC là hệ thống triển khai giao thức **IBC v2** để kết nối Ethereum và Cosmos, sử dụng **Groth16 zero-knowledge proof** để xác minh chữ ký validator Tendermint trên Ethereum. Hệ thống gồm 3 tầng ngôn ngữ: Solidity (smart contract), Go (relayer/operator), và Rust (CosmWasm light client).

```
┌─────────────────────────────────────────────────────────────────────┐
│                        KIẾN TRÚC HỆ THỐNG                          │
│                                                                     │
│  ┌──────────────┐          ┌──────────────┐         ┌────────────┐ │
│  │ Cosmos Chain  │◄── IBC ──►│  Go Relayer  │◄── IBC ──►│  Ethereum  │ │
│  │  (CometBFT)  │  packets  │  (Operator)  │  packets  │   (EVM)   │ │
│  │              │          │              │          │            │ │
│  │ ┌──────────┐ │          │ ┌──────────┐ │          │ ┌────────┐ │ │
│  │ │08-wasm   │ │  Beacon  │ │ Groth16  │ │  ZK      │ │ICS26   │ │ │
│  │ │ETH Light │◄├──header──┤ │ Prover   │ ├──proof──►│ │Router  │ │ │
│  │ │Client    │ │          │ │ (gnark)  │ │          │ │        │ │ │
│  │ └──────────┘ │          │ └──────────┘ │          │ │ICS07   │ │ │
│  └──────────────┘          └──────────────┘          │ │ICS20   │ │ │
│                                                       │ └────────┘ │ │
│                                                       └────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

## 2. Ba Thành Phần Chính

### 2.1. Smart Contract trên Ethereum

Hệ thống contract xử lý toàn bộ logic IBC phía Ethereum:

```
ICS26Router (UUPS Proxy) ← Điểm vào chính cho mọi gói tin IBC
  ├─ ICS20Transfer (UUPS Proxy) ← Cầu nối token (ICS-20)
  │   ├─ IBCERC20 (Beacon Proxy) ← Token ERC20 đại diện tài sản từ Cosmos
  │   └─ Escrow (Beacon Proxy) ← Giữ token gốc khi chuyển đi
  └─ Groth16ICS07Tendermint ← Light client xác minh trạng thái Cosmos
      ├─ WrapperVerifier ← Chuyển đổi Ed25519 → BN254 cho ZK
      ├─ Groth16Verifier ← Xác minh Groth16 proof on-chain
      ├─ Membership ← Xác minh Merkle proof (ICS-23)
      ├─ Misbehaviour ← Phát hiện validator gian lận
      └─ UpdateClient ← Xác nhận header Tendermint mới
```

**Vai trò từng contract:**

| Contract | Proxy | Vai trò |
|----------|-------|---------|
| `ICS26Router` | UUPS | Định tuyến packet IBC, lưu commitment, quản lý light client |
| `ICS20Transfer` | UUPS | Xử lý chuyển token: lock/unlock, mint/burn |
| `IBCERC20` | Beacon | Token ERC20 wrap cho tài sản từ chuỗi khác |
| `Escrow` | Beacon | Giữ token gốc, có rate limiting chống mint vô hạn |
| `Groth16ICS07Tendermint` | Không | Light client Tendermint, lưu consensus state |
| `WrapperVerifier` | Không | Giải nén Ed25519, tính SHA512, encode cho Groth16 |
| `Groth16Verifier` | Không | Xác minh proof bằng pairing trên BN254 |

### 2.2. Go Relayer (Operator)

Relayer là dịch vụ Go chạy liên tục, kết nối hai chuỗi:

```
relayer/
├── cmd/main.go         ← CLI: start, create-clients, genesis, fixtures
├── services/           ← Vòng lặp chính, batch builder, worker
├── subscriber/         ← Lắng nghe sự kiện Cosmos (WebSocket) + ETH (log filter)
├── transaction/        ← Ký và gửi giao dịch lên cả hai chuỗi
├── client/             ← Truy vấn RPC Cosmos + Ethereum + Beacon API
├── prover/             ← Sinh Groth16 proof (gnark, Ed25519)
├── keys/               ← Quản lý private key
├── utils/              ← Tiện ích: IBC path, byte conversion
└── bindings/           ← Go binding tự động sinh từ ABI Solidity
```

**Chức năng chính:**
- Lắng nghe sự kiện `send_packet`, `write_acknowledgement` từ Cosmos
- Lắng nghe sự kiện `SendPacket`, `WriteAcknowledgement` từ Ethereum
- Cập nhật light client hai chiều (Cosmos trên ETH, ETH trên Cosmos)
- Sinh ZK proof cho chữ ký validator
- Tạo và gửi Merkle proof cho packet commitment/acknowledgement
- Gom batch để xử lý hiệu quả

### 2.3. Wasm Light Client trên Cosmos

Light client Ethereum chạy trên Cosmos dưới dạng **CosmWasm contract** (kiểu `08-wasm`):
- Xác minh Beacon Chain sync committee signature
- Lưu trạng thái: `LatestSlot`, `LatestExecutionBlockNumber`, sync committee
- Nhận `MsgUpdateClient` chứa Beacon header + sync aggregate
- Xác minh Ethereum storage proof (Merkle Patricia Trie) cho acknowledgement

## 3. Mô Hình Tin Cậy (Trust Model)

```
                    Cosmos Chain                    Ethereum Chain
                    ─────────────                   ──────────────
Đồng thuận:        CometBFT (≥2/3 VP)             PoS + Beacon (≥2/3 sync committee)
Light client:      08-wasm (ETH LC trên Cosmos)    Groth16ICS07 (Cosmos LC trên ETH)
Bằng chứng:        ETH storage proof (MPT)          Cosmos Merkle proof (IAVL + ICS23)
Xác minh chữ ký:   Sync committee BLS              Ed25519 → Groth16 ZK proof
```

**Các giả định bảo mật:**
1. **Validator set trung thực**: ≥ 2/3 voting power không gian lận
2. **Groth16 circuit đúng đắn**: Mạch `PreHashCircuit` xác minh chính xác Ed25519
3. **Trusted setup an toàn**: Proving key và verifying key không bị thao túng
4. **Relayer không cần tin tưởng**: Chỉ chuyển tiếp, không thể giả mạo proof

## 4. Luồng Dữ Liệu Chính

### 4.1. Cosmos → Ethereum (Gửi token)

```
User (Cosmos)                 Relayer                      Ethereum
     │                           │                             │
     │── MsgTransfer ──►         │                             │
     │   (send_packet event)     │                             │
     │                           │◄── WebSocket ───            │
     │                           │                             │
     │                           │── UpdateClient ──────────►  │
     │                           │   (ZK proof + header)       │
     │                           │                             │
     │                           │── RecvPacket ────────────►  │
     │                           │   (Merkle proof + packet)   │
     │                           │                         mint/unlock token
     │                           │                             │
     │                           │   (WriteAck event) ◄────────│
     │                           │                             │
     │    ◄── MsgAcknowledgement │                             │
     │        (ETH storage proof)│                             │
     │                           │                             │
   Hoàn tất                                                Hoàn tất
```

### 4.2. Ethereum → Cosmos (Gửi token)

```
User (Ethereum)               Relayer                      Cosmos
     │                           │                             │
     │── sendTransfer() ──►      │                             │
     │   (SendPacket event)      │                             │
     │                           │◄── Event filter ────        │
     │                           │                             │
     │                           │── MsgUpdateClient ───────►  │
     │                           │   (Beacon header)           │
     │                           │                             │
     │                           │── MsgRecvPacket ─────────►  │
     │                           │   (ETH storage proof)       │
     │                           │                         mint/unlock token
     │                           │                             │
     │                           │   (write_ack event) ◄───────│
     │                           │                             │
     │    ◄── AckPacket ─────────│                             │
     │        (Cosmos Merkle proof)                            │
     │                           │                             │
   Hoàn tất                                                Hoàn tất
```

## 5. Chuỗi Xác Minh ZK Proof

Điểm khác biệt chính của fast-ibc so với upstream `solidity-ibc-eureka`: thay vì dùng SP1 (RISC-V zkVM), hệ thống dùng **gnark Groth16** để xác minh chữ ký Ed25519.

```
                         Off-chain (Go Relayer)
┌─────────────────────────────────────────────────────────┐
│  Chữ ký Ed25519 (R ∥ S, 64 bytes)                      │
│  Public key (A, 32 bytes)                               │
│  Message (vote sign bytes)                              │
│         │                                               │
│         ▼                                               │
│  Giải nén Ed25519 → tọa độ Weierstrass (aX, aY, rX, rY)│
│  Tính H = SHA512(R ∥ A ∥ msg) off-chain                │
│         │                                               │
│         ▼                                               │
│  Xây dựng witness cho PreHashCircuit                    │
│  groth16.Prove(r1cs, pk, witness) → proof               │
│  Xác minh local: groth16.Verify(proof, vk, pubWitness) │
└──────────────────────────┬──────────────────────────────┘
                           │ gửi lên Ethereum
                           ▼
                    On-chain (Solidity)
┌─────────────────────────────────────────────────────────┐
│  WrapperVerifier.verifyProof():                         │
│    1. Giải nén Ed25519 point (Montgomery ladder)        │
│    2. Tính SHA512 on-chain (~100k gas)                  │
│    3. Reduce 512-bit hash mod L (Ed25519 scalar order)  │
│    4. Encode thành 24 field limbs (6 × 4 limbs 64-bit) │
│         │                                               │
│         ▼                                               │
│  Groth16Verifier.verifyProof():                         │
│    - BN254 pairing check (precompile 0x08)              │
│    - ECAdd, ECMul (precompile 0x06, 0x07)               │
│    → true/false                                          │
└─────────────────────────────────────────────────────────┘
```

## 6. Quyền Truy Cập và Bảo Mật

### Vai trò trên Ethereum (AccessManager)

| Vai trò | Giá trị | Quyền hạn |
|---------|---------|-----------|
| `ADMIN_ROLE` | 0 | Nâng cấp contract, quản trị hệ thống |
| `RELAYER_ROLE` | 1 | Gửi recvPacket, ackPacket, timeoutPacket, updateClient |
| `PAUSER_ROLE` | 2 | Tạm dừng khẩn cấp toàn bộ chuyển token |
| `UNPAUSER_ROLE` | 3 | Khôi phục hoạt động sau tạm dừng |
| `RATE_LIMITER_ROLE` | 5 | Cấu hình giới hạn rút token hàng ngày |

### Cơ chế bảo vệ

- **Chống replay**: Lưu packet receipt, gửi lại là no-op
- **Chống double-spend**: Commitment storage + sequence number
- **Chống mint vô hạn**: Rate limiting trên Escrow và IBCERC20
- **Phát hiện gian lận**: Misbehaviour contract đóng băng light client
- **Chống reentrancy**: `ReentrancyGuardTransientUpgradeable`
- **Tạm dừng khẩn cấp**: Pause/unpause trên ICS20Transfer

## 7. Cấu Hình Triển Khai

### File cấu hình JSON

```json
{
  "modules": [
    {
      "name": "cosmos_to_eth",
      "config": {
        "tm_rpc_url": "http://localhost:26657",
        "ics26_address": "0x...",
        "eth_rpc_url": "http://localhost:8545",
        "eth_ws_url": "ws://localhost:8546",
        "ics07_client": "0x..."
      }
    },
    {
      "name": "eth_to_cosmos",
      "config": {
        "eth_beacon_api_url": "http://localhost:5052",
        "signer_address": "cosmos1..."
      }
    }
  ]
}
```

### Biến môi trường (.env)

| Biến | Mục đích |
|------|----------|
| `ETH_PRIVATE_KEY` | Khóa riêng để ký giao dịch Ethereum |
| `COSMOS_PRIVATE_KEY` | Khóa riêng để ký giao dịch Cosmos |
| `COSMOS_CHAIN_ID` | Chain ID của mạng Cosmos |
| `PROVER_R1CS_PATH` | Đường dẫn tới constraint system (6.1MB) |
| `PROVER_PK_PATH` | Đường dẫn tới proving key (34MB) |
| `PROVER_VK_PATH` | Đường dẫn tới verifying key (1.4KB) |

## 8. Yêu Cầu Hạ Tầng

| Thành phần | Yêu cầu |
|------------|---------|
| Go Relayer | Go 1.21+, RAM ≥ 4GB (cho proof generation) |
| Cosmos Node | CometBFT RPC + WebSocket (port 26657) |
| Ethereum Node | JSON-RPC + WebSocket (port 8545/8546) |
| Beacon Node | REST API (port 5052) |
| Smart Contract | Foundry, Bun (không dùng npm/yarn) |
| E2E Testing | Docker + Kurtosis |
