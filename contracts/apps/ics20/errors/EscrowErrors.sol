// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/// @title EscrowErrors
/// @notice Interface for escrow-related errors
interface EscrowErrors {
    /// @notice Unauthorized function call
    /// @param caller The caller of the function
    error EscrowUnauthorized(address caller);
}
