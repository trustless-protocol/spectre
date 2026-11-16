package transaction

import (
	"context"
	"errors"
	"log"
	"math"
	"strconv"
	"strings"
	"testing"

	"relayer/services"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// setCosmosCaps sets both caps for one test. A negative value leaves that
// variable unset.
func setCosmosCaps(t *testing.T, maxGas, maxFee int64) {
	t.Helper()
	set := func(name string, value int64) {
		if value < 0 {
			t.Setenv(name, "")
			return
		}
		t.Setenv(name, strconv.FormatInt(value, 10))
	}
	set(EnvCosmosMaxGas, maxGas)
	set(EnvCosmosMaxFee, maxFee)
}

// sendOne submits a single-message, splittable batch at sequence 10.
func sendOne(t *testing.T, mock *cosmosOOGRPC) (uint64, int, error) {
	t.Helper()
	return (&Handler{}).sendCosmosTxBatchWithSplitting(
		context.Background(), newCosmosOOGEndpoint(t, mock), []sdk.Msg{cosmosOOGTestMessage()}, 7, 10, true)
}

// The audit's case: the gas limit, and with it the fee, is whatever the RPC's
// Simulate reports. The mock reports no block MaxGas, which is the CometBFT
// default (-1), so before the cap nothing bounded it (#575 M4).
func TestCosmosSpendCaps_SimulateCannotRaiseTheGas(t *testing.T) {
	t.Run("gas the node says is above the cap is not signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 50_000_000)}

		nextSequence, succeeded, err := sendOne(t, mock)
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		// Not permanent: a permanent failure makes the module drop the packet.
		if errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("err = %v must not be permanent", err)
		}
		if mock.broadcasts != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", mock.broadcasts)
		}
		if nextSequence != 10 || succeeded != 0 {
			t.Fatalf("result = sequence:%d succeeded:%d, want 10:0: nothing was signed", nextSequence, succeeded)
		}
		if !strings.Contains(err.Error(), EnvCosmosMaxGas) {
			t.Fatalf("err = %v, want it to name the variable to change", err)
		}
	})

	t.Run("headroom over the cap is held at the cap", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		// 18M simulated: 1.3x is 23.4M, above the cap, but the cap still covers it.
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 18_000_000)}

		if _, _, err := sendOne(t, mock); err != nil {
			t.Fatalf("send: %v", err)
		}
		if len(mock.gasLimits) != 1 || mock.gasLimits[0] != 20_000_000 {
			t.Fatalf("broadcast gas limits = %v, want one transaction at the cap", mock.gasLimits)
		}
		if mock.feeAmounts[0] != 20_000_000 {
			t.Fatalf("fee = %d, want it to follow the capped gas limit", mock.feeAmounts[0])
		}
	})

	t.Run("gas under the cap is unchanged", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 1_000_000)}

		if _, _, err := sendOne(t, mock); err != nil {
			t.Fatalf("send: %v", err)
		}
		if len(mock.gasLimits) != 1 || mock.gasLimits[0] != 1_300_000 {
			t.Fatalf("broadcast gas limits = %v, want 1.3x the simulated gas", mock.gasLimits)
		}
	})

	t.Run("without a cap the node still decides", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, -1, -1)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 50_000_000)}

		if _, _, err := sendOne(t, mock); err != nil {
			t.Fatalf("send: %v", err)
		}
		if len(mock.gasLimits) != 1 || mock.gasLimits[0] != 65_000_000 {
			t.Fatalf("broadcast gas limits = %v, want 1.3x the simulated gas: an unset cap is no cap", mock.gasLimits)
		}
	})

	// A splittable batch over the cap is halved, exactly as one over the block
	// limit is. Each half is simulated again on its own.
	t.Run("a batch over the cap is split", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 15_000_000}

		nextSequence, succeeded, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, true)
		if err != nil {
			t.Fatalf("send: %v", err)
		}
		if nextSequence != 12 || succeeded != 2 {
			t.Fatalf("result = sequence:%d succeeded:%d, want 12:2", nextSequence, succeeded)
		}
		if len(mock.gasLimits) != 2 || mock.gasLimits[0] != 19_500_000 || mock.gasLimits[1] != 19_500_000 {
			t.Fatalf("broadcast gas limits = %v, want two transactions at 1.3x 15M", mock.gasLimits)
		}
	})

	t.Run("an atomic batch over the cap is not signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 15_000_000}

		_, _, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, false)
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		if mock.broadcasts != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", mock.broadcasts)
		}
		// It names what the node said, not that plus the headroom (39000000).
		if !strings.Contains(err.Error(), "simulated gas 30000000 ") {
			t.Fatalf("err = %v, want it to name the 30000000 simulated gas", err)
		}
	})

	// #575 M4: an atomic batch -- a client update and the packets proven against
	// it -- was refused on its headroom: 16M simulated is 20.8M with it, over a
	// 20M cap the batch fits under. One message in the same spot is signed at the
	// cap; so is the batch now.
	t.Run("an atomic batch under the cap is signed at the cap", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 8_000_000}

		nextSequence, succeeded, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, false)
		if err != nil {
			t.Fatalf("send: %v", err)
		}
		if nextSequence != 11 || succeeded != 2 {
			t.Fatalf("result = sequence:%d succeeded:%d, want 11:2: one transaction", nextSequence, succeeded)
		}
		if len(mock.gasLimits) != 1 || mock.gasLimits[0] != 20_000_000 {
			t.Fatalf("broadcast gas limits = %v, want one transaction at the cap", mock.gasLimits)
		}
	})

	// At the cap there is no more gas to give. Retrying with the identical limit
	// is the loop the ladder exists to prevent.
	t.Run("out of gas at the cap ends the ladder", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 18_000_000), deliverTxOOG: 100}

		_, _, err := sendOne(t, mock)
		if !errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("err = %v, want the ladder to report it exhausted", err)
		}
		if len(mock.gasLimits) != 1 {
			t.Fatalf("broadcast gas limits = %v, want one attempt: a second at the same cap cannot differ", mock.gasLimits)
		}
	})
}

// L13: the default fee is the gas limit as an int64 coin amount. A Simulate
// answer of 2^63 or more became a negative amount, and sdk.NewCoin panics on one.
func TestCosmosSend_AbsurdSimulatedGasDoesNotPanic(t *testing.T) {
	for _, gas := range []uint64{1 << 63, math.MaxUint64, math.MaxInt64} {
		t.Run(strconv.FormatUint(gas, 10), func(t *testing.T) {
			configureCosmosOOGTest(t)
			setCosmosCaps(t, -1, -1)
			mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, gas)}

			nextSequence, succeeded, err := sendOne(t, mock)
			if err == nil {
				t.Fatal("send succeeded with a gas limit no fee amount can hold")
			}
			if errors.Is(err, services.ErrPermanentRelayFailure) {
				t.Fatalf("err = %v must not be permanent: it says the RPC is wrong, not the packet", err)
			}
			if mock.broadcasts != 0 || nextSequence != 10 || succeeded != 0 {
				t.Fatalf("broadcasts=%d sequence=%d succeeded=%d, want nothing signed", mock.broadcasts, nextSequence, succeeded)
			}
		})
	}

	t.Run("with a cap the same answer is a cap refusal", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 20_000_000, 0)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 1<<63)}

		if _, _, err := sendOne(t, mock); !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
	})
}

func TestCosmosSpendCaps_Fee(t *testing.T) {
	t.Run("a fee above the cap is not signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 0, 1_000_000)
		// Default fee = gas limit = 1.3M, above the 1M fee cap.
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 1_000_000)}

		nextSequence, succeeded, err := sendOne(t, mock)
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		if mock.broadcasts != 0 || nextSequence != 10 || succeeded != 0 {
			t.Fatalf("broadcasts=%d sequence=%d succeeded=%d, want nothing signed", mock.broadcasts, nextSequence, succeeded)
		}
		if !strings.Contains(err.Error(), EnvCosmosMaxFee) {
			t.Fatalf("err = %v, want it to name the variable to change", err)
		}
	})

	t.Run("a fee at the cap is signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 0, 1_300_000)
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 1_000_000)}

		if _, _, err := sendOne(t, mock); err != nil {
			t.Fatalf("send: %v", err)
		}
		if len(mock.feeAmounts) != 1 || mock.feeAmounts[0] != 1_300_000 {
			t.Fatalf("fees = %v, want 1300000", mock.feeAmounts)
		}
	})

	// COSMOS_FEE_AMOUNT is per message, so a batch's fee grows with its size and
	// halving the batch halves the fee.
	t.Run("a batch whose fee is above the cap is split", func(t *testing.T) {
		configureCosmosOOGTest(t)
		t.Setenv("COSMOS_FEE_AMOUNT", "1000")
		setCosmosCaps(t, 0, 1500)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 100_000}

		nextSequence, succeeded, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, true)
		if err != nil {
			t.Fatalf("send: %v", err)
		}
		if nextSequence != 12 || succeeded != 2 {
			t.Fatalf("result = sequence:%d succeeded:%d, want 12:2", nextSequence, succeeded)
		}
		if len(mock.feeAmounts) != 2 || mock.feeAmounts[0] != 1000 || mock.feeAmounts[1] != 1000 {
			t.Fatalf("fees = %v, want two transactions at 1000", mock.feeAmounts)
		}
	})

	t.Run("an atomic batch whose fee is above the cap is not signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		t.Setenv("COSMOS_FEE_AMOUNT", "1000")
		setCosmosCaps(t, 0, 1500)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 100_000}

		_, _, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, false)
		if !errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
		}
		if mock.broadcasts != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", mock.broadcasts)
		}
	})
}

// captureStdLog returns what fn wrote to the standard logger. It reads the
// buffer only after the logger is restored, so a late write cannot race it.
func captureStdLog(t *testing.T, fn func()) (out string) {
	t.Helper()
	var buf strings.Builder
	prev, flags := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer func() {
		log.SetOutput(prev)
		log.SetFlags(flags)
		out = buf.String()
	}()
	fn()
	return ""
}

// A cap refusing one message is the operator's to act on: an ATTENTION error.
// An atomic batch is only ever a folded relay, which the relay module then sends
// apart, so its refusal is a warning; a part refused again says ATTENTION
// itself. One ATTENTION line per batch on a path that works dilutes the label
// (#575 M4).
func TestCosmosCapRefusal_LogLevel(t *testing.T) {
	cases := []struct {
		name           string
		maxGas, maxFee int64
		feePerMsg      string // COSMOS_FEE_AMOUNT; empty keeps the gas-matching fee
		msgs           int
		simPerMsg      uint64
		want           string
	}{
		{"one message over the gas cap", 20_000_000, 0, "", 1, 30_000_000, "ERROR [SendCosmosTx][ATTENTION] simulated gas"},
		{"atomic batch over the gas cap", 20_000_000, 0, "", 2, 15_000_000, "WARN [SendCosmosTx] atomic batch of 2: simulated gas"},
		{"one message over the fee cap", 0, 1_000_000, "", 1, 1_000_000, "ERROR [SendCosmosTx][ATTENTION] fee"},
		{"atomic batch over the fee cap", 0, 1500, "1000", 2, 100_000, "WARN [SendCosmosTx] atomic batch of 2: fee"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			configureCosmosOOGTest(t)
			setCosmosCaps(t, tc.maxGas, tc.maxFee)
			if tc.feePerMsg != "" {
				t.Setenv("COSMOS_FEE_AMOUNT", tc.feePerMsg)
			}
			mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: tc.simPerMsg}
			msgs := make([]sdk.Msg, tc.msgs)
			for i := range msgs {
				msgs[i] = cosmosOOGTestMessage()
			}

			var err error
			out := captureStdLog(t, func() {
				// A single message may split; a batch here is atomic.
				_, _, err = (&Handler{}).sendCosmosTxBatchWithSplitting(
					context.Background(), newCosmosOOGEndpoint(t, mock), msgs, 7, 10, tc.msgs == 1)
			})
			if !errors.Is(err, services.ErrSpendCapExceeded) {
				t.Fatalf("err = %v, want ErrSpendCapExceeded", err)
			}
			if !strings.Contains(out, tc.want) {
				t.Fatalf("log does not contain %q:\n%s", tc.want, out)
			}
			if tc.msgs > 1 && strings.Contains(out, "ATTENTION") {
				t.Fatalf("an atomic batch the relay sends apart must not log ATTENTION:\n%s", out)
			}
		})
	}
}

// The ceiling can be the block's gas limit rather than the cap. An atomic batch
// is judged against it the same way: signed at the limit when its simulated gas
// fits, refused -- and not as a spending cap -- when it does not.
func TestCosmosAtomicBatch_BlockGasCeiling(t *testing.T) {
	send := func(t *testing.T, mock *cosmosOOGRPC) error {
		t.Helper()
		_, _, err := (&Handler{}).sendCosmosTxBatchWithSplitting(
			context.Background(), newCosmosOOGEndpoint(t, mock),
			[]sdk.Msg{cosmosOOGTestMessage(), cosmosOOGTestMessage()}, 7, 10, false)
		return err
	}

	t.Run("simulated gas under the block limit is signed at the limit", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 0, 0)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 8_000_000, blockMaxGas: 20_000_000}

		if err := send(t, mock); err != nil {
			t.Fatalf("send: %v", err)
		}
		if len(mock.gasLimits) != 1 || mock.gasLimits[0] != 20_000_000 {
			t.Fatalf("broadcast gas limits = %v, want one transaction at the block limit", mock.gasLimits)
		}
	})

	t.Run("simulated gas over the block limit is not signed", func(t *testing.T) {
		configureCosmosOOGTest(t)
		setCosmosCaps(t, 0, 0)
		mock := &cosmosOOGRPC{t: t, simulateGasPerMsg: 15_000_000, blockMaxGas: 20_000_000}

		err := send(t, mock)
		if err == nil || errors.Is(err, services.ErrSpendCapExceeded) {
			t.Fatalf("err = %v, want a refusal that is not a spending cap", err)
		}
		if mock.broadcasts != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0: no block can hold it", mock.broadcasts)
		}
	})
}

func TestCosmosSpendCaps_GasCeiling(t *testing.T) {
	cases := []struct {
		name                string
		capGas, block, want uint64
	}{
		{"no cap, no block limit", 0, 0, 0},
		{"no cap", 0, 30, 30},
		{"cap only: the block reports none", 20, 0, 20},
		{"cap tighter", 20, 30, 20},
		{"block tighter", 40, 30, 30},
	}
	for _, tc := range cases {
		if got := (cosmosSpendCaps{maxGas: tc.capGas}).gasCeiling(tc.block); got != tc.want {
			t.Errorf("%s: gasCeiling(%d) with cap %d = %d, want %d", tc.name, tc.block, tc.capGas, got, tc.want)
		}
	}
}

func TestCosmosSpendCaps_Environment(t *testing.T) {
	t.Run("start needs both, and names what is missing", func(t *testing.T) {
		setCosmosCaps(t, -1, -1)
		err := RequireCosmosSpendCaps()
		if err == nil || !strings.Contains(err.Error(), EnvCosmosMaxGas) || !strings.Contains(err.Error(), EnvCosmosMaxFee) {
			t.Fatalf("err = %v, want both variables named", err)
		}
	})

	t.Run("0 is an accepted decision", func(t *testing.T) {
		setCosmosCaps(t, 0, 0)
		if err := RequireCosmosSpendCaps(); err != nil {
			t.Fatalf("explicit 0 refused: %v", err)
		}
		caps, err := loadCosmosSpendCaps()
		if err != nil || caps != (cosmosSpendCaps{}) {
			t.Fatalf("caps = %+v err = %v, want no caps", caps, err)
		}
	})

	// 2^63 parses as a uint64 and cannot be a fee amount.
	for _, bad := range []string{"-1", "1.5", "20M", "9223372036854775808"} {
		t.Run("rejects "+bad, func(t *testing.T) {
			t.Setenv(EnvCosmosMaxGas, bad)
			t.Setenv(EnvCosmosMaxFee, "1")
			if err := RequireCosmosSpendCaps(); err == nil {
				t.Fatalf("%q accepted as a cap", bad)
			}
			if _, err := loadCosmosSpendCaps(); err == nil {
				t.Fatalf("%q accepted as a cap by the send path", bad)
			}
		})
	}

	// A malformed cap must stop the send, not silently mean "no cap".
	t.Run("a malformed cap stops the send", func(t *testing.T) {
		configureCosmosOOGTest(t)
		t.Setenv(EnvCosmosMaxGas, "20M")
		t.Setenv(EnvCosmosMaxFee, "0")
		mock := &cosmosOOGRPC{t: t, simulation: simulatedGasResponse(t, 1_000_000)}
		if _, _, err := sendOne(t, mock); err == nil {
			t.Fatal("send went ahead with an unreadable cap")
		}
		if mock.broadcasts != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", mock.broadcasts)
		}
	})
}
