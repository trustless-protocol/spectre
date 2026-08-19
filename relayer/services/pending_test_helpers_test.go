package services

import "time"

// currentPacketForTest lets unit tests exercise the identity-checked mutation
// APIs without keeping the unsafe key-only production methods around.
func (t *PendingPacketTracker) currentPacketForTest(sourceClient string, sequence uint64) (pendingPacketInfo, bool) {
	for _, info := range t.GetAll() {
		if info.Packet.SourceClient == sourceClient && info.Packet.Sequence == sequence {
			return info, true
		}
	}
	return pendingPacketInfo{}, false
}

func (t *PendingPacketTracker) removeCurrentForTest(sourceClient string, sequence uint64) {
	if info, ok := t.currentPacketForTest(sourceClient, sequence); ok {
		t.RemoveIfCurrent(info)
	}
}

func (t *PendingPacketTracker) recordTimeoutFailureForTest(sourceClient string, sequence uint64, now time.Time) bool {
	info, ok := t.currentPacketForTest(sourceClient, sequence)
	return ok && t.RecordTimeoutFailureIfCurrent(info, now)
}

func (t *PendingPacketTracker) deferTimeoutRetryForTest(sourceClient string, sequence uint64, now time.Time) int {
	if info, ok := t.currentPacketForTest(sourceClient, sequence); ok {
		return t.DeferTimeoutRetryIfCurrent(info, now)
	}
	return 0
}

func (t *PendingPacketTracker) clearDeferralsForTest(sourceClient string, sequence uint64) {
	if info, ok := t.currentPacketForTest(sourceClient, sequence); ok {
		t.ClearDeferralsIfCurrent(info)
	}
}
