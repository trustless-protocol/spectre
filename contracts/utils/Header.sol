pragma solidity ^0.8.0;

import { IICS07TendermintMsgs } from "../light-clients/msgs/IICS07TendermintMsgs.sol";
import { Encode } from "./Encode.sol";
library Header {

    function hashValSet(
        IICS07TendermintMsgs.ValidatorSet memory valset
    ) pure public returns (bytes32) {
        bytes[] memory validatorBytes = new bytes[](valset.validators.length);
        for (uint256 i = 0; i < valset.validators.length; i++) {
            bytes memory validatorHash = Encode.encodeValidator(
                IICS07TendermintMsgs.SimpleValidator({
                    pubKey: valset.validators[i].pubKey,
                    votingPower: valset.validators[i].votingPower
                })
            );
            validatorBytes[i] = validatorHash;

        }

        bytes32 root = merkleHash(validatorBytes);
        return root;
    }

    function hashHeader(
        IICS07TendermintMsgs.BlockHeader memory header
    ) pure public returns (bytes32) {
        // CometBFT Header.Hash() ALWAYS hashes 14 fields via merkle.
        // Missing fields are encoded as empty bytes (nil in Go).
        bytes[] memory headerBytes = new bytes[](14);

        // Field 0: Version (proto.Marshal)
        headerBytes[0] = Encode.encodeVersion(header.version);

        // Field 1: ChainID (cdcEncode → StringValue)
        headerBytes[1] = Encode.cdcEncodeString(header.chainId);

        // Field 2: Height (cdcEncode → Int64Value)
        headerBytes[2] = Encode.cdcEncodeInt64(uint256(header.height));

        // Field 3: Time (StdTimeMarshal → Timestamp)
        headerBytes[3] = Encode.encodeTimestamp(header.time);

        // Field 4: LastBlockId (proto.Marshal, empty if not present)
        if (header.hasLastBlockId) {
            headerBytes[4] = Encode.encodeBlockId(header.lastBlockId);
        } else {
            headerBytes[4] = new bytes(0);
        }

        // Field 5: LastCommitHash (cdcEncode → BytesValue)
        if (header.hasLastCommitHash) {
            headerBytes[5] = Encode.cdcEncodeBytes32(header.lastCommitHash);
        } else {
            headerBytes[5] = new bytes(0);
        }

        // Field 6: DataHash (cdcEncode → BytesValue)
        if (header.hasDataHash) {
            headerBytes[6] = Encode.cdcEncodeBytes32(header.dataHash);
        } else {
            headerBytes[6] = new bytes(0);
        }

        // Field 7-9: Always present hashes (cdcEncode → BytesValue)
        headerBytes[7] = Encode.cdcEncodeBytes32(header.validatorsHash);
        headerBytes[8] = Encode.cdcEncodeBytes32(header.nextValidatorsHash);
        headerBytes[9] = Encode.cdcEncodeBytes32(header.consensusHash);

        // Field 10: AppHash (cdcEncode → BytesValue)
        headerBytes[10] = Encode.cdcEncodeBytes32(header.appHash);

        // Field 11: LastResultsHash (cdcEncode → BytesValue)
        if (header.hasLastResultsHash) {
            headerBytes[11] = Encode.cdcEncodeBytes32(header.lastResultsHash);
        } else {
            headerBytes[11] = new bytes(0);
        }

        // Field 12: EvidenceHash (cdcEncode → BytesValue)
        if (header.hasEvidenceHash) {
            headerBytes[12] = Encode.cdcEncodeBytes32(header.evidenceHash);
        } else {
            headerBytes[12] = new bytes(0);
        }

        // Field 13: ProposerAddress (cdcEncode → BytesValue)
        headerBytes[13] = Encode.cdcEncodeBytes(header.proposerAddress);

        return merkleHash(headerBytes);
    }

    function merkleHash(
        bytes[] memory bytesArray
    ) public pure returns (bytes32) {
        if (bytesArray.length == 0) {
            return bytes32(0);
        }

        // tmhash(0x00 || leaf) — 1-byte leaf prefix per Tendermint spec
        if (bytesArray.length == 1) {
            return sha256(abi.encodePacked(bytes1(0x00), bytesArray[0]));
        }

        uint256 split = nextPowerOfTwo(bytesArray.length) / 2;
        bytes32 left = merkleHash(getSlice(bytesArray, 0, split));
        bytes32 right = merkleHash(getSlice(bytesArray, split, bytesArray.length));

        // tmhash(0x01 || left || right) — 1-byte inner prefix per Tendermint spec
        return sha256(abi.encodePacked(bytes1(0x01), left, right));
    }

    function getSlice(bytes[] memory bytesArray, uint256 from, uint256 to)
        public pure returns (bytes[] memory result) {
        require(from <= to && to <= bytesArray.length, "Invalid range");
        
        uint256 length = to - from;
        result = new bytes[](length);
        
        assembly {
            let src := add(add(bytesArray, 0x20), mul(from, 0x20))
            let dst := add(result, 0x20)
            let size := mul(length, 0x20)
            
            for { let i := 0 } lt(i, size) { i := add(i, 0x20) } {
                mstore(add(dst, i), mload(add(src, i)))
            }
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
}