package services

import (
	"crypto/ed25519"
	"testing"
	"time"

	relayerclient "relayer/client"
	"relayer/prover/witness"

	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/types"
)

func TestExtractValidatorSignatures_NilLightBlock(t *testing.T) {
	_, err := ExtractValidatorSignatures(nil, "test-chain", 4)
	if err == nil {
		t.Fatal("expected error for nil light block")
	}
}

func TestExtractValidatorSignatures_NilCommit(t *testing.T) {
	lb := &relayerclient.LightBlock{
		SignedHeader: types.SignedHeader{Commit: nil},
	}
	_, err := ExtractValidatorSignatures(lb, "test-chain", 4)
	if err == nil {
		t.Fatal("expected error for nil commit")
	}
}

func TestExtractValidatorSignatures_EmptyValidators(t *testing.T) {
	lb := &relayerclient.LightBlock{
		SignedHeader: types.SignedHeader{
			Commit: &types.Commit{Signatures: []types.CommitSig{}},
		},
		ValSet: types.ValidatorSet{Validators: []*types.Validator{}},
	}
	_, err := ExtractValidatorSignatures(lb, "test-chain", 4)
	if err == nil {
		t.Fatal("expected error for empty validators")
	}
}

// makeSignedCommit builds a LightBlock with n validators whose voting powers are
// given by `powers`. A subset signs per `signs[i]`; others are marked absent.
func makeSignedCommit(t *testing.T, chainID string, powers []int64, signs []bool) *relayerclient.LightBlock {
	t.Helper()
	if len(powers) != len(signs) {
		t.Fatal("powers/signs length mismatch")
	}
	flags := make([]types.BlockIDFlag, len(signs))
	for i, signed := range signs {
		flags[i] = types.BlockIDFlagAbsent
		if signed {
			flags[i] = types.BlockIDFlagCommit
		}
	}
	return makeCommitWithFlags(t, chainID, powers, flags)
}

// makeCommitWithFlags builds a LightBlock whose validator i (by `powers` order)
// casts the precommit named by flags[i]. COMMIT and NIL votes carry a signature
// CometBFT accepts for that vote; any other flag carries 64 filler bytes, since
// CometBFT cannot even compute sign bytes for it.
func makeCommitWithFlags(t *testing.T, chainID string, powers []int64, flags []types.BlockIDFlag) *relayerclient.LightBlock {
	t.Helper()
	if len(powers) != len(flags) {
		t.Fatal("powers/flags length mismatch")
	}
	vals := make([]*types.Validator, len(powers))
	privByAddr := make(map[string]cmted25519.PrivKey, len(powers))
	flagByAddr := make(map[string]types.BlockIDFlag, len(powers))
	for i := range powers {
		priv := cmted25519.GenPrivKey()
		addr := priv.PubKey().Address()
		vals[i] = &types.Validator{
			Address:     addr,
			PubKey:      priv.PubKey(),
			VotingPower: powers[i],
		}
		privByAddr[addr.String()] = priv
		flagByAddr[addr.String()] = flags[i]
	}
	// NewValidatorSet sorts validators (power desc, then address). CometBFT
	// requires commit.Signatures[i] to be index-aligned with
	// valSet.Validators[i], so we must build the signatures in the *sorted*
	// order — not the original `vals` order — otherwise verification pairs each
	// signature with the wrong validator (issue: flaky quorum test).
	valSet := types.NewValidatorSet(vals)

	blockID := types.BlockID{
		Hash:          make([]byte, 32),
		PartSetHeader: types.PartSetHeader{Total: 1, Hash: make([]byte, 32)},
	}
	header := &types.Header{ChainID: chainID, Height: 100, Time: time.Now().UTC()}
	commit := &types.Commit{
		Height:     100,
		Round:      0,
		BlockID:    blockID,
		Signatures: make([]types.CommitSig, len(valSet.Validators)),
	}
	for i, v := range valSet.Validators {
		commit.Signatures[i] = types.CommitSig{
			BlockIDFlag:      types.BlockIDFlagAbsent,
			ValidatorAddress: v.Address,
			Timestamp:        header.Time,
		}
	}
	for i, v := range valSet.Validators {
		flag := flagByAddr[v.Address.String()]
		if flag == types.BlockIDFlagAbsent {
			continue
		}
		cs := commit.Signatures[i]
		cs.BlockIDFlag = flag
		commit.Signatures[i] = cs
		if flag != types.BlockIDFlagCommit && flag != types.BlockIDFlagNil {
			cs.Signature = make([]byte, ed25519.SignatureSize)
			commit.Signatures[i] = cs
			continue
		}
		sig, err := privByAddr[v.Address.String()].Sign(commit.VoteSignBytes(chainID, int32(i)))
		if err != nil {
			t.Fatalf("sign[%d]: %v", i, err)
		}
		cs.Signature = sig
		commit.Signatures[i] = cs
	}

	return &relayerclient.LightBlock{
		SignedHeader: types.SignedHeader{Header: header, Commit: commit},
		ValSet:       *valSet,
		BlockHeight:  100,
	}
}

func TestExtractValidatorSignatures_ReachesQuorum(t *testing.T) {
	chainID := "test-chain"
	// 4 validators with equal power; 3 sign = 75% ≥ 2/3.
	lb := makeSignedCommit(t, chainID, []int64{10, 10, 10, 10}, []bool{true, true, true, false})
	got, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Signatures) < 3 {
		t.Fatalf("expected at least 3 signers, got %d", len(got.Signatures))
	}
	for i, s := range got.Signatures {
		if len(s.Signature) != ed25519.SignatureSize {
			t.Fatalf("sig[%d] wrong length %d", i, len(s.Signature))
		}
		// Reconstruct signed bytes from commit to verify locally. The extractor
		// already did this check, but we want the test to fail loudly if the
		// reconstruction logic ever drifts.
		voteBytes := lb.SignedHeader.Commit.VoteSignBytes(chainID, int32(s.Index))
		if !ed25519.Verify(s.PublicKey, voteBytes, s.Signature) {
			t.Fatalf("sig[%d] failed local verify", i)
		}
	}
}

func TestExtractValidatorSignatures_InsufficientPower(t *testing.T) {
	chainID := "test-chain"
	// 4 validators; only 2 of 4 sign = 50% < 2/3.
	lb := makeSignedCommit(t, chainID, []int64{10, 10, 10, 10}, []bool{true, true, false, false})
	_, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err == nil {
		t.Fatal("expected insufficient voting power error")
	}
}

func TestExtractValidatorSignatures_AllAbsent(t *testing.T) {
	chainID := "test-chain"
	lb := makeSignedCommit(t, chainID, []int64{10}, []bool{false})
	_, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err == nil {
		t.Fatal("expected error when no signatures present")
	}
}

func TestExtractValidatorSignatures_GreedyPicksSmallestPrefix(t *testing.T) {
	chainID := "test-chain"
	// One dominant validator (70%) covers quorum alone; extractor should stop at 1.
	lb := makeSignedCommit(t, chainID, []int64{70, 10, 10, 10}, []bool{true, true, true, true})
	got, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Signatures) != 1 {
		t.Fatalf("expected 1 signer for dominant validator, got %d", len(got.Signatures))
	}
	if got.Signatures[0].Power != 70 {
		t.Fatalf("expected dominant power 70, got %d", got.Signatures[0].Power)
	}
}

func TestExtractValidatorSignatures_AllowedIndices(t *testing.T) {
	chainID := "test-chain"
	lb := makeSignedCommit(t, chainID, []int64{30, 25, 25, 20}, []bool{true, true, true, true})

	allowed := map[uint32]bool{1: true, 2: true, 3: true}
	got, err := ExtractValidatorSignatures(lb, chainID, 4, allowed)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var accumulated int64
	for _, sig := range got.Signatures {
		if !allowed[uint32(sig.Index)] {
			t.Fatalf("selected disallowed validator index %d", sig.Index)
		}
		accumulated += sig.Power
	}
	quorum := lb.ValSet.TotalVotingPower()*2/3 + 1
	if accumulated < quorum {
		t.Fatalf("selected power %d below quorum %d", accumulated, quorum)
	}

	_, err = ExtractValidatorSignatures(lb, chainID, 4, map[uint32]bool{3: true})
	if err == nil {
		t.Fatal("expected insufficient voting power with too-small allowed subset")
	}
}

func TestExtractValidatorSignatures_SortsSelectedSignersByIndex(t *testing.T) {
	sigs := []witness.ValidatorSignature{
		{Index: 2, Power: 70},
		{Index: 0, Power: 20},
		{Index: 1, Power: 10},
	}
	sortSelectedSignaturesByIndex(sigs)
	if sigs[0].Index != 0 || sigs[1].Index != 1 || sigs[2].Index != 2 {
		t.Fatalf("expected sorted indices [0 1 2], got [%d %d %d]", sigs[0].Index, sigs[1].Index, sigs[2].Index)
	}
}

// A nil precommit is protocol-legal and its signature verifies, but it commits to
// no block: the circuit and SpectreClient both reject it. The top-power validator
// is always first in the greedy selection, so selecting its nil vote would fail
// every proof for as long as it keeps voting nil (#575 M1).
func TestExtractValidatorSignatures_SkipsNilPrecommit(t *testing.T) {
	chainID := "test-chain"
	commit, nilVote := types.BlockIDFlagCommit, types.BlockIDFlagNil
	// Distinct powers, so validator set index 0 is the 30-power nil voter.
	lb := makeCommitWithFlags(t, chainID, []int64{30, 25, 25, 20},
		[]types.BlockIDFlag{nilVote, commit, commit, commit})

	got, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err != nil {
		t.Fatalf("quorum is reachable without the nil voter, got: %v", err)
	}
	for _, list := range [][]witness.ValidatorSignature{got.Signatures, got.Candidates} {
		for _, sig := range list {
			if flag := lb.SignedHeader.Commit.Signatures[sig.Index].BlockIDFlag; flag != commit {
				t.Fatalf("validator index %d with BlockIDFlag %d was offered as a proof signer", sig.Index, flag)
			}
		}
	}
	if len(got.Signatures) != 3 {
		t.Fatalf("expected the 3 COMMIT voters, got %d signers", len(got.Signatures))
	}
}

// Without the nil voter the COMMIT power is 60 of 100: no quorum, and the
// extractor must say so instead of counting the nil vote toward it.
func TestExtractValidatorSignatures_NilPrecommitDoesNotCountTowardQuorum(t *testing.T) {
	chainID := "test-chain"
	commit, nilVote := types.BlockIDFlagCommit, types.BlockIDFlagNil
	lb := makeCommitWithFlags(t, chainID, []int64{40, 20, 20, 20},
		[]types.BlockIDFlag{nilVote, commit, commit, commit})

	if _, err := ExtractValidatorSignatures(lb, chainID, 4); err == nil {
		t.Fatal("expected insufficient voting power: a nil precommit is not a proof signer")
	}
}

// A flag CometBFT does not define makes Commit.VoteSignBytes panic. An RPC can
// send one, so the entry must be skipped before sign bytes are computed.
func TestExtractValidatorSignatures_SkipsUnknownBlockIDFlag(t *testing.T) {
	chainID := "test-chain"
	commit := types.BlockIDFlagCommit
	lb := makeCommitWithFlags(t, chainID, []int64{30, 25, 25, 20},
		[]types.BlockIDFlag{types.BlockIDFlag(0), commit, commit, commit})

	got, err := ExtractValidatorSignatures(lb, chainID, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Signatures) != 3 {
		t.Fatalf("expected the 3 COMMIT voters, got %d signers", len(got.Signatures))
	}
}
