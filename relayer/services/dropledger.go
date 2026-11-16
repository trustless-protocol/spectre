package services

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"relayer/chain"
)

// Bounds of the drop ledger.
const (
	// dropBackoffMin is the hold after the first drop. Two flush passes, so a
	// refusal caused by a stale node clears quickly.
	dropBackoffMin = 10 * time.Minute
	// dropBackoffMax is the longest hold: a packet that keeps failing costs at
	// most one retry (and the client update that precedes it) per this interval.
	dropBackoffMax = 6 * time.Hour
	// dropLedgerMaxEntries bounds the file. The oldest entry is evicted first.
	dropLedgerMaxEntries = 4096
	// dropLedgerMaxAge forgets an entry nothing has offered again for this long.
	dropLedgerMaxAge = 30 * 24 * time.Hour
	dropReasonMaxLen = 512

	dropLedgerVersion = 1
)

// DroppedPacket is one packet the relay gave up on after a permanent failure.
type DroppedPacket struct {
	Type     chain.EventType `json:"type"`
	Sequence uint64          `json:"sequence"`
	// Packet is the adapter-encoded packet (chain.Event.Raw). It is kept so a
	// later pass can ask the destination whether the packet settled meanwhile.
	Packet         []byte    `json:"packet"`
	FirstDroppedAt time.Time `json:"first_dropped_at"`
	LastDroppedAt  time.Time `json:"last_dropped_at"`
	Attempts       int       `json:"attempts"`
	NextRetryAt    time.Time `json:"next_retry_at"`
	Reason         string    `json:"reason"`
}

type dropLedgerState struct {
	Version int             `json:"version"`
	Dropped []DroppedPacket `json:"dropped"`
}

type dropKey struct {
	typ    chain.EventType
	digest [sha256.Size]byte
}

// DropLedger is the durable record of packets the relay dropped after a
// permanent failure.
//
// A dropped packet is not gone: its commitment is still outstanding on the
// source, so the enumeration backstop lists it again on the next pass. With no
// record of the drop, every pass advanced the client and re-submitted a packet
// already shown to fail (#575 M2). The ledger holds a dropped packet out of the
// relay path until its backoff expires.
//
// The backoff is bounded, not final. A dropped acknowledgement or timeout has no
// other exit, so it is retried for as long as it stays outstanding, only less
// often: 10 min, doubling, capped at 6 h.
//
// Safe for concurrent use.
type DropLedger struct {
	mu      sync.Mutex
	path    string // "" keeps the ledger in memory only
	entries map[dropKey]DroppedPacket
	write   func(path string, data []byte) error
}

// NewDropLedger loads the ledger at path, or starts an empty one when the file
// does not exist. An empty path keeps the ledger in memory only.
func NewDropLedger(path string) (*DropLedger, error) {
	ledger := &DropLedger{path: path, entries: make(map[dropKey]DroppedPacket), write: writeStateFile}
	if path == "" {
		return ledger, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return ledger, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read dropped-packet ledger %s: %w", path, err)
	}
	var state dropLedgerState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("decode dropped-packet ledger %s: %w", path, err)
	}
	if state.Version != dropLedgerVersion {
		return nil, fmt.Errorf("dropped-packet ledger %s has version %d, want %d", path, state.Version, dropLedgerVersion)
	}
	cutoff := time.Now().Add(-dropLedgerMaxAge)
	for _, entry := range state.Dropped {
		if entry.LastDroppedAt.Before(cutoff) {
			continue
		}
		ledger.entries[keyOfDrop(entry.Type, entry.Packet)] = entry
	}
	return ledger, nil
}

// DropLedgerFileName is the ledger file of one relay direction inside a relay
// state directory. direction is a file-safe name such as "cosmos-to-eth".
func DropLedgerFileName(direction string) string {
	return "dropped-" + direction + ".json"
}

// OpenDropLedger opens the drop ledger of one relay direction in this instance's
// state directory. Without a state directory the ledger lives in memory only.
func (s *Services) OpenDropLedger(direction string) (*DropLedger, error) {
	if s.stateDir == "" {
		return NewDropLedger("")
	}
	return NewDropLedger(filepath.Join(s.stateDir, DropLedgerFileName(direction)))
}

func keyOfDrop(typ chain.EventType, packet []byte) dropKey {
	return dropKey{typ: typ, digest: sha256.Sum256(packet)}
}

// dropBackoff is the hold after the attempts-th drop.
func dropBackoff(attempts int) time.Duration {
	hold := dropBackoffMin
	for i := 1; i < attempts && hold < dropBackoffMax; i++ {
		hold *= 2
	}
	if hold > dropBackoffMax {
		hold = dropBackoffMax
	}
	return hold
}

// Held reports whether the packet was dropped and is still inside its backoff,
// and when it may be retried.
func (l *DropLedger) Held(typ chain.EventType, packet []byte, now time.Time) (time.Time, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	entry, ok := l.entries[keyOfDrop(typ, packet)]
	if !ok || !now.Before(entry.NextRetryAt) {
		return time.Time{}, false
	}
	return entry.NextRetryAt, true
}

// Record notes a drop and returns the resulting entry. Each drop of the same
// packet lengthens its hold.
//
// The entry is kept in memory even when the write fails, and the error is
// returned: the hold then lasts until the process restarts, which costs one
// extra retry rather than a retry on every pass.
func (l *DropLedger) Record(typ chain.EventType, sequence uint64, packet []byte, reason string, now time.Time) (DroppedPacket, error) {
	if len(reason) > dropReasonMaxLen {
		reason = reason[:dropReasonMaxLen] + "..."
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	key := keyOfDrop(typ, packet)
	entry, known := l.entries[key]
	if !known {
		entry = DroppedPacket{
			Type:           typ,
			Sequence:       sequence,
			Packet:         append([]byte(nil), packet...),
			FirstDroppedAt: now,
		}
		l.evictOldestLocked()
	}
	entry.Attempts++
	entry.LastDroppedAt = now
	entry.NextRetryAt = now.Add(dropBackoff(entry.Attempts))
	entry.Reason = reason
	l.entries[key] = entry
	return entry, l.persistLocked()
}

// evictOldestLocked makes room for one more entry.
func (l *DropLedger) evictOldestLocked() {
	if len(l.entries) < dropLedgerMaxEntries {
		return
	}
	var (
		oldestKey dropKey
		oldest    time.Time
		found     bool
	)
	for key, entry := range l.entries {
		if !found || entry.LastDroppedAt.Before(oldest) {
			oldestKey, oldest, found = key, entry.LastDroppedAt, true
		}
	}
	if found {
		delete(l.entries, oldestKey)
	}
}

// Clear forgets a packet: it was delivered, or the destination reports it
// settled. It reports whether there was an entry to forget.
func (l *DropLedger) Clear(typ chain.EventType, packet []byte) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	key := keyOfDrop(typ, packet)
	if _, ok := l.entries[key]; !ok {
		return false, nil
	}
	delete(l.entries, key)
	return true, l.persistLocked()
}

// Snapshot lists every dropped packet, oldest drop first.
func (l *DropLedger) Snapshot() []DroppedPacket {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.sortedLocked(func(DroppedPacket) bool { return true })
}

// Due lists the packets whose backoff has expired.
func (l *DropLedger) Due(now time.Time) []DroppedPacket {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.sortedLocked(func(entry DroppedPacket) bool { return !now.Before(entry.NextRetryAt) })
}

func (l *DropLedger) sortedLocked(keep func(DroppedPacket) bool) []DroppedPacket {
	out := make([]DroppedPacket, 0, len(l.entries))
	for _, entry := range l.entries {
		if keep(entry) {
			out = append(out, entry)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].FirstDroppedAt.Equal(out[j].FirstDroppedAt) {
			return out[i].FirstDroppedAt.Before(out[j].FirstDroppedAt)
		}
		if out[i].Type != out[j].Type {
			return out[i].Type < out[j].Type
		}
		return out[i].Sequence < out[j].Sequence
	})
	return out
}

func (l *DropLedger) persistLocked() error {
	if l.path == "" {
		return nil
	}
	state := dropLedgerState{Version: dropLedgerVersion, Dropped: l.sortedLocked(func(DroppedPacket) bool { return true })}
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("encode dropped-packet ledger: %w", err)
	}
	if err := l.write(l.path, data); err != nil {
		return fmt.Errorf("persist dropped-packet ledger %s: %w", l.path, err)
	}
	return nil
}
