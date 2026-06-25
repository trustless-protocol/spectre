// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Groth16ICS07Tendermint } from "../../contracts/light-clients/Groth16ICS07Tendermint.sol";
import { UpdateClient } from "../../contracts/programs/UpdateClient.sol";
import { Misbehaviour } from "../../contracts/programs/Misbehaviour.sol";
import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { Header } from "../../contracts/utils/Header.sol";
import { IUpdateClientMsgs } from "../../contracts/light-clients/msgs/IUpdateClientMsgs.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";
import { ILightClientMsgs } from "../../contracts/msgs/ILightClientMsgs.sol";
import { IGroth16ICS07TendermintErrors } from "../../contracts/light-clients/errors/IGroth16ICS07TendermintErrors.sol";

contract CacheAlwaysTrueVerifier {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure returns (bool) {
        return true;
    }
}

contract CacheAlwaysFailingVerifier {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure returns (bool) {
        revert();
    }
}

contract UpdateClientCacheTest is Test {
    WrapperVerifier internal wrapper;
    UpdateClient internal updateClientImpl;
    CacheAlwaysTrueVerifier internal stubBucket;
    CacheAlwaysFailingVerifier internal failingStub;

    address constant STUB_MEMBERSHIP = address(0xBABE);
    address constant STUB_MISBEHAVIOUR = address(0xBEEF);

    string constant CHAIN_ID = "cosmoshub-0";
    uint64 constant TRUSTED_HEIGHT = 1000;
    uint64 constant HEIGHT_1001 = 1001;
    uint64 constant HEIGHT_1002 = 1002;
    uint128 constant TRUSTED_TS_NS = 1_700_000_000 * 1e9;
    uint128 constant TS_1001_NS = 1_700_000_010 * 1e9;
    uint128 constant TS_1002_NS = 1_700_000_020 * 1e9;
    uint32 constant TRUSTING_PERIOD = 14 days;
    uint32 constant UNBONDING_PERIOD = 21 days;
    uint16 constant BUCKET = 4;

    function setUp() public {
        vm.warp(1_700_000_030);

        wrapper = new WrapperVerifier(address(this));
        stubBucket = new CacheAlwaysTrueVerifier();
        failingStub = new CacheAlwaysFailingVerifier();
        updateClientImpl = new UpdateClient();
        wrapper.setBucket(BUCKET, address(stubBucket), CacheAlwaysTrueVerifier.verifyProof.selector);
    }

    function test_constructor_pinsInitialValidatorSet() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        Groth16ICS07Tendermint ics07 = _deployLightClient(_trustedConsensus(pinned), pinned);

        (uint32[] memory indices, bytes32[] memory pubkeys, uint64[] memory votingPowers) =
            ics07.getPinnedValidatorSet();
        assertEq(indices.length, pinned.validators.length, "pinned length");
        for (uint256 i = 0; i < pinned.validators.length; i++) {
            assertEq(indices[i], i, "pinned index");
            assertEq(pubkeys[i], pinned.validators[i].pubKey, "pinned pubkey");
            assertEq(votingPowers[i], pinned.validators[i].votingPower, "pinned power");
        }
    }

    function test_updateClient_acceptsPinnedQuorum() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(trustedCS, pinned, 3);
        ILightClientMsgs.UpdateResult result = ics07.updateClient(abi.encode(msg_));

        assertEq(uint8(result), uint8(ILightClientMsgs.UpdateResult.Update));
        IICS07TendermintMsgs.ClientState memory updated =
            abi.decode(ics07.getClientState(), (IICS07TendermintMsgs.ClientState));
        assertEq(updated.latestHeight.revisionHeight, HEIGHT_1001, "latest height");
    }

    function test_updateClient_reverts_whenPinnedQuorumIsInsufficient() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(trustedCS, pinned, 2);

        vm.expectRevert(
            abi.encodeWithSelector(IGroth16ICS07TendermintErrors.InsufficientVotingPower.selector, 200, 400)
        );
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_reverts_whenPinnedIndexIsDuplicated() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(trustedCS, pinned, 3);
        msg_.pinnedValidatorIndices[1] = 0;
        msg_.signerPubkeys[1] = pinned.validators[0].pubKey;

        vm.expectRevert(abi.encodeWithSelector(IGroth16ICS07TendermintErrors.DuplicateSigner.selector, 0));
        ics07.updateClient(abi.encode(msg_));
    }

    function test_updateClient_reverts_whenPinnedIndexIsOutOfRange() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(trustedCS, pinned, 3);
        msg_.pinnedValidatorIndices[0] = 4;

        vm.expectRevert(abi.encodeWithSelector(IGroth16ICS07TendermintErrors.SignerIndexOutOfRange.selector, 4));
        ics07.updateClient(abi.encode(msg_));
    }

    function test_reAnchorPinnedSet_updatesPinnedSet() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsgWithNextHash(trustedCS, pinnedA, Header.hashValSet(pinnedB), 3);
        ics07.reAnchorPinnedSet(msg_, pinnedB);

        (, bytes32[] memory pubkeys,) = ics07.getPinnedValidatorSet();
        assertEq(pubkeys[0], pinnedB.validators[0].pubKey, "re-anchored pubkey");
    }

    function test_reAnchorPinnedSet_reverts_whenNewSetHashDoesNotMatchHeader() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        bytes32 expected = Header.hashValSet(pinnedA);
        bytes32 actual = Header.hashValSet(pinnedB);
        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsgWithNextHash(trustedCS, pinnedA, expected, 3);

        vm.expectRevert(
            abi.encodeWithSelector(IGroth16ICS07TendermintErrors.MismatchedValidatorHashes.selector, expected, actual)
        );
        ics07.reAnchorPinnedSet(msg_, pinnedB);
    }

    function test_reAnchorPinnedSet_doesNotUpdatePinnedSetWhenHeaderIsMisbehaviour() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        IUpdateClientMsgs.MsgUpdateClient memory good = _buildMsg(trustedCS, pinnedA, 3);
        ics07.updateClient(abi.encode(good));

        IUpdateClientMsgs.MsgUpdateClient memory conflicting =
            _buildMsgWithNextHash(trustedCS, pinnedA, Header.hashValSet(pinnedB), 3);
        ics07.reAnchorPinnedSet(conflicting, pinnedB);

        (, bytes32[] memory pubkeys,) = ics07.getPinnedValidatorSet();
        assertEq(pubkeys[0], pinnedA.validators[0].pubKey, "misbehaviour re-anchor must not mutate pin");
        IICS07TendermintMsgs.ClientState memory updated =
            abi.decode(ics07.getClientState(), (IICS07TendermintMsgs.ClientState));
        assertTrue(updated.isFrozen, "misbehaviour still freezes");
    }

    function test_reAnchorPinnedSet_reverts_whenNoOpHeaderIsNotLatestHeight() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ValidatorSet memory pinnedC = _buildValSet(4, 0);
        pinnedC.validators[3].pubKey = bytes32(uint256(0xC0FFEE));
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        IUpdateClientMsgs.MsgUpdateClient memory reanchorB =
            _buildMsgWithNextHash(trustedCS, pinnedA, Header.hashValSet(pinnedB), 3);
        ics07.reAnchorPinnedSet(reanchorB, pinnedB);
        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 = _consensusFromMsg(reanchorB);

        IUpdateClientMsgs.MsgUpdateClient memory reanchorC =
            _buildMsgAt(HEIGHT_1001, HEIGHT_1002, trustedCS1001, pinnedB, Header.hashValSet(pinnedC), TS_1002_NS, 3);
        ics07.reAnchorPinnedSet(reanchorC, pinnedC);

        vm.expectRevert(
            abi.encodeWithSelector(
                IGroth16ICS07TendermintErrors.NonMonotonicHeightUpdate.selector, HEIGHT_1002, HEIGHT_1001
            )
        );
        ics07.reAnchorPinnedSet(reanchorB, pinnedB);
    }

    function test_updateClient_reverts_whenProofVerificationFails() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        wrapper.setBucket(BUCKET, address(failingStub), CacheAlwaysFailingVerifier.verifyProof.selector);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ = _buildMsg(trustedCS, pinned, 3);
        vm.expectRevert(IGroth16ICS07TendermintErrors.ProofVerificationFailed.selector);
        ics07.updateClient(abi.encode(msg_));
    }

    function test_reAnchorPinnedSet_reverts_whenProofVerificationFails() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        wrapper.setBucket(BUCKET, address(failingStub), CacheAlwaysFailingVerifier.verifyProof.selector);

        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsgWithNextHash(trustedCS, pinnedA, Header.hashValSet(pinnedB), 3);
        vm.expectRevert(IGroth16ICS07TendermintErrors.ProofVerificationFailed.selector);
        ics07.reAnchorPinnedSet(msg_, pinnedB);
    }

    function test_updateClient_reverts_whenTrustedHeightNotStored() public {
        IICS07TendermintMsgs.ValidatorSet memory pinned = _buildValSet(4, 0);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinned);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinned);

        uint64 unstoredHeight = TRUSTED_HEIGHT - 1;
        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsgAt(unstoredHeight, HEIGHT_1001, trustedCS, pinned, Header.hashValSet(pinned), TS_1001_NS, 3);
        vm.expectRevert(IGroth16ICS07TendermintErrors.ConsensusStateNotFound.selector);
        ics07.updateClient(abi.encode(msg_));
    }

    function test_reAnchorPinnedSet_reverts_whenTrustedHeightNotStored() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClient(trustedCS, pinnedA);

        uint64 unstoredHeight = TRUSTED_HEIGHT - 1;
        IUpdateClientMsgs.MsgUpdateClient memory msg_ =
            _buildMsgAt(unstoredHeight, HEIGHT_1001, trustedCS, pinnedA, Header.hashValSet(pinnedB), TS_1001_NS, 3);
        vm.expectRevert(IGroth16ICS07TendermintErrors.ConsensusStateNotFound.selector);
        ics07.reAnchorPinnedSet(msg_, pinnedB);
    }

    function test_misbehaviourUsesHistoricalPinnedSetAfterReAnchor() public {
        IICS07TendermintMsgs.ValidatorSet memory pinnedA = _buildValSet(4, 0);
        IICS07TendermintMsgs.ValidatorSet memory pinnedB = _buildValSet(4, 100);
        IICS07TendermintMsgs.ConsensusState memory trustedCS = _trustedConsensus(pinnedA);
        Groth16ICS07Tendermint ics07 = _deployLightClientWithMisbehaviour(trustedCS, pinnedA);

        IUpdateClientMsgs.MsgUpdateClient memory reanchorB =
            _buildMsgWithNextHash(trustedCS, pinnedA, Header.hashValSet(pinnedB), 3);
        ics07.reAnchorPinnedSet(reanchorB, pinnedB);
        IICS07TendermintMsgs.ConsensusState memory trustedCS1001 = _consensusFromMsg(reanchorB);

        IUpdateClientMsgs.MsgUpdateClient memory update1002 =
            _buildMsgAt(HEIGHT_1001, HEIGHT_1002, trustedCS1001, pinnedB, Header.hashValSet(pinnedB), TS_1002_NS, 3);
        ics07.updateClient(abi.encode(update1002));
        IICS07TendermintMsgs.ConsensusState memory trustedCS1002 = _consensusFromMsg(update1002);

        uint64 height1003 = HEIGHT_1002 + 1;
        uint128 ts1003Ns = TS_1002_NS + 10 * 1e9;

        IICS07TendermintMsgs.Header memory h1 =
            _buildHeader(HEIGHT_1002, height1003, pinnedB, Header.hashValSet(pinnedB), ts1003Ns, 3);
        IICS07TendermintMsgs.Header memory h2 =
            _buildHeader(HEIGHT_1002, height1003, pinnedB, Header.hashValSet(pinnedA), ts1003Ns, 3);

        IMisbehaviourMsgs.MsgSubmitMisbehaviour memory msg_;
        msg_.clientState = _clientState();
        msg_.misbehaviour = IMisbehaviourMsgs.Misbehaviour({
            client_id: IICS07TendermintMsgs.ChainId({ id: CHAIN_ID, revisionNumber: 0 }), header1: h1, header2: h2
        });
        msg_.trustedConsensusState1 = trustedCS1002;
        msg_.trustedConsensusState2 = trustedCS1002;
        msg_.time = ts1003Ns;
        msg_.proof1 = _buildMisbehaviourProof(pinnedB, 3);
        msg_.proof2 = _buildMisbehaviourProof(pinnedB, 3);

        ics07.misbehaviour(abi.encode(msg_));
        IICS07TendermintMsgs.ClientState memory updated =
            abi.decode(ics07.getClientState(), (IICS07TendermintMsgs.ClientState));
        assertTrue(updated.isFrozen, "historical pinned evidence freezes");
    }

    function _deployLightClient(
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.ValidatorSet memory pinned
    )
        internal
        returns (Groth16ICS07Tendermint)
    {
        return new Groth16ICS07Tendermint(
            address(wrapper),
            STUB_MEMBERSHIP,
            STUB_MISBEHAVIOUR,
            address(updateClientImpl),
            abi.encode(_clientState()),
            keccak256(abi.encode(trustedCS)),
            pinned,
            address(0)
        );
    }

    function _deployLightClientWithMisbehaviour(
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.ValidatorSet memory pinned
    )
        internal
        returns (Groth16ICS07Tendermint)
    {
        Misbehaviour misbehaviourImpl = new Misbehaviour();
        return new Groth16ICS07Tendermint(
            address(wrapper),
            STUB_MEMBERSHIP,
            address(misbehaviourImpl),
            address(updateClientImpl),
            abi.encode(_clientState()),
            keccak256(abi.encode(trustedCS)),
            pinned,
            address(0)
        );
    }

    function _buildMsg(
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.ValidatorSet memory pinned,
        uint16 activeCount
    )
        internal
        pure
        returns (IUpdateClientMsgs.MsgUpdateClient memory)
    {
        return _buildMsgWithNextHash(trustedCS, pinned, Header.hashValSet(pinned), activeCount);
    }

    function _buildMsgWithNextHash(
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.ValidatorSet memory pinned,
        bytes32 nextValidatorsHash,
        uint16 activeCount
    )
        internal
        pure
        returns (IUpdateClientMsgs.MsgUpdateClient memory msg_)
    {
        return _buildMsgAt(TRUSTED_HEIGHT, HEIGHT_1001, trustedCS, pinned, nextValidatorsHash, TS_1001_NS, activeCount);
    }

    function _buildMsgAt(
        uint64 trustedHeight,
        uint64 newHeight,
        IICS07TendermintMsgs.ConsensusState memory trustedCS,
        IICS07TendermintMsgs.ValidatorSet memory pinned,
        bytes32 nextValidatorsHash,
        uint128 headerTime,
        uint16 activeCount
    )
        internal
        pure
        returns (IUpdateClientMsgs.MsgUpdateClient memory msg_)
    {
        IICS07TendermintMsgs.Header memory header =
            _buildHeader(trustedHeight, newHeight, pinned, nextValidatorsHash, headerTime, activeCount);

        uint32[] memory idx = new uint32[](BUCKET);
        uint32[] memory pinnedIdx = new uint32[](BUCKET);
        bytes32[] memory pks = new bytes32[](BUCKET);
        bool[] memory act = new bool[](BUCKET);
        for (uint256 i = 0; i < BUCKET; i++) {
            idx[i] = uint32(i);
            pinnedIdx[i] = uint32(i);
            if (i < activeCount) {
                pks[i] = pinned.validators[i].pubKey;
                act[i] = true;
            }
        }

        msg_.clientState = _clientState();
        msg_.trustedConsensusState = trustedCS;
        msg_.proposedHeader = header;
        msg_.time = TS_1002_NS;
        msg_.proof = [uint256(0), 0, 0, 0, 0, 0, 0, 0];
        msg_.commitments = [uint256(0), 0];
        msg_.commitmentPok = [uint256(0), 0];
        msg_.bucket = BUCKET;
        msg_.signerIndices = idx;
        msg_.pinnedValidatorIndices = pinnedIdx;
        msg_.signerPubkeys = pks;
        msg_.active = act;
    }

    function _buildMisbehaviourProof(
        IICS07TendermintMsgs.ValidatorSet memory pinned,
        uint16 activeCount
    )
        internal
        pure
        returns (IMisbehaviourMsgs.BatchProof memory proof_)
    {
        proof_.proof = [uint256(0), 0, 0, 0, 0, 0, 0, 0];
        proof_.commitments = [uint256(0), 0];
        proof_.commitmentPok = [uint256(0), 0];
        proof_.bucket = BUCKET;
        proof_.signerIndices = new uint32[](BUCKET);
        proof_.pinnedValidatorIndices = new uint32[](BUCKET);
        proof_.signerPubkeys = new bytes32[](BUCKET);
        proof_.active = new bool[](BUCKET);
        for (uint256 i = 0; i < BUCKET; i++) {
            proof_.signerIndices[i] = uint32(i);
            proof_.pinnedValidatorIndices[i] = uint32(i);
            if (i < activeCount) {
                proof_.signerPubkeys[i] = pinned.validators[i].pubKey;
                proof_.active[i] = true;
            }
        }
    }

    function _consensusFromMsg(IUpdateClientMsgs.MsgUpdateClient memory msg_)
        internal
        pure
        returns (IICS07TendermintMsgs.ConsensusState memory)
    {
        return IICS07TendermintMsgs.ConsensusState({
            timestamp: msg_.proposedHeader.signedHeader.header.time,
            root: msg_.proposedHeader.signedHeader.header.appHash,
            nextValidatorsHash: msg_.proposedHeader.signedHeader.header.nextValidatorsHash
        });
    }

    function _buildHeader(
        uint64 trustedHeight,
        uint64 newHeight,
        IICS07TendermintMsgs.ValidatorSet memory currentValSet,
        bytes32 nextValidatorsHash,
        uint128 headerTime,
        uint16 activeCount
    )
        internal
        pure
        returns (IICS07TendermintMsgs.Header memory header)
    {
        IICS07TendermintMsgs.BlockHeader memory bh;
        bh.chainId = CHAIN_ID;
        bh.height = newHeight;
        bh.time = headerTime;
        bh.appHash = bytes32(uint256(0xCC0000 + newHeight));
        bh.validatorsHash = Header.hashValSet(currentValSet);
        bh.nextValidatorsHash = nextValidatorsHash;
        bytes32 headerHash = Header.hashHeader(bh);

        IICS07TendermintMsgs.BlockCommit memory bc = IICS07TendermintMsgs.BlockCommit({
            height: newHeight,
            round: 0,
            blockId: IICS07TendermintMsgs.BlockId({
                hashData: headerHash,
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: bytes32(uint256(0x9A57)) })
            }),
            commitSigs: _buildCommitSigs(currentValSet, activeCount, headerTime)
        });

        header = IICS07TendermintMsgs.Header({
            signedHeader: IICS07TendermintMsgs.SignedHeader({ header: bh, commit: bc }),
            trustedHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: trustedHeight })
        });
    }

    function _buildCommitSigs(
        IICS07TendermintMsgs.ValidatorSet memory vs,
        uint16 activeCount,
        uint128 headerTime
    )
        internal
        pure
        returns (IICS07TendermintMsgs.CommitSig[] memory sigs)
    {
        sigs = new IICS07TendermintMsgs.CommitSig[](vs.validators.length);
        for (uint256 i = 0; i < vs.validators.length; i++) {
            if (i < activeCount) {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: vs.validators[i].valAddress,
                        timestamp: headerTime,
                        hasSignature: false,
                        signature: ""
                    })
                });
            } else {
                sigs[i] = IICS07TendermintMsgs.CommitSig({
                    flag: IICS07TendermintMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT,
                    data: IICS07TendermintMsgs.CommitSigData({
                        validatorAddress: "", timestamp: 0, hasSignature: false, signature: ""
                    })
                });
            }
        }
    }

    function _buildValSet(
        uint16 valCount,
        uint256 seed
    )
        internal
        pure
        returns (IICS07TendermintMsgs.ValidatorSet memory vs)
    {
        IICS07TendermintMsgs.ValidatorInfo[] memory vals = new IICS07TendermintMsgs.ValidatorInfo[](valCount);
        uint64 total = 0;
        for (uint256 i = 0; i < valCount; i++) {
            vals[i] = IICS07TendermintMsgs.ValidatorInfo({
                valAddress: abi.encodePacked(uint160(seed + i + 1)),
                pubKey: bytes32(uint256(0xA0000000 + seed + i + 1)),
                votingPower: 100,
                proposerPriority: 0
            });
            total += 100;
        }
        vs = IICS07TendermintMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _trustedConsensus(IICS07TendermintMsgs.ValidatorSet memory pinned)
        internal
        pure
        returns (IICS07TendermintMsgs.ConsensusState memory)
    {
        return IICS07TendermintMsgs.ConsensusState({
            timestamp: TRUSTED_TS_NS, root: bytes32(uint256(0xAAA1)), nextValidatorsHash: Header.hashValSet(pinned)
        });
    }

    function _clientState() internal pure returns (IICS07TendermintMsgs.ClientState memory) {
        return IICS07TendermintMsgs.ClientState({
            chainId: CHAIN_ID,
            trustLevel: IICS07TendermintMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: IICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT }),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            zkAlgorithm: IICS07TendermintMsgs.SupportedZkAlgorithm.Groth16,
            clockDrift: 1800
        });
    }
}
