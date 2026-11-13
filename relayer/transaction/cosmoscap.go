// Operator caps on what one Cosmos relay transaction may spend (#575 M4).
//
// This is the Cosmos mirror of feecap.go. A relay transaction's gas limit is the
// GasUsed the RPC's Simulate query reports, times a headroom factor, and by
// default its fee equals that gas limit. The only bound was the block's MaxGas,
// read from the same RPC and absent when it is -1 or the query fails. So the RPC
// chose the gas, and through it the fee.
package transaction

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	services "relayer/services"
)

const (
	// EnvCosmosMaxGas caps the gas limit of one relay transaction.
	EnvCosmosMaxGas = "COSMOS_MAX_GAS"
	// EnvCosmosMaxFee caps the fee of one relay transaction, in base units of
	// COSMOS_FEE_DENOM.
	EnvCosmosMaxFee = "COSMOS_MAX_FEE"
)

// cosmosSpendCaps are the operator's per-transaction caps. Zero is no cap.
type cosmosSpendCaps struct {
	maxGas uint64
	maxFee int64
}

// envSpendCap reads one cap as a non-negative integer no larger than max. set
// reports whether the variable is present; a present 0 is an explicit "no cap".
func envSpendCap(name string, max uint64) (value uint64, set bool, err error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return 0, false, nil
	}
	value, err = strconv.ParseUint(raw, 10, 64)
	if err != nil || value > max {
		return 0, false, fmt.Errorf("invalid %s=%q: want an integer from 0 to %d (0 = no cap)", name, raw, max)
	}
	return value, true, nil
}

// loadCosmosSpendCaps reads both caps. An unset variable is no cap: only `start`
// requires them (RequireCosmosSpendCaps).
func loadCosmosSpendCaps() (cosmosSpendCaps, error) {
	// A gas limit above MaxInt64 cannot be turned into a fee amount, so it is not
	// a usable cap either.
	maxGas, _, err := envSpendCap(EnvCosmosMaxGas, math.MaxInt64)
	if err != nil {
		return cosmosSpendCaps{}, err
	}
	maxFee, _, err := envSpendCap(EnvCosmosMaxFee, math.MaxInt64)
	if err != nil {
		return cosmosSpendCaps{}, err
	}
	return cosmosSpendCaps{maxGas: maxGas, maxFee: int64(maxFee)}, nil
}

// RequireCosmosSpendCaps is `start`'s gate: both caps must be set, 0 included.
func RequireCosmosSpendCaps() error {
	var missing []string
	for _, name := range []string{EnvCosmosMaxGas, EnvCosmosMaxFee} {
		_, set, err := envSpendCap(name, math.MaxInt64)
		if err != nil {
			return err
		}
		if !set {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"%s not set: `start` signs Cosmos transactions with the gas the RPC's Simulate reports and needs a cap for each, per transaction (0 = no cap)",
			strings.Join(missing, " and "))
	}
	return nil
}

// gasCeiling is the most gas one transaction may carry: the tighter of the
// block's limit (0 when the chain reports none) and the operator's cap.
func (c cosmosSpendCaps) gasCeiling(maxBlockGas uint64) uint64 {
	switch {
	case c.maxGas == 0:
		return maxBlockGas
	case maxBlockGas == 0 || c.maxGas < maxBlockGas:
		return c.maxGas
	default:
		return maxBlockGas
	}
}

// gasAboveCap reports a transaction that cannot succeed under the operator's gas
// cap: the node already says it uses more. Broadcasting it at the cap would pay
// for a certain out-of-gas.
func (c cosmosSpendCaps) gasAboveCap(simulatedGas uint64) error {
	if c.maxGas == 0 || simulatedGas <= c.maxGas {
		return nil
	}
	return fmt.Errorf("simulated gas %d is above %s=%d; not signed: %w",
		simulatedGas, EnvCosmosMaxGas, c.maxGas, services.ErrSpendCapExceeded)
}

// feeAboveCap reports a fee the operator's cap does not allow.
func (c cosmosSpendCaps) feeAboveCap(fee int64, denom string) error {
	if c.maxFee == 0 || fee <= c.maxFee {
		return nil
	}
	return fmt.Errorf("fee %d%s is above %s=%d; not signed: %w",
		fee, denom, EnvCosmosMaxFee, c.maxFee, services.ErrSpendCapExceeded)
}
