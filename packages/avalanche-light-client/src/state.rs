//! Client and consensus state models.

use alloy_primitives::{Address, B256};
use cosmwasm_std::Binary;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::error::Error;

/// Hard cap on stored validators — far above the primary network (~600) but a
/// bound against unbounded state growth through a hostile rotation payload.
pub const MAX_VALIDATORS: usize = 4096;

/// One pinned validator: compressed BLS key + stake weight, in the canonical
/// order the P-Chain serves (the signer bitset indexes into this order).
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct WireValidator {
    /// 48-byte compressed BLS12-381 G1 public key.
    pub public_key: Binary,
    /// Stake weight.
    pub weight: u64,
}

/// The pinned canonical validator-set snapshot.
///
/// `total_weight` may exceed the sum of the listed weights: validators without
/// a registered BLS key are excluded from the list but still dilute the
/// quorum, exactly as in avalanchego's canonical-set construction. The stored
/// order is trusted as canonical — a wrong order or membership fails closed,
/// because the bitset then selects keys whose aggregate cannot verify.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ValidatorSetState {
    /// Canonically ordered validators with BLS keys.
    pub validators: Vec<WireValidator>,
    /// Total stake weight of the snapshot, keyless validators included.
    pub total_weight: u64,
    /// P-Chain height the snapshot was taken at; rotations must increase it.
    pub p_chain_height: u64,
}

impl ValidatorSetState {
    /// Validates size, key lengths, and weight accounting.
    pub fn validate(&self) -> Result<(), Error> {
        if self.validators.is_empty() {
            return Err(Error::Invalid("validator set must not be empty"));
        }
        if self.validators.len() > MAX_VALIDATORS {
            return Err(Error::Invalid("validator set exceeds the size cap"));
        }
        let mut sum: u64 = 0;
        for validator in &self.validators {
            if validator.public_key.len() != avalanche_warp::PUBLIC_KEY_LEN {
                return Err(Error::Invalid("validator key must be 48 bytes"));
            }
            if validator.weight == 0 {
                return Err(Error::Invalid("validator weight must be non-zero"));
            }
            sum = sum
                .checked_add(validator.weight)
                .ok_or(Error::Invalid("validator weights overflow"))?;
        }
        if self.total_weight < sum {
            return Err(Error::Invalid(
                "total weight is below the sum of listed validators",
            ));
        }
        Ok(())
    }

    /// Converts to the verification crate's borrowed form.
    pub fn to_warp_validators(&self) -> Result<Vec<avalanche_warp::Validator>, Error> {
        self.validators
            .iter()
            .map(|validator| {
                Ok(avalanche_warp::Validator {
                    public_key: validator
                        .public_key
                        .as_slice()
                        .try_into()
                        .map_err(|_| Error::Invalid("validator key must be 48 bytes"))?,
                    weight: validator.weight,
                })
            })
            .collect()
    }
}

/// Active Avalanche client state.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ClientState {
    /// Highest verified C-Chain height.
    pub latest_height: u64,
    /// Set when misbehaviour froze the client.
    pub frozen_height: Option<u64>,
    /// Avalanche network id (1 mainnet, 5 Fuji, 12345 local) — replay domain
    /// of every warp signature.
    pub network_id: u32,
    /// 32-byte source blockchain id of the C-Chain (P-Chain create-tx id).
    pub source_chain_id: Binary,
    /// EVM chain id the header stream must belong to (43114/43113/43112).
    pub evm_chain_id: u64,
    /// ICS26Router address on the C-Chain.
    #[schemars(with = "String")]
    pub router: Address,
    /// Solidity commitment mapping slot in the router.
    #[schemars(with = "String")]
    pub commitment_slot: B256,
    /// Stake-quorum numerator (Avalanche convention: 67).
    pub quorum_num: u64,
    /// Stake-quorum denominator (Avalanche convention: 100).
    pub quorum_den: u64,
    /// Pinned canonical validator set.
    pub validator_set: ValidatorSetState,
}

impl ClientState {
    /// Validates bootstrap and loaded invariants.
    pub fn validate(&self) -> Result<(), Error> {
        l2_client::state::Height::new(self.latest_height).map_err(Error::Evm)?;
        if self.frozen_height == Some(0) {
            return Err(Error::Invalid("frozen height must be non-zero when set"));
        }
        if self.network_id == 0 {
            return Err(Error::Invalid("network id must be non-zero"));
        }
        if self.source_chain_id.len() != 32 || self.source_chain_id.as_slice() == [0_u8; 32] {
            return Err(Error::Invalid(
                "source chain id must be 32 non-zero bytes",
            ));
        }
        if self.evm_chain_id == 0 {
            return Err(Error::Invalid("EVM chain id must be non-zero"));
        }
        if self.router.is_zero() {
            return Err(Error::Invalid("router address must be non-zero"));
        }
        if self.quorum_den == 0 || self.quorum_num == 0 || self.quorum_num > self.quorum_den {
            return Err(Error::Invalid("quorum must satisfy 0 < num <= den"));
        }
        self.validator_set.validate()
    }

    /// The 32-byte source chain id as a fixed array.
    pub fn source_chain_id_bytes(&self) -> Result<[u8; 32], Error> {
        self.source_chain_id
            .as_slice()
            .try_into()
            .map_err(|_| Error::Invalid("source chain id must be 32 bytes"))
    }
}

/// State retained for one verified C-Chain block. The field set matches the
/// shared attested-client consensus state so packet verification (storage-root
/// proofs) and host semantics stay identical; `height`/`block_hash` name the
/// chain's own terms.
#[derive(Clone, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct ConsensusState {
    /// Execution state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// Router account storage root (packet proofs verify against this).
    #[schemars(with = "String")]
    pub ibc_storage_root: B256,
    /// Header timestamp in nanoseconds.
    pub timestamp_nanos: u64,
    /// C-Chain block number.
    pub height: u64,
    /// Coreth block hash the quorum signed.
    #[schemars(with = "String")]
    pub block_hash: B256,
    /// Parent block hash, for adjacent-height linkage checks.
    #[schemars(with = "String")]
    pub parent_hash: B256,
    /// Host time (seconds) this state was first accepted.
    pub first_accepted_at: u64,
}

impl ConsensusState {
    /// Validates stored-state invariants.
    pub fn validate(&self) -> Result<(), Error> {
        l2_client::state::Height::new(self.height).map_err(Error::Evm)?;
        if self.state_root.is_zero() || self.ibc_storage_root.is_zero() || self.block_hash.is_zero()
        {
            return Err(Error::Invalid(
                "consensus roots and block hash must be non-zero",
            ));
        }
        Ok(())
    }

    /// Returns whether another state disagrees on block identity.
    #[must_use]
    pub fn conflicts_with(&self, other: &Self) -> bool {
        self.block_hash != other.block_hash
            || self.parent_hash != other.parent_hash
            || self.state_root != other.state_root
            || self.ibc_storage_root != other.ibc_storage_root
    }
}

/// A verified, normalized update (produced only by verification).
#[derive(Clone, Debug, PartialEq, Eq)]
pub struct Header {
    /// C-Chain block number.
    pub height: u64,
    /// Execution state root.
    pub state_root: B256,
    /// Router storage root proven against `state_root`.
    pub router_storage_root: B256,
    /// Header timestamp in seconds.
    pub timestamp_seconds: u64,
    /// Coreth block hash the quorum signed.
    pub block_hash: B256,
    /// Parent block hash.
    pub parent_hash: B256,
}

impl Header {
    /// Converts the verified header into stored consensus state.
    pub fn consensus_state(&self, accepted_at: u64) -> Result<ConsensusState, Error> {
        let state = ConsensusState {
            state_root: self.state_root,
            ibc_storage_root: self.router_storage_root,
            timestamp_nanos: self
                .timestamp_seconds
                .checked_mul(1_000_000_000)
                .ok_or(Error::Invalid("timestamp overflow"))?,
            height: self.height,
            block_hash: self.block_hash,
            parent_hash: self.parent_hash,
            first_accepted_at: accepted_at,
        };
        state.validate()?;
        Ok(state)
    }

    /// Checks for a same-height conflict between two verified headers.
    pub fn conflicts_with(&self, other: &Self) -> Result<bool, Error> {
        Ok(self.height == other.height
            && self
                .consensus_state(0)?
                .conflicts_with(&other.consensus_state(0)?))
    }
}
