package transaction

import (
	"crypto/ecdsa"
	"errors"
	"testing"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type fakeSigner struct {
	ethCalls    int
	cosmosCalls int
	err         error
}

func (s *fakeSigner) EthKey() (*ecdsa.PrivateKey, error) {
	s.ethCalls++
	return nil, s.err
}

func (s *fakeSigner) CosmosKeyBytes() ([]byte, error) {
	s.cosmosCalls++
	return nil, s.err
}

// Classifying an INCLUDED Cosmos out-of-gas as transient is only safe if the
// retry differs from the attempt that failed. It does not by itself: a transient
// failure charges no retry budget (chargesRetryBudget requires permanent), is
// never marked Solo, and is re-queued at the HEAD of the queue. So a repeatable
// simulation under-estimate re-broadcasts the identical tx forever — burning a
// sequence and fees each round and blocking every packet behind it.
//
// planCosmosOutOfGas is what makes each retry a different attempt. These tests
// pin that every path shrinks something and that the shrinking bottoms out.

func TestPlanCosmosOutOfGasSplitsABatchBeforeTouchingHeadroom(t *testing.T) {
	// Splitting is preferred while it is available: the halves are simulated
	// separately, which fixes the cause (one estimate covering too much work)
	// rather than papering over it with a bigger multiplier.
	for _, msgCount := range []int{2, 3, 16} {
		if got := planCosmosOutOfGas(msgCount, 0, 1_000_000, 30_000_000, true); got != oogSplitBatch {
			t.Fatalf("planCosmosOutOfGas(msgCount=%d) = %v, want oogSplitBatch", msgCount, got)
		}
	}
	// Even at the last headroom step, a multi-message batch still splits.
	last := len(cosmosGasHeadroom) - 1
	if got := planCosmosOutOfGas(4, last, 1_000_000, 30_000_000, true); got != oogSplitBatch {
		t.Fatalf("multi-message at the last headroom step = %v, want oogSplitBatch", got)
	}
}

func TestPlanCosmosOutOfGasEscalatesASingleMessage(t *testing.T) {
	for step := 0; step < len(cosmosGasHeadroom)-1; step++ {
		if got := planCosmosOutOfGas(1, step, 1_000_000, 30_000_000, true); got != oogEscalateHeadroom {
			t.Fatalf("planCosmosOutOfGas(1 msg, step=%d) = %v, want oogEscalateHeadroom", step, got)
		}
	}
}

// The two floors. Without them "transient" is unbounded, which is the defect.
func TestPlanCosmosOutOfGasBottomsOutAsPermanent(t *testing.T) {
	last := len(cosmosGasHeadroom) - 1
	if got := planCosmosOutOfGas(1, last, 1_000_000, 30_000_000, true); got != oogPermanent {
		t.Fatalf("single message at the last headroom step = %v, want oogPermanent — "+
			"the ladder must end, or the retry loop never does", got)
	}
	// At the block ceiling a larger factor buys nothing: the next attempt would be
	// the same transaction. That is permanent however many steps are left.
	if got := planCosmosOutOfGas(1, 0, 30_000_000, 30_000_000, true); got != oogPermanent {
		t.Fatalf("single message already at the block gas limit = %v, want oogPermanent", got)
	}
	if got := planCosmosOutOfGas(1, 0, 31_000_000, 30_000_000, true); got != oogPermanent {
		t.Fatalf("single message above the block gas limit = %v, want oogPermanent", got)
	}
}

// A chain that does not report a block max (MaxGas unset) must not be read as
// "the ceiling is 0, everything is permanent" — that would dead-letter valid
// packets on the first out-of-gas.
func TestPlanCosmosOutOfGasWithNoBlockLimitStillEscalates(t *testing.T) {
	if got := planCosmosOutOfGas(1, 0, 1_000_000, 0, true); got != oogEscalateHeadroom {
		t.Fatalf("no block max reported = %v, want oogEscalateHeadroom", got)
	}
}

// Drive the planner the way the send path does and show the sequence of actions
// terminates from every starting point. This is the property the whole change
// exists for; the individual cases above only pin its pieces.
func TestCosmosOutOfGasRetriesTerminate(t *testing.T) {
	const blockMax = uint64(30_000_000)
	for _, msgCount := range []int{1, 2, 5, 32} {
		msgs, step := msgCount, 0
		// Generous bound: any run longer than this is a loop, not a ladder.
		limit := 64
		for i := 0; ; i++ {
			if i > limit {
				t.Fatalf("no termination from msgCount=%d after %d steps", msgCount, limit)
			}
			switch planCosmosOutOfGas(msgs, step, 1_000_000, blockMax, true) {
			case oogSplitBatch:
				// Worst case for termination: the failing half is the larger one.
				msgs -= msgs / 2
				step = 0 // halves are simulated afresh
			case oogEscalateHeadroom:
				step++
			default:
				if msgs != 1 {
					t.Fatalf("terminated as permanent while %d messages could still be split", msgs)
				}
				goto done
			}
		}
	done:
	}
}

// The ladder must actually climb. A flat or descending entry would make an
// "escalation" re-send the same transaction, restoring the loop while looking
// like it had been fixed.
func TestCosmosGasHeadroomIsStrictlyIncreasing(t *testing.T) {
	if len(cosmosGasHeadroom) < 2 {
		t.Fatal("a single-entry ladder cannot escalate; the retry would repeat the failed attempt")
	}
	if cosmosGasHeadroom[0] != 1.3 {
		t.Fatalf("first step = %v, want the historical 1.3 so unaffected batches are unchanged",
			cosmosGasHeadroom[0])
	}
	for i := 1; i < len(cosmosGasHeadroom); i++ {
		if cosmosGasHeadroom[i] <= cosmosGasHeadroom[i-1] {
			t.Fatalf("step %d (%v) does not exceed step %d (%v)",
				i, cosmosGasHeadroom[i], i-1, cosmosGasHeadroom[i-1])
		}
	}
}

// The predicate comes from the SDK's registered error rather than the literals
// "sdk" and 11, so it tracks whatever the chain actually sends.
func TestCosmosOutOfGasMatchesTheSDKError(t *testing.T) {
	cs, code := errortypes.ErrOutOfGas.Codespace(), errortypes.ErrOutOfGas.ABCICode()
	if !cosmosOutOfGas(cs, code) {
		t.Fatalf("cosmosOutOfGas(%q, %d) = false, want true for the SDK's own ErrOutOfGas", cs, code)
	}
	if cosmosOutOfGas("wasm", code) {
		t.Fatal("right code in another codespace must not be read as out of gas")
	}
	if cosmosOutOfGas(cs, errortypes.ErrInsufficientFee.ABCICode()) {
		t.Fatal("another sdk code must not be read as out of gas")
	}
	if cosmosOutOfGas("", 0) {
		t.Fatal("a success must not be read as out of gas")
	}
}

// Startup validation must go THROUGH the seam.
//
// It used to read ETH_PRIVATE_KEY and COSMOS_PRIVATE_KEY itself and parse them
// with keys.RestoreKey, which left a ninth key source outside the seam that
// RLY-07 introduced -- and one a KMS-backed signer could never satisfy, since it
// does not hand back a raw private key to parse. An installed signer must be the
// only thing consulted, with the environment untouched.
func TestValidateKeysUsesTheInstalledSignerNotTheEnvironment(t *testing.T) {
	t.Setenv("ETH_PRIVATE_KEY", "")
	t.Setenv("COSMOS_PRIVATE_KEY", "")

	want := errors.New("kms unavailable")
	fake := &fakeSigner{err: want}
	h := &Handler{}
	h.SetSigner(fake)

	if err := h.ValidateKeys(); !errors.Is(err, want) {
		t.Fatalf("ValidateKeys() = %v, want the signer's own error %v; a validator that reads the "+
			"environment cannot check a KMS-backed key at all", err, want)
	}
	if fake.ethCalls == 0 {
		t.Fatal("the ETH key was never requested from the signer")
	}
}

// Both keys are checked, not just the first: a relayer that starts with a broken
// Cosmos key fails on its first send instead of at startup, which is the whole
// point of validating.
func TestValidateKeysChecksBothKeys(t *testing.T) {
	priv, err := ethcrypto.GenerateKey()
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	fake := &okEthBadCosmosSigner{eth: priv}
	h := &Handler{}
	h.SetSigner(fake)

	if err := h.ValidateKeys(); err == nil {
		t.Fatal("ValidateKeys() = nil with an unusable Cosmos key; the failure would surface on the first send")
	}
	if !fake.cosmosAsked {
		t.Fatal("the Cosmos key was never requested")
	}
}

type okEthBadCosmosSigner struct {
	eth         *ecdsa.PrivateKey
	cosmosAsked bool
}

func (s *okEthBadCosmosSigner) EthKey() (*ecdsa.PrivateKey, error) { return s.eth, nil }

func (s *okEthBadCosmosSigner) CosmosKeyBytes() ([]byte, error) {
	s.cosmosAsked = true
	return nil, errors.New("COSMOS_PRIVATE_KEY environment variable is required in .env file")
}
