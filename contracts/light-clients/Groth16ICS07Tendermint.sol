// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable gas-strict-inequalities

import { IICS07TendermintMsgs } from "./msgs/IICS07TendermintMsgs.sol";
import { IUpdateClientMsgs } from "./msgs/IUpdateClientMsgs.sol";
import { IMembershipMsgs } from "./msgs/IMembershipMsgs.sol";
import { IMisbehaviourMsgs } from "./msgs/IMisbehaviourMsgs.sol";
import { IGroth16Msgs } from "./msgs/IGroth16Msgs.sol";
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
import { Multicall } from "@openzeppelin-contracts/utils/Multicall.sol";
import { TransientSlot } from "@openzeppelin-contracts/utils/TransientSlot.sol";
import { AccessControl } from "@openzeppelin-contracts/access/AccessControl.sol";

/// @title Groth16 ICS07 Tendermint Light Client
/// @author srdtrk
/// @notice This contract implements an ICS07 IBC tendermint light client using gnark Groth16.
contract Groth16ICS07Tendermint is
    IGroth16ICS07TendermintErrors,
    IGroth16ICS07Tendermint,
    ILightClient,
    Multicall,
    AccessControl
{
    using TransientSlot for *;

    IVerifier private immutable VERIFIER;
    IMembership private immutable MEMBERSHIP;
    IMisbehaviour private immutable MISBEHAVIOUR;
    IUpdateClient private immutable UPDATE_CLIENT;

    /// @notice The ICS07Tendermint client state, exposed through `getClientState()`.
    IICS07TendermintMsgs.ClientState private clientState;
    /// @notice The mapping from height to consensus state keccak256 hashes.
    /// @dev Revision number need not be keyed as it is not allowed to change.
    mapping(uint64 height => bytes32 hash) private _consensusStateHashes;
    /// @notice Latest packed validator metadata keyed by the CometBFT validators hash.
    bytes32 private _cachedValidatorSetHash;
    bytes32 private _cachedValidatorsHashLeaf;
    address private _cachedValidatorSetPointer;
    uint64 private _cachedValidatorTotalVotingPower;
    uint16 private _cachedValidatorEntryCount;
    uint256 private _validatorPubKeyOverrideBits;
    uint256 private _validatorVotingPowerOverrideBits;
    mapping(uint16 index => bytes32 pubKey) private _validatorPubKeyOverrides;
    mapping(uint16 index => uint64 votingPower) private _validatorVotingPowerOverrides;

    uint32 private constant VALIDATOR_CACHE_MAGIC = 0x56414c34; // "VAL4"
    /// @dev Hard active-validator limit for this client. The delta cache stores per-index
    ///      override bitmaps in uint256 words, and the packed SSTORE2 cache stays below
    ///      EIP-170 with margin at this bound.
    uint16 private constant MAX_VALIDATOR_COUNT = 180;
    uint16 private constant MAX_DELTA_LEAF_COUNT = 16;
    uint256 private constant VALIDATOR_CACHE_HEADER_LEN = 48;
    uint256 private constant VALIDATOR_CACHE_ENTRY_LEN = 44;
    uint256 private constant VALIDATOR_CACHE_NODE_LEN = 32;

    struct ValidatorCacheHeader {
        uint64 totalVotingPower;
        uint16 entryCount;
        bytes32 validatorsHashLeaf;
        uint16 nodeCount;
    }

    struct DeltaRecomputeContext {
        bytes32 validatorsHash;
        bytes baseData;
        uint16 validatorCount;
        IUpdateClientMsgs.ValidatorSetDelta delta;
        bytes32[] newLeafHashes;
        uint256 pubKeyBits;
        uint256 votingPowerBits;
        uint256 overrideBits;
    }

    uint16 private constant ALLOWED_CLOCK_DRIFT = 30 minutes;

    bytes32 private constant PROOF_SUBMITTER_ROLE = keccak256("PROOF_SUBMITTER_ROLE");

    /// @inheritdoc IGroth16ICS07Tendermint
    bytes32 public immutable MISBEHAVIOUR_SUBMITTER_ROLE = keccak256("MISBEHAVIOUR_SUBMITTER_ROLE");

    /// @notice keccak256 of the client's chain ID, cached at construction.
    /// @dev The chain ID never changes for the lifetime of the client, so the
    ///      per-update equality check reads this immutable instead of hashing
    ///      the storage string on every call.
    bytes32 internal immutable CHAIN_ID_HASH;
    /// @notice Header.Hash() leaf for the client's chain ID, cached at construction.
    bytes32 internal immutable CHAIN_ID_LEAF_HASH;

    /// @notice The constructor sets the program verification key and the initial client and consensus states.
    /// @param verifier The address of the Groth16 verifier contract.
    /// @param _clientState The encoded initial client state.
    /// @param _consensusState The encoded initial consensus state.
    /// @param roleManager Manages the proof submitters and can submit proofs. Should be the ICS26Router if used in IBC.
    constructor(
        address verifier,
        address membership_,
        address misbehaviour_,
        address updateClient_,
        bytes memory _clientState,
        bytes32 _consensusState,
        address roleManager
    ) {
        clientState = abi.decode(_clientState, (IICS07TendermintMsgs.ClientState));
        CHAIN_ID_HASH = keccak256(bytes(clientState.chainId));
        CHAIN_ID_LEAF_HASH = Header.chainIdLeafHash(clientState.chainId);

        // updateClient/misbehaviour now read the chain-ID revision from
        // clientState.latestHeight.revisionNumber instead of re-parsing the
        // chain-ID string on every call. Assert the two agree at construction so
        // a misconfigured client fails fast at deploy time rather than silently
        // using the wrong revision (defense-in-depth, suggested in review of #75).
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
            clientState.trustingPeriod + ALLOWED_CLOCK_DRIFT <= clientState.unbondingPeriod,
            TrustingPeriodTooLong(clientState.trustingPeriod, clientState.unbondingPeriod)
        );

        if (roleManager == address(0)) {
            _grantRole(PROOF_SUBMITTER_ROLE, address(0)); // Allow anyone to submit proofs
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, address(0)); // Allow anyone to submit misbehaviour
        } else {
            _grantRole(DEFAULT_ADMIN_ROLE, roleManager); // Allow the role manager to manage roles
            _grantRole(PROOF_SUBMITTER_ROLE, roleManager); // Allow the role manager to submit proofs
            _grantRole(MISBEHAVIOUR_SUBMITTER_ROLE, roleManager); // Allow the role manager to submit misbehaviour
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
    function getCachedValidatorSet(bytes32 validatorsHash)
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers)
    {
        address pointer = _validatorCachePointer(validatorsHash);
        if (pointer == address(0)) {
            return (new uint32[](0), new bytes32[](0), new uint64[](0));
        }
        bytes memory cacheData = SSTORE2.read(pointer, 0, _validatorCacheEntryDataLen(_cachedValidatorEntryCount));
        ValidatorCacheHeader memory cacheHeader = _readValidatorCacheHeader(validatorsHash, cacheData);

        indices = new uint32[](cacheHeader.entryCount);
        pubkeys = new bytes32[](cacheHeader.entryCount);
        votingPowers = new uint64[](cacheHeader.entryCount);
        uint256 pubKeyBits = _validatorPubKeyOverrideBits;
        uint256 votingPowerBits = _validatorVotingPowerOverrideBits;
        for (uint16 i = 0; i < cacheHeader.entryCount; i++) {
            indices[i] = uint32(i);
            (votingPowers[i], pubkeys[i]) = _resolveCurrentCachedValidator(
                validatorsHash, cacheData, cacheHeader.entryCount, pubKeyBits, votingPowerBits, i
            );
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
        (
            bool currentCached,
            bytes32 currentValidatorsHash,
            bytes32 currentValidatorsHashLeaf,
            bool cacheCurrentValidatorSet
        ) = _prepareUpdateClientMessage(msg_);

        IUpdateClientMsgs.UpdateClientOutput memory output = currentCached
            ? UPDATE_CLIENT.updateClientCachedCurrentWithHeaderCache(
                msg_, CHAIN_ID_LEAF_HASH, currentValidatorsHashLeaf
            )
            : UPDATE_CLIENT.updateClientResolved(msg_);

        _validateUpdateClientOutput(output);

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult(output);
        uint64 totalVotingPower =
            currentCached ? _verifyCachedBatchAndQuorum(currentValidatorsHash, msg_) : _verifyFullBatchAndQuorum(msg_);
        if (cacheCurrentValidatorSet) {
            _cacheValidatorSet(currentValidatorsHash, msg_, totalVotingPower);
        }
        if (updateResult == ILightClientMsgs.UpdateResult.Update) {
            // adding the new consensus state to the mapping
            if (output.newHeight.revisionHeight > clientState.latestHeight.revisionHeight) {
                clientState.latestHeight = output.newHeight;
            }
            _consensusStateHashes[output.newHeight.revisionHeight] = keccak256(abi.encode(output.newConsensusState));
        } else if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
            clientState.isFrozen = true;
        } else if (updateResult == ILightClientMsgs.UpdateResult.NoOp) {
            return ILightClientMsgs.UpdateResult.NoOp;
        }
        return updateResult;
    }

    function _prepareUpdateClientMessage(IUpdateClientMsgs.MsgUpdateClient memory msg_)
        private
        returns (
            bool currentCached,
            bytes32 currentValidatorsHash,
            bytes32 currentValidatorsHashLeaf,
            bool cacheCurrentValidatorSet
        )
    {
        currentValidatorsHash = msg_.proposedHeader.signedHeader.header.validatorsHash;
        bytes32 trustedNextValidatorsHash = msg_.trustedConsensusState.nextValidatorsHash;

        bool adjacent = _isAdjacentUpdate(msg_.proposedHeader);
        (currentCached, currentValidatorsHashLeaf) = _getUsableValidatorHashLeaf(currentValidatorsHash);

        if (!currentCached) {
            if (msg_.currentValidatorSetDelta.baseValidatorsHash != bytes32(0)) {
                currentValidatorsHashLeaf =
                    _cacheDeltaValidatorSet(currentValidatorsHash, msg_.currentValidatorSetDelta);
                currentCached = true;
                msg_.proposedHeader.validatorSet = _emptyValidatorSet();
            } else {
                _validateSuppliedValidatorSetHash(currentValidatorsHash, msg_.proposedHeader.validatorSet);
                cacheCurrentValidatorSet = true;
            }
        } else {
            msg_.proposedHeader.validatorSet = _emptyValidatorSet();
        }

        if (adjacent) {
            msg_.proposedHeader.trustedNextValidatorSet = _emptyValidatorSet();
            return (currentCached, currentValidatorsHash, currentValidatorsHashLeaf, cacheCurrentValidatorSet);
        }

        if (!currentCached && trustedNextValidatorsHash == currentValidatorsHash) {
            msg_.proposedHeader.trustedNextValidatorSet = msg_.proposedHeader.validatorSet;
        } else if (!currentCached) {
            _validateSuppliedValidatorSetHash(trustedNextValidatorsHash, msg_.proposedHeader.trustedNextValidatorSet);
        }
        // When `currentCached == true` and the update is non-adjacent, the
        // caller-supplied `trustedNextValidatorSet` is intentionally left
        // unvalidated here. Validation is delegated to
        // `UPDATE_CLIENT.updateClientCachedCurrentWithHeaderCache` →
        // `verifyHeaderCachedCurrent`, which checks
        // `Header.hashValSet(trustedNextValidatorSet) ==
        // trustedConsensusState.nextValidatorsHash` before any header decision.
        // Do NOT remove that delegated check without restoring an equivalent
        // guard here, otherwise non-adjacent cached updates would accept an
        // arbitrary next-validator set.

        return (currentCached, currentValidatorsHash, currentValidatorsHashLeaf, cacheCurrentValidatorSet);
    }

    function _isAdjacentUpdate(IICS07TendermintMsgs.Header memory header) private pure returns (bool) {
        return header.signedHeader.header.height == header.trustedHeight.revisionHeight + 1;
    }

    function _validateSuppliedValidatorSetHash(
        bytes32 expectedHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet
    )
        private
        pure
    {
        uint256 validatorCount = validatorSet.validators.length;
        if (validatorCount == 0) {
            revert ValidatorSetCacheMiss(expectedHash);
        }
        if (validatorCount > MAX_VALIDATOR_COUNT) {
            revert ValidatorCountExceedsLimit(validatorCount, MAX_VALIDATOR_COUNT);
        }
        bytes32 actualHash = Header.hashValSet(validatorSet);
        require(actualHash == expectedHash, MismatchedValidatorHashes(expectedHash, actualHash));
    }

    function _cacheDeltaValidatorSet(
        bytes32 validatorsHash,
        IUpdateClientMsgs.ValidatorSetDelta memory delta
    )
        private
        returns (bytes32 validatorsHashLeaf)
    {
        address basePointer = _validatorCachePointer(delta.baseValidatorsHash);
        if (basePointer == address(0)) {
            revert ValidatorSetCacheMiss(delta.baseValidatorsHash);
        }

        bytes memory baseData = SSTORE2.read(basePointer);
        ValidatorCacheHeader memory baseHeader = _readValidatorCacheHeader(delta.baseValidatorsHash, baseData);
        require(
            baseHeader.entryCount == _cachedValidatorEntryCount, CachedValidatorSetCorrupted(delta.baseValidatorsHash)
        );

        uint256 changedCount = delta.leafCount;
        require(
            changedCount > 0 && changedCount <= MAX_DELTA_LEAF_COUNT,
            CachedValidatorSetCorrupted(delta.baseValidatorsHash)
        );

        uint256 pubKeyBits = _validatorPubKeyOverrideBits;
        uint256 votingPowerBits = _validatorVotingPowerOverrideBits;
        uint256 newTotalVotingPower = uint256(_cachedValidatorTotalVotingPower);
        bytes32[] memory newLeafHashes = new bytes32[](changedCount);
        uint32 previousIndex = 0;
        for (uint256 i = 0; i < changedCount; i++) {
            uint32 changedIndex = delta.indices[i];
            require(changedIndex < baseHeader.entryCount, SignerIndexOutOfRange(changedIndex));
            if (i > 0) {
                require(changedIndex > previousIndex, CachedValidatorSetCorrupted(delta.baseValidatorsHash));
            }
            previousIndex = changedIndex;

            (uint64 oldVotingPower, bytes32 oldPubKey) = _resolveCurrentCachedValidator(
                delta.baseValidatorsHash, baseData, baseHeader.entryCount, pubKeyBits, votingPowerBits, changedIndex
            );
            uint64 changedVotingPower = delta.votingPowers[i];
            bytes32 changedPubKey = delta.pubKeys[i];
            require(
                oldVotingPower != changedVotingPower || oldPubKey != changedPubKey,
                CachedValidatorSetCorrupted(delta.baseValidatorsHash)
            );
            newTotalVotingPower = newTotalVotingPower - uint256(oldVotingPower) + uint256(changedVotingPower);
            newLeafHashes[i] = Header.simpleValidatorLeafHash(changedPubKey, changedVotingPower);
        }
        require(newTotalVotingPower <= type(uint64).max, CachedValidatorSetCorrupted(delta.baseValidatorsHash));

        bytes32 computedRoot = _recomputeCurrentDeltaRoot(
            delta.baseValidatorsHash, baseData, baseHeader.entryCount, delta, newLeafHashes, pubKeyBits, votingPowerBits
        );
        require(computedRoot == validatorsHash, MismatchedValidatorHashes(validatorsHash, computedRoot));

        for (uint256 i = 0; i < changedCount; i++) {
            uint16 changedIndex = uint16(delta.indices[i]);
            uint256 mask = uint256(1) << changedIndex;
            uint64 changedVotingPower = delta.votingPowers[i];
            bytes32 changedPubKey = delta.pubKeys[i];
            (uint64 baseVotingPower, bytes32 basePubKey) =
                _resolveCachedValidator(delta.baseValidatorsHash, baseData, baseHeader, changedIndex);

            if (changedVotingPower == baseVotingPower) {
                votingPowerBits &= ~mask;
            } else {
                _validatorVotingPowerOverrides[changedIndex] = changedVotingPower;
                votingPowerBits |= mask;
            }

            if (changedPubKey == basePubKey) {
                pubKeyBits &= ~mask;
            } else {
                _validatorPubKeyOverrides[changedIndex] = changedPubKey;
                pubKeyBits |= mask;
            }
        }

        validatorsHashLeaf = Header.bytes32LeafHash(validatorsHash);
        _validatorPubKeyOverrideBits = pubKeyBits;
        _validatorVotingPowerOverrideBits = votingPowerBits;
        _cachedValidatorTotalVotingPower = uint64(newTotalVotingPower);
        _cachedValidatorsHashLeaf = validatorsHashLeaf;
        _cachedValidatorSetHash = validatorsHash;
    }

    function _cacheValidatorSet(
        bytes32 validatorsHash,
        IUpdateClientMsgs.MsgUpdateClient memory msg_,
        uint64 totalVotingPower
    )
        private
    {
        if (_hasUsableValidatorCache(validatorsHash)) {
            return;
        }

        _setValidatorCache(
            validatorsHash, _buildValidatorSetCache(validatorsHash, msg_.proposedHeader.validatorSet, totalVotingPower)
        );
    }

    function _buildValidatorSetCache(
        bytes32 validatorsHash,
        IICS07TendermintMsgs.ValidatorSet memory validatorSet,
        uint64 totalVotingPower
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

        bytes32[] memory leafHashes = new bytes32[](validatorCount);
        for (uint256 i = 0; i < vals.length; i++) {
            leafHashes[i] = Header.simpleValidatorLeafHash(vals[i].pubKey, vals[i].votingPower);
        }
        uint256 nodeCount = validatorCount * 2 - 1;
        bytes32[] memory nodeHashes = new bytes32[](nodeCount);
        _writeMerkleNodes(leafHashes, 0, validatorCount, nodeHashes, 0);

        data = new bytes(
            VALIDATOR_CACHE_HEADER_LEN + validatorCount * VALIDATOR_CACHE_ENTRY_LEN + nodeCount
                * VALIDATOR_CACHE_NODE_LEN
        );
        _writeUint32(data, 0, VALIDATOR_CACHE_MAGIC);
        _writeUint64(data, 4, totalVotingPower);
        _writeUint16(data, 12, validatorCount);
        _writeBytes32(data, 14, Header.bytes32LeafHash(validatorsHash));
        _writeUint16(data, 46, nodeCount);

        uint256 offset = VALIDATOR_CACHE_HEADER_LEN;
        for (uint256 i = 0; i < vals.length; i++) {
            IICS07TendermintMsgs.ValidatorInfo memory val = vals[i];
            _writeUint32(data, offset, uint32(i));
            _writeUint64(data, offset + 4, val.votingPower);
            _writeBytes32(data, offset + 12, val.pubKey);
            offset += VALIDATOR_CACHE_ENTRY_LEN;
        }
        for (uint256 i = 0; i < nodeHashes.length; i++) {
            _writeBytes32(data, offset, nodeHashes[i]);
            offset += VALIDATOR_CACHE_NODE_LEN;
        }
    }

    function _emptyValidatorSet() private pure returns (IICS07TendermintMsgs.ValidatorSet memory validatorSet) {
        validatorSet.validators = new IICS07TendermintMsgs.ValidatorInfo[](0);
    }

    /// @dev Enforces 2/3+ voting power over the active signers in msg_.signerIndices,
    ///      then dispatches to the bucket's Groth16 verifier via the wrapper. Padding
    ///      slots (active=false) carry deterministic dummy data and are skipped here.
    ///      The SharedBlock passed to the verifier is built directly from
    ///      msg_.proposedHeader so the on-chain quorum check and the in-circuit
    ///      reconstruction agree on the signed bytes.
    function _verifyFullBatchAndQuorum(IUpdateClientMsgs.MsgUpdateClient memory msg_)
        internal
        returns (uint64 totalVotingPower)
    {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = msg_.proposedHeader.validatorSet.validators;
        uint256 numVals = vals.length;
        require(
            msg_.signerIndices.length == msg_.bucket && msg_.signerPubkeys.length == msg_.bucket
                && msg_.timestampSeconds.length == msg_.bucket && msg_.timestampNanos.length == msg_.bucket
                && msg_.active.length == msg_.bucket,
            BatchLengthMismatch()
        );

        for (uint256 i = 0; i < numVals; i++) {
            totalVotingPower += vals[i].votingPower;
        }

        uint64 accumulated = 0;
        bool hasPrevSigner = false;
        uint32 prevIdx = 0;
        for (uint256 i = 0; i < msg_.signerIndices.length; i++) {
            if (!msg_.active[i]) {
                continue; // dummy padding slot
            }
            uint32 idx = msg_.signerIndices[i];
            require(idx < numVals, SignerIndexOutOfRange(idx));
            if (hasPrevSigner) {
                require(idx > prevIdx, DuplicateSigner(idx));
            }
            require(vals[idx].pubKey == msg_.signerPubkeys[i], PubkeyMismatch(idx));
            hasPrevSigner = true;
            prevIdx = idx;
            accumulated += vals[idx].votingPower;
        }
        require(
            uint256(accumulated) * 3 > uint256(totalVotingPower) * 2,
            InsufficientVotingPower(accumulated, totalVotingPower)
        );

        _verifyBatchProof(msg_);
    }

    function _verifyCachedBatchAndQuorum(
        bytes32 validatorsHash,
        IUpdateClientMsgs.MsgUpdateClient memory msg_
    )
        internal
        returns (uint64 totalVotingPower)
    {
        require(
            msg_.signerIndices.length == msg_.bucket && msg_.signerPubkeys.length == msg_.bucket
                && msg_.timestampSeconds.length == msg_.bucket && msg_.timestampNanos.length == msg_.bucket
                && msg_.active.length == msg_.bucket,
            BatchLengthMismatch()
        );

        address pointer = _validatorCachePointer(validatorsHash);
        if (pointer == address(0)) {
            revert ValidatorSetCacheMiss(validatorsHash);
        }

        bytes memory cacheData = SSTORE2.read(pointer, 0, _validatorCacheEntryDataLen(_cachedValidatorEntryCount));
        ValidatorCacheHeader memory cacheHeader = _readValidatorCacheHeader(validatorsHash, cacheData);
        totalVotingPower = _cachedValidatorTotalVotingPower;
        uint64 accumulated = 0;
        bool hasPrevSigner = false;
        uint32 prevIdx = 0;
        uint256 pubKeyBits = _validatorPubKeyOverrideBits;
        uint256 votingPowerBits = _validatorVotingPowerOverrideBits;
        for (uint256 i = 0; i < msg_.signerIndices.length; i++) {
            if (!msg_.active[i]) {
                continue; // dummy padding slot
            }

            uint32 idx = msg_.signerIndices[i];
            if (hasPrevSigner) {
                require(idx > prevIdx, DuplicateSigner(idx));
            }
            hasPrevSigner = true;
            prevIdx = idx;

            uint64 votingPower;
            bytes32 pubKey;
            (votingPower, pubKey) = _resolveCurrentCachedValidator(
                validatorsHash, cacheData, cacheHeader.entryCount, pubKeyBits, votingPowerBits, idx
            );
            require(pubKey == msg_.signerPubkeys[i], PubkeyMismatch(idx));
            accumulated += votingPower;
        }

        require(
            uint256(accumulated) * 3 > uint256(totalVotingPower) * 2,
            InsufficientVotingPower(accumulated, totalVotingPower)
        );

        _verifyBatchProof(msg_);
    }

    function _hasUsableValidatorCache(bytes32 validatorsHash) private view returns (bool) {
        return _validatorCachePointer(validatorsHash) != address(0) && _cachedValidatorEntryCount != 0;
    }

    function _getUsableValidatorHashLeaf(bytes32 validatorsHash)
        private
        view
        returns (bool usable, bytes32 validatorsHashLeaf)
    {
        if (_validatorCachePointer(validatorsHash) == address(0)) {
            return (false, bytes32(0));
        }
        return (true, _cachedValidatorsHashLeaf);
    }

    function _validatorCachePointer(bytes32 validatorsHash) private view returns (address) {
        return validatorsHash == _cachedValidatorSetHash ? _cachedValidatorSetPointer : address(0);
    }

    function _setValidatorCache(bytes32 validatorsHash, bytes memory cacheData) private {
        ValidatorCacheHeader memory cacheHeader = _readValidatorCacheHeader(validatorsHash, cacheData);
        _cachedValidatorSetPointer = SSTORE2.write(cacheData);
        _cachedValidatorSetHash = validatorsHash;
        _cachedValidatorsHashLeaf = cacheHeader.validatorsHashLeaf;
        _cachedValidatorTotalVotingPower = cacheHeader.totalVotingPower;
        _cachedValidatorEntryCount = cacheHeader.entryCount;
        _validatorPubKeyOverrideBits = 0;
        _validatorVotingPowerOverrideBits = 0;
    }

    function _verifyBatchProof(IUpdateClientMsgs.MsgUpdateClient memory msg_) private {
        IICS07TendermintMsgs.BlockCommit memory commit = msg_.proposedHeader.signedHeader.commit;
        IVerifier.SharedBlock memory shared = IVerifier.SharedBlock({
            height: commit.height,
            round: uint64(commit.round),
            blockIDHash: commit.blockId.hashData,
            partSetTotal: commit.blockId.partSetHeader.total,
            partSetHash: commit.blockId.partSetHeader.hashData,
            chainID: bytes(msg_.proposedHeader.signedHeader.header.chainId)
        });

        require(
            VERIFIER.verifyBatchProof(
                msg_.bucket,
                msg_.proof,
                msg_.commitments,
                msg_.commitmentPok,
                msg_.signerPubkeys,
                msg_.timestampSeconds,
                msg_.timestampNanos,
                msg_.active,
                shared
            ),
            ProofVerificationFailed()
        );
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
            cacheHeader.entryCount == 0 || cacheHeader.entryCount > MAX_VALIDATOR_COUNT
                || cacheHeader.nodeCount != cacheHeader.entryCount * 2 - 1
                || !_isValidValidatorCacheDataLength(cacheData.length, cacheHeader.entryCount, cacheHeader.nodeCount)
        ) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
    }

    function _validatorCacheEntryDataLen(uint16 entryCount) private pure returns (uint256) {
        return VALIDATOR_CACHE_HEADER_LEN + uint256(entryCount) * VALIDATOR_CACHE_ENTRY_LEN;
    }

    function _isValidValidatorCacheDataLength(
        uint256 dataLength,
        uint16 entryCount,
        uint16 nodeCount
    )
        private
        pure
        returns (bool)
    {
        uint256 entryDataLen = _validatorCacheEntryDataLen(entryCount);
        return dataLength == entryDataLen || dataLength == entryDataLen + uint256(nodeCount) * VALIDATOR_CACHE_NODE_LEN;
    }

    function _resolveCachedValidator(
        bytes32 validatorsHash,
        bytes memory cacheData,
        ValidatorCacheHeader memory cacheHeader,
        uint32 index
    )
        private
        pure
        returns (uint64 votingPower, bytes32 pubKey)
    {
        if (index >= cacheHeader.entryCount) {
            revert CachedSignerNotFound(validatorsHash, index);
        }

        uint256 offset = VALIDATOR_CACHE_HEADER_LEN + uint256(index) * VALIDATOR_CACHE_ENTRY_LEN;
        if (_readUint32(cacheData, offset) != index) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
        votingPower = _readUint64(cacheData, offset + 4);
        pubKey = _readBytes32(cacheData, offset + 12);
    }

    function _resolveCurrentCachedValidator(
        bytes32 validatorsHash,
        bytes memory baseData,
        uint16 entryCount,
        uint256 pubKeyBits,
        uint256 votingPowerBits,
        uint32 index
    )
        private
        view
        returns (uint64 votingPower, bytes32 pubKey)
    {
        if (index >= entryCount) {
            revert CachedSignerNotFound(validatorsHash, index);
        }

        uint16 shortIndex = uint16(index);
        uint256 offset = VALIDATOR_CACHE_HEADER_LEN + uint256(index) * VALIDATOR_CACHE_ENTRY_LEN;
        if (_readUint32(baseData, offset) != index) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
        uint256 mask = uint256(1) << uint256(index);
        votingPower = (votingPowerBits & mask) == 0
            ? _readUint64(baseData, offset + 4)
            : _validatorVotingPowerOverrides[shortIndex];
        pubKey = (pubKeyBits & mask) == 0 ? _readBytes32(baseData, offset + 12) : _validatorPubKeyOverrides[shortIndex];
    }

    function _recomputeCurrentDeltaRoot(
        bytes32 validatorsHash,
        bytes memory baseData,
        uint16 validatorCount,
        IUpdateClientMsgs.ValidatorSetDelta memory delta,
        bytes32[] memory newLeafHashes,
        uint256 pubKeyBits,
        uint256 votingPowerBits
    )
        private
        view
        returns (bytes32 root)
    {
        DeltaRecomputeContext memory ctx = DeltaRecomputeContext({
            validatorsHash: validatorsHash,
            baseData: baseData,
            validatorCount: validatorCount,
            delta: delta,
            newLeafHashes: newLeafHashes,
            pubKeyBits: pubKeyBits,
            votingPowerBits: votingPowerBits,
            overrideBits: pubKeyBits | votingPowerBits
        });
        root = _recomputeCurrentDeltaRange(ctx, 0, validatorCount, 0, 0, delta.leafCount);
    }

    function _recomputeCurrentDeltaRange(
        DeltaRecomputeContext memory ctx,
        uint256 from,
        uint256 to,
        uint256 nodeIndex,
        uint256 changedFrom,
        uint256 changedTo
    )
        private
        view
        returns (bytes32 nodeHash)
    {
        if (changedFrom == changedTo && !_rangeHasOverride(ctx.overrideBits, from, to)) {
            return _getBaseNodeHash(ctx.validatorsHash, ctx.baseData, ctx.validatorCount, nodeIndex);
        }

        if (nodeIndex > type(uint16).max) {
            revert BatchLengthMismatch();
        }

        uint256 length = to - from;
        if (length == 1) {
            if (changedFrom != changedTo) {
                if (changedTo != changedFrom + 1 || from != ctx.delta.indices[changedFrom]) {
                    revert CachedValidatorSetCorrupted(ctx.validatorsHash);
                }
                return ctx.newLeafHashes[changedFrom];
            }

            (uint64 votingPower, bytes32 pubKey) = _resolveCurrentCachedValidator(
                ctx.validatorsHash, ctx.baseData, ctx.validatorCount, ctx.pubKeyBits, ctx.votingPowerBits, uint32(from)
            );
            return Header.simpleValidatorLeafHash(pubKey, votingPower);
        }

        uint256 split = _splitPoint(length);
        uint256 mid = from + split;
        uint256 leftNodeIndex = nodeIndex + 1;
        uint256 rightNodeIndex = nodeIndex + 1 + (split * 2 - 1);
        uint256 rightChangedFrom = changedFrom;
        while (rightChangedFrom < changedTo && ctx.delta.indices[rightChangedFrom] < mid) {
            rightChangedFrom++;
        }

        bytes32 left;
        bytes32 right;
        left = _recomputeCurrentDeltaRange(ctx, from, mid, leftNodeIndex, changedFrom, rightChangedFrom);
        right = _recomputeCurrentDeltaRange(ctx, mid, to, rightNodeIndex, rightChangedFrom, changedTo);

        return Header.innerHash(left, right);
    }

    function _rangeHasOverride(uint256 overrideBits, uint256 from, uint256 to) private pure returns (bool) {
        uint256 width = to - from;
        uint256 mask = ((uint256(1) << width) - 1) << from;
        return (overrideBits & mask) != 0;
    }

    function _getBaseNodeHash(
        bytes32 validatorsHash,
        bytes memory baseData,
        uint16 entryCount,
        uint256 nodeIndex
    )
        private
        pure
        returns (bytes32)
    {
        uint256 offset = VALIDATOR_CACHE_HEADER_LEN + uint256(entryCount) * VALIDATOR_CACHE_ENTRY_LEN + nodeIndex
            * VALIDATOR_CACHE_NODE_LEN;
        if (offset + VALIDATOR_CACHE_NODE_LEN > baseData.length) {
            revert CachedValidatorSetCorrupted(validatorsHash);
        }
        return _readBytes32(baseData, offset);
    }

    function _writeMerkleNodes(
        bytes32[] memory leafHashes,
        uint256 from,
        uint256 to,
        bytes32[] memory nodeHashes,
        uint256 nodeIndex
    )
        private
        pure
        returns (bytes32 nodeHash)
    {
        uint256 length = to - from;
        if (length == 1) {
            nodeHash = leafHashes[from];
            nodeHashes[nodeIndex] = nodeHash;
            return nodeHash;
        }

        uint256 split = _splitPoint(length);
        bytes32 left = _writeMerkleNodes(leafHashes, from, from + split, nodeHashes, nodeIndex + 1);
        bytes32 right = _writeMerkleNodes(leafHashes, from + split, to, nodeHashes, nodeIndex + 1 + (split * 2 - 1));
        nodeHash = Header.innerHash(left, right);
        nodeHashes[nodeIndex] = nodeHash;
    }

    function _splitPoint(uint256 n) private pure returns (uint256 split) {
        split = 1;
        while ((split << 1) < n) {
            split <<= 1;
        }
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
    /// @dev If the proof is empty, then we assume that the proof was cached earlier in the same tx.
    /// @dev The proof is cached in the transient storage.
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

    /// @dev The misbehavior is verfied in the gnark program. Here we only check the public values which contain the
    /// trusted headers.
    /// @inheritdoc ILightClient
    function misbehaviour(bytes calldata) external view notFrozen onlyMisbehaviourSubmitter {
        revert FeatureNotSupported();
    }

    /// @inheritdoc IGroth16ICS07Tendermint
    function unfreeze() external override(IGroth16ICS07Tendermint, ILightClient) onlyRole(DEFAULT_ADMIN_ROLE) {
        require(clientState.isFrozen, ClientNotFrozen());
        clientState.isFrozen = false;
    }

    /// @inheritdoc ILightClient
    function upgradeClient(bytes calldata) external pure {
        // NOTE: This feature will not be supported. (#130)
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

        // We avoid the cost of caching for single kv pairs, as reusing the proof is not necessary
        if (kvPairs.length > 1) {
            _cacheKvPairs(height.revisionHeight, kvPairs, trustedConsensusState.timestamp);
        }
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
            block.timestamp - _nanosToSeconds(time) <= ALLOWED_CLOCK_DRIFT,
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

    /// @notice Caches the key-value pairs to the transient storage with the timestamp.
    /// @param proofHeight The height of the proof.
    /// @param kvPairs The key-value pairs.
    /// @param timestamp The timestamp of the trusted consensus state in unix nanoseconds.
    /// @dev WARNING: Transient store is not reverted even if a message within a transaction reverts.
    /// @dev WARNING: This function must be called after all proof and validation checks.
    function _cacheKvPairs(uint64 proofHeight, IMembershipMsgs.KVPair[] memory kvPairs, uint256 timestamp) private {
        for (uint256 i = 0; i < kvPairs.length; i++) {
            bytes32 kvPairHash = keccak256(abi.encode(proofHeight, kvPairs[i]));
            kvPairHash.asUint256().tstore(timestamp);
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
