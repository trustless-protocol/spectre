// Package server exposes the attestor as a gRPC sidecar: read-only queries
// over the attested-root feed plus a frontier-watch stream. The relayer
// connects to this service and relays only against roots the attestor has
// affirmed — the same consume-a-verifier relationship it has with the prover.
package server

import (
	"context"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
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

// runModeHead maps the caller's requested finality to the replica head the
// answer must be judged against, so a verification is gated at the same level
// the relayer's height gate used.
func runModeHead(mode attestorpb.RunMode) (opstack.Head, error) {
	switch mode {
	case attestorpb.RunMode_RUN_MODE_UNSAFE:
		return opstack.HeadUnsafe, nil
	case attestorpb.RunMode_RUN_MODE_SAFE:
		return opstack.HeadSafe, nil
	case attestorpb.RunMode_RUN_MODE_FINALIZED:
		return opstack.HeadFinalized, nil
	default:
		return "", status.Error(codes.InvalidArgument, "run_mode must be unsafe, safe or finalized")
	}
}

// chainForVerify resolves which attestor a VerifyStateRoot request means.
//
// src_chain is optional so an older caller keeps working: when this daemon serves
// exactly one chain there is nothing to disambiguate. With several configured it
// is required, because guessing would verify against the wrong replica and answer
// valid=false for a perfectly good header — which the relayer reads as a
// divergence and stops on, rather than as the misroute it is.
func (s *Server) chainForVerify(srcChain string) (*opstack.OpStackAttestor, error) {
	if srcChain != "" {
		return s.chain(srcChain)
	}
	switch len(s.chains) {
	case 0:
		return nil, status.Error(codes.FailedPrecondition, "no chains are configured")
	case 1:
		for _, a := range s.chains {
			return a, nil
		}
	}
	names := make([]string, 0, len(s.chains))
	for name := range s.chains {
		names = append(names, name)
	}
	sort.Strings(names)
	return nil, status.Errorf(codes.InvalidArgument,
		"src_chain is required: this attestor serves %d chains (%s)",
		len(s.chains), strings.Join(names, ", "))
}

// VerifyStateRoot answers whether the caller's block identity is what this
// attestor's replica holds at that height.
//
// The relayer calls it once per header build: the height gate (AttestedUpTo)
// bounds how far it may relay, and this bounds what it relays at that height.
// Without it a relayer whose L2 RPC reorged past the frontier — or points at a
// different endpoint than the replica — would package a block the attestor
// never verified, and the wasm client would accept it, since the client checks
// only the header against itself and the router proof against the header.
func (s *Server) VerifyStateRoot(ctx context.Context, req *attestorpb.VerifyStateRootRequest) (*attestorpb.VerifyStateRootResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	if len(req.GetExpectedStateRoot()) != common.HashLength {
		return nil, status.Errorf(codes.InvalidArgument,
			"expected_state_root must contain exactly %d bytes", common.HashLength)
	}
	if n := len(req.GetExpectedBlockHash()); n != 0 && n != common.HashLength {
		return nil, status.Errorf(codes.InvalidArgument,
			"expected_block_hash must be empty or contain exactly %d bytes", common.HashLength)
	}
	head, err := runModeHead(req.GetRunMode())
	if err != nil {
		return nil, err
	}
	attestorForChain, err := s.chainForVerify(req.GetSrcChain())
	if err != nil {
		return nil, err
	}

	headBlock, err := attestorForChain.HeadAt(ctx, head)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "read replica %s head: %v", head, err)
	}
	// "Not derived yet" is not "wrong block": the caller retries the first and
	// must treat the second as a divergence, so they cannot share an answer.
	if req.GetBlockNumber() > headBlock {
		return nil, status.Errorf(codes.FailedPrecondition,
			"block %d is above the replica %s head %d", req.GetBlockNumber(), head, headBlock)
	}

	commitment, err := attestorForChain.CommitmentAt(ctx, req.GetBlockNumber())
	if err != nil {
		return nil, status.Errorf(codes.Unavailable,
			"query replica block %d: %v", req.GetBlockNumber(), err)
	}

	valid := commitment.StateRoot == common.BytesToHash(req.GetExpectedStateRoot())
	if expected := req.GetExpectedBlockHash(); len(expected) != 0 {
		valid = valid && commitment.BlockHash == common.BytesToHash(expected)
	}

	return &attestorpb.VerifyStateRootResponse{
		Valid:       valid,
		BlockNumber: commitment.BlockNumber,
		BlockHash:   append([]byte(nil), commitment.BlockHash[:]...),
		StateRoot:   append([]byte(nil), commitment.StateRoot[:]...),
	}, nil
}
