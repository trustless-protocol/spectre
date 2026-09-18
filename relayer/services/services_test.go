package services

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	client "relayer/client"
	"strings"
	"sync"
	"testing"
	"time"

	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

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

// --- the ICS-24 addressing the ETH side proves against ---
//
// EthPath and EthIBCStorageKey decide WHICH storage slot a membership proof is
// taken at. Get them wrong and the proof is well-formed but about the wrong
// location: it either fails verification, or -- worse for a non-membership proof
// -- succeeds in showing that some unrelated slot is empty.
//
// Both are a cross-language contract with contracts/core/libraries/ICS24Host.sol, so the
// expectations here are transcribed from that file rather than from the Go code
// they check. A test that derived them the same way the code does would agree
// with any mistake.

// The layout is abi.encodePacked(clientId, uint8(pathType), uint64BigEndian(seq))
// — ICS24Host.sol:25-88. Path types come from the same three functions:
//
//	1  packetCommitmentPathCalldata
//	2  packetReceiptCommitmentPathCalldata
//	3  packetAcknowledgementCommitmentPathCalldata
func TestEthPath(t *testing.T) {
	const (
		pathCommitment = byte(1)
		pathReceipt    = byte(2)
		pathAck        = byte(3)
	)

	t.Run("matches the Solidity packed layout", func(t *testing.T) {
		tests := []struct {
			name     string
			clientID string
			sequence uint64
			pathType byte
			want     string // hex, transcribed from the Solidity encoding
		}{
			{
				name: "commitment", clientID: "client-0", sequence: 1, pathType: pathCommitment,
				// "client-0" = 636c69656e742d30, then 01, then 8-byte BE 1
				want: "636c69656e742d30" + "01" + "0000000000000001",
			},
			{
				name: "receipt", clientID: "client-0", sequence: 1, pathType: pathReceipt,
				want: "636c69656e742d30" + "02" + "0000000000000001",
			},
			{
				name: "acknowledgement", clientID: "client-0", sequence: 1, pathType: pathAck,
				want: "636c69656e742d30" + "03" + "0000000000000001",
			},
			{
				name: "big-endian sequence", clientID: "c", sequence: 0x0102030405060708, pathType: pathCommitment,
				// Byte order is the whole point: little-endian would address a
				// completely different slot and still look like a valid proof.
				want: "63" + "01" + "0102030405060708",
			},
			{
				name: "max sequence", clientID: "c", sequence: ^uint64(0), pathType: pathCommitment,
				want: "63" + "01" + "ffffffffffffffff",
			},
			{
				name: "empty client id", clientID: "", sequence: 7, pathType: pathReceipt,
				want: "" + "02" + "0000000000000007",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := hex.EncodeToString(EthPath(tt.clientID, tt.sequence, tt.pathType))
				if got != tt.want {
					t.Fatalf("path = %s, want %s", got, tt.want)
				}
			})
		}
	})

	t.Run("does not alias the caller's client id", func(t *testing.T) {
		// EthPath appends to []byte(clientID). If that conversion ever stopped
		// copying, a second call could scribble over the first call's path — and
		// the two calls in a recv/ack pair happen microseconds apart.
		first := EthPath("client-0", 1, 1)
		second := EthPath("client-0", 2, 1)
		if hex.EncodeToString(first) == hex.EncodeToString(second) {
			t.Fatal("two paths at different sequences are identical")
		}
		if got, want := hex.EncodeToString(first), "636c69656e742d30010000000000000001"; got != want {
			t.Fatalf("the first path changed after a second call: %s, want %s", got, want)
		}
	})
}

// The storage slot is keccak256(keccak256(path) ‖ IBCSTORE_STORAGE_SLOT) — the
// standard Solidity mapping-slot derivation against the ERC-7201 store root in
// contracts/core/store/IBCStoreUpgradeable.sol:26.
func TestEthIBCStorageKey(t *testing.T) {
	t.Run("the storage root matches the Solidity constant", func(t *testing.T) {
		// Transcribed from IBCStoreUpgradeable.sol:26. If this fails, the two
		// languages disagree about where the IBC store lives and every ETH proof
		// is taken at the wrong root.
		const solidityIBCStoreSlot = "0x1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600"
		if ICS26_IBC_STORAGE_SLOT != solidityIBCStoreSlot {
			t.Fatalf("Go storage slot %s does not match Solidity %s", ICS26_IBC_STORAGE_SLOT, solidityIBCStoreSlot)
		}
	})

	t.Run("distinct paths address distinct slots", func(t *testing.T) {
		// The three path types share a client id and a sequence, so if the type
		// byte were dropped from the hash all three would collide -- a receipt
		// proof would answer a commitment query.
		seen := map[string]string{}
		for _, pathType := range []byte{1, 2, 3} {
			key := EthIBCStorageKey(EthPath("client-0", 42, pathType)).Hex()
			if prev, dup := seen[key]; dup {
				t.Fatalf("path type %d addresses the same slot as %s", pathType, prev)
			}
			seen[key] = string(rune('0' + pathType))
		}
		if len(seen) != 3 {
			t.Fatalf("3 path types produced %d distinct slots", len(seen))
		}
	})

	t.Run("derives the slot Solidity would", func(t *testing.T) {
		// A known-answer test, because the two structural checks below hold even
		// with the derivation gutted: dropping the keccak of the path, or leaving
		// the store root out of the second hash, still yields distinct and
		// deterministic slots. Mutation caught exactly that.
		//
		// The expected value is computed independently with Foundry's cast, not by
		// this package's own code:
		//
		//	P=636c69656e742d3001000000000000002a          # client-0 ‖ 0x01 ‖ BE(42)
		//	H=$(cast keccak "0x$P")
		//	cast keccak "${H}1260944489272988d9df285149b5aa1b0f48f2136d6f416159f840a3e0747600"
		const want = "0x0596148b513b849040f1f6302b374a4442916c30ef7088676dc098b3a0a706b5"
		got := EthIBCStorageKey(EthPath("client-0", 42, 1)).Hex()
		if got != want {
			t.Fatalf("storage key = %s, want %s\n(recompute with cast before changing this: the value is the contract with Solidity)", got, want)
		}
	})

	t.Run("the same path is always the same slot", func(t *testing.T) {
		// Not a tautology: the derivation reads a package-level constant, and a
		// caller-order dependence there would make proofs intermittently wrong.
		a := EthIBCStorageKey(EthPath("client-0", 42, 1))
		b := EthIBCStorageKey(EthPath("client-0", 42, 1))
		if a != b {
			t.Fatalf("the same path derived two slots: %s and %s", a.Hex(), b.Hex())
		}
	})
}

// ToEthPacket copies a packet field by field between two type systems, the same
// shape of code as the misbehaviour header conversion — and with the same failure
// mode: a field left out does not fail to compile, it just arrives zero and the
// router rejects the packet on chain.
//
// So this walks both structs and compares VALUES, field by field, rather than
// listing fields by hand and repeating the omission risk one level up.
//
// The first version asserted non-zeroness plus a hand-written port check, and
// the fixture used "transfer" for both ports — which made that check read
// "transfer" != "transfer" and pass under any swap. Review caught it; a mutation
// swapping SourcePort and DestPort in ToEthPacket survived the whole suite.
// Every field below therefore carries a value distinct from every other field of
// its type, so a conversion that pairs the wrong two cannot come out equal.
func TestToEthPacket_CarriesEveryField(t *testing.T) {
	src := channeltypesv2.Packet{
		Sequence:          7,
		SourceClient:      "cosmos-client-0",
		DestinationClient: "eth-router-0",
		TimeoutTimestamp:  1_700_000_000,
		Payloads: []channeltypesv2.Payload{{
			SourcePort:      "transfer-src",
			DestinationPort: "transfer-dst",
			Version:         "ics20-2",
			Encoding:        "application/x-solidity-abi",
			Value:           []byte{0x01, 0x02},
		}},
	}
	// Guard the fixture itself: a field added to the proto packet that this test
	// does not set would leave the conversion of it unchecked.
	assertNoZeroFields(t, reflect.ValueOf(src), "source packet")

	got := ToEthPacket(src)

	// ToEthPacket renames two fields on the way across, so the walk is told the
	// mapping rather than being loosened to skip them.
	assertFieldsCarriedAcross(t, reflect.ValueOf(src), reflect.ValueOf(got), "Packet", map[string]string{
		"DestinationClient": "DestClient",
		"DestinationPort":   "DestPort",
	})
}

// An empty payload list must stay empty rather than becoming nil: the ABI encoder
// treats the two the same, but a nil slice where the caller expected a length is
// how an off-by-one creeps into batch building.
func TestToEthPacket_EmptyPayloads(t *testing.T) {
	got := ToEthPacket(channeltypesv2.Packet{Sequence: 1})
	if got.Payloads == nil {
		t.Fatal("payloads is nil, want an empty slice")
	}
	if len(got.Payloads) != 0 {
		t.Fatalf("payloads = %d, want 0", len(got.Payloads))
	}
}

// TestShouldRelayCosmosTimeoutToEth pins the direction-detection logic that
// decides whether a CosmosTimeout event needs a follow-up timeoutPacket on
// ETH. The old guard compared SourceClient against the router client ID — a
// mismatch in IBC v2 because:
//
//   - For a Cosmos→ETH packet, the Cosmos-emitted event has
//     SourceClient="08-wasm-0" (Cosmos's client identifying the ETH chain) and
//     DestinationClient="cosmoshub-1" (ETH's client identifying Cosmos).
//   - For an ETH→Cosmos packet, it's mirrored: source=cosmoshub-1, dest=08-wasm-0.
//
// `CosmosRouterClientID` in the relayer config is the ETH-side router's name
// for this Cosmos chain ("cosmoshub-1"), so the correct field to compare is
// `DestinationClient`. The old comparison silently inverted both branches.
func TestShouldRelayCosmosTimeoutToEth(t *testing.T) {
	tests := []struct {
		name                 string
		packet               *channeltypesv2.Packet
		cosmosRouterClientID string
		want                 bool
	}{
		{
			name: "Cosmos-originated packet skips ETH relay",
			packet: &channeltypesv2.Packet{
				SourceClient:      "08-wasm-0",   // Cosmos's client for ETH
				DestinationClient: "cosmoshub-1", // ETH's client for Cosmos
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false, // Cosmos's MsgTimeout already refunded the sender
		},
		{
			name: "ETH-originated packet is relayed",
			packet: &channeltypesv2.Packet{
				SourceClient:      "cosmoshub-1", // ETH's client for Cosmos
				DestinationClient: "08-wasm-0",   // Cosmos's client for ETH
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true, // ETH still holds the commitment, needs ICS26Router.timeoutPacket
		},
		{
			name:                 "nil packet is not relayed",
			packet:               nil,
			cosmosRouterClientID: "cosmoshub-1",
			want:                 false,
		},
		{
			name: "unrelated client IDs are relayed (treated as ETH-originated)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "cosmoshub-1",
			want:                 true,
		},
		{
			name: "empty cosmosRouterClientID is not relayed (misconfiguration)",
			packet: &channeltypesv2.Packet{
				SourceClient:      "other-client",
				DestinationClient: "yet-another",
			},
			cosmosRouterClientID: "",
			want:                 false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ShouldRelayCosmosTimeoutToEth(tc.packet, tc.cosmosRouterClientID)
			if got != tc.want {
				t.Fatalf("ShouldRelayCosmosTimeoutToEth=%v want=%v", got, tc.want)
			}
		})
	}
}

func TestScanForEVMTimeouts(t *testing.T) {
	t.Run("tracker removal", func(t *testing.T) {
		expiredPacket := channeltypesv2.Packet{
			SourceClient:      "l2-client-0",
			DestinationClient: "08-wasm-1",
			Sequence:          7,
			TimeoutTimestamp:  uint64(time.Now().Add(-time.Minute).Unix()),
		}

		tests := []struct {
			name             string
			commitment       bool
			commitmentErr    error
			timeoutSucceeded bool
			wantRemaining    bool
		}{
			{
				name:          "commitment check error keeps pending packet",
				commitmentErr: errors.New("rpc unavailable"),
				wantRemaining: true,
			},
			{
				name:             "timeout submit failure keeps pending packet",
				commitment:       true,
				timeoutSucceeded: false,
				wantRemaining:    true,
			},
			{
				name:          "cleared commitment removes pending packet",
				commitment:    false,
				wantRemaining: false,
			},
			{
				name:             "successful timeout removes pending packet",
				commitment:       true,
				timeoutSucceeded: true,
				wantRemaining:    false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				tracker := NewPendingPacketTracker()
				tracker.Add(expiredPacket, 123)

				timeoutCalled := false
				svc := New(nil, nil, DefaultConfig())
				svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
					tag:     "Test",
					tracker: tracker,
					prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
						return &client.LightBlock{}, true
					},
					hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
						return tc.commitment, tc.commitmentErr
					},
					timeoutSend: func(_ context.Context, _ evmTimeoutDeps, packet EthPacket, _ *client.LightBlock) timeoutSendOutcome {
						timeoutCalled = true
						if packet.BlockNumber != 123 {
							t.Fatalf("timeout packet block number = %d, want 123", packet.BlockNumber)
						}
						if tc.timeoutSucceeded {
							return timeoutSent
						}
						return timeoutDeferred
					},
				})

				gotRemaining := tracker.Len() == 1
				if gotRemaining != tc.wantRemaining {
					t.Fatalf("tracker remaining = %v, want %v", gotRemaining, tc.wantRemaining)
				}
				wantTimeoutCalled := tc.commitmentErr == nil && tc.commitment
				if timeoutCalled != wantTimeoutCalled {
					t.Fatalf("timeoutCalled = %v, want %v", timeoutCalled, wantTimeoutCalled)
				}
			})
		}
	})

	t.Run("shared update failure runs once and defers all", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		for _, sequence := range []uint64{1, 2, 3} {
			tracker.Add(channeltypesv2.Packet{
				SourceClient: "l2-client-0", Sequence: sequence,
				TimeoutTimestamp: uint64(time.Now().Add(-time.Minute).Unix()),
			}, sequence)
		}

		prepareCalls := 0
		timeoutCalls := 0
		svc := New(nil, nil, DefaultConfig())
		svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
			tag:     "Test",
			tracker: tracker,
			prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
				prepareCalls++
				return nil, false
			},
			hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
				return true, nil
			},
			timeoutSend: func(context.Context, evmTimeoutDeps, EthPacket, *client.LightBlock) timeoutSendOutcome {
				timeoutCalls++
				return timeoutSent
			},
		})

		if prepareCalls != 1 {
			t.Fatalf("shared client update attempts = %d, want 1", prepareCalls)
		}
		if timeoutCalls != 0 {
			t.Fatalf("timeout submissions after failed shared update = %d, want 0", timeoutCalls)
		}
		for _, info := range tracker.GetAll() {
			if info.Deferrals != 1 || info.TimeoutAttempts != 0 {
				t.Fatalf("seq=%d retry state = deferrals:%d attempts:%d, want 1:0",
					info.Packet.Sequence, info.Deferrals, info.TimeoutAttempts)
			}
		}
	})

	t.Run("clears shared deferrals after recovery", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		packet := channeltypesv2.Packet{
			SourceClient: "l2-client-0", Sequence: 1,
			TimeoutTimestamp: uint64(time.Now().Add(-time.Minute).Unix()),
		}
		tracker.Add(packet, 123)
		// Make the packet due again with a saturated shared-outage history.
		past := time.Now().Add(-time.Hour)
		for i := 0; i < 6; i++ {
			tracker.deferTimeoutRetryForTest(packet.SourceClient, packet.Sequence, past)
		}

		svc := New(nil, nil, DefaultConfig())
		svc.scanForEVMTimeouts(context.Background(), evmTimeoutDeps{}, evmTimeoutScanOptions{
			tag:     "Test",
			tracker: tracker,
			prepareTimeouts: func(context.Context, evmTimeoutDeps) (*client.LightBlock, bool) {
				return &client.LightBlock{}, true
			},
			hasPendingCommitment: func(context.Context, evmTimeoutDeps, channeltypesv2.Packet) (bool, error) {
				return true, nil
			},
			timeoutSend: func(context.Context, evmTimeoutDeps, EthPacket, *client.LightBlock) timeoutSendOutcome {
				return timeoutDeferred
			},
		})

		got := tracker.GetAll()
		if len(got) != 1 || got[0].Deferrals != 1 {
			t.Fatalf("deferrals after recovered shared prerequisite = %#v, want a fresh single deferral", got)
		}
	})
}

func TestHasEthIBCPathValue_HonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, err := HasEthIBCPathValue(ctx, EVMEndpoint{Client: client, Contracts: EVMContracts{Router: common.HexToAddress("0x1")}}, []byte("path"))
		done <- err
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("error = %v, want context.Canceled", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("storage query ignored cancellation")
	}
}

func TestScanForCosmosTimeouts_HonorsCancellation(t *testing.T) {
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { <-release }))
	t.Cleanup(server.Close)
	t.Cleanup(func() { close(release) })
	client, err := ethclient.Dial(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)

	svc := New(nil, nil, DefaultConfig())
	svc.TrackCosmosPending(channeltypesv2.Packet{
		SourceClient:      "07-tendermint-0",
		DestinationClient: "client-0",
		Sequence:          1,
		TimeoutTimestamp:  1,
	}, 100)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		svc.scanForCosmosTimeouts(ctx, CosmosEndpoint{}, EVMEndpoint{
			Client:    client,
			Contracts: EVMContracts{Router: common.HexToAddress("0x1")},
		}, "08-wasm-0")
		close(done)
	}()
	time.Sleep(30 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Cosmos timeout scan ignored cancellation")
	}
}

// TestUpdateMsgsFailed is the second P1. SendCosmosTxBatch submits
// [updateMsgs..., timeoutMsgs...] and SucceededCount is a prefix length, so
// fewer successes than update messages proves no timeout was ever executed.
func TestUpdateMsgsFailed(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		updateMsgCount int
		want           bool
	}{
		{
			name:           "update prefix incomplete: no timeout could have run",
			err:            &BatchPartialError{SucceededCount: 1, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           true,
		},
		{
			name:           "update prefix complete: the failure is in the timeouts",
			err:            &BatchPartialError{SucceededCount: 2, Err: errors.New("boom")},
			updateMsgCount: 2,
			want:           false,
		},
		{
			name:           "whole tx failed with update msgs present: nothing executed",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 1,
			want:           true,
		},
		{
			name:           "no update msgs: the timeouts are the only suspects",
			err:            errors.New("broadcast failed"),
			updateMsgCount: 0,
			want:           false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := updateMsgsFailed(tc.err, tc.updateMsgCount); got != tc.want {
				t.Fatalf("updateMsgsFailed = %v, want %v", got, tc.want)
			}
		})
	}
}

func timeoutRetryPacket(sequence uint64) channeltypesv2.Packet {
	return channeltypesv2.Packet{SourceClient: "client-0", Sequence: sequence, TimeoutTimestamp: 1}
}

type timeoutBatchHandler struct {
	TransactionHandler
	singletons [][]any
}

func (h *timeoutBatchHandler) SendCosmosTxBatch(_ context.Context, _ CosmosEndpoint, msgs []any) error {
	h.singletons = append(h.singletons, msgs)
	if len(msgs) != 1 {
		return errors.New("test only accepts singleton probes")
	}
	if msgs[0] == "poison" {
		return ErrPermanentRelayFailure
	}
	return nil
}

func TestHandleCosmosTimeoutBatchFailure(t *testing.T) {
	t.Run("permanent cosmos timeout batch isolates the poison packet", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		poison := timeoutRetryPacket(11)
		healthy := timeoutRetryPacket(12)
		tracker.Add(poison, 10)
		tracker.Add(healthy, 10)

		handler := &timeoutBatchHandler{}
		svc := New(handler, nil, DefaultConfig())
		processed := []pendingPacketInfo{{Packet: poison}, {Packet: healthy}}
		svc.handleCosmosTimeoutBatchFailure(context.Background(), CosmosEndpoint{}, tracker, nil, []any{"poison", "healthy"}, processed, ErrPermanentRelayFailure)

		if len(handler.singletons) != 2 {
			t.Fatalf("singleton probes = %d, want 2", len(handler.singletons))
		}
		if tracker.Len() != 1 {
			t.Fatalf("tracked packets = %d, want only poison packet", tracker.Len())
		}
		remaining := tracker.GetAll()
		if len(remaining) != 1 || remaining[0].Packet.Sequence != poison.Sequence || remaining[0].TimeoutAttempts != 1 {
			t.Fatalf("remaining tracker state = %#v, want one charged poison packet", remaining)
		}
	})

	t.Run("update prefix failure defers all", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		first := timeoutRetryPacket(31)
		second := timeoutRetryPacket(32)
		tracker.Add(first, 10)
		tracker.Add(second, 10)

		svc := New(nil, nil, DefaultConfig())
		svc.handleCosmosTimeoutBatchFailure(
			context.Background(), CosmosEndpoint{}, tracker,
			[]any{"update"}, []any{"first", "second"},
			[]pendingPacketInfo{{Packet: first}, {Packet: second}},
			&BatchPartialError{SucceededCount: 0, Err: ErrPermanentRelayFailure},
		)

		for _, info := range tracker.GetAll() {
			if info.TimeoutAttempts != 0 || info.Deferrals != 1 {
				t.Fatalf("seq=%d retry state = attempts:%d deferrals:%d, want 0:1", info.Packet.Sequence, info.TimeoutAttempts, info.Deferrals)
			}
		}
	})

	t.Run("singleton probe defers transient failure", func(t *testing.T) {
		tracker := NewPendingPacketTracker()
		packet := timeoutRetryPacket(41)
		tracker.Add(packet, 10)

		svc := New(&transientTimeoutBatchHandler{}, nil, DefaultConfig())
		svc.handleCosmosTimeoutBatchFailure(
			context.Background(), CosmosEndpoint{}, tracker,
			nil, []any{"timeout"}, []pendingPacketInfo{{Packet: packet}}, ErrPermanentRelayFailure,
		)

		info := tracker.GetAll()
		if len(info) != 1 || info[0].TimeoutAttempts != 0 || info[0].Deferrals != 1 {
			t.Fatalf("singleton transient retry state = %#v, want one deferred uncharged packet", info)
		}
	})
}

type transientTimeoutBatchHandler struct{ TransactionHandler }

func (h *transientTimeoutBatchHandler) SendCosmosTxBatch(_ context.Context, _ CosmosEndpoint, _ []any) error {
	return errors.New("cosmos rpc unavailable")
}

func TestResolveCosmosTimeoutProofState(t *testing.T) {
	t.Run("defers shared failures", func(t *testing.T) {
		for _, tc := range []struct {
			name             string
			buildUpdate      func() (*EthClientUpdateResult, error)
			fetchClientState func() (*client.EthereumClientState, error)
		}{
			{
				name:        "update build",
				buildUpdate: func() (*EthClientUpdateResult, error) { return nil, errors.New("beacon unavailable") },
				fetchClientState: func() (*client.EthereumClientState, error) {
					t.Fatal("client state fetch must not run")
					return nil, nil
				},
			},
			{
				name:             "client state fetch",
				buildUpdate:      func() (*EthClientUpdateResult, error) { return &EthClientUpdateResult{}, nil },
				fetchClientState: func() (*client.EthereumClientState, error) { return nil, errors.New("cosmos unavailable") },
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				tracker := NewPendingPacketTracker()
				expired := make([]pendingPacketInfo, 0, 2)
				for _, sequence := range []uint64{1, 2} {
					packet := timeoutRetryPacket(sequence)
					tracker.Add(packet, sequence)
					expired = append(expired, pendingPacketInfo{Packet: packet, BlockNumber: sequence})
				}

				svc := New(nil, nil, DefaultConfig())
				_, _, ok := svc.resolveCosmosTimeoutProofState(tracker, expired, tc.buildUpdate, tc.fetchClientState)
				if ok {
					t.Fatal("failed shared prerequisite reported success")
				}
				for _, info := range tracker.GetAll() {
					if info.Deferrals != 1 || info.TimeoutAttempts != 0 {
						t.Fatalf("seq=%d retry state = deferrals:%d attempts:%d, want 1:0",
							info.Packet.Sequence, info.Deferrals, info.TimeoutAttempts)
					}
				}
			})
		}
	})

	// TestResolveCosmosTimeoutProofState/prerequisite failures defer without charging is the P1 from
	// review, raised twice: both shared-prerequisite exits returned bare.
	//
	// Neither failure is attributable to any packet -- they fail identically for the
	// whole scan -- so nothing may be charged. But NOT deferring is what made the
	// bug expensive and invisible at once: the next 30s scan rebuilds the client
	// state, finality and update construction from scratch forever, and because
	// Deferrals stays at zero the STUCK line and the stuck gauge cannot report it.
	t.Run("prerequisite failures defer without charging", func(t *testing.T) {
		t.Parallel()

		for _, tc := range []struct {
			name             string
			buildUpdate      func() (*EthClientUpdateResult, error)
			fetchClientState func() (*client.EthereumClientState, error)
		}{
			{
				name:        "BuildEthClientUpdateMsgs fails",
				buildUpdate: func() (*EthClientUpdateResult, error) { return nil, fmt.Errorf("beacon api unavailable") },
				fetchClientState: func() (*client.EthereumClientState, error) {
					t.Fatal("client state must not be fetched after the update build failed")
					return nil, nil
				},
			},
			{
				// updateResult with a nil EthClientState sends the scan to the
				// fallback fetch -- the second exit, which survived the first fix.
				name:             "GetEthereumClientState fallback fails",
				buildUpdate:      func() (*EthClientUpdateResult, error) { return &EthClientUpdateResult{}, nil },
				fetchClientState: func() (*client.EthereumClientState, error) { return nil, fmt.Errorf("cosmos rpc down") },
			},
		} {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				s, expired := prereqServices(t, 11, 12)
				tracker := s.BatchBuilder.PendingTracker

				_, _, ok := s.resolveCosmosTimeoutProofState(tracker, expired, tc.buildUpdate, tc.fetchClientState)
				if ok {
					t.Fatal("a failed shared prerequisite must abort the scan")
				}

				// Deferred: nothing may be due again immediately, or the next 30s
				// scan rebuilds the client state, finality and update construction.
				if due := tracker.GetDue(time.Now()); len(due) != 0 {
					t.Fatalf("%d packet(s) still due immediately; the next scan would rebuild everything", len(due))
				}

				// Not charged: a shared failure is not evidence about any packet.
				if got := tracker.DeadLetteredTimeouts(); got != 0 {
					t.Fatalf("DeadLetteredTimeouts() = %d; a shared prerequisite failure must not charge packets", got)
				}
				if tracker.Len() != len(expired) {
					t.Fatalf("tracker holds %d packet(s), want %d; nothing was delivered", tracker.Len(), len(expired))
				}
			})
		}
	})

	// TestResolveCosmosTimeoutProofState/repeated failures become visible: deferrals never
	// dead-letter, so the stuck gauge is the only thing that can surface a
	// prerequisite that never recovers. Without the deferral it reads zero forever.
	t.Run("repeated failures become visible", func(t *testing.T) {
		t.Parallel()

		s, expired := prereqServices(t, 31)
		failing := func() (*EthClientUpdateResult, error) { return nil, fmt.Errorf("beacon api unavailable") }
		unreachable := func() (*client.EthereumClientState, error) {
			t.Fatal("unreachable")
			return nil, nil
		}

		for range deferralAlertThreshold {
			s.resolveCosmosTimeoutProofState(s.BatchBuilder.PendingTracker, expired, failing, unreachable)
		}

		count, worst := s.BatchBuilder.PendingTracker.StuckTimeouts()
		if count != 1 || worst != deferralAlertThreshold {
			t.Fatalf("StuckTimeouts() = (%d, %d), want (1, %d): a prerequisite that never recovers must become visible",
				count, worst, deferralAlertThreshold)
		}
	})

	// A healthy resolve must not touch the tracker at all.
	t.Run("success leaves packets due", func(t *testing.T) {
		t.Parallel()

		s, expired := prereqServices(t, 41)
		_, state, ok := s.resolveCosmosTimeoutProofState(s.BatchBuilder.PendingTracker, expired, okUpdate, func() (*client.EthereumClientState, error) {
			t.Fatal("fallback must not run when the update carries a client state")
			return nil, nil
		})
		if !ok || state == nil {
			t.Fatal("a successful resolve must report ok with a client state")
		}
		if due := s.BatchBuilder.PendingTracker.GetDue(time.Now()); len(due) != len(expired) {
			t.Fatalf("GetDue() = %d, want %d: a successful resolve must not defer anything", len(due), len(expired))
		}
	})
}

// prereqServices builds a Services holding n tracked packets, and returns the
// pendingPacketInfo set the scan would be carrying at the point the shared
// prerequisites are resolved.
func prereqServices(t *testing.T, seqs ...uint64) (*Services, []pendingPacketInfo) {
	t.Helper()
	s := &Services{BatchBuilder: NewBatchBuilder()}
	expired := make([]pendingPacketInfo, 0, len(seqs))
	for _, seq := range seqs {
		pkt := channeltypesv2.Packet{
			Sequence:          seq,
			SourceClient:      "08-wasm-0",
			DestinationClient: "client-0",
			TimeoutTimestamp:  1,
		}
		s.BatchBuilder.PendingTracker.Add(pkt, 100)
		expired = append(expired, pendingPacketInfo{Packet: pkt, BlockNumber: 100})
	}
	return s, expired
}

func okUpdate() (*EthClientUpdateResult, error) {
	return &EthClientUpdateResult{EthClientState: &client.EthereumClientState{LatestSlot: 1}}, nil
}

func TestRecoveryState(t *testing.T) {
	t.Run("missing file and partial checkpoints", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		path := filepath.Join(t.TempDir(), "state", "cursors.json")
		store, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := store.Get(); ok {
			t.Fatal("missing file unexpectedly contained a source")
		}
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20, "eth_send": 30}); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 10, "eth_write_ack": 40}); err != nil {
			t.Fatal(err)
		}
		got, ok := store.Get()
		if !ok || !got.equal(RecoveryCursors{"cosmos": 10, "eth_send": 30, "eth_write_ack": 40}) {
			t.Fatalf("merged cursors = %+v, ok=%t", got, ok)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if info.Mode().Perm() != stateFilePerm {
			t.Fatalf("state mode = %o, want %o", info.Mode().Perm(), stateFilePerm)
		}

		resetLoadedRecoveryStoresForTest()
		reloaded, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		got, ok = reloaded.Get()
		if !ok || got["cosmos"] != 10 || got["eth_send"] != 30 || got["eth_write_ack"] != 40 {
			t.Fatalf("reloaded cursors = %+v, ok=%t", got, ok)
		}
	})

	t.Run("checkpoint can move backward", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		path := filepath.Join(t.TempDir(), "state.json")
		store, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 100, "eth_send": 200}); err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 50}); err != nil {
			t.Fatal(err)
		}

		got, ok := store.Get()
		want := RecoveryCursors{"cosmos": 50, "eth_send": 200}
		if !ok || !got.equal(want) {
			t.Fatalf("checkpoint cursors = %+v, ok=%t, want %+v", got, ok, want)
		}

		resetLoadedRecoveryStoresForTest()
		reloaded, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		got, ok = reloaded.Get()
		if !ok || !got.equal(want) {
			t.Fatalf("reloaded checkpoint cursors = %+v, ok=%t, want %+v", got, ok, want)
		}
	})

	t.Run("failed write does not advance memory", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		dir := t.TempDir()
		path := filepath.Join(dir, "state.json")
		store, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 10}); err != nil {
			t.Fatal(err)
		}

		blockedPath := filepath.Join(dir, "blocked")
		if err := os.Mkdir(blockedPath, 0o700); err != nil {
			t.Fatal(err)
		}
		store.path = blockedPath
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20}); err == nil {
			t.Fatal("Save succeeded when the target path was a directory")
		}
		got, ok := store.Get()
		if !ok || got["cosmos"] != 10 {
			t.Fatalf("cursor after failed write = %+v, ok=%t, want height 10", got, ok)
		}

		store.path = path
		if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": 20}); err != nil {
			t.Fatal(err)
		}
		resetLoadedRecoveryStoresForTest()
		reloaded, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		got, ok = reloaded.Get()
		if !ok || got["cosmos"] != 20 {
			t.Fatalf("reloaded cursor after retry = %+v, ok=%t, want height 20", got, ok)
		}
	})

	t.Run("rejects corrupt and unsupported files", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			data []byte
		}{
			{name: "corrupt", data: []byte("{")},
			{name: "wrong version", data: []byte(`{"version":1,"sources":{}}`)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				resetLoadedRecoveryStoresForTest()
				path := filepath.Join(t.TempDir(), "state.json")
				if err := os.WriteFile(path, tc.data, 0o600); err != nil {
					t.Fatal(err)
				}
				if _, err := LoadRecoveryState(path); err == nil {
					t.Fatal("LoadRecoveryState succeeded for invalid state")
				}
			})
		}
	})

	// One process now keeps one record, but two recovery streams write to it: the
	// Cosmos subscriber owns CosmosHeight, the ETH subscriber owns EthSendBlock and
	// EthWriteAckBlk. They run concurrently and each save is a read-modify-write of
	// the whole file, so a lost update here means a stream silently resumes from an
	// older cursor and rescans -- or worse, the other stream's newer cursor is
	// rolled back and its window is skipped.
	//
	// Memoization by path is what makes the two share one view; without it each
	// would hold its own copy and the last writer would erase the other's field.
	t.Run("memoization shares one view between streams", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		path := filepath.Join(t.TempDir(), "state.json")
		cosmosView, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		ethView, err := LoadRecoveryState(path)
		if err != nil {
			t.Fatal(err)
		}
		if cosmosView != ethView {
			t.Fatal("the same path returned two stores; the two streams would each keep their own cursors")
		}

		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			if err := cosmosView.SaveCheckpoint(RecoveryCursors{"cosmos": 11}); err != nil {
				t.Errorf("cosmos save: %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			if err := ethView.SaveCheckpoint(RecoveryCursors{"eth_send": 22, "eth_write_ack": 33}); err != nil {
				t.Errorf("eth save: %v", err)
			}
		}()
		wg.Wait()

		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		var file recoveryStateFile
		if err := json.Unmarshal(data, &file); err != nil {
			t.Fatal(err)
		}
		want := RecoveryCursors{"cosmos": 11, "eth_send": 22, "eth_write_ack": 33}
		if !file.Cursors.equal(want) {
			t.Fatalf("persisted cursors = %+v, want %+v; one stream's write erased the other's", file.Cursors, want)
		}
	})

	t.Run("write leaves no temp files", func(t *testing.T) {
		resetLoadedRecoveryStoresForTest()
		dir := t.TempDir()
		store, err := LoadRecoveryState(filepath.Join(dir, "state.json"))
		if err != nil {
			t.Fatal(err)
		}
		for i := range 5 {
			if err := store.SaveCheckpoint(RecoveryCursors{"cosmos": uint64(i + 1)}); err != nil {
				t.Fatalf("save %d: %v", i, err)
			}
		}
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".tmp") {
				t.Fatalf("temporary recovery file %q was not removed", entry.Name())
			}
		}
	})
}
