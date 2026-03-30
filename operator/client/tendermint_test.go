package client

import (
	"math/big"
	"testing"

	updateClientContract "operator/bindings/UpdateClient"

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
