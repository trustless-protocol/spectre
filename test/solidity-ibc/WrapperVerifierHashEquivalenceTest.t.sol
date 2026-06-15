// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { Encode } from "../../contracts/utils/Encode.sol";
import { IVerifier } from "../../contracts/interfaces/IVerifier.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";

contract WrapperVerifierHashHarness is WrapperVerifier {
    constructor() WrapperVerifier(address(this)) { }

    function hashWitness(
        uint16 bucket,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        IVerifier.SharedBlock calldata shared
    )
        external
        pure
        returns (bytes32)
    {
        return _hashWitness(bucket, pubkeys, timestampSeconds, timestampNanos, active, shared);
    }
}

contract WrapperVerifierHashEquivalenceTest is Test {
    uint16 internal constant MAX_MSG_LEN = 192;
    uint256 internal constant WITNESS_SLOT_MSG_OFFSET = 1 + 32 + 2;
    uint256 internal constant WITNESS_SLOT_LEN = WITNESS_SLOT_MSG_OFFSET + MAX_MSG_LEN;

    WrapperVerifierHashHarness internal harness;

    function setUp() public {
        harness = new WrapperVerifierHashHarness();
    }

    function test_hashWitness_matchesReference_bucket16_mixedSlots() public {
        uint16 bucket = 16;
        bytes32[] memory pubkeys = new bytes32[](bucket);
        uint64[] memory tsS = new uint64[](bucket);
        uint32[] memory tsN = new uint32[](bucket);
        bool[] memory active = new bool[](bucket);

        for (uint256 i = 0; i < bucket; i++) {
            pubkeys[i] = bytes32(uint256(0xA0 + i));
            tsS[i] = uint64(1_700_000_000 + i);
            tsN[i] = uint32((i % 5) * 1_000_000);
            active[i] = i < 11;
        }

        IVerifier.SharedBlock memory shared = IVerifier.SharedBlock({
            height: 53,
            round: 0,
            blockIDHash: bytes32(uint256(0x1234)),
            partSetTotal: 1,
            partSetHash: bytes32(uint256(0x5678)),
            chainID: bytes("test-ibc-eth")
        });

        assertEq(
            harness.hashWitness(bucket, pubkeys, tsS, tsN, active, shared),
            _hashWitnessReference(bucket, pubkeys, tsS, tsN, active, shared)
        );
    }

    function test_hashWitness_matchesReference_edgeVarintsAndLongChainId() public {
        uint16 bucket = 4;
        bytes32[] memory pubkeys = new bytes32[](bucket);
        uint64[] memory tsS = new uint64[](bucket);
        uint32[] memory tsN = new uint32[](bucket);
        bool[] memory active = new bool[](bucket);

        pubkeys[0] = bytes32(uint256(1));
        pubkeys[1] = bytes32(uint256(2));
        pubkeys[2] = bytes32(uint256(3));
        pubkeys[3] = bytes32(uint256(4));

        tsS[0] = 0;
        tsS[1] = 127;
        tsS[2] = 128;
        tsS[3] = 16_384;

        tsN[0] = 0;
        tsN[1] = 127;
        tsN[2] = 128;
        tsN[3] = 1_000_000;

        active[0] = true;
        active[1] = false;
        active[2] = true;
        active[3] = false;

        bytes memory chainId = new bytes(80);
        for (uint256 i = 0; i < chainId.length; i++) {
            chainId[i] = bytes1(uint8(65 + (i % 26)));
        }

        IVerifier.SharedBlock memory shared = IVerifier.SharedBlock({
            height: 128,
            round: 127,
            blockIDHash: bytes32(uint256(0xCAFE)),
            partSetTotal: 128,
            partSetHash: bytes32(uint256(0xBEEF)),
            chainID: chainId
        });

        assertEq(
            harness.hashWitness(bucket, pubkeys, tsS, tsN, active, shared),
            _hashWitnessReference(bucket, pubkeys, tsS, tsN, active, shared)
        );
    }

    function _hashWitnessReference(
        uint16 bucket,
        bytes32[] memory pubkeys,
        uint64[] memory timestampSeconds,
        uint32[] memory timestampNanos,
        bool[] memory active,
        IVerifier.SharedBlock memory shared
    )
        internal
        pure
        returns (bytes32)
    {
        bytes memory buf = new bytes(pubkeys.length * WITNESS_SLOT_LEN);
        uint256 offset = 0;
        for (uint256 i = 0; i < pubkeys.length; i++) {
            bytes memory msgBytes = active[i]
                ? _voteSignBytesReference(shared, timestampSeconds[i], timestampNanos[i])
                : abi.encodePacked("fast-ibc-dummy", bucket, uint16(i));
            require(msgBytes.length <= MAX_MSG_LEN, "reference msg too long");

            _storeByte(buf, offset, active[i] ? 1 : 0);
            _storeBytes32(buf, offset + 1, pubkeys[i]);
            _storeByte(buf, offset + 33, msgBytes.length >> 8);
            _storeByte(buf, offset + 34, msgBytes.length);
            _copyBytes(buf, offset + WITNESS_SLOT_MSG_OFFSET, msgBytes);
            offset += WITNESS_SLOT_LEN;
        }
        return sha256(buf);
    }

    function _voteSignBytesReference(
        IVerifier.SharedBlock memory shared,
        uint64 tsSec,
        uint32 tsNanos
    )
        internal
        pure
        returns (bytes memory)
    {
        IICS07TendermintMsgs.BlockId memory blockId = IICS07TendermintMsgs.BlockId({
            hashData: shared.blockIDHash,
            partSetHeader: IICS07TendermintMsgs.PartSetHeader({
                total: shared.partSetTotal, hashData: shared.partSetHash
            })
        });
        bytes memory encodedBlockId = Encode.encodeBlockId(blockId);

        uint256 encodedTsLen = 0;
        if (tsSec > 0) {
            encodedTsLen += 1 + _varintLen(uint256(tsSec));
        }
        if (tsNanos > 0) {
            encodedTsLen += 1 + _varintLen(uint256(tsNanos));
        }

        bytes memory encodedTs = new bytes(encodedTsLen);
        uint256 tsOffset = 0;
        if (tsSec > 0) {
            _storeByte(encodedTs, tsOffset, 0x08);
            tsOffset = _writeVarint(encodedTs, tsOffset + 1, uint256(tsSec));
        }
        if (tsNanos > 0) {
            _storeByte(encodedTs, tsOffset, 0x10);
            _writeVarint(encodedTs, tsOffset + 1, uint256(tsNanos));
        }

        uint256 encodedLen = 2;
        if (shared.height > 0) {
            encodedLen += 9;
        }
        if (shared.round > 0) {
            encodedLen += 9;
        }
        if (encodedBlockId.length > 0) {
            encodedLen += 1 + _varintLen(encodedBlockId.length) + encodedBlockId.length;
        }
        // Timestamp is always present in CanonicalVote (matches Go/CometBFT).
        encodedLen += 1 + _varintLen(encodedTs.length) + encodedTs.length;
        if (shared.chainID.length > 0) {
            encodedLen += 1 + _varintLen(shared.chainID.length) + shared.chainID.length;
        }

        bytes memory out = new bytes(_varintLen(encodedLen) + encodedLen);
        uint256 offset = _writeVarint(out, 0, encodedLen);
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
            offset = _copyBytes(out, offset, encodedBlockId);
        }
        _storeByte(out, offset, 0x2A);
        offset = _writeVarint(out, offset + 1, encodedTs.length);
        offset = _copyBytes(out, offset, encodedTs);
        if (shared.chainID.length > 0) {
            _storeByte(out, offset, 0x32);
            offset = _writeVarint(out, offset + 1, shared.chainID.length);
            _copyBytes(out, offset, shared.chainID);
        }

        return out;
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
                mstore(add(dstPtr, fullWords), or(and(srcWord, not(keepMask)), and(dstWord, keepMask)))
            }
        }
        return dstOffset + len;
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
