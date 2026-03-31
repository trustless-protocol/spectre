// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// solhint-disable gas-custom-errors

// solhint-disable-next-line no-global-import
import "forge-std/console.sol";
import { Test, stdStorage, StdStorage } from "forge-std/Test.sol";
import { stdJson } from "forge-std/StdJson.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { Membership } from "../../contracts/programs/Membership.sol";
import { Misbehaviour } from "../../contracts/programs/Misbehaviour.sol";
import { UpdateClient } from "../../contracts/programs/UpdateClient.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";
import { IUpdateClientAndMembershipMsgs } from "../../contracts/light-clients/msgs/IUcAndMembershipMsgs.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { IGroth16Msgs } from "../../contracts/light-clients/msgs/IGroth16Msgs.sol";
import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { IGroth16ICS07TendermintErrors } from "../../contracts/light-clients/errors/IGroth16ICS07TendermintErrors.sol";
import { MockGroth16Verifier } from "@groth16-contracts/MockGroth16Verifier.sol";
import { Groth16Verifier as PlonkVerifier } from "@groth16-contracts/v5.0.0/PlonkVerifier.sol";
import { Groth16Verifier as Groth16Verifier } from "@groth16-contracts/v5.0.0/Groth16Verifier.sol";

struct Groth16ICS07GenesisFixtureJson {
    bytes trustedClientState;
    bytes trustedConsensusState;
    bytes32 updateClientVkey;
    bytes32 membershipVkey;
    bytes32 ucAndMembershipVkey;
    bytes32 misbehaviourVkey;
}

abstract contract Groth16ICS07TendermintTest is
    Test,
    IICS02ClientMsgs,
    IGroth16Msgs,
    IICS07TendermintMsgs,
    IUpdateClientMsgs,
    IMembershipMsgs,
    IUpdateClientAndMembershipMsgs,
    IGroth16ICS07TendermintErrors,
    ILightClientMsgs
{
    using stdJson for string;
    using stdStorage for StdStorage;

    Groth16ICS07Tendermint public ics07Tendermint;
    Groth16ICS07Tendermint public mockIcs07Tendermint;

    Groth16ICS07GenesisFixtureJson internal genesisFixture;

    string internal constant FIXTURE_DIR = "/test/groth16-ics07/fixtures/";

    function setUpTest(string memory fileName, address roleManager) public {
        genesisFixture = loadGenesisFixture(fileName);

        ConsensusState memory trustedConsensusState = abi.decode(genesisFixture.trustedConsensusState, (ConsensusState));

        bytes32 trustedConsensusHash = keccak256(abi.encode(trustedConsensusState));
        ClientState memory trustedClientState = abi.decode(genesisFixture.trustedClientState, (ClientState));

        address verifier;
        if (trustedClientState.zkAlgorithm == SupportedZkAlgorithm.Plonk) {
            verifier = address(new PlonkVerifier());
        } else if (trustedClientState.zkAlgorithm == SupportedZkAlgorithm.Groth16) {
            verifier = address(new Groth16Verifier());
        } else {
            revert("Unsupported zk algorithm");
        }

        address membership = address(new Membership());
        address misbehaviour = address(new Misbehaviour());
        address updateClient = address(new UpdateClient());

        ics07Tendermint = new Groth16ICS07Tendermint(
            // genesisFixture.updateClientVkey,
            // genesisFixture.membershipVkey,
            // genesisFixture.ucAndMembershipVkey,
            // genesisFixture.misbehaviourVkey,
            verifier,
            membership,
            misbehaviour,
            updateClient,
            genesisFixture.trustedClientState,
            trustedConsensusHash,
            roleManager
        );

        mockIcs07Tendermint = new Groth16ICS07Tendermint(
            // genesisFixture.updateClientVkey,
            // genesisFixture.membershipVkey,
            // genesisFixture.ucAndMembershipVkey,
            // genesisFixture.misbehaviourVkey,
            address(new MockGroth16Verifier()),
            membership,
            misbehaviour,
            updateClient,
            genesisFixture.trustedClientState,
            trustedConsensusHash,
            roleManager
        );

        ClientState memory clientState = abi.decode(mockIcs07Tendermint.getClientState(), (ClientState));
        assert(keccak256(abi.encode(clientState)) == keccak256(genesisFixture.trustedClientState));

        bytes32 consensusHash = mockIcs07Tendermint.getConsensusStateHash(clientState.latestHeight.revisionHeight);
        assert(consensusHash == trustedConsensusHash);
    }

    function loadGenesisFixture(string memory fileName) public view returns (Groth16ICS07GenesisFixtureJson memory) {
        string memory root = vm.projectRoot();
        string memory path = string.concat(root, FIXTURE_DIR, fileName);
        string memory json = vm.readFile(path);
        bytes memory trustedClientState = json.readBytes(".trustedClientState");
        bytes memory trustedConsensusState = json.readBytes(".trustedConsensusState");
        bytes32 updateClientVkey = json.readBytes32(".updateClientVkey");
        bytes32 membershipVkey = json.readBytes32(".membershipVkey");
        bytes32 ucAndMembershipVkey = json.readBytes32(".ucAndMembershipVkey");
        bytes32 misbehaviourVkey = json.readBytes32(".misbehaviourVkey");

        Groth16ICS07GenesisFixtureJson memory fix = Groth16ICS07GenesisFixtureJson({
            trustedClientState: trustedClientState,
            trustedConsensusState: trustedConsensusState,
            updateClientVkey: updateClientVkey,
            membershipVkey: membershipVkey,
            ucAndMembershipVkey: ucAndMembershipVkey,
            misbehaviourVkey: misbehaviourVkey
        });

        return fix;
    }

    function _nanosToSeconds(uint256 nanos) internal pure returns (uint256) {
        return nanos / 1e9;
    }

    struct FixtureTestCase {
        string name;
        string fileName;
    }
}
