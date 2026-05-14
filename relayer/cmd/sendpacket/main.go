package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	"cosmossdk.io/math"
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
	transfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypesv2 "github.com/cosmos/ibc-go/v10/modules/core/04-channel/v2/types"
)

func main() {
	var (
		receiver       = flag.String("receiver", "", "Ethereum receiver address (0x...)")
		amount         = flag.String("amount", "1000000000", "Amount to transfer")
		denom          = flag.String("denom", "stake", "Cosmos denom to transfer")
		timeoutSeconds = flag.Int("timeout-seconds", 1800, "Timeout in seconds from now")
		clientID       = flag.String("client-id", "08-wasm-0", "Source client ID on Cosmos")
		memo           = flag.String("memo", "", "Packet memo")
		node           = flag.String("node", "http://127.0.0.1:26657", "CometBFT RPC URL")
		chainID        = flag.String("chain-id", "test-ibc-eth", "Cosmos chain ID")
		count          = flag.Int("count", 1, "Number of identical packets to send in one tx")
		gasLimit       = flag.Uint64("gas-limit", 0, "Gas limit (0 = auto: 200k per msg)")
		feeAmount      = flag.Int64("fee-amount", 0, "Fee amount (0 = auto: gasLimit * 1 denom)")
		feeDenom       = flag.String("fee-denom", "stake", "Fee denom")
	)
	flag.Parse()

	if *receiver == "" {
		fmt.Fprintln(os.Stderr, "--receiver is required")
		os.Exit(1)
	}

	privKeyHex := strings.TrimPrefix(os.Getenv("COSMOS_PRIVATE_KEY"), "0x")
	if privKeyHex == "" {
		fmt.Fprintln(os.Stderr, "COSMOS_PRIVATE_KEY environment variable is required")
		os.Exit(1)
	}

	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to decode private key: %v\n", err)
		os.Exit(1)
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	signerAddr := sdk.AccAddress(privKey.PubKey().Address())

	client, err := rpchttp.New(*node, "/websocket")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create RPC client: %v\n", err)
		os.Exit(1)
	}

	// Query account info
	accountNumber, sequence, err := queryAccountInfo(client, signerAddr.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to query account info: %v\n", err)
		os.Exit(1)
	}

	// Setup encoding config
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)
	channeltypesv2.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	txConfig := authtx.NewTxConfig(cdc, authtx.DefaultSignModes)

	// Build payloads
	timeout := uint64(time.Now().Add(time.Duration(*timeoutSeconds) * time.Second).Unix())
	var msgs []sdk.Msg

	for i := 0; i < *count; i++ {
		transferPayload := transfertypes.FungibleTokenPacketData{
			Denom:    *denom,
			Amount:   *amount,
			Sender:   signerAddr.String(),
			Receiver: strings.ToLower(*receiver),
			Memo:     *memo,
		}
		encodedPayload, err := transfertypes.EncodeABIFungibleTokenPacketData(&transferPayload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to encode packet data: %v\n", err)
			os.Exit(1)
		}

		payload := channeltypesv2.Payload{
			SourcePort:      transfertypes.PortID,
			DestinationPort: transfertypes.PortID,
			Version:         transfertypes.V1,
			Encoding:        transfertypes.EncodingABI,
			Value:           encodedPayload,
		}

		msg := &channeltypesv2.MsgSendPacket{
			SourceClient:     *clientID,
			TimeoutTimestamp: timeout,
			Payloads:         []channeltypesv2.Payload{payload},
			Signer:           signerAddr.String(),
		}
		msgs = append(msgs, msg)
	}

	// Determine gas and fee
	gl := *gasLimit
	if gl == 0 {
		gl = uint64(200000 * len(msgs))
	}
	fa := *feeAmount
	if fa == 0 {
		fa = int64(gl)
	}

	// Build transaction
	txBuilder := txConfig.NewTxBuilder()
	if err := txBuilder.SetMsgs(msgs...); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set messages: %v\n", err)
		os.Exit(1)
	}
	txBuilder.SetGasLimit(gl)
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(*feeDenom, math.NewInt(fa))))

	// Empty signature for sign bytes generation
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
		fmt.Fprintf(os.Stderr, "failed to set empty signature: %v\n", err)
		os.Exit(1)
	}

	signerData := authsigning.SignerData{
		Address:       signerAddr.String(),
		ChainID:       *chainID,
		AccountNumber: accountNumber,
		Sequence:      sequence,
		PubKey:        pubKey,
	}

	signBytes, err := authsigning.GetSignBytesAdapter(
		context.Background(),
		txConfig.SignModeHandler(),
		sdksigning.SignMode_SIGN_MODE_DIRECT,
		signerData,
		txBuilder.GetTx(),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get sign bytes: %v\n", err)
		os.Exit(1)
	}

	sigRaw, err := privKey.Sign(signBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to sign transaction: %v\n", err)
		os.Exit(1)
	}

	sigV2 := sdksigning.SignatureV2{
		PubKey: pubKey,
		Data: &sdksigning.SingleSignatureData{
			SignMode:  sdksigning.SignMode_SIGN_MODE_DIRECT,
			Signature: sigRaw,
		},
		Sequence: sequence,
	}
	if err := txBuilder.SetSignatures(sigV2); err != nil {
		fmt.Fprintf(os.Stderr, "failed to set signatures: %v\n", err)
		os.Exit(1)
	}

	txBytes, err := txConfig.TxEncoder()(txBuilder.GetTx())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode transaction: %v\n", err)
		os.Exit(1)
	}

	result, err := client.BroadcastTxSync(context.Background(), txBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to broadcast transaction: %v\n", err)
		os.Exit(1)
	}

	if result.Code != 0 {
		fmt.Fprintf(os.Stderr, "transaction failed with code %d: %s\n", result.Code, result.Log)
		os.Exit(1)
	}

	fmt.Printf("%s\n", result.Hash.String())
}

func queryAccountInfo(client *rpchttp.HTTP, address string) (uint64, uint64, error) {
	queryReq := &authtypes.QueryAccountRequest{Address: address}
	reqBytes, err := proto.Marshal(queryReq)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to marshal query request: %w", err)
	}

	result, err := client.ABCIQuery(context.Background(), "/cosmos.auth.v1beta1.Query/Account", reqBytes)
	if err != nil {
		return 0, 0, fmt.Errorf("ABCI query failed: %w", err)
	}
	if result.Response.Code != 0 {
		return 0, 0, fmt.Errorf("query failed with code %d: %s", result.Response.Code, result.Response.Log)
	}

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	authtypes.RegisterInterfaces(interfaceRegistry)

	var queryResp authtypes.QueryAccountResponse
	if err := proto.Unmarshal(result.Response.Value, &queryResp); err != nil {
		return 0, 0, fmt.Errorf("failed to unmarshal query response: %w", err)
	}

	var account sdk.AccountI
	if err := interfaceRegistry.UnpackAny(queryResp.Account, &account); err != nil {
		return 0, 0, fmt.Errorf("failed to unpack account: %w", err)
	}

	return account.GetAccountNumber(), account.GetSequence(), nil
}
