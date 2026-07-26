package arbitrum

import (
	"math/big"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func TestAttestedRootStoreFrontiersAndPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "attested-roots.json")
	store, err := LoadAttestedRootStore(path, "arbitrum-one", 1_000)
	if err != nil {
		t.Fatalf("load new store: %v", err)
	}

	finalized := testStoreProposal(1, 90)
	if err := store.RecordProposal(finalized); err != nil {
		t.Fatalf("record finalized proposal: %v", err)
	}
	if err := store.RecordAssertionAttestation(
		finalized,
		testStoreCommitment(90),
		time.Unix(1_700_000_000, 0),
		false,
	); err != nil {
		t.Fatalf("record finalized attestation: %v", err)
	}

	provisional := testStoreProposal(2, 100)
	if err := store.RecordProposal(provisional); err != nil {
		t.Fatalf("record provisional proposal: %v", err)
	}
	if err := store.RecordAssertionAttestation(
		provisional,
		testStoreCommitment(100),
		time.Unix(1_700_000_100, 0),
		true,
	); err != nil {
		t.Fatalf("record provisional attestation: %v", err)
	}

	root, found := store.HighestAttested(false)
	if !found || root.L2BlockNumber != 90 || root.Provisional {
		t.Fatalf("confirmed frontier: found=%t root=%+v", found, root)
	}
	root, found = store.HighestAttested(true)
	if !found || root.L2BlockNumber != 100 || !root.Provisional {
		t.Fatalf("provisional frontier: found=%t root=%+v", found, root)
	}
	root, found = store.HighestAttestedAtOrBelow(95, true)
	if !found || root.L2BlockNumber != 90 {
		t.Fatalf("frontier at or below 95: found=%t root=%+v", found, root)
	}

	store.SetNextL1Block(1_234)
	if err := store.Save(); err != nil {
		t.Fatalf("save store: %v", err)
	}
	reloaded, err := LoadAttestedRootStore(path, "arbitrum-one", 0)
	if err != nil {
		t.Fatalf("reload store: %v", err)
	}
	if reloaded.NextL1Block() != 1_234 {
		t.Fatalf("reloaded cursor: got %d want 1234", reloaded.NextL1Block())
	}
	root, found = reloaded.HighestAttested(true)
	if !found || root.AssertionHash != provisional.AssertionHash {
		t.Fatalf("reloaded frontier: found=%t root=%+v", found, root)
	}
}

func TestAttestedRootStoreRejectsIdentityMismatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "attested-roots.json")
	store, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("load store: %v", err)
	}
	if err := store.BindSourceIdentity("l1:1/l2:42161/rollup:0x1"); err != nil {
		t.Fatalf("bind source identity: %v", err)
	}
	if err := store.Save(); err != nil {
		t.Fatalf("save store: %v", err)
	}
	if _, err := LoadAttestedRootStore(path, "arbitrum-sepolia", 1); err == nil {
		t.Fatal("expected persisted src_chain mismatch")
	}
	reloaded, err := LoadAttestedRootStore(path, "arbitrum-one", 1)
	if err != nil {
		t.Fatalf("reload identity-bound store: %v", err)
	}
	if err := reloaded.BindSourceIdentity("l1:1/l2:42161/rollup:0x2"); err == nil {
		t.Fatal("expected persisted source identity mismatch")
	}
}

func testStoreProposal(id byte, height uint64) ProposedAssertion {
	return ProposedAssertion{
		AssertionHash:    common.Hash{31: id},
		ParentHash:       common.Hash{30: id},
		L2BlockHash:      testStoreCommitment(height).BlockHash,
		InboxAccumulator: common.Hash{29: id},
		L1BlockNumber:    height + 1_000,
		Confirmed:        height < 100,
	}
}

func testStoreCommitment(height uint64) BlockCommitment {
	return BlockCommitment{
		BlockNumber: height,
		BlockHash:   common.BigToHash(newBigInt(height)),
		StateRoot:   common.BigToHash(newBigInt(height + 10_000)),
	}
}

func newBigInt(value uint64) *big.Int {
	return new(big.Int).SetUint64(value)
}
