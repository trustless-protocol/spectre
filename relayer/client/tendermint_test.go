package client

import (
	"encoding/json"
	"math/big"
	"testing"

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
		clientState := updateClientContract.IICS07TendermintMsgsClientState{
			ChainId: "cosmoshub-4",
			TrustLevel: updateClientContract.IICS07TendermintMsgsTrustThreshold{
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
			ZkAlgorithm:     uint8(Groth16),
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
	original := updateClientContract.IICS07TendermintMsgsClientState{
		ChainId: "cosmoshub-4",
		TrustLevel: updateClientContract.IICS07TendermintMsgsTrustThreshold{
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
		ZkAlgorithm:     uint8(Groth16),
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
	if decoded.ZkAlgorithm != original.ZkAlgorithm {
		t.Errorf("ZkAlgorithm: got %d, want %d", decoded.ZkAlgorithm, original.ZkAlgorithm)
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

func TestEncodeUpdateClientMsgMatchesGeneratedABI(t *testing.T) {
	var validatorsHash [32]byte
	validatorsHash[0] = 0x88
	validatorsHash[1] = 0xbe

	var deltaBaseHash [32]byte
	deltaBaseHash[0] = 0xaa

	zero8 := [8]*big.Int{}
	for i := range zero8 {
		zero8[i] = big.NewInt(0)
	}
	zero2 := [2]*big.Int{big.NewInt(0), big.NewInt(0)}
	var deltaIndices [16]uint32
	var deltaPubKeys [16][32]byte
	var deltaVotingPowers [16]uint64
	deltaIndices[0] = 7
	deltaPubKeys[0] = [32]byte{0x07}
	deltaVotingPowers[0] = 110

	msg := updateClientContract.IUpdateClientMsgsMsgUpdateClient{
		ClientState: updateClientContract.IICS07TendermintMsgsClientState{
			ChainId: "test-0",
			TrustLevel: updateClientContract.IICS07TendermintMsgsTrustThreshold{
				Numerator:   1,
				Denominator: 3,
			},
			LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
				RevisionNumber: 0,
				RevisionHeight: 19,
			},
			TrustingPeriod:  1209600,
			UnbondingPeriod: 1814400,
			ZkAlgorithm:     uint8(Groth16),
			ClockDrift:      15,
		},
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
		Time:             big.NewInt(1700000002000000000),
		Proof:            zero8,
		Commitments:      zero2,
		CommitmentPok:    zero2,
		Bucket:           16,
		SignerIndices:    []uint32{0, 1},
		SignerPubkeys:    [][32]byte{{0x01}, {0x02}},
		TimestampSeconds: []uint64{1700000001, 1700000001},
		TimestampNanos:   []uint32{0, 1},
		Active:           []bool{true, true},
		CurrentValidatorSetDelta: updateClientContract.IUpdateClientMsgsValidatorSetDelta{
			BaseValidatorsHash: deltaBaseHash,
			LeafCount:          1,
			Indices:            deltaIndices,
			PubKeys:            deltaPubKeys,
			VotingPowers:       deltaVotingPowers,
		},
	}

	encoded, err := EncodeUpdateClientMsg(msg)
	if err != nil {
		t.Fatalf("encode update client msg: %v", err)
	}

	contractABI, err := updateClientContract.ContractUpdateClientMetaData.GetAbi()
	if err != nil {
		t.Fatalf("parse generated ABI: %v", err)
	}
	unpacked, err := contractABI.Methods["updateClient"].Inputs.Unpack(encoded)
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
	var decoded updateClientContract.IUpdateClientMsgsMsgUpdateClient
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
	if decoded.CurrentValidatorSetDelta.BaseValidatorsHash != deltaBaseHash {
		t.Fatalf(
			"delta base hash decoded incorrectly: got %x want %x",
			decoded.CurrentValidatorSetDelta.BaseValidatorsHash,
			deltaBaseHash,
		)
	}
	if decoded.CurrentValidatorSetDelta.LeafCount != 1 ||
		decoded.CurrentValidatorSetDelta.Indices[0] != 7 ||
		decoded.CurrentValidatorSetDelta.PubKeys[0] != deltaPubKeys[0] ||
		decoded.CurrentValidatorSetDelta.VotingPowers[0] != 110 {
		t.Fatalf("delta decoded incorrectly: %+v", decoded.CurrentValidatorSetDelta)
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

func TestSupportedZkAlgorithmString(t *testing.T) {
	tests := []struct {
		name string
		alg  SupportedZkAlgorithm
		want string
	}{
		{name: "Groth16", alg: Groth16, want: "Groth16"},
		{name: "Plonk", alg: Plonk, want: "Plonk"},
		{name: "unknown value", alg: SupportedZkAlgorithm(99), want: "Unknown"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.alg.String()
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
