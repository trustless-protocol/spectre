// This file owns the config end to end: the schema, loading it, validating it,
// and editing the file in place.
//
// They are one file on purpose. They read and write the same structure, so
// splitting them would mean touching three files every time a key is added.
//
// The in-place editing deserves a note. create-clients writes client ids back
// into config.json and must PRESERVE the operator's formatting, which is why
// there is a byte-level JSON parser here (skipJSONSpace, scanJSONStringEnd,
// skipJSONValue, findJSONObjectMember, insertJSONObjectMember) rather than a
// round trip through encoding/json -- that would rewrite the whole file. A
// hand-written parser is also where subtle bugs live, and #310 was one in this
// exact area, so it is worth having it findable by file name instead of buried
// among two thousand lines of command definitions.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	routerContract "relayer/bindings/ICS26Router"
	"relayer/chain"
	tendermintClient "relayer/client"
	"relayer/services"
	"relayer/transaction"
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
	case chain.Cosmos, chain.Ethereum, chain.OPStack, chain.Arbitrum, chain.Avalanche:
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
		// Avalanche's C-Chain (an L1, not a rollup) deliberately reuses the same
		// two directions as the rollups: the Cosmos->EVM leg is chain-agnostic
		// groth16+EVM, and the return leg is the shared attested-header client
		// with a coreth header profile.
		case (src == chain.OPStack || src == chain.Arbitrum || src == chain.Avalanche) && dst == chain.Cosmos:
			return dirL2ToCosmos, false, nil
		case src == chain.Cosmos && (dst == chain.OPStack || dst == chain.Arbitrum || dst == chain.Avalanche):
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
	type finalityOwner struct {
		module string
		head   string
	}
	seenAttestorKeyFinality := make(map[[32]byte]finalityOwner)
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
			one.kind = chain.ChainType(m.SrcChain) // opstack | arbitrum | avalanche
			if err := one.validateWith(requireL2WasmClientID); err != nil {
				return nil, fmt.Errorf("module %q: %w", m.Name, err)
			}
			head, err := parseHeadKind(one.HeadKind)
			if err != nil {
				return nil, fmt.Errorf("module %q: %w", m.Name, err)
			}
			for _, publicKey := range one.Attestors.PublicKeys {
				var key [32]byte
				copy(key[:], publicKey)
				if owner, exists := seenAttestorKeyFinality[key]; exists && owner.head != head.String() {
					return nil, fmt.Errorf(
						"modules %q (%s) and %q (%s) reuse attestor public key %x across different head_kind values; provision disjoint keys per finality tier",
						owner.module, owner.head, m.Name, head.String(), key,
					)
				}
				seenAttestorKeyFinality[key] = finalityOwner{module: m.Name, head: head.String()}
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

// validateSingleRelayPair enforces that one `start` process serves exactly one
// source → destination pair.
//
// The rule is not a simplification for its own sake — running several paths in one
// process is what forces every shared structure between them: one nonce cache
// across unrelated chains (#320), one loopErrCh where a single dead loop takes the
// others down (#321), and one state directory two paths write into. A process per
// path removes all three by construction, so the blast radius of a failure is that
// path alone.
//
// A Cosmos→L2 destination and its L2→Cosmos return leg are ONE pair: they relay
// opposite directions of the same path, and the return leg resolves its timeout
// path from the forward one (findL2TimeoutReturnPath). The Cosmos↔ETH pair needs no
// equivalent grouping — runAdapterEngine already drives both directions from a
// single cosmos_to_eth module, with eth_to_cosmos supplying only the beacon URL.
//
// This runs before the prover load so a config mistake surfaces at once rather than
// after the bucket registry is read.
func validateSingleRelayPair(cfg *appConfig, configPath string) error {
	if cfg == nil {
		return fmt.Errorf("no relay source configured in %s (need a cosmos_to_eth, cosmos_to_l2, or l2_to_cosmos module)", configPath)
	}

	declared := declaredRelayPaths(cfg)
	if len(declared) == 0 {
		return fmt.Errorf("no relay source configured in %s (need a cosmos_to_eth, cosmos_to_l2, or l2_to_cosmos module)", configPath)
	}

	// Count paths, not modules: the L2 forward and return legs pair up, so the
	// number of Cosmos↔L2 paths is however many legs the longer side declares.
	l2Paths := len(cfg.CosmosToL2Configs)
	if n := len(cfg.L2ToCosmosConfigs); n > l2Paths {
		l2Paths = n
	}
	if len(cfg.CosmosToEthConfigs)+l2Paths <= 1 {
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "%s declares more than one relay path, but one relayer process serves exactly one source → destination pair.\n", configPath)
	b.WriteString("declared:\n")
	for _, d := range declared {
		fmt.Fprintf(&b, "  %s\n", d)
	}
	b.WriteString("a cosmos_to_l2 destination and its l2_to_cosmos return leg are one path; everything else is a path of its own.\n")
	b.WriteString("keep one path per config file and start one process per file:\n")
	b.WriteString("  ./relayer start --config <one-path>.json")
	return fmt.Errorf("%s", b.String())
}

// declaredRelayPaths lists every relay module the config declares, each identified
// by the field an operator edits it by. Modules are listed individually rather than
// paired up: which forward leg a return leg belongs to is resolved at runtime from
// the L2 chain id, so pairing them here by file order would name the wrong partner
// exactly when the counts disagree — the case this error exists for.
func declaredRelayPaths(cfg *appConfig) []string {
	out := make([]string, 0, len(cfg.CosmosToEthConfigs)+len(cfg.CosmosToL2Configs)+len(cfg.L2ToCosmosConfigs))
	for i := range cfg.CosmosToEthConfigs {
		out = append(out, fmt.Sprintf("cosmos_to_eth ics26_client_id=%q", cfg.CosmosToEthConfigs[i].ICS26ClientID))
	}
	for i := range cfg.CosmosToL2Configs {
		name := ""
		if i < len(cfg.CosmosToL2Names) {
			name = cfg.CosmosToL2Names[i]
		}
		out = append(out, fmt.Sprintf("cosmos_to_l2 ics26_client_id=%q%s", cfg.CosmosToL2Configs[i].ICS26ClientID, moduleNameSuffix(name)))
	}
	for i := range cfg.L2ToCosmosConfigs {
		out = append(out, fmt.Sprintf("l2_to_cosmos attestor_src_chain=%q", cfg.L2ToCosmosConfigs[i].AttestorSrcChain))
	}
	return out
}

// relayPathID names the single relay path this process serves, compactly enough
// to sit on every log line.
//
// A1 makes this well defined: start refuses a config declaring more than one
// path, so the answer is a process-level constant and can be stamped once as a
// log prefix. Direction is NOT a process-level constant and deliberately absent
// here -- a path is bidirectional (cmd/run_adapters.go runs a cosmos->eth and an
// eth->cosmos module side by side), so the direction belongs on the labels of
// the code that actually knows it.
//
// The client id is part of the identity because three processes all relaying
// cosmos->eth from different source clients is a valid deployment, and their
// logs are merged by journald or Loki with no file boundary left to tell them
// apart. That merge is the problem this solves.
func relayPathID(cfg *appConfig) string {
	switch {
	case len(cfg.CosmosToEthConfigs) > 0:
		return "cosmos<->eth/" + cfg.CosmosToEthConfigs[0].ICS26ClientID
	case len(cfg.CosmosToL2Configs) > 0:
		// Prefer the rollup's own name over a generic "l2": the return leg names
		// the chain family (opstack / arbitrum / base) and an operator running two
		// rollups needs to tell them apart.
		kind := "l2"
		if len(cfg.L2ToCosmosConfigs) > 0 && cfg.L2ToCosmosConfigs[0].AttestorSrcChain != "" {
			kind = cfg.L2ToCosmosConfigs[0].AttestorSrcChain
		}
		return "cosmos<->" + kind + "/" + cfg.CosmosToL2Configs[0].ICS26ClientID
	case len(cfg.L2ToCosmosConfigs) > 0:
		return cfg.L2ToCosmosConfigs[0].AttestorSrcChain + "->cosmos"
	}
	return "relayer"
}

// moduleNameSuffix appends the config module's label when it has one, so an
// operator who labelled their modules sees the label they wrote.
func moduleNameSuffix(name string) string {
	if strings.TrimSpace(name) == "" {
		return ""
	}
	return fmt.Sprintf(" (module %q)", name)
}

// validateStartupConfig runs every check that is cheap, read-only, and able to
// prove this config wrong before the relay loops open.
//
// It exists because the alternative is the most expensive failure shape there
// is: the process starts, loads the prover, relays for hours, and only then
// reveals that ics26_client_id names a client the router never heard of — by
// which time packets are already half-relayed. Everything checked here costs one
// read and no gas, so none of it has a reason to surface later.
//
// Two rules decide what belongs, and both are load-bearing:
//
//   - A DEFINITE DISAGREEMENT is fatal. The endpoint answered, and the answer
//     contradicts the config. No amount of retrying fixes that.
//   - A NON-ANSWERING endpoint is not. That is an availability problem, not a
//     configuration one; the relay loops already retry, and refusing to boot on
//     a blip would turn it into an outage. Same doctrine verifyL2ChainID
//     follows, stated there and reused here.
//
// Waiting states — beacon finality, the attestor frontier — are deliberately out
// of scope. They are operational conditions, not wrong config, and blocking
// startup on them would make a healthy cold start look like a broken one.
//
// Every finding is collected before returning rather than failing on the first.
// A config with three wrong keys should take one run to diagnose, not three.
//
// It runs after validateSingleRelayPair, so every module list below holds at
// most one entry. The lists are still walked rather than indexed: an empty list
// is the normal shape for the direction this process does not serve, and a
// first-entry-only check is the bug validateRelayStartupConfig was written to
// fix.
//
// NOT HERE, deliberately — the attestor half. Two checks belong to this
// validation by intent and cannot be written correctly yet:
//
//   - Asking the attestor at startup whether it can verify state roots, instead
//     of warning once per process and relaying on. The mechanism meant to carry
//     the answer — a capabilities list on ChainInfo — is superseded by the
//     per-header Ed25519 attestations in #417, so building against the old shape
//     would add a check for a field the new attestor does not send.
//   - Rejecting an include_provisional that cannot take effect. Whether it is a
//     no-op depends on the FEED's attestation_head, which is the attestor's
//     answer, not a config key here. Keying the warning off head_kind instead
//     would look right and fire on the wrong configs.
//
// Both land once #417 settles the attestor's startup contract.
func validateStartupConfig(stdCtx context.Context, cfg *appConfig) error {
	if cfg == nil {
		return nil
	}

	var findings []configFinding
	findings = append(findings, validateStateDir()...)
	findings = append(findings, validateCosmosChainID(stdCtx, cfg)...)
	findings = append(findings, validateEVMDeployments(stdCtx, cfg)...)

	if len(findings) == 0 {
		log.Printf("[start] config validated")
		return nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "config validation found %d problem(s); the relayer stops here rather than "+
		"surfacing them at the first packet:", len(findings))
	for _, f := range findings {
		fmt.Fprintf(&b, "\n  - %s: %s", f.key, f.detail)
	}
	return fmt.Errorf("%s", b.String())
}

// configFinding is one definite disagreement between the config and the world,
// named by the config key or environment variable at fault.
type configFinding struct {
	key    string
	detail string
}

// configProbeTimeout bounds one read taken while validating. Same value and same
// reason as l2ChainIDProbeTimeout: a rate-limited public endpoint can take
// seconds to answer, and a slow answer must not be mistaken for a wrong one.
const configProbeTimeout = 10 * time.Second

// stateDirProbeFile is the basename written and immediately removed to prove the
// state directory is writable. It carries the pid so two relayers probing the
// same directory cannot delete each other's probe.
const stateDirProbeFile = ".relayer-write-probe"

// validateStateDir proves the recovery-cursor directory can actually be written
// before anything depends on it.
//
// An unwritable state directory does not stop the relayer today: the cursor save
// fails per tick, deep inside the recovery path, and the process keeps relaying
// with an in-memory cursor that dies with it. The next restart then silently
// falls back to the lookback window — the exact event-loss shape persisted
// cursors exist to close.
//
// This creates the directory rather than inspecting an ancestor's mode bits.
// Reading permissions and predicting the outcome is a guess, and guessing is the
// failure this check removes; the run creates that same directory moments later
// anyway. Nothing on a chain changes, and nothing another process can observe
// beyond a directory that was about to exist.
func validateStateDir() []configFinding {
	path := services.RecoveryStatePath()
	key := "RELAYER_RECOVERY_STATE_FILE"
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return []configFinding{{
			key: key,
			detail: fmt.Sprintf("recovery state file %q needs directory %q, which cannot be created: %v. "+
				"Without it every recovery cursor is lost on restart and the relayer falls back to its "+
				"lookback window, which is the event-loss window the cursors exist to close", path, dir, err),
		}}
	}

	probe := filepath.Join(dir, fmt.Sprintf("%s.%d", stateDirProbeFile, os.Getpid()))
	f, err := os.OpenFile(probe, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600)
	if err != nil {
		return []configFinding{{
			key: key,
			detail: fmt.Sprintf("directory %q for recovery state file %q is not writable: %v. "+
				"Recovery cursors would silently fail to persist and every restart would fall back "+
				"to the lookback window", dir, path, err),
		}}
	}
	closeErr := f.Close()
	removeErr := os.Remove(probe)
	if closeErr != nil || removeErr != nil {
		return []configFinding{{
			key: key,
			detail: fmt.Sprintf("directory %q for recovery state file %q accepted a file but not its "+
				"completion (close: %v, remove: %v)", dir, path, closeErr, removeErr),
		}}
	}
	return nil
}

// validateCosmosChainID compares COSMOS_CHAIN_ID against the chain the configured
// Cosmos endpoint actually serves.
//
// The variable is read only when a Cosmos transaction is built
// (transaction/handler.go), so today a wrong or missing value surfaces at the
// first Cosmos send — after the packet was received, the proof was built and the
// fee was budgeted.
//
// A Cosmos↔L2 path declares tm_rpc_url on both of its legs, pointing at the same
// chain, so the endpoints are de-duplicated and each distinct one is asked once.
func validateCosmosChainID(stdCtx context.Context, cfg *appConfig) []configFinding {
	urls := cosmosEndpointsInConfig(cfg)
	if len(urls) == 0 {
		return nil
	}

	want := strings.TrimSpace(os.Getenv("COSMOS_CHAIN_ID"))
	if want == "" {
		return []configFinding{{
			key: "COSMOS_CHAIN_ID",
			detail: "not set, but this config relays through Cosmos. Every Cosmos transaction is signed " +
				"with it, so the relayer would start, receive packets, build proofs, and only then fail " +
				"on the first send",
		}}
	}

	var findings []configFinding
	for _, url := range urls {
		got, ok := probeCosmosNetwork(stdCtx, url)
		if !ok {
			log.Printf("[start] could not read the chain id from %s; skipping the COSMOS_CHAIN_ID check "+
				"for that endpoint (declared %q). The relay loops keep retrying it.", url, want)
			continue
		}
		if got != want {
			findings = append(findings, configFinding{
				key: "COSMOS_CHAIN_ID / tm_rpc_url",
				detail: fmt.Sprintf("COSMOS_CHAIN_ID is %q but tm_rpc_url %s serves chain %q. Transactions "+
					"signed for %q are rejected by %q, so no Cosmos send from this process can succeed",
					want, url, got, want, got),
			})
			continue
		}
		log.Printf("[start] cosmos chain id verified for %s: %s", url, got)
	}
	return findings
}

// cosmosEndpointsInConfig returns each distinct tm_rpc_url in the config, in file
// order. Every relay direction here has a Cosmos side, so all three module lists
// contribute one.
func cosmosEndpointsInConfig(cfg *appConfig) []string {
	var urls []string
	seen := map[string]bool{}
	add := func(u string) {
		if u == "" || seen[u] {
			return
		}
		seen[u] = true
		urls = append(urls, u)
	}
	for i := range cfg.CosmosToEthConfigs {
		add(cfg.CosmosToEthConfigs[i].TmRpcUrl)
	}
	for i := range cfg.CosmosToL2Configs {
		add(cfg.CosmosToL2Configs[i].TmRpcUrl)
	}
	for i := range cfg.L2ToCosmosConfigs {
		add(cfg.L2ToCosmosConfigs[i].TmRpcUrl)
	}
	return urls
}

// probeCosmosNetwork returns the chain id the endpoint reports. ok is false when
// it did not answer at all — the tolerated case, kept distinct from an empty
// answer so "no answer" is never compared as if it were one.
func probeCosmosNetwork(stdCtx context.Context, url string) (network string, ok bool) {
	ctx, cancel := context.WithTimeout(stdCtx, configProbeTimeout)
	defer cancel()

	c, err := tendermintClient.DialCosmosRPC(url, "/websocket", configProbeTimeout)
	if err != nil {
		return "", false
	}
	status, err := c.Status(ctx)
	if err != nil || status == nil {
		return "", false
	}
	return status.NodeInfo.Network, true
}

// evmDeployment is one EVM endpoint and everything this config claims is
// deployed on it.
type evmDeployment struct {
	// spectreAddr and wasmClientID are what the router's registration must agree
	// with. Empty means "not declared for this leg", and the check is skipped.
	spectreAddr  string
	wasmClientID string
	label        string
	rpcURL       string
	routerAddr   string
	routerCliID  string
	// routerKey and clientIDKey are the config keys the router address and the
	// client id came from. They differ per module kind — ics26_address /
	// ics26_client_id on the outbound legs, rollup_profile.common.l2_router /
	// l2_ics26_client_id on the return leg — and a finding must name the key the
	// operator will grep for, not the one this file happens to know.
	routerKey   string
	clientIDKey string
	// addresses maps a config key to the address it names. Only well-formed,
	// non-empty addresses are listed: an unset optional address is a different
	// question, already answered by the required-field checks in loadConfig.
	addresses map[string]string
}

// validateEVMDeployments proves that every address this config points at holds
// code, and that each router knows the client id its module relays through.
//
// Both failures are invisible until the first packet and then
// indistinguishable: a call to an address with no code and a call for an
// unknown client id both come back as an opaque revert with no reason string.
// Reading eth_getCode and one getClient view call costs nothing and separates
// them.
func validateEVMDeployments(stdCtx context.Context, cfg *appConfig) []configFinding {
	var findings []configFinding

	for _, d := range evmDeploymentsInConfig(cfg) {
		ctx, cancel := context.WithTimeout(stdCtx, configProbeTimeout)
		client, err := tendermintClient.DialEthRPC(ctx, d.rpcURL, configProbeTimeout)
		if err != nil {
			cancel()
			log.Printf("[start] could not dial %s for %s; skipping its contract checks. "+
				"The relay loops keep retrying it.", d.rpcURL, d.label)
			continue
		}

		reachable := true
		for _, key := range sortedKeys(d.addresses) {
			addr := d.addresses[key]
			code, err := client.CodeAt(ctx, common.HexToAddress(addr), nil)
			if err != nil {
				// One unanswered read means the endpoint stopped answering, not
				// that the contract is missing. Stop probing it rather than
				// reporting every remaining address as absent.
				reachable = false
				log.Printf("[start] %s: could not read code at %s (%s) from %s; skipping the remaining "+
					"contract checks for this endpoint", d.label, key, addr, d.rpcURL)
				break
			}
			if len(code) == 0 {
				findings = append(findings, configFinding{
					key: fmt.Sprintf("%s.%s", d.label, key),
					detail: fmt.Sprintf("%s holds no code on %s. Every call the relayer makes to it returns "+
						"empty, which reaches the log as an unexplained revert at the first packet",
						addr, d.rpcURL),
				})
			}
		}

		if reachable {
			findings = append(findings, validateRouterClientID(ctx, client, d, findings)...)
		}
		client.Close()
		cancel()
	}
	return findings
}

// validateRouterClientID asks the router for the client id this module relays
// through. It is skipped when the router itself was already reported as
// code-less, because "the router has no code" and "the router does not know this
// client" would then be one fact reported twice, sending an operator looking for
// a second problem.
// validateRouterClientID asks the router whether it knows d.routerCliID, and
// whether the client it resolves to is wired to THIS deployment.
//
// found is read-only and is used for one thing: suppressing the check when the
// router address itself is already a finding. The return value carries only
// what this function discovered, because the caller appends it -- returning
// found would report every earlier finding a second time. It did: a module with
// five missing contracts came out as ten problems, the router check handing
// back the five it was given.
func validateRouterClientID(ctx context.Context, client *ethclient.Client, d evmDeployment, found []configFinding) []configFinding {
	if d.routerAddr == "" || d.routerCliID == "" {
		return nil
	}
	routerKey := fmt.Sprintf("%s.%s", d.label, d.routerKey)
	for _, f := range found {
		if f.key == routerKey {
			return nil
		}
	}

	caller, err := routerContract.NewContractICS26RouterCaller(common.HexToAddress(d.routerAddr), client)
	if err != nil {
		return nil
	}
	addr, err := caller.GetClient(&bind.CallOpts{Context: ctx}, d.routerCliID)
	if err == nil {
		log.Printf("[start] %s: router %s resolves client id %q to %s",
			d.label, d.routerAddr, d.routerCliID, addr.Hex())
		// Resolving is not the same as being wired to THIS deployment. A client id
		// registered against a different light client, or counterparty-wired to a
		// different Cosmos client, resolves perfectly well and then sends every
		// update and proof at one client while the router verifies another. The
		// pure half of that decision already exists and is tested; only the RPC to
		// feed it was missing here.
		if d.spectreAddr == "" || d.wasmClientID == "" {
			return nil
		}
		// The address half first, and unconditionally: it is already in hand, and a
		// counterparty read that fails must not hide a mismatch we have already
		// proved. Reading the counterparty is a second RPC and gets the usual
		// treatment for a call that did not answer -- logged, not fatal.
		expected := common.HexToAddress(d.spectreAddr)
		counterparty := ""
		if cp, cpErr := caller.GetCounterparty(&bind.CallOpts{Context: ctx}, d.routerCliID); cpErr == nil {
			counterparty = cp.ClientId
		} else if addr == expected {
			log.Printf("[start] %s: could not read the counterparty for client id %q (%v); "+
				"skipping the counterparty half of the wiring check",
				d.label, d.routerCliID, cpErr)
			return nil
		} else {
			counterparty = d.wasmClientID // let the address mismatch be the finding
		}
		if wErr := routerWiringIsReusable(
			routerWiring{client: addr, counterparty: counterparty},
			expected, d.wasmClientID,
		); wErr != nil {
			return []configFinding{{
				key: d.clientIDKey,
				detail: fmt.Sprintf("%s: the router at %s knows client id %q, but %s",
					d.label, d.routerAddr, d.routerCliID, wErr),
			}}
		}
		return nil
	}

	// getClient never returns the zero address: it reverts with
	// IBCClientNotFound (contracts/utils/ICS02ClientUpgradeable.sol:78-80). So a
	// revert IS the answer, and a transport error is not — telling them apart is
	// what keeps this check inside the doctrine on validateStartupConfig.
	if !isEVMRevert(err) {
		log.Printf("[start] %s: could not ask router %s for client id %q (%v); skipping that check. "+
			"The relay loops keep retrying it.", d.label, d.routerAddr, d.routerCliID, err)
		return nil
	}
	return []configFinding{{
		key: fmt.Sprintf("%s.%s", d.label, d.clientIDKey),
		detail: fmt.Sprintf("router %s on %s has no client registered under id %q (%v). Every packet this "+
			"module relays is addressed to that id, so none of them can be delivered",
			d.routerAddr, d.rpcURL, d.routerCliID, err),
	}}
}

// evmDeploymentsInConfig lists one entry per configured EVM endpoint.
//
// cosmos_to_eth and cosmos_to_l2 share a struct and a contract set — to the
// relayer an L2 is just another EVM chain hosting a SpectreClient and a router —
// so they are gathered identically and only the label differs. The l2_to_cosmos
// return leg reaches the same two things on the L2, it only names them
// differently: the address comes from rollup_profile.common.l2_router rather
// than an ics26_address key, and the endpoint is l2_rpc_url. Checking only the
// outbound legs would leave the return path with the failure this validation
// exists to remove, which is the one-sided fix this repo keeps finding.
func evmDeploymentsInConfig(cfg *appConfig) []evmDeployment {
	var out []evmDeployment
	// routerCliID and wasmClientID are passed in rather than derived here,
	// because which of them the RELAY will use differs by path and validating
	// the other one proves nothing. See the two loops below.
	gather := func(c cosmosToEthConfig, label, routerCliID, wasmClientID string) {
		if c.EthRpcUrl == "" {
			return
		}
		addrs := map[string]string{}
		for key, addr := range map[string]string{
			"ics26_address":      c.ICS26Address,
			"spectre_client":     c.SpectreClient,
			"signature_verifier": c.SignatureVerifier,
			"membership":         c.Membership,
			"misbehaviour":       c.Misbehaviour,
			"update_client":      c.UpdateClient,
		} {
			if common.IsHexAddress(addr) {
				addrs[key] = addr
			}
		}
		if len(addrs) == 0 {
			return
		}
		out = append(out, evmDeployment{
			label:        label,
			rpcURL:       c.EthRpcUrl,
			routerAddr:   addrs["ics26_address"],
			routerCliID:  routerCliID,
			routerKey:    "ics26_address",
			clientIDKey:  "ics26_client_id",
			spectreAddr:  addrs["spectre_client"],
			wasmClientID: wasmClientID,
			addresses:    addrs,
		})
	}
	for i := range cfg.CosmosToEthConfigs {
		c := cfg.CosmosToEthConfigs[i]
		// The EFFECTIVE ids, not the ones in the file: buildSource applies these
		// same overrides (build_source.go:90-100), so validating the raw value
		// proves a config startup then does not use. A stale ICS26_CLIENT_ID in
		// .env would pass here and fail at the first packet -- the exact failure
		// this whole function exists to move forward.
		gather(c, "cosmos_to_eth",
			envOrDefault("ICS26_CLIENT_ID", c.ICS26ClientID),
			envOrDefault("COSMOS_WASM_CLIENT_ID", c.CosmosWasmClientID))
	}
	for i := range cfg.CosmosToL2Configs {
		c := cfg.CosmosToL2Configs[i]
		// NOT overridden, and the asymmetry is the point: buildCosmosToL2Dest
		// reads c2l.ICS26ClientID directly and never consults the environment
		// (build_cosmos_to_l2.go:57,112,125). Applying the override here validated
		// an id this path will never use -- so a stale value in .env failed a
		// perfectly good JSON config, and a correct one hid a bad JSON id. Found
		// in review, and it is the mirror of the bug the cosmos_to_eth branch
		// above fixes: the rule is that validation reads whatever the BUILDER
		// reads, per path.
		gather(c, "cosmos_to_l2", c.ICS26ClientID, c.CosmosWasmClientID)
	}
	for i := range cfg.L2ToCosmosConfigs {
		src := cfg.L2ToCosmosConfigs[i]
		if src.L2RpcUrl == "" {
			continue
		}
		router, err := l2RouterFromProfile(src.RollupProfile)
		if err != nil {
			// A malformed profile is already a load-time error; nothing to add.
			continue
		}
		out = append(out, evmDeployment{
			label:       "l2_to_cosmos",
			rpcURL:      src.L2RpcUrl,
			routerAddr:  router.Hex(),
			routerCliID: src.L2ICS26ClientID,
			routerKey:   "rollup_profile.common.l2_router",
			clientIDKey: "l2_ics26_client_id",
			addresses:   map[string]string{"rollup_profile.common.l2_router": router.Hex()},
		})
	}
	return out
}

// isEVMRevert reports whether the node executed the call and the contract
// reverted, as opposed to the call never reaching or never leaving the node.
//
// The distinction decides whether a failure is fatal, so it is made on the two
// signals that mean the EVM ran: a JSON-RPC error carrying revert DATA, or
// geth's "execution reverted" text for a revert that returned none.
//
// Implementing rpc.DataError is NOT one of those signals, even though it is the
// obvious test. go-ethereum's jsonError satisfies that interface for every
// JSON-RPC error a node returns, so "header not found" on a pruned node, a
// -32005 rate limit and a real revert all match it. Treating those as a definite
// disagreement would refuse to boot against a throttled endpoint — precisely the
// blip-into-outage this validation must not cause. Only non-nil ErrorData
// narrows it back to a call that actually executed.
func isEVMRevert(err error) bool {
	if err == nil {
		return false
	}
	var dataErr rpc.DataError
	if errors.As(err, &dataErr) && dataErr.ErrorData() != nil {
		return true
	}
	return strings.Contains(strings.ToLower(err.Error()), "execution reverted")
}

// sortedKeys keeps findings and log lines in a stable order, so two runs against
// the same broken config produce the same message.
func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
