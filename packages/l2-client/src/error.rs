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
    #[error("duplicate storage proof for key {0}")]
    DuplicateStorageProof(String),
    #[error("invalid canonical EVM header: {0}")]
    InvalidEvmHeader(&'static str),
    #[error("invalid L2 header: {0}")]
    InvalidHeader(&'static str),
    #[error("state conflict at height {0}")]
    Conflict(u64),
    #[error("client is frozen at height {0}")]
    Frozen(u64),
    #[error("invalid fixture provenance: {0}")]
    InvalidFixtureProvenance(&'static str),
    #[error("fixture checksum does not match its provenance manifest")]
    FixtureChecksumMismatch,
    #[error("IBC host query failed: {0}")]
    HostQuery(String),
    #[error("IBC host returned an unexpected type URL: {0}")]
    UnexpectedHostTypeUrl(String),
    #[error("pinned Ethereum Wasm checksum does not match host state")]
    L1ChecksumMismatch,
    #[error("pinned Ethereum client is not active: {0}")]
    L1ClientInactive(String),
    #[error("Ethereum consensus state slot {actual} does not match requested slot {expected}")]
    L1ConsensusSlotMismatch { expected: u64, actual: u64 },
    #[error("protobuf decoding failed: {0}")]
    ProtobufDecode(#[from] prost::DecodeError),
    #[error("JSON decoding failed: {0}")]
    JsonDecode(#[from] serde_json::Error),
    #[error("storage encoding failed: {0}")]
    StorageEncode(#[from] prost::EncodeError),
}
