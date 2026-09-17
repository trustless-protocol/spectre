package transaction

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"relayer/services"
)

// abandonedTxFixture is the state the give-up paths are in: a signed payload
// already broadcast at some nonce, and a node that will be asked to displace it.
type abandonedTxFixture struct {
	node    *mockJSONRPC
	client  *ethclient.Client
	key     *ecdsa.PrivateKey
	from    common.Address
	chainID *big.Int
	tx      *types.Transaction
}

func newAbandonedTxFixture(t *testing.T, nonce uint64) abandonedTxFixture {
	t.Helper()

	node := newMockJSONRPC()
	srv := httptest.NewServer(node)
	t.Cleanup(srv.Close)

	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial mock RPC: %v", err)
	}
	t.Cleanup(client.Close)

	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	chainID := big.NewInt(1)
	payload := types.NewTx(&types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: big.NewInt(2_000_000_000),
		GasFeeCap: big.NewInt(4_000_000_000),
		Gas:       100_000,
		To:        &common.Address{0x1},
		Data:      []byte{0xde, 0xad, 0xbe, 0xef},
	})
	signed, err := types.SignTx(payload, types.LatestSignerForChainID(chainID), key)
	if err != nil {
		t.Fatalf("sign payload: %v", err)
	}

	return abandonedTxFixture{
		node:    node,
		client:  client,
		key:     key,
		from:    crypto.PubkeyToAddress(key.PublicKey),
		chainID: chainID,
		tx:      signed,
	}
}

func (f abandonedTxFixture) endpoint() services.EVMEndpoint {
	return services.EVMEndpoint{Client: f.client}
}

// A transaction the relayer stops waiting for does not stop existing: it stays
// queued, and the next attempt takes a LOWER nonce, so closing the gap lets the
// payload execute after the relayer already resubmitted it -- the action runs
// twice (#410). Cancelling means putting a no-op at that nonce first.
func TestCancelAbandonedTx(t *testing.T) {
	t.Run("puts a zero-value self-transfer at the abandoned nonce", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)

		outcome, err := (&Handler{}).cancelAbandonedTx(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID)
		if err != nil {
			t.Fatalf("cancel abandoned tx: %v", err)
		}
		if outcome != cancelBroadcast {
			t.Fatalf("outcome = %v, want cancelBroadcast: the node accepted it, which is all this can claim", outcome)
		}

		sent := f.node.sentTransactions()
		if len(sent) != 1 {
			t.Fatalf("broadcast %d transactions, want exactly the cancellation", len(sent))
		}
		cancel := sent[0]

		// The nonce is the whole point: a cancellation anywhere else leaves the
		// payload executable.
		if cancel.Nonce() != f.tx.Nonce() {
			t.Errorf("cancellation nonce = %d, want the abandoned %d", cancel.Nonce(), f.tx.Nonce())
		}
		if cancel.To() == nil || *cancel.To() != f.from {
			t.Errorf("cancellation to = %v, want a self-transfer to %s", cancel.To(), f.from.Hex())
		}
		if cancel.Value().Sign() != 0 {
			t.Errorf("cancellation value = %s, want zero", cancel.Value())
		}
		if len(cancel.Data()) != 0 {
			t.Errorf("cancellation carries %d bytes of calldata, want none", len(cancel.Data()))
		}

		// Below the original's fee the node keeps the original and the
		// cancellation is the transaction that gets dropped.
		if cancel.GasTipCap().Cmp(f.tx.GasTipCap()) <= 0 {
			t.Errorf("cancellation tip cap %s does not exceed the abandoned %s; it would not displace it",
				cancel.GasTipCap(), f.tx.GasTipCap())
		}
		if cancel.GasFeeCap().Cmp(f.tx.GasFeeCap()) <= 0 {
			t.Errorf("cancellation fee cap %s does not exceed the abandoned %s; it would not displace it",
				cancel.GasFeeCap(), f.tx.GasFeeCap())
		}
	})

	t.Run("a nonce the account is already past needs no cancellation", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)
		f.node.sendErrFor = func(*types.Transaction) error { return errors.New("nonce too low") }

		// A consumed nonce says something took it, and the likeliest something is
		// the payload itself -- which would mean the operation SUCCEEDED. Reported
		// by @DongLieu: this used to report success, and the caller then resubmitted
		// a payload that had already executed. It is unresolved, not cancelled.
		outcome, err := (&Handler{}).cancelAbandonedTx(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID)
		if outcome != cancelUnresolved {
			t.Fatalf("outcome = %v, want cancelUnresolved: the nonce may have been consumed by the payload", outcome)
		}
		if err == nil {
			t.Fatal("a consumed nonce reported as a clean cancellation; the payload may have executed")
		}
	})

	t.Run("a cancellation already in the pool is not a failure", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)
		f.node.sendErrFor = func(*types.Transaction) error { return errors.New("already known") }

		outcome, err := (&Handler{}).cancelAbandonedTx(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID)
		if err != nil {
			t.Fatalf("an already-broadcast cancellation reported as a failure: %v", err)
		}
		if outcome != cancelBroadcast {
			t.Fatalf("outcome = %v, want cancelBroadcast", outcome)
		}
	})

	t.Run("any other rejection is reported", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)
		f.node.sendErrFor = func(*types.Transaction) error { return errors.New("replacement transaction underpriced") }

		outcome, err := (&Handler{}).cancelAbandonedTx(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID)
		if outcome != cancelUnresolved {
			t.Fatalf("outcome = %v, want cancelUnresolved", outcome)
		}
		if err == nil {
			t.Fatal("a rejected cancellation reported as success; the payload can still execute")
		}
		if !strings.Contains(err.Error(), "42") {
			t.Errorf("error = %v, want it to name the nonce left outstanding", err)
		}
	})
}

// onGiveUp is the one exit the wait loop takes when it stops waiting, so the
// classification of the result is decided in a single place.
//
// That classification changed after review. Reported by @DongLieu: a broadcast
// cancellation is acceptance by ONE node, not inclusion, so "the cancellation
// landed" was never a thing this code could know -- and returning a plain
// retryable error let the payload be resubmitted while the original was still
// live. Every ending now carries ErrUnresolvedEVMNonce.
func TestOnGiveUp(t *testing.T) {
	cause := errors.New("transaction wait mined timed out")
	noHashes := []common.Hash(nil)

	t.Run("a broadcast cancellation is still unresolved", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)

		err := (&Handler{}).onGiveUp(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID, noHashes, cause)
		if !errors.Is(err, cause) {
			t.Fatalf("error = %v, want the wait failure still reachable", err)
		}
		if !errors.Is(err, services.ErrUnresolvedEVMNonce) {
			t.Fatalf("error = %v, want unresolved: one node accepting the replacement does not decide the nonce", err)
		}
		if !errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("error = %v, want permanent: nothing derived from this attempt may be retried", err)
		}
	})

	t.Run("a failed cancellation is unresolved too", func(t *testing.T) {
		f := newAbandonedTxFixture(t, 42)
		f.node.sendErrFor = func(*types.Transaction) error { return errors.New("replacement transaction underpriced") }

		err := (&Handler{}).onGiveUp(
			context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID, noHashes, cause)
		if !errors.Is(err, services.ErrUnresolvedEVMNonce) {
			t.Fatalf("error = %v, want unresolved", err)
		}
		if !errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("error = %v, want permanent", err)
		}
		// Permanent must not cost the diagnosis of what actually went wrong.
		if !errors.Is(err, cause) {
			t.Errorf("error = %v, want the original wait failure still reachable", err)
		}
	})

	t.Run("invalidates the cached nonce whether or not the cancellation landed", func(t *testing.T) {
		for _, tc := range []struct {
			name    string
			sendErr error
		}{
			{name: "cancellation landed"},
			{name: "cancellation rejected", sendErr: errors.New("replacement transaction underpriced")},
		} {
			t.Run(tc.name, func(t *testing.T) {
				f := newAbandonedTxFixture(t, 42)
				if tc.sendErr != nil {
					f.node.sendErrFor = func(*types.Transaction) error { return tc.sendErr }
				}

				h := &Handler{}
				h.mu.Lock()
				state := h.senderState(f.chainID.String(), f.from)
				state.nonce = 42
				state.nonceValid = true
				h.mu.Unlock()

				_ = h.onGiveUp(context.Background(), f.endpoint(), f.key, f.tx, f.from, f.chainID, nil, cause)

				h.mu.Lock()
				stillValid := h.senderState(f.chainID.String(), f.from).nonceValid
				h.mu.Unlock()
				if stillValid {
					t.Fatal("cached nonce still valid; the next attempt would reuse a nonce the chain disagrees with")
				}
			})
		}
	})
}

// The #410 regression at the exit that reaches it soonest: the pending-nonce
// probe sees a gap, gives up, and the caller then resubmits the payload at the
// lower nonce the chain actually needs. Without a cancellation broadcast at the
// abandoned nonce first, closing the gap executes the payload a second time.
func TestExecuteWithRetryAndResubmission_AbandonedNonceIsCancelled(t *testing.T) {
	oldTimeout := ethTxReceiptTimeout
	oldPoll := ethTxReceiptPollInterval
	ethTxReceiptTimeout = 60 * time.Millisecond
	ethTxReceiptPollInterval = 10 * time.Millisecond
	defer func() {
		ethTxReceiptTimeout = oldTimeout
		ethTxReceiptPollInterval = oldPoll
	}()

	node := newMockJSONRPC()
	node.nonce = 7 // what the chain needs; the cache below is far ahead of it
	srv := httptest.NewServer(node)
	defer srv.Close()
	client, err := ethclient.Dial(srv.URL)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	privKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	from := crypto.PubkeyToAddress(privKey.PublicKey)

	const abandonedNonce = 50
	h := &Handler{}
	h.mu.Lock()
	state := h.senderState("1", from)
	state.nonce = abandonedNonce
	state.nonceValid = true
	h.mu.Unlock()

	// A bound contract call signs AND broadcasts, and the hazard only exists
	// because the payload really did reach the mempool -- so the stand-in has to
	// broadcast too, or the test would be asserting against a chain that never
	// saw it.
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
			GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
			Data: []byte{0xde, 0xad, 0xbe, 0xef},
		})
		signed, signErr := auth.Signer(auth.From, tx)
		if signErr != nil {
			return nil, signErr
		}
		if sendErr := client.SendTransaction(auth.Context, signed); sendErr != nil {
			return nil, sendErr
		}
		return signed, nil
	}

	_, _, _, err = h.executeWithRetryAndResubmission(
		context.Background(), services.EVMEndpoint{Client: client}, privKey, 100_000, senderFn,
	)
	if err == nil {
		t.Fatal("execute succeeded; want the nonce-gap failure")
	}

	sent := node.sentTransactions()
	if len(sent) != 2 {
		t.Fatalf("broadcast %d transactions, want the payload and its cancellation", len(sent))
	}
	payload, cancel := sent[0], sent[1]
	if payload.Nonce() != abandonedNonce {
		t.Fatalf("payload nonce = %d, want the cached %d", payload.Nonce(), abandonedNonce)
	}
	if cancel.Nonce() != abandonedNonce {
		t.Fatalf("cancellation nonce = %d, want the abandoned %d: anywhere else and the payload still executes",
			cancel.Nonce(), abandonedNonce)
	}
	if len(cancel.Data()) != 0 || cancel.Value().Sign() != 0 {
		t.Fatalf("the transaction left at nonce %d still carries %d bytes of calldata and value %s; it is the payload, not a no-op",
			abandonedNonce, len(cancel.Data()), cancel.Value())
	}
	if cancel.To() == nil || *cancel.To() != from {
		t.Errorf("cancellation to = %v, want a self-transfer to %s", cancel.To(), from.Hex())
	}
}
