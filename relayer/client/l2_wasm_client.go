package client

import (
	"context"
	"encoding/json"
	"fmt"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

// GetL2PinnedEthClientID reads the Ethereum client id an L2 (rollup) wasm client is
// anchored to, from the client's own state on Cosmos.
//
// This is the authoritative copy. The relayer config carries the same id so header
// builders know which client's L1 block to pin proofs to, but the contract only ever
// consults the one baked into its client state at creation. When the two disagree
// nothing fails at creation — it fails later, on every update, with an error that
// names neither client: the contract asks the host for the Ethereum consensus state
// at a slot that only exists on the *other* client, ibc-go answers NotFound as a gRPC
// status error, and wasmd redacts it to "codespace: undefined, code: 1".
func GetL2PinnedEthClientID(cosmosClient *rpchttp.HTTP, l2ClientID string) (string, error) {
	queryReq := &clienttypes.QueryClientStateRequest{ClientId: l2ClientID}
	reqBytes, err := proto.Marshal(queryReq)
	if err != nil {
		return "", fmt.Errorf("marshal client state query for %s: %w", l2ClientID, err)
	}

	result, err := cosmosClient.ABCIQuery(context.Background(), "/ibc.core.client.v1.Query/ClientState", reqBytes)
	if err != nil {
		return "", fmt.Errorf("query client state %s: %w", l2ClientID, err)
	}
	if result.Response.Code != 0 {
		return "", fmt.Errorf("query client state %s failed with code %d: %s",
			l2ClientID, result.Response.Code, result.Response.Log)
	}

	var queryResp clienttypes.QueryClientStateResponse
	if err := proto.Unmarshal(result.Response.Value, &queryResp); err != nil {
		return "", fmt.Errorf("unmarshal client state response for %s: %w", l2ClientID, err)
	}
	if queryResp.ClientState == nil {
		return "", fmt.Errorf("client %s has no client state on Cosmos", l2ClientID)
	}

	var wasmClientState ibcwasmtypes.ClientState
	if err := proto.Unmarshal(queryResp.ClientState.Value, &wasmClientState); err != nil {
		return "", fmt.Errorf("unmarshal wasm client state for %s: %w", l2ClientID, err)
	}

	// Only the pinned id is decoded — the rest of the profile is chain-specific and
	// belongs to the wasm client, not the relayer.
	var state struct {
		Profile struct {
			Common struct {
				EthereumClient struct {
					ClientID string `json:"client_id"`
				} `json:"ethereum_client"`
			} `json:"common"`
		} `json:"profile"`
	}
	if err := json.Unmarshal(wasmClientState.Data, &state); err != nil {
		return "", fmt.Errorf("decode L2 client state for %s: %w", l2ClientID, err)
	}
	id := state.Profile.Common.EthereumClient.ClientID
	if id == "" {
		return "", fmt.Errorf("L2 client %s carries no profile.common.ethereum_client.client_id", l2ClientID)
	}
	return id, nil
}
