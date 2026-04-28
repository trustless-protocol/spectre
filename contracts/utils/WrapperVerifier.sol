// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IVerifier} from "../interfaces/IVerifier.sol";
import {IICS07TendermintMsgs} from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import {Encode} from "./Encode.sol";

/// @title WrapperVerifier
/// @notice Rebuilds each validator's canonical-vote bytes on-chain, hashes the
///         witness with SHA-256, and forwards the 32-byte digest to the
///         per-bucket gnark Groth16 verifier. One `BucketVerifier` must be
///         registered per supported bucket via `setBucket` before
///         `verifyBatchProof` can succeed for that size.
///
///         The in-circuit hash (see prover/hash_witness.go) reproduces the same
///         byte layout; any drift between this file, hash_witness.go, and
///         BatchCircuit.Define() fails the proof's public-input assertion.
contract WrapperVerifier is IVerifier {
    /// Per-slot canonical-vote buffer width baked into the circuit. Must match
    /// prover.MaxMsgLen exactly — anything larger fails LengthExceeded; shorter
    /// payloads are right-padded with zeros.
    uint16 constant MAX_MSG_LEN = 192;

    /// @param verifier Address of the per-bucket gnark-generated Groth16 verifier.
    /// @param selector 4-byte function selector of that verifier's `verifyProof`
    ///                 (always `uint256[8],uint256[2],uint256[2],uint256[32]`
    ///                 — every bucket has L=32 under hash-aggregate).
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
        bytes32[2][] calldata signatures,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    ) external view override returns (bool) {
        if (
            signatures.length != bucket
                || pubkeys.length != bucket
                || timestampSeconds.length != bucket
                || timestampNanos.length != bucket
                || active.length != bucket
        ) revert LengthMismatch();

        BucketVerifier memory bv = buckets[bucket];
        if (bv.verifier == address(0)) revert UnknownBucket(bucket);

        bytes32 h = _hashWitness(bucket, signatures, pubkeys, timestampSeconds, timestampNanos, active, shared);

        uint256[32] memory publicInputs;
        for (uint256 i = 0; i < 32; i++) {
            publicInputs[i] = uint256(uint8(h[i]));
        }

        bytes memory cd =
            abi.encodePacked(bv.selector, proof, commitments, commitmentPok, publicInputs);
        (bool ok,) = bv.verifier.staticcall(cd);
        return ok;
    }

    /// @dev Canonical witness byte layout — must match Go's
    ///      prover.ComputeWitnessHash and BatchCircuit.Define exactly:
    ///        per slot: active(1) || R(32) || S(32) || A(32) || msgLen(2 BE) || msg[MAX_MSG_LEN padded]
    ///      For active slots msg = canonical-vote bytes rebuilt from
    ///      (shared, timestampSeconds[i], timestampNanos[i]). For inactive
    ///      (padding) slots msg = `dummyMagic || bucket(2 BE) || slot(2 BE)`,
    ///      reproduced byte-for-byte from prover/dummy.go.
    function _hashWitness(
        uint16 bucket,
        bytes32[2][] calldata signatures,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    ) internal pure returns (bytes32) {
        bytes memory buf;
        for (uint256 i = 0; i < signatures.length; i++) {
            bytes memory msgBytes;
            if (active[i]) {
                msgBytes = _voteSignBytes(shared, timestampSeconds[i], timestampNanos[i]);
            } else {
                msgBytes = _dummyMsgBytes(bucket, uint16(i));
            }
            if (msgBytes.length > MAX_MSG_LEN) revert MsgTooLong(msgBytes.length);
            bytes memory padded = new bytes(MAX_MSG_LEN);
            for (uint256 j = 0; j < msgBytes.length; j++) {
                padded[j] = msgBytes[j];
            }
            buf = abi.encodePacked(
                buf,
                active[i] ? bytes1(0x01) : bytes1(0x00), // active (1)
                signatures[i][0],          // R (32)
                signatures[i][1],          // S (32)
                pubkeys[i],                // A (32)
                uint16(msgBytes.length),   // msgLen (2 BE)
                padded                     // msg (MAX_MSG_LEN padded)
            );
        }
        return sha256(buf);
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
