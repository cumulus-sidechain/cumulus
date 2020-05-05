package bft

import (
	"bytes"
	"crypto/sha256"
)

type Block struct {
	Height    uint64
	Timestamp int64
	Hash      []byte
	PrevHash  []byte
	Payload   []byte
	*QC
}

// type ProposedBlock struct {
// 	ProposerID uint32
// 	ViewNumber uint32
// 	*Block
// 	proposerSignature *Signature
// }

func (b Block) GetPayload() []byte {
	return b.Payload
}

func (b Block) GetBlockHash() []byte {
	return b.Hash
}

func (b *Block) computeBlockHash() {
	data := bytes.Join(
		[][]byte{
			b.PrevHash[:],
			b.Payload,
		},
		[]byte{},
	)

	blockHash := sha256.Sum256(data)

	b.Hash = blockHash[:]
}

func (b Block) GetPreviousBlockHash() []byte {
	return b.PrevHash
}

func (b Block) GetHeight() uint64 {
	return b.Height
}

func (b Block) VerifyBlock(localPrevBlock *Block) bool {
	// Verify if block height is previous block height (maintained locally) + 1
	if b.Height != localPrevBlock.Height+1 {
		return false
	}
	return true
}

// new block should be built on top of current highest QC
func NewBlock(blockPayload []byte, highestQC *QC, timestamp int64) *Block {
	b := &Block{
		Height:    highestQC.Height + 1,
		PrevHash:  highestQC.BlockHash,
		Payload:   blockPayload,
		QC:        highestQC,
		Timestamp: timestamp,
	}

	b.computeBlockHash()

	return b
}

func NewGenesisBlock() *Block {
	b := &Block{
		Height:   0,
		PrevHash: []byte{},
		Payload:  nil,
		QC:       nil,
	}

	b.computeBlockHash()

	return b
}
