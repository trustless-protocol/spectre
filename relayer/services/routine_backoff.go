package services

import "time"

const (
	routineFailureBackoffMin = time.Minute
	routineFailureBackoffMax = 15 * time.Minute
)

type routineBackoff struct {
	current     time.Duration
	nextAttempt time.Time
}

func (b *routineBackoff) Ready(now time.Time) bool {
	return b.nextAttempt.IsZero() || !now.Before(b.nextAttempt)
}

func (b *routineBackoff) RecordFailure(now time.Time) time.Duration {
	if b.current == 0 {
		b.current = routineFailureBackoffMin
	} else {
		b.current *= 2
		if b.current > routineFailureBackoffMax {
			b.current = routineFailureBackoffMax
		}
	}
	b.nextAttempt = now.Add(b.current)
	return b.current
}

func (b *routineBackoff) RecordSuccess() {
	b.current = 0
	b.nextAttempt = time.Time{}
}
