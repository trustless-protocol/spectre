// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IICS07TendermintMsgs } from "../msgs/IICS07TendermintMsgs.sol";
import { ISpectreClientErrors } from "../errors/ISpectreClientErrors.sol";
import { SpectreStore } from "./SpectreStore.sol";
import { Header } from "../../utils/Header.sol";
import { SSTORE2 } from "../../utils/SSTORE2.sol";

/// @title ValidatorSetLib
/// @notice Codec + lookup helpers for the SSTORE2-packed pinned validator cache, plus the
///         per-signer commit-signature checks used by SpectreClient's quorum accounting.
/// @dev The byte layout MUST stay identical to the relayer prover — do not "improve" it.
library ValidatorSetLib {
    uint32 internal constant VALIDATOR_CACHE_MAGIC = 0x56414c34; // "VAL4"
    uint16 internal constant MAX_VALIDATOR_COUNT = 180;
    uint256 internal constant VALIDATOR_CACHE_HEADER_LEN = 48;
    uint256 internal constant VALIDATOR_CACHE_ENTRY_LEN = 44;

    struct ValidatorCacheHeader {
        uint64 totalVotingPower;
        uint16 entryCount;
        bytes32 validatorsHashLeaf;
        uint16 nodeCount;
    }

    /// @notice Packs a validator set into the SSTORE2 cache byte layout.
    function buildCache(
        bytes32 validatorsHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet
    )
        internal
        pure
        returns (bytes memory data)
    {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = validatorSet.validators;
        uint256 validatorCount = vals.length;
        if (validatorCount == 0) {
            revert ISpectreClientErrors.BatchLengthMismatch();
        }
        if (validatorCount > MAX_VALIDATOR_COUNT) {
            revert ISpectreClientErrors.ValidatorCountExceedsLimit(validatorCount, MAX_VALIDATOR_COUNT);
        }

        uint256 totalVotingPower = 0;
        data = new bytes(VALIDATOR_CACHE_HEADER_LEN + validatorCount * VALIDATOR_CACHE_ENTRY_LEN);
        _writeUint32(data, 0, VALIDATOR_CACHE_MAGIC);
        _writeUint16(data, 12, validatorCount);
        _writeBytes32(data, 14, Header.bytes32LeafHash(validatorsHash));
        _writeUint16(data, 46, 0);

        uint256 offset = VALIDATOR_CACHE_HEADER_LEN;
        for (uint256 i = 0; i < vals.length; i++) {
            IICS07TendermintMsgs.ValidatorInfo memory val = vals[i];
            totalVotingPower += val.votingPower;
            _writeUint32(data, offset, uint32(i));
            _writeUint64(data, offset + 4, val.votingPower);
            _writeBytes32(data, offset + 12, val.pubKey);
            offset += VALIDATOR_CACHE_ENTRY_LEN;
        }
        require(totalVotingPower <= type(uint64).max, ISpectreClientErrors.BatchLengthMismatch());
        _writeUint64(data, 4, uint64(totalVotingPower));
    }

    /// @notice Reads and validates the cache pointed to by a snapshot.
    function readCache(SpectreStore.PinnedValidatorSetSnapshot memory snapshot)
        internal
        view
        returns (ValidatorCacheHeader memory cacheHeader, bytes memory cacheData)
    {
        if (snapshot.pointer == address(0)) {
            revert ISpectreClientErrors.ValidatorSetCacheMiss(snapshot.validatorsHash);
        }
        cacheData = SSTORE2.read(snapshot.pointer, 0, validatorCacheEntryDataLen(snapshot.entryCount));
        cacheHeader = readValidatorCacheHeader(snapshot.validatorsHash, cacheData);
    }

    function readValidatorCacheHeader(
        bytes32 validatorsHash,
        bytes memory cacheData
    )
        internal
        pure
        returns (ValidatorCacheHeader memory cacheHeader)
    {
        if (cacheData.length < VALIDATOR_CACHE_HEADER_LEN) {
            revert ISpectreClientErrors.CachedValidatorSetCorrupted(validatorsHash);
        }
        if (_readUint32(cacheData, 0) != VALIDATOR_CACHE_MAGIC) {
            revert ISpectreClientErrors.CachedValidatorSetCorrupted(validatorsHash);
        }

        cacheHeader = ValidatorCacheHeader({
            totalVotingPower: _readUint64(cacheData, 4),
            entryCount: _readUint16(cacheData, 12),
            validatorsHashLeaf: _readBytes32(cacheData, 14),
            nodeCount: _readUint16(cacheData, 46)
        });

        if (
            cacheHeader.entryCount == 0 || cacheHeader.entryCount > MAX_VALIDATOR_COUNT || cacheHeader.nodeCount != 0
                || cacheData.length != validatorCacheEntryDataLen(cacheHeader.entryCount)
                || cacheHeader.validatorsHashLeaf != Header.bytes32LeafHash(validatorsHash)
        ) {
            revert ISpectreClientErrors.CachedValidatorSetCorrupted(validatorsHash);
        }
    }

    function validatorCacheEntryDataLen(uint16 entryCount) internal pure returns (uint256) {
        return VALIDATOR_CACHE_HEADER_LEN + uint256(entryCount) * VALIDATOR_CACHE_ENTRY_LEN;
    }

    function resolveValidator(
        bytes32 validatorsHash,
        bytes memory cacheData,
        uint16 entryCount,
        uint32 index
    )
        internal
        pure
        returns (uint64 votingPower, bytes32 pubKey)
    {
        if (index >= entryCount) {
            revert ISpectreClientErrors.CachedSignerNotFound(validatorsHash, index);
        }

        uint256 offset = VALIDATOR_CACHE_HEADER_LEN + uint256(index) * VALIDATOR_CACHE_ENTRY_LEN;
        if (_readUint32(cacheData, offset) != index) {
            revert ISpectreClientErrors.CachedValidatorSetCorrupted(validatorsHash);
        }
        votingPower = _readUint64(cacheData, offset + 4);
        pubKey = _readBytes32(cacheData, offset + 12);
    }

    /// @notice Requires that every active proof signer corresponds to a COMMIT slot in the header commit.
    function requireProofSignersCommitSigs(
        IICS07TendermintMsgs.CommitSig[] memory commitSigs,
        uint32[] memory signerIndices,
        bool[] memory active
    )
        internal
        pure
    {
        require(signerIndices.length == active.length, ISpectreClientErrors.BatchLengthMismatch());
        for (uint256 i = 0; i < signerIndices.length; i++) {
            if (!active[i]) {
                continue;
            }
            _requireProofSignerCommitSig(commitSigs, signerIndices[i]);
        }
    }

    function _requireProofSignerCommitSig(
        IICS07TendermintMsgs.CommitSig[] memory commitSigs,
        uint32 signerIndex
    )
        private
        pure
    {
        require(signerIndex < commitSigs.length, ISpectreClientErrors.SignerIndexOutOfRange(signerIndex));
        require(
            commitSigs[signerIndex].flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
            ISpectreClientErrors.ProofSignerCommitSigMismatch(signerIndex)
        );
    }

    function _writeUint16(bytes memory data, uint256 offset, uint256 value) private pure {
        assembly ("memory-safe") {
            let ptr := add(add(data, 0x20), offset)
            mstore8(ptr, shr(8, value))
            mstore8(add(ptr, 1), value)
        }
    }

    function _writeUint32(bytes memory data, uint256 offset, uint32 value) private pure {
        assembly ("memory-safe") {
            let ptr := add(add(data, 0x20), offset)
            mstore8(ptr, shr(24, value))
            mstore8(add(ptr, 1), shr(16, value))
            mstore8(add(ptr, 2), shr(8, value))
            mstore8(add(ptr, 3), value)
        }
    }

    function _writeUint64(bytes memory data, uint256 offset, uint64 value) private pure {
        assembly ("memory-safe") {
            let ptr := add(add(data, 0x20), offset)
            mstore8(ptr, shr(56, value))
            mstore8(add(ptr, 1), shr(48, value))
            mstore8(add(ptr, 2), shr(40, value))
            mstore8(add(ptr, 3), shr(32, value))
            mstore8(add(ptr, 4), shr(24, value))
            mstore8(add(ptr, 5), shr(16, value))
            mstore8(add(ptr, 6), shr(8, value))
            mstore8(add(ptr, 7), value)
        }
    }

    function _writeBytes32(bytes memory data, uint256 offset, bytes32 value) private pure {
        assembly ("memory-safe") {
            mstore(add(add(data, 0x20), offset), value)
        }
    }

    function _readUint16(bytes memory data, uint256 offset) private pure returns (uint16 value) {
        assembly ("memory-safe") {
            value := shr(240, mload(add(add(data, 0x20), offset)))
        }
    }

    function _readUint32(bytes memory data, uint256 offset) private pure returns (uint32 value) {
        assembly ("memory-safe") {
            value := shr(224, mload(add(add(data, 0x20), offset)))
        }
    }

    function _readUint64(bytes memory data, uint256 offset) private pure returns (uint64 value) {
        assembly ("memory-safe") {
            value := shr(192, mload(add(add(data, 0x20), offset)))
        }
    }

    function _readBytes32(bytes memory data, uint256 offset) private pure returns (bytes32 value) {
        assembly ("memory-safe") {
            value := mload(add(add(data, 0x20), offset))
        }
    }
}
