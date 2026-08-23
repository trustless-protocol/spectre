// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { Membership } from "contracts/light-clients/spectre/modules/Membership.sol";
import { MembershipMsgs } from "contracts/light-clients/spectre/messages/MembershipMsgs.sol";

/// @dev Documents the intentional 32-byte commitment-value invariant (issue #111).
/// IBC v2 only ever proves 32-byte hash commitments on this path; a non-32-byte
/// value must be rejected up front rather than silently mis-verified. If the
/// invariant is ever relaxed (to support raw variable-length values), these tests
/// should be updated together with a careful leaf-level generalization.
contract MembershipValueLengthTest is Test {
    Membership internal m;

    bytes32 internal constant APP_HASH = bytes32(uint256(0xABCD));

    function setUp() public {
        m = new Membership();
    }

    function _kvWithValue(bytes memory value) internal pure returns (MembershipMsgs.KVPair[] memory) {
        bytes[] memory path = new bytes[](2);
        path[0] = bytes("ibc");
        path[1] = bytes("commitments/key");
        MembershipMsgs.KVPair[] memory kvs = new MembershipMsgs.KVPair[](1);
        kvs[0] = MembershipMsgs.KVPair({ path: path, value: value });
        return kvs;
    }

    // A single dummy existence proof — enough to pass the proofs.length > 0 check;
    // verification reverts at the value-length gate before the proof is used.
    function _oneProof() internal pure returns (MembershipMsgs.MerkleProof[] memory) {
        MembershipMsgs.CommitmentProof[] memory proofs = new MembershipMsgs.CommitmentProof[](1);
        proofs[0].proofType = MembershipMsgs.ProofType.EXIST;
        MembershipMsgs.MerkleProof[] memory mps = new MembershipMsgs.MerkleProof[](1);
        mps[0] = MembershipMsgs.MerkleProof({ proofs: proofs });
        return mps;
    }

    function test_membership_rejectsShortValue() public {
        MembershipMsgs.KVPair[] memory kvs = _kvWithValue(new bytes(16)); // 16 bytes
        vm.expectRevert(Membership.InvalidValueLength.selector);
        m.verifyMembership(APP_HASH, kvs, _oneProof());
    }

    function test_membership_rejectsLongValue() public {
        MembershipMsgs.KVPair[] memory kvs = _kvWithValue(new bytes(64)); // 64 bytes
        vm.expectRevert(Membership.InvalidValueLength.selector);
        m.verifyMembership(APP_HASH, kvs, _oneProof());
    }

    function test_membership_thirtyTwoByteValuePassesLengthGate() public {
        // A 32-byte value must NOT revert with InvalidValueLength — it proceeds to
        // real proof verification (which fails on the dummy proof with a different
        // error). We only assert the length gate is not the rejection reason.
        MembershipMsgs.KVPair[] memory kvs = _kvWithValue(new bytes(32));
        try m.verifyMembership(APP_HASH, kvs, _oneProof()) {
        // proof would not actually verify with dummy data
        }
        catch (bytes memory reason) {
            require(bytes4(reason) != Membership.InvalidValueLength.selector, "32-byte value must pass the length gate");
        }
    }
}
