package relayer

// EthCosmosConfigInfo is a struct that holds the configuration information for the Eth to Cosmos config template
type EthCosmosConfigInfo struct {
	// Ethereum chain identifier
	EthChainID string
	// Cosmos chain identifier
	CosmosChainID string
	// Tendermint RPC URL
	TmRPC string
	// ICS26 Router address
	ICS26Address string
	// Ethereum RPC URL
	EthRPC string
	// Ethereum Beacon API URL
	BeaconAPI string
	// Signer address cosmos
	SignerAddress string
	// Whether we use the mock client in Cosmos
	MockWasmClient bool
	// ICS07 Tendermint light client address (set after create-clients)
	ICS07Client string
	// WrapperVerifier contract address
	WrapperVerifier string
	// Membership program contract address
	Membership string
	// Misbehaviour program contract address
	Misbehaviour string
	// UpdateClient program contract address
	UpdateClient string
}

func CreateEthCosmosModules(
	configInfo EthCosmosConfigInfo,
) []ModuleConfig {
	return []ModuleConfig{
		{
			Name:     ModuleEthToCosmosCompat,
			SrcChain: configInfo.EthChainID,
			DstChain: configInfo.CosmosChainID,
			Config: ethToCosmosConfig{
				TmRpcUrl:        configInfo.TmRPC,
				Ics26Address:    configInfo.ICS26Address,
				EthRpcUrl:       configInfo.EthRPC,
				EthBeaconApiUrl: configInfo.BeaconAPI,
				SignerAddress:   configInfo.SignerAddress,
				Mock:            configInfo.MockWasmClient,
			},
		},
		{
			Name:     ModuleCosmosToEth,
			SrcChain: configInfo.CosmosChainID,
			DstChain: configInfo.EthChainID,
			Config: CosmosToEthModuleConfig{
				TmRpcUrl:        configInfo.TmRPC,
				Ics26Address:    configInfo.ICS26Address,
				EthRpcUrl:       configInfo.EthRPC,
				ICS07Client:     configInfo.ICS07Client,
				WrapperVerifier: configInfo.WrapperVerifier,
				Membership:      configInfo.Membership,
				Misbehaviour:    configInfo.Misbehaviour,
				UpdateClient:    configInfo.UpdateClient,
			},
		},
	}
}
