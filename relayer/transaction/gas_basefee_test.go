package transaction

import (
	"math/big"
	"testing"
)

// TestLiftAboveBaseFee covers the pricing failure seen on Arbitrum Sepolia:
//
//	max fee per gas less than block base fee: maxFeePerGas: 21486000, baseFee: 22414000
//
// SuggestGasPrice reflects the chain when it is called; on a chain whose base fee
// is climbing the value is already stale by the time EstimateGas uses it, and the
// node rejects the call rather than queueing it. Local devnets never show this —
// their base fee sits at the floor — so it only appears against a public network.
func TestLiftAboveBaseFee(t *testing.T) {
	cases := []struct {
		name      string
		suggested *big.Int
		baseFee   *big.Int
		want      *big.Int
	}{
		{
			// The exact numbers from the failure.
			name:      "suggestion below base fee is lifted clear",
			suggested: big.NewInt(21_486_000), baseFee: big.NewInt(22_414_000),
			want: big.NewInt(44_828_000), // 2x base fee
		},
		{
			name:      "suggestion already clear is untouched",
			suggested: big.NewInt(100_000_000), baseFee: big.NewInt(22_414_000),
			want: big.NewInt(100_000_000),
		},
		{
			name:      "exactly at the floor is untouched",
			suggested: big.NewInt(44_828_000), baseFee: big.NewInt(22_414_000),
			want: big.NewInt(44_828_000),
		},
		{
			name:      "nil suggestion takes the floor",
			suggested: nil, baseFee: big.NewInt(1_000),
			want: big.NewInt(2_000),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := liftAboveBaseFee(tc.suggested, tc.baseFee)
			if got.Cmp(tc.want) != 0 {
				t.Fatalf("liftAboveBaseFee(%v, %v) = %v, want %v", tc.suggested, tc.baseFee, got, tc.want)
			}
		})
	}
}

// A chain with no base fee (pre-1559) must pass the suggestion through untouched
// rather than being handed a nil-derived value.
func TestLiftAboveBaseFeeWithoutBaseFee(t *testing.T) {
	suggested := big.NewInt(5_000)
	if got := liftAboveBaseFee(suggested, nil); got != suggested {
		t.Fatalf("a nil base fee must return the suggestion unchanged, got %v", got)
	}
}
