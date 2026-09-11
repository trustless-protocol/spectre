// Package grpc is the shared read-only AttestorService adapter. It only knows
// immutable src_chain routes to core ports; it never imports an OP or
// Arbitrum package and never triggers a plugin write path.
package grpc

import (
	"context"
	"fmt"
	"sort"

	"attestor/core"
	attestorpb "attestor/types/attestor"

	googlegrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const signatureLength = 64

// Options controls documented compatibility behavior at the transport
// boundary. New callers should always provide src_chain.
type Options struct {
	AllowLegacyEmptySrcChain bool
}

// Server implements AttestorService using immutable core port routes.
type Server struct {
	attestorpb.UnimplementedAttestorServiceServer
	routes       map[string]core.Ports
	sourceChains []string
	options      Options
}

// New validates and copies the composition root's routing table. Subsequent
// mutation by the caller cannot make a request route to a different plugin.
func New(routes map[string]core.Ports, options Options) (*Server, error) {
	if len(routes) == 0 {
		return nil, fmt.Errorf("attestor routing table must not be empty")
	}
	copied := make(map[string]core.Ports, len(routes))
	names := make([]string, 0, len(routes))
	for srcChain, ports := range routes {
		if srcChain == "" {
			return nil, fmt.Errorf("attestor route src_chain must not be empty")
		}
		if ports.Status == nil {
			return nil, fmt.Errorf("attestor route %q status reader must not be nil", srcChain)
		}
		copied[srcChain] = ports
		names = append(names, srcChain)
	}
	sort.Strings(names)
	return &Server{routes: copied, sourceChains: names, options: options}, nil
}

// Register attaches this adapter to a gRPC server.
func (s *Server) Register(server *googlegrpc.Server) {
	attestorpb.RegisterAttestorServiceServer(server, s)
}

// Info reports each immutable route and only reads plugin-local status/feed
// snapshots. It does not call a runner or refresh an upstream replica.
func (s *Server) Info(ctx context.Context, request *attestorpb.InfoRequest) (*attestorpb.InfoResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	response := &attestorpb.InfoResponse{Chains: make([]*attestorpb.ChainInfo, 0, len(s.sourceChains))}
	for _, srcChain := range s.sourceChains {
		ports := s.routes[srcChain]
		feedStatus, err := ports.Status.Status(ctx)
		if err != nil {
			return nil, mapError(err)
		}
		chain := &attestorpb.ChainInfo{
			SrcChain:           srcChain,
			AttestationHead:    string(feedStatus.AttestationHead),
			ReplicaSeen:        feedStatus.ReplicaSeen,
			ReplicaUnsafeL2:    feedStatus.ReplicaUnsafe,
			ReplicaSafeL2:      feedStatus.ReplicaSafe,
			ReplicaFinalizedL2: feedStatus.ReplicaFinalized,
			PendingGames:       feedStatus.Pending,
		}
		if ports.Feed != nil {
			if root, found, err := ports.Feed.AttestedUpTo(ctx, core.AttestationPolicy{}); err != nil {
				return nil, mapError(err)
			} else if found {
				chain.AttestedUpTo = rootToProto(root)
			}
			if root, found, err := ports.Feed.AttestedUpTo(ctx, core.AttestationPolicy{IncludeProvisional: true}); err != nil {
				return nil, mapError(err)
			} else if found {
				chain.AttestedUpToProvisional = rootToProto(root)
			}
		}
		response.Chains = append(response.Chains, chain)
	}
	return response, nil
}

func (s *Server) AttestedUpTo(ctx context.Context, request *attestorpb.AttestedUpToRequest) (*attestorpb.AttestedUpToResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	_, ports, err := s.route(request.GetSrcChain())
	if err != nil {
		return nil, err
	}
	if ports.Feed == nil {
		return nil, status.Error(codes.Unimplemented, "attestor route does not expose a commitment feed")
	}
	root, found, err := ports.Feed.AttestedUpTo(ctx, core.AttestationPolicy{IncludeProvisional: request.GetIncludeProvisional()})
	if err != nil {
		return nil, mapError(err)
	}
	response := &attestorpb.AttestedUpToResponse{Found: found}
	if found {
		response.Root = rootToProto(root)
	}
	return response, nil
}

func (s *Server) AttestedRootAtOrBelow(ctx context.Context, request *attestorpb.AttestedRootAtOrBelowRequest) (*attestorpb.AttestedRootAtOrBelowResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	_, ports, err := s.route(request.GetSrcChain())
	if err != nil {
		return nil, err
	}
	if ports.Feed == nil {
		return nil, status.Error(codes.Unimplemented, "attestor route does not expose a commitment feed")
	}
	root, found, err := ports.Feed.AttestedRootAtOrBelow(
		ctx,
		request.GetL2BlockNumber(),
		core.AttestationPolicy{IncludeProvisional: request.GetIncludeProvisional()},
	)
	if err != nil {
		return nil, mapError(err)
	}
	response := &attestorpb.AttestedRootAtOrBelowResponse{Found: found}
	if found {
		response.Root = rootToProto(root)
	}
	return response, nil
}

func (s *Server) WatchAttested(request *attestorpb.WatchAttestedRequest, stream googlegrpc.ServerStreamingServer[attestorpb.WatchAttestedResponse]) error {
	if request == nil {
		return status.Error(codes.InvalidArgument, "request must not be nil")
	}
	_, ports, err := s.route(request.GetSrcChain())
	if err != nil {
		return err
	}
	if ports.Watcher == nil {
		return status.Error(codes.Unimplemented, "attestor route does not support frontier watching")
	}
	updates, err := ports.Watcher.WatchFrontier(stream.Context(), core.WatchRequest{
		Policy: core.AttestationPolicy{IncludeProvisional: request.GetIncludeProvisional()},
	})
	if err != nil {
		return mapError(err)
	}
	for {
		select {
		case <-stream.Context().Done():
			return nil
		case update, ok := <-updates:
			if !ok {
				if stream.Context().Err() != nil {
					return nil
				}
				return status.Error(codes.Unavailable, "attestor frontier watch ended")
			}
			if !update.Found {
				continue
			}
			if err := stream.Send(&attestorpb.WatchAttestedResponse{Root: rootToProto(update.Root)}); err != nil {
				return err
			}
		}
	}
}

func (s *Server) VerifyStateRoot(ctx context.Context, request *attestorpb.VerifyStateRootRequest) (*attestorpb.VerifyStateRootResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	_, ports, err := s.route(request.GetSrcChain())
	if err != nil {
		return nil, err
	}
	if ports.Verifier == nil {
		return nil, status.Error(codes.Unimplemented, "attestor route does not expose a block verifier")
	}
	decoded, err := requestFromProto(request)
	if err != nil {
		return nil, err
	}
	verdict, err := ports.Verifier.VerifyStateRoot(ctx, decoded)
	if err != nil {
		return nil, mapError(err)
	}
	if verdict.Valid && len(verdict.Signature) != signatureLength {
		return nil, status.Errorf(codes.Internal, "positive attestation verdict has %d-byte signature, want %d", len(verdict.Signature), signatureLength)
	}
	response := &attestorpb.VerifyStateRootResponse{
		Valid:       verdict.Valid,
		BlockNumber: verdict.BlockNumber,
		BlockHash:   append([]byte(nil), verdict.BlockHash[:]...),
		StateRoot:   append([]byte(nil), verdict.StateRoot[:]...),
	}
	if verdict.Valid {
		response.AttestationSignature = append([]byte(nil), verdict.Signature...)
	}
	return response, nil
}

func (s *Server) route(requested string) (string, core.Ports, error) {
	if requested != "" {
		ports, found := s.routes[requested]
		if !found {
			return "", core.Ports{}, status.Errorf(codes.NotFound, "source chain %q is not configured", requested)
		}
		return requested, ports, nil
	}
	if !s.options.AllowLegacyEmptySrcChain {
		return "", core.Ports{}, status.Error(codes.InvalidArgument, "src_chain must not be empty")
	}
	switch len(s.sourceChains) {
	case 0:
		return "", core.Ports{}, status.Error(codes.NotFound, "no chains are configured")
	case 1:
		name := s.sourceChains[0]
		return name, s.routes[name], nil
	default:
		return "", core.Ports{}, status.Errorf(codes.InvalidArgument, "src_chain is required: this attestor serves %d chains (%s)", len(s.sourceChains), joinSourceChains(s.sourceChains))
	}
}

func requestFromProto(request *attestorpb.VerifyStateRootRequest) (core.BlockIdentityRequest, error) {
	if len(request.GetExpectedStateRoot()) != 32 {
		return core.BlockIdentityRequest{}, status.Error(codes.InvalidArgument, "expected_state_root must contain exactly 32 bytes")
	}
	var decoded core.BlockIdentityRequest
	decoded.BlockNumber = request.GetBlockNumber()
	copy(decoded.ExpectedStateRoot[:], request.GetExpectedStateRoot())
	if blockHash := request.GetExpectedBlockHash(); len(blockHash) != 0 {
		if len(blockHash) != 32 {
			return core.BlockIdentityRequest{}, status.Error(codes.InvalidArgument, "expected_block_hash must be empty or contain exactly 32 bytes")
		}
		var expected [32]byte
		copy(expected[:], blockHash)
		decoded.ExpectedBlockHash = &expected
	}
	switch request.GetRunMode() {
	case attestorpb.RunMode_RUN_MODE_UNSAFE:
		decoded.RunMode = core.RunModeUnsafe
	case attestorpb.RunMode_RUN_MODE_SAFE:
		decoded.RunMode = core.RunModeSafe
	case attestorpb.RunMode_RUN_MODE_FINALIZED:
		decoded.RunMode = core.RunModeFinalized
	default:
		return core.BlockIdentityRequest{}, status.Error(codes.InvalidArgument, "run_mode must be unsafe, safe or finalized")
	}
	return decoded, nil
}

func rootToProto(root core.AttestedRoot) *attestorpb.AttestedRoot {
	response := &attestorpb.AttestedRoot{
		L2BlockNumber: root.L2BlockNumber,
		Root:          append([]byte(nil), root.Root...),
		Source:        root.Source,
		Provisional:   root.Provisional,
		AttestedAt:    root.AttestedAt.Unix(),
	}
	switch {
	case len(root.AssertionHash) != 0:
		response.Provenance = &attestorpb.AttestedRoot_AssertionHash{AssertionHash: append([]byte(nil), root.AssertionHash...)}
	case root.GameIndex != nil:
		response.Provenance = &attestorpb.AttestedRoot_GameIndex{GameIndex: *root.GameIndex}
	}
	return response
}

func mapError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}
	switch core.ErrorKindOf(err) {
	case core.ErrorInvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	case core.ErrorNotFound:
		return status.Error(codes.NotFound, err.Error())
	case core.ErrorFailedPrecondition:
		return status.Error(codes.FailedPrecondition, err.Error())
	case core.ErrorUnavailable:
		return status.Error(codes.Unavailable, err.Error())
	case core.ErrorUnimplemented:
		return status.Error(codes.Unimplemented, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func joinSourceChains(sourceChains []string) string {
	if len(sourceChains) == 0 {
		return ""
	}
	joined := sourceChains[0]
	for _, sourceChain := range sourceChains[1:] {
		joined += ", " + sourceChain
	}
	return joined
}
