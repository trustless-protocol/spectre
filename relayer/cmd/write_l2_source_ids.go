package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Write-back for the two client ids an l2_to_cosmos module cannot know up front.
//
// Both are assigned by ibc-go's global client sequence when MsgCreateClient lands,
// so config.example.json can only carry placeholders:
//
//   - config.l2_wasm_client_id — the Cosmos-side wasm client tracking the L2
//   - config.rollup_profile.common.ethereum_client.client_id — the shared Ethereum
//     client the L2 client authenticates L1 state through
//
// The second one is the dangerous one. create-clients-cosmos already injects the
// real id into the L2 client's *client state* (that copy is what the contract
// queries), so a stale id left in the config does not fail at creation — it fails
// later, and silently: the relayer builds headers pinned to whatever client the
// config names, while the contract queries the client it was actually created
// with. When those differ, every ConsensusState lookup misses, ibc-go answers
// NotFound, and wasmd redacts it to "codespace: undefined, code: 1" — an error
// that names neither client. Worse in production: the pinned client is the one
// nobody advances, so it expires while the relayer keeps a different one fresh.
//
// Writing both ids back is what closes that gap; see validateL2SourceEthClient for
// the startup check that catches configs edited by hand.

// writeL2SourceClientIDs persists the created client ids into the l2_to_cosmos
// module whose config.l2_ics26_client_id equals l2ICS26ClientID. Empty ids are
// skipped, so callers can write whichever they have.
func writeL2SourceClientIDs(configPath, l2ICS26ClientID, l2WasmClientID, l1ClientID string) error {
	if l2ICS26ClientID == "" {
		return fmt.Errorf("l2_ics26_client_id is required to locate the l2_to_cosmos module")
	}
	data, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}
	writes := []struct {
		path  []string
		value string
	}{
		{[]string{"l2_wasm_client_id"}, l2WasmClientID},
		{[]string{"rollup_profile", "common", "ethereum_client", "client_id"}, l1ClientID},
	}
	for _, w := range writes {
		if w.value == "" {
			continue
		}
		if data, err = replaceL2SourceMember(data, l2ICS26ClientID, w.path, w.value); err != nil {
			return err
		}
	}
	if err := os.WriteFile(configPath, data, configFilePerm); err != nil {
		return err
	}
	return os.Chmod(configPath, configFilePerm)
}

// replaceL2SourceMember sets config.<path...> on the l2_to_cosmos module whose
// config.l2_ics26_client_id equals l2ICS26ClientID, preserving the surrounding JSON
// formatting. It mirrors replaceConfigMemberForSource (same module classification
// rule, same byte-splicing) but walks a nested path and matches on the L2-side
// ICS26 client id, which is how an l2_to_cosmos module is identified.
func replaceL2SourceMember(data []byte, l2ICS26ClientID string, path []string, value string) ([]byte, error) {
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

		match, configStart, configEnd, err := matchL2Source(data, moduleStart, l2ICS26ClientID)
		if err != nil {
			return nil, err
		}
		if match {
			return setJSONPath(data, configStart, configEnd, path, encodedValue)
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
	return nil, fmt.Errorf("l2_to_cosmos module with l2_ics26_client_id %q not found in config", l2ICS26ClientID)
}

// matchL2Source reports whether the module at moduleStart is the l2_to_cosmos module
// for l2ICS26ClientID, returning its config object range. Classification uses the
// same canonical src_chain/dst_chain rule as loadConfig, so a free-form module name
// does not hide the module from write-back.
func matchL2Source(data []byte, moduleStart int, l2ICS26ClientID string) (bool, int, int, error) {
	name, err := moduleStringMember(data, moduleStart, "name")
	if err != nil {
		return false, 0, 0, err
	}
	srcChain, err := moduleStringMember(data, moduleStart, "src_chain")
	if err != nil {
		return false, 0, 0, err
	}
	dstChain, err := moduleStringMember(data, moduleStart, "dst_chain")
	if err != nil {
		return false, 0, 0, err
	}
	dir, _, cerr := classifyModule(configModule{Name: name, SrcChain: srcChain, DstChain: dstChain})
	if cerr != nil || dir != dirL2ToCosmos {
		return false, 0, 0, nil
	}

	configStart, configEnd, ok, err := findJSONObjectMember(data, moduleStart, "config")
	if err != nil {
		return false, 0, 0, err
	}
	if !ok || configStart >= len(data) || data[configStart] != '{' {
		return false, 0, 0, fmt.Errorf("%s.config is not an object", dirL2ToCosmos)
	}
	id, err := configStringMember(data, configStart, "l2_ics26_client_id")
	if err != nil {
		return false, 0, 0, err
	}
	return id == l2ICS26ClientID, configStart, configEnd, nil
}

// configStringMember reads an optional string member of the config object at
// configStart, returning "" when absent.
func configStringMember(data []byte, configStart int, key string) (string, error) {
	start, end, ok, err := findJSONObjectMember(data, configStart, key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", nil
	}
	var s string
	if err := json.Unmarshal(data[start:end], &s); err != nil {
		return "", fmt.Errorf("parse config %s: %w", key, err)
	}
	return s, nil
}

// setJSONPath sets path (one key per nesting level) to encodedValue inside the
// object spanning [objectStart, objectEnd), splicing bytes so unrelated formatting
// and comments-free JSON layout survive. Missing intermediate objects are created.
func setJSONPath(data []byte, objectStart, objectEnd int, path []string, encodedValue []byte) ([]byte, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty config path")
	}
	key := path[0]
	valueStart, valueEnd, ok, err := findJSONObjectMember(data, objectStart, key)
	if err != nil {
		return nil, err
	}

	if len(path) == 1 {
		if !ok {
			return insertJSONObjectMember(data, objectStart, objectEnd, key, encodedValue)
		}
		out := make([]byte, 0, len(data)-valueEnd+valueStart+len(encodedValue))
		out = append(out, data[:valueStart]...)
		out = append(out, encodedValue...)
		out = append(out, data[valueEnd:]...)
		return out, nil
	}

	if !ok {
		return insertJSONObjectMember(data, objectStart, objectEnd, key, nestedJSONObject(path[1:], encodedValue))
	}
	if valueStart >= len(data) || data[valueStart] != '{' {
		return nil, fmt.Errorf("config path element %q is not an object", key)
	}
	return setJSONPath(data, valueStart, valueEnd, path[1:], encodedValue)
}

// nestedJSONObject builds {"k1":{"k2":...{"kn":value}}} for a path whose parent is
// missing from the config entirely.
func nestedJSONObject(path []string, encodedValue []byte) []byte {
	out := encodedValue
	for i := len(path) - 1; i >= 0; i-- {
		// %q is used rather than json.Marshal so this cannot fail: falling back to a
		// partial value here would splice malformed JSON into the operator's config.
		key := []byte(strconv.Quote(path[i]))
		nested := make([]byte, 0, len(key)+len(out)+4)
		nested = append(nested, '{')
		nested = append(nested, key...)
		nested = append(nested, ':', ' ')
		nested = append(nested, out...)
		nested = append(nested, '}')
		out = nested
	}
	return out
}
