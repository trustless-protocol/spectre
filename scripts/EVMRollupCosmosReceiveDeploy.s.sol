// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { stdJson } from "forge-std/StdJson.sol";
import { Script } from "forge-std/Script.sol";

import { AccessManager } from "@openzeppelin-contracts/access/manager/AccessManager.sol";
import { ERC1967Proxy } from "@openzeppelin-contracts/proxy/ERC1967/ERC1967Proxy.sol";
import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";

import { ICS20Transfer } from "../contracts/ICS20Transfer.sol";
import { ICS26Router } from "../contracts/ICS26Router.sol";
import { SignatureVerifier } from "../contracts/light-clients/SignatureVerifier.sol";
import { SpectreClient } from "../contracts/light-clients/SpectreClient.sol";
import { Membership } from "../contracts/light-clients/modules/Membership.sol";
import { Misbehaviour } from "../contracts/light-clients/modules/Misbehaviour.sol";
import { UpdateClient } from "../contracts/light-clients/modules/UpdateClient.sol";
import { IICS07TendermintMsgs } from "../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../contracts/msgs/IICS02ClientMsgs.sol";
import { Groth16Verifier_N4 } from "../contracts/verifiers/Groth16Verifier_N4.sol";
import { IBCERC20 } from "../contracts/utils/IBCERC20.sol";
import { Escrow } from "../contracts/utils/Escrow.sol";
import { ICS20Lib } from "../contracts/utils/ICS20Lib.sol";
import { RelayerHelper } from "../contracts/utils/RelayerHelper.sol";
import { DeployAccessManagerWithRoles } from "./deployments/DeployAccessManagerWithRoles.sol";

/// @notice Deploys the EVM rollup-side contracts needed to verify Cosmos state:
///         recvPacket, ackPacket, and timeoutPacket all verify through SpectreClient.
/// @dev The script grants ICS26Router only SpectreClient's proof submitter role, because
///      recv/ack/timeout proof checks are invoked by the router contract. Admin and
///      misbehaviour-submitter roles stay with ROLE_MANAGER unless permissionless mode is explicit.
contract EVMRollupCosmosReceiveDeploy is Script, DeployAccessManagerWithRoles {
    using stdJson for string;

    string internal constant DEFAULT_SOURCE_CLIENT_ID = "evm-rollup-cosmos-0";
    string internal constant DEFAULT_COUNTERPARTY_CLIENT_ID = "cosmos-evm-rollup-0";
    bytes32 internal constant PROOF_SUBMITTER_ROLE = keccak256("PROOF_SUBMITTER_ROLE");
    bytes32 internal constant DEFAULT_ADMIN_ROLE = 0x00;

    struct Deployment {
        address accessManager;
        address ics26Router;
        address ics20Transfer;
        address relayerHelper;
        address signatureVerifier;
        address verifierN4;
        address membership;
        address updateClient;
        address misbehaviour;
        address spectreClient;
        address roleManager;
        bool permissionlessSpectreRoles;
    }

    struct SpectreBootstrap {
        bytes clientState;
        bytes32 consensusStateHash;
        IICS07TendermintMsgs.ValidatorSet validatorSet;
    }

    function run() public returns (string memory) {
        string memory sourceClientId = vm.envOr("SOURCE_CLIENT_ID", DEFAULT_SOURCE_CLIENT_ID);
        string memory counterpartyClientId = vm.envOr("COUNTERPARTY_CLIENT_ID", DEFAULT_COUNTERPARTY_CLIENT_ID);
        address permit2 = vm.envOr("PERMIT2", address(0));
        address roleManager = vm.envOr("ROLE_MANAGER", msg.sender);
        bool permissionlessSpectreRoles = vm.envOr("SPECTRE_PERMISSIONLESS_ROLES", false);
        string memory outputPath = vm.envOr("EVM_ROLLUP_RECEIVE_DEPLOYMENT_OUT", string(""));
        SpectreBootstrap memory bootstrap = _loadSpectreBootstrap();

        require(
            permissionlessSpectreRoles || roleManager != address(0),
            "ROLE_MANAGER zero requires SPECTRE_PERMISSIONLESS_ROLES=true"
        );

        vm.startBroadcast();

        Deployment memory deployment = _deployReceiveStack(permit2, roleManager, permissionlessSpectreRoles, bootstrap);
        _configureSourceClient(deployment, sourceClientId, counterpartyClientId);

        vm.stopBroadcast();

        string memory finalJson = _deploymentJson(deployment, sourceClientId, counterpartyClientId);

        if (bytes(outputPath).length != 0) {
            vm.writeFile(outputPath, finalJson);
        }

        return finalJson;
    }

    function _loadSpectreBootstrap() private view returns (SpectreBootstrap memory bootstrap) {
        bootstrap.clientState = vm.envBytes("SPECTRE_CLIENT_STATE");
        bootstrap.consensusStateHash = vm.envBytes32("SPECTRE_CONSENSUS_STATE_HASH");
        bootstrap.validatorSet = abi.decode(vm.envBytes("SPECTRE_VALIDATOR_SET"), (IICS07TendermintMsgs.ValidatorSet));
    }

    function _deployReceiveStack(
        address permit2,
        address configuredRoleManager,
        bool permissionlessSpectreRoles,
        SpectreBootstrap memory bootstrap
    )
        private
        returns (Deployment memory deployment)
    {
        SignatureVerifier signatureVerifier = new SignatureVerifier(msg.sender);
        Groth16Verifier_N4 verifierN4 = new Groth16Verifier_N4();
        signatureVerifier.setBucket(4, address(verifierN4), Groth16Verifier_N4.verifyProof.selector);

        Membership membership = new Membership();
        UpdateClient updateClient = new UpdateClient(address(signatureVerifier));
        Misbehaviour misbehaviour = new Misbehaviour(address(signatureVerifier));

        AccessManager accessManager = new AccessManager(msg.sender);

        ICS26Router routerLogic = new ICS26Router();
        ERC1967Proxy routerProxy =
            new ERC1967Proxy(address(routerLogic), abi.encodeCall(ICS26Router.initialize, (address(accessManager))));

        ICS20Transfer transferLogic = new ICS20Transfer();
        ERC1967Proxy transferProxy = new ERC1967Proxy(
            address(transferLogic),
            abi.encodeCall(
                ICS20Transfer.initialize,
                (address(routerProxy), address(new Escrow()), address(new IBCERC20()), permit2, address(accessManager))
            )
        );

        RelayerHelper relayerHelper = new RelayerHelper(address(routerProxy));

        accessManagerSetTargetRoles(accessManager, address(routerProxy), address(transferProxy), false);

        address[] memory relayers = new address[](1);
        relayers[0] = msg.sender;
        accessManagerSetRoles(
            accessManager, relayers, new address[](0), new address[](0), msg.sender, msg.sender, msg.sender
        );

        address constructorRoleManager = permissionlessSpectreRoles ? address(0) : msg.sender;
        SpectreClient spectreClient = new SpectreClient(
            address(updateClient),
            address(membership),
            address(misbehaviour),
            bootstrap.clientState,
            bootstrap.consensusStateHash,
            bootstrap.validatorSet,
            constructorRoleManager
        );

        if (!permissionlessSpectreRoles) {
            spectreClient.grantRole(PROOF_SUBMITTER_ROLE, address(routerProxy));
            if (configuredRoleManager != msg.sender) {
                spectreClient.grantRole(DEFAULT_ADMIN_ROLE, configuredRoleManager);
                spectreClient.grantRole(spectreClient.MISBEHAVIOUR_SUBMITTER_ROLE(), configuredRoleManager);
                spectreClient.grantRole(PROOF_SUBMITTER_ROLE, configuredRoleManager);
            }
        }

        deployment = Deployment({
            accessManager: address(accessManager),
            ics26Router: address(routerProxy),
            ics20Transfer: address(transferProxy),
            relayerHelper: address(relayerHelper),
            signatureVerifier: address(signatureVerifier),
            verifierN4: address(verifierN4),
            membership: address(membership),
            updateClient: address(updateClient),
            misbehaviour: address(misbehaviour),
            spectreClient: address(spectreClient),
            roleManager: permissionlessSpectreRoles ? address(0) : configuredRoleManager,
            permissionlessSpectreRoles: permissionlessSpectreRoles
        });
    }

    function _configureSourceClient(
        Deployment memory deployment,
        string memory sourceClientId,
        string memory counterpartyClientId
    )
        private
    {
        bytes[] memory cosmosMerklePrefix = new bytes[](2);
        cosmosMerklePrefix[0] = bytes("ibc");
        cosmosMerklePrefix[1] = bytes("");

        ICS26Router(deployment.ics26Router)
            .addClient(
                sourceClientId,
                IICS02ClientMsgs.CounterpartyInfo({ clientId: counterpartyClientId, merklePrefix: cosmosMerklePrefix }),
                deployment.spectreClient
            );
        ICS26Router(deployment.ics26Router).addIBCApp(ICS20Lib.DEFAULT_PORT_ID, deployment.ics20Transfer);
    }

    function _deploymentJson(
        Deployment memory deployment,
        string memory sourceClientId,
        string memory counterpartyClientId
    )
        private
        returns (string memory)
    {
        string memory json = "json";
        json.serialize("chainId", block.chainid);
        json.serialize("deployer", Strings.toHexString(msg.sender));
        json.serialize("accessManager", Strings.toHexString(deployment.accessManager));
        json.serialize("ics26Router", Strings.toHexString(deployment.ics26Router));
        json.serialize("ics20Transfer", Strings.toHexString(deployment.ics20Transfer));
        json.serialize("relayerHelper", Strings.toHexString(deployment.relayerHelper));
        json.serialize("signatureVerifier", Strings.toHexString(deployment.signatureVerifier));
        json.serialize("verifierN4", Strings.toHexString(deployment.verifierN4));
        json.serialize("membership", Strings.toHexString(deployment.membership));
        json.serialize("updateClient", Strings.toHexString(deployment.updateClient));
        json.serialize("misbehaviour", Strings.toHexString(deployment.misbehaviour));
        json.serialize("spectreClient", Strings.toHexString(deployment.spectreClient));
        json.serialize("roleManager", Strings.toHexString(deployment.roleManager));
        json.serialize("permissionlessSpectreRoles", deployment.permissionlessSpectreRoles);
        json.serialize("sourceClientId", sourceClientId);
        json.serialize("counterpartyClientId", counterpartyClientId);
        json.serialize("supportsAcknowledgements", true);
        json.serialize("supportsTimeouts", true);
        string memory finalJson = json.serialize("supportedSignatureBucket", uint256(4));

        return finalJson;
    }
}
