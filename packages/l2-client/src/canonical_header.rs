//! Strict canonical EVM execution-header encoding for pinned L2 fork configurations.

use alloy_primitives::{keccak256, Address, Bloom, Bytes, FixedBytes, B256, U256};
use alloy_rlp::Header as RlpHeader;
use schemars::JsonSchema;
use serde::{Deserialize, Serialize};

use crate::error::Error;

/// Fork-specific field set accepted for a canonical execution header.
#[derive(Clone, Copy, Debug, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(rename_all = "snake_case")]
pub enum ExecutionHeaderFork {
    /// EIP-1559 fields through the London fork.
    London,
    /// Adds the withdrawals root.
    Shanghai,
    /// Adds EIP-4844 blob fields and the parent beacon block root.
    Cancun,
    /// Adds the requests hash.
    Prague,
    /// Avalanche C-Chain (coreth) header: the fixed `ext_data_hash` field plus
    /// coreth's cascading optional tail. One variant covers every coreth
    /// upgrade (Granite, Helicon, ...) because presence of the tail fields is
    /// taken from the header itself and the RLP mirrors coreth's encoder,
    /// which writes an absent optional as an empty string whenever any later
    /// optional is present. Coreth headers never carry a withdrawals root or
    /// a requests hash.
    Coreth,
}

/// Complete EVM execution header used to bind rollup commitments to an L2 state root.
#[derive(Clone, Debug, Default, PartialEq, Eq, Serialize, Deserialize, JsonSchema)]
#[serde(deny_unknown_fields)]
pub struct CanonicalEvmHeader {
    /// Parent block hash.
    #[schemars(with = "String")]
    pub parent_hash: B256,
    /// Ommer-list hash.
    #[schemars(with = "String")]
    pub ommers_hash: B256,
    /// Block beneficiary.
    #[schemars(with = "String")]
    pub beneficiary: Address,
    /// L2 execution state root.
    #[schemars(with = "String")]
    pub state_root: B256,
    /// Transaction trie root.
    #[schemars(with = "String")]
    pub transactions_root: B256,
    /// Receipt trie root.
    #[schemars(with = "String")]
    pub receipts_root: B256,
    /// 256-byte log bloom.
    #[schemars(with = "String")]
    pub logs_bloom: Bloom,
    /// PoW difficulty retained in post-merge headers for canonical compatibility.
    #[schemars(with = "String")]
    pub difficulty: U256,
    /// L2 block number.
    pub number: u64,
    /// Block gas limit.
    pub gas_limit: u64,
    /// Gas used by the block.
    pub gas_used: u64,
    /// Unix timestamp in seconds.
    pub timestamp: u64,
    /// Opaque extra data, bounded to the Ethereum protocol maximum.
    #[schemars(with = "String")]
    pub extra_data: Bytes,
    /// Mix hash / post-merge prev-randao field.
    #[schemars(with = "String")]
    pub mix_hash: B256,
    /// Eight-byte nonce.
    #[schemars(with = "String")]
    pub nonce: FixedBytes<8>,
    /// EIP-1559 base fee.
    #[schemars(with = "Option<String>")]
    pub base_fee_per_gas: Option<U256>,
    /// Shanghai withdrawals root.
    #[schemars(with = "Option<String>")]
    pub withdrawals_root: Option<B256>,
    /// Cancun blob gas used.
    pub blob_gas_used: Option<u64>,
    /// Cancun excess blob gas.
    pub excess_blob_gas: Option<u64>,
    /// Cancun parent beacon block root.
    #[schemars(with = "Option<String>")]
    pub parent_beacon_block_root: Option<B256>,
    /// Prague requests hash.
    #[schemars(with = "Option<String>")]
    pub requests_hash: Option<B256>,

    // Coreth (Avalanche C-Chain) fields. All are None for every Ethereum-family
    // fork and skipped when serializing, so deployed OP/Base/Arbitrum clients —
    // which parse with deny_unknown_fields — never see them on the wire.
    /// Coreth atomic-transaction data hash (fixed field, present on every
    /// coreth header).
    #[schemars(with = "Option<String>")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub ext_data_hash: Option<B256>,
    /// Coreth atomic-transaction gas used (ApricotPhase4).
    #[schemars(with = "Option<String>")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub ext_data_gas_used: Option<U256>,
    /// Coreth block gas cost (ApricotPhase4).
    #[schemars(with = "Option<String>")]
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub block_gas_cost: Option<U256>,
    /// Coreth millisecond timestamp (Granite, ACP-226).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub time_milliseconds: Option<u64>,
    /// Coreth minimum delay excess (Granite, ACP-226).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub min_delay_excess: Option<u64>,
    /// Coreth gas target exponent (Helicon).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub target_exponent: Option<u64>,
    /// Coreth minimum price exponent (Helicon, ACP-283).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub min_price_exponent: Option<u64>,
    /// Coreth settled height (Helicon, ACP-194 streaming asynchronous execution).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub settled_height: Option<u64>,
    /// Coreth settled gas unix time (Helicon, ACP-194).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub settled_gas_unix: Option<u64>,
    /// Coreth settled gas numerator (Helicon, ACP-194).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub settled_gas_numerator: Option<u64>,
    /// Coreth settled excess (Helicon, ACP-194).
    #[serde(default, skip_serializing_if = "Option::is_none")]
    pub settled_excess: Option<u64>,
}

impl CanonicalEvmHeader {
    /// Ethereum's consensus maximum extra-data size.
    pub const MAX_EXTRA_DATA_BYTES: usize = 32;

    /// Returns the checked block number.
    pub const fn number(&self) -> u64 {
        self.number
    }

    /// Returns the checked Unix timestamp in seconds.
    pub const fn timestamp(&self) -> u64 {
        self.timestamp
    }

    /// Returns the execution state root committed by this header.
    pub const fn state_root(&self) -> B256 {
        self.state_root
    }

    /// Validates the fork-specific optional-field combination and fixed protocol limits.
    pub fn validate_for_fork(&self, fork: ExecutionHeaderFork) -> Result<(), Error> {
        if self.extra_data.len() > Self::MAX_EXTRA_DATA_BYTES {
            return Err(Error::InvalidEvmHeader("extra_data exceeds 32 bytes"));
        }
        if self.gas_used > self.gas_limit {
            return Err(Error::InvalidEvmHeader("gas_used exceeds gas_limit"));
        }

        let london = self.base_fee_per_gas.is_some();
        let shanghai = self.withdrawals_root.is_some();
        let cancun = self.blob_gas_used.is_some()
            && self.excess_blob_gas.is_some()
            && self.parent_beacon_block_root.is_some();
        let partial_cancun = self.blob_gas_used.is_some()
            || self.excess_blob_gas.is_some()
            || self.parent_beacon_block_root.is_some();

        let coreth_tail = self.ext_data_gas_used.is_some()
            || self.block_gas_cost.is_some()
            || self.time_milliseconds.is_some()
            || self.min_delay_excess.is_some()
            || self.target_exponent.is_some()
            || self.min_price_exponent.is_some()
            || self.settled_height.is_some()
            || self.settled_gas_unix.is_some()
            || self.settled_gas_numerator.is_some()
            || self.settled_excess.is_some();
        let coreth = self.ext_data_hash.is_some() || coreth_tail;

        let valid = match fork {
            ExecutionHeaderFork::London => {
                london && !shanghai && !partial_cancun && self.requests_hash.is_none() && !coreth
            }
            ExecutionHeaderFork::Shanghai => {
                london && shanghai && !partial_cancun && self.requests_hash.is_none() && !coreth
            }
            ExecutionHeaderFork::Cancun => {
                london && shanghai && cancun && self.requests_hash.is_none() && !coreth
            }
            ExecutionHeaderFork::Prague => {
                london && shanghai && cancun && self.requests_hash.is_some() && !coreth
            }
            // Coreth: the fixed ext_data_hash must be present; withdrawals and
            // requests roots never exist on the C-Chain; the blob triple is
            // all-or-none exactly as on Ethereum (ACP-131 activates all three).
            ExecutionHeaderFork::Coreth => {
                london
                    && self.ext_data_hash.is_some()
                    && !shanghai
                    && (cancun || !partial_cancun)
                    && self.requests_hash.is_none()
            }
        };

        if !valid {
            return Err(Error::InvalidEvmHeader(
                "optional fields do not match the pinned fork",
            ));
        }
        Ok(())
    }

    /// RLP-encodes this header in the field order selected by the pinned fork.
    pub fn rlp_bytes(&self, fork: ExecutionHeaderFork) -> Result<Vec<u8>, Error> {
        self.validate_for_fork(fork)?;
        if fork == ExecutionHeaderFork::Coreth {
            return self.coreth_rlp_bytes();
        }
        let mut fields = vec![
            alloy_rlp::encode(self.parent_hash),
            alloy_rlp::encode(self.ommers_hash),
            alloy_rlp::encode(self.beneficiary),
            alloy_rlp::encode(self.state_root),
            alloy_rlp::encode(self.transactions_root),
            alloy_rlp::encode(self.receipts_root),
            alloy_rlp::encode(self.logs_bloom),
            alloy_rlp::encode(self.difficulty),
            alloy_rlp::encode(self.number),
            alloy_rlp::encode(self.gas_limit),
            alloy_rlp::encode(self.gas_used),
            alloy_rlp::encode(self.timestamp),
            alloy_rlp::encode(self.extra_data.clone()),
            alloy_rlp::encode(self.mix_hash),
            alloy_rlp::encode(self.nonce),
            alloy_rlp::encode(
                self.base_fee_per_gas
                    .ok_or(Error::InvalidEvmHeader("missing base fee"))?,
            ),
        ];

        if matches!(
            fork,
            ExecutionHeaderFork::Shanghai
                | ExecutionHeaderFork::Cancun
                | ExecutionHeaderFork::Prague
        ) {
            fields.push(alloy_rlp::encode(
                self.withdrawals_root
                    .ok_or(Error::InvalidEvmHeader("missing withdrawals root"))?,
            ));
        }
        if matches!(
            fork,
            ExecutionHeaderFork::Cancun | ExecutionHeaderFork::Prague
        ) {
            fields.push(alloy_rlp::encode(
                self.blob_gas_used
                    .ok_or(Error::InvalidEvmHeader("missing blob gas used"))?,
            ));
            fields.push(alloy_rlp::encode(
                self.excess_blob_gas
                    .ok_or(Error::InvalidEvmHeader("missing excess blob gas"))?,
            ));
            fields
                .push(alloy_rlp::encode(self.parent_beacon_block_root.ok_or(
                    Error::InvalidEvmHeader("missing parent beacon root"),
                )?));
        }
        if matches!(fork, ExecutionHeaderFork::Prague) {
            fields.push(alloy_rlp::encode(
                self.requests_hash
                    .ok_or(Error::InvalidEvmHeader("missing requests hash"))?,
            ));
        }

        let payload_length = fields.iter().map(Vec::len).sum();
        let mut output = Vec::with_capacity(payload_length + 9);
        RlpHeader {
            list: true,
            payload_length,
        }
        .encode(&mut output);
        for field in fields {
            output.extend(field);
        }
        Ok(output)
    }

    /// RLP-encodes a coreth (Avalanche C-Chain) header: the 15 shared fields,
    /// the fixed `ext_data_hash`, then coreth's cascading optional tail. The
    /// tail mirrors coreth's generated encoder exactly: a field is emitted iff
    /// it or any later tail field is present, and an absent field emitted that
    /// way encodes as the RLP empty string (0x80). Trailing absent fields are
    /// omitted entirely.
    fn coreth_rlp_bytes(&self) -> Result<Vec<u8>, Error> {
        let mut fields = vec![
            alloy_rlp::encode(self.parent_hash),
            alloy_rlp::encode(self.ommers_hash),
            alloy_rlp::encode(self.beneficiary),
            alloy_rlp::encode(self.state_root),
            alloy_rlp::encode(self.transactions_root),
            alloy_rlp::encode(self.receipts_root),
            alloy_rlp::encode(self.logs_bloom),
            alloy_rlp::encode(self.difficulty),
            alloy_rlp::encode(self.number),
            alloy_rlp::encode(self.gas_limit),
            alloy_rlp::encode(self.gas_used),
            alloy_rlp::encode(self.timestamp),
            alloy_rlp::encode(self.extra_data.clone()),
            alloy_rlp::encode(self.mix_hash),
            alloy_rlp::encode(self.nonce),
            alloy_rlp::encode(
                self.ext_data_hash
                    .ok_or(Error::InvalidEvmHeader("missing ext data hash"))?,
            ),
        ];

        // Coreth optional tail, in coreth's declaration order. An entry is the
        // field's RLP when present, None when absent.
        let tail: Vec<Option<Vec<u8>>> = vec![
            self.base_fee_per_gas.map(|v| alloy_rlp::encode(v)),
            self.ext_data_gas_used.map(|v| alloy_rlp::encode(v)),
            self.block_gas_cost.map(|v| alloy_rlp::encode(v)),
            self.blob_gas_used.map(|v| alloy_rlp::encode(v)),
            self.excess_blob_gas.map(|v| alloy_rlp::encode(v)),
            self.parent_beacon_block_root.map(|v| alloy_rlp::encode(v)),
            self.time_milliseconds.map(|v| alloy_rlp::encode(v)),
            self.min_delay_excess.map(|v| alloy_rlp::encode(v)),
            self.target_exponent.map(|v| alloy_rlp::encode(v)),
            self.min_price_exponent.map(|v| alloy_rlp::encode(v)),
            self.settled_height.map(|v| alloy_rlp::encode(v)),
            self.settled_gas_unix.map(|v| alloy_rlp::encode(v)),
            self.settled_gas_numerator.map(|v| alloy_rlp::encode(v)),
            self.settled_excess.map(|v| alloy_rlp::encode(v)),
        ];
        let last_present = tail.iter().rposition(Option::is_some);
        if let Some(last) = last_present {
            for entry in &tail[..=last] {
                fields.push(entry.clone().unwrap_or_else(|| vec![alloy_rlp::EMPTY_STRING_CODE]));
            }
        }

        let payload_length = fields.iter().map(Vec::len).sum();
        let mut output = Vec::with_capacity(payload_length + 9);
        RlpHeader {
            list: true,
            payload_length,
        }
        .encode(&mut output);
        for field in fields {
            output.extend(field);
        }
        Ok(output)
    }

    /// Returns `keccak256(rlp(header))` for the pinned fork configuration.
    pub fn hash(&self, fork: ExecutionHeaderFork) -> Result<B256, Error> {
        Ok(keccak256(self.rlp_bytes(fork)?))
    }
}
