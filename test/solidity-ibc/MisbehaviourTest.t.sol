// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Misbehaviour } from "../../contracts/programs/Misbehaviour.sol";
import { Header } from "../../contracts/utils/Header.sol";
import { IMisbehaviourMsgs } from "../../contracts/light-clients/msgs/IMisbehaviourMsgs.sol";
import { IICS07TendermintMsgs } from "../../contracts/light-clients/msgs/IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../contracts/msgs/IICS02ClientMsgs.sol";

contract MisbehaviourHarness is Misbehaviour {
    function validateBasicHarness(IMisbehaviourMsgs.Misbehaviour memory misbehaviour_) external pure {
        validateBasic(misbehaviour_);
    }
}

contract MisbehaviourTest is Test {
    MisbehaviourHarness internal harness;

    function setUp() public {
        harness = new MisbehaviourHarness();
    }

    function test_validateBasic_detects_conflicting_block_ids_even_if_app_hash_matches() public view {
        IMisbehaviourMsgs.Misbehaviour memory misbehaviour_ = _buildMisbehaviour({
            appHash1: bytes32(uint256(0xAA)),
            appHash2: bytes32(uint256(0xAA)),
            blockIdHash1: bytes32(uint256(0x1111)),
            blockIdHash2: bytes32(uint256(0x2222))
        });

        harness.validateBasicHarness(misbehaviour_);
    }

    function test_validateBasic_reverts_when_block_id_hash_matches() public {
        IMisbehaviourMsgs.Misbehaviour memory misbehaviour_ = _buildMisbehaviour({
            appHash1: bytes32(uint256(0xAA)),
            appHash2: bytes32(uint256(0xBB)),
            blockIdHash1: bytes32(uint256(0x1111)),
            blockIdHash2: bytes32(uint256(0x1111))
        });

        vm.expectRevert(Misbehaviour.MisbehaviourNotDetected.selector);
        harness.validateBasicHarness(misbehaviour_);
    }

    function _buildMisbehaviour(
        bytes32 appHash1,
        bytes32 appHash2,
        bytes32 blockIdHash1,
        bytes32 blockIdHash2
    ) internal pure returns (IMisbehaviourMsgs.Misbehaviour memory misbehaviour_) {
        misbehaviour_.client_id = IICS07TendermintMsgs.ChainId({ id: "cosmoshub-1", revisionNumber: 1 });
        misbehaviour_.header1 = _buildHeader(appHash1, blockIdHash1);
        misbehaviour_.header2 = _buildHeader(appHash2, blockIdHash2);
    }

    function _buildHeader(
        bytes32 appHash,
        bytes32 blockIdHash
    ) internal pure returns (IICS07TendermintMsgs.Header memory header) {
        IICS07TendermintMsgs.ValidatorSet memory validatorSet = _validatorSet();
        bytes32 validatorsHash = Header.hashValSet(validatorSet);

        header.validatorSet = validatorSet;
        header.trustedNextValidatorSet = validatorSet;
        header.trustedHeight = IICS02ClientMsgs.Height({ revisionNumber: 1, revisionHeight: 9 });
        header.signedHeader.header = IICS07TendermintMsgs.BlockHeader({
            version: IICS07TendermintMsgs.Version({ blockVersion: 1, appVersion: 1 }),
            chainId: "cosmoshub-1",
            height: 10,
            time: 1_700_000_000_000_000_000,
            hasLastBlockId: false,
            lastBlockId: IICS07TendermintMsgs.BlockId({
                hashData: bytes32(0),
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 0, hashData: bytes32(0) })
            }),
            hasLastCommitHash: false,
            lastCommitHash: bytes32(0),
            hasDataHash: false,
            dataHash: bytes32(0),
            validatorsHash: validatorsHash,
            nextValidatorsHash: validatorsHash,
            consensusHash: bytes32(uint256(0xC0)),
            appHash: appHash,
            hasLastResultsHash: false,
            lastResultsHash: bytes32(0),
            hasEvidenceHash: false,
            evidenceHash: bytes32(0),
            proposerAddress: hex"01"
        });
        header.signedHeader.commit = IICS07TendermintMsgs.BlockCommit({
            height: 10,
            round: 0,
            blockId: IICS07TendermintMsgs.BlockId({
                hashData: blockIdHash,
                partSetHeader: IICS07TendermintMsgs.PartSetHeader({ total: 1, hashData: bytes32(uint256(0x1234)) })
            }),
            commitSigs: new IICS07TendermintMsgs.CommitSig[](0)
        });
    }

    function _validatorSet() internal pure returns (IICS07TendermintMsgs.ValidatorSet memory validatorSet) {
        IICS07TendermintMsgs.ValidatorInfo[] memory validators = new IICS07TendermintMsgs.ValidatorInfo[](1);
        validators[0] = IICS07TendermintMsgs.ValidatorInfo({
            valAddress: hex"01",
            pubKey: bytes32(uint256(0xABC)),
            votingPower: 10,
            proposerPriority: 0
        });

        validatorSet = IICS07TendermintMsgs.ValidatorSet({
            validators: validators,
            hasProposer: false,
            proposer: IICS07TendermintMsgs.ValidatorInfo({
                valAddress: new bytes(0),
                pubKey: bytes32(0),
                votingPower: 0,
                proposerPriority: 0
            }),
            totalVotingPower: 10
        });
    }
}
