use cosmwasm_std::Binary;

use serde::{Deserialize, Serialize};

use crate::msg::{
    CheckForMisbehaviourResult, ClientMessage, EvmAccountProof, IbcHeight,
    IndexedAttestorSignature, MigrateMsg, SignedAttestedL2Header, StatusResult, SudoMsg,
    TimestampAtHeightResult, UpdateStateResult,
};
use crate::{canonical_header::CanonicalEvmHeader, error::Error};

#[derive(Clone, Debug, PartialEq, Eq, Deserialize, Serialize)]
#[serde(deny_unknown_fields)]
struct UnsignedHeader {
    l2_header: CanonicalEvmHeader,
    router_proof: EvmAccountProof,
}

#[test]
fn decodes_the_shared_go_to_rust_client_message_fixture() {
    let fixture = include_bytes!("../../../test/fixtures/wasm-contracts/l2-client-message.json");
    let message: ClientMessage<UnsignedHeader> = serde_json::from_slice(fixture).unwrap();
    let ClientMessage::Header(header) = message else {
        panic!("shared fixture must contain a header");
    };
    assert_eq!(header.l2_header.number, 42);
    assert_eq!(header.router_proof.proof, vec![vec![1, 2], vec![3]]);
}

#[test]
fn signed_header_uses_base64_signatures_and_numeric_proof_nodes() {
    let fixture = include_bytes!("../../../test/fixtures/wasm-contracts/l2-client-message.json");
    let unsigned: ClientMessage<UnsignedHeader> = serde_json::from_slice(fixture).unwrap();
    let ClientMessage::Header(unsigned) = unsigned else {
        panic!("shared fixture must contain a header");
    };
    let signed = ClientMessage::Header(SignedAttestedL2Header {
        l2_header: unsigned.l2_header,
        router_proof: unsigned.router_proof,
        attestor_signature: vec![IndexedAttestorSignature {
            attestor_index: 0,
            signature: Binary::from(vec![7; 64]),
        }],
    });

    let json = serde_json::to_string(&signed).unwrap();
    assert!(json.contains(r#""attestor_signature":["#));
    assert!(!json.contains(r#""signatures""#));
    assert!(json.contains(r#""signature":"BwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBwcHBw==""#));
    assert!(json.contains(r#""proof":[[1,2],[3]]"#));
    let decoded: ClientMessage<SignedAttestedL2Header> = serde_json::from_str(&json).unwrap();
    assert_eq!(decoded, signed);
}

#[test]
fn unsigned_header_does_not_decode_as_signed_header() {
    let fixture = include_bytes!("../../../test/fixtures/wasm-contracts/l2-client-message.json");
    assert!(serde_json::from_slice::<ClientMessage<SignedAttestedL2Header>>(fixture).is_err());
}

#[test]
fn migrate_wire_is_strict_snake_case() {
    assert_eq!(
        serde_json::to_string(&MigrateMsg::KeepAttestors {}).unwrap(),
        r#"{"keep_attestors":{}}"#
    );
    assert_eq!(
        serde_json::to_string(&MigrateMsg::ReplaceAttestors {
            public_keys: vec![Binary::from(vec![1; 32])],
            threshold: 1,
        })
        .unwrap(),
        r#"{"replace_attestors":{"public_keys":["AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE="],"threshold":1}}"#
    );
    assert!(serde_json::from_slice::<MigrateMsg>(
        br#"{"replace_attestors":{"public_keys":["AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE="],"threshold":1,"unexpected":true}}"#,
    )
    .is_err());
}

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
