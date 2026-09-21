// warp-spike measures whether a stake-weighted BLS quorum of the Avalanche
// primary network can be aggregated from outside Avalanche — the one
// assumption the warp light-client design depends on that nothing had
// demonstrated. It drives Ava Labs' signature-aggregator sidecar for
// block-hash payloads over live C-Chain blocks and INDEPENDENTLY re-verifies
// every aggregate with avalanchego's own warp verification code against the
// canonical validator set, so a "success" is a verified aggregate, not a
// trusted HTTP response.
//
// It deliberately lives outside relayer/go.mod: it imports avalanchego, which
// the relayer must not depend on.
//
// Usage (run the signature-aggregator first, pointed at the same network):
//
//	go run . -count 20 -interval 30s -csv out.csv
package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/constants"
	"github.com/ava-labs/avalanchego/utils/set"
	"github.com/ava-labs/avalanchego/vms/platformvm"
	platformapi "github.com/ava-labs/avalanchego/vms/platformvm/api"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp/payload"
)

func main() {
	baseURL := flag.String("base-url", "https://api.avax-test.network", "Avalanche API node (info + P-Chain + C-Chain)")
	aggregatorURL := flag.String("aggregator", "http://127.0.0.1:18080", "signature-aggregator base URL")
	quorum := flag.Uint64("quorum", 67, "quorum percentage to request and verify")
	count := flag.Int("count", 1, "number of aggregation attempts")
	interval := flag.Duration("interval", 30*time.Second, "delay between attempts")
	csvPath := flag.String("csv", "", "append results as CSV (default stdout only)")
	flag.Parse()

	ctx := context.Background()
	if err := run(ctx, *baseURL, *aggregatorURL, *quorum, *count, *interval, *csvPath); err != nil {
		fmt.Fprintln(os.Stderr, "warp-spike:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, baseURL, aggregatorURL string, quorum uint64, count int, interval time.Duration, csvPath string) error {
	var networkID uint32
	if err := rpcCall(ctx, baseURL+"/ext/info", "info.getNetworkID", map[string]any{}, &struct {
		NetworkID *jsonUint32 `json:"networkID"`
	}{(*jsonUint32)(&networkID)}); err != nil {
		return fmt.Errorf("info.getNetworkID: %w", err)
	}
	var chainIDStr struct {
		BlockchainID string `json:"blockchainID"`
	}
	if err := rpcCall(ctx, baseURL+"/ext/info", "info.getBlockchainID", map[string]any{"alias": "C"}, &chainIDStr); err != nil {
		return fmt.Errorf("info.getBlockchainID: %w", err)
	}
	cChainID, err := ids.FromString(chainIDStr.BlockchainID)
	if err != nil {
		return fmt.Errorf("parse C-Chain id: %w", err)
	}

	pClient := platformvm.NewClient(baseURL)

	var out *os.File
	if csvPath != "" {
		out, err = os.OpenFile(csvPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return err
		}
		defer out.Close()
	}
	emit := func(line string) {
		fmt.Println(line)
		if out != nil {
			fmt.Fprintln(out, line)
		}
	}
	emit("timestamp,block_number,latency_ms,ok,signers,signed_weight_pct,total_validators,error")

	successes := 0
	for i := 0; i < count; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(interval):
			}
		}
		row, ok := attempt(ctx, baseURL, aggregatorURL, networkID, cChainID, pClient, quorum)
		emit(row)
		if ok {
			successes++
		}
	}
	fmt.Printf("# %d/%d attempts produced a verified %d%% aggregate\n", successes, count, quorum)
	return nil
}

// attempt runs one aggregate-and-verify cycle and returns a CSV row.
func attempt(ctx context.Context, baseURL, aggregatorURL string, networkID uint32, cChainID ids.ID, pClient *platformvm.Client, quorum uint64) (string, bool) {
	now := time.Now().UTC().Format(time.RFC3339)
	fail := func(blockNumber uint64, latency time.Duration, err error) (string, bool) {
		msg := strings.ReplaceAll(err.Error(), ",", ";")
		return fmt.Sprintf("%s,%d,%d,false,0,0,0,%s", now, blockNumber, latency.Milliseconds(), msg), false
	}

	// 1. Latest accepted C-Chain block (finalized == accepted on Avalanche).
	var blk struct {
		Number string `json:"number"`
		Hash   string `json:"hash"`
	}
	if err := ethCall(ctx, baseURL+"/ext/bc/C/rpc", "eth_getBlockByNumber", []any{"finalized", false}, &blk); err != nil {
		return fail(0, 0, fmt.Errorf("eth_getBlockByNumber: %w", err))
	}
	var blockNumber uint64
	fmt.Sscanf(blk.Number, "0x%x", &blockNumber)
	hashBytes, err := hex.DecodeString(strings.TrimPrefix(blk.Hash, "0x"))
	if err != nil || len(hashBytes) != 32 {
		return fail(blockNumber, 0, fmt.Errorf("bad block hash %q", blk.Hash))
	}

	// 2. Unsigned warp message: Hash payload over the EVM block hash, built
	//    with avalanchego's own codecs.
	hashPayload, err := payload.NewHash(ids.ID(hashBytes))
	if err != nil {
		return fail(blockNumber, 0, err)
	}
	unsigned, err := warp.NewUnsignedMessage(networkID, cChainID, hashPayload.Bytes())
	if err != nil {
		return fail(blockNumber, 0, err)
	}

	// 3. Aggregate through the sidecar.
	start := time.Now()
	reqBody, _ := json.Marshal(map[string]any{
		"message":           "0x" + hex.EncodeToString(unsigned.Bytes()),
		"signing-subnet-id": constants.PrimaryNetworkID.String(),
		"quorum-percentage": quorum,
	})
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, aggregatorURL+"/aggregate-signatures", bytes.NewReader(reqBody))
	if err != nil {
		return fail(blockNumber, 0, err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fail(blockNumber, time.Since(start), err)
	}
	defer resp.Body.Close()
	var aggReply struct {
		SignedMessage string `json:"signed-message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&aggReply); err != nil || aggReply.SignedMessage == "" {
		return fail(blockNumber, time.Since(start), fmt.Errorf("aggregator status %d: %v", resp.StatusCode, err))
	}
	latency := time.Since(start)

	// 4. Independent verification: parse the signed message and check the
	//    BLS aggregate against the canonical primary-network validator set,
	//    with avalanchego's own Verify.
	signedBytes, err := hex.DecodeString(strings.TrimPrefix(aggReply.SignedMessage, "0x"))
	if err != nil {
		return fail(blockNumber, latency, err)
	}
	signed, err := warp.ParseMessage(signedBytes)
	if err != nil {
		return fail(blockNumber, latency, fmt.Errorf("parse signed message: %w", err))
	}
	if !bytes.Equal(signed.UnsignedMessage.Bytes(), unsigned.Bytes()) {
		return fail(blockNumber, latency, fmt.Errorf("aggregator returned a different unsigned message"))
	}
	warpSets, err := pClient.GetAllValidatorsAt(ctx, platformapi.ProposedHeight)
	if err != nil {
		return fail(blockNumber, latency, fmt.Errorf("getAllValidatorsAt: %w", err))
	}
	warpSet, ok := warpSets[constants.PrimaryNetworkID]
	if !ok {
		return fail(blockNumber, latency, fmt.Errorf("no primary network validator set in reply"))
	}
	if err := signed.Signature.Verify(&signed.UnsignedMessage, networkID, warpSet, quorum, 100); err != nil {
		return fail(blockNumber, latency, fmt.Errorf("aggregate does not verify: %w", err))
	}

	// 5. Reporting: signer count + exact signed weight from the bitset.
	bitSetSig, ok := signed.Signature.(*warp.BitSetSignature)
	if !ok {
		return fail(blockNumber, latency, fmt.Errorf("unexpected signature type %T", signed.Signature))
	}
	signerIndices := set.BitsFromBytes(bitSetSig.Signers)
	signers := 0
	var signedWeight uint64
	for i, v := range warpSet.Validators {
		if signerIndices.Contains(i) {
			signers++
			signedWeight += v.Weight
		}
	}
	pct := float64(signedWeight) * 100 / float64(warpSet.TotalWeight)
	return fmt.Sprintf("%s,%d,%d,true,%d,%.2f,%d,", now, blockNumber, latency.Milliseconds(), signers, pct, len(warpSet.Validators)), true
}

// jsonUint32 decodes the string-encoded numbers the info API returns.
type jsonUint32 uint32

func (j *jsonUint32) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		var n uint32
		if err := json.Unmarshal(b, &n); err != nil {
			return err
		}
		*j = jsonUint32(n)
		return nil
	}
	var n uint32
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return err
	}
	*j = jsonUint32(n)
	return nil
}

// rpcCall performs a JSON-RPC 2.0 call and decodes result into out.
func rpcCall(ctx context.Context, url, method string, params any, out any) error {
	body, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "method": method, "params": params, "id": 1})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if envelope.Error != nil {
		return fmt.Errorf("%s", envelope.Error.Message)
	}
	return json.Unmarshal(envelope.Result, out)
}

// ethCall is rpcCall with positional params (the eth namespace).
func ethCall(ctx context.Context, url, method string, params []any, out any) error {
	return rpcCall(ctx, url, method, params, out)
}
