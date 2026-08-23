// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Membership } from "contracts/light-clients/spectre/modules/Membership.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";

/// @dev Exposes proof validation/root helpers for focused inner-spec regression tests.
contract MembershipInnerSpecHarness is Membership {
    function exposedCheckExistenceProof(
        MembershipMsgs.ExistenceProof memory proof,
        MembershipMsgs.ProofSpec memory spec
    )
        external
        view
    {
        checkExistenceProof(proof, spec);
    }

    function exposedCalculateExistenceRoot(MembershipMsgs.ExistenceProof memory proof) external view returns (bytes32) {
        return calculateExistenceRoot(proof);
    }
}

contract MembershipInnerSpecTest is Test {
    MembershipInnerSpecHarness internal m;

    function setUp() public {
        m = new MembershipInnerSpecHarness();
    }

    function _singleInnerPath(
        MembershipMsgs.HashOp hashOp,
        bytes memory prefix,
        bytes memory suffix
    )
        internal
        pure
        returns (MembershipMsgs.InnerOp[] memory)
    {
        MembershipMsgs.InnerOp[] memory path = new MembershipMsgs.InnerOp[](1);
        path[0] = MembershipMsgs.InnerOp({ hashOp: hashOp, prefix: prefix, suffix: suffix });
        return path;
    }

    function _iavlProof(
        bytes memory key,
        bytes memory value
    )
        internal
        pure
        returns (MembershipMsgs.ExistenceProof memory)
    {
        return MembershipMsgs.ExistenceProof({
            key: key,
            value: value,
            leaf: MembershipMsgs.LeafOp({
                hashOp: MembershipMsgs.HashOp.SHA256,
                prehashKey: MembershipMsgs.HashOp.NO_HASH,
                prehashValue: MembershipMsgs.HashOp.SHA256,
                prefix: hex"000000"
            }),
            path: new MembershipMsgs.InnerOp[](0)
        });
    }

    function _tendermintProof(
        bytes memory key,
        bytes memory value,
        MembershipMsgs.InnerOp[] memory path
    )
        internal
        pure
        returns (MembershipMsgs.ExistenceProof memory)
    {
        return MembershipMsgs.ExistenceProof({
            key: key,
            value: value,
            leaf: MembershipMsgs.LeafOp({
                hashOp: MembershipMsgs.HashOp.SHA256,
                prehashKey: MembershipMsgs.HashOp.NO_HASH,
                prehashValue: MembershipMsgs.HashOp.SHA256,
                prefix: hex"00"
            }),
            path: path
        });
    }

    function _kvPairs(bytes memory ibcKey, bytes memory value) internal pure returns (MembershipMsgs.KVPair[] memory) {
        bytes[] memory path = new bytes[](2);
        path[0] = bytes("ibc");
        path[1] = ibcKey;

        MembershipMsgs.KVPair[] memory kvs = new MembershipMsgs.KVPair[](1);
        kvs[0] = MembershipMsgs.KVPair({ path: path, value: value });
        return kvs;
    }

    function _merkleProofs(
        MembershipMsgs.ExistenceProof memory iavlProof,
        MembershipMsgs.ExistenceProof memory tendermintProof
    )
        internal
        pure
        returns (MembershipMsgs.MerkleProof[] memory)
    {
        MembershipMsgs.CommitmentProof[] memory proofs = new MembershipMsgs.CommitmentProof[](2);
        proofs[0].proofType = MembershipMsgs.ProofType.EXIST;
        proofs[0].existenceProof = iavlProof;
        proofs[1].proofType = MembershipMsgs.ProofType.EXIST;
        proofs[1].existenceProof = tendermintProof;

        MembershipMsgs.MerkleProof[] memory merkleProofs = new MembershipMsgs.MerkleProof[](1);
        merkleProofs[0] = MembershipMsgs.MerkleProof({ proofs: proofs });
        return merkleProofs;
    }

    function test_checkExistenceProof_rejectsTendermintNoHashInnerOp() public {
        bytes32 appHash = bytes32(uint256(0xABCD));
        MembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        MembershipMsgs.ExistenceProof memory proof = _tendermintProof(
            bytes("ibc"),
            abi.encodePacked(bytes32(uint256(0xCAFE))),
            _singleInnerPath(MembershipMsgs.HashOp.NO_HASH, abi.encodePacked(appHash), "")
        );

        vm.expectRevert(Membership.UnexpectedInnerHashOp.selector);
        m.exposedCheckExistenceProof(proof, spec);
    }

    function test_membership_rejectsForgedTendermintNoHashRoot() public {
        bytes memory ibcKey = bytes("commitments/forged");
        bytes memory commitmentValue = abi.encodePacked(bytes32(uint256(0xCAFE)));
        bytes32 appHash = bytes32(uint256(0xABCD));

        MembershipMsgs.ExistenceProof memory iavlProof = _iavlProof(ibcKey, commitmentValue);
        bytes32 fakeIbcRoot = m.exposedCalculateExistenceRoot(iavlProof);
        MembershipMsgs.ExistenceProof memory tendermintProof = _tendermintProof(
            bytes("ibc"),
            abi.encodePacked(fakeIbcRoot),
            _singleInnerPath(MembershipMsgs.HashOp.NO_HASH, abi.encodePacked(appHash), "")
        );

        vm.expectRevert(Membership.UnexpectedInnerHashOp.selector);
        m.verifyMembership(appHash, _kvPairs(ibcKey, commitmentValue), _merkleProofs(iavlProof, tendermintProof));
    }

    function test_membership_acceptsValidTendermintSha256InnerOp() public {
        bytes memory ibcKey = bytes("commitments/real");
        bytes memory commitmentValue = abi.encodePacked(bytes32(uint256(0xBEEF)));

        MembershipMsgs.ExistenceProof memory iavlProof = _iavlProof(ibcKey, commitmentValue);
        bytes32 fakeIbcRoot = m.exposedCalculateExistenceRoot(iavlProof);
        MembershipMsgs.ExistenceProof memory tendermintProof = _tendermintProof(
            bytes("ibc"),
            abi.encodePacked(fakeIbcRoot),
            _singleInnerPath(MembershipMsgs.HashOp.SHA256, hex"01", abi.encodePacked(bytes32(uint256(0x1234))))
        );
        bytes32 appHash = m.exposedCalculateExistenceRoot(tendermintProof);

        m.verifyMembership(appHash, _kvPairs(ibcKey, commitmentValue), _merkleProofs(iavlProof, tendermintProof));
    }
}
