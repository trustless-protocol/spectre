// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IICS02ClientMsgs } from "../../msgs/IICS02ClientMsgs.sol";
import { IICS02Client } from "../../interfaces/IICS02Client.sol";
import { IICS02ClientErrors } from "../../errors/IICS02ClientErrors.sol";
import { IAccessManaged } from "@openzeppelin-contracts/access/manager/IAccessManaged.sol";
import { IAccessManager } from "@openzeppelin-contracts/access/manager/IAccessManager.sol";
import { IBCRolesLib } from "../../utils/IBCRolesLib.sol";
import { ICS02ClientStore } from "../../utils/ICS02ClientStore.sol";
import { IClientMigrationProposer } from "../interfaces/IClientMigrationProposer.sol";
import { IClientMigrationModule, ClientMigrationModuleIds } from "../interfaces/IClientMigrationModule.sol";

/// @title Client Migration Proposer
/// @notice Commits a delayed client migration in the router's storage context via delegatecall.
contract ClientMigrationProposer is IClientMigrationProposer, IICS02ClientErrors {
    uint48 internal constant MIN_CLIENT_MIGRATION_DELAY = 48 hours;

    address private immutable SELF;

    constructor() {
        SELF = address(this);
    }

    modifier onlyDelegated() {
        require(address(this) != SELF, IBCClientMigrationDirectCallNotAllowed());
        _;
    }

    /// @inheritdoc IClientMigrationModule
    function moduleId() external pure returns (bytes32) {
        return ClientMigrationModuleIds.PROPOSER;
    }

    /// @inheritdoc IClientMigrationProposer
    function proposeClientMigration(
        string calldata clientId,
        IICS02ClientMsgs.CounterpartyInfo calldata counterpartyInfo,
        address client
    )
        external
        onlyDelegated
    {
        ICS02ClientStore.Layout storage $ = ICS02ClientStore.load();
        require(address($.clients[clientId]) != address(0), IBCClientNotFound(clientId));

        IAccessManager manager = IAccessManager(IAccessManaged(address(this)).authority());
        require(!manager.isTargetClosed(address(this)), IBCUnauthorizedMigrator(clientId, msg.sender));

        (bool isMember, uint32 executionDelay) =
            manager.hasRole(IBCRolesLib.getLightClientMigratorRole(clientId), msg.sender);
        require(isMember, IBCUnauthorizedMigrator(clientId, msg.sender));
        require(executionDelay != 0, IBCClientMigrationDelayRequired(clientId));

        ICS02ClientStore.ClientMigration storage migration = $.migrations[clientId];
        bool active = migration.digest != bytes32(0) && block.timestamp <= migration.expireAfter;
        require(!active, IBCClientMigrationAlreadyProposed(clientId));

        uint256 delay = executionDelay > MIN_CLIENT_MIGRATION_DELAY ? executionDelay : MIN_CLIENT_MIGRATION_DELAY;
        uint48 executeAfter = uint48(block.timestamp + delay);
        migration.digest = _migrationDigest(clientId, counterpartyInfo, client);
        migration.executeAfter = executeAfter;
        migration.expireAfter = uint48(uint256(executeAfter) + manager.expiration());
        migration.proposer = msg.sender;

        emit IICS02Client.ICS02ClientMigrationProposed(
            clientId, migration.digest, executeAfter, migration.expireAfter, migration.proposer
        );
    }

    function _migrationDigest(
        string calldata clientId,
        IICS02ClientMsgs.CounterpartyInfo calldata counterpartyInfo,
        address client
    )
        private
        pure
        returns (bytes32)
    {
        return keccak256(abi.encode(clientId, counterpartyInfo.clientId, counterpartyInfo.merklePrefix, client));
    }
}
