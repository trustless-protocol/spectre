package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// chainIDServer answers eth_chainId with the given hex-quantity id.
func chainIDServer(t *testing.T, hexID string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":1,"result":%q}`, hexID)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func profileWithChainID(id uint64) json.RawMessage {
	return json.RawMessage(fmt.Sprintf(`{
		"common": {
			"l2_chain_id": %d,
			"l2_router": "0x6c4729db04a00b4a76d854df35980e1646dff4ba",
			"commitment_slot": "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600",
			"profile_version": "arbitrum_attestor_v1",
			"l2_header_fork": "london",
			"attestor_public_key": "0x1111111111111111111111111111111111111111111111111111111111111111",
			"attestation_head": "safe"
		}
	}`, id))
}

// TestValidateL2ChainIDsRejectsAMismatch is the case this check exists for.
//
// It is exactly what happened on a real bring-up: the profile still named the
// local devnet (412346) while l2_rpc_url pointed at Arbitrum Sepolia (421614),
// and everything relayed anyway — transactions take their chain id from
// eth_chainId, and nothing on either side of the bridge reads l2_chain_id. The
// operator only noticed because the transactions turned up on a chain they had
// not configured.
func TestValidateL2ChainIDsRejectsAMismatch(t *testing.T) {
	t.Parallel()

	srv := chainIDServer(t, "0x66eee") // 421614, Arbitrum Sepolia
	cfg := l2ToCosmosConfig{
		AttestorSrcChain: "arbitrum-sepolia",
		L2RpcUrl:         srv.URL,
		RollupProfile:    profileWithChainID(412346), // the devnet
	}

	err := validateL2ChainIDs(context.Background(), []l2ToCosmosConfig{cfg})
	if err == nil {
		t.Fatal("a profile naming a different chain than the RPC serves must not start")
	}
	// Both numbers must appear: the operator has to know which one to change.
	for _, want := range []string{"412346", "421614", "arbitrum-sepolia"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error must name %s so the mismatch is actionable, got: %v", want, err)
		}
	}
}

func TestValidateL2ChainIDsAcceptsAMatch(t *testing.T) {
	t.Parallel()

	srv := chainIDServer(t, "0x66eee") // 421614
	cfg := l2ToCosmosConfig{
		AttestorSrcChain: "arbitrum-sepolia",
		L2RpcUrl:         srv.URL,
		RollupProfile:    profileWithChainID(421614),
	}

	if err := validateL2ChainIDs(context.Background(), []l2ToCosmosConfig{cfg}); err != nil {
		t.Fatalf("a correct config must start: %v", err)
	}
}

// TestValidateL2ChainIDsToleratesAnUnreachableEndpoint separates "wrong" from
// "unknown".
//
// The comparison reports a briefly down node as unknown rather than as a
// profile mismatch. The EVM signer lock remains a separate startup gate and
// refuses to submit without an endpoint-reported chain id.
func TestValidateL2ChainIDsToleratesAnUnreachableEndpoint(t *testing.T) {
	t.Parallel()

	cfg := l2ToCosmosConfig{
		AttestorSrcChain: "arbitrum-sepolia",
		L2RpcUrl:         "http://127.0.0.1:1", // nothing listening
		RollupProfile:    profileWithChainID(421614),
	}

	if err := validateL2ChainIDs(context.Background(), []l2ToCosmosConfig{cfg}); err != nil {
		t.Fatalf("an unreachable endpoint must not be reported as a profile mismatch: %v", err)
	}
}

// A profile with no l2_chain_id IS fatal, unlike an unreachable endpoint: the
// Rust verifier declares the field non-optional under deny_unknown_fields, so
// such a profile cannot create a client at all.
func TestValidateL2ChainIDsRejectsAMissingChainID(t *testing.T) {
	t.Parallel()

	srv := chainIDServer(t, "0x66eee")
	cfg := l2ToCosmosConfig{
		AttestorSrcChain: "arbitrum-sepolia",
		L2RpcUrl:         srv.URL,
		RollupProfile: json.RawMessage(`{
			"common": {"l2_router": "0x6c4729db04a00b4a76d854df35980e1646dff4ba"}
		}`),
	}

	err := validateL2ChainIDs(context.Background(), []l2ToCosmosConfig{cfg})
	if err == nil || !strings.Contains(err.Error(), "l2_chain_id") {
		t.Fatalf("a profile without l2_chain_id must be rejected by name, got: %v", err)
	}
}

// Every source is checked, not just the first — the same shape of bug as the
// eth_ws_url guard that only looked at c2eList[0] (#298).
func TestValidateL2ChainIDsChecksEverySource(t *testing.T) {
	t.Parallel()

	good := chainIDServer(t, "0x66eee") // 421614
	bad := chainIDServer(t, "0x14a34")  // 84532, Base Sepolia

	sources := []l2ToCosmosConfig{
		{AttestorSrcChain: "arbitrum-sepolia", L2RpcUrl: good.URL, RollupProfile: profileWithChainID(421614)},
		{AttestorSrcChain: "base-sepolia", L2RpcUrl: bad.URL, RollupProfile: profileWithChainID(8453)},
	}

	err := validateL2ChainIDs(context.Background(), sources)
	if err == nil {
		t.Fatal("a mismatch on the second source must still be caught")
	}
	if !strings.Contains(err.Error(), "base-sepolia") {
		t.Fatalf("error must name the offending source, got: %v", err)
	}
	if strings.Contains(err.Error(), "arbitrum-sepolia") {
		t.Fatalf("error must not blame the correctly-configured source, got: %v", err)
	}
}

// TestValidateL2ChainIDsRejectsAnOutOfRangeAnswer is the hole the first five
// tests left open, reported by @neitdung and confirmed by @DongLieu.
//
// probeL2ChainID used to return (0, false) — the "endpoint did not answer"
// signal — for any chain id above the uint64 range. The endpoint HAD answered,
// with a value that cannot possibly equal the declared id (the profile field is
// uint64), yet the caller logged "could not read eth_chainId ... skipping" and
// let startup proceed. That routed a definite disagreement into the "unknown,
// tolerate" branch, which is the opposite of the rule this check enforces.
func TestValidateL2ChainIDsRejectsAnOutOfRangeAnswer(t *testing.T) {
	t.Parallel()

	srv := chainIDServer(t, "0x10000000000000000") // 2^64, one past uint64
	cfg := l2ToCosmosConfig{
		AttestorSrcChain: "repro",
		L2RpcUrl:         srv.URL,
		RollupProfile:    profileWithChainID(1),
	}

	err := validateL2ChainIDs(context.Background(), []l2ToCosmosConfig{cfg})
	if err == nil {
		t.Fatal("an answered-but-unmatchable chain id must be a mismatch, not a skipped check")
	}
	if !strings.Contains(err.Error(), "18446744073709551616") {
		t.Fatalf("the error must report the value the RPC actually served: %v", err)
	}
}

// TestVerifyL2ChainIDNoAnswerPolicyDiffersByCaller pins the deliberate asymmetry
// between the two profile-check callers: start reports an unreachable endpoint
// as unverified, while create-clients refuses because it is about to commit the
// profile into a client that cannot be repaired. start's EVM signer lock remains
// a separate gate and will refuse to submit with no chain id.
func TestVerifyL2ChainIDNoAnswerPolicyDiffersByCaller(t *testing.T) {
	t.Parallel()

	// Nothing is listening here, so the probe cannot get an answer.
	const dead = "http://127.0.0.1:1"

	t.Run("start tolerates", func(t *testing.T) {
		t.Parallel()
		verified, err := verifyL2ChainID(context.Background(), l2ChainIDCheck{
			label: "l2_to_cosmos config", rpcURL: dead, want: 1, tolerateNoAnswer: true,
		})
		if err != nil {
			t.Fatalf("an unreachable endpoint must remain unverified rather than mismatched: %v", err)
		}
		if verified {
			t.Fatal("a tolerated no-answer must report verified=false; " +
				"reporting true would let the caller log a check that never ran")
		}
	})

	t.Run("create-clients refuses", func(t *testing.T) {
		t.Parallel()
		_, err := verifyL2ChainID(context.Background(), l2ChainIDCheck{
			label: "l2-config", rpcURL: dead, want: 1,
		})
		if err == nil {
			t.Fatal("create-clients must not commit an unverified profile on-chain")
		}
		if !strings.Contains(err.Error(), "cannot be verified") {
			t.Fatalf("the error must say the check could not run, not that the chain mismatched: %v", err)
		}
	})
}

// TestVerifyL2ChainIDRejectsAMismatchForCreateClients covers the create-clients
// side of the durable-damage case: the profile is committed into the client
// verbatim, so a wrong chain id there survives every later check.
func TestVerifyL2ChainIDRejectsAMismatchForCreateClients(t *testing.T) {
	t.Parallel()

	srv := chainIDServer(t, "0x66eee") // 421614
	_, err := validateL2ClientChainID(context.Background(), &l2ClientConfig{
		L2RPCURL:      srv.URL,
		RollupProfile: profileWithChainID(412346), // the devnet
	})
	if err == nil {
		t.Fatal("create-clients must reject a profile naming a different chain than its RPC serves")
	}
	for _, want := range []string{"412346", "421614", "l2-config"} {
		if !strings.Contains(err.Error(), want) {
			t.Fatalf("error must contain %q: %v", want, err)
		}
	}
}

// TestValidateL2ChainIDsDoesNotClaimAnUnrunCheck pins the operator-facing half.
//
// The skip branch logged "skipping the l2_chain_id check", and then the success
// line ran unconditionally right after it — so an unreachable endpoint produced
// "skipping" immediately followed by "l2_chain_id verified", which reads as the
// check having passed. The two are now mutually exclusive by construction: the
// caller branches on whether a comparison actually happened, rather than
// remembering to skip a log line.
func TestValidateL2ChainIDsDoesNotClaimAnUnrunCheck(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	// Nothing is listening, so the profile comparison gets no answer. The later
	// EVM signer lock will stop startup rather than run without a nonce guard.
	sources := []l2ToCosmosConfig{{
		AttestorSrcChain: "unreachable",
		L2RpcUrl:         "http://127.0.0.1:1",
		RollupProfile:    json.RawMessage(`{"common":{"l2_chain_id":421614}}`),
	}}
	if err := validateL2ChainIDs(context.Background(), sources); err != nil {
		t.Fatalf("an unreachable endpoint must remain unverified rather than mismatched: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "skipping the l2_chain_id check") {
		t.Fatalf("the skip must be reported; log was:\n%s", out)
	}
	if strings.Contains(out, "l2_chain_id verified") {
		t.Fatalf("the log claims a check that never ran; an operator reading this would "+
			"believe the declared chain id was confirmed. log was:\n%s", out)
	}
}

// A cancelled context must stop startup, not look like an unreachable node.
//
// The previous commit made the probe cancellable by passing runCtx instead of
// cmd.Context(). But a cancelled probe returns the same nil as a node that did
// not answer, and `start` tolerates a no-answer — so Ctrl-C made the check pass
// and startup walked on into the rest of its initialisation, to fail later on
// something unrelated. Threading the signal context in is only half the fix.
func TestVerifyL2ChainIDDoesNotTolerateCancellation(t *testing.T) {
	t.Parallel()

	stdCtx, cancel := context.WithCancel(context.Background())
	cancel()

	// tolerateNoAnswer is the `start` policy — the one branch that could swallow it.
	verified, err := verifyL2ChainID(stdCtx, l2ChainIDCheck{
		label: "l2_to_cosmos config", rpcURL: "http://127.0.0.1:1", want: 1, tolerateNoAnswer: true,
	})
	if err == nil {
		t.Fatal("a cancelled probe was reported as a tolerated no-answer; startup would continue " +
			"past Ctrl-C and fail later somewhere unrelated")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("err = %v, want it to wrap context.Canceled so the caller can tell a shutdown "+
			"from a real failure", err)
	}
	if verified {
		t.Fatal("verified must stay false when no comparison happened")
	}
}

// The tolerated path must still work when the context is healthy, or this fix
// would turn every unreachable endpoint into a startup failure.
func TestVerifyL2ChainIDStillToleratesAnUnreachableNode(t *testing.T) {
	t.Parallel()

	verified, err := verifyL2ChainID(context.Background(), l2ChainIDCheck{
		label: "l2_to_cosmos config", rpcURL: "http://127.0.0.1:1", want: 1, tolerateNoAnswer: true,
	})
	if err != nil {
		t.Fatalf("an unreachable endpoint must not stop startup: %v", err)
	}
	if verified {
		t.Fatal("nothing was compared, so verified must be false")
	}
}
