package arbitrum

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	// SourceGame marks an OP commitment selected from a dispute game.
	SourceGame = "game"

	// SourceAssertion marks an Arbitrum commitment selected from a RollupCore
	// assertion or legacy node and independently checked against the local
	// Nitro replica.
	SourceAssertion = "assertion"

	attestedRootStateVersion = 1
)

// AttestedRoot is one independently verified commitment exposed to relayers.
// Rollup-specific identifiers are kept beside the generic height/root fields so
// the header builder can resolve the L1 object it must prove.
type AttestedRoot struct {
	L2BlockNumber uint64      `json:"l2_block_number"`
	Root          common.Hash `json:"root"`
	Source        string      `json:"source"`
	GameIndex     uint64      `json:"game_index,omitempty"`
	AssertionHash common.Hash `json:"assertion_hash,omitempty"`
	Provisional   bool        `json:"provisional,omitempty"`
	AttestedAt    time.Time   `json:"attested_at"`

	// Arbitrum-only verification metadata. These fields are persisted so a
	// provisional entry can be rechecked and promoted without trusting the
	// relayer to provide its provenance.
	L2BlockHash        common.Hash `json:"l2_block_hash,omitempty"`
	L1BlockNumber      uint64      `json:"l1_block_number,omitempty"`
	AssertionConfirmed bool        `json:"assertion_confirmed,omitempty"`
	LegacyNodeNumber   uint64      `json:"legacy_node_number,omitempty"`
}

// ProposedAssertion is a finalized-L1 BoLD AssertionCreated or legacy
// NodeCreated event waiting for the local Nitro replica to reach and verify its
// committed L2 block.
type ProposedAssertion struct {
	AssertionHash    common.Hash `json:"assertion_hash"`
	ParentHash       common.Hash `json:"parent_hash"`
	L2BlockHash      common.Hash `json:"l2_block_hash"`
	InboxAccumulator common.Hash `json:"inbox_accumulator"`
	L1BlockNumber    uint64      `json:"l1_block_number"`
	Confirmed        bool        `json:"confirmed,omitempty"`
	LegacyNodeNumber uint64      `json:"legacy_node_number,omitempty"`
}

// AssertionMismatch is a terminal local verdict that an L1 assertion did not
// match the canonical block independently derived by Nitro.
type AssertionMismatch struct {
	AssertionHash common.Hash `json:"assertion_hash"`
	L2BlockHash   common.Hash `json:"l2_block_hash"`
	Reason        string      `json:"reason"`
	RecordedAt    time.Time   `json:"recorded_at"`
}

type persistedAttestedRootState struct {
	Version        uint32              `json:"version"`
	SrcChain       string              `json:"src_chain"`
	SourceIdentity string              `json:"source_identity,omitempty"`
	NextL1Block    uint64              `json:"next_l1_block"`
	Proposals      []ProposedAssertion `json:"proposals,omitempty"`
	Attested       []AttestedRoot      `json:"attested,omitempty"`
	Mismatches     []AssertionMismatch `json:"mismatches,omitempty"`
}

// BindSourceIdentity pins the durable feed to one concrete rollup deployment.
// Existing pre-identity state is upgraded in place; a conflicting identity
// fails closed so stale roots cannot be served after a config change.
func (s *AttestedRootStore) BindSourceIdentity(identity string) error {
	if identity == "" {
		return errors.New("attested-root source identity must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.state.SourceIdentity == "" {
		s.state.SourceIdentity = identity
		return nil
	}
	if s.state.SourceIdentity != identity {
		return fmt.Errorf(
			"attested-root source identity %q does not match configured %q",
			s.state.SourceIdentity,
			identity,
		)
	}
	return nil
}

// AttestedRootStore persists the assertion ingest cursor and independently
// verified feed. Writers mutate it from the Arbitrum attestation loop; gRPC
// readers use the RWMutex-protected frontier accessors.
type AttestedRootStore struct {
	mu    sync.RWMutex
	path  string
	state persistedAttestedRootState
}

// NewAttestedRootStore constructs an in-memory feed, primarily for tests.
func NewAttestedRootStore(srcChain string, startL1Block uint64) (*AttestedRootStore, error) {
	if srcChain == "" {
		return nil, errors.New("attested-root store src_chain must not be empty")
	}
	return &AttestedRootStore{
		state: persistedAttestedRootState{
			Version:     attestedRootStateVersion,
			SrcChain:    srcChain,
			NextL1Block: startL1Block,
		},
	}, nil
}

// LoadAttestedRootStore loads a durable feed. A missing file starts at
// startL1Block; corrupt or identity-mismatched files fail closed.
func LoadAttestedRootStore(path, srcChain string, startL1Block uint64) (*AttestedRootStore, error) {
	store, err := NewAttestedRootStore(srcChain, startL1Block)
	if err != nil {
		return nil, err
	}
	if path == "" {
		return nil, errors.New("attested-root state path must not be empty")
	}
	store.path = path

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			return nil, fmt.Errorf("create attested-root state directory: %w", err)
		}
		return store, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read attested-root state %q: %w", path, err)
	}

	var state persistedAttestedRootState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("decode attested-root state %q: %w", path, err)
	}
	if state.Version != attestedRootStateVersion {
		return nil, fmt.Errorf(
			"unsupported attested-root state version %d, want %d",
			state.Version,
			attestedRootStateVersion,
		)
	}
	if state.SrcChain != srcChain {
		return nil, fmt.Errorf(
			"attested-root state src_chain %q does not match configured %q",
			state.SrcChain,
			srcChain,
		)
	}
	store.state = state
	return store, nil
}

// Save atomically replaces the persistent state file. In-memory stores have no
// path and therefore require no persistence.
func (s *AttestedRootStore) Save() error {
	if s == nil {
		return errors.New("attested-root store is nil")
	}
	s.mu.RLock()
	path := s.path
	data, err := json.MarshalIndent(s.state, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return fmt.Errorf("marshal attested-root state: %w", err)
	}
	if path == "" {
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("create attested-root state directory: %w", err)
	}
	file, err := os.CreateTemp(dir, filepath.Base(path)+".tmp-*")
	if err != nil {
		return fmt.Errorf("create temporary attested-root state: %w", err)
	}
	temporaryPath := file.Name()
	cleanup := func() {
		_ = file.Close()
		_ = os.Remove(temporaryPath)
	}
	if err := file.Chmod(0o600); err != nil {
		cleanup()
		return fmt.Errorf("restrict temporary attested-root state: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		cleanup()
		return fmt.Errorf("write temporary attested-root state: %w", err)
	}
	if err := file.Sync(); err != nil {
		cleanup()
		return fmt.Errorf("sync temporary attested-root state: %w", err)
	}
	if err := file.Close(); err != nil {
		_ = os.Remove(temporaryPath)
		return fmt.Errorf("close temporary attested-root state: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		_ = os.Remove(temporaryPath)
		return fmt.Errorf("replace attested-root state: %w", err)
	}
	return nil
}

// SrcChain returns the immutable source-chain identifier.
func (s *AttestedRootStore) SrcChain() string {
	if s == nil {
		return ""
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.SrcChain
}

// NextL1Block is the next finalized L1 block whose assertion logs must be
// ingested.
func (s *AttestedRootStore) NextL1Block() uint64 {
	if s == nil {
		return 0
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state.NextL1Block
}

// SetNextL1Block advances the finalized-log cursor.
func (s *AttestedRootStore) SetNextL1Block(next uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state.NextL1Block = next
}

// RecordProposal idempotently records a hash-checked AssertionCreated event.
func (s *AttestedRootStore) RecordProposal(proposal ProposedAssertion) error {
	if proposal.AssertionHash == (common.Hash{}) {
		return errors.New("proposal assertion hash must not be zero")
	}
	if proposal.L2BlockHash == (common.Hash{}) {
		return errors.New("proposal L2 block hash must not be zero")
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for index := range s.state.Proposals {
		existing := &s.state.Proposals[index]
		if existing.AssertionHash != proposal.AssertionHash {
			continue
		}
		if existing.ParentHash != proposal.ParentHash ||
			existing.L2BlockHash != proposal.L2BlockHash ||
			existing.InboxAccumulator != proposal.InboxAccumulator ||
			existing.LegacyNodeNumber != proposal.LegacyNodeNumber {
			return fmt.Errorf("assertion %s was observed with conflicting event data", proposal.AssertionHash)
		}
		existing.Confirmed = existing.Confirmed || proposal.Confirmed
		return nil
	}
	for index := range s.state.Attested {
		existing := &s.state.Attested[index]
		if existing.AssertionHash == proposal.AssertionHash {
			if existing.L2BlockHash != proposal.L2BlockHash {
				return fmt.Errorf("attested assertion %s changed L2 block hash", proposal.AssertionHash)
			}
			existing.AssertionConfirmed = existing.AssertionConfirmed || proposal.Confirmed
			return nil
		}
	}
	s.state.Proposals = append(s.state.Proposals, proposal)
	return nil
}

// MarkLegacyNodeConfirmed records a finalized legacy NodeConfirmed event. The
// event identifies its proposal by node number rather than node hash.
func (s *AttestedRootStore) MarkLegacyNodeConfirmed(
	nodeNumber uint64,
	l2BlockHash common.Hash,
) error {
	if nodeNumber == 0 {
		return errors.New("confirmed legacy node number must not be zero")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for index := range s.state.Proposals {
		proposal := &s.state.Proposals[index]
		if proposal.LegacyNodeNumber != nodeNumber {
			continue
		}
		if proposal.L2BlockHash != l2BlockHash {
			return fmt.Errorf(
				"confirmed legacy node %d block hash %s does not match proposal %s",
				nodeNumber,
				l2BlockHash,
				proposal.L2BlockHash,
			)
		}
		proposal.Confirmed = true
		return nil
	}
	for index := range s.state.Attested {
		entry := &s.state.Attested[index]
		if entry.LegacyNodeNumber != nodeNumber {
			continue
		}
		if entry.L2BlockHash != l2BlockHash {
			return fmt.Errorf(
				"confirmed legacy node %d block hash %s does not match attested %s",
				nodeNumber,
				l2BlockHash,
				entry.L2BlockHash,
			)
		}
		entry.AssertionConfirmed = true
		return nil
	}
	// A confirmation can legitimately refer to a node created before the
	// configured lookback. It is irrelevant to this feed.
	return nil
}

// MarkAssertionConfirmed records a finalized AssertionConfirmed event.
func (s *AttestedRootStore) MarkAssertionConfirmed(assertionHash, l2BlockHash common.Hash) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index := range s.state.Proposals {
		proposal := &s.state.Proposals[index]
		if proposal.AssertionHash != assertionHash {
			continue
		}
		if proposal.L2BlockHash != l2BlockHash {
			return fmt.Errorf(
				"confirmed assertion %s block hash %s does not match proposal %s",
				assertionHash,
				l2BlockHash,
				proposal.L2BlockHash,
			)
		}
		proposal.Confirmed = true
		return nil
	}
	for index := range s.state.Attested {
		entry := &s.state.Attested[index]
		if entry.AssertionHash != assertionHash {
			continue
		}
		if entry.L2BlockHash != l2BlockHash {
			return fmt.Errorf(
				"confirmed assertion %s block hash %s does not match attested %s",
				assertionHash,
				l2BlockHash,
				entry.L2BlockHash,
			)
		}
		entry.AssertionConfirmed = true
		return nil
	}
	// A confirmation can legitimately refer to an assertion created before the
	// configured lookback. It is irrelevant to this feed.
	return nil
}

// Proposals returns a stable copy of unresolved assertion proposals.
func (s *AttestedRootStore) Proposals() []ProposedAssertion {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ProposedAssertion, len(s.state.Proposals))
	copy(out, s.state.Proposals)
	return out
}

// ProvisionalAssertions returns assertion feed entries awaiting promotion.
func (s *AttestedRootStore) ProvisionalAssertions() []AttestedRoot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []AttestedRoot
	for _, entry := range s.state.Attested {
		if entry.Source == SourceAssertion && entry.Provisional {
			out = append(out, entry)
		}
	}
	return out
}

// RecordAssertionAttestation moves a proposal into the verified feed.
func (s *AttestedRootStore) RecordAssertionAttestation(
	proposal ProposedAssertion,
	commitment BlockCommitment,
	now time.Time,
	provisional bool,
) error {
	if commitment.BlockHash != proposal.L2BlockHash {
		return fmt.Errorf(
			"proposal %s commits to %s, local commitment is %s",
			proposal.AssertionHash,
			proposal.L2BlockHash,
			commitment.BlockHash,
		)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeProposalLocked(proposal.AssertionHash)
	for index := range s.state.Attested {
		entry := &s.state.Attested[index]
		if entry.AssertionHash != proposal.AssertionHash {
			continue
		}
		if entry.L2BlockHash != commitment.BlockHash || entry.Root != commitment.StateRoot {
			return fmt.Errorf("assertion %s changed its local commitment", proposal.AssertionHash)
		}
		entry.Provisional = provisional
		entry.AssertionConfirmed = entry.AssertionConfirmed || proposal.Confirmed
		return nil
	}
	s.state.Attested = append(s.state.Attested, AttestedRoot{
		L2BlockNumber:      commitment.BlockNumber,
		Root:               commitment.StateRoot,
		Source:             SourceAssertion,
		AssertionHash:      proposal.AssertionHash,
		Provisional:        provisional,
		AttestedAt:         now.UTC(),
		L2BlockHash:        commitment.BlockHash,
		L1BlockNumber:      proposal.L1BlockNumber,
		AssertionConfirmed: proposal.Confirmed,
		LegacyNodeNumber:   proposal.LegacyNodeNumber,
	})
	return nil
}

// PromoteAssertion marks an assertion irreversible after both the L1
// confirmation and Nitro-finalized recheck have succeeded.
func (s *AttestedRootStore) PromoteAssertion(assertionHash common.Hash) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for index := range s.state.Attested {
		entry := &s.state.Attested[index]
		if entry.AssertionHash == assertionHash {
			entry.Provisional = false
			entry.AssertionConfirmed = true
			return true
		}
	}
	return false
}

// RemoveAssertion removes an assertion from both pending and attested state.
func (s *AttestedRootStore) RemoveAssertion(assertionHash common.Hash) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeProposalLocked(assertionHash)
	kept := s.state.Attested[:0]
	for _, entry := range s.state.Attested {
		if entry.AssertionHash != assertionHash {
			kept = append(kept, entry)
		}
	}
	s.state.Attested = kept
}

// RemoveLegacyNode removes a rejected legacy node from pending and attested
// state. NodeRejected identifies the node by number and does not carry its
// node hash.
func (s *AttestedRootStore) RemoveLegacyNode(nodeNumber uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	proposals := s.state.Proposals[:0]
	for _, proposal := range s.state.Proposals {
		if proposal.LegacyNodeNumber != nodeNumber {
			proposals = append(proposals, proposal)
		}
	}
	s.state.Proposals = proposals
	attested := s.state.Attested[:0]
	for _, entry := range s.state.Attested {
		if entry.LegacyNodeNumber != nodeNumber {
			attested = append(attested, entry)
		}
	}
	s.state.Attested = attested
}

// RecordAssertionMismatch removes a rejected assertion and retains an alerting
// record for operators.
func (s *AttestedRootStore) RecordAssertionMismatch(
	assertionHash, l2BlockHash common.Hash,
	reason string,
	now time.Time,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removeProposalLocked(assertionHash)
	kept := s.state.Attested[:0]
	for _, entry := range s.state.Attested {
		if entry.AssertionHash != assertionHash {
			kept = append(kept, entry)
		}
	}
	s.state.Attested = kept
	for _, mismatch := range s.state.Mismatches {
		if mismatch.AssertionHash == assertionHash {
			return
		}
	}
	s.state.Mismatches = append(s.state.Mismatches, AssertionMismatch{
		AssertionHash: assertionHash,
		L2BlockHash:   l2BlockHash,
		Reason:        reason,
		RecordedAt:    now.UTC(),
	})
}

// Mismatches returns a stable copy of terminal mismatch records.
func (s *AttestedRootStore) Mismatches() []AssertionMismatch {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]AssertionMismatch, len(s.state.Mismatches))
	copy(out, s.state.Mismatches)
	return out
}

// HighestAttested returns the highest qualifying feed entry.
func (s *AttestedRootStore) HighestAttested(includeProvisional bool) (AttestedRoot, bool) {
	return s.HighestAttestedAtOrBelow(^uint64(0), includeProvisional)
}

// HighestAttestedAtOrBelow returns the highest qualifying entry no greater
// than the inclusive height bound.
func (s *AttestedRootStore) HighestAttestedAtOrBelow(
	l2BlockNumber uint64,
	includeProvisional bool,
) (AttestedRoot, bool) {
	if s == nil {
		return AttestedRoot{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	var best AttestedRoot
	found := false
	for _, entry := range s.state.Attested {
		if entry.L2BlockNumber > l2BlockNumber {
			continue
		}
		if entry.Provisional && !includeProvisional {
			continue
		}
		if !found || betterAttestedRoot(entry, best) {
			best = entry
			found = true
		}
	}
	return best, found
}

func betterAttestedRoot(candidate, current AttestedRoot) bool {
	if candidate.L2BlockNumber != current.L2BlockNumber {
		return candidate.L2BlockNumber > current.L2BlockNumber
	}
	if candidate.Provisional != current.Provisional {
		return !candidate.Provisional
	}
	return bytes.Compare(candidate.AssertionHash[:], current.AssertionHash[:]) < 0
}

func (s *AttestedRootStore) removeProposalLocked(assertionHash common.Hash) {
	kept := s.state.Proposals[:0]
	for _, proposal := range s.state.Proposals {
		if proposal.AssertionHash != assertionHash {
			kept = append(kept, proposal)
		}
	}
	s.state.Proposals = kept
}
