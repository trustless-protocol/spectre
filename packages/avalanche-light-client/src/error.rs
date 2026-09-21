//! Client error type.

/// Every failure the Avalanche client surfaces to the host.
#[derive(Debug, thiserror::Error)]
pub enum Error {
    /// Shared EVM machinery (header shape, MPT proofs, packet paths, host
    /// envelopes) failed.
    #[error(transparent)]
    Evm(#[from] l2_client::error::Error),

    /// Warp aggregate verification failed.
    #[error(transparent)]
    Warp(#[from] avalanche_warp::Error),

    /// Wire (de)serialization failed.
    #[error(transparent)]
    Serde(#[from] serde_json::Error),

    /// Host protobuf envelope handling failed.
    #[error(transparent)]
    Proto(#[from] prost::EncodeError),

    /// CosmWasm standard-library failure (JSON responses, host API).
    #[error(transparent)]
    Std(#[from] cosmwasm_std::StdError),

    /// A structural invariant of the client or a message does not hold.
    #[error("invalid client input: {0}")]
    Invalid(&'static str),

    /// The client is frozen at the given height.
    #[error("client is frozen at height {0}")]
    Frozen(u64),

    /// A replacement validator set does not advance the P-Chain height.
    #[error("validator set rotation must advance the P-Chain height ({stored} -> {proposed})")]
    NonMonotonicRotation {
        /// Snapshot height of the pinned set.
        stored: u64,
        /// Snapshot height of the proposed set.
        proposed: u64,
    },

    /// The BLS querier could not be reached or answered malformed data.
    #[error("bls querier: {0}")]
    BlsQuerier(String),
}
