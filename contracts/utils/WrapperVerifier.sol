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
    uint16 constant MAX_MSG_LEN = 192;
    uint256 constant WITNESS_SLOT_LEN = 1 + 32 + 2 + MAX_MSG_LEN;

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
        bytes memory buf = new bytes(pubkeys.length * WITNESS_SLOT_LEN);
        uint256 offset = 0;
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes memory msgBytes;
            if (active[i]) {
                msgBytes = _voteSignBytes(shared, timestampSeconds[i], timestampNanos[i]);
            } else {
                msgBytes = _dummyMsgBytes(bucket, uint16(i));
            }
            uint256 msgLen = msgBytes.length;
            if (msgLen > MAX_MSG_LEN) revert MsgTooLong(msgLen);

            _storeByte(buf, offset, active[i] ? 1 : 0);
            _storeBytes32(buf, offset + 1, pubkeys[i]);
            _storeByte(buf, offset + 33, msgLen >> 8);
            _storeByte(buf, offset + 34, msgLen);
            _copyBytes(buf, offset + 35, msgBytes);
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

    function _copyBytes(bytes memory dst, uint256 dstOffset, bytes memory src) private pure {
        assembly {
            let len := mload(src)
            let dstPtr := add(add(dst, 0x20), dstOffset)
            let srcPtr := add(src, 0x20)
            let fullWords := and(len, not(31))

            for { let i := 0 } lt(i, fullWords) { i := add(i, 0x20) } {
                mstore(add(dstPtr, i), mload(add(srcPtr, i)))
            }

            let rem := and(len, 31)
            if rem {
                let srcWord := mload(add(srcPtr, fullWords))
                let dstWord := mload(add(dstPtr, fullWords))
                let keepMask := sub(shl(mul(sub(32, rem), 8), 1), 1)
                mstore(
                    add(dstPtr, fullWords),
                    or(and(srcWord, not(keepMask)), and(dstWord, keepMask))
                )
            }
        }
    }

    /// @dev Reproduces prover.DummyMsgBytes for a padding slot. Must stay in
    ///      lock-step with relayer/prover/dummy.go — any drift breaks the
    ///      circuit's hash assertion.
    function _dummyMsgBytes(uint16 bucket, uint16 slot) internal pure returns (bytes memory) {
        return abi.encodePacked("fast-ibc-dummy", bucket, slot);
    }

    /// @dev Build the cometbft canonical-vote bytes for a single validator
    ///      from (sharedBlock, ts) — equivalent to Encode.voteSignBytes but
    ///      taking the SharedBlock + explicit timestamp instead of a
    ///      BlockCommit. Output is byte-identical to cometbft's
    ///      `Commit.VoteSignBytes(chainID, valIdx)` when fed equivalent inputs.
    function _voteSignBytes(
        IVerifier.SharedBlock calldata shared,
        uint64 tsSec,
        uint32 tsNanos
    ) internal pure returns (bytes memory) {
        IICS07TendermintMsgs.BlockId memory blockId = IICS07TendermintMsgs.BlockId({
            hashData: shared.blockIDHash,
            partSetHeader: IICS07TendermintMsgs.PartSetHeader({
                total: shared.partSetTotal,
                hashData: shared.partSetHash
            })
        });
        bytes memory encodedBlockId = Encode.encodeBlockId(blockId);

        // Compose Timestamp{seconds, nanos} — gogoproto omits zero scalars.
        bytes memory encodedTs;
        if (tsSec > 0) {
            encodedTs = abi.encodePacked(uint8(0x08), Encode.encodeVarint(uint256(tsSec)));
        }
        if (tsNanos > 0) {
            encodedTs = abi.encodePacked(encodedTs, uint8(0x10), Encode.encodeVarint(uint256(tsNanos)));
        }

        bytes memory encoded;
        // Field 1: type = PRECOMMIT (2), tag 0x08
        encoded = abi.encodePacked(uint8(0x08), uint8(0x02));

        // Field 2: height, sfixed64, tag 0x11
        if (shared.height > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x11), Encode.encodeSfixed64(int64(uint64(shared.height))));
        }

        // Field 3: round, sfixed64, tag 0x19
        if (shared.round > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x19), Encode.encodeSfixed64(int64(uint64(shared.round))));
        }

        // Field 4: block_id, length-delimited, tag 0x22
        if (encodedBlockId.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x22), Encode.encodeVarint(encodedBlockId.length), encodedBlockId);
        }

        // Field 5: timestamp, length-delimited, tag 0x2a
        if (encodedTs.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x2A), Encode.encodeVarint(encodedTs.length), encodedTs);
        }

        // Field 6: chain_id, length-delimited, tag 0x32
        if (shared.chainID.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x32), Encode.encodeVarint(shared.chainID.length), shared.chainID);
        }

        // Wrap with overall varint length prefix (cometbft MarshalDelimited).
        return abi.encodePacked(Encode.encodeVarint(encoded.length), encoded);
    }
}
