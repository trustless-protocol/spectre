// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { Misbehaviour } from "../../contracts/programs/Misbehaviour.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { IVerifier } from "../../contracts/interfaces/IVerifier.sol";
import { IMembership } from "../../contracts/interfaces/IMembership.sol";
import { IMisbehaviour } from "../../contracts/interfaces/IMisbehaviour.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";
import { IUpdateClient } from "../../contracts/interfaces/IUpdateClient.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { Header as HeaderLib } from "../../contracts/utils/Header.sol";
import { IGroth16ICS07TendermintErrors } from "../../contracts/light-clients/errors/IGroth16ICS07TendermintErrors.sol";
import { IAccessControl } from "@openzeppelin-contracts/access/IAccessControl.sol";

contract MockVerifierForMisbehaviour is IVerifier {
    function verifyBatchProof(
        uint16,
        uint256[8] calldata,
        uint256[2] calldata,
        uint256[2] calldata,
        bytes32[] calldata,
        uint64[] calldata,
        uint32[] calldata,
        bool[] calldata,
        IVerifier.SharedBlock calldata
    )
        external
        pure
        returns (bool)
    {
        return true;
    }
}

contract DummyMembershipForMisbehaviour is IMembership {
    function membership(
        bytes32,
        IMembershipMsgs.KVPair[] calldata,
        IMembershipMsgs.MerkleProof[] calldata
    )
        external
        pure { }
}

contract DummyUpdateClientForMisbehaviour is IUpdateClient {
    function updateClient(IUpdateClientMsgs.MsgUpdateClient calldata)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        revert("unused");
    }

    function updateClientResolved(IUpdateClientMsgs.MsgUpdateClient calldata)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        revert("unused");
    }

    function updateClientCachedCurrent(IUpdateClientMsgs.MsgUpdateClient calldata)
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        revert("unused");
    }

    function updateClientCachedCurrentWithHeaderCache(
        IUpdateClientMsgs.MsgUpdateClient calldata,
        bytes32,
        bytes32
    )
        external
        pure
        returns (IUpdateClientMsgs.UpdateClientOutput memory)
    {
        revert("unused");
    }
}

contract Groth16ICS07MisbehaviourHarness is Groth16ICS07Tendermint {
    constructor(
        address verifier,
        address membership_,
        address misbehaviour_,
        address updateClient_,
        bytes memory clientState_,
        bytes32 consensusStateHash_,
        address roleManager
    )
        Groth16ICS07Tendermint(
            verifier, membership_, misbehaviour_, updateClient_, clientState_, consensusStateHash_, roleManager
        )
    { }
}

contract MisbehaviourTest is Test, IICS07TendermintMsgs {
    Groth16ICS07MisbehaviourHarness internal lightClient;
    Misbehaviour internal misbehaviourVerifier;
    MockVerifierForMisbehaviour internal mockVerifier;
    DummyMembershipForMisbehaviour internal dummyMembership;
    DummyUpdateClientForMisbehaviour internal dummyUpdateClient;

    ClientState internal clientState_;
    ConsensusState internal consensusState_;
    bytes32 internal consensusStateHash_;
    ValidatorSet internal valset_;
    bytes32 internal valSetHash_;

    uint128 internal constant TRUSTED_TIME_NANOS = 1_700_000_000_000_000_000;
    uint128 internal constant HEADER_TIME_NANOS = 1_700_000_100_000_000_000;
    string internal constant CHAIN_ID = "test-chain-0";

    function setUp() public {
        misbehaviourVerifier = new Misbehaviour();
        mockVerifier = new MockVerifierForMisbehaviour();
        dummyMembership = new DummyMembershipForMisbehaviour();
        dummyUpdateClient = new DummyUpdateClientForMisbehaviour();

        ValidatorInfo[] memory vals = new ValidatorInfo[](4);
        vals[0] = _validator("val0", bytes32(uint256(1)), 25);
        vals[1] = _validator("val1", bytes32(uint256(2)), 25);
        vals[2] = _validator("val2", bytes32(uint256(3)), 25);
        vals[3] = _validator("val3", bytes32(uint256(4)), 25);

        valset_ = ValidatorSet({
            validators: vals, hasProposer: false, proposer: _validator("", bytes32(0), 0), totalVotingPower: 100
        });
        valSetHash_ = HeaderLib.hashValSet(valset_);

        consensusState_ = ConsensusState({
            timestamp: TRUSTED_TIME_NANOS, root: bytes32(uint256(0xABCD)), nextValidatorsHash: valSetHash_
        });
        consensusStateHash_ = keccak256(abi.encode(consensusState_));

        IICS02ClientMsgs.Height memory latestHeight = IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: 10 });

        clientState_ = ClientState({
            chainId: CHAIN_ID,
            trustLevel: TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: latestHeight,
            trustingPeriod: 3600,
            unbondingPeriod: 7200,
            isFrozen: false,
            zkAlgorithm: SupportedZkAlgorithm.Groth16
        });

        bytes memory encodedClientState = abi.encode(clientState_);

        lightClient = new Groth16ICS07MisbehaviourHarness(
            address(mockVerifier),
            address(dummyMembership),
            address(misbehaviourVerifier),
            address(dummyUpdateClient),
            encodedClientState,
            consensusStateHash_,
            address(0)
        );

        vm.warp(1_700_000_500);
    }

    function test_misbehaviour_alwaysRevertsWithFeatureNotSupported() public {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(encoded);
    }

    function test_misbehaviour_revertsWithGarbageBytes() public {
        bytes memory garbage = bytes("this is not a valid encoded message");
        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(garbage);
    }

    function test_misbehaviour_revertsWithEmptyBytes() public {
        bytes memory empty = bytes("");
        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(empty);
    }

    function test_misbehaviour_cannotFreezeClientWithValidMessage() public {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(encoded);

        (,,,,, bool isFrozen,) = lightClient.clientState();
        assertFalse(isFrozen, "client must not be frozen when misbehaviour is disabled");
    }

    function test_misbehaviour_cannotFreezeClientWithFakeSignatures() public {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();

        for (uint256 i = 0; i < msg_.misbehaviour.header1.signedHeader.commit.commitSigs.length; i++) {
            msg_.misbehaviour.header1.signedHeader.commit.commitSigs[i].data.signature = bytes("fake_signature");
        }
        for (uint256 i = 0; i < msg_.misbehaviour.header2.signedHeader.commit.commitSigs.length; i++) {
            msg_.misbehaviour.header2.signedHeader.commit.commitSigs[i].data.signature = bytes("fake_signature");
        }

        bytes memory encoded = abi.encode(msg_);
        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(encoded);

        (,,,,, bool isFrozen,) = lightClient.clientState();
        assertFalse(isFrozen, "client must not be frozen with fake signatures");
    }

    function test_misbehaviour_cannotFreezeClientWithEmptySignatures() public {
        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();

        for (uint256 i = 0; i < msg_.misbehaviour.header1.signedHeader.commit.commitSigs.length; i++) {
            msg_.misbehaviour.header1.signedHeader.commit.commitSigs[i].data.signature = bytes("");
            msg_.misbehaviour.header1.signedHeader.commit.commitSigs[i].data.hasSignature = false;
        }
        for (uint256 i = 0; i < msg_.misbehaviour.header2.signedHeader.commit.commitSigs.length; i++) {
            msg_.misbehaviour.header2.signedHeader.commit.commitSigs[i].data.signature = bytes("");
            msg_.misbehaviour.header2.signedHeader.commit.commitSigs[i].data.hasSignature = false;
        }

        bytes memory encoded = abi.encode(msg_);
        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(encoded);

        (,,,,, bool isFrozen,) = lightClient.clientState();
        assertFalse(isFrozen, "client must not be frozen with empty signatures");
    }

    function test_misbehaviour_cannotFreezeClientFromAnyAddress() public {
        address attacker = makeAddr("attacker");

        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        vm.prank(attacker);
        vm.expectRevert(IGroth16ICS07TendermintErrors.FeatureNotSupported.selector);
        lightClient.misbehaviour(encoded);

        (,,,,, bool isFrozen,) = lightClient.clientState();
        assertFalse(isFrozen, "client must not be frozen by arbitrary attacker");
    }

    function test_misbehaviour_revertsForUnauthorizedCallerWhenManaged() public {
        address governance = makeAddr("governance");
        address unauthorized = makeAddr("unauthorized");
        bytes memory encodedClientState = abi.encode(clientState_);
        Groth16ICS07MisbehaviourHarness managedClient = new Groth16ICS07MisbehaviourHarness(
            address(mockVerifier),
            address(dummyMembership),
            address(misbehaviourVerifier),
            address(dummyUpdateClient),
            encodedClientState,
            consensusStateHash_,
            governance
        );

        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        bytes32 misbehaviourRole = managedClient.MISBEHAVIOUR_SUBMITTER_ROLE();

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, misbehaviourRole
            )
        );
        managedClient.misbehaviour(encoded);

        (,,,,, bool isFrozen,) = managedClient.clientState();
        assertFalse(isFrozen, "client must not be frozen by unauthorized caller");
    }

    function test_misbehaviour_distinctRoleFromProofSubmitter() public view {
        assertNotEq(
            lightClient.PROOF_SUBMITTER_ROLE(),
            lightClient.MISBEHAVIOUR_SUBMITTER_ROLE(),
            "MISBEHAVIOUR_SUBMITTER_ROLE must be distinct from PROOF_SUBMITTER_ROLE"
        );
    }

    function test_unfreeze_revertsWhenNotFrozen() public {
        address governance = makeAddr("governance");
        bytes memory encodedClientState = abi.encode(clientState_);
        Groth16ICS07MisbehaviourHarness managedClient = new Groth16ICS07MisbehaviourHarness(
            address(mockVerifier),
            address(dummyMembership),
            address(misbehaviourVerifier),
            address(dummyUpdateClient),
            encodedClientState,
            consensusStateHash_,
            governance
        );

        vm.prank(governance);
        vm.expectRevert(IGroth16ICS07TendermintErrors.ClientNotFrozen.selector);
        managedClient.unfreeze();
    }

    function test_unfreeze_revertsForNonAdmin() public {
        address governance = makeAddr("governance");
        address unauthorized = makeAddr("unauthorized");
        bytes memory encodedClientState = abi.encode(clientState_);
        Groth16ICS07MisbehaviourHarness managedClient = new Groth16ICS07MisbehaviourHarness(
            address(mockVerifier),
            address(dummyMembership),
            address(misbehaviourVerifier),
            address(dummyUpdateClient),
            encodedClientState,
            consensusStateHash_,
            governance
        );

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, bytes32(0))
        );
        managedClient.unfreeze();
    }

    function _buildValidMisbehaviourMsg() internal view returns (IMisbehaviourMsgs.MsgSubmitMisbehaviour memory) {
        IICS07TendermintMsgs.Header memory header1 = _buildHeader(15, bytes32(uint256(0xAAA1)));
        IICS07TendermintMsgs.Header memory header2 = _buildHeader(12, bytes32(uint256(0xAAA2)));

        IMisbehaviourMsgs.Misbehaviour memory misbehaviour_ = IMisbehaviourMsgs.Misbehaviour({
            client_id: ChainId({ id: CHAIN_ID, revisionNumber: 0 }), header1: header1, header2: header2
        });

        return IMisbehaviourMsgs.MsgSubmitMisbehaviour({
            clientState: clientState_,
            misbehaviour: misbehaviour_,
            trustedConsensusState1: consensusState_,
            trustedConsensusState2: consensusState_,
            time: uint128(block.timestamp) * 1_000_000_000
        });
    }

    function _buildHeader(
        uint64 height,
        bytes32 lastBlockIdHash
    )
        internal
        view
        returns (IICS07TendermintMsgs.Header memory)
    {
        BlockHeader memory blockHeader = BlockHeader({
            version: Version({ blockVersion: 11, appVersion: 0 }),
            chainId: CHAIN_ID,
            height: height,
            time: HEADER_TIME_NANOS,
            hasLastBlockId: true,
            lastBlockId: BlockId({
                hashData: lastBlockIdHash,
                partSetHeader: PartSetHeader({ total: 1, hashData: bytes32(uint256(0x1234)) })
            }),
            hasLastCommitHash: false,
            lastCommitHash: bytes32(0),
            hasDataHash: false,
            dataHash: bytes32(0),
            validatorsHash: valSetHash_,
            nextValidatorsHash: valSetHash_,
            consensusHash: bytes32(uint256(0xCAFE)),
            appHash: bytes32(uint256(0xBEEF)),
            hasLastResultsHash: false,
            lastResultsHash: bytes32(0),
            hasEvidenceHash: false,
            evidenceHash: bytes32(0),
            proposerAddress: bytes("prop")
        });

        bytes32 headerHash = HeaderLib.hashHeader(blockHeader);

        CommitSig[] memory sigs = new CommitSig[](4);
        sigs[0] = _commitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val0", HEADER_TIME_NANOS);
        sigs[1] = _commitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val1", HEADER_TIME_NANOS);
        sigs[2] = _commitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val2", HEADER_TIME_NANOS);
        sigs[3] = _commitSig(CommitSigFlag.BLOCK_ID_FLAG_COMMIT, "val3", HEADER_TIME_NANOS);

        BlockCommit memory commit = BlockCommit({
            height: height,
            round: 0,
            blockId: BlockId({
                hashData: headerHash, partSetHeader: PartSetHeader({ total: 1, hashData: bytes32(uint256(0x5678)) })
            }),
            commitSigs: sigs
        });

        SignedHeader memory signedHeader = SignedHeader({ header: blockHeader, commit: commit });

        IICS02ClientMsgs.Height memory trustedHeight =
            IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: 10 });

        return IICS07TendermintMsgs.Header({
            signedHeader: signedHeader,
            validatorSet: valset_,
            trustedHeight: trustedHeight,
            trustedNextValidatorSet: valset_
        });
    }

    function _validator(string memory addr, bytes32 pubkey, uint64 power) internal pure returns (ValidatorInfo memory) {
        return ValidatorInfo({ valAddress: bytes(addr), pubKey: pubkey, votingPower: power, proposerPriority: 0 });
    }

    function _commitSig(
        CommitSigFlag flag,
        string memory addr,
        uint128 timestamp
    )
        internal
        pure
        returns (CommitSig memory)
    {
        return CommitSig({
            flag: flag,
            data: CommitSigData({
                validatorAddress: bytes(addr),
                timestamp: timestamp,
                hasSignature: flag != CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                signature: bytes("")
            })
        });
    }
}
