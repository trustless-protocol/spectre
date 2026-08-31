package services

import (
	"testing"
	"time"
)

// NOTE ON MERGE ORDER: PR #424 also creates services/services_test.go, for the
// ICS-24 addressing tests. The two additions are disjoint, so whichever lands
// second resolves the conflict by concatenating -- there is nothing to reconcile.
// Both belong here: E7 pairs a test file with exactly one production file, and
// both cover services.go.

// tendermintClientExpiry is the mirror of ethClientExpiry in chain/cosmos, and
// the two are shaped differently on purpose: a Tendermint client carries one
// explicit trusting period, while a beacon client carries none and its expiry has
// to be inferred from sync-committee geometry.
//
// What they share is the clock. Both measure from the state the client currently
// TRUSTS, never from now.
func TestTendermintClientExpiry(t *testing.T) {
	trustedAt := time.Unix(1_700_000_000, 0).UTC()

	t.Run("trusting period after the trusted header time", func(t *testing.T) {
		got := tendermintClientExpiry(trustedAt, 1209600) // 14 days
		want := trustedAt.Add(14 * 24 * time.Hour)
		if !got.Equal(want) {
			t.Fatalf("expiry = %s, want %s", got, want)
		}
	})

	t.Run("measured from the trusted header, not from now", func(t *testing.T) {
		// A value derived from wall-clock time would report every client as
		// healthy forever, because it moves with the check. This is the property
		// the anti-expiry refresh depends on.
		const period = uint32(3600)
		older := tendermintClientExpiry(trustedAt, period)
		newer := tendermintClientExpiry(trustedAt.Add(time.Hour), period)
		if gap := newer.Sub(older); gap != time.Hour {
			t.Fatalf("an hour of header progress moved the expiry by %s, want 1h", gap)
		}
	})

	t.Run("the period is seconds", func(t *testing.T) {
		// The field is uint32 seconds on the wire. Reading it as any other unit
		// silently mis-sizes the refresh window -- as nanoseconds a 14-day period
		// becomes 1.2ms and the client is always "about to expire".
		if got := tendermintClientExpiry(trustedAt, 60); !got.Equal(trustedAt.Add(time.Minute)) {
			t.Fatalf("60 gave %s, want one minute after %s", got, trustedAt)
		}
	})

	t.Run("a zero trusting period expires at the trusted header itself", func(t *testing.T) {
		// Not treated as "no expiry": zero here means the client state is
		// unusable, and reporting it as already expired sends the refresh routine
		// at it immediately rather than letting it sit unnoticed.
		if got := tendermintClientExpiry(trustedAt, 0); !got.Equal(trustedAt) {
			t.Fatalf("expiry = %s, want the trusted header time %s", got, trustedAt)
		}
	})
}
