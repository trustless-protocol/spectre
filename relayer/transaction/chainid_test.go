package transaction

import (
	"errors"
	"strings"
	"testing"

	services "relayer/services"
)

// The signature names the chain it is valid on, and that name came from the RPC
// on every send. An endpoint that changes its answer must not get a signature
// for the new chain (#575, "RPC chain id per send").
func TestEVMChainID_EndpointIsHeldToItsFirstAnswer(t *testing.T) {
	setFeeCaps(t, -1, -1)
	f := newFeeCapFixture(t)
	h := &Handler{}

	if err := f.send(t, h); err != nil {
		t.Fatalf("first send: %v", err)
	}

	f.node.mu.Lock()
	f.node.chainID = 8453 // the same key is funded here too
	f.node.mu.Unlock()

	err := f.send(t, h)
	if !errors.Is(err, ErrEVMChainIDChanged) {
		t.Fatalf("err = %v, want ErrEVMChainIDChanged", err)
	}
	// It says the endpoint is wrong, not the packet: nothing may be dropped on it.
	if errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("err = %v must not be permanent", err)
	}
	if !strings.Contains(err.Error(), "8453") || !strings.Contains(err.Error(), "signs for 1") {
		t.Fatalf("err = %v, want both chain ids named", err)
	}
	if sent := f.node.sentTransactions(); len(sent) != 1 {
		t.Fatalf("broadcast %d transactions, want only the first", len(sent))
	}

	// The endpoint answering correctly again is accepted: the pin is the id, not
	// a latch.
	f.node.mu.Lock()
	f.node.chainID = 1
	f.node.mu.Unlock()
	if err := f.send(t, h); err != nil {
		t.Fatalf("send after the endpoint recovered: %v", err)
	}
}

// Under `start` the first answer is not trusted on its own either: it must be a
// chain the process took a nonce-domain lock for.
func TestEVMChainID_FirstAnswerMustBeALockedChain(t *testing.T) {
	setFeeCaps(t, -1, -1)

	t.Run("a chain outside the locked set is refused", func(t *testing.T) {
		f := newFeeCapFixture(t)
		h := &Handler{}
		h.AllowEVMChainIDs([]string{"11155111"})

		err := f.send(t, h) // the node reports chain id 1
		if !errors.Is(err, ErrEVMChainIDChanged) {
			t.Fatalf("err = %v, want ErrEVMChainIDChanged", err)
		}
		if sent := f.node.sentTransactions(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
		}
		if n := f.nonceQueries(); n != 0 {
			t.Fatalf("queried the nonce %d time(s); a refused send must not touch the nonce cache", n)
		}
	})

	t.Run("a locked chain is accepted", func(t *testing.T) {
		f := newFeeCapFixture(t)
		h := &Handler{}
		h.AllowEVMChainIDs([]string{"11155111", "1"})

		if err := f.send(t, h); err != nil {
			t.Fatalf("send: %v", err)
		}
	})

	// A refused first answer must not become the pin.
	t.Run("a refused answer is not pinned", func(t *testing.T) {
		f := newFeeCapFixture(t)
		h := &Handler{}
		h.AllowEVMChainIDs([]string{"1"})

		f.node.mu.Lock()
		f.node.chainID = 8453
		f.node.mu.Unlock()
		if err := f.send(t, h); !errors.Is(err, ErrEVMChainIDChanged) {
			t.Fatalf("err = %v, want ErrEVMChainIDChanged", err)
		}

		f.node.mu.Lock()
		f.node.chainID = 1
		f.node.mu.Unlock()
		if err := f.send(t, h); err != nil {
			t.Fatalf("send on the locked chain after a refusal: %v", err)
		}
	})
}
