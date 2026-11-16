// Which chain an EVM transaction is signed for (#575, "RPC chain id per send").
//
// The id comes from eth_chainId and was read again on every send. A signature is
// valid on exactly the chain it names, so an endpoint that starts answering a
// different id is handed transactions valid on that other chain, at a fee it
// also chose. That matters whenever the relayer key is funded on more than one
// EVM chain, which the runbooks allow.
//
// An endpoint is now held to the id it first reported, and under `start` that
// first id must be one the process took a nonce-domain lock for.
package transaction

import (
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

// ErrEVMChainIDChanged marks a transaction that was not signed because its
// endpoint reported a chain id this process does not sign for. It is not a
// permanent failure: it says the endpoint is wrong, not the packet.
var ErrEVMChainIDChanged = errors.New("EVM endpoint reports an unexpected chain id")

// AllowEVMChainIDs names the chains this process may sign for: the ids `start`
// resolved when it took the EVM nonce-domain locks. A command that does not call
// it holds each endpoint to the first id that endpoint reports.
func (h *Handler) AllowEVMChainIDs(ids []string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.allowedChainIDs = make(map[string]struct{}, len(ids))
	for _, id := range ids {
		h.allowedChainIDs[id] = struct{}{}
	}
}

// checkEVMChainID refuses a chain id the endpoint did not report before, or one
// outside the allowed set. It pins the id on first use.
func (h *Handler) checkEVMChainID(client *ethclient.Client, reported *big.Int) error {
	if reported == nil {
		return fmt.Errorf("endpoint returned no chain id: %w", ErrEVMChainIDChanged)
	}
	id := reported.String()

	h.mu.Lock()
	defer h.mu.Unlock()
	if pinned, ok := h.pinnedChainIDs[client]; ok {
		if pinned != id {
			return fmt.Errorf("endpoint reported chain id %s but this process signs for %s on it; not signed: %w",
				id, pinned, ErrEVMChainIDChanged)
		}
		return nil
	}
	if h.allowedChainIDs != nil {
		if _, ok := h.allowedChainIDs[id]; !ok {
			allowed := make([]string, 0, len(h.allowedChainIDs))
			for allowedID := range h.allowedChainIDs {
				allowed = append(allowed, allowedID)
			}
			sort.Strings(allowed)
			return fmt.Errorf("endpoint reported chain id %s, not a chain this process locked a nonce domain for at startup (%s); not signed: %w",
				id, strings.Join(allowed, ", "), ErrEVMChainIDChanged)
		}
	}
	if h.pinnedChainIDs == nil {
		h.pinnedChainIDs = make(map[*ethclient.Client]string)
	}
	h.pinnedChainIDs[client] = id
	return nil
}
