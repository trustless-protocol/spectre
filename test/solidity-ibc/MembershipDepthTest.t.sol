// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Membership } from "../../contracts/programs/Membership.sol";
import { IMembershipMsgs } from "../../contracts/light-clients/msgs/IMembershipMsgs.sol";

/// @dev Exposes the internal checkExistenceProof so the depth bound (issue #110)
/// can be unit-tested directly.
contract MembershipDepthHarness is Membership {
    function exposedCheckExistenceProof(
        IMembershipMsgs.ExistenceProof memory proof,
        IMembershipMsgs.ProofSpec memory spec
    )
        external
        view
    {
        checkExistenceProof(proof, spec);
    }
}

contract MembershipDepthTest is Test {
    MembershipDepthHarness internal m;

    uint256 internal constant MAX_DEPTH = 128; // mirrors Membership.MAX_PROOF_DEPTH

    function setUp() public {
        m = new MembershipDepthHarness();
    }

    // Build an existence proof whose leaf matches the Tendermint spec (so the leaf
    // checks pass and we reach the depth gate) with a path of `pathLen` inner ops.
    // The Tendermint spec skips the IAVL inner-op loop, so the path contents are
    // irrelevant past the depth gate.
    function _proofWithPath(
        uint256 pathLen,
        IMembershipMsgs.ProofSpec memory spec
    )
        internal
        pure
        returns (IMembershipMsgs.ExistenceProof memory)
    {
        IMembershipMsgs.InnerOp[] memory path = new IMembershipMsgs.InnerOp[](pathLen);
        for (uint256 i = 0; i < pathLen; i++) {
            path[i] = IMembershipMsgs.InnerOp({ hashOp: IMembershipMsgs.HashOp.SHA256, prefix: hex"01", suffix: "" });
        }
        return IMembershipMsgs.ExistenceProof({
            key: bytes("key"),
            value: bytes("value"),
            leaf: IMembershipMsgs.LeafOp({
                hashOp: spec.leafOp.hashOp,
                prehashKey: spec.leafOp.prehashKey,
                prehashValue: spec.leafOp.prehashValue,
                prefix: spec.leafOp.prefix
            }),
            path: path
        });
    }

    function test_checkExistenceProof_revertsWhenPathExceedsMaxDepth() public {
        IMembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        IMembershipMsgs.ExistenceProof memory proof = _proofWithPath(MAX_DEPTH + 1, spec);

        vm.expectRevert(abi.encodeWithSelector(Membership.ProofPathTooLong.selector, MAX_DEPTH + 1, MAX_DEPTH));
        m.exposedCheckExistenceProof(proof, spec);
    }

    function test_checkExistenceProof_allowsPathAtMaxDepth() public view {
        IMembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        IMembershipMsgs.ExistenceProof memory proof = _proofWithPath(MAX_DEPTH, spec);

        // Reaching here without ProofPathTooLong means the boundary (== MAX_PROOF_DEPTH)
        // is accepted; the Tendermint spec skips the IAVL inner-op loop afterwards.
        m.exposedCheckExistenceProof(proof, spec);
    }

    function test_checkExistenceProof_allowsTypicalShallowPath() public view {
        IMembershipMsgs.ProofSpec memory spec = m.tendermintSpec();
        IMembershipMsgs.ExistenceProof memory proof = _proofWithPath(8, spec);
        m.exposedCheckExistenceProof(proof, spec);
    }
}
