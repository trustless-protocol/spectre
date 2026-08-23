// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/// @notice Canonical identities for migration modules accepted by the router.
library ClientMigrationModuleIds {
    bytes32 internal constant PROPOSER = keccak256("fast-ibc.client-migration.proposer.v1");
    bytes32 internal constant EXECUTOR = keccak256("fast-ibc.client-migration.executor.v1");
}

/// @notice Runtime identity exposed by a migration delegatecall module.
interface IClientMigrationModule {
    function moduleId() external pure returns (bytes32);
}
