package transaction

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"relayer/rpcmock"
	services "relayer/services"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// revertDataHex is revert data for a custom error with no arguments.
func revertDataHex(signature string) string {
	selector := errorSelector(signature)
	return hexutil.Encode(selector[:])
}

// A packet transaction that reverts used to be discovered from its receipt, after
// its gas was paid, and the flush backstop offered the same packet again every
// pass (#575 M2). The eth_call before the broadcast is what makes a poison packet
// cost nothing.
func TestPreflightRefusesARevertingPacket(t *testing.T) {
	t.Run("single packet", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 0, GasUsed: 16_000_000, BlockNumber: 100})
		node.FailCallWith(3, "execution reverted", revertDataHex("AppCallbackReverted()"))
		h, _ := testHandler(t)

		err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg())
		if !errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("err = %v, want a permanent relay failure so the module drops the packet", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s) for a packet the simulation showed reverting, want 0", len(sent))
		}
	})

	t.Run("multicall", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 0, GasUsed: 16_000_000, BlockNumber: 100})
		node.FailCallWith(3, "execution reverted", revertDataHex("AppCallbackReverted()"))
		h, _ := testHandler(t)

		msgs := []any{recvPacketMsg(), recvPacketMsg()}
		err := h.SendEthTxBatch(testCtx(t), testEndpoint(t, node), "client-0", msgs)
		if !errors.Is(err, services.ErrPermanentRelayFailure) {
			t.Fatalf("err = %v, want a permanent relay failure so the module isolates the packets", err)
		}
		if !strings.Contains(err.Error(), "multicall labels=") {
			t.Fatalf("err = %v, want it to name the multicall's inner calls", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 0 {
			t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
		}
	})
}

// The receipt path keeps two reverts apart from a poison packet, because neither
// says anything about the packet. The simulation must classify them the same way,
// or a frozen client would start dropping packets.
func TestPreflightKeepsTheReceiptPathErrorClasses(t *testing.T) {
	cases := []struct {
		name string
		data string
		want error
	}{
		{"validator cache race", revertDataHex("ValidatorSetCacheMiss(bytes32)") + strings.Repeat("00", 32), services.ErrValidatorCacheRace},
		{"frozen client", revertDataHex("FrozenClientState()"), services.ErrClientNotActive},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			node := rpcmock.NewEVMNode(t)
			node.FailCallWith(3, "execution reverted", tc.data)
			h, _ := testHandler(t)

			err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg())
			if !errors.Is(err, tc.want) {
				t.Fatalf("err = %v, want %v", err, tc.want)
			}
			if errors.Is(err, services.ErrPermanentRelayFailure) {
				t.Fatalf("err = %v must not be permanent: the packet was not shown to be the cause", err)
			}
			if sent := node.SentRawTxs(); len(sent) != 0 {
				t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
			}
		})
	}
}

// Anything short of a revert is not evidence about the packet. Refusing on it
// would turn an RPC hiccup into a dropped packet, so the send goes ahead exactly
// as it did before the simulation existed.
func TestPreflightSendsWhenTheSimulationIsInconclusive(t *testing.T) {
	t.Run("simulation passes", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 21_000, BlockNumber: 100})
		h, _ := testHandler(t)

		if err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("send: %v", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 1 {
			t.Fatalf("broadcast %d transaction(s), want 1", len(sent))
		}
	})

	t.Run("transport error", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 21_000, BlockNumber: 100})
		node.FailCall(errors.New("connection reset by peer"))
		h, _ := testHandler(t)

		if err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("send: %v", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 1 {
			t.Fatalf("broadcast %d transaction(s), want 1: a transport error is not a revert", len(sent))
		}
	})

	t.Run("bare revert that passes with more gas", func(t *testing.T) {
		// A nested out-of-gas reverts with no data at the transaction's own gas
		// limit and passes with more. That is the gas ladder's case, not a poison
		// packet.
		node := rpcmock.NewEVMNode(t)
		node.SetBlockNumber(50)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 21_000, BlockNumber: 51})
		node.FailCallAt(50, 1, 3, "execution reverted")
		h, _ := testHandler(t)

		if err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("send: %v", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 1 {
			t.Fatalf("broadcast %d transaction(s), want 1", len(sent))
		}
	})

	t.Run("node behind this relayer's last mined block", func(t *testing.T) {
		// A lagging replica has not seen the client update this process just
		// mined, so it reports the packet's proof height as missing. Its revert
		// describes its own staleness.
		node := rpcmock.NewEVMNode(t)
		node.SetBlockNumber(10)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 21_000, BlockNumber: 100})
		h, _ := testHandler(t)
		endpoint := testEndpoint(t, node)

		if err := h.SendEthTx(testCtx(t), endpoint, "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("first send: %v", err)
		}
		node.FailCallWith(3, "execution reverted", revertDataHex("ConsensusStateNotFound()"))
		if err := h.SendEthTx(testCtx(t), endpoint, "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("second send: %v", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 2 {
			t.Fatalf("broadcast %d transaction(s), want 2", len(sent))
		}
	})
}

// A bare revert that still reverts with more gas is a require() or a revert(),
// not a gas shortage: more gas cannot fix it, so it is refused like any revert.
func TestPreflightRefusesABareRevertThatMoreGasDoesNotFix(t *testing.T) {
	node := rpcmock.NewEVMNode(t)
	node.SetBlockNumber(50)
	node.MineEverything(rpcmock.RawReceipt{Status: 0, GasUsed: 3_000_000, BlockNumber: 51})
	node.FailCallAt(50, 2, 3, "execution reverted")
	h, _ := testHandler(t)

	err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg())
	if !errors.Is(err, services.ErrPermanentRelayFailure) {
		t.Fatalf("err = %v, want a permanent relay failure", err)
	}
	if sent := node.SentRawTxs(); len(sent) != 0 {
		t.Fatalf("broadcast %d transaction(s), want 0", len(sent))
	}
}

// A client update is not simulated: it is not what M2 is about, and its failure
// is already held by the module rather than retried every pass.
func TestPacketCallInputOnlyCoversPackets(t *testing.T) {
	if _, _, ok, err := packetCallInput(services.CosmosClientUpdateBuildResult{}); ok || err != nil {
		t.Fatalf("client update: ok=%v err=%v, want it left out of the simulation", ok, err)
	}
	input, label, ok, err := packetCallInput(recvPacketMsg())
	if err != nil || !ok {
		t.Fatalf("recvPacket: ok=%v err=%v", ok, err)
	}
	if !strings.HasPrefix(hexutil.Encode(input), "0x"+routerSelector(t, "recvPacket")) {
		t.Fatalf("input does not start with the recvPacket selector: %x", input[:4])
	}
	if !strings.HasPrefix(label, "recvPacket seq=") {
		t.Fatalf("label = %q", label)
	}
}

func TestPreflightProbeGas(t *testing.T) {
	cases := []struct {
		name            string
		gas, blockLimit uint64
		want            uint64
	}{
		{"multiplied", 1_000_000, 0, 8_000_000},
		{"clamped to the block limit", 16_000_000, 30_000_000, 30_000_000},
		{"overflow saturates", ^uint64(0) / 2, 0, ^uint64(0)},
		{"already at the block limit", 30_000_000, 30_000_000, 30_000_000},
	}
	for _, tc := range cases {
		if got := preflightProbeGas(tc.gas, tc.blockLimit); got != tc.want {
			t.Errorf("%s: preflightProbeGas(%d, %d) = %d, want %d", tc.name, tc.gas, tc.blockLimit, got, tc.want)
		}
	}
}

// L12: blockNumber is optional in a receipt. The line that logged a confirmed
// transaction dereferenced it, so a node omitting it crashed the relayer after
// the transaction had already succeeded.
func TestSendEthTx_ReceiptWithoutABlockNumber(t *testing.T) {
	t.Run("single message", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 21_000, OmitBlockNumber: true})
		h, _ := testHandler(t)

		if err := h.SendEthTx(testCtx(t), testEndpoint(t, node), "client-0", recvPacketMsg()); err != nil {
			t.Fatalf("send: %v", err)
		}
		if sent := node.SentRawTxs(); len(sent) != 1 {
			t.Fatalf("broadcast %d transaction(s), want 1", len(sent))
		}
	})

	t.Run("multicall", func(t *testing.T) {
		node := rpcmock.NewEVMNode(t)
		node.MineEverything(rpcmock.RawReceipt{Status: 1, GasUsed: 50_000, OmitBlockNumber: true})
		h, _ := testHandler(t)

		if err := h.SendEthTxBatch(testCtx(t), testEndpoint(t, node), "client-0", []any{recvPacketMsg(), recvPacketMsg()}); err != nil {
			t.Fatalf("send: %v", err)
		}
	})
}

func TestReceiptBlock(t *testing.T) {
	if got := receiptBlock(nil); got != "unknown" {
		t.Errorf("receiptBlock(nil) = %q", got)
	}
	if got := receiptBlock(&types.Receipt{}); got != "unknown" {
		t.Errorf("receipt without a block number = %q, want unknown", got)
	}
	if got := receiptBlock(&types.Receipt{BlockNumber: big.NewInt(1234)}); got != "1234" {
		t.Errorf("receiptBlock = %q, want 1234", got)
	}
}
