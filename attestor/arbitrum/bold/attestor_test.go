package bold

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	arbitrum "attestor/arbitrum"
	"github.com/ethereum/go-ethereum/common"
)

func TestAssertionAttestorPromotesOnlyAfterRollupAndNitroFinalization(t *testing.T) {
	store, err := arbitrum.NewAttestedRootStore("arbitrum-one", 100)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	proposal := arbitrum.ProposedAssertion{
		AssertionHash:    common.HexToHash("0x01"),
		ParentHash:       common.HexToHash("0x02"),
		L2BlockHash:      common.HexToHash("0x03"),
		InboxAccumulator: common.HexToHash("0x04"),
		L1BlockNumber:    100,
	}
	source := &assertionTestSource{
		finalizedBlock: 100,
		proposals:      []arbitrum.ProposedAssertion{proposal},
		status:         assertionStatusPending,
	}
	commitment := arbitrum.BlockCommitment{
		BlockNumber: 500,
		BlockHash:   proposal.L2BlockHash,
		StateRoot:   common.HexToHash("0x05"),
	}
	resolver := &assertionTestResolver{
		results: map[arbitrum.RunMode]arbitrum.BlockCommitment{
			arbitrum.RunModeSafe: commitment,
		},
		errs: map[arbitrum.RunMode]error{
			arbitrum.RunModeFinalized: fmt.Errorf(
				"%w: finalized head has not reached block",
				arbitrum.ErrCommitmentNotReady,
			),
		},
	}
	loop, err := NewAssertionAttestor(
		source,
		resolver,
		store,
		AssertionAttestorConfig{PollInterval: time.Second, MaxL1BlockRange: 10},
	)
	if err != nil {
		t.Fatalf("create assertion attestor: %v", err)
	}
	loop.now = func() time.Time { return time.Unix(1_700_000_000, 0) }

	if err := loop.SyncOnce(context.Background()); err != nil {
		t.Fatalf("sync pending assertion: %v", err)
	}
	if _, found := store.HighestAttested(false); found {
		t.Fatal("pending assertion entered finalized frontier")
	}
	root, found := store.HighestAttested(true)
	if !found || !root.Provisional || root.L2BlockNumber != 500 {
		t.Fatalf("provisional frontier: found=%t root=%+v", found, root)
	}

	source.status = assertionStatusConfirmed
	resolver.errs[arbitrum.RunModeFinalized] = nil
	resolver.results[arbitrum.RunModeFinalized] = commitment
	if err := loop.SyncOnce(context.Background()); err != nil {
		t.Fatalf("sync confirmed assertion: %v", err)
	}
	root, found = store.HighestAttested(false)
	if !found || root.Provisional || root.AssertionHash != proposal.AssertionHash {
		t.Fatalf("confirmed frontier: found=%t root=%+v", found, root)
	}
}

func TestAssertionAttestorExcludesNitroMismatch(t *testing.T) {
	store, err := arbitrum.NewAttestedRootStore("arbitrum-one", 10)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	proposal := arbitrum.ProposedAssertion{
		AssertionHash: common.HexToHash("0x11"),
		L2BlockHash:   common.HexToHash("0x12"),
	}
	source := &assertionTestSource{
		finalizedBlock: 10,
		proposals:      []arbitrum.ProposedAssertion{proposal},
		status:         assertionStatusPending,
	}
	resolver := &assertionTestResolver{
		errs: map[arbitrum.RunMode]error{
			arbitrum.RunModeSafe: fmt.Errorf(
				"%w: assertion block is not canonical",
				arbitrum.ErrCommitmentMismatch,
			),
		},
	}
	loop, err := NewAssertionAttestor(
		source,
		resolver,
		store,
		AssertionAttestorConfig{PollInterval: time.Second},
	)
	if err != nil {
		t.Fatalf("create assertion attestor: %v", err)
	}
	if err := loop.SyncOnce(context.Background()); err != nil {
		t.Fatalf("sync mismatched assertion: %v", err)
	}
	if _, found := store.HighestAttested(true); found {
		t.Fatal("mismatched assertion entered feed")
	}
	mismatches := store.Mismatches()
	if len(mismatches) != 1 || mismatches[0].AssertionHash != proposal.AssertionHash {
		t.Fatalf("mismatch records: %+v", mismatches)
	}
}

type assertionTestSource struct {
	finalizedBlock uint64
	proposals      []arbitrum.ProposedAssertion
	confirmations  []ConfirmedAssertion
	status         uint8
}

func (s *assertionTestSource) FinalizedBlockNumber(context.Context) (uint64, error) {
	return s.finalizedBlock, nil
}

func (s *assertionTestSource) ChainIDs(
	context.Context,
	uint64,
) (*big.Int, *big.Int, error) {
	return big.NewInt(1), big.NewInt(42161), nil
}

func (s *assertionTestSource) Assertions(
	context.Context,
	uint64,
	uint64,
) ([]arbitrum.ProposedAssertion, []ConfirmedAssertion, error) {
	proposals := append([]arbitrum.ProposedAssertion(nil), s.proposals...)
	confirmations := append([]ConfirmedAssertion(nil), s.confirmations...)
	s.proposals = nil
	s.confirmations = nil
	return proposals, confirmations, nil
}

func (s *assertionTestSource) AssertionStatus(
	context.Context,
	common.Hash,
	uint64,
) (uint8, error) {
	return s.status, nil
}

type assertionTestResolver struct {
	results map[arbitrum.RunMode]arbitrum.BlockCommitment
	errs    map[arbitrum.RunMode]error
}

func (r *assertionTestResolver) ResolveCanonicalBlockHash(
	_ context.Context,
	_ common.Hash,
	mode arbitrum.RunMode,
) (arbitrum.BlockCommitment, error) {
	return r.results[mode], r.errs[mode]
}
