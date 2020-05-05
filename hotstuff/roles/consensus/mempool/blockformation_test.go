package mempool_test

import (
	"testing"

	"github.com/cumulus-sidechain/cumulus/hotstuff/roles/consensus/mempool"
	. "github.com/onsi/gomega"
)

func TestNewDummyBlockPayloadGenerator(t *testing.T) {
	// create a payload generator using random seed
	payloadGenerator := mempool.NewDummyBlockPayloadGenerator(10)
	// generate data
	data := payloadGenerator.Payload()
	// set seed of payloadGenerator to 3 and re-generate data
	seedData := payloadGenerator.SetSeed(int64(3)).Payload()
	t.Log(data)
	t.Log(seedData)
}
