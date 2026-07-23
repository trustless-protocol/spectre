package l2rollup

import (
	"testing"

	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
)

// TestBuildWasmUpdateClient verifies the L2 header is wrapped in a wasm
// ClientMessage inside a MsgUpdateClient — the exact wrapper Dũng's client (and
// the existing beacon client) consumes.
func TestBuildWasmUpdateClient(t *testing.T) {
	l2Header := []byte(`{"height":123,"state_root":"0xabc"}`)
	msg, err := buildWasmUpdateClient("cosmos1signer", "l2-op-0", l2Header)
	if err != nil {
		t.Fatalf("build: %v", err)
	}
	if msg.ClientId != "l2-op-0" || msg.Signer != "cosmos1signer" {
		t.Fatalf("wrong envelope: client=%q signer=%q", msg.ClientId, msg.Signer)
	}
	// The ClientMessage Any must decode back to a wasm ClientMessage carrying the
	// L2 header verbatim.
	cm, ok := msg.ClientMessage.GetCachedValue().(*ibcwasmtypes.ClientMessage)
	if !ok {
		t.Fatalf("client message is not a wasm ClientMessage: %T", msg.ClientMessage.GetCachedValue())
	}
	if string(cm.Data) != string(l2Header) {
		t.Fatalf("L2 header not preserved: got %q", cm.Data)
	}
}
