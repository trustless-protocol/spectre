package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	proto "github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"relayer/chain"
	tendermintClient "relayer/client"
	"relayer/prover"
	"relayer/relay"
	"relayer/services"
	"relayer/transaction"
	utils "relayer/utils"
)

func isShutdownErr(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

func stopRelaysAndCleanup(cancel func(), wait func(), cleanups []func()) {
	cancel()
	wait()
	for _, cleanup := range cleanups {
		cleanup()
	}
}

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
	flagForceRotation  = "force-rotation"
	flagTargetHeight   = "target-height"
	flagMaxHops        = "max-hops"
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
	SpectreClient         string `json:"spectre_client"`
	SignatureVerifier     string `json:"signature_verifier"`
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
	RotationThreshold     string `json:"rotation_threshold"`
	RefreshInterval       uint64 `json:"refresh_interval_seconds"`
}

type ethToCosmosConfig struct {
	// BeaconUrl is the only field consumed for the ETH→Cosmos direction; the
	// Cosmos RPC, Eth RPC and ICS26 address all come from the cosmos_to_eth
	// module (used for both directions).
	BeaconUrl string `json:"eth_beacon_api_url"`
}

type configModule struct {
	// Name is a free-form human label (logs / diagnostics). For legacy configs it
	// doubles as the direction selector when dst_chain is absent (back-compat).
	Name string `json:"name"`
	// SrcChain / DstChain are chain kinds (chain.ChainType: "cosmos", "ethereum",
	// "opstack", "arbitrum") in the canonical schema; the (src, dst) pair selects
	// the relay direction. In a legacy config (no dst_chain) src_chain instead
	// carries the ics26_client_id fallback.
	SrcChain string          `json:"src_chain"`
	DstChain string          `json:"dst_chain"`
	Builder  string          `json:"builder"`
	Config   json.RawMessage `json:"config"`
}

// moduleDirection is the relay direction a config module resolves to.
type moduleDirection string

const (
	dirCosmosToEth moduleDirection = "cosmos_to_eth"
	dirEthToCosmos moduleDirection = "eth_to_cosmos"
	dirL2ToCosmos  moduleDirection = "l2_to_cosmos"
	dirCosmosToL2  moduleDirection = "cosmos_to_l2"
)

// isChainKind reports whether s names a known chain family (chain.ChainType).
func isChainKind(s string) bool {
	switch chain.ChainType(s) {
	case chain.Cosmos, chain.Ethereum, chain.OPStack, chain.Arbitrum:
		return true
	default:
		return false
	}
}

// classifyModule resolves a module's relay direction. A module is CANONICAL when
// BOTH src_chain and dst_chain are known chain kinds — the (src, dst) pair selects
// the direction through the table below, so a new direction is added by extending
// the table + its adapter, not by editing a name switch. Otherwise the module is
// LEGACY and dispatches on `name`; this is deliberate back-compat, because an old
// config may carry a junk dst_chain (a since-removed field) or a src_chain that is
// really the ics26_client_id — neither is a valid chain kind, so such a module
// correctly stays on the legacy path. isLegacy keeps the src_chain ->
// ics26_client_id fallback that the canonical schema no longer overloads. An
// unrecognized module is a hard error — never silently dropped.
func classifyModule(m configModule) (dir moduleDirection, isLegacy bool, err error) {
	if isChainKind(m.SrcChain) && isChainKind(m.DstChain) {
		src := chain.ChainType(m.SrcChain)
		dst := chain.ChainType(m.DstChain)
		switch {
		case src == chain.Cosmos && dst == chain.Ethereum:
			return dirCosmosToEth, false, nil
		case src == chain.Ethereum && dst == chain.Cosmos:
			return dirEthToCosmos, false, nil
		case (src == chain.OPStack || src == chain.Arbitrum) && dst == chain.Cosmos:
			return dirL2ToCosmos, false, nil
		case src == chain.Cosmos && (dst == chain.OPStack || dst == chain.Arbitrum):
			return dirCosmosToL2, false, nil
		default:
			return "", false, fmt.Errorf("module %q: unsupported direction src_chain=%q dst_chain=%q", m.Name, m.SrcChain, m.DstChain)
		}
	}
	switch m.Name {
	case string(dirCosmosToEth):
		return dirCosmosToEth, true, nil
	case string(dirEthToCosmos):
		return dirEthToCosmos, true, nil
	default:
		return "", false, fmt.Errorf(
			"module %q: unrecognized module; use a legacy name (cosmos_to_eth / eth_to_cosmos) or declare src_chain + dst_chain as chain kinds (cosmos, ethereum, opstack, arbitrum)", m.Name)
	}
}

type batchConfig struct {
	BatchSize          uint8  `json:"batch_size"`
	BatchPeriodSeconds uint64 `json:"batch_period_seconds"`
}

// jsonConfig deliberately has no `server` member. One existed — log_level, address,
// port — and nothing ever read it: the relayer opens no listener at all, so there
// was no HTTP endpoint, no configurable bind address, and no way for log_level to
// affect anything. Configs that still carry the block keep loading, because
// encoding/json ignores unknown members and nothing here sets
// DisallowUnknownFields (pinned by TestLoadConfig_IgnoresLegacyServerBlock).
type jsonConfig struct {
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
	// CosmosToL2Configs holds every configured Cosmos→L2 destination in file order.
	// Same schema as cosmos_to_eth (the groth16 path covers Cosmos → any EVM), but
	// pointed at an L2 rollup's exec RPC + its deployed SpectreClient/ICS26Router;
	// `start` runs one independent cosmos→l2 relay loop per entry.
	CosmosToL2Configs []cosmosToEthConfig
	// CosmosToL2Names mirrors CosmosToL2Configs and carries the config module labels
	// used in startup diagnostics.
	CosmosToL2Names []string
	// HasEthToCosmos records whether the direction was explicitly configured.
	// It distinguishes an absent module from a configured module with an empty
	// beacon endpoint, which must fail validation instead of silently disabling
	// the ETH→Cosmos relay.
	HasEthToCosmos    bool
	EthToCosmosConfig ethToCosmosConfig
	// L2ToCosmosConfigs holds every configured L2 (opstack/arbitrum) → Cosmos source
	// in file order; `start` runs one independent relay module per entry.
	L2ToCosmosConfigs []l2ToCosmosConfig
	BatchConfig       services.BatchConfig
}

// writeConfigMember rewrites configPath in place, setting
// modules[name=="cosmos_to_eth" && matches sourceClientID].config.<member> =
// value. An empty sourceClientID targets the first cosmos_to_eth module. Other
// fields and existing JSON formatting are preserved outside the value.
func writeConfigMember(configPath, sourceClientID, member, value string) error {
	return writeConfigMemberIn(configPath, sourceClientID, member, value, dirCosmosToEth, dirCosmosToL2)
}

// writeConfigMemberIn is writeConfigMember restricted to the given module directions.
// Callers that know which side they are writing MUST name it: a cosmos_to_eth and a
// cosmos_to_l2 module can share an ics26_client_id, so an unrestricted search resolves
// the ambiguity by file order rather than by intent.
func writeConfigMemberIn(configPath, sourceClientID, member, value string, dirs ...moduleDirection) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	out, err := replaceConfigMemberForSourceIn(data, sourceClientID, member, value, dirs)
	if err != nil {
		return err
	}
	if err := os.WriteFile(configPath, out, configFilePerm); err != nil {
		return err
	}
	return os.Chmod(configPath, configFilePerm)
}

// assertConfigMemberWritable reports whether writeConfigMemberIn would find a
// destination for member, without touching the file.
//
// It resolves through the same replaceConfigMemberForSourceIn the real write uses
// rather than re-implementing the lookup, so the check cannot drift from the write
// it is guarding. A separately-written precheck that disagreed would be worse than
// none: it would pass, the write would still fail, and the failure would once again
// land after the irreversible step.
func assertConfigMemberWritable(configPath, sourceClientID, member string, dirs ...moduleDirection) error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	_, err = replaceConfigMemberForSourceIn(data, sourceClientID, member, "", dirs)
	return err
}

// unrecordedClientError explains a client that exists on chain but is not in the
// config, and names the one recovery that does not create a second one.
//
// Re-running create-clients-cosmos is the natural reaction to a non-zero exit and
// the wrong one: it creates another client, pays for it again, and orphans this one
// with nothing advancing it until it expires (#310).
func unrecordedClientError(configPath, wasmClientID string, cause error) error {
	return fmt.Errorf(
		"created Ethereum light client %s on Cosmos but could not record it in %s: %w\n"+
			"  The client exists and is usable. Set cosmos_wasm_client_id=%s on the cosmos_to_eth "+
			"module by hand, then run create-clients-eth.\n"+
			"  Do NOT re-run create-clients-cosmos: it would create a second client and leave %s orphaned",
		wasmClientID, configPath, cause, wasmClientID, wasmClientID)
}

// writeSpectreClientAddress persists the deployed SpectreClient light-client
// address (ETH side) back into the sourceClientID module.
func writeSpectreClientAddress(configPath, sourceClientID, addr string) error {
	return writeConfigMember(configPath, sourceClientID, "spectre_client", addr)
}

// writeEthWasmClientID persists the created 08-wasm Ethereum light-client id into
// the cosmos_to_eth module for sourceClientID. The id is assigned by ibc-go's global
// client sequence, so it cannot be known until MsgCreateClient lands — write it back
// rather than requiring the operator to predict it.
//
// The direction is pinned deliberately. A cosmos_to_eth and a cosmos_to_l2 module may
// legitimately carry the SAME ics26_client_id (config.example.json ships exactly that:
// both say "cosmoshub-1"), so a direction-agnostic search matches whichever appears
// first in the file. That silently wrote the L2 client id into the Ethereum module and
// left the L2 module untouched, breaking the counterparty wiring create-clients-eth
// reads next.
func writeEthWasmClientID(configPath, sourceClientID, id string) error {
	return writeConfigMemberIn(configPath, sourceClientID, "cosmos_wasm_client_id", id, dirCosmosToEth)
}

// writeL2WasmClientID persists the created L2 wasm client id into the cosmos_to_l2
// module for sourceClientID. See writeEthWasmClientID for why the direction is pinned.
func writeL2WasmClientID(configPath, sourceClientID, id string) error {
	return writeConfigMemberIn(configPath, sourceClientID, "cosmos_wasm_client_id", id, dirCosmosToL2)
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
	return replaceConfigMemberForSourceIn(data, sourceClientID, member, value,
		[]moduleDirection{dirCosmosToEth, dirCosmosToL2})
}

// replaceConfigMemberForSourceIn is replaceConfigMemberForSource restricted to dirs.
func replaceConfigMemberForSourceIn(data []byte, sourceClientID, member, value string, dirs []moduleDirection) ([]byte, error) {
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

		// Classify this module by the SAME rule loadConfig/classifyModule uses
		// (canonical src_chain/dst_chain first, legacy name fallback), NOT by a raw
		// name == "cosmos_to_eth" check. A canonical config's name is a free-form
		// label (e.g. "cosmos-to-eth"), so a name check would miss it and write-back
		// would fail with "not found" after create-clients mutated chain state.
		name, err := moduleStringMember(data, moduleStart, "name")
		if err != nil {
			return nil, err
		}
		srcChain, err := moduleStringMember(data, moduleStart, "src_chain")
		if err != nil {
			return nil, err
		}
		dstChain, err := moduleStringMember(data, moduleStart, "dst_chain")
		if err != nil {
			return nil, err
		}
		dir, _, cerr := classifyModule(configModule{Name: name, SrcChain: srcChain, DstChain: dstChain})
		// create-clients-eth writes spectre_client back into either a cosmos_to_eth
		// source or a cosmos_to_l2 destination (both hold a SpectreClient on the EVM
		// side, created by the same flow), so accept both directions here.
		accepted := false
		for _, d := range dirs {
			if dir == d {
				accepted = true
				break
			}
		}
		if cerr == nil && accepted {
			configStart, configEnd, ok, err := findJSONObjectMember(data, moduleStart, "config")
			if err != nil {
				return nil, err
			}
			if !ok || configStart >= len(data) || data[configStart] != '{' {
				return nil, fmt.Errorf("%s.config is not an object", dir)
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

// moduleStringMember reads an optional string member of the raw JSON module object
// at moduleStart, returning "" when the member is absent (so classification can run
// on whatever subset of name/src_chain/dst_chain the config carries).
func moduleStringMember(data []byte, moduleStart int, key string) (string, error) {
	start, end, ok, err := findJSONObjectMember(data, moduleStart, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	var s string
	if err := json.Unmarshal(data[start:end], &s); err != nil {
		return "", fmt.Errorf("parse module %s: %w", key, err)
	}
	return s, nil
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

// loadConfig reads a config that is expected to be complete — every client already
// created and recorded. This is what `start`, `update-client` and `create-clients-eth`
// want.
func loadConfig(configPath string) (*appConfig, error) {
	return loadConfigWith(configPath, true)
}

// loadConfigForClientCreation reads a config that is about to have its L2 wasm
// clients created, so their ids are legitimately still empty (#309). Only
// create-clients-cosmos uses this: it is the command that fills them in.
func loadConfigForClientCreation(configPath string) (*appConfig, error) {
	return loadConfigWith(configPath, false)
}

func loadConfigWith(configPath string, requireL2WasmClientID bool) (*appConfig, error) {
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
	var c2eNames []string
	var c2l2List []cosmosToEthConfig
	var c2l2Names []string
	var l2List []l2ToCosmosConfig
	var e2c ethToCosmosConfig
	var hasEthToCosmos bool
	batch := services.DefaultConfig().BatchConfig
	if jc.Batch.BatchSize != 0 {
		batch.BatchSize = jc.Batch.BatchSize
	}
	if jc.Batch.BatchPeriodSeconds != 0 {
		batch.BatchPeriods = time.Duration(jc.Batch.BatchPeriodSeconds) * time.Second
	}
	for _, m := range jc.Modules {
		dir, isLegacy, err := classifyModule(m)
		if err != nil {
			return nil, err
		}
		switch dir {
		case dirCosmosToEth:
			var one cosmosToEthConfig
			if err := json.Unmarshal(m.Config, &one); err != nil {
				return nil, fmt.Errorf("module %q: parse cosmos_to_eth config: %w", m.Name, err)
			}
			// Legacy configs overloaded src_chain as the ics26_client_id fallback;
			// the canonical schema keeps src_chain as the chain kind, so only fall
			// back for a legacy module.
			if one.ICS26ClientID == "" && isLegacy {
				one.ICS26ClientID = m.SrcChain
			}
			c2eList = append(c2eList, one)
			c2eNames = append(c2eNames, m.Name)
		case dirEthToCosmos:
			hasEthToCosmos = true
			if err := json.Unmarshal(m.Config, &e2c); err != nil {
				return nil, fmt.Errorf("module %q: parse eth_to_cosmos config: %w", m.Name, err)
			}
		case dirCosmosToL2:
			// Same schema as cosmos_to_eth, pointed at the L2 (groth16 → any EVM).
			var one cosmosToEthConfig
			if err := json.Unmarshal(m.Config, &one); err != nil {
				return nil, fmt.Errorf("module %q: parse cosmos_to_l2 config: %w", m.Name, err)
			}
			c2l2List = append(c2l2List, one)
			c2l2Names = append(c2l2Names, m.Name)
		case dirL2ToCosmos:
			var one l2ToCosmosConfig
			if err := json.Unmarshal(m.Config, &one); err != nil {
				return nil, fmt.Errorf("module %q: parse l2_to_cosmos config: %w", m.Name, err)
			}
			one.kind = chain.ChainType(m.SrcChain) // opstack | arbitrum
			if err := one.validateWith(requireL2WasmClientID); err != nil {
				return nil, fmt.Errorf("module %q: %w", m.Name, err)
			}
			l2List = append(l2List, one)
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
		if err := validateCosmosToEthConfig(c2eList[i], moduleFieldPrefix(c2eNames[i], dirCosmosToEth)); err != nil {
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
	//
	// eth_ws_url is deliberately NOT required here — see validateRelayStartupConfig,
	// which `start` calls. loadConfig is shared with create-clients-cosmos,
	// create-clients-eth and update-client, none of which ever read EthWsUrl
	// (build_source.go is the only consumer, and it runs from `start` alone), so
	// rejecting a missing websocket at load time blocks client creation on a
	// field that command will not use.
	if hasEthToCosmos {
		if e2c.BeaconUrl == "" {
			return nil, fmt.Errorf("eth_to_cosmos is configured but eth_to_cosmos.eth_beacon_api_url is empty; the ETH→Cosmos direction reads finality and light-client updates from the beacon REST API and cannot run without it")
		}
		if err := validateURL(e2c.BeaconUrl, "eth_to_cosmos.eth_beacon_api_url"); err != nil {
			return nil, err
		}
	}

	// Validate every cosmos_to_l2 destination (same schema/checks as cosmos_to_eth);
	// distinct destinations must not collide on the ICS-26 client id.
	seenL2ClientIDs := make(map[string]struct{}, len(c2l2List))
	for i := range c2l2List {
		if err := validateCosmosToEthConfig(c2l2List[i], moduleFieldPrefix(c2l2Names[i], dirCosmosToL2)); err != nil {
			return nil, err
		}
		id := c2l2List[i].ICS26ClientID
		if id != "" {
			if _, dup := seenL2ClientIDs[id]; dup {
				return nil, fmt.Errorf("duplicate cosmos_to_l2 ics26_client_id %q; each destination needs a distinct client id", id)
			}
			seenL2ClientIDs[id] = struct{}{}
		}
	}

	return &appConfig{
		CosmosToEthConfig:  c2e,
		CosmosToEthConfigs: c2eList,
		CosmosToL2Configs:  c2l2List,
		CosmosToL2Names:    c2l2Names,
		HasEthToCosmos:     hasEthToCosmos,
		EthToCosmosConfig:  e2c,
		L2ToCosmosConfigs:  l2List,
		BatchConfig:        batch,
	}, nil
}

// validateRelayStartupConfig rejects a config that would start the relay loop
// half-dead. It is called by `start` only — NOT by loadConfig, which the
// client-creation commands share.
//
// Configuring eth_to_cosmos states the intent to relay that direction, and it
// cannot run without a websocket: SubscribeEth needs eth_subscribe, which HTTP
// cannot serve. Without this check the relayer starts, logs a single
// "eth websocket URL is not configured" line while the prover is still loading,
// and then runs half-dead — the forward direction works, nothing ever picks up an
// acknowledgement, and no later line says why. Fail before the prover load instead.
//
// Every source is checked, not just the first: runAdapterEngine spawns an
// eth->cosmos leg per cosmos_to_eth source, and each leg subscribes with that
// source's own eth_ws_url (build_source.go passes c2e.EthWsUrl into its context).
// Validating only CosmosToEthConfigs[0] left every later source free to fail
// exactly this way.
func validateRelayStartupConfig(cfg *appConfig) error {
	if cfg == nil || !cfg.HasEthToCosmos {
		return nil
	}
	for i := range cfg.CosmosToEthConfigs {
		if cfg.CosmosToEthConfigs[i].EthWsUrl != "" {
			continue
		}
		return fmt.Errorf(
			"eth_to_cosmos is configured but cosmos_to_eth[%d] (ics26_client_id %q) has an empty eth_ws_url; "+
				"the ETH→Cosmos direction subscribes to ICS26Router events over eth_subscribe "+
				"and cannot run without a ws:// or wss:// endpoint",
			i, cfg.CosmosToEthConfigs[i].ICS26ClientID)
	}
	return nil
}

// resolveBeaconURL returns the L1 beacon endpoint to create/advance the Ethereum
// light client with. Only the eth_to_cosmos module declares one now: L2 sources no
// longer pin an Ethereum client, so an L2-only deployment needs no beacon at all.
func resolveBeaconURL(cfg *appConfig) string {
	if cfg.EthToCosmosConfig.BeaconUrl != "" {
		return cfg.EthToCosmosConfig.BeaconUrl
	}
	return ""
}

func moduleFieldPrefix(moduleName string, fallback moduleDirection) string {
	if name := strings.TrimSpace(moduleName); name != "" {
		return name
	}
	return string(fallback)
}

// validateCosmosToEthConfig checks the URLs and hex addresses of one Cosmos→ETH
// source or Cosmos→L2 destination. Empty (unpopulated) sources pass so a config
// with only an eth_to_cosmos module still loads.
func validateCosmosToEthConfig(c2e cosmosToEthConfig, fieldPrefix string) error {
	if c2e.TmRpcUrl == "" && c2e.EthRpcUrl == "" && c2e.ICS26Address == "" {
		return nil
	}
	if err := validateURL(c2e.TmRpcUrl, fieldPrefix+".tm_rpc_url"); err != nil {
		return err
	}
	if err := validateURL(c2e.EthRpcUrl, fieldPrefix+".eth_rpc_url"); err != nil {
		return err
	}
	if c2e.EthWsUrl != "" {
		if err := validateURL(c2e.EthWsUrl, fieldPrefix+".eth_ws_url"); err != nil {
			return err
		}
	}
	if err := validateHexAddress(c2e.ICS26Address, fieldPrefix+".ics26_address"); err != nil {
		return err
	}
	for _, f := range []struct {
		val, name string
	}{
		{c2e.SpectreClient, fieldPrefix + ".spectre_client"},
		{c2e.SignatureVerifier, fieldPrefix + ".signature_verifier"},
		{c2e.Membership, fieldPrefix + ".membership"},
		{c2e.Misbehaviour, fieldPrefix + ".misbehaviour"},
		{c2e.UpdateClient, fieldPrefix + ".update_client"},
	} {
		if f.val != "" {
			if err := validateHexAddress(f.val, f.name); err != nil {
				return err
			}
		}
	}
	if c2e.RotationThreshold != "" {
		if _, err := services.ParseRotationThreshold(c2e.RotationThreshold); err != nil {
			return fmt.Errorf("%s.rotation_threshold: %w", fieldPrefix, err)
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

	ethClient, err := tendermintClient.DialEthRPC(ctx, cfg.CosmosToEthConfig.EthRpcUrl, tendermintClient.DefaultRPCTimeout)
	if err != nil {
		return fmt.Errorf("ethereum rpc unavailable at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}
	defer ethClient.Close()

	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("ethereum rpc not responding at %s: %w", cfg.CosmosToEthConfig.EthRpcUrl, err)
	}

	cosmosClient, err := tendermintClient.DialCosmosRPC(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket", tendermintClient.DefaultRPCTimeout)
	if err != nil {
		return fmt.Errorf("failed to create cosmos rpc client for %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}
	status, err := cosmosClient.Status(ctx)
	if err != nil {
		return fmt.Errorf("cosmos rpc unavailable at %s: %w", cfg.CosmosToEthConfig.TmRpcUrl, err)
	}

	if beaconURL := resolveBeaconURL(cfg); beaconURL != "" {
		if _, err := tendermintClient.GetBeaconGenesis(ctx, beaconURL); err != nil {
			return fmt.Errorf("beacon api unavailable at %s: %w", beaconURL, err)
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

// validateStartupKeys uses the transaction signer seam so startup cannot retain
// a separate environment-only signing-key path.
func validateStartupKeys() error {
	return (&transaction.Handler{}).ValidateKeys()
}

func cosmosWasmClientIDOrDefault(cfg *appConfig) string {
	return envOrDefault("COSMOS_WASM_CLIENT_ID", cfg.CosmosToEthConfig.CosmosWasmClientID)
}

// selectSource returns a shallow copy of cfg whose singular CosmosToEthConfig is
// the source identified by clientID, so the single-source create-clients flow
// (preflight, run funcs, config write-back) targets the chosen source. An empty
// clientID selects the sole source, or errors when several are configured.
// Legacy configs without a parsed slice fall back to the singular config.
// selectSource picks the create-clients target by ics26_client_id. It searches both
// cosmos_to_eth sources and cosmos_to_l2 destinations: the L2 destination's ETH-side
// (its SpectreClient + ICS26Router on the L2) is created by the same create-clients-eth
// flow — a cosmos_to_l2 entry is the same config struct pointed at the L2 exec RPC, so
// the selected entry is placed in CosmosToEthConfig and the downstream deploy/write-back
// treat it uniformly.
func selectSource(cfg *appConfig, clientID string) (*appConfig, error) {
	eth := cfg.CosmosToEthConfigs
	l2 := cfg.CosmosToL2Configs
	total := len(eth) + len(l2)

	// Legacy single-config path (no modules parsed into the lists).
	if total == 0 {
		if clientID != "" && clientID != cfg.CosmosToEthConfig.ICS26ClientID {
			return nil, fmt.Errorf("no cosmos_to_eth or cosmos_to_l2 source with ics26_client_id %q", clientID)
		}
		return cfg, nil
	}

	pick := func(c cosmosToEthConfig) *appConfig {
		out := *cfg
		out.CosmosToEthConfig = c
		return &out
	}

	if clientID == "" {
		if total == 1 {
			if len(eth) == 1 {
				return pick(eth[0]), nil
			}
			return pick(l2[0]), nil
		}
		ids := make([]string, 0, total)
		for i := range eth {
			ids = append(ids, eth[i].ICS26ClientID)
		}
		for i := range l2 {
			ids = append(ids, l2[i].ICS26ClientID)
		}
		return nil, fmt.Errorf("config has %d cosmos_to_eth/cosmos_to_l2 sources (%s); pass --source <ics26_client_id> to pick one",
			total, strings.Join(ids, ", "))
	}
	for i := range eth {
		if eth[i].ICS26ClientID == clientID {
			return pick(eth[i]), nil
		}
	}
	for i := range l2 {
		if l2[i].ICS26ClientID == clientID {
			return pick(l2[i]), nil
		}
	}
	return nil, fmt.Errorf("no cosmos_to_eth or cosmos_to_l2 source with ics26_client_id %q", clientID)
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
		CreateClientsCosmos(zLogger),
		CreateClientsEth(zLogger),
		UpdateClient(zLogger),
		SubmitMisbehaviour(zLogger),
		Genesis(zLogger),
	)

	if err := rootCmd.Execute(); err != nil {
		logger.Fatal(err)
	}
}

// buildCreateClientsDeps dials both chains and returns scoped relay dependencies
// wired for the create-clients flow. wasmClientID is the eth-light-client-on-
// Cosmos id used as the on-chain counterparty; pass "" for the Cosmos step,
// which discovers it. The returned cosmosClient has its WebSocket started — the
// caller must Stop it.
func buildCreateClientsDeps(logger *zap.Logger, cfg *appConfig, wasmClientID string) (services.RelayDeps, *rpchttp.HTTP, error) {
	logger.Sugar().Infof("create-clients: dialing ethereum rpc %s", cfg.CosmosToEthConfig.EthRpcUrl)
	ethClient, err := tendermintClient.DialEthRPC(context.Background(), cfg.CosmosToEthConfig.EthRpcUrl, tendermintClient.DefaultRPCTimeout)
	if err != nil {
		return services.RelayDeps{}, nil, fmt.Errorf("failed to connect to Ethereum: %w", err)
	}

	logger.Sugar().Infof("create-clients: creating cosmos rpc client %s", cfg.CosmosToEthConfig.TmRpcUrl)
	cosmosClient, err := tendermintClient.DialCosmosRPC(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket", tendermintClient.DefaultRPCTimeout)
	if err != nil {
		return services.RelayDeps{}, nil, fmt.Errorf("failed to create Cosmos RPC client: %w", err)
	}

	cosmosRouterClientID := cosmosRouterClientIDOrDefault(cfg)
	if cosmosRouterClientID == "" {
		return services.RelayDeps{}, nil, fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
	}

	deps := services.RelayDeps{
		Cosmos: services.CosmosEndpoint{Client: cosmosClient},
		EVM: services.EVMEndpoint{Client: ethClient, BeaconAPIURL: resolveBeaconURL(cfg), Contracts: services.EVMContracts{
			Router: common.HexToAddress(cfg.CosmosToEthConfig.ICS26Address), SignatureVerifier: common.HexToAddress(cfg.CosmosToEthConfig.SignatureVerifier),
			Membership: common.HexToAddress(cfg.CosmosToEthConfig.Membership), Misbehaviour: common.HexToAddress(cfg.CosmosToEthConfig.Misbehaviour),
			UpdateClient: common.HexToAddress(cfg.CosmosToEthConfig.UpdateClient), RoleManager: common.HexToAddress(roleManagerOrDefault(cfg)),
		}},
		IDs: services.ClientIDs{CosmosOnEVM: cosmosRouterClientID, EVMOnCosmos: wasmClientID},
	}

	if err := cosmosClient.Start(); err != nil {
		return services.RelayDeps{}, nil, fmt.Errorf("failed to start Cosmos WS client: %w", err)
	}
	return deps, cosmosClient, nil
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
	if resolveBeaconURL(cfg) == "" {
		return "", fmt.Errorf("an L1 beacon endpoint is required to create the Ethereum light client on Cosmos: " +
			"set eth_beacon_api_url on the eth_to_cosmos module")
	}

	// Resolve where the id will be recorded BEFORE spending an on-chain
	// MsgCreateClient. ibc-go assigns the id and the client cannot be un-created, so
	// discovering the destination is missing afterwards leaves a paid-for client that
	// nothing points at and nothing advances until it expires (#310).
	//
	// selectSource picks a cosmos_to_l2 module into CosmosToEthConfig when the config
	// has no cosmos_to_eth source, and the write-back is pinned to cosmos_to_eth — so
	// an eth_to_cosmos + L2 config passes every other guard and fails only here.
	// Cheap to know, irreversible to learn late.
	if err := assertConfigMemberWritable(configPath, cfg.CosmosToEthConfig.ICS26ClientID,
		"cosmos_wasm_client_id", dirCosmosToEth); err != nil {
		return "", fmt.Errorf("refusing to create an Ethereum light client that could not be recorded: %w", err)
	}

	deps, cosmosClient, err := buildCreateClientsDeps(logger, cfg, "")
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
	wasmClientID, err := worker.CreateEthClient(context.Background(), deps.Cosmos, deps.EVM, deps.IDs.CosmosOnEVM, wasmChecksum)
	if err != nil {
		return "", fmt.Errorf("failed to create Ethereum client on Cosmos: %w", err)
	}
	logger.Sugar().Infof("Ethereum light client created on Cosmos: clientID=%s", wasmClientID)

	if err := writeEthWasmClientID(configPath, cfg.CosmosToEthConfig.ICS26ClientID, wasmClientID); err != nil {
		return "", unrecordedClientError(configPath, wasmClientID, err)
	}
	logger.Sugar().Infof("create-clients-cosmos: wrote cosmos_wasm_client_id=%s into %s", wasmClientID, configPath)
	return wasmClientID, nil
}

// runCreateClientsEth deploys the Cosmos (ICS07 Tendermint) light client on
// Ethereum, registers wasmClientID as its counterparty, and persists the
// spectre_client address back into config. wasmClientID must already be known
// (created by the Cosmos step) so the on-chain counterparty is wired to the
// real id. Idempotent: if spectre_client already has deployed code, it is reused.
func runCreateClientsEth(logger *zap.Logger, cfg *appConfig, configPath, wasmClientID, trustLevel string, trustingPeriod uint32) (common.Address, error) {
	if wasmClientID == "" {
		return common.Address{}, fmt.Errorf("cosmos_wasm_client_id is empty; run create-clients-cosmos first")
	}

	deps, cosmosClient, err := buildCreateClientsDeps(logger, cfg, wasmClientID)
	if err != nil {
		return common.Address{}, err
	}
	defer cosmosClient.Stop()

	// Idempotency: skip the deploy when the configured spectre_client is a light
	// client this run can actually keep, so re-running after a partial failure
	// does not redeploy ICS07. Reusability is a property of the client's state,
	// not of the address having code — see reusableSpectreClientAt.
	if cfg.CosmosToEthConfig.SpectreClient != "" {
		addr := common.HexToAddress(cfg.CosmosToEthConfig.SpectreClient)
		ok, reason := reusableSpectreClientAt(deps.EVM, cosmosClient, addr, cfg.CosmosToEthConfig.ICS26ClientID, wasmClientID)
		if ok {
			logger.Sugar().Infof(
				"create-clients-eth: reusing the SpectreClient already deployed at %s; skipping deploy", addr.Hex())
			return addr, nil
		}
		logger.Sugar().Warnf(
			"create-clients-eth: not reusing spectre_client %s — %s; deploying a new one",
			addr.Hex(), reason)
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
	ics07Addr, err := worker.CreateCosmosClient(context.Background(), deps.Cosmos, deps.EVM, deps.IDs, proofType, trustingPeriod, 0, trustLevel, clockDrift)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create Cosmos client on Ethereum: %w", err)
	}
	if (ics07Addr == common.Address{}) {
		return common.Address{}, fmt.Errorf("ics07 address missing after deploy")
	}

	if err := writeSpectreClientAddress(configPath, cfg.CosmosToEthConfig.ICS26ClientID, ics07Addr.Hex()); err != nil {
		return common.Address{}, fmt.Errorf("persist ics07 address to %s: %w", configPath, err)
	}
	logger.Sugar().Infof("create-clients-eth: wrote spectre_client=%s into %s", ics07Addr.Hex(), configPath)
	return ics07Addr, nil
}

// CreateClients runs the full two-chain setup as a one-shot (devnet bring-up):
// create-clients-cosmos first (so the auto-assigned wasm client id is known),
// then create-clients-eth wired to that id. No id guessing, no assertion.
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
			// This command creates the L2 wasm clients, so their ids are allowed
			// to still be empty here — it is what fills them in (#309).
			cfg, err := loadConfigForClientCreation(configPath)
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
			logConfigTarget(logger, "create-clients-cosmos", configPath, cfg)
			if err := preflightCreateClients(cfg); err != nil {
				return err
			}
			wasmChecksum, err := cmd.Flags().GetString(flagWasmChecksum)
			if err != nil {
				return fmt.Errorf("failed to get wasm checksum: %w", err)
			}
			l2Configs, err := cmd.Flags().GetStringArray(flagL2Config)
			if err != nil {
				return err
			}

			// One invocation creates one kind of client. --l2-config means "create these
			// L2 clients"; without it the command creates the Ethereum client for the
			// Cosmos<->Ethereum path.
			//
			// The two used to be combined, because every L2 client was anchored to an
			// Ethereum client whose id had to be injected into its rollup profile. The
			// attestor-trusted profile has no ethereum_client member, so that reason is
			// gone — and combining them is actively harmful. runCreateClientsCosmos keys
			// its cosmos_wasm_client_id write-back on cfg.CosmosToEthConfig, which
			// selectSource may have set to a cosmos_to_l2 module: with --source naming an
			// L2 (say arb-client-0) the Ethereum MsgCreateClient lands and the write-back
			// then fails because no cosmos_to_eth module carries that id, leaving an
			// orphaned client on chain. When an L2 module instead shares its id with the
			// Ethereum one (config.example.json ships cosmos-to-eth and cosmos-to-op both
			// on cosmoshub-1), the write-back succeeds and repoints the Ethereum module at
			// a brand-new client the Ethereum-side SpectreClient was never registered
			// against — every recvPacket then reverts on a counterparty mismatch and the
			// original client is orphaned to expire.
			//
			// Splitting the invocations removes the ambiguity instead of trying to detect
			// it, and matches the one-command-per-chain direction of #255.
			if len(l2Configs) == 0 {
				if _, err := runCreateClientsCosmos(logger, cfg, configPath, wasmChecksum); err != nil {
					return err
				}
				return nil
			}
			logger.Sugar().Infof(
				"create-clients-cosmos: creating %d L2 client(s); the Ethereum client is not touched "+
					"(run this command without --l2-config to create it)", len(l2Configs))
			// Then create one L2 wasm client on Cosmos per --l2-config. All Cosmos-side,
			// same signer — one command per chain (see #255).
			for _, l2Path := range l2Configs {
				l2cfg, err := loadL2ClientConfig(l2Path)
				if err != nil {
					return err
				}
				l2ClientID, err := runCreateClientsL2(logger, cfg, l2cfg)
				if err != nil {
					return fmt.Errorf("create L2 client from %s: %w", l2Path, err)
				}
				// A cosmos_to_l2 module's cosmos_wasm_client_id is the Cosmos client
				// tracking the L2. create-clients-eth reads it to register the
				// SpectreClient's counterparty, so a wrong id there makes every
				// recvPacket revert on a counterparty mismatch. The module is keyed by
				// the L2-side client id (the rollup router's client that tracks Cosmos).
				target := l2cfg.CounterpartyClientID
				if target == "" {
					logger.Sugar().Warnf(
						"create-clients-cosmos[l2]: %s has no counterparty_client_id, cannot locate its module; "+
							"set cosmos_wasm_client_id=%s manually before create-clients-eth",
						l2Path, l2ClientID)
					continue
				}
				if err := writeL2WasmClientID(configPath, target, l2ClientID); err != nil {
					return fmt.Errorf("persist cosmos_wasm_client_id for %s: %w", l2Path, err)
				}
				logger.Sugar().Infof(
					"create-clients-cosmos[l2]: wrote cosmos_wasm_client_id=%s into the %q module",
					l2ClientID, target)

				// The return path (l2_to_cosmos) queries the same client, and its id is
				// not predictable before MsgCreateClient lands — so write it back here
				// rather than leaving the operator to copy it across two modules.
				if err := writeL2SourceClientIDs(configPath, target, l2ClientID); err != nil {
					return fmt.Errorf("persist l2_to_cosmos client ids for %s: %w", l2Path, err)
				}
				logger.Sugar().Infof(
					"create-clients-cosmos[l2]: wrote l2_wasm_client_id=%s into the %q l2_to_cosmos module",
					l2ClientID, target)
			}
			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().String(flagWasmChecksum, "", "wasm checksum for Ethereum light client (hex)")
	cmd.Flags().String(flagSource, "", "ics26_client_id of the cosmos_to_eth source to target (required when several are configured)")
	cmd.Flags().StringArray(flagL2Config, nil, "path to an L2 client config JSON; repeatable, one per L2 rollup source — creates its L2 wasm client on Cosmos")
	return cmd
}

// CreateClientsEth deploys only the Cosmos light client on Ethereum (ICS07) and
// persists spectre_client. Requires cosmos_wasm_client_id (run create-clients-cosmos first).
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
			logConfigTarget(logger, "create-clients-eth", configPath, cfg)
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
	cmd.Flags().String(flagSource, "", "ics26_client_id of the cosmos_to_eth source or cosmos_to_l2 destination to target (required when several are configured)")
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

			ethClient, err := tendermintClient.DialEthRPC(context.Background(), cfg.CosmosToEthConfig.EthRpcUrl, tendermintClient.DefaultRPCTimeout)
			if err != nil {
				return fmt.Errorf("failed to connect to Ethereum: %w", err)
			}

			cosmosClient, err := tendermintClient.DialCosmosRPC(cfg.CosmosToEthConfig.TmRpcUrl, "/websocket", tendermintClient.DefaultRPCTimeout)
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

			cosmosRouterClientID := cosmosRouterClientIDOrDefault(cfg)
			if cosmosRouterClientID == "" {
				return fmt.Errorf("cosmos router client ID (ICS26_CLIENT_ID or ics26_client_id) is required and cannot be empty")
			}
			roleManager := roleManagerOrDefault(cfg)
			if cfg.CosmosToEthConfig.SpectreClient == "" {
				return fmt.Errorf("spectre_client address is required in cosmos_to_eth config")
			}
			deps := services.RelayDeps{
				Cosmos: services.CosmosEndpoint{Client: cosmosClient},
				EVM: services.EVMEndpoint{Client: ethClient, BeaconAPIURL: cfg.EthToCosmosConfig.BeaconUrl, Contracts: services.EVMContracts{
					Router: common.HexToAddress(cfg.CosmosToEthConfig.ICS26Address), SignatureVerifier: common.HexToAddress(cfg.CosmosToEthConfig.SignatureVerifier),
					Membership: common.HexToAddress(cfg.CosmosToEthConfig.Membership), Misbehaviour: common.HexToAddress(cfg.CosmosToEthConfig.Misbehaviour),
					UpdateClient: common.HexToAddress(cfg.CosmosToEthConfig.UpdateClient), RoleManager: common.HexToAddress(roleManager), SpectreClient: common.HexToAddress(cfg.CosmosToEthConfig.SpectreClient),
				}},
				IDs: services.ClientIDs{CosmosOnEVM: cosmosRouterClientID, EVMOnCosmos: cosmosWasmClientID},
			}

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
			if cfg.CosmosToEthConfig.RotationThreshold != "" {
				cosmosConfig.RotationThreshold = cfg.CosmosToEthConfig.RotationThreshold
			}
			if cfg.CosmosToEthConfig.RefreshInterval != 0 {
				cosmosConfig.RefreshInterval = time.Duration(cfg.CosmosToEthConfig.RefreshInterval) * time.Second
				cosmosConfig.RefreshIntervalConfigured = true
			}
			if envVal := os.Getenv("FETCH_TIMEOUT"); envVal != "" {
				if d, err := strconv.Atoi(envVal); err == nil && d > 0 {
					cosmosConfig.FetchTimeout = time.Duration(d) * time.Second
				}
			}
			cosmosConfig.BatchConfig = cfg.BatchConfig
			deps.Config = cosmosConfig

			if err := cosmosClient.Start(); err != nil {
				return fmt.Errorf("failed to start Cosmos WS client: %w", err)
			}
			defer cosmosClient.Stop()

			trustedBlock, err := cmd.Flags().GetInt64(flagTrustedBlock)
			if err != nil {
				return fmt.Errorf("failed to get trusted block: %w", err)
			}
			forceRotation, err := cmd.Flags().GetBool(flagForceRotation)
			if err != nil {
				return fmt.Errorf("failed to get force-rotation flag: %w", err)
			}
			targetHeight, err := cmd.Flags().GetInt64(flagTargetHeight)
			if err != nil {
				return fmt.Errorf("failed to get target-height flag: %w", err)
			}
			maxHops, err := cmd.Flags().GetInt64(flagMaxHops)
			if err != nil {
				return fmt.Errorf("failed to get max-hops flag: %w", err)
			}
			if maxHops <= 0 {
				return fmt.Errorf("update-client: --max-hops must be greater than 0")
			}

			worker := services.NewWorker(&transaction.Handler{}, p)

			// RLY-01: validator churn can force the builder onto a multi-hop
			// path that only advances the client to an intermediate height
			// per call. Loop only while the builder reports that it actually
			// selected a hop; an ordinary update to the live latest height is
			// complete even if the chain produces another block while the ETH
			// transaction is landing.
			var latestBlock *tendermintClient.LightBlock
			for hop := int64(0); ; hop++ {
				if hop >= maxHops {
					return fmt.Errorf("update-client: reached --max-hops=%d without catching up", maxHops)
				}
				result, err := worker.BuildCosmosClientUpdateMsg(
					deps.Cosmos,
					deps.EVM,
					cosmosConfig.FetchTimeout,
					cosmosConfig.RotationThreshold,
					cosmosConfig.ProofType,
					trustedBlock,
					cosmosConfig.TrustLevel,
					forceRotation,
					targetHeight,
				)
				if err != nil {
					return fmt.Errorf("failed to build Cosmos client update on Ethereum: %w", err)
				}
				if result == nil || result.LightBlock == nil {
					return fmt.Errorf("update-client: builder returned nil light block")
				}
				if result.HasMsg {
					if err := worker.TxHandler.SendEthTx(cmd.Context(), deps.EVM, deps.IDs.CosmosOnEVM, *result); err != nil {
						return fmt.Errorf("failed to update Cosmos client on Ethereum: %w", err)
					}
				}
				latestBlock = result.LightBlock
				logger.Sugar().Infof("update-client hop %d complete: height=%d kind=%d isHop=%t hopTarget=%d hasMsg=%t",
					hop, latestBlock.BlockHeight, result.Kind, result.IsHop, result.HopTarget, result.HasMsg)
				if !result.IsHop {
					break
				}
				trustedBlock = 0 // re-derive from on-chain state next iteration
			}
			logger.Sugar().Infof("update-client complete: latest_height=%d", latestBlock.BlockHeight)

			return nil
		},
	}
	cmd.Flags().String(flagConfigPath, "config.json", "path to JSON config file")
	cmd.Flags().Int64(flagTrustedBlock, 0, "trusted Cosmos block height hint; 0 reads from the on-chain client")
	cmd.Flags().Bool(flagGPUProve, false, "use the ICICLE GPU backend for proving (or set GPU_PROVE=1); requires an icicle-enabled build")
	cmd.Flags().Bool(flagForceRotation, false, "force pinned-set rotation regardless of overlap threshold (operator stopgap for RLY-01)")
	cmd.Flags().Int64(flagTargetHeight, 0, "override the update target height; 0 uses the chain's current latest height (operator stopgap for RLY-01)")
	cmd.Flags().Int64(flagMaxHops, 16, "maximum multi-hop iterations before giving up in one invocation")
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

			runCtx, stopSignals := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
			defer stopSignals()

			// Load JSON config
			cfg, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}
			// Checks that only apply to relaying, so they live here rather than in
			// the shared loadConfig. Run before the prover load: a config error
			// should surface immediately, not after the bucket registry is read.
			if err := validateRelayStartupConfig(cfg); err != nil {
				return err
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
			l2Dests := cfg.CosmosToL2Configs
			l2Sources := cfg.L2ToCosmosConfigs
			if len(sources) == 0 && len(l2Dests) == 0 && len(l2Sources) == 0 {
				return fmt.Errorf("no relay source configured in %s (need a cosmos_to_eth, cosmos_to_l2, or l2_to_cosmos module)", configPath)
			}
			// validateL2TimeoutReturnPathConfigs is gone with #328: the return path
			// is now resolved per dest by matching the L2 chain id at build time
			// (l2TimeoutReturnPathConfigForDest), so a separate up-front pass would
			// only duplicate it. This check stays because it answers a different
			// question -- whether each l2_to_cosmos source's configured chain id is
			// the chain its RPC actually serves -- and it must run before anything
			// dials.
			// runCtx, not cmd.Context(): the probe wraps whatever it is given in a
			// 10s timeout, so on cmd.Context() a Ctrl-C during a hung L2 RPC waited
			// out 10s per source before the process could exit.
			if err := validateL2ChainIDs(runCtx, l2Sources); err != nil {
				return err
			}
			// Env overrides (ICS26_CLIENT_ID, COSMOS_WASM_CLIENT_ID, ROLE_MANAGER)
			// name a single source; only honor them when exactly one is
			// configured, otherwise they would wrongly apply to every source.
			allowEnvOverride := len(sources) == 1

			// One shared TransactionHandler across all sources. What sharing buys is
			// the mutex, not the cache: two sources can target the SAME (chain,
			// signer) pair, and only a single h.mu can serialize their nonce
			// allocation. Per-source handlers would each take their own lock and
			// hand the same nonce to both.
			//
			// Sharing the cache is safe because it is keyed by (chainID, address)
			// — see transaction.evmNonceKey — so sources submitting to DIFFERENT
			// EVM chains keep separate nonces for the same signer. Before that
			// keying the cache was a single counter, and a nonce allocated against
			// one L2 was reused against another as a future nonce: accepted into
			// the queued pool, never minable, and never surfaced as an error (#320).
			//
			// The Cosmos side re-queries the account sequence fresh under cosmosMu,
			// so sharing is safe there too.
			txHandler := &transaction.Handler{}

			// One independent relay loop per Cosmos→ETH source. Each has its own
			// Tendermint RPC, SpectreClient and router client id; they share the
			// prover, the TransactionHandler, the ETH beacon endpoint, and the
			// same ICS26Router (ETH events are partitioned by the per-source
			// router client id filter).
			var wg sync.WaitGroup
			total := len(sources) + len(l2Dests) + len(l2Sources)
			loopErrCh := make(chan error, total)
			cleanups := make([]func(), 0, total)
			relayCtx, cancelRelays := context.WithCancel(runCtx)
			defer cancelRelays()
			l2ReturnPaths := make([]l2TimeoutReturnPathConfig, 0, len(l2Dests))
			onceCleanup := func(cleanup func()) func() {
				var once sync.Once
				return func() {
					if cleanup != nil {
						once.Do(cleanup)
					}
				}
			}
			stopAndCleanup := func() {
				stopRelaysAndCleanup(cancelRelays, wg.Wait, cleanups)
			}
			for i := range sources {
				svc, deps, cleanup, err := buildCosmosToEthSource(
					logger, sources[i], cfg.EthToCosmosConfig, cfg.BatchConfig, p, txHandler,
					buildCosmosToEthSourceOptions{
						allowEnvOverride:   allowEnvOverride,
						startSubscriptions: true,
					},
				)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("cosmos_to_eth source %q: %w", sources[i].ICS26ClientID, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				wg.Add(1)
				go func(svc *services.Services, deps services.RelayDeps, cleanup func()) {
					defer wg.Done()
					defer cleanup()
					// Chain-adapter RelayModule engine — the sole relay engine since
					// the legacy services.StartLoop was removed after the cutover.
					if err := runAdapterEngine(relayCtx, svc, deps); err != nil {
						loopErrCh <- fmt.Errorf("cosmos_to_eth source %q: %w", deps.IDs.CosmosOnEVM, err)
					}
				}(svc, deps, cleanup)
			}

			// One independent relay loop per Cosmos→L2 destination: the same groth16
			// pipeline as Cosmos→ETH pointed at the L2's SpectreClient/ICS26Router,
			// with no reverse beacon direction.
			//
			// Built in two passes: first resolve every dest's return path and match
			// every l2_to_cosmos source, then launch the goroutines below. This way
			// a return-path mismatch aborts before any Cosmos→L2 connection opens,
			// instead of after the relay loop is already dialing/subscribing.
			type builtL2Dest struct {
				svc     *services.Services
				deps    services.RelayDeps
				cleanup func()
			}
			builtL2Dests := make([]builtL2Dest, 0, len(l2Dests))
			for i := range l2Dests {
				svc, deps, cleanup, err := buildCosmosToL2Dest(
					logger, l2Dests[i], cfg.BatchConfig, p, txHandler,
				)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("cosmos_to_l2 dest %q: %w", l2Dests[i].ICS26ClientID, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				name := ""
				if i < len(cfg.CosmosToL2Names) {
					name = cfg.CosmosToL2Names[i]
				}
				// Only resolve the return-path chain id when there's an
				// l2_to_cosmos source to match it against; otherwise a
				// forward-only deployment pays an avoidable L2 RPC round
				// trip (and startup-abort risk) for a path nothing consults.
				if len(l2Sources) > 0 {
					returnPath, err := l2TimeoutReturnPathConfigForDest(runCtx, name, l2Dests[i], l2TimeoutReturnPath{svc: svc, deps: deps})
					if err != nil {
						stopAndCleanup()
						return fmt.Errorf("cosmos_to_l2 dest %q: %w", l2Dests[i].ICS26ClientID, err)
					}
					l2ReturnPaths = append(l2ReturnPaths, returnPath)
				}
				// Accumulate rather than launching here: nothing dials, subscribes
				// or submits until every dest and source has built and matched.
				builtL2Dests = append(builtL2Dests, builtL2Dest{svc: svc, deps: deps, cleanup: cleanup})
			}

			// One independent relay module per L2->Cosmos source (opstack/arbitrum).
			// Each dials its own L1/L2/Cosmos clients + attestor sidecar; they share
			// the TransactionHandler (same Cosmos signer -> shared sequence path).
			type builtL2Source struct {
				module   *relay.Module
				cleanup  func()
				srcChain string
			}
			builtL2Sources := make([]builtL2Source, 0, len(l2Sources))
			for i := range l2Sources {
				timeoutReturn, err := findL2TimeoutReturnPath(runCtx, l2Sources[i], l2ReturnPaths)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("l2_to_cosmos source %q: %w", l2Sources[i].AttestorSrcChain, err)
				}
				module, cleanup, err := buildL2ToCosmosModule(logger, l2Sources[i], txHandler, timeoutReturn)
				if err != nil {
					stopAndCleanup()
					return fmt.Errorf("l2_to_cosmos source %q: %w", l2Sources[i].AttestorSrcChain, err)
				}
				cleanup = onceCleanup(cleanup)
				cleanups = append(cleanups, cleanup)
				builtL2Sources = append(builtL2Sources, builtL2Source{module: module, cleanup: cleanup, srcChain: l2Sources[i].AttestorSrcChain})
			}

			// Every Cosmos→L2 dest and L2→Cosmos source above built and matched
			// cleanly — only now do we start dialing/subscribing/submitting.
			for _, d := range builtL2Dests {
				wg.Add(1)
				go func(svc *services.Services, deps services.RelayDeps, cleanup func()) {
					defer wg.Done()
					defer cleanup()
					if err := runCosmosToL2Engine(relayCtx, svc, deps); err != nil {
						loopErrCh <- fmt.Errorf("cosmos_to_l2 dest %q: %w", deps.IDs.CosmosOnEVM, err)
					}
				}(d.svc, d.deps, d.cleanup)
			}
			for _, s := range builtL2Sources {
				wg.Add(1)
				go func(module *relay.Module, cleanup func(), srcChain string) {
					defer wg.Done()
					defer cleanup()
					if err := runL2Engine(relayCtx, module); err != nil {
						loopErrCh <- fmt.Errorf("l2_to_cosmos source %q: %w", srcChain, err)
					}
				}(s.module, s.cleanup, s.srcChain)
			}

			logger.Sugar().Infof("Relayer started: relaying %d Cosmos→ETH + %d Cosmos→L2 + %d L2→Cosmos source(s)", len(sources), len(l2Dests), len(l2Sources))
			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case err := <-loopErrCh:
				stopAndCleanup()
				return err
			case <-done:
				// A worker sends its error before its deferred wg.Done. If all
				// workers have returned, prefer that buffered error over treating
				// the coincident done signal as a clean shutdown.
				select {
				case err := <-loopErrCh:
					return err
				default:
					return nil
				}
			case <-runCtx.Done():
				logger.Sugar().Infof("Relayer shutdown requested: %v", runCtx.Err())
				stopAndCleanup()
				logger.Sugar().Info("Relayer clients stopped; exiting")
			}

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
			tendermintRpcClient, err := tendermintClient.DialCosmosRPC(tendermintRpcEndpoint, "/websocket", tendermintClient.DefaultRPCTimeout)
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
