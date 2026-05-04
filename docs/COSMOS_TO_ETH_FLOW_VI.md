# Luong Hoat Dong IBC Relayer: Cosmos <-> Ethereum

Tai lieu mo ta chi tiet luong hoat dong cua relayer khi chuyen tiep goi tin (packet) IBC giua chuoi Cosmos va Ethereum, bao gom ca hai huong va vong doi day du cua moi giao dich (RecvPacket + Acknowledgement).

## Tong Quan Kien Truc

```
Cosmos Chain                     Relayer                              Ethereum
+-----------+    WebSocket     +----------------+    ETH JSON-RPC    +-------------------+
| MsgTransfer| ------------->  | SubscribeCosmos |                   | ICS26Router       |
| send_packet|    event        | BatchBuilder    |                   |   RecvPacket()    |
+------------+                 | StartLoop       | ----------------> | ICS07Tendermint   |
                               | UpdateClient    |    SendEthTx()   |   UpdateClient()  |
                               | GenerateProof   |                   |   VerifyMembership|
                               +----------------+                   +-------------------+
```

Luong chinh gom 5 giai doan:
1. **Lang nghe su kien** tu Cosmos qua WebSocket
2. **Gom batch** cac packet lai de xu ly hieu qua
3. **Cap nhat light client** tren Ethereum (bao gom sinh ZK proof)
4. **Chung minh Merkle** cho packet commitment tren Cosmos
5. **Gui giao dich** `RecvPacket` len Ethereum

---

## Giai Doan 1: Lang Nghe Su Kien Cosmos

### `subscriber.SubscribeCosmos(ctx, batchBuilder)`
**File:** `relayer/subscriber/event.go:33`

Ham nay chay trong mot goroutine rieng, lang nghe 3 loai su kien tu Cosmos node qua WebSocket:

| Su kien | Filter CometBFT | Muc dich |
|---------|-----------------|----------|
| `send_packet` | `tm.event='Tx' AND message.action='/ibc.applications.transfer.v1.MsgTransfer'` | User gui token tu Cosmos |
| `write_acknowledgement` | `tm.event='Tx' AND message.action='/ibc.core.channel.v2.MsgRecvPacket'` | Cosmos nhan packet tu ETH va ghi ack |
| `timeout_packet` | `tm.event='Tx' AND message.action='/ibc.applications.transfer.v1.MsgTimeout'` | Packet het han |

**Logic xu ly su kien `send_packet`** (dong 53-79):
1. Lay du lieu event tu truong `send_packet.encoded_packet_hex`
2. Giai ma hex string thanh bytes: `hex.DecodeString(packetEncodedStr)`
3. Unmarshal protobuf thanh `channeltypesv2.Packet`: `proto.Unmarshal(packetBytes, &packet)`
4. Chen vao batch builder voi loai `Send`: `batchBuilder.InsertPacket(services.Packet{PacketType: services.Send, Packet: &packet})`

**Cau truc Packet sau khi parse:**
```
channeltypesv2.Packet {
    Sequence:          uint64      // So thu tu packet
    SourceClient:      string      // Client ID phia Cosmos (vd: "cosmoshub-1")
    DestinationClient: string      // Client ID phia ETH (vd: "08-wasm-0")
    TimeoutTimestamp:  uint64      // Thoi gian het han (Unix nano)
    Payloads:          []Payload   // Du lieu chuyen (port, version, encoding, value)
}
```

---

## Giai Doan 2: Gom Batch

### `BatchBuilder.InsertPacket(packet)`
**File:** `relayer/services/batch.go:43`

Them packet vao danh sach cho xu ly. Su dung `sync.Mutex` de dam bao an toan khi nhieu goroutine ghi dong thoi.

```go
func (b *BatchBuilder) InsertPacket(packet Packet) {
    b.mtx.Lock()
    b.packets = append(b.packets, packet)
    b.mtx.Unlock()
}
```

### `BatchBuilder.CheckBatch(config, ch)`
**File:** `relayer/services/batch.go:54`

Chay trong goroutine rieng, kiem tra moi 3 giay xem co nen gui batch di xu ly khong. Co 2 dieu kien de flush batch:

1. **Theo thoi gian:** So luong packet < `BatchSize` (mac dinh 10) VA da qua `BatchPeriods` (mac dinh 3 giay) ke tu batch truoc
2. **Theo so luong:** So luong packet > `BatchSize`

Khi mot trong hai dieu kien thoa, gui `BatchPackets` vao channel de vong lap chinh xu ly.

---

## Giai Doan 3: Vong Lap Xu Ly Chinh

### `Services.StartLoop(ctx)`
**File:** `relayer/services/main.go:82`

Day la ham trung tam dieu phoi toan bo relayer. No khoi dong 4 goroutine:

| Goroutine | Chuc nang |
|-----------|-----------|
| `SubscribeCosmos` | Lang nghe su kien Cosmos |
| `SubscribeEth` | Lang nghe su kien Ethereum |
| Routine scheduler | Cap nhat light client dinh ky (24h) |
| Batch checker | Kiem tra batch moi 3 giay |

**Xu ly khi nhan duoc batch** (dong 132-499):

Voi moi batch nhan tu channel `s.BatchPackets`:

1. **Cho 2 block** (6 giay) de packet commitment duoc bao gom vao `AppHash` cua Cosmos
2. **Cap nhat Cosmos light client tren Ethereum** (giai doan 3.1)
3. **Lay ETH block time** de kiem tra timeout
4. **Duyet tung packet** trong batch va xu ly theo loai

---

## Giai Doan 3.1: Cap Nhat Cosmos Light Client Tren Ethereum

### `Worker.UpdateCosmosClient(ctx, proofType, trustedBlock, trustLevel)`
**File:** `relayer/services/routine.go:67`

Day la buoc quan trong nhat — cap nhat trang thai cua chuoi Cosmos tren smart contract Ethereum de Ethereum biet trang thai moi nhat cua Cosmos.

**Logic chi tiet:**

**Buoc 1 — Xac dinh trusted block** (dong 68-106):
- Query `Status()` tu Cosmos node de lay `LatestBlockHeight`
- Neu `trustedBlock == 0` (lan dau): doc `clientState` tu smart contract ICS07 tren Ethereum de lay `LatestHeight.RevisionHeight`
- Neu `trustedBlock >= LatestBlockHeight`: client da cap nhat roi, tra ve light block hien tai
- Neu `trustedBlock < LatestBlockHeight`: can cap nhat

**Buoc 2 — Lay du lieu light block** (dong 108-117):
```go
trustedLightBlock = GetLightBlock(cosmosClient, trustedBlock)    // block da tin tuong
latestLightBlock  = GetLightBlock(cosmosClient, latestBlockHeight) // block moi nhat
```

Moi `LightBlock` gom: `SignedHeader` (header + commit signatures), `ValidatorSet`, `NextValidatorSet`.

**Buoc 3 — Tinh toan cac tham so** (dong 119-165):
- `unbondingPeriod`: lay tu Cosmos node (thoi gian unstake)
- `trustingPeriod = unbondingPeriod * 2/3`: thoi gian tin tuong
- `chainId`, `revision`: trich tu header cua trusted block
- `trustThreshold`: parse tu chuoi "1/3"
- Xay dung `clientState` va `consensusState` cho smart contract

**Buoc 4 — Tao header de xuat** (dong 167):
```go
proposedHeader = latestLightBlock.IntoHeader(*trustedLightBlock)
```
Chuyen `LightBlock` thanh dinh dang `Header` ma smart contract chap nhan, bao gom: `SignedHeader`, `ValidatorSet`, `TrustedValidatorSet`, `TrustedHeight`.

**Buoc 5 — Trich xuat chu ky validator** (dong 173-178):
```go
valSig = prover.ExtractValidatorSignature(latestLightBlock, chainId)
```

### `prover.ExtractValidatorSignature(lightBlock, chainID)`
**File:** `relayer/prover/extractor.go:21`

Tim chu ky Ed25519 hop le dau tien tu commit cua block:

1. Duyet qua cac `commit.Signatures`
2. Bo qua cac chu ky co `BlockIDFlag == BlockIDFlagAbsent`
3. Voi moi chu ky:
   - Lay public key cua validator tuong ung (32 bytes Ed25519)
   - Lay du lieu chu ky `sig.Signature` (64 bytes = R || S)
   - Tinh `VoteSignBytes` (du lieu canonical ma validator da ky)
   - **Xac minh chu ky**: `ed25519.Verify(pubKey, voteData, sigData)`
4. Tra ve `ValidatorSignature{Signature, PublicKey, SignBytes}` dau tien hop le

> Luu y: Hien tai chi chung minh 1 chu ky validator. Voi mang multi-validator can mo rong de chung minh >= 2/3 voting power.

**Buoc 6 — Sinh ZK Proof (Groth16)** (dong 180-183):
```go
proof, commitments, commitmentPok = w.Prover.GenerateProof(valSig.Signature, valSig.PublicKey, valSig.SignBytes)
```

### `EcipProver.GenerateProof(sig, pub, msg)`
**File:** `relayer/prover/prover.go:70`

Sinh Groth16 zero-knowledge proof de chung minh rng chu ky Ed25519 la hop le ma khong tiet lo private key. Day la phan nang nhat ve tinh toan.

**Logic chi tiet:**

1. **Validate dau vao**: sig phai 64 bytes, pub phai 32 bytes
2. **Tach signature**: `R = sig[:32]` (diem tren duong cong), `S = sig[32:]` (scalar)
3. **Giai nen diem Ed25519 sang toa do Weierstrass**:
   ```go
   aX, aY = DecompressPoint(pub)  // Public key -> (x, y)
   rX, rY = DecompressPoint(R)    // R point -> (x, y)
   ```
4. **Tinh hash off-chain**: `H = SHA512(R || pub || msg)` roi chuyen sang scalar tren truong Fr25519
5. **Xay dung witness** cho mach (circuit):
   ```
   assignment = PreHashCircuit {
       Sig:  { R: (rX, rY), S: s }
       Hash: h
       Pub:  { A: (aX, aY) }
   }
   ```
6. **Tao witness**: `frontend.NewWitness(&assignment, BN254.ScalarField())`
7. **Sinh proof**: `groth16.Prove(r1cs, pk, witness, ...)`
   - Su dung proving key `pk.bin` (34MB, da pre-compute)
   - Su dung constraint system `r1cs.bin` (6.1MB)
   - Dung keccak256 cho commitment hash (tuong thich Solidity)
8. **Xac minh local**: `groth16.Verify(proof, vk, pubWitness, ...)` — kiem tra proof truoc khi gui on-chain
9. **Chuyen sang dinh dang Solidity**: `ProofToBigInts(proof)` — trich xuat 8 big.Int tu cac diem G1/G2 cua proof

**Ket qua tra ve:**
- `proof [8]*big.Int`: A(G1), B(G2), C(G1) — bang chung ZK
- `commitments [2]*big.Int`: Pedersen commitment
- `commitmentPok [2]*big.Int`: Proof of knowledge cho commitment

**Buoc 7 — Gui transaction `updateClient` len Ethereum** (dong 185-200):
```go
msg = IUpdateClientMsgsMsgUpdateClient{
    ClientState:           clientState,
    TrustedConsensusState: consensusState,
    Time:                  now,
    ProposedHeader:        proposedHeader,
    Proof:                 proof,         // ZK proof
    Commitments:           commitments,
    CommitmentPok:         commitmentPok,
}
w.TxHandler.SendEthTx(ctx, msg)
```

Smart contract `Groth16ICS07Tendermint.UpdateClient()` se:
- Xac minh Groth16 proof on-chain (kiem tra chu ky validator)
- Cap nhat `clientState` va `consensusState` luu tru
- Sau do cac lenh `VerifyMembership` co the su dung consensus state moi nay

---

## Giai Doan 4: Chung Minh Merkle Cho Packet (Case `Send`)

Sau khi light client da cap nhat, relayer can chung minh rang packet commitment ton tai tren Cosmos state tree.

**File:** `relayer/services/main.go:173-257`

### Buoc 4.1 — Kiem tra timeout
```go
if ethBlockTime >= packet.TimeoutTimestamp {
    // Packet da het han, bo qua
    continue
}
```

### Buoc 4.2 — Tinh IBC commitment path

### `utils.IbcCommitmentPath(packet, appendByte)`
**File:** `relayer/utils/abci.go:15`

Tao duong dan Merkle trong IBC store de truy van proof:

```go
func IbcCommitmentPath(packet, appendByte) [][]byte {
    sequenceBytes = BigEndian(packet.Sequence)  // 8 bytes
    path = packet.SourceClient + appendByte + sequenceBytes
    return [][]byte{"ibc", path}
}
```

- Voi `Send` packet: `appendByte = []byte{1}` (packet commitment)
- Voi `Ack` packet: `appendByte = []byte{3}` (acknowledgement)

Vi du: Client "cosmoshub-1", sequence 5 -> path = `["ibc", "cosmoshub-1\x01\x00\x00\x00\x00\x00\x00\x00\x05"]`

### Buoc 4.3 — Lay Merkle proof tu Cosmos

### `client.ProvePath(cosmosClient, height, path)`
**File:** `relayer/client/tendermint.go:476`

Truy van ABCI cua Cosmos node de lay Merkle proof:

1. Xay dung query path: `store/ibc/key`
2. Goi `ABCIQueryWithOptions` tai `height - 1` (proof phai tu block truoc)
3. Kiem tra:
   - `Response.Code == 0` (thanh cong)
   - `Response.Height == height - 1` (dung height)
   - `Response.Key` khop voi `path[1]`
4. Chuyen `ProofOps` thanh `MerkleProof`: `commitmenttypes.ConvertProofs(result.Response.ProofOps)`
5. Tra ve `(value, merkleProof)` — value la packet commitment hash, proof la bang chung Merkle

### Buoc 4.4 — Parse commitment proof

Chuyen moi `CommitmentProof` tu dinh dang ICS23 sang dinh dang ma smart contract chap nhan:
```go
for _, p := range proof.Proofs {
    commitmentProof = client.ParseCommitmentProof(p)
    merkleProof.Proofs = append(merkleProof.Proofs, commitmentProof)
}
```

### Buoc 4.5 — Xay dung Membership Message

```go
membershipMsg = ILightClientMsgsMsgVerifyMembership{
    Height:     { RevisionHeight: latestBlockHeight },
    KvPairs:    [{ Path: ibcPath, Value: BytesToBytes32(value) }],
    MerkleProofs: [merkleProof],
    AppHash:    BytesToBytes32(latestLightBlock.SignedHeader.AppHash),
    TrustedConsensusState: {
        Timestamp:          latestLightBlock.Time (nanoseconds),
        Root:               AppHash,
        NextValidatorsHash: NextValidatorsHash,
    },
    MembershipType: 0,  // verify membership (khong phai non-membership)
}
```

### Buoc 4.6 — ABI encode membership msg

```go
calldata = tendermintAbiJson.Pack("verifyMembership", membershipMsg)
calldata = calldata[4:]  // Cat bo 4 byte function selector
```

Cat function selector vi ICS26Router se goi `abi.decode` truc tiep, khong phai goi ham `verifyMembership()`.

---

## Giai Doan 5: Gui RecvPacket Len Ethereum

### Buoc 5.1 — Xay dung message RecvPacket

```go
payloads = convert(packet.Payloads)  // Chuyen tu Cosmos format sang ETH format

msgRecvPacket = IICS26RouterMsgsMsgRecvPacket{
    Packet: {
        Sequence:         packet.Sequence,
        SourceClient:     packet.SourceClient,
        DestClient:       packet.DestinationClient,
        TimeoutTimestamp: packet.TimeoutTimestamp,
        Payloads:         payloads,
    },
    MembershipMsg: calldata,  // Membership proof da encode
}
```

### Buoc 5.2 — Gui giao dich Ethereum

### `Handler.SendEthTx(ctx, msg)`
**File:** `relayer/transaction/handler.go:154`

**Logic chi tiet:**

1. **Load private key** tu bien moi truong `ETH_PRIVATE_KEY`
2. **Tao transactor**:
   ```go
   privateKey = keys.RestoreKey(privKeyHex)
   fromAddress = crypto.PubkeyToAddress(publicKey)
   nonce = ethClient.PendingNonceAt(fromAddress)
   gasPrice = ethClient.SuggestGasPrice()
   auth = bind.NewKeyedTransactorWithChainID(privateKey, chainId)
   auth.GasLimit = 3000000
   ```
3. **Khoi tao contract instances**:
   ```go
   ics07Tendermint = NewContractGroth16ICS07Tendermint(clientContract, ethClient)
   ics26Router     = NewContractICS26Router(routerContract, ethClient)
   ```
4. **Phan loai message va gui tx** (switch msg type):
   - `MsgUpdateClient` -> `ics07Tendermint.UpdateClient(auth, encodedData)`
   - `MsgVerifyMembership` -> `ics07Tendermint.VerifyMembership(auth, msg)`
   - `MsgVerifyNonMembership` -> `ics07Tendermint.VerifyNonMembership(auth, msg)`
   - `MsgRecvPacket` -> `ics26Router.RecvPacket(auth, msg)`
5. **Cho receipt**: `bind.WaitMined(ctx, ethClient, tx)`
6. **Kiem tra status**: `receipt.Status == 1` (thanh cong) hoac `0` (that bai/revert)

---

## Luong Xu Ly ETH -> Cosmos: RecvPacket

Khi user gui packet tu Ethereum, relayer can chuyen tiep len Cosmos.

### Giai doan 1: Lang nghe su kien ETH

### `subscriber.SubscribeEth(ctx, batchBuilder)`
**File:** `relayer/subscriber/event.go:180`

Ham nay chay trong goroutine rieng, dang ky lang nghe 4 loai event tu smart contract `ICS26Router` tren Ethereum qua WebSocket:

| Event | Channel | Muc dich |
|-------|---------|----------|
| `SendPacket` | `sendPacketCh` | User gui packet tu ETH |
| `WriteAcknowledgement` | `writeAckCh` | ETH nhan packet tu Cosmos va ghi ack |
| `AckPacket` | `ackPacketCh` | ETH nhan ack (chua xu ly — TODO) |
| `TimeoutPacket` | `timeoutPacketCh` | Packet het han (chua xu ly — TODO) |

**Logic khoi tao** (dong 181-227):
1. Tao `ContractICS26RouterFilterer` tu contract address + WebSocket client
2. Tao 4 channel tuong ung voi 4 loai event
3. Goi `WatchSendPacket`, `WatchWriteAcknowledgement`, `WatchAckPacket`, `WatchTimeoutPacket` de dang ky subscription
4. Moi subscription co `defer Unsubscribe()` de don dep khi ham return

**Xu ly event `SendPacket`** (dong 232-238):
```go
case ev := <-sendPacketCh:
    cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
    batchBuilder.InsertPacket(services.Packet{
        PacketType: services.Send,
        Packet:     &cosmosPacket,
    })
```

### `EthPacketToCosmosPacket(ethPacket, sequence)`
**File:** `relayer/subscriber/event.go:158`

Chuyen doi packet tu dinh dang Solidity struct sang dinh dang Cosmos protobuf:

```go
func EthPacketToCosmosPacket(ethPacket, sequence) channeltypesv2.Packet {
    // Chuyen tung payload: SourcePort, DestPort -> DestinationPort, Version, Encoding, Value
    for _, p := range ethPacket.Payloads {
        payloads = append(payloads, channeltypesv2.Payload{
            SourcePort:      p.SourcePort,
            DestinationPort: p.DestPort,    // ten truong khac nhau giua 2 dinh dang
            Version:         p.Version,
            Encoding:        p.Encoding,
            Value:           p.Value,
        })
    }
    return channeltypesv2.Packet{
        Sequence:          sequence.Uint64(),
        SourceClient:      ethPacket.SourceClient,
        DestinationClient: ethPacket.DestClient,
        TimeoutTimestamp:  ethPacket.TimeoutTimestamp,
        Payloads:          payloads,
    }
}
```

> **Luu y:** Hien tai event `SendPacket` tu ETH duoc chen vao batch voi `PacketType: Send`, cung loai voi packet tu Cosmos. Trong `StartLoop`, case `Send` chi xu ly huong Cosmos->ETH (goi `ProvePath` tren Cosmos node va `SendEthTx`). Viec xu ly huong ETH->Cosmos (gui `MsgRecvPacket` len Cosmos) can duoc bo sung them.

---

## Luong Xu Ly Ack (Cosmos Ack -> Ethereum)

Sau khi Cosmos nhan packet tu ETH va xu ly xong, no ghi `write_acknowledgement`. Relayer bat su kien nay va gui `AckPacket` ve Ethereum de ETH biet ket qua xu ly.

**File:** `relayer/services/main.go:258-339`

### Giai doan 1: Bat su kien `write_acknowledgement` tu Cosmos

**File:** `relayer/subscriber/event.go:80-123`

Khi Cosmos xu ly `MsgRecvPacket` (nhan packet tu ETH), no emit event `write_acknowledgement` chua:
- `write_acknowledgement.encoded_packet_hex`: packet goc (hex-encoded protobuf)
- `write_acknowledgement.encoded_acknowledgement_hex`: ket qua xu ly (hex-encoded protobuf)

**Logic parse:**
1. Decode hex -> bytes cho ca packet va acknowledgement
2. Unmarshal protobuf: `proto.Unmarshal(packetBytes, &packet)` va `proto.Unmarshal(ackBytes, &acknowledgement)`
3. Kiem tra `len(acknowledgement.AppAcknowledgements) > 0`
4. Chen vao batch:
   ```go
   batchBuilder.InsertPacket(services.Packet{
       PacketType: services.Ack,
       Packet:     &packet,
       AckBytes:   acknowledgement.AppAcknowledgements,
   })
   ```

### Giai doan 2: Xu ly trong StartLoop (case `Ack`)

**File:** `relayer/services/main.go:258-339`

Luong nay tuong tu voi case `Send` (Cosmos->ETH RecvPacket), nhung co mot so khac biet quan trong:

**Buoc 2.1 — Kiem tra ack bytes:**
```go
if len(packet.AckBytes) == 0 {
    // Khong co acknowledgement data, bo qua
    continue
}
```

**Buoc 2.2 — Tinh IBC acknowledgement path:**
```go
ibcPath := utils.IbcCommitmentPath(*packet.Packet, []byte{3})  // 0x03 = ack suffix
```

Khac voi `Send` dung `[]byte{1}` (commitment suffix), `Ack` dung `[]byte{3}` (acknowledgement suffix).

Vi du: Client "cosmoshub-1", sequence 5 -> path = `["ibc", "cosmoshub-1\x03\x00\x00\x00\x00\x00\x00\x00\x05"]`

**Buoc 2.3 — Lay Merkle proof tu Cosmos:**
```go
value, proof, err := client.ProvePath(ctx.CosmosClient(), latestLightBlock.BlockHeight, ibcPath)
```
Giong nhu case `Send`, truy van ABCI de chung minh acknowledgement ton tai tren Cosmos state tree.

**Buoc 2.4 — Xay dung Membership Message:**
Giong hoan toan case `Send` — xay dung `ILightClientMsgsMsgVerifyMembership` voi `MembershipType: 0`, su dung `AppHash` va `TrustedConsensusState` tu `latestLightBlock`.

**Buoc 2.5 — ABI encode:**
```go
calldata, err := tendermintAbiJson.Pack("verifyMembership", membershipMsg)
calldata = calldata[4:]  // Cat function selector
```

**Buoc 2.6 — Xay dung MsgAckPacket** (khac voi MsgRecvPacket):
```go
msgAckPacket := contractICS26Router.IICS26RouterMsgsMsgAckPacket{
    Packet: contractICS26Router.IICS26RouterMsgsPacket{
        Sequence:         packet.Packet.Sequence,
        SourceClient:     packet.Packet.SourceClient,
        DestClient:       packet.Packet.DestinationClient,
        TimeoutTimestamp: packet.Packet.TimeoutTimestamp,
        Payloads:         payloads,
    },
    Acknowledgement: packet.AckBytes[0],   // Byte array ket qua xu ly
    MembershipMsg:   calldata,              // Membership proof da encode
}
```

**So sanh MsgRecvPacket vs MsgAckPacket:**

| Truong | MsgRecvPacket | MsgAckPacket |
|--------|---------------|--------------|
| `Packet` | Co | Co |
| `MembershipMsg` | Co (proof commitment) | Co (proof ack) |
| `Acknowledgement` | Khong | Co (ack bytes) |
| Path suffix | `0x01` (commitment) | `0x03` (acknowledgement) |

**Buoc 2.7 — Gui AckPacket len Ethereum:**
```go
s.worker.TxHandler.SendEthTx(ctx, msgAckPacket)
```

Smart contract `ICS26Router.AckPacket()` se:
- Goi ICS07 `VerifyMembership` de kiem tra ack ton tai tren Cosmos
- Xu ly ack: cap nhat trang thai packet (da hoan tat)
- Giai phong token escrow neu la ICS20 transfer

---

## Luong Xu Ly WriteAck (ETH WriteAck -> Cosmos MsgAcknowledgement)

Sau khi Ethereum nhan `RecvPacket` va xu ly xong, no emit event `WriteAcknowledgement`. Relayer bat su kien nay va gui `MsgAcknowledgement` ve Cosmos de hoan tat chu ky IBC.

**File:** `relayer/services/main.go:416-491`

Day la luong phuc tap nhat vi phai cho Beacon Chain finality (khong nhu Cosmos co instant finality).

### Giai doan 1: Bat su kien `WriteAcknowledgement` tu Ethereum

**File:** `relayer/subscriber/event.go:240-248`

```go
case ev := <-writeAckCh:
    cosmosPacket := EthPacketToCosmosPacket(ev.Packet, ev.Sequence)
    batchBuilder.InsertPacket(services.Packet{
        PacketType:  services.WriteAck,
        Packet:      &cosmosPacket,
        AckBytes:    ev.Acknowledgements,     // Ack data tu ETH
        BlockNumber: ev.Raw.BlockNumber,      // Block number de cho finality
    })
```

Diem khac biet: packet co them `BlockNumber` — can thiet de biet khi nao Beacon Chain da finalize block nay.

### Giai doan 2: Xu ly trong StartLoop (case `WriteAck`)

**Buoc 2.1 — Lay Cosmos signer address:**

### `Handler.CosmosSignerAddress()`
**File:** `relayer/transaction/handler.go:548`

```go
func (h *Handler) CosmosSignerAddress() (string, error) {
    privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
    privKeyBytes := hex.DecodeString(privKeyHex)
    privKey := secp256k1.PrivKey{Key: privKeyBytes}
    return sdk.AccAddress(privKey.PubKey().Address()).String(), nil
}
```

Derive Cosmos bech32 address tu private key de dien vao truong `Signer` cua message.

**Buoc 2.2 — Cho Beacon Chain finality** (dong 423-445):

Ethereum khong co instant finality nhu Cosmos. Block phai duoc Beacon Chain finalize truoc khi co the tao proof hop le.

```go
for attempt := 0; attempt < 60; attempt++ {     // Toi da 60 lan
    time.Sleep(10 * time.Second)                 // Moi lan cho 10 giay

    finalityUpdate := client.GetFinalityUpdate(beaconAPIURL)
    execBlock := finalityUpdate.FinalizedHeader.Execution.BlockNumber

    if execBlock >= packet.BlockNumber {         // Block da finalized!
        finalized = true
        break
    }
    // Chua finalized, tiep tuc cho...
}
if !finalized {
    // Timeout sau 10 phut, bo qua packet nay
    continue
}
```

Beacon finality mat khoang 12-15 phut (2 epoch). Vong lap cho toi da ~10 phut (60 * 10s).

**Buoc 2.3 — Cap nhat ETH Light Client tren Cosmos:**

### `Worker.UpdateEthClient(ctx)`
**File:** `relayer/services/routine.go:364`

Cap nhat trang thai Ethereum tren Cosmos chain de Cosmos biet trang thai moi nhat cua ETH.

**Logic chi tiet:**

1. **Lay ETH client state hien tai tu Cosmos:**
   ```go
   ethClientState := client.GetEthereumClientState(cosmosClient, ethClientID)
   trustedSlot := ethClientState.LatestSlot
   ```

   ### `client.GetEthereumClientState(cosmosClient, clientID)`
   **File:** `relayer/client/ethereum.go:644`

   Truy van IBC client state tren Cosmos:
   - Goi ABCI query: `/ibc.core.client.v1.Query/ClientState`
   - Unmarshal: `QueryClientStateResponse` -> `ibcwasmtypes.ClientState` -> `EthereumClientState` (JSON)
   - Tra ve struct chua: `LatestSlot`, `LatestExecutionBlockNumber`, `SyncCommitteeSize`, `GenesisTime`...

2. **Lay finality update tu Beacon API:**
   ```go
   finalityUpdate := client.GetFinalityUpdate(beaconAPIURL)
   ```

3. **Kiem tra sync committee participation** (bao mat):
   ```go
   participation := CountSyncCommitteeParticipants(finalityUpdate.SyncAggregate.SyncCommitteeBits)
   if participation * 3 < syncCommitteeSize * 2 {
       return error  // Can >= 2/3 sync committee tham gia
   }
   ```

4. **Kiem tra da cap nhat chua:**
   ```go
   if finalizedSlot <= trustedSlot {
       return nil  // Da up-to-date
   }
   ```

5. **Xu ly chuyen doi sync committee period:**

   ### `Worker.updateEthClientWithPeriodCrossing(...)`
   **File:** `relayer/services/routine.go:416`

   Ethereum thay doi sync committee moi ~27 gio (1 period = 256 epoch). Neu can cap nhat qua nhieu period, phai gui nhieu `MsgUpdateClient` lien tiep.

   **Logic:**
   - Tinh `trustedPeriod` va `targetPeriod` tu slot numbers
   - Lay `LightClientUpdates` cho tat ca period tu `trustedPeriod` -> `targetPeriod`
   - Voi moi period crossing: xay dung `MsgUpdateClient` chua `EthereumHeader` voi `NextSyncCommittee`
   - Them mot `MsgUpdateClient` cuoi cung cho finality update moi nhat
   - **Cho Cosmos chain dong bo**: Wasm light client tren Cosmos tinh `current_slot` tu `block.time`. Neu `current_slot < signature_slot`, giao dich se bi reject. Vong lap cho toi khi Cosmos bat kip:
     ```go
     for range 60 {
         currentSlot := ethClientState.ComputeSlotAtTimestamp(cosmosTime)
         if currentSlot > sigSlot { break }
         time.Sleep(5 * time.Second)
     }
     ```
   - Gui tat ca messages trong 1 batch: `SendCosmosTxBatch(ctx, msgs)`

6. **Xay dung MsgUpdateClient:**

   ### `buildMsgUpdateClient(signer, clientID, header)`
   **File:** `relayer/services/routine.go:561`

   ```go
   headerBytes := json.Marshal(header)           // EthereumHeader -> JSON
   clientMessage := ibcwasmtypes.ClientMessage{Data: headerBytes}
   clientMessageAny := codectypes.NewAnyWithValue(clientMessage)
   return &clienttypes.MsgUpdateClient{
       ClientId:      clientID,        // vd: "08-wasm-0"
       ClientMessage: clientMessageAny,
       Signer:        signerAddr,
   }
   ```

**Buoc 2.4 — Xac nhan ETH client da cap nhat du:**
```go
ethClientState := client.GetEthereumClientState(cosmosClient, ethClientID)
proofBlockNumber := ethClientState.LatestExecutionBlockNumber
proofSlot := ethClientState.LatestSlot

if proofBlockNumber < packet.BlockNumber {
    // ETH client chua cap nhat toi block chua event, bo qua
    continue
}
```

**Buoc 2.5 — Tinh ack_path va lay ETH storage proof:**

```go
// ack_path = destClientID + 0x03 + sequence(8 bytes big-endian)
seqBytes := make([]byte, 8)
binary.BigEndian.PutUint64(seqBytes, packet.Packet.Sequence)
ackPath := append([]byte(packet.Packet.DestinationClient), 0x03)
ackPath = append(ackPath, seqBytes...)
```

### `client.GetEthMembershipProof(ethClient, contractAddr, ackPath, slot, blockNumber)`
**File:** `relayer/client/ethereum.go:774`

Lay Merkle Patricia proof tu EVM storage de chung minh ack da duoc ghi tren Ethereum.

**Logic:**
1. **Tinh storage key** (theo layout cua Solidity mapping):
   ```go
   pathHash := crypto.Keccak256(ackPath)
   storageKey := crypto.Keccak256Hash(pathHash, slot.Bytes())
   // slot = ICS26_IBC_STORAGE_SLOT (vi tri cua mapping trong storage)
   ```

2. **Goi `eth_getProof`** — RPC method cua Ethereum de lay Merkle proof:
   ```go
   client.CallContext(ctx, &result, "eth_getProof",
       contractAddr,                     // Dia chi ICS26Router
       []string{storageKey.Hex()},       // Storage key can chung minh
       toBlockNumArg(blockNumber),       // Tai block number da finalize
   )
   ```

3. **Xay dung proof JSON** tuong thich voi Wasm light client tren Cosmos:
   ```go
   proof := MembershipProof{
       AccountProof: {
           StorageRoot: result.StorageHash.Hex(),  // Merkle root cua contract storage
           Proof:       result.AccountProof,        // MPT proof: stateRoot -> storageRoot
       },
       StorageProof: {
           Key:   storageKey.Hex(),                 // Slot duoc chung minh
           Value: sp.Value.String(),                // Gia tri tai slot (ack hash)
           Proof: sp.Proof,                         // MPT proof: storageRoot -> value
       },
   }
   return json.Marshal(proof)
   ```

**So sanh 2 loai proof:**

| | Cosmos Merkle Proof | ETH Storage Proof |
|---|---|---|
| **Loai cay** | IAVL + ICS23 CommitmentProof | Merkle Patricia Trie |
| **Ham query** | `ProvePath` (ABCIQuery) | `GetEthMembershipProof` (eth_getProof) |
| **Root** | AppHash (tren Cosmos block header) | StorageRoot (tren ETH state) |
| **Xac minh boi** | Smart contract ICS07 tren ETH | Wasm light client tren Cosmos |

**Buoc 2.6 — Gui MsgAcknowledgement len Cosmos:**

```go
ackMsg := &channeltypesv2.MsgAcknowledgement{
    Packet:          *packet.Packet,
    Acknowledgement: channeltypesv2.Acknowledgement{
        AppAcknowledgements: packet.AckBytes,    // Ack bytes tu ETH event
    },
    ProofAcked:  proofBytes,                     // ETH storage proof (JSON)
    ProofHeight: clienttypes.Height{
        RevisionNumber: 0,
        RevisionHeight: proofSlot,               // Beacon slot (khong phai block number)
    },
    Signer: signerAddr,                          // Cosmos relayer address
}
```

### `Handler.SendCosmosTx(ctx, msg)`
**File:** `relayer/transaction/handler.go:558`

Ky va gui giao dich len Cosmos chain.

**Logic chi tiet:**

1. **Load private key** tu `COSMOS_PRIVATE_KEY` env var
2. **Derive signer address**: `sdk.AccAddress(privKey.PubKey().Address())`
3. **Doc cau hinh**: `COSMOS_CHAIN_ID`, `COSMOS_GAS_LIMIT` (mac dinh 200,000), `COSMOS_FEE_DENOM` (mac dinh "stake")
4. **Xu ly dac biet cho MsgUpdateClient**: Tang gas len 2,000,000 vi wasm verification ton nhieu gas
5. **Query account info** tu chain:
   ```go
   accountNumber, sequence := queryAccountInfo(ctx, signerAddr)
   ```
   ### `Handler.queryAccountInfo(ctx, address)`
   **File:** `relayer/transaction/handler.go:917`

   Goi ABCI query `/cosmos.auth.v1beta1.Query/Account` de lay `accountNumber` va `sequence` (nonce) cua tai khoan.

6. **Setup encoding**: Dang ky cac interface (crypto, auth, channel, client) vao `InterfaceRegistry`
7. **Xay dung transaction**:
   ```go
   txBuilder := txConfig.NewTxBuilder()
   txBuilder.SetMsgs(sdkMsg)
   txBuilder.SetGasLimit(gasLimit)
   txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, feeAmount)))
   ```
8. **Ky transaction** (2 buoc):
   - Buoc 1: Set empty signature de sinh sign bytes
   - Buoc 2: Ky sign bytes bang private key, set actual signature
   ```go
   signBytes := authsigning.GetSignBytesAdapter(ctx, handler, SIGN_MODE_DIRECT, signerData, tx)
   sigRaw := privKey.Sign(signBytes)
   txBuilder.SetSignatures(sigV2WithRealSig)
   ```
9. **Encode va broadcast**:
   ```go
   txBytes := txConfig.TxEncoder()(txBuilder.GetTx())
   result := cosmosClient.BroadcastTxCommit(ctx, txBytes)   // Cho tx vao block
   ```
10. **Kiem tra ket qua**:
    - `CheckTx.Code != 0`: Tx bi reject tai mempool (validation fail)
    - `TxResult.Code != 0`: Tx bi reject tai DeliverTx (execution fail)
    - Ca hai `== 0`: Thanh cong

### `Handler.SendCosmosTxBatch(ctx, msgs)`
**File:** `relayer/transaction/handler.go:730`

Gui nhieu message trong 1 giao dich Cosmos (dung cho `UpdateEthClient` khi can nhieu `MsgUpdateClient`).

**Khac voi `SendCosmosTx`:**
- Nhan `[]any` thay vi `any` — nhieu message
- Gas = `baseGas * len(msgs)` (tinh theo so luong message)
- Fee = `feeAmount * len(msgs)` (ti le voi gas)
- Neu co `MsgUpdateClient` trong batch: `baseGas = max(baseGas, 2000000)`
- Dung `txBuilder.SetMsgs(sdkMsgs...)` de set nhieu message cung luc

---

## So Do Tong The: Chu Ky Day Du Cosmos -> ETH

```
Thoi gian -->

=== GIAI DOAN A: Cosmos gui packet, ETH nhan (RecvPacket) ===

[1] User gui MsgTransfer tren Cosmos
         |
[2] Cosmos emit "send_packet" event
         |
[3] SubscribeCosmos nhan event qua WebSocket
         |
[4] Parse hex -> protobuf -> channeltypesv2.Packet
         |
[5] InsertPacket vao BatchBuilder (PacketType: Send)
         |
[6] CheckBatch flush batch (sau 3s hoac >= 10 packets)
         |
[7] StartLoop nhan batch tu channel
         |
[8] Sleep 6s (cho AppHash cap nhat)
         |
[9] UpdateCosmosClient:
    [9a] GetLightBlock(trustedHeight) + GetLightBlock(latestHeight)
    [9b] ExtractValidatorSignature (Ed25519)
    [9c] GenerateProof (Groth16 ZK proof)
    [9d] SendEthTx(MsgUpdateClient) -> ICS07.UpdateClient()
         |
[10] Voi moi packet trong batch:
    [10a] Kiem tra timeout
    [10b] IbcCommitmentPath(packet, 0x01)
    [10c] ProvePath -> ABCIQuery lay Merkle proof
    [10d] Xay dung MembershipMsg + ABI encode
    [10e] SendEthTx(MsgRecvPacket) -> ICS26Router.RecvPacket()

=== GIAI DOAN B: ETH ghi ack, Cosmos nhan (WriteAck) ===

         |
[11] Ethereum ICS26Router:
    [11a] VerifyMembership qua ICS07 (kiem tra Merkle proof)
    [11b] Xu ly payload qua ICS20Transfer
    [11c] Mint/unlock token cho nguoi nhan
    [11d] Emit WriteAcknowledgement event
         |
[12] SubscribeEth nhan WriteAck event (PacketType: WriteAck + BlockNumber)
         |
[13] Cho Beacon Chain finality (~12-15 phut, toi da 60 lan * 10s)
    [13a] GetFinalityUpdate(beaconAPIURL)
    [13b] So sanh finalizedBlock >= eventBlock
         |
[14] UpdateEthClient tren Cosmos:
    [14a] GetEthereumClientState -> trustedSlot
    [14b] GetFinalityUpdate -> finalizedSlot
    [14c] Kiem tra sync committee participation >= 2/3
    [14d] updateEthClientWithPeriodCrossing (xu ly chuyen doi period)
    [14e] SendCosmosTxBatch(MsgUpdateClient...)
         |
[15] Xac nhan ETH client.LatestExecutionBlockNumber >= eventBlock
         |
[16] GetEthMembershipProof:
    [16a] Tinh storage key: keccak256(keccak256(ackPath) ++ storageSlot)
    [16b] Goi eth_getProof -> AccountProof + StorageProof
         |
[17] SendCosmosTx(MsgAcknowledgement)
         |
[18] Cosmos xu ly ack -> hoan tat Cosmos->ETH transfer
```

---

## So Do Tong The: Luong ETH -> Cosmos (Ack)

```
Thoi gian -->

=== GIAI DOAN A: ETH gui packet, Cosmos nhan ===

[1] User gui packet tren Ethereum (ICS26Router.SendPacket)
         |
[2] ETH emit SendPacket event
         |
[3] SubscribeEth nhan event
         |
[4] EthPacketToCosmosPacket: chuyen dinh dang Solidity -> Cosmos
         |
[5] InsertPacket vao BatchBuilder (PacketType: Send)
         |
    (* Luu y: Hien tai case Send trong StartLoop chi xu ly huong
       Cosmos->ETH. Luong ETH->Cosmos MsgRecvPacket can bo sung *)

=== GIAI DOAN B: Cosmos ghi ack, ETH nhan (AckPacket) ===

[6] Cosmos xu ly MsgRecvPacket tu ETH
         |
[7] Cosmos emit "write_acknowledgement" event
         |
[8] SubscribeCosmos nhan event (PacketType: Ack + AckBytes)
    [8a] Parse packet hex -> protobuf
    [8b] Parse acknowledgement hex -> protobuf
    [8c] Kiem tra AppAcknowledgements non-empty
         |
[9] InsertPacket vao BatchBuilder
         |
[10] StartLoop nhan batch -> UpdateCosmosClient (ZK proof + SendEthTx)
         |
[11] Voi moi packet Ack trong batch:
    [11a] IbcCommitmentPath(packet, 0x03)   // ack suffix
    [11b] ProvePath -> ABCIQuery lay Merkle proof cho ack
    [11c] Xay dung MembershipMsg (giong Send)
    [11d] ABI encode verifyMembership + cat function selector
         |
[12] Xay dung MsgAckPacket:
    [12a] Packet info (sequence, clients, timeout, payloads)
    [12b] Acknowledgement bytes (ket qua xu ly cua Cosmos)
    [12c] MembershipMsg (proof da encode)
         |
[13] SendEthTx(MsgAckPacket) -> ICS26Router.AckPacket()
         |
[14] ETH xu ly ack -> giai phong escrow -> hoan tat ETH->Cosmos transfer
```

---

## Bang Tong Hop Cac Ham Chinh

### Subscriber & Batch

| Ham | File | Mo ta |
|-----|------|-------|
| `SubscribeCosmos` | `subscriber/event.go:33` | Lang nghe 3 loai su kien Cosmos (send_packet, write_ack, timeout) qua WebSocket |
| `SubscribeEth` | `subscriber/event.go:180` | Lang nghe 4 loai event tu ICS26Router contract tren ETH |
| `EthPacketToCosmosPacket` | `subscriber/event.go:158` | Chuyen packet tu Solidity struct sang Cosmos protobuf |
| `InsertPacket` | `services/batch.go:43` | Them packet vao batch (thread-safe voi sync.Mutex) |
| `CheckBatch` | `services/batch.go:54` | Flush batch khi du thoi gian (3s) hoac so luong (>10) |

### Dieu phoi & Xu ly

| Ham | File | Mo ta |
|-----|------|-------|
| `StartLoop` | `services/main.go:82` | Vong lap chinh: khoi dong subscriber, batch checker, xu ly 4 loai packet |
| `IbcCommitmentPath` | `utils/abci.go:15` | Tao IBC Merkle path: `["ibc", sourceClient + suffix + sequence]` |

### Cap nhat Light Client (Cosmos tren ETH)

| Ham | File | Mo ta |
|-----|------|-------|
| `UpdateCosmosClient` | `services/routine.go:67` | Lay light block, sinh ZK proof, gui UpdateClient len ETH |
| `ExtractValidatorSignature` | `prover/extractor.go:21` | Tim chu ky Ed25519 hop le dau tien tu commit |
| `GenerateProof` | `prover/prover.go:70` | Sinh Groth16 ZK proof cho chu ky Ed25519 (nang nhat) |
| `ProofToBigInts` | `prover/prover.go:166` | Chuyen Groth16 proof sang [8]*big.Int tuong thich Solidity |

### Cap nhat Light Client (ETH tren Cosmos)

| Ham | File | Mo ta |
|-----|------|-------|
| `UpdateEthClient` | `services/routine.go:364` | Lay Beacon finality update, kiem tra sync committee, gui MsgUpdateClient |
| `updateEthClientWithPeriodCrossing` | `services/routine.go:416` | Xu ly cap nhat qua nhieu sync committee period |
| `buildMsgUpdateClient` | `services/routine.go:561` | Xay dung MsgUpdateClient (wasm) cho ETH light client |
| `GetEthereumClientState` | `client/ethereum.go:644` | Query ETH client state tu Cosmos (ABCI -> wasm -> JSON) |
| `GetFinalityUpdate` | `client/ethereum.go` | Goi Beacon API lay finality update moi nhat |

### Chung minh (Proof)

| Ham | File | Mo ta |
|-----|------|-------|
| `ProvePath` | `client/tendermint.go:476` | Truy van ABCI lay ICS23 Merkle proof tu Cosmos state tree |
| `ParseCommitmentProof` | `client/tendermint.go` | Chuyen ICS23 CommitmentProof sang dinh dang Solidity struct |
| `GetEthMembershipProof` | `client/ethereum.go:774` | Goi eth_getProof lay MPT proof tu EVM storage |

### Gui giao dich

| Ham | File | Mo ta |
|-----|------|-------|
| `SendEthTx` | `transaction/handler.go:154` | Ky ETH tx, gui, cho receipt, lay revert reason neu fail |
| `SendCosmosTx` | `transaction/handler.go:558` | Ky Cosmos tx (SIGN_MODE_DIRECT), BroadcastTxCommit |
| `SendCosmosTxBatch` | `transaction/handler.go:730` | Gui nhieu msg trong 1 Cosmos tx (gas = base * len) |
| `CosmosSignerAddress` | `transaction/handler.go:548` | Derive bech32 address tu COSMOS_PRIVATE_KEY |
| `queryAccountInfo` | `transaction/handler.go:917` | Query account number + sequence tu chain |

### Tong hop loai Packet

| PacketType | Nguon | Huong | Xu ly |
|-----------|-------|-------|-------|
| `Send` (0) | Cosmos `send_packet` | Cosmos -> ETH | RecvPacket tren ETH (Cosmos Merkle proof) |
| `Ack` (1) | Cosmos `write_acknowledgement` | ETH -> Cosmos -> ETH | AckPacket tren ETH (Cosmos Merkle proof) |
| `Timeout` (2) | Cosmos `timeout_packet` | - | NonMembership proof tren ETH |
| `WriteAck` (3) | ETH `WriteAcknowledgement` | Cosmos -> ETH -> Cosmos | MsgAcknowledgement tren Cosmos (ETH storage proof) |
