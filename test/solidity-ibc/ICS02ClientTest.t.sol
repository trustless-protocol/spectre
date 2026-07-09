// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable custom-errors,max-line-length,gas-small-strings

import { Test } from "forge-std/Test.sol";

import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";

import { IICS02Client } from "../../contracts/interfaces/IICS02Client.sol";
import { ILightClient } from "../../contracts/interfaces/ILightClient.sol";
import { IGroth16ICS07Tendermint } from "../../contracts/light-clients/IGroth16ICS07Tendermint.sol";
import { IAccessManaged } from "@openzeppelin-contracts/access/manager/IAccessManaged.sol";
import { IICS02ClientErrors } from "../../contracts/errors/IICS02ClientErrors.sol";

import { ICS02ClientUpgradeable } from "../../contracts/utils/ICS02ClientUpgradeable.sol";
import { ERC1967Proxy } from "@openzeppelin-contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { ICS26Router } from "../../contracts/ICS26Router.sol";
import { TestHelper } from "./utils/TestHelper.sol";
import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { IBCRolesLib } from "../../contracts/utils/IBCRolesLib.sol";

contract ICS02ClientTest is Test {
    ICS02ClientUpgradeable public ics02Client;
    AccessManager public accessManager;

    address public lightClient = makeAddr("lightClient");

    bytes[] public merklePrefix = [bytes("ibc"), bytes("")];
    bytes[] public randomPrefix = [bytes("test"), bytes("prefix")];

    string public clientIdentifier;

    address public idCustomizer = makeAddr("idCustomizer");
    address public relayer = makeAddr("relayer");
    address public misbehaviourSubmitter = makeAddr("misbehaviourSubmitter");

    TestHelper public th = new TestHelper();

    function setUp() public {
        ICS26Router ics26RouterLogic = new ICS26Router();

        accessManager = new AccessManager(address(this));

        ERC1967Proxy routerProxy = new ERC1967Proxy(
            address(ics26RouterLogic), abi.encodeCall(ICS26Router.initialize, (address(accessManager)))
        );

        ics02Client = ICS02ClientUpgradeable(address(routerProxy));

        accessManager.setTargetFunctionRole(
            address(ics02Client), IBCRolesLib.ics26IdCustomizerSelectors(), IBCRolesLib.ID_CUSTOMIZER_ROLE
        );
        accessManager.setTargetFunctionRole(
            address(ics02Client), IBCRolesLib.ics26RelayerSelectors(), IBCRolesLib.RELAYER_ROLE
        );
        accessManager.setTargetFunctionRole(
            address(ics02Client), IBCRolesLib.ics26MisbehaviourSelectors(), IBCRolesLib.MISBEHAVIOUR_SUBMITTER_ROLE
        );

        accessManager.grantRole(IBCRolesLib.ID_CUSTOMIZER_ROLE, idCustomizer, 0);
        accessManager.grantRole(IBCRolesLib.RELAYER_ROLE, relayer, 0);
        accessManager.grantRole(IBCRolesLib.MISBEHAVIOUR_SUBMITTER_ROLE, misbehaviourSubmitter, 0);

        string memory counterpartyId = "42-dummy-01";
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(counterpartyId, merklePrefix);
        vm.expectEmit();
        emit IICS02Client.ICS02ClientAdded(th.FIRST_CLIENT_ID(), counterpartyInfo, lightClient);
        clientIdentifier = ics02Client.addClient(counterpartyInfo, lightClient);

        ILightClient fetchedLightClient = ics02Client.getClient(clientIdentifier);
        assertNotEq(address(fetchedLightClient), address(0), "client not found");

        IICS02ClientMsgs.CounterpartyInfo memory fetchedCounterparty = ics02Client.getCounterparty(clientIdentifier);
        assertEq(fetchedCounterparty.clientId, counterpartyId, "counterparty not set correctly");
    }

    function test_success_customClientId() public {
        string memory customClientId = "custom-client-id";
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(customClientId, merklePrefix);
        vm.prank(idCustomizer);
        string memory newId = ics02Client.addClient(customClientId, counterpartyInfo, lightClient);
        assertEq(customClientId, newId, "custom client id not set correctly");
    }

    function test_failure_customClientId() public {
        vm.startPrank(idCustomizer);
        // client id is not custom (starts with "client-")
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(clientIdentifier, merklePrefix);
        vm.expectRevert(abi.encodeWithSelector(IICS02ClientErrors.IBCInvalidClientId.selector, clientIdentifier));
        ics02Client.addClient(clientIdentifier, counterpartyInfo, lightClient);

        // reuse of client id
        string memory customClientId = "custom-client-id";
        ics02Client.addClient(customClientId, counterpartyInfo, lightClient);
        vm.expectRevert(abi.encodeWithSelector(IICS02ClientErrors.IBCClientAlreadyExists.selector, customClientId));
        ics02Client.addClient(customClientId, counterpartyInfo, lightClient);

        // unauthorized id customizer
        vm.stopPrank();
        address unauthorized = makeAddr("unauthorized");
        vm.prank(unauthorized);
        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        ics02Client.addClient(customClientId, counterpartyInfo, lightClient);
    }

    function test_MigrateClient() public {
        address unauthorized = makeAddr("unauthorized");
        address clientMigrator = makeAddr("clientMigrator");

        // Grant the per-clientId migrator role for clientIdentifier to clientMigrator.
        uint64 migratorRole = ics02Client.getLightClientMigratorRole(clientIdentifier);
        accessManager.grantRole(migratorRole, clientMigrator, 0);

        string memory counterpartyId = "42-dummy-01";
        address newLightClient = makeAddr("newLightClient");
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(counterpartyId, randomPrefix);

        // An address without the per-clientId migrator role is rejected.
        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(IICS02ClientErrors.IBCUnauthorizedMigrator.selector, clientIdentifier, unauthorized)
        );
        ics02Client.migrateClient(clientIdentifier, counterpartyInfo, newLightClient);

        // The granted migrator can migrate that specific clientId.
        vm.prank(clientMigrator);
        ics02Client.migrateClient(clientIdentifier, counterpartyInfo, newLightClient);
        ILightClient fetchedLightClient = ics02Client.getClient(clientIdentifier);
        assertEq(address(fetchedLightClient), newLightClient, "client not migrated");

        IICS02ClientMsgs.CounterpartyInfo memory fetchedCounterparty = ics02Client.getCounterparty(clientIdentifier);
        assertEq(fetchedCounterparty.clientId, counterpartyId, "counterparty not migrated");
        assertEq(fetchedCounterparty.merklePrefix, randomPrefix, "counterparty not migrated");
        assertEq(ics02Client.getNextClientSeq(), 1, "client seq not incremented");

        // A grant for one clientId must not authorize migration of any other clientId.
        // Use the unrestricted addClient (no custom id) to register a second default client.
        string memory otherClientId = ics02Client.addClient(counterpartyInfo, lightClient);

        IICS02ClientMsgs.CounterpartyInfo memory otherCounterparty =
            IICS02ClientMsgs.CounterpartyInfo(counterpartyId, randomPrefix);
        vm.prank(clientMigrator);
        vm.expectRevert(
            abi.encodeWithSelector(IICS02ClientErrors.IBCUnauthorizedMigrator.selector, otherClientId, clientMigrator)
        );
        ics02Client.migrateClient(otherClientId, otherCounterparty, newLightClient);
    }

    function test_failure_MigrateClient_withDelay() public {
        address delayedMigrator = makeAddr("delayedMigrator");

        // Grant the per-clientId migrator role for clientIdentifier to delayedMigrator with a non-zero execution delay.
        uint64 migratorRole = ics02Client.getLightClientMigratorRole(clientIdentifier);
        accessManager.grantRole(migratorRole, delayedMigrator, 60); // 60 seconds delay

        string memory counterpartyId = "42-dummy-01";
        address newLightClient = makeAddr("newLightClient");
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(counterpartyId, randomPrefix);

        // Even though they have the role, because it has a non-zero delay, the migration should revert immediately.
        vm.prank(delayedMigrator);
        vm.expectRevert(
            abi.encodeWithSelector(
                IICS02ClientErrors.IBCUnauthorizedMigrator.selector, clientIdentifier, delayedMigrator
            )
        );
        ics02Client.migrateClient(clientIdentifier, counterpartyInfo, newLightClient);
    }

    function test_failure_MigrateClient_whenClosed() public {
        address clientMigrator = makeAddr("clientMigrator");

        // Grant the per-clientId migrator role for clientIdentifier to clientMigrator.
        uint64 migratorRole = ics02Client.getLightClientMigratorRole(clientIdentifier);
        accessManager.grantRole(migratorRole, clientMigrator, 0);

        string memory counterpartyId = "42-dummy-01";
        address newLightClient = makeAddr("newLightClient");
        IICS02ClientMsgs.CounterpartyInfo memory counterpartyInfo =
            IICS02ClientMsgs.CounterpartyInfo(counterpartyId, randomPrefix);

        // Close the target contract via AccessManager.
        accessManager.setTargetClosed(address(ics02Client), true);

        // A granted migrator calling migrateClient should revert when the target is closed.
        vm.prank(clientMigrator);
        vm.expectRevert(
            abi.encodeWithSelector(
                IICS02ClientErrors.IBCUnauthorizedMigrator.selector, clientIdentifier, clientMigrator
            )
        );
        ics02Client.migrateClient(clientIdentifier, counterpartyInfo, newLightClient);
    }

    function test_Misbehaviour() public {
        bytes memory misbehaviourMsg = "testMisbehaviourMsg";
        bytes memory misbehaviourCall = abi.encodeCall(ILightClient.misbehaviour, (misbehaviourMsg));
        vm.mockCall(lightClient, misbehaviourCall, bytes(""));

        vm.expectCall(lightClient, misbehaviourCall);
        vm.prank(misbehaviourSubmitter);
        ics02Client.submitMisbehaviour(clientIdentifier, misbehaviourMsg);
    }

    function test_failure_submitMisbehaviour() public {
        address unauthorized = makeAddr("unauthorized");
        bytes memory misbehaviourMsg = "testMisbehaviourMsg";

        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        vm.prank(unauthorized);
        ics02Client.submitMisbehaviour(clientIdentifier, misbehaviourMsg);
    }

    function test_unfreezeClient() public {
        bytes memory unfreezeCall = abi.encodeCall(ILightClient.unfreeze, ());
        vm.mockCall(lightClient, unfreezeCall, bytes(""));

        vm.expectEmit();
        emit IICS02Client.ICS02ClientUnfrozen(clientIdentifier);
        ics02Client.unfreezeClient(clientIdentifier);
    }

    function test_failure_unfreezeClient() public {
        address unauthorized = makeAddr("unauthorized");

        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        vm.prank(unauthorized);
        ics02Client.unfreezeClient(clientIdentifier);
    }

    function test_success_updateClient() public {
        bytes memory updateMsg = "testUpdateMsg";
        bytes memory updateCall = abi.encodeCall(ILightClient.updateClient, (updateMsg));
        vm.mockCall(lightClient, updateCall, abi.encode(ILightClientMsgs.UpdateResult(0)));

        vm.expectCall(lightClient, updateCall);
        vm.prank(relayer);
        ics02Client.updateClient(clientIdentifier, updateMsg);
    }

    function test_success_reAnchorPinnedSet() public {
        IUpdateClientMsgs.MsgUpdateClient memory updateMsg;
        IICS07TendermintMsgs.ValidatorSet memory newPinnedValidatorSet;
        bytes memory reAnchorCall =
            abi.encodeCall(IGroth16ICS07Tendermint.reAnchorPinnedSet, (updateMsg, newPinnedValidatorSet));
        vm.mockCall(lightClient, reAnchorCall, bytes(""));

        vm.expectCall(lightClient, reAnchorCall);
        vm.prank(relayer);
        ics02Client.reAnchorPinnedSet(clientIdentifier, updateMsg, newPinnedValidatorSet);
    }

    function test_failure_reAnchorPinnedSet() public {
        address unauthorized = makeAddr("unauthorized");
        IUpdateClientMsgs.MsgUpdateClient memory updateMsg;
        IICS07TendermintMsgs.ValidatorSet memory newPinnedValidatorSet;

        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        vm.prank(unauthorized);
        ics02Client.reAnchorPinnedSet(clientIdentifier, updateMsg, newPinnedValidatorSet);
    }

    function test_failure_updateClient() public {
        address unauthorized = makeAddr("unauthorized");
        bytes memory updateMsg = "testUpdateMsg";

        vm.expectRevert(abi.encodeWithSelector(IAccessManaged.AccessManagedUnauthorized.selector, unauthorized));
        vm.prank(unauthorized);
        ics02Client.updateClient(clientIdentifier, updateMsg);
    }
}
