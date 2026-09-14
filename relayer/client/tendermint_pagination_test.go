package client

import (
	"context"
	"fmt"
	"testing"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	commettypes "github.com/cometbft/cometbft/types"
)

// fakePager is a test double for the CometBFT validators RPC. It paginates a
// backing slice and records the perPage values it was asked for, so tests can
// assert we never fall back to the nil/default-30 behavior (issue #105).
type fakePager struct {
	all         []*commettypes.Validator
	perPageSeen []int
	failOnPage  int // 0 = never fail
}

func (f *fakePager) Validators(_ context.Context, _ *int64, page, perPage *int) (*coretypes.ResultValidators, error) {
	if perPage == nil {
		return nil, fmt.Errorf("perPage must not be nil (would trigger default-30 truncation)")
	}
	f.perPageSeen = append(f.perPageSeen, *perPage)
	if f.failOnPage != 0 && *page == f.failOnPage {
		return nil, fmt.Errorf("simulated RPC failure on page %d", *page)
	}
	pp := *perPage
	start := min((*page-1)*pp, len(f.all))
	end := min(start+pp, len(f.all))
	return &coretypes.ResultValidators{
		Validators: f.all[start:end],
		Count:      end - start,
		Total:      len(f.all),
	}, nil
}

func makeVals(n int) []*commettypes.Validator {
	vals := make([]*commettypes.Validator, n)
	for i := range vals {
		vals[i] = &commettypes.Validator{
			Address:     fmt.Appendf(nil, "addr-%04d", i),
			VotingPower: int64(i + 1),
		}
	}
	return vals
}

func TestFetchAllValidators_PagesFullSet(t *testing.T) {
	cases := []struct {
		name      string
		total     int
		wantPages int
	}{
		{"empty", 0, 1},
		{"under one page", 17, 1},
		{"exactly 30 (old default boundary)", 30, 1},
		{"just over old default", 31, 1},
		{"exactly one max page", 100, 1},
		{"just over one page", 101, 2},
		{"large set (180 validators)", 180, 2},
		{"multi page (250)", 250, 3},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := &fakePager{all: makeVals(tc.total)}
			got, err := fetchAllValidators(context.Background(), f, 42)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(got) != tc.total {
				t.Fatalf("got %d validators, want full set of %d (truncation regression)", len(got), tc.total)
			}
			if len(f.perPageSeen) != tc.wantPages {
				t.Fatalf("made %d RPC pages, want %d", len(f.perPageSeen), tc.wantPages)
			}
			for i, pp := range f.perPageSeen {
				if pp != cometBFTMaxPerPage {
					t.Fatalf("page %d requested perPage=%d, want explicit %d (nil would default to 30)", i+1, pp, cometBFTMaxPerPage)
				}
			}
		})
	}
}

func TestFetchAllValidators_PropagatesError(t *testing.T) {
	f := &fakePager{all: makeVals(150), failOnPage: 2}
	_, err := fetchAllValidators(context.Background(), f, 42)
	if err == nil {
		t.Fatal("expected error when a page fetch fails, got nil")
	}
}
