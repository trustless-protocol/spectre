// Package server exposes the attestor as a gRPC sidecar: read-only queries
// over the attested-root feed plus a frontier-watch stream. The relayer
// connects to this service and relays only against roots the attestor has
// affirmed — the same consume-a-verifier relationship it has with the prover.
package server

import (
	"context"
	"math"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"attestor/opstack"
	attestorpb "attestor/types/attestor"
)

const defaultWatchPoll = 2 * time.Second

// Server implements attestorpb.AttestorServiceServer over the in-process
// attestors of one `attest` run. All methods are read-only: the Run goroutines
// stay the sole writers, and every read goes through the stores' RWMutex
// accessors or atomic snapshots.
type Server struct {
	attestorpb.UnimplementedAttestorServiceServer
	chains map[string]*opstack.OpStackAttestor
	// WatchPoll is the frontier polling cadence of WatchAttested streams
	// (defaults to 2s; overridable for tests).
	WatchPoll time.Duration
}

// New builds a server over the given attestors, keyed by src_chain.
func New(chains map[string]*opstack.OpStackAttestor) *Server {
	return &Server{chains: chains}
}

// Register attaches the service to a gRPC server.
func (s *Server) Register(g *grpc.Server) {
	attestorpb.RegisterAttestorServiceServer(g, s)
}

func (s *Server) chain(srcChain string) (*opstack.OpStackAttestor, error) {
	if a, ok := s.chains[srcChain]; ok {
		return a, nil
	}
	return nil, status.Errorf(codes.NotFound, "unknown src_chain %q", srcChain)
}

func toPB(r opstack.AttestedRoot) *attestorpb.AttestedRoot {
	root := make([]byte, len(r.Root))
	copy(root, r.Root[:])
	return &attestorpb.AttestedRoot{
		L2BlockNumber: r.L2BlockNumber,
		Root:          root,
		Source:        r.Source,
		GameIndex:     r.GameIndex,
		Provisional:   r.Provisional,
		AttestedAt:    r.AttestedAt.Unix(),
	}
}

// frontier returns the highest qualifying feed entry.
func frontier(a *opstack.OpStackAttestor, includeProvisional bool) (opstack.AttestedRoot, bool) {
	return a.Store().HighestAttestedAtOrBelow(math.MaxUint64, includeProvisional)
}

func (s *Server) Info(_ context.Context, _ *attestorpb.InfoRequest) (*attestorpb.InfoResponse, error) {
	resp := &attestorpb.InfoResponse{}
	for _, a := range s.chains {
		info := &attestorpb.ChainInfo{
			SrcChain:        a.SrcChain(),
			AttestationHead: string(a.Head()),
			PendingGames:    uint64(len(a.Store().Pending())),
		}
		if st, ok := a.LastSyncStatus(); ok {
			info.ReplicaSeen = true
			info.ReplicaUnsafeL2 = st.UnsafeL2
			info.ReplicaSafeL2 = st.SafeL2
			info.ReplicaFinalizedL2 = st.FinalizedL2
		}
		if entry, ok := frontier(a, false); ok {
			info.AttestedUpTo = toPB(entry)
		}
		if entry, ok := frontier(a, true); ok {
			info.AttestedUpToProvisional = toPB(entry)
		}
		resp.Chains = append(resp.Chains, info)
	}
	return resp, nil
}

func (s *Server) AttestedUpTo(_ context.Context, req *attestorpb.AttestedUpToRequest) (*attestorpb.AttestedUpToResponse, error) {
	a, err := s.chain(req.SrcChain)
	if err != nil {
		return nil, err
	}
	entry, ok := frontier(a, req.IncludeProvisional)
	if !ok {
		return &attestorpb.AttestedUpToResponse{Found: false}, nil
	}
	return &attestorpb.AttestedUpToResponse{Found: true, Root: toPB(entry)}, nil
}

func (s *Server) AttestedRootAtOrBelow(_ context.Context, req *attestorpb.AttestedRootAtOrBelowRequest) (*attestorpb.AttestedRootAtOrBelowResponse, error) {
	a, err := s.chain(req.SrcChain)
	if err != nil {
		return nil, err
	}
	entry, ok := a.Store().HighestAttestedAtOrBelow(req.L2BlockNumber, req.IncludeProvisional)
	if !ok {
		return &attestorpb.AttestedRootAtOrBelowResponse{Found: false}, nil
	}
	return &attestorpb.AttestedRootAtOrBelowResponse{Found: true, Root: toPB(entry)}, nil
}

// WatchAttested streams frontier advances: the current frontier on subscribe,
// then a message whenever the highest qualifying entry grows. It intentionally
// streams the frontier, not every feed entry — a game-verified root landing
// below an already-streamed derived root does not re-emit.
func (s *Server) WatchAttested(req *attestorpb.WatchAttestedRequest, stream grpc.ServerStreamingServer[attestorpb.WatchAttestedResponse]) error {
	a, err := s.chain(req.SrcChain)
	if err != nil {
		return err
	}
	poll := s.WatchPoll
	if poll == 0 {
		poll = defaultWatchPoll
	}
	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	var last uint64
	sent := false
	for {
		if entry, ok := frontier(a, req.IncludeProvisional); ok && (!sent || entry.L2BlockNumber > last) {
			if err := stream.Send(&attestorpb.WatchAttestedResponse{Root: toPB(entry)}); err != nil {
				return err
			}
			last = entry.L2BlockNumber
			sent = true
		}
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
		}
	}
}
