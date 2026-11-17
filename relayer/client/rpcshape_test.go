package client

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"relayer/rpcmock"

	abci "github.com/cometbft/cometbft/abci/types"
	cmted25519 "github.com/cometbft/cometbft/crypto/ed25519"
	"github.com/cometbft/cometbft/p2p"
	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/gogoproto/proto"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	ics23 "github.com/cosmos/ics23/go"
)

// These tests are the #575 panic class: an RPC answer that is well-formed JSON
// but not the shape the relayer assumed. Each one used to take the process down.

func validatorWithPower(power int64) *commettypes.Validator {
	key := cmted25519.GenPrivKey().PubKey()
	return &commettypes.Validator{Address: key.Address(), PubKey: key, VotingPower: power}
}

// L9: CometBFT's NewValidatorSet panics on a list it cannot accept, and it runs
// before the hash check that would refuse a wrong list.
func TestNewValidatorSet_MalformedListIsAnErrorNotAPanic(t *testing.T) {
	good := validatorWithPower(10)
	duplicate := *good

	cases := []struct {
		name string
		list []*commettypes.Validator
	}{
		{"duplicate address", []*commettypes.Validator{good, &duplicate}},
		{"zero voting power", []*commettypes.Validator{good, validatorWithPower(0)}},
		{"negative voting power", []*commettypes.Validator{good, validatorWithPower(-5)}},
		{"nil entry", []*commettypes.Validator{good, nil}},
		{"total power over the limit", []*commettypes.Validator{
			validatorWithPower(commettypes.MaxTotalVotingPower), validatorWithPower(commettypes.MaxTotalVotingPower),
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			set, err := newValidatorSet(tc.list)
			if err == nil || set != nil {
				t.Fatalf("set=%v err=%v, want the list refused with an error", set, err)
			}
		})
	}

	t.Run("a valid list is the set CometBFT builds", func(t *testing.T) {
		list := []*commettypes.Validator{validatorWithPower(30), validatorWithPower(20), validatorWithPower(10)}
		set, err := newValidatorSet(list)
		if err != nil {
			t.Fatalf("valid list refused: %v", err)
		}
		want := commettypes.NewValidatorSet(list)
		if string(set.Hash()) != string(want.Hash()) || set.TotalVotingPower() != 60 {
			t.Fatalf("set differs from NewValidatorSet: hash %X vs %X", set.Hash(), want.Hash())
		}
	})
}

// A commit answer with a null header or a null commit decodes without error, and
// GetLightBlock read header.ProposerAddress on the next line.
func TestGetLightBlock_CommitWithoutHeaderIsAnError(t *testing.T) {
	for _, tc := range []struct {
		name   string
		header *commettypes.Header
		commit *commettypes.Commit
	}{
		{"no header", nil, &commettypes.Commit{Height: 5}},
		{"no commit", &commettypes.Header{Height: 5}, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			node := rpcmock.NewCosmosNode(t)
			node.Handle("status", func(json.RawMessage) (any, error) {
				return &coretypes.ResultStatus{NodeInfo: p2p.DefaultNodeInfo{}}, nil
			})
			node.Handle("commit", func(json.RawMessage) (any, error) {
				return &coretypes.ResultCommit{SignedHeader: commettypes.SignedHeader{Header: tc.header, Commit: tc.commit}}, nil
			})

			_, err := GetLightBlock(context.Background(), node.Client(t), 5)
			if err == nil || !strings.Contains(err.Error(), "no header or no commit") {
				t.Fatalf("err = %v, want the malformed commit refused", err)
			}
			if n := node.Calls("validators"); n != 0 {
				t.Fatalf("asked for validators %d time(s) after an unusable commit", n)
			}
		})
	}
}

// L11: a client-state query can succeed and carry no client state.
func TestClientStateQueryWithoutAClientState(t *testing.T) {
	node := rpcmock.NewCosmosNode(t)
	node.Handle("abci_query", func(json.RawMessage) (any, error) {
		value, err := proto.Marshal(&clienttypes.QueryClientStateResponse{})
		if err != nil {
			return nil, err
		}
		return &coretypes.ResultABCIQuery{Response: abci.ResponseQuery{Value: value}}, nil
	})
	client := node.Client(t)

	if _, err := GetEthereumClientState(context.Background(), client, "08-wasm-0"); err == nil ||
		!strings.Contains(err.Error(), "no client state") {
		t.Fatalf("GetEthereumClientState err = %v, want the empty answer refused", err)
	}
	if _, err := GetWasmClientLatestHeight(context.Background(), client, "08-wasm-0"); err == nil ||
		!strings.Contains(err.Error(), "no client state") {
		t.Fatalf("GetWasmClientLatestHeight err = %v, want the empty answer refused", err)
	}
}

// The period and slot arithmetic divides by three client-state fields. The state
// is read back from an RPC, so a zero there is an integer divide by zero.
func TestEthereumClientState_Validate(t *testing.T) {
	valid := EthereumClientState{SlotsPerEpoch: 32, EpochsPerSyncCommitteePeriod: 256, SecondsPerSlot: 12}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid state refused: %v", err)
	}
	for name, mutate := range map[string]func(*EthereumClientState){
		"slots_per_epoch":                  func(cs *EthereumClientState) { cs.SlotsPerEpoch = 0 },
		"epochs_per_sync_committee_period": func(cs *EthereumClientState) { cs.EpochsPerSyncCommitteePeriod = 0 },
		"seconds_per_slot":                 func(cs *EthereumClientState) { cs.SecondsPerSlot = 0 },
	} {
		state := valid
		mutate(&state)
		if err := state.Validate(); err == nil || !strings.Contains(err.Error(), name) {
			t.Errorf("%s = 0: err = %v, want it refused and named", name, err)
		}
	}
	if err := (*EthereumClientState)(nil).Validate(); err == nil {
		t.Error("a nil state must be refused")
	}
}

// L10: three shapes of ICS23 proof that parsed to a nil pointer or panicked.
func TestParseCommitmentProof_MalformedShapes(t *testing.T) {
	exist := &ics23.ExistenceProof{Key: []byte("k"), Value: []byte("v"), Leaf: &ics23.LeafOp{Hash: ics23.HashOp_SHA256}}
	inner := &ics23.InnerOp{Hash: ics23.HashOp_SHA256, Prefix: []byte{1}}

	compressed := func(entries []*ics23.CompressedBatchEntry, lookup []*ics23.InnerOp) *ics23.CommitmentProof {
		return &ics23.CommitmentProof{Proof: &ics23.CommitmentProof_Compressed{
			Compressed: &ics23.CompressedBatchProof{Entries: entries, LookupInners: lookup},
		}}
	}
	compressedExist := func(path ...int32) *ics23.CompressedBatchEntry {
		return &ics23.CompressedBatchEntry{Proof: &ics23.CompressedBatchEntry_Exist{
			Exist: &ics23.CompressedExistenceProof{Key: []byte("k"), Value: []byte("v"), Leaf: exist.Leaf, Path: path},
		}}
	}

	refused := map[string]*ics23.CommitmentProof{
		// Returned (nil, nil): a nil proof and no error.
		"batch entry of neither kind": {Proof: &ics23.CommitmentProof_Batch{
			Batch: &ics23.BatchProof{Entries: []*ics23.BatchEntry{{}}},
		}},
		"batch with no entries":   {Proof: &ics23.CommitmentProof_Batch{Batch: &ics23.BatchProof{}}},
		"existence proof missing": {Proof: &ics23.CommitmentProof_Exist{}},
		"non-existence missing":   {Proof: &ics23.CommitmentProof_Nonexist{}},
		// Panicked inside ics23.Decompress: index out of range.
		"compressed step past the lookup table": compressed([]*ics23.CompressedBatchEntry{compressedExist(0, 7)}, []*ics23.InnerOp{inner}),
		"compressed negative step":              compressed([]*ics23.CompressedBatchEntry{compressedExist(-1)}, []*ics23.InnerOp{inner}),
		"compressed step with no lookup table":  compressed([]*ics23.CompressedBatchEntry{compressedExist(0)}, nil),
		// Panicked inside ics23.Decompress: nil dereference.
		"compressed entry of neither kind": compressed([]*ics23.CompressedBatchEntry{{}}, []*ics23.InnerOp{inner}),
		"compressed nil entry":             compressed([]*ics23.CompressedBatchEntry{nil}, []*ics23.InnerOp{inner}),
		"compressed proof missing":         {Proof: &ics23.CommitmentProof_Compressed{}},
		"compressed non-existence step out of range": compressed([]*ics23.CompressedBatchEntry{{
			Proof: &ics23.CompressedBatchEntry_Nonexist{Nonexist: &ics23.CompressedNonExistenceProof{
				Key: []byte("k"), Right: &ics23.CompressedExistenceProof{Key: []byte("r"), Path: []int32{3}},
			}},
		}}, []*ics23.InnerOp{inner}),
	}
	for name, proof := range refused {
		t.Run(name, func(t *testing.T) {
			parsed, err := ParseCommitmentProof(proof)
			if err == nil {
				t.Fatalf("parsed = %+v with no error, want the proof refused", parsed)
			}
			if parsed != nil {
				t.Fatalf("parsed = %+v alongside err = %v", parsed, err)
			}
		})
	}

	// A well-formed compressed proof still parses, to what its first entry says.
	t.Run("a well-formed compressed proof still parses", func(t *testing.T) {
		parsed, err := ParseCommitmentProof(compressed(
			[]*ics23.CompressedBatchEntry{compressedExist(0, 1, 0)},
			[]*ics23.InnerOp{inner, {Hash: ics23.HashOp_SHA256, Suffix: []byte{2}}},
		))
		if err != nil {
			t.Fatalf("well-formed compressed proof refused: %v", err)
		}
		if parsed.ProofType != ProofType_EXIST || string(parsed.ExistenceProof.Key) != "k" || len(parsed.ExistenceProof.Path) != 3 {
			t.Fatalf("parsed = %+v, want the existence proof with its 3 path steps", parsed)
		}
		if len(parsed.ExistenceProof.Path[1].Suffix) != 1 || parsed.ExistenceProof.Path[1].Suffix[0] != 2 {
			t.Fatalf("path step 1 = %+v, want the second lookup entry", parsed.ExistenceProof.Path[1])
		}
	})
}
