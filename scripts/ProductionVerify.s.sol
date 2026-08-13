// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Script } from "forge-std/Script.sol";
import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { TimelockController } from "@openzeppelin-contracts/governance/TimelockController.sol";
import { ICS26Router } from "../contracts/ICS26Router.sol";
import { ICS20Transfer } from "../contracts/ICS20Transfer.sol";
import { SignatureVerifier } from "../contracts/light-clients/SignatureVerifier.sol";
import { IBCRolesLib } from "../contracts/utils/IBCRolesLib.sol";
import { ProductionConfigLib } from "./deployments/ProductionConfigLib.sol";
import { IGroth16Verifier } from "../contracts/light-clients/interfaces/IGroth16Verifier.sol";

/// @notice Read-only post-deployment checks for ProductionDeploy.
/// @dev Governance must be a timelock whose minimum delay covers the
///      AccessManager security delay; ADMIN_ROLE is intentionally immediate.
contract ProductionVerify is Script {
    uint32 internal constant DEFAULT_DELAY = 2 days;

    function run() external view returns (bool) {
        AccessManager manager = AccessManager(vm.envAddress("ACCESS_MANAGER"));
        SignatureVerifier verifier = SignatureVerifier(vm.envAddress("SIGNATURE_VERIFIER"));
        address router = vm.envAddress("ICS26_ROUTER");
        address transfer = vm.envAddress("ICS20_TRANSFER");
        address governance = vm.envAddress("GOVERNANCE_ADMIN");
        address upgrader = vm.envAddress("UPGRADER_ACCOUNT");
        address bootstrap = vm.envAddress("BOOTSTRAP_ACCOUNT");
        address relayer = vm.envAddress("RELAYER_ACCOUNT");
        address pauser1 = vm.envAddress("PAUSER_1");
        address pauser2 = vm.envAddress("PAUSER_2");
        address unpauser = vm.envAddress("UNPAUSER_ACCOUNT");
        address watcher = vm.envAddress("MISBEHAVIOUR_WATCHER");
        uint256 configuredDelay = vm.envOr("SECURITY_DELAY", uint256(DEFAULT_DELAY));

        require(configuredDelay >= DEFAULT_DELAY && configuredDelay <= type(uint32).max, "invalid security delay");
        uint32 delay = uint32(configuredDelay);

        require(address(manager).code.length != 0, "access manager has no code");
        require(router.code.length != 0 && transfer.code.length != 0, "IBC proxy has no code");
        require(governance.code.length != 0, "governance admin has no code");
        require(upgrader.code.length != 0, "upgrader must be a contract");
        require(
            TimelockController(payable(governance)).getMinDelay() >= configuredDelay,
            "governance timelock delay too short"
        );
        require(verifier.OWNER() == address(manager), "signature verifier owner mismatch");
        require(ICS26Router(router).authority() == address(manager), "router authority mismatch");
        require(ICS20Transfer(transfer).authority() == address(manager), "transfer authority mismatch");

        address[] memory principals = ProductionConfigLib.principalList(
            bootstrap, governance, upgrader, relayer, pauser1, pauser2, unpauser, watcher
        );
        ProductionConfigLib.requireDistinct(principals, "role separation failure");

        bytes4 verifierSelector = IGroth16Verifier.verifyProof.selector;
        address[] memory verifiers = ProductionConfigLib.verifierList(
            vm.envAddress("VERIFIER_N4"),
            vm.envAddress("VERIFIER_N8"),
            vm.envAddress("VERIFIER_N16"),
            vm.envAddress("VERIFIER_N32"),
            vm.envAddress("VERIFIER_N64"),
            vm.envAddress("VERIFIER_N128")
        );
        // Re-checked here, not only at deploy time: a deployment made before this precondition
        // existed must not pass verification. See ProductionConfigLib.verifierList.
        ProductionConfigLib.requireDistinct(verifiers, "verifier buckets must be distinct contracts");

        uint16[6] memory buckets = [uint16(4), 8, 16, 32, 64, 128];
        for (uint256 i = 0; i < buckets.length; ++i) {
            _checkBucket(verifier, buckets[i], verifiers[i], verifierSelector);
        }

        _checkRoleAbsent(manager, IBCRolesLib.ADMIN_ROLE, bootstrap, "bootstrap admin role still active");
        _checkRole(manager, IBCRolesLib.ADMIN_ROLE, governance, "governance admin role missing");
        _checkRoleDelay(manager, IBCRolesLib.UPGRADER_ROLE, upgrader, delay, "upgrader role delay is too short");
        _checkRoleExactDelay(manager, IBCRolesLib.RELAYER_ROLE, relayer, 0, "relayer role delay mismatch");
        _checkRoleExactDelay(manager, IBCRolesLib.PAUSER_ROLE, pauser1, 0, "pauser 1 role delay mismatch");
        _checkRoleExactDelay(manager, IBCRolesLib.PAUSER_ROLE, pauser2, 0, "pauser 2 role delay mismatch");
        _checkRoleDelay(manager, IBCRolesLib.UNPAUSER_ROLE, unpauser, delay, "unpauser role delay is too short");
        _checkRoleExactDelay(
            manager, IBCRolesLib.MISBEHAVIOUR_SUBMITTER_ROLE, watcher, 0, "watcher role delay mismatch"
        );
        _checkRoleDelay(
            manager, IBCRolesLib.ID_CUSTOMIZER_ROLE, governance, delay, "ID customizer role delay is too short"
        );

        _checkSelectors(manager, router, IBCRolesLib.ics26RelayerSelectors(), IBCRolesLib.RELAYER_ROLE);
        _checkSelectors(manager, router, IBCRolesLib.ics26IdCustomizerSelectors(), IBCRolesLib.ID_CUSTOMIZER_ROLE);
        _checkSelectors(manager, router, IBCRolesLib.pauserSelectors(), IBCRolesLib.PAUSER_ROLE);
        _checkSelectors(manager, router, IBCRolesLib.unpauserSelectors(), IBCRolesLib.UNPAUSER_ROLE);
        _checkSelectors(
            manager, router, IBCRolesLib.ics26MisbehaviourSelectors(), IBCRolesLib.MISBEHAVIOUR_SUBMITTER_ROLE
        );
        _checkSelectors(manager, router, IBCRolesLib.uupsUpgradeSelectors(), IBCRolesLib.UPGRADER_ROLE);
        _checkSelectors(manager, transfer, IBCRolesLib.pauserSelectors(), IBCRolesLib.PAUSER_ROLE);
        _checkSelectors(manager, transfer, IBCRolesLib.unpauserSelectors(), IBCRolesLib.UNPAUSER_ROLE);
        _checkSelectors(manager, transfer, IBCRolesLib.upgraderSelectors(), IBCRolesLib.UPGRADER_ROLE);

        bytes4[] memory setBucketSelector = new bytes4[](1);
        setBucketSelector[0] = SignatureVerifier.setBucket.selector;
        _checkSelectors(manager, address(verifier), setBucketSelector, IBCRolesLib.UPGRADER_ROLE);
        return true;
    }

    function _checkBucket(SignatureVerifier verifier, uint16 bucket, address expected, bytes4 selector) private view {
        (address implementation, bytes4 configuredSelector) = verifier.buckets(bucket);
        require(implementation == expected, "bucket verifier address mismatch");
        require(implementation.code.length != 0, "bucket verifier has no code");
        require(configuredSelector == selector, "bucket selector mismatch");
    }

    function _checkRole(AccessManager manager, uint64 role, address account, string memory message) private view {
        (bool member,) = manager.hasRole(role, account);
        require(member, message);
    }

    function _checkRoleDelay(
        AccessManager manager,
        uint64 role,
        address account,
        uint32 minimumDelay,
        string memory message
    )
        private
        view
    {
        (bool member, uint32 delay) = manager.hasRole(role, account);
        require(member && delay >= minimumDelay, message);
    }

    function _checkRoleExactDelay(
        AccessManager manager,
        uint64 role,
        address account,
        uint32 expectedDelay,
        string memory message
    )
        private
        view
    {
        (bool member, uint32 delay) = manager.hasRole(role, account);
        require(member && delay == expectedDelay, message);
    }

    function _checkRoleAbsent(AccessManager manager, uint64 role, address account, string memory message) private view {
        (bool member,) = manager.hasRole(role, account);
        require(!member, message);
    }

    function _checkSelectors(
        AccessManager manager,
        address target,
        bytes4[] memory selectors,
        uint64 role
    )
        private
        view
    {
        for (uint256 i = 0; i < selectors.length; ++i) {
            require(manager.getTargetFunctionRole(target, selectors[i]) == role, "target role mismatch");
        }
    }
}
