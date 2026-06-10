package prover

import (
	"fmt"
	"math/big"
)

// HeaderBatchProof is the relayer-side representation of one Solidity
// IMisbehaviourMsgs.BatchProof. PaddedSignatures carries the bucket-sized
// signer metadata that must be sent with the proof.
type HeaderBatchProof struct {
	Bucket        int
	PaddedSigs    []ValidatorSignature
	Proof         [8]*big.Int
	Commitments   [2]*big.Int
	CommitmentPok [2]*big.Int
}

// MisbehaviourProof contains the two independent batch proofs needed for
// standalone Tendermint equivocation: one proof for each conflicting header.
type MisbehaviourProof struct {
	Header1 HeaderBatchProof
	Header2 HeaderBatchProof
}

func (p *EcipProver) GenerateHeaderBatchProof(sigs []ValidatorSignature) (HeaderBatchProof, error) {
	bucket, paddedSigs, proof, commitments, commitmentPok, err := p.GenerateProof(sigs)
	if err != nil {
		return HeaderBatchProof{}, err
	}

	return HeaderBatchProof{
		Bucket:        bucket,
		PaddedSigs:    paddedSigs,
		Proof:         proof,
		Commitments:   commitments,
		CommitmentPok: commitmentPok,
	}, nil
}

func (p *EcipProver) GenerateMisbehaviourProof(
	header1Sigs []ValidatorSignature,
	header2Sigs []ValidatorSignature,
) (MisbehaviourProof, error) {
	if len(header1Sigs) == 0 {
		return MisbehaviourProof{}, fmt.Errorf("header1 proof: no signatures provided")
	}
	if len(header2Sigs) == 0 {
		return MisbehaviourProof{}, fmt.Errorf("header2 proof: no signatures provided")
	}

	header1Proof, err := p.GenerateHeaderBatchProof(header1Sigs)
	if err != nil {
		return MisbehaviourProof{}, fmt.Errorf("header1 proof: %w", err)
	}

	header2Proof, err := p.GenerateHeaderBatchProof(header2Sigs)
	if err != nil {
		return MisbehaviourProof{}, fmt.Errorf("header2 proof: %w", err)
	}

	return MisbehaviourProof{
		Header1: header1Proof,
		Header2: header2Proof,
	}, nil
}
