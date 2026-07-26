package opstack

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

// SyncStatus carries the replica heads the attestor gates on. Unsafe is
// sequencer gossip; Safe is derived from L1 batch data (reorgable with L1);
// Finalized is derived from finalized L1 data — "tx finality", irreversible.
type SyncStatus struct {
	UnsafeL2    uint64
	SafeL2      uint64
	FinalizedL2 uint64
}

// ReplicaClient talks to the op-node of the verify-mode replica.
type ReplicaClient struct {
	c *rpc.Client
}

func DialReplica(ctx context.Context, url string) (*ReplicaClient, error) {
	c, err := rpc.DialContext(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to dial op-node rpc %s: %w", url, err)
	}
	return &ReplicaClient{c: c}, nil
}

func (r *ReplicaClient) Close() {
	r.c.Close()
}

// l2BlockRef is the subset of op-node's L2BlockRef the attestor reads.
type l2BlockRef struct {
	Number uint64 `json:"number"`
}

type syncStatusResponse struct {
	UnsafeL2    l2BlockRef `json:"unsafe_l2"`
	SafeL2      l2BlockRef `json:"safe_l2"`
	FinalizedL2 l2BlockRef `json:"finalized_l2"`
}

// SyncStatus returns the replica's current unsafe, safe and finalized L2
// heads via optimism_syncStatus.
func (r *ReplicaClient) SyncStatus(ctx context.Context) (SyncStatus, error) {
	var resp syncStatusResponse
	if err := r.c.CallContext(ctx, &resp, "optimism_syncStatus"); err != nil {
		return SyncStatus{}, fmt.Errorf("optimism_syncStatus failed: %w", err)
	}
	return SyncStatus{
		UnsafeL2:    resp.UnsafeL2.Number,
		SafeL2:      resp.SafeL2.Number,
		FinalizedL2: resp.FinalizedL2.Number,
	}, nil
}

type outputAtBlockResponse struct {
	OutputRoot hexutil.Bytes `json:"outputRoot"`
}

// OutputAtBlock returns the replica's self-computed output root at the given
// L2 block via optimism_outputAtBlock.
func (r *ReplicaClient) OutputAtBlock(ctx context.Context, l2Block uint64) ([32]byte, error) {
	var resp outputAtBlockResponse
	if err := r.c.CallContext(ctx, &resp, "optimism_outputAtBlock", hexutil.Uint64(l2Block)); err != nil {
		return [32]byte{}, fmt.Errorf("optimism_outputAtBlock(%d) failed: %w", l2Block, err)
	}
	if len(resp.OutputRoot) != 32 {
		return [32]byte{}, fmt.Errorf("optimism_outputAtBlock(%d) returned %d-byte output root, want 32", l2Block, len(resp.OutputRoot))
	}
	var root [32]byte
	copy(root[:], resp.OutputRoot)
	return root, nil
}
