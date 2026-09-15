package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmtp2p "github.com/cometbft/cometbft/p2p"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	"github.com/ethereum/go-ethereum/common"

	"relayer/rpcmock"
)

// One relayer process serves exactly one source → destination pair (03-Relayer
// §2.1 requirement 5). Before this guard, `start` looped over every configured
// module and ran them in one process, which is what made a shared nonce cache and
// a shared loopErrCh necessary in the first place.
func TestValidateSingleRelayPair(t *testing.T) {
	c2e := func(ids ...string) []cosmosToEthConfig {
		out := make([]cosmosToEthConfig, 0, len(ids))
		for _, id := range ids {
			out = append(out, cosmosToEthConfig{ICS26ClientID: id})
		}
		return out
	}
	l2src := func(chains ...string) []l2ToCosmosConfig {
		out := make([]l2ToCosmosConfig, 0, len(chains))
		for _, c := range chains {
			out = append(out, l2ToCosmosConfig{AttestorSrcChain: c})
		}
		return out
	}

	tests := []struct {
		name string
		cfg  appConfig
		// wantNames are the identities the error must name so the operator knows
		// which paths to split apart; empty means the config must be accepted.
		wantNames []string
	}{
		{
			name: "one cosmos_to_eth",
			cfg:  appConfig{CosmosToEthConfigs: c2e("chain-a")},
		},
		{
			name: "one cosmos_to_l2 with its return leg is one pair",
			cfg: appConfig{
				CosmosToL2Configs: c2e("08-wasm-3"),
				L2ToCosmosConfigs: l2src("opstack"),
			},
		},
		{
			name: "forward-only cosmos_to_l2",
			cfg:  appConfig{CosmosToL2Configs: c2e("08-wasm-3")},
		},
		{
			name: "two cosmos_to_eth sources",
			cfg:  appConfig{CosmosToEthConfigs: c2e("chain-a", "chain-b")},
			// Naming both is the point: "config has 2 sources" leaves the operator
			// guessing which two, and file order is not what they edited by.
			wantNames: []string{"chain-a", "chain-b"},
		},
		{
			name:      "two cosmos_to_l2 destinations",
			cfg:       appConfig{CosmosToL2Configs: c2e("08-wasm-3", "08-wasm-4")},
			wantNames: []string{"08-wasm-3", "08-wasm-4"},
		},
		{
			name:      "two l2_to_cosmos sources",
			cfg:       appConfig{L2ToCosmosConfigs: l2src("opstack", "arbitrum")},
			wantNames: []string{"opstack", "arbitrum"},
		},
		{
			name: "cosmos_to_eth alongside cosmos_to_l2 is two pairs",
			cfg: appConfig{
				CosmosToEthConfigs: c2e("chain-a"),
				CosmosToL2Configs:  c2e("08-wasm-3"),
			},
			wantNames: []string{"chain-a", "08-wasm-3"},
		},
		{
			name: "cosmos_to_eth alongside an l2_to_cosmos return leg is two pairs",
			cfg: appConfig{
				CosmosToEthConfigs: c2e("chain-a"),
				L2ToCosmosConfigs:  l2src("opstack"),
			},
			wantNames: []string{"chain-a", "opstack"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateSingleRelayPair(&tt.cfg, "config.json")
			if len(tt.wantNames) == 0 {
				if err != nil {
					t.Fatalf("single-pair config rejected: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("start accepted more than one relay path in one process")
			}
			for _, name := range tt.wantNames {
				if !strings.Contains(err.Error(), name) {
					t.Errorf("error must name path %q so the operator knows what to split; got: %v", name, err)
				}
			}
			// Telling the operator what is wrong without telling them the way out
			// leaves them editing the config at random.
			if !strings.Contains(err.Error(), "--config") {
				t.Errorf("error must show the split command; got: %v", err)
			}
		})
	}
}

// `start` with nothing to relay is a config mistake, not an idle process.
func TestValidateSingleRelayPairRejectsEmptyConfig(t *testing.T) {
	err := validateSingleRelayPair(&appConfig{}, "/etc/relayer/config.json")
	if err == nil {
		t.Fatal("start accepted a config with no relay module")
	}
	if !strings.Contains(err.Error(), "/etc/relayer/config.json") {
		t.Errorf("error must name the config file it read; got: %v", err)
	}
}

// The shipped examples encode the split the guard now enforces: one runnable file
// per path, plus a catalogue that carries every module and is not meant to be
// started. Without this test the two drift silently — the catalogue is the file
// whose name most invites `--config config.example.json`.
func TestShippedExamplesCarryOnePathEach(t *testing.T) {
	t.Parallel()

	for file := range pathExamples {
		t.Run(file, func(t *testing.T) {
			t.Parallel()
			cfg, err := loadConfig(filepath.Join("..", file))
			if err != nil {
				t.Fatalf("loadConfig(%s): %v", file, err)
			}
			if err := validateSingleRelayPair(cfg, file); err != nil {
				t.Fatalf("%s is a runnable example but start rejects it: %v", file, err)
			}
		})
	}

	t.Run("catalogue is not runnable", func(t *testing.T) {
		t.Parallel()
		cfg, err := loadConfig(filepath.Join("..", "config.example.json"))
		if err != nil {
			t.Fatalf("loadConfig(config.example.json): %v", err)
		}
		err = validateSingleRelayPair(cfg, "config.example.json")
		if err == nil {
			t.Fatal("the catalogue lists every path; start must send the operator to a per-path file instead")
		}
		if !strings.Contains(err.Error(), "--config") {
			t.Errorf("the rejection must show the way out; got: %v", err)
		}
	})
}

// --- state directory ---------------------------------------------------------

func TestValidateStateDirAcceptsAWritableDirAndLeavesNothingBehind(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(dir, "state", "recovery-cursors.json"))

	if findings := validateStateDir(); len(findings) != 0 {
		t.Fatalf("expected no findings for a writable dir, got %+v", findings)
	}

	// The probe must not survive. A leftover file in the state directory is not
	// merely untidy: the next run's O_EXCL probe would fail on it and report a
	// writable directory as unwritable.
	entries, err := os.ReadDir(filepath.Join(dir, "state"))
	if err != nil {
		t.Fatalf("read state dir: %v", err)
	}
	if len(entries) != 0 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("probe left files behind: %v", names)
	}
}

func TestValidateStateDirRejectsADirectoryItCannotWrite(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignores the mode bits this test relies on")
	}
	parent := t.TempDir()
	readOnly := filepath.Join(parent, "ro")
	if err := os.Mkdir(readOnly, 0o500); err != nil {
		t.Fatalf("create read-only dir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(readOnly, 0o700) })

	t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(readOnly, "nested", "recovery-cursors.json"))

	findings := validateStateDir()
	if len(findings) != 1 {
		t.Fatalf("expected exactly one finding, got %+v", findings)
	}
	if findings[0].key != "RELAYER_RECOVERY_STATE_FILE" {
		t.Fatalf("finding names %q, want the env var an operator would fix", findings[0].key)
	}
}

// --- cosmos chain id ---------------------------------------------------------

func TestValidateCosmosChainIDFailsBeforeDialingWhenTheVariableIsUnset(t *testing.T) {
	node := rpcmock.NewCosmosNode(t)
	node.Handle("status", func(json.RawMessage) (any, error) {
		return statusResult("some-chain"), nil
	})
	t.Setenv("COSMOS_CHAIN_ID", "")

	findings := validateCosmosChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{TmRpcUrl: node.URL()}},
	})

	if len(findings) != 1 || findings[0].key != "COSMOS_CHAIN_ID" {
		t.Fatalf("expected one COSMOS_CHAIN_ID finding, got %+v", findings)
	}
	// Nothing to compare against means nothing to ask, and asking anyway would
	// make an unset variable look like an endpoint problem in the log.
	if calls := node.Calls("status"); calls != 0 {
		t.Fatalf("queried the node %d times with no declared chain id; want 0", calls)
	}
}

func TestValidateCosmosChainIDReportsADefiniteDisagreement(t *testing.T) {
	node := rpcmock.NewCosmosNode(t)
	node.Handle("status", func(json.RawMessage) (any, error) {
		return statusResult("served-chain"), nil
	})
	t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

	findings := validateCosmosChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{TmRpcUrl: node.URL()}},
	})

	if len(findings) != 1 {
		t.Fatalf("expected one finding, got %+v", findings)
	}
	// Both halves must appear: an operator has to know which of the two to fix,
	// and a message naming only one of them cannot tell them.
	for _, want := range []string{"declared-chain", "served-chain", node.URL()} {
		if !strings.Contains(findings[0].detail, want) {
			t.Fatalf("finding %q does not mention %q", findings[0].detail, want)
		}
	}
}

func TestValidateCosmosChainIDAcceptsAMatch(t *testing.T) {
	node := rpcmock.NewCosmosNode(t)
	node.Handle("status", func(json.RawMessage) (any, error) {
		return statusResult("test-ibc-eth"), nil
	})
	t.Setenv("COSMOS_CHAIN_ID", "test-ibc-eth")

	if findings := validateCosmosChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{TmRpcUrl: node.URL()}},
	}); len(findings) != 0 {
		t.Fatalf("expected no findings for a matching chain id, got %+v", findings)
	}
}

func TestValidateCosmosChainIDToleratesAnEndpointThatDoesNotAnswer(t *testing.T) {
	// A node with no "status" handler answers "method not found" — the shape of
	// an endpoint that is reachable but useless. Startup must survive it: an
	// unavailable node is an operations problem, and the relay loops retry.
	node := rpcmock.NewCosmosNode(t)
	t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

	if findings := validateCosmosChainID(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{TmRpcUrl: node.URL()}},
	}); len(findings) != 0 {
		t.Fatalf("a non-answering endpoint must not fail startup, got %+v", findings)
	}
}

func TestCosmosEndpointsAreGatheredFromEveryModuleListWithoutDuplicates(t *testing.T) {
	cfg := &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{TmRpcUrl: "http://a"}, {TmRpcUrl: "http://b"}},
		CosmosToL2Configs:  []cosmosToEthConfig{{TmRpcUrl: "http://a"}},
		L2ToCosmosConfigs:  []l2ToCosmosConfig{{TmRpcUrl: "http://c"}, {TmRpcUrl: ""}},
	}
	got := cosmosEndpointsInConfig(cfg)
	want := []string{"http://a", "http://b", "http://c"}
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Fatalf("gathered %v, want %v", got, want)
	}
}

// --- EVM deployments ---------------------------------------------------------

// deployedCode is any non-empty code answer; the check only asks whether there
// is code at all, never what it is.
const deployedCode = "0x60006000"

func evmConfigFor(node *rpcmock.EVMNode) cosmosToEthConfig {
	return cosmosToEthConfig{
		EthRpcUrl:         node.URL(),
		ICS26Address:      "0x1111111111111111111111111111111111111111",
		ICS26ClientID:     "client-0",
		SpectreClient:     "0x2222222222222222222222222222222222222222",
		SignatureVerifier: "0x3333333333333333333333333333333333333333",
		Membership:        "0x4444444444444444444444444444444444444444",
		Misbehaviour:      "0x5555555555555555555555555555555555555555",
		UpdateClient:      "0x6666666666666666666666666666666666666666",
	}
}

func deployAll(node *rpcmock.EVMNode, c cosmosToEthConfig) {
	for _, a := range []string{
		c.ICS26Address, c.SpectreClient, c.SignatureVerifier,
		c.Membership, c.Misbehaviour, c.UpdateClient,
	} {
		node.SetCode(common.HexToAddress(a), deployedCode)
	}
	// getClient's return value: one ABI-encoded non-zero address.
	node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))
}

func TestValidateEVMDeploymentsReportsEveryAddressThatHoldsNoCode(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	// Only the router is deployed, so the other five are the findings — and the
	// router being present is what lets the client-id check run at all.
	node.SetCode(common.HexToAddress(cfg.ICS26Address), deployedCode)
	node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))

	findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	})

	got := map[string]bool{}
	for _, f := range findings {
		got[f.key] = true
	}
	for _, key := range []string{
		"cosmos_to_eth.spectre_client", "cosmos_to_eth.signature_verifier",
		"cosmos_to_eth.membership", "cosmos_to_eth.misbehaviour",
		"cosmos_to_eth.update_client",
	} {
		if !got[key] {
			t.Fatalf("no finding for %s; got %+v", key, findings)
		}
	}
	if got["cosmos_to_eth.ics26_address"] {
		t.Fatalf("reported the deployed router as code-less: %+v", findings)
	}
	// This test used a presence map, which is why it could not see the router
	// check handing back the findings it was given -- five problems reported as
	// ten. Count as well as check.
	if len(findings) != len(got) {
		t.Fatalf("%d findings for %d distinct keys; some are reported twice: %+v",
			len(findings), len(got), findings)
	}
}

// validateRouterClientID receives the findings so far so it can stay quiet when
// the router address is itself a finding. It must not hand them back: the caller
// appends its return value, so returning the accumulator reports everything
// twice. Found in review -- a module with five missing contracts came out as ten
// problems.
func TestFindingsAreNotReportedTwice(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	// Router deployed (so the client-id check runs) and nothing else, which is
	// the shape that accumulates findings BEFORE the router check.
	node.SetCode(common.HexToAddress(cfg.ICS26Address), deployedCode)
	node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))

	findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	})

	seen := map[string]int{}
	for _, f := range findings {
		seen[f.key]++
	}
	for key, n := range seen {
		if n > 1 {
			t.Fatalf("%q reported %d times (%d findings, %d distinct): an operator "+
				"reading this cannot tell how many problems there actually are",
				key, n, len(findings), len(seen))
		}
	}
}

// The mirror of TestValidateEVMDeploymentsValidatesTheIDStartupWillActuallyUse,
// and the bug is its exact inverse: buildCosmosToL2Dest reads ICS26ClientID
// straight from the JSON and never consults the environment, so applying the
// override to a cosmos_to_l2 entry validates an id that path will never use.
// A stale value in .env then fails a good config, and a correct one hides a bad
// JSON id. Validation must read whatever the BUILDER reads, per path.
func TestCosmosToL2ValidatesTheConfigIDNotTheEnvOverride(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	cfg.ICS26ClientID = "client-from-json"
	deployAll(node, cfg)
	// Every getClient reverts, so whichever id is asked about becomes a finding
	// naming that id -- which is how we see WHICH id was validated.
	node.FailCall(errors.New("execution reverted"))
	t.Setenv("ICS26_CLIENT_ID", "stale-from-dotenv")

	findings := validateEVMDeployments(context.Background(),
		&appConfig{CosmosToL2Configs: []cosmosToEthConfig{cfg}})

	if len(findings) == 0 {
		t.Fatal("no finding at all; the cosmos_to_l2 entry was not validated")
	}
	for _, f := range findings {
		if strings.Contains(f.detail, "stale-from-dotenv") {
			t.Fatalf("validated the env override on a cosmos_to_l2 entry: %s\n"+
				"buildCosmosToL2Dest never reads it, so this rejects a config the relay "+
				"would have run fine", f.detail)
		}
	}
	for _, f := range findings {
		if strings.Contains(f.detail, "client-from-json") {
			return
		}
	}
	t.Fatalf("no finding names the JSON client id, so nothing proves which id was checked: %+v", findings)
}

func TestValidateEVMDeploymentsAcceptsAFullyDeployedModule(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	deployAll(node, cfg)

	if findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	}); len(findings) != 0 {
		t.Fatalf("expected no findings for a deployed module, got %+v", findings)
	}
}

func TestValidateEVMDeploymentsReportsAClientIDTheRouterDoesNotKnow(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	deployAll(node, cfg)
	// getClient reverts with IBCClientNotFound for an unregistered id
	// (contracts/utils/ICS02ClientUpgradeable.sol:78-80).
	node.FailCall(errors.New("execution reverted"))

	findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	})

	if len(findings) != 1 || findings[0].key != "cosmos_to_eth.ics26_client_id" {
		t.Fatalf("expected one ics26_client_id finding, got %+v", findings)
	}
	if !strings.Contains(findings[0].detail, "client-0") {
		t.Fatalf("finding %q does not name the client id", findings[0].detail)
	}
}

func TestValidateEVMDeploymentsToleratesATransportErrorFromTheRouterCall(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	deployAll(node, cfg)
	// Not a revert. The node did not execute the call, so this says nothing
	// about the client id and must not stop startup.
	node.FailCall(errors.New("connection reset by peer"))

	if findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	}); len(findings) != 0 {
		t.Fatalf("a transport error must not fail startup, got %+v", findings)
	}
}

func TestValidateEVMDeploymentsDoesNotReportTheClientIDWhenTheRouterItselfIsMissing(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := cosmosToEthConfig{
		EthRpcUrl:     node.URL(),
		ICS26Address:  "0x1111111111111111111111111111111111111111",
		ICS26ClientID: "client-0",
	}
	// Nothing deployed, and getClient would revert. "The router has no code" and
	// "the router does not know this client" are then the same fact, and
	// reporting both sends an operator looking for a second problem.
	node.FailCall(errors.New("execution reverted"))

	findings := validateEVMDeployments(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	})

	if len(findings) != 1 || findings[0].key != "cosmos_to_eth.ics26_address" {
		t.Fatalf("expected only the missing-router finding, got %+v", findings)
	}
}

func TestValidateEVMDeploymentsToleratesAnUnreachableEndpoint(t *testing.T) {
	cfg := cosmosToEthConfig{
		// Reserved TEST-NET-1 address (RFC 5737): routable syntax, no listener.
		EthRpcUrl:     "http://192.0.2.1:8545",
		ICS26Address:  "0x1111111111111111111111111111111111111111",
		ICS26ClientID: "client-0",
	}
	ctx, cancel := context.WithTimeout(context.Background(), configProbeTimeout)
	defer cancel()

	if findings := validateEVMDeployments(ctx, &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	}); len(findings) != 0 {
		t.Fatalf("an unreachable endpoint must not fail startup, got %+v", findings)
	}
}

func TestEVMDeploymentsSkipUnsetAndMalformedAddresses(t *testing.T) {
	deployments := evmDeploymentsInConfig(&appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{
			EthRpcUrl:     "http://eth",
			ICS26Address:  "0x1111111111111111111111111111111111111111",
			SpectreClient: "",            // unset
			Membership:    "not-an-addr", // malformed
		}},
		CosmosToL2Configs: []cosmosToEthConfig{{
			EthRpcUrl:    "http://l2",
			ICS26Address: "0x2222222222222222222222222222222222222222",
		}},
	})

	if len(deployments) != 2 {
		t.Fatalf("expected one deployment per EVM endpoint, got %d", len(deployments))
	}
	if len(deployments[0].addresses) != 1 {
		t.Fatalf("unset/malformed addresses must not be probed, got %v", deployments[0].addresses)
	}
	if deployments[0].label != "cosmos_to_eth" || deployments[1].label != "cosmos_to_l2" {
		t.Fatalf("labels are %q and %q; an operator greps the config by these",
			deployments[0].label, deployments[1].label)
	}
}

func TestValidateEVMDeploymentsCoversTheL2ReturnLegToo(t *testing.T) {
	// The mirror of TestValidateEVMDeploymentsReportsEveryAddressThatHoldsNoCode. An
	// l2_to_cosmos module reaches a router and a client id just like the outbound
	// legs; it only names them differently, and checking one direction while
	// leaving the other exposed is the asymmetry this repo keeps rediscovering.
	node := rpcmock.NewEVMNode(t)
	router := "0x7777777777777777777777777777777777777777"
	cfg := &appConfig{
		L2ToCosmosConfigs: []l2ToCosmosConfig{{
			L2RpcUrl:        node.URL(),
			L2ICS26ClientID: "l2-client-0",
			RollupProfile: json.RawMessage(
				`{"common":{"l2_chain_id":1,"l2_router":"` + router + `"}}`),
		}},
	}

	// Undeployed router AND a reverting getClient. Both facts have the same
	// cause, so only the router may be reported — and the dedupe that ensures
	// it has to match on the L2 leg's own key name, not the outbound one.
	node.FailCall(jsonRPCError{msg: "execution reverted", data: "0x"})
	findings := validateEVMDeployments(context.Background(), cfg)
	if len(findings) != 1 || findings[0].key != "l2_to_cosmos.rollup_profile.common.l2_router" {
		t.Fatalf("expected only the undeployed L2 router to be reported, got %+v", findings)
	}
	node.FailCall(nil)

	// Deployed: the router now resolves and there is nothing to report.
	node.SetCode(common.HexToAddress(router), deployedCode)
	node.SetCallResult("0x000000000000000000000000" + strings.Repeat("ab", 20))
	if findings := validateEVMDeployments(context.Background(), cfg); len(findings) != 0 {
		t.Fatalf("expected no findings once the L2 router is deployed, got %+v", findings)
	}

	// And an unregistered client id on that router is reported against the L2
	// key, not the outbound one.
	node.FailCall(jsonRPCError{msg: "execution reverted", data: "0x"})
	findings = validateEVMDeployments(context.Background(), cfg)
	if len(findings) != 1 || findings[0].key != "l2_to_cosmos.l2_ics26_client_id" {
		t.Fatalf("expected an l2_ics26_client_id finding, got %+v", findings)
	}
}

func TestEVMDeploymentsSkipAModuleWithNoRPC(t *testing.T) {
	if d := evmDeploymentsInConfig(&appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{{ICS26Address: "0x1111111111111111111111111111111111111111"}},
	}); len(d) != 0 {
		t.Fatalf("a module with no eth_rpc_url has nothing to probe, got %+v", d)
	}
}

// --- revert classification ---------------------------------------------------

func TestIsEVMRevertSeparatesExecutionFromTransport(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"geth revert text", errors.New("execution reverted"), true},
		{"revert text with a reason", errors.New("execution reverted: IBCClientNotFound"), true},
		{"mixed case", errors.New("Execution Reverted"), true},
		{"wrapped revert", fmt.Errorf("call getClient: %w", errors.New("execution reverted")), true},
		{"dial failure", errors.New("dial tcp 127.0.0.1:8545: connect: connection refused"), false},
		{"timeout", context.DeadlineExceeded, false},
		{"cancelled", context.Canceled, false},

		// The rpc.DataError branch. go-ethereum wraps EVERY JSON-RPC error in a
		// type satisfying it, so the interface alone cannot separate a revert
		// from a node complaint — only the presence of data can.
		{"revert with data", jsonRPCError{msg: "execution reverted", data: "0x08c379a0"}, true},
		{"pruned node", jsonRPCError{msg: "header not found"}, false},
		{"rate limited", jsonRPCError{msg: "your app has exceeded its compute units"}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isEVMRevert(tc.err); got != tc.want {
				t.Fatalf("isEVMRevert(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// The path is what the log prefix carries, so it has to be derivable from config
// alone and stable for the life of the process. A1 makes that true by refusing a
// config with more than one path.
func TestRelayPathID(t *testing.T) {
	for _, tc := range []struct {
		name string
		cfg  *appConfig
		want string
	}{
		{
			name: "cosmos to eth carries the source client id",
			cfg:  &appConfig{CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "cosmoshub-1"}}},
			want: "cosmos<->eth/cosmoshub-1",
		},
		{
			// The rollup names itself: an operator running an OP and an Arbitrum
			// path needs to tell the two apart, and "l2" does not.
			name: "cosmos to l2 names the rollup from the return leg",
			cfg: &appConfig{
				CosmosToL2Configs: []cosmosToEthConfig{{ICS26ClientID: "cosmoshub-1"}},
				L2ToCosmosConfigs: []l2ToCosmosConfig{{AttestorSrcChain: "arbitrum"}},
			},
			want: "cosmos<->arbitrum/cosmoshub-1",
		},
		{
			name: "cosmos to l2 without a return leg still identifies the path",
			cfg:  &appConfig{CosmosToL2Configs: []cosmosToEthConfig{{ICS26ClientID: "cosmoshub-1"}}},
			want: "cosmos<->l2/cosmoshub-1",
		},
		{
			name: "return leg only",
			cfg:  &appConfig{L2ToCosmosConfigs: []l2ToCosmosConfig{{AttestorSrcChain: "opstack"}}},
			want: "opstack->cosmos",
		},
		{
			// Never empty: an empty prefix would silently drop the identity from
			// every line rather than being obviously wrong.
			name: "nothing configured still yields a label",
			cfg:  &appConfig{},
			want: "relayer",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := relayPathID(tc.cfg); got != tc.want {
				t.Fatalf("relayPathID = %q, want %q", got, tc.want)
			}
		})
	}
}

// --- aggregation -------------------------------------------------------------

func TestValidateStartupConfigReportsEveryProblemInASingleMessage(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(dir, "recovery-cursors.json"))
	t.Setenv("COSMOS_CHAIN_ID", "declared-chain")

	cosmos := rpcmock.NewCosmosNode(t)
	cosmos.Handle("status", func(json.RawMessage) (any, error) {
		return statusResult("served-chain"), nil
	})
	eth := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(eth)
	cfg.TmRpcUrl = cosmos.URL()

	err := validateStartupConfig(context.Background(), &appConfig{
		CosmosToEthConfigs: []cosmosToEthConfig{cfg},
	})
	if err == nil {
		t.Fatal("expected config validation to fail")
	}

	// One run must diagnose the whole config. Failing on the first problem turns
	// a three-key mistake into three restarts.
	msg := err.Error()
	for _, want := range []string{
		"COSMOS_CHAIN_ID", "served-chain",
		"cosmos_to_eth.ics26_address", "cosmos_to_eth.update_client",
	} {
		if !strings.Contains(msg, want) {
			t.Fatalf("message does not mention %q:\n%s", want, msg)
		}
	}
}

func TestValidateStartupConfigPassesOnAConfigWithNothingToCheck(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("RELAYER_RECOVERY_STATE_FILE", filepath.Join(dir, "recovery-cursors.json"))

	if err := validateStartupConfig(context.Background(), &appConfig{}); err != nil {
		t.Fatalf("an empty config has nothing to disagree with: %v", err)
	}
}

func TestSortedKeysIsStable(t *testing.T) {
	got := sortedKeys(map[string]string{"c": "", "a": "", "b": ""})
	if fmt.Sprint(got) != fmt.Sprint([]string{"a", "b", "c"}) {
		t.Fatalf("sortedKeys returned %v", got)
	}
}

// jsonRPCError mirrors go-ethereum's jsonError: it satisfies rpc.DataError for
// every node-returned error, carrying data only when the call actually reverted
// with a reason.
type jsonRPCError struct {
	msg  string
	data any
}

func (e jsonRPCError) Error() string  { return e.msg }
func (e jsonRPCError) ErrorData() any { return e.data }
func (e jsonRPCError) ErrorCode() int { return -32000 }

// statusResult is the minimal /status reply: validation reads only the
// network name, so the stub models only that.
func statusResult(network string) *coretypes.ResultStatus {
	return &coretypes.ResultStatus{
		NodeInfo: cmtp2p.DefaultNodeInfo{Network: network},
	}
}

// Reported by @DongLieu on #467. getClient succeeding proves only that the
// router knows SOME client under that id. It says nothing about whether that
// client is the spectre_client this config deploys against, nor whether its
// counterparty is the Cosmos client this config relays to. A wrong-but-
// registered id passes startup, and then every update and proof targets one
// client while the router verifies another.
//
// routerWiringIsReusable already decides this and is already tested; only the
// RPC to feed it was missing here. Per-selector call results are what let the
// stub answer getClient and getCounterparty differently, which SetCallResult
// alone could not express.
func TestValidateEVMDeploymentsRejectsARouterWiredToAnotherClient(t *testing.T) {
	const (
		getClientSelector       = "7eb78932"
		getCounterpartySelector = "b0777bfa"
	)
	abiAddress := func(hexAddr string) string {
		return "0x000000000000000000000000" + strings.TrimPrefix(hexAddr, "0x")
	}

	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	deployAll(node, cfg)
	// The counterparty the router reports, ABI-encoded as (string,bytes[]): a
	// dynamic head, then the Cosmos client id.
	node.SetCallResultFor(getCounterpartySelector, abiCounterparty("08-wasm-0"))

	t.Run("the client id points at a different light client", func(t *testing.T) {
		node.SetCallResultFor(getClientSelector, abiAddress("0x9999999999999999999999999999999999999999"))
		local := cfg
		local.CosmosWasmClientID = "08-wasm-0"
		findings := validateEVMDeployments(context.Background(), &appConfig{CosmosToEthConfigs: []cosmosToEthConfig{local}})
		if len(findings) == 0 {
			t.Fatal("a router pointing that client id at another light client was accepted; " +
				"updates would target spectre_client while the router verifies something else")
		}
	})

	t.Run("the counterparty is a different Cosmos client", func(t *testing.T) {
		node.SetCallResultFor(getClientSelector, abiAddress(cfg.SpectreClient))
		local := cfg
		local.CosmosWasmClientID = "08-wasm-7" // the router says 08-wasm-0
		findings := validateEVMDeployments(context.Background(), &appConfig{CosmosToEthConfigs: []cosmosToEthConfig{local}})
		if len(findings) == 0 {
			t.Fatal("a router counterparty-wired to another Cosmos client was accepted; " +
				"recvPacket would reject every packet's counterparty")
		}
	})

	t.Run("correctly wired passes", func(t *testing.T) {
		node.SetCallResultFor(getClientSelector, abiAddress(cfg.SpectreClient))
		local := cfg
		local.CosmosWasmClientID = "08-wasm-0"
		if findings := validateEVMDeployments(context.Background(), &appConfig{CosmosToEthConfigs: []cosmosToEthConfig{local}}); len(findings) != 0 {
			t.Fatalf("a correctly wired router was rejected: %v", findings)
		}
	})
}

// Reported by @DongLieu on #467. build_source.go:102 reads ICS26_CLIENT_ID over
// the config, so validating the raw config value proves something startup then
// does not use: a stale id in .env passes here and fails at the first packet,
// which is exactly the failure this validation exists to move forward.
func TestValidateEVMDeploymentsValidatesTheIDStartupWillActuallyUse(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	cfg := evmConfigFor(node)
	deployAll(node, cfg)
	// The router knows only the id in the config file: an id it does not know
	// reverts, which is how getClient reports IBCClientNotFound.
	node.FailCall(errors.New("execution reverted"))

	t.Setenv("ICS26_CLIENT_ID", "stale-from-dotenv")
	findings := validateEVMDeployments(context.Background(), &appConfig{CosmosToEthConfigs: []cosmosToEthConfig{cfg}})
	if len(findings) == 0 {
		t.Fatal("the override was ignored and the config value validated instead; a stale " +
			"ICS26_CLIENT_ID in .env would pass startup and fail at the first packet")
	}
	for _, f := range findings {
		if strings.Contains(f.detail, "stale-from-dotenv") {
			return
		}
	}
	t.Fatalf("no finding names the effective client id: %v", findings)
}

// abiCounterparty encodes what getCounterparty returns: a single dynamic STRUCT
// {string clientId; bytes[] merklePrefix}. A dynamic struct return is itself
// referenced by an offset, so the head is one word pointing at the tuple, and
// only then do the tuple's own two offsets begin. Getting that outer word wrong
// is why the first version of this test could not decode.
func abiCounterparty(clientID string) string {
	word := func(v int) string { return fmt.Sprintf("%064x", v) }
	padded := hex.EncodeToString([]byte(clientID))
	for len(padded)%64 != 0 {
		padded += "0"
	}
	tuple := word(0x40) + // clientId offset, relative to the tuple
		word(0x60+len(padded)/2) + // merklePrefix offset, past the string
		word(len(clientID)) + padded +
		word(0) // empty bytes[]
	return "0x" + word(0x20) + tuple
}

// Two processes relaying cosmos->eth from different source clients are a valid
// deployment, and their lines are merged by journald or Loki with no file
// boundary left. The prefix is the only thing that tells them apart, so it must
// actually differ.
func TestRelayPathIDDistinguishesTwoCosmosSources(t *testing.T) {
	a := relayPathID(&appConfig{CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "cosmoshub-1"}}})
	b := relayPathID(&appConfig{CosmosToEthConfigs: []cosmosToEthConfig{{ICS26ClientID: "osmosis-1"}}})
	if a == b {
		t.Fatalf("both sources produced %q; merged logs cannot be told apart", a)
	}
}
