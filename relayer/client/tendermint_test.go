package client

import (
	"encoding/json"
	"math/big"
	"testing"

	misbehaviourContract "relayer/bindings/Misbehaviour"
	updateClientContract "relayer/bindings/UpdateClient"

	ics23 "github.com/cosmos/ics23/go"
)

func TestParseTrustThreshold(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantNum uint8
		wantDen uint8
		wantErr bool
	}{
		{name: "valid 2/3", input: "2/3", wantNum: 2, wantDen: 3},
		{name: "valid 1/3", input: "1/3", wantNum: 1, wantDen: 3},
		{name: "invalid format dash", input: "2-3", wantErr: true},
		{name: "invalid format no separator", input: "abc", wantErr: true},
		{name: "invalid numerator", input: "abc/3", wantErr: true},
		{name: "invalid denominator", input: "2/abc", wantErr: true},
		{name: "zero denominator", input: "2/0", wantErr: true},
		{name: "numerator exceeds uint8", input: "300/6", wantErr: true},
		{name: "denominator exceeds uint8", input: "2/256", wantErr: true},
		{name: "uint8 max accepted", input: "255/255", wantNum: 255, wantDen: 255},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseTrustThreshold(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q, got nil", tc.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", tc.input, err)
			}
			if got.Numerator != tc.wantNum {
				t.Errorf("numerator: got %d, want %d", got.Numerator, tc.wantNum)
			}
			if got.Denominator != tc.wantDen {
				t.Errorf("denominator: got %d, want %d", got.Denominator, tc.wantDen)
			}
		})
	}
}

func TestEncodeMisbehaviourContractMsg(t *testing.T) {
	proof := misbehaviourContract.ISpectreClientMsgsBatchProof{}
	for i := range proof.Proof {
		proof.Proof[i] = big.NewInt(0)
	}
	for i := range proof.Commitments {
		proof.Commitments[i] = big.NewInt(0)
		proof.CommitmentPok[i] = big.NewInt(0)
	}
	header := misbehaviourContract.IICS07TendermintMsgsHeader{}
	header.SignedHeader.Header.Time = big.NewInt(0)
	msg := misbehaviourContract.ISpectreClientMsgsMsgSubmitMisbehaviour{
		Misbehaviour: misbehaviourContract.ISpectreClientMsgsMisbehaviour{
			Header1: header,
			Header2: header,
		},
		TrustedConsensusState1: misbehaviourContract.IICS07TendermintMsgsConsensusState{Timestamp: big.NewInt(1)},
		TrustedConsensusState2: misbehaviourContract.IICS07TendermintMsgsConsensusState{Timestamp: big.NewInt(2)},
		Time:                   big.NewInt(3),
		Proof1:                 proof,
		Proof2:                 proof,
	}
	encoded, err := EncodeMisbehaviourContractMsg(msg)
	if err != nil {
		t.Fatalf("encode generated misbehaviour contract message: %v", err)
	}
	if len(encoded) == 0 {
		t.Fatal("expected non-empty ABI payload")
	}
}

func TestParseLeafOp(t *testing.T) {
	t.Run("nil input returns zero value", func(t *testing.T) {
		result := ParseLeafOp(nil)
		if result.HashOp != 0 {
			t.Errorf("HashOp: got %d, want 0", result.HashOp)
		}
		if result.PrehashKey != 0 {
			t.Errorf("PrehashKey: got %d, want 0", result.PrehashKey)
		}
		if result.PrehashValue != 0 {
			t.Errorf("PrehashValue: got %d, want 0", result.PrehashValue)
		}
		if result.Prefix != nil {
			t.Errorf("Prefix: got %v, want nil", result.Prefix)
		}
	})

	t.Run("valid leaf op", func(t *testing.T) {
		leafOp := &ics23.LeafOp{
			Hash:         ics23.HashOp_SHA256,
			PrehashKey:   ics23.HashOp_NO_HASH,
			PrehashValue: ics23.HashOp_SHA256,
			Prefix:       []byte{0x00},
		}
		result := ParseLeafOp(leafOp)
		if result.HashOp != uint8(ics23.HashOp_SHA256) {
			t.Errorf("HashOp: got %d, want %d", result.HashOp, uint8(ics23.HashOp_SHA256))
		}
		if result.PrehashKey != uint8(ics23.HashOp_NO_HASH) {
			t.Errorf("PrehashKey: got %d, want %d", result.PrehashKey, uint8(ics23.HashOp_NO_HASH))
		}
		if result.PrehashValue != uint8(ics23.HashOp_SHA256) {
			t.Errorf("PrehashValue: got %d, want %d", result.PrehashValue, uint8(ics23.HashOp_SHA256))
		}
		if len(result.Prefix) != 1 || result.Prefix[0] != 0x00 {
			t.Errorf("Prefix: got %v, want [0x00]", result.Prefix)
		}
	})
}

func TestParseInnerOp(t *testing.T) {
	t.Run("nil input returns zero value", func(t *testing.T) {
		result := ParseInnerOp(nil)
		if result.HashOp != 0 {
			t.Errorf("HashOp: got %d, want 0", result.HashOp)
		}
		if result.Prefix != nil {
			t.Errorf("Prefix: got %v, want nil", result.Prefix)
		}
		if result.Suffix != nil {
			t.Errorf("Suffix: got %v, want nil", result.Suffix)
		}
	})

	t.Run("valid inner op", func(t *testing.T) {
		innerOp := &ics23.InnerOp{
			Hash:   ics23.HashOp_SHA256,
			Prefix: []byte{0x01},
			Suffix: []byte{0x02},
		}
		result := ParseInnerOp(innerOp)
		if result.HashOp != uint8(ics23.HashOp_SHA256) {
			t.Errorf("HashOp: got %d, want %d", result.HashOp, uint8(ics23.HashOp_SHA256))
		}
		if len(result.Prefix) != 1 || result.Prefix[0] != 0x01 {
			t.Errorf("Prefix: got %v, want [0x01]", result.Prefix)
		}
		if len(result.Suffix) != 1 || result.Suffix[0] != 0x02 {
			t.Errorf("Suffix: got %v, want [0x02]", result.Suffix)
		}
	})
}

func TestEncodeClientState(t *testing.T) {
	t.Run("valid client state encodes without error", func(t *testing.T) {
		clientState := ClientState{
			ChainId: "cosmoshub-4",
			TrustLevel: TrustThreshold{
				Numerator:   2,
				Denominator: 3,
			},
			LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
				RevisionNumber: 4,
				RevisionHeight: 100,
			},
			TrustingPeriod:  1209600,
			UnbondingPeriod: 1814400,
			IsFrozen:        false,
			ClockDrift:      15,
		}
		encoded, err := EncodeClientState(clientState)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(encoded) == 0 {
			t.Fatal("encoded result should not be empty")
		}
	})
}

func TestDecodeClientState(t *testing.T) {
	original := ClientState{
		ChainId: "cosmoshub-4",
		TrustLevel: TrustThreshold{
			Numerator:   2,
			Denominator: 3,
		},
		LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
			RevisionNumber: 4,
			RevisionHeight: 100,
		},
		TrustingPeriod:  1209600,
		UnbondingPeriod: 1814400,
		IsFrozen:        false,
		ClockDrift:      15,
	}

	encoded, err := EncodeClientState(original)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}

	decoded, err := DecodeClientState(encoded)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}

	if decoded.ChainId != original.ChainId {
		t.Errorf("ChainId: got %q, want %q", decoded.ChainId, original.ChainId)
	}
	if decoded.TrustLevel.Numerator != original.TrustLevel.Numerator {
		t.Errorf("TrustLevel.Numerator: got %d, want %d", decoded.TrustLevel.Numerator, original.TrustLevel.Numerator)
	}
	if decoded.TrustLevel.Denominator != original.TrustLevel.Denominator {
		t.Errorf("TrustLevel.Denominator: got %d, want %d", decoded.TrustLevel.Denominator, original.TrustLevel.Denominator)
	}
	if decoded.LatestHeight.RevisionNumber != original.LatestHeight.RevisionNumber {
		t.Errorf("LatestHeight.RevisionNumber: got %d, want %d", decoded.LatestHeight.RevisionNumber, original.LatestHeight.RevisionNumber)
	}
	if decoded.LatestHeight.RevisionHeight != original.LatestHeight.RevisionHeight {
		t.Errorf("LatestHeight.RevisionHeight: got %d, want %d", decoded.LatestHeight.RevisionHeight, original.LatestHeight.RevisionHeight)
	}
	if decoded.TrustingPeriod != original.TrustingPeriod {
		t.Errorf("TrustingPeriod: got %d, want %d", decoded.TrustingPeriod, original.TrustingPeriod)
	}
	if decoded.UnbondingPeriod != original.UnbondingPeriod {
		t.Errorf("UnbondingPeriod: got %d, want %d", decoded.UnbondingPeriod, original.UnbondingPeriod)
	}
	if decoded.IsFrozen != original.IsFrozen {
		t.Errorf("IsFrozen: got %v, want %v", decoded.IsFrozen, original.IsFrozen)
	}
	if decoded.ClockDrift != original.ClockDrift {
		t.Errorf("ClockDrift: got %d, want %d", decoded.ClockDrift, original.ClockDrift)
	}
}

func TestEncodeConsensusState(t *testing.T) {
	t.Run("valid consensus state encodes without error", func(t *testing.T) {
		consensusState := updateClientContract.IICS07TendermintMsgsConsensusState{
			Timestamp:          big.NewInt(1700000000000),
			Root:               [32]byte{0x01, 0x02, 0x03},
			NextValidatorsHash: [32]byte{0x04, 0x05, 0x06},
		}
		encoded, err := EncodeConsensusState(consensusState)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(encoded) == 0 {
			t.Fatal("encoded result should not be empty")
		}
	})
}

// TestEncodeUpdateApplicationStateMsgMatchesGeneratedABI checks that the manual
// ABI tuple (updateApplicationStateMsgType) produces bytes that the abigen-
// generated UpdateClient.verifyHeader input tuple can round-trip — i.e. the
// hand-built encoding matches the generated ABI exactly.
func TestEncodeUpdateApplicationStateMsgMatchesGeneratedABI(t *testing.T) {
	var validatorsHash [32]byte
	validatorsHash[0] = 0x88
	validatorsHash[1] = 0xbe

	zero8 := [8]*big.Int{}
	for i := range zero8 {
		zero8[i] = big.NewInt(0)
	}
	zero2 := [2]*big.Int{big.NewInt(0), big.NewInt(0)}

	msg := updateClientContract.ISpectreClientMsgsMsgUpdateApplicationState{
		TrustedConsensusState: updateClientContract.IICS07TendermintMsgsConsensusState{
			Timestamp:          big.NewInt(1700000000000000000),
			NextValidatorsHash: validatorsHash,
		},
		ProposedHeader: updateClientContract.IICS07TendermintMsgsHeader{
			SignedHeader: updateClientContract.IICS07TendermintMsgsSignedHeader{
				Header: updateClientContract.IICS07TendermintMsgsBlockHeader{
					ChainId:            "test-0",
					Time:               big.NewInt(1700000001000000000),
					ValidatorsHash:     validatorsHash,
					NextValidatorsHash: validatorsHash,
				},
				Commit: updateClientContract.IICS07TendermintMsgsBlockCommit{
					Height: 58,
				},
			},
			TrustedHeight: updateClientContract.IICS02ClientMsgsHeight{
				RevisionNumber: 0,
				RevisionHeight: 19,
			},
		},
		Time: big.NewInt(1700000002000000000),
		Proof: updateClientContract.ISpectreClientMsgsBatchProof{
			Proof:                  zero8,
			Commitments:            zero2,
			CommitmentPok:          zero2,
			Bucket:                 16,
			SignerIndices:          []uint32{0, 1},
			PinnedValidatorIndices: []uint32{7, 8},
			SignerPubkeys:          [][32]byte{{0x01}, {0x02}},
			Active:                 []bool{true, true},
		},
	}

	encoded, err := EncodeUpdateApplicationStateMsg(msg)
	if err != nil {
		t.Fatalf("encode update application state msg: %v", err)
	}

	contractABI, err := updateClientContract.ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		t.Fatalf("parse generated ABI: %v", err)
	}
	unpacked, err := contractABI.Methods["verifyHeader"].Inputs.Unpack(encoded)
	if err != nil {
		t.Fatalf("generated ABI failed to unpack encoded msg: %v", err)
	}
	if len(unpacked) != 1 {
		t.Fatalf("unpacked %d values, want 1", len(unpacked))
	}

	jsonBytes, err := json.Marshal(unpacked[0])
	if err != nil {
		t.Fatalf("marshal unpacked tuple: %v", err)
	}
	var decoded updateClientContract.ISpectreClientMsgsMsgUpdateApplicationState
	if err := json.Unmarshal(jsonBytes, &decoded); err != nil {
		t.Fatalf("unmarshal decoded update msg: %v", err)
	}
	if decoded.ProposedHeader.SignedHeader.Header.ValidatorsHash != validatorsHash {
		t.Fatalf(
			"validatorsHash decoded incorrectly: got %x want %x",
			decoded.ProposedHeader.SignedHeader.Header.ValidatorsHash,
			validatorsHash,
		)
	}
	if len(decoded.Proof.PinnedValidatorIndices) != 2 || decoded.Proof.PinnedValidatorIndices[0] != 7 || decoded.Proof.PinnedValidatorIndices[1] != 8 {
		t.Fatalf("pinned validator indices decoded incorrectly: %+v", decoded.Proof.PinnedValidatorIndices)
	}
}

func TestBytesToBytes32(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  [32]byte
	}{
		{name: "nil returns zero array", input: nil, want: [32]byte{}},
		{name: "empty returns zero array", input: []byte{}, want: [32]byte{}},
		{
			name:  "short input pads with zeros",
			input: []byte{0xAA, 0xBB},
			want: func() [32]byte {
				var b [32]byte
				b[0], b[1] = 0xAA, 0xBB
				return b
			}(),
		},
		{
			name: "exact 32 bytes copies fully",
			input: func() []byte {
				b := make([]byte, 32)
				for i := range b {
					b[i] = byte(i)
				}
				return b
			}(),
			want: func() [32]byte {
				var b [32]byte
				for i := range b {
					b[i] = byte(i)
				}
				return b
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := bytesToBytes32(tc.input)
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
