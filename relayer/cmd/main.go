package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	proto "github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	tendermintClient "relayer/client"
	"relayer/keys"
	"relayer/prover"
	"relayer/services"
	"relayer/transaction"
	utils "relayer/utils"
)

const (
	flagConfigPath     = "config"
	flagGPUProve       = "gpu-prove"
	flagOnlyOnce       = "only-once"
	flagProofType      = "proof-type"
	flagOutput         = "output"
	flagOutputPath     = "output-path"
	flagTrustLevel     = "trust-level"
	flagTrustingPeriod = "trusting-period"
	flagClockDrift     = "clock-drift"
	flagTrustedBlock   = "trusted-block"
	flagWasmChecksum   = "wasm-checksum"
	flagSource         = "source"
	flagBenchmark      = "benchmark"
	configFilePerm     = 0o600
)

// --- Config types for JSON config file ---

type cosmosToEthConfig struct {
	TmRpcUrl              string `json:"tm_rpc_url"`
	ICS26Address          string `json:"ics26_address"`
	ICS26ClientID         string `json:"ics26_client_id"`
	CosmosWasmClientID    string `json:"cosmos_wasm_client_id"`
	EthRpcUrl             string `json:"eth_rpc_url"`
	EthWsUrl              string `json:"eth_ws_url"`
	ICS07Client           string `json:"ics07_client"`
	WrapperVerifier       string `json:"wrapper_verifier"`
	Membership            string `json:"membership"`
	Misbehaviour          string `json:"misbehaviour"`
	UpdateClient          string `json:"update_client"`
	TrustingPeriod        uint32 `json:"trusting_period"`
	TrustLevel            string `json:"trust_level"`
	ProofType             string `json:"proof_type"`
	ClockDrift            uint32 `json:"clock_drift"`
	BeaconFinalityRetries uint32 `json:"beacon_finality_retries"`
	AppHashWaitRetries    uint32 `json:"app_hash_wait_retries"`
	AppHashWaitInterval   uint64 `json:"app_hash_wait_interval_seconds"`
	FetchTimeout          uint64 `json:"fetch_timeout"`
}

type ethToCosmosConfig struct {
	// BeaconUrl is the only field consumed for the ETH→Cosmos direction; the
	// Cosmos RPC, Eth RPC and ICS26 address all come from the cosmos_to_eth
	// module (used for both directions).
	BeaconUrl string `json:"eth_beacon_api_url"`
}

type configModule struct {
	Name     string          `json:"name"`
	SrcChain string          `json:"src_chain"`
	Config   json.RawMessage `json:"config"`
}

type serverConfig struct {
	LogLevel string `json:"log_level"`
	Address  string `json:"address"`
	Port     uint64 `json:"port"`
}

type batchConfig struct {
	BatchSize          uint8  `json:"batch_size"`
	BatchPeriodSeconds uint64 `json:"batch_period_seconds"`
}

type jsonConfig struct {
	Server                serverConfig     `json:"server"`
	Batch                 batchConfig      `json:"batch"`
	DeprecatedBatchConfig *json.RawMessage `json:"batch_config"`
	Modules               []configModule   `json:"modules"`
}

type appConfig struct {
	// CosmosToEthConfig is the first configured Cosmos→ETH source. Retained for
	// the single-source commands (create-clients*, update-client) and the env
	// override helpers, which operate on one source at a time. It equals
	// CosmosToEthConfigs[0] when at least one source is configured.
	CosmosToEthConfig cosmosToEthConfig
	// CosmosToEthConfigs holds every configured Cosmos→ETH source in file order.
	// `start` runs one independent relay loop per entry, so a second Cosmos
	// source is just another `cosmos_to_eth` module — same circuit, prover and
	// ETH contracts, a different Tendermint RPC and ICS-07 client.
	CosmosToEthConfigs []cosmosToEthConfig
	EthToCosmosConfig  ethToCosmosConfig
	BatchConfig        services.BatchConfig
}

// writeConfigMember rewrites configPath in place, setting
// modules[name=="cosmos_to_eth" && matches sourceClientID].config.<member> =
// value. An empty sourceClientID targets the first cosmos_to_eth module. Other
// fields and existing JSON formatting are preserved outside the value.
func writeConfigMember(configPath, sourceClientID, member, value string) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	out, err := replaceConfigMemberForSource(data, sourceClientID, member, value)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath, out, configFilePerm); err != nil {
		return err
	}
	return os.Chmod(configPath, configFilePerm)
}

// writeICS07Address persists the deployed ICS07 Tendermint light-client address
// (ETH side) back into the sourceClientID module.
func writeICS07Address(configPath, sourceClientID, addr string) error {
	return writeConfigMember(configPath, sourceClientID, "ics07_client", addr)
}

// writeWasmClientID persists the created 08-wasm Ethereum light-client id
// (Cosmos side) back into the sourceClientID module. The id is assigned by
// ibc-go's global client sequence, so it cannot be known until MsgCreateClient
// lands — write it back rather than requiring the operator to predict it.
func writeWasmClientID(configPath, sourceClientID, id string) error {
	return writeConfigMember(configPath, sourceClientID, "cosmos_wasm_client_id", id)
}

// replaceConfigMember targets the first cosmos_to_eth module. Retained for the
// single-source call sites and tests.
func replaceConfigMember(data []byte, member, value string) ([]byte, error) {
	return replaceConfigMemberForSource(data, "", member, value)
}

// moduleSourceClientID resolves the source id of a cosmos_to_eth module: its
// config.ics26_client_id, falling back to the module-level src_chain — mirroring
// how loadConfig defaults ICS26ClientID.
func moduleSourceClientID(data []byte, moduleStart, configStart int) (string, error) {
	if s, e, ok, err := findJSONObjectMember(data, configStart, "ics26_client_id"); err != nil {
		return "", err
	} else if ok {
		var id string
		if err := json.Unmarshal(data[s:e], &id); err != nil {
			return "", fmt.Errorf("parse ics26_client_id: %w", err)
		}
		if id != "" {
			return id, nil
		}
	}
	if s, e, ok, err := findJSONObjectMember(data, moduleStart, "src_chain"); err != nil {
		return "", err
	} else if ok {
		var id string
		if err := json.Unmarshal(data[s:e], &id); err != nil {
			return "", fmt.Errorf("parse src_chain: %w", err)
		}
		return id, nil
	}
	return "", nil
}

// replaceConfigMemberForSource sets config.<member> on the cosmos_to_eth module
// whose source id equals sourceClientID (or the first cosmos_to_eth module when
// sourceClientID is empty), preserving the surrounding JSON formatting.
func replaceConfigMemberForSource(data []byte, sourceClientID, member, value string) ([]byte, error) {
	encodedValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	rootStart := skipJSONSpace(data, 0)
	if rootStart >= len(data) || data[rootStart] != '{' {
		return nil, fmt.Errorf("config root is not an object")
	}
	modulesStart, modulesEnd, ok, err := findJSONObjectMember(data, rootStart, "modules")
	if err != nil {
		return nil, err
	}
	if !ok || modulesStart >= len(data) || data[modulesStart] != '[' {
		return nil, fmt.Errorf("config has no modules array")
	}

	i := skipJSONSpace(data, modulesStart+1)
	for i < modulesEnd {
		if data[i] == ']' {
			break
		}
		if data[i] != '{' {
			return nil, fmt.Errorf("module entry is not an object")
		}

		moduleStart := i
		moduleEnd, err := skipJSONValue(data, moduleStart)
		if err != nil {
			return nil, err
		}

		nameStart, nameEnd, ok, err := findJSONObjectMember(data, moduleStart, "name")
		if err != nil {
			return nil, err
		}
		if ok {
			var name string
			if err := json.Unmarshal(data[nameStart:nameEnd], &name); err != nil {
				return nil, fmt.Errorf("parse module name: %w", err)
			}
			if name == "cosmos_to_eth" {
				configStart, configEnd, ok, err := findJSONObjectMember(data, moduleStart, "config")
				if err != nil {
					return nil, err
				}
				if !ok || configStart >= len(data) || data[configStart] != '{' {
					return nil, fmt.Errorf("cosmos_to_eth.config is not an object")
				}

				matches := sourceClientID == ""
				if !matches {
					id, err := moduleSourceClientID(data, moduleStart, configStart)
					if err != nil {
						return nil, err
					}
					matches = id == sourceClientID
				}
				if matches {
					valueStart, valueEnd, ok, err := findJSONObjectMember(data, configStart, member)
					if err != nil {
						return nil, err
					}
					if ok {
						out := make([]byte, 0, len(data)-valueEnd+valueStart+len(encodedValue))
						out = append(out, data[:valueStart]...)
						out = append(out, encodedValue...)
						out = append(out, data[valueEnd:]...)
						return out, nil
					}
					return insertJSONObjectMember(data, configStart, configEnd, member, encodedValue)
				}
			}
		}

		i = skipJSONSpace(data, moduleEnd)
		if i < len(data) && data[i] == ',' {
			i = skipJSONSpace(data, i+1)
			continue
		}
		if i < len(data) && data[i] == ']' {
			break
		}
	}

	if sourceClientID != "" {
		return nil, fmt.Errorf("cosmos_to_eth source %q not found in config", sourceClientID)
	}
	return nil, fmt.Errorf("module cosmos_to_eth not found in config")
}

func skipJSONSpace(data []byte, i int) int {
	for i < len(data) {
		switch data[i] {
		case ' ', '\n', '\r', '\t':
			i++
		default:
			return i
		}
	}
	return i
}

func scanJSONStringEnd(data []byte, i int) (int, error) {
	if i >= len(data) || data[i] != '"' {
		return 0, fmt.Errorf("expected JSON string")
	}
	escaped := false
	for j := i + 1; j < len(data); j++ {
		if escaped {
			escaped = false
			continue
		}
		switch data[j] {
		case '\\':
			escaped = true
		case '"':
			return j + 1, nil
		}
	}
	return 0, fmt.Errorf("unterminated JSON string")
}

func skipJSONValue(data []byte, i int) (int, error) {
	i = skipJSONSpace(data, i)
	if i >= len(data) {
		return 0, fmt.Errorf("unexpected end of JSON")
	}

	switch data[i] {
	case '"':
		return scanJSONStringEnd(data, i)
	case '{':
		j := skipJSONSpace(data, i+1)
		if j < len(data) && data[j] == '}' {
			return j + 1, nil
		}
		for {
			keyEnd, err := scanJSONStringEnd(data, j)
			if err != nil {
				return 0, err
			}
			j = skipJSONSpace(data, keyEnd)
			if j >= len(data) || data[j] != ':' {
				return 0, fmt.Errorf("expected ':' after object key")
			}
			j, err = skipJSONValue(data, j+1)
			if err != nil {
				return 0, err
			}
			j = skipJSONSpace(data, j)
			if j >= len(data) {
				return 0, fmt.Errorf("unterminated JSON object")
			}
			if data[j] == '}' {
				return j + 1, nil
			}
			if data[j] != ',' {
				return 0, fmt.Errorf("expected ',' or '}' in object")
			}
			j = skipJSONSpace(data, j+1)
		}
	case '[':
		j := skipJSONSpace(data, i+1)
		if j < len(data) && data[j] == ']' {
			return j + 1, nil
		}
		for {
			var err error
			j, err = skipJSONValue(data, j)
			if err != nil {
				return 0, err
			}
			j = skipJSONSpace(data, j)
			if j >= len(data) {
				return 0, fmt.Errorf("unterminated JSON array")
			}
			if data[j] == ']' {
				return j + 1, nil
			}
			if data[j] != ',' {
				return 0, fmt.Errorf("expected ',' or ']' in array")
			}
			j = skipJSONSpace(data, j+1)
		}
	default:
		j := i
		for j < len(data) {
			switch data[j] {
			case ' ', '\n', '\r', '\t', ',', '}', ']':
				return j, nil
			default:
				j++
			}
		}
		return j, nil
	}
}

func findJSONObjectMember(data []byte, objectStart int, key string) (int, int, bool, error) {
	if objectStart >= len(data) || data[objectStart] != '{' {
		return 0, 0, false, fmt.Errorf("expected JSON object")
	}
	i := skipJSONSpace(data, objectStart+1)
	if i < len(data) && data[i] == '}' {
		return 0, 0, false, nil
	}
	for {
		keyStart := i
		keyEnd, err := scanJSONStringEnd(data, keyStart)
		if err != nil {
			return 0, 0, false, err
		}
		var got string
		if err := json.Unmarshal(data[keyStart:keyEnd], &got); err != nil {
			return 0, 0, false, fmt.Errorf("parse object key: %w", err)
		}
		i = skipJSONSpace(data, keyEnd)
		if i >= len(data) || data[i] != ':' {
			return 0, 0, false, fmt.Errorf("expected ':' after object key")
		}
		valueStart := skipJSONSpace(data, i+1)
		valueEnd, err := skipJSONValue(data, valueStart)
		if err != nil {
			return 0, 0, false, err
		}
		if got == key {
			return valueStart, valueEnd, true, nil
		}
		i = skipJSONSpace(data, valueEnd)
		if i >= len(data) {
			return 0, 0, false, fmt.Errorf("unterminated JSON object")
		}
		if data[i] == '}' {
			return 0, 0, false, nil
		}
		if data[i] != ',' {
			return 0, 0, false, fmt.Errorf("expected ',' or '}' in object")
		}
		i = skipJSONSpace(data, i+1)
	}
}

func insertJSONObjectMember(data []byte, objectStart, objectEnd int, key string, encodedValue []byte) ([]byte, error) {
	if objectEnd <= objectStart || objectEnd > len(data) || data[objectEnd-1] != '}' {
		return nil, fmt.Errorf("invalid object range")
	}
	encodedKey, err := json.Marshal(key)
	if err != nil {
		return nil, err
	}
	closeIndex := objectEnd - 1
	empty := skipJSONSpace(data, objectStart+1) == closeIndex

	separator := []byte(",")
	if empty {
		separator = nil
	}
	insert := make([]byte, 0, len(separator)+len(encodedKey)+len(encodedValue)+2)
	insert = append(insert, separator...)
	insert = append(insert, encodedKey...)
	insert = append(insert, ':', ' ')
	insert = append(insert, encodedValue...)

	out := make([]byte, 0, len(data)+len(insert))
	out = append(out, data[:closeIndex]...)
	out = append(out, insert...)
	out = append(out, data[closeIndex:]...)
	return out, nil
}

func loadConfig(configPath string) (*appConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var jc jsonConfig
	if err := json.Unmarshal(data, &jc); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	if jc.DeprecatedBatchConfig != nil {
		return nil, fmt.Errorf("batch_config is deprecated; use batch")
	}

	var c2e cosmosToEthConfig
	var c2eList []cosmosToEthConfig
	var e2c ethToCosmosConfig
	batch := services.DefaultConfig().BatchConfig
	if jc.Batch.BatchSize != 0 {
		batch.BatchSize = jc.Batch.BatchSize
	}
	if jc.Batch.BatchPeriodSeconds != 0 {
		batch.BatchPeriods = time.Duration(jc.Batch.BatchPeriodSeconds) * time.Second
	}
	for _, m := range jc.Modules {
		switch m.Name {
		case "cosmos_to_eth":
			var one cosmosToEthConfig
			if err := json.Unmarshal(m.Config, &one); err != nil {
				return nil, fmt.Errorf("failed to parse cosmos_to_eth config: %w", err)
			}
			if one.ICS26ClientID == "" {
				one.ICS26ClientID = m.SrcChain
			}
			c2eList = append(c2eList, one)
		case "eth_to_cosmos":
			if err := json.Unmarshal(m.Config, &e2c); err != nil {
				return nil, fmt.Errorf("failed to parse eth_to_cosmos config: %w", err)
			}
		}
	}
	if len(c2eList) > 0 {
		c2e = c2eList[0]
	}

	// Validate every cosmos_to_eth source. Distinct sources must not collide on
	// the ICS-26 client id, or their ETH event streams (filtered by that id)
	// would cross-feed.
	seenClientIDs := make(map[string]struct{}, len(c2eList))
	for i := range c2eList {
		if err := validateCosmosToEthConfig(c2eList[i]); err != nil {
			return nil, err
		}
		id := c2eList[i].ICS26ClientID
		if id != "" {
			if _, dup := seenClientIDs[id]; dup {
				return nil, fmt.Errorf("duplicate cosmos_to_eth ics26_client_id %q; each source needs a distinct client id", id)
			}
			seenClientIDs[id] = struct{}{}
		}
	}

	// Validate eth_to_cosmos config if populated. Only eth_beacon_api_url is
	// consumed for the ETH→Cosmos direction.
	if e2c.BeaconUrl != "" {
		if err := validateURL(e2c.BeaconUrl, "eth_to_cosmos.eth_beacon_api_url"); err != nil {
			return nil, err
		}
	}

	return &appConfig{
		CosmosToEthConfig:  c2e,
		CosmosToEthConfigs: c2eList,
		EthToCosmosConfig:  e2c,
		BatchConfig:        batch,
	}, nil
}

// validateCosmosToEthConfig checks the URLs and hex addresses of one Cosmos→ETH
// source. Empty (unpopulated) sources pass so a config with only an
// eth_to_cosmos module still loads.
func validateCosmosToEthConfig(c2e cosmosToEthConfig) error {
	if c2e.TmRpcUrl == "" && c2e.EthRpcUrl == "" && c2e.ICS26Address == "" {
		return nil
	}
	if err := validateURL(c2e.TmRpcUrl, "cosmos_to_eth.tm_rpc_url"); err != nil {
		return err
	}
	if err := validateURL(c2e.EthRpcUrl, "cosmos_to_eth.eth_rpc_url"); err != nil {
		return err
	}
	if c2e.EthWsUrl != "" {
		if err := validateURL(c2e.EthWsUrl, "cosmos_to_eth.eth_ws_url"); err != nil {
			return err
		}
	}
	if err := validateHexAddress(c2e.ICS26Address, "cosmos_to_eth.ics26_address"); err != nil {
		return err
	}
	for _, f := range []struct {
		val, name string
	}{
		{c2e.ICS07Client, "cosmos_to_eth.ics07_client"},
		{c2e.WrapperVerifier, "cosmos_to_eth.wrapper_verifier"},
		{c2e.Membership, "cosmos_to_eth.membership"},
		{c2e.Misbehaviour, "cosmos_to_eth.misbehaviour"},
		{c2e.UpdateClient, "cosmos_to_eth.update_client"},
	} {
		if f.val != "" {
			if err := validateHexAddress(f.val, f.name); err != nil {
				return err
			}
		}
	}
	return nil
}

func validateURL(rawURL, fieldName string) error {
	if rawURL == "" {
		return fmt.Errorf("%s is empty", fieldName)
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL %q for %s: %w", rawURL, fieldName, err)
	}
	if u.Scheme == "" {
		return fmt.Errorf("URL %q for %s is missing scheme", rawURL, fieldName)
	}
	return nil
}

func validateHexAddress(addr, fieldName string) error {
	if addr == "" {
		return fmt.Errorf("%s is empty", fieldName)
	}
	if !common.IsHexAddress(addr) {
		return fmt.Errorf("invalid hex address %q for %s", addr, fieldName)
	}
	return nil
}

func preflightCreateClients(cfg *appConfig) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ethClient, err := ethclient.DialContext(ctx, cfg.CosmosToEthConfig.EthRpcUrl)
	if err != nil {
		return fmt.Errorf("ethereum rpc unavailable at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}
	defer ethClient.Close()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("ethereum rpc not responding at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}

	cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
	if err != nil {
		return fmt.Errorf("failed to create cosmos rpc client for %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}
	status, err := cosmosClient.Status(ctx)
	if err != nil {
		return fmt.Errorf("cosmos rpc unavailable at %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}

	if cfg.EthToCosmosConfig.BeaconUrl != "" {
		if _, err := tendermintClient.GetBeaconGenesis(ctx, cfg.EthToCosmosConfig.BeaconUrl); err != nil {
			return fmt.Errorf("beacon api unavailable at %s: %w", cfg.EthToCosmosConfig.BeaconUrl, err)
		}
	}

	log.Printf("[create-clients] preflight OK: cosmos_height=%d eth_chain_id=%s", status.SyncInfo.LatestBlockHeight, chainID.String())
	return nil
}

func cosmosHasWasmChecksum(cosmosClient *rpchttp.HTTP, checksum string) (bool, error) {
	checksum = strings.ToLower(strings.TrimPrefix(checksum, "0x"))

	reqBytes, err := proto.Marshal(&ibcwasmtypes.QueryChecksumsRequest{})
	if err != nil {
		return false, fmt.Errorf("failed to marshal checksum query: %w", err)
	}

	result, err := cosmosClient.ABCIQuery(context.Background(), "/ibc.lightclients.wasm.v1.Query/Checksums", reqBytes)
	if err != nil {
		return false, fmt.Errorf("failed to query wasm checksums: %w", err)
	}
	if result.Response.Code != 0 {
		return false, fmt.Errorf("checksum query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	var resp ibcwasmtypes.QueryChecksumsResponse
	if err := proto.Unmarshal(result.Response.Value, &resp); err != nil {
		return false, fmt.Errorf("failed to unmarshal checksum query response: %w", err)
	}

	for _, existing := range resp.Checksums {
		if strings.ToLower(existing) == checksum {
			return true, nil
		}
	}

	return false, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func roleManagerOrDefault(cfg *appConfig) string {
	return envOrDefault("ROLE_MANAGER", cfg.CosmosToEthConfig.ICS26Address)
}

func cosmosRouterClientIDOrDefault(cfg *appConfig) string {
	return envOrDefault("ICS26_CLIENT_ID", cfg.CosmosToEthConfig.ICS26ClientID)
}

func validateStartupKeys() error {
	ethPrivKey := os.Getenv("ETH_PRIVATE_KEY")
	if ethPrivKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	if _, err := keys.RestoreKey(ethPrivKey); err != nil {
		return fmt.Errorf("failed to restore ETH private key: %w", err)
	}

	if _, err := (&transaction.Handler{}).CosmosSignerAddress(); err != nil {
		if os.Getenv("COSMOS_PRIVATE_KEY") == "" {
			return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
		}
		return fmt.Errorf("failed to decode COSMOS_PRIVATE_KEY: %w", err)
	}

	return nil
}

func cosmosWasmClientIDOrDefault(cfg *appConfig) string {
	return envOrDefault("COSMOS_WASM_CLIENT_ID", cfg.CosmosToEthConfig.CosmosWasmClientID)
}

// selectSource returns a shallow copy of cfg whose singular CosmosToEthConfig is
// the source identified by clientID, so the single-source create-clients flow
// (preflight, run funcs, config write-back) targets the chosen source. An empty
// clientID selects the sole source, or errors when several are configured.
// Legacy configs without a parsed slice fall back to the singular config.
func selectSource(cfg *appConfig, clientID string) (*appConfig, error) {
	sources := cfg.CosmosToEthConfigs
	if len(sources) == 0 {
		if clientID != "" && clientID != cfg.CosmosToEthConfig.ICS26ClientID {
			return nil, fmt.Errorf("no cosmos_to_eth source with ics26_client_id %q", clientID)
		}
		return cfg, nil
	}
	if clientID == "" {
		if len(sources) == 1 {
			out := *cfg
			out.CosmosToEthConfig = sources[0]
			return &out, nil
		}
		ids := make([]string, len(sources))
		for i := range sources {
			ids[i] = sources[i].ICS26ClientID
		}
		return nil, fmt.Errorf("config has %d cosmos_to_eth sources (%s); pass --source <ics26_client_id> to pick one",
			len(sources), strings.Join(ids, ", "))
	}
	for i := range sources {
		if sources[i].ICS26ClientID == clientID {
			out := *cfg
			out.CosmosToEthConfig = sources[i]
			return &out, nil
		}
	}
	return nil, fmt.Errorf("no cosmos_to_eth source with ics26_client_id %q", clientID)
}

// proofBackendFromFlags resolves the GPU/CPU backend from --gpu-prove or the
// GPU_PROVE env var. Returns (backend, true) when an explicit selection was
// made, otherwise (nil, false) so the caller falls back to env-only defaults.
func proofBackendFromFlags(cmd *cobra.Command) (prover.ProofBackend, bool, error) {
	flagSet := cmd.Flags().Changed(flagGPUProve)
	envSet := prover.GPUProveEnvEnabled()
	if !flagSet && !envSet {
		return nil, false, nil
	}

	useGPU := envSet
	if flagSet {
		v, err := cmd.Flags().GetBool(flagGPUProve)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get gpu prove flag: %w", err)
		}
		useGPU = v
	}

	backend, err := prover.NewProofBackend(useGPU)
	if err != nil {
		return nil, false, err
	}
	return backend, true, nil
}

// --- Main ---

func main() {
	zLogger, _ := zap.NewProduction(zap.AddStacktrace(zap.DPanicLevel))
	defer zLogger.Sync()

	logger := zLogger.Sugar()

	rootCmd := &cobra.Command{
		Use:   "relayer [command]",
		Short: "fast-ibc operator — relay IBC packets between Cosmos and Ethereum",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.AddCommand(
		Start(zLogger),
		CreateClients(zLogger),
		CreateClientsCosmos(zLogger),
		CreateClientsEth(zLogger),
		UpdateClient(zLogger),
		Genesis(zLogger),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

// buildCreateClientsContext dials both chains and returns a services.Context
// wired for the create-clients flow. wasmClientID is the eth-light-client-on-
// Cosmos id used as the on-chain counterparty; pass "" for the Cosmos step,
// which discovers it. The returned cosmosClient has its WebSocket started — the
// caller must Stop it.
func buildCreateClientsContext(logger *zap.Logger, cfg *appConfig, wasmClientID string) (services.Context, *rpchttp.HTTP, error) {
	logger.Sugar().Infof("create-clients: dialing ethereum rpc %s", cfg.CosmosToEthConfig.EthRpcUrl)
	ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
	if err != nil {
		return services.Context{}, nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	logger.Sugar().Infof("create-clients: creating cosmos rpc client %s", cfg.CosmosToEthConfig.TmRpcUrl)
	cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
	if err != nil {
		return services.Context{}, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	cosmosRouterClientID := cosmosRouterClientIDOrDefault(cfg)
	if cosmosRouterClientID == "" {
		return services.Context{}, nil, fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
	}

	ctx := services.NewCtxWithBeacon(
		cosmosClient, ethClient, nil,
		"",
		cfg.EthToCosmosConfig.BeaconUrl,
		wasmClientID,
	)
	ctx.SetCosmosRouterClientID(cosmosRouterClientID)
	ctx.SetAddresses(
		cfg.CosmosToEthConfig.ICS26Address,
		cfg.CosmosToEthConfig.WrapperVerifier,
		cfg.CosmosToEthConfig.Membership,
		cfg.CosmosToEthConfig.Misbehaviour,
		cfg.CosmosToEthConfig.UpdateClient,
		roleManagerOrDefault(cfg),
	)

	if err := cosmosClient.Start(); err != nil {
		return services.Context{}, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
	}
	return ctx, cosmosClient, nil
}

// runCreateClientsCosmos creates the Ethereum light client on Cosmos — the side
// whose client id is auto-assigned by ibc-go's global client sequence —
// registers the counterparty, and persists the resulting cosmos_wasm_client_id
// back into config. Returns the created id. Run this BEFORE the ETH step so the
// real id (not a guessed one) can be wired into the ETH router counterparty.
func runCreateClientsCosmos(logger *zap.Logger, cfg *appConfig, configPath, wasmChecksum string) (string, error) {
	if wasmChecksum == "" {
		wasmChecksum = os.Getenv("WASM_CHECKSUM")
	}
	if wasmChecksum == "" {
		return "", fmt.Errorf("--wasm-checksum (or WASM_CHECKSUM) is required to create the Ethereum light client on Cosmos")
	}
	if cfg.EthToCosmosConfig.BeaconUrl == "" {
		return "", fmt.Errorf("eth_to_cosmos.eth_beacon_api_url is required to create the Ethereum light client on Cosmos")
	}

	ctx, cosmosClient, err := buildCreateClientsContext(logger, cfg, "")
	if err != nil {
		return "", err
	}
	defer cosmosClient.Stop()

	// Validate the checksum is already stored on Cosmos before mutating anything.
	logger.Sugar().Infof("create-clients-cosmos: validating wasm checksum on Cosmos: %s", wasmChecksum)
	ok, err := cosmosHasWasmChecksum(cosmosClient, wasmChecksum)
	if err != nil {
		return "", fmt.Errorf("failed to validate wasm checksum on Cosmos: %w", err)
	}
	if !ok {
		return "", fmt.Errorf("wasm checksum %s has not been previously stored on Cosmos", wasmChecksum)
	}

	if existing := cosmosWasmClientIDOrDefault(cfg); existing != "" {
		logger.Sugar().Warnf("create-clients-cosmos: config already has cosmos_wasm_client_id=%s; a new client will be created and the value overwritten", existing)
	}

	worker := services.NewWorker(&transaction.Handler{}, nil)
	logger.Sugar().Infof("Creating Ethereum light client on Cosmos (checksum=%s)...", wasmChecksum)
	wasmClientID, err := worker.CreateEthClient(ctx, wasmChecksum)
	if err != nil {
		return "", fmt.Errorf("failed to create Ethereum client on Cosmos: %w", err)
	}
	logger.Sugar().Infof("Ethereum light client created on Cosmos: clientID=%s", wasmClientID)

	if err := writeWasmClientID(configPath, cfg.CosmosToEthConfig.ICS26ClientID, wasmClientID); err != nil {
		return "", fmt.Errorf("persist cosmos_wasm_client_id to %s: %w", configPath, err)
	}
	logger.Sugar().Infof("create-clients-cosmos: wrote cosmos_wasm_client_id=%s into %s", wasmClientID, configPath)
	return wasmClientID, nil
}

// runCreateClientsEth deploys the Cosmos (ICS07 Tendermint) light client on
// Ethereum, registers wasmClientID as its counterparty, and persists the
// ics07_client address back into config. wasmClientID must already be known
// (created by the Cosmos step) so the on-chain counterparty is wired to the
// real id. Idempotent: if ics07_client already has deployed code, it is reused.
func runCreateClientsEth(logger *zap.Logger, cfg *appConfig, configPath, wasmClientID, trustLevel string, trustingPeriod uint32) (common.Address, error) {
	if wasmClientID == "" {
		return common.Address{}, fmt.Errorf("cosmos_wasm_client_id is empty; run create-clients-cosmos first")
	}

	ctx, cosmosClient, err := buildCreateClientsContext(logger, cfg, wasmClientID)
	if err != nil {
		return common.Address{}, err
	}
	defer cosmosClient.Stop()

	// Idempotency: skip the deploy if a contract already lives at the configured
	// ics07_client address, so re-running after a partial failure does not
	// redeploy ICS07.
	if cfg.CosmosToEthConfig.ICS07Client != "" {
		addr := common.HexToAddress(cfg.CosmosToEthConfig.ICS07Client)
		if code, err := ctx.EthClient().CodeAt(context.Background(), addr, nil); err == nil && len(code) > 0 {
			logger.Sugar().Infof("create-clients-eth: ics07_client already deployed at %s; skipping deploy", addr.Hex())
			return addr, nil
		}
	}

	if trustingPeriod == 0 {
		unbondingPeriod, err := tendermintClient.GetUnbondingTime(cosmosClient)
		if err != nil {
			return common.Address{}, fmt.Errorf("failed to fetch unbonding time: %w", err)
		}
		trustingPeriod = 2 * uint32(unbondingPeriod) / 3
	}

	proofType := cfg.CosmosToEthConfig.ProofType
	if proofType == "" {
		proofType = "groth16"
	}

	clockDrift := cfg.CosmosToEthConfig.ClockDrift
	if clockDrift == 0 {
		clockDrift = tendermintClient.DefaultClockDrift
	}

	worker := services.NewWorker(&transaction.Handler{}, nil)
	logger.Sugar().Infof(
		"Creating Cosmos light client on Ethereum (trustingPeriod=%d, trustLevel=%s, proofType=%s, clockDrift=%d, counterparty=%s)...",
		trustingPeriod, trustLevel, proofType, clockDrift, wasmClientID,
	)
	ics07Addr, err := worker.CreateCosmosClient(ctx, proofType, trustingPeriod, 0, trustLevel, clockDrift)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create Cosmos client on Ethereum: %w", err)
	}
	if (ics07Addr == common.Address{}) {
		return common.Address{}, fmt.Errorf("ics07 address missing after deploy")
	}

	if err := writeICS07Address(configPath, cfg.CosmosToEthConfig.ICS26ClientID, ics07Addr.Hex()); err != nil {
		return common.Address{}, fmt.Errorf("persist ics07 address to %s: %w", configPath, err)
	}
	logger.Sugar().Infof("create-clients-eth: wrote ics07_client=%s into %s", ics07Addr.Hex(), configPath)
	return ics07Addr, nil
}

// CreateClients runs the full two-chain setup as a one-shot (devnet bring-up):
// create-clients-cosmos first (so the auto-assigned wasm client id is known),
// then create-clients-eth wired to that id. No id guessing, no assertion.
func CreateClients(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-clients",
		Short: "deploy light clients on both chains (runs create-clients-cosmos then create-clients-eth)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}
			_ = godotenv.Load()
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			logger.Sugar().Infof("create-clients: config loaded from %s", configPath)
			source, err := cmd.Flags().GetString(flagSource)
			if err != nil {
				return fmt.Errorf("failed to get source flag: %w", err)
			}
			cfg, err = selectSource(cfg, source)
			if err != nil {
				return err
			}
			logger.Sugar().Infof("create-clients: targeting source %q", cfg.CosmosToEthConfig.ICS26ClientID)
			if err := preflightCreateClients(cfg); err != nil {
				return err
			}
			logger.Sugar().Info("create-clients: preflight passed")

			wasmChecksum, err := cmd.Flags().GetString(flagWasmChecksum)
			if err != nil {
				return fmt.Errorf("failed to get wasm checksum: %w", err)
			}
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level: %w", err)
			}
			trustingPeriod, err := cmd.Flags().GetUint32(flagTrustingPeriod)
			if err != nil {
				return fmt.Errorf("failed to get trusting period: %w", err)
			}

			// 1. Cosmos side first — discovers the real wasm client id.
			wasmClientID, err := runCreateClientsCosmos(logger, cfg, configPath, wasmChecksum)
			if err != nil {
				return err
			}
			// 2. ETH side — wired to the discovered wasm client id.
			if _, err := runCreateClientsEth(logger, cfg, configPath, wasmClientID, trustLevel, trustingPeriod); err != nil {
				return err
			}

			logger.Sugar().Info("=== Setup Complete ===")
			logger.Sugar().Infof("client ids/addresses persisted to %s; ready for 'start'", configPath)
			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagTrustLevel, "2/3", "trust level for Cosmos light client (e.g., 1/3, 2/3)")
	cmd.Flags().Uint32(flagTrustingPeriod, 0, "trusting period in seconds for Cosmos light client (default: 2/3 of chain unbonding period)")
	cmd.Flags().String(flagWasmChecksum, "", "wasm checksum for Ethereum light client (hex)")
	cmd.Flags().String(flagSource, "", "ics26_client_id of the cosmos_to_eth source to target (required when several are configured)")
	return cmd
}

// CreateClientsCosmos creates only the Ethereum light client on Cosmos and
// persists cosmos_wasm_client_id. Run this before create-clients-eth.
func CreateClientsCosmos(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-clients-cosmos",
		Short: "create the Ethereum light client on Cosmos (writes cosmos_wasm_client_id)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}
			_ = godotenv.Load()
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			source, err := cmd.Flags().GetString(flagSource)
			if err != nil {
				return fmt.Errorf("failed to get source flag: %w", err)
			}
			cfg, err = selectSource(cfg, source)
			if err != nil {
				return err
			}
			if err := preflightCreateClients(cfg); err != nil {
				return err
			}
			wasmChecksum, err := cmd.Flags().GetString(flagWasmChecksum)
			if err != nil {
				return fmt.Errorf("failed to get wasm checksum: %w", err)
			}
			_, err = runCreateClientsCosmos(logger, cfg, configPath, wasmChecksum)
			return err
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagWasmChecksum, "", "wasm checksum for Ethereum light client (hex)")
	cmd.Flags().String(flagSource, "", "ics26_client_id of the cosmos_to_eth source to target (required when several are configured)")
	return cmd
}

// CreateClientsEth deploys only the Cosmos light client on Ethereum (ICS07) and
// persists ics07_client. Requires cosmos_wasm_client_id (run create-clients-cosmos first).
func CreateClientsEth(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-clients-eth",
		Short: "deploy the Cosmos light client on Ethereum (requires cosmos_wasm_client_id)",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}
			_ = godotenv.Load()
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			source, err := cmd.Flags().GetString(flagSource)
			if err != nil {
				return fmt.Errorf("failed to get source flag: %w", err)
			}
			cfg, err = selectSource(cfg, source)
			if err != nil {
				return err
			}
			if err := preflightCreateClients(cfg); err != nil {
				return err
			}
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level: %w", err)
			}
			trustingPeriod, err := cmd.Flags().GetUint32(flagTrustingPeriod)
			if err != nil {
				return fmt.Errorf("failed to get trusting period: %w", err)
			}
			wasmClientID := cosmosWasmClientIDOrDefault(cfg)
			if wasmClientID == "" {
				return fmt.Errorf("cosmos_wasm_client_id is empty in config; run create-clients-cosmos first")
			}
			_, err = runCreateClientsEth(logger, cfg, configPath, wasmClientID, trustLevel, trustingPeriod)
			return err
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagTrustLevel, "2/3", "trust level for Cosmos light client (e.g., 1/3, 2/3)")
	cmd.Flags().Uint32(flagTrustingPeriod, 0, "trusting period in seconds for Cosmos light client (default: 2/3 of chain unbonding period)")
	cmd.Flags().String(flagSource, "", "ics26_client_id of the cosmos_to_eth source to target (required when several are configured)")
	return cmd
}

// UpdateClient advances the Cosmos light client on Ethereum once.
func UpdateClient(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client",
		Short: "update the Cosmos light client on Ethereum once",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}

			_ = godotenv.Load()

			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			ethClient, err := ethclient.Dial(cfg.CosmosToEthConfig.EthRpcUrl)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}

			cosmosClient, err := rpchttp.New(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create Cosmos RPC client: %w", err)
			}

			binDir := envOrDefault("PROVER_BIN_DIR", "./bin")
			selectedBackend, hasBackendOverride, err := proofBackendFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("failed to resolve proof backend: %w", err)
			}

			var p *prover.EcipProver
			if hasBackendOverride {
				logger.Sugar().Infof("update-client: overriding proof backend via flags: %s", selectedBackend.Name())
				p, err = prover.NewProverWithBackend(binDir, selectedBackend)
			} else {
				p, err = prover.NewProver(binDir)
			}
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			cosmosWasmClientID := cosmosWasmClientIDOrDefault(cfg)
			if cosmosWasmClientID == "" {
				return fmt.Errorf("cosmos_wasm_client_id is required in cosmos_to_eth config")
			}

			ctx := services.NewCtxWithBeacon(
				cosmosClient, ethClient, nil,
				"",
				cfg.EthToCosmosConfig.BeaconUrl,
				cosmosWasmClientID,
			)
			cosmosRouterClientID := cosmosRouterClientIDOrDefault(cfg)
			if cosmosRouterClientID == "" {
				return fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
			}
			ctx.SetCosmosRouterClientID(cosmosRouterClientID)

			roleManager := roleManagerOrDefault(cfg)
			ctx.SetAddresses(
				cfg.CosmosToEthConfig.ICS26Address,
				cfg.CosmosToEthConfig.WrapperVerifier,
				cfg.CosmosToEthConfig.Membership,
				cfg.CosmosToEthConfig.Misbehaviour,
				cfg.CosmosToEthConfig.UpdateClient,
				roleManager,
			)
			if cfg.CosmosToEthConfig.ICS07Client == "" {
				return fmt.Errorf("ics07_client address is required in cosmos_to_eth config")
			}
			ctx.SetClient(common.HexToAddress(cfg.CosmosToEthConfig.ICS07Client))

			cosmosConfig := services.DefaultConfig()
			if cfg.CosmosToEthConfig.TrustingPeriod != 0 {
				cosmosConfig.TrustingPeriod = cfg.CosmosToEthConfig.TrustingPeriod
			}
			if cfg.CosmosToEthConfig.TrustLevel != "" {
				cosmosConfig.TrustLevel = cfg.CosmosToEthConfig.TrustLevel
			}
			if cfg.CosmosToEthConfig.ProofType != "" {
				cosmosConfig.ProofType = cfg.CosmosToEthConfig.ProofType
			}
			if cfg.CosmosToEthConfig.ClockDrift != 0 {
				cosmosConfig.ClockDrift = cfg.CosmosToEthConfig.ClockDrift
			}
			if cfg.CosmosToEthConfig.BeaconFinalityRetries != 0 {
				cosmosConfig.BeaconFinalityRetries = cfg.CosmosToEthConfig.BeaconFinalityRetries
			}
			if cfg.CosmosToEthConfig.AppHashWaitRetries != 0 {
				cosmosConfig.AppHashWaitRetries = cfg.CosmosToEthConfig.AppHashWaitRetries
			}
			if cfg.CosmosToEthConfig.AppHashWaitInterval != 0 {
				cosmosConfig.AppHashWaitInterval = time.Duration(cfg.CosmosToEthConfig.AppHashWaitInterval) * time.Second
			}
			if cfg.CosmosToEthConfig.FetchTimeout != 0 {
				cosmosConfig.FetchTimeout = time.Duration(cfg.CosmosToEthConfig.FetchTimeout) * time.Second
			}
			if envVal := os.Getenv("FETCH_TIMEOUT"); envVal != "" {
				if d, err := strconv.Atoi(envVal); err == nil && d > 0 {
					cosmosConfig.FetchTimeout = time.Duration(d) * time.Second
				}
			}
			cosmosConfig.BatchConfig = cfg.BatchConfig
			ctx.Config = cosmosConfig

			if err := cosmosClient.Start(); err != nil {
				return fmt.Errorf("failed to start Cosmos WS client: %w", err)
			}
			defer cosmosClient.Stop()

			trustedBlock, err := cmd.Flags().GetInt64(flagTrustedBlock)
			if err != nil {
				return fmt.Errorf("failed to get trusted block: %w", err)
			}

			worker := services.NewWorker(&transaction.Handler{}, p)
			latestBlock, err := worker.UpdateCosmosClient(ctx, cosmosConfig.ProofType, trustedBlock, cosmosConfig.TrustLevel)
			if err != nil {
				return fmt.Errorf("failed to update Cosmos client on Ethereum: %w", err)
			}
			logger.Sugar().Infof("update-client complete: latest_height=%d", latestBlock.BlockHeight)

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().Int64(flagTrustedBlock, 0, "trusted Cosmos block height hint; 0 reads from the on-chain client")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU backend for proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	return cmd
}

// Start starts the relay service loop.
// It loads config from a JSON file, connects to Cosmos and Ethereum,
// and subscribes to send_packet events to relay packets.
func Start(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start the relay loop",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath, err := cmd.Flags().GetString(flagConfigPath)
			if err != nil {
				return fmt.Errorf("failed to get config path: %w", err)
			}

			// Load .env for prover paths, private keys, etc.
			_ = godotenv.Load()

			benchmarkFlag, err := cmd.Flags().GetBool(flagBenchmark)
			if err != nil {
				return fmt.Errorf("failed to get benchmark flag: %w", err)
			}
			utils.SetBenchEnabled(benchmarkFlag || utils.BenchEnabled())
			if utils.BenchEnabled() {
				log.Printf("[benchmark] enabled: detailed gas/timing logs are active")
			}

			if err := validateStartupKeys(); err != nil {
				return err
			}

			// Load JSON config
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Load the prover once and share it across every source loop — the
			// bucket registry (r1cs/pk/vk) is read-only after load, so concurrent
			// GenerateProof calls are safe. On a GPU backend the calls serialize
			// on the device; correctness is unaffected.
			binDir := envOrDefault("PROVER_BIN_DIR", "./bin")
			selectedBackend, hasBackendOverride, err := proofBackendFromFlags(cmd)
			if err != nil {
				return fmt.Errorf("failed to resolve proof backend: %w", err)
			}

			var p *prover.EcipProver
			if hasBackendOverride {
				logger.Sugar().Infof("start: overriding proof backend via flags: %s", selectedBackend.Name())
				p, err = prover.NewProverWithBackend(binDir, selectedBackend)
			} else {
				p, err = prover.NewProver(binDir)
			}
			if err != nil {
				return fmt.Errorf("failed to load prover: %w", err)
			}

			sources := cfg.CosmosToEthConfigs
			if len(sources) == 0 {
				return fmt.Errorf("no cosmos_to_eth source configured in %s", configPath)
			}
			// Env overrides (ICS26_CLIENT_ID, COSMOS_WASM_CLIENT_ID, ROLE_MANAGER)
			// name a single source; only honor them when exactly one is
			// configured, otherwise they would wrongly apply to every source.
			allowEnvOverride := len(sources) == 1

			// One shared TransactionHandler across all sources: they submit from
			// the same ETH signer, so its per-account nonce cache + mutex must be
			// shared to serialize nonce allocation (per-source handlers would
			// collide on nonce). The Cosmos side re-queries the account sequence
			// fresh under cosmosMu, so sharing is safe there too.
			txHandler := &transaction.Handler{}

			// One independent relay loop per Cosmos→ETH source. Each has its own
			// Tendermint RPC, ICS-07 client and router client id; they share the
			// prover, the TransactionHandler, the ETH beacon endpoint, and the
			// same ICS26Router (ETH events are partitioned by the per-source
			// router client id filter).
			var wg sync.WaitGroup
			for i := range sources {
				svc, srcCtx, cleanup, err := startCosmosToEthSource(
					logger, sources[i], cfg.EthToCosmosConfig, cfg.BatchConfig, p, txHandler, allowEnvOverride,
				)
				if err != nil {
					return fmt.Errorf("cosmos_to_eth source %q: %w", sources[i].ICS26ClientID, err)
				}
				wg.Add(1)
				go func(svc *services.Services, srcCtx services.Context, cleanup func()) {
					defer wg.Done()
					defer cleanup()
					svc.StartLoop(srcCtx)
				}(svc, srcCtx, cleanup)
			}
			logger.Sugar().Infof("Relayer started: relaying %d Cosmos→ETH source(s)", len(sources))
			wg.Wait()

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU backend for proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	cmd.Flags().Bool(flagBenchmark, false, "enable detailed benchmark gas/timing logs (or set RELAYER_BENCHMARK=1)")
	return cmd
}

// Genesis generates the genesis state for a new client.
func Genesis(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "genesis",
		Short: "genesis",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := godotenv.Load()
			if err != nil {
				return fmt.Errorf("error loading .env file: %v", err)
			}
			tendermintRpcEndpoint := os.Getenv("TENDERMINT_RPC_URL")
			if tendermintRpcEndpoint == "" {
				return fmt.Errorf("TENDERMINT_RPC_URL environment variable is required in .env file")
			}
			tendermintRpcClient, err := rpchttp.New(tendermintRpcEndpoint, "/websocket")
			if err != nil {
				return fmt.Errorf("failed to create RPC client: %w", err)
			}
			trustedBlock, err := cmd.Flags().GetInt64(flagTrustedBlock)
			if err != nil {
				return fmt.Errorf("failed to get trusted block: %w", err)
			}
			trustingPeriod, err := cmd.Flags().GetUint32(flagTrustingPeriod)
			if err != nil {
				return fmt.Errorf("failed to get trusting period: %w", err)
			}
			trustLevel, err := cmd.Flags().GetString(flagTrustLevel)
			if err != nil {
				return fmt.Errorf("failed to get trust level from flag: %w", err)
			}
			proofType, err := cmd.Flags().GetString(flagProofType)
			if err != nil {
				return fmt.Errorf("failed to get proof type from flag: %w", err)
			}
			clockDrift, err := cmd.Flags().GetUint32(flagClockDrift)
			if err != nil {
				return fmt.Errorf("failed to get clock drift from flag: %w", err)
			}

			genesis, err := tendermintClient.GetGenesis(tendermintRpcClient, trustedBlock, trustingPeriod, trustLevel, proofType, clockDrift)
			if err != nil {
				return fmt.Errorf("failed to get genesis: %w", err)
			}

			outputType, err := cmd.Flags().GetString(flagOutput)
			if err != nil {
				return fmt.Errorf("failed to get output type from flag: %w", err)
			}

			data, err := json.Marshal(genesis)
			if err != nil {
				return fmt.Errorf("failed to marshal genesis state: %w", err)
			}

			switch outputType {
			case "json":
				fmt.Println(string(data))
			case "file":
				outputDir, err := cmd.Flags().GetString(flagOutputPath)
				if err != nil {
					return fmt.Errorf("failed to get output path from flag: %w", err)
				}
				if err := os.WriteFile(outputDir, data, 0o600); err != nil {
					return fmt.Errorf("failed to write genesis state to file: %w", err)
				}
			default:
				return fmt.Errorf("unsupported output type: %s, supported types are: json, file", outputType)
			}

			return nil
		},
	}
	cmd.Flags().String(flagProofType, "groth16", "the type of proof to use (groth16, plonk)")
	cmd.Flags().Int64(flagTrustedBlock, 0, "the trusted block height, if <height> is 0 then catch latest block")
	cmd.Flags().String(flagOutput, "json", "the output structure for the genesis state (json, file)")
	cmd.Flags().String(flagOutputPath, "./data/genesis.json", "the path to the output file for the genesis state")
	cmd.Flags().String(flagTrustLevel, "2/3", "the trust level for the genesis state (e.g., 2/3)")
	cmd.Flags().Uint32(flagTrustingPeriod, 0, "the trusting period for the genesis state")
	cmd.Flags().Uint32(flagClockDrift, tendermintClient.DefaultClockDrift, "allowed clock drift in seconds between proven consensus state and verifying chain time")
	return cmd
}
