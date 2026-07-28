//! Direct Base creation tests.

use alloy_primitives::{Address, B256};
use base_verifier::config::{Config, OutputRootFormat};
use cosmwasm_std::{
    testing::{message_info, mock_dependencies, mock_env},
    Binary,
};
use l2_client::canonical_header::ExecutionHeaderFork;
use l2_client::msg::{IbcHeight, QueryMsg};
use l2_client::state::{ClientState, CommonProfile, ConsensusState, L1Client};

use crate::{
    contract::{instantiate, query},
    msg::InstantiateMsg,
};

#[test]
fn creation_stores_runtime_profile_and_normalizes_acceptance_time() {
    let mut deps = mock_dependencies();
    let client = ClientState {
        latest_height: 9,
        frozen_height: Some(8),
        finality_policy: l2_client::state::FinalityPolicy::default(),
        freshness_policy: l2_client::state::FreshnessPolicy::default(),
        last_finalized_update_at: None,
        profile: Config {
            common: CommonProfile {
                l1_chain_id: 0,
                l2_chain_id: 0,
                ethereum_client: L1Client {
                    client_id: String::new(),
                    wasm_checksum: vec![],
                },
                l2_router: Address::ZERO,
                commitment_slot: B256::ZERO,
                rollup_version: String::new(),
            },
            dispute_game_factory: Address::ZERO,
            game_list_slot: B256::ZERO,
            root_claim_bytecode_offset: 0,
            output_root_format: OutputRootFormat::LegacyV0,
            l2_header_fork: ExecutionHeaderFork::Prague,
        },
    };
    let consensus = ConsensusState {
        state_root: B256::ZERO,
        ibc_storage_root: B256::ZERO,
        timestamp_nanos: 12,
        l2_height: 9,
        l2_block_hash: B256::ZERO,
        parent_hash: B256::ZERO,
        l1_origin_number: 0,
        l1_origin_hash: B256::ZERO,
        finality_level: l2_client::state::FinalityLevel::Unsafe,
        proposal_status: l2_client::state::ProposalStatus::Pending,
        first_accepted_at: 0,
        finality_reached_at: 0,
        evidence_hash: B256::ZERO,
        rollup_commitment: B256::ZERO,
    };
    let info = message_info(&deps.api.addr_make("creator"), &[]);
    let env = mock_env();
    let accepted_at = env.block.time.seconds();
    instantiate(
        deps.as_mut(),
        env,
        info,
        InstantiateMsg {
            client_state: Binary::from(serde_json::to_vec(&client).unwrap()),
            consensus_state: Binary::from(serde_json::to_vec(&consensus).unwrap()),
            checksum: Binary::from(vec![7; 32]),
        },
    )
    .unwrap();
    assert_eq!(
        l2_client::runtime::client_state::<Config>(deps.as_ref().storage).unwrap(),
        client
    );
    let mut expected_consensus = consensus;
    expected_consensus.first_accepted_at = accepted_at;
    expected_consensus.finality_reached_at = accepted_at;
    assert_eq!(
        l2_client::runtime::consensus_state(deps.as_ref().storage, 9).unwrap(),
        expected_consensus
    );
    assert_eq!(
        query(
            deps.as_ref(),
            mock_env(),
            QueryMsg::TimestampAtHeight {
                height: IbcHeight {
                    revision_number: 0,
                    revision_height: 9,
                },
            },
        )
        .unwrap(),
        Binary::from(br#"{"timestamp":12}"#.as_slice())
    );
    assert_eq!(
        query(deps.as_ref(), mock_env(), QueryMsg::Status {}).unwrap(),
        Binary::from(br#"{"status":"Frozen"}"#.as_slice())
    );
}
