package witness

import (
	"fmt"
	"slices"
)

// SupportedBuckets is the protocol's bucket family: one circuit and one
// on-chain Groth16Verifier_N{N} exist per member. It is a protocol constant,
// not a deployment choice — adding a member means a new circuit + verifier
// contract, not editing this line. Which buckets are actually USABLE at
// runtime is a different question the engine answers by discovering artifacts
// on disk (gnark.New scans bin/n*/); nothing else declares the loaded set.
var SupportedBuckets = []int{4, 8, 16, 32, 64, 128}

// IsSupported reports whether n is a member of the protocol bucket family.
func IsSupported(n int) bool { return slices.Contains(SupportedBuckets, n) }

// SmallestBucketGEQ returns the smallest bucket in the given set that fits n
// signers. The set is whatever the caller actually has (the engine's loaded
// buckets), so a set with holes — say {4, 16} with 8 unbuilt — still resolves
// (6 signers → 16).
func SmallestBucketGEQ(n int, buckets []int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("signer count must be positive, got %d", n)
	}
	best := 0
	for _, b := range buckets {
		if b >= n && (best == 0 || b < best) {
			best = b
		}
	}
	if best == 0 {
		return 0, fmt.Errorf("%d signers exceed the largest available bucket (available: %v)", n, buckets)
	}
	return best, nil
}
