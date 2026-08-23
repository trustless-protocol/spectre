// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { IntegrationTest } from "test/integration/IntegrationTest.t.sol";
import { SpectreClient } from "contracts/light-clients/spectre/SpectreClient.sol";
import { SignatureVerifier } from "contracts/light-clients/spectre/SignatureVerifier.sol";
import { UpdateClient } from "contracts/light-clients/spectre/modules/UpdateClient.sol";
import { Header } from "contracts/light-clients/spectre/libraries/Header.sol";
import { SpectreClientMsgs } from "contracts/light-clients/spectre/messages/SpectreClientMsgs.sol";
import { SpectreMsgs } from "contracts/light-clients/spectre/messages/SpectreMsgs.sol";
import { ICS02ClientMsgs } from "contracts/core/messages/ICS02ClientMsgs.sol";
import { LightClientMsgs } from "contracts/light-clients/messages/LightClientMsgs.sol";
import { ICS26RouterMsgs } from "contracts/core/messages/ICS26RouterMsgs.sol";
import { ICS20TransferMsgs } from "contracts/apps/ics20/messages/ICS20TransferMsgs.sol";
import { ICS20Lib } from "contracts/apps/ics20/libraries/ICS20Lib.sol";
import { Strings } from "@openzeppelin-contracts/utils/Strings.sol";
import { console } from "forge-std/Test.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";
import { ICS24Host } from "contracts/core/libraries/ICS24Host.sol";

contract DummyMembership {
    function verifyMembership(
        bytes32,
        MembershipMsgs.KVPair[] calldata,
        MembershipMsgs.MerkleProof[] calldata
    )
        external
        pure { }
}

contract AlwaysTrueVerifier {
    function verifyProof(bytes calldata, uint256[2] calldata) external pure returns (bool) {
        return true;
    }
}

contract RecvPacketGasTest is IntegrationTest {
    SignatureVerifier wrapper;
    UpdateClient updateClientImpl;
    AlwaysTrueVerifier stubBucket;

    address constant STUB_MEMBERSHIP = address(0xBABE);
    address constant STUB_MISBEHAVIOUR = address(0xBEEF);

    string constant CHAIN_ID = "cosmoshub-0";
    uint64 constant TRUSTED_HEIGHT = 1000;
    uint64 constant NEW_HEIGHT = 1001;
    uint128 constant TRUSTED_TS_NS = 1_700_000_000 * 1e9;
    uint128 constant NEW_TS_NS = 1_700_000_010 * 1e9;
    uint32 constant TRUSTING_PERIOD = 14 days;
    uint32 constant UNBONDING_PERIOD = 21 days;

    struct BucketConfig {
        uint16 bucket; // padded slot count
        uint16 valCount; // validators
        uint16 activeCount; // real signers
    }

    function _cfg(uint16 bucket) internal pure returns (BucketConfig memory) {
        if (bucket == 4) return BucketConfig(4, 4, 3);
        if (bucket == 8) return BucketConfig(8, 10, 7);
        if (bucket == 16) return BucketConfig(16, 20, 14);
        if (bucket == 32) return BucketConfig(32, 30, 21);
        if (bucket == 64) return BucketConfig(64, 60, 42);
        if (bucket == 128) return BucketConfig(128, 120, 84);
        revert("unknown bucket");
    }

    function _clientState() internal pure returns (SpectreMsgs.ClientState memory) {
        return SpectreMsgs.ClientState({
            chainId: CHAIN_ID,
            trustLevel: SpectreMsgs.TrustThreshold({ numerator: 1, denominator: 3 }),
            latestHeight: ICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT }),
            trustingPeriod: TRUSTING_PERIOD,
            unbondingPeriod: UNBONDING_PERIOD,
            isFrozen: false,
            clockDrift: 1800
        });
    }

    function _buildValSet(uint16 valCount) internal pure returns (SpectreMsgs.ValidatorSet memory vs) {
        SpectreMsgs.ValidatorInfo[] memory vals = new SpectreMsgs.ValidatorInfo[](valCount);
        uint64 total = 0;
        for (uint256 i = 0; i < valCount; i++) {
            vals[i] = SpectreMsgs.ValidatorInfo({
                valAddress: abi.encodePacked(uint160(i + 1)),
                pubKey: bytes32(uint256(0xA0000000 + i + 1)),
                votingPower: 100,
                proposerPriority: 0
            });
            total += 100;
        }
        vs = SpectreMsgs.ValidatorSet({
            validators: vals, hasProposer: false, proposer: vals[0], totalVotingPower: total
        });
    }

    function _buildCommitSigs(
        SpectreMsgs.ValidatorSet memory vs,
        uint16 activeCount
    )
        internal
        pure
        returns (SpectreMsgs.CommitSig[] memory sigs)
    {
        sigs = new SpectreMsgs.CommitSig[](vs.validators.length);
        for (uint256 i = 0; i < vs.validators.length; i++) {
            // The commit slot names the validator it belongs to, or the quorum check rejects the
            // proof citing it (ZK-09). Commit and pinned set are the same validators in the same
            // order here, so slot i is validator i.
            bytes20 addr = bytes20(sha256(abi.encodePacked(vs.validators[i].pubKey)));
            if (i < activeCount) {
                sigs[i] = SpectreMsgs.CommitSig({
                    flag: SpectreMsgs.CommitSigFlag.BLOCK_ID_FLAG_COMMIT, validatorAddress: addr
                });
            } else {
                sigs[i] = SpectreMsgs.CommitSig({
                    flag: SpectreMsgs.CommitSigFlag.BLOCK_ID_FLAG_ABSENT, validatorAddress: addr
                });
            }
        }
    }

    function _buildSelfConsistent(BucketConfig memory cfg)
        internal
        returns (
            SpectreMsgs.Header memory header,
            SpectreMsgs.ConsensusState memory trustedCS,
            SpectreMsgs.ValidatorSet memory vs
        )
    {
        vs = _buildValSet(cfg.valCount);
        bytes32 valSetHash = Header.hashValSet(vs);

        SpectreMsgs.BlockHeader memory bh;
        bh.chainId = CHAIN_ID;
        bh.height = NEW_HEIGHT;
        bh.time = NEW_TS_NS;
        bh.appHash = bytes32(uint256(0xCCC));
        bh.validatorsHash = valSetHash;
        bh.nextValidatorsHash = valSetHash;
        bytes32 headerHash = Header.hashHeader(bh);

        SpectreMsgs.BlockCommit memory bc = SpectreMsgs.BlockCommit({
            height: NEW_HEIGHT,
            round: 0,
            blockId: SpectreMsgs.BlockId({
                hashData: headerHash,
                partSetHeader: SpectreMsgs.PartSetHeader({ total: 1, hashData: bytes32(uint256(0x9A57)) })
            }),
            commitSigs: _buildCommitSigs(vs, cfg.activeCount)
        });

        header = SpectreMsgs.Header({
            signedHeader: SpectreMsgs.SignedHeader({ header: bh, commit: bc }),
            trustedHeight: ICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: TRUSTED_HEIGHT })
        });

        trustedCS = SpectreMsgs.ConsensusState({
            timestamp: TRUSTED_TS_NS, root: bytes32(uint256(0xAAA)), nextValidatorsHash: valSetHash
        });
    }

    function _measureRecvPacket(uint16 bucket) internal {
        // Warp to make the clock drift check happy
        vm.warp(1_700_000_020);

        wrapper = new SignatureVerifier(address(this));
        stubBucket = new AlwaysTrueVerifier();
        updateClientImpl = new UpdateClient(address(wrapper));

        wrapper.setBucket(bucket, address(stubBucket), AlwaysTrueVerifier.verifyProof.selector);

        BucketConfig memory cfg = _cfg(bucket);
        (
            SpectreMsgs.Header memory header,
            SpectreMsgs.ConsensusState memory trustedCS,
            SpectreMsgs.ValidatorSet memory vs
        ) = _buildSelfConsistent(cfg);

        SpectreMsgs.ClientState memory cs = _clientState();

        DummyMembership stubMembership = new DummyMembership();

        // Deploy the real light client with stubMembership so it succeeds
        SpectreClient realClient = new SpectreClient(
            address(updateClientImpl),
            address(stubMembership),
            STUB_MISBEHAVIOUR,
            abi.encode(cs),
            trustedCS,
            vs,
            address(0)
        );

        // Add this light client to the router
        string memory clientID = ics26Router.addClient(
            ICS02ClientMsgs.CounterpartyInfo(counterpartyId, merklePrefix), address(realClient)
        );

        // 1. Build the update message
        uint32[] memory idx = new uint32[](bucket);
        bytes32[] memory pks = new bytes32[](bucket);
        bool[] memory act = new bool[](bucket);
        uint32[] memory pinnedValidatorIndices = new uint32[](bucket);
        for (uint256 i = 0; i < bucket; i++) {
            pinnedValidatorIndices[i] = uint32(i);
            if (i < cfg.activeCount) {
                idx[i] = uint32(i);
                pks[i] = vs.validators[i].pubKey;
                act[i] = true;
            }
        }

        SpectreClientMsgs.MsgUpdateApplicationState memory m;
        m.trustedConsensusState = trustedCS;
        m.proposedHeader = header;
        m.time = NEW_TS_NS;
        m.proof.proof = [uint256(0), 0, 0, 0, 0, 0, 0, 0];
        m.proof.commitments = [uint256(0), 0];
        m.proof.commitmentPok = [uint256(0), 0];
        m.proof.bucket = bucket;
        m.proof.signerIndices = idx;
        m.proof.signerPubkeys = pks;
        m.proof.active = act;
        m.proof.pinnedValidatorIndices = pinnedValidatorIndices;
        bytes memory encodedUpdate = abi.encode(m);

        // First, update the client state to seed the consensus state height 1001
        realClient.updateApplicationState(encodedUpdate);

        // 2. Build the packet to receive
        string memory foreignDenom = "uatom";
        string memory senderStr = "cosmos1mhmwgrfrcrdex5gnr0vcqt90wknunsxej63feh";
        address receiver = makeAddr("receiver_of_foreign_denom");
        string memory receiverStr = Strings.toHexString(receiver);

        ICS20TransferMsgs.FungibleTokenPacketData memory packetData =
            _getPacketData(senderStr, receiverStr, foreignDenom);
        ICS26RouterMsgs.Payload[] memory payloads = _getPayloads(abi.encode(packetData));
        ICS26RouterMsgs.Packet memory recvPacket = ICS26RouterMsgs.Packet({
            sequence: 1,
            sourceClient: counterpartyId,
            destClient: clientID,
            timeoutTimestamp: uint64(block.timestamp + 1000),
            payloads: payloads
        });

        // Construct the verifyMembershipMsg with height 1001 matching the new height we just updated!
        bytes[] memory pathParts = new bytes[](1);
        pathParts[0] = ICS24Host.packetCommitmentPathCalldata(recvPacket.sourceClient, recvPacket.sequence);
        bytes[] memory fullPath = ICS24Host.prefixedPath(merklePrefix, pathParts[0]);

        bytes32 commitmentBz = ICS24Host.packetCommitmentBytes32(recvPacket);
        bytes memory valBytes = abi.encodePacked(commitmentBz);

        MembershipMsgs.KVPair[] memory kvPairs = new MembershipMsgs.KVPair[](1);
        kvPairs[0] = MembershipMsgs.KVPair({ path: fullPath, value: valBytes });

        MembershipMsgs.MerkleProof[] memory merkleProofs = new MembershipMsgs.MerkleProof[](1);

        LightClientMsgs.MsgVerifyMembership memory membershipMsg = LightClientMsgs.MsgVerifyMembership({
            height: ICS02ClientMsgs.Height({ revisionNumber: 0, revisionHeight: NEW_HEIGHT }),
            kvPairs: kvPairs,
            merkleProofs: merkleProofs,
            appHash: bytes32(uint256(0xCCC)), // bh.appHash from above
            trustedConsensusState: SpectreMsgs.ConsensusState({
                timestamp: NEW_TS_NS, root: bytes32(uint256(0xCCC)), nextValidatorsHash: Header.hashValSet(vs)
            }),
            membershipType: MembershipMsgs.MembershipType.Membership,
            path: new bytes[](0),
            value: bytes("")
        });

        ICS26RouterMsgs.MsgRecvPacket memory msgRecvPacket =
            ICS26RouterMsgs.MsgRecvPacket({ packet: recvPacket, membershipMsg: abi.encode(membershipMsg) });

        // 3. Measure recvPacket gas
        uint256 g0 = gasleft();
        ics26Router.recvPacket(msgRecvPacket);
        uint256 gasUsed = g0 - gasleft();

        console.log("=========================================");
        console.log("Bucket Size:", bucket);
        console.log("recvPacket gas used:", gasUsed);
        console.log("=========================================");
    }

    function test_gas_recvPacket_n004() public {
        _measureRecvPacket(4);
    }

    function test_gas_recvPacket_n008() public {
        _measureRecvPacket(8);
    }

    function test_gas_recvPacket_n016() public {
        _measureRecvPacket(16);
    }

    function test_gas_recvPacket_n032() public {
        _measureRecvPacket(32);
    }

    function test_gas_recvPacket_n064() public {
        _measureRecvPacket(64);
    }

    function test_gas_recvPacket_n128() public {
        _measureRecvPacket(128);
    }
}
