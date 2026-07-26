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

	"attestor/optimism/opstack"
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
	response := &attestorpb.AttestedRoot{
		L2BlockNumber: r.L2BlockNumber,
		Root:          root,
		Source:        r.Source,
		Provisional:   r.Provisional,
		AttestedAt:    r.AttestedAt.Unix(),
	}
	if r.Source == "" || r.Source == opstack.SourceGame {
		// Game index zero is valid, so set the oneof wrapper explicitly.
		response.Provenance = &attestorpb.AttestedRoot_GameIndex{
			GameIndex: r.GameIndex,
		}
	}
	return response
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

// frontierChanged reports whether the frontier's relay-relevant identity —
// height, root, provisional status — differs. Source/GameIndex changes at an
// identical (height, root) are not re-emitted: a game root landing at an
// already-streamed derived height attests the same bytes.
func frontierChanged(last, cur opstack.AttestedRoot) bool {
	return cur.L2BlockNumber != last.L2BlockNumber || cur.Root != last.Root || cur.Provisional != last.Provisional
}

// WatchAttested streams the frontier: the current entry on subscribe, then a
// message whenever its identity (height, root, provisional flag) changes.
// Usually that is growth, but with include_provisional set the finalized
// recheck can also correct the streamed root in place, confirm it, or revoke
// it — revocation retreats the frontier and the new lower entry is re-emitted.
// It intentionally streams the frontier, not every feed entry — a root landing
// below the streamed frontier does not emit.
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

	var last opstack.AttestedRoot
	sent := false
	for {
		entry, ok := frontier(a, req.IncludeProvisional)
		if ok && (!sent || frontierChanged(last, entry)) {
			if err := stream.Send(&attestorpb.WatchAttestedResponse{Root: toPB(entry)}); err != nil {
				return err
			}
			last = entry
			sent = true
		}
		if !ok {
			// The frontier vanished (every qualifying entry was provisional and
			// revoked); re-announce whatever appears next, even if lower.
			sent = false
		}
		select {
		case <-stream.Context().Done():
			return nil
		case <-ticker.C:
		}
	}
}
