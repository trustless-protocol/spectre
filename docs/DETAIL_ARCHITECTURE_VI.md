# Kiến Trúc Chi Tiết Relayer (Detail-Level Architecture)

## 1. Khởi Động Hệ Thống

### 1.1. CLI Entry Point

**File:** `relayer/cmd/main.go`

Relayer sử dụng framework **Cobra** với 4 lệnh chính:

| Lệnh | Mục đích | Khi nào dùng |
|-------|----------|-------------|
| `start` | Chạy vòng lặp relay chính | Vận hành liên tục |
| `create-clients` | Triển khai light client hai chiều | Một lần khi thiết lập |
| `genesis` | Tạo trạng thái genesis cho client mới | Khởi tạo ban đầu |
| `fixtures membership` | Kiểm tra Merkle proof | Debug/test |

### 1.2. Trình Tự Khởi Động (Lệnh `start`)

```
1. godotenv.Load()
   └─ Đọc .env: private key, đường dẫn prover, chain ID

2. loadConfig(configPath)
   └─ Đọc JSON config → CosmosToEthConfig + EthToCosmosConfig

3. Kết nối Ethereum
   ├─ ethclient.Dial(eth_rpc_url)       → ethClient (HTTP, cho giao dịch)
   └─ ethclient.Dial(eth_ws_url)        → ethWsClient (WebSocket, cho event)

4. Kết nối Cosmos
   └─ rpchttp.New(tm_rpc_url)           → cosmosClient (HTTP + WebSocket)

5. Khởi tạo Prover
   └─ prover.NewProver(r1csPath, pkPath, vkPath)
      ├─ Đọc r1cs.bin (6.1MB) → constraint system
      ├─ Đọc pk.bin (34MB) → proving key
      └─ Đọc vk.bin (1.4KB) → verifying key

6. Tạo Context
   └─ services.NewCtxWithBeacon(cosmosClient, ethClient, ethWsClient, beaconURL, ethClientID)
      ├─ Khởi tạo Logger
      ├─ Tạo 2 Timestamp struct (cho ETH và Cosmos)
      └─ Lưu tham chiếu tới tất cả client

7. Cấu hình địa chỉ contract
   └─ ctx.SetAddresses(ics26, verifier, membership, misbehaviour, updateClient, roleManager)
      └─ ctx.SetClient(ics07Address)

8. Khởi động Cosmos WebSocket
   └─ cosmosClient.WSEvents.Start()

9. Tạo Services và chạy
   └─ services.New(subscriber, txHandler, prover, ethConfig, cosmosConfig)
      └─ svc.StartLoop(ctx)
```

## 2. Kiến Trúc Service Layer

### 2.1. Sơ Đồ Phụ Thuộc (Dependency Graph)

```
┌─────────────────────────────────────────────────┐
│                   Services                       │
│  ┌─────────────┐  ┌────────────┐  ┌──────────┐ │
│  │ EventListener│  │   Worker   │  │BatchBuilder│ │
│  │ (interface)  │  │            │  │          │ │
│  └──────┬───────┘  │ ┌────────┐│  └────┬─────┘ │
│         │          │ │TxHandler││       │       │
│  ┌──────▼───────┐  │ │(interf.)││  ┌────▼─────┐ │
│  │  Subscriber  │  │ └────┬───┘│  │ []Packet │ │
│  │  - Cosmos WS │  │ ┌────▼───┐│  │ sync.Mtx │ │
│  │  - ETH filter│  │ │Prover  ││  └──────────┘ │
│  └──────────────┘  │ │(interf.)││               │
│                    │ └────┬───┘│               │
│                    └──────┼────┘               │
│                    ┌──────▼────┐                │
│                    │Handler    │                │
│                    │EcipProver │                │
│                    └───────────┘                │
└─────────────────────────────────────────────────┘
```

### 2.2. Interfaces

**File:** `relayer/services/main.go:35-55`

```go
type TransactionHandler interface {
    CreateCosmosClientContract(ctx Context, clientState, consensusHash []byte) error
    CreateEthClient(ctx Context, clientState, consensusState exported.ClientState) (string, error)
    SendEthTx(ctx Context, msg any) error
    SendCosmosTx(ctx Context, msg any) error
    SendCosmosTxBatch(ctx Context, msgs []any) error
    CosmosSignerAddress() (string, error)
}

type Prover interface {
    GenerateProof(sig, pub, msg []byte) (
        proof [8]*big.Int,
        commitments [2]*big.Int,
        commitmentPok [2]*big.Int,
        err error,
    )
}

type EventListener interface {
    SubscribeCosmos(ctx Context, batchBuilder *BatchBuilder)
    SubscribeEth(ctx Context, batchBuilder *BatchBuilder)
}
```

### 2.3. Context Object

**File:** `relayer/services/context.go`

`Context` là đối tượng trung tâm được truyền qua toàn bộ hệ thống:

```
Context
├── Logger *log.Logger              ← Ghi log chuẩn
├── Config Config                   ← Cấu hình runtime
│
├── latestEthTimestamp *Timestamp   ← Theo dõi cập nhật Cosmos LC trên ETH
│   ├── mtx sync.Mutex             ← Bảo vệ đọc/ghi đồng thời
│   ├── LatestUpdateTime time.Time  ← Thời điểm cập nhật cuối
│   └── LatestUpdateHeight uint64   ← Block height cuối cùng đã cập nhật
│
├── latestCosmosTimestamp *Timestamp ← Theo dõi cập nhật ETH LC trên Cosmos
│
├── cosmosClient *rpchttp.HTTP      ← RPC + WebSocket tới Cosmos node
├── ethClient *ethclient.Client     ← HTTP RPC tới Ethereum node
├── ethWsClient *ethclient.Client   ← WebSocket tới Ethereum node
├── beaconAPIURL string             ← REST API tới Beacon node
├── ethClientID string              ← ID của ETH light client trên Cosmos
│
└── Địa chỉ contract (tất cả *common.Address):
    ├── ics26Router                 ← ICS26Router trên Ethereum
    ├── ics07Client                 ← Groth16ICS07Tendermint trên Ethereum
    ├── verifier                    ← WrapperVerifier
    ├── membership                  ← Contract xác minh Merkle proof
    ├── misbehaviour                ← Contract phát hiện gian lận
    ├── updateClient                ← Contract xác nhận header
    └── roleManager                 ← Quản lý quyền (address(0) = ai cũng được)
```

### 2.4. Cấu Hình Runtime

**File:** `relayer/services/config.go`

```
Config
├── KeyPath string                  ← (chưa được sử dụng)
├── IntervalParams IntervalConfig
│   ├── blockHeight uint8           ← Mặc định: 5
│   └── blockTime time.Duration     ← Mặc định: 20 giây
├── IntervalType IntervalType       ← "blockHeight" hoặc "timestamp"
└── BatchConfig BatchConfig
    ├── BatchSize uint8             ← Mặc định: 10 packet/batch
    └── BatchPeriods time.Duration  ← Mặc định: 3 giây
```

## 3. Mô Hình Đồng Thời (Goroutine Model)

### 3.1. Tổng Quan Goroutine

`StartLoop()` khởi tạo 4 goroutine chạy song song + 1 vòng lặp chính:

```
┌─ Goroutine 1: SubscribeCosmos ──────────────────────────────┐
│  Lắng nghe 3 loại sự kiện Cosmos qua WebSocket              │
│  → InsertPacket vào BatchBuilder                             │
└──────────────────────────────────────────────────────────────┘

┌─ Goroutine 2: SubscribeEth ─────────────────────────────────┐
│  Lắng nghe 4 loại event từ ICS26Router qua WebSocket ETH    │
│  → InsertPacket vào BatchBuilder                             │
└──────────────────────────────────────────────────────────────┘

┌─ Goroutine 3: Routine Scheduler ────────────────────────────┐
│  Mỗi giây kiểm tra:                                         │
│  - Cosmos LC trên ETH cần cập nhật? (mỗi 24 giờ)            │
│  - ETH LC trên Cosmos cần cập nhật? (mỗi 24 giờ)            │
└──────────────────────────────────────────────────────────────┘

┌─ Goroutine 4: Batch Checker ────────────────────────────────┐
│  Mỗi 3 giây kiểm tra BatchBuilder:                          │
│  - Đã đủ số lượng (> BatchSize)? → flush                     │
│  - Đã quá thời gian (> BatchPeriods)? → flush                │
│  → Gửi BatchPackets vào channel                              │
└──────────────────────────────────────────────────────────────┘

┌─ Main Loop: Packet Handler ─────────────────────────────────┐
│  Chờ nhận batch từ channel                                   │
│  → Chờ 6 giây (AppHash cập nhật)                             │
│  → UpdateCosmosClient (ZK proof)                             │
│  → Xử lý từng packet theo loại: Send, Ack, Timeout, WriteAck│
└──────────────────────────────────────────────────────────────┘
```

### 3.2. Đồng Bộ Hóa (Synchronization)

| Cơ chế | Vị trí | Bảo vệ |
|--------|--------|--------|
| `sync.Mutex` | `BatchBuilder.mtx` | Đọc/ghi danh sách packet |
| `sync.Mutex` | `Timestamp.mtx` | Đọc/ghi thời gian cập nhật LC |
| Channel | `Services.BatchPackets` | Giao tiếp giữa Batch Checker và Main Loop |
| Channel (CometBFT) | `WSEvents.Subscribe` | Nhận sự kiện từ Cosmos |
| Channel (go-ethereum) | `WatchOpts` | Nhận event từ Ethereum |

### 3.3. Luồng Dữ Liệu Giữa Các Goroutine

```
SubscribeCosmos ──┐                                    ┌── Main Loop
                  ├── InsertPacket ──► BatchBuilder ──►│    (xử lý)
SubscribeEth ─────┘      (mutex)          │            └──────────
                                          │
                              CheckBatch ──┘
                              (mỗi 3 giây)
                                  │
                                  ▼
                           channel BatchPackets
                                  │
                                  ▼
                             Main Loop
                          (blocking receive)
```

## 4. Hệ Thống Sự Kiện (Event Subscription)

### 4.1. Sự Kiện Cosmos

**File:** `relayer/subscriber/event.go:33`

| Sự kiện | Bộ lọc CometBFT | Loại Packet | Trường dữ liệu |
|---------|-----------------|-------------|-----------------|
| `send_packet` | `message.action='/ibc.applications.transfer.v1.MsgTransfer'` | `Send` | `send_packet.encoded_packet_hex` |
| `write_acknowledgement` | `message.action='/ibc.core.channel.v2.MsgRecvPacket'` | `Ack` | `write_acknowledgement.encoded_packet_hex` + `encoded_acknowledgement_hex` |
| `timeout_packet` | `message.action='/ibc.applications.transfer.v1.MsgTimeout'` | `Timeout` | `timeout_packet.encoded_packet_hex` |

**Xử lý sự kiện:**
1. Trích xuất trường hex-encoded từ `event.Events`
2. Giải mã hex → bytes: `hex.DecodeString()`
3. Unmarshal protobuf: `proto.Unmarshal() → channeltypesv2.Packet`
4. Gọi `batchBuilder.InsertPacket()` với loại tương ứng

### 4.2. Sự Kiện Ethereum

**File:** `relayer/subscriber/event.go:180`

| Event | Channel | Loại Packet | Dữ liệu đặc biệt |
|-------|---------|-------------|-------------------|
| `SendPacket` | `sendPacketCh` | `Send` | Packet struct từ Solidity |
| `WriteAcknowledgement` | `writeAckCh` | `WriteAck` | + `Acknowledgements` + `BlockNumber` |
| `AckPacket` | `ackPacketCh` | _(chưa xử lý)_ | Chỉ log |
| `TimeoutPacket` | `timeoutPacketCh` | _(chưa xử lý)_ | Chỉ log |

**Chuyển đổi format:**
```
EthPacketToCosmosPacket():
  Solidity struct                    →  Cosmos protobuf
  ────────────────                       ──────────────
  ev.Packet.SourceClient             →  packet.SourceClient
  ev.Packet.DestClient               →  packet.DestinationClient  (tên trường khác)
  ev.Packet.Payloads[i].DestPort     →  payload.DestinationPort   (tên trường khác)
  ev.Sequence (big.Int)              →  packet.Sequence (uint64)
```

## 5. Batch Builder

### 5.1. Cấu Trúc

**File:** `relayer/services/batch.go`

```
BatchBuilder
├── mtx sync.Mutex        ← Bảo vệ truy cập đồng thời
├── timestamp time.Time   ← Thời điểm batch bắt đầu gom
└── packets []Packet      ← Danh sách packet chờ xử lý

Packet
├── PacketType PacketType          ← Send(0), Ack(1), Timeout(2), WriteAck(3)
├── Packet *channeltypesv2.Packet  ← Dữ liệu IBC packet
├── AckBytes [][]byte              ← Acknowledgement data (cho Ack và WriteAck)
└── BlockNumber uint64             ← ETH block (cho WriteAck, cần chờ finality)
```

### 5.2. Điều Kiện Flush Batch

```
CheckBatch(config, channel):
  Khóa mutex
  NẾU (số packet < BatchSize) VÀ (đã quá BatchPeriods):
    → Gửi batch vào channel, xóa batch
  HOẶC NẾU (số packet > BatchSize):
    → Gửi batch vào channel, xóa batch
  Mở khóa mutex
```

**Giá trị mặc định:** BatchSize = 10, BatchPeriods = 3 giây

## 6. Xử Lý Packet Trong Main Loop

### 6.1. Tiền Xử Lý Chung (Cho Mọi Batch)

```
1. Sleep 6 giây (chờ 2 block Cosmos để AppHash chứa commitment)

2. UpdateCosmosClient(ctx, "groth16", trustedHeight, "1/3")
   ├─ Lấy trusted block + latest block từ Cosmos
   ├─ Trích xuất chữ ký Ed25519 của validator
   ├─ Sinh Groth16 proof
   ├─ Gửi MsgUpdateClient lên Ethereum
   └─ Trả về latestLightBlock (dùng cho bước tiếp)

3. Lấy ETH block time hiện tại (để kiểm tra timeout)

4. Duyệt từng packet trong batch
```

### 6.2. Case `Send` — Cosmos → ETH RecvPacket

```
Kiểm tra timeout → Tính IBC path (suffix 0x01)
  → ProvePath (Cosmos ABCI query)
    → Kiểm tra value non-empty
      → Parse CommitmentProof → MerkleProof struct
        → Xây dựng MembershipMsg
          → ABI encode + cắt 4 byte selector
            → Xây dựng MsgRecvPacket
              → SendEthTx → ICS26Router.RecvPacket()
```

### 6.3. Case `Ack` — Cosmos Ack → ETH AckPacket

```
Kiểm tra ackBytes non-empty → Tính IBC path (suffix 0x03)
  → ProvePath (Cosmos ABCI query)
    → Parse CommitmentProof → MerkleProof struct
      → Xây dựng MembershipMsg (giống Send)
        → ABI encode + cắt 4 byte selector
          → Xây dựng MsgAckPacket (có thêm Acknowledgement bytes)
            → SendEthTx → ICS26Router.AckPacket()
```

### 6.4. Case `Timeout` — NonMembership Proof

```
Tính IBC path (suffix 0x03)
  → ProvePath (Cosmos ABCI query)
    → Parse CommitmentProof → MerkleProof struct
      → Xây dựng MsgVerifyNonMembership
        → ABI encode "verifyMembership" (⚠️ bug: nên là "verifyNonMembership")
          → Xây dựng MsgRecvPacket (⚠️ nên là timeout msg)
            → SendEthTx
```

### 6.5. Case `WriteAck` — ETH Ack → Cosmos MsgAcknowledgement

Luồng phức tạp nhất vì phải chờ Beacon Chain finality:

```
Lấy CosmosSignerAddress
  │
  ▼
Chờ Beacon finality (tối đa 60 × 10 giây = 10 phút)
  │ Gọi GetFinalityUpdate() lặp lại
  │ So sánh finalizedBlock >= eventBlock
  │
  ▼
UpdateEthClient(ctx)
  │ ├─ GetEthereumClientState từ Cosmos
  │ ├─ GetFinalityUpdate từ Beacon API
  │ ├─ Kiểm tra sync committee participation >= 2/3
  │ ├─ Xử lý chuyển đổi sync committee period
  │ ├─ Chờ Cosmos chain đồng bộ slot
  │ └─ SendCosmosTxBatch(MsgUpdateClient...)
  │
  ▼
Kiểm tra ETH client.LatestExecutionBlockNumber >= eventBlock
  │
  ▼
Tính ack_path: destClientID + 0x03 + sequence(8 bytes)
  │
  ▼
GetEthMembershipProof(ethClient, routerAddr, ackPath, storageSlot, blockNumber)
  │ ├─ Tính storage key: keccak256(keccak256(ackPath) ∥ storageSlot)
  │ ├─ Gọi eth_getProof → AccountProof + StorageProof
  │ └─ Marshal thành JSON (MembershipProof struct)
  │
  ▼
Xây dựng và gửi MsgAcknowledgement
  └─ SendCosmosTx(ctx, ackMsg)
```

## 7. Cập Nhật Light Client

### 7.1. Cosmos Light Client trên Ethereum

**Hàm:** `Worker.UpdateCosmosClient()` — `services/routine.go:67`

```
                        ┌──────────────────┐
                        │  Cosmos RPC Node  │
                        └────────┬─────────┘
                                 │
              ┌──────────────────┼──────────────────┐
              │                  │                   │
              ▼                  ▼                   ▼
      Status()            GetLightBlock()     GetUnbondingTime()
      (latest height)     (trusted + latest)  (staking params)
              │                  │                   │
              └──────────────────┼──────────────────┘
                                 │
                                 ▼
                    ExtractValidatorSignature()
                    ├─ Duyệt commit.Signatures
                    ├─ Bỏ qua BlockIDFlagAbsent
                    ├─ Lấy pubKey (32 bytes Ed25519)
                    ├─ Lấy signature (64 bytes = R ∥ S)
                    ├─ Tính VoteSignBytes (canonical)
                    └─ ed25519.Verify() → xác nhận cục bộ
                                 │
                                 ▼
                       GenerateProof(sig, pub, msg)
                       ├─ Tách R, S từ signature
                       ├─ DecompressPoint(pub) → (aX, aY)
                       ├─ DecompressPoint(R) → (rX, rY)
                       ├─ SHA512(R ∥ pub ∥ msg) → hash
                       ├─ Xây dựng witness (PreHashCircuit)
                       ├─ groth16.Prove(r1cs, pk, witness)
                       ├─ groth16.Verify(proof, vk, pubWitness)
                       └─ ProofToBigInts() → [8]*big.Int
                                 │
                                 ▼
                    SendEthTx(MsgUpdateClient)
                    └─ ICS07.UpdateClient(encodedMsg)
```

### 7.2. Ethereum Light Client trên Cosmos

**Hàm:** `Worker.UpdateEthClient()` — `services/routine.go:364`

```
                      ┌──────────────┐
                      │  Beacon API  │
                      └──────┬───────┘
                             │
              ┌──────────────┼──────────────┐
              │              │              │
              ▼              ▼              ▼
    GetFinalityUpdate()  GetLCUpdates()  GetBootstrap()
    (finalized slot)     (period range)  (sync committee)
              │              │              │
              └──────────────┼──────────────┘
                             │
                             ▼
                 Kiểm tra participation >= 2/3
                             │
                             ▼
              ┌──────────────────────────────┐
              │  Xử lý period crossing:      │
              │  Với mỗi period boundary:    │
              │    buildMsgUpdateClient()     │
              │    (EthereumHeader JSON)      │
              │                              │
              │  Thêm finality update cuối   │
              └──────────────┬───────────────┘
                             │
                             ▼
                  Chờ Cosmos chain đồng bộ slot
                  (current_slot > signature_slot)
                             │
                             ▼
                  SendCosmosTxBatch(msgs)
```

## 8. Gửi Giao Dịch

### 8.1. Giao Dịch Ethereum

**Hàm:** `Handler.SendEthTx()` — `transaction/handler.go:154`

```
Load ETH_PRIVATE_KEY từ env
  │
  ▼
Derive fromAddress từ public key
  │
  ▼
Lấy nonce + gasPrice từ node
  │
  ▼
Tạo auth transactor (bind.NewKeyedTransactorWithChainID)
  ├─ GasLimit: 3,000,000 (cố định)
  ├─ GasPrice: từ SuggestGasPrice
  └─ Value: 0 ETH
  │
  ▼
Phân loại message và gọi contract:
  ├─ MsgUpdateClient    → ics07Tendermint.UpdateClient(auth, encodedData)
  ├─ MsgVerifyMembership → ics07Tendermint.VerifyMembership(auth, msg)
  ├─ MsgVerifyNonMembership → ics07Tendermint.VerifyNonMembership(auth, msg)
  └─ MsgRecvPacket      → ics26Router.RecvPacket(auth, msg)
  │
  ▼
bind.WaitMined() — chờ transaction vào block
  │
  ▼
Kiểm tra receipt.Status:
  ├─ Status == 1: Thành công
  └─ Status == 0: Thất bại → eth_call replay để lấy revert reason
```

### 8.2. Giao Dịch Cosmos

**Hàm:** `Handler.SendCosmosTx()` — `transaction/handler.go:558`

```
Load COSMOS_PRIVATE_KEY từ env
  │
  ▼
Derive signerAddr (bech32) từ secp256k1 public key
  │
  ▼
Đọc cấu hình: COSMOS_CHAIN_ID, COSMOS_GAS_LIMIT, COSMOS_FEE_DENOM
  │
  ▼
queryAccountInfo() — ABCI query /cosmos.auth.v1beta1.Query/Account
  └─ Trả về: accountNumber, sequence (nonce)
  │
  ▼
Đăng ký interface (InterfaceRegistry):
  crypto, auth, channel v2, client, wasm
  │
  ▼
Xây dựng transaction (2 bước ký):
  ├─ Bước 1: SetSignatures(emptySig) → để tạo sign bytes
  ├─ GetSignBytesAdapter(SIGN_MODE_DIRECT) → signBytes
  ├─ privKey.Sign(signBytes) → sigRaw
  └─ Bước 2: SetSignatures(realSig) → signature thật
  │
  ▼
Encode: txConfig.TxEncoder()(tx) → txBytes
  │
  ▼
BroadcastTxCommit(txBytes) — chờ vào block
  │
  ▼
Kiểm tra:
  ├─ CheckTx.Code != 0: Reject tại mempool
  ├─ TxResult.Code != 0: Reject tại DeliverTx
  └─ Cả hai == 0: Thành công
```

### 8.3. So Sánh Hai Loại Giao Dịch

| Đặc điểm | Ethereum | Cosmos |
|-----------|----------|--------|
| **Ký** | ECDSA secp256k1 (go-ethereum) | secp256k1 (cosmos-sdk) |
| **Nonce** | `PendingNonceAt()` | `queryAccountInfo() → sequence` |
| **Gas** | Cố định 3M (relay) / 10M (deploy) | Mặc định 200K / 2M cho MsgUpdateClient |
| **Fee** | gasPrice × gasLimit (ETH) | gasLimit × 1 (stake) |
| **Broadcast** | `bind.WaitMined()` (async + poll) | `BroadcastTxCommit()` (đồng bộ) |
| **Lỗi** | Replay `eth_call` lấy revert reason | CheckTx + DeliverTx code/log |
| **Private key** | `ETH_PRIVATE_KEY` (env) | `COSMOS_PRIVATE_KEY` (env) |

## 9. Hệ Thống Proof

### 9.1. Cosmos Merkle Proof (ICS-23)

Dùng khi chứng minh packet commitment/acknowledgement tồn tại trên Cosmos state tree.

```
ProvePath(cosmosClient, height, path)
  │
  ▼
ABCIQueryWithOptions:
  ├─ path: "store/ibc/key"
  ├─ data: path[1] (ackPath hoặc commitmentPath)
  ├─ height: targetHeight - 1  (proof phải từ block trước)
  └─ prove: true
  │
  ▼
Kết quả:
  ├─ Response.Value: commitment hash (32 bytes)
  ├─ Response.ProofOps: danh sách Merkle proof operations
  └─ ConvertProofs() → MerkleProof (gồm IAVL proof + Tendermint proof)
```

**Encoding IBC path:**
```
IbcCommitmentPath(packet, suffix):
  path = packet.SourceClient + suffix + BigEndian(sequence, 8 bytes)
  return ["ibc", path]

  suffix 0x01 = packet commitment (cho RecvPacket)
  suffix 0x03 = acknowledgement (cho AckPacket)
```

### 9.2. Ethereum Storage Proof (Merkle Patricia Trie)

Dùng khi chứng minh acknowledgement tồn tại trong EVM storage của ICS26Router.

```
GetEthMembershipProof(ethClient, contractAddr, ackPath, slot, blockNumber)
  │
  ▼
Tính storage key (theo Solidity mapping layout):
  storageKey = keccak256(keccak256(ackPath) ∥ slot)
  │
  ▼
Gọi eth_getProof RPC:
  ├─ address: ICS26Router contract
  ├─ storageKeys: [storageKey]
  └─ blockNumber: đã finalize
  │
  ▼
Kết quả JSON:
  {
    "account_proof": {
      "storage_root": "0x...",        // Merkle root của contract storage
      "proof": ["0x...", ...]         // MPT proof: stateRoot → storageRoot
    },
    "storage_proof": {
      "key": "0x...",                 // Storage slot
      "value": "0x...",              // Giá trị tại slot (ack hash)
      "proof": ["0x...", ...]        // MPT proof: storageRoot → value
    }
  }
```

### 9.3. So Sánh Hai Loại Proof

| Đặc điểm | Cosmos Merkle Proof | ETH Storage Proof |
|-----------|--------------------|--------------------|
| **Loại cây** | IAVL + ICS-23 CommitmentProof | Merkle Patricia Trie (MPT) |
| **Hàm query** | `ProvePath()` (ABCI Query) | `GetEthMembershipProof()` (eth_getProof) |
| **Root** | AppHash (trong block header Cosmos) | StorageRoot (trong state trie ETH) |
| **Xác minh bởi** | Smart contract ICS07 trên Ethereum | Wasm light client trên Cosmos |
| **Tính toàn vẹn** | Liên kết với consensus state đã xác minh bằng ZK proof | Liên kết với Beacon finalized slot |

## 10. Prover (Groth16 ZK Proof)

### 10.1. Khởi Tạo

**Hàm:** `prover.NewProver()` — `prover/prover.go:30`

```
Đọc 3 file artifact:
  ├─ r1cs.bin (6.1MB) → Constraint system (mô tả mạch logic)
  ├─ pk.bin (34MB) → Proving key (bí mật cho phép sinh proof)
  └─ vk.bin (1.4KB) → Verifying key (công khai, để xác minh)
```

### 10.2. Sinh Proof

**Hàm:** `EcipProver.GenerateProof()` — `prover/prover.go:70`

```
Đầu vào:
  sig [64]byte = R(32 bytes) ∥ S(32 bytes)   ← Chữ ký Ed25519
  pub [32]byte                                 ← Public key nén
  msg []byte                                   ← Vote sign bytes

Bước 1: Tách signature
  R = sig[:32]                                 ← Điểm trên đường cong
  S = sig[32:]                                 ← Scalar (edwards25519)

Bước 2: Giải nén điểm Ed25519 → tọa độ Weierstrass
  (aX, aY) = DecompressPoint(pub)              ← Public key
  (rX, rY) = DecompressPoint(R)                ← Điểm R

Bước 3: Tính hash off-chain
  H = SHA512(R ∥ pub ∥ msg)
  h = ScalarFromUniformBytes(H)                ← Reduce mod Ed25519 order

Bước 4: Xây dựng witness
  PreHashCircuit {
    Sig:  { R: (rX, rY), S: s }
    Hash: h
    Pub:  { A: (aX, aY) }
  }

Bước 5: Sinh proof
  witness = frontend.NewWitness(assignment, BN254.ScalarField)
  proof = groth16.Prove(r1cs, pk, witness, keccak256HashMode)

Bước 6: Xác minh cục bộ
  pubWitness = witness.Public()
  groth16.Verify(proof, vk, pubWitness)        ← Nếu fail → không gửi on-chain

Bước 7: Chuyển đổi sang format Solidity
  ProofToBigInts(proof):
    out[0..1] = A (G1 point)
    out[2..5] = B (G2 point, thứ tự A1,A0 cho EIP-197)
    out[6..7] = C (G1 point)
    commitments[0..1] = Pedersen commitment
    commitmentPok[0..1] = Proof of knowledge

Đầu ra:
  proof [8]*big.Int         ← A, B, C của Groth16
  commitments [2]*big.Int   ← Pedersen commitment
  commitmentPok [2]*big.Int ← PoK cho commitment
```

## 11. Tổng Hợp Theo File

| File | Dòng | Vai trò chính |
|------|------|---------------|
| `cmd/main.go` | ~650 | CLI, wiring, khởi động hệ thống |
| `services/main.go` | ~500 | Vòng lặp chính, xử lý 4 loại packet |
| `services/routine.go` | ~580 | Worker: UpdateCosmosClient, UpdateEthClient |
| `services/context.go` | ~160 | Context object, getter/setter |
| `services/batch.go` | ~72 | BatchBuilder: gom, kiểm tra, flush packet |
| `services/config.go` | ~50 | Cấu hình runtime (batch size, interval) |
| `subscriber/event.go` | ~278 | Lắng nghe sự kiện Cosmos + ETH |
| `transaction/handler.go` | ~970 | Ký và gửi giao dịch hai chuỗi |
| `transaction/packet.go` | ~25 | Helper xây dựng RecvPacket payload |
| `client/tendermint.go` | ~750 | Truy vấn Cosmos: light block, proof, genesis |
| `client/ethereum.go` | ~820 | Truy vấn ETH: Beacon API, storage proof, SSZ |
| `prover/prover.go` | ~207 | Sinh Groth16 proof, chuyển đổi format |
| `prover/extractor.go` | ~73 | Trích xuất chữ ký validator từ block |
| `keys/keys.go` | ~20 | Khôi phục ETH private key từ hex |
| `utils/abci.go` | ~24 | BytesToBytes32, IbcCommitmentPath |
