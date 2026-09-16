package attestorgrpc

import (
	"context"
	"crypto/ed25519"
	"net"
	"testing"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"
	"relayer/chain/l2rollup"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

type contractAttestorServer struct {
	attestorpb.UnimplementedAttestorServiceServer
	frontier   *attestorpb.AttestedUpToResponse
	verify     *attestorpb.VerifyStateRootResponse
	verifyErr  error
	verifySeen *attestorpb.VerifyStateRootRequest
}

func (s *contractAttestorServer) AttestedUpTo(_ context.Context, _ *attestorpb.AttestedUpToRequest) (*attestorpb.AttestedUpToResponse, error) {
	return s.frontier, nil
}

func (s *contractAttestorServer) VerifyStateRoot(_ context.Context, request *attestorpb.VerifyStateRootRequest) (*attestorpb.VerifyStateRootResponse, error) {
	s.verifySeen = request
	if s.verifyErr != nil {
		return nil, s.verifyErr
	}
	return s.verify, nil
}

func newGeneratedContractClient(t *testing.T, server *contractAttestorServer) *Client {
	t.Helper()
	listener := bufconn.Listen(1 << 20)
	grpcServer := googlegrpc.NewServer()
	attestorpb.RegisterAttestorServiceServer(grpcServer, server)
	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(grpcServer.Stop)
	connection, err := googlegrpc.NewClient(
		"passthrough:///attestor-contract",
		googlegrpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return listener.DialContext(ctx)
		}),
		googlegrpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { _ = connection.Close() })
	return &Client{conn: connection, rpc: attestorpb.NewAttestorServiceClient(connection)}
}

// TestClientVerifyStateRootTransportContract pins Client as a pure transport:
// request fields reach the wire unmodified, a valid signature round-trips,
// gRPC status codes propagate, and a bad signature from the attestor passes
// through verbatim rather than being swallowed or normalised — this client
// does not verify signatures itself. The fail-closed boundary that rejects a
// bad signature lives in chain/l2rollup (header_attested.go), not here.
func TestClientVerifyStateRootTransportContract(t *testing.T) {
	const (
		chainID = uint64(8453)
		height  = uint64(42)
	)
	stateRoot := make([]byte, 32)
	stateRoot[0] = 0xaa
	blockHash := make([]byte, 32)
	blockHash[0] = 0xbb
	var router [attestation.RouterLength]byte
	router[0] = 0xcc
	signer, err := attestation.NewSigner(chainID, "0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f20")
	if err != nil {
		t.Fatalf("new signer: %v", err)
	}
	setHash, err := (attestation.AttestorConfig{PublicKeys: [][]byte{signer.PublicKey()}, Threshold: 1}).SetHash()
	if err != nil {
		t.Fatalf("attestor set hash: %v", err)
	}
	sign := func(l2Router [attestation.RouterLength]byte, attestorSetHash [attestation.HashLength]byte) []byte {
		t.Helper()
		signature, signErr := signer.Sign(l2Router[:], attestorSetHash[:], height, blockHash, stateRoot)
		if signErr != nil {
			t.Fatalf("sign identity: %v", signErr)
		}
		return signature
	}
	statement, err := attestation.SigningBytes(chainID, router[:], setHash[:], height, blockHash, stateRoot)
	if err != nil {
		t.Fatalf("build attestation statement: %v", err)
	}
	response := func(valid bool, signature []byte) *attestorpb.VerifyStateRootResponse {
		return &attestorpb.VerifyStateRootResponse{
			Valid: valid, BlockNumber: height, BlockHash: blockHash, StateRoot: stateRoot, AttestationSignature: signature,
		}
	}
	request := l2rollup.VerificationRequest{
		SrcChain: "base-mainnet", BlockNumber: height, StateRoot: stateRoot, BlockHash: blockHash,
		RunMode: attestorpb.RunMode_RUN_MODE_FINALIZED, L2Router: router, AttestorSetHash: setHash,
	}

	server := &contractAttestorServer{
		frontier: &attestorpb.AttestedUpToResponse{Found: false},
		verify:   response(true, sign(router, setHash)),
	}
	client := newGeneratedContractClient(t, server)
	if root, found, err := client.AttestedUpTo(context.Background(), "base-mainnet", false); err != nil || found || root != nil {
		t.Fatalf("found=false frontier = (%+v, %t, %v)", root, found, err)
	}

	verdict, err := client.VerifyStateRoot(context.Background(), request)
	if err != nil || !verdict.Valid {
		t.Fatalf("valid signed verdict = (%+v, %v)", verdict, err)
	}
	if verdict.BlockNumber != height || string(verdict.StateRoot) != string(stateRoot) || string(verdict.BlockHash) != string(blockHash) {
		t.Fatalf("valid signed verdict identity = %+v", verdict)
	}
	if !ed25519.Verify(signer.PublicKey(), statement, verdict.Signature) {
		t.Fatal("relayer rejected a valid generated signature")
	}
	if got := server.verifySeen; got == nil || got.GetSrcChain() != "base-mainnet" || got.GetBlockNumber() != height ||
		string(got.GetExpectedStateRoot()) != string(stateRoot) || string(got.GetExpectedBlockHash()) != string(blockHash) ||
		got.GetRunMode() != attestorpb.RunMode_RUN_MODE_FINALIZED || string(got.GetL2Router()) != string(router[:]) ||
		string(got.GetAttestorSetHash()) != string(setHash[:]) {
		t.Fatalf("VerifyStateRoot request = %+v", got)
	}

	if t.Run("canonical mismatch", func(t *testing.T) {
		server.verify = response(false, nil)
		server.verifyErr = nil
		result, callErr := client.VerifyStateRoot(context.Background(), request)
		if callErr != nil || result.Valid {
			t.Fatalf("VerifyStateRoot = (%+v, %v), want valid=false", result, callErr)
		}
	}) {
		// Bad-signature responses: the client must pass the signature through
		// exactly as the attestor sent it, neither swallowing nor normalising
		// it, so that the real verifier (outside this package) still sees
		// what was actually on the wire.
		for _, tc := range []struct {
			name      string
			signature []byte
		}{
			{name: "missing signature", signature: nil},
			{name: "wrong signature", signature: make([]byte, attestation.SignatureLength)},
			{name: "wrong attestation context", signature: sign([attestation.RouterLength]byte{0xff}, setHash)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				server.verify = response(true, tc.signature)
				server.verifyErr = nil
				result, callErr := client.VerifyStateRoot(context.Background(), request)
				if callErr != nil || !result.Valid {
					t.Fatalf("VerifyStateRoot = (%+v, %v), want valid=true (transport does not verify)", result, callErr)
				}
				if string(result.Signature) != string(tc.signature) {
					t.Fatalf("signature passthrough = %x, want verbatim %x", result.Signature, tc.signature)
				}
			})
		}
	}

	for _, expectedCode := range []codes.Code{codes.FailedPrecondition, codes.Unavailable} {
		server.verifyErr = status.Error(expectedCode, "test upstream failure")
		_, err := client.VerifyStateRoot(context.Background(), request)
		if status.Code(err) != expectedCode {
			t.Fatalf("transport status = %v, want %v (err=%v)", status.Code(err), expectedCode, err)
		}
	}
}
