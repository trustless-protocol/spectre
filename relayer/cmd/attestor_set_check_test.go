package main

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"testing"

	"relayer/rpcmock"

	"attestor/types/attestation"

	abci "github.com/cometbft/cometbft/abci/types"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

func testAttestorSet(t *testing.T, n int, threshold uint16) attestation.AttestorConfig {
	t.Helper()
	keys := make([][]byte, 0, n)
	for i := 0; i < n; i++ {
		pub, _, err := ed25519.GenerateKey(nil)
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, pub)
	}
	sort.Slice(keys, func(i, j int) bool { return bytes.Compare(keys[i], keys[j]) < 0 })
	set := attestation.AttestorConfig{PublicKeys: keys, Threshold: threshold}
	if _, err := set.SetHash(); err != nil {
		t.Fatalf("test set does not hash: %v", err)
	}
	return set
}

// l2ClientNode serves the client state of one L2 wasm client holding set.
func l2ClientNode(t *testing.T, set attestation.AttestorConfig) *rpcmock.CosmosNode {
	t.Helper()
	data, err := json.Marshal(map[string]any{
		"latest_height": 10, "frozen_height": nil, "profile": map[string]any{}, "attestors": set,
	})
	if err != nil {
		t.Fatal(err)
	}
	wasmState, err := codectypes.NewAnyWithValue(&ibcwasmtypes.ClientState{Data: data})
	if err != nil {
		t.Fatal(err)
	}
	response, err := (&clienttypes.QueryClientStateResponse{ClientState: wasmState}).Marshal()
	if err != nil {
		t.Fatal(err)
	}
	node := rpcmock.NewCosmosNode(t)
	node.Handle("abci_query", func(json.RawMessage) (any, error) {
		return &coretypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: response}}, nil
	})
	return node
}

func l2AttestorConfig(url string, set attestation.AttestorConfig) *appConfig {
	return &appConfig{
		HasL2ToCosmos: true,
		L2ToCosmosConfig: l2ToCosmosConfig{
			TmRpcUrl: url, L2WasmClientID: "08-wasm-3", Attestors: set,
		},
	}
}

// #575: the client can rotate its attestor set after creation. A config still on
// the old set collects signatures the client refuses, and only refused headers
// said so.
func TestValidateL2AttestorSet(t *testing.T) {
	onChain := testAttestorSet(t, 3, 2)

	t.Run("the configured set is the client's", func(t *testing.T) {
		node := l2ClientNode(t, onChain)
		if findings := validateL2AttestorSet(context.Background(), l2AttestorConfig(node.URL(), onChain)); len(findings) != 0 {
			t.Fatalf("findings = %+v, want none", findings)
		}
	})

	t.Run("a rotated set is a finding", func(t *testing.T) {
		node := l2ClientNode(t, onChain)
		stale := testAttestorSet(t, 3, 2)
		findings := validateL2AttestorSet(context.Background(), l2AttestorConfig(node.URL(), stale))
		if len(findings) != 1 || findings[0].key != "l2_to_cosmos.attestors" || !strings.Contains(findings[0].detail, "08-wasm-3") {
			t.Fatalf("findings = %+v, want one naming the client", findings)
		}
	})

	t.Run("a different threshold is a finding", func(t *testing.T) {
		node := l2ClientNode(t, onChain)
		other := onChain
		other.Threshold = 3
		if findings := validateL2AttestorSet(context.Background(), l2AttestorConfig(node.URL(), other)); len(findings) != 1 {
			t.Fatalf("findings = %+v, want one", findings)
		}
	})

	t.Run("an unreadable client is not a finding", func(t *testing.T) {
		node := rpcmock.NewCosmosNode(t)
		node.Handle("abci_query", func(json.RawMessage) (any, error) { return nil, errors.New("node is catching up") })
		if findings := validateL2AttestorSet(context.Background(), l2AttestorConfig(node.URL(), onChain)); len(findings) != 0 {
			t.Fatalf("findings = %+v, want none: availability is not evidence", findings)
		}
	})

	t.Run("no L2 source, nothing to check", func(t *testing.T) {
		if findings := validateL2AttestorSet(context.Background(), &appConfig{}); findings != nil {
			t.Fatalf("findings = %+v", findings)
		}
	})
}
