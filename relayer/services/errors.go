package services

import "errors"

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
