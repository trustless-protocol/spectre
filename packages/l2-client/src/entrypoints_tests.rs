use alloy_primitives::{Address, Bloom, Bytes, FixedBytes, B256, U256};
use cosmwasm_std::{
    testing::{mock_dependencies, mock_env},
    Binary,
};
use serde::{Deserialize, Serialize};

use crate::{
    canonical_header::{CanonicalEvmHeader, ExecutionHeaderFork},
    entrypoints,
    error::Error,
    msg::{AttestedL2Header, ClientMessage, EvmAccountProof, QueryMsg, SudoMsg},
    runtime,
    state::{ClientState, CommonProfile, ConsensusState, Header, Height, RuntimeProfile},
    L2LightClient,
};

#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize)]
struct Profile(CommonProfile);

impl RuntimeProfile for Profile {
    fn common(&self) -> &CommonProfile {
        &self.0
    }

    fn expected_profile_version() -> &'static str {
        "entrypoint_test_v1"
    }
}

struct Adapter;

impl L2LightClient for Adapter {
    type Profile = Profile;

    fn verify(_profile: &Profile, header: &AttestedL2Header) -> Result<Header, Error> {
        Ok(normalize(header))
    }
}

fn profile() -> Profile {
    Profile(CommonProfile {
        l2_chain_id: 2,
        l2_router: Address::with_last_byte(1),
        commitment_slot: B256::with_last_byte(2),
        profile_version: "entrypoint_test_v1".into(),
        l2_header_fork: ExecutionHeaderFork::London,
    })
}

fn attested(block: u8) -> AttestedL2Header {
    AttestedL2Header {
        l2_header: CanonicalEvmHeader {
            parent_hash: B256::with_last_byte(4),
            ommers_hash: B256::with_last_byte(2),
            beneficiary: Address::with_last_byte(3),
            state_root: B256::with_last_byte(block),
            transactions_root: B256::with_last_byte(block + 1),
            receipts_root: B256::with_last_byte(6),
            logs_bloom: Bloom::ZERO,
            difficulty: U256::ZERO,
            number: 5,
            gas_limit: 30_000_000,
            gas_used: 21_000,
            timestamp: 1_700_000_000,
            extra_data: Bytes::new(),
            mix_hash: B256::with_last_byte(block + 2),
            nonce: FixedBytes::ZERO,
            base_fee_per_gas: Some(U256::from(9)),
            withdrawals_root: None,
            blob_gas_used: None,
            excess_blob_gas: None,
            parent_beacon_block_root: None,
            requests_hash: None,
        },
        router_proof: EvmAccountProof { proof: vec![] },
    }
}

fn normalize(header: &AttestedL2Header) -> Header {
    Header {
        height: Height::new(header.l2_header.number).unwrap(),
        state_root: header.l2_header.state_root,
        router_storage_root: header.l2_header.transactions_root,
        timestamp_seconds: header.l2_header.timestamp,
        l2_block_hash: header.l2_header.mix_hash,
        parent_hash: header.l2_header.parent_hash,
    }
}

#[test]
fn verify_client_message_rejects_non_conflicting_misbehaviour_envelopes() {
    let mut deps = mock_dependencies();
    let trusted = normalize(&attested(5));
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        profile: profile(),
    };
    let consensus: ConsensusState = trusted.consensus_state(0).unwrap();
    runtime::instantiate(deps.as_mut().storage, &client, &consensus, vec![9], 10).unwrap();

    let identical = attested(5);
    let mut different_height = attested(8);
    different_height.l2_header.number = 6;

    for message in [
        ClientMessage::Misbehaviour {
            header_1: identical.clone(),
            header_2: identical,
        },
        ClientMessage::Misbehaviour {
            header_1: attested(5),
            header_2: different_height,
        },
    ] {
        let client_message = Binary::from(serde_json::to_vec(&message).unwrap());

        assert!(matches!(
            entrypoints::query::<Adapter>(
                deps.as_ref(),
                QueryMsg::VerifyClientMessage { client_message }
            ),
            Err(Error::InvalidHeader("headers do not prove misbehaviour"))
        ));
    }
}

/// Every host path must refuse a conflict without authenticated attestations.
#[test]
fn host_paths_reject_conflicts_without_authenticated_attestations() {
    let mut deps = mock_dependencies();
    let trusted = normalize(&attested(5));
    let client = ClientState {
        latest_height: 5,
        frozen_height: None,
        profile: profile(),
    };
    let consensus: ConsensusState = trusted.consensus_state(0).unwrap();
    runtime::instantiate(deps.as_mut().storage, &client, &consensus, vec![9], 10).unwrap();

    let client_message = Binary::from(
        serde_json::to_vec(&ClientMessage::Misbehaviour {
            header_1: attested(5),
            header_2: attested(9),
        })
        .unwrap(),
    );

    for error in [
        entrypoints::query::<Adapter>(
            deps.as_ref(),
            QueryMsg::VerifyClientMessage {
                client_message: client_message.clone(),
            },
        )
        .unwrap_err(),
        entrypoints::query::<Adapter>(
            deps.as_ref(),
            QueryMsg::CheckForMisbehaviour {
                client_message: client_message.clone(),
            },
        )
        .unwrap_err(),
        entrypoints::sudo::<Adapter>(
            deps.as_mut(),
            &mock_env(),
            SudoMsg::UpdateStateOnMisbehaviour { client_message },
        )
        .unwrap_err(),
    ] {
        assert!(
            matches!(
                error,
                Error::InvalidHeader("conflicting headers do not carry authenticated attestations")
            ),
            "expected a conflict without authenticated attestations to be refused, got {error:?}"
        );
    }
    assert_eq!(
        runtime::client_state::<Profile>(deps.as_ref().storage)
            .unwrap()
            .frozen_height,
        None,
        "a conflict without authenticated attestations must not freeze the client"
    );
}
