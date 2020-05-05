package bft

import (
	"bytes"
	"log"
)

type AggregatedSignature struct {
	RawSignature []byte
	Signers      []byte
}

type QC struct {
	*AggregatedSignature
	Height    uint64
	BlockHash []byte
}

func AggregateSignature(sigs []*Signature) *AggregatedSignature {
	rawSigs := []byte{}
	signers := []byte{}

	return &AggregatedSignature{
		RawSignature: rawSigs,
		Signers:      signers,
	}
}

func (qc *QC) AggregateVotes(blockHash []byte, votes []*Vote) {
	var signatures = []*Signature{}
	for _, vote := range votes {
		if !bytes.Equal(vote.BlockHash, blockHash) {
			log.Println("votes used for QC are not for the same block")
		}
		signatures = append(signatures, vote.Signature)
	}

	qc.AggregatedSignature = AggregateSignature(signatures)
}

func NewQC(height uint64, blockHash []byte, votes []*Vote) (*QC, error) {
	qc := QC{
		Height:    height,
		BlockHash: blockHash,
	}

	qc.AggregateVotes(blockHash, votes)

	return &qc, nil
}

func NewGenesisQC(genesisBlockHash []byte) *QC {
	return &QC{
		Height:              0,
		BlockHash:           genesisBlockHash,
		AggregatedSignature: &AggregatedSignature{},
	}
}
