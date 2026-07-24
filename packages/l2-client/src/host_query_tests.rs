//! Unit tests for pinned-Ethereum host-response validation.

use alloy_primitives::B256;
use ethereum_light_client::consensus_state::ConsensusState as EthereumConsensusState;
use ethereum_types::consensus::sync_committee::SummarizedSyncCommittee;
use ibc_proto::{
    google::protobuf::Any,
    ibc::{
        core::client::v1::{
            QueryClientStateResponse, QueryClientStatusResponse, QueryConsensusStateResponse,
        },
        lightclients::wasm::v1::{
            ClientState as WasmClientState, ConsensusState as WasmConsensusState,
        },
    },
};
use prost::Message;

use crate::{
    error::Error,
    host_query::{decode_consensus_state, validate_client_state, validate_client_status},
    state::L1Client,
};

#[test]
fn rejects_missing_pinned_client_state() {
    let error = validate_client_state(
        QueryClientStateResponse {
            client_state: None,
            proof: vec![],
            proof_height: None,
        },
        &L1Client {
            client_id: "08-wasm-0".into(),
            wasm_checksum: vec![7; 32],
        },
    )
    .unwrap_err();
    assert!(matches!(error, Error::HostQuery(_)));
}

#[test]
fn accepts_matching_pinned_client_state() {
    let l1 = pinned_l1();
    validate_client_state(
        QueryClientStateResponse {
            client_state: Some(Any {
                type_url: "/ibc.lightclients.wasm.v1.ClientState".into(),
                value: WasmClientState {
                    data: vec![],
                    checksum: l1.wasm_checksum.clone(),
                    latest_height: None,
                }
                .encode_to_vec(),
            }),
            proof: vec![],
            proof_height: None,
        },
        &l1,
    )
    .unwrap();
}

#[test]
fn rejects_client_state_with_wrong_checksum_or_type() {
    let l1 = pinned_l1();
    let wrong_checksum = validate_client_state(
        QueryClientStateResponse {
            client_state: Some(Any {
                type_url: "/ibc.lightclients.wasm.v1.ClientState".into(),
                value: WasmClientState {
                    data: vec![],
                    checksum: vec![8; 32],
                    latest_height: None,
                }
                .encode_to_vec(),
            }),
            proof: vec![],
            proof_height: None,
        },
        &l1,
    )
    .unwrap_err();
    assert!(matches!(wrong_checksum, Error::L1ChecksumMismatch));

    let wrong_type = validate_client_state(
        QueryClientStateResponse {
            client_state: Some(Any {
                type_url: "/unexpected.ClientState".into(),
                value: vec![],
            }),
            proof: vec![],
            proof_height: None,
        },
        &l1,
    )
    .unwrap_err();
    assert!(matches!(wrong_type, Error::UnexpectedHostTypeUrl(_)));
}

#[test]
fn rejects_inactive_pinned_client() {
    let error = validate_client_status(QueryClientStatusResponse {
        status: "Frozen".into(),
    })
    .unwrap_err();
    assert!(matches!(error, Error::L1ClientInactive(status) if status == "Frozen"));
}

#[test]
fn rejects_missing_beacon_consensus_state() {
    let error = decode_consensus_state(
        QueryConsensusStateResponse {
            consensus_state: None,
            proof: vec![],
            proof_height: None,
        },
        42,
    )
    .unwrap_err();
    assert!(matches!(error, Error::HostQuery(_)));
}

#[test]
fn accepts_matching_beacon_consensus_state() {
    let expected_slot = 42;
    let state = EthereumConsensusState {
        slot: expected_slot,
        state_root: B256::repeat_byte(1),
        timestamp: 123,
        current_sync_committee: SummarizedSyncCommittee::default(),
        next_sync_committee: None,
    };
    let result = decode_consensus_state(
        consensus_response(serde_json::to_vec(&state).unwrap()),
        expected_slot,
    )
    .unwrap();
    assert_eq!(result.consensus, state);
}

#[test]
fn rejects_malformed_or_wrong_slot_beacon_consensus_state() {
    let malformed =
        decode_consensus_state(consensus_response(b"not-json".to_vec()), 42).unwrap_err();
    assert!(matches!(malformed, Error::JsonDecode(_)));

    let state = EthereumConsensusState {
        slot: 43,
        state_root: B256::repeat_byte(1),
        timestamp: 123,
        current_sync_committee: SummarizedSyncCommittee::default(),
        next_sync_committee: None,
    };
    let mismatch =
        decode_consensus_state(consensus_response(serde_json::to_vec(&state).unwrap()), 42)
            .unwrap_err();
    assert!(matches!(
        mismatch,
        Error::L1ConsensusSlotMismatch {
            expected: 42,
            actual: 43
        }
    ));
}

fn pinned_l1() -> L1Client {
    L1Client {
        client_id: "08-wasm-0".into(),
        wasm_checksum: vec![7; 32],
    }
}

fn consensus_response(data: Vec<u8>) -> QueryConsensusStateResponse {
    QueryConsensusStateResponse {
        consensus_state: Some(Any {
            type_url: "/ibc.lightclients.wasm.v1.ConsensusState".into(),
            value: WasmConsensusState { data }.encode_to_vec(),
        }),
        proof: vec![],
        proof_height: None,
    }
}
