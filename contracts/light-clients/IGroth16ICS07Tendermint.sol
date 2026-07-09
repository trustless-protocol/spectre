// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/// @title IGroth16ICS07Tendermint
/// @notice IGroth16ICS07Tendermint is the interface for the ICS07 Tendermint light client
interface IGroth16ICS07Tendermint {
    /// @notice Emitted when the consensus state advances.
    /// @param height The new latest revision height.
    event ClientUpdated(uint64 height);

    /// @notice Emitted when the pinned validator set is re-anchored.
    /// @param height The height whose `nextValidatorsHash` the new set was verified against.
    /// @param validatorsHash The hash of the newly pinned validator set.
    event PinnedSetReAnchored(uint64 height, bytes32 validatorsHash);

    /// @notice Emitted when the client freezes, either by detecting conflicting
    /// state during an update or via a submitted misbehaviour proof.
    event ClientFrozen();

    /// @notice Emitted when the client is unfrozen by the admin.
    event ClientUnfrozen();

    /// @notice The role identifier for the misbehaviour submitter role
    /// @dev The misbehaviour submitter role is used to whitelist addresses that can submit misbehaviour reports
    /// @dev If `address(0)` has this role, then anyone can submit misbehaviour reports
    /// @dev This role is distinct from PROOF_SUBMITTER_ROLE to allow separate access control
    /// @return The role identifier
    function MISBEHAVIOUR_SUBMITTER_ROLE() external view returns (bytes32);

    /// @notice Unfreezes a previously frozen client.
    /// @dev Can only be called by an address with the `DEFAULT_ADMIN_ROLE`.
    function unfreeze() external;

    /// @notice Returns pinned validator metadata.
    /// @return indices      Validator indices in the pinned set, sorted ascending.
    /// @return pubkeys      Ed25519 public keys for each pinned validator.
    /// @return votingPowers Voting power for each pinned validator.
    function getPinnedValidatorSet()
        external
        view
        returns (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers);

    /// @notice Replaces the pinned validator set after verifying a checkpoint commit.
    /// @param reAnchorMsg The ABI-encoded (MsgUpdateClient, ValidatorSet) re-anchor message.
    function reAnchorPinnedSet(bytes calldata reAnchorMsg) external;
}
