// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/// @title IGroth16ICS07Tendermint
/// @notice IGroth16ICS07Tendermint is the interface for the ICS07 Tendermint light client
interface IGroth16ICS07Tendermint {
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

    /// @notice Returns cached validator metadata for the given validators hash.
    /// @param validatorsHash The CometBFT validators hash.
    /// @return indices      Validator indices in the original set, sorted ascending.
    /// @return pubkeys      Ed25519 public keys for each cached validator.
    /// @return votingPowers Voting power for each cached validator.
    function getCachedValidatorSet(bytes32 validatorsHash)
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers);
}
