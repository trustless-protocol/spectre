// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { IGroth16Msgs } from "../../contracts/light-clients/msgs/IGroth16Msgs.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";

import { MockGroth16Verifier } from "@groth16-contracts/MockGroth16Verifier.sol";

abstract contract Groth16ICS07MockTest is Test {
    string public constant MOCK_CHAIN_ID = "mock-chain";
    bytes32 public constant MOCK_VKEY = keccak256("MOCK_VKEY");
    bytes32 public constant MOCK_ROOT = keccak256("MOCK_ROOT");
    bytes32 public constant MOCK_VAL_HASH = keccak256("MOCK_VAL_HASH");

    address public roleManager = makeAddr("roleManager");
    address public proofSubmitter = makeAddr("proofSubmitter");

    Groth16ICS07Tendermint public ics07Tendermint;

    bytes[] public membershipPath = [bytes("ibc"), bytes("path")];
    bytes public membershipValue = bytes("value");

    function setUp() public {
        bytes32 firstConsensusStateHash = keccak256(abi.encode(newMockConsensusState(1)));

        ics07Tendermint = new Groth16ICS07Tendermint(
            // MOCK_VKEY,
            // MOCK_VKEY,
            // MOCK_VKEY,
            // MOCK_VKEY,
            address(new MockGroth16Verifier()),
            address(0),
            address(0),
            address(0),
            abi.encode(mockClientState(1)),
            firstConsensusStateHash,
            roleManager
        );

        bytes32 proofSubmitterRole = ics07Tendermint.PROOF_SUBMITTER_ROLE();
        vm.prank(roleManager);
        ics07Tendermint.grantRole(proofSubmitterRole, proofSubmitter);
    }

    function mockClientState(uint64 height) public pure returns (IICS07TendermintMsgs.ClientState memory) {
        return IICS07TendermintMsgs.ClientState({
            chainId: MOCK_CHAIN_ID,
            trustLevel: IICS07TendermintMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: height }),
            trustingPeriod: 1 weeks,
            unbondingPeriod: 2 weeks,
            isFrozen: false,
            zkAlgorithm: IICS07TendermintMsgs.SupportedZkAlgorithm.Groth16
        });
    }

    /// @notice Create a new mock consensus state
    /// @param timestamp The timestamp of the consensus state in unix nanoseconds
    /// @return The new consensus state
    function newMockConsensusState(uint128 timestamp)
        public
        pure
        returns (IICS07TendermintMsgs.ConsensusState memory)
    {
        return IICS07TendermintMsgs.ConsensusState({
            timestamp: timestamp,
            root: MOCK_ROOT,
            nextValidatorsHash: MOCK_VAL_HASH
        });
    }

    function newUpdateClientMsg() public view returns (bytes memory) {
        IICS07TendermintMsgs.ClientState memory clientState =
            abi.decode(ics07Tendermint.getClientState(), (IICS07TendermintMsgs.ClientState));
        IICS02ClientMsgs.Height memory trustedHeight =
            IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: clientState.latestHeight.revisionHeight });
        ++clientState.latestHeight.revisionHeight;

    //     IUpdateClientMsgs.UpdateClientOutput memory output = IUpdateClientMsgs.UpdateClientOutput({
    //         clientState: clientState,
    //         trustedConsensusState: newMockConsensusState(trustedHeight.revisionHeight),
    //         newConsensusState: newMockConsensusState(clientState.latestHeight.revisionHeight),
    //         time: uint128(block.timestamp * 1e9),
    //         trustedHeight: trustedHeight,
    //         newHeight: clientState.latestHeight
    //     });

    //     return abi.encode(
    //         IUpdateClientMsgs.MsgUpdateClient({
    //             groth16Proof: IGroth16Msgs.Groth16Proof({ vKey: MOCK_VKEY, publicValues: abi.encode(output), proof: bytes("") })
    //         })
    //     );
    // }

    // function newMembershipMsg(uint64 height) public view returns (ILightClientMsgs.MsgVerifyMembership memory) {
    //     IMembershipMsgs.MembershipOutput memory output =
    //         IMembershipMsgs.MembershipOutput({ commitmentRoot: MOCK_ROOT, kvPairs: new IMembershipMsgs.KVPair[](1) });
    //     output.kvPairs[0] = IMembershipMsgs.KVPair({ path: membershipPath, value: membershipValue });

    //     IMembershipMsgs.Groth16MembershipProof memory groth16Proof = IMembershipMsgs.Groth16MembershipProof({
    //         groth16Proof: IGroth16Msgs.Groth16Proof({ vKey: MOCK_VKEY, publicValues: abi.encode(output), proof: bytes("") }),
    //         trustedConsensusState: newMockConsensusState(height)
    //     });

    //     IMembershipMsgs.MembershipProof memory proof = IMembershipMsgs.MembershipProof({
    //         proofType: IMembershipMsgs.MembershipProofType.Groth16MembershipProof,
    //         proof: abi.encode(groth16Proof)
    //     });

    //     return ILightClientMsgs.MsgVerifyMembership({
    //         proof: abi.encode(proof),
    //         proofHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: height }),
    //         path: membershipPath,
    //         value: membershipValue
    //     });
    // }

    // function newNonMembershipMsg(uint64 height) public view returns (ILightClientMsgs.MsgVerifyNonMembership memory) {
    //     IMembershipMsgs.MembershipOutput memory output =
    //         IMembershipMsgs.MembershipOutput({ commitmentRoot: MOCK_ROOT, kvPairs: new IMembershipMsgs.KVPair[](1) });
    //     output.kvPairs[0] = IMembershipMsgs.KVPair({ path: membershipPath, value: bytes("") });

    //     IMembershipMsgs.Groth16MembershipProof memory groth16Proof = IMembershipMsgs.Groth16MembershipProof({
    //         groth16Proof: IGroth16Msgs.Groth16Proof({ vKey: MOCK_VKEY, publicValues: abi.encode(output), proof: bytes("") }),
    //         trustedConsensusState: newMockConsensusState(height)
    //     });

    //     IMembershipMsgs.MembershipProof memory proof = IMembershipMsgs.MembershipProof({
    //         proofType: IMembershipMsgs.MembershipProofType.Groth16MembershipProof,
    //         proof: abi.encode(groth16Proof)
    //     });

    //     return ILightClientMsgs.MsgVerifyNonMembership({
    //         proof: abi.encode(proof),
    //         proofHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: height }),
    //         path: membershipPath
    //     });
    // }

    // function newMisbehaviourMsg() public view returns (bytes memory) {
    //     IMisbehaviourMsgs.MisbehaviourOutput memory output = IMisbehaviourMsgs.MisbehaviourOutput({
    //         clientState: abi.decode(ics07Tendermint.getClientState(), (IICS07TendermintMsgs.ClientState)),
    //         time: uint128(block.timestamp * 1e9),
    //         trustedHeight1: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: 1 }),
    //         trustedHeight2: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: 1 }),
    //         trustedConsensusState1: newMockConsensusState(1),
    //         trustedConsensusState2: newMockConsensusState(1)
    //     });

    //     return abi.encode(
    //         IMisbehaviourMsgs.MsgSubmitMisbehaviour({
    //             groth16Proof: IGroth16Msgs.Groth16Proof({ vKey: MOCK_VKEY, publicValues: abi.encode(output), proof: bytes("") })
    //         })
    //     );
    }
}
