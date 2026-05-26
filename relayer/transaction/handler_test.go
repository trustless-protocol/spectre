package transaction

import (
	"bytes"
	"encoding/hex"
	"testing"

	contractICS26Router "relayer/bindings/ICS26Router"
)

// These tests cover the calldata packing layer that SendEthTxBatch relies on.
// We do not exercise the RPC submission path here — that requires a full
// services.Context + ethclient mock and is validated end-to-end via the
// benchmark run described in /Users/ducnt/.claude/plans/...
//
// Selector reference (from `cast sig` against the ICS26Router ABI):
//   recvPacket(...)              → 0x596e00b9
//   ackPacket(...)               → 0xfdbd955d
//   timeoutPacket(...)           → 0x223e357a
//   updateClient(string,bytes)   → 0x6fbf8079
//   multicall(bytes[])           → 0xac9650d8

func mustDecodeHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("decode %q: %v", s, err)
	}
	return b
}

func samplePacket() contractICS26Router.IICS26RouterMsgsPacket {
	return contractICS26Router.IICS26RouterMsgsPacket{
		Sequence:         42,
		SourceClient:     "cosmoshub-1",
		DestClient:       "08-wasm-0",
		TimeoutTimestamp: 1_700_000_000,
		Payloads: []contractICS26Router.IICS26RouterMsgsPayload{{
			SourcePort:  "transfer",
			DestPort:    "transfer",
			Version:     "ics20-2",
			Encoding:    "abi",
			Value:       []byte{0x01, 0x02, 0x03},
		}},
	}
}

func TestSelectorsForBatchedMsgs(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	cases := []struct {
		name     string
		msg      any
		method   string
		selector string
	}{
		{
			name: "recvPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
				Packet:        pkt,
				MembershipMsg: []byte{0xaa, 0xbb},
			},
			method:   "recvPacket",
			selector: "596e00b9",
		},
		{
			name: "ackPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgAckPacket{
				Packet:          pkt,
				Acknowledgement: []byte{0xde, 0xad},
				MembershipMsg:   []byte{0xbe, 0xef},
			},
			method:   "ackPacket",
			selector: "fdbd955d",
		},
		{
			name: "timeoutPacket",
			msg: contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket{
				Packet:           pkt,
				NonMembershipMsg: []byte{0xca, 0xfe},
			},
			method:   "timeoutPacket",
			selector: "223e357a",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := parsedABI.Pack(tc.method, tc.msg)
			if err != nil {
				t.Fatalf("pack: %v", err)
			}
			if len(data) < 4 {
				t.Fatalf("calldata too short: %d", len(data))
			}
			gotSelector := hex.EncodeToString(data[:4])
			if gotSelector != tc.selector {
				t.Fatalf("selector mismatch: got %s, want %s", gotSelector, tc.selector)
			}
		})
	}
}

// TestMulticallWraps verifies that wrapping per-msg calldata into a multicall
// produces the expected outer selector + that each inner blob is recoverable
// after Unpack. This is the round-trip that on-chain MulticallUpgradeable
// would do.
func TestMulticallWraps(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	recv, err := parsedABI.Pack("recvPacket", contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet:        pkt,
		MembershipMsg: []byte{0x01},
	})
	if err != nil {
		t.Fatalf("pack recv: %v", err)
	}
	ack, err := parsedABI.Pack("ackPacket", contractICS26Router.IICS26RouterMsgsMsgAckPacket{
		Packet:          pkt,
		Acknowledgement: []byte{0x02},
		MembershipMsg:   []byte{0x03},
	})
	if err != nil {
		t.Fatalf("pack ack: %v", err)
	}

	innerCalldata := [][]byte{recv, ack}
	outer, err := parsedABI.Pack("multicall", innerCalldata)
	if err != nil {
		t.Fatalf("pack multicall: %v", err)
	}
	wantSelector := mustDecodeHex(t, "ac9650d8")
	if !bytes.Equal(outer[:4], wantSelector) {
		t.Fatalf("multicall selector mismatch: got %x, want %x", outer[:4], wantSelector)
	}

	// Decode the outer multicall payload and verify we recover the same two
	// inner calldata blobs in order.
	method, err := parsedABI.MethodById(outer[:4])
	if err != nil {
		t.Fatalf("MethodById: %v", err)
	}
	if method.Name != "multicall" {
		t.Fatalf("unexpected method: %s", method.Name)
	}
	unpacked, err := method.Inputs.Unpack(outer[4:])
	if err != nil {
		t.Fatalf("unpack inputs: %v", err)
	}
	if len(unpacked) != 1 {
		t.Fatalf("expected 1 top-level arg, got %d", len(unpacked))
	}
	gotInner, ok := unpacked[0].([][]byte)
	if !ok {
		t.Fatalf("unexpected inner type: %T", unpacked[0])
	}
	if len(gotInner) != 2 {
		t.Fatalf("expected 2 inner calls, got %d", len(gotInner))
	}
	if !bytes.Equal(gotInner[0], recv) {
		t.Fatalf("inner[0] mismatch")
	}
	if !bytes.Equal(gotInner[1], ack) {
		t.Fatalf("inner[1] mismatch")
	}
}

// TestSelectorForUpdateClient pins the selector for ICS26Router.updateClient,
// which V2 batches into the same multicall as recvPacket/ackPacket/
// timeoutPacket.
func TestSelectorForUpdateClient(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	data, err := parsedABI.Pack("updateClient", "cosmoshub-1", []byte{0xde, 0xad, 0xbe, 0xef})
	if err != nil {
		t.Fatalf("pack: %v", err)
	}
	if len(data) < 4 {
		t.Fatalf("calldata too short: %d", len(data))
	}
	got := hex.EncodeToString(data[:4])
	if got != "6fbf8079" {
		t.Fatalf("selector mismatch: got %s, want 6fbf8079", got)
	}
}

// TestMulticallWithUpdateClient verifies that an updateClient + packet combo
// produces a valid multicall payload — the layout we expect from the V2
// handleCosmos refactor.
func TestMulticallWithUpdateClient(t *testing.T) {
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		t.Fatalf("GetAbi: %v", err)
	}
	pkt := samplePacket()

	updateData, err := parsedABI.Pack("updateClient", "cosmoshub-1", []byte{0x01, 0x02})
	if err != nil {
		t.Fatalf("pack updateClient: %v", err)
	}
	recvData, err := parsedABI.Pack("recvPacket", contractICS26Router.IICS26RouterMsgsMsgRecvPacket{
		Packet:        pkt,
		MembershipMsg: []byte{0x03},
	})
	if err != nil {
		t.Fatalf("pack recvPacket: %v", err)
	}

	outer, err := parsedABI.Pack("multicall", [][]byte{updateData, recvData})
	if err != nil {
		t.Fatalf("pack multicall: %v", err)
	}
	if !bytes.Equal(outer[:4], mustDecodeHex(t, "ac9650d8")) {
		t.Fatalf("multicall selector mismatch: got %x", outer[:4])
	}

	method, err := parsedABI.MethodById(outer[:4])
	if err != nil {
		t.Fatalf("MethodById: %v", err)
	}
	unpacked, err := method.Inputs.Unpack(outer[4:])
	if err != nil {
		t.Fatalf("unpack: %v", err)
	}
	inner := unpacked[0].([][]byte)
	if len(inner) != 2 {
		t.Fatalf("expected 2 inner calls, got %d", len(inner))
	}
	// First inner must be updateClient — atomicity demands this so packet
	// proofs verify against the just-applied client state.
	if !bytes.Equal(inner[0][:4], mustDecodeHex(t, "6fbf8079")) {
		t.Fatalf("inner[0] is not updateClient: %x", inner[0][:4])
	}
	if !bytes.Equal(inner[1][:4], mustDecodeHex(t, "596e00b9")) {
		t.Fatalf("inner[1] is not recvPacket: %x", inner[1][:4])
	}
}
