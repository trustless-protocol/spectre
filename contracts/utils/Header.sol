pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { Encode } from "./Encode.sol";
library Header {

    function hashValSet(
        IICS07TendermintMsgs.ValidatorSet memory valset
    ) pure public returns (bytes32) {
        uint256 validatorCount = valset.validators.length;
        if (validatorCount == 0) {
            return bytes32(0);
        }

        bytes32[] memory leafHashes = new bytes32[](validatorCount);
        for (uint256 i = 0; i < valset.validators.length; i++) {
            bytes memory validatorBytes = Encode.encodeValidator(
                IICS07TendermintMsgs.SimpleValidator({
                    pubKey: valset.validators[i].pubKey,
                    votingPower: valset.validators[i].votingPower
                })
            );
            leafHashes[i] = _leafHash(validatorBytes);
        }

        return _merkleHashRange(leafHashes, 0, validatorCount);
    }

    function hashHeader(
        IICS07TendermintMsgs.BlockHeader memory header
    ) pure public returns (bytes32) {
        // CometBFT Header.Hash() ALWAYS hashes 14 fields via merkle.
        // Missing fields are encoded as empty bytes (nil in Go).
        bytes32[] memory leafHashes = new bytes32[](14);

        leafHashes[0] = _leafHash(Encode.encodeVersion(header.version));
        leafHashes[1] = _leafHash(Encode.cdcEncodeString(header.chainId));
        leafHashes[2] = _leafHash(Encode.cdcEncodeInt64(uint256(header.height)));
        leafHashes[3] = _leafHash(Encode.encodeTimestamp(header.time));
        leafHashes[4] = header.hasLastBlockId ? _leafHash(Encode.encodeBlockId(header.lastBlockId)) : _emptyLeafHash();
        leafHashes[5] = header.hasLastCommitHash ? _leafHash(Encode.cdcEncodeBytes32(header.lastCommitHash)) : _emptyLeafHash();
        leafHashes[6] = header.hasDataHash ? _leafHash(Encode.cdcEncodeBytes32(header.dataHash)) : _emptyLeafHash();
        leafHashes[7] = _leafHash(Encode.cdcEncodeBytes32(header.validatorsHash));
        leafHashes[8] = _leafHash(Encode.cdcEncodeBytes32(header.nextValidatorsHash));
        leafHashes[9] = _leafHash(Encode.cdcEncodeBytes32(header.consensusHash));
        leafHashes[10] = _leafHash(Encode.cdcEncodeBytes32(header.appHash));
        leafHashes[11] = header.hasLastResultsHash ? _leafHash(Encode.cdcEncodeBytes32(header.lastResultsHash)) : _emptyLeafHash();
        leafHashes[12] = header.hasEvidenceHash ? _leafHash(Encode.cdcEncodeBytes32(header.evidenceHash)) : _emptyLeafHash();
        leafHashes[13] = _leafHash(Encode.cdcEncodeBytes(header.proposerAddress));

        return _merkleHashRange(leafHashes, 0, leafHashes.length);
    }

    function merkleHash(
        bytes[] memory bytesArray
    ) public pure returns (bytes32) {
        if (bytesArray.length == 0) {
            return bytes32(0);
        }

        bytes32[] memory leafHashes = new bytes32[](bytesArray.length);
        for (uint256 i = 0; i < bytesArray.length; i++) {
            leafHashes[i] = _leafHash(bytesArray[i]);
        }

        return _merkleHashRange(leafHashes, 0, bytesArray.length);
    }

    function getSlice(bytes[] memory bytesArray, uint256 from, uint256 to)
        public pure returns (bytes[] memory result) {
        require(from <= to && to <= bytesArray.length, "Invalid range");

        uint256 length = to - from;
        result = new bytes[](length);

        assembly {
            mcopy(add(result, 0x20), add(add(bytesArray, 0x20), mul(from, 0x20)), mul(length, 0x20))
        }
    }

    function nextPowerOfTwo(uint256 n) public pure returns (uint256) {
        if (n == 0) return 1;

        // Handle the case where n is already a power of 2
        if (n & (n - 1) == 0) return n;

        // Find the next power of 2
        uint256 power = 1;
        while (power < n) {
            power <<= 1;
        }
        return power;
    }

    function _merkleHashRange(
        bytes32[] memory leafHashes,
        uint256 from,
        uint256 to
    ) private pure returns (bytes32) {
        uint256 length = to - from;
        if (length == 1) {
            return leafHashes[from];
        }

        uint256 split = _splitPoint(length);
        bytes32 left = _merkleHashRange(leafHashes, from, from + split);
        bytes32 right = _merkleHashRange(leafHashes, from + split, to);
        return _innerHash(left, right);
    }

    function _splitPoint(uint256 n) private pure returns (uint256 split) {
        split = 1;
        while ((split << 1) < n) {
            split <<= 1;
        }
    }

    function _leafHash(bytes memory leaf) private pure returns (bytes32) {
        bytes memory prefixed = new bytes(leaf.length + 1);
        prefixed[0] = bytes1(0x00);
        _copyBytes(prefixed, 1, leaf);
        return sha256(prefixed);
    }

    function _emptyLeafHash() private pure returns (bytes32) {
        bytes memory prefixed = new bytes(1);
        prefixed[0] = bytes1(0x00);
        return sha256(prefixed);
    }

    function _innerHash(bytes32 left, bytes32 right) private pure returns (bytes32) {
        bytes memory prefixed = new bytes(65);
        prefixed[0] = bytes1(0x01);
        assembly {
            mstore(add(prefixed, 0x21), left)
            mstore(add(prefixed, 0x41), right)
        }
        return sha256(prefixed);
    }

    function _copyBytes(bytes memory out, uint256 dstOffset, bytes memory src) private pure returns (uint256) {
        uint256 len = src.length;
        if (len == 0) {
            return dstOffset;
        }

        assembly {
            mcopy(add(add(out, 0x20), dstOffset), add(src, 0x20), len)
        }

        return dstOffset + len;
    }
}
