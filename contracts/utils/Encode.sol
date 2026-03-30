pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";
library Encode {
    function encodeVarint(uint256 value) public pure returns (bytes memory) {
        if (value < 128) {
            return abi.encodePacked(uint8(value));
        }
        
        bytes memory result;
        while (value >= 128) {
            result = abi.encodePacked(result, uint8((value & 0x7F) | 0x80));
            value >>= 7;
        }
        result = abi.encodePacked(result, uint8(value));
        return result;
    }

    function encodeString(string memory value) public pure returns (bytes memory) {
        bytes memory valueBytes = bytes(value);
        uint256 length = valueBytes.length;
        
        // Encode length as varint
        bytes memory lengthBytes = encodeVarint(length);

        // Concatenate length and value
        return abi.encodePacked(lengthBytes, valueBytes);
    }

    /// @notice Encodes a nanosecond timestamp as a protobuf google.protobuf.Timestamp.
    /// @dev Omits fields with zero value per proto3 rules.
    function encodeTimestamp(uint128 nanos) public pure returns (bytes memory) {
        uint128 secs = nanos / 1_000_000_000;
        uint128 ns = nanos % 1_000_000_000;
        bytes memory encoded = new bytes(0);

        // Field 1: seconds (tag = 1, wire type = 0 for varint)
        if (secs > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x08)); // tag: (1 << 3) | 0
            encoded = abi.encodePacked(encoded, encodeVarint(uint256(secs)));
        }
        // Field 2: nanos (tag = 2, wire type = 0 for varint)
        if (ns > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x10)); // tag: (2 << 3) | 0
            encoded = abi.encodePacked(encoded, encodeVarint(uint256(ns)));
        }

        return encoded;
    }

    function encodeValidator(
        IICS07TendermintMsgs.SimpleValidator memory validator
    ) public pure returns (bytes memory) {
        bytes memory encoded = new bytes(0);

        // Field 1: pub_key (tag = 1, wire type = 2 for length-delimited)
        // PubKey is a nested message: crypto.PublicKey{Sum: &PublicKey_Ed25519{Ed25519: pubKey}}
        // Inner encoding: Ed25519 oneof field 1 (tag=0x0A), length=32, data
        bytes memory pubKeyInner = abi.encodePacked(uint8(0x0A), uint8(32), validator.pubKey);
        encoded = abi.encodePacked(encoded, uint8(0x0A)); // tag: (1 << 3) | 2
        encoded = abi.encodePacked(encoded, encodeVarint(pubKeyInner.length)); // length of nested PublicKey message
        encoded = abi.encodePacked(encoded, pubKeyInner); // nested PublicKey message

        // Field 2: voting_power (tag = 2, wire type = 0 for varint)
        // Proto3: skip zero-value fields
        if (validator.votingPower > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x10)); // tag: (2 << 3) | 0
            encoded = abi.encodePacked(encoded, encodeVarint(uint256(validator.votingPower)));
        }

        return encoded;
    }

    function encodeVersion(IICS07TendermintMsgs.Version memory version) public pure returns (bytes memory) {
        bytes memory encoded = new bytes(0);

        // Proto3: skip zero-value fields
        if (version.blockVersion > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x08)); // tag: (1 << 3) | 0
            encoded = abi.encodePacked(encoded, encodeVarint(uint256(version.blockVersion)));
        }

        if (version.appVersion > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x10)); // tag: (2 << 3) | 0
            encoded = abi.encodePacked(encoded, encodeVarint(uint256(version.appVersion)));
        }

        return encoded;
    }

    function encodeBlockId(IICS07TendermintMsgs.BlockId memory blockId) public pure returns (bytes memory) {
        bytes memory encoded = new bytes(0);
        
        // Field 1: hashData (tag = 1, wire type = 2 for bytes)
        encoded = abi.encodePacked(encoded, uint8(0x0A)); // tag: (1 << 3) | 2
        encoded = abi.encodePacked(encoded, uint8(32)); // 32 bytes length
        encoded = abi.encodePacked(encoded, blockId.hashData);
        
        // Field 2: partSetHeader (tag = 2, wire type = 2 for message)
        bytes memory partSetHeaderEncoded = encodePartSetHeader(blockId.partSetHeader);
        encoded = abi.encodePacked(encoded, uint8(0x12)); // tag: (2 << 3) | 2
        encoded = abi.encodePacked(encoded, encodeVarint(partSetHeaderEncoded.length));
        encoded = abi.encodePacked(encoded, partSetHeaderEncoded);
        
        return encoded;
    }

    /// @notice Wraps a string in gogoproto StringValue{Value: str} for header hashing.
    /// Matches CometBFT's cdcEncode(string) used in Header.Hash().
    function cdcEncodeString(string memory value) public pure returns (bytes memory) {
        bytes memory valueBytes = bytes(value);
        if (valueBytes.length == 0) return new bytes(0);
        // StringValue field 1 (tag=0x0A, wire type 2) + varint(len) + bytes
        return abi.encodePacked(uint8(0x0A), encodeVarint(valueBytes.length), valueBytes);
    }

    /// @notice Wraps an int64 in gogoproto Int64Value{Value: n} for header hashing.
    /// Matches CometBFT's cdcEncode(int64) used in Header.Hash().
    function cdcEncodeInt64(uint256 value) public pure returns (bytes memory) {
        if (value == 0) return new bytes(0);
        // Int64Value field 1 (tag=0x08, wire type 0) + varint(value)
        return abi.encodePacked(uint8(0x08), encodeVarint(value));
    }

    /// @notice Wraps variable-length bytes in gogoproto BytesValue{Value: bz} for header hashing.
    /// Matches CometBFT's cdcEncode([]byte) used in Header.Hash().
    function cdcEncodeBytes(bytes memory value) public pure returns (bytes memory) {
        if (value.length == 0) return new bytes(0);
        // BytesValue field 1 (tag=0x0A, wire type 2) + varint(len) + bytes
        return abi.encodePacked(uint8(0x0A), encodeVarint(value.length), value);
    }

    /// @notice Wraps a bytes32 hash in gogoproto BytesValue{Value: hash} for header hashing.
    /// Matches CometBFT's cdcEncode(HexBytes) used in Header.Hash().
    function cdcEncodeBytes32(bytes32 value) public pure returns (bytes memory) {
        if (value == bytes32(0)) return new bytes(0);
        // BytesValue field 1 (tag=0x0A, wire type 2) + length 32 + hash
        return abi.encodePacked(uint8(0x0A), uint8(32), value);
    }

    /// @notice Encodes a signed 64-bit integer as 8-byte little-endian (protobuf sfixed64).
    function encodeSfixed64(int64 value) public pure returns (bytes memory) {
        bytes memory result = new bytes(8);
        uint64 v = uint64(value);
        result[0] = bytes1(uint8(v));
        result[1] = bytes1(uint8(v >> 8));
        result[2] = bytes1(uint8(v >> 16));
        result[3] = bytes1(uint8(v >> 24));
        result[4] = bytes1(uint8(v >> 32));
        result[5] = bytes1(uint8(v >> 40));
        result[6] = bytes1(uint8(v >> 48));
        result[7] = bytes1(uint8(v >> 56));
        return result;
    }

    function encodePartSetHeader(IICS07TendermintMsgs.PartSetHeader memory partSetHeader) public pure returns (bytes memory) {
        bytes memory encoded = new bytes(0);
        
        // Field 1: total (tag = 1, wire type = 0 for varint)
        encoded = abi.encodePacked(encoded, uint8(0x08)); // tag: (1 << 3) | 0
        encoded = abi.encodePacked(encoded, encodeVarint(uint256(partSetHeader.total)));
        
        // Field 2: hashData (tag = 2, wire type = 2 for bytes)
        encoded = abi.encodePacked(encoded, uint8(0x12)); // tag: (2 << 3) | 2
        encoded = abi.encodePacked(encoded, uint8(32)); // 32 bytes length
        encoded = abi.encodePacked(encoded, partSetHeader.hashData);
        
        return encoded;
    }

    /// @notice Encodes a CanonicalVote as protobuf bytes (equivalent to CometBFT's VoteSignBytes).
    /// @param commit The block commit.
    /// @param chainId The chain ID string.
    /// @param valIdx The validator index into commitSigs.
    /// @return The protobuf-encoded canonical vote bytes.
    function voteSignBytes(
        IICS07TendermintMsgs.BlockCommit memory commit,
        string memory chainId,
        uint32 valIdx
    ) public pure returns (bytes memory) {
        IICS07TendermintMsgs.CommitSig memory commitSig = commit.commitSigs[valIdx];

        bool useCommitBlockId = commitSig.flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT;
        bytes memory encodedBlockId = useCommitBlockId ? encodeBlockId(commit.blockId) : new bytes(0);
        bytes memory encodedTimestamp = encodeTimestamp(commitSig.data.timestamp);
        bytes memory chainIdBytes = bytes(chainId);

        bytes memory encoded = new bytes(0);

        // Field 1: type = PrecommitType (2), varint, tag 0x08
        encoded = abi.encodePacked(encoded, uint8(0x08), uint8(0x02));

        // Field 2: height, sfixed64 (fixed 8-byte little-endian), tag 0x11
        if (commit.height > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x11));
            encoded = abi.encodePacked(encoded, encodeSfixed64(int64(uint64(commit.height))));
        }

        // Field 3: round, sfixed64 (fixed 8-byte little-endian), tag 0x19
        if (commit.round > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x19));
            encoded = abi.encodePacked(encoded, encodeSfixed64(int64(uint64(commit.round))));
        }

        // Field 4: block_id, length-delimited, tag 0x22
        if (encodedBlockId.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x22));
            encoded = abi.encodePacked(encoded, encodeVarint(encodedBlockId.length));
            encoded = abi.encodePacked(encoded, encodedBlockId);
        }

        // Field 5: timestamp, length-delimited, tag 0x2a
        if (encodedTimestamp.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x2A));
            encoded = abi.encodePacked(encoded, encodeVarint(encodedTimestamp.length));
            encoded = abi.encodePacked(encoded, encodedTimestamp);
        }

        // Field 6: chain_id, length-delimited string, tag 0x32
        if (chainIdBytes.length > 0) {
            encoded = abi.encodePacked(encoded, uint8(0x32));
            encoded = abi.encodePacked(encoded, encodeString(chainId));
        }

        // Wrap with varint length prefix (MarshalDelimited for Amino compatibility)
        return abi.encodePacked(encodeVarint(encoded.length), encoded);
    }
}
