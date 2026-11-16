package transaction

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"net/http/httptest"
	"strings"
	"testing"

	services "relayer/services"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func gwei(n int64) *big.Int { return new(big.Int).Mul(big.NewInt(n), big.NewInt(1_000_000_000)) }

// setFeeCaps sets both caps for one test, in gwei. A negative value leaves that
// variable unset.
func setFeeCaps(t *testing.T, maxFeeGwei, maxTipGwei int64) {
	t.Helper()
	set := func(name string, value int64) {
		if value < 0 {
			t.Setenv(name, "")
			return
		}
		t.Setenv(name, gwei(value).String())
	}
	set(EnvEVMMaxFeePerGas, maxFeeGwei)
	set(EnvEVMMaxPriorityFeePerGas, maxTipGwei)
}

// feeCapFixture drives the real send path against a node that quotes whatever
// fees the test tells it to and mines every transaction it accepts.
type feeCapFixture struct {
	node   *mockJSONRPC
	client *ethclient.Client
	key    *ecdsa.PrivateKey
	from   common.Address
}

func newFeeCapFixture(t *testing.T) *feeCapFixture {
	t.Helper()
	fastReceipts(t)
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
	return &feeCapFixture{node: node, client: client, key: key, from: crypto.PubkeyToAddress(key.PublicKey)}
}

// send runs one transaction through executeWithRetryAndResubmission. The stand-in
// for a bound contract call prices the transaction from auth exactly as the
// bindings do, signs it through auth.Signer and broadcasts it.
func (f *feeCapFixture) send(t *testing.T, h *Handler) error {
	t.Helper()
	senderFn := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		var tx *types.Transaction
		if auth.GasPrice != nil {
			tx = types.NewTx(&types.LegacyTx{
				Nonce: auth.Nonce.Uint64(), GasPrice: auth.GasPrice, Gas: auth.GasLimit, To: &common.Address{0x1},
			})
		} else {
			tx = types.NewTx(&types.DynamicFeeTx{
				ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: auth.GasTipCap,
				GasFeeCap: auth.GasFeeCap, Gas: auth.GasLimit, To: &common.Address{0x1},
			})
		}
		signed, err := auth.Signer(auth.From, tx)
		if err != nil {
			return nil, err
		}
		f.node.mu.Lock()
		f.node.receiptResps[signed.Hash()] = &types.Receipt{Status: 1, GasUsed: 21_000, BlockNumber: big.NewInt(100)}
		f.node.mu.Unlock()
		if err := f.client.SendTransaction(auth.Context, signed); err != nil {
			return nil, err
		}
		return signed, nil
	}
	_, _, _, err := h.executeWithRetryAndResubmission(testCtx(t), services.EVMEndpoint{Client: f.client}, f.key, 100_000, senderFn)
	return err
}

func (f *feeCapFixture) nonceQueries() int {
	f.node.mu.Lock()
	defer f.node.mu.Unlock()
	return f.node.nonceCalls[f.from]
}

// The audit's case: the RPC chooses the price. With caps set, what it quotes no
// longer decides what the relayer signs (#575 M3).
func TestFeeCaps_HostileRPCCannotRaiseTheFee(t *testing.T) {
	t.Run("a hostile tip is held at the tip cap", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.baseFee = gwei(10)
		f.node.tipPrice = gwei(1_000_000) // 0.001 ETH per gas

		if err := f.send(t, &Handler{}); err != nil {
			t.Fatalf("send: %v", err)
		}
		sent := f.node.sentTransactions()
		if len(sent) != 1 {
			t.Fatalf("broadcast %d transactions, want 1", len(sent))
		}
		if got := sent[0].GasTipCap(); got.Cmp(gwei(5)) != 0 {
			t.Errorf("tip = %s wei, want the cap %s", got, gwei(5))
		}
		// 2 x base fee + tip = 25 gwei, under the fee cap, so it is used as is.
		if got := sent[0].GasFeeCap(); got.Cmp(gwei(25)) != 0 {
			t.Errorf("fee cap = %s wei, want %s", got, gwei(25))
		}
	})

	t.Run("the fee cap bounds the base-fee headroom", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.baseFee = gwei(60) // 2 x 60 + tip is above the cap; 60 + tip is not
		f.node.tipPrice = gwei(2)

		if err := f.send(t, &Handler{}); err != nil {
			t.Fatalf("send: %v", err)
		}
		sent := f.node.sentTransactions()
		if len(sent) != 1 {
			t.Fatalf("broadcast %d transactions, want 1: the cap still covers the current base fee", len(sent))
		}
		if got := sent[0].GasFeeCap(); got.Cmp(gwei(100)) != 0 {
			t.Errorf("fee cap = %s wei, want the cap %s", got, gwei(100))
		}
	})

	t.Run("a base fee above the cap is not signed", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.baseFee = gwei(900_000)
		f.node.tipPrice = gwei(2)

		err := f.send(t, &Handler{})
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		// Not permanent: a permanent failure makes the module drop the packets.
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("err = %v must not be permanent; the send works once fees fall", err)
		}
		if sent := f.node.sentTransactions(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
		}
		if n := f.nonceQueries(); n != 0 {
			t.Fatalf("queried the nonce %d time(s); a refused send must not touch the nonce cache", n)
		}
		if !strings.Contains(err.Error(), EnvEVMMaxFeePerGas) {
			t.Fatalf("err = %v, want it to name the variable to change", err)
		}
	})

	t.Run("without caps the RPC still decides", func(t *testing.T) {
		setFeeCaps(t, -1, -1)
		f := newFeeCapFixture(t)
		f.node.baseFee = gwei(10)
		f.node.tipPrice = gwei(1_000_000)

		if err := f.send(t, &Handler{}); err != nil {
			t.Fatalf("send: %v", err)
		}
		if got := f.node.sentTransactions()[0].GasTipCap(); got.Cmp(gwei(1_000_000)) != 0 {
			t.Errorf("tip = %s wei, want the suggested %s: an unset cap is no cap", got, gwei(1_000_000))
		}
	})

	t.Run("an explicit 0 is no cap", func(t *testing.T) {
		setFeeCaps(t, 0, 0)
		f := newFeeCapFixture(t)
		f.node.baseFee = gwei(10)
		f.node.tipPrice = gwei(1_000_000)

		if err := f.send(t, &Handler{}); err != nil {
			t.Fatalf("send: %v", err)
		}
		if got := f.node.sentTransactions()[0].GasTipCap(); got.Cmp(gwei(1_000_000)) != 0 {
			t.Errorf("tip = %s wei, want the suggested %s", got, gwei(1_000_000))
		}
	})
}

func TestFeeCaps_LegacyPricedChain(t *testing.T) {
	t.Run("a market price above the cap is not signed", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.noTipCap = true
		f.node.gasPrice = gwei(900_000)

		h := &Handler{}
		err := f.send(t, h)
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		if sent := f.node.sentTransactions(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
		}
		// The nonce was read, and nothing was signed at it: it is still the next one.
		if !nonceCacheIsValid(h, "1", f.key) {
			t.Fatal("a refused send invalidated the cached nonce")
		}
	})

	t.Run("a market price under the cap is sent", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.noTipCap = true
		f.node.gasPrice = gwei(20)

		if err := f.send(t, &Handler{}); err != nil {
			t.Fatalf("send: %v", err)
		}
		sent := f.node.sentTransactions()
		if len(sent) != 1 || sent[0].Type() != types.LegacyTxType {
			t.Fatalf("sent = %v, want one legacy transaction", sent)
		}
		if sent[0].GasPrice().Cmp(gwei(100)) > 0 {
			t.Errorf("gas price %s wei is above the cap", sent[0].GasPrice())
		}
	})
}

// A stuck transaction is outbid by 15% or by a fresh suggestion, whichever is
// higher. Both come from numbers the RPC supplied, so the replacement is where an
// uncapped fee would grow without limit.
func TestFeeCaps_ReplacementsStayUnderTheCaps(t *testing.T) {
	newAuth := func(t *testing.T, f *feeCapFixture) *bind.TransactOpts {
		t.Helper()
		auth, err := bind.NewKeyedTransactorWithChainID(f.key, big.NewInt(1))
		if err != nil {
			t.Fatalf("transactor: %v", err)
		}
		return auth
	}
	stuck := func(tip, fee *big.Int) *types.Transaction {
		return types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: 3, GasTipCap: tip, GasFeeCap: fee, Gas: 21_000, To: &common.Address{0x1},
		})
	}

	t.Run("a bump is held at the caps", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.gasPrice = gwei(900_000) // what the node now suggests for both price and tip

		bumped, err := (&Handler{}).bumpGasAndResubmit(context.Background(), services.EVMEndpoint{Client: f.client},
			stuck(gwei(2), gwei(40)), newAuth(t, f), 1)
		if err != nil {
			t.Fatalf("bump: %v", err)
		}
		if bumped.GasFeeCap().Cmp(gwei(100)) != 0 || bumped.GasTipCap().Cmp(gwei(5)) != 0 {
			t.Fatalf("bumped to feeCap=%s tip=%s wei, want the caps %s and %s",
				bumped.GasFeeCap(), bumped.GasTipCap(), gwei(100), gwei(5))
		}
	})

	t.Run("a bump under the caps is unchanged", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)

		bumped, err := (&Handler{}).bumpGasAndResubmit(context.Background(), services.EVMEndpoint{Client: f.client},
			stuck(gwei(2), gwei(40)), newAuth(t, f), 1)
		if err != nil {
			t.Fatalf("bump: %v", err)
		}
		// 15% over 40 and 2 gwei.
		if bumped.GasFeeCap().Cmp(gwei(46)) != 0 || bumped.GasTipCap().Cmp(big.NewInt(2_300_000_000)) != 0 {
			t.Fatalf("bumped to feeCap=%s tip=%s wei, want 46 and 2.3 gwei", bumped.GasFeeCap(), bumped.GasTipCap())
		}
	})

	t.Run("a legacy bump is held at the fee cap", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.gasPrice = gwei(900_000)
		legacy := types.NewTx(&types.LegacyTx{Nonce: 3, GasPrice: gwei(40), Gas: 21_000, To: &common.Address{0x1}})

		bumped, err := (&Handler{}).bumpGasAndResubmit(context.Background(), services.EVMEndpoint{Client: f.client},
			legacy, newAuth(t, f), 1)
		if err != nil {
			t.Fatalf("bump: %v", err)
		}
		if bumped.GasPrice().Cmp(gwei(100)) != 0 {
			t.Fatalf("bumped to %s wei, want the cap %s", bumped.GasPrice(), gwei(100))
		}
	})

	// The cancellation is signed outside auth.Signer, so it has its own gate.
	t.Run("a cancellation is held at the caps", func(t *testing.T) {
		setFeeCaps(t, 100, 5)
		f := newFeeCapFixture(t)
		f.node.gasPrice = gwei(900_000)

		if _, err := (&Handler{}).cancelAbandonedTx(context.Background(), services.EVMEndpoint{Client: f.client},
			f.key, stuck(gwei(2), gwei(40)), f.from, big.NewInt(1)); err != nil {
			t.Fatalf("cancel: %v", err)
		}
		sent := f.node.sentTransactions()
		if len(sent) != 1 {
			t.Fatalf("broadcast %d transactions, want the cancellation", len(sent))
		}
		if sent[0].GasFeeCap().Cmp(gwei(100)) != 0 || sent[0].GasTipCap().Cmp(gwei(5)) != 0 {
			t.Fatalf("cancellation at feeCap=%s tip=%s wei, want the caps", sent[0].GasFeeCap(), sent[0].GasTipCap())
		}
	})
}

// The gate in front of the signature is what holds if a future path prices a
// transaction without going through the caps.
func TestFeeCaps_SignerRefusesATransactionAboveTheCaps(t *testing.T) {
	setFeeCaps(t, 100, 5)
	f := newFeeCapFixture(t)
	f.node.baseFee = gwei(10)
	f.node.tipPrice = gwei(1)

	overpriced := func(auth *bind.TransactOpts) (*types.Transaction, error) {
		tx := types.NewTx(&types.DynamicFeeTx{
			ChainID: big.NewInt(1), Nonce: auth.Nonce.Uint64(), GasTipCap: gwei(1),
			GasFeeCap: gwei(101), Gas: auth.GasLimit, To: &common.Address{0x1},
		})
		return auth.Signer(auth.From, tx)
	}
	_, _, _, err := (&Handler{}).executeWithRetryAndResubmission(testCtx(t), services.EVMEndpoint{Client: f.client}, f.key, 100_000, overpriced)
	if !errors.Is(err, services.ErrSpendCapExceeded) {
		t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
	}
}

func TestEVMFeeCaps_DynamicFees(t *testing.T) {
	caps := evmFeeCaps{maxFee: gwei(100), maxTip: gwei(5)}
	cases := []struct {
		name             string
		baseFee, tip     *big.Int
		wantTip, wantFee *big.Int
		refused          bool
	}{
		{name: "under both caps", baseFee: gwei(10), tip: gwei(2), wantTip: gwei(2), wantFee: gwei(22)},
		{name: "tip held at its cap", baseFee: gwei(10), tip: gwei(50), wantTip: gwei(5), wantFee: gwei(25)},
		{name: "headroom held at the fee cap", baseFee: gwei(60), tip: gwei(2), wantTip: gwei(2), wantFee: gwei(100)},
		{name: "exactly coverable", baseFee: gwei(95), tip: gwei(5), wantTip: gwei(5), wantFee: gwei(100)},
		{name: "base fee plus tip above the cap", baseFee: gwei(96), tip: gwei(5), refused: true},
		{name: "base fee alone above the cap", baseFee: gwei(101), tip: big.NewInt(0), refused: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tip, fee, err := caps.dynamicFees(tc.baseFee, tc.tip)
			if tc.refused {
				if !errors.Is(err, services.ErrSpendCapExceeded) {
					t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected refusal: %v", err)
			}
			if tip.Cmp(tc.wantTip) != 0 || fee.Cmp(tc.wantFee) != 0 {
				t.Fatalf("tip=%s fee=%s, want tip=%s fee=%s", tip, fee, tc.wantTip, tc.wantFee)
			}
		})
	}

	t.Run("no caps is the old formula", func(t *testing.T) {
		tip, fee, err := evmFeeCaps{}.dynamicFees(gwei(500), gwei(300))
		if err != nil || tip.Cmp(gwei(300)) != 0 || fee.Cmp(gwei(1300)) != 0 {
			t.Fatalf("tip=%s fee=%s err=%v, want 300 and 1300 gwei", tip, fee, err)
		}
	})

	// A tip cap above the fee cap is a config mistake, not a reason to sign an
	// invalid transaction.
	t.Run("clamp never leaves the tip above the fee cap", func(t *testing.T) {
		tip, fee := evmFeeCaps{maxFee: gwei(10), maxTip: gwei(50)}.clampDynamic(gwei(40), gwei(80))
		if tip.Cmp(fee) > 0 || fee.Cmp(gwei(10)) != 0 {
			t.Fatalf("tip=%s fee=%s, want tip <= fee = 10 gwei", tip, fee)
		}
	})
}

func TestEVMFeeCaps_Environment(t *testing.T) {
	t.Run("start needs both, and names what is missing", func(t *testing.T) {
		setFeeCaps(t, -1, -1)
		err := RequireEVMFeeCaps()
		if err == nil || !strings.Contains(err.Error(), EnvEVMMaxFeePerGas) || !strings.Contains(err.Error(), EnvEVMMaxPriorityFeePerGas) {
			t.Fatalf("err = %v, want both variables named", err)
		}

		setFeeCaps(t, 100, -1)
		err = RequireEVMFeeCaps()
		if err == nil || !strings.Contains(err.Error(), EnvEVMMaxPriorityFeePerGas) || strings.Contains(err.Error(), EnvEVMMaxFeePerGas+" and") {
			t.Fatalf("err = %v, want only the missing variable named", err)
		}
	})

	t.Run("0 is an accepted decision", func(t *testing.T) {
		setFeeCaps(t, 0, 0)
		if err := RequireEVMFeeCaps(); err != nil {
			t.Fatalf("explicit 0 refused: %v", err)
		}
		caps, err := loadEVMFeeCaps()
		if err != nil || caps.maxFee != nil || caps.maxTip != nil {
			t.Fatalf("caps = %+v err = %v, want no caps", caps, err)
		}
	})

	for _, bad := range []string{"-1", "1.5", "100gwei", "0x10", "1e9"} {
		t.Run("rejects "+bad, func(t *testing.T) {
			t.Setenv(EnvEVMMaxFeePerGas, bad)
			t.Setenv(EnvEVMMaxPriorityFeePerGas, "1")
			if err := RequireEVMFeeCaps(); err == nil {
				t.Fatalf("%q accepted as a cap", bad)
			}
			if _, err := loadEVMFeeCaps(); err == nil {
				t.Fatalf("%q accepted as a cap by the send path", bad)
			}
		})
	}

	// A malformed cap must stop the send, not silently mean "no cap".
	t.Run("a malformed cap stops the send", func(t *testing.T) {
		t.Setenv(EnvEVMMaxFeePerGas, "100gwei")
		t.Setenv(EnvEVMMaxPriorityFeePerGas, "1")
		f := newFeeCapFixture(t)
		if err := f.send(t, &Handler{}); err == nil {
			t.Fatal("send went ahead with an unreadable cap")
		}
		if sent := f.node.sentTransactions(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
		}
	})
}
