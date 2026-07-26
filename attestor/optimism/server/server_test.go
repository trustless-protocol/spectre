package server_test

import (
	"context"
	"net"
	"path/filepath"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"attestor/optimism/client"
	"attestor/optimism/opstack"
	"attestor/optimism/server"
	attestorpb "attestor/types/attestor"
)

func newTypedClient(t *testing.T, addr string) *client.Client {
	t.Helper()
	cl, err := client.Dial(addr)
	if err != nil {
		t.Fatalf("client.Dial: %v", err)
	}
	t.Cleanup(func() { cl.Close() })
	return cl
}

func root(b byte) [32]byte {
	var r [32]byte
	r[0] = b
	return r
}

// newTestAttestor builds an attestor whose loop never runs — the server reads
// only the store and status snapshots, so nil game/replica sources are fine.
func newTestAttestor(t *testing.T, srcChain string) *opstack.OpStackAttestor {
	t.Helper()
	store, err := opstack.LoadStore(filepath.Join(t.TempDir(), "state.json"))
	if err != nil {
		t.Fatalf("LoadStore: %v", err)
	}
	cfg := opstack.Config{
		SrcChain:        srcChain,
		AttestationHead: opstack.HeadFinalized,
	}
	return opstack.New(cfg, nil, nil, store, nil, nil, nil)
}

// dial starts the server on a bufconn listener and returns a connected pb
// client plus the raw conn for the client-package test.
func dial(t *testing.T, s *server.Server) *grpc.ClientConn {
	t.Helper()
	lis := bufconn.Listen(1 << 20)
	grpcSrv := grpc.NewServer()
	s.Register(grpcSrv)
	go func() {
		_ = grpcSrv.Serve(lis)
	}()
	t.Cleanup(grpcSrv.Stop)
	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("grpc.NewClient: %v", err)
	}
	t.Cleanup(func() { conn.Close() })
	return conn
}

func TestQueriesRoundTrip(t *testing.T) {
	a := newTestAttestor(t, "op-test")
	now := time.Unix(4000, 0)
	a.Store().AppendDerived(100, root(0x01), now, false)
	a.Store().AppendDerived(300, root(0x03), now, true) // provisional

	conn := dial(t, server.New(map[string]*opstack.OpStackAttestor{"op-test": a}))
	c := attestorpb.NewAttestorServiceClient(conn)
	ctx := context.Background()

	up, err := c.AttestedUpTo(ctx, &attestorpb.AttestedUpToRequest{SrcChain: "op-test"})
	if err != nil {
		t.Fatalf("AttestedUpTo: %v", err)
	}
	if !up.GetFound() || up.GetRoot().GetL2BlockNumber() != 100 || up.GetRoot().GetProvisional() {
		t.Fatalf("confirmed frontier = %+v, want block 100 non-provisional", up)
	}
	upProv, err := c.AttestedUpTo(ctx, &attestorpb.AttestedUpToRequest{SrcChain: "op-test", IncludeProvisional: true})
	if err != nil {
		t.Fatalf("AttestedUpTo(provisional): %v", err)
	}
	if !upProv.GetFound() || upProv.GetRoot().GetL2BlockNumber() != 300 {
		t.Fatalf("provisional frontier = %+v, want block 300", upProv)
	}

	at, err := c.AttestedRootAtOrBelow(ctx, &attestorpb.AttestedRootAtOrBelowRequest{SrcChain: "op-test", L2BlockNumber: 250})
	if err != nil {
		t.Fatalf("AttestedRootAtOrBelow: %v", err)
	}
	if !at.GetFound() || at.GetRoot().GetL2BlockNumber() != 100 || at.GetRoot().GetRoot()[0] != 0x01 {
		t.Fatalf("at-or-below 250 = %+v, want block 100 root 0x01…", at)
	}
	miss, err := c.AttestedRootAtOrBelow(ctx, &attestorpb.AttestedRootAtOrBelowRequest{SrcChain: "op-test", L2BlockNumber: 50})
	if err != nil {
		t.Fatalf("AttestedRootAtOrBelow(50): %v", err)
	}
	if miss.GetFound() {
		t.Fatalf("at-or-below 50 = %+v, want not found", miss)
	}

	if _, err := c.AttestedUpTo(ctx, &attestorpb.AttestedUpToRequest{SrcChain: "nope"}); err == nil {
		t.Fatal("unknown chain must error")
	}

	info, err := c.Info(ctx, &attestorpb.InfoRequest{})
	if err != nil {
		t.Fatalf("Info: %v", err)
	}
	if len(info.GetChains()) != 1 || info.GetChains()[0].GetSrcChain() != "op-test" ||
		info.GetChains()[0].GetAttestationHead() != "finalized" || info.GetChains()[0].GetReplicaSeen() {
		t.Fatalf("Info = %+v, want one op-test chain, head finalized, replica unseen", info)
	}
}

func TestWatchStreamsFrontierAdvances(t *testing.T) {
	a := newTestAttestor(t, "op-test")
	now := time.Unix(4000, 0)
	a.Store().AppendDerived(100, root(0x01), now, false)

	srv := server.New(map[string]*opstack.OpStackAttestor{"op-test": a})
	srv.WatchPoll = 10 * time.Millisecond
	conn := dial(t, srv)
	c := attestorpb.NewAttestorServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream, err := c.WatchAttested(ctx, &attestorpb.WatchAttestedRequest{SrcChain: "op-test"})
	if err != nil {
		t.Fatalf("WatchAttested: %v", err)
	}

	first, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv(initial): %v", err)
	}
	if first.GetRoot().GetL2BlockNumber() != 100 {
		t.Fatalf("initial frontier = %+v, want block 100", first)
	}

	// A provisional append must NOT advance the confirmed watch; a confirmed
	// one must.
	a.Store().AppendDerived(200, root(0x02), now, true)
	a.Store().AppendDerived(300, root(0x03), now, false)
	next, err := stream.Recv()
	if err != nil {
		t.Fatalf("Recv(advance): %v", err)
	}
	if next.GetRoot().GetL2BlockNumber() != 300 {
		t.Fatalf("advanced frontier = %+v, want block 300 (skipping provisional 200)", next)
	}
}

func TestClientPackageRoundTrip(t *testing.T) {
	// The client package dials real addresses; reuse its conversion path by
	// pointing a raw pb client at bufconn above, and exercise the typed client
	// against a real TCP listener here.
	a := newTestAttestor(t, "op-test")
	a.Store().AppendDerived(100, root(0x0a), time.Unix(4000, 0), false)

	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	grpcSrv := grpc.NewServer()
	server.New(map[string]*opstack.OpStackAttestor{"op-test": a}).Register(grpcSrv)
	go func() { _ = grpcSrv.Serve(lis) }()
	t.Cleanup(grpcSrv.Stop)

	cl := newTypedClient(t, lis.Addr().String())
	ctx := context.Background()

	cur, ok, err := cl.AttestedUpTo(ctx, "op-test", false)
	if err != nil || !ok || cur.Height != 100 || cur.ID != root(0x0a) {
		t.Fatalf("client AttestedUpTo = (%+v, %t, %v), want height 100", cur, ok, err)
	}
	entry, ok, err := cl.AttestedRootAtOrBelow(ctx, "op-test", 500, false)
	if err != nil || !ok || entry.L2BlockNumber != 100 || entry.Source != opstack.SourceDerived {
		t.Fatalf("client AtOrBelow = (%+v, %t, %v), want derived entry at 100", entry, ok, err)
	}
}
