package l2rollup

import "testing"

// spans reproduces the loop scanPacketLogs and filterAssertionCreated share, so the
// boundary arithmetic is pinned once. Both call sites must cover [from,to] exactly:
// no gap (an event would be missed silently) and no span wider than the configured
// chunk (the provider rejects it, which is the bug this exists to prevent).
func spans(from, to, chunk uint64) [][2]uint64 {
	if chunk == 0 || to < from || to-from < chunk {
		return [][2]uint64{{from, to}}
	}
	var out [][2]uint64
	for start := from; start <= to; start += chunk {
		end := start + chunk - 1
		if end > to {
			end = to
		}
		out = append(out, [2]uint64{start, end})
	}
	return out
}

func TestLogScanSpansCoverTheRangeExactly(t *testing.T) {
	cases := []struct {
		name            string
		from, to, chunk uint64
		wantSpans       int
	}{
		{"unchunked when chunk is zero", 100, 5000, 0, 1},
		{"unchunked when the range fits", 100, 105, 10, 1},
		{"exact multiple", 0, 29, 10, 3},
		{"ragged tail", 0, 25, 10, 3},
		{"single block", 42, 42, 10, 1},
		{"alchemy free tier over a startup lookback", 0, 255, 10, 26},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := spans(tc.from, tc.to, tc.chunk)
			if len(got) != tc.wantSpans {
				t.Fatalf("span count = %d, want %d: %v", len(got), tc.wantSpans, got)
			}
			if got[0][0] != tc.from {
				t.Fatalf("first span starts at %d, want %d", got[0][0], tc.from)
			}
			if last := got[len(got)-1]; last[1] != tc.to {
				t.Fatalf("last span ends at %d, want %d", last[1], tc.to)
			}
			for i, s := range got {
				if s[1] < s[0] {
					t.Fatalf("span %d is inverted: %v", i, s)
				}
				if tc.chunk > 0 && s[1]-s[0]+1 > tc.chunk {
					t.Fatalf("span %d spans %d blocks, over the %d cap: %v",
						i, s[1]-s[0]+1, tc.chunk, s)
				}
				if i > 0 && s[0] != got[i-1][1]+1 {
					t.Fatalf("gap or overlap between %v and %v", got[i-1], s)
				}
			}
		})
	}
}

// A chunk of 1 is degenerate but must still be correct rather than looping forever
// or skipping the last block.
func TestLogScanSpansWithChunkOfOne(t *testing.T) {
	got := spans(10, 14, 1)
	if len(got) != 5 {
		t.Fatalf("span count = %d, want 5: %v", len(got), got)
	}
	for i, s := range got {
		if s[0] != s[1] || s[0] != uint64(10+i) {
			t.Fatalf("span %d = %v, want [%d,%d]", i, s, 10+i, 10+i)
		}
	}
}

func TestWithLogScanChunkIsRecorded(t *testing.T) {
	s := (&Source{}).WithLogScanChunk(10)
	if s.logScanChunk != 10 {
		t.Fatalf("logScanChunk = %d, want 10", s.logScanChunk)
	}
	if (&Source{}).WithLogScanChunk(0).logScanChunk != 0 {
		t.Fatal("zero must stay zero — that is the unchunked default")
	}
}
