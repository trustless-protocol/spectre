package runner

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cometbft/cometbft/rpc/client/http"

	tendermintContract "operator/bindings/SP1ICS07Tendermint"
	tendermintClient "operator/client"
)

func RunMembership(client *http.HTTP, keyPathsStr string, trustedBlock int64, isBase64 bool) ([]tendermintContract.IMembershipMsgsKVPair, []tendermintContract.IMembershipMsgsMerkleProof, error) {
	keyPaths := strings.Split(keyPathsStr, ",")

	kvPairs := []tendermintContract.IMembershipMsgsKVPair{}
	proofs := []tendermintContract.IMembershipMsgsMerkleProof{}
	for _, keyPath := range keyPaths {
		path := make([][]byte, 0, 2)
		if isBase64 {
			parts := strings.Split(keyPath, "\\")
			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("invalid key path format: %s", keyPath)
			}
			for _, p := range parts {
				decoded, err := base64.StdEncoding.DecodeString(p)
				if err != nil {
					return nil, nil, err
				}
				path = append(path, decoded)
			}
		} else {
			path = [][]byte{
				[]byte("ibc"),
				[]byte(keyPath),
			}
		}

		value, proof, err := tendermintClient.ProvePath(client, trustedBlock, path)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to prove path: %w", err)
		}

		kvPair := tendermintContract.IMembershipMsgsKVPair{
			Path:  path,
			Value: bytesToBytes32(value),
		}

		kvPairs = append(kvPairs, kvPair)
		merkleProof := tendermintContract.IMembershipMsgsMerkleProof{
			Proofs: []tendermintContract.IMembershipMsgsCommitmentProof{},
		}
		for _, p := range proof.Proofs {
			commitmentProof, err := tendermintClient.ParseCommitmentProof(p)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to parse commitment proof: %w", err)
			}
			merkleProof.Proofs = append(merkleProof.Proofs, *commitmentProof)
		}
		proofs = append(proofs, merkleProof)
	}

	return kvPairs, proofs, nil
}

func bytesToBytes32(data []byte) [32]byte {
	var result [32]byte
	copy(result[:], data)
	return result
}
