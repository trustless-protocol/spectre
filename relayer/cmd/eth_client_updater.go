package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"

	"relayer/services"
)

const (
	// ethClientUpdateLag is how far behind the L1 head the pinned Ethereum client
	// may fall before a header build refreshes it. It is a slack window, not a target:
	// the client only ever advances to finalized headers, so demanding zero lag would
	// refresh on every build. Well inside a default execution node's state window.
	ethClientUpdateLag = 64

	// ethClientUpdateCooldown keeps a burst of header builds from each triggering
	// their own update: the client advances in finalized steps, so a second attempt
	// moments later would rebuild the same message for nothing.
	ethClientUpdateCooldown = 20 * time.Second

	// A submitted update can be rejected because its signature slot is ahead of the
	// slot the client computes from the *Cosmos* block timestamp, which trails real
	// time by a few seconds. That is transient and self-correcting: resubmitting the
	// same message once the host clock catches up succeeds, whereas dropping it and
	// fetching a newer update later just reproduces the skew.
	ethClientSkewRetries = 3
	ethClientSkewBackoff = 4 * time.Second

	// ethClientPollInterval is how often the anti-staleness loop CHECKS the pinned
	// client while no packets are flowing. Each tick costs one eth_blockNumber and one
	// client-state read and does nothing unless the client is already more than
	// ethClientUpdateLag behind, so this is a ceiling on staleness, not an update rate.
	//
	// The loop exists because the on-demand path only runs during a header build. With
	// no traffic nothing advances the client at all, and staleness is not merely a
	// latency problem: once the pinned block leaves the execution node's state window,
	// eth_getProof at that block fails, so the NEXT packet cannot be relayed either.
	// Left long enough the client expires outright. Observed idle on a devnet at 260
	// blocks behind, with proofs already failing "historical state is not available".
	//
	// Sizing is bounded by the L1 state window (~128 blocks on a default geth), NOT by
	// the trusting period — expiry is far more forgiving. Two things eat that budget:
	// the client only ever advances to FINALIZED headers, so it sits ~64 blocks behind
	// the head by construction, and it can drift a further interval's worth between
	// ticks (~40 blocks at mainnet block times). Worst case here is therefore ~104
	// blocks behind — inside the window, but without much room. Shorten this before
	// raising ethClientUpdateLag.
	ethClientPollInterval = 8 * time.Minute
)

// ethClientUpdater advances the Ethereum light client an L2 module is pinned to,
// on demand, at the moment a header build needs it.
//
// An L2 header builder proves the rollup contracts at the L1 block that client
// trusts and cannot pick another, because the L2 wasm client verifies those proofs
// against the L1 state root read back from the same client. Freshness is therefore a
// precondition of building, not a background nicety: a stalled client puts the proof
// block outside the execution node's servable state, and a game or assertion covering
// a recent L2 block does not exist at an old L1 block at all (#276).
//
// A deployment that also relays ETH→Cosmos advances the client as a side effect of
// that direction; an L2-only deployment has nothing doing it, which is what this
// covers.
type ethClientUpdater struct {
	worker    *services.Worker
	txHandler services.TransactionHandler
	svcCtx    services.Context
	l1Head    func(context.Context) (uint64, error)

	mu          sync.Mutex
	lastAttempt time.Time
}

func newEthClientUpdater(
	worker *services.Worker,
	txHandler services.TransactionHandler,
	svcCtx services.Context,
	l1Head func(context.Context) (uint64, error),
) *ethClientUpdater {
	return &ethClientUpdater{worker: worker, txHandler: txHandler, svcCtx: svcCtx, l1Head: l1Head}
}

// updateIfStale advances the client when it has fallen more than ethClientUpdateLag
// behind the L1 head. Errors are returned for logging but are not fatal to the
// caller: the client may still be recent enough to prove against, or another relay
// direction may be advancing it.
func (r *ethClientUpdater) updateIfStale(ctx context.Context, l1ClientID string, trustedBlock uint64) error {
	head, err := r.l1Head(ctx)
	if err != nil {
		return fmt.Errorf("eth client %s: read L1 head: %w", l1ClientID, err)
	}
	if head <= trustedBlock+ethClientUpdateLag {
		return nil
	}

	r.mu.Lock()
	if time.Since(r.lastAttempt) < ethClientUpdateCooldown {
		r.mu.Unlock()
		return nil
	}
	r.lastAttempt = time.Now()
	r.mu.Unlock()

	result, err := r.worker.BuildEthClientUpdateMsgs(r.svcCtx)
	if err != nil {
		return fmt.Errorf("eth client %s: build update (trusted %d, L1 head %d): %w",
			l1ClientID, trustedBlock, head, err)
	}
	if result == nil || len(result.Msgs) == 0 {
		return nil // already covers the latest finalized header
	}
	return r.submitWithSkewRetry(ctx, l1ClientID, result.Msgs)
}

// submitWithSkewRetry resubmits the same messages through the shared Cosmos batch
// path (serializing with every other sender on that account) when the only thing
// wrong is the host clock trailing the update's signature slot.
func (r *ethClientUpdater) submitWithSkewRetry(ctx context.Context, l1ClientID string, msgs []any) error {
	var err error
	for attempt := 0; attempt <= ethClientSkewRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(ethClientSkewBackoff):
			}
		}
		err = r.txHandler.SendCosmosTxBatch(ctx, r.svcCtx, msgs)
		if err == nil {
			return nil
		}
		if !isSignatureSlotSkew(err) {
			return fmt.Errorf("eth client %s: submit update: %w", l1ClientID, err)
		}
	}
	return fmt.Errorf("eth client %s: submit update, still ahead of the host clock after %d retries: %w",
		l1ClientID, ethClientSkewRetries, err)
}

// isSignatureSlotSkew reports whether the client rejected an update only because its
// signature slot is ahead of the slot computed from the Cosmos block timestamp.
func isSignatureSlotSkew(err error) bool {
	return err != nil && strings.Contains(err.Error(), "is more recent than the calculated current slot")
}

// buildEthClientUpdater wires a refresher for one L2 module's pinned client, and
// returns a cleanup for the Cosmos RPC client it dials.
func buildEthClientUpdater(
	cfg l2ToCosmosConfig,
	l1ClientID string,
	txHandler services.TransactionHandler,
	l1Head func(context.Context) (uint64, error),
) (*ethClientUpdater, func(), error) {
	cosmosClient, err := rpchttp.New(cfg.TmRpcUrl, "/websocket")
	if err != nil {
		return nil, nil, fmt.Errorf("eth client update: create Cosmos rpc client: %w", err)
	}
	if err := cosmosClient.Start(); err != nil {
		return nil, nil, fmt.Errorf("eth client update: start Cosmos rpc client: %w", err)
	}
	svcCtx := services.NewCtxWithBeacon(cosmosClient, nil, nil, "", cfg.EthBeaconAPIURL, l1ClientID)
	worker := services.NewWorker(txHandler, nil)
	return newEthClientUpdater(worker, txHandler, svcCtx, l1Head),
		func() { _ = cosmosClient.Stop() }, nil
}
