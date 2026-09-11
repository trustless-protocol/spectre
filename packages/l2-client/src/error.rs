//! Typed failures for shared L2 client operations.

/// Error returned while decoding, storing, or querying shared L2 client data.
#[derive(Debug, thiserror::Error)]
#[allow(missing_docs, clippy::module_name_repetitions)]
pub enum Error {
    #[error("CosmWasm operation failed: {0}")]
    Std(#[from] cosmwasm_std::StdError),
    #[error("client state is missing")]
    ClientStateMissing,
    #[error("consensus state is missing at height {0}")]
    ConsensusStateMissing(u64),
    #[error("invalid revision number {0}; only revision zero is supported")]
    InvalidRevision(u64),
    #[error("invalid zero L2 height")]
    ZeroHeight,
    #[error("L2 timestamp overflow")]
    TimestampOverflow,
    #[error("proof validation failed: {0}")]
    Proof(String),
    #[error("proof exceeds configured {limit} limit of {maximum}")]
    ProofLimit {
        /// Name of the bounded proof dimension.
        limit: &'static str,
        /// Largest permitted value for that dimension.
        maximum: usize,
    },
    #[error("invalid canonical EVM header: {0}")]
    InvalidEvmHeader(&'static str),
    #[error("invalid L2 header: {0}")]
    InvalidHeader(&'static str),
    #[error("state conflict at height {height}")]
    StateConflict {
        /// Height where the conflicting trusted blocks were observed.
        height: u64,
    },
    #[error("profile version is not supported by this artifact")]
    InvalidProfileVersion,
    #[error("L2 parent hash does not extend the trusted chain")]
    ParentHashMismatch,
    #[error("client is frozen at height {0}")]
    Frozen(u64),
    #[error(
        "non-zero delay is unsupported: delay_time_period={delay_time_period}, delay_block_period={delay_block_period}"
    )]
    UnsupportedNonZeroDelay {
        /// Requested time delay.
        delay_time_period: u64,
        /// Requested block delay.
        delay_block_period: u64,
    },
    #[error("lifecycle operation is unsupported: {operation}")]
    UnsupportedLifecycleOperation {
        /// Stable host operation name.
        operation: &'static str,
    },
    #[error("invalid fixture provenance: {0}")]
    InvalidFixtureProvenance(&'static str),
    #[error("fixture checksum does not match its provenance manifest")]
    FixtureChecksumMismatch,
    #[error("protobuf decoding failed: {0}")]
    ProtobufDecode(#[from] prost::DecodeError),
    #[error("JSON decoding failed: {0}")]
    JsonDecode(#[from] serde_json::Error),
    #[error("storage encoding failed: {0}")]
    StorageEncode(#[from] prost::EncodeError),
}
