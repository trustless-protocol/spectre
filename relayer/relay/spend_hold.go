// This file holds the spending-cap hold: what the module does when an operator
// cap (ETH_MAX_FEE_PER_GAS, ETH_MAX_PRIORITY_FEE_PER_GAS, COSMOS_MAX_GAS,
// COSMOS_MAX_FEE) refuses to sign a transaction.
//
// The caps are checked at signing, after the client update and the proofs are
// built. A refused batch is re-queued and comes back within 15 seconds, so a cap
// that stayed binding rebuilt both -- a Groth16 proof on a Cosmos->EVM path --
// and logged a refusal every 15 seconds (#575 M3). Nothing reached a chain and
// nothing was lost; the cost was the prover and the log.
//
// So the module holds the direction, on the same 1m -> 15m ladder as its other
// holds. Packets stay queued and tracked; the timeout scanner still sees them.
package relay

import (
	"errors"
	"log"
	"time"

	"relayer/services"
)

// spendHeldOff reports whether the direction is inside a spending-cap hold. Once
// the hold expires one pass goes through: a refusal re-arms it, and a send that
// lands lifts it.
//
// Callers must hold batchMu.
func (m *Module) spendHeldOff(now time.Time) bool {
	if m.spendProbeAt.IsZero() {
		return false
	}
	if now.Before(m.spendProbeAt) {
		return true
	}
	m.spendProbeAt = time.Time{} // try this pass; holdIfSpendCapped re-arms it
	return false
}

// holdIfSpendCapped arms or widens the hold when err is a spending-cap refusal,
// and reports whether it did. what names the send that was refused.
//
// Callers must hold batchMu.
func (m *Module) holdIfSpendCapped(err error, what string) bool {
	if !errors.Is(err, services.ErrSpendCapExceeded) {
		return false
	}
	m.spendBackoff = nextPeriodicUpdateBackoff(m.spendBackoff)
	m.spendProbeAt = time.Now().Add(m.spendBackoff)
	log.Printf("WARN [%s Relay] %s: %v; nothing was sent. Relaying is held for %s, then tried again. "+
		"Queued packets are kept and still time out normally.", m.name, what, err, m.spendBackoff)
	return true
}

// noteSpendAllowed lifts the hold after a send landed.
//
// Callers must hold batchMu.
func (m *Module) noteSpendAllowed() {
	if m.spendBackoff == 0 && m.spendProbeAt.IsZero() {
		return
	}
	m.spendBackoff = 0
	m.spendProbeAt = time.Time{}
	log.Printf("[%s Relay] a send landed within the spending caps; hold lifted", m.name)
}
