//! A self-contained, signed Ethereum light client chain for tests.
//!
//! The recorded fixtures are signed by the real sync committee, so a test cannot choose how many
//! members sign an update. This chain is signed by a small committee whose secret keys the test
//! holds. The signatures are real BLS signatures (milagro) over the real signing root; only the
//! beacon data around them is synthetic.

use alloy_primitives::{Bytes, FixedBytes, B256};
use ethereum_types::consensus::{
    bls::{BlsPublicKey, BlsSignature},
    domain::{compute_domain, DomainType},
    fork::{Fork, ForkParameters, Version},
    light_client_header::{
        BeaconBlockHeader, ExecutionPayloadHeader, LightClientHeader, LightClientUpdate,
    },
    merkle::{floorlog2, EXECUTION_PAYLOAD_GINDEX, FINALIZED_ROOT_GINDEX_ELECTRA},
    signing_data::compute_signing_root,
    sync_committee::{SyncAggregate, SyncCommittee},
};
use sha2::{Digest, Sha256};
use tree_hash::TreeHash;

use crate::{
    client_state::ClientState, consensus_state::ConsensusState,
    sync_protocol_helpers::get_subtree_index, test_utils::bls_verifier::aggreagate,
};

/// The number of members of the synthetic sync committee.
pub const COMMITTEE_SIZE: usize = 8;
/// The slot of the consensus state the chain is anchored at.
pub const ANCHOR_SLOT: u64 = 100;
/// The finalized slot of every update the chain produces.
pub const FINALIZED_SLOT: u64 = 200;
/// The attested slot of every update the chain produces.
pub const ATTESTED_SLOT: u64 = 264;
/// The signature slot of every update the chain produces.
pub const SIGNATURE_SLOT: u64 = ATTESTED_SLOT + 1;

const SECONDS_PER_SLOT: u64 = 12;

/// A synthetic chain with an 8-member sync committee, all inside one sync committee period.
pub struct SyntheticChain {
    secret_keys: Vec<milagro_bls::SecretKey>,
    committee: SyncCommittee,
}

impl Default for SyntheticChain {
    fn default() -> Self {
        Self::new()
    }
}

impl SyntheticChain {
    /// Creates the chain with a deterministic committee.
    ///
    /// # Panics
    /// Panics if a deterministic key cannot be derived, which would be a bug in the test setup.
    #[must_use]
    pub fn new() -> Self {
        let secret_keys = (1..=COMMITTEE_SIZE)
            .map(|i| {
                milagro_bls::SecretKey::key_generate(&[u8::try_from(i).unwrap(); 32], &[]).unwrap()
            })
            .collect::<Vec<_>>();
        let pubkeys = secret_keys
            .iter()
            .map(|sk| FixedBytes::from(milagro_bls::PublicKey::from_secret_key(sk).as_bytes()))
            .collect::<Vec<BlsPublicKey>>();
        let aggregate_pubkey = aggreagate(&pubkeys).unwrap();

        Self {
            secret_keys,
            committee: SyncCommittee {
                pubkeys,
                aggregate_pubkey,
            },
        }
    }

    /// The full sync committee.
    #[must_use]
    pub const fn committee(&self) -> &SyncCommittee {
        &self.committee
    }

    /// A client state whose latest slot is [`ANCHOR_SLOT`].
    ///
    /// `min_sync_committee_participants` is set to a third of the committee, rounded up, which is
    /// what the relayer configures.
    #[must_use]
    pub fn client_state(&self) -> ClientState {
        let fork = |version: u8| Fork {
            version: Version::from([version, 0, 0, 0]),
            epoch: 0,
        };
        ClientState {
            chain_id: 1,
            genesis_validators_root: B256::repeat_byte(7),
            min_sync_committee_participants: (COMMITTEE_SIZE as u64).div_ceil(3),
            sync_committee_size: COMMITTEE_SIZE as u64,
            seconds_per_slot: SECONDS_PER_SLOT,
            slots_per_epoch: 32,
            epochs_per_sync_committee_period: 256,
            latest_slot: ANCHOR_SLOT,
            fork_parameters: ForkParameters {
                altair: fork(1),
                bellatrix: fork(2),
                capella: fork(3),
                deneb: fork(4),
                electra: fork(5),
                ..ForkParameters::default()
            },
            ..ClientState::default()
        }
    }

    /// The consensus state at [`ANCHOR_SLOT`], with the synthetic committee as the current one.
    #[must_use]
    pub fn consensus_state(&self) -> ConsensusState {
        ConsensusState {
            slot: ANCHOR_SLOT,
            state_root: B256::repeat_byte(1),
            timestamp: ANCHOR_SLOT * SECONDS_PER_SLOT,
            current_sync_committee: self.committee.to_summarized_sync_committee(),
            next_sync_committee: None,
        }
    }

    /// A timestamp after [`SIGNATURE_SLOT`], so no update is from the future.
    #[must_use]
    pub const fn current_timestamp(&self) -> u64 {
        (SIGNATURE_SLOT + 100) * SECONDS_PER_SLOT
    }

    /// An update attesting an execution state root, signed by the first `signers` committee members.
    ///
    /// Two updates that differ only in `execution_state_root` attest to conflicting execution
    /// states at [`ATTESTED_SLOT`].
    ///
    /// # Panics
    /// Panics if `signers` exceeds [`COMMITTEE_SIZE`].
    #[must_use]
    pub fn update(&self, execution_state_root: B256, signers: usize) -> LightClientUpdate {
        assert!(
            signers <= COMMITTEE_SIZE,
            "more signers than committee members"
        );

        let finalized_header = header(FINALIZED_SLOT, B256::ZERO, B256::repeat_byte(0xaa));
        let finality_branch = branch(floorlog2(FINALIZED_ROOT_GINDEX_ELECTRA));
        let attested_state_root = fold(
            finalized_header.beacon.tree_hash_root(),
            &finality_branch,
            get_subtree_index(FINALIZED_ROOT_GINDEX_ELECTRA),
        );
        let attested_header = header(ATTESTED_SLOT, attested_state_root, execution_state_root);

        let client_state = self.client_state();
        let fork_version = client_state
            .fork_parameters
            .compute_fork_version(client_state.compute_epoch_at_slot(SIGNATURE_SLOT - 1));
        let domain = compute_domain(
            DomainType::SYNC_COMMITTEE,
            Some(fork_version),
            Some(client_state.genesis_validators_root),
            client_state.fork_parameters.genesis_fork_version,
        );
        let signing_root = compute_signing_root(&attested_header.beacon, domain);

        let signatures = self.secret_keys[..signers]
            .iter()
            .map(|sk| milagro_bls::Signature::new(signing_root.as_slice(), sk))
            .collect::<Vec<_>>();
        let signature_refs = signatures.iter().collect::<Vec<_>>();
        let aggregate = milagro_bls::AggregateSignature::aggregate(&signature_refs);

        let bits = u8::try_from((1_u16 << signers) - 1).unwrap();

        LightClientUpdate {
            attested_header,
            next_sync_committee: None,
            next_sync_committee_branch: None,
            finalized_header,
            finality_branch,
            sync_aggregate: SyncAggregate {
                sync_committee_bits: Bytes::from(vec![bits]),
                sync_committee_signature: BlsSignature::from(aggregate.as_bytes()),
            },
            signature_slot: SIGNATURE_SLOT,
        }
    }
}

/// A pre-Gloas-format header whose body root commits to an execution header with the given state root.
fn header(slot: u64, state_root: B256, execution_state_root: B256) -> LightClientHeader {
    let execution = ExecutionPayloadHeader {
        state_root: execution_state_root,
        ..ExecutionPayloadHeader::default()
    };
    let execution_branch = branch(floorlog2(EXECUTION_PAYLOAD_GINDEX));
    let body_root = fold(
        execution.tree_hash_root(),
        &execution_branch,
        get_subtree_index(EXECUTION_PAYLOAD_GINDEX),
    );

    LightClientHeader {
        beacon: BeaconBlockHeader {
            slot,
            state_root,
            body_root,
            ..BeaconBlockHeader::default()
        },
        execution: Some(execution),
        execution_block_hash: None,
        execution_branch,
    }
}

/// An arbitrary but fixed Merkle branch of `depth` nodes.
fn branch(depth: usize) -> Vec<B256> {
    (1..=depth)
        .map(|i| B256::repeat_byte(u8::try_from(i).unwrap()))
        .collect()
}

/// Folds `leaf` up `branch`, the inverse of `validate_merkle_branch`.
fn fold(leaf: B256, branch: &[B256], index: u64) -> B256 {
    branch.iter().enumerate().fold(leaf, |value, (i, node)| {
        let mut hasher = Sha256::new();
        if (index >> i) & 1 == 1 {
            hasher.update(node);
            hasher.update(value);
        } else {
            hasher.update(value);
            hasher.update(node);
        }
        B256::from_slice(&hasher.finalize())
    })
}
