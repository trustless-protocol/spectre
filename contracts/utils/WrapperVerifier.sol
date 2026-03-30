// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import {IVerifier} from "../interfaces/IVerifier.sol";
import {IGroth16Verifier} from "../interfaces/IVerifier.sol";
import {LibBytes} from "./LibBytes.sol";

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
    //sqrt of -1
    uint256 constant sqrtm1 =
        0x2b8324804fc1df0b2b4d00993dfbd7a72f431806ad2fe478c4ee1b274a0ea0b0;

    uint256 constant MINUS_2 =
        0x7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffeb;
    // -2 mod(n), used to speed up inversion operations
    uint256 constant MINUS_2MODN =
        0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3eb;
    // address of the ModExp precompiled contract (Arbitrary-precision exponentiation under modulo)
    address constant MODEXP_PRECOMPILE =
        0x0000000000000000000000000000000000000005;

    IGroth16Verifier public immutable GROTH16_VERFIER;

    constructor(IGroth16Verifier _groth16Verifier) {
        GROTH16_VERFIER = _groth16Verifier;
    }

    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        bytes32[2] calldata signature,
        bytes32 pubkey,
        bytes calldata message
    ) external override returns (bool) {
        // SHA-512 input uses raw Ed25519 little-endian bytes — correct as-is
        bytes32 R = signature[0];
        bytes memory hashData = new bytes(64 + message.length);
        assembly ("memory-safe") {
            mstore(add(hashData, 32), R)
            mstore(add(hashData, 64), pubkey)
            calldatacopy(add(hashData, 96), message.offset, message.length)
        }

        uint256[2] memory hashResult = sha512(hashData);
        // Ed25519 (RFC 8032) interprets SHA-512 output as 512-bit little-endian integer.
        // Reverse each 256-bit half and swap to get LE interpretation.
        uint256 lo = reverseBytes(hashResult[0]); // first 32 SHA bytes reversed = low 256 bits
        uint256 hi = reverseBytes(hashResult[1]); // last 32 SHA bytes reversed = high 256 bits
        uint256 H = red512Modq([hi, lo]);

        // Ed25519 encodes scalars and points as little-endian bytes.
        // Solidity bytes32→uint256 is big-endian. Reverse to get correct values.
        uint256 S = reverseBytes(uint256(signature[1]));

        // Decompress Ed25519 points (returns Edwards coordinates to match Go prover)
        (uint256 rX, uint256 rY) = decodePoint(reverseBytes(uint256(R)));
        (uint256 aX, uint256 aY) = decodePoint(reverseBytes(uint256(pubkey)));

        uint256[24] memory publicInputs;
        uint256 offset = 0;
        offset = _writeLimbs(publicInputs, offset, rX);
        offset = _writeLimbs(publicInputs, offset, rY);
        offset = _writeLimbs(publicInputs, offset, S);
        offset = _writeLimbs(publicInputs, offset, H);
        offset = _writeLimbs(publicInputs, offset, aX);
        offset = _writeLimbs(publicInputs, offset, aY);

        try
            GROTH16_VERFIER.verifyProof(
                proof,
                commitments,
                commitmentPok,
                publicInputs
            )
        {
            return true;
        } catch {
            return false;
        }
    }

    function _writeLimbs(
        uint256[24] memory out,
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

    function red512Modq(uint256[2] memory val) public pure returns (uint256) {
        //2²⁵²+27742317777372353535851937790883648493
        //0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed
        //0xffffffffffffffffffffffffffffffec6ef5bf4737dcf70d6ec31748d98951d is 2^256%N
        return
            addmod(
                mulmod(
                    val[0],
                    0xffffffffffffffffffffffffffffffec6ef5bf4737dcf70d6ec31748d98951d,
                    0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed
                ),
                val[1],
                0x1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed
            );
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
            //  result :=addmod(result,0,p)
        }
        if (mulmod(result, result, p) != self) {
            result = mulmod(result, sqrtm1, p);
        }
        if (mulmod(result, result, p) != self) {
            revert();
        }
        return result;
    }

    /// @notice Calculate the modular inverse of a given integer, which is the inverse of this integer modulo n.
    /// @dev Uses the ModExp precompiled contract at address 0x05 for fast computation using little Fermat theorem
    /// @param self The integer of which to find the modular inverse
    /// @return result The modular inverse of the input integer. If the modular inverse doesn't exist, it revert the tx

    function nModInv(uint256 self) internal returns (uint256 result) {
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
            mstore(add(pointer, 0x80), MINUS_2MODN)
            // We save the point of the last argument, it will be override by the result
            // of the precompile call in order to avoid paying for the memory expansion properly
            let _result := add(pointer, 0xa0)
            // Define the modulus (M)
            mstore(_result, N)

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

    function sha512(
        bytes memory message
    ) internal view returns (uint256[2] memory) {
        uint64[8] memory h = [
            0x6a09e667f3bcc908,
            0xbb67ae8584caa73b,
            0x3c6ef372fe94f82b,
            0xa54ff53a5f1d36f1,
            0x510e527fade682d1,
            0x9b05688c2b3e6c1f,
            0x1f83d9abfb41bd6b,
            0x5be0cd19137e2179
        ];
        sha2(message, h);
        return [
            uint256(
                bytes32(
                    abi.encodePacked(
                        bytes8(h[0]),
                        bytes8(h[1]),
                        bytes8(h[2]),
                        bytes8(h[3])
                    )
                )
            ),
            uint256(
                bytes32(
                    abi.encodePacked(
                        bytes8(h[4]),
                        bytes8(h[5]),
                        bytes8(h[6]),
                        bytes8(h[7])
                    )
                )
            )
        ];
    }

    function sha2(bytes memory message, uint64[8] memory h) internal pure {
        uint64[80] memory k = [
            0x428a2f98d728ae22,
            0x7137449123ef65cd,
            0xb5c0fbcfec4d3b2f,
            0xe9b5dba58189dbbc,
            0x3956c25bf348b538,
            0x59f111f1b605d019,
            0x923f82a4af194f9b,
            0xab1c5ed5da6d8118,
            0xd807aa98a3030242,
            0x12835b0145706fbe,
            0x243185be4ee4b28c,
            0x550c7dc3d5ffb4e2,
            0x72be5d74f27b896f,
            0x80deb1fe3b1696b1,
            0x9bdc06a725c71235,
            0xc19bf174cf692694,
            0xe49b69c19ef14ad2,
            0xefbe4786384f25e3,
            0x0fc19dc68b8cd5b5,
            0x240ca1cc77ac9c65,
            0x2de92c6f592b0275,
            0x4a7484aa6ea6e483,
            0x5cb0a9dcbd41fbd4,
            0x76f988da831153b5,
            0x983e5152ee66dfab,
            0xa831c66d2db43210,
            0xb00327c898fb213f,
            0xbf597fc7beef0ee4,
            0xc6e00bf33da88fc2,
            0xd5a79147930aa725,
            0x06ca6351e003826f,
            0x142929670a0e6e70,
            0x27b70a8546d22ffc,
            0x2e1b21385c26c926,
            0x4d2c6dfc5ac42aed,
            0x53380d139d95b3df,
            0x650a73548baf63de,
            0x766a0abb3c77b2a8,
            0x81c2c92e47edaee6,
            0x92722c851482353b,
            0xa2bfe8a14cf10364,
            0xa81a664bbc423001,
            0xc24b8b70d0f89791,
            0xc76c51a30654be30,
            0xd192e819d6ef5218,
            0xd69906245565a910,
            0xf40e35855771202a,
            0x106aa07032bbd1b8,
            0x19a4c116b8d2d0c8,
            0x1e376c085141ab53,
            0x2748774cdf8eeb99,
            0x34b0bcb5e19b48a8,
            0x391c0cb3c5c95a63,
            0x4ed8aa4ae3418acb,
            0x5b9cca4f7763e373,
            0x682e6ff3d6b2b8a3,
            0x748f82ee5defb2fc,
            0x78a5636f43172f60,
            0x84c87814a1f0ab72,
            0x8cc702081a6439ec,
            0x90befffa23631e28,
            0xa4506cebde82bde9,
            0xbef9a3f7b2c67915,
            0xc67178f2e372532b,
            0xca273eceea26619c,
            0xd186b8c721c0c207,
            0xeada7dd6cde0eb1e,
            0xf57d4f7fee6ed178,
            0x06f067aa72176fba,
            0x0a637dc5a2c898a6,
            0x113f9804bef90dae,
            0x1b710b35131c471b,
            0x28db77f523047d84,
            0x32caab7b40c72493,
            0x3c9ebe0a15c9bebc,
            0x431d67c49c100d4c,
            0x4cc5d4becb3e42b6,
            0x597f299cfc657e2a,
            0x5fcb6fab3ad6faec,
            0x6c44198c4a475817
        ];

        bytes memory padding = padMessage(message);
        require(padding.length % 128 == 0, "PADDING_ERROR");
        uint64[80] memory w;
        uint64[16] memory blocks;
        uint256 messageLength = (message.length / 128) * 128;
        unchecked {
            for (
                uint256 i = 0;
                i < (messageLength + padding.length);
                i += 128
            ) {
                if (i < messageLength) {
                    getBlock(message, blocks, i);
                } else {
                    getBlock(padding, blocks, i - messageLength);
                }
                for (uint256 j = 0; j < 16; ++j) {
                    w[j] = blocks[j];
                }

                for (uint256 j = 16; j < 80; ++j) {
                    w[j] =
                        gamma1(w[j - 2]) +
                        w[j - 7] +
                        gamma0(w[j - 15]) +
                        w[j - 16];
                }

                compress(h, k, w);
            }
        }
    }

    function compress(
        uint64[8] memory h,
        uint64[80] memory k,
        uint64[80] memory w
    ) internal pure {
        uint64[8] memory temp;
        for (uint256 j = 0; j < 8; ++j) {
            temp[j] = h[j];
        }
        // SHA-512 requires wrapping uint64 arithmetic; unchecked does NOT
        // propagate from the caller (sha2), so it must be declared here.
        unchecked {
            for (uint256 j = 0; j < 80; ++j) {
                uint64 t1 = temp[7] +
                    sigma1(temp[4]) +
                    ch(temp[4], temp[5], temp[6]) +
                    k[j] +
                    w[j];
                uint64 t2 = sigma0(temp[0]) + maj(temp[0], temp[1], temp[2]);

                temp[7] = temp[6];
                temp[6] = temp[5];
                temp[5] = temp[4];
                temp[4] = temp[3] + t1;
                temp[3] = temp[2];
                temp[2] = temp[1];
                temp[1] = temp[0];
                temp[0] = t1 + t2;
            }
            for (uint256 j = 0; j < 8; ++j) {
                h[j] += temp[j];
            }
        }
    }

    function padMessage(
        bytes memory message
    ) internal pure returns (bytes memory) {
        uint256 messageLength = message.length;
        bytes8 bitLength = bytes8(uint64(messageLength * 8));
        uint256 mdi = messageLength % 128;
        uint256 paddingLength;
        if (mdi < 112) {
            paddingLength = 119 - mdi;
        } else {
            paddingLength = 247 - mdi;
        }
        bytes memory padding = new bytes(paddingLength);
        bytes memory tail = LibBytes.slice(
            message,
            messageLength - mdi,
            messageLength
        );
        return abi.encodePacked(tail, bytes1(0x80), padding, bitLength);
    }

    function getBlock(
        bytes memory message,
        uint64[16] memory blocks,
        uint256 index
    ) internal pure {
        for (uint256 i = 0; i < 16; ++i) {
            blocks[i] = uint64(LibBytes.readBytes8(message, index + i * 8));
        }
    }

    function ch(uint64 x, uint64 y, uint64 z) internal pure returns (uint64) {
        return (x & y) ^ (~x & z);
    }

    function maj(uint64 x, uint64 y, uint64 z) internal pure returns (uint64) {
        return (x & y) ^ (x & z) ^ (y & z);
    }

    function sigma0(uint64 x) internal pure returns (uint64) {
        return (rotateRight(x, 28) ^ rotateRight(x, 34) ^ rotateRight(x, 39));
    }

    function sigma1(uint64 x) internal pure returns (uint64) {
        return (rotateRight(x, 14) ^ rotateRight(x, 18) ^ rotateRight(x, 41));
    }

    function gamma0(uint64 x) internal pure returns (uint64) {
        return (rotateRight(x, 1) ^ rotateRight(x, 8) ^ (x >> 7));
    }

    function gamma1(uint64 x) internal pure returns (uint64) {
        return (rotateRight(x, 19) ^ rotateRight(x, 61) ^ (x >> 6));
    }

    function rotateRight(uint64 x, uint64 n) internal pure returns (uint64) {
        return (x << (64 - n)) | (x >> n);
    }
}
