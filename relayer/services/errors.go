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

// ErrUnresolvedEVMNonce marks a failure whose transaction may STILL EXECUTE.
//
// It is not the same claim as ErrPermanentRelayFailure, and conflating the two
// is what made this necessary. Permanent means "this attempt cannot succeed --
// the chain rejected it deterministically", and relay/batch.go acts on that by
// ISOLATING a multi-packet batch: it resubmits each packet on its own so one
// poison packet does not discard its valid siblings. That is right for a revert
// and wrong here.
//
// Reported by @DongLieu on #485. When the relayer gives up waiting for a
// transaction, the abandoned nonce has one of three outcomes and the node has
// told us none of them: the cancellation replaced it, the original still
// executes, or the original already did. Resubmitting the payload -- as a whole
// batch, as isolated packets, or as a timeout retry -- can run the operation a
// second time while the first is still live, and the invalidated nonce cache
// means the second attempt may take a LOWER nonce and mine first.
//
// So a failure carrying this sentinel means: do not resubmit anything on this
// path until the outcome at that nonce is established. The packets stay in the
// pending tracker, which is the durable record, and settle by their terminal
// event or are refunded by the timeout scanner.
var ErrUnresolvedEVMNonce = errors.New("the outcome of an abandoned EVM nonce is unresolved; the payload must not be resubmitted")

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

// CommittedBeforeFailure reads the delivered prefix out of a batch failure.
//
// A split batch can commit its leading sub-batches and then fail, and that count
// is the only thing that tells a caller how much of its batch to re-queue.
// Dropping it forces the caller to choose between re-submitting what was already
// delivered and discarding what was not.
//
// Any other error means nothing was delivered: the send either never reached a
// transaction or the single transaction it built failed as a whole. It lives here
// rather than in an adapter because both Cosmos-bound destinations need it and
// both already depend on this package.
func CommittedBeforeFailure(err error) int {
	var partial *BatchPartialError
	if errors.As(err, &partial) && partial.SucceededCount > 0 {
		return partial.SucceededCount
	}
	return 0
}
