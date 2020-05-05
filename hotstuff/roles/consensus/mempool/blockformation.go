package mempool

import (
	"fmt"
	"math/rand"
	"time"
)

// HoC == Hash Of Collection

type BlockPayloadGenerator interface {
	Payload() []byte // should return an array of 1000 dummy HoC
}

// DummyBlockPayloadGenerator implements BlockPayloadGenerator interface
type DummyBlockPayloadGenerator struct {
	rng  *rand.Rand
	size uint
}

func (p *DummyBlockPayloadGenerator) SetSeed(seed int64) *DummyBlockPayloadGenerator {
	p.rng = rand.New(rand.NewSource(seed))
	return p // return reference to DummyBlockPayload for chaining
}

func (p DummyBlockPayloadGenerator) Payload() []byte {
	// return the dummyData of the input DummyBlockPayload p
	// the dummyData should be changed everytime Payload() is called
	arr := make([]byte, p.size)
	if _, err := p.rng.Read(arr); err != nil {
		// Handle err
		fmt.Println("Cannot generate random data...")
	}
	return arr
}

func NewDummyBlockPayloadGenerator(size uint) *DummyBlockPayloadGenerator {
	// factor method to create a DummyBlockPayload with random data
	// size: size of payload in bytes
	seed := time.Now().UnixNano() //random seed used for this DummyBlockPayload
	rng := rand.New(rand.NewSource(seed))
	return &DummyBlockPayloadGenerator{rng, size}
}
