package cosmos

import (
	"encoding/json"
	"testing"

	relayerclient "relayer/client"
)

func mustBeaconHeader(t *testing.T, signatureSlot string) []byte {
	t.Helper()
	payload, err := json.Marshal(relayerclient.EthereumHeader{
		ConsensusUpdate: relayerclient.LightClientUpdate{SignatureSlot: signatureSlot},
	})
	if err != nil {
		t.Fatalf("marshal header: %v", err)
	}
	return payload
}

func TestHighestBeaconSignatureSlot_UsesLatestHeader(t *testing.T) {
	slot, err := highestBeaconSignatureSlot([][]byte{
		mustBeaconHeader(t, "123"),
		mustBeaconHeader(t, "130"),
		mustBeaconHeader(t, "127"),
	})
	if err != nil {
		t.Fatalf("highest slot: %v", err)
	}
	if slot != 130 {
		t.Fatalf("slot = %d, want 130", slot)
	}
}

func TestHighestBeaconSignatureSlot_RejectsMalformedPayload(t *testing.T) {
	if _, err := highestBeaconSignatureSlot([][]byte{[]byte("not-json")}); err == nil {
		t.Fatal("malformed header must fail")
	}
}

func TestHighestBeaconSignatureSlot_RejectsInvalidSlot(t *testing.T) {
	if _, err := highestBeaconSignatureSlot([][]byte{mustBeaconHeader(t, "not-a-slot")}); err == nil {
		t.Fatal("invalid signature slot must fail")
	}
}
