package l2rollup

import (
	"bytes"
	"encoding/binary"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// buildAfterState lays out the 6-word AssertionState exactly as the RollupCore ABI
// encodes it, so the test exercises decodeAssertionCreated's offsets independently of
// the builder.
func buildAfterState(blockHash, sendRoot, endHistoryRoot ethcommon.Hash, inboxPos, posInMsg uint64, status byte) []byte {
	afterState := make([]byte, 6*evmWord)
	copy(afterState[0:evmWord], blockHash.Bytes())
	copy(afterState[evmWord:2*evmWord], sendRoot.Bytes())
	binary.BigEndian.PutUint64(afterState[3*evmWord-8:3*evmWord], inboxPos)
	binary.BigEndian.PutUint64(afterState[4*evmWord-8:4*evmWord], posInMsg)
	afterState[5*evmWord-1] = status
	copy(afterState[5*evmWord:6*evmWord], endHistoryRoot.Bytes())
	return afterState
}

// buildAssertionCreatedLog assembles a well-formed AssertionCreated log whose topic-1
// hash is the genuine BoLD assertion hash of the supplied fields.
func buildAssertionCreatedLog(parent, inboxAcc ethcommon.Hash, afterState []byte) types.Log {
	data := make([]byte, assertionCreatedDataWords*evmWord)
	copy(data[13*evmWord:19*evmWord], afterState)
	copy(data[19*evmWord:20*evmWord], inboxAcc.Bytes())
	assertionHash := boldAssertionHash(parent, afterState, inboxAcc)
	return types.Log{
		Topics: []ethcommon.Hash{assertionCreatedTopic, assertionHash, parent},
		Data:   data,
	}
}

func TestAssertionStorageSlot(t *testing.T) {
	assertionHash := ethcommon.HexToHash("0x1111111111111111111111111111111111111111111111111111111111111111")
	mappingSlot := ethcommon.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000075")

	got := assertionStorageSlot(assertionHash, mappingSlot)
	want := crypto.Keccak256Hash(assertionHash.Bytes(), mappingSlot.Bytes())
	if got != want {
		t.Fatalf("assertionStorageSlot = %s, want keccak(key||slot) = %s", got, want)
	}
}

func TestDecodeAssertionCreatedRoundTrip(t *testing.T) {
	parent := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000aa")
	blockHash := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000bb")
	sendRoot := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000cc")
	endHistoryRoot := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000dd")
	inboxAcc := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000ee")
	const inboxPos, posInMsg uint64 = 42, 7

	afterState := buildAfterState(blockHash, sendRoot, endHistoryRoot, inboxPos, posInMsg, 1) // finished
	log := buildAssertionCreatedLog(parent, inboxAcc, afterState)

	claim, gotBlockHash, gotAssertionHash, err := decodeAssertionCreated(log)
	if err != nil {
		t.Fatalf("decodeAssertionCreated: %v", err)
	}

	if gotBlockHash != blockHash {
		t.Errorf("blockHash = %s, want %s", gotBlockHash, blockHash)
	}
	if gotAssertionHash != log.Topics[1] {
		t.Errorf("assertionHash = %s, want %s", gotAssertionHash, log.Topics[1])
	}
	if !bytes.Equal(claim.ParentAssertionHash, parent.Bytes()) {
		t.Errorf("parent = %x, want %x", claim.ParentAssertionHash, parent.Bytes())
	}
	gs := claim.AfterState.GlobalState
	if !bytes.Equal(gs.Bytes32Vals[0], blockHash.Bytes()) || !bytes.Equal(gs.Bytes32Vals[1], sendRoot.Bytes()) {
		t.Errorf("global state bytes32 = %x, want [%x %x]", gs.Bytes32Vals, blockHash.Bytes(), sendRoot.Bytes())
	}
	if gs.U64Vals[0] != inboxPos || gs.U64Vals[1] != posInMsg {
		t.Errorf("u64 vals = %v, want [%d %d]", gs.U64Vals, inboxPos, posInMsg)
	}
	if claim.AfterState.MachineStatus != MachineStatusFinished {
		t.Errorf("machine status = %q, want finished", claim.AfterState.MachineStatus)
	}
	if !bytes.Equal(claim.AfterState.EndHistoryRoot, endHistoryRoot.Bytes()) {
		t.Errorf("end history root = %x, want %x", claim.AfterState.EndHistoryRoot, endHistoryRoot.Bytes())
	}
	if !bytes.Equal(claim.InboxAcc, inboxAcc.Bytes()) {
		t.Errorf("inbox acc = %x, want %x", claim.InboxAcc, inboxAcc.Bytes())
	}
}

func TestDecodeAssertionCreatedRejectsHashMismatch(t *testing.T) {
	parent := ethcommon.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000aa")
	afterState := buildAfterState(
		ethcommon.HexToHash("0xbb"), ethcommon.HexToHash("0xcc"), ethcommon.HexToHash("0xdd"), 1, 2, 1,
	)
	log := buildAssertionCreatedLog(parent, ethcommon.HexToHash("0xee"), afterState)
	// Corrupt the committed block hash without recomputing the topic → hash self-check must fail.
	copy(log.Data[13*evmWord:14*evmWord], ethcommon.HexToHash("0xdeadbeef").Bytes())

	if _, _, _, err := decodeAssertionCreated(log); err == nil {
		t.Fatal("expected hash-mismatch error, got nil")
	}
}

func TestDecodeAssertionCreatedRejectsMalformed(t *testing.T) {
	// Wrong topic count.
	if _, _, _, err := decodeAssertionCreated(types.Log{Topics: []ethcommon.Hash{assertionCreatedTopic}}); err == nil {
		t.Error("expected error for wrong topic count")
	}
	// Wrong data length.
	shortData := types.Log{
		Topics: []ethcommon.Hash{assertionCreatedTopic, {}, {}},
		Data:   make([]byte, 10*evmWord),
	}
	if _, _, _, err := decodeAssertionCreated(shortData); err == nil {
		t.Error("expected error for wrong data length")
	}
	// Unknown machine status byte.
	afterState := buildAfterState(ethcommon.HexToHash("0xbb"), ethcommon.HexToHash("0xcc"), ethcommon.HexToHash("0xdd"), 0, 0, 9)
	badStatus := buildAssertionCreatedLog(ethcommon.HexToHash("0xaa"), ethcommon.HexToHash("0xee"), afterState)
	if _, _, _, err := decodeAssertionCreated(badStatus); err == nil {
		t.Error("expected error for unknown machine status")
	}
}

func TestMachineStatusFromByte(t *testing.T) {
	cases := map[byte]MachineStatus{0: MachineStatusRunning, 1: MachineStatusFinished, 2: MachineStatusErrored}
	for b, want := range cases {
		got, err := machineStatusFromByte(b)
		if err != nil || got != want {
			t.Errorf("machineStatusFromByte(%d) = (%q, %v), want (%q, nil)", b, got, err, want)
		}
	}
	if _, err := machineStatusFromByte(3); err == nil {
		t.Error("expected error for status 3")
	}
}
