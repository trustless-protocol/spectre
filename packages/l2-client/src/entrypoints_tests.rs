use alloy_primitives::{Address, Bloom, Bytes, FixedBytes, B256, U256};
use cosmwasm_std::{
    testing::{mock_dependencies, mock_env},
    Binary, Storage,
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

fn storage_snapshot(storage: &dyn Storage, heights: &[u64]) -> Vec<Option<Vec<u8>>> {
    std::iter::once(storage.get(crate::store::HOST_CLIENT_STATE_KEY.as_bytes()))
        .chain(
            heights
                .iter()
                .map(|height| storage.get(crate::store::consensus_db_key(*height).as_bytes())),
        )
        .collect()
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

#[test]
fn non_zero_delay_precedes_path_height_state_and_proof_validation() {
    let mut deps = mock_dependencies();
    let before = storage_snapshot(deps.as_ref().storage, &[0, 5]);

    for (message, expected_time, expected_blocks) in [
        (
            SudoMsg::VerifyMembership {
                height: crate::msg::IbcHeight {
                    revision_number: 9,
                    revision_height: 0,
                },
                delay_time_period: 7,
                delay_block_period: 0,
                proof: Binary::from(b"not-json"),
                merkle_path: crate::msg::MerklePath { key_path: vec![] },
                value: Binary::default(),
            },
            7,
            0,
        ),
        (
            SudoMsg::VerifyNonMembership {
                height: crate::msg::IbcHeight {
                    revision_number: 9,
                    revision_height: 0,
                },
                delay_time_period: 0,
                delay_block_period: 3,
                proof: Binary::from(b"not-json"),
                merkle_path: crate::msg::MerklePath { key_path: vec![] },
            },
            0,
            3,
        ),
    ] {
        assert!(matches!(
            entrypoints::sudo::<Adapter>(deps.as_mut(), &mock_env(), message),
            Err(Error::UnsupportedNonZeroDelay {
                delay_time_period,
                delay_block_period,
            }) if delay_time_period == expected_time && delay_block_period == expected_blocks
        ));
        assert_eq!(storage_snapshot(deps.as_ref().storage, &[0, 5]), before);
    }
}

#[test]
fn status_bytes_and_frozen_update_outcome_are_stable() {
    for (frozen_height, expected) in [
        (None, br#"{"status":"Active"}"#.as_slice()),
        (Some(5), br#"{"status":"Frozen"}"#.as_slice()),
    ] {
        let mut deps = mock_dependencies();
        let trusted = normalize(&attested(5));
        let client = ClientState {
            latest_height: 5,
            frozen_height,
            profile: profile(),
        };
        let consensus = trusted.consensus_state(0).unwrap();
        runtime::instantiate(deps.as_mut().storage, &client, &consensus, vec![9], 10).unwrap();

        let status = entrypoints::query::<Adapter>(deps.as_ref(), QueryMsg::Status {}).unwrap();
        assert_eq!(status.as_slice(), expected);

        if let Some(height) = frozen_height {
            let before = storage_snapshot(deps.as_ref().storage, &[5, 6]);
            let client_message =
                Binary::from(serde_json::to_vec(&ClientMessage::Header(attested(6))).unwrap());
            assert!(matches!(
                entrypoints::sudo::<Adapter>(
                    deps.as_mut(),
                    &mock_env(),
                    SudoMsg::UpdateState { client_message },
                ),
                Err(Error::Frozen(actual)) if actual == height
            ));
            assert_eq!(storage_snapshot(deps.as_ref().storage, &[5, 6]), before);
        }
    }
}
