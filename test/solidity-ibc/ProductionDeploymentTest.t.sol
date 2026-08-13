// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";
import { stdJson } from "forge-std/StdJson.sol";
import { TimelockController } from "@openzeppelin-contracts/governance/TimelockController.sol";
import { ProductionDeploy } from "../../scripts/ProductionDeploy.s.sol";
import { ProductionVerify } from "../../scripts/ProductionVerify.s.sol";

contract ProductionVerifierMock {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure { }
}

contract ProductionUpgraderMock { }

/// @notice Exercises the production deployment and verification scripts together.
contract ProductionDeploymentTest is Test {
    using stdJson for string;

    uint256 internal constant SECURITY_DELAY = 2 days;

    /// @notice Populates a configuration the scripts must accept. Negative tests below start
    ///         from this and break exactly one thing, so a revert can only come from that.
    function _setValidEnv() internal {
        address[] memory proposers = new address[](1);
        address[] memory executors = new address[](1);
        proposers[0] = address(this);
        executors[0] = address(this);
        address governance = address(new TimelockController(SECURITY_DELAY, proposers, executors, address(this)));

        _setAddressEnv("BOOTSTRAP_ACCOUNT", address(0xB001));
        _setAddressEnv("GOVERNANCE_ADMIN", governance);
        _setAddressEnv("UPGRADER_ACCOUNT", address(new ProductionUpgraderMock()));
        _setAddressEnv("RELAYER_ACCOUNT", address(0xB002));
        _setAddressEnv("PAUSER_1", address(0xB003));
        _setAddressEnv("PAUSER_2", address(0xB004));
        _setAddressEnv("UNPAUSER_ACCOUNT", address(0xB005));
        _setAddressEnv("MISBEHAVIOUR_WATCHER", address(0xB006));

        // Six SEPARATE verifiers. Each bucket's verifier is generated from its own verifying
        // key, so sharing one address across buckets is a misconfiguration the scripts must
        // reject — a fixture that reuses one mock cannot exercise that check.
        _setAddressEnv("VERIFIER_N4", address(new ProductionVerifierMock()));
        _setAddressEnv("VERIFIER_N8", address(new ProductionVerifierMock()));
        _setAddressEnv("VERIFIER_N16", address(new ProductionVerifierMock()));
        _setAddressEnv("VERIFIER_N32", address(new ProductionVerifierMock()));
        _setAddressEnv("VERIFIER_N64", address(new ProductionVerifierMock()));
        _setAddressEnv("VERIFIER_N128", address(new ProductionVerifierMock()));

        vm.setEnv("SECURITY_DELAY", vm.toString(SECURITY_DELAY));
    }

    /// @notice Deploy, verify, and the two preconditions that cannot be checked after the fact.
    ///
    /// All of it lives in ONE test function on purpose. `vm.setEnv` writes the PROCESS
    /// environment, which is shared across test functions while each function gets its own EVM
    /// state — so a second function reads the first function's addresses and finds no code at
    /// them. Forge also runs functions in parallel, which makes the clobbering nondeterministic.
    /// Splitting these apart is what made this file fail intermittently.
    function testProductionDeployVerifyAndRejectBrokenConfigs() public {
        // 1. A valid configuration deploys, and the result passes verification.
        _setValidEnv();

        string memory deployment = new ProductionDeploy().run();
        _setAddressEnv("ACCESS_MANAGER", deployment.readAddress(".accessManager"));
        _setAddressEnv("SIGNATURE_VERIFIER", deployment.readAddress(".signatureVerifier"));
        _setAddressEnv("ICS26_ROUTER", deployment.readAddress(".ics26Router"));
        _setAddressEnv("ICS20_TRANSFER", deployment.readAddress(".ics20Transfer"));

        assertTrue(new ProductionVerify().run());

        ProductionVerify verifier = new ProductionVerify();
        address configuredUpgrader = vm.envAddress("UPGRADER_ACCOUNT");
        address configuredUnpauser = vm.envAddress("UNPAUSER_ACCOUNT");

        // Verification must reject old deployments that predate deploy-time role separation.
        _setAddressEnv("UNPAUSER_ACCOUNT", vm.envAddress("PAUSER_1"));
        vm.expectRevert("role separation failure");
        verifier.run();
        _setAddressEnv("UNPAUSER_ACCOUNT", configuredUnpauser);

        // The upgrade role must remain contract-owned after deployment, not only at bootstrap.
        _setAddressEnv("UPGRADER_ACCOUNT", address(0xB007));
        vm.expectRevert("upgrader must be a contract");
        verifier.run();
        _setAddressEnv("UPGRADER_ACCOUNT", configuredUpgrader);

        ProductionDeploy deployer = new ProductionDeploy();

        // 2. Two buckets sharing one verifier. Not a harmless duplicate: each bucket's verifier
        // is generated from its own verifying key, so every proof in the mis-pointed bucket
        // fails on-chain against a VK built for a different signer count. With six
        // near-identical env vars this is a plausible copy-paste, and nothing else catches it —
        // the address is well-formed, the contract has code, and the selector matches.
        _setValidEnv();
        _setAddressEnv("VERIFIER_N8", vm.envAddress("VERIFIER_N4"));
        vm.expectRevert("verifier buckets must be distinct contracts");
        deployer.run();

        // 3. An unpauser that is also a pauser. UNPAUSER_ROLE carries the security delay and
        // PAUSER_ROLE does not, precisely so stopping the system is immediate while restarting
        // it is slow and visible. One key holding both collapses that distinction.
        _setValidEnv();
        _setAddressEnv("UNPAUSER_ACCOUNT", vm.envAddress("PAUSER_1"));
        vm.expectRevert("role separation failure");
        deployer.run();
    }

    function _setAddressEnv(string memory name, address value) internal {
        vm.setEnv(name, vm.toString(value));
    }
}
