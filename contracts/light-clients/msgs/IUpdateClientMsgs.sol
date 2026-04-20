// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import { IGroth16Msgs } from "./IGroth16Msgs.sol";
import { IICS07TendermintMsgs } from "./IICS07TendermintMsgs.sol";
import { IICS02ClientMsgs } from "../../msgs/IICS02ClientMsgs.sol";

/// @title Update Client Program Messages
/// @author srdtrk
/// @notice Defines shared types for the update client program.
interface IUpdateClientMsgs {
    /// @notice The message that is submitted to the updateClient function.
    /// @param clientState The client state.
    /// @param trustedConsensusState The trusted consensus state.
    /// @param proposedHeader The proposed header with validator signatures.
    /// @param time The current time in unix nanoseconds.
    /// @param proof The Groth16 proof for the batched Ed25519 signature verification (8 uint256s).
    /// @param commitments The proof commitments (2 uint256s).
    /// @param commitmentPok The proof of knowledge for commitments (2 uint256s).
    /// @param bucket The validator-count bucket this proof was generated against; selects the per-bucket verifier on-chain.
    /// @param signerIndices Validator indices in proposedHeader.validatorSet. Length == bucket.
    /// @param signatures (R || S) bytes per slot. signatures[i][0] = R, signatures[i][1] = S. Length == bucket.
    /// @param signerPubkeys Compressed Ed25519 public keys per slot, matching signerIndices. Length == bucket.
    /// @param signMessages Canonical vote sign-bytes per slot. Length == bucket.
    struct MsgUpdateClient {
        IICS07TendermintMsgs.ClientState clientState;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState;
        IICS07TendermintMsgs.Header proposedHeader;
        uint128 time;
        uint256[8] proof;
        uint256[2] commitments;
        uint256[2] commitmentPok;
        uint16 bucket;
        uint32[] signerIndices;
        bytes32[2][] signatures;
        bytes32[] signerPubkeys;
        bytes[] signMessages;
    }

    /// @notice The public value output for the gnark update client program.
    /// @param clientState The client state that was used to verify the header.
    /// @param trustedConsensusState The trusted consensus state.
    /// @param newConsensusState The new consensus state with the verified header.
    /// @param time The time which the header was verified in unix nanoseconds.
    /// @param trustedHeight The trusted height.
    /// @param newHeight The new height.
    struct UpdateClientOutput {
        IICS07TendermintMsgs.ClientState clientState;
        IICS07TendermintMsgs.ConsensusState trustedConsensusState;
        IICS07TendermintMsgs.ConsensusState newConsensusState;
        uint128 time;
        IICS02ClientMsgs.Height trustedHeight;
        IICS02ClientMsgs.Height newHeight;
    }
}
