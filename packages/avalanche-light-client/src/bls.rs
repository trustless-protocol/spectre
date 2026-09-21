//! BLS verification via the chain's custom querier.
//!
//! The wire shape is deliberately identical to the Ethereum light client's
//! custom query (`aggregate_verify { public_keys, message, signature }`): the
//! host chain routes one custom-query handler, Avalanche warp and Ethereum
//! consensus use the same BLS ciphersuite and domain-separation tag, and the
//! handler's fast-aggregate-verify takes arbitrary message bytes — so a chain
//! that already hosts the Ethereum client verifies warp aggregates with zero
//! chain-side changes.

use cosmwasm_std::{Binary, CustomQuery, QuerierWrapper, QueryRequest};
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::error::Error;

/// The custom BLS query, wire-compatible with the Ethereum client's.
#[derive(Clone, Debug, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum AvalancheCustomQuery {
    /// Verify an aggregate BLS signature by `public_keys` over `message`.
    AggregateVerify {
        /// 48-byte compressed public keys of the signers.
        public_keys: Vec<Binary>,
        /// The raw signed bytes (the unsigned warp message).
        message: Binary,
        /// The 96-byte aggregate signature.
        signature: Binary,
    },
}

impl CustomQuery for AvalancheCustomQuery {}

/// [`avalanche_warp::BlsVerify`] over the chain's custom querier.
pub struct QuerierBlsVerifier<'a> {
    /// The CosmWasm querier.
    pub querier: QuerierWrapper<'a, AvalancheCustomQuery>,
}

impl avalanche_warp::BlsVerify for QuerierBlsVerifier<'_> {
    fn fast_aggregate_verify(
        &self,
        public_keys: &[&[u8; avalanche_warp::PUBLIC_KEY_LEN]],
        message: &[u8],
        signature: &[u8; avalanche_warp::SIGNATURE_LEN],
    ) -> Result<(), avalanche_warp::Error> {
        let request: QueryRequest<AvalancheCustomQuery> =
            QueryRequest::Custom(AvalancheCustomQuery::AggregateVerify {
                public_keys: public_keys
                    .iter()
                    .map(|key| Binary::from(key.to_vec()))
                    .collect(),
                message: Binary::from(message.to_vec()),
                signature: Binary::from(signature.to_vec()),
            });
        // A transport failure and a negative verdict both fail closed; the
        // querier error text is unavailable through the fixed trait error, and
        // the host logs carry it.
        let valid: bool = self
            .querier
            .query(&request)
            .map_err(|_| avalanche_warp::Error::InvalidSignature)?;
        if valid {
            Ok(())
        } else {
            Err(avalanche_warp::Error::InvalidSignature)
        }
    }
}

/// Surfaces querier construction problems as client errors.
pub fn querier_error(context: &str) -> Error {
    Error::BlsQuerier(context.into())
}
