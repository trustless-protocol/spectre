// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

/// @notice Production deployment bootstrap. Unlike E2ETestDeploy, this script
/// never gives the verifier or final AccessManager authority to the broadcaster.

import { stdJson } from "forge-std/StdJson.sol";
import { Script } from "forge-std/Script.sol";
import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { TimelockController } from "@openzeppelin-contracts/governance/TimelockController.sol";
import { ERC1967Proxy } from "@openzeppelin-contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { ICS26Router } from "../contracts/ICS26Router.sol";
import { ICS20Transfer } from "../contracts/ICS20Transfer.sol";
import { SignatureVerifier } from "../contracts/light-clients/SignatureVerifier.sol";
import { Membership } from "../contracts/light-clients/modules/Membership.sol";
import { UpdateClient } from "../contracts/light-clients/modules/UpdateClient.sol";
import { Misbehaviour } from "../contracts/light-clients/modules/Misbehaviour.sol";
import { ClientMigrationProposer } from "../contracts/light-clients/modules/ClientMigrationProposer.sol";
import { ClientMigrationExecutor } from "../contracts/light-clients/modules/ClientMigrationExecutor.sol";
import { Escrow } from "../contracts/utils/Escrow.sol";
import { IBCERC20 } from "../contracts/utils/IBCERC20.sol";
import { ICS20Lib } from "../contracts/utils/ICS20Lib.sol";
import { IICS07TendermintMsgs } from "../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { DeployAccessManagerWithRoles } from "./deployments/DeployAccessManagerWithRoles.sol";
import { ProductionConfigLib } from "./deployments/ProductionConfigLib.sol";
import { IBCRolesLib } from "../contracts/utils/IBCRolesLib.sol";
import { IGroth16Verifier } from "../contracts/light-clients/interfaces/IGroth16Verifier.sol";

/// @dev GOVERNANCE_ADMIN must be an OpenZeppelin TimelockController (or a
///      compatible contract exposing `getMinDelay()`) whose delay is at least
///      SECURITY_DELAY. AccessManager ADMIN_ROLE intentionally has no delay;
///      governance timelocking protects changes made through that role.
contract ProductionDeploy is Script, IICS07TendermintMsgs, DeployAccessManagerWithRoles {
    using stdJson for string;

    uint32 internal constant DEFAULT_DELAY = 2 days;

    function run() public returns (string memory) {
        address bootstrap = vm.envAddress("BOOTSTRAP_ACCOUNT");
        address governance = vm.envAddress("GOVERNANCE_ADMIN");
        address upgrader = vm.envAddress("UPGRADER_ACCOUNT");
        address relayer = vm.envAddress("RELAYER_ACCOUNT");
        address pauser1 = vm.envAddress("PAUSER_1");
        address pauser2 = vm.envAddress("PAUSER_2");
        address unpauser = vm.envAddress("UNPAUSER_ACCOUNT");
        address watcher = vm.envAddress("MISBEHAVIOUR_WATCHER");
        address verifierN4 = vm.envAddress("VERIFIER_N4");
        address verifierN8 = vm.envAddress("VERIFIER_N8");
        address verifierN16 = vm.envAddress("VERIFIER_N16");
        address verifierN32 = vm.envAddress("VERIFIER_N32");
        address verifierN64 = vm.envAddress("VERIFIER_N64");
        address verifierN128 = vm.envAddress("VERIFIER_N128");
        uint256 configuredDelay = vm.envOr("SECURITY_DELAY", uint256(DEFAULT_DELAY));

        require(bootstrap != address(0), "invalid bootstrap account");
        require(governance != address(0) && governance != bootstrap, "invalid governance admin");
        require(upgrader != address(0) && upgrader != bootstrap, "invalid upgrader");
        require(governance.code.length != 0, "governance admin must be a contract");
        require(upgrader.code.length != 0, "upgrader must be a contract");
        require(relayer != address(0) && watcher != address(0), "missing operational account");
        require(pauser1 != address(0) && pauser2 != address(0), "need two pausers");
        require(unpauser != address(0), "missing safety account");

        address[] memory verifiers = ProductionConfigLib.verifierList(
            verifierN4, verifierN8, verifierN16, verifierN32, verifierN64, verifierN128
        );
        for (uint256 i = 0; i < verifiers.length; ++i) {
            require(verifiers[i].code.length != 0, "all verifier buckets require deployed code");
        }
        ProductionConfigLib.requireDistinct(verifiers, "verifier buckets must be distinct contracts");

        require(configuredDelay >= DEFAULT_DELAY && configuredDelay <= type(uint32).max, "invalid security delay");
        require(
            TimelockController(payable(governance)).getMinDelay() >= configuredDelay,
            "governance timelock delay too short"
        );

        // Every privileged account must be a separate key. Checking the whole set at once
        // rather than a handful of pairs closes the gap that matters most: an UNPAUSER_ACCOUNT
        // equal to a pauser would let one key both pause and unpause, defeating the delay
        // deliberately placed on unpausing — stopping the system is meant to be fast, and
        // restarting it slow and visible. The same reasoning applies to every other pair, so
        // none is left to be rediscovered at deploy time.
        address[] memory principals = ProductionConfigLib.principalList(
            bootstrap, governance, upgrader, relayer, pauser1, pauser2, unpauser, watcher
        );
        ProductionConfigLib.requireDistinct(principals, "role separation failure");

        uint32 delay = uint32(configuredDelay);

        // Pin the simulated and broadcast sender to the account that receives
        // the temporary AccessManager admin role. This avoids relying on
        // Forge's independent default-sender selection.
        vm.startBroadcast(bootstrap);

        AccessManager accessManager = new AccessManager(bootstrap);
        // AccessManager is the immutable verifier owner. Initial bucket setup is
        // executed through the manager while bootstrap is still its admin.
        SignatureVerifier signatureVerifier = new SignatureVerifier(address(accessManager));
        bytes4 verifierSelector = IGroth16Verifier.verifyProof.selector;
        _registerBucket(accessManager, signatureVerifier, 4, verifierN4, verifierSelector);
        _registerBucket(accessManager, signatureVerifier, 8, verifierN8, verifierSelector);
        _registerBucket(accessManager, signatureVerifier, 16, verifierN16, verifierSelector);
        _registerBucket(accessManager, signatureVerifier, 32, verifierN32, verifierSelector);
        _registerBucket(accessManager, signatureVerifier, 64, verifierN64, verifierSelector);
        _registerBucket(accessManager, signatureVerifier, 128, verifierN128, verifierSelector);

        address membership = address(new Membership());
        address updateClient = address(new UpdateClient(address(signatureVerifier)));
        address misbehaviour = address(new Misbehaviour(address(signatureVerifier)));
        address clientMigrationProposer = address(new ClientMigrationProposer());
        address clientMigrationExecutor = address(new ClientMigrationExecutor());
        address routerLogic = address(new ICS26Router(clientMigrationProposer, clientMigrationExecutor));
        address transferLogic = address(new ICS20Transfer());

        ERC1967Proxy routerProxy =
            new ERC1967Proxy(routerLogic, abi.encodeCall(ICS26Router.initialize, (address(accessManager))));
        ERC1967Proxy transferProxy = new ERC1967Proxy(
            transferLogic,
            abi.encodeCall(
                ICS20Transfer.initialize,
                (
                    address(routerProxy),
                    address(new Escrow()),
                    address(new IBCERC20()),
                    address(0),
                    address(accessManager)
                )
            )
        );

        // Register the default app while the selector still has its default
        // ADMIN_ROLE. Once target roles are installed, the bootstrap account no
        // longer has permission to call addIBCApp directly.
        ICS26Router(address(routerProxy)).addIBCApp(ICS20Lib.DEFAULT_PORT_ID, address(transferProxy));

        accessManagerSetTargetRoles(accessManager, address(routerProxy), address(transferProxy), false);
        accessManagerSetProductionUpgradeRoles(
            accessManager, address(routerProxy), address(transferProxy), address(signatureVerifier)
        );

        accessManager.grantRole(IBCRolesLib.RELAYER_ROLE, relayer, 0);
        accessManager.grantRole(IBCRolesLib.PAUSER_ROLE, pauser1, 0);
        accessManager.grantRole(IBCRolesLib.PAUSER_ROLE, pauser2, 0);
        accessManager.grantRole(IBCRolesLib.UNPAUSER_ROLE, unpauser, delay);
        accessManager.grantRole(IBCRolesLib.UPGRADER_ROLE, upgrader, delay);
        accessManager.grantRole(IBCRolesLib.MISBEHAVIOUR_SUBMITTER_ROLE, watcher, 0);
        accessManager.grantRole(IBCRolesLib.ID_CUSTOMIZER_ROLE, governance, delay);

        // The bootstrap key is removed before the deployment receipt is emitted.
        accessManager.grantRole(IBCRolesLib.ADMIN_ROLE, governance, 0);
        accessManager.renounceRole(IBCRolesLib.ADMIN_ROLE, bootstrap);
        vm.stopBroadcast();

        string memory json = "json";
        json.serialize("accessManager", vm.toString(address(accessManager)));
        json.serialize("signatureVerifier", vm.toString(address(signatureVerifier)));
        json.serialize("membership", vm.toString(membership));
        json.serialize("updateClient", vm.toString(updateClient));
        json.serialize("misbehaviour", vm.toString(misbehaviour));
        json.serialize("clientMigrationProposer", vm.toString(clientMigrationProposer));
        json.serialize("clientMigrationExecutor", vm.toString(clientMigrationExecutor));
        json.serialize("ics26Router", vm.toString(address(routerProxy)));
        return json.serialize("ics20Transfer", vm.toString(address(transferProxy)));
    }

    function _registerBucket(
        AccessManager accessManager,
        SignatureVerifier signatureVerifier,
        uint16 bucket,
        address verifier,
        bytes4 selector
    )
        internal
    {
        accessManager.execute(
            address(signatureVerifier), abi.encodeCall(SignatureVerifier.setBucket, (bucket, verifier, selector))
        );
    }
}
