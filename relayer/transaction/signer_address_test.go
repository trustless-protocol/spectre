package transaction

import (
	"strings"
	"testing"

	sdkbech32 "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ethereum/go-ethereum/common"
)

// Devnet-only sample key (relayer/.env); safe to embed in tests.
const testCosmosPrivKeyHex = "a81f9eb900c02f35d28ec80d1da79dde52cc4cc37d6fbaef5768ab59f32cbfdd"

func TestEthSignerAddress(t *testing.T) {
	t.Setenv("ETH_PRIVATE_KEY", testCosmosPrivKeyHex)

	address, err := (&Handler{}).EthSignerAddress()
	if err != nil {
		t.Fatalf("EthSignerAddress: %v", err)
	}
	if want := common.HexToAddress("0xd58ed941051839DB5ffe8735905881d3C7460cE1"); address != want {
		t.Fatalf("EthSignerAddress = %s, want %s", address.Hex(), want.Hex())
	}
}

func TestCosmosSignerAddressPrefix(t *testing.T) {
	tests := []struct {
		name       string
		prefix     string // COSMOS_ADDRESS_PREFIX; "" = unset
		wantPrefix string
	}{
		{name: "default is cosmos", prefix: "", wantPrefix: "cosmos"},
		{name: "realio", prefix: "realio", wantPrefix: "realio"},
		{name: "osmo", prefix: "osmo", wantPrefix: "osmo"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("COSMOS_PRIVATE_KEY", testCosmosPrivKeyHex)
			t.Setenv("COSMOS_ADDRESS_PREFIX", tc.prefix)

			addr, err := (&Handler{}).CosmosSignerAddress()
			if err != nil {
				t.Fatalf("CosmosSignerAddress: %v", err)
			}
			if !strings.HasPrefix(addr, tc.wantPrefix+"1") {
				t.Fatalf("address %q does not start with %q", addr, tc.wantPrefix+"1")
			}
			hrp, _, err := sdkbech32.DecodeAndConvert(addr)
			if err != nil {
				t.Fatalf("address %q is not valid bech32: %v", addr, err)
			}
			if hrp != tc.wantPrefix {
				t.Fatalf("decoded hrp = %q, want %q", hrp, tc.wantPrefix)
			}
		})
	}
}

// The prefix must only change the encoding, never the underlying account bytes.
func TestCosmosSignerAddressPrefixSameAccount(t *testing.T) {
	t.Setenv("COSMOS_PRIVATE_KEY", testCosmosPrivKeyHex)

	t.Setenv("COSMOS_ADDRESS_PREFIX", "")
	defaultAddr, err := (&Handler{}).CosmosSignerAddress()
	if err != nil {
		t.Fatalf("CosmosSignerAddress (default): %v", err)
	}

	t.Setenv("COSMOS_ADDRESS_PREFIX", "realio")
	realioAddr, err := (&Handler{}).CosmosSignerAddress()
	if err != nil {
		t.Fatalf("CosmosSignerAddress (realio): %v", err)
	}

	_, defaultBytes, err := sdkbech32.DecodeAndConvert(defaultAddr)
	if err != nil {
		t.Fatalf("decode %q: %v", defaultAddr, err)
	}
	_, realioBytes, err := sdkbech32.DecodeAndConvert(realioAddr)
	if err != nil {
		t.Fatalf("decode %q: %v", realioAddr, err)
	}
	if string(defaultBytes) != string(realioBytes) {
		t.Fatalf("account bytes differ between prefixes: %x vs %x", defaultBytes, realioBytes)
	}
}
