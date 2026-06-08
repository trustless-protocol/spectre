// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import { Test } from "forge-std/Test.sol";

import { WrapperVerifier } from "../../contracts/utils/WrapperVerifier.sol";
import { IVerifier } from "../../contracts/interfaces/IVerifier.sol";

/// @dev Exposes the PRODUCTION canonical-vote builder used inside
/// `WrapperVerifier._hashWitness` for active slots, so it can be golden-vectored
/// directly. It reuses the exact same internal routines the on-chain proof path
/// uses (`_sharedBlockIdBytes` / `_buildVotePrefix` / `_buildChainSuffix` /
/// `_writeVoteSignBytes`) — no re-implementation.
contract WrapperVerifierVoteHarness is WrapperVerifier {
    constructor() WrapperVerifier(address(0xdead)) { }

    /// @notice Build the length-prefixed CanonicalVote bytes exactly as an active
    /// witness slot does, and return them for byte-for-byte comparison.
    function buildCanonicalVote(
        IVerifier.SharedBlock calldata shared,
        uint64 tsSec,
        uint32 tsNanos
    )
        external
        pure
        returns (bytes memory)
    {
        bytes memory encodedBlockId = _sharedBlockIdBytes(shared);
        bytes memory commonVotePrefix = _buildVotePrefix(shared, encodedBlockId);
        bytes memory chainSuffix = _buildChainSuffix(shared.chainID);

        bytes memory buf = new bytes(MAX_MSG_LEN);
        uint256 msgLen = _writeVoteSignBytes(buf, 0, commonVotePrefix, chainSuffix, tsSec, tsNanos);

        bytes memory out = new bytes(msgLen);
        assembly {
            mcopy(add(out, 0x20), add(buf, 0x20), msgLen)
        }
        return out;
    }
}

/// @notice Golden-vector lock for the **production** canonical-vote encoder
/// (`WrapperVerifier`'s inline `_writeVoteSignBytes` path that the on-chain proof
/// verification actually uses), against Go/CometBFT `Commit.VoteSignBytes()`.
///
/// `EncodeTest` golden-vectors `Encode.voteSignBytes`, but that library function
/// is not on the proof-verification path — `WrapperVerifier._hashWitness` rebuilds
/// the canonical vote with its own (gas-optimized, direct-write) routine. These
/// vectors are the SAME CometBFT bytes asserted in EncodeTest, applied to the
/// production encoder so a divergence between it and CometBFT is caught directly
/// (not only via a Solidity-vs-Solidity equivalence reference). Companion to #109.
contract WrapperVerifierVoteGoldenTest is Test {
    WrapperVerifierVoteHarness internal h;

    // Same fixtures as EncodeTest: block hash, part-set hash, chain id, timestamp.
    bytes32 internal constant BLOCK_HASH = 0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f;
    bytes32 internal constant PARTSET_HASH = 0x101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f;
    uint64 internal constant TS_SECONDS = 1_704_067_200;

    function setUp() public {
        h = new WrapperVerifierVoteHarness();
    }

    function _shared(uint64 height, uint64 round, bytes memory chainId)
        internal
        pure
        returns (IVerifier.SharedBlock memory)
    {
        return IVerifier.SharedBlock({
            height: height,
            round: round,
            blockIDHash: BLOCK_HASH,
            partSetTotal: 1,
            partSetHash: PARTSET_HASH,
            chainID: chainId
        });
    }

    /// height 100, round 0, ts 1704067200s, chain "test-chain".
    function test_prod_voteSignBytes_full() public view {
        assertEq(
            h.buildCanonicalVote(_shared(100, 0, bytes("test-chain")), TS_SECONDS, 0),
            hex"69080211640000000000000022480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a06088081c8ac06320a746573742d636861696e"
        );
    }

    /// height 0 + round 0 are omitted (proto3 scalar zero); timestamp still present.
    function test_prod_voteSignBytes_zeroHeightZeroRound() public view {
        assertEq(
            h.buildCanonicalVote(_shared(0, 0, bytes("test-chain")), TS_SECONDS, 0),
            hex"60080222480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a06088081c8ac06320a746573742d636861696e"
        );
    }

    /// height 128, round 127 (multi-byte sfixed64 values).
    function test_prod_voteSignBytes_nonZeroHeightAndRound() public view {
        assertEq(
            h.buildCanonicalVote(_shared(128, 127, bytes("test-chain")), TS_SECONDS, 0),
            hex"720802118000000000000000197f0000000000000022480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a06088081c8ac06320a746573742d636861696e"
        );
    }

    /// THE FIX: zero timestamp must still emit the field as `2a00` (CometBFT
    /// non-nullable stdtime), not omit it.
    function test_prod_voteSignBytes_zeroTimestamp() public view {
        assertEq(
            h.buildCanonicalVote(_shared(100, 0, bytes("test-chain")), 0, 0),
            hex"63080211640000000000000022480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a00320a746573742d636861696e"
        );
    }

    /// empty chain id => chain_id field omitted.
    function test_prod_voteSignBytes_noChainId() public view {
        assertEq(
            h.buildCanonicalVote(_shared(100, 0, bytes("")), TS_SECONDS, 0),
            hex"5d080211640000000000000022480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a06088081c8ac06"
        );
    }

    /// height 0, round 5: height field omitted (proto3 zero), round field present.
    /// This is the canonical CometBFT zero-omission cross-check: a non-zero round with
    /// zero height must encode the round but not the height.
    function test_prod_voteSignBytes_zeroHeightNonZeroRound() public view {
        assertEq(
            h.buildCanonicalVote(_shared(0, 5, bytes("test-chain")), TS_SECONDS, 0),
            hex"69080219050000000000000022480a20202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f122408011220101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f2a06088081c8ac06320a746573742d636861696e"
        );
    }
}
