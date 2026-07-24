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
fn creation_stores_runtime_profile_and_consensus_without_queries_or_validation() {
    let mut deps = mock_dependencies();
    let client = ClientState {
        latest_height: 9,
        frozen_height: Some(8),
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
    };
    let info = message_info(&deps.api.addr_make("creator"), &[]);
    instantiate(
        deps.as_mut(),
        mock_env(),
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
    assert_eq!(
        l2_client::runtime::consensus_state(deps.as_ref().storage, 9).unwrap(),
        consensus
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
