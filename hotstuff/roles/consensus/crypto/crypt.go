package crypto

import "github.com/dapperlabs/flow-consensus-research/hotstuff/roles/consensus/bft"

func SignMsg(msg interface{}, signer uint) *bft.Signature {
	// var sig bytes.Buffer
	// enc := gob.NewEncoder(&sig)
	// err := enc.Encode(msg)
	// if err != nil {
	// 	log.Fatal("encode error:", err)
	// }
	// rawSig := sig.Bytes()
	rawSig := []byte{}
	return &bft.Signature{
		RawSignature: rawSig,
		Signer:       uint32(signer),
	}
}

func VerifySignature(rawData interface{}, sig bft.Signature) bool {
	return true
}
