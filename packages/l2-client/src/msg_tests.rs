use cosmwasm_std::Binary;

use crate::error::Error;
use crate::msg::{
    CheckForMisbehaviourResult, ClientMessage, IbcHeight, StatusResult, SudoMsg,
    TimestampAtHeightResult, UpdateStateResult,
};

#[test]
fn decodes_the_flat_host_membership_payload() {
    let message: SudoMsg = serde_json::from_slice(
        br#"{"verify_membership":{"height":{"revision_number":0,"revision_height":42},"delay_time_period":1,"delay_block_period":2,"proof":"AQ==","merkle_path":{"key_path":["Y2xpZW50LTA=" ]},"value":"Ag=="}}"#,
    )
    .unwrap();

    match message {
        SudoMsg::VerifyMembership {
            height,
            delay_time_period,
            delay_block_period,
            proof,
            merkle_path,
            value,
        } => {
            assert_eq!(height.revision_height, 42);
            assert_eq!(delay_time_period, 1);
            assert_eq!(delay_block_period, 2);
            assert_eq!(proof, Binary::from([1]));
            assert_eq!(merkle_path.key_path, vec![Binary::from(b"client-0")]);
            assert_eq!(value, Binary::from([2]));
        }
        _ => panic!("unexpected sudo message"),
    }
}

#[test]
fn serializes_host_result_wrappers() {
    assert_eq!(
        serde_json::to_string(&UpdateStateResult {
            heights: vec![IbcHeight {
                revision_number: 0,
                revision_height: 42,
            }],
        })
        .unwrap(),
        r#"{"heights":[{"revision_number":0,"revision_height":42}]}"#
    );
    assert_eq!(
        serde_json::to_string(&CheckForMisbehaviourResult {
            found_misbehaviour: true,
        })
        .unwrap(),
        r#"{"found_misbehaviour":true}"#
    );
    assert_eq!(
        serde_json::to_string(&TimestampAtHeightResult { timestamp: 7 }).unwrap(),
        r#"{"timestamp":7}"#
    );
    assert_eq!(
        serde_json::to_string(&StatusResult {
            status: "Active".into(),
        })
        .unwrap(),
        r#"{"status":"Active"}"#
    );
}

#[test]
fn decodes_misbehaviour_client_message_envelopes() {
    let message: ClientMessage<u64> =
        serde_json::from_slice(br#"{"type":"misbehaviour","value":{"header_1":7,"header_2":7}}"#)
            .unwrap();
    assert_eq!(
        message,
        ClientMessage::Misbehaviour {
            header_1: 7,
            header_2: 7,
        }
    );
}

#[test]
fn decodes_unsupported_lifecycle_messages_without_losing_payloads() {
    let upgrade: SudoMsg = serde_json::from_slice(
        br#"{"verify_upgrade_and_update_state":{"upgrade_client_state":"AQ==","upgrade_consensus_state":"Ag==","proof_upgrade_client":"Aw==","proof_upgrade_consensus_state":"BA=="}}"#,
    )
    .unwrap();
    assert!(matches!(
        upgrade,
        SudoMsg::VerifyUpgradeAndUpdateState {
            upgrade_client_state,
            upgrade_consensus_state,
            proof_upgrade_client,
            proof_upgrade_consensus_state,
        } if upgrade_client_state.as_slice() == [1]
            && upgrade_consensus_state.as_slice() == [2]
            && proof_upgrade_client.as_slice() == [3]
            && proof_upgrade_consensus_state.as_slice() == [4]
    ));

    let recovery: SudoMsg = serde_json::from_slice(br#"{"migrate_client_store":{}}"#).unwrap();
    assert!(matches!(recovery, SudoMsg::MigrateClientStore {}));
}

#[test]
fn unsupported_operation_error_text_is_stable() {
    assert_eq!(
        Error::UnsupportedNonZeroDelay {
            delay_time_period: 7,
            delay_block_period: 3,
        }
        .to_string(),
        "non-zero delay is unsupported: delay_time_period=7, delay_block_period=3"
    );
    assert_eq!(
        Error::UnsupportedLifecycleOperation {
            operation: "verify_upgrade_and_update_state",
        }
        .to_string(),
        "lifecycle operation is unsupported: verify_upgrade_and_update_state"
    );
}
