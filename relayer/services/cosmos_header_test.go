package services

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto/ed25519"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	commettypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	updateClientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
)

type mockTxHandler struct {
	TransactionHandler
}

type mockProver struct {
	Prover
}

func encodeSolidityBytes(data []byte) string {
	length := len(data)
	paddedLen := ((length + 31) / 32) * 32
	padded := make([]byte, paddedLen)
	copy(padded, data)

	var buf bytes.Buffer
	var offset [32]byte
	offset[31] = 0x20
	buf.Write(offset[:])

	var lenBuf [32]byte
	binary.BigEndian.PutUint64(lenBuf[24:], uint64(length))
	buf.Write(lenBuf[:])

	buf.Write(padded)

	return "0x" + hex.EncodeToString(buf.Bytes())
}

func base64Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func TestHandleCosmosHeaderFetchFailure(t *testing.T) {
	privKey := ed25519.GenPrivKey()
	pubKey := privKey.PubKey()
	val := commettypes.NewValidator(pubKey, 100)
	valSet100 := commettypes.NewValidatorSet([]*commettypes.Validator{val})
	valSet101 := commettypes.NewValidatorSet([]*commettypes.Validator{val})

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(interfaceRegistry)
	stakingParams := stakingtypes.QueryParamsResponse{
		Params: stakingtypes.Params{
			UnbondingTime: 1000 * time.Second,
		},
	}
	stakingParamsBz, err := cdc.Marshal(&stakingParams)
	if err != nil {
		t.Fatalf("failed to marshal staking params: %v", err)
	}

	cases := []struct {
		name         string
		simulateZero bool
	}{
		{
			name:         "HeaderByNumber returns RPC error",
			simulateZero: false,
		},
		{
			name:         "HeaderByNumber returns zero timestamp",
			simulateZero: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				var req struct {
					JSONRPC string      `json:"jsonrpc"`
					ID      interface{} `json:"id"`
					Method  string      `json:"method"`
					Params  interface{} `json:"params"`
				}
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}

				switch req.Method {
				case "status":
					res := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]interface{}{
							"node_info": map[string]interface{}{
								"id":      "1234567890abcdef1234567890abcdef12345678",
								"network": "test-ibc-eth",
							},
							"sync_info": map[string]interface{}{
								"latest_block_height": "100",
								"latest_block_time":   "2026-06-10T00:00:00Z",
							},
						},
					}
					json.NewEncoder(w).Encode(res)

				case "commit":
					res := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]interface{}{
							"signed_header": map[string]interface{}{
								"header": map[string]interface{}{
									"version": map[string]interface{}{
										"block": "11",
										"app":   "0",
									},
									"chain_id":             "test-ibc-eth",
									"height":               "100",
									"time":                 "2026-06-10T00:00:00Z",
									"validators_hash":      hex.EncodeToString(valSet100.Hash()),
									"next_validators_hash": hex.EncodeToString(valSet101.Hash()),
								},
								"commit": map[string]interface{}{
									"height": "100",
									"round":  0,
									"block_id": map[string]interface{}{
										"hash": "1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
									},
								},
							},
						},
					}
					json.NewEncoder(w).Encode(res)

				case "validators":
					var heightStr string
					if paramsMap, ok := req.Params.(map[string]interface{}); ok {
						if h, exists := paramsMap["height"]; exists {
							heightStr = fmt.Sprintf("%v", h)
						}
					}
					var valSet *commettypes.ValidatorSet
					if heightStr == "101" {
						valSet = valSet101
					} else {
						valSet = valSet100
					}

					valsJson := []map[string]interface{}{}
					for _, v := range valSet.Validators {
						valsJson = append(valsJson, map[string]interface{}{
							"address": hex.EncodeToString(v.Address),
							"pub_key": map[string]interface{}{
								"type":  "tendermint/PubKeyEd25519",
								"value": base64Encode(v.PubKey.Bytes()),
							},
							"voting_power":      fmt.Sprintf("%d", v.VotingPower),
							"proposer_priority": fmt.Sprintf("%d", v.ProposerPriority),
						})
					}

					res := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]interface{}{
							"block_height": heightStr,
							"validators":   valsJson,
							"count":        fmt.Sprintf("%d", len(valsJson)),
							"total":        fmt.Sprintf("%d", len(valsJson)),
						},
					}
					json.NewEncoder(w).Encode(res)

				case "abci_query":
					res := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result": map[string]interface{}{
							"response": map[string]interface{}{
								"code":  0,
								"value": base64Encode(stakingParamsBz),
							},
						},
					}
					json.NewEncoder(w).Encode(res)

				case "eth_call":
					clientState := relayerclient.ClientState{
						ChainId:    "test-ibc-eth",
						TrustLevel: relayerclient.TrustThreshold{Numerator: 2, Denominator: 3},
						LatestHeight: updateClientContract.IICS02ClientMsgsHeight{
							RevisionNumber: 0,
							RevisionHeight: 100,
						},
						TrustingPeriod:  600,
						UnbondingPeriod: 1000,
						IsFrozen:        false,
					}
					encodedState, _ := relayerclient.EncodeClientState(clientState)
					solBytesHex := encodeSolidityBytes(encodedState)

					res := map[string]interface{}{
						"jsonrpc": "2.0",
						"id":      req.ID,
						"result":  solBytesHex,
					}
					json.NewEncoder(w).Encode(res)

				case "eth_getBlockByNumber":
					if tc.simulateZero {
						res := map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      req.ID,
							"result": map[string]interface{}{
								"number":     "0x64",
								"hash":       "0x0000000000000000000000000000000000000000000000000000000000000000",
								"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
								"timestamp":  "0x0",
							},
						}
						json.NewEncoder(w).Encode(res)
					} else {
						res := map[string]interface{}{
							"jsonrpc": "2.0",
							"id":      req.ID,
							"error": map[string]interface{}{
								"code":    -32000,
								"message": "simulated header fetch failure",
							},
						}
						json.NewEncoder(w).Encode(res)
					}

				default:
					http.Error(w, fmt.Sprintf("unsupported method %s", req.Method), http.StatusBadRequest)
				}
			}))
			defer mockServer.Close()

			rpcClient, err := rpc.DialHTTPWithClient(mockServer.URL, mockServer.Client())
			if err != nil {
				t.Fatalf("failed to dial mock server: %v", err)
			}
			ethClient := ethclient.NewClient(rpcClient)

			cosmosClient, err := rpchttp.NewWithClient(mockServer.URL, "/websocket", mockServer.Client())
			if err != nil {
				t.Fatalf("failed to create cosmos client: %v", err)
			}

			ctx := NewCtx(cosmosClient, ethClient)
			ctx.Config.FetchTimeout = time.Second * 15
			ctx.SetClient(common.HexToAddress("0x1111111111111111111111111111111111111111"))

			txHandler := &mockTxHandler{}
			prover := &mockProver{}
			s := New(nil, txHandler, prover, Config{}, Config{})

			s.lastCosmosAppHashHeight = 1000 // avoid waiting for AppHash

			batch := CosmosBatch{
				Packets: []CosmosPacket{
					makeCosmosPacket(1, CosmosSend),
				},
			}

			s.handleCosmos(ctx, batch)

			// Check that packets were re-queued transiently
			s.BatchBuilder.cosmosMtx.Lock()
			defer s.BatchBuilder.cosmosMtx.Unlock()

			if len(s.BatchBuilder.cosmosPackets) != 1 {
				t.Fatalf("expected 1 packet re-queued, got %d", len(s.BatchBuilder.cosmosPackets))
			}
			if s.BatchBuilder.cosmosPackets[0].Packet.Sequence != 1 {
				t.Fatalf("expected sequence 1, got %d", s.BatchBuilder.cosmosPackets[0].Packet.Sequence)
			}
		})
	}
}
