// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.28;

import { IICS02ClientMsgs } from "../../msgs/IICS02ClientMsgs.sol";

/// @title ICS07 Tendermint Messages
/// @author srdtrk
/// @notice Defines shared types for ICS07Tendermint implementations.
interface IICS07TendermintMsgs {
    /// @notice Fraction of validator overlap needed to update header
    /// @param numerator Numerator of the fraction
    /// @param denominator Denominator of the fraction
    struct TrustThreshold {
        uint8 numerator;
        uint8 denominator;
    }

    /// @notice Defines the ICS07Tendermint ClientState for ibc-lite
    /// @param chainId Chain ID
    /// @param trustLevel Fraction of validator overlap needed to update header
    /// @param latestHeight Latest height the client was updated to
    /// @param trustingPeriod duration of the period since the LatestTimestamp during which the
    /// submitted headers are valid for upgrade in seconds.
    /// @param unbondingPeriod duration of the staking unbonding period in seconds
    /// @param isFrozen whether or not client is frozen (due to misbehavior)
    /// @param zkAlgorithm The zk algorithm supported by this contract (for the relayers).
    struct ClientState {
        string chainId;
        TrustThreshold trustLevel;
        IICS02ClientMsgs.Height latestHeight;
        uint32 trustingPeriod;
        uint32 unbondingPeriod;
        bool isFrozen;
        SupportedZkAlgorithm zkAlgorithm;
    }

    /// @notice Defines the Tendermint light client's consensus state at some height.
    /// @param timestamp timestamp that corresponds to the counterparty block height
    /// in which the ConsensusState was generated. (in unix nanoseconds)
    /// @param root commitment root (i.e app hash)
    /// @param nextValidatorsHash next validators hash
    struct ConsensusState {
        uint128 timestamp;
        bytes32 root;
        bytes32 nextValidatorsHash;
    }

    /// @notice Defines the supported zk algorithms
    enum SupportedZkAlgorithm {
        Groth16,
        Plonk
    }

     struct Header {
        SignedHeader signedHeader;
        ValidatorSet validatorSet;
        IICS02ClientMsgs.Height trustedHeight;
        ValidatorSet trustedNextValidatorSet;
    }

    struct SignedHeader {
        BlockHeader header;
        BlockCommit commit;
    }

    struct ValidatorSet {
        ValidatorInfo[] validators;
        bool hasProposer;
        ValidatorInfo proposer;
        uint64 totalVotingPower;
    }

    struct ValidatorInfo {
        bytes valAddress;
        bytes32 pubKey;
        uint64 votingPower;
        int64 proposerPriority;
    }

    struct SimpleValidator {
        bytes32 pubKey;
        uint64 votingPower;
    }

    struct BlockHeader {
        Version version;
        string chainId;
        uint64 height;
        uint128 time;
        bool hasLastBlockId;
        BlockId lastBlockId;
        bool hasLastCommitHash;
        bytes32 lastCommitHash;
        bool hasDataHash;
        bytes32 dataHash;
        bytes32 validatorsHash;
        bytes32 nextValidatorsHash;
        bytes32 consensusHash;
        bytes32 appHash;
        bool hasLastResultsHash;
        bytes32 lastResultsHash;
        bool hasEvidenceHash;
        bytes32 evidenceHash;
        bytes proposerAddress;
    }

    struct BlockCommit {
        uint64 height;
        uint32 round;
        BlockId blockId;
        CommitSig[] commitSigs;
    }

    struct BlockId {
        bytes32 hashData;
        PartSetHeader partSetHeader;
    }

    struct PartSetHeader {
        uint32 total;
        bytes32 hashData;
    }

    enum CommitSigFlag {
        /// no vote was received from a validator.
        BLOCK_ID_FLAG_ABSENT,
        /// voted for the Commit.BlockID.
        BLOCK_ID_FLAG_COMMIT,
        /// voted for nil.
        BLOCK_ID_FLAG_NIL
    }

    struct CommitSigData {
        bytes validatorAddress;
        uint128 timestamp;
        bool hasSignature;
        bytes signature;
    }

    struct CommitSig {
        CommitSigFlag flag;
        CommitSigData data;
    }

    struct ChainId {
        string id;
        uint64 revisionNumber;
    }

    struct Options {
        TrustThreshold trustThreshold;
        uint32 trustingPeriod;
        uint32 clockDrift;
    }

    struct Version {
        uint64 blockVersion;
        uint64 appVersion;
    }

    struct ClientConsensusStatePath {
        string clientId;
        uint64 revisionNumber;
        uint64 revisionHeight;
    }

    struct UntrustedBlockState {
        SignedHeader signedHeader;
        ValidatorSet validatorSet;
    }

    struct TrustedBlockState{
        string chainId;
        uint128 headerTime;
        uint64 height;
        ValidatorSet nextValidatorSet;
        bytes32 nextValidatorHash;
    }

    struct VotingPowerTally {
        /// Total voting power
        uint64 total;
        /// Tallied voting power
        uint64 tallied;
        /// Trust threshold for voting power
        TrustThreshold trustThreshold;
    }

    struct NonAbsentCommitVote {
        SignedVote signedVote;
        /// Flag indicating whether the signature has already been verified.
        bool verified;
    }

    struct NonAbsentCommitVotes {
        /// Votes sorted by validator address.
        NonAbsentCommitVote[] votes;
        /// Internal buffer for storing sign_bytes.
        ///
        /// The buffer is reused for each canonical vote so that we allocate it
        /// once.
        bytes32 signBytes;
    }


    struct CanonicalVote {
        /// Type of vote (prevote or precommit)
        VoteType voteType;

        /// Block height
        uint64 height;

        /// Round
        uint64 round;

        /// Block ID
        BlockId blockId;

        /// Timestamp
        uint128 timestamp;

        /// Chain ID
        ChainId chainId;
    }

    struct SignedVote {
        CanonicalVote vote;
        bytes validatorAddress;
        bytes32 signature;
    }

    enum VoteType {
        /// Votes for blocks which validators observe are valid for a given round
        Prevote,

        /// Votes to commit to a particular block for a given round
        Precommit
    }

    enum Verdict {
        /// Verification succeeded, the block is valid.
        SUCCESS,
        /// The minimum voting power threshold is not reached,
        /// the block cannot be trusted yet.
        NOT_ENOUGH_TRUST,
        /// Verification failed, the block is invalid.
        INVALID
    }

}
