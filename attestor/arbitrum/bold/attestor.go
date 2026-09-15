package bold

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/big"
	"time"

	arbitrum "attestor/arbitrum"
	"github.com/ethereum/go-ethereum/common"
)

const (
	defaultAssertionPollInterval = 12 * time.Second
	defaultMaxL1BlockRange       = uint64(2_000)
)

// CommitmentResolver checks whether an assertion's L2 block hash is canonical
// in the configured Nitro view.
type CommitmentResolver interface {
	ResolveCanonicalBlockHash(
		context.Context,
		common.Hash,
		arbitrum.RunMode,
	) (arbitrum.BlockCommitment, error)
}

// AssertionAttestorConfig controls finalized-L1 assertion ingestion.
type AssertionAttestorConfig struct {
	PollInterval    time.Duration
	MaxL1BlockRange uint64
}

// AssertionAttestor turns finalized-L1 RollupCore events into a Nitro-checked
// commitment feed.
type AssertionAttestor struct {
	source  AssertionSource
	runtime CommitmentResolver
	store   *arbitrum.AttestedRootStore
	config  AssertionAttestorConfig
	now     func() time.Time
}

// NewAssertionAttestor constructs an assertion ingestion loop.
func NewAssertionAttestor(
	source AssertionSource,
	runtime CommitmentResolver,
	store *arbitrum.AttestedRootStore,
	config AssertionAttestorConfig,
) (*AssertionAttestor, error) {
	if source == nil {
		return nil, errors.New("assertion source must not be nil")
	}
	if runtime == nil {
		return nil, errors.New("Nitro commitment resolver must not be nil")
	}
	if store == nil {
		return nil, errors.New("attested-root store must not be nil")
	}
	if config.PollInterval == 0 {
		config.PollInterval = defaultAssertionPollInterval
	}
	if config.PollInterval < 0 {
		return nil, errors.New("assertion poll interval must be greater than zero")
	}
	if config.MaxL1BlockRange == 0 {
		config.MaxL1BlockRange = defaultMaxL1BlockRange
	}
	return &AssertionAttestor{
		source:  source,
		runtime: runtime,
		store:   store,
		config:  config,
		now:     time.Now,
	}, nil
}

// ValidateChainIDs fails closed if either RPC endpoint or RollupCore belongs to
// a different chain than the configured deployment.
func (a *AssertionAttestor) ValidateChainIDs(
	ctx context.Context,
	expectedL1ChainID uint64,
	expectedL2ChainID uint64,
) error {
	finalizedL1Block, err := a.source.FinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}
	l1ChainID, l2ChainID, err := a.source.ChainIDs(ctx, finalizedL1Block)
	if err != nil {
		return err
	}
	if l1ChainID == nil || l2ChainID == nil {
		return errors.New("assertion source returned a nil chain ID")
	}
	expectedL1 := new(big.Int).SetUint64(expectedL1ChainID)
	expectedL2 := new(big.Int).SetUint64(expectedL2ChainID)
	if l1ChainID.Cmp(expectedL1) != 0 {
		return fmt.Errorf(
			"L1 chain ID mismatch: connected=%s configured=%s",
			l1ChainID,
			expectedL1,
		)
	}
	if l2ChainID.Cmp(expectedL2) != 0 {
		return fmt.Errorf(
			"RollupCore L2 chain ID mismatch: contract=%s configured=%s",
			l2ChainID,
			expectedL2,
		)
	}
	return nil
}

// Run continuously advances the verified assertion feed until cancellation.
// Transient L1 or Nitro failures leave the last persisted frontier intact and
// are retried on the next tick.
func (a *AssertionAttestor) Run(ctx context.Context) {
	sync := func() {
		if err := a.SyncOnce(ctx); err != nil && ctx.Err() == nil {
			log.Printf("Arbitrum assertion attestation failed: %v", err)
		}
	}
	sync()

	ticker := time.NewTicker(a.config.PollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			sync()
		}
	}
}

// SyncOnce ingests every newly finalized L1 assertion event and reconciles
// pending/provisional entries against both RollupCore status and Nitro.
func (a *AssertionAttestor) SyncOnce(ctx context.Context) error {
	finalizedL1Block, err := a.source.FinalizedBlockNumber(ctx)
	if err != nil {
		return err
	}

	for from := a.store.NextL1Block(); from <= finalizedL1Block; {
		to := finalizedL1Block
		if remaining := finalizedL1Block - from; remaining >= a.config.MaxL1BlockRange {
			to = from + a.config.MaxL1BlockRange - 1
		}
		proposals, confirmations, err := a.source.Assertions(ctx, from, to)
		if err != nil {
			return err
		}
		if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
			for _, proposal := range proposals {
				if err := store.RecordProposal(proposal); err != nil {
					return false, err
				}
			}
			for _, confirmation := range confirmations {
				if err := store.MarkAssertionConfirmed(
					confirmation.AssertionHash,
					confirmation.L2BlockHash,
				); err != nil {
					return false, err
				}
			}
			store.SetNextL1Block(to + 1)
			return true, nil
		}); err != nil {
			return err
		}
		if to == finalizedL1Block {
			break
		}
		from = to + 1
	}

	if err := a.reconcileProposals(ctx, finalizedL1Block); err != nil {
		return err
	}
	if err := a.reconcileProvisional(ctx, finalizedL1Block); err != nil {
		return err
	}
	return nil
}

func (a *AssertionAttestor) reconcileProposals(
	ctx context.Context,
	finalizedL1Block uint64,
) error {
	for _, proposal := range a.store.Proposals() {
		status, err := a.source.AssertionStatus(
			ctx,
			proposal,
			finalizedL1Block,
		)
		if err != nil {
			return err
		}
		switch status {
		case assertionStatusNone:
			if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
				store.RemoveAssertion(proposal.AssertionHash)
				return true, nil
			}); err != nil {
				return err
			}
			continue
		case assertionStatusConfirmed:
			proposal.Confirmed = true
		case assertionStatusPending:
			if proposal.Confirmed {
				return fmt.Errorf(
					"assertion %s has a finalized confirmation event but pending status",
					proposal.AssertionHash,
				)
			}
		default:
			return fmt.Errorf(
				"assertion %s has unsupported status %d",
				proposal.AssertionHash,
				status,
			)
		}

		commitment, provisional, err := a.resolveProposal(ctx, proposal)
		switch {
		case errors.Is(err, arbitrum.ErrCommitmentNotReady):
			if proposal.Confirmed {
				if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
					return true, store.MarkAssertionConfirmed(proposal.AssertionHash, proposal.L2BlockHash)
				}); err != nil {
					return err
				}
			}
			continue
		case errors.Is(err, arbitrum.ErrCommitmentMismatch):
			if err := a.rejectAssertion(proposal.AssertionHash, proposal.L2BlockHash, err); err != nil {
				return err
			}
			continue
		case err != nil:
			return err
		}
		if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
			return true, store.RecordAssertionAttestation(proposal, commitment, a.now(), provisional)
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *AssertionAttestor) resolveProposal(
	ctx context.Context,
	proposal arbitrum.ProposedAssertion,
) (arbitrum.BlockCommitment, bool, error) {
	if proposal.Confirmed {
		commitment, err := a.runtime.ResolveCanonicalBlockHash(
			ctx,
			proposal.L2BlockHash,
			arbitrum.RunModeFinalized,
		)
		if err == nil {
			return commitment, false, nil
		}
		if !errors.Is(err, arbitrum.ErrCommitmentNotReady) {
			return arbitrum.BlockCommitment{}, false, err
		}
	}
	commitment, err := a.runtime.ResolveCanonicalBlockHash(
		ctx,
		proposal.L2BlockHash,
		arbitrum.RunModeSafe,
	)
	if err != nil {
		return arbitrum.BlockCommitment{}, false, err
	}
	return commitment, true, nil
}

func (a *AssertionAttestor) reconcileProvisional(
	ctx context.Context,
	finalizedL1Block uint64,
) error {
	for _, entry := range a.store.ProvisionalAssertions() {
		status, err := a.source.AssertionStatus(
			ctx,
			arbitrum.ProposedAssertion{
				AssertionHash: entry.AssertionHash,
				L2BlockHash:   entry.L2BlockHash,
			},
			finalizedL1Block,
		)
		if err != nil {
			return err
		}
		if status == assertionStatusNone {
			if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
				store.RemoveAssertion(entry.AssertionHash)
				return true, nil
			}); err != nil {
				return err
			}
			continue
		}
		if status != assertionStatusPending && status != assertionStatusConfirmed {
			return fmt.Errorf(
				"assertion %s has unsupported status %d",
				entry.AssertionHash,
				status,
			)
		}

		_, err = a.runtime.ResolveCanonicalBlockHash(
			ctx,
			entry.L2BlockHash,
			arbitrum.RunModeSafe,
		)
		switch {
		case errors.Is(err, arbitrum.ErrCommitmentNotReady):
			continue
		case errors.Is(err, arbitrum.ErrCommitmentMismatch):
			if err := a.rejectAssertion(entry.AssertionHash, entry.L2BlockHash, err); err != nil {
				return err
			}
			continue
		case err != nil:
			return err
		}
		if status != assertionStatusConfirmed {
			continue
		}
		_, err = a.runtime.ResolveCanonicalBlockHash(
			ctx,
			entry.L2BlockHash,
			arbitrum.RunModeFinalized,
		)
		switch {
		case errors.Is(err, arbitrum.ErrCommitmentNotReady):
			if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
				return true, store.MarkAssertionConfirmed(entry.AssertionHash, entry.L2BlockHash)
			}); err != nil {
				return err
			}
			continue
		case errors.Is(err, arbitrum.ErrCommitmentMismatch):
			if err := a.rejectAssertion(entry.AssertionHash, entry.L2BlockHash, err); err != nil {
				return err
			}
			continue
		case err != nil:
			return err
		}
		if err := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
			if err := store.MarkAssertionConfirmed(entry.AssertionHash, entry.L2BlockHash); err != nil {
				return false, err
			}
			store.PromoteAssertion(entry.AssertionHash)
			return true, nil
		}); err != nil {
			return err
		}
	}
	return nil
}

func (a *AssertionAttestor) rejectAssertion(
	assertionHash common.Hash,
	l2BlockHash common.Hash,
	err error,
) error {
	if commitErr := a.store.Commit(func(store *arbitrum.AttestedRootStore) (bool, error) {
		store.RecordAssertionMismatch(assertionHash, l2BlockHash, err.Error(), a.now())
		return true, nil
	}); commitErr != nil {
		return commitErr
	}
	log.Printf(
		"Arbitrum assertion rejected: assertion_hash=%s l2_block_hash=%s reason=%v",
		assertionHash,
		l2BlockHash,
		err,
	)
	return nil
}
