package bft

type ViewChange struct {
	*QC
	*Signature
}

type BlockProposal struct {
	*Block
	*Signature
}

type Vote struct {
	BlockHash []byte
	*Signature
}

type Signature struct {
	RawSignature []byte
	Signer       uint32
}
