package server

import (
	"context"
	"errors"
	"sort"

	"attestor/arbitrum"
	attestorpb "attestor/types/attestor"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AttestorServer serves the verified commitment frontier and lower-level
// state-root checks backed by the attestor-owned Nitro node.
type AttestorServer struct {
	attestorpb.UnimplementedAttestorServiceServer
	runtime *arbitrum.RuntimeState
	feeds   map[string]AttestedRootReader
}

// AttestedRootReader is the chain-specific verified feed exposed through the
// shared oracle-shaped attestor API.
type AttestedRootReader interface {
	HighestAttested(includeProvisional bool) (arbitrum.AttestedRoot, bool)
	HighestAttestedAtOrBelow(
		l2BlockNumber uint64,
		includeProvisional bool,
	) (arbitrum.AttestedRoot, bool)
}

// NewAttestorServer constructs a server backed by a Nitro block reader.
func NewAttestorServer(reader arbitrum.NitroHeaderReader) (*AttestorServer, error) {
	runtime, err := arbitrum.NewRuntimeState(reader)
	if err != nil {
		return nil, err
	}
	return NewAttestorServerWithRuntime(runtime)
}

// NewAttestorServerWithRuntime constructs a server backed by the shared
// unsafe/safe/finalized Nitro runtime tracker.
func NewAttestorServerWithRuntime(runtime *arbitrum.RuntimeState) (*AttestorServer, error) {
	return NewAttestorServerWithRuntimeAndFeeds(runtime, nil)
}

// NewAttestorServerWithRuntimeAndFeeds constructs a server and registers
// its independently verified per-chain commitment feeds.
func NewAttestorServerWithRuntimeAndFeeds(
	runtime *arbitrum.RuntimeState,
	feeds map[string]AttestedRootReader,
) (*AttestorServer, error) {
	if runtime == nil {
		return nil, errors.New("attestor runtime state must not be nil")
	}
	copiedFeeds := make(map[string]AttestedRootReader, len(feeds))
	for srcChain, feed := range feeds {
		if srcChain == "" {
			return nil, errors.New("attested-root feed src_chain must not be empty")
		}
		if feed == nil {
			return nil, errors.New("attested-root feed must not be nil")
		}
		copiedFeeds[srcChain] = feed
	}
	return &AttestorServer{runtime: runtime, feeds: copiedFeeds}, nil
}

// Info reports every configured source chain and the current replica frontier.
func (s *AttestorServer) Info(
	_ context.Context,
	request *attestorpb.InfoRequest,
) (*attestorpb.InfoResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	if s == nil || s.runtime == nil {
		return nil, status.Error(codes.FailedPrecondition, "attestor server is not initialized")
	}
	snapshot := s.runtime.Snapshot()
	srcChains := make([]string, 0, len(s.feeds))
	for srcChain := range s.feeds {
		srcChains = append(srcChains, srcChain)
	}
	sort.Strings(srcChains)

	response := &attestorpb.InfoResponse{
		Chains: make([]*attestorpb.ChainInfo, 0, len(srcChains)),
	}
	for _, srcChain := range srcChains {
		feed := s.feeds[srcChain]
		chain := &attestorpb.ChainInfo{
			SrcChain: srcChain,
			// Arbitrum does not have a daemon-wide gating head. State-root
			// verification selects unsafe/safe/finalized per request, so leave
			// attestation_head unset instead of reporting a fabricated value.
			ReplicaSeen: snapshot.UnsafeObserved || snapshot.SafeObserved || snapshot.FinalizedSeen,
		}
		if snapshot.UnsafeObserved {
			chain.ReplicaUnsafeL2 = snapshot.Unsafe.BlockNumber
		}
		if snapshot.SafeObserved {
			chain.ReplicaSafeL2 = snapshot.Safe.BlockNumber
		}
		if snapshot.FinalizedSeen {
			chain.ReplicaFinalizedL2 = snapshot.Finalized.BlockNumber
		}
		if root, found := feed.HighestAttested(false); found {
			chain.AttestedUpTo = attestedRootToProto(root)
		}
		if root, found := feed.HighestAttested(true); found {
			chain.AttestedUpToProvisional = attestedRootToProto(root)
		}
		response.Chains = append(response.Chains, chain)
	}
	return response, nil
}

// AttestedUpTo returns the highest independently verified commitment for the
// requested source chain.
func (s *AttestorServer) AttestedUpTo(
	_ context.Context,
	request *attestorpb.AttestedUpToRequest,
) (*attestorpb.AttestedUpToResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	feed, err := s.feed(request.GetSrcChain())
	if err != nil {
		return nil, err
	}
	root, found := feed.HighestAttested(request.GetIncludeProvisional())
	response := &attestorpb.AttestedUpToResponse{Found: found}
	if found {
		response.Root = attestedRootToProto(root)
	}
	return response, nil
}

// AttestedRootAtOrBelow returns the best independently verified commitment no
// greater than the caller's inclusive L2 height bound.
func (s *AttestorServer) AttestedRootAtOrBelow(
	_ context.Context,
	request *attestorpb.AttestedRootAtOrBelowRequest,
) (*attestorpb.AttestedRootAtOrBelowResponse, error) {
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	feed, err := s.feed(request.GetSrcChain())
	if err != nil {
		return nil, err
	}
	root, found := feed.HighestAttestedAtOrBelow(
		request.GetL2BlockNumber(),
		request.GetIncludeProvisional(),
	)
	response := &attestorpb.AttestedRootAtOrBelowResponse{Found: found}
	if found {
		response.Root = attestedRootToProto(root)
	}
	return response, nil
}

func (s *AttestorServer) feed(srcChain string) (AttestedRootReader, error) {
	if s == nil {
		return nil, status.Error(codes.FailedPrecondition, "attestor server is not initialized")
	}
	if srcChain == "" {
		return nil, status.Error(codes.InvalidArgument, "src_chain must not be empty")
	}
	feed, ok := s.feeds[srcChain]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "source chain %q is not configured", srcChain)
	}
	return feed, nil
}

func attestedRootToProto(root arbitrum.AttestedRoot) *attestorpb.AttestedRoot {
	response := &attestorpb.AttestedRoot{
		L2BlockNumber: root.L2BlockNumber,
		Root:          append([]byte(nil), root.Root[:]...),
		Source:        root.Source,
		Provisional:   root.Provisional,
		AttestedAt:    root.AttestedAt.Unix(),
	}
	switch {
	case root.AssertionHash != (common.Hash{}):
		response.Provenance = &attestorpb.AttestedRoot_AssertionHash{
			AssertionHash: append([]byte(nil), root.AssertionHash[:]...),
		}
	case root.Source == arbitrum.SourceGame:
		// A game index of zero is valid. The oneof records its presence even
		// though zero is the scalar default.
		response.Provenance = &attestorpb.AttestedRoot_GameIndex{
			GameIndex: root.GameIndex,
		}
	}
	return response
}

// VerifyStateRoot compares the caller's expected root and optional block hash
// with Nitro's canonical block header at the requested L2 height.
func (s *AttestorServer) VerifyStateRoot(
	ctx context.Context,
	request *attestorpb.VerifyStateRootRequest,
) (*attestorpb.VerifyStateRootResponse, error) {
	if s == nil || s.runtime == nil {
		return nil, status.Error(codes.FailedPrecondition, "state-root verifier is not initialized")
	}
	if request == nil {
		return nil, status.Error(codes.InvalidArgument, "request must not be nil")
	}
	if len(request.GetExpectedStateRoot()) != common.HashLength {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"expected_state_root must contain exactly %d bytes",
			common.HashLength,
		)
	}
	if length := len(request.GetExpectedBlockHash()); length != 0 && length != common.HashLength {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"expected_block_hash must be empty or contain exactly %d bytes",
			common.HashLength,
		)
	}
	runMode, err := runModeFromProto(request.GetRunMode())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	head, err := s.runtime.Head(ctx, runMode)
	if errors.Is(err, ethereum.NotFound) {
		return nil, status.Errorf(codes.Unavailable, "Nitro %s head is not available", runMode)
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "query Nitro %s head: %v", runMode, err)
	}
	if request.GetBlockNumber() > head.BlockNumber {
		return nil, status.Errorf(
			codes.FailedPrecondition,
			"block %d is above Nitro %s head %d",
			request.GetBlockNumber(),
			runMode,
			head.BlockNumber,
		)
	}

	commitment, err := s.runtime.CommitmentAt(ctx, request.GetBlockNumber())
	if errors.Is(err, ethereum.NotFound) {
		return nil, status.Errorf(codes.NotFound, "Nitro block %d was not found", request.GetBlockNumber())
	}
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "query Nitro block %d: %v", request.GetBlockNumber(), err)
	}

	blockHash := commitment.BlockHash
	stateRoot := commitment.StateRoot
	valid := stateRoot == common.BytesToHash(request.GetExpectedStateRoot())
	if expectedBlockHash := request.GetExpectedBlockHash(); len(expectedBlockHash) != 0 {
		valid = valid && blockHash == common.BytesToHash(expectedBlockHash)
	}

	return &attestorpb.VerifyStateRootResponse{
		Valid:       valid,
		BlockNumber: commitment.BlockNumber,
		BlockHash:   append([]byte(nil), blockHash[:]...),
		StateRoot:   append([]byte(nil), stateRoot[:]...),
	}, nil
}
