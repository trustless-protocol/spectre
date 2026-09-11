package attestorgrpc

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"

	attestorpb "attestor/types/attestor"
	"relayer/chain/l2rollup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// codeServer answers every RPC with one gRPC code, so a test can name a code and
// read back what the relayer decided about it.
type codeServer struct {
	attestorpb.UnimplementedAttestorServiceServer
	code codes.Code
}

func (s *codeServer) VerifyStateRoot(context.Context, *attestorpb.VerifyStateRootRequest) (*attestorpb.VerifyStateRootResponse, error) {
	return nil, status.Error(s.code, "from the server")
}

func (s *codeServer) AttestedUpTo(context.Context, *attestorpb.AttestedUpToRequest) (*attestorpb.AttestedUpToResponse, error) {
	return nil, status.Error(s.code, "from the server")
}

// dialCodeServer runs a real gRPC server and a real Client over it. The transport
// is the point: mapStatus reads a code off the wire, and a unit test that hands
// it a locally built status asserts the assumption instead of the behaviour.
func dialCodeServer(t *testing.T, code codes.Code) *Client {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	server := grpc.NewServer()
	attestorpb.RegisterAttestorServiceServer(server, &codeServer{code: code})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return listener.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	return &Client{conn: conn, rpc: attestorpb.NewAttestorServiceClient(conn)}
}

// The table this pins is a CONTRACT WITH THE ATTESTOR, not a local preference:
// each code is produced by attestor/grpc for one documented reason.
//
// This is HALF of the classification, and deliberately so -- the two halves live
// in different packages because the seam does:
//
//	code -> sentinel   here, over a real connection (mapStatus)
//	sentinel -> outcome  TestClassifyAttestorFailure in the parent package
//
// classifyAttestorFailure is unexported and attestorgrpc imports l2rollup, so one
// test cannot span both without either widening the API or inverting the
// dependency. Naming the pair is better than doing either.
//
// Reported by @neitdung on #473: every FailedPrecondition was read as "replica
// behind" while the server also returned it for a missing block verifier, a
// missing commitment feed, an unconfigured signer, an absent chain table and an
// unsigned positive verdict. Those are configuration and invariant failures that
// no amount of waiting fixes, so the permanent-failure hold this PR added could
// never fire for any of them. The attestor now gives them disjoint codes, and
// this table is the relayer half of that agreement.
func TestServerCodesClassifyIntoTheRightOutcome(t *testing.T) {
	cases := []struct {
		name     string
		code     codes.Code
		sentinel error
		// why names the server condition that produces this code, so a reader can
		// find the other side of the contract.
		why string
	}{
		{
			name: "replica has not reached the block", code: codes.FailedPrecondition,
			sentinel: l2rollup.ErrAttestorReplicaBehind,
			why:      "adapter: block N is above the replica head -- the ONLY meaning left for this code",
		},
		{
			name: "route serves no verifier, feed or signer", code: codes.Unimplemented,
			sentinel: l2rollup.ErrAttestorUnimplemented,
			why:      "server: route does not expose a block verifier / commitment feed; adapter: signer is not configured",
		},
		{
			name: "src_chain is not served", code: codes.NotFound,
			sentinel: l2rollup.ErrAttestorUnknownRoute,
			why:      "server: unknown src_chain, and no chains are configured at all",
		},
		{
			name: "the request itself is wrong", code: codes.InvalidArgument,
			sentinel: l2rollup.ErrAttestorBadRequest,
			why:      "server: malformed input or bad run_mode",
		},
		{
			name: "upstream said nothing", code: codes.Unavailable,
			sentinel: l2rollup.ErrAttestorUnavailable,
			why:      "server: L1/L2 RPC did not answer",
		},
		{
			name: "deadline is silence, not a no", code: codes.DeadlineExceeded,
			sentinel: l2rollup.ErrAttestorUnavailable,
			why:      "transport: the attestor said nothing about the block",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			client := dialCodeServer(t, tc.code)

			// Both RPCs, because they do not share a path: VerifyStateRoot
			// intercepts Unimplemented ahead of mapStatus to report the
			// VerifyStateRoot-specific sentinel (which wraps the general one), so a
			// test that drove only VerifyStateRoot would never execute mapStatus's
			// Unimplemented branch at all. It was mutated to ErrAttestorReplicaBehind
			// and this test stayed green until AttestedUpTo was added here.
			calls := map[string]func() error{
				"VerifyStateRoot": func() error {
					_, err := client.VerifyStateRoot(context.Background(), "op-sepolia", 100,
						make([]byte, 32), make([]byte, 32), attestorpb.RunMode_RUN_MODE_FINALIZED)
					return err
				},
				"AttestedUpTo": func() error {
					_, _, err := client.AttestedUpTo(context.Background(), "op-sepolia", false)
					return err
				},
			}
			for name, call := range calls {
				err := call()
				if err == nil {
					t.Fatalf("%s: no error for %v (%s)", name, tc.code, tc.why)
				}
				if !errors.Is(err, tc.sentinel) {
					t.Fatalf("%s: %v did not map to %v: %v\n  server side: %s", name, tc.code, tc.sentinel, err, tc.why)
				}
				// The operator has to be able to read the server's own words.
				if !strings.Contains(err.Error(), "from the server") {
					t.Fatalf("%s: the status message was dropped: %v", name, err)
				}
			}
		})
	}
}

// An unmapped code must stay transient: keeping packets and retrying is the safe
// default for a failure nobody has classified, and Internal in particular is the
// server's catch-all for anything unclassified.
func TestAnUnmappedCodeStaysTransient(t *testing.T) {
	client := dialCodeServer(t, codes.Internal)
	_, err := client.VerifyStateRoot(context.Background(), "op-sepolia", 100, make([]byte, 32), make([]byte, 32), attestorpb.RunMode_RUN_MODE_FINALIZED)
	if err == nil {
		t.Fatal("Internal produced no error")
	}
	for _, permanent := range []error{
		l2rollup.ErrAttestorUnimplemented, l2rollup.ErrAttestorUnknownRoute, l2rollup.ErrAttestorBadRequest,
	} {
		if errors.Is(err, permanent) {
			t.Fatalf("Internal was mapped to %v; a catch-all code must not hold a source off", permanent)
		}
	}
}
