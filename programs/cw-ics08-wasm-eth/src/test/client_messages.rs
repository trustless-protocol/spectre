//! End-to-end checks, through the contract entry points, of which consensus state a client message
//! may be verified against and how much of the sync committee must sign misbehaviour evidence.

use alloy_primitives::B256;
use cosmwasm_std::{
    coins, from_json,
    testing::{message_info, mock_env, MockApi, MockQuerier, MockStorage},
    Binary, Env, OwnedDeps, Storage, Timestamp,
};
use ethereum_light_client::{
    consensus_state::ConsensusState,
    error::EthereumIBCError,
    header::{ActiveSyncCommittee, Header},
    test_utils::synthetic::{SyntheticChain, ANCHOR_SLOT, COMMITTEE_SIZE, FINALIZED_SLOT_STEP},
};

use crate::{
    contract::{instantiate, query, sudo},
    custom_query::EthereumCustomQuery,
    msg::{
        CheckForMisbehaviourMsg, CheckForMisbehaviourResult, EthereumMisbehaviourMsg,
        InstantiateMsg, QueryMsg, SudoMsg, UpdateStateMsg, VerifyClientMessageMsg,
    },
    state::{consensus_db_key, get_eth_client_state},
    test::helpers::mk_deps,
    ContractError,
};

type Deps = OwnedDeps<MockStorage, MockApi, MockQuerier<EthereumCustomQuery>, EthereumCustomQuery>;

/// The slots a chain of two updates on top of the bootstrap consensus state ends up at.
const SECOND_SLOT: u64 = ANCHOR_SLOT + FINALIZED_SLOT_STEP;
const THIRD_SLOT: u64 = SECOND_SLOT + FINALIZED_SLOT_STEP;

struct Setup {
    chain: SyntheticChain,
    deps: Deps,
    env: Env,
}

impl Setup {
    /// A client bootstrapped on the synthetic chain, with the clock past every slot it uses.
    fn new() -> Self {
        let chain = SyntheticChain::new();
        let mut deps = mk_deps();
        let mut env = mock_env();
        env.block.time = Timestamp::from_seconds(chain.current_timestamp());
        let creator = deps.api.addr_make("creator");

        instantiate(
            deps.as_mut(),
            env.clone(),
            message_info(&creator, &coins(1, "uatom")),
            InstantiateMsg {
                client_state: serde_json::to_vec(&chain.client_state()).unwrap().into(),
                consensus_state: serde_json::to_vec(&chain.consensus_state()).unwrap().into(),
                checksum: b"checksum".into(),
            },
        )
        .unwrap();

        Self { chain, deps, env }
    }

    /// A client that has advanced twice, so its latest consensus state is at `THIRD_SLOT` and the
    /// bootstrap one at `ANCHOR_SLOT` is two updates old.
    fn advanced() -> Self {
        let mut setup = Self::new();
        for trusted_slot in [ANCHOR_SLOT, SECOND_SLOT] {
            let header = setup.header(trusted_slot, B256::repeat_byte(0x11));
            setup.verify_header(&header).unwrap();
            setup.update_state(&header).unwrap();
        }
        assert_eq!(setup.latest_slot(), THIRD_SLOT);
        setup
    }

    fn header(&self, trusted_slot: u64, execution_state_root: B256) -> Header {
        Header {
            active_sync_committee: ActiveSyncCommittee::Current(self.chain.committee().clone()),
            consensus_update: self.chain.update_after(
                trusted_slot,
                execution_state_root,
                COMMITTEE_SIZE,
            ),
            trusted_slot,
            execution_header_rlp: None,
        }
    }

    /// Misbehaviour evidence anchored at `trusted_slot`, with each update signed by `signers`.
    fn evidence(&self, trusted_slot: u64, signers: usize) -> EthereumMisbehaviourMsg {
        EthereumMisbehaviourMsg {
            trusted_slot,
            sync_committee: ActiveSyncCommittee::Current(self.chain.committee().clone()),
            update_1: self
                .chain
                .update_after(trusted_slot, B256::repeat_byte(0x11), signers),
            update_2: self
                .chain
                .update_after(trusted_slot, B256::repeat_byte(0x22), signers),
        }
    }

    fn latest_slot(&self) -> u64 {
        get_eth_client_state(&self.deps.storage)
            .unwrap()
            .latest_slot
    }

    fn consensus_state(&self, slot: u64) -> ConsensusState {
        let any = self
            .deps
            .storage
            .get(consensus_db_key(slot).as_bytes())
            .unwrap();
        crate::state::get_eth_consensus_state(&self.deps.storage, slot).unwrap_or_else(|_| {
            panic!("no consensus state at {slot}: {any:?}");
        })
    }

    fn query(&self, msg: QueryMsg) -> Result<Binary, ContractError> {
        query(self.deps.as_ref(), self.env.clone(), msg)
    }

    fn verify_header(&self, header: &Header) -> Result<Binary, ContractError> {
        self.query(QueryMsg::VerifyClientMessage(VerifyClientMessageMsg {
            client_message: serde_json::to_vec(header).unwrap().into(),
        }))
    }

    fn verify_evidence(&self, evidence: &EthereumMisbehaviourMsg) -> Result<Binary, ContractError> {
        self.query(QueryMsg::VerifyClientMessage(VerifyClientMessageMsg {
            client_message: serde_json::to_vec(evidence).unwrap().into(),
        }))
    }

    fn check_evidence(
        &self,
        evidence: &EthereumMisbehaviourMsg,
    ) -> Result<CheckForMisbehaviourResult, ContractError> {
        self.query(QueryMsg::CheckForMisbehaviour(CheckForMisbehaviourMsg {
            client_message: serde_json::to_vec(evidence).unwrap().into(),
        }))
        .map(|bz| from_json(bz).unwrap())
    }

    fn update_state(&mut self, header: &Header) -> Result<(), ContractError> {
        sudo(
            self.deps.as_mut(),
            self.env.clone(),
            SudoMsg::UpdateState(UpdateStateMsg {
                client_message: serde_json::to_vec(header).unwrap().into(),
            }),
        )
        .map(|_| ())
    }
}

#[test]
fn a_chain_of_headers_each_trusting_the_previous_one_is_accepted() {
    // This is how the relayer submits a period crossing: every header trusts the finalized slot of
    // the one before it, which is the latest consensus state by the time it is verified.
    let setup = Setup::advanced();

    assert_eq!(setup.consensus_state(SECOND_SLOT).slot, SECOND_SLOT);
    assert_eq!(setup.consensus_state(THIRD_SLOT).slot, THIRD_SLOT);
}

#[test]
fn a_header_trusting_an_older_consensus_state_is_rejected() {
    let setup = Setup::advanced();
    // Valid against the bootstrap consensus state, and signed by its whole committee.
    let stale = setup.header(ANCHOR_SLOT, B256::repeat_byte(0x33));

    assert!(matches!(
        setup.verify_header(&stale),
        Err(ContractError::TrustedSlotNotLatest {
            trusted_slot: ANCHOR_SLOT,
            latest_slot: THIRD_SLOT,
        })
    ));
}

#[test]
fn a_stale_header_cannot_overwrite_a_stored_consensus_state() {
    let mut setup = Setup::advanced();
    let key = consensus_db_key(SECOND_SLOT);
    let before = setup.deps.storage.get(key.as_bytes()).unwrap();
    // It finalizes SECOND_SLOT again, but with another execution state root.
    let stale = setup.header(ANCHOR_SLOT, B256::repeat_byte(0x33));

    assert!(matches!(
        setup.update_state(&stale),
        Err(ContractError::TrustedSlotNotLatest { .. })
    ));

    assert_eq!(setup.deps.storage.get(key.as_bytes()).unwrap(), before);
    assert_eq!(setup.latest_slot(), THIRD_SLOT);
}

#[test]
fn misbehaviour_anchored_at_an_older_consensus_state_is_rejected() {
    let setup = Setup::advanced();
    // Fully signed evidence that is valid against the bootstrap consensus state.
    let stale = setup.evidence(ANCHOR_SLOT, COMMITTEE_SIZE);

    assert!(matches!(
        setup.check_evidence(&stale),
        Err(ContractError::TrustedSlotNotLatest {
            trusted_slot: ANCHOR_SLOT,
            latest_slot: THIRD_SLOT,
        })
    ));
    assert!(matches!(
        setup.verify_evidence(&stale),
        Err(ContractError::TrustedSlotNotLatest { .. })
    ));
}

#[test]
fn supermajority_misbehaviour_anchored_at_the_latest_consensus_state_is_found() {
    let setup = Setup::advanced();
    let evidence = setup.evidence(THIRD_SLOT, 6);

    setup.verify_evidence(&evidence).unwrap();
    assert!(setup.check_evidence(&evidence).unwrap().found_misbehaviour);
}

#[test]
fn minority_misbehaviour_cannot_freeze_the_client() {
    let setup = Setup::advanced();
    // 3 of 8 meets min_sync_committee_participants, the only floor the evidence used to face.
    let evidence = setup.evidence(THIRD_SLOT, 3);

    assert!(matches!(
        setup.check_evidence(&evidence),
        Err(ContractError::VerifyClientMessageFailed(
            EthereumIBCError::NotEnoughSignatures
        ))
    ));
    assert!(matches!(
        setup.verify_evidence(&evidence),
        Err(ContractError::VerifyClientMessageFailed(
            EthereumIBCError::NotEnoughSignatures
        ))
    ));
}
