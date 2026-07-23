package codec

import (
	"math/big"
	"testing"

	spectreContract "relayer/bindings/SpectreClient"
	updateclientContract "relayer/bindings/UpdateClient"
	relayerclient "relayer/client"

	clienttypes "github.com/cosmos/ibc-go/v10/modules/core/02-client/types"
)

// TestCosmosUpdate_GobRoundTrip guards the Cosmos->ETH codec: the ETH
// destination must recover the exact typed message the groth16 builder produced,
// including the big.Int proof fields, across the opaque []byte boundary.
func TestCosmosUpdate_GobRoundTrip(t *testing.T) {
	const kindConsensusUpdate = 1 // services.ConsensusUpdate
	appMsg := updateclientContract.ISpectreClientMsgsMsgUpdateApplicationState{
		Time: big.NewInt(1_700_000_000),
		Proof: updateclientContract.ISpectreClientMsgsBatchProof{
			Proof:                  [8]*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5), big.NewInt(6), big.NewInt(7), big.NewInt(8)},
			Commitments:            [2]*big.Int{big.NewInt(9), big.NewInt(10)},
			CommitmentPok:          [2]*big.Int{big.NewInt(11), big.NewInt(12)},
			Bucket:                 4,
			SignerIndices:          []uint32{0, 1, 2, 3},
			PinnedValidatorIndices: []uint32{0, 1, 2, 3},
			SignerPubkeys:          [][32]byte{{0x01}, {0x02}, {0x03}, {0x04}},
			Active:                 []bool{true, true, true, false},
		},
	}
	newValSet := spectreContract.IICS07TendermintMsgsValidatorSet{TotalVotingPower: 42}

	payload, err := EncodeCosmosUpdate(kindConsensusUpdate, appMsg, newValSet)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	kind, got, gotVal, err := DecodeCosmosUpdate(payload)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if kind != kindConsensusUpdate {
		t.Fatalf("kind: want %d, got %d", kindConsensusUpdate, kind)
	}
	if got.Time.Cmp(appMsg.Time) != 0 {
		t.Fatalf("Time: want %s, got %s", appMsg.Time, got.Time)
	}
	for i := range appMsg.Proof.Proof {
		if got.Proof.Proof[i].Cmp(appMsg.Proof.Proof[i]) != 0 {
			t.Fatalf("Proof[%d]: want %s, got %s", i, appMsg.Proof.Proof[i], got.Proof.Proof[i])
		}
	}
	if got.Proof.Bucket != 4 || len(got.Proof.Active) != 4 || got.Proof.Active[3] {
		t.Fatalf("proof metadata mangled: %+v", got.Proof)
	}
	if gotVal.TotalVotingPower != 42 {
		t.Fatalf("NewValSet.TotalVotingPower: want 42, got %d", gotVal.TotalVotingPower)
	}
}

// TestBeaconUpdate_GobRoundTrip guards the ETH->Cosmos codec: the Cosmos
// destination must recover the MsgUpdateClient messages plus the EthClientState
// + SigSlot it needs for the pre-submit catch-up.
func TestBeaconUpdate_GobRoundTrip(t *testing.T) {
	msgs := []any{
		&clienttypes.MsgUpdateClient{ClientId: "08-wasm-0", Signer: "cosmos1abc"},
		&clienttypes.MsgUpdateClient{ClientId: "08-wasm-0", Signer: "cosmos1abc"},
	}
	cs := relayerclient.EthereumClientState{
		ChainID:            1,
		LatestSlot:         123456,
		SyncCommitteeSize:  512,
		IbcContractAddress: "0xabc",
	}

	payload, err := EncodeBeaconUpdate(msgs, cs, 1_700_000_000, 123500)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	gotMsgs, gotCS, sigSlot, err := DecodeBeaconUpdate(payload)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(gotMsgs) != 2 {
		t.Fatalf("want 2 msgs, got %d", len(gotMsgs))
	}
	m0, ok := gotMsgs[0].(*clienttypes.MsgUpdateClient)
	if !ok || m0.ClientId != "08-wasm-0" || m0.Signer != "cosmos1abc" {
		t.Fatalf("msg[0] mangled: %+v", gotMsgs[0])
	}
	if sigSlot != 123500 {
		t.Fatalf("sigSlot: want 123500, got %d", sigSlot)
	}
	if gotCS.LatestSlot != 123456 || gotCS.SyncCommitteeSize != 512 || gotCS.IbcContractAddress != "0xabc" {
		t.Fatalf("client state mangled: %+v", gotCS)
	}
}
