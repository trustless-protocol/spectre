// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Membership } from "contracts/light-clients/spectre/modules/Membership.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";

/// @dev Exposes the internal checkExistenceProof so the depth bound (issue #110)
/// can be unit-tested directly.
contract MembershipDepthHarness is Membership {
    function exposedCheckExistenceProof(
        MembershipMsgs.ExistenceProof memory proof,
        MembershipMsgs.ProofSpec memory spec
    )
        external
        view
    {
        checkExistenceProof(proof, spec);
    }

    /// @dev calculateExistenceRoot is the function that actually loops over the
    /// path, and the non-membership path reaches it (via calculateNonExistenceRoot)
    /// BEFORE checkExistenceProof. The cap must live here.
    function exposedCalculateExistenceRoot(MembershipMsgs.ExistenceProof memory proof) external view returns (bytes32) {
        return calculateExistenceRoot(proof);
    }
}

contract MembershipDepthTest is Test {
    MembershipDepthHarness internal m;

    uint256 internal constant MAX_DEPTH = 128; // mirrors Membership.MAX_PROOF_DEPTH

    function setUp() public {
        m = new MembershipDepthHarness();
    }

    // Build an existence proof whose leaf and inner ops match the Tendermint spec
    // (so validation reaches the depth gate) with a path of `pathLen` inner ops.
    function _proofWithPath(
        uint256 pathLen,
        MembershipMsgs.ProofSpec memory spec
    )
        internal
        pure
        returns (MembershipMsgs.ExistenceProof memory)
    {
        MembershipMsgs.InnerOp[] memory path = new MembershipMsgs.InnerOp[](pathLen);
        for (uint256 i = 0; i < pathLen; i++) {
            path[i] = MembershipMsgs.InnerOp({ hashOp: MembershipMsgs.HashOp.SHA256, prefix: hex"01", suffix: "" });
        }
        return MembershipMsgs.ExistenceProof({
            key: bytes("key"),
            value: bytes("value"),
            leaf: MembershipMsgs.LeafOp({
                hashOp: spec.leafOp.hashOp,
                prehashKey: spec.leafOp.prehashKey,
                prehashValue: spec.leafOp.prehashValue,
                prefix: spec.leafOp.prefix
            }),
            path: path
        });
    }

    function test_checkExistenceProof_revertsWhenPathExceedsMaxDepth() public {
        MembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        MembershipMsgs.ExistenceProof memory proof = _proofWithPath(MAX_DEPTH + 1, spec);

        vm.expectRevert(abi.encodeWithSelector(Membership.ProofPathTooLong.selector, MAX_DEPTH + 1, MAX_DEPTH));
        m.exposedCheckExistenceProof(proof, spec);
    }

    function test_checkExistenceProof_allowsPathAtMaxDepth() public view {
        MembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        MembershipMsgs.ExistenceProof memory proof = _proofWithPath(MAX_DEPTH, spec);

        // Reaching here without ProofPathTooLong means the boundary (== MAX_PROOF_DEPTH)
        // is accepted.
        m.exposedCheckExistenceProof(proof, spec);
    }

    function test_checkExistenceProof_allowsTypicalShallowPath() public view {
        MembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        MembershipMsgs.ExistenceProof memory proof = _proofWithPath(8, spec);
        m.exposedCheckExistenceProof(proof, spec);
    }

    // Direct proof that the cap lives in calculateExistenceRoot (the loop source),
    // not only in checkExistenceProof. Before the fix this would loop over the path
    // and return a root; now it reverts before the loop.
    function test_calculateExistenceRoot_revertsWhenPathExceedsMaxDepth() public {
        MembershipMsgs.InnerOp[] memory longPath = new MembershipMsgs.InnerOp[](MAX_DEPTH + 1);
        for (uint256 i = 0; i <= MAX_DEPTH; i++) {
            longPath[i] = MembershipMsgs.InnerOp({ hashOp: MembershipMsgs.HashOp.SHA256, prefix: hex"01", suffix: "" });
        }
        MembershipMsgs.ExistenceProof memory proof = MembershipMsgs.ExistenceProof({
            key: bytes("k"),
            value: bytes("v"),
            leaf: MembershipMsgs.LeafOp({
                hashOp: MembershipMsgs.HashOp.SHA256,
                prehashKey: MembershipMsgs.HashOp.NO_HASH,
                prehashValue: MembershipMsgs.HashOp.SHA256,
                prefix: hex"00"
            }),
            path: longPath
        });
        vm.expectRevert(abi.encodeWithSelector(Membership.ProofPathTooLong.selector, MAX_DEPTH + 1, MAX_DEPTH));
        m.exposedCalculateExistenceRoot(proof);
    }

    // The non-membership path computes calculateNonExistenceRoot (-> calculateExistenceRoot,
    // which loops over the path) BEFORE checkExistenceProof runs. This guards against an
    // over-long nonExistenceProof.left/right path bypassing the cap via the public entrypoint.
    function test_membership_nonExistence_revertsWhenLeftPathExceedsMaxDepth() public {
        MembershipMsgs.InnerOp[] memory longPath = new MembershipMsgs.InnerOp[](MAX_DEPTH + 1);
        for (uint256 i = 0; i <= MAX_DEPTH; i++) {
            longPath[i] = MembershipMsgs.InnerOp({ hashOp: MembershipMsgs.HashOp.SHA256, prefix: hex"01", suffix: "" });
        }

        MembershipMsgs.ExistenceProof memory left = MembershipMsgs.ExistenceProof({
            key: bytes("k"),
            value: bytes("v"),
            leaf: MembershipMsgs.LeafOp({
                hashOp: MembershipMsgs.HashOp.SHA256,
                prehashKey: MembershipMsgs.HashOp.NO_HASH,
                prehashValue: MembershipMsgs.HashOp.SHA256,
                prefix: hex"00"
            }),
            path: longPath
        });

        MembershipMsgs.ExistenceProof memory emptyProof; // unused right neighbour
        MembershipMsgs.NonExistenceProof memory ne = MembershipMsgs.NonExistenceProof({
            key: bytes("absent"), hasLeft: true, left: left, hasRight: false, right: emptyProof
        });

        // proofs.length must equal proofSpecs.length (2) to reach calculateNonExistenceRoot.
        MembershipMsgs.CommitmentProof[] memory proofs = new MembershipMsgs.CommitmentProof[](2);
        proofs[0].proofType = MembershipMsgs.ProofType.NON_EXIST;
        proofs[0].nonExistenceProof = ne;
        proofs[1].proofType = MembershipMsgs.ProofType.EXIST; // unused; reverts before this is read

        MembershipMsgs.MerkleProof[] memory mps = new MembershipMsgs.MerkleProof[](1);
        mps[0] = MembershipMsgs.MerkleProof({ proofs: proofs });

        bytes[] memory path = new bytes[](2);
        path[0] = bytes("ibc");
        path[1] = bytes("receipts/key");
        MembershipMsgs.KVPair[] memory kvs = new MembershipMsgs.KVPair[](1);
        kvs[0] = MembershipMsgs.KVPair({ path: path, value: bytes("") }); // empty value -> non-membership

        vm.expectRevert(abi.encodeWithSelector(Membership.ProofPathTooLong.selector, MAX_DEPTH + 1, MAX_DEPTH));
        m.verifyMembership(bytes32(uint256(0xABCD)), kvs, mps);
    }
}
