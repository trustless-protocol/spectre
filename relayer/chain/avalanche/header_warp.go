package avalanche

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"relayer/chain"
	"relayer/chain/l2rollup"
	relayerclient "relayer/client"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// warpHeaderBuilder assembles the update the Avalanche warp wasm client
// accepts: the canonical coreth header at the target height, the primary
// network's aggregate BLS signature over that block's hash (collected through
// the signature-aggregator sidecar's ACP-118 fan-out), and a router account
// proof fetched at the header's settled height (see corethProofQueryHeight —
// under asynchronous execution the header's state root commits the settled
// state, not the same-height one).
//
// Trust shape: the client verifies a ≥ quorum stake-weighted aggregate from
// its pinned canonical validator set. The builder holds no keys and attests
// nothing — a wrong header or a stale aggregate fails closed on-chain.
type warpHeaderBuilder struct {
	l2            *ethclient.Client
	router        ethcommon.Address
	networkID     uint32
	sourceChainID [32]byte
	aggregatorURL string
	quorum        uint64
	httpClient    *http.Client
	name          string
}

// NewWarpHeaderBuilder wires the Avalanche C-Chain warp header builder.
func NewWarpHeaderBuilder(l2 *ethclient.Client, router ethcommon.Address, networkID uint32, sourceChainID [32]byte, aggregatorURL string, quorum uint64, name string) (l2rollup.HeaderBuilder, error) {
	if l2 == nil {
		return nil, fmt.Errorf("%s: C-Chain client must not be nil", name)
	}
	if networkID == 0 {
		return nil, fmt.Errorf("%s: network id must not be zero", name)
	}
	if sourceChainID == [32]byte{} {
		return nil, fmt.Errorf("%s: source chain id must not be zero", name)
	}
	if aggregatorURL == "" {
		return nil, fmt.Errorf("%s: aggregator url must not be empty", name)
	}
	if quorum == 0 || quorum > 100 {
		return nil, fmt.Errorf("%s: quorum %d must be in 1..100", name, quorum)
	}
	return &warpHeaderBuilder{
		l2:            l2,
		router:        router,
		networkID:     networkID,
		sourceChainID: sourceChainID,
		aggregatorURL: aggregatorURL,
		quorum:        quorum,
		httpClient:    &http.Client{Timeout: 45 * time.Second},
		name:          name,
	}, nil
}

func (w *warpHeaderBuilder) Name() string { return w.name }

// BuildHeader packages the update for exactly request.Height. Aggregation has
// no expiry — validators re-sign any accepted block — so a retry simply
// re-aggregates; nothing is cached between attempts.
func (w *warpHeaderBuilder) BuildHeader(ctx context.Context, request l2rollup.HeaderRequest) (l2rollup.ClientMessage, uint64, error) {
	wireHeader, blockHash, _, proofHeight, err := readCorethHeader(ctx, w.l2.Client(), new(big.Int).SetUint64(request.Height))
	if err != nil {
		return nil, 0, chain.Transient(fmt.Errorf("%s: header at %d: %w", w.name, request.Height, err))
	}

	unsigned := buildBlockHashWarpMessage(w.networkID, w.sourceChainID, blockHash)
	bitSet, signature, err := w.aggregate(ctx, unsigned)
	if err != nil {
		return nil, 0, chain.Transient(fmt.Errorf("%s: aggregate for block %d: %w", w.name, request.Height, err))
	}

	routerProof, err := relayerclient.EthGetProof(ctx, w.l2, w.router, nil, new(big.Int).SetUint64(proofHeight))
	if err != nil {
		return nil, 0, chain.Transient(fmt.Errorf("%s: router proof at settled height %d: %w", w.name, proofHeight, err))
	}

	return &WarpSignedHeader{
		Header:       wireHeader,
		SignerBitSet: bitSet,
		Signature:    signature[:],
		RouterProof:  l2rollup.EvmAccountProof{Proof: routerProof.AccountProof},
	}, request.Height, nil
}

// aggregate collects the primary network's aggregate signature over the
// unsigned message through the signature-aggregator sidecar and re-splits the
// returned signed message, refusing an answer over different unsigned bytes.
func (w *warpHeaderBuilder) aggregate(ctx context.Context, unsigned []byte) ([]byte, [warpSignatureLen]byte, error) {
	var signature [warpSignatureLen]byte
	body, err := json.Marshal(map[string]any{
		"message":           "0x" + hex.EncodeToString(unsigned),
		"signing-subnet-id": PrimaryNetworkID,
		"quorum-percentage": w.quorum,
	})
	if err != nil {
		return nil, signature, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, w.aggregatorURL+"/aggregate-signatures", bytes.NewReader(body))
	if err != nil {
		return nil, signature, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := w.httpClient.Do(request)
	if err != nil {
		return nil, signature, err
	}
	defer response.Body.Close()
	var reply struct {
		SignedMessage string `json:"signed-message"`
	}
	if err := json.NewDecoder(response.Body).Decode(&reply); err != nil || reply.SignedMessage == "" {
		return nil, signature, fmt.Errorf("aggregator status %d: %v", response.StatusCode, err)
	}
	raw, err := hex.DecodeString(trimHexPrefix(reply.SignedMessage))
	if err != nil {
		return nil, signature, fmt.Errorf("decode signed message: %w", err)
	}
	unsignedEcho, bitSet, signature, err := splitSignedWarpMessage(raw)
	if err != nil {
		return nil, signature, err
	}
	if !bytes.Equal(unsignedEcho, unsigned) {
		return nil, signature, fmt.Errorf("aggregator answered over different unsigned bytes")
	}
	return bitSet, signature, nil
}

func trimHexPrefix(s string) string {
	if len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X') {
		return s[2:]
	}
	return s
}
