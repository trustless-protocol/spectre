// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IVerifier} from "../interfaces/IVerifier.sol";
import {IICS07TendermintMsgs} from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import {Encode} from "./Encode.sol";

/// @title WrapperVerifier
/// @notice Rebuilds each validator's canonical-vote bytes on-chain, hashes the
///         witness with SHA-256, and forwards the 32-byte digest packed into
///         two 128-bit public field elements to the per-bucket gnark Groth16
///         verifier. One `BucketVerifier` must be registered per supported
///         bucket via `setBucket` before `verifyBatchProof` can succeed for
///         that size.
///
///         The in-circuit hash (see prover/hash_witness.go) reproduces the same
///         byte layout; any drift between this file, hash_witness.go, and
///         BatchCircuit.Define() fails the proof's public-input assertion.
contract WrapperVerifier is IVerifier {
    /// Per-slot canonical-vote buffer width baked into the circuit. Must match
    /// prover.MaxMsgLen exactly — anything larger fails LengthExceeded; shorter
    /// payloads are right-padded with zeros.
    // Must match relayer/prover/dummy.go::dummyMagic exactly. The circuit
    // asserts SHA-256 over this layout, so drift invalidates every proof.
    bytes14 constant DUMMY_MAGIC = 0x666173742d6962632d64756d6d79; // "fast-ibc-dummy"
    uint16 constant MAX_MSG_LEN = 192;
    uint256 constant WITNESS_SLOT_MSG_OFFSET = 1 + 32 + 2;
    uint256 constant WITNESS_SLOT_LEN = WITNESS_SLOT_MSG_OFFSET + MAX_MSG_LEN;

    /// @param verifier Address of the per-bucket gnark-generated Groth16 verifier.
    /// @param selector 4-byte function selector of that verifier's `verifyProof`.
    ///                 Upstream gnark PR #1554 (merged Feb 2026) switched the
    ///                 BN254 Solidity verifier signature from
    ///                 `verifyProof(uint256[8],uint256[2],uint256[2],uint256[N])`
    ///                 to `verifyProof(bytes,uint256[N])` where the bytes blob
    ///                 packs (proof || commitments || commitmentPok) = 384 bytes.
    ///                 `N` now equals 2 because the SHA-256 digest is exposed
    ///                 as two 128-bit public field elements.
    struct BucketVerifier {
        address verifier;
        bytes4 selector;
    }

    address public immutable OWNER;
    mapping(uint16 => BucketVerifier) public buckets;

    error UnknownBucket(uint16 bucket);
    error LengthMismatch();
    error MsgTooLong(uint256 length);
    error NotOwner();

    constructor(address owner) {
        OWNER = owner;
    }

    function setBucket(uint16 bucket, address verifier, bytes4 selector) external {
        if (msg.sender != OWNER) revert NotOwner();
        buckets[bucket] = BucketVerifier({verifier: verifier, selector: selector});
    }

    /// @inheritdoc IVerifier
    function verifyBatchProof(
        uint16 bucket,
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    ) external view override returns (bool) {
        if (
            pubkeys.length != bucket
                || timestampSeconds.length != bucket
                || timestampNanos.length != bucket
                || active.length != bucket
        ) revert LengthMismatch();

        BucketVerifier memory bv = buckets[bucket];
        if (bv.verifier == address(0)) revert UnknownBucket(bucket);

        bytes32 h = _hashWitness(bucket, pubkeys, timestampSeconds, timestampNanos, active, shared);

        uint256[2] memory publicInputs = _digestPublicInputs(h);

        bytes memory proofBytes = abi.encodePacked(proof, commitments, commitmentPok);
        bytes memory cd = abi.encodeWithSelector(bv.selector, proofBytes, publicInputs);
        (bool ok,) = bv.verifier.staticcall(cd);
        return ok;
    }

    /// @dev Canonical witness byte layout — must match Go's
    ///      prover.ComputeWitnessHash and BatchCircuit.Define exactly:
    ///        per slot: active(1) || A(32) || msgLen(2 BE) || msg[MAX_MSG_LEN padded]
    ///      A is hashed because Solidity uses pubkeys[i] to look up validator
    ///      voting power; binding it stops calldata pubkey swaps. R/S are not
    ///      in calldata or hash — the Groth16 proof itself binds them via the
    ///      in-circuit Ed25519 verify, and no on-chain logic consumes them.
    ///      For active slots msg = canonical-vote bytes rebuilt from
    ///      (shared, timestampSeconds[i], timestampNanos[i]). For inactive
    ///      (padding) slots msg = `dummyMagic || bucket(2 BE) || slot(2 BE)`,
    ///      reproduced byte-for-byte from prover/dummy.go.
    function _hashWitness(
        uint16 bucket,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    ) internal pure returns (bytes32) {
        bytes memory encodedBlockId = _sharedBlockIdBytes(shared);
        bytes memory commonVotePrefix = _buildVotePrefix(shared, encodedBlockId);
        bytes memory chainSuffix = _buildChainSuffix(shared.chainID);

        bytes memory buf = new bytes(pubkeys.length * WITNESS_SLOT_LEN);
        uint256 offset = 0;
        for (uint256 i = 0; i < pubkeys.length; i++) {
            uint256 msgLen;
            if (active[i]) {
                msgLen = _writeVoteSignBytes(
                    buf,
                    offset + WITNESS_SLOT_MSG_OFFSET,
                    commonVotePrefix,
                    chainSuffix,
                    timestampSeconds[i],
                    timestampNanos[i]
                );
            } else {
                msgLen = _writeDummyMsgBytes(buf, offset + WITNESS_SLOT_MSG_OFFSET, bucket, uint16(i));
            }

            _storeByte(buf, offset, active[i] ? 1 : 0);
            _storeBytes32(buf, offset + 1, pubkeys[i]);
            _storeByte(buf, offset + 33, msgLen >> 8);
            _storeByte(buf, offset + 34, msgLen);
            offset += WITNESS_SLOT_LEN;
        }
        return sha256(buf);
    }

    function _digestPublicInputs(bytes32 h) internal pure returns (uint256[2] memory publicInputs) {
        uint256 digest = uint256(h);
        publicInputs[0] = digest >> 128;
        publicInputs[1] = digest & type(uint128).max;
    }

    function _storeByte(bytes memory dst, uint256 offset, uint256 value) private pure {
        assembly {
            mstore8(add(add(dst, 0x20), offset), value)
        }
    }

    function _storeBytes32(bytes memory dst, uint256 offset, bytes32 value) private pure {
        assembly {
            mstore(add(add(dst, 0x20), offset), value)
        }
    }

    function _copyBytes(bytes memory dst, uint256 dstOffset, bytes memory src) private pure returns (uint256) {
        uint256 len;
        assembly {
            len := mload(src)
            mcopy(add(add(dst, 0x20), dstOffset), add(src, 0x20), len)
        }
        return dstOffset + len;
    }

    function _sharedBlockIdBytes(
        IVerifier.SharedBlock calldata shared
    ) internal pure returns (bytes memory) {
        IICS07TendermintMsgs.BlockId memory blockId = IICS07TendermintMsgs.BlockId({
            hashData: shared.blockIDHash,
            partSetHeader: IICS07TendermintMsgs.PartSetHeader({
                total: shared.partSetTotal,
                hashData: shared.partSetHash
            })
        });
        return Encode.encodeBlockId(blockId);
    }

    /// @dev Build the canonical-vote prefix shared by every active slot:
    ///      type || height || round || block_id. Timestamp and chain_id are
    ///      appended separately so the per-slot path only writes the fields
    ///      that actually vary.
    function _buildVotePrefix(
        IVerifier.SharedBlock calldata shared,
        bytes memory encodedBlockId
    ) internal pure returns (bytes memory out) {
        uint256 outLen = 2; // field 1: tag + PRECOMMIT value
        if (shared.height > 0) {
            outLen += 9;
        }
        if (shared.round > 0) {
            outLen += 9;
        }
        if (encodedBlockId.length > 0) {
            outLen += 1 + _varintLen(encodedBlockId.length) + encodedBlockId.length;
        }

        out = new bytes(outLen);
        uint256 offset = 0;
        _storeByte(out, offset, 0x08);
        _storeByte(out, offset + 1, 0x02);
        offset += 2;

        if (shared.height > 0) {
            _storeByte(out, offset, 0x11);
            offset = _copyBytes(out, offset + 1, Encode.encodeSfixed64(int64(uint64(shared.height))));
        }

        if (shared.round > 0) {
            _storeByte(out, offset, 0x19);
            offset = _copyBytes(out, offset + 1, Encode.encodeSfixed64(int64(uint64(shared.round))));
        }

        if (encodedBlockId.length > 0) {
            _storeByte(out, offset, 0x22);
            offset = _writeVarint(out, offset + 1, encodedBlockId.length);
            _copyBytes(out, offset, encodedBlockId);
        }
    }

    function _buildChainSuffix(bytes calldata chainId) internal pure returns (bytes memory out) {
        if (chainId.length == 0) {
            return new bytes(0);
        }

        out = new bytes(1 + _varintLen(chainId.length) + chainId.length);
        uint256 offset = 0;
        _storeByte(out, offset, 0x32);
        offset = _writeVarint(out, offset + 1, chainId.length);
        _copyCalldataBytes(out, offset, chainId);
    }

    /// @dev Writes the cometbft canonical-vote bytes directly into the final
    ///      witness buffer, avoiding per-slot temporary byte-array allocations.
    function _writeVoteSignBytes(
        bytes memory dst,
        uint256 dstOffset,
        bytes memory commonVotePrefix,
        bytes memory chainSuffix,
        uint64 tsSec,
        uint32 tsNanos
    ) internal pure returns (uint256 msgLen) {
        // Compose Timestamp{seconds, nanos} — gogoproto omits zero scalars.
        uint256 encodedTsLen = 0;
        if (tsSec > 0) {
            encodedTsLen += 1 + _varintLen(uint256(tsSec));
        }
        if (tsNanos > 0) {
            encodedTsLen += 1 + _varintLen(uint256(tsNanos));
        }

        uint256 encodedLen = commonVotePrefix.length + chainSuffix.length;
        if (encodedTsLen > 0) {
            encodedLen += 1 + _varintLen(encodedTsLen) + encodedTsLen;
        }

        msgLen = _varintLen(encodedLen) + encodedLen;
        if (msgLen > MAX_MSG_LEN) revert MsgTooLong(msgLen);

        uint256 offset = _writeVarint(dst, dstOffset, encodedLen);
        offset = _copyBytes(dst, offset, commonVotePrefix);

        if (encodedTsLen > 0) {
            _storeByte(dst, offset, 0x2A);
            offset = _writeVarint(dst, offset + 1, encodedTsLen);
            if (tsSec > 0) {
                _storeByte(dst, offset, 0x08);
                offset = _writeVarint(dst, offset + 1, uint256(tsSec));
            }
            if (tsNanos > 0) {
                _storeByte(dst, offset, 0x10);
                offset = _writeVarint(dst, offset + 1, uint256(tsNanos));
            }
        }

        _copyBytes(dst, offset, chainSuffix);
    }

    /// @dev Reproduces prover.DummyMsgBytes for a padding slot. Must stay in
    ///      lock-step with relayer/prover/dummy.go — any drift breaks the
    ///      circuit's hash assertion.
    function _writeDummyMsgBytes(
        bytes memory dst,
        uint256 dstOffset,
        uint16 bucket,
        uint16 slot
    ) internal pure returns (uint256 msgLen) {
        msgLen = 18; // len("fast-ibc-dummy") + uint16(bucket) + uint16(slot)
        if (msgLen > MAX_MSG_LEN) revert MsgTooLong(msgLen);

        _storeBytes32(dst, dstOffset, bytes32(DUMMY_MAGIC));
        uint256 suffixOffset = dstOffset + 14;
        _storeByte(dst, suffixOffset, bucket >> 8);
        _storeByte(dst, suffixOffset + 1, bucket);
        _storeByte(dst, suffixOffset + 2, slot >> 8);
        _storeByte(dst, suffixOffset + 3, slot);
    }

    function _copyCalldataBytes(bytes memory dst, uint256 dstOffset, bytes calldata src) private pure returns (uint256) {
        assembly {
            calldatacopy(add(add(dst, 0x20), dstOffset), src.offset, src.length)
        }
        return dstOffset + src.length;
    }

    function _varintLen(uint256 value) private pure returns (uint256 len) {
        len = 1;
        while (value >= 128) {
            value >>= 7;
            unchecked {
                ++len;
            }
        }
    }

    function _writeVarint(bytes memory out, uint256 offset, uint256 value) private pure returns (uint256) {
        while (value >= 128) {
            _storeByte(out, offset, uint8((value & 0x7F) | 0x80));
            unchecked {
                ++offset;
            }
            value >>= 7;
        }
        _storeByte(out, offset, uint8(value));
        return offset + 1;
    }
}
