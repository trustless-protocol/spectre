# So sánh: Trước và sau khi support multi-validator (Cosmos → ETH updateClient)

Doc này mô tả chính xác input chảy từ Cosmos block đến Ethereum Groth16Verifier, ở hai trạng thái:

- **Trước (single-sig)**: chỉ chứng minh 1 validator ký block.
- **Sau (multi-sig, approach B + hash-aggregate)**: chứng minh N validator ký cùng block, N đạt 2/3 voting power, canonical vote rebuild trong circuit, hash-aggregate 1 Fr public input.

---

## 1. Tổng quan luồng Cosmos → ETH

```
Cosmos chain                  Relayer (Go)                    Ethereum
─────────────                 ──────────────                   ──────────
Tendermint block                                               Groth16ICS07Tendermint
 ├─ SignedHeader              ① query block                      ├─ updateClient(msg)
 │   ├─ Header                ② extract sigs                        │
 │   └─ Commit.Signatures     ③ build witness                       ├─ _verifyBatchAndQuorum
 │       (per-validator       ④ Groth16.Prove                       │   ├─ WrapperVerifier.verifyBatchProof
 │        Ed25519 sigs +      ⑤ build MsgUpdateClient               │   │   ├─ rehash witness → digest
 │        timestamps)         ⑥ SendEthTx                           │   │   └─ Groth16Verifier_N{N}
 └─ ValidatorSet                                                    │   │        (pairing check)
                                                                    │   └─ voting power quorum check
                                                                    └─ commit new consensus state
```

---

## 2. Input chi tiết ở từng bước

### 2.1 Trước (single-sig)

| Bước | Input | Nguồn |
|---|---|---|
| Sign bytes | `voteBytes = Commit.VoteSignBytes(chainID, valIdx)` | cometbft canonical vote của 1 validator |
| Prover circuit public inputs (Fr) | `R.X, R.Y, S, H, A.X, A.Y` limb-packed = 24 Fr, + msg 256 bytes | đều là 1 slot |
| MsgUpdateClient bổ sung | (none — chỉ proof+commitments+commitmentPok) | — |
| On-chain verify | `WrapperVerifier.verifyProof(proof, commitments, commitmentPok, sig, pubkey, signMessage)` | 1 sig, 1 msg bytes |

Lỗ hổng: không check 2/3 voting power; 1 validator ký là đủ → **không đạt trust model IBC**.

### 2.2 Sau (multi-sig, approach B + hash-aggregate)

| Bước | Input | Nguồn |
|---|---|---|
| Sign bytes | mỗi validator có riêng `voteBytes_i = Commit.VoteSignBytes(chainID, valIdx_i)` | khác nhau ở Timestamp_i |
| Circuit witness | per-slot: `(R_i, A_i, S_i, TsSec_i, TsNanos_i)`; shared: `height, round, blockIDHash, partSetTotal, partSetHash, chainID, chainIDLen` | private inputs |
| Circuit public inputs | `Hash [32]uints.U8` = SHA-256 của toàn bộ witness (1 Fr public input duy nhất, mở rộng thành 32 byte) | — |
| MsgUpdateClient bổ sung | `bucket, signerIndices[], signatures[], signerPubkeys[], timestampSeconds[], timestampNanos[]` | relayer build |
| On-chain verify | `WrapperVerifier.verifyBatchProof(bucket, proof, commitments, commitmentPok, signatures, pubkeys, tsSec, tsNanos, sharedBlock)` | rehash trên chain → SHA-256 digest → dispatch bucket verifier |

Đạt trust model: Solidity check `Σ votingPower(unique signerIndices) ≥ ⌈2/3 · totalVotingPower⌉ + 1` và từng `signerPubkeys[i] == validators[signerIndices[i]].pubKey`.

---

## 3. Circuit

### 3.1 Trước (`PreHashCircuit` — single-sig)

```go
type PreHashCircuit struct {
    Sig  eddsa.Signature[Fp, Fr]  `gnark:",public"`
    Msg  [MaxMsgLen]uints.U8      `gnark:",public"`
    Pub  eddsa.PublicKey[Fp, Fr]  `gnark:",public"`
}
func (c *PreHashCircuit) Define(api frontend.API) error {
    return eddsa.Verify(api, c.Sig, c.Msg[:], c.Pub, Config{FromWei: false})
}
```

- **Số slot**: 1 cố định.
- **Public inputs**: ~24 Fr + `MaxMsgLen` Fr (~280 Fr).
- **Constraint count**: ~400–500k.
- **Verifier size**: vài KB (không vấn đề EIP-170).
- **Msg**: đưa nguyên canonical vote bytes vào circuit làm public input.

### 3.2 Sau (`BatchCircuit` — multi-sig)

```go
type BatchCircuit[Base, Scalars emulated.FieldParams] struct {
    Hash [32]uints.U8 `gnark:",public"`   // SHA-256 commitment to all witness

    Sig       []eddsa.Signature[Base, Scalars]   // private, length = N (bucket size)
    Pub       []eddsa.PublicKey[Base, Scalars]   // private
    TsSeconds []frontend.Variable                // private
    TsNanos   []frontend.Variable                // private

    // Shared block data (private)
    Height, Round       frontend.Variable
    BlockIDHash         [32]uints.U8
    PartSetTotal        frontend.Variable
    PartSetHash         [32]uints.U8
    ChainID             [canonvote.MaxChainIDLen]uints.U8
    ChainIDLen          frontend.Variable
}

func (c *BatchCircuit) Define(api) {
    // 1) Ed25519 batch verify qua VerifyBatchWithCanonicalVote:
    //    - Rebuild canonical vote trong circuit từ (shared, ts_i) cho từng slot
    //    - FixedLengthSum SHA-512 trên msg reconstructed
    //    - Fiat-Shamir aggregate N Ed25519 eqs → 1 ecip check
    // 2) (optional — hash-aggregate) Hash witness bytes → assert == Hash public
}
```

- **Số slot**: fixed theo bucket N ∈ `{3, 4, 8, 16, 32, 64, 128}` — **một circuit per bucket**.
- **Public inputs**: 32 Fr (chỉ 1 SHA-256 digest) → verifier contract ~5-8KB, nằm dưới EIP-170 24KB.
- **Constraint count** (xấp xỉ):
  - N=3: ~1.57M
  - N=4: ~2.04M (không hash-aggregate) / ~2.21M (có hash-aggregate)
  - N=8: ~4M+
- **Canonical vote reconstruction**: proto-encode cv bytes (Type|Height|Round|BlockID|Timestamp|ChainID) bằng gadget `canonvote.CanonicalVoteBytes`, byte-match cometbft `VoteSignBytes`.
- **Hash-aggregate**: giảm public inputs từ ~22·N + 116 Fr (approach B gốc) → 32 Fr cố định. Verifier Solidity rehash calldata trên chain.

### 3.3 Khác biệt then chốt

| Đặc điểm | Trước | Sau |
|---|---|---|
| Số signatures | 1 | N đến 128 (bucket) |
| Trust model | không đạt IBC | 2/3 voting power |
| Msg trong public input | có, ~256 Fr | không — rebuild trong circuit |
| Public input count | ~280 Fr | 32 Fr (hash-aggregate) |
| Verifier Solidity size | vài KB | ~5-8KB (constant theo N) |
| Circuit/(pk,vk) artifact | 1 bộ | 7 bộ (mỗi bucket) |
| Padding | không cần | copy slot 0 (⚠ đang gây lỗi ecip duplicate — xem mục 7) |

---

## 4. Relayer logic

### 4.1 Trước (`services/routine.go`)

```go
// Extract 1 signature
sig, pub, msg, err := prover.ExtractValidatorSignature(lightBlock, chainID)
// Prove
proof, commitments, commitmentPok, err := w.Prover.GenerateProof(sig, pub, msg)
// Build & send
msg := MsgUpdateClient{ClientState, TrustedConsensusState, ProposedHeader, Time, Proof, Commitments, CommitmentPok}
w.TxHandler.SendEthTx(ctx, msg)
```

### 4.2 Sau

```go
// 1) Extract TẤT CẢ sig hợp lệ, sort theo voting power desc, chọn prefix đủ 2/3
extracted, err := prover.ExtractValidatorSignatures(lightBlock, chainID)
// extracted.Signatures = slice sigs, extracted.Shared = block data

// 2) Tính bucket, pad slot dư với copy slot 0 (tạm thời — sẽ đổi sang distinct)
bucket := smallestBucketGEQ(len(extracted.Signatures))
paddedSigs := PadSigsToBucket(extracted.Signatures, bucket)

// 3) Hash witness off-chain (đồng bộ SHA-256 với circuit)
hash := ComputeWitnessHash(extracted.Shared, paddedSigs)

// 4) Prove
_, proof, commitments, commitmentPok, err := w.Prover.GenerateProof(extracted.Shared, extracted.Signatures)

// 5) Pack msg với arrays per-slot
msg := MsgUpdateClient{
    ClientState, TrustedConsensusState, ProposedHeader, Time,
    Proof, Commitments, CommitmentPok,
    Bucket:           uint16(bucket),
    SignerIndices:    []uint32{...},   // index trong validator set
    Signatures:       [][2][32]byte{}, // (R, S) per slot
    SignerPubkeys:    [][32]byte{},
    TimestampSeconds: []uint64{},
    TimestampNanos:   []uint32{},
}
```

### 4.3 Khác biệt prover

| Phần | Trước | Sau |
|---|---|---|
| Extractor | trả 1 sig đầu tiên | trả tất cả sigs + shared block data |
| Witness builder | 1 slot | N slots (padded đến bucket) |
| Prover artifact | 1 `(r1cs, pk, vk)` | 7 bộ, load theo `binDir/n{N}/` |
| Bucket selection | n/a | `smallestBucketGEQ(quorumCount)` |

---

## 5. Solidity logic

### 5.1 `MsgUpdateClient` message layout

**Trước**:
```solidity
struct MsgUpdateClient {
    ClientState clientState;
    ConsensusState trustedConsensusState;
    Header proposedHeader;
    uint128 time;
    uint256[8] proof;
    uint256[2] commitments;
    uint256[2] commitmentPok;
}
```

**Sau**: thêm 6 field cho multi-validator:
```solidity
struct MsgUpdateClient {
    // ...trước giữ nguyên...
    uint16 bucket;                  // circuit bucket đã dùng
    uint32[] signerIndices;         // index validator set cho từng slot
    bytes32[2][] signatures;        // (R, S) per slot
    bytes32[] signerPubkeys;
    uint64[] timestampSeconds;
    uint32[] timestampNanos;
}
```

### 5.2 `Groth16ICS07Tendermint._verifyBatchAndQuorum`

**Trước** (pseudo):
```solidity
verifier.verifyProof(proof, commitments, commitmentPok, sig, pubkey, msg);
// không check voting power
```

**Sau**:
```solidity
// 1) Sum voting power của UNIQUE signerIndices
uint64 power = 0;
bool[] seen = new bool[](validators.length);
for (uint i = 0; i < bucket; i++) {
    uint32 idx = signerIndices[i];
    if (seen[idx]) continue;               // skip padding duplicate
    seen[idx] = true;
    require(validators[idx].pubKey == signerPubkeys[i], PubkeyMismatch());
    power += validators[idx].votingPower;
}
require(power * 3 > totalVotingPower * 2, InsufficientVotingPower());

// 2) Verify Groth16 batch proof — dispatch bucket
WRAPPER.verifyBatchProof(
    bucket, proof, commitments, commitmentPok,
    signatures, signerPubkeys, timestampSeconds, timestampNanos,
    sharedBlock  // derived từ proposedHeader
);
```

### 5.3 `WrapperVerifier.verifyBatchProof` (mới)

```solidity
function verifyBatchProof(
    uint16 bucket,
    uint256[8] calldata proof,
    uint256[2] calldata commitments,
    uint256[2] calldata commitmentPok,
    bytes32[2][] calldata signatures,
    bytes32[] calldata pubkeys,
    uint64[] calldata tsSec,
    uint32[] calldata tsNanos,
    SharedBlock calldata shared
) external view returns (bool) {
    // 1) Require array lengths == bucket
    // 2) Rehash calldata: digest = SHA-256(signatures || pubkeys || tsSec || tsNanos || shared_fields)
    //    - layout khớp chính xác circuit.go + hash_witness.go
    // 3) Dispatch to Groth16Verifier_N{bucket}:
    //    selector = bytes4(keccak256("verifyProof(uint256[8],uint256[2],uint256[2],uint256[32])"))
    //    publicInputs[32] = digest bytes as uint256
    //    staticcall(verifier, abi.encodePacked(selector, proof, commitments, commitmentPok, publicInputs))
}
```

---

## 6. Input đầy đủ từ Cosmos đến Ethereum verify

Trong mục này mỗi step **đặt tên cho output** và step sau **khai báo rõ dùng output nào làm input** — để thấy chính xác dòng dữ liệu không đứt gãy.

### Step 1 — Relayer đọc từ Cosmos

**Input**: height H = `latestBlockHeight`, chainID từ config.

**Gọi**: `ctx.CosmosClient().LightBlock(ctx, H)` qua Tendermint RPC.

**Output** (gọi là `LB`):
```
LB = LightBlock {
  SignedHeader {
    Header {
      ChainID = "test-ibc-eth"
      Height  = 1082
      Time    = 2026-04-22T05:16:29.993821000Z
      AppHash = 0x...
      ValidatorsHash, NextValidatorsHash, ...
    }
    Commit {
      Height       = 1082             ← LB.Commit.Height
      Round        = 0                ← LB.Commit.Round
      BlockID {
        Hash = 0xa448... (32B)        ← LB.Commit.BlockID.Hash
        PartSetHeader {
          Total = 1                   ← LB.Commit.BlockID.PartSetHeader.Total
          Hash  = 0xb9d3... (32B)     ← LB.Commit.BlockID.PartSetHeader.Hash
        }
      }
      Signatures[]                    ← LB.Commit.Signatures[i]
        [0]: BlockIDFlagAbsent, Timestamp=time.Time{}, Signature=nil
        [1]: BlockIDFlagCommit, Signature=sig_1 (64B), Timestamp=t_1
        [2]: BlockIDFlagCommit, Signature=sig_2 (64B), Timestamp=t_2
        [3]: BlockIDFlagCommit, Signature=sig_3 (64B), Timestamp=t_3
    }
  }
  ValSet.Validators[i] = { PubKey, VotingPower, Address }    ← LB.ValSet
}
```

`LB` được pass nguyên sang Step 2 và cả sang Step 4 (làm `proposedHeader` trong MsgUpdateClient).

### Step 2 — Extractor (`relayer/prover/extractor.go`)

**Input**: `LB` từ Step 1, `chainID` từ config.

**Logic**:
1. Với mỗi `LB.Commit.Signatures[i]`:
   - Skip nếu `BlockIDFlag == Absent` (validator không ký).
   - Compute `voteBytes = LB.Commit.VoteSignBytes(chainID, i)` — canonical vote mà validator i đã ký.
   - Verify `ed25519.Verify(LB.ValSet.Validators[i].PubKey, voteBytes, sig.Signature)`.
   - Nếu pass → append vào `candidates`.
2. `sort(candidates, by Power desc)`.
3. Lấy prefix nhỏ nhất đủ `power_sum ≥ ⌈2/3 · totalPower⌉ + 1` → đó là `quorumPrefix`.

**Output** (gọi là `EXT`):
```
EXT = ExtractorResult {
  Shared = SharedBlockData {
    Height       = LB.Commit.Height                             (= 1082)
    Round        = int64(LB.Commit.Round)                        (= 0)
    BlockIDHash  = LB.Commit.BlockID.Hash                        (32B)
    PartSetTotal = LB.Commit.BlockID.PartSetHeader.Total          (= 1)
    PartSetHash  = LB.Commit.BlockID.PartSetHeader.Hash           (32B)
    ChainID      = chainID                                        ("test-ibc-eth")
  }
  Signatures = [quorumPrefix]ValidatorSignature {
    // Với mỗi validator i được chọn vào prefix:
    Signature        = LB.Commit.Signatures[i].Signature          (64B R||S)
    PublicKey        = LB.ValSet.Validators[i].PubKey.Bytes()     (32B compressed)
    Index            = i                                           (index trong ValSet)
    Power            = LB.ValSet.Validators[i].VotingPower
    TimestampSeconds = LB.Commit.Signatures[i].Timestamp.Unix()
    TimestampNanos   = LB.Commit.Signatures[i].Timestamp.Nanosecond()
  }
}
```

Lưu ý: `EXT.Shared` và `EXT.Signatures[i].Timestamp*` đủ để **tái tạo lại `voteBytes_i`** mà validator i đã ký — Step 3 và circuit sẽ dùng điều này.

### Step 3 — Prover (`relayer/prover/prover.go`)

**Input**: `EXT` từ Step 2.

**Bước 3a — Bucket selection + padding**:
```
bucket := smallestBucketGEQ(len(EXT.Signatures))     // ví dụ quorumPrefix=3 → bucket=4
paddedSigs := PadSigsToBucket(EXT.Signatures, bucket)
  // paddedSigs[0..len(EXT.Signatures)-1] = EXT.Signatures (copy từng slot)
  // paddedSigs[len(EXT.Signatures)..bucket-1] = paddedSigs[0] (⚠ duplicate)
```

`paddedSigs` là slice `[bucket]ValidatorSignature` độ dài cố định = bucket N của circuit.

**Bước 3b — Compute witness hash off-chain** (để làm public input):
```
H_BYTES = ComputeWitnessHash(EXT.Shared, paddedSigs) = SHA-256 của chuỗi byte:
  for i = 0 .. bucket-1:
      paddedSigs[i].Signature[:32]                       (R, 32B)
      paddedSigs[i].Signature[32:]                       (S, 32B)
      paddedSigs[i].PublicKey                            (A, 32B)
      uint64BE(paddedSigs[i].TimestampSeconds)           (8B)
      uint32BE(paddedSigs[i].TimestampNanos)             (4B)
  uint64BE(EXT.Shared.Height)                            (8B)
  uint32BE(EXT.Shared.Round)                             (4B)
  EXT.Shared.BlockIDHash                                 (32B)
  uint32BE(EXT.Shared.PartSetTotal)                      (4B)
  EXT.Shared.PartSetHash                                 (32B)
  EXT.Shared.ChainID padded đến MaxChainIDLen            (48B)
  uint32BE(len(EXT.Shared.ChainID))                      (4B)
→ H_BYTES = 32-byte SHA-256 digest
```

Layout này **phải giống hệt** in-circuit reconstruction và Solidity `_hashWitness` — nếu khác 1 byte → proof không verify được.

**Bước 3c — Build circuit assignment** (điền witness vào `BatchCircuit`):
```
assignment = BatchCircuit {
  // PUBLIC input (duy nhất):
  Hash[0..31] = H_BYTES[0..31]          ← từ 3b

  // PRIVATE inputs — per slot i ∈ [0, bucket):
  Sig[i].R = emulated.AffinePoint{ valueOf(rX), valueOf(rY) }
            where (rX, rY) = utils.DecompressPoint(paddedSigs[i].Signature[:32])
  Sig[i].S = emulated.ValueOf( utils.ScalarToBigInt( paddedSigs[i].Signature[32:] ) )
  Pub[i].A = emulated.AffinePoint{ valueOf(aX), valueOf(aY) }
            where (aX, aY) = utils.DecompressPoint(paddedSigs[i].PublicKey)
  TsSeconds[i] = paddedSigs[i].TimestampSeconds
  TsNanos[i]   = paddedSigs[i].TimestampNanos

  // PRIVATE inputs — shared:
  Height       = EXT.Shared.Height
  Round        = EXT.Shared.Round
  BlockIDHash  = EXT.Shared.BlockIDHash              (32 uints.U8)
  PartSetTotal = EXT.Shared.PartSetTotal
  PartSetHash  = EXT.Shared.PartSetHash              (32 uints.U8)
  ChainID      = pad(EXT.Shared.ChainID, MaxChainIDLen, 0)    (MaxChainIDLen uints.U8)
  ChainIDLen   = len(EXT.Shared.ChainID)
}
```

**Bước 3d — Prove**:
```
witness      = frontend.NewWitness(assignment, BN254)
pk, r1cs, vk = p.byBucket[bucket]     // load từ binDir/n{N}/
proof        = groth16.Prove(r1cs, pk, witness)
groth16.Verify(proof, vk, witness.Public())   // local sanity check
```

**Bước 3e — Extract proof components**:
```
(PROOF, COMMITMENTS, COMMITMENT_POK) = ProofToBigInts(proof)
  PROOF          = [8]*big.Int   ( Ar.X, Ar.Y, Bs.X.A1, Bs.X.A0, Bs.Y.A1, Bs.Y.A0, Krs.X, Krs.Y )
  COMMITMENTS    = [2]*big.Int   ( commitments[0].X, commitments[0].Y )
  COMMITMENT_POK = [2]*big.Int   ( CommitmentPok.X, CommitmentPok.Y )
```

**Output của Step 3**: `(bucket, PROOF, COMMITMENTS, COMMITMENT_POK)`.

Note: `H_BYTES` **không** cần gửi lên chain — Solidity tự compute lại từ calldata (Step 7). Nhưng calldata phải chứa đúng paddedSigs + shared đã dùng — tức Step 4 phải đồng bộ với `paddedSigs` ở Step 3a.

### Step 4 — Relayer build MsgUpdateClient

```
MsgUpdateClient {
  ClientState = {ChainId, TrustLevel, LatestHeight, TrustingPeriod, UnbondingPeriod, IsFrozen, ZkAlgorithm}
  TrustedConsensusState = {Timestamp, Root, NextValidatorsHash}
  ProposedHeader = latestLightBlock.IntoHeader(trustedLightBlock)
    (gồm cả Commit.Signatures[], Validators[], TrustedNextValidatorSet...)
  Time = time.Now().UnixNano() as big.Int (uint128 ABI)

  Proof         = proof
  Commitments   = commitments
  CommitmentPok = commitmentPok

  Bucket           = uint16(bucket)
  SignerIndices    = [bucket]uint32    ← từ paddedSigs[i].Index
  Signatures       = [bucket][2][32]byte (R, S)
  SignerPubkeys    = [bucket][32]byte
  TimestampSeconds = [bucket]uint64
  TimestampNanos   = [bucket]uint32
}
```

### Step 5 — TransactionHandler.SendEthTx

```
data := abi.Pack(MsgUpdateClient)
tx := ics07Tendermint.UpdateClient(auth, data)
```

### Step 6 — On-chain `Groth16ICS07Tendermint.updateClient`

```
1) Decode msg
2) Check client state frozen/period/height consistency
3) _verifyBatchAndQuorum:
     a) Derive SharedBlock từ msg.proposedHeader:
          height       = proposedHeader.signedHeader.commit.height
          round        = uint32(proposedHeader.signedHeader.commit.round)
          blockIDHash  = proposedHeader.signedHeader.commit.blockId.hashData
          partSetTotal = proposedHeader.signedHeader.commit.blockId.partSetHeader.total
          partSetHash  = proposedHeader.signedHeader.commit.blockId.partSetHeader.hashData
          chainID      = proposedHeader.signedHeader.header.chainId (bytes)
     b) Voting power check:
          power = 0; seen = bool[validators.length]
          for i in [0, bucket):
             idx = signerIndices[i]
             if !seen[idx]:
                require validators[idx].pubKey == signerPubkeys[i]
                power += validators[idx].votingPower
                seen[idx] = true
          require power * 3 > totalVotingPower * 2
     c) WrapperVerifier.verifyBatchProof(bucket, proof, commitments, commitmentPok,
                                         signatures, signerPubkeys, tsSec, tsNanos, shared)
```

### Step 7 — `WrapperVerifier.verifyBatchProof`

```
1) length check: signatures/pubkeys/tsSec/tsNanos đều == bucket, chainID ≤ MAX_CHAIN_ID_LEN.
2) Rehash calldata:
     digest = SHA-256(
       for each slot i ∈ [0, bucket):
         signatures[i][0] (R, 32B) || signatures[i][1] (S, 32B) || signerPubkeys[i] (A, 32B)
         || tsSec[i] (8B BE) || tsNanos[i] (4B BE)
       height (8B BE) || round (4B BE) || blockIDHash (32B) || partSetTotal (4B BE)
       || partSetHash (32B) || chainID (pad đến MAX_CHAIN_ID_LEN) || chainIDLen (4B BE)
     )
3) publicInputs[32] = digest bytes
4) bucketVerifier = buckets[bucket]
5) staticcall(bucketVerifier.verifier,
              abi.encodePacked(bucketVerifier.selector,
                               proof, commitments, commitmentPok, publicInputs))
6) return success
```

### Step 8 — `Groth16Verifier_N{bucket}.verifyProof`

```
verifyProof(uint256[8] proof, uint256[2] commitments, uint256[2] commitmentPok, uint256[32] input)
  (auto-generated by gnark)
  thực hiện BN254 pairing check:
     e(A, B) == e(alpha_G1, beta_G2) * e(Input(vk_x, input), gamma_G2) * e(C, delta_G2)
  Trả revert nếu proof không hợp lệ, return bình thường nếu pass.
```

Proof hợp lệ ↔ tồn tại witness (R_i, S_i, A_i, ts_i, shared) sao cho:
- `SHA-256(witness bytes) == digest` (binds to calldata)
- Batch Ed25519 verify qua ecip: mỗi (R_i, S_i, A_i) ký trên canonical vote rebuild từ (shared, ts_i)

→ Ethereum tin chắc N validator trong validator set đã ký block với tổng voting power ≥ 2/3.

---

## 7. Vấn đề đang debug (2026-04-23)

**Padding bằng copy slot 0 phá ecip batch verify**:
- ECIP constructFunction pair points qua `line(p, q)`. Khi p == q (duplicate witness do padding), đi vào nhánh tangent doubling — tạo divisor với double-zero tại P_0 không khớp formula verifier kỳ vọng.
- Empirical: smoke test (4 sigs distinct) pass; E2E (3 real + 1 padded duplicate) fail tại constraint cuối của ecip aggregate check.

**Các option đang cân nhắc**:

| Option | Mô tả | Trade-off |
|---|---|---|
| A | Extractor trả all valid sigs; pad bằng sigs validator khác (distinct) | Không support chain nhỏ khi T < bucket |
| B | Thêm bucket nhỏ `{3, 4, 8, 16, ...}` để giảm pad | Nhiều circuit setup hơn |
| C | Thêm `active[]` flag vào `VerifyBatchWithCanonicalVote`, scalar 0 cho padded → ecip selector skip | Modify ecip-gnark, +N Fr public input |
| D | Fix `constructFunction` trong ecip-gnark để handle duplicate đúng | Rủi ro cao, ảnh hưởng cả single-sig |

Option C được đánh giá tốt nhất. Chưa implement.

---

## 8. Tổng kết

| Khía cạnh | Trước | Sau |
|---|---|---|
| Trust model | 1 sig (không đạt IBC) | 2/3 voting power qua N sigs |
| Circuit | fixed 1 slot, msg public input | N slots theo bucket, canonvote rebuild, hash-aggregate |
| Prover artifacts | 1 bộ | 7 bộ (per bucket) |
| Public inputs verifier | ~280 Fr | 32 Fr (SHA-256 digest) |
| Verifier Solidity | vài KB | ~5-8KB/bucket (đạt EIP-170) |
| MsgUpdateClient | 7 field | 13 field (thêm bucket + 5 arrays) |
| On-chain check | chỉ proof verify | + voting power quorum + pubkey binding |
| Blocker hiện tại | — | padding duplicate → cần Option C |
