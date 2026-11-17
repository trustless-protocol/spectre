// Operator caps on what an EVM transaction may pay per gas (#575 M3).
//
// Every fee this package signs is read from the RPC: the tip
// (eth_maxPriorityFeePerGas), the base fee (latest header), the legacy price
// (eth_gasPrice), and each replacement is the larger of a 15% bump and a fresh
// suggestion. Nothing bounded them, so whoever answers the RPC chose the price of
// every transaction the relayer signed.
//
// The rule: nothing is signed above a cap. A fee above its cap is held at the
// cap. One transaction is refused instead: a first send the caps cannot get
// mined now (base fee plus tip above the fee cap). Sending it would only park a
// nonce until the give-up path has to cancel it.
package transaction

import (
	"fmt"
	"math/big"
	"os"
	"strings"

	services "relayer/services"

	"github.com/ethereum/go-ethereum/core/types"
)

const (
	// EnvEVMMaxFeePerGas caps maxFeePerGas (EIP-1559) and gasPrice (legacy), in wei.
	EnvEVMMaxFeePerGas = "ETH_MAX_FEE_PER_GAS"
	// EnvEVMMaxPriorityFeePerGas caps maxPriorityFeePerGas, in wei.
	EnvEVMMaxPriorityFeePerGas = "ETH_MAX_PRIORITY_FEE_PER_GAS"
)

// evmFeeCaps are the operator's per-gas caps in wei. A nil cap is no cap.
type evmFeeCaps struct {
	maxFee *big.Int
	maxTip *big.Int
}

// envWeiCap reads one cap. set reports whether the variable is present. A present
// 0 is an explicit "no cap" and returns a nil value.
func envWeiCap(name string) (value *big.Int, set bool, err error) {
	raw := strings.TrimSpace(os.Getenv(name))
	if raw == "" {
		return nil, false, nil
	}
	parsed, ok := new(big.Int).SetString(raw, 10)
	if !ok || parsed.Sign() < 0 {
		return nil, false, fmt.Errorf("invalid %s=%q: want a non-negative integer in wei per gas (0 = no cap)", name, raw)
	}
	if parsed.Sign() == 0 {
		return nil, true, nil
	}
	return parsed, true, nil
}

// loadEVMFeeCaps reads both caps. An unset variable is no cap: only `start`
// requires them (RequireEVMFeeCaps), so a one-shot command keeps working without.
func loadEVMFeeCaps() (evmFeeCaps, error) {
	maxFee, _, err := envWeiCap(EnvEVMMaxFeePerGas)
	if err != nil {
		return evmFeeCaps{}, err
	}
	maxTip, _, err := envWeiCap(EnvEVMMaxPriorityFeePerGas)
	if err != nil {
		return evmFeeCaps{}, err
	}
	return evmFeeCaps{maxFee: maxFee, maxTip: maxTip}, nil
}

// RequireEVMFeeCaps is `start`'s gate: both caps must be set. Running without a
// cap is allowed, but only as a decision written down as 0.
func RequireEVMFeeCaps() error {
	var missing []string
	for _, name := range []string{EnvEVMMaxFeePerGas, EnvEVMMaxPriorityFeePerGas} {
		_, set, err := envWeiCap(name)
		if err != nil {
			return err
		}
		if !set {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf(
			"%s not set: `start` signs EVM transactions at fees read from the RPC and needs a cap for each, in wei per gas (0 = no cap)",
			strings.Join(missing, " and "))
	}
	return nil
}

// minBig returns the smaller of v and limit. A nil limit is no limit.
func minBig(v, limit *big.Int) *big.Int {
	if limit != nil && v.Cmp(limit) > 0 {
		return limit
	}
	return v
}

// dynamicFees prices a first EIP-1559 send: the suggested tip, and a fee cap of
// twice the base fee plus that tip, each held at its cap. It refuses when the fee
// cap does not cover the current base fee plus the tip, because that transaction
// cannot be mined now.
func (c evmFeeCaps) dynamicFees(baseFee, suggestedTip *big.Int) (tipCap, feeCap *big.Int, err error) {
	tipCap = minBig(suggestedTip, c.maxTip)
	feeCap = new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)
	if c.maxFee == nil {
		return tipCap, feeCap, nil
	}
	if needed := new(big.Int).Add(baseFee, tipCap); needed.Cmp(c.maxFee) > 0 {
		return nil, nil, fmt.Errorf("base fee %s + tip %s wei per gas is above %s=%s; not signed: %w",
			baseFee, tipCap, EnvEVMMaxFeePerGas, c.maxFee, services.ErrSpendCapExceeded)
	}
	return tipCap, minBig(feeCap, c.maxFee), nil
}

// legacyPrice prices a first legacy send. market is the price the network asks
// now; price is market raised by the same-nonce replacement floor. It refuses
// when the market price itself is above the cap, and otherwise holds price at it.
func (c evmFeeCaps) legacyPrice(market, price *big.Int) (*big.Int, error) {
	if c.maxFee != nil && market.Cmp(c.maxFee) > 0 {
		return nil, fmt.Errorf("gas price %s wei is above %s=%s; not signed: %w",
			market, EnvEVMMaxFeePerGas, c.maxFee, services.ErrSpendCapExceeded)
	}
	return minBig(price, c.maxFee), nil
}

// clampDynamic holds a replacement's fees at the caps. The result may be too low
// to replace the transaction it targets; the node decides that, and refusing
// costs nothing.
func (c evmFeeCaps) clampDynamic(tipCap, feeCap *big.Int) (*big.Int, *big.Int) {
	tipCap, feeCap = minBig(tipCap, c.maxTip), minBig(feeCap, c.maxFee)
	if tipCap.Cmp(feeCap) > 0 {
		tipCap = feeCap // a tip above the fee cap is not a valid transaction
	}
	return tipCap, feeCap
}

// clampGasPrice holds a legacy replacement's price at the fee cap.
func (c evmFeeCaps) clampGasPrice(price *big.Int) *big.Int {
	return minBig(price, c.maxFee)
}

// check is the last gate before a signature. The pricing above already keeps a
// transaction inside the caps, so this only fires if a new path forgets to.
func (c evmFeeCaps) check(tx *types.Transaction) error {
	if c.maxFee != nil && tx.GasFeeCap().Cmp(c.maxFee) > 0 {
		return fmt.Errorf("refusing to sign: fee %s wei per gas is above %s=%s: %w",
			tx.GasFeeCap(), EnvEVMMaxFeePerGas, c.maxFee, services.ErrSpendCapExceeded)
	}
	// A legacy or access-list transaction has no separate tip: its whole price is
	// the fee, checked above.
	hasTip := tx.Type() != types.LegacyTxType && tx.Type() != types.AccessListTxType
	if hasTip && c.maxTip != nil && tx.GasTipCap().Cmp(c.maxTip) > 0 {
		return fmt.Errorf("refusing to sign: priority fee %s wei per gas is above %s=%s: %w",
			tx.GasTipCap(), EnvEVMMaxPriorityFeePerGas, c.maxTip, services.ErrSpendCapExceeded)
	}
	return nil
}
