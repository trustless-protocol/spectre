// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IVerifier} from "../interfaces/IVerifier.sol";
import {LibBytes} from "./LibBytes.sol";

/// @title WrapperVerifier
/// @notice Decompresses Ed25519 points, packs per-slot public inputs, and
///         dispatches to the per-bucket gnark Groth16 verifier selected by
///         `bucket`. One `BucketVerifier` must be registered per supported
///         bucket before `verifyBatchProof` can succeed for that size.
contract WrapperVerifier is IVerifier {
    uint256 constant mask =
        0x8000000000000000000000000000000000000000000000000000000000000000;
    uint256 constant a =
        0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffec;
    uint256 constant p =
        0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffed;
    uint256 constant d =
        0x52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a3;
    uint256 constant I =
        0x2b8324804fc1df0b2b4d00993dfbd7a72f431806ad2fe478c4ee1b274a0ea0b0;
    uint256 constant pp3div8 =
        0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe;
    uint256 constant N =
        0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed;
    uint256 constant sqrtm1 =
        0x2b8324804fc1df0b2b4d00993dfbd7a72f431806ad2fe478c4ee1b274a0ea0b0;

    uint256 constant MINUS_2 =
        0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeb;
    uint256 constant MINUS_2MODN =
        0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb;
    address constant MODEXP_PRECOMPILE =
        0x0000000000000000000000000000000000000005;

    // MAX_CHAIN_ID_LEN mirrors canonvote.MaxChainIDLen — the padded width
    // every chainID buffer is extended to before entering the circuit.
    uint16 constant MAX_CHAIN_ID_LEN = 48;

    // Per-slot public-input word count. Approach B: BatchCircuit exposes
    // (R.X, R.Y, S, A.X, A.Y) as 5 field elements × 4 limbs = 20 words, plus
    // (timestampSeconds, timestampNanos) as two single-word Variables = 22.
    // Message bytes are NOT public input — the circuit reconstructs the
    // canonical vote from shared block data + timestamps.
    uint16 constant SLOT_WORDS = 22;

    // Shared-block suffix word count: each byte in the 32-byte hashes and
    // the 48-byte padded chainID becomes its own Fr public input (each U8 in
    // the circuit is a single Variable), matching BatchCircuit's field order.
    //   height (1) + round (1) + blockIDHash (32) + partSetTotal (1)
    //   + partSetHash (32) + chainID (48) + chainIDLen (1) = 116 words.
    uint16 constant SHARED_WORDS = 1 + 1 + 32 + 1 + 32 + MAX_CHAIN_ID_LEN + 1;

    /// @param verifier Address of the per-bucket gnark-generated Groth16 verifier.
    /// @param selector 4-byte function selector of that verifier's `verifyProof`.
    ///                 Differs per bucket because the `uint256[L]` input length
    ///                 is part of the type signature.
    struct BucketVerifier {
        address verifier;
        bytes4 selector;
    }

    address public immutable OWNER;
    mapping(uint16 => BucketVerifier) public buckets;

    error UnknownBucket(uint16 bucket);
    error LengthMismatch();
    error ChainIDTooLong(uint256 length);
    error NotOwner();

    constructor(address owner) {
        OWNER = owner;
    }

    function setBucket(uint16 bucket, address verifier, bytes4 selector) external {
        if (msg.sender != OWNER) revert NotOwner();
        buckets[bucket] = BucketVerifier({verifier: verifier, selector: selector});
    }

    /// @inheritdoc IVerifier
    function verifyBatchProof(
        uint16 bucket,
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        bytes32[2][] calldata signatures,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        IVerifier.SharedBlock calldata shared
    ) external override returns (bool) {
        if (
            signatures.length != bucket
                || pubkeys.length != bucket
                || timestampSeconds.length != bucket
                || timestampNanos.length != bucket
        ) revert LengthMismatch();
        if (shared.chainID.length > MAX_CHAIN_ID_LEN) revert ChainIDTooLong(shared.chainID.length);
        BucketVerifier memory bv = buckets[bucket];
        if (bv.verifier == address(0)) revert UnknownBucket(bucket);

        uint256 inputLen = uint256(bucket) * SLOT_WORDS + SHARED_WORDS;
        uint256[] memory publicInputs = new uint256[](inputLen);
        uint256 off = _packPerSlot(publicInputs, signatures, pubkeys, timestampSeconds, timestampNanos);
        _packShared(publicInputs, off, shared);

        return _dispatch(bv, proof, commitments, commitmentPok, publicInputs);
    }

    /// @dev Pack every slot's (R.X, R.Y, S, A.X, A.Y, tsSec, tsNanos) —
    ///      matching BatchCircuit's per-slot field ordering.
    function _packPerSlot(
        uint256[] memory publicInputs,
        bytes32[2][] calldata signatures,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos
    ) internal returns (uint256 off) {
        for (uint256 i = 0; i < signatures.length; i++) {
            bytes32 R = signatures[i][0];
            uint256 S = reverseBytes(uint256(signatures[i][1]));
            (uint256 rX, uint256 rY) = decodePoint(reverseBytes(uint256(R)));
            (uint256 aX, uint256 aY) = decodePoint(reverseBytes(uint256(pubkeys[i])));

            off = _writeLimbs(publicInputs, off, rX);
            off = _writeLimbs(publicInputs, off, rY);
            off = _writeLimbs(publicInputs, off, S);
            off = _writeLimbs(publicInputs, off, aX);
            off = _writeLimbs(publicInputs, off, aY);

            publicInputs[off++] = uint256(timestampSeconds[i]);
            publicInputs[off++] = uint256(timestampNanos[i]);
        }
    }

    /// @dev Pack the SharedBlockData suffix: scalar Variables + per-byte U8s
    ///      for the 32-byte hashes and the MAX_CHAIN_ID_LEN-padded chainID.
    function _packShared(
        uint256[] memory publicInputs,
        uint256 off,
        IVerifier.SharedBlock calldata shared
    ) internal pure {
        publicInputs[off++] = uint256(shared.height);
        publicInputs[off++] = uint256(shared.round);
        for (uint256 i = 0; i < 32; i++) {
            publicInputs[off + i] = uint256(uint8(shared.blockIDHash[i]));
        }
        off += 32;
        publicInputs[off++] = uint256(shared.partSetTotal);
        for (uint256 i = 0; i < 32; i++) {
            publicInputs[off + i] = uint256(uint8(shared.partSetHash[i]));
        }
        off += 32;
        for (uint256 i = 0; i < MAX_CHAIN_ID_LEN; i++) {
            publicInputs[off + i] = i < shared.chainID.length
                ? uint256(uint8(shared.chainID[i]))
                : 0;
        }
        off += MAX_CHAIN_ID_LEN;
        publicInputs[off] = shared.chainID.length;
    }

    /// @dev Encode calldata matching gnark's fixed-size input ABI
    ///      (`verifyProof(uint256[8], uint256[2], uint256[2], uint256[L])`)
    ///      by laying out words contiguously after the selector, then staticcall.
    ///      abi.encodePacked inlines fixed-size arrays word-by-word which is what
    ///      gnark's signature expects (no length prefix, no offset indirection).
    function _dispatch(
        BucketVerifier memory bv,
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        uint256[] memory publicInputs
    ) internal view returns (bool ok) {
        bytes memory cd = abi.encodePacked(
            bv.selector,
            proof,
            commitments,
            commitmentPok,
            publicInputs
        );
        (ok,) = bv.verifier.staticcall(cd);
    }

    function _writeLimbs(
        uint256[] memory out,
        uint256 offset,
        uint256 x
    ) internal pure returns (uint256) {
        out[offset + 0] = x & 0xffffffffffffffff;
        out[offset + 1] = (x >> 64) & 0xffffffffffffffff;
        out[offset + 2] = (x >> 128) & 0xffffffffffffffff;
        out[offset + 3] = (x >> 192) & 0xffffffffffffffff;
        return offset + 4;
    }

    function decodePoint(
        uint256 compressedPoint
    ) internal returns (uint256 xEdwards, uint256 yEdwards) {
        uint256 signBit = compressedPoint >> 255;
        uint256 yTwisted = compressedPoint & ~mask;

        // Twisted Edwards curve: -x^2 + y^2 = 1 + d*x^2*y^2
        // Solve for x^2: x^2 = (y^2 - 1) / (d*y^2 + 1)
        uint256 y2 = mulmod(yTwisted, yTwisted, p);
        uint256 x2 = mulmod(
            addmod(y2, a, p),
            pModInv(addmod(mulmod(d, y2, p), 1, p)),
            p
        );
        uint256 xTwisted = sqrtMod(x2);
        if (xTwisted % 2 != signBit) {
            xTwisted = p - xTwisted;
        }

        // Return Edwards coordinates (matching Go prover's utils.DecompressPoint)
        return (xTwisted, yTwisted);
    }

    /// @notice Reverse the byte order of a 256-bit value (convert between LE and BE).
    function reverseBytes(uint256 v) internal pure returns (uint256) {
        // Swap bytes pairwise
        v = ((v & 0xFF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00) >> 8) |
            ((v & 0x00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF) << 8);
        // Swap 2-byte pairs
        v = ((v & 0xFFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000) >> 16) |
            ((v & 0x0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF0000FFFF) << 16);
        // Swap 4-byte pairs
        v = ((v & 0xFFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000) >> 32) |
            ((v & 0x00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF00000000FFFFFFFF) << 32);
        // Swap 8-byte pairs
        v = ((v & 0xFFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF0000000000000000) >> 64) |
            ((v & 0x0000000000000000FFFFFFFFFFFFFFFF0000000000000000FFFFFFFFFFFFFFFF) << 64);
        // Swap 16-byte halves
        v = (v >> 128) | (v << 128);
        return v;
    }

    function sqrtMod(uint256 self) internal returns (uint256 result) {
        assembly ("memory-safe") {
            // load the free memory pointer value
            let pointer := mload(0x40)

            // Define length of base (Bsize)
            mstore(pointer, 0x20)
            // Define the exponent size (Esize)
            mstore(add(pointer, 0x20), 0x20)
            // Define the modulus size (Msize)
            mstore(add(pointer, 0x40), 0x20)
            // Define variables base (B)
            mstore(add(pointer, 0x60), self)
            // Define the exponent (E)
            mstore(add(pointer, 0x80), pp3div8)
            // We save the point of the last argument, it will be override by the result
            // of the precompile call in order to avoid paying for the memory expansion properly
            let _result := add(pointer, 0xa0)
            // Define the modulus (M)
            mstore(_result, p)

            // Call the precompiled ModExp (0x05) https://www.evm.codes/precompiled#0x05
            if iszero(
                call(
                    not(0), // amount of gas to send
                    MODEXP_PRECOMPILE, // target
                    0x00, // value in wei
                    pointer, // argsOffset
                    0xc0, // argsSize (6 * 32 bytes)
                    _result, // retOffset (we override M to avoid paying for the memory expansion)
                    0x20 // retSize (32 bytes)
                )
            ) {
                revert(0, 0)
            }

            result := mload(_result)
        }
        if (mulmod(result, result, p) != self) {
            result = mulmod(result, sqrtm1, p);
        }
        if (mulmod(result, result, p) != self) {
            revert();
        }
        return result;
    }

    /// @notice Calculate the modular inverse of a given integer, which is the inverse of this integer modulo p.
    /// @dev Uses the ModExp precompiled contract at address 0x05 for fast computation using little Fermat theorem
    /// @param self The integer of which to find the modular inverse
    /// @return result The modular inverse of the input integer. If the modular inverse doesn't exist, it revert the tx
    function pModInv(uint256 self) internal returns (uint256 result) {
        assembly ("memory-safe") {
            // load the free memory pointer value
            let pointer := mload(0x40)

            // Define length of base (Bsize)
            mstore(pointer, 0x20)
            // Define the exponent size (Esize)
            mstore(add(pointer, 0x20), 0x20)
            // Define the modulus size (Msize)
            mstore(add(pointer, 0x40), 0x20)
            // Define variables base (B)
            mstore(add(pointer, 0x60), self)
            // Define the exponent (E)
            mstore(add(pointer, 0x80), MINUS_2)
            // We save the point of the last argument, it will be override by the result
            // of the precompile call in order to avoid paying for the memory expansion properly
            let _result := add(pointer, 0xa0)
            // Define the modulus (M)
            mstore(_result, p)

            // Call the precompiled ModExp (0x05) https://www.evm.codes/precompiled#0x05
            if iszero(
                call(
                    not(0), // amount of gas to send
                    MODEXP_PRECOMPILE, // target
                    0x00, // value in wei
                    pointer, // argsOffset
                    0xc0, // argsSize (6 * 32 bytes)
                    _result, // retOffset (we override M to avoid paying for the memory expansion)
                    0x20 // retSize (32 bytes)
                )
            ) {
                revert(0, 0)
            }

            // we return the value in the last memory word created by the function
            result := mload(_result)
        }
    }
}
