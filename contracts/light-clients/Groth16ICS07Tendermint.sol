// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable gas-strict-inequalities

import { IICS07TendermintMsgs } from "./msgs/IICS07TendermintMsgs.sol";
import { IUpdateClientMsgs } from "./msgs/IUpdateClientMsgs.sol";
import { IMembershipMsgs } from "./msgs/IMembershipMsgs.sol";
import { IMisbehaviourMsgs } from "./msgs/IMisbehaviourMsgs.sol";
import { ILightClientMsgs } from "../msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";

import { IGroth16ICS07TendermintErrors } from "./errors/IGroth16ICS07TendermintErrors.sol";
import { IGroth16ICS07Tendermint } from "./IGroth16ICS07Tendermint.sol";
import { IMembership } from "../interfaces/IMembership.sol";
import { IMisbehaviour } from "../interfaces/IMisbehaviour.sol";
import { IUpdateClient } from "../interfaces/IUpdateClient.sol";
import { ILightClient } from "../interfaces/ILightClient.sol";
import { IVerifier } from "../interfaces/IVerifier.sol";

import { Paths } from "./utils/Paths.sol";
import { Encode } from "../utils/Encode.sol";
import { Header } from "../utils/Header.sol";
import { ChainId } from "../utils/ChainId.sol";
import { SSTORE2 } from "../utils/SSTORE2.sol";
import { AccessControl } from "@openzeppelin-contracts/access/AccessControl.sol";

/// @title Groth16 ICS07 Tendermint Light Client
/// @author srdtrk
/// @notice This contract implements an ICS07 IBC tendermint light client using gnark Groth16.
contract Groth16ICS07Tendermint is IGroth16ICS07TendermintErrors, IGroth16ICS07Tendermint, ILightClient, AccessControl {
    IVerifier private immutable VERIFIER;
    IMembership private immutable MEMBERSHIP;
    IMisbehaviour private immutable MISBEHAVIOUR;
    IUpdateClient private immutable UPDATE_CLIENT;

    /// @notice The ICS07Tendermint client state, exposed through `getClientState()`.
    IICS07TendermintMsgs.ClientState private clientState;
    /// @notice The mapping from height to consensus state keccak256 hashes.
    /// @dev Revision number need not be keyed as it is not allowed to change.
    mapping(uint64 height => bytes32 hash) private _consensusStateHashes;

    bytes32 private _pinnedValidatorsHash;
    address private _pinnedValidatorSetPointer;
    uint64 private _pinnedTotalVotingPower;
    uint16 private _pinnedEntryCount;
    mapping(uint64 height => PinnedValidatorSetSnapshot snapshot) private _pinnedValidatorSetSnapshots;
    uint64[] private _pinnedSnapshotHeights;

    uint32 private constant VALIDATOR_CACHE_MAGIC = 0x56414c34; // "VAL4"
    uint16 private constant MAX_VALIDATOR_COUNT = 180;
    uint256 private constant VALIDATOR_CACHE_HEADER_LEN = 48;
    uint256 private constant VALIDATOR_CACHE_ENTRY_LEN = 44;

    struct ValidatorCacheHeader {
        uint64 totalVotingPower;
        uint16 entryCount;
        bytes32 validatorsHashLeaf;
        uint16 nodeCount;
    }

    struct PinnedValidatorSetSnapshot {
        bytes32 validatorsHash;
        address pointer;
        uint64 totalVotingPower;
        uint16 entryCount;
    }

    bytes32 private constant PROOF_SUBMITTER_ROLE = keccak256("PROOF_SUBMITTER_ROLE");

    /// @inheritdoc IGroth16ICS07Tendermint
    bytes32 public immutable MISBEHAVIOUR_SUBMITTER_ROLE = keccak256("MISBEHAVIOUR_SUBMITTER_ROLE");

    /// @notice keccak256 of the client's chain ID, cached at construction.
    bytes32 internal immutable CHAIN_ID_HASH;
    /// @notice Header.Hash() leaf for the client's chain ID, cached at construction.
    bytes32 internal immutable CHAIN_ID_LEAF_HASH;

    /// @notice The constructor sets the program verification key and the initial client, consensus, and pinned
    /// validator states.
    constructor(
        address verifier,
        address membership_,
        address misbehaviour_,
        address updateClient_,
        bytes memory _clientState,
        bytes32 _consensusState,
        IICS07TendermintMsgs.ValidatorSet memory initialPinnedValidatorSet,
        address roleManager
    ) {
        clientState = abi.decode(_clientState, (IICS07TendermintMsgs.ClientState));
        CHAIN_ID_HASH = keccak256(bytes(clientState.chainId));
        CHAIN_ID_LEAF_HASH = Header.chainIdLeafHash(clientState.chainId);

        uint64 parsedRevision = ChainId.get(clientState.chainId).revisionNumber;
        require(
            parsedRevision == clientState.latestHeight.revisionNumber,
            MismatchedRevisionHeights(parsedRevision, clientState.latestHeight.revisionNumber)
        );
        _consensusStateHashes[clientState.latestHeight.revisionHeight] = _consensusState;

        VERIFIER = IVerifier(verifier);
        MEMBERSHIP = IMembership(membership_);
        MISBEHAVIOUR = IMisbehaviour(misbehaviour_);
        UPDATE_CLIENT = IUpdateClient(updateClient_);

        require(
            clientState.trustingPeriod + clientState.clockDrift <= clientState.unbondingPeriod,
            TrustingPeriodTooLong(clientState.trustingPeriod, clientState.unbondingPeriod)
        );

        _setPinnedValidatorSet(initialPinnedValidatorSet);
        _storePinnedValidatorSetSnapshot(clientState.latestHeight.revisionHeight);

        if (roleManager == address(0)) {
            _grantRole(PROOF_SUBMITTER_ROLE, address(0));
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, address(0));
        } else {
            _grantRole(DEFAULT_ADMIN_ROLE, roleManager);
            _grantRole(PROOF_SUBMITTER_ROLE, roleManager);
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, roleManager);
        }
    }

    /// @inheritdoc ILightClient
    function getClientState() external view returns (bytes memory) {
        return abi.encode(clientState);
    }

    function _getConsensusStateHash(uint64 revisionHeight) private view returns (bytes32) {
        bytes32 hash = _consensusStateHashes[revisionHeight];
        require(hash != 0, ConsensusStateNotFound());
        return hash;
    }

    /// @inheritdoc IGroth16ICS07Tendermint
    function getPinnedValidatorSet()
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers)
    {
        ValidatorCacheHeader memory cacheHeader;
        bytes memory cacheData;
        (cacheHeader, cacheData) = _readPinnedValidatorCache();

        indices = new uint32[](cacheHeader.entryCount);
        pubkeys = new bytes32[](cacheHeader.entryCount);
        votingPowers = new uint64[](cacheHeader.entryCount);
        for (uint16 i = 0; i < cacheHeader.entryCount; i++) {
            indices[i] = uint32(i);
            (votingPowers[i], pubkeys[i]) =
                _resolvePinnedValidator(_pinnedValidatorsHash, cacheData, cacheHeader.entryCount, i);
        }
    }

    /// @dev This function verifies the public values and forwards the proof to the Groth16 verifier.
    /// @inheritdoc ILightClient
    function updateClient(bytes calldata updateClientMsg)
        external
        notFrozen
        onlyProofSubmitter
        returns (ILightClientMsgs.UpdateResult)
    {
        IUpdateClientMsgs.MsgUpdateClient memory msg_ = abi.decode(updateClientMsg, (IUpdateClientMsgs.MsgUpdateClient));

        IUpdateClientMsgs.UpdateClientOutput memory output = UPDATE_CLIENT.updateClient(msg_);
        _validateUpdateClientOutput(output);

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult(output);
        _verifyPinnedBatchAndQuorum(msg_);

        if (updateResult == ILightClientMsgs.UpdateResult.Update) {
            require(
                output.newHeight.revisionHeight > clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate(clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            clientState.latestHeight = output.newHeight;
            _consensusStateHashes[output.newHeight.revisionHeight] = keccak256(abi.encode(output.newConsensusState));
            emit ClientUpdated(output.newHeight.revisionHeight);
        } else if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
            clientState.isFrozen = true;
            emit ClientFrozen();
        } else if (updateResult == ILightClientMsgs.UpdateResult.NoOp) {
            return ILightClientMsgs.UpdateResult.NoOp;
        }
        return updateResult;
    }

    function reAnchorPinnedSet(bytes calldata reAnchorMsg) external notFrozen onlyProofSubmitter {
        (
            IUpdateClientMsgs.MsgUpdateClient memory msg_,
            IICS07TendermintMsgs.ValidatorSet memory newPinnedValidatorSet
        ) = abi.decode(reAnchorMsg, (IUpdateClientMsgs.MsgUpdateClient, IICS07TendermintMsgs.ValidatorSet));
        IUpdateClientMsgs.UpdateClientOutput memory output = UPDATE_CLIENT.updateClient(msg_);
        _validateUpdateClientOutput(output);

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult(output);
        _verifyPinnedBatchAndQuorum(msg_);

        if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
            clientState.isFrozen = true;
            emit ClientFrozen();
            return;
        }

        bytes32 newValidatorsHash = Header.hashValSet(newPinnedValidatorSet);
        require(
            newValidatorsHash == msg_.proposedHeader.signedHeader.header.nextValidatorsHash,
            MismatchedValidatorHashes(msg_.proposedHeader.signedHeader.header.nextValidatorsHash, newValidatorsHash)
        );

        if (updateResult == ILightClientMsgs.UpdateResult.Update) {
            require(
                output.newHeight.revisionHeight > clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate(clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            _setPinnedValidatorSet(newPinnedValidatorSet);
            clientState.latestHeight = output.newHeight;
            _consensusStateHashes[output.newHeight.revisionHeight] = keccak256(abi.encode(output.newConsensusState));
            _storePinnedValidatorSetSnapshot(output.newHeight.revisionHeight);
            emit ClientUpdated(output.newHeight.revisionHeight);
            emit PinnedSetReAnchored(output.newHeight.revisionHeight, newValidatorsHash);
        } else if (updateResult == ILightClientMsgs.UpdateResult.NoOp) {
            require(
                output.newHeight.revisionHeight == clientState.latestHeight.revisionHeight,
                NonMonotonicHeightUpdate(clientState.latestHeight.revisionHeight, output.newHeight.revisionHeight)
            );
            _setPinnedValidatorSet(newPinnedValidatorSet);
            _storePinnedValidatorSetSnapshot(output.newHeight.revisionHeight);
            emit PinnedSetReAnchored(output.newHeight.revisionHeight, newValidatorsHash);
        }
    }

    function _verifyPinnedBatchAndQuorum(IUpdateClientMsgs.MsgUpdateClient memory msg_)
        internal
        returns (uint64 totalVotingPower, uint64 accumulatedVotingPower)
    {
        require(
            msg_.signerIndices.length == msg_.bucket && msg_.pinnedValidatorIndices.length == msg_.bucket
                && msg_.signerPubkeys.length == msg_.bucket && msg_.active.length == msg_.bucket,
            BatchLengthMismatch()
        );

        ValidatorCacheHeader memory cacheHeader;
        bytes memory cacheData;
        (cacheHeader, cacheData) = _readPinnedValidatorCache();
        totalVotingPower = _pinnedTotalVotingPower;
        uint256 seenPinned = 0;
        bool hasPrevCommitSigner = false;
        uint32 prevCommitSigner = 0;

        for (uint256 i = 0; i < msg_.bucket; i++) {
            if (!msg_.active[i]) {
                continue;
            }

            uint32 commitIdx = msg_.signerIndices[i];
            if (hasPrevCommitSigner) {
                require(commitIdx > prevCommitSigner, DuplicateSigner(commitIdx));
            }
            hasPrevCommitSigner = true;
            prevCommitSigner = commitIdx;

            uint32 pinnedIdx = msg_.pinnedValidatorIndices[i];
            require(pinnedIdx < cacheHeader.entryCount, SignerIndexOutOfRange(pinnedIdx));
            uint256 mask = uint256(1) << pinnedIdx;
            require((seenPinned & mask) == 0, DuplicateSigner(pinnedIdx));
            seenPinned |= mask;

            (uint64 votingPower, bytes32 pubKey) =
                _resolvePinnedValidator(_pinnedValidatorsHash, cacheData, cacheHeader.entryCount, pinnedIdx);
            require(pubKey == msg_.signerPubkeys[i], PubkeyMismatch(pinnedIdx));
            accumulatedVotingPower += votingPower;
        }

        require(
            uint256(accumulatedVotingPower) * 3 > uint256(totalVotingPower) * 2,
            InsufficientVotingPower(accumulatedVotingPower, totalVotingPower)
        );

        _requireProofSignersCommitSigs(
            msg_.proposedHeader.signedHeader.commit.commitSigs, msg_.signerIndices, msg_.active
        );
        _verifyUpdateBatchProof(msg_);
    }

    function _verifyMisbehaviourBatchAndQuorum(
        IICS07TendermintMsgs.Header memory header,
        IMisbehaviourMsgs.BatchProof memory proof_
    )
        internal
        returns (uint64 totalVotingPower, uint64 accumulatedVotingPower)
    {
        require(
            proof_.signerIndices.length == proof_.bucket && proof_.pinnedValidatorIndices.length == proof_.bucket
                && proof_.signerPubkeys.length == proof_.bucket && proof_.active.length == proof_.bucket,
            BatchLengthMismatch()
        );

        PinnedValidatorSetSnapshot memory snapshot = _pinnedValidatorSetSnapshotAt(header.trustedHeight.revisionHeight);
        ValidatorCacheHeader memory cacheHeader;
        bytes memory cacheData;
        (cacheHeader, cacheData) = _readPinnedValidatorCache(snapshot);
        totalVotingPower = snapshot.totalVotingPower;
        uint256 seenPinned = 0;
        bool hasPrevCommitSigner = false;
        uint32 prevCommitSigner = 0;

        for (uint256 i = 0; i < proof_.bucket; i++) {
            if (!proof_.active[i]) {
                continue;
            }

            uint32 commitIdx = proof_.signerIndices[i];
            if (hasPrevCommitSigner) {
                require(commitIdx > prevCommitSigner, DuplicateSigner(commitIdx));
            }
            hasPrevCommitSigner = true;
            prevCommitSigner = commitIdx;

            uint32 pinnedIdx = proof_.pinnedValidatorIndices[i];
            require(pinnedIdx < cacheHeader.entryCount, SignerIndexOutOfRange(pinnedIdx));
            uint256 mask = uint256(1) << pinnedIdx;
            require((seenPinned & mask) == 0, DuplicateSigner(pinnedIdx));
            seenPinned |= mask;

            (uint64 votingPower, bytes32 pubKey) =
                _resolvePinnedValidator(snapshot.validatorsHash, cacheData, cacheHeader.entryCount, pinnedIdx);
            require(pubKey == proof_.signerPubkeys[i], PubkeyMismatch(pinnedIdx));
            accumulatedVotingPower += votingPower;
        }

        require(
            uint256(accumulatedVotingPower) * 3 > uint256(totalVotingPower) * 2,
            InsufficientVotingPower(accumulatedVotingPower, totalVotingPower)
        );

        _requireProofSignersCommitSigs(header.signedHeader.commit.commitSigs, proof_.signerIndices, proof_.active);
        _verifyMisbehaviourBatchProof(header, proof_);
    }

    function _requireProofSignersCommitSigs(
        IICS07TendermintMsgs.CommitSig[] memory commitSigs,
        uint32[] memory signerIndices,
        bool[] memory active
    )
        private
        pure
    {
        require(signerIndices.length == active.length, BatchLengthMismatch());
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
        require(signerIndex < commitSigs.length, SignerIndexOutOfRange(signerIndex));
        require(
            commitSigs[signerIndex].flag == IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
            ProofSignerCommitSigMismatch(signerIndex)
        );
    }

    function _verifyUpdateBatchProof(IUpdateClientMsgs.MsgUpdateClient memory msg_) private {
        _verifyBatchProofForHeader(
            msg_.proposedHeader,
            msg_.bucket,
            msg_.proof,
            msg_.commitments,
            msg_.commitmentPok,
            msg_.signerPubkeys,
            msg_.active
        );
    }

    function _verifyMisbehaviourBatchProof(
        IICS07TendermintMsgs.Header memory header,
        IMisbehaviourMsgs.BatchProof memory proof_
    )
        private
    {
        _verifyBatchProofForHeader(
            header,
            proof_.bucket,
            proof_.proof,
            proof_.commitments,
            proof_.commitmentPok,
            proof_.signerPubkeys,
            proof_.active
        );
    }

    function _verifyBatchProofForHeader(
        IICS07TendermintMsgs.Header memory header,
        uint16 bucket,
        uint256[8] memory proof,
        uint256[2] memory commitments,
        uint256[2] memory commitmentPok,
        bytes32[] memory signerPubkeys,
        bool[] memory active
    )
        private
    {
        IVerifier.SharedBlock memory shared = _sharedBlockFromHeader(header);

        require(
            VERIFIER.verifyBatchProof(bucket, proof, commitments, commitmentPok, signerPubkeys, active, shared),
            ProofVerificationFailed()
        );
    }

    function _sharedBlockFromHeader(IICS07TendermintMsgs.Header memory header)
        private
        pure
        returns (IVerifier.SharedBlock memory shared)
    {
        IICS07TendermintMsgs.BlockCommit memory commit = header.signedHeader.commit;
        require(commit.height == header.signedHeader.header.height, InvalidHeaderHeight(commit.height));

        shared = IVerifier.SharedBlock({
            height: commit.height, round: uint64(commit.round), blockIDHash: commit.blockId.hashData
        });
    }

    function _setPinnedValidatorSet(IICS07TendermintMsgs.ValidatorSet memory validatorSet) private {
        bytes32 validatorsHash = Header.hashValSet(validatorSet);
        bytes memory cacheData = _buildPinnedValidatorSetCache(validatorsHash, validatorSet);
        ValidatorCacheHeader memory cacheHeader = _readValidatorCacheHeader(validatorsHash, cacheData);
        _pinnedValidatorSetPointer = SSTORE2.write(cacheData);
        _pinnedValidatorsHash = validatorsHash;
        _pinnedTotalVotingPower = cacheHeader.totalVotingPower;
        _pinnedEntryCount = cacheHeader.entryCount;
    }

    function _storePinnedValidatorSetSnapshot(uint64 height) private {
        bool exists = _pinnedValidatorSetSnapshots[height].pointer != address(0);
        _pinnedValidatorSetSnapshots[height] = PinnedValidatorSetSnapshot({
            validatorsHash: _pinnedValidatorsHash,
            pointer: _pinnedValidatorSetPointer,
            totalVotingPower: _pinnedTotalVotingPower,
            entryCount: _pinnedEntryCount
        });
        if (!exists) {
            _pinnedSnapshotHeights.push(height);
        }
    }

    function _pinnedValidatorSetSnapshotAt(uint64 height)
        private
        view
        returns (PinnedValidatorSetSnapshot memory snapshot)
    {
        uint256 len = _pinnedSnapshotHeights.length;
        if (len == 0) {
            revert ValidatorSetCacheMiss(bytes32(0));
        }
        uint256 left = 0;
        uint256 right = len - 1;
        while (left < right) {
            uint256 mid = (left + right + 1) / 2;
            if (_pinnedSnapshotHeights[mid] <= height) {
                left = mid;
            } else {
                right = mid - 1;
            }
        }
        uint64 checkpoint = _pinnedSnapshotHeights[left];
        require(checkpoint <= height, ValidatorSetCacheMiss(bytes32(0)));
        snapshot = _pinnedValidatorSetSnapshots[checkpoint];
        if (snapshot.pointer == address(0)) {
            revert ValidatorSetCacheMiss(snapshot.validatorsHash);
        }
    }

    function _buildPinnedValidatorSetCache(
        bytes32 validatorsHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet
    )
        private
        pure
        returns (bytes memory data)
    {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = validatorSet.validators;
        uint256 validatorCount = vals.length;
        if (validatorCount == 0) {
            revert BatchLengthMismatch();
        }
        if (validatorCount > MAX_VALIDATOR_COUNT) {
            revert ValidatorCountExceedsLimit(validatorCount, MAX_VALIDATOR_COUNT);
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
        require(totalVotingPower <= type(uint64).max, BatchLengthMismatch());
        _writeUint64(data, 4, uint64(totalVotingPower));
    }

    function _readPinnedValidatorCache()
        private
        view
        returns (ValidatorCacheHeader memory cacheHeader, bytes memory cacheData)
    {
        return _readPinnedValidatorCache(
            PinnedValidatorSetSnapshot({
                validatorsHash: _pinnedValidatorsHash,
                pointer: _pinnedValidatorSetPointer,
                totalVotingPower: _pinnedTotalVotingPower,
                entryCount: _pinnedEntryCount
            })
        );
    }

    function _readPinnedValidatorCache(PinnedValidatorSetSnapshot memory snapshot)
        private
        view
        returns (ValidatorCacheHeader memory cacheHeader, bytes memory cacheData)
    {
        if (snapshot.pointer == address(0)) {
            revert ValidatorSetCacheMiss(snapshot.validatorsHash);
        }
        cacheData = SSTORE2.read(snapshot.pointer, 0, _validatorCacheEntryDataLen(snapshot.entryCount));
        cacheHeader = _readValidatorCacheHeader(snapshot.validatorsHash, cacheData);
    }

    function _readValidatorCacheHeader(
        bytes32 validatorsHash,
        bytes memory cacheData
    )
        private
        pure
        returns (ValidatorCacheHeader memory cacheHeader)
    {
        if (cacheData.length < VALIDATOR_CACHE_HEADER_LEN) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
        if (_readUint32(cacheData, 0) != VALIDATOR_CACHE_MAGIC) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }

        cacheHeader = ValidatorCacheHeader({
            totalVotingPower: _readUint64(cacheData, 4),
            entryCount: _readUint16(cacheData, 12),
            validatorsHashLeaf: _readBytes32(cacheData, 14),
            nodeCount: _readUint16(cacheData, 46)
        });

        if (
            cacheHeader.entryCount == 0 || cacheHeader.entryCount > MAX_VALIDATOR_COUNT || cacheHeader.nodeCount != 0
                || cacheData.length != _validatorCacheEntryDataLen(cacheHeader.entryCount)
        ) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
    }

    function _validatorCacheEntryDataLen(uint16 entryCount) private pure returns (uint256) {
        return VALIDATOR_CACHE_HEADER_LEN + uint256(entryCount) * VALIDATOR_CACHE_ENTRY_LEN;
    }

    function _resolvePinnedValidator(
        bytes32 validatorsHash,
        bytes memory cacheData,
        uint16 entryCount,
        uint32 index
    )
        private
        pure
        returns (uint64 votingPower, bytes32 pubKey)
    {
        if (index >= entryCount) {
            revert CachedSignerNotFound(validatorsHash, index);
        }

        uint256 offset = VALIDATOR_CACHE_HEADER_LEN + uint256(index) * VALIDATOR_CACHE_ENTRY_LEN;
        if (_readUint32(cacheData, offset) != index) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
        votingPower = _readUint64(cacheData, offset + 4);
        pubKey = _readBytes32(cacheData, offset + 12);
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

    /// @inheritdoc ILightClient
    function verifyMembership(ILightClientMsgs.MsgVerifyMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        require(msg_.value.length > 0, EmptyValue());
        return _membership(
            msg_.height,
            msg_.kvPairs,
            msg_.merkleProofs,
            msg_.appHash,
            msg_.trustedConsensusState,
            msg_.membershipType,
            msg_.path,
            msg_.value
        );
    }

    /// @inheritdoc ILightClient
    function verifyNonMembership(ILightClientMsgs.MsgVerifyNonMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        return _membership(
            msg_.height,
            msg_.kvPairs,
            msg_.merkleProofs,
            msg_.appHash,
            msg_.trustedConsensusState,
            msg_.membershipType,
            msg_.path,
            bytes("")
        );
    }

    /// @notice The entrypoint for verifying (non)membership proof.
    /// @dev This is a non-membership proof if the value is empty.
    /// @param height The height of the proof.
    /// @param kvPairs The path and value of the key-value pair.
    /// @param merkleProofs The merkle proofs of membership.
    /// @param appHash The final hash value that needs to verify.
    /// @param membershipType Membership type.
    /// @return The timestamp of the trusted consensus state in unix seconds.
    function _membership(
        IICS02ClientMsgs.Height calldata height,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs,
        bytes32 appHash,
        IICS07TendermintMsgs.ConsensusState calldata trustedConsensusState,
        IMembershipMsgs.MembershipType membershipType,
        bytes[] calldata kvPath,
        bytes memory kvValue
    )
        private
        returns (uint256)
    {
        if (membershipType == IMembershipMsgs.MembershipType.Membership) {
            return _handleMembership(height, kvPairs, merkleProofs, appHash, trustedConsensusState, kvPath, kvValue);
        }

        revert UnknownMembershipType(uint8(membershipType));
    }

    /// @dev Standalone misbehaviour is accepted only when both conflicting
    /// headers carry proof-backed >2/3 Ed25519 quorums over the pinned validator set.
    /// @inheritdoc ILightClient
    function misbehaviour(bytes calldata misbehaviourMsg) external notFrozen onlyMisbehaviourSubmitter {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ =
            abi.decode(misbehaviourMsg, (IMisbehaviourMsgs.MsgSubmitMisbehaviour));

        IMisbehaviourMsgs.MisbehaviourOutput memory output = MISBEHAVIOUR.misbehaviour(
            msg_.clientState, msg_.misbehaviour, msg_.trustedConsensusState1, msg_.trustedConsensusState2, msg_.time
        );
        _validateMisbehaviourOutput(
            output, msg_.clientState, msg_.trustedConsensusState1, msg_.trustedConsensusState2, msg_.time
        );

        _verifyMisbehaviourBatchAndQuorum(msg_.misbehaviour.header1, msg_.proof1);
        _verifyMisbehaviourBatchAndQuorum(msg_.misbehaviour.header2, msg_.proof2);

        clientState.isFrozen = true;
        emit ClientFrozen();
    }

    /// @inheritdoc IGroth16ICS07Tendermint
    function unfreeze() external override(IGroth16ICS07Tendermint, ILightClient) onlyRole(DEFAULT_ADMIN_ROLE) {
        require(clientState.isFrozen, ClientNotFrozen());
        clientState.isFrozen = false;
        emit ClientUnfrozen();
    }

    /// @inheritdoc ILightClient
    function upgradeClient(bytes calldata) external pure {
        // NOTE: This feature will not be supported.
        revert FeatureNotSupported();
    }

    /// @notice Handles the `Groth16MembershipProof` proof type.
    /// @param height The height of the proof.
    /// @param kvPairs The path and value of the key-value pair that need verify.
    /// @param merkleProofs The merkle proofs of membership.
    /// @param appHash The final hash value that needs to verify.
    /// @param kvPath The path of the key-value pair in storage.
    /// @param kvValue The value of the key-value pair in storage.
    /// @return The timestamp of the trusted consensus state.
    function _handleMembership(
        IICS02ClientMsgs.Height calldata height,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs,
        bytes32 appHash,
        IICS07TendermintMsgs.ConsensusState calldata trustedConsensusState,
        bytes[] calldata kvPath,
        bytes memory kvValue
    )
        private
        returns (uint256)
    {
        require(
            kvPairs.length > 0 && kvPairs.length <= type(uint16).max,
            LengthIsOutOfRange(kvPairs.length, 1, type(uint16).max)
        );

        // validate provided inputs
        _validateMembershipInput(appHash, height.revisionHeight, trustedConsensusState);

        {
            // loop through the key-value pairs and validate them
            // if provided kv pairs inputs don't contains path and value
            // from contract state return error
            // if provided proofs contain kv path but value not match return an
            // error
            bool found = false;
            for (uint256 i = 0; i < kvPairs.length; ++i) {
                if (!Paths.equal(kvPairs[i].path, kvPath)) {
                    continue;
                }

                bytes memory value = kvPairs[i].value;
                require(
                    value.length == kvValue.length && keccak256(value) == keccak256(kvValue),
                    MembershipProofValueMismatch(kvValue, value)
                );

                found = true;
                break;
            }
            require(found, MembershipProofKeyNotFound(kvPath));
        }

        //verify membership of input proofs
        MEMBERSHIP.membership(appHash, kvPairs, merkleProofs);

        return _getTimestampInSeconds(trustedConsensusState);
    }

    /// @notice Validates the Membership input known values.
    /// @param commitmentRoot The commitment root of the provided input.
    /// @param proofHeight The height of the proof.
    /// @param trustedConsensusState The trusted consensus state
    function _validateMembershipInput(
        bytes32 commitmentRoot,
        uint64 proofHeight,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    )
        private
        view
    {
        bytes32 trustedConsensusStateHash = keccak256(abi.encode(trustedConsensusState));
        bytes32 storedConsensusStateHash = _getConsensusStateHash(proofHeight);
        require(
            trustedConsensusStateHash == storedConsensusStateHash,
            ConsensusStateHashMismatch(storedConsensusStateHash, trustedConsensusStateHash)
        );

        require(
            commitmentRoot == trustedConsensusState.root,
            ConsensusStateRootMismatch(trustedConsensusState.root, commitmentRoot)
        );

        _validateConsensusStateTrustingPeriod(trustedConsensusState.timestamp);
    }

    /// @notice Validates the trusted consensus state timestamp for membership proofs.
    /// @param consensusStateTimestamp The trusted consensus state timestamp in unix nanoseconds.
    function _validateConsensusStateTrustingPeriod(uint128 consensusStateTimestamp) private view {
        uint256 consensusStateTimestampSeconds = _nanosToSeconds(consensusStateTimestamp);
        require(
            // Membership proof freshness is defined against the destination chain clock.
            // forge-lint: disable-next-line(block-timestamp)
            consensusStateTimestampSeconds <= block.timestamp,
            ProofIsInTheFuture(block.timestamp, consensusStateTimestampSeconds)
        );

        // forge-lint: disable-next-line(unsafe-typecast)
        uint128 durationSinceConsensusState = uint128(block.timestamp - consensusStateTimestampSeconds);
        require(
            durationSinceConsensusState < clientState.trustingPeriod,
            InsufficientTrustingPeriod(durationSinceConsensusState, uint128(clientState.trustingPeriod))
        );
    }

    /// @notice Validates the Groth16ICS07UpdateClientOutput public values.
    /// @param output The public values.
    function _validateUpdateClientOutput(IUpdateClientMsgs.UpdateClientOutput memory output) private view {
        _validateClientStateAndTime(output.clientState, output.time);

        bytes32 outputConsensusStateHash = keccak256(abi.encode(output.trustedConsensusState));
        bytes32 storedConsensusStateHash = _getConsensusStateHash(output.trustedHeight.revisionHeight);
        require(
            outputConsensusStateHash == storedConsensusStateHash,
            ConsensusStateHashMismatch(storedConsensusStateHash, outputConsensusStateHash)
        );
    }

    /// @notice Validates the Groth16ICS07MisbehaviourOutput public values.
    /// @param output The public values.
    function _validateMisbehaviourOutput(
        IMisbehaviourMsgs.MisbehaviourOutput memory output,
        IICS07TendermintMsgs.ClientState memory clientState_,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState1,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState2,
        uint128 time
    )
        private
        view
    {
        _validateClientStateAndTime(clientState_, time);

        // make sure the trusted consensus state from header 1 is known (trusted) by matching it with the one in the
        // mapping
        bytes32 outputConsensusStateHash1 = keccak256(abi.encode(trustedConsensusState1));
        bytes32 storedConsensusStateHash1 = _getConsensusStateHash(output.trustedHeight1.revisionHeight);
        require(
            outputConsensusStateHash1 == storedConsensusStateHash1,
            ConsensusStateHashMismatch(storedConsensusStateHash1, outputConsensusStateHash1)
        );

        // make sure the trusted consensus state from header 2 is known (trusted) by matching it with the one in the
        // mapping
        bytes32 outputConsensusStateHash2 = keccak256(abi.encode(trustedConsensusState2));
        bytes32 storedConsensusStateHash2 = _getConsensusStateHash(output.trustedHeight2.revisionHeight);
        require(
            outputConsensusStateHash2 == storedConsensusStateHash2,
            ConsensusStateHashMismatch(storedConsensusStateHash2, outputConsensusStateHash2)
        );
    }

    /// @notice Validates the client state and time.
    /// @dev This function does not check the equality of the latest height and isFrozen.
    /// @param publicClientState The public client state.
    /// @param time The time in unix nanoseconds.
    function _validateClientStateAndTime(
        IICS07TendermintMsgs.ClientState memory publicClientState,
        uint128 time
    )
        private
        view
    {
        require(_nanosToSeconds(time) <= block.timestamp, ProofIsInTheFuture(block.timestamp, _nanosToSeconds(time)));
        require(
            block.timestamp - _nanosToSeconds(time) <= clientState.clockDrift,
            ProofIsTooOld(block.timestamp, _nanosToSeconds(time))
        );

        // Check client state equality
        // NOTE: We do not check the equality of latest height and isFrozen, this is because:
        // 1. Latest height can be updated by a frontrunner relayer in order to DOS the proof of another relayer.
        // 2. Each external call has the `notFrozen` modifier which checks if the client is frozen.
        // 3. The revision number is not allowed to change with us checking the chain-id and the implementation in the
        // gnark program.
        require(
            keccak256(bytes(publicClientState.chainId)) == CHAIN_ID_HASH,
            ChainIdMismatch(clientState.chainId, publicClientState.chainId)
        );
        require(
            publicClientState.trustLevel.numerator == clientState.trustLevel.numerator
                && publicClientState.trustLevel.denominator == clientState.trustLevel.denominator,
            TrustThresholdMismatch(
                clientState.trustLevel.numerator,
                clientState.trustLevel.denominator,
                publicClientState.trustLevel.numerator,
                publicClientState.trustLevel.denominator
            )
        );
        require(
            publicClientState.trustingPeriod == clientState.trustingPeriod,
            TrustingPeriodMismatch(clientState.trustingPeriod, publicClientState.trustingPeriod)
        );
        require(
            publicClientState.unbondingPeriod == clientState.unbondingPeriod,
            UnbondingPeriodMismatch(clientState.unbondingPeriod, publicClientState.unbondingPeriod)
        );
        require(
            publicClientState.clockDrift == clientState.clockDrift,
            ClockDriftMismatch(clientState.clockDrift, publicClientState.clockDrift)
        );
    }

    /// @notice Checks for basic misbehaviour.
    /// @dev This function checks if the consensus state at the new height is different than the one in the mapping
    /// @dev or if the timestamp is not increasing.
    /// @dev If any of these conditions are met, it returns a Misbehaviour UpdateResult.
    /// @param output The public values of the update client program.
    /// @return The result of the update.
    function _checkUpdateResult(IUpdateClientMsgs.UpdateClientOutput memory output)
        private
        view
        returns (ILightClientMsgs.UpdateResult)
    {
        bytes32 consensusStateHash = _consensusStateHashes[output.newHeight.revisionHeight];
        if (consensusStateHash == bytes32(0)) {
            // No consensus state at the new height, so no misbehaviour
            return ILightClientMsgs.UpdateResult.Update;
        } else if (
            consensusStateHash != keccak256(abi.encode(output.newConsensusState))
                || output.trustedConsensusState.timestamp >= output.newConsensusState.timestamp
        ) {
            // The consensus state at the new height is different than the one in the mapping
            // or the timestamp is not increasing
            return ILightClientMsgs.UpdateResult.Misbehaviour;
        } else {
            // The consensus state at the new height is the same as the one in the mapping
            return ILightClientMsgs.UpdateResult.NoOp;
        }
    }

    /// @notice Returns the timestamp of the trusted consensus state in unix seconds.
    /// @param consensusState The consensus state.
    /// @return The timestamp of the trusted consensus state in unix seconds.
    function _getTimestampInSeconds(IICS07TendermintMsgs.ConsensusState memory consensusState)
        private
        pure
        returns (uint256)
    {
        return _nanosToSeconds(consensusState.timestamp);
    }

    /// @notice Converts nanoseconds to seconds.
    /// @param nanos The nanoseconds.
    /// @return The seconds.
    function _nanosToSeconds(uint256 nanos) private pure returns (uint256) {
        return nanos / 1e9;
    }

    /// @notice Modifier to check if the client is not frozen.
    modifier notFrozen() {
        require(!clientState.isFrozen, FrozenClientState());
        _;
    }

    /// @notice Modifier to check if the caller has the proof submitter role or if the role is permitted for anyone.
    modifier onlyProofSubmitter() {
        if (!hasRole(PROOF_SUBMITTER_ROLE, address(0))) {
            _checkRole(PROOF_SUBMITTER_ROLE);
        }
        _;
    }

    /// @notice Modifier to check if the caller has the misbehaviour submitter role or if the role is permitted for
    /// anyone.
    modifier onlyMisbehaviourSubmitter() {
        if (!hasRole(MISBEHAVIOUR_SUBMITTER_ROLE, address(0))) {
            _checkRole(MISBEHAVIOUR_SUBMITTER_ROLE);
        }
        _;
    }
}
