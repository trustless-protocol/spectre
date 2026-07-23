package services

import (
	"errors"
	"fmt"
)

// ErrPermanentRelayFailure marks a relay failure as deterministic / on-chain:
// the tx was actually included and reverted (ETH receipt.Status==0, or a Cosmos
// DeliverTx non-zero code). Such a failure will recur on retry with the same
// inputs, so the adapter module DROPS the packet (via chain.Permanent) instead
// of re-queueing it forever — a genuinely stuck packet is refunded by the
// timeout scanner once it expires undelivered.
//
// Infrastructure failures (RPC down, beacon finality not yet reached, light-
// client build failure, broadcast timeout, account-sequence mismatch at CheckTx)
// are NOT wrapped with this — they are transient and must be re-queued (with a
// waiting backoff), so a valid packet is never dropped just because the infra
// was briefly unavailable (issue #80 review).
//
// Tx handlers (relayer/transaction) wrap their revert errors with this via
// fmt.Errorf("...: %w", ErrPermanentRelayFailure); callers classify with
// errors.Is.
var ErrPermanentRelayFailure = errors.New("permanent relay failure (on-chain revert)")

// ErrValidatorCacheRace marks an updateClient revert caused by the validator
// cache changing after the relayer built a cached/delta update. Rebuilding the
// update from fresh on-chain state, or retrying with the full validator set, can
// make progress, so it is treated as transient (re-queued), not permanent.
var ErrValidatorCacheRace = errors.New("validator cache race")

// CosmosTxFailure preserves ABCI failure metadata across package boundaries so
// relay logic can classify specific failures without scraping formatted strings.
type CosmosTxFailure struct {
	Stage     string
	Code      uint32
	Codespace string
	Log       string
	Data      []byte
	Err       error
}

func (e *CosmosTxFailure) Error() string {
	if e == nil {
		return "<nil>"
	}
	msg := fmt.Sprintf("cosmos tx failed at %s with code %d", e.Stage, e.Code)
	if e.Codespace != "" {
		msg += fmt.Sprintf(" codespace=%s", e.Codespace)
	}
	if e.Log != "" {
		msg += fmt.Sprintf(": %s", e.Log)
	}
	return msg
}

func (e *CosmosTxFailure) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Err
}

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
