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
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"

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
			one.kind = chain.ChainType(m.SrcChain) // opstack | arbitrum
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
