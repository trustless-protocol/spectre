// SPDX-License-Identifier: MIT
pragma solidity ^0.8.4;

/// @title IVerifier - Ed25519 signature verification via Groth16 proof
/// @notice Wrapper that recomputes the hash-aggregate witness commitment from
///         calldata and forwards that 32-byte digest, packed into two public
///         field elements, to the underlying Groth16 verifier selected by
///         bucket.
interface IVerifier {
    /// @notice Shared CanonicalVote fields the witness commitment binds across
    ///         every validator signing the same Tendermint block. As of #199 the
    ///         on-chain witness no longer rebuilds the full per-slot canonical
    ///         vote: it commits only the common prefix (Type|Height), whether the
    ///         round field is present, and the 32-byte block hash. Timestamp,
    ///         chainID and the BlockID part-set header are left out of the hash —
    ///         the in-circuit Ed25519 verify still binds the full signed bytes,
    ///         and chainID is enforced separately on-chain (header.chainId ==
    ///         clientState.chainId).
    struct SharedBlock {
        uint64 height;
        uint64 round;
        bytes32 blockIDHash;
    }

    /// @notice Verify a batch of Ed25519 signatures via a single Groth16 proof.
    /// @dev The relayer picks `bucket` (∈ {4,8,16}) as the smallest circuit
    ///      that fits the number of signers needed to reach 2/3 voting power. Slots
    ///      beyond the true signer count are padded with deterministic dummy
    ///      signatures whose `active[i] = false`. The on-circuit ECIP gate
    ///      zeroes their contribution; the on-chain quorum check skips them.
    ///
    ///      The verifier wrapper hashes the per-slot pubkeys + active flags plus
    ///      the shared common-prefix fields into a single SHA-256 digest. The
    ///      circuit recomputes the same digest internally and exposes it as
    ///      two 128-bit BN254 public inputs, so generated per-bucket
    ///      verifiers use a fixed `uint256[2]` public input ABI.
    /// @param bucket Validator-count bucket (selects which per-bucket Groth16Verifier to dispatch to)
    /// @param proof The Groth16 proof (8 uint256s: Ar, Bs, Krs)
    /// @param commitments The proof commitments (2 uint256s)
    /// @param commitmentPok The proof of knowledge for commitments (2 uint256s)
    /// @param pubkeys Per-slot Ed25519 compressed public key. Length == bucket.
    /// @param active Per-slot real-signer flag (true = real validator, false = dummy padding). Length == bucket.
    /// @param shared Common-prefix CanonicalVote fields shared across all slots.
    function verifyBatchProof(
        uint16 bucket,
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        bytes32[] calldata pubkeys,
        bool[] calldata active,
        SharedBlock calldata shared
    )
        external
        returns (bool);
}

/// @title IGroth16Verifier - Raw single-signature Groth16 proof verifier
/// @notice Legacy raw verifier interface for the non-batch verifier.
interface IGroth16Verifier {
    function verifyProof(
        uint256[8] calldata proof,
        uint256[2] calldata commitments,
        uint256[2] calldata commitmentPok,
        uint256[24] calldata input
    )
        external
        view;
}
