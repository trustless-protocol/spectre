package services

import (
	"errors"
	"fmt"
)

// ErrPermanentRelayFailure marks a relay failure as deterministic / on-chain:
// the tx was actually included and reverted (ETH receipt.Status==0, or a Cosmos
// DeliverTx non-zero code). Such a failure will recur on retry with the same
// inputs, so the packet's retry budget should be consumed and it should
// eventually be dead-lettered.
//
// Infrastructure failures (RPC down, beacon finality not yet reached, light-
// client build failure, broadcast timeout, account-sequence mismatch at CheckTx)
// are NOT wrapped with this — they are transient and must be re-queued without
// consuming the budget, so a valid packet is never dropped just because the
// infra was briefly unavailable (issue #80 review).
//
// Tx handlers (relayer/transaction) wrap their revert errors with this via
// fmt.Errorf("...: %w", ErrPermanentRelayFailure); callers classify with
// errors.Is.
var ErrPermanentRelayFailure = errors.New("permanent relay failure (on-chain revert)")

// ErrValidatorCacheRace marks an updateClient revert caused by the validator
// cache changing after the relayer built a cached/delta update. Rebuilding the
// update from fresh on-chain state, or retrying with the full validator set, can
// make progress, so packet retry budgets should not be consumed.
var ErrValidatorCacheRace = errors.New("validator cache race")

// BatchPartialError is returned when a split batch is partially successful.
// It tracks how many messages at the beginning of the batch were successfully
// committed before a subsequent sub-batch encountered an error.
type BatchPartialError struct {
	SucceededCount int
	Err            error
}

func (e *BatchPartialError) Error() string {
	return fmt.Sprintf("partial batch failure: %d messages succeeded: %v", e.SucceededCount, e.Err)
}

func (e *BatchPartialError) Unwrap() error {
	return e.Err
}
