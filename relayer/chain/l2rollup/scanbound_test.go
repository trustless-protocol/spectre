package l2rollup

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"relayer/chain"

	attestorpb "attestor/types/attestor"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// #575 L14: the scanner cut [cursor, head] into spans with `start += chunk`, where
// head is whatever the endpoint reports.

// Near the top of the uint64 range that addition wrapped to a small number, the
// loop condition stayed true, and the loop never ended. The watchdog turns "does
// not return" into a failure instead of a hung run.
func TestChunkSpans_TopOfTheRangeTerminates(t *testing.T) {
	const top = uint64(math.MaxUint64)
	cases := []struct {
		name            string
		from, to, chunk uint64
		want            [][2]uint64
	}{
		{"ragged tail ending at the maximum", top - 25, top, 10, [][2]uint64{{top - 25, top - 16}, {top - 15, top - 6}, {top - 5, top}}},
		{"exact multiple ending at the maximum", top - 19, top, 10, [][2]uint64{{top - 19, top - 10}, {top - 9, top}}},
		{"chunk of one at the maximum", top - 1, top, 1, [][2]uint64{{top - 1, top - 1}, {top, top}}},
		{"the whole range in the widest chunk", 0, top, top, [][2]uint64{{0, top - 1}, {top, top}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			done := make(chan [][2]uint64, 1)
			go func() { done <- chunkSpans(tc.from, tc.to, tc.chunk) }()
			select {
			case got := <-done:
				if fmt.Sprint(got) != fmt.Sprint(tc.want) {
					t.Fatalf("spans = %v, want %v", got, tc.want)
				}
			case <-time.After(200 * time.Millisecond):
				t.Fatal("chunkSpans did not return: the span arithmetic wrapped and the loop cannot end")
			}
		})
	}
}

func TestScanWindowEnd(t *testing.T) {
	cases := []struct {
		name              string
		from, head, chunk uint64
		want              uint64
	}{
		{"head inside the window", 1000, 1500, 0, 1500},
		{"head exactly one window away", 1000, 1000 + maxScanBlocksPerTick - 1, 0, 1000 + maxScanBlocksPerTick - 1},
		{"head past the window, unchunked", 1000, 5_000_000, 0, 1000 + maxScanBlocksPerTick - 1},
		{"a small chunk caps the spans, not the blocks", 1000, 5_000_000, 10, 1000 + 10*maxLogSpansPerTick - 1},
		{"a wide chunk caps the blocks", 1000, 5_000_000, 50_000, 1000 + maxScanBlocksPerTick - 1},
		{"an implausible head", 1000, 1 << 62, 10, 1000 + 10*maxLogSpansPerTick - 1},
		{"the maximum head", math.MaxUint64 - 5, math.MaxUint64, 10, math.MaxUint64},
		{"the maximum head, far cursor", 1000, math.MaxUint64, 0, 1000 + maxScanBlocksPerTick - 1},
		{"head behind the cursor", 1000, 900, 10, 900},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := scanWindowEnd(tc.from, tc.head, tc.chunk)
			if got != tc.want {
				t.Fatalf("scanWindowEnd(%d, %d, %d) = %d, want %d", tc.from, tc.head, tc.chunk, got, tc.want)
			}
			// What the bound is for: the spans one tick is cut into.
			if spans := uint64(len(chunkSpans(tc.from, got, tc.chunk))); spans > maxLogSpansPerTick {
				t.Fatalf("one tick is cut into %d spans, over the %d cap", spans, maxLogSpansPerTick)
			}
		})
	}
}

// One implausible head must not size the scan. With a chunk configured, the old
// loop asked for one span per chunk all the way to the reported head, allocated
// up front.
func TestSubscribe_AnImplausibleHeadDoesNotSizeTheScan(t *testing.T) {
	restore := l2SubscribeInterval
	l2SubscribeInterval = 20 * time.Millisecond
	t.Cleanup(func() { l2SubscribeInterval = restore })

	const (
		startHead = uint64(1000)
		chunk     = uint64(50_000)
	)
	router := ethcommon.HexToAddress("0x2222222222222222222222222222222222222222")
	node := &idleL2Node{router: router, logBlock: startHead}
	node.setHead(startHead)
	src := &Source{
		eth:                node.start(t),
		headKind:           Unsafe,
		l2ClientID:         idleTestClientID,
		cosmosWasmClientID: idleTestDestID,
		router:             router,
		logScanChunk:       chunk,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		done <- src.Subscribe(ctx, func(context.Context, []chain.Event) []int { return nil })
	}()

	waitForScan(t, node)
	seeded := len(node.scannedRanges())
	node.setHead(1 << 62)

	// Let a few bounded ticks run against the implausible head.
	deadline := time.Now().Add(3 * time.Second)
	for len(node.scannedRanges()) < seeded+4 && time.Now().Before(deadline) {
		time.Sleep(5 * time.Millisecond)
	}
	cancel()
	<-done

	ranges := node.scannedRanges()
	if len(ranges) < seeded+4 {
		t.Fatalf("scanned ranges %v: the scan stopped making progress after the head jumped", ranges)
	}
	var furthest uint64
	for _, r := range ranges {
		if r[1] < r[0] || r[1]-r[0]+1 > chunk {
			t.Fatalf("scanned range %v is not one chunk of at most %d blocks", r, chunk)
		}
		if r[1] > furthest {
			furthest = r[1]
		}
	}
	// SendPacket and WriteAcknowledgement are read over the same span, so every
	// distinct span appears twice; the cursor must still have advanced by whole
	// windows only.
	ticks := uint64(len(ranges)-seeded)/2/(maxScanBlocksPerTick/chunk) + 1
	if limit := startHead + ticks*maxScanBlocksPerTick; furthest > limit {
		t.Fatalf("the cursor reached block %d after about %d tick(s); a bounded scan stays under %d", furthest, ticks, limit)
	}
}

// A block number wider than 64 bits is not a head. Uint64 on it returns the low
// 64 bits, which would be scanned to as if they were real.
func TestSourceHead_RejectsANumberThatIsNotAUint64(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 2^64 + 100: its low 64 bits read as block 100.
		header := strings.Replace(idleHeaderJSON(0), `"number":"0x0"`, `"number":"0x10000000000000064"`, 1)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"jsonrpc":"2.0","id":1,"result":%s}`, header)
	}))
	t.Cleanup(srv.Close)
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(client.Close)

	head, err := (&Source{eth: client, headKind: Unsafe}).head(context.Background())
	if err == nil {
		t.Fatalf("head = %d, want the oversized block number refused", head)
	}
}

// panicClient is an attestor client that panics on an answer it cannot decode.
type panicClient struct{ *fakeAttestor }

func (panicClient) AttestedUpTo(context.Context, string, bool) (*attestorpb.AttestedRoot, bool, error) {
	var m map[string]int
	m["x"]++
	return nil, false, nil
}

// A panic in one endpoint's client must count as that endpoint failing. The
// quorum tolerates failed endpoints; a crashed process does not.
func TestAttestedUpTo_APanickingEndpointIsAFailedEndpoint(t *testing.T) {
	good := func(index uint16, height uint64) SigningAttestor {
		return SigningAttestor{Index: index, Client: &fakeAttestor{root: &attestorpb.AttestedRoot{L2BlockNumber: height}, found: true}}
	}
	frontier, err := NewQuorumAttestationFrontier([]SigningAttestor{
		good(0, 120), good(1, 110), {Index: 2, Client: panicClient{&fakeAttestor{}}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	root, found, err := frontier.AttestedUpTo(context.Background(), "op-sepolia", false)
	if err != nil || !found || root.GetL2BlockNumber() == 0 {
		t.Fatalf("root=%v found=%v err=%v, want the two healthy endpoints to make the quorum", root, found, err)
	}

	// With only the panicking endpoint left to reach the threshold, it is a
	// reported failure, not a crash.
	frontier, err = NewQuorumAttestationFrontier([]SigningAttestor{
		good(0, 120), {Index: 2, Client: panicClient{&fakeAttestor{}}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := frontier.AttestedUpTo(context.Background(), "op-sepolia", false); err == nil || !strings.Contains(err.Error(), "panicked") {
		t.Fatalf("err = %v, want the endpoint's panic reported as its failure", err)
	}
}

// #575 L2: the frontier was the minimum of the FIRST threshold answers, so one
// fast endpoint reporting a low height pinned RelayableHeight.
func TestAttestedUpTo_AFastLowEndpointDoesNotPinTheFrontier(t *testing.T) {
	restore := frontierQuorumGrace
	frontierQuorumGrace = 2 * time.Second
	t.Cleanup(func() { frontierQuorumGrace = restore })

	release := make(chan struct{})
	root := func(height uint64) *attestorpb.AttestedRoot { return &attestorpb.AttestedRoot{L2BlockNumber: height} }
	frontier, err := NewQuorumAttestationFrontier([]SigningAttestor{
		{Index: 0, Client: &fakeAttestor{root: root(10), found: true}}, // answers at once, far behind
		{Index: 1, Client: &fakeAttestor{root: root(120), found: true, frontierBlock: release}},
		{Index: 2, Client: &fakeAttestor{root: root(110), found: true, frontierBlock: release}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		time.Sleep(50 * time.Millisecond) // the low answer is in before the other two
		close(release)
	}()

	start := time.Now()
	got, found, err := frontier.AttestedUpTo(context.Background(), "op-sepolia", false)
	if err != nil || !found {
		t.Fatalf("found=%v err=%v", found, err)
	}
	// Two endpoints attest through 110: that is the quorum frontier.
	if got.GetL2BlockNumber() != 110 {
		t.Fatalf("frontier = %d, want 110, the 2nd-highest answer", got.GetL2BlockNumber())
	}
	// Every endpoint answered, so the grace was not waited out.
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Fatalf("took %s: it waited for the grace although every endpoint had answered", elapsed)
	}
}

// A hung endpoint costs at most the grace, not the caller's whole deadline.
func TestAttestedUpTo_AHungEndpointCostsAtMostTheGrace(t *testing.T) {
	restore := frontierQuorumGrace
	frontierQuorumGrace = 50 * time.Millisecond
	t.Cleanup(func() { frontierQuorumGrace = restore })

	root := func(height uint64) *attestorpb.AttestedRoot { return &attestorpb.AttestedRoot{L2BlockNumber: height} }
	frontier, err := NewQuorumAttestationFrontier([]SigningAttestor{
		{Index: 0, Client: &fakeAttestor{frontierBlock: make(chan struct{})}}, // never answers
		{Index: 1, Client: &fakeAttestor{root: root(100), found: true}},
		{Index: 2, Client: &fakeAttestor{root: root(90), found: true}},
	}, 2)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := time.Now()
	got, found, err := frontier.AttestedUpTo(ctx, "op-sepolia", false)
	if err != nil || !found || got.GetL2BlockNumber() != 90 {
		t.Fatalf("frontier=%v found=%v err=%v, want 90 from the two that answered", got, found, err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("took %s: one hung endpoint held the call to the caller's deadline", elapsed)
	}
}

func TestQuorumFrontier(t *testing.T) {
	roots := []*attestorpb.AttestedRoot{{L2BlockNumber: 90}, {L2BlockNumber: 130}, {L2BlockNumber: 100}, {L2BlockNumber: 120}}
	for threshold, want := range map[int]uint64{1: 130, 2: 120, 3: 100, 4: 90} {
		if got := quorumFrontier(roots, threshold).GetL2BlockNumber(); got != want {
			t.Errorf("threshold %d: frontier %d, want %d", threshold, got, want)
		}
	}
	if roots[0].GetL2BlockNumber() != 90 {
		t.Error("quorumFrontier reordered its input")
	}
}
