// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IUpdateClientMsgs } from "./msgs/IUpdateClientMsgs.sol";
import { IICS07TendermintMsgs } from "./msgs/IICS07TendermintMsgs.sol";

/// @title IGroth16ICS07Tendermint
/// @notice IGroth16ICS07Tendermint is the interface for the ICS07 Tendermint light client
interface IGroth16ICS07Tendermint {
    /// @notice The role identifier for the misbehaviour submitter role
    /// @dev The misbehaviour submitter role is used to whitelist addresses that can submit misbehaviour reports
    /// @dev If `address(0)` has this role, then anyone can submit misbehaviour reports
    /// @dev This role is distinct from PROOF_SUBMITTER_ROLE to allow separate access control
    /// @return The role identifier
    function MISBEHAVIOUR_SUBMITTER_ROLE() external view returns (bytes32);

    /// @notice Unfreezes a previously frozen client.
    /// @dev Can only be called by an address with the `DEFAULT_ADMIN_ROLE`.
    function unfreeze() external;

    // /// @notice Immutable update client program verification key.
    // /// @return The verification key for the update client program.
    // function UPDATE_CLIENT_PROGRAM_VKEY() external view returns (bytes32);

    // /// @notice Immutable membership program verification key.
    // /// @return The verification key for the membership program.
    // function MEMBERSHIP_PROGRAM_VKEY() external view returns (bytes32);

    // /// @notice Immutable update client and membership program verification key.
    // /// @return The verification key for the update client and membership program.
    // function UPDATE_CLIENT_AND_MEMBERSHIP_PROGRAM_VKEY() external view returns (bytes32);

    // /// @notice Immutable misbehaviour program verification key.
    // /// @return The verification key for the misbehaviour program.
    // function MISBEHAVIOUR_PROGRAM_VKEY() external view returns (bytes32);

    /// @notice Returns pinned validator metadata.
    /// @return indices      Validator indices in the pinned set, sorted ascending.
    /// @return pubkeys      Ed25519 public keys for each pinned validator.
    /// @return votingPowers Voting power for each pinned validator.
    function getPinnedValidatorSet()
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers);

    /// @notice Replaces the pinned validator set after verifying a checkpoint commit.
    function reAnchorPinnedSet(
        IUpdateClientMsgs.MsgUpdateClient calldata updateMsg,
        IICS07TendermintMsgs.ValidatorSet calldata newPinnedValidatorSet
    ) external;
}
