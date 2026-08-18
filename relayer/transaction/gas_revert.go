package transaction

import (
	"context"
	"errors"
	"log"
	"math/big"
	"math/bits"
	"strings"
	"time"

	"relayer/services"

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
