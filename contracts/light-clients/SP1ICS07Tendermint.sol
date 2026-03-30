// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable gas-strict-inequalities

import { IICS07TendermintMsgs } from "./msgs/IICS07TendermintMsgs.sol";
import { IUpdateClientMsgs } from "./msgs/IUpdateClientMsgs.sol";
import { IMembershipMsgs } from "./msgs/IMembershipMsgs.sol";
import { IUpdateClientAndMembershipMsgs } from "./msgs/IUcAndMembershipMsgs.sol";
import { IMisbehaviourMsgs } from "./msgs/IMisbehaviourMsgs.sol";
import { ISP1Msgs } from "./msgs/ISP1Msgs.sol";
import { ILightClientMsgs } from "../msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../msgs/IICS02ClientMsgs.sol";

import { ISP1ICS07TendermintErrors } from "./errors/ISP1ICS07TendermintErrors.sol";
import { ISP1ICS07Tendermint } from "./ISP1ICS07Tendermint.sol";
import { IMembership } from "../interfaces/IMembership.sol";
import { IMisbehaviour } from "../interfaces/IMisbehaviour.sol";
import { IUpdateClient } from "../interfaces/IUpdateClient.sol";
import { ILightClient } from "../interfaces/ILightClient.sol";
import { IVerifier } from "../interfaces/IVerifier.sol";

import { Paths } from "./utils/Paths.sol";
import { Encode } from "../utils/Encode.sol";
import { Multicall } from "@openzeppelin-contracts/utils/Multicall.sol";
import { TransientSlot } from "@openzeppelin-contracts/utils/TransientSlot.sol";
import { AccessControl } from "@openzeppelin-contracts/access/AccessControl.sol";

/// @title SP1 ICS07 Tendermint Light Client
/// @author srdtrk
/// @notice This contract implements an ICS07 IBC tendermint light client using SP1.
contract SP1ICS07Tendermint is
    ISP1ICS07TendermintErrors,
    ISP1ICS07Tendermint,
    ILightClient,
    Multicall,
    AccessControl
{
    using TransientSlot for *;

    /// @inheritdoc ISP1ICS07Tendermint
    IVerifier public immutable VERIFIER;
    IMembership public immutable MEMBERSHIP;
    IMisbehaviour public immutable MISBEHAVIOUR;
    IUpdateClient public immutable UPDATE_CLIENT;

    /// @notice The ICS07Tendermint client state
    // natlint-disable-next-line MissingInheritdoc
    IICS07TendermintMsgs.ClientState public clientState;
    /// @notice The mapping from height to consensus state keccak256 hashes.
    /// @dev Revision number need not be keyed as it is not allowed to change.
    mapping(uint64 height => bytes32 hash) private _consensusStateHashes;

    /// @inheritdoc ISP1ICS07Tendermint
    uint16 public constant ALLOWED_SP1_CLOCK_DRIFT = 30 minutes;

    /// @inheritdoc ISP1ICS07Tendermint
    bytes32 public constant PROOF_SUBMITTER_ROLE = keccak256("PROOF_SUBMITTER_ROLE");

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
        _consensusStateHashes[clientState.latestHeight.revisionHeight] = _consensusState;

        VERIFIER = IVerifier(verifier);
        MEMBERSHIP = IMembership(membership_);
        MISBEHAVIOUR = IMisbehaviour(misbehaviour_);
        UPDATE_CLIENT = IUpdateClient(updateClient_);

        require(
            clientState.trustingPeriod + ALLOWED_SP1_CLOCK_DRIFT <= clientState.unbondingPeriod,
            TrustingPeriodTooLong(clientState.trustingPeriod, clientState.unbondingPeriod)
        );

        if (roleManager == address(0)) {
            _grantRole(PROOF_SUBMITTER_ROLE, address(0)); // Allow anyone to submit proofs
        } else {
            _grantRole(DEFAULT_ADMIN_ROLE, roleManager); // Allow the role manager to manage roles
            _grantRole(PROOF_SUBMITTER_ROLE, roleManager); // Allow the role manager to submit proofs
        }
    }

    /// @inheritdoc ILightClient
    function getClientState() external view returns (bytes memory) {
        return abi.encode(clientState);
    }

    /// @inheritdoc ISP1ICS07Tendermint
    function getConsensusStateHash(uint64 revisionHeight) public view returns (bytes32) {
        bytes32 hash = _consensusStateHashes[revisionHeight];
        require(hash != 0, ConsensusStateNotFound());
        return hash;
    }

    /// @dev This function verifies the public values and forwards the proof to the SP1 verifier.
    /// @inheritdoc ILightClient
    function updateClient(
        bytes calldata updateClientMsg
    )
        external
        notFrozen
        onlyProofSubmitter
        returns (ILightClientMsgs.UpdateResult)
    {
        IUpdateClientMsgs.MsgUpdateClient memory msg_ = abi.decode(updateClientMsg, (IUpdateClientMsgs.MsgUpdateClient));
        IUpdateClientMsgs.UpdateClientOutput memory output =
            UPDATE_CLIENT.updateClient(
                msg_
            );

        _validateUpdateClientOutput(output);

        ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult(output);
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

        // TODO: multi signatures
	    string memory chainId = msg_.proposedHeader.signedHeader.header.chainId;
        IICS07TendermintMsgs.BlockCommit memory untrustedHeaderCommit = msg_.proposedHeader.signedHeader.commit;
        bytes32 pubkey = msg_.proposedHeader.validatorSet.validators[0].pubKey;
        bytes memory sig = untrustedHeaderCommit.commitSigs[0].data.signature;
        require(sig.length == 64, "invalid signature length");
        bytes32[2] memory signature;
        assembly {
            mstore(signature, mload(add(sig, 32)))
            mstore(add(signature, 32), mload(add(sig, 64)))
        }
        bool proofValid = VERIFIER.verifyProof(
            msg_.proof,
            msg_.commitments,
            msg_.commitmentPok,
            signature,
            pubkey,
            Encode.voteSignBytes(untrustedHeaderCommit, chainId, 0)
        );
        require(proofValid, ProofVerificationFailed());

        return updateResult;
    }

    /// @inheritdoc ILightClient
    function verifyMembership(ILightClientMsgs.MsgVerifyMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        require(msg_.appHash.length > 0, EmptyValue());
        return _membership(msg_.height, msg_.kvPairs, msg_.merkleProofs, msg_.appHash, msg_.trustedConsensusState, msg_.membershipType);
    }

    /// @inheritdoc ILightClient
    function verifyNonMembership(ILightClientMsgs.MsgVerifyNonMembership calldata msg_)
        external
        notFrozen
        onlyProofSubmitter
        returns (uint256)
    {
        return _membership(msg_.height, msg_.kvPairs, msg_.merkleProofs, msg_.appHash, msg_.trustedConsensusState, msg_.membershipType);
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
        IMembershipMsgs.MembershipType membershipType
    )
        private
        returns (uint256)
    {
        // TODO: cached proof
        // if (merkleProofs.length == 0) {
        //     // cached proof
        //     return _nanosToSeconds(_getCachedKvPair(height.revisionHeight, IMembershipMsgs.KVPair(path, value)));
        // }

        if (membershipType == IMembershipMsgs.MembershipType.Membership) {
            return _handleMembership(height, kvPairs, merkleProofs, appHash, trustedConsensusState);
        } else if (membershipType == IMembershipMsgs.MembershipType.MembershipAndUpdateClient) {
            // return _handleSP1UpdateClientAndMembership(height, membershipProof.proof, path, value);
        }

        // unreachable
        revert UnknownMembershipType(uint8(membershipType));
    }

    /// @dev The misbehavior is verfied in the sp1 program. Here we only check the public values which contain the
    /// trusted headers.
    /// @inheritdoc ILightClient
    function misbehaviour(
        bytes calldata misbehaviourMsg
    ) external notFrozen onlyProofSubmitter {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = abi.decode(misbehaviourMsg, (IMisbehaviourMsgs.MsgSubmitMisbehaviour));
        IMisbehaviourMsgs.MisbehaviourOutput memory output = MISBEHAVIOUR.misbehaviour(
            msg_.clientState,
            msg_.misbehaviour,
            msg_.trustedConsensusState1,
            msg_.trustedConsensusState2,
            msg_.time
        );

        _validateMisbehaviourOutput(
            output,
            msg_.clientState,
            msg_.trustedConsensusState1,
            msg_.trustedConsensusState2,
            msg_.time
        );

        // _verifyProof(msgSubmitMisbehaviour.sp1Proof);

        // If the misbehaviour and proof is valid, the client needs to be frozen
        clientState.isFrozen = true;
    }

    /// @inheritdoc ILightClient
    function upgradeClient(bytes calldata) external view notFrozen onlyProofSubmitter {
        // NOTE: This feature will not be supported. (#130)
        revert FeatureNotSupported();
    }

    /// @notice Handles the `SP1MembershipProof` proof type.
    /// @param height The height of the proof.
    /// @param kvPairs The path and value of the key-value pair.
    /// @param merkleProofs The merkle proofs of membership.
    /// @param appHash The final hash value that needs to verify.
    /// @return The timestamp of the trusted consensus state.
    function _handleMembership(
        IICS02ClientMsgs.Height calldata height,
        IMembershipMsgs.KVPair[] calldata kvPairs,
        IMembershipMsgs.MerkleProof[] calldata merkleProofs,
        bytes32 appHash,
        IICS07TendermintMsgs.ConsensusState calldata trustedConsensusState
    )
        private
        returns (uint256)
    {
        require(
            kvPairs.length > 0 && kvPairs.length <= type(uint16).max,
            LengthIsOutOfRange(kvPairs.length, 1, type(uint16).max)
        );

        IMembershipMsgs.MembershipOutput memory output =
            MEMBERSHIP.membership(appHash, kvPairs, merkleProofs);
        _validateMembershipOutput(output.commitmentRoot, height.revisionHeight, trustedConsensusState);

        // We avoid the cost of caching for single kv pairs, as reusing the proof is not necessary
        if (output.kvPairs.length > 1) {
            _cacheKvPairs(height.revisionHeight, output.kvPairs, trustedConsensusState.timestamp);
        }
        return _getTimestampInSeconds(trustedConsensusState);
    }

    /// @notice The entrypoint for handling the `SP1MembershipAndUpdateClientProof` proof type.
    /// @dev This function verifies the public values and forwards the proof to the SP1 verifier.
    /// @param proofHeight The height of the proof.
    /// @param proofBytes The encoded proof.
    /// @param kvPath The path of the key-value pair.
    /// @param kvValue The value of the key-value pair.
    /// @return The timestamp of the new consensus state.
    // solhint-disable-next-line code-complexity,function-max-lines
    function _handleSP1UpdateClientAndMembership(
        IICS02ClientMsgs.Height calldata proofHeight,
        bytes memory proofBytes,
        bytes[] calldata kvPath,
        bytes32 kvValue
    )
        private
        returns (uint256)
    {
        // validate proof and deserialize output
        IUpdateClientAndMembershipMsgs.UcAndMembershipOutput memory output;
        {
            IMembershipMsgs.SP1MembershipAndUpdateClientProof memory proof =
                abi.decode(proofBytes, (IMembershipMsgs.SP1MembershipAndUpdateClientProof));
            // require(
            //     proof.sp1Proof.vKey == UPDATE_CLIENT_AND_MEMBERSHIP_PROGRAM_VKEY,
            //     VerificationKeyMismatch(UPDATE_CLIENT_AND_MEMBERSHIP_PROGRAM_VKEY, proof.sp1Proof.vKey)
            // );

            output = abi.decode(proof.sp1Proof.publicValues, (IUpdateClientAndMembershipMsgs.UcAndMembershipOutput));
            require(
                output.kvPairs.length > 0 && output.kvPairs.length <= type(uint16).max,
                LengthIsOutOfRange(output.kvPairs.length, 1, type(uint16).max)
            );

            require(
                proofHeight.revisionHeight == output.updateClientOutput.newHeight.revisionHeight
                    && proofHeight.revisionNumber == output.updateClientOutput.newHeight.revisionNumber,
                ProofHeightMismatch(
                    proofHeight.revisionNumber,
                    proofHeight.revisionHeight,
                    output.updateClientOutput.newHeight.revisionNumber,
                    output.updateClientOutput.newHeight.revisionHeight
                )
            );

            _validateUpdateClientOutput(output.updateClientOutput);

            // TODO: verify proof with input
            // _verifyProof(proof.sp1Proof);
        }

        // check update result
        {
            ILightClientMsgs.UpdateResult updateResult = _checkUpdateResult(output.updateClientOutput);
            if (updateResult == ILightClientMsgs.UpdateResult.Update) {
                // adding the new consensus state to the mapping
                if (proofHeight.revisionHeight > clientState.latestHeight.revisionHeight) {
                    clientState.latestHeight = proofHeight;
                }
                _consensusStateHashes[proofHeight.revisionHeight] =
                    keccak256(abi.encode(output.updateClientOutput.newConsensusState));
            } else if (updateResult == ILightClientMsgs.UpdateResult.Misbehaviour) {
                revert CannotHandleMisbehavior();
            } // else: NoOp
        }

        // loop through the key-value pairs and validate them
        {
            bool found = false;
            for (uint256 i = 0; i < output.kvPairs.length; i++) {
                if (!Paths.equal(output.kvPairs[i].path, kvPath)) {
                    continue;
                }

                bytes32 value = output.kvPairs[i].value;
                require(
                    value == kvValue,
                    MembershipProofValueMismatch(kvValue, value)
                );

                found = true;
                break;
            }
            require(found, MembershipProofKeyNotFound(kvPath));
        }

        _validateMembershipOutput(
            output.updateClientOutput.newConsensusState.root,
            output.updateClientOutput.newHeight.revisionHeight,
            output.updateClientOutput.newConsensusState
        );

        // We avoid the cost of caching for single kv pairs, as reusing the proof is not necessary
        if (output.kvPairs.length > 1) {
            _cacheKvPairs(
                proofHeight.revisionHeight, output.kvPairs, output.updateClientOutput.newConsensusState.timestamp
            );
        }
        return _getTimestampInSeconds(output.updateClientOutput.newConsensusState);
    }

    /// @notice Validates the MembershipOutput public values.
    /// @param outputCommitmentRoot The commitment root of the output.
    /// @param proofHeight The height of the proof.
    /// @param trustedConsensusState The trusted consensus state
    function _validateMembershipOutput(
        bytes32 outputCommitmentRoot,
        uint64 proofHeight,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState
    )
        private
        view
    {
        bytes32 trustedConsensusStateHash = keccak256(abi.encode(trustedConsensusState));
        bytes32 storedConsensusStateHash = getConsensusStateHash(proofHeight);
        require(
            trustedConsensusStateHash == storedConsensusStateHash,
            ConsensusStateHashMismatch(storedConsensusStateHash, trustedConsensusStateHash)
        );

        require(
            outputCommitmentRoot == trustedConsensusState.root,
            ConsensusStateRootMismatch(trustedConsensusState.root, outputCommitmentRoot)
        );
    }

    /// @notice Validates the SP1ICS07UpdateClientOutput public values.
    /// @param output The public values.
    function _validateUpdateClientOutput(IUpdateClientMsgs.UpdateClientOutput memory output) private view {
        _validateClientStateAndTime(output.clientState, output.time);

        bytes32 outputConsensusStateHash = keccak256(abi.encode(output.trustedConsensusState));
        bytes32 storedConsensusStateHash = getConsensusStateHash(output.trustedHeight.revisionHeight);
        require(
            outputConsensusStateHash == storedConsensusStateHash,
            ConsensusStateHashMismatch(storedConsensusStateHash, outputConsensusStateHash)
        );
    }

    /// @notice Validates the SP1ICS07MisbehaviourOutput public values.
    /// @param output The public values.
    function _validateMisbehaviourOutput(
        IMisbehaviourMsgs.MisbehaviourOutput memory output,
        IICS07TendermintMsgs.ClientState memory clientState_,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState1,
        IICS07TendermintMsgs.ConsensusState memory trustedConsensusState2,
        uint128 time
    ) private view {
        _validateClientStateAndTime(clientState_, time);

        // make sure the trusted consensus state from header 1 is known (trusted) by matching it with the one in the
        // mapping
        bytes32 outputConsensusStateHash1 = keccak256(abi.encode(trustedConsensusState1));
        bytes32 storedConsensusStateHash1 = getConsensusStateHash(output.trustedHeight1.revisionHeight);
        require(
            outputConsensusStateHash1 == storedConsensusStateHash1,
            ConsensusStateHashMismatch(storedConsensusStateHash1, outputConsensusStateHash1)
        );

        // make sure the trusted consensus state from header 2 is known (trusted) by matching it with the one in the
        // mapping
        bytes32 outputConsensusStateHash2 = keccak256(abi.encode(trustedConsensusState2));
        bytes32 storedConsensusStateHash2 = getConsensusStateHash(output.trustedHeight2.revisionHeight);
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
            block.timestamp - _nanosToSeconds(time) <= ALLOWED_SP1_CLOCK_DRIFT,
            ProofIsTooOld(block.timestamp, _nanosToSeconds(time))
        );

        // Check client state equality
        // NOTE: We do not check the equality of latest height and isFrozen, this is because:
        // 1. Latest height can be updated by a frontrunner relayer in order to DOS the proof of another relayer.
        // 2. Each external call has the `notFrozen` modifier which checks if the client is frozen.
        // 3. The revision number is not allowed to change with us checking the chain-id and the implementation in the
        // sp1 program.
        require(
            bytes(publicClientState.chainId).length == bytes(clientState.chainId).length
                && keccak256(bytes(publicClientState.chainId)) == keccak256(bytes(clientState.chainId)),
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

    /// @notice Verifies the SP1 proof
    /// @param proof The SP1 proof.
    /// @dev WARNING: proof.vKey must be verified before calling this function.
    // TODO: verify proof with given input
    // function _verifyProof(ISP1Msgs.SP1Proof memory proof) private view {
    //     VERIFIER.verifyProof(proof.vKey, proof.publicValues, proof.proof);
    // }

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

    /// @notice Gets the timestamp of the cached key-value pair from the transient storage.
    /// @param proofHeight The height of the proof.
    /// @param kvPair The key-value pair.
    /// @return The timestamp of the cached key-value pair in unix nanoseconds.
    function _getCachedKvPair(
        uint64 proofHeight,
        IMembershipMsgs.KVPair memory kvPair
    )
        private
        view
        returns (uint256)
    {
        bytes32 kvPairHash = keccak256(abi.encode(proofHeight, kvPair));
        uint256 timestamp = kvPairHash.asUint256().tload();
        require(timestamp != 0, KeyValuePairNotInCache(kvPair.path, kvPair.value));
        return timestamp;
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
}
