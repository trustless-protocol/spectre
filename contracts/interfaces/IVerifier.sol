// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

/// @title IVerifier - Ed25519 signature verification via Groth16 proof
/// @notice Wrapper that recomputes the hash-aggregate witness commitment from
///         calldata and forwards that 32-byte digest to the underlying Groth16
///         verifier selected by bucket.
interface IVerifier {
    /// @notice Shared CanonicalVote fields that are identical across every
    ///         validator signing the same Tendermint block.
    struct SharedBlock {
        uint64 height;
        uint64 round;
        bytes32 blockIDHash;
        uint32 partSetTotal;
        bytes32 partSetHash;
        bytes chainID;
    }

    /// @notice Verify a batch of Ed25519 signatures via a single Groth16 proof.
    /// @dev The relayer picks `bucket` (∈ {4,8,16,32,64,128}) as the smallest circuit
    ///      that fits the number of signers needed to reach 2/3 voting power. Slots
    ///      beyond the true signer count are padded with deterministic dummy
    ///      signatures whose `active[i] = false`. The on-circuit ECIP gate
    ///      zeroes their contribution; the on-chain quorum check skips them.
    ///
    ///      The verifier wrapper hashes all per-slot witness data plus the
    ///      shared canonical-vote fields into a single SHA-256 digest. The
    ///      circuit recomputes the same digest internally, so generated
    ///      per-bucket verifiers use a fixed `uint256[32]` public input ABI.
    /// @param bucket Validator-count bucket (selects which per-bucket Groth16Verifier to dispatch to)
    /// @param proof The Groth16 proof (8 uint256s: Ar, Bs, Krs)
    /// @param commitments The proof commitments (2 uint256s)
    /// @param commitmentPok The proof of knowledge for commitments (2 uint256s)
    /// @param signatures Per-slot Ed25519 signature: signatures[i][0] = R, signatures[i][1] = S. Length == bucket.
    /// @param pubkeys Per-slot Ed25519 compressed public key. Length == bucket.
    /// @param timestampSeconds Per-slot google.protobuf.Timestamp.seconds. Length == bucket.
    /// @param timestampNanos Per-slot google.protobuf.Timestamp.nanos. Length == bucket.
    /// @param active Per-slot real-signer flag (true = real validator, false = dummy padding). Length == bucket.
    /// @param shared CanonicalVote fields shared across all slots.
    function verifyBatchProof(
        uint16 bucket,
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        bytes32[2][] calldata signatures,
        bytes32[] calldata pubkeys,
        uint64[] calldata timestampSeconds,
        uint32[] calldata timestampNanos,
        bool[] calldata active,
        SharedBlock calldata shared
    ) external returns (bool);
}

/// @title IGroth16Verifier - Raw single-signature Groth16 proof verifier
/// @notice Legacy raw verifier interface for the non-batch verifier.
interface IGroth16Verifier {
    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        uint256[24] calldata input
    ) external view;
}
