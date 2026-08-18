package transaction

import (
	"math"

	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
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
