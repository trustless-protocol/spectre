package services

import (
	"testing"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

// TestShouldRelayCosmosTimeoutToEth pins the direction-detection logic that
// decides whether a CosmosTimeout event needs a follow-up timeoutPacket on
// ETH. The old guard compared SourceClient against the router client ID — a
// mismatch in IBC v2 Eureka because:
//
//   - For a Cosmos→ETH packet, the Cosmos-emitted event has
//     SourceClient="08-wasm-0" (Cosmos's client identifying the ETH chain) and
//     DestinationClient="cosmoshub-1" (ETH's client identifying Cosmos).
//   - For an ETH→Cosmos packet, it's mirrored: source=cosmoshub-1, dest=08-wasm-0.
//
// `CosmosRouterClientID` in the relayer config is the ETH-side router's name
// for this Cosmos chain ("cosmoshub-1"), so the correct field to compare is
// `DestinationClient`. The old comparison silently inverted both branches.
func TestShouldRelayCosmosTimeoutToEth(t *testing.T) {
	tests := []struct {
		name                 string
		packet               *channeltypesv2.Packet
		cosmosRouterClientID string
		want                 bool
	}{
		{
			name: "Cosmos-originated packet skips ETH relay",
			packet: &channeltypesv2.Packet{
				SourceClient:      "08-wasm-0",   // Cosmos's client for ETH
				DestinationClient: "cosmoshub-1", // ETH's client for Cosmos
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false, // Cosmos's MsgTimeout already refunded the sender
		},
		{
			name: "ETH-originated packet is relayed",
			packet: &channeltypesv2.Packet{
				SourceClient:      "cosmoshub-1", // ETH's client for Cosmos
				DestinationClient: "08-wasm-0",   // Cosmos's client for ETH
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true, // ETH still holds the commitment, needs ICS26Router.timeoutPacket
		},
		{
			name:                 "nil packet is not relayed",
			packet:               nil,
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false,
		},
		{
			name: "unrelated client IDs are relayed (treated as ETH-originated)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true,
		},
		{
			name: "empty cosmosRouterClientID is not relayed (misconfiguration)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "",
			want:                 false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldRelayCosmosTimeoutToEth(tc.packet, tc.cosmosRouterClientID)
			if got != tc.want {
				t.Fatalf("shouldRelayCosmosTimeoutToEth=%v want=%v", got, tc.want)
			}
		})
	}
}
