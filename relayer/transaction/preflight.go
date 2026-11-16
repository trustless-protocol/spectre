// Pre-flight simulation of an EVM packet transaction.
//
// The Cosmos half Simulates before it broadcasts (sendCosmosTxBatchAtHeadroom).
// The EVM half did not: a packet transaction that reverts was only discovered
// from its receipt, after its gas was paid, and the flush backstop then offered
// the same packet again on every pass (#575 M2). This file is the EVM mirror of
// that simulation.
package transaction

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/bits"
	"strings"

	contractICS26Router "relayer/bindings/ICS26Router"
	services "relayer/services"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// packetCallInput ABI-encodes a single packet message as the ICS26Router call
// SendEthTx is about to make. ok is false for anything that is not a packet
// (a client update, a direct light-client call): those are not simulated.
func packetCallInput(msg any) (input []byte, label string, ok bool, err error) {
	var method string
	switch m := msg.(type) {
	case contractICS26Router.IICS26RouterMsgsMsgRecvPacket:
		method, label = "recvPacket", fmt.Sprintf("recvPacket seq=%d", m.Packet.Sequence)
	case contractICS26Router.IICS26RouterMsgsMsgAckPacket:
		method, label = "ackPacket", fmt.Sprintf("ackPacket seq=%d", m.Packet.Sequence)
	case contractICS26Router.IICS26RouterMsgsMsgTimeoutPacket:
		method, label = "timeoutPacket", fmt.Sprintf("timeoutPacket seq=%d", m.Packet.Sequence)
	default:
		return nil, "", false, nil
	}
	parsedABI, err := contractICS26Router.ContractICS26RouterMetaData.GetAbi()
	if err != nil {
		return nil, "", false, fmt.Errorf("[SendEthTx] failed to load ICS26Router ABI: %w", err)
	}
	input, err = parsedABI.Pack(method, msg)
	if err != nil {
		return nil, "", false, fmt.Errorf("[SendEthTx] pack %s: %w", label, err)
	}
	return input, label, true, nil
}

// noteMinedBlock records the highest block one of this process's transactions was
// mined in, per endpoint. preflightEthCall will not trust a node that is behind it.
func (h *Handler) noteMinedBlock(client *ethclient.Client, block *big.Int) {
	if client == nil || block == nil || !block.IsUint64() {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.minedBlocks == nil {
		h.minedBlocks = make(map[*ethclient.Client]uint64)
	}
	if n := block.Uint64(); n > h.minedBlocks[client] {
		h.minedBlocks[client] = n
	}
}

func (h *Handler) minedBlockFloor(client *ethclient.Client) uint64 {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.minedBlocks[client]
}

// preflightEthCall executes a packet transaction as an eth_call at the node's
// head and returns an error when it would revert, so the caller does not
// broadcast it. No gas is spent and no nonce is taken.
//
// It returns the same error classes the receipt path returns for a mined revert
// (ErrValidatorCacheRace, ErrClientNotActive, ErrPermanentRelayFailure), so every
// layer above reacts as it already does.
//
// Only a positive revert signal stops the send:
//   - revert data;
//   - a bare "execution reverted" that still reverts with more gas. A nested
//     out-of-gas looks the same at the transaction's own gas limit, passes the
//     larger probe, and is left to the gas ladder.
//
// Everything else is not evidence, and the transaction is sent as before: a
// transport error, a top-frame out-of-gas, an unknown error shape, and a node
// that has not yet reached this process's own last mined block (a lagging
// replica would report the client update that just landed as missing).
//
// The call runs at the head block's timestamp, which is older than the block the
// transaction would mine in. SpectreClient only rejects a consensus state ahead
// of block.timestamp + clockDrift (30 s by default), so the gap matters only for
// a client created with a clock drift near the chain's block time; a refusal
// there is retried by the flush backstop.
func (h *Handler) preflightEthCall(
	stdCtx context.Context,
	endpoint services.EVMEndpoint,
	from, to common.Address,
	input []byte,
	gasLimit uint64,
	label string,
) error {
	callCtx, cancel := context.WithTimeout(stdCtx, ethRevertProbeTimeout)
	defer cancel()

	header, err := endpoint.EthClient().HeaderByNumber(callCtx, nil)
	if err != nil || header == nil || header.Number == nil {
		return nil
	}
	if floor := h.minedBlockFloor(endpoint.EthClient()); header.Number.IsUint64() && header.Number.Uint64() < floor {
		log.Printf("[SendEthTx] pre-flight skipped for %s: node head %s is behind this relayer's last mined block %d",
			label, header.Number, floor)
		return nil
	}

	callMsg := ethereum.CallMsg{From: from, To: &to, Gas: gasLimit, Data: input}
	_, callErr := endpoint.EthClient().CallContract(callCtx, callMsg, header.Number)
	if callErr == nil {
		return nil
	}

	if !revertDataPresent(callErr) {
		if !isBareExecutionRevert(callErr) {
			return nil
		}
		probeGas := preflightProbeGas(gasLimit, header.GasLimit)
		if probeGas > gasLimit {
			probeMsg := callMsg
			probeMsg.Gas = probeGas
			_, probeErr := endpoint.EthClient().CallContract(callCtx, probeMsg, header.Number)
			if probeErr == nil || !(revertDataPresent(probeErr) || isBareExecutionRevert(probeErr)) {
				return nil
			}
		}
	}

	if name, ok := validatorCacheRaceErrorName(callErr); ok {
		return fmt.Errorf("%s would revert with %s at block %s; not broadcast: %w",
			label, name, header.Number, services.ErrValidatorCacheRace)
	}
	// Before the generic ending, for the same reason as on the receipt path: a
	// frozen client is not a poison packet.
	if name, ok := clientNotActiveErrorName(callErr); ok {
		return fmt.Errorf("%s would revert with %s at block %s; not broadcast: %w",
			label, name, header.Number, services.ErrClientNotActive)
	}
	cause := ""
	if decoded, ok := decodeRevert(callErr); ok {
		cause = " with " + decoded
	}
	log.Printf("WARN [SendEthTx] pre-flight: %s would revert%s at block %s; not broadcast, no gas spent",
		label, cause, header.Number)
	return fmt.Errorf("%s would revert%s at block %s; not broadcast, no gas spent: %w",
		label, cause, header.Number, services.ErrPermanentRelayFailure)
}

// isBareExecutionRevert reports a revert that carries no data: the node executed
// the call and it reverted, but nothing says whether a require() or a nested
// out-of-gas did it.
func isBareExecutionRevert(err error) bool {
	return err != nil && !revertDataPresent(err) && strings.Contains(strings.ToLower(err.Error()), "execution reverted")
}

// preflightProbeGas is the larger gas the ambiguous case is replayed with:
// gasProbeMultiplier times the limit, clamped to the block gas limit.
func preflightProbeGas(gasLimit, blockGasLimit uint64) uint64 {
	hi, probe := bits.Mul64(gasLimit, gasProbeMultiplier)
	if hi != 0 {
		probe = ^uint64(0)
	}
	if blockGasLimit > 0 && probe > blockGasLimit {
		probe = blockGasLimit
	}
	return probe
}
