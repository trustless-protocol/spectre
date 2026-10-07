package services

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	spectreContract "relayer/bindings/SpectreClient"
	updateClientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"
	"relayer/rpcmock"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var pinnedSnapshotTestClient = common.HexToAddress("0x00000000000000000000000000000000000000dd")

// pinnedSnapshotNode answers getClientState and getPinnedSnapshot the way a
// SpectreClient does after a few updateApplicationState calls and no rotation:
// latestHeight has moved past the snapshot height.
func pinnedSnapshotNode(t *testing.T, snapshotHeight uint64, snapshotTime time.Time, trustingPeriod uint32) (*rpcmock.EVMNode, EVMEndpoint) {
	t.Helper()
	spectreABI, err := spectreContract.ContractSpectreClientMetaData.GetAbi()
	if err != nil {
		t.Fatal(err)
	}
	node := rpcmock.NewEVMNode(t)

	state, err := relayerclient.EncodeClientState(relayerclient.ClientState{
		ChainId:        "test-1",
		TrustLevel:     relayerclient.TrustThreshold{Numerator: 2, Denominator: 3},
		LatestHeight:   updateClientContract.IICS02ClientMsgsHeight{RevisionNumber: 1, RevisionHeight: snapshotHeight + 500},
		TrustingPeriod: trustingPeriod,
	})
	if err != nil {
		t.Fatal(err)
	}
	stateOut, err := spectreABI.Methods["getClientState"].Outputs.Pack(state)
	if err != nil {
		t.Fatal(err)
	}
	node.SetCallResultFor(hexutil.Encode(spectreABI.Methods["getClientState"].ID), hexutil.Encode(stateOut))

	snapOut, err := spectreABI.Methods["getPinnedSnapshot"].Outputs.Pack(snapshotHeight, big.NewInt(snapshotTime.UnixNano()))
	if err != nil {
		t.Fatal(err)
	}
	node.SetCallResultFor(hexutil.Encode(spectreABI.Methods["getPinnedSnapshot"].ID), hexutil.Encode(snapOut))

	return node, EVMEndpoint{Client: node.Client(t), Contracts: EVMContracts{SpectreClient: pinnedSnapshotTestClient}}
}

func TestFetchPinnedSnapshotTime(t *testing.T) {
	ctx := context.Background()
	snapshotTime := time.Date(2026, 8, 20, 12, 43, 56, 0, time.UTC)

	t.Run("reads the snapshot timestamp", func(t *testing.T) {
		_, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, 57600)
		got, supported, err := fetchPinnedSnapshotTime(ctx, evm)
		if err != nil || !supported || !got.Equal(snapshotTime) {
			t.Fatalf("got %v supported=%v err=%v; want %v supported=true", got, supported, err, snapshotTime)
		}
	})

	// A client deployed before the getter existed reverts on the unknown
	// selector. That is "fall back", not an error the caller has to retry.
	t.Run("revert means a pre-#541 client", func(t *testing.T) {
		node, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, 57600)
		node.FailCall(errors.New("execution reverted"))
		_, supported, err := fetchPinnedSnapshotTime(ctx, evm)
		if err != nil || supported {
			t.Fatalf("supported=%v err=%v; want unsupported with no error", supported, err)
		}
	})

	// A revert that carries data arrives as code 3 with hex data: still a revert.
	t.Run("code-3 revert with data means a pre-#541 client", func(t *testing.T) {
		node, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, 57600)
		node.FailCallWith(3, "execution reverted", "0x08c379a0")
		_, supported, err := fetchPinnedSnapshotTime(ctx, evm)
		if err != nil || supported {
			t.Fatalf("supported=%v err=%v; want unsupported with no error", supported, err)
		}
	})

	// Review of d85c8cd: a throttled endpoint whose error carries diagnostic
	// data was read as a revert, and the expiry silently fell back to the
	// latest-height anchor. Data alone is not a revert signal.
	t.Run("rate limit with data is an error, not a fallback", func(t *testing.T) {
		node, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, 57600)
		node.FailCallWith(-32005, "limit exceeded", map[string]any{"retryAfter": 1})
		_, supported, err := fetchPinnedSnapshotTime(ctx, evm)
		if err == nil || supported {
			t.Fatalf("supported=%v err=%v; want an error so the caller retries", supported, err)
		}
	})

	// A node that never ran the call must not be read as an old contract:
	// downgrading to the latest-height anchor on a blip is the #541 bug.
	t.Run("transport failure is an error", func(t *testing.T) {
		node, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, 57600)
		node.FailCall(errors.New("header not found"))
		if _, _, err := fetchPinnedSnapshotTime(ctx, evm); err == nil {
			t.Fatal("want an error for a non-revert RPC failure")
		}
	})
}

// The #541 shape: latestHeight advanced by app-state updates, pinned set never
// rotated. Expiry must follow the snapshot, and must not need the Cosmos RPC
// (the zero CosmosEndpoint would panic if it were touched).
func TestCosmosClientExpiry_AnchorsOnPinnedSnapshot(t *testing.T) {
	snapshotTime := time.Date(2026, 8, 20, 12, 43, 56, 0, time.UTC)
	const trustingPeriod = 57600 // 16h
	_, evm := pinnedSnapshotNode(t, 6593231, snapshotTime, trustingPeriod)

	expiresAt, period, err := CosmosClientExpiry(context.Background(), CosmosEndpoint{}, evm)
	if err != nil {
		t.Fatal(err)
	}
	if want := snapshotTime.Add(16 * time.Hour); !expiresAt.Equal(want) {
		t.Fatalf("expiresAt = %v, want %v (snapshot + trusting period)", expiresAt, want)
	}
	if period != 16*time.Hour {
		t.Fatalf("period = %v, want 16h", period)
	}
}

func TestSnapshotTimestampToTime(t *testing.T) {
	want := time.Date(2026, 8, 20, 12, 43, 56, 0, time.UTC)
	if got, err := snapshotTimestampToTime(big.NewInt(want.UnixNano())); err != nil || !got.Equal(want) {
		t.Fatalf("got %v, %v; want %v", got, err, want)
	}
	for name, nanos := range map[string]*big.Int{
		"nil":      nil,
		"zero":     big.NewInt(0),
		"negative": big.NewInt(-1),
		"overflow": new(big.Int).Lsh(big.NewInt(1), 64),
	} {
		if _, err := snapshotTimestampToTime(nanos); err == nil {
			t.Fatalf("%s: want an error", name)
		}
	}
}

func TestRotationDueIn(t *testing.T) {
	now := time.Date(2026, 8, 21, 0, 0, 0, 0, time.UTC)
	anchor := now.Add(-10 * time.Hour)
	if got := rotationDueIn(anchor, 15*time.Hour, now); got != 5*time.Hour {
		t.Fatalf("due = %v, want 5h", got)
	}
	if got := rotationDueIn(anchor, 10*time.Hour, now); got != 0 {
		t.Fatalf("exactly due = %v, want 0", got)
	}
	if got := rotationDueIn(anchor, time.Hour, now); got != 0 {
		t.Fatalf("overdue = %v, want 0 (rotate now)", got)
	}
}
