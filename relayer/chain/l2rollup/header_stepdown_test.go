package l2rollup

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"attestor/types/attestation"
	attestorpb "attestor/types/attestor"
	"relayer/chain"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
)

// startL2Node serves an L2 execution node holding blocks 1..head, each with a
// state root derived from its number. A block above head is answered null, which
// ethclient reads as ethereum.NotFound.
func startL2Node(t *testing.T, head uint64) *ethclient.Client {
	t.Helper()
	zero := "0x" + strings.Repeat("0", 64)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     json.RawMessage   `json:"id"`
			Method string            `json:"method"`
			Params []json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var result any
		switch req.Method {
		case "eth_getBlockByNumber":
			var tag string
			_ = json.Unmarshal(req.Params[0], &tag)
			number, err := strconv.ParseUint(strings.TrimPrefix(tag, "0x"), 16, 64)
			if err == nil && number > 0 && number <= head {
				result = map[string]any{
					"parentHash": zero, "sha3Uncles": zero, "miner": "0x" + strings.Repeat("0", 40),
					"stateRoot":        ethcommon.BigToHash(new(big.Int).SetUint64(number)).Hex(),
					"transactionsRoot": zero, "receiptsRoot": zero, "logsBloom": "0x" + strings.Repeat("0", 512),
					"difficulty": "0x0", "number": hexutil.EncodeUint64(number), "gasLimit": "0x1c9c380",
					"gasUsed": "0x0", "timestamp": hexutil.EncodeUint64(1_700_000_000 + number), "extraData": "0x",
					"mixHash": zero, "nonce": "0x0000000000000000", "baseFeePerGas": "0x1",
				}
			}
		case "eth_getProof":
			result = map[string]any{"accountProof": []string{}, "storageHash": zero, "storageProof": []any{}}
		default:
			http.Error(w, "unexpected method "+req.Method, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"jsonrpc": "2.0", "id": req.ID, "result": result})
	}))
	t.Cleanup(srv.Close)
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(client.Close)
	return client
}

// replicaAttestor signs any block up to its replica head and reports claim as its
// frontier. A claim above head is the over-reporting endpoint.
type replicaAttestor struct {
	head      uint64
	claim     uint64
	hang      bool  // AttestedUpTo never answers
	verifyErr error // VerifyStateRoot fails with it, whatever the block
	signer    attestation.Signer

	mu        sync.Mutex
	asked     []uint64 // blocks VerifyStateRoot was asked about
	frontiers int      // AttestedUpTo calls
}

func (r *replicaAttestor) AttestedUpTo(ctx context.Context, _ string, _ bool) (*attestorpb.AttestedRoot, bool, error) {
	r.mu.Lock()
	r.frontiers++
	r.mu.Unlock()
	if r.hang {
		<-ctx.Done()
		return nil, false, ctx.Err()
	}
	return &attestorpb.AttestedRoot{L2BlockNumber: r.claim}, true, nil
}

func (r *replicaAttestor) AttestedRootAtOrBelow(context.Context, string, uint64, bool) (*attestorpb.AttestedRoot, bool, error) {
	return nil, false, errors.New("not used by the header builder")
}

func (r *replicaAttestor) VerifyStateRoot(_ context.Context, request VerificationRequest) (SignedVerdict, error) {
	r.mu.Lock()
	r.asked = append(r.asked, request.BlockNumber)
	r.mu.Unlock()
	if r.verifyErr != nil {
		return SignedVerdict{}, r.verifyErr
	}
	if request.BlockNumber > r.head {
		return SignedVerdict{}, fmt.Errorf("%w: block %d is above the replica head %d", ErrAttestorReplicaBehind, request.BlockNumber, r.head)
	}
	signature, err := r.signer.Sign(request.L2Router[:], request.AttestorSetHash[:], request.BlockNumber, request.BlockHash, request.StateRoot)
	if err != nil {
		return SignedVerdict{}, err
	}
	return SignedVerdict{Valid: true, BlockNumber: request.BlockNumber, BlockHash: request.BlockHash, StateRoot: request.StateRoot, Signature: signature}, nil
}

func (r *replicaAttestor) calls() ([]uint64, int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]uint64(nil), r.asked...), r.frontiers
}

// stepDownBuilder wires a threshold builder to the replicas against an L2 node
// whose head is l2Head. An endpoint's index is its key's place in the set.
func stepDownBuilder(t *testing.T, threshold uint16, l2Head uint64, replicas ...*replicaAttestor) (*attestedHeaderBuilder, map[*replicaAttestor]uint16) {
	t.Helper()
	for i, replica := range replicas {
		signer, err := attestation.NewSigner(8453, fmt.Sprintf("%064x", i+40))
		if err != nil {
			t.Fatal(err)
		}
		replica.signer = signer
	}
	sorted := append([]*replicaAttestor(nil), replicas...)
	sort.Slice(sorted, func(i, j int) bool {
		return string(sorted[i].signer.PublicKey()) < string(sorted[j].signer.PublicKey())
	})
	set := attestation.AttestorConfig{Threshold: threshold}
	endpoints := make([]SigningAttestor, len(sorted))
	index := make(map[*replicaAttestor]uint16, len(sorted))
	for i, replica := range sorted {
		set.PublicKeys = append(set.PublicKeys, replica.signer.PublicKey())
		endpoints[i] = SigningAttestor{Index: uint16(i), Client: replica}
		index[replica] = uint16(i)
	}
	built, err := NewAttestedHeaderBuilder(startL2Node(t, l2Head), ethcommon.HexToAddress("0x1111111111111111111111111111111111111111"),
		8453, set, endpoints, "op", attestorpb.RunMode_RUN_MODE_UNSAFE, true, "l2-op")
	if err != nil {
		t.Fatal(err)
	}
	return built.(*attestedHeaderBuilder), index
}

func signerIndices(t *testing.T, message ClientMessage) (uint64, []uint16) {
	t.Helper()
	header, ok := message.(*SignedAttestedL2Header)
	if !ok {
		t.Fatalf("message is %T, want *SignedAttestedL2Header", message)
	}
	indices := make([]uint16, 0, len(header.AttestorSignature))
	for _, signature := range header.AttestorSignature {
		indices = append(indices, signature.AttestorIndex)
	}
	return header.L2Header.Number, indices
}

// #575 review: frontier answers are not signed. One endpoint that reports more
// than its replica holds lifts the threshold-th highest answer to a block only
// threshold-1 honest replicas have reached, and the header failed there on every
// pass where the honest frontiers differed.
func TestBuildHeader_StepsDownPastAnOverReportingAttestor(t *testing.T) {
	fast := &replicaAttestor{head: 1000, claim: 1000}
	slow := &replicaAttestor{head: 990, claim: 990}
	liar := &replicaAttestor{head: 0, claim: 5000}
	b, index := stepDownBuilder(t, 2, 1100, fast, slow, liar)

	frontier, err := NewQuorumAttestationFrontier(b.attestors, 2)
	if err != nil {
		t.Fatal(err)
	}
	root, _, err := frontier.AttestedUpTo(context.Background(), "op", true)
	if err != nil || root.GetL2BlockNumber() != 1000 {
		t.Fatalf("frontier = %v (%v), want 1000: the liar lifts it to the fast replica's block", root, err)
	}

	message, committed, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 1000})
	if err != nil {
		t.Fatalf("BuildHeader(1000): %v", err)
	}
	number, signers := signerIndices(t, message)
	if committed != 990 || number != 990 {
		t.Fatalf("committed %d, header %d; want 990, the slow replica's frontier", committed, number)
	}
	want := []uint16{index[fast], index[slow]}
	sort.Slice(want, func(i, j int) bool { return want[i] < want[j] })
	if fmt.Sprint(signers) != fmt.Sprint(want) {
		t.Fatalf("signed by %v, want the two honest replicas %v", signers, want)
	}
}

// The healthy case costs nothing extra: no frontier round, one signature round.
func TestBuildHeader_AnHonestQuorumDoesNotStepDown(t *testing.T) {
	replicas := []*replicaAttestor{{head: 1000, claim: 1000}, {head: 995, claim: 995}, {head: 990, claim: 990}}
	b, _ := stepDownBuilder(t, 2, 1100, replicas...)

	_, committed, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 995})
	if err != nil || committed != 995 {
		t.Fatalf("BuildHeader(995) = %d, %v; want 995", committed, err)
	}
	for _, replica := range replicas {
		if _, frontiers := replica.calls(); frontiers != 0 {
			t.Fatalf("the builder asked for frontiers %d time(s) on a block the quorum signed", frontiers)
		}
	}
}

// With threshold 1 the frontier is the highest answer, so a liar can name a block
// the L2 RPC does not have yet.
func TestBuildHeader_StepsDownWhenTheL2RPCHasNoBlock(t *testing.T) {
	honest := &replicaAttestor{head: 1000, claim: 1000}
	liar := &replicaAttestor{head: 0, claim: 5000}
	b, _ := stepDownBuilder(t, 1, 1100, honest, liar)

	_, committed, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 5000})
	if err != nil || committed != 1000 {
		t.Fatalf("BuildHeader(5000) = %d, %v; want 1000", committed, err)
	}
}

// A failure no lower block can fix is returned as it is, without a frontier round.
func TestBuildHeader_DoesNotStepDownPastAPermanentFailure(t *testing.T) {
	badRequest := fmt.Errorf("%w: invalid run mode", ErrAttestorBadRequest)

	t.Run("a refused request", func(t *testing.T) {
		refusing := &replicaAttestor{head: 1000, claim: 1000, verifyErr: badRequest}
		b, _ := stepDownBuilder(t, 1, 1100, refusing)
		_, _, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 1000})
		if err == nil || !chain.IsPermanent(err) || !errors.Is(err, ErrAttestorBadRequest) {
			t.Fatalf("err = %v, want the permanent bad-request failure", err)
		}
		if asked, frontiers := refusing.calls(); len(asked) != 1 || frontiers != 0 {
			t.Fatalf("asked %v, %d frontier round(s); want one signature request and no step-down", asked, frontiers)
		}
	})

	// One replica is behind, but the other can never sign: no block completes it.
	t.Run("a quorum that cannot complete", func(t *testing.T) {
		behind := &replicaAttestor{head: 0, claim: 990}
		refusing := &replicaAttestor{head: 1000, claim: 1000, verifyErr: badRequest}
		b, _ := stepDownBuilder(t, 2, 1100, behind, refusing)
		_, _, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 1000})
		if err == nil || !chain.IsPermanent(err) || !errors.Is(err, ErrAttestorReplicaBehind) {
			t.Fatalf("err = %v, want the permanent quorum failure", err)
		}
		for _, replica := range []*replicaAttestor{behind, refusing} {
			if _, frontiers := replica.calls(); frontiers != 0 {
				t.Fatalf("%d frontier round(s) for a quorum no block completes", frontiers)
			}
		}
	})
}

// Each endpoint names at most one candidate; when none is signed the first
// failure is returned, still retryable.
func TestBuildHeader_ReturnsTheCauseWhenNoLowerFrontierIsSigned(t *testing.T) {
	fast := &replicaAttestor{head: 1000, claim: 1000}
	stuck := &replicaAttestor{head: 0, claim: 990}
	liar := &replicaAttestor{head: 0, claim: 5000}
	b, _ := stepDownBuilder(t, 2, 1100, fast, stuck, liar)

	_, _, err := b.BuildHeader(context.Background(), HeaderRequest{Height: 1000})
	if err == nil || chain.IsPermanent(err) || !strings.Contains(err.Error(), "L2 block 1000") {
		t.Fatalf("err = %v, want the retryable shortfall at 1000", err)
	}
	if asked, _ := fast.calls(); fmt.Sprint(asked) != "[1000 990]" {
		t.Fatalf("the fast replica was asked about %v, want [1000 990]: one round per lower frontier", asked)
	}
}

// A hung endpoint costs the step-down its candidate, not the build.
func TestBuildHeader_AHungEndpointDoesNotStallTheStepDown(t *testing.T) {
	previous := stepDownFrontierTimeout
	stepDownFrontierTimeout = 50 * time.Millisecond
	t.Cleanup(func() { stepDownFrontierTimeout = previous })

	b, _ := stepDownBuilder(t, 2, 1100,
		&replicaAttestor{head: 1000, claim: 1000}, &replicaAttestor{head: 990, claim: 990}, &replicaAttestor{head: 0, hang: true})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, committed, err := b.BuildHeader(ctx, HeaderRequest{Height: 1000})
	if err != nil || committed != 990 {
		t.Fatalf("BuildHeader(1000) = %d, %v; want 990", committed, err)
	}
}

func TestStepDownCandidates(t *testing.T) {
	frontiers := []endpointFrontier{{0, 1000}, {1, 990}, {2, 5000}, {3, 990}, {4, 0}, {5, 980}}
	if got := stepDownCandidates(frontiers, 1000); fmt.Sprint(got) != "[990 980]" {
		t.Fatalf("candidates = %v, want [990 980]: distinct, below the failed block, highest first", got)
	}
	if got := describeFrontiers(frontiers[:3]); got != "attestor 0 at 1000, attestor 1 at 990, attestor 2 at 5000" {
		t.Fatalf("describeFrontiers = %q", got)
	}
}
