# Luồng Prove/Verify Mới Với Hash-Aggregate

## 1. Mục tiêu

Hash-aggregate dùng một `sha256` digest để bind toàn bộ dữ liệu mà circuit dùng, thay vì expose từng field làm Groth16 public input riêng lẻ.

Vấn đề cũ:

```text
publicInputs = perSignerInputs + sharedBlockInputs
L = 22 * bucket + 116

bucket = 4  -> 204 public inputs
bucket = 8  -> 292 public inputs
bucket = 16 -> 468 public inputs
```

Groth16 verifier Solidity phình theo số public inputs vì generated verifier phải embed verification-key points và chạy public-input MSM cho từng input.

Hash-aggregate đổi thành:

```text
publicInputs = bytes(sha256(witnessBytes))
L = 32
```

Tức verifier chỉ nhận 32 field elements, mỗi element là 1 byte của SHA-256 digest.

## 2. Ý tưởng chính

Có 3 nơi tính cùng một hash:

```text
Go prover:
  H = sha256(witness bytes)

Circuit:
  H' = sha256(private witness bytes)
  assert H' == H

Solidity WrapperVerifier:
  H_onchain = sha256(calldata bytes)
  đưa H_onchain vào Groth16 verifier
```

Nếu proof pass thì:

```text
hash(private witness trong circuit) == hash(calldata on-chain)
```

Với giả định SHA-256 không collision, circuit đã verify đúng dữ liệu mà contract thấy.

## 3. Input/Output Theo Tầng

### 3.1. Relayer/Go Prover

Input:

```text
shared:
  height
  round
  blockIDHash
  partSetTotal
  partSetHash
  chainID

signatures:
  signer index
  signature R/S
  public key
  timestamp seconds
  timestamp nanos
  voting power metadata ở ngoài flow proof
```

Output gửi lên on-chain:

```text
bucket
proof[8]
commitments[2]
commitmentPok[2]
signatures[]          // wrapper hash lại để bind với private witness
signerPubkeys[]       // light client dùng để quorum-check, wrapper cũng hash lại
timestampSeconds[]
timestampNanos[]
shared block fields
```

### 3.2. Circuit

Public input thật của Groth16:

```text
Hash[32]
```

Private witness:

```text
Sig[]                 // Ed25519 R/S
Pub[]                 // Ed25519 compressed pubkey, decompressed inside assignment
TsSeconds[]
TsNanos[]
Height
Round
BlockIDHash
PartSetTotal
PartSetHash
ChainID[48]
ChainIDLen
```

Circuit làm 2 việc:

```text
1. Rebuild witness bytes từ private witness rồi sha256.
2. Verify Ed25519 signatures against canonical-vote content.
```

### 3.3. Solidity WrapperVerifier

Input từ `Groth16ICS07Tendermint`:

```solidity
verifyBatchProof(
    bucket,
    proof,
    commitments,
    commitmentPok,
    signatures,
    signerPubkeys,
    timestampSeconds,
    timestampNanos,
    shared
)
```

Output:

```text
true  -> Groth16 proof hợp lệ
false -> proof/verifier staticcall fail
```

Wrapper không verify Ed25519 trực tiếp. Wrapper chỉ:

```text
1. Recompute sha256 từ calldata.
2. Convert digest thành uint256[32].
3. Dispatch sang per-bucket Groth16 verifier.
```

### 3.4. Groth16ICS07Tendermint

Trước khi gọi wrapper, contract vẫn check quorum:

```text
signerIndices length == bucket
signatures length == bucket
signerPubkeys length == bucket
timestamps length == bucket

for each signer:
  signer index nằm trong validator set
  signerPubkeys[i] == validatorSet[signerIndices[i]].pubKey
  cộng voting power unique signer

require accumulatedVotingPower > 2/3 totalVotingPower
```

Sau đó mới gọi:

```text
WrapperVerifier.verifyBatchProof(...)
```

## 4. Witness Hash Layout

Hash layout phải giống hệt ở Go, circuit, và Solidity. Chỉ cần khác endian, padding, hoặc thứ tự field là proof sẽ fail.

```text
for each slot i in [0, bucket):
  R                32 bytes   signatures[i][0], Ed25519 R compressed, raw LE bytes
  S                32 bytes   signatures[i][1], Ed25519 S scalar, raw LE bytes
  A                32 bytes   signerPubkeys[i], Ed25519 compressed pubkey, raw LE bytes
  ts_sec            8 bytes   timestampSeconds[i], big-endian uint64
  ts_nanos          4 bytes   timestampNanos[i], big-endian uint32

shared:
  height            8 bytes   big-endian uint64
  round             4 bytes   big-endian uint32
  blockIDHash      32 bytes
  partSetTotal      4 bytes   big-endian uint32
  partSetHash      32 bytes
  chainID          48 bytes   right-zero-padded
  chainIDLen        4 bytes   actual chainID length, big-endian uint32
```

Formula:

```text
H = sha256(
  slot[0] || slot[1] || ... || slot[bucket-1] ||
  height ||
  round ||
  blockIDHash ||
  partSetTotal ||
  partSetHash ||
  paddedChainID48 ||
  chainIDLen
)
```

Lưu ý: hash hiện tại không chứa toàn bộ validator set. Nó chỉ chứa `signerPubkeys` của các signer slots. Toàn bộ validator set vẫn nằm trong `MsgUpdateClient.proposedHeader.validatorSet` để `Groth16ICS07Tendermint` check quorum/voting power.

## 5. Vì Sao Signature Vẫn Nằm Trong Hash?

`Sig` là private witness trong circuit, không còn là Groth16 public input.

Nhưng nếu không hash `Sig`, có thể xảy ra:

```text
on-chain calldata:
  sig = A

circuit private witness:
  sig = B
```

Contract không đọc được private witness, nên không biết circuit đã dùng `B`.

Khi hash chứa `R || S`:

```text
Wrapper hash calldata sig A -> H
Circuit hash private sig -> H'
assert H == H'
```

Proof chỉ pass nếu `sig private` trong circuit đúng với `sig` trong calldata.

## 6. Code Mô Phỏng

### 6.1. Go Prover

```go
func GenerateProof(shared SharedBlockData, sigs []ValidatorSignature) {
    bucket := SmallestBucketGEQ(len(sigs))
    paddedSigs := padSigsToBucket(sigs, bucket)

    // Public input của Groth16.
    hash := sha256(encodeWitnessBytes(shared, paddedSigs))

    // Full witness: hash public, còn lại private.
    assignment := BatchCircuit{
        Hash: hash,

        Sig: paddedSigs.Signatures,
        Pub: paddedSigs.PubKeys,
        TsSeconds: paddedSigs.TimestampSeconds,
        TsNanos: paddedSigs.TimestampNanos,

        Height: shared.Height,
        Round: shared.Round,
        BlockIDHash: shared.BlockIDHash,
        PartSetTotal: shared.PartSetTotal,
        PartSetHash: shared.PartSetHash,
        ChainID: paddedChainID,
        ChainIDLen: len(shared.ChainID),
    }

    proof := groth16.Prove(r1cs, pk, assignment)
    return proof
}
```

Code thật:

```text
relayer/prover/hash_witness.go
relayer/prover/prover.go
```

### 6.2. Circuit

```go
type BatchCircuit struct {
    Hash [32]uints.U8 `gnark:",public"`

    Sig []Signature
    Pub []PublicKey
    TsSeconds []frontend.Variable
    TsNanos []frontend.Variable

    Height frontend.Variable
    Round frontend.Variable
    BlockIDHash [32]uints.U8
    PartSetTotal frontend.Variable
    PartSetHash [32]uints.U8
    ChainID [48]uints.U8
    ChainIDLen frontend.Variable
}

func (c *BatchCircuit) Define(api frontend.API) error {
    bytes := serializePrivateWitness(c)

    digest := sha256(bytes)
    for i := 0; i < 32; i++ {
        assert(c.Hash[i] == digest[i])
    }

    verifyEd25519Batch(
        c.Sig,
        c.Pub,
        canonicalVote(c.Height, c.Round, c.BlockIDHash, ...),
        c.TsSeconds,
        c.TsNanos,
    )
}
```

Điểm quan trọng:

```text
Hash là public.
Sig/Pub/Header/Timestamp là private.
Nhưng private data bị bind bởi Hash.
```

Code thật:

```text
relayer/prover/circuit.go
```

### 6.3. Solidity WrapperVerifier

```solidity
function verifyBatchProof(
    uint16 bucket,
    uint256[8] calldata proof,
    uint256[2] calldata commitments,
    uint256[2] calldata commitmentPok,
    bytes32[2][] calldata signatures,
    bytes32[] calldata pubkeys,
    uint64[] calldata timestampSeconds,
    uint32[] calldata timestampNanos,
    SharedBlock calldata shared
) external view returns (bool) {
    require(signatures.length == bucket);
    require(pubkeys.length == bucket);
    require(timestampSeconds.length == bucket);
    require(timestampNanos.length == bucket);

    bytes memory witnessBytes = encodeWitnessBytesFromCalldata(...);
    bytes32 h = sha256(witnessBytes);

    uint256[32] memory input;
    for (uint256 i = 0; i < 32; i++) {
        input[i] = uint256(uint8(h[i]));
    }

    return verifier.staticcall(
        abi.encodePacked(
            selector,
            proof,
            commitments,
            commitmentPok,
            input
        )
    );
}
```

Code thật:

```text
contracts/utils/WrapperVerifier.sol
```

### 6.4. Light Client Quorum Check

```solidity
function _verifyBatchAndQuorum(MsgUpdateClient memory msg_) internal {
    for i in msg_.signerIndices:
        validator = validatorSet[msg_.signerIndices[i]]

        require(validator.pubKey == msg_.signerPubkeys[i])

        if first time seeing this validator:
            accumulated += validator.votingPower

    require(accumulated * 3 > totalVotingPower * 2)

    shared = buildSharedBlock(msg_.proposedHeader)

    require(
        wrapper.verifyBatchProof(
            msg_.bucket,
            msg_.proof,
            msg_.commitments,
            msg_.commitmentPok,
            msg_.signatures,
            msg_.signerPubkeys,
            msg_.timestampSeconds,
            msg_.timestampNanos,
            shared
        )
    )
}
```

Code thật:

```text
contracts/light-clients/Groth16ICS07Tendermint.sol
```

## 7. So Sánh Cũ Và Mới

### Cũ: Public Inputs Mở Rộng Theo Bucket

```text
Groth16 public input:
  R.x/R.y/S/A.x/A.y/timestamp per signer
  height/round/block hashes/chainID

Size:
  22 * bucket + 116
```

Ưu điểm:

```text
Circuit nhẹ hơn vì không cần SHA-256 commitment.
```

Nhược điểm:

```text
Verifier Solidity rất lớn.
Mỗi bucket lớn hơn làm contract phình.
Dễ vượt EIP-170 24KB contract size limit.
```

### Mới: Hash-Aggregate

```text
Groth16 public input:
  Hash[32]
```

Ưu điểm:

```text
Verifier Solidity nhỏ hơn nhiều.
Public input count cố định, không tăng theo bucket.
Deploy được per-bucket verifier dễ hơn.
```

Tradeoff:

```text
Circuit/prover nặng hơn vì phải tính SHA-256 trong circuit.
Phải đảm bảo byte layout giống hệt giữa Go, circuit, Solidity.
```

## 8. Điều Cần Nhớ Khi Review

Không nên nói:

```text
Valset/header/sig là public input.
```

Nên nói:

```text
Valset/header/sig là semantic public data hoặc calldata-visible data.
Groth16 public input thật chỉ là hash commitment.
```

Hash không làm dữ liệu secret. Dữ liệu vẫn nằm trong calldata. Hash chỉ làm cho Groth16 verifier không phải nhận từng field riêng lẻ, nhưng proof vẫn bị ràng buộc với đúng calldata đó.

```text
Hash = public commitment kích thước cố định.
```

## 9. Checklist Khi Regenerate

Sau khi đổi circuit sang hash-aggregate, phải regenerate:

```text
relayer/prover artifacts:
  r1cs.bin
  pk.bin
  vk.bin

Solidity verifiers:
  contracts/verifiers/Groth16Verifier_N4.sol
  contracts/verifiers/Groth16Verifier_N8.sol
  ...
```

Verifier generated đúng phải có ABI:

```solidity
verifyProof(
    uint256[8] calldata proof,
    uint256[2] calldata commitments,
    uint256[2] calldata commitmentPok,
    uint256[32] calldata input
)
```

Nếu verifier vẫn là:

```solidity
uint256[204] calldata input
uint256[292] calldata input
```

thì đó là verifier cũ, chưa khớp hash-aggregate.
