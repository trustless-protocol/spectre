// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Bytes } from "@openzeppelin-contracts/utils/Bytes.sol";

/// @title Byte Prefix
/// @notice Domain-neutral prefix comparison shared by protocol packages.
library BytesPrefix {
    /// @notice Checks whether `value` begins with `prefix`.
    function hasPrefix(bytes memory value, bytes memory prefix) internal pure returns (bool) {
        if (value.length < prefix.length) {
            return false;
        }
        return keccak256(Bytes.slice(value, 0, prefix.length)) == keccak256(prefix);
    }
}
