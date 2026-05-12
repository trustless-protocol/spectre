package prover

import "fmt"

// Buckets lists the supported validator-count buckets. Each bucket has its own
// compiled circuit and (pk, vk) pair on disk under bin/n{N}/. The prover picks
// the smallest bucket that fits the required signer count for a block.
var Buckets = []int{4}

// SmallestBucketGEQ returns the smallest bucket size ≥ n. It errors when n
// exceeds the largest bucket — the operator must add a larger bucket and
// recompile.
func SmallestBucketGEQ(n int) (int, error) {
	return SmallestBucketGEQFrom(Buckets, n)
}

// SmallestBucketGEQFrom returns the smallest bucket size ≥ n from the provided
// ordered bucket list.
func SmallestBucketGEQFrom(buckets []int, n int) (int, error) {
	if len(buckets) == 0 {
		return 0, fmt.Errorf("no buckets configured")
	}
	for _, b := range buckets {
		if b >= n {
			return b, nil
		}
	}
	return 0, fmt.Errorf("signer count %d exceeds largest bucket %d", n, buckets[len(buckets)-1])
}

// MaxBucket returns the largest configured bucket.
func MaxBucket() int {
	return Buckets[len(Buckets)-1]
}
