// Package attestorgrpc binds the L2 attestor sidecar's gRPC AttestorService
// (#240/#258) to the l2rollup.AttestorClient interface the source gates on. It is a
// thin transport layer: the generated proto AttestedRoot is returned as-is (the
// attestor module owns that shape), so there is no mapping to keep in sync.
package attestorgrpc

import (
	"context"
	"fmt"

	attestorpb "attestor/types/attestor"
	"relayer/chain/l2rollup"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// Client adapts the generated AttestorServiceClient to l2rollup.AttestorClient.
type Client struct {
	conn *grpc.ClientConn
	rpc  attestorpb.AttestorServiceClient
}

var _ l2rollup.AttestorClient = (*Client)(nil)

// Dial connects to the attestor sidecar at addr (host:port). The transport is
// insecure (the sidecar is a co-located trusted process); the caller must Close it.
func Dial(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("attestorgrpc: dial %s: %w", addr, err)
	}
	return &Client{conn: conn, rpc: attestorpb.NewAttestorServiceClient(conn)}, nil
}

// Close releases the gRPC connection.
func (c *Client) Close() error { return c.conn.Close() }

// AttestedUpTo returns the attestor's highest independently-re-derived frontier.
func (c *Client) AttestedUpTo(ctx context.Context, srcChain string, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error) {
	resp, err := c.rpc.AttestedUpTo(ctx, &attestorpb.AttestedUpToRequest{
		SrcChain:           srcChain,
		IncludeProvisional: includeProvisional,
	})
	if err != nil {
		return nil, false, fmt.Errorf("attestorgrpc: AttestedUpTo(%s): %w", srcChain, err)
	}
	if !resp.GetFound() {
		return nil, false, nil
	}
	if resp.GetRoot() == nil {
		return nil, false, fmt.Errorf("attestorgrpc: AttestedUpTo(%s) found=true but root is nil", srcChain)
	}
	return resp.GetRoot(), true, nil
}

// AttestedRootAtOrBelow returns the best attested root at or below l2BlockNumber.
func (c *Client) AttestedRootAtOrBelow(ctx context.Context, srcChain string, l2BlockNumber uint64, includeProvisional bool) (*attestorpb.AttestedRoot, bool, error) {
	resp, err := c.rpc.AttestedRootAtOrBelow(ctx, &attestorpb.AttestedRootAtOrBelowRequest{
		SrcChain:           srcChain,
		L2BlockNumber:      l2BlockNumber,
		IncludeProvisional: includeProvisional,
	})
	if err != nil {
		return nil, false, fmt.Errorf("attestorgrpc: AttestedRootAtOrBelow(%s, %d): %w", srcChain, l2BlockNumber, err)
	}
	if !resp.GetFound() {
		return nil, false, nil
	}
	if resp.GetRoot() == nil {
		return nil, false, fmt.Errorf("attestorgrpc: AttestedRootAtOrBelow(%s, %d) found=true but root is nil", srcChain, l2BlockNumber)
	}
	return resp.GetRoot(), true, nil
}

// VerifyStateRoot compares one block identity against the attestor's replica.
//
// Every in-tree attestor serves this RPC. An attestor binary older than that
// answers Unimplemented, which is reported as ErrVerifyStateRootUnsupported so the
// caller can tell "this attestor cannot answer" apart from "this attestor says no"
// — the two must not collapse, because the first is a version skew and the second
// is a divergence.
func (c *Client) VerifyStateRoot(ctx context.Context, srcChain string, l2BlockNumber uint64, stateRoot, blockHash []byte, runMode attestorpb.RunMode) (l2rollup.VerifiedStateRoot, error) {
	resp, err := c.rpc.VerifyStateRoot(ctx, &attestorpb.VerifyStateRootRequest{
		SrcChain:          srcChain,
		BlockNumber:       l2BlockNumber,
		ExpectedStateRoot: stateRoot,
		ExpectedBlockHash: blockHash,
		RunMode:           runMode,
	})
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return l2rollup.VerifiedStateRoot{}, fmt.Errorf("%w: %v", l2rollup.ErrVerifyStateRootUnsupported, err)
		}
		return l2rollup.VerifiedStateRoot{}, fmt.Errorf("attestorgrpc: VerifyStateRoot(%d): %w", l2BlockNumber, err)
	}
	return l2rollup.VerifiedStateRoot{
		Valid:     resp.GetValid(),
		Signature: append([]byte(nil), resp.GetAttestationSignature()...),
	}, nil
}
