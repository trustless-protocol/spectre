// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { SpectreClient } from "../../contracts/light-clients/SpectreClient.sol";
import { Misbehaviour } from "../../contracts/light-clients/modules/Misbehaviour.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { ISpectreClientMsgs } from "../../contracts/light-clients/msgs/ISpectreClientMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { ISignatureVerifier } from "../../contracts/light-clients/interfaces/ISignatureVerifier.sol";
import { IMembership } from "../../contracts/light-clients/interfaces/IMembership.sol";
import { IUpdateClient } from "../../contracts/light-clients/interfaces/IUpdateClient.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";
import { Header as HeaderLib } from "../../contracts/utils/Header.sol";
import { ISpectreClientErrors } from "../../contracts/light-clients/errors/ISpectreClientErrors.sol";
import { IAccessControl } from "@openzeppelin-contracts/access/IAccessControl.sol";

contract MockVerifierForMisbehaviour is ISignatureVerifier {
    bool internal result = true;
    uint256 public calls;

    function setResult(bool result_) external {
        result = result_;
    }

    function verifyBatchProof(
        uint16,
        uint256[8] calldata,
        uint256[2] calldata,
        uint256[2] calldata,
        bytes32[] calldata,
        bool[] calldata,
        ISignatureVerifier.SharedBlock calldata
    )
        external
        returns (bool)
    {
        calls++;
        return result;
    }
}

contract DummyMembershipForMisbehaviour is IMembership {
    function verifyMembership(
        bytes32,
        IMembershipMsgs.KVPair[] calldata,
        IMembershipMsgs.MerkleProof[] calldata
    )
        external
        pure { }
}

contract DummyUpdateClientForMisbehaviour is IUpdateClient {
    function verifyHeader(ISpectreClientMsgs.MsgUpdateApplicationState calldata)
        external
        pure
        returns (ISpectreClientMsgs.VerifyHeaderOutput memory)
    {
        revert("unused");
    }
}

contract MisbehaviourTest is Test, IICS07TendermintMsgs {
    SpectreClient internal lightClient;
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

    struct LegacyMsgSubmitMisbehaviour {
        IICS07TendermintMsgs.ClientState clientState;
        ISpectreClientMsgs.Misbehaviour misbehaviour;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState1;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState2;
        uint128 time;
    }

    function setUp() public {
        mockVerifier = new MockVerifierForMisbehaviour();
        misbehaviourVerifier = new Misbehaviour(address(mockVerifier));
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
            zkAlgorithm: SupportedZkAlgorithm.Groth16,
            clockDrift: 15
        });

        lightClient = _deploy(address(0));

        vm.warp(1_700_000_500);
    }

    function _deploy(address roleManager) internal returns (SpectreClient) {
        return new SpectreClient(
            address(dummyUpdateClient),
            address(dummyMembership),
            address(misbehaviourVerifier),
            abi.encode(clientState_),
            consensusStateHash_,
            valset_,
            roleManager
        );
    }

    function test_misbehaviour_freezesWithProofBackedQuorums() public {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        lightClient.misbehaviour(encoded);

        assertTrue(_isFrozen(lightClient), "client must freeze after both header proofs verify");
        assertEq(mockVerifier.calls(), 2, "both header proofs must be verified");
    }

    function test_misbehaviour_revertsWithGarbageBytes() public {
        bytes memory garbage = bytes("this is not a valid encoded message");
        vm.expectRevert();
        lightClient.misbehaviour(garbage);
        assertFalse(_isFrozen(lightClient), "garbage bytes must not freeze");
    }

    function test_misbehaviour_revertsWithEmptyBytes() public {
        bytes memory empty = bytes("");
        vm.expectRevert();
        lightClient.misbehaviour(empty);
        assertFalse(_isFrozen(lightClient), "empty bytes must not freeze");
    }

    function test_misbehaviour_cannotFreezeLegacyMessageWithoutProofs() public {
        bytes memory encoded = _buildLegacyMisbehaviourMsg();
        vm.expectRevert();
        lightClient.misbehaviour(encoded);

        assertFalse(_isFrozen(lightClient), "legacy no-proof message must not freeze");
        assertEq(mockVerifier.calls(), 0, "legacy message must not reach verifier");
    }

    function test_misbehaviour_cannotFreezeWhenVerifierRejectsProof() public {
        mockVerifier.setResult(false);
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        vm.expectRevert(ISpectreClientErrors.ProofVerificationFailed.selector);
        lightClient.misbehaviour(encoded);

        assertFalse(_isFrozen(lightClient), "rejected proof must not freeze");
    }

    function test_misbehaviour_cannotFreezeWithSubQuorumProofMetadata() public {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.proof1 = _proof(2);

        bytes memory encoded = abi.encode(msg_);
        vm.expectRevert(abi.encodeWithSelector(ISpectreClientErrors.InsufficientVotingPower.selector, 50, 100));
        lightClient.misbehaviour(encoded);

        assertFalse(_isFrozen(lightClient), "sub-quorum proof metadata must not freeze");
        assertEq(mockVerifier.calls(), 0, "sub-quorum metadata must not reach verifier");
    }

    function test_misbehaviour_cannotFreezeWithStaleTime() public {
        // Long-range freeze attempt: `time` is consistent with the (old) trusted state's trusting
        // period but lags the wall clock beyond clockDrift. Must revert before any other check.
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.time = 1_700_000_200 * 1e9; // wall clock is 1_700_000_500, clockDrift is 15s

        vm.expectRevert(
            abi.encodeWithSelector(ISpectreClientErrors.ProofIsTooOld.selector, 1_700_000_500, 1_700_000_200)
        );
        lightClient.misbehaviour(abi.encode(msg_));

        assertFalse(_isFrozen(lightClient), "stale time must not freeze");
        assertEq(mockVerifier.calls(), 0, "stale time must not reach verifier");
    }

    function test_misbehaviour_cannotFreezeWithFutureTime() public {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.time = 1_700_000_501 * 1e9; // one second ahead of the wall clock

        vm.expectRevert(
            abi.encodeWithSelector(ISpectreClientErrors.ProofIsInTheFuture.selector, 1_700_000_500, 1_700_000_501)
        );
        lightClient.misbehaviour(abi.encode(msg_));

        assertFalse(_isFrozen(lightClient), "future time must not freeze");
        assertEq(mockVerifier.calls(), 0, "future time must not reach verifier");
    }

    function test_misbehaviour_rejectsSpoofedTrustedOverlapForNonAdjacentHeaders() public {
        ValidatorSet memory attackerValSet = _attackerValidatorSet();
        ValidatorSet memory spoofedTrustedValSet = _trustedValidatorSetWithSpoofedAddresses();

        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.misbehaviour.header1 =
            _buildHeaderWithSets(15, bytes32(uint256(0xAAA1)), attackerValSet, spoofedTrustedValSet);
        msg_.misbehaviour.header2 =
            _buildHeaderWithSets(15, bytes32(uint256(0xAAA2)), attackerValSet, spoofedTrustedValSet);
        msg_.proof1 = _proofForValidatorSet(attackerValSet, 3);
        msg_.proof2 = _proofForValidatorSet(attackerValSet, 3);

        vm.expectRevert(abi.encodeWithSelector(ISpectreClientErrors.PubkeyMismatch.selector, uint32(0)));
        lightClient.misbehaviour(abi.encode(msg_));

        assertFalse(_isFrozen(lightClient), "spoofed trusted overlap must not freeze");
    }

    function test_misbehaviour_rejectsProofSignerAbsentFromCommitSigs() public {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.misbehaviour.header1.signedHeader.commit.commitSigs[0] =
            _commitSig(CommitSigFlag.BLOCK_ID_FLAG_ABSENT, "", 0);

        vm.expectRevert(abi.encodeWithSelector(ISpectreClientErrors.ProofSignerCommitSigMismatch.selector, uint32(0)));
        lightClient.misbehaviour(abi.encode(msg_));

        assertFalse(_isFrozen(lightClient), "proof signer absent from commit must not freeze");
        assertEq(mockVerifier.calls(), 0, "mismatched commit metadata must not reach verifier");
    }

    function test_misbehaviour_revertsForMismatchedHeaderHeights() public {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        msg_.misbehaviour.header2 = _buildHeader(16, bytes32(uint256(0xAAA2)));

        bytes memory encoded = abi.encode(msg_);
        vm.expectRevert(
            abi.encodeWithSelector(ISpectreClientErrors.MismatchedMisbehaviourHeaderHeights.selector, 15, 16)
        );
        lightClient.misbehaviour(encoded);

        assertFalse(_isFrozen(lightClient), "different-height headers must not freeze");
        assertEq(mockVerifier.calls(), 0, "different-height headers must not reach verifier");
    }

    function test_misbehaviour_revertsForUnauthorizedCallerWhenManaged() public {
        address governance = makeAddr("governance");
        address unauthorized = makeAddr("unauthorized");
        SpectreClient managedClient = _deploy(governance);

        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        bytes memory encoded = abi.encode(msg_);

        bytes32 misbehaviourRole = managedClient.MISBEHAVIOUR_SUBMITTER_ROLE();

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(
                IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, misbehaviourRole
            )
        );
        managedClient.misbehaviour(encoded);

        assertFalse(_isFrozen(managedClient), "client must not be frozen by unauthorized caller");
    }

    function test_misbehaviour_distinctRoleFromProofSubmitter() public view {
        assertNotEq(
            keccak256("PROOF_SUBMITTER_ROLE"),
            lightClient.MISBEHAVIOUR_SUBMITTER_ROLE(),
            "MISBEHAVIOUR_SUBMITTER_ROLE must be distinct from PROOF_SUBMITTER_ROLE"
        );
    }

    function test_unfreeze_revertsWhenNotFrozen() public {
        address governance = makeAddr("governance");
        SpectreClient managedClient = _deploy(governance);

        vm.prank(governance);
        vm.expectRevert(ISpectreClientErrors.ClientNotFrozen.selector);
        managedClient.unfreeze();
    }

    function test_unfreeze_revertsForNonAdmin() public {
        address governance = makeAddr("governance");
        address unauthorized = makeAddr("unauthorized");
        SpectreClient managedClient = _deploy(governance);

        vm.prank(unauthorized);
        vm.expectRevert(
            abi.encodeWithSelector(IAccessControl.AccessControlUnauthorizedAccount.selector, unauthorized, bytes32(0))
        );
        managedClient.unfreeze();
    }

    function _buildValidMisbehaviourMsg() internal view returns (ISpectreClientMsgs.MsgSubmitMisbehaviour memory) {
        IICS07TendermintMsgs.Header memory header1 = _buildHeader(15, bytes32(uint256(0xAAA1)));
        IICS07TendermintMsgs.Header memory header2 = _buildHeader(15, bytes32(uint256(0xAAA2)));

        ISpectreClientMsgs.Misbehaviour memory misbehaviour_ = ISpectreClientMsgs.Misbehaviour({
            clientId: ChainId({ id: CHAIN_ID, revisionNumber: 0 }), header1: header1, header2: header2
        });

        return ISpectreClientMsgs.MsgSubmitMisbehaviour({
            misbehaviour: misbehaviour_,
            trustedConsensusState1: consensusState_,
            trustedConsensusState2: consensusState_,
            time: uint128(block.timestamp) * 1_000_000_000,
            proof1: _proof(3),
            proof2: _proof(3)
        });
    }

    function _buildLegacyMisbehaviourMsg() internal view returns (bytes memory) {
        ISpectreClientMsgs.MsgSubmitMisbehaviour memory msg_ = _buildValidMisbehaviourMsg();
        return abi.encode(
            LegacyMsgSubmitMisbehaviour({
                clientState: clientState_,
                misbehaviour: msg_.misbehaviour,
                trustedConsensusState1: msg_.trustedConsensusState1,
                trustedConsensusState2: msg_.trustedConsensusState2,
                time: msg_.time
            })
        );
    }

    function _proof(uint256 activeCount) internal view returns (ISpectreClientMsgs.BatchProof memory proof_) {
        return _proofForValidatorSet(valset_, activeCount);
    }

    function _proofForValidatorSet(
        ValidatorSet memory validatorSet,
        uint256 activeCount
    )
        internal
        pure
        returns (ISpectreClientMsgs.BatchProof memory proof_)
    {
        proof_.bucket = 4;
        proof_.signerIndices = new uint32[](4);
        proof_.pinnedValidatorIndices = new uint32[](4);
        proof_.signerPubkeys = new bytes32[](4);
        proof_.active = new bool[](4);

        for (uint32 i = 0; i < 4; i++) {
            proof_.signerIndices[i] = i;
            proof_.pinnedValidatorIndices[i] = i;
            proof_.signerPubkeys[i] = validatorSet.validators[i].pubKey;
            proof_.active[i] = i < activeCount;
        }
    }

    function _isFrozen(SpectreClient client) internal view returns (bool) {
        ClientState memory state = abi.decode(client.getClientState(), (ClientState));
        return state.isFrozen;
    }

    function _buildHeader(
        uint64 height,
        bytes32 lastBlockIdHash
    )
        internal
        view
        returns (IICS07TendermintMsgs.Header memory)
    {
        return _buildHeaderWithSets(height, lastBlockIdHash, valset_, valset_);
    }

    function _buildHeaderWithSets(
        uint64 height,
        bytes32 lastBlockIdHash,
        ValidatorSet memory currentValSet,
        ValidatorSet memory trustedNextValSet
    )
        internal
        pure
        returns (IICS07TendermintMsgs.Header memory)
    {
        bytes32 currentValSetHash = HeaderLib.hashValSet(currentValSet);
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
            validatorsHash: currentValSetHash,
            nextValidatorsHash: currentValSetHash,
            consensusHash: bytes32(uint256(0xCAFE)),
            appHash: bytes32(uint256(0xBEEF)),
            hasLastResultsHash: false,
            lastResultsHash: bytes32(0),
            hasEvidenceHash: false,
            evidenceHash: bytes32(0),
            proposerAddress: bytes("prop")
        });

        bytes32 headerHash = HeaderLib.hashHeader(blockHeader);

        BlockCommit memory commit = BlockCommit({
            height: height,
            round: 0,
            blockId: BlockId({
                hashData: headerHash, partSetHeader: PartSetHeader({ total: 1, hashData: bytes32(uint256(0x5678)) })
            }),
            commitSigs: _commitSigsForValidatorSet(currentValSet)
        });

        SignedHeader memory signedHeader = SignedHeader({ header: blockHeader, commit: commit });

        IICS02ClientMsgs.Height memory trustedHeight =
            IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: 10 });

        return IICS07TendermintMsgs.Header({ signedHeader: signedHeader, trustedHeight: trustedHeight });
    }

    function _attackerValidatorSet() internal pure returns (ValidatorSet memory validatorSet) {
        ValidatorInfo[] memory vals = new ValidatorInfo[](4);
        vals[0] = _validator("spoof0", bytes32(uint256(0xA1)), 25);
        vals[1] = _validator("spoof1", bytes32(uint256(0xA2)), 25);
        vals[2] = _validator("spoof2", bytes32(uint256(0xA3)), 25);
        vals[3] = _validator("spoof3", bytes32(uint256(0xA4)), 25);

        validatorSet = ValidatorSet({
            validators: vals, hasProposer: false, proposer: _validator("", bytes32(0), 0), totalVotingPower: 100
        });
    }

    function _trustedValidatorSetWithSpoofedAddresses() internal view returns (ValidatorSet memory validatorSet) {
        ValidatorInfo[] memory vals = new ValidatorInfo[](4);
        vals[0] = _validator("spoof0", valset_.validators[0].pubKey, 25);
        vals[1] = _validator("spoof1", valset_.validators[1].pubKey, 25);
        vals[2] = _validator("spoof2", valset_.validators[2].pubKey, 25);
        vals[3] = _validator("spoof3", valset_.validators[3].pubKey, 25);

        validatorSet = ValidatorSet({
            validators: vals, hasProposer: false, proposer: _validator("", bytes32(0), 0), totalVotingPower: 100
        });
    }

    function _commitSigsForValidatorSet(ValidatorSet memory validatorSet)
        internal
        pure
        returns (CommitSig[] memory sigs)
    {
        sigs = new CommitSig[](validatorSet.validators.length);
        for (uint256 i = 0; i < validatorSet.validators.length; i++) {
            sigs[i] = CommitSig({
                flag: CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
                data: CommitSigData({
                    validatorAddress: validatorSet.validators[i].valAddress,
                    timestamp: HEADER_TIME_NANOS,
                    hasSignature: true,
                    signature: bytes("")
                })
            });
        }
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
