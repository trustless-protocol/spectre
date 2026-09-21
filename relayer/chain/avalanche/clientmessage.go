// Package avalanche is the Avalanche C-Chain provider for the relayer.
// Avalanche is an L1, not a rollup: it reuses the shared attested-EVM
// machinery in chain/l2rollup (the Source, the wire header, the proof
// envelopes, the relay engines) but carries its own trust model — every
// update is authenticated by the primary network's stake-weighted warp BLS
// aggregate, collected through the signature-aggregator sidecar. This package
// holds everything Avalanche-specific: the coreth header reading/hashing, the
// warp message primitives, the warp header builder, and the settled-height
// proof mapping the asynchronous-execution upgrades require.
package avalanche

import (
	"encoding/json"
	"fmt"

	"relayer/chain/l2rollup"
)

// WarpSignedHeader mirrors avalanche-light-client `WarpSignedHeader` — the
// update shape the Avalanche warp wasm client accepts: a canonical coreth
// header, the primary network's aggregate BLS signature over the block-hash
// warp message (signer bitset + 96-byte signature, base64 JSON like every
// CosmWasm Binary), and the router account proof at the header's settled
// height.
type WarpSignedHeader struct {
	Header       l2rollup.CanonicalEvmHeader `json:"header"`
	SignerBitSet []byte                      `json:"signer_bit_set"`
	Signature    []byte                      `json:"signature"`
	RouterProof  l2rollup.EvmAccountProof    `json:"router_proof"`
}

// EncodeClientMessage marshals the warp header into the shared ClientMessage
// envelope {"type":"header","value":<header>} (serde tag="type",
// content="value"), the same envelope the attested clients use.
func (h *WarpSignedHeader) EncodeClientMessage() ([]byte, error) {
	if len(h.SignerBitSet) == 0 || len(h.Signature) != warpSignatureLen {
		return nil, fmt.Errorf("avalanche client message: warp header requires a signer bitset and a %d-byte aggregate signature", warpSignatureLen)
	}
	value, err := json.Marshal(h)
	if err != nil {
		return nil, fmt.Errorf("avalanche client message: marshal header value: %w", err)
	}
	envelope := struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}{Type: "header", Value: value}
	out, err := json.Marshal(envelope)
	if err != nil {
		return nil, fmt.Errorf("avalanche client message: marshal envelope: %w", err)
	}
	return out, nil
}

// WarpSignedHeader satisfies the shared ClientMessage contract.
var _ l2rollup.ClientMessage = (*WarpSignedHeader)(nil)
