package transaction

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"operator/keys"
	"os"
	"strings"
	"time"

	contractICS26Router "operator/bindings/ICS26Router"
	routerContract "operator/bindings/ICS26Router"
	tendermintContract "operator/bindings/SP1ICS07Tendermint"
	updateclient "operator/bindings/UpdateClient"
	operatorclient "operator/client"
	services "operator/services"
	utils "operator/utils"

	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdksigning "github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/gogoproto/proto"
	ibcwasmtypes "github.com/cosmos/ibc-go/modules/light-clients/08-wasm/v10/types"
	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
	clienttypesv2 "github.com/cosmos/ibc-go/v10/modules/core/02-client/v2/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
	exported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

type Handler struct {
}

func (h *Handler) CreateCosmosClientContract(ctx services.Context, clientState, consensusHash []byte) error {
	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to restore private key: %w", err)
	}

	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := ctx.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainIdInt, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("invalid chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return fmt.Errorf("failed to create auth transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // in wei
	auth.GasLimit = uint64(10000000) // in units
	auth.GasPrice = gasPrice

	address, tx, _, err := tendermintContract.DeployContractSP1ICS07Tendermint(
		auth,
		ctx.EthClient(),
		*ctx.VerifierContract(),
		*ctx.MembershipContract(),
		*ctx.MisbehaviourContract(),
		*ctx.UpdateClientContract(),
		clientState,
		utils.BytesToBytes32(consensusHash),
		*ctx.RoleManagerAddress(),
	)

	if err != nil {
		return fmt.Errorf("failed to deploy ics07 contract: %w", err)
	}
	log.Printf("[CreateCosmosClient] Deploy tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed waiting for deploy receipt: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("deploy tx %s reverted (gasUsed=%d)", tx.Hash().Hex(), receipt.GasUsed)
	}
	log.Printf("[CreateCosmosClient] ICS07 deployed at %s (block %d, gasUsed=%d)", address.String(), receipt.BlockNumber.Uint64(), receipt.GasUsed)
	ctx.SetClient(address)

	// Grant PROOF_SUBMITTER_ROLE to ICS26Router so it can call verifyMembership
	ics07Instance, err := tendermintContract.NewContractSP1ICS07Tendermint(address, ctx.EthClient())
	if err != nil {
		return fmt.Errorf("failed to instantiate ICS07 contract: %w", err)
	}
	proofSubmitterRole := crypto.Keccak256Hash([]byte("PROOF_SUBMITTER_ROLE"))
	nonce, err = ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce for grantRole: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	tx, err = ics07Instance.GrantRole(auth, proofSubmitterRole, *ctx.RouterContract())
	if err != nil {
		return fmt.Errorf("failed to grant PROOF_SUBMITTER_ROLE to ICS26Router: %w", err)
	}
	log.Printf("[CreateCosmosClient] GrantRole tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receipt, err = bind.WaitMined(context.Background(), ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed waiting for GrantRole receipt: %w", err)
	}
	if receipt.Status == 0 {
		return fmt.Errorf("GrantRole tx %s reverted (gasUsed=%d)", tx.Hash().Hex(), receipt.GasUsed)
	}
	log.Printf("[CreateCosmosClient] PROOF_SUBMITTER_ROLE granted to ICS26Router (gasUsed=%d)", receipt.GasUsed)

	ics26Router, err := routerContract.NewContractICS26Router(*ctx.RouterContract(), ctx.EthClient())
	if err != nil {
		return err
	}

	nonce, err = ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))

	tx, err = ics26Router.AddClient(
		auth,
		"cosmoshub-1",
		routerContract.IICS02ClientMsgsCounterpartyInfo{
			ClientId:     "08-wasm-0",
			MerklePrefix: [][]byte{[]byte(exported.StoreKey), []byte("")},
		},
		*ctx.ClientContract(),
	)

	if err != nil {
		return fmt.Errorf("failed to add client to ICS26Router: %w", err)
	}
	log.Printf("[CreateCosmosClient] AddClient tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receipt, err = bind.WaitMined(context.Background(), ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed waiting for AddClient receipt: %w", err)
	}
	if receipt.Status == 0 {
		log.Printf("[CreateCosmosClient] AddClient tx reverted (gasUsed=%d) — client may already exist, continuing...", receipt.GasUsed)
		return nil
	}
	log.Printf("[CreateCosmosClient] AddClient confirmed (block %d, gasUsed=%d)", receipt.BlockNumber.Uint64(), receipt.GasUsed)

	return nil
}

func (h *Handler) SendEthTx(ctx services.Context, msg any) error {
	privKey := os.Getenv("ETH_PRIVATE_KEY")
	if privKey == "" {
		return fmt.Errorf("ETH_PRIVATE_KEY environment variable is required in .env file")
	}
	privateKey, err := keys.RestoreKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to restore private key: %w", err)
	}

	publicKey, err := keys.PublicKey(privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fromAddress := crypto.PubkeyToAddress(*publicKey)
	nonce, err := ctx.EthClient().PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice, err := ctx.EthClient().SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainIdInt, err := ctx.EthClient().ChainID(context.Background())
	if err != nil {
		return fmt.Errorf("invalid chain id: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainIdInt)
	if err != nil {
		return fmt.Errorf("failed to create auth transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(3000000) // in units
	auth.GasPrice = gasPrice

	ics07Tendermint, err := tendermintContract.NewContractSP1ICS07Tendermint(
		*ctx.ClientContract(),
		ctx.EthClient(),
	)
	if err != nil {
		return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
	}

	icS26Router, err := contractICS26Router.NewContractICS26Router(
		*ctx.RouterContract(),
		ctx.EthClient(),
	)
	if err != nil {
		return fmt.Errorf("failed to create ICS07 Tendermint contract: %w", err)
	}

	var tx *types.Transaction
	switch msg := msg.(type) {
	case updateclient.IUpdateClientMsgsMsgUpdateClient:
		log.Printf("[SendEthTx] Encoding updateClient msg...")
		data, err := operatorclient.EncodeUpdateClientMsg(msg)
		if err != nil {
			return fmt.Errorf("failed to encode updateClient msg: %w", err)
		}

		log.Printf("[SendEthTx] Sending updateClient tx...")
		tx, err = ics07Tendermint.UpdateClient(auth, data)
		if err != nil {
			return fmt.Errorf("failed to send updateClient tx: %w", err)
		}
	case tendermintContract.ILightClientMsgsMsgVerifyMembership:
		log.Printf("[SendEthTx] Sending verifyMembership tx...")
		tx, err = ics07Tendermint.VerifyMembership(auth, msg)
		if err != nil {
			return fmt.Errorf("failed to verify membership: %w", err)
		}
	case tendermintContract.ILightClientMsgsMsgVerifyNonMembership:
		log.Printf("[SendEthTx] Sending verifyNonMembership tx...")
		tx, err = ics07Tendermint.VerifyNonMembership(auth, msg)
		if err != nil {
			return fmt.Errorf("failed to verify non-membership: %w", err)
		}
	case contractICS26Router.IICS26RouterMsgsMsgRecvPacket:
		log.Printf("[SendEthTx] Sending recvPacket tx...")
		tx, err = icS26Router.RecvPacket(auth, msg)
		if err != nil {
			return fmt.Errorf("failed to recv packet: %w", err)
		}
	default:
		return fmt.Errorf("unsupported message type: %T", msg)
	}

	log.Printf("[SendEthTx] Tx sent: %s. Waiting for receipt...", tx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), ctx.EthClient(), tx)
	if err != nil {
		return fmt.Errorf("failed waiting for tx receipt: %w", err)
	}
	if receipt.Status == 0 {
		// Try to get revert reason by replaying the tx via eth_call
		callMsg := ethereum.CallMsg{
			From:     fromAddress,
			To:       tx.To(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice(),
			Value:    tx.Value(),
			Data:     tx.Data(),
		}
		_, callErr := ctx.EthClient().CallContract(context.Background(), callMsg, receipt.BlockNumber)
		if callErr != nil {
			log.Printf("[SendEthTx] Revert reason: %v", callErr)
			// Extract hex-encoded revert data for custom error decoding
			type dataErr interface {
				ErrorData() interface{}
			}
			if de, ok := callErr.(dataErr); ok {
				log.Printf("[SendEthTx] Revert data (hex): %v", de.ErrorData())
			}
		}
		return fmt.Errorf("tx %s reverted (status=0, gasUsed=%d)", tx.Hash().Hex(), receipt.GasUsed)
	}
	log.Printf("[SendEthTx] Tx %s confirmed in block %d (gasUsed=%d)", tx.Hash().Hex(), receipt.BlockNumber.Uint64(), receipt.GasUsed)

	return nil
}

func (h *Handler) CreateEthClient(svcCtx services.Context, clientState exported.ClientState, consensusState exported.ConsensusState) error {

	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())
	fmt.Println("signerAddr: ", signerAddr.String())

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration
	gasLimit := uint64(200000) // Default gas limit
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		if _, err := fmt.Sscanf(gasStr, "%d", &gasLimit); err != nil {
			return fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	feeAmount := int64(10000000) // Default fee amount
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		if _, err := fmt.Sscanf(feeStr, "%d", &feeAmount); err != nil {
			return fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
	}

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	clienttypes.RegisterInterfaces(interfaceRegistry)
	ibcwasmtypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	msg, err := clienttypes.NewMsgCreateClient(clientState, consensusState, signerAddr.String())
	if err != nil {
		return err
	}

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msg); err != nil {
		return fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction
	result, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("MsgCreateClient broadcast successfully. Hash: %s", result.Hash.String())

	// Wait for MsgCreateClient tx to be included in a block before sending the next tx
	if err := h.waitForTx(svcCtx, result.Hash, 30*time.Second); err != nil {
		return fmt.Errorf("failed waiting for MsgCreateClient tx: %w", err)
	}

	// Build and broadcast MsgRegisterCounterparty as a separate transaction
	registerMsg := clienttypesv2.NewMsgRegisterCounterparty(
		"08-wasm-0",
		[][]byte{[]byte(exported.StoreKey), []byte("")},
		"cosmoshub-1",
		signerAddr.String(),
	)

	// Re-query account info (sequence incremented after first tx)
	accountNumber, sequence, err = h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info for register counterparty: %w", err)
	}

	txBuilder2 := txConfig.NewTxBuilder()
	if err := txBuilder2.SetMsgs(registerMsg); err != nil {
		return fmt.Errorf("failed to set register counterparty message: %w", err)
	}
	txBuilder2.SetGasLimit(gasLimit)
	txBuilder2.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	emptySig2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder2.SetSignatures(emptySig2); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	signerData2 := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	signBytes2, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData2,
		txBuilder2.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes for register counterparty: %w", err)
	}

	sigRaw2, err := privKey.Sign(signBytes2)
	if err != nil {
		return fmt.Errorf("failed to sign register counterparty tx: %w", err)
	}

	sigV22 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw2,
		},
		Sequence: sequence,
	}
	if err := txBuilder2.SetSignatures(sigV22); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	txBytes2, err := txConfig.TxEncoder()(txBuilder2.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode register counterparty tx: %w", err)
	}

	result2, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes2)
	if err != nil {
		return fmt.Errorf("failed to broadcast register counterparty tx: %w", err)
	}

	if result2.Code != 0 {
		return fmt.Errorf("register counterparty tx failed with code %d: %s", result2.Code, result2.Log)
	}

	log.Printf("MsgRegisterCounterparty broadcast successfully. Hash: %s", result2.Hash.String())

	return nil
}

func (h *Handler) SendCosmosTx(svcCtx services.Context, msg any) error {
	protoMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("message must be a proto.Message")
	}
	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration
	gasLimit := uint64(200000) // Default gas limit
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		if _, err := fmt.Sscanf(gasStr, "%d", &gasLimit); err != nil {
			return fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	feeAmount := int64(1000) // Default fee amount
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		if _, err := fmt.Sscanf(feeStr, "%d", &feeAmount); err != nil {
			return fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
	}

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()

	// Convert the proto.Message to sdk.Msg
	sdkMsg, ok := protoMsg.(sdk.Msg)
	if !ok {
		return fmt.Errorf("message does not implement sdk.Msg interface")
	}

	if err := txBuilder.SetMsgs(sdkMsg); err != nil {
		return fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction
	result, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("Transaction broadcast successfully. Hash: %s", result.Hash.String())

	return nil
}

// SendCosmosTxBatch sends multiple messages in a single Cosmos transaction
func (h *Handler) SendCosmosTxBatch(svcCtx services.Context, msgs []any) error {
	if len(msgs) == 0 {
		return nil
	}

	// Convert all messages to sdk.Msg
	var sdkMsgs []sdk.Msg
	for i, msg := range msgs {
		protoMsg, ok := msg.(proto.Message)
		if !ok {
			return fmt.Errorf("message %d must be a proto.Message", i)
		}
		sdkMsg, ok := protoMsg.(sdk.Msg)
		if !ok {
			return fmt.Errorf("message %d does not implement sdk.Msg interface", i)
		}
		sdkMsgs = append(sdkMsgs, sdkMsg)
	}

	// Get the private key from environment variable
	privKeyHex := os.Getenv("COSMOS_PRIVATE_KEY")
	if privKeyHex == "" {
		return fmt.Errorf("COSMOS_PRIVATE_KEY environment variable is required in .env file")
	}

	// Decode the private key
	privKeyBytes, err := hex.DecodeString(strings.TrimPrefix(privKeyHex, "0x"))
	if err != nil {
		return fmt.Errorf("failed to decode private key: %w", err)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())

	// Get chain configuration from environment
	chainID := os.Getenv("COSMOS_CHAIN_ID")
	if chainID == "" {
		return fmt.Errorf("COSMOS_CHAIN_ID environment variable is required in .env file")
	}

	// Get gas and fee configuration - use higher gas for batch transactions
	gasLimit := uint64(200000) * uint64(len(sdkMsgs)) // Scale gas with number of messages
	if gasStr := os.Getenv("COSMOS_GAS_LIMIT"); gasStr != "" {
		var baseGas uint64
		if _, err := fmt.Sscanf(gasStr, "%d", &baseGas); err != nil {
			return fmt.Errorf("failed to parse COSMOS_GAS_LIMIT: %w", err)
		}
		gasLimit = baseGas * uint64(len(sdkMsgs))
	}

	feeDenom := os.Getenv("COSMOS_FEE_DENOM")
	if feeDenom == "" {
		feeDenom = "stake" // Default fee denom
	}

	feeAmount := int64(1000) * int64(len(sdkMsgs)) // Scale fee with number of messages
	if feeStr := os.Getenv("COSMOS_FEE_AMOUNT"); feeStr != "" {
		var baseFee int64
		if _, err := fmt.Sscanf(feeStr, "%d", &baseFee); err != nil {
			return fmt.Errorf("failed to parse COSMOS_FEE_AMOUNT: %w", err)
		}
		feeAmount = baseFee * int64(len(sdkMsgs))
	}

	// Query account info (account number and sequence) from the chain
	accountNumber, sequence, err := h.queryAccountInfo(svcCtx, signerAddr.String())
	if err != nil {
		return fmt.Errorf("failed to query account info: %w", err)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Build the transaction
	txBuilder := txConfig.NewTxBuilder()

	if err := txBuilder.SetMsgs(sdkMsgs...); err != nil {
		return fmt.Errorf("failed to set messages: %w", err)
	}

	txBuilder.SetGasLimit(gasLimit)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(feeDenom, sdkmath.NewInt(feeAmount))))

	// First, set an empty signature to populate signer info for sign bytes generation
	pubKey := privKey.PubKey()
	emptySig := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: nil,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(emptySig); err != nil {
		return fmt.Errorf("failed to set empty signature: %w", err)
	}

	// Create signer data
	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	// Get sign bytes using the adapter function
	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		return fmt.Errorf("failed to get sign bytes: %w", err)
	}

	// Sign the bytes
	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Set the actual signature
	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}

	if err := txBuilder.SetSignatures(sigV2); err != nil {
		return fmt.Errorf("failed to set signatures: %w", err)
	}

	// Encode the transaction
	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		return fmt.Errorf("failed to encode transaction: %w", err)
	}

	// Broadcast the transaction
	result, err := svcCtx.CosmosClient().BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		return fmt.Errorf("failed to broadcast transaction: %w", err)
	}

	if result.Code != 0 {
		return fmt.Errorf("transaction failed with code %d: %s", result.Code, result.Log)
	}

	log.Printf("Batch transaction broadcast successfully. Hash: %s, Messages: %d", result.Hash.String(), len(sdkMsgs))

	return nil
}

// queryAccountInfo queries the account number and sequence for the given address
func (h *Handler) queryAccountInfo(svcCtx services.Context, address string) (uint64, uint64, error) {
	// Build the query request
	queryReq := &authtypes.QueryAccountRequest{
		Address: address,
	}

	reqBytes, err := proto.Marshal(queryReq)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal query request: %w", err)
	}

	// Query path for auth account
	queryPath := "/cosmos.auth.v1beta1.Query/Account"

	// Make ABCI query
	result, err := svcCtx.CosmosClient().ABCIQuery(context.Background(), queryPath, reqBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("ABCI query failed: %w", err)
	}

	if result.Response.Code != 0 {
		return 0, 0, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	// Setup interface registry to decode the account
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	// Decode the response
	var queryResp authtypes.QueryAccountResponse
	if err := proto.Unmarshal(result.Response.Value, &queryResp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal query response: %w", err)
	}

	// Unpack the account from Any type
	var account sdk.AccountI
	if err := interfaceRegistry.UnpackAny(queryResp.Account, &account); err != nil {
		return 0, 0, fmt.Errorf("failed to unpack account: %w", err)
	}

	return account.GetAccountNumber(), account.GetSequence(), nil
}

// waitForTx polls the chain until the transaction with the given hash is included in a block or the timeout expires.
func (h *Handler) waitForTx(svcCtx services.Context, txHash []byte, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		result, err := svcCtx.CosmosClient().Tx(context.Background(), txHash, false)
		if err == nil && result != nil && result.Height > 0 {
			log.Printf("Tx %X confirmed at height %d", txHash, result.Height)
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("timeout waiting for tx %X to be included in a block", txHash)
}
