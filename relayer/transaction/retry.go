// This file answers one question for both chains: what must the NEXT attempt
// change?
//
// It replaces gas_revert.go and cosmos_oog.go, which were named after the
// situations that produced them -- a revert attributed to gas on EVM, an
// out-of-gas on Cosmos -- rather than after what they decide. Naming by
// situation is why they were two files: every new situation would have added a
// third. Naming by the decision means one file, and it is the same decision C1
// turns into a Retry type, so the file already has the name that type will want.
//
// Merging also puts the two mirrors next to each other, which is what the
// symmetry rule in non-functional requirement #5 asks for.
package transaction

import (
	"context"
	"errors"
	"log"
	"math"
	"math/big"
	"math/bits"
	"strings"
	"time"

	"relayer/services"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum"
)

const (
	outOfGasThresholdBasisPoints = 9800
	gasProbeMultiplier           = 8
	ethRevertProbeTimeout        = 15 * time.Second
)

var evmGasHeadroomBasisPoints = []uint64{10_000, 15_000, 20_000, 30_000}

type ethRevertClass int

const (
	revertLogical ethRevertClass = iota
	revertOutOfGas
)

func applyEVMGasHeadroom(gas uint64, step int) uint64 {
	if step < 0 || step >= len(evmGasHeadroomBasisPoints) {
		step = 0
	}
	factor := evmGasHeadroomBasisPoints[step]
	hi, lo := bits.Mul64(gas, factor)
	if hi >= 10_000 {
		return ^uint64(0)
	}
	scaled, _ := bits.Div64(hi, lo, 10_000)
	return scaled
}

// classifyEthRevert replays a reverted call at the same block with more gas.
// The replay, not receipt utilization, is the primary OOG signal: EIP-150 can
// leave a nested multicall failure below a 99 percent utilization threshold.
func classifyEthRevert(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	callMsg ethereum.CallMsg,
	blockNumber *big.Int,
	replayErr error,
	gasUsed, gasLimit uint64,
) ethRevertClass {
	if revertDataPresent(replayErr) {
		return revertLogical
	}

	probeGas := gasLimit
	if probeGas == 0 {
		probeGas = gasUsed
	}
	probeGas *= gasProbeMultiplier

	probeCtx, cancel := context.WithTimeout(stdCtx, ethRevertProbeTimeout)
	defer cancel()
	if header, err := endpoint.EthClient().HeaderByNumber(probeCtx, blockNumber); err == nil && header.GasLimit > 0 && probeGas > header.GasLimit {
		log.Printf("[EthTxSender] revert probe: clamping %d gas to block limit %d", probeGas, header.GasLimit)
		probeGas = header.GasLimit
	}

	probeMsg := callMsg
	probeMsg.Gas = probeGas
	// Fees are irrelevant to eth_call and can cause a large diagnostic probe to
	// be rejected for account balance rather than actually executed.
	probeMsg.GasPrice = nil
	probeMsg.GasFeeCap = nil
	probeMsg.GasTipCap = nil

	_, probeErr := endpoint.EthClient().CallContract(probeCtx, probeMsg, blockNumber)
	switch {
	case probeErr == nil:
		return revertOutOfGas
	case revertDataPresent(probeErr):
		return revertLogical
	case isExecutionOutOfGas(probeErr):
		return revertOutOfGas
	case isGasAllowanceRejection(probeErr), isProbeTransportError(probeErr):
		if outOfGasByUtilization(gasUsed, gasLimit) {
			return revertOutOfGas
		}
		return revertLogical
	default:
		return revertLogical
	}
}

func outOfGasByUtilization(gasUsed, gasLimit uint64) bool {
	return gasLimit > 0 && gasUsed <= gasLimit && gasUsed*10_000 >= gasLimit*outOfGasThresholdBasisPoints
}

func revertDataPresent(err error) bool {
	if err == nil {
		return false
	}
	var dataErr interface{ ErrorData() interface{} }
	if !errors.As(err, &dataErr) {
		return false
	}
	switch data := dataErr.ErrorData().(type) {
	case nil:
		return false
	case string:
		return data != "" && data != "0x"
	default:
		return true
	}
}

func isExecutionOutOfGas(err error) bool {
	return err != nil && !revertDataPresent(err) && strings.Contains(strings.ToLower(err.Error()), "out of gas")
}

func isGasAllowanceRejection(err error) bool {
	if err == nil || revertDataPresent(err) {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, phrase := range []string{"gas required exceeds", "exceeds block gas limit", "intrinsic gas too low"} {
		if strings.Contains(msg, phrase) {
			return true
		}
	}
	return false
}

func isProbeTransportError(err error) bool {
	if err == nil || revertDataPresent(err) {
		return false
	}
	msg := strings.ToLower(err.Error())
	for _, phrase := range []string{
		"execution reverted", "out of gas", "invalid opcode",
		"gas required exceeds", "exceeds block gas limit", "intrinsic gas too low",
	} {
		if strings.Contains(msg, phrase) {
			return false
		}
	}
	return true
}

type gasLimitKnob string

const (
	knobEthGasLimit       gasLimitKnob = "ETH_GAS_LIMIT"
	knobMulticallGasLimit gasLimitKnob = "ETH_MULTICALL_GAS_LIMIT"
	knobMisbehaviour      gasLimitKnob = "ETH_MISBEHAVIOUR_GAS_LIMIT"
	knobDeployEstimate    gasLimitKnob = "the deploy gas estimate (estimateCosmosClientDeployGas)"
	knobFixedClientLimit  gasLimitKnob = "the hardcoded client-call gas limit in handler.go"
)

var cosmosGasHeadroom = []float64{1.3, 2.0, 3.0}

// applyCosmosGasHeadroom scales both simulated and fallback limits. Applying it
// only after successful simulation made a simulation-failure retry identical:
// an included out-of-gas singleton climbed the ladder while reusing the same
// fallback gas on every step.
func applyCosmosGasHeadroom(gas uint64, headroomStep int) uint64 {
	if headroomStep < 0 || headroomStep >= len(cosmosGasHeadroom) {
		headroomStep = 0
	}
	factor := cosmosGasHeadroom[headroomStep]
	maxUint64 := ^uint64(0)
	if gas > 0 && float64(gas) >= float64(maxUint64)/factor {
		return maxUint64
	}
	return uint64(float64(gas) * factor)
}

// applyCosmosFeeHeadroom preserves the configured fee on the historical first
// attempt, then scales it by the same relative factor as the gas ladder. This
// keeps the configured gas-price ratio from shrinking on CheckTx OOG retries.
func applyCosmosFeeHeadroom(fee int64, headroomStep int) int64 {
	if fee <= 0 {
		return fee
	}
	if headroomStep < 0 || headroomStep >= len(cosmosGasHeadroom) {
		headroomStep = 0
	}
	ratio := cosmosGasHeadroom[headroomStep] / cosmosGasHeadroom[0]
	const maxInt64 = int64(^uint64(0) >> 1)
	if float64(fee) >= float64(maxInt64)/ratio {
		return maxInt64
	}
	return int64(math.Ceil(float64(fee) * ratio))
}

type cosmosOutOfGasAction int

const (
	oogSplitBatch cosmosOutOfGasAction = iota
	oogEscalateHeadroom
	oogPermanent
)

func cosmosOutOfGas(codespace string, code uint32) bool {
	return codespace == errortypes.ErrOutOfGas.Codespace() && code == errortypes.ErrOutOfGas.ABCICode()
}

// planCosmosOutOfGas always changes the next attempt. Non-atomic batches are
// halved before retrying; atomic batches retain their all-or-nothing shape and
// instead move up the finite headroom ladder.
func planCosmosOutOfGas(msgCount, headroomStep int, finalGasLimit, maxBlockGas uint64, allowSplit bool) cosmosOutOfGasAction {
	if allowSplit && msgCount > 1 {
		return oogSplitBatch
	}
	if headroomStep+1 < len(cosmosGasHeadroom) && (maxBlockGas == 0 || finalGasLimit < maxBlockGas) {
		return oogEscalateHeadroom
	}
	return oogPermanent
}
