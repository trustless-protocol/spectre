package opstack

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

func testGame(idx uint64, l2Block uint64, claim [32]byte) ProposedRoot {
	return ProposedRoot{
		GameIndex:     idx,
		GameAddress:   common.BytesToAddress([]byte{byte(idx + 1)}),
		RootClaim:     claim,
		L2BlockNumber: l2Block,
		L1Timestamp:   1000 + idx,
	}
}

func TestStoreRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := LoadStore(path)
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	if !s.Fresh() {
		t.Fatal("missing state file must yield a fresh store")
	}
	s.Bootstrap(7)

	now := time.Unix(4000, 0).UTC()
	g0, g1, g2 := testGame(7, 100, root(0x11)), testGame(8, 200, root(0x12)), testGame(9, 300, root(0x13))
	s.AdvanceIngest(8, &g0)
	s.AdvanceIngest(9, &g1)
	s.AdvanceIngest(10, &g2)
	s.RecordMatch(g0, now, false)
	s.RecordMatch(g1, now, true) // provisional, enters recheck
	if err := s.Save(); err != nil {
		t.Fatalf("Save: %v", err)
	}

	re, err := LoadStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if re.Fresh() {
		t.Fatal("reloaded store must not be fresh")
	}
	if got := re.NextGameIndex(); got != 10 {
		t.Fatalf("reloaded cursor = %d, want 10", got)
	}
	if got := re.Pending(); len(got) != 1 || got[0].GameIndex != 9 {
		t.Fatalf("reloaded pending = %+v, want only game 9", got)
	}
	if got := re.RecheckList(); len(got) != 1 || got[0].Game.GameIndex != 8 {
		t.Fatalf("reloaded recheck = %+v, want only game 8", got)
	}
	if cur, ok := re.AttestedUpTo(false); !ok || cur.Height != 100 {
		t.Fatalf("reloaded confirmed AttestedUpTo = (%+v, %t), want height 100", cur, ok)
	}
	if cur, ok := re.AttestedUpTo(true); !ok || cur.Height != 200 {
		t.Fatalf("reloaded provisional AttestedUpTo = (%+v, %t), want height 200", cur, ok)
	}
}

func TestStoreCorruptFileFailsLoud(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadStore(path); err == nil {
		t.Fatal("corrupt state file must be a hard error, not a silent reset")
	}
}

func TestHighestAttestedAtOrBelow(t *testing.T) {
	s, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatal(err)
	}
	now := time.Unix(4000, 0)
	s.RecordMatch(testGame(0, 100, root(0x21)), now, false)
	s.RecordMatch(testGame(1, 200, root(0x22)), now, true) // provisional
	s.RecordMatch(testGame(2, 300, root(0x23)), now, false)

	if _, ok := s.HighestAttestedAtOrBelow(99, false); ok {
		t.Fatal("nothing attested at or below 99")
	}
	if got, ok := s.HighestAttestedAtOrBelow(100, false); !ok || got.L2BlockNumber != 100 {
		t.Fatalf("at 100: (%+v, %t), want the boundary entry itself", got, ok)
	}
	// Provisional entry at 200 must be skipped without opt-in.
	if got, ok := s.HighestAttestedAtOrBelow(250, false); !ok || got.L2BlockNumber != 100 {
		t.Fatalf("at 250 confirmed-only: (%+v, %t), want height 100", got, ok)
	}
	if got, ok := s.HighestAttestedAtOrBelow(250, true); !ok || got.L2BlockNumber != 200 {
		t.Fatalf("at 250 with provisional: (%+v, %t), want height 200", got, ok)
	}
	if got, ok := s.HighestAttestedAtOrBelow(1000, false); !ok || got.L2BlockNumber != 300 {
		t.Fatalf("at 1000: (%+v, %t), want height 300", got, ok)
	}
}

func TestSaveIsAtomicReplacement(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	s, err := LoadStore(path)
	if err != nil {
		t.Fatal(err)
	}
	s.Bootstrap(1)
	if err := s.Save(); err != nil {
		t.Fatalf("first save: %v", err)
	}
	s.Bootstrap(2)
	if err := s.Save(); err != nil {
		t.Fatalf("second save: %v", err)
	}
	re, err := LoadStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if got := re.NextGameIndex(); got != 2 {
		t.Fatalf("cursor = %d, want 2", got)
	}
	// No temp files may survive a successful save.
	entries, err := os.ReadDir(filepath.Dir(path))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		names := make([]string, 0, len(entries))
		for _, e := range entries {
			names = append(names, e.Name())
		}
		t.Fatalf("state dir contains %v, want only the state file", names)
	}
}

// Regression: derived feed entries persist GameIndex 0, which collides with
// factory game index 0. Resolving game 0's provisional verdict must never
// touch derived entries — confirming must not clear their provisional flag,
// and correcting must not delete them.
func TestResolveGame0DoesNotTouchDerivedEntries(t *testing.T) {
	now := time.Unix(4000, 0).UTC()
	game0 := testGame(0, 50, root(0xbb))

	t.Run("confirm keeps derived provisional", func(t *testing.T) {
		s, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
		if err != nil {
			t.Fatalf("LoadStore: %v", err)
		}
		s.AppendDerived(100, root(0xaa), now, true) // provisional derived root
		s.RecordMatch(game0, now, true)             // provisional verdict for game 0

		s.ConfirmProvisional(0)

		if got, ok := s.HighestAttestedAtOrBelow(100, false); ok && got.source() == SourceDerived {
			t.Fatalf("derived root at 100 lost its provisional flag via game-0 confirmation: %+v", got)
		}
		if got, ok := s.HighestAttestedAtOrBelow(50, false); !ok || got.GameIndex != 0 || got.source() != SourceGame {
			t.Fatalf("game 0's own entry should be confirmed, got %+v ok=%t", got, ok)
		}
	})

	t.Run("correct keeps derived entries", func(t *testing.T) {
		s, err := LoadStore(filepath.Join(t.TempDir(), "state.json"))
		if err != nil {
			t.Fatalf("LoadStore: %v", err)
		}
		s.AppendDerived(100, root(0xaa), now, false) // confirmed derived root
		s.RecordMatch(game0, now, true)

		// Finalized recheck contradicts game 0's provisional match.
		s.CorrectProvisional(RecheckEntry{Game: game0, ProvisionalMatch: true}, root(0xcc), now)

		got, ok := s.HighestAttestedAtOrBelow(100, false)
		if !ok || got.source() != SourceDerived || got.L2BlockNumber != 100 {
			t.Fatalf("correcting game 0 must not delete the derived entry at 100, got %+v ok=%t", got, ok)
		}
		if entry, ok := s.HighestAttestedAtOrBelow(50, true); ok && entry.source() == SourceGame {
			t.Fatalf("game 0's wrong provisional entry should have been pulled, got %+v", entry)
		}
	})
}
