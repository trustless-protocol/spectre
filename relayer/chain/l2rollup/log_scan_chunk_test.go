package l2rollup

import (
	"errors"
	"fmt"
	"testing"

	"relayer/chain"
	relayerclient "relayer/client"
)

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
			got := chunkSpans(tc.from, tc.to, tc.chunk)
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
	got := chunkSpans(10, 14, 1)
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

// The third NeedsChange branch: a provider that refuses the span.
//
// The scanner has no way to ask what the cap is — the comment on logScanChunk
// records the shape of the answer, "a plain 400 that names neither the setting
// nor a workable value" — so the only way to find a working span is to try a
// smaller one. Retrying the same one is not a retry at all: it is the identical
// refused request, forever, and on the L2 path that means every packet on the
// chain stops being seen.
func TestScanNarrowsWhenTheProviderRefusesTheSpan(t *testing.T) {
	const cap = uint64(10)
	span := &relayerclient.LogSpan{} // unconfigured: one call for the whole range
	var asked []uint64

	// A scan that never narrows does not fail, it spins: same width, same refusal,
	// forever. Bounding the calls turns that into a red test instead of a hung one
	// — an unbounded test here reports a broken narrowing as a wall-clock timeout,
	// which is indistinguishable from a slow machine.
	events, _, err := scanNarrowing(0, 99, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		width := to - from + 1
		asked = append(asked, width)
		if len(asked) > 64 {
			t.Fatalf("the scan asked 64 times without settling; widths so far: %v", asked)
		}
		if width > cap {
			return nil, nil, errors.New("query returned more than 10000 results")
		}
		return []chain.Event{{Height: from}}, nil, nil
	})
	if err != nil {
		t.Fatalf("scan never found a span the provider would serve: %v", err)
	}
	if span.Chunk == 0 || span.Chunk > cap {
		t.Fatalf("span settled at %d, want a value at or below the provider cap %d", span.Chunk, cap)
	}
	// It must have asked for something SMALLER each time it was refused; asking
	// the same width twice is the loop this branch exists to prevent.
	for i := 1; i < len(asked); i++ {
		if asked[i] >= asked[i-1] && asked[i-1] > cap {
			t.Fatalf("width %d did not shrink after %d was refused: %v", asked[i], asked[i-1], asked)
		}
	}
	// And the whole range is still covered once a working span is found: one event
	// per piece, and the pieces tile [0,99] exactly (the last one is short unless
	// the span divides 100).
	wantPieces := uint64(len(chunkSpans(0, 99, span.Chunk)))
	if uint64(len(events)) != wantPieces {
		t.Fatalf("got %d events over [0,99] at span %d, want %d; the range was not covered exactly",
			len(events), span.Chunk, wantPieces)
	}
}

// The narrowing must SURVIVE the call, or every tick pays one rejected request to
// rediscover the same cap and the log says nothing after the first line.
func TestNarrowedSpanIsRemembered(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 64}
	// Every scan in this test is bounded. A scan that does not narrow repeats the
	// same refused width forever, and an unbounded one turns that into a hang —
	// which reads as a slow machine, not as the defect it is.
	attempts := 0
	refuse := func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		attempts++
		if attempts > 64 {
			return nil, nil, errors.New("BOUND: the scan repeated a width it had already been refused")
		}
		if to-from+1 > 8 {
			return nil, nil, errors.New("block range too large")
		}
		return nil, nil, nil
	}
	if _, _, err := scanNarrowing(0, 63, span, refuse); err != nil {
		t.Fatalf("first scan: %v", err)
	}
	attempts = 0
	settled := span.Chunk
	if settled > 8 {
		t.Fatalf("span settled at %d, want <= 8", settled)
	}

	calls := 0
	if _, _, err := scanNarrowing(64, 127, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		calls++
		if calls > 64 {
			t.Fatal("the scan is repeating a width it already had refused")
		}
		return refuse(from, to)
	}); err != nil {
		t.Fatalf("second scan: %v", err)
	}
	if span.Chunk != settled {
		t.Fatalf("span moved to %d on a scan with no refusals, want it to stay at %d", span.Chunk, settled)
	}
	if wantCalls := 64 / int(settled); calls != wantCalls {
		t.Fatalf("second scan made %d calls, want %d — it re-discovered the cap instead of remembering it", calls, wantCalls)
	}
}

// The bottom of the ladder. A provider still refusing a ONE-block span is not a
// sizing problem, and halving it further is how a "retry" becomes an infinite
// loop. chain.Climb is what ends it, and it must end as permanent.
func TestScanStopsWhenTheSpanCannotShrinkFurther(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 4}
	// The bound is the half of this test that fails LOUDLY. Its own subject is
	// "the scan halves forever", and a scan that never narrows does not report a
	// wrong answer -- it reports none, and the assertions below are never
	// reached. 4 -> 2 -> 1 -> bottom is three rejections; anything past a
	// generous multiple of that is the loop this test exists to rule out.
	calls := 0
	_, _, err := scanNarrowing(0, 7, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		calls++
		if calls > 8 {
			t.Fatalf("the scan asked for [%d,%d] a %dth time without reaching the bottom of the ladder; it is not narrowing",
				from, to, calls)
		}
		return nil, nil, errors.New("block range too large")
	})
	if err == nil {
		t.Fatal("a provider refusing every span was reported as a successful scan")
	}
	if !chain.IsPermanent(err) {
		t.Fatalf("error = %v; a span that cannot shrink further must be permanent, or the scan halves forever", err)
	}
	if r, known := chain.Classify(err); !known || r.Outcome() != chain.OutcomePermanent {
		t.Fatalf("classified as %s, want permanent", r.Outcome())
	}
}

// A failure that is NOT a range rejection must pass straight through. Narrowing
// on a transport error would shrink the span for a reason that has nothing to do
// with size, and would end at the bottom of the ladder reporting permanent — a
// dropped scan for what was a reconnect.
func TestScanDoesNotNarrowOnAnOrdinaryFailure(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 32}
	wantErr := errors.New("dial tcp: connection refused")
	calls := 0
	_, _, err := scanNarrowing(0, 63, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		calls++
		return nil, nil, wantErr
	})
	if !errors.Is(err, wantErr) {
		t.Fatalf("error = %v, want the transport cause", err)
	}
	if chain.IsPermanent(err) {
		t.Fatalf("a transport failure was reported permanent: %v", err)
	}
	if span.Chunk != 32 {
		t.Fatalf("span narrowed to %d on a transport failure", span.Chunk)
	}
	if calls != 1 {
		t.Fatalf("made %d calls, want 1 — it retried a failure that retrying cannot fix here", calls)
	}
}

// Wrapping must not hide the rejection: the scanner wraps each piece with the
// span it was refused for before anything classifies it.
func TestRejectionIsFoundThroughTheSpanWrapper(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 8}
	calls := 0
	_, _, _ = scanNarrowing(0, 15, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		calls++
		if calls > 64 {
			t.Fatal("the scan never bottomed out")
		}
		return nil, nil, fmt.Errorf("filter SendPacket: %w", errors.New("logs matched by query exceeds limit of 10000"))
	})
	if span.Chunk == 8 {
		t.Fatal("the rejection was not recognised through the wrapping, so nothing narrowed")
	}
}

// #453 stated this invariant on the chunk loop it wrote and never asserted it:
// "One set across every span: a send in span 1 can be settled by a terminal log
// in span 3, and chunking must not hide that from the filter." That loop is now
// scanNarrowing, so the property has to be carried here -- and a stated
// invariant with nothing enforcing it is a missing assertion, not documentation.
//
// What it costs if it breaks: dropSettled sees only the last piece's keys, so a
// send whose AckPacket landed in a later piece is relayed for nothing and
// tracked again -- after the settle removed it and after the cursor moved past
// the terminal log that would have cleared it. Nothing settles it a second time.
func TestNarrowingUnionsSettledAcrossPieces(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 10}
	// One key per piece, named by the piece that produced it, so a set built from
	// only the last piece is distinguishable from a set built from all of them.
	_, settled, err := scanNarrowing(0, 29, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		return nil, map[settledKey]struct{}{settledKey(fmt.Sprintf("piece@%d", from)): {}}, nil
	})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	for _, want := range []settledKey{"piece@0", "piece@10", "piece@20"} {
		if _, ok := settled[want]; !ok {
			t.Fatalf("%q is missing from the settled set %v; a packet settled in one piece "+
				"would be relayed again and re-tracked with nothing left to clear it", want, settled)
		}
	}
}

// The union is per ATTEMPT, not per call. scanNarrowing restarts the WHOLE range
// after a narrowing, so keys from the attempt that failed describe pieces that
// are about to be read again -- and if that re-read no longer reports one (a
// reorg unwound the terminal event), the stale key would drop a packet that is
// live again from the batch.
//
// The scenario has to have a piece SUCCEED before the refusal, or it proves
// nothing: a piece that errors never reaches the union at all -- the loop breaks
// first -- so its keys are discarded whether the set is reset or hoisted. My
// first version of this test refused on the first piece and a hoisted set walked
// straight through it.
//
// [0,63] at span 32: piece [0,31] succeeds and reports "stale", then [32,63] is
// refused. Narrowed to 16, all four pieces succeed and report "fresh". Reset
// gives {fresh}; hoisted gives {stale, fresh}.
func TestNarrowingDoesNotCarrySettledAcrossARestart(t *testing.T) {
	span := &relayerclient.LogSpan{Chunk: 32}
	refused := false
	// Same bound, same reason as TestScanStopsWhenTheSpanCannotShrinkFurther: the
	// stub refuses every piece wider than 16, so a scan that does not narrow asks
	// forever and the assertions below are never reached. Two calls at span 32
	// plus four at 16 is six; the cap is a generous multiple.
	calls := 0
	_, settled, err := scanNarrowing(0, 63, span, func(from, to uint64) ([]chain.Event, map[settledKey]struct{}, error) {
		calls++
		if calls > 24 {
			t.Fatalf("the scan asked for [%d,%d] a %dth time; the span is not narrowing", from, to, calls)
		}
		if to-from+1 > 16 {
			if from == 0 {
				return nil, map[settledKey]struct{}{"unwound-by-the-reorg": {}}, nil
			}
			refused = true
			return nil, nil, errors.New("query returned more than 10000 results")
		}
		return nil, map[settledKey]struct{}{"still-settled": {}}, nil
	})
	if err != nil {
		t.Fatalf("scan: %v", err)
	}
	if !refused {
		t.Fatal("the provider was never refused, so no restart happened and this asserts nothing")
	}
	if _, stale := settled["unwound-by-the-reorg"]; stale {
		t.Fatal("a key from a piece that SUCCEEDED in the failed attempt survived the restart; " +
			"the set was hoisted out of the retry loop")
	}
	if _, ok := settled["still-settled"]; !ok {
		t.Fatal("the completed attempt's keys were lost")
	}
}
